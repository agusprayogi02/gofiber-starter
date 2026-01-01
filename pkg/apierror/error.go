package apierror

import (
	"errors"
	"fmt"
	"time"

	"starter-gofiber/dto"
	"starter-gofiber/variables"

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

type TooManyRequestsError struct {
	Message string
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

func (e *TooManyRequestsError) Error() string {
	return e.Message
}

// ErrorHelper handles errors and returns appropriate HTTP response
// Note: This function depends on dto and utils/date - will need to update after those are moved
func ErrorHelper(c *fiber.Ctx, err error) error {
	var order string
	rest := dto.ErrorResponse{
		Message:   err.Error(),
		Timestamp: TimeNow(),
	}

	switch err := err.(type) {
	case *NotFoundError:
		rest.Code = fiber.StatusNotFound
		order = err.Order
	case *BadRequestError:
		rest.Code = fiber.StatusBadRequest
		order = err.Order
	case *InternalServerError:
		rest.Code = fiber.StatusInternalServerError
		order = err.Order
	case *UnauthorizedError:
		rest.Code = fiber.StatusUnauthorized
		order = err.Order
	case *ForbiddenError:
		rest.Code = fiber.StatusForbidden
		order = err.Order
	case *UnprocessableEntityError:
		rest.Code = fiber.StatusUnprocessableEntity
		order = err.Order
		rest.Data = err.Data
	case *TooManyRequestsError:
		rest.Code = fiber.StatusTooManyRequests
		order = err.Order
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

func TimeNow() string {
	return time.Now().Format(variables.FORMAT_TIME)
}
