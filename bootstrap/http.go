package bootstrap

import (
	"errors"

	"ptm-indonesia/config"
	"ptm-indonesia/controllers"
	"ptm-indonesia/helper"
	"ptm-indonesia/model"
	"ptm-indonesia/routes"
	"ptm-indonesia/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/sirupsen/logrus"
)

func NewFiberApp(
	cfg *config.AppConfig,
	logger *logrus.Logger,
	localizer *helper.Localizer,
	requestValidator *validation.RequestValidator,
	healthController *controllers.HealthController,
	authController *controllers.AuthController,
	userController *controllers.UserController,
) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:         cfg.App.Name,
		ReadTimeout:     cfg.App.ReadTimeout,
		WriteTimeout:    cfg.App.WriteTimeout,
		IdleTimeout:     cfg.App.IdleTimeout,
		StructValidator: requestValidator,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			locale := localizer.Resolve(c.Get("Accept-Language"))

			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
					Errors: requestValidator.TranslateValidationErrors(locale, validationErrors),
				})
			}

			var fiberError *fiber.Error
			if errors.As(err, &fiberError) {
				return c.Status(fiberError.Code).JSON(model.ErrorResponse{
					Errors: []string{fiberError.Message},
				})
			}

			logger.WithError(err).Error("request failed")

			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Errors: []string{localizer.MustLocalize(locale, "error.internal_server")},
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.App.CORSOrigins,
		AllowMethods:     []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	api := app.Group("/api/v1")
	api.Use(helper.NewRateLimiter(cfg, localizer))
	routes.RegisterAPIRoutes(api, healthController, authController, userController)

	return app
}
