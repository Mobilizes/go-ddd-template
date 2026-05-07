package middleware

import (
	"errors"
	"mob/ddd-template/internal/presentation/dto"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func IsLogin(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	tokenString, ok := parseBearerToken(authHeader)
	if !ok {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_DENIED_ACCESS, "invalid token format")
		return ctx.Status(fiber.StatusUnauthorized).JSON(res)
	}

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil || !token.Valid {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_DENIED_ACCESS, "invalid or expired token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(res)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_DENIED_ACCESS, "invalid token claims")
		return ctx.Status(fiber.StatusUnauthorized).JSON(res)
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_DENIED_ACCESS, "failed to login with access token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(res)
	}

	ctx.Locals("accessToken", token)
	ctx.Locals("userId", userID)

	return ctx.Next()
}

func parseBearerToken(authHeader string) (string, bool) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return "", false
	}

	return parts[1], true
}
