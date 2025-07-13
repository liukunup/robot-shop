package jwt

import (
	v1 "backend/api/v1"
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const (
	BearerPrefix       = "Bearer " // 请勿删除空格
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
	DefaultTokenExpiry = 1 * time.Hour
	RandomBytes        = 32
)

type JWT struct {
	key                []byte
	signingMethod      jwt.SigningMethod
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	defaultTokenExpiry time.Duration
	tokenStore         TokenStore
}

type AccessClaims struct {
	UserID  uint   `json:"userid"`
	TokenID string `json:"jti,omitempty"` // 唯一标识符
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID   uint   `json:"userid"`
	TokenID  string `json:"jti,omitempty"`    // 唯一标识符
	FamilyID string `json:"family,omitempty"` // Token簇ID
	jwt.RegisteredClaims
}

type ResetPasswordClaims struct {
	Email   string `json:"email"`
	TokenID string `json:"jti,omitempty"` // 唯一标识符
	jwt.RegisteredClaims
}

func NewJwt(conf *viper.Viper) *JWT {
	key := conf.GetString("security.jwt.key")
	if len(key) < 32 {
		panic("jwt key length must be at least 32")
	}

	// tokenStore := conf.GetString("security.jwt.token_store")
	// switch tokenStore {
	// case "redis":
	// 	tokenStore = NewRedisTokenStore(conf)
	// case "memory":
	// 	tokenStore = NewMemoryTokenStore()
	// default:
	// 	panic("unsupported token store type")
	// }

	tokenStore := NewMemoryTokenStore()

	return &JWT{
		key:                []byte(key),
		signingMethod:      jwt.SigningMethodHS256,
		accessTokenExpiry:  AccessTokenExpiry,
		refreshTokenExpiry: RefreshTokenExpiry,
		defaultTokenExpiry: DefaultTokenExpiry,
		tokenStore:         tokenStore,
	}
}

// 生成 Access Token & Refresh Token
func (j *JWT) GenerateTokenPair(ctx context.Context, uid uint, familyID string) (*v1.TokenPair, error) {
	// 生成 Access Token
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
	accessTokenStr, err := accessToken.SignedString(j.key)
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token
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
	refreshTokenStr, err := refreshToken.SignedString(j.key)
	if err != nil {
		return nil, err
	}

	// 存储 Refresh Token
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

// 刷新 Access Token
func (j *JWT) RefreshAccessToken(ctx context.Context, refreshToken string) (*v1.TokenPair, error) {
	// 验证 Refresh Token
	refreshClaims, err := j.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 检查 Refresh Token 是否有效
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

	// 使旧 Refresh Token 失效
	if j.tokenStore != nil {
		err = j.tokenStore.InvalidateRefreshToken(ctx, refreshClaims.TokenID)
		if err != nil {
			return nil, err
		}
	}

	// 生成新的 Token Pair
	return j.GenerateTokenPair(ctx, refreshClaims.UserID, refreshClaims.FamilyID)
}

// 验证 Access Token
func (j *JWT) ValidateAccessToken(ctx context.Context, accessToken string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, v1.ErrUnexpectedSigningMethod
		}
		return j.key, nil
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

func (j *JWT) InvalidateRefreshTokenFamilyByUserID(ctx context.Context, uid uint) error {
	if j.tokenStore != nil {
		return j.tokenStore.InvalidateRefreshTokenFamilyByUserID(ctx, uid)
	}
	return nil
}

// 解析 Refresh Token
func (j *JWT) parseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, v1.ErrUnexpectedSigningMethod
		}
		return j.key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, v1.ErrInvalidRefreshToken
}

// 生成 Token ID
func generateTokenID() (string, error) {
	b := make([]byte, RandomBytes/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
