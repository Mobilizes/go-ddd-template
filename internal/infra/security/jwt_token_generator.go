package security

import (
	"crypto/rand"
	"encoding/base64"
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

func (g *JWTTokenGenerator) GenerateAccessToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(g.expiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.secretKey))
}

func (g *JWTTokenGenerator) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
