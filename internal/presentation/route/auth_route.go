package route

import (
	"mob/ddd-template/internal/presentation/handler"
	"mob/ddd-template/internal/presentation/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuth(h handler.AuthHandler, r *fiber.App) {
	authGroup := r.Group("/api")
	{
		authGroup.Post("/login", h.Login)
		authGroup.Post("/refresh", middleware.IsLogin, h.Refresh)
		authGroup.Post("/logout", middleware.IsLogin, h.Logout)
		authGroup.Post("/logout-all", middleware.IsLogin, h.LogoutAll)
	}
}
