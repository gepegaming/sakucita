package server

import (
	"errors"

	"sakucita/internal/domain"
	"sakucita/internal/shared/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func fiberErrorHandler(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := any("internal server error")

	// error validation
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		status = fiber.StatusUnprocessableEntity
		message = utils.GenerateMessageValidation(err)

		return c.Status(status).JSON(domain.ErrorResponse{
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  message,
		})
	}

	// handle errro dari domain
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.Code).JSON(domain.ErrorResponse{
			Message: appErr.Message,
			Errors:  appErr.Error(),
		})
	}

	// fiber error
	var fe *fiber.Error
	if errors.As(err, &fe) {
		status = fe.Code
		message = fe.Message
		return c.Status(status).JSON(domain.ErrorResponse{
			Errors: message,
		})
	}

	// fallback
	return c.Status(status).JSON(domain.ErrorResponse{
		Message: fiber.ErrInternalServerError.Message,
		Errors:  message,
	})
}
