package helper

import (
	"ptm-indonesia/config"
	"ptm-indonesia/model"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func NewRateLimiter(cfg *config.AppConfig, localizer *Localizer) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.RateLimit.MaxRequests,
		Expiration: cfg.RateLimit.Window,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		Next: func(c fiber.Ctx) bool {
			return c.Path() == "/api/v1/health"
		},
		LimitReached: func(c fiber.Ctx) error {
			locale := localizer.Resolve(c.Get("Accept-Language"))
			return c.Status(fiber.StatusTooManyRequests).JSON(model.ErrorResponse{
				Errors: []string{localizer.MustLocalize(locale, "error.rate_limit")},
			})
		},
	})
}
