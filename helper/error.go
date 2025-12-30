package helper

import (
	"errors"
	"fmt"

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
type ForbiddenError struct {
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

func (e *ForbiddenError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func (e *UnprocessableEntityError) Error() string {
	return e.Message
}

func ErrorHelper(c *fiber.Ctx, err error) error {
	var order string
	rest := dto.ErrorResponse{
		Message:   err.Error(),
		Timestamp: TimeNow(),
	}

	switch err.(type) {
	case *NotFoundError:
		rest.Code = fiber.StatusNotFound
		order = err.(*NotFoundError).Order
	case *BadRequestError:
		rest.Code = fiber.StatusBadRequest
		order = err.(*BadRequestError).Order
	case *InternalServerError:
		rest.Code = fiber.StatusInternalServerError
		order = err.(*InternalServerError).Order
	case *UnauthorizedError:
		rest.Code = fiber.StatusUnauthorized
		order = err.(*UnauthorizedError).Order
	case *ForbiddenError:
		rest.Code = fiber.StatusForbidden
		order = err.(*ForbiddenError).Order
	case *UnprocessableEntityError:
		rest.Code = fiber.StatusUnprocessableEntity
		order = err.(*UnprocessableEntityError).Order
		rest.Data = err.(*UnprocessableEntityError).Data
	default:
		var e *fiber.Error
		if errors.As(err, &e) {
			rest.Code = e.Code
			order = fmt.Sprintf("Handling Error: %s", e.Message)
		}
	}
	rest.Order = &order

	return c.Status(rest.Code).JSON(rest)
}
