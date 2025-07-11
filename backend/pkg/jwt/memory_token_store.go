package jwt

import (
	"sync"
	"time"
)

// The MemoryTokenStore is an in-memory implementation of TokenStore.
type MemoryTokenStore struct {
	refreshTokens map[string]refreshTokenEntry // store for Refresh Tokens
	revokedTokens map[string]time.Time         // store for revoked Access Tokens
	mu            sync.RWMutex
}

func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{
		refreshTokens: make(map[string]refreshTokenEntry),
		revokedTokens: make(map[string]time.Time),
	}
}

func (s *MemoryTokenStore) StoreRefreshToken(tokenID string, familyID string, userID uint, expiry time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.refreshTokens[tokenID] = refreshTokenEntry{
		TokenID:   tokenID,
		FamilyID:  familyID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiry),
		Valid:     true,
	}
	return nil
}

func (s *MemoryTokenStore) IsRefreshTokenValid(tokenID string, familyID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查 Token 是否存在
	entry, exists := s.refreshTokens[tokenID]
	if !exists {
		return false, nil
	}

	// 检查 Token 是否有效且未过期
	if !entry.Valid || time.Now().After(entry.ExpiresAt) {
		return false, nil
	}

	// 检查 FamilyID 是否匹配
	if familyID != "" && entry.FamilyID != familyID {
		return false, nil
	}

	return true, nil
}

func (s *MemoryTokenStore) InvalidateRefreshToken(tokenID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, exists := s.refreshTokens[tokenID]; exists {
		// 使单个 Token 无效
		entry.Valid = false
		s.refreshTokens[tokenID] = entry
	}
	return nil
}

func (s *MemoryTokenStore) InvalidateRefreshTokenFamily(familyID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for tokenID, entry := range s.refreshTokens {
		// 使所有属于该 FamilyID 的 Token 无效
		if entry.FamilyID == familyID {
			entry.Valid = false
			s.refreshTokens[tokenID] = entry
		}
	}
	return nil
}

func (s *MemoryTokenStore) RevokeAccessToken(tokenID string, expiry time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 撤销时加入黑名单
	s.revokedTokens[tokenID] = expiry
	return nil
}

func (s *MemoryTokenStore) IsAccessTokenRevoked(tokenID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查 Token 是否在黑名单中
	expiry, exists := s.revokedTokens[tokenID]
	if !exists {
		return false, nil
	}

	// 删除已过期的黑名单记录
	if time.Now().After(expiry) {
		s.mu.RUnlock()
		s.mu.Lock()
		delete(s.revokedTokens, tokenID)
		s.mu.Unlock()
		s.mu.RLock()
		return false, nil
	}

	return true, nil
}

// 清理过期的 Token 记录
func (s *MemoryTokenStore) CleanupExpiredTokens() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// 清理过期的 Refresh Token
	for tokenID, entry := range s.refreshTokens {
		if now.After(entry.ExpiresAt) {
			delete(s.refreshTokens, tokenID)
		}
	}

	// 清理过期的 Access Token 黑名单
	for tokenID, expiry := range s.revokedTokens {
		if now.After(expiry) {
			delete(s.revokedTokens, tokenID)
		}
	}
}
