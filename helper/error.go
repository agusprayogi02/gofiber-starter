package helper

import (
	"errors"

	"starter-gofiber/dto"

	"github.com/gofiber/fiber/v2"
)

type NotFoundError struct {
	Message string
	Order   string
}

type BadRequestError struct {
	Message string
	Order   string
}

type InternalServerError struct {
	Message string
	Order   string
}

type UnauthorizedError struct {
	Message string
	Order   string
}

type UnprocessableEntityError struct {
	Message string
	Data    interface{}
	Order   string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func (e *UnprocessableEntityError) Error() string {
	return e.Message
}

func ErrorHelper(c *fiber.Ctx, err error) error {
	var statusCode int
	var order string

	switch err.(type) {
	case *NotFoundError:
		statusCode = fiber.StatusNotFound
		order = err.(*NotFoundError).Order
	case *BadRequestError:
		statusCode = fiber.StatusBadRequest
		order = err.(*BadRequestError).Order
	case *InternalServerError:
		statusCode = fiber.StatusInternalServerError
		order = err.(*InternalServerError).Order
	case *UnauthorizedError:
		statusCode = fiber.StatusUnauthorized
		order = err.(*UnauthorizedError).Order
	case *UnprocessableEntityError:
		statusCode = fiber.StatusUnprocessableEntity
		order = err.(*UnprocessableEntityError).Order
	default:
		var e *fiber.Error
		if errors.As(err, &e) {
			statusCode = e.Code
			order = e.Message
		}
	}

	rest := dto.ErrorResponse{
		Code:      statusCode,
		Order:     &order,
		Message:   err.Error(),
		Timestamp: TimeNow(),
	}
	if statusCode == fiber.StatusUnprocessableEntity {
		rest.Data = err.(*UnprocessableEntityError).Data
	}

	return c.Status(statusCode).JSON(rest)
}
