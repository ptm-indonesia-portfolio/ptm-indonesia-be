package bootstrap

import (
	"ptm-indonesia/config"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type HTTPApplication struct {
	Config *config.AppConfig
	App    *fiber.App
	Logger *logrus.Logger
}

func NewHTTPApplication(cfg *config.AppConfig, app *fiber.App, logger *logrus.Logger) *HTTPApplication {
	return &HTTPApplication{
		Config: cfg,
		App:    app,
		Logger: logger,
	}
}
