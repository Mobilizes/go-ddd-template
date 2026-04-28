package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct {
	secretKey string
	expiry    time.Duration
}

func NewJWTTokenGenerator(secretKey string, expiry time.Duration) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		secretKey: secretKey,
		expiry:    expiry,
	}
}

func (g *JWTTokenGenerator) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(g.expiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.secretKey))
}
