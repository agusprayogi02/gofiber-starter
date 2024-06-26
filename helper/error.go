package helper

import (
	"github.com/gofiber/fiber/v2"
	"starter-gofiber/dto"
)

type NotFoundError struct {
	Message string
}

type BadRequestError struct {
	Message string
}

type InternalServerError struct {
	Message string
}

type UnauthorizedError struct {
	Message string
}

type UnprocessableEntityError struct {
	Message string
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

	switch err.(type) {
	case *NotFoundError:
		statusCode = fiber.StatusNotFound
	case *BadRequestError:
		statusCode = fiber.StatusBadRequest
	case *InternalServerError:
		statusCode = fiber.StatusInternalServerError
	case *UnauthorizedError:
		statusCode = fiber.StatusUnauthorized
	case *UnprocessableEntityError:
		statusCode = fiber.StatusUnprocessableEntity
	default:
		statusCode = fiber.StatusInternalServerError
	}

	return c.Status(statusCode).JSON(dto.ErrorResponse{
		Message:   err.Error(),
		Timestamp: TimeNow(),
	})
}
