package routes

import (
	"ptm-indonesia/controllers"

	"github.com/gofiber/fiber/v3"
)

func RegisterAPIRoutes(
	api fiber.Router,
	healthController *controllers.HealthController,
	authController *controllers.AuthController,
	userController *controllers.UserController,
) {
	RegisterHealthRoutes(api, healthController)
	RegisterAuthRoutes(api, authController)
	RegisterUserRoutes(api, userController)
}
