package routes

import (
	"ptm-indonesia/controllers"

	"github.com/gofiber/fiber/v3"
)

func RegisterHealthRoutes(router fiber.Router, healthController *controllers.HealthController) {
	router.Get("/health", healthController.Check)
}
