package jwt

import (
	v1 "backend/api/v1"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const (
	BearerPrefix       = "Bearer " // 请勿删除空格
	DefaultTokenExpiry = 24 * time.Hour
	ResetTokenExpiry   = 1 * time.Hour
	ResetTokenLength   = 32
)

type JWT struct {
	key               []byte
	signingMethod     jwt.SigningMethod
	accessTokenExpiry time.Duration
	resetTokenExpiry  time.Duration
	tokenStore        TokenStore // 用于黑名单/白名单功能
}

type TokenStore interface {
	IsRevoked(tokenID string) (bool, error)
	Revoke(tokenID string, expiry time.Time) error
}

type MyCustomClaims struct {
	UserId  uint
	TokenID string `json:"jti,omitempty"` // 唯一标识符，可用于撤销
	jwt.RegisteredClaims
}

type ResetTokenClaims struct {
	Email   string `json:"email"`
	TokenID string `json:"jti,omitempty"` // 唯一标识符，可用于撤销
	jwt.RegisteredClaims
}

func NewJwt(conf *viper.Viper, store TokenStore) (*JWT, error) {
	key := conf.GetString("security.jwt.key")
	if len(key) < 32 {
		return nil, v1.ErrInvalidKeyLength
	}

	tokenExpiry := conf.GetDuration("security.jwt.expiry")
	if tokenExpiry == 0 {
		tokenExpiry = DefaultTokenExpiry
	}

	return &JWT{
		key:               []byte(key),
		signingMethod:     jwt.SigningMethodHS256,
		accessTokenExpiry: tokenExpiry,
		resetTokenExpiry:  ResetTokenExpiry,
		tokenStore:        store,
	}, nil
}

func (j *JWT) GenToken(userId uint, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "",
			Subject:   "",
			ID:        "",
			Audience:  []string{},
		},
	})

	// Sign and get the complete encoded token as a string using the key
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (*MyCustomClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("token is empty")
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (j *JWT) RevokeToken(tokenID string, expiry time.Time) error {
	if j.tokenStore == nil {
		return nil
	}
	return j.tokenStore.Revoke(tokenID, expiry)
}
