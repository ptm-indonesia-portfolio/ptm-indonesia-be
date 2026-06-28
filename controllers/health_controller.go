package controllers

import (
	"context"

	"ptm-indonesia/helper"
	servicesContract "ptm-indonesia/services/contract"

	"github.com/gofiber/fiber/v3"
)

type HealthController struct {
	healthService servicesContract.HealthService
	responder     *helper.Responder
	localizer     *helper.Localizer
}

func NewHealthController(
	healthService servicesContract.HealthService,
	responder *helper.Responder,
	localizer *helper.Localizer,
) *HealthController {
	return &HealthController{
		healthService: healthService,
		responder:     responder,
		localizer:     localizer,
	}
}

func (h *HealthController) Check(c fiber.Ctx) error {
	locale := h.localizer.Resolve(c.Get("Accept-Language"))
	response := h.healthService.Check(context.Background())

	statusCode := fiber.StatusOK
	messageID := "health.success"

	if response.Database != "up" {
		statusCode = fiber.StatusInternalServerError
		messageID = "health.failed"
	}

	return h.responder.Success(
		c,
		statusCode,
		h.localizer.MustLocalize(locale, messageID),
		response,
	)
}
