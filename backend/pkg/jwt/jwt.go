package jwt

import (
	v1 "backend/api/v1"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const (
	BearerPrefix             = "Bearer " // 请勿删除空格
	AccessTokenExpiry        = 15 * time.Minute
	RefreshTokenExpiry       = 30 * 24 * time.Hour
	ResetPasswordTokenExpiry = 1 * time.Hour
	RandomBytes              = 32
)

type JWT struct {
	secret                   []byte
	signingMethod            jwt.SigningMethod
	accessTokenExpiry        time.Duration
	refreshTokenExpiry       time.Duration
	resetPasswordTokenExpiry time.Duration
	tokenStore               TokenStore
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

func NewJwt(conf *viper.Viper, store TokenStore) *JWT {
	secret := conf.GetString("security.jwt.secret")
	if len(secret) < 32 {
		panic("jwt secret length must be at least 32")
	}

	return &JWT{
		secret:                   []byte(secret),
		signingMethod:            jwt.SigningMethodHS256,
		accessTokenExpiry:        AccessTokenExpiry,
		refreshTokenExpiry:       RefreshTokenExpiry,
		resetPasswordTokenExpiry: ResetPasswordTokenExpiry,
		tokenStore:               store,
	}
}

func (j *JWT) GenerateTokenPair(uid uint, familyID string) (*v1.TokenPair, error) {
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
	accessTokenStr, err := accessToken.SignedString(j.secret)
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token
	refreshTokenID, err := generateTokenID()
	if err != nil {
		return nil, err
	}
	familyID, err := generateTokenID()
	if err != nil {
		return nil, err
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
	refreshTokenStr, err := refreshToken.SignedString(j.secret)
	if err != nil {
		return nil, err
	}

	// 存储 Refresh Token 信息
	if j.tokenStore != nil {
		err = j.tokenStore.StoreRefreshToken(refreshTokenID, familyID, uid, j.refreshTokenExpiry)
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
func (j *JWT) RefreshAccessToken(refreshToken string) (*v1.TokenPair, error) {
	// 验证 Refresh Token
	refreshClaims, err := j.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 检查 Refresh Token 是否有效
	if j.tokenStore != nil {
		var valid bool
		valid, err = j.tokenStore.IsRefreshTokenValid(refreshClaims.TokenID, refreshClaims.FamilyID)
		if err != nil {
			return nil, err
		}
		if !valid {
			return nil, errors.New("invalid refresh token")
		}
	}

	// 使旧 Refresh Token 失效(可选，根据安全需求)
	if j.tokenStore != nil {
		err = j.tokenStore.InvalidateRefreshToken(refreshClaims.TokenID)
		if err != nil {
			return nil, err
		}
	}

	// 生成新的 Token Pair
	return j.GenerateTokenPair(refreshClaims.UserID)
}

// 验证 Access Token
func (j *JWT) ValidateAccessToken(accessToken string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid access token")
}

// 解析 Refresh Token
func (j *JWT) parseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}

// 生成 Token ID
func generateTokenID() (string, error) {
	b := make([]byte, RandomBytes/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
