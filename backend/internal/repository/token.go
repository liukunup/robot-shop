package repository

import (
	v1 "backend/api/v1"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	keyPrefix          = "token"
	nsAccess           = "access"
	nsRefresh          = "refresh"
	nsFamily           = "family"
	nsUserFamilies     = "user_families"
	maxPendingSyncSize = 1000000
)

type refreshTokenEntry struct {
	TokenID   string
	FamilyID  string    // Family ID for grouping tokens
	UserID    uint      // User ID associated with the token
	ExpiresAt time.Time // Expiration time of the token
	Valid     bool      // Indicates if the token is valid
}

type TokenStore interface {
	// RefreshToken managed by whitelist
	StoreRefreshToken(ctx context.Context, tokenID string, familyID string, userID uint, expiry time.Duration) error
	IsRefreshTokenValid(ctx context.Context, tokenID string, familyID string) (bool, error)
	InvalidateRefreshToken(ctx context.Context, tokenID string) error
	InvalidateRefreshTokenByFamilyID(ctx context.Context, familyID string) error
	InvalidateRefreshTokenByUserID(ctx context.Context, userID uint) error

	// AccessToken managed by blacklist
	RevokeAccessToken(ctx context.Context, tokenID string, expiry time.Duration) error
	IsAccessTokenRevoked(ctx context.Context, tokenID string) (bool, error)
}

func NewTokenStore(
	repository *Repository,
) TokenStore {
	s := &tokenStore{
		Repository:  repository,
		pendingSync: make(map[string]struct{}),
		stopChan:    make(chan struct{}),
	}

	// Start health check ticker
	s.healthTicker = time.NewTicker(10 * time.Second)
	go s.healthCheck()

	return s
}

type tokenStore struct {
	*Repository
	pendingSync  map[string]struct{} // 记录待同步数据
	pendingMu    sync.Mutex          // 数据同步锁
	mu           sync.RWMutex        // Redis连接状态锁
	isRedisDown  bool                // Redis连接状态
	healthTicker *time.Ticker
	stopChan     chan struct{}
}

// Generate a unique key for each token based on namespace and ID
func (s *tokenStore) key(namespace, id string) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefix, namespace, id)
}

// 写操作：双写 Ristretto + Redis
func (s *tokenStore) StoreRefreshToken(ctx context.Context, tokenID string, familyID string, userID uint, expiry time.Duration) error {
	data := refreshTokenEntry{
		TokenID:   tokenID,
		FamilyID:  familyID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiry),
		Valid:     true,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal refresh token data failed: %w", err)
	}

	key := s.key(nsRefresh, tokenID)

	// 1. 写入 Ristretto
	cached := s.cache.SetWithTTL(key, jsonData, int64(len(jsonData)), expiry)

	s.mu.RLock()
	redisDown := s.isRedisDown
	s.mu.RUnlock()

	if redisDown {
		s.pendingMu.Lock()
		if len(s.pendingSync) >= maxPendingSyncSize {
			return fmt.Errorf("pending sync size exceeded %d, data will not be stored", maxPendingSyncSize)
		}
		s.pendingSync[key] = struct{}{} // 标记为待同步
		s.pendingMu.Unlock()

		if !cached {
			return fmt.Errorf("failed to store refresh token")
		}
		return nil
	}

	// 2. 写入 Redis
	// 使用 Pipelined 确保原子性和性能
	_, err = s.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SetEx(ctx, key, jsonData, expiry)

		// 维护 FamilyID 到 TokenID 的映射（用于快速查找）
		if familyID != "" {
			familyKey := s.key(nsFamily, familyID)
			pipe.SAdd(ctx, familyKey, tokenID)
			pipe.Expire(ctx, familyKey, expiry+15*time.Minute) // 比 Token 多保留15分钟
		}

		// 维护 UserID 到 FamilyID 的映射（用于快速查找）
		userFamiliesKey := s.key(nsUserFamilies, fmt.Sprintf("%d", userID))
		pipe.SAdd(ctx, userFamiliesKey, familyID)
		pipe.Expire(ctx, userFamiliesKey, expiry+24*time.Hour) // 比 Token 多保留24小时
		return nil
	})

	if err != nil {
		s.mu.Lock()
		s.isRedisDown = true
		s.mu.Unlock()

		if !cached {
			return fmt.Errorf("failed to store refresh token")
		}
	}
	return nil
}

// 读操作：优先检查 Ristretto，再检查 Redis 并回填
func (s *tokenStore) IsRefreshTokenValid(ctx context.Context, tokenID string, familyID string) (bool, error) {
	key := s.key(nsRefresh, tokenID)

	// 1. 检查 Ristretto
	if data, found := s.cache.Get(key); found {
		var token refreshTokenEntry
		if err := json.Unmarshal(data.([]byte), &token); err != nil {
			return false, fmt.Errorf("unmarshal token data from cache failed: %w", err)
		}
		// 检查有效性、FamilyID匹配和过期时间
		return token.Valid &&
			(familyID == "" || token.FamilyID == familyID) &&
			time.Now().Before(token.ExpiresAt), nil
	}

	s.mu.RLock()
	redisDown := s.isRedisDown
	s.mu.RUnlock()

	if redisDown {
		return false, v1.ErrRedisUnavailable
	}

	// 2. 检查 Redis
	data, err := s.rdb.Get(ctx, key).Bytes()
	switch {
	case err == redis.Nil:
		return false, nil // 不存在视为无效

	case err != nil:
		s.mu.Lock()
		s.isRedisDown = true
		s.mu.Unlock()

		return false, v1.ErrRedisUnavailable
	}

	var token refreshTokenEntry
	if err := json.Unmarshal(data, &token); err != nil {
		return false, fmt.Errorf("unmarshal token data from redis failed: %w", err)
	}

	// 回填 Ristretto
	s.cache.SetWithTTL(key, data, int64(len(data)), time.Until(token.ExpiresAt))

	// 检查有效性、FamilyID匹配和过期时间
	return token.Valid &&
		(familyID == "" || token.FamilyID == familyID) &&
		time.Now().Before(token.ExpiresAt), nil
}

func (s *tokenStore) InvalidateRefreshToken(ctx context.Context, tokenID string) error {
	key := s.key(nsRefresh, tokenID)

	// 1. 从 Ristretto 中删除即可
	s.cache.Del(key)

	s.mu.RLock()
	redisDown := s.isRedisDown
	s.mu.RUnlock()

	if redisDown {
		return v1.ErrRedisUnavailable
	}

	// 2. 将 Redis 更新为无效状态
	// 使用Lua脚本保证原子性
	script := `
	local key = KEYS[1]
	local data = redis.call('GET', key)
	if not data then return 0 end

	local token = cjson.decode(data)
	token['valid'] = false
	redis.call('SET', key, cjson.encode(token))
	return 1
	`

	_, err := s.rdb.Eval(ctx, script, []string{key}).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("lua script failed: %w", err)
	}
	return nil
}

func (s *tokenStore) InvalidateRefreshTokenByFamilyID(ctx context.Context, familyID string) error {
	// 检查 Redis 状态
	s.mu.RLock()
	redisDown := s.isRedisDown
	s.mu.RUnlock()

	if redisDown {
		return v1.ErrRedisUnavailable
	}

	familyKey := s.key(nsFamily, familyID)

	// 分页扫描避免阻塞
	var cursor uint64
	for {
		var tokenIDs []string
		var err error
		tokenIDs, cursor, err = s.rdb.SScan(ctx, familyKey, cursor, "", 100).Result()
		if err != nil {
			return fmt.Errorf("scan family failed: %w", err)
		}

		// 批量失效
		for _, tokenID := range tokenIDs {
			if err := s.InvalidateRefreshToken(ctx, tokenID); err != nil {
				if errors.Is(err, v1.ErrRedisUnavailable) {
					// Redis 不可用时不报错（降级处理）
					s.logger.Warn("redis unavailable, skip token invalidation: %v", zap.Error(err))
					continue
				}
				return fmt.Errorf("invalidate token %s failed: %w", tokenID, err)
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}

func (s *tokenStore) InvalidateRefreshTokenByUserID(ctx context.Context, userID uint) error {
	// 检查 Redis 状态
	s.mu.RLock()
	redisDown := s.isRedisDown
	s.mu.RUnlock()

	if redisDown {
		return v1.ErrRedisUnavailable
	}

	userFamiliesKey := s.key(nsUserFamilies, fmt.Sprintf("%d", userID))

	// 分页扫描避免阻塞
	var cursor uint64
	for {
		var familyIDs []string
		var err error
		familyIDs, cursor, err = s.rdb.SScan(ctx, userFamiliesKey, cursor, "", 100).Result()
		if err != nil {
			return fmt.Errorf("scan user families failed: %w", err)
		}

		// 批量失效
		for _, familyID := range familyIDs {
			if err := s.InvalidateRefreshTokenByFamilyID(ctx, familyID); err != nil {
				if errors.Is(err, v1.ErrRedisUnavailable) {
					// Redis 不可用时不报错（降级处理）
					s.logger.Warn("redis unavailable, skip family invalidation: %v", zap.Error(err))
					continue
				}
				return fmt.Errorf("invalidate family %s failed: %w", familyID, err)
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}

// 写操作：双写 Ristretto + Redis
func (s *tokenStore) RevokeAccessToken(ctx context.Context, tokenID string, expiry time.Duration) error {
	// KV
	key := s.key(nsAccess, tokenID)
	val := "1" // 标记为已撤销

	// 1. 写入 Ristretto
	cached := s.cache.SetWithTTL(key, val, int64(len(val)), expiry)

	s.mu.RLock()
	redisDown := s.isRedisDown
	s.mu.RUnlock()

	if redisDown {
		s.pendingMu.Lock()
		if len(s.pendingSync) >= maxPendingSyncSize {
			return fmt.Errorf("pending sync size exceeded %d, data will not be stored", maxPendingSyncSize)
		}
		s.pendingSync[key] = struct{}{} // 标记为待同步
		s.pendingMu.Unlock()

		if !cached {
			return fmt.Errorf("failed to store access token")
		}
		return nil
	}

	// 2. 写入 Redis
	// 使用 SetNX 确保原子性，避免重复撤销
	ok, err := s.rdb.SetNX(ctx, key, val, expiry).Result()
	if err != nil {
		s.mu.Lock()
		s.isRedisDown = true
		s.mu.Unlock()

		if !cached {
			return fmt.Errorf("failed to store access token")
		}
		return nil
	}
	if !ok {
		return v1.ErrTokenAlreadyRevoked
	}
	return nil
}

// 读操作：优先检查 Ristretto 缓存，再检查 Redis 并回填数据
func (s *tokenStore) IsAccessTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	// KV
	key := s.key(nsAccess, tokenID)
	val := "1" // 标记为已撤销

	// 1. 优先检查 Ristretto
	_, found := s.cache.Get(key)
	if found {
		return true, nil // 如果在缓存中找到，说明已撤销
	}

	s.mu.RLock()
	redisDown := s.isRedisDown
	s.mu.RUnlock()

	// Redis 服务不可用 且 缓存也没找到
	if redisDown && !found {
		return false, v1.ErrRedisUnavailable
	}

	// 2. 继续检查 Redis
	exists, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		s.mu.Lock()
		s.isRedisDown = true
		s.mu.Unlock()

		// Redis 服务不可用 且 缓存也没找到
		if !found {
			return false, v1.ErrRedisUnavailable
		}
	}

	// 回填 Ristretto 缓存
	expiry, err := s.rdb.TTL(ctx, key).Result()
	if err != nil {
		s.mu.Lock()
		s.isRedisDown = true
		s.mu.Unlock()

		// 仅仅是查询 TTL 报错，并不影响业务
		return exists > 0, err
	}
	s.cache.SetWithTTL(key, val, int64(len(val)), expiry) // 保持与 Redis 一致的过期时间

	return exists > 0, nil
}

func (s *tokenStore) syncPendingToRedis(ctx context.Context) error {
	s.pendingMu.Lock()
	defer s.pendingMu.Unlock()

	if len(s.pendingSync) == 0 {
		return nil
	}

	// 分批处理避免大事务
	batchSize := 100
	keys := make([]string, 0, len(s.pendingSync))
	for key := range s.pendingSync {
		keys = append(keys, key)
	}

	for i := 0; i < len(keys); i += batchSize {
		end := i + batchSize
		if end > len(keys) {
			end = len(keys)
		}
		batch := keys[i:end]

		// 使用pipeline批量处理
		_, err := s.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			for _, key := range batch {
				// 从 Ristretto 中获取数据和过期时间
				data, found := s.cache.Get(key)
				if !found {
					continue
				}
				expiry, found := s.cache.GetTTL(key)
				if !found {
					continue
				}

				// 根据 key 的前缀判断类型
				switch {
				case strings.HasPrefix(key, s.key(nsRefresh, "")):
					var token refreshTokenEntry
					if err := json.Unmarshal(data.([]byte), &token); err != nil {
						continue
					}

					pipe.SetEx(ctx, key, data, expiry) // 保持与 Ristretto 一致的过期时间

					// 维护 FamilyID 到 TokenID 的映射（用于快速查找）
					if token.FamilyID != "" {
						familyKey := s.key(nsFamily, token.FamilyID)
						pipe.SAdd(ctx, familyKey, token.TokenID)
						pipe.Expire(ctx, familyKey, expiry+15*time.Minute) // 比 Token 多保留15分钟
					}

					// 维护 UserID 到 FamilyID 的映射（用于快速查找）
					userFamiliesKey := s.key(nsUserFamilies, fmt.Sprintf("%d", token.UserID))
					pipe.SAdd(ctx, userFamiliesKey, token.FamilyID)
					pipe.Expire(ctx, userFamiliesKey, expiry+24*time.Hour) // 比 Token 多保留24小时

				case strings.HasPrefix(key, s.key(nsAccess, "")):
					pipe.SetEx(ctx, key, data, expiry) // 保持与 Ristretto 一致的过期时间
				}
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("batch sync failed: %w", err)
		}

		// 清理已同步的数据
		for _, key := range batch {
			delete(s.pendingSync, key)
		}
	}
	return nil
}

func (s *tokenStore) healthCheck() {
	for {
		select {
		case <-s.healthTicker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			_, err := s.rdb.Ping(ctx).Result()
			cancel()

			s.mu.Lock()
			if err != nil {
				if !s.isRedisDown {
					s.logger.WithContext(ctx).Warn("redis disconnected, entering fallback mode", zap.Error(err))
				}
				s.isRedisDown = true
			} else {
				if s.isRedisDown {
					s.logger.WithContext(ctx).Info("redis reconnected, resuming normal operations")
					s.isRedisDown = false

					// 执行增量同步
					if err := s.syncPendingToRedis(ctx); err != nil {
						s.logger.WithContext(ctx).Error("failed to sync pending data to redis", zap.Error(err))
					}
				}
			}
			s.mu.Unlock()
		case <-s.stopChan:
			return
		}
	}
}

func (s *tokenStore) Close() {
	if s.healthTicker != nil {
		s.healthTicker.Stop()
	}
	close(s.stopChan)
}
