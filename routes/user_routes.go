package routes

import (
	"ptm-indonesia/controllers"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(router fiber.Router, userController *controllers.UserController) {
	users := router.Group("/users", userController.RequireLogin)
	users.Get("", userController.List)
	users.Post("", userController.Create)
	users.Get("/:id", userController.Detail)
	users.Put("/:id", userController.Update)
	users.Delete("/:id", userController.Delete)
}
