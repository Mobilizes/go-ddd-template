package route

import (
	"mob/ddd-template/internal/presentation/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuth(h handler.AuthHandler, r *fiber.App) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.Post("/login", h.Login)
		authGroup.Post("/refresh", h.Refresh)
		authGroup.Post("/logout", h.Logout)
		authGroup.Post("/logout-all", h.LogoutAll)
	}
}
