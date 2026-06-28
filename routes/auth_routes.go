package routes

import (
	"ptm-indonesia/controllers"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(router fiber.Router, authController *controllers.AuthController) {
	auth := router.Group("/auth")
	auth.Get("/google/login", authController.GoogleLogin)
	auth.Get("/google/callback", authController.GoogleCallback)
	auth.Post("/refresh", authController.Refresh)
	auth.Get("/me", authController.Me)
	auth.Post("/logout", authController.Logout)
}
