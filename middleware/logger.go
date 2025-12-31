package middleware

import (
	"time"

	"starter-gofiber/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestLogger middleware logs all HTTP requests with structured logging
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID
		requestID := uuid.New().String()
		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)

		// Record start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get response status
		status := c.Response().StatusCode()

		// Log request
		helper.LogRequest(
			c.Method(),
			c.Path(),
			status,
			duration,
			zap.String("request_id", requestID),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.Int("body_size", len(c.Response().Body())),
		)

		return err
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
