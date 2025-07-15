package repository

import (
	"backend/pkg/log"
	"backend/pkg/zapgorm2"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ctxTxKey = "TxKey"

type Repository struct {
	db *gorm.DB
	e  *casbin.SyncedEnforcer
	//rdb    *redis.UniversalClient
	logger *log.Logger
}

func NewRepository(
	logger *log.Logger,
	db *gorm.DB,
	e *casbin.SyncedEnforcer,
	// rdb *redis.UniversalClient,
) *Repository {
	return &Repository{
		db: db,
		e:  e,
		//rdb:    rdb,
		logger: logger,
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	logger := zapgorm2.New(l.Logger)
	driver := conf.GetString("data.db.user.driver")
	dsn := conf.GetString("data.db.user.dsn")

	// GORM doc: https://gorm.io/docs/connecting_to_the_database.html
	switch driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger,
		})
	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		panic("unknown db driver")
	}
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	// Connection Pool config
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

func NewCasbinEnforcer(conf *viper.Viper, l *log.Logger, db *gorm.DB) *casbin.SyncedEnforcer {
	a, _ := gormadapter.NewAdapterByDB(db)
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)

	if err != nil {
		panic(err)
	}
	e, _ := casbin.NewSyncedEnforcer(m, a)
	e.StartAutoLoadPolicy(10 * time.Second) // 每10秒自动加载策略，防止启动多服务进程策略不一致

	// Enable Logger, decide whether to show it in terminal
	//e.EnableLog(true)

	// Save the policy back to DB.
	e.EnableAutoSave(true)

	return e
}

// 多级缓存：内存缓存 -> Redis
type MultiLevelCache struct {
	redisClient  redis.UniversalClient                 // Redis 客户端
	memoryCache  *ristretto.Cache[string, interface{}] // Ristretto 内存缓存
	mu           sync.RWMutex                          // 用于保护内存缓存的并发访问
	isRedisDown  bool                                  // Redis 故障标志
	pendingSync  *sync.Map                             // 用于故障期间暂存待同步数据
	healthTicker *time.Ticker                          // 用于健康检查的定时器
	stopChan     chan struct{}                         // 用于停止健康检查的通道
	logger       *log.Logger
}

// 创建多级缓存实例
func NewMultiLevelCache(conf *viper.Viper, logger *log.Logger) *MultiLevelCache {
	// 初始化 Redis 客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.GetStringSlice("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	isRedisDown := false
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		isRedisDown = true
		logger.Warn("redis connection failed, entering fallback mode", zap.Error(err))
	}

	// 初始化 Ristretto 内存缓存
	memCache, err := ristretto.NewCache(&ristretto.Config[string, interface{}]{
		NumCounters: 1e7,     // 键的跟踪数量
		MaxCost:     1 << 30, // 最大缓存容量(1GB)
		BufferItems: 64,      // 性能优化参数
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create Ristretto cache: %v", err))
	}

	mlc := &MultiLevelCache{
		redisClient: rdb,
		memoryCache: memCache,
		isRedisDown: isRedisDown,
		pendingSync: &sync.Map{},
		stopChan:    make(chan struct{}),
		logger:      logger,
	}

	// 启动健康检查(每10秒检查一次 Redis 连接)
	mlc.healthTicker = time.NewTicker(10 * time.Second)
	go mlc.healthCheck()

	return mlc
}

// 读操作：内存缓存优先 + Redis 回填数据
func (m *MultiLevelCache) Get(ctx context.Context, key string) ([]byte, error) {
	// 1. 先尝试从内存缓存读取
	if val, found := m.memoryCache.Get(key); found {
		return val.([]byte), nil
	}

	// 2. 如果内存中没有，尝试从 Redis 读取
	m.mu.RLock()
	redisAvailable := !m.isRedisDown
	m.mu.RUnlock()

	if redisAvailable {
		val, err := m.redisClient.Get(ctx, key).Bytes()
		if err == nil {
			// 将 Redis 数据回填到内存缓存
			m.memoryCache.Set(key, val, int64(len(val)))
			return val, nil
		}
		if err != redis.Nil {
			m.logger.WithContext(ctx).Error("redis Get error", zap.String("key", key), zap.Error(err))
		}
	}

	return nil, redis.Nil
}

// 写操作：双写 内存缓存 + Redis
func (m *MultiLevelCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	// 1. 总是先写入内存缓存
	m.memoryCache.SetWithTTL(key, value, int64(len(value)), ttl)

	m.mu.RLock()
	redisDown := m.isRedisDown
	m.mu.RUnlock()

	// 2. 写入 Redis 或记录待同步
	if redisDown {
		m.pendingSync.Store(key, value)
		return nil
	}

	if err := m.redisClient.Set(ctx, key, value, ttl).Err(); err != nil {
		m.mu.Lock()
		m.isRedisDown = true
		m.pendingSync.Store(key, value)
		m.mu.Unlock()
		return err
	}

	return nil
}

// 删除操作：双删 内存缓存 + Redis
func (m *MultiLevelCache) Delete(ctx context.Context, key string) error {
	// 1. 从内存缓存中删除
	m.memoryCache.Del(key)

	m.mu.RLock()
	redisDown := m.isRedisDown
	m.mu.RUnlock()

	// 2. 如果 Redis 可用，直接删除；否则记录待同步
	if redisDown {
		m.pendingSync.Delete(key)
		return nil
	}

	if err := m.redisClient.Del(ctx, key).Err(); err != nil {
		m.mu.Lock()
		m.isRedisDown = true
		m.pendingSync.Delete(key)
		m.mu.Unlock()
		return err
	}

	return nil
}

// Redis 定时健康检查
func (m *MultiLevelCache) healthCheck() {
	for {
		select {
		case <-m.healthTicker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err := m.redisClient.Ping(ctx).Result()
			cancel()

			m.mu.Lock()
			if err != nil {
				if !m.isRedisDown {
					m.logger.Warn("Redis disconnected, entering fallback mode", zap.Error(err))
				}
				m.isRedisDown = true
			} else {
				if m.isRedisDown {
					m.logger.Info("Redis reconnected, starting data sync")
					m.isRedisDown = false
					go m.syncPendingData()
				}
			}
			m.mu.Unlock()
		case <-m.stopChan:
			return
		}
	}
}

// syncPendingData 将故障期间的数据同步到Redis
func (m *MultiLevelCache) syncPendingData() {
	// 分批同步，避免阻塞
	const batchSize = 100
	ctx := context.Background()

	m.pendingSync.Range(func(key, value interface{}) bool {
		keys := make([]string, 0, batchSize)
		values := make([][]byte, 0, batchSize)

		// 收集一批数据
		m.pendingSync.Range(func(k, v interface{}) bool {
			keys = append(keys, k.(string))
			values = append(values, v.([]byte))
			return len(keys) < batchSize
		})

		if len(keys) == 0 {
			return false
		}

		// 使用Pipeline批量写入
		pipe := m.redisClient.Pipeline()
		for i := 0; i < len(keys); i++ {
			pipe.Set(ctx, keys[i], values[i], 0)
		}

		if _, err := pipe.Exec(ctx); err != nil {
			m.logger.WithContext(ctx).Error("Failed to sync batch to Redis", zap.Error(err))
			return false
		}

		// 删除已同步的key
		for _, key := range keys {
			m.pendingSync.Delete(key)
		}

		m.logger.WithContext(ctx).Info("Successfully synced to Redis", zap.Int("count", len(keys)))
		return true
	})
}

// 关闭资源
func (m *MultiLevelCache) Close() error {
	m.healthTicker.Stop()
	close(m.stopChan)
	m.memoryCache.Close()
	m.redisClient.Close()
	return nil
}
