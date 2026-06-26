package helper

import (
	"ptm-indonesia/model"

	"github.com/gofiber/fiber/v3"
)

type Responder struct{}

func NewResponder() *Responder {
	return &Responder{}
}

func (r *Responder) Success(c fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(model.SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func (r *Responder) Errors(c fiber.Ctx, statusCode int, errors []string) error {
	return c.Status(statusCode).JSON(model.ErrorResponse{
		Errors: errors,
	})
}
