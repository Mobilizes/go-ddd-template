package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func RegisterHealthCheck(r *fiber.App) {
	r.Get(healthcheck.LivenessEndpoint, healthcheck.New())
}
