package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserClaims struct {
	UserID string
	// TODO Add user role
}
type TokenManager struct {
	signingKey []byte
}

func NewTokenManager(signingKey string) *TokenManager {
	return &TokenManager{signingKey: []byte(signingKey)}
}

func (tm *TokenManager) Generate(claims UserClaims) (string, error) {
	jwtClaims := jwt.MapClaims{
		"sub": claims.UserID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(tm.signingKey)
}
