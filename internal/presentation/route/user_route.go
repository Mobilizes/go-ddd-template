package route

import (
	"mob/ddd-template/internal/presentation/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterUser(h handler.UserHandler, r *fiber.App) {
	userGroup := r.Group("/api/user")
	{
		userGroup.Post("/", h.Create)
		userGroup.Get("/", h.GetAll)
		userGroup.Get("/:id", h.GetById)
		userGroup.Delete("/:id", h.Delete)
	}
}
