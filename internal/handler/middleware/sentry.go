package middleware

import (
	"starter-gofiber/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// SentryMiddleware catches panics and sends errors to Sentry
func SentryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to error
				var err error
				switch x := r.(type) {
				case string:
					err = fiber.NewError(fiber.StatusInternalServerError, x)
				case error:
					err = x
				default:
					err = fiber.NewError(fiber.StatusInternalServerError, "Unknown panic")
				}

				// Capture error in Sentry with context
				logger.CaptureException(err)

				// Re-panic to let Fiber's recover handle it
				panic(r)
			}
		}()

		// Continue chain
		err := c.Next()
		// Capture errors from handlers
		if err != nil {
			logger.CaptureException(err)
		}

		return err
	}
}
