package middleware

import (
	"starter-gofiber/pkg/logger"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestLogger middleware logs all HTTP requests using fiberzap
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID before fiberzap logging
		requestID := uuid.New().String()
		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)

		// Use fiberzap middleware for structured HTTP logging
		return fiberzap.New(fiberzap.Config{
			Logger: logger.Logger,
			Fields: []string{"ip", "latency", "status", "method", "url", "error"},
		})(c)
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(c *fiber.Ctx) string {
	requestID := c.Locals("requestID")
	if requestID == nil {
		return ""
	}
	return requestID.(string)
}
