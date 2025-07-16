package repository

import (
	"backend/pkg/log"
	"backend/pkg/zapgorm2"
	"context"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/glebarez/sqlite"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ctxTxKey = "TxKey"

type Repository struct {
	db     *gorm.DB
	e      *casbin.SyncedEnforcer
	cache  *ristretto.Cache[string, interface{}]
	rdb    redis.UniversalClient
	m      *MinIO
	logger *log.Logger
}

func NewRepository(
	db *gorm.DB,
	e *casbin.SyncedEnforcer,
	cache *ristretto.Cache[string, interface{}],
	rdb redis.UniversalClient,
	m *MinIO,
	logger *log.Logger,
) *Repository {
	return &Repository{
		db:     db,
		e:      e,
		cache:  cache,
		rdb:    rdb,
		m:      m,
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

func NewCache() *ristretto.Cache[string, interface{}] {
	cache, err := ristretto.NewCache(&ristretto.Config[string, interface{}]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(fmt.Errorf("failed to create Ristretto cache: %w", err))
	}

	return cache
}

func NewRedis(conf *viper.Viper, log *log.Logger) redis.UniversalClient {
	// Use UniversalClient to support both single and cluster mode
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.GetStringSlice("data.redis.addrs"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		_ = rdb.Close() // close the client if ping fails
		log.WithContext(ctx).Warn("failed to connect to Redis", zap.Error(err))
	}

	return rdb
}

type minioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	Secure    bool
}

type MinIO struct {
	client *minio.Client
	bucket string
}

func NewMinIO(conf *viper.Viper, log *log.Logger) *MinIO {
	cfg := &minioConfig{
		Endpoint:  conf.GetString("storage.minio.endpoint"),
		AccessKey: conf.GetString("storage.minio.access_key"),
		SecretKey: conf.GetString("storage.minio.secret_key"),
		Bucket:    conf.GetString("storage.minio.bucket"),
		Region:    conf.GetString("storage.minio.region"),
		Secure:    conf.GetBool("storage.minio.secure"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Region: cfg.Region,
		Secure: cfg.Secure,
	})
	if err != nil {
		log.WithContext(ctx).Warn("failed to initialize MinIO client", zap.Error(err))
		return nil
	}

	// 检查桶是否存在（不存在则创建）
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		log.WithContext(ctx).Warn("failed to check bucket existence", zap.Error(err))
		return nil
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			log.WithContext(ctx).Warn("failed to create bucket", zap.Error(err))
			return nil
		}
	}

	return &MinIO{
		client: client,
		bucket: cfg.Bucket,
	}
}
