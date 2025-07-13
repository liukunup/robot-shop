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
	NsRefresh      = "refresh"
	NsFamily       = "family"
	NsUserFamilies = "user_families"
	NsAccess       = "access"
)

type RedisTokenStore struct {
	client    redis.UniversalClient
	keyPrefix string
}

func NewRedisTokenStore(client redis.UniversalClient, keyPrefix string) TokenStore {
	return &RedisTokenStore{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// 生成带 Namespace 的 Redis Key
func (s *RedisTokenStore) key(namespace, id string) string {
	return fmt.Sprintf("%s:%s:%s", s.keyPrefix, namespace, id)
}

func (s *RedisTokenStore) StoreRefreshToken(ctx context.Context, tokenID string, familyID string, userID uint, expiry time.Duration) error {
	data := refreshTokenEntry{
		TokenID:   tokenID,
		FamilyID:  familyID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiry),
		Valid:     true,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal token data failed: %w", err)
	}

	_, err = s.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		// 存储 Refresh Token
		pipe.SetEx(ctx, s.key(NsRefresh, tokenID), jsonData, expiry)
		// 存储 Family
		if familyID != "" {
			familyKey := s.key(NsFamily, familyID)
			pipe.SAdd(ctx, familyKey, tokenID)
			pipe.Expire(ctx, familyKey, expiry+15*time.Minute) // 比 Token 多保留一段时间
		}
		// 维护 UserID 到 FamilyID 的映射（用于快速查找）
		userFamiliesKey := s.key(NsUserFamilies, fmt.Sprintf("%d", userID))
		pipe.SAdd(ctx, userFamiliesKey, familyID)
		pipe.Expire(ctx, userFamiliesKey, expiry+24*time.Hour) // 比 Token 多保留一段时间
		return nil
	})

	if err != nil {
		return fmt.Errorf("%w: %v", v1.ErrRedisUnavailable, err)
	}
	return nil
}

func (s *RedisTokenStore) IsRefreshTokenValid(ctx context.Context, tokenID string, familyID string) (bool, error) {
	data, err := s.client.Get(ctx, s.key(NsRefresh, tokenID)).Bytes()
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

	// 检查有效性、FamilyID匹配和过期时间
	return token.Valid &&
		(familyID == "" || token.FamilyID == familyID) &&
		time.Now().Before(token.ExpiresAt), nil
}

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

	_, err := s.client.Eval(ctx, script, []string{s.key(NsRefresh, tokenID)}).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("lua script failed: %w", err)
	}
	return nil
}

func (s *RedisTokenStore) InvalidateRefreshTokenFamily(ctx context.Context, familyID string) error {
	familyKey := s.key(NsFamily, familyID)

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

func (s *RedisTokenStore) InvalidateRefreshTokenFamilyByUserID(ctx context.Context, userID uint) error {
	userFamiliesKey := s.key(NsUserFamilies, fmt.Sprintf("%d", userID))

	// 分页扫描避免阻塞
	var cursor uint64
	for {
		var familyIDs []string
		var err error
		familyIDs, cursor, err = s.client.SScan(ctx, userFamiliesKey, cursor, "", 100).Result()
		if err != nil {
			return fmt.Errorf("scan user families failed: %w", err)
		}

		// 批量失效
		for _, familyID := range familyIDs {
			if err := s.InvalidateRefreshTokenFamily(ctx, familyID); err != nil {
				return fmt.Errorf("invalidate family %s failed: %w", familyID, err)
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (s *RedisTokenStore) RevokeAccessToken(ctx context.Context, tokenID string, expiry time.Duration) error {
	// 使用SET NX模式避免覆盖更晚的撤销记录
	ok, err := s.client.SetNX(ctx, s.key(NsAccess, tokenID), "1", expiry).Result()
	if err != nil {
		return fmt.Errorf("%w: %v", v1.ErrRedisUnavailable, err)
	}
	if !ok {
		return errors.New("token already revoked with later expiry")
	}
	return nil
}

func (s *RedisTokenStore) IsAccessTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	exists, err := s.client.Exists(ctx, s.key(NsAccess, tokenID)).Result()
	if err != nil {
		return false, fmt.Errorf("%w: %v", v1.ErrRedisUnavailable, err)
	}
	return exists > 0, nil
}
