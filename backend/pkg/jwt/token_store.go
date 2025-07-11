package jwt

import (
	"context"
	"time"
)

// RefreshTokenEntry represents a single entry in the refresh token store.
type refreshTokenEntry struct {
	TokenID   string    // Unique identifier for the refresh token
	FamilyID  string    // Family ID for grouping tokens
	UserID    uint      // User ID associated with the token
	ExpiresAt time.Time // Expiration time of the token
	Valid     bool      // Indicates if the token is valid
}

// TokenStore defines the interface for storing and managing JWT tokens.
type TokenStore interface {
	// Refresh Token
	StoreRefreshToken(ctx context.Context, tokenID string, familyID string, userID uint, expiry time.Duration) error
	IsRefreshTokenValid(ctx context.Context, tokenID string, familyID string) (bool, error)
	InvalidateRefreshToken(ctx context.Context, tokenID string) error
	InvalidateRefreshTokenFamily(ctx context.Context, familyID string) error
	// Access Token
	RevokeAccessToken(ctx context.Context, tokenID string, expiry time.Time) error
	IsAccessTokenRevoked(ctx context.Context, tokenID string) (bool, error)
}
