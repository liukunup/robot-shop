package jwt

import (
	v1 "backend/api/v1"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	NsRefresh = "refresh"
	NsFamily  = "family"
	NsUserFamilies = "user_families"
)

type RedisTokenStore struct {
	client        redis.UniversalClient
	keyPrefix     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewRedisTokenStore(client redis.UniversalClient, keyPrefix string) *RedisTokenStore {
	return &RedisTokenStore{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// 生成带命名空间的Redis键
func (s *RedisTokenStore) key(namespace, id string) string {
	return fmt.Sprintf("%s:%s:%s", s.keyPrefix, namespace, id)
}

// StoreRefreshToken 存储Refresh Token (带自动过期和家族管理)
func (s *RedisTokenStore) StoreRefreshToken(ctx context.Context, tokenID, familyID string, userID uint, expiry time.Duration) error {
	data := refreshTokenEntry{
		TokenID:   tokenID,
		FamilyID:  familyID,
		UserID:    userID,
		Valid:     true,
		ExpiresAt: time.Now().Add(expiry),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal token data failed: %w", err)
	}

	// 使用Pipeline批量操作
	_, err = s.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		// 存储Token数据
		pipe.SetEx(ctx, s.key(NsRefresh, tokenID), jsonData, expiry)

		// 添加到家族集合
		if familyID != "" {
			familyKey := s.key("family", familyID)
			pipe.SAdd(ctx, familyKey, tokenID)
			pipe.Expire(ctx, familyKey, expiry+10*time.Minute) // 比Token多保留一段时间
		}

		// 维护用户到家族的映射（用于快速查找）
		userFamiliesKey := s.key("user_families", fmt.Sprintf("%d", userID))
		pipe.SAdd(ctx, userFamiliesKey, familyID)
		pipe.Expire(ctx, userFamiliesKey, expiry+24*time.Hour)
		return nil
	})

	if err != nil {
		return fmt.Errorf("%w: %v", v1.ErrRedisUnavailable, err)
	}
	return nil
}

// IsRefreshTokenValid 检查Refresh Token有效性
func (s *RedisTokenStore) IsRefreshTokenValid(ctx context.Context, tokenID, familyID string) (bool, error) {
	data, err := s.client.Get(ctx, s.key("refresh", tokenID)).Bytes()
	switch {
	case err == redis.Nil:
		return false, nil // 不存在视为无效
	case err != nil:
		return false, fmt.Errorf("%w: %v", v1.ErrRedisUnavailable, err)
	}

	var token refreshTokenEntry
	if err := json.Unmarshal(data, &token); err != nil {
		return false, fmt.Errorf("unmarshal token data failed: %w", err)
	}

	// 检查有效性、家族匹配和过期时间
	now := time.Now().UnixNano()
	return token.Valid && 
		(familyID == "" || token.FamilyID == familyID) && 
		now < token.ExpiresAt, nil
}

// InvalidateRefreshToken 使单个Token失效
func (s *RedisTokenStore) InvalidateRefreshToken(ctx context.Context, tokenID string) error {
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

	_, err := s.client.Eval(ctx, script, []string{s.key("refresh", tokenID)}).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("lua script failed: %w", err)
	}
	return nil
}

// InvalidateRefreshTokenFamily 使整个家族失效
func (s *RedisTokenStore) InvalidateRefreshTokenFamily(ctx context.Context, familyID string) error {
	familyKey := s.key("family", familyID)

	// 分页扫描避免阻塞
	var cursor uint64
	for {
		var tokenIDs []string
		var err error
		tokenIDs, cursor, err = s.client.SScan(ctx, familyKey, cursor, "", 100).Result()
		if err != nil {
			return fmt.Errorf("scan family failed: %w", err)
		}

		// 批量失效
		for _, tokenID := range tokenIDs {
			if err := s.InvalidateRefreshToken(ctx, tokenID); err != nil {
				return fmt.Errorf("invalidate token %s failed: %w", tokenID, err)
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}

// RevokeAccessToken 撤销Access Token
func (s *RedisTokenStore) RevokeAccessToken(ctx context.Context, tokenID string, expiry time.Time) error {
	// 使用SET NX模式避免覆盖更晚的撤销记录
	ok, err := s.client.SetNX(ctx, s.key("revoked_access", tokenID), "1", time.Until(expiry)).Result()
	if err != nil {
		return fmt.Errorf("%w: %v", v1.ErrRedisUnavailable, err)
	}
	if !ok {
		return errors.New("token already revoked with later expiry")
	}
	return nil
}

// IsAccessTokenRevoked 检查Access Token是否被撤销
func (s *RedisTokenStore) IsAccessTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	exists, err := s.client.Exists(ctx, s.key("revoked_access", tokenID)).Result()
	if err != nil {
		return false, fmt.Errorf("%w: %v", v1.ErrRedisUnavailable, err)
	}
	return exists > 0, nil
}

// GetUserFamilies 获取用户的所有Token家族（用于全设备退出）
func (s *RedisTokenStore) GetUserFamilies(ctx context.Context, userID uint) ([]string, error) {
	key := s.key("user_families", fmt.Sprintf("%d", userID))
	return s.client.SMembers(ctx, key).Result()
}

// CleanStaleData 清理过期数据（后台任务）
func (s *RedisTokenStore) CleanStaleData(ctx context.Context) error {
	// Redis会自动清理过期的Key，此方法用于处理额外关系
	// 实现略（可使用SCAN遍历清理无效的家族集合等）
	return nil
}
