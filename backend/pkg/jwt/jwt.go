package jwt

import (
	v1 "backend/api/v1"
	"backend/internal/repository"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const (
	bearerPrefix       = "Bearer " // 请勿删除空格
	accessTokenExpiry  = 15 * time.Minute
	refreshTokenExpiry = 7 * 24 * time.Hour
	defaultTokenExpiry = 1 * time.Hour
	randomBytes        = 32
)

type JWT struct {
	secretKey          []byte
	signingMethod      jwt.SigningMethod
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	defaultTokenExpiry time.Duration
	tokenStore         repository.TokenStore
}

type AccessClaims struct {
	UserID  uint   `json:"userid"`
	TokenID string `json:"jti,omitempty"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID   uint   `json:"userid"`
	TokenID  string `json:"jti,omitempty"`
	FamilyID string `json:"family,omitempty"` // Token的组ID
	jwt.RegisteredClaims
}

type ResetPasswordClaims struct {
	Email   string `json:"email"`
	TokenID string `json:"jti,omitempty"`
	jwt.RegisteredClaims
}

func NewJwt(conf *viper.Viper, tokenStore repository.TokenStore) *JWT {
	key := conf.GetString("security.jwt.key")
	if len(key) < 32 {
		panic("jwt key length must be at least 32")
	}

	return &JWT{
		secretKey:          []byte(key),
		signingMethod:      jwt.SigningMethodHS256,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
		defaultTokenExpiry: defaultTokenExpiry,
		tokenStore:         tokenStore,
	}
}

// 生成 AccessToken & RefreshToken
func (j *JWT) GenerateTokenPair(ctx context.Context, uid uint, familyID string) (*v1.TokenPair, error) {
	// 生成 AccessToken
	accessTokenID, err := generateTokenID()
	if err != nil {
		return nil, err
	}
	accessClaims := AccessClaims{
		UserID:  uid,
		TokenID: accessTokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        accessTokenID,
		},
	}
	accessToken := jwt.NewWithClaims(j.signingMethod, accessClaims)
	accessTokenStr, err := accessToken.SignedString(j.secretKey)
	if err != nil {
		return nil, err
	}

	// 生成 RefreshToken
	refreshTokenID, err := generateTokenID()
	if err != nil {
		return nil, err
	}
	if familyID == "" {
		familyID, err = generateTokenID()
		if err != nil {
			return nil, err
		}
	}
	refreshClaims := RefreshClaims{
		UserID:   uid,
		TokenID:  refreshTokenID,
		FamilyID: familyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        refreshTokenID,
		},
	}
	refreshToken := jwt.NewWithClaims(j.signingMethod, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(j.secretKey)
	if err != nil {
		return nil, err
	}

	// 存储 RefreshToken
	if j.tokenStore != nil {
		err = j.tokenStore.StoreRefreshToken(ctx, refreshTokenID, familyID, uid, j.refreshTokenExpiry)
		if err != nil {
			return nil, err
		}
	}

	return &v1.TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(j.accessTokenExpiry.Seconds()),
	}, nil
}

// 刷新 AccessToken
func (j *JWT) RefreshAccessToken(ctx context.Context, refreshToken string) (*v1.TokenPair, error) {
	// 验证 RefreshToken
	refreshClaims, err := j.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 检查 RefreshToken 是否有效
	if j.tokenStore != nil {
		var valid bool
		valid, err = j.tokenStore.IsRefreshTokenValid(ctx, refreshClaims.TokenID, refreshClaims.FamilyID)
		if err != nil {
			return nil, err
		}
		if !valid {
			return nil, v1.ErrInvalidRefreshToken
		}
	}

	// 使旧 RefreshToken 失效
	if j.tokenStore != nil {
		err = j.tokenStore.InvalidateRefreshToken(ctx, refreshClaims.TokenID)
		if err != nil {
			return nil, err
		}
	}

	// 生成新的 TokenPair
	return j.GenerateTokenPair(ctx, refreshClaims.UserID, refreshClaims.FamilyID)
}

// 验证 AccessToken
func (j *JWT) ValidateAccessToken(ctx context.Context, accessToken string) (*AccessClaims, error) {
	accessToken = strings.TrimPrefix(accessToken, bearerPrefix)
	token, err := jwt.ParseWithClaims(accessToken, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, v1.ErrUnexpectedSigningMethod
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, v1.ErrInvalidAccessToken
	}

	if j.tokenStore != nil {
		revoked, err := j.tokenStore.IsAccessTokenRevoked(ctx, claims.TokenID)
		if err != nil {
			return nil, err
		}
		if revoked {
			return nil, v1.ErrInvalidAccessToken
		}
	}

	return claims, nil
}

func (j *JWT) InvalidateRefreshTokenByUserID(ctx context.Context, uid uint) error {
	if j.tokenStore != nil {
		return j.tokenStore.InvalidateRefreshTokenByUserID(ctx, uid)
	}
	return fmt.Errorf("token store is not initialized")
}

// 解析 RefreshToken
func (j *JWT) parseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, v1.ErrUnexpectedSigningMethod
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, v1.ErrInvalidRefreshToken
}

// 生成 ResetPasswordToken
func (j *JWT) GenerateResetPasswordToken(email string) (string, error) {
	tokenID, err := generateTokenID()
	if err != nil {
		return "", err
	}

	claims := ResetPasswordClaims{
		Email:   email,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.defaultTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(j.signingMethod, claims)
	return token.SignedString(j.secretKey)
}

// 生成 Token ID
func generateTokenID() (string, error) {
	b := make([]byte, randomBytes/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
