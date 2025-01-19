package http

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Error  string   `json:"error"`
	Errors []string `json:"errors,omitempty"`
}

var ErrInternalServer = fiber.NewError(fiber.StatusInternalServerError, "internal server error")

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	var v validator.ValidationErrors
	if errors.As(err, &v) {
		code = fiber.StatusBadRequest
		var errors []string
		for _, err := range v {
			errors = append(errors, fmt.Sprintf("%s: %s", err.Field(), err.Tag()))
		}
		return c.Status(code).JSON(ErrorResponse{Error: "bad request", Errors: errors})
	}

	if e != nil {
		if code == fiber.StatusInternalServerError {
			log.Error().Stack().Err(err).Msg("failed to handle request")
		}
		return c.Status(code).JSON(ErrorResponse{Error: e.Error()})
	}
	return nil
}
