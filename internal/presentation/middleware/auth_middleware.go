package middleware

import (
	"mob/ddd-template/internal/presentation/dto"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func IsLogin(ctx fiber.Ctx) error {
	authHeader := ctx.GetHeaders()["Authorization"]
	if len(authHeader) != 1 {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_DENIED_ACCESS, "invalid token format")
		return ctx.Status(fiber.StatusUnauthorized).JSON(res)
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_DENIED_ACCESS, "invalid token format")
		return ctx.Status(fiber.StatusUnauthorized).JSON(res)
	}

	return ctx.Next()
}
