package helper

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

// InitSentry initializes Sentry error tracking
func InitSentry(dsn, environment string) error {
	if dsn == "" {
		Info("Sentry DSN not configured, error tracking disabled")
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      environment,
		Debug:            environment == "dev",
		TracesSampleRate: 1.0,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Filter out certain errors if needed
			return event
		},
	})
	if err != nil {
		return err
	}

	Info("Sentry initialized successfully")
	return nil
}

// CaptureError sends an error to Sentry
func CaptureError(err error, ctx ...*fiber.Ctx) {
	if err == nil {
		return
	}

	// Build scope with context
	sentry.WithScope(func(scope *sentry.Scope) {
		// Add fiber context if provided
		if len(ctx) > 0 && ctx[0] != nil {
			c := ctx[0]

			// Add request context (not using SetRequest as it expects http.Request)
			scope.SetContext("request", map[string]interface{}{
				"url":    c.OriginalURL(),
				"method": c.Method(),
				"headers": map[string]string{
					"User-Agent":   c.Get("User-Agent"),
					"Content-Type": c.Get("Content-Type"),
					"X-Request-ID": GetRequestID(c),
				},
			})

			// Add user information if available
			userID := c.Locals("user_id")
			if userID != nil {
				scope.SetUser(sentry.User{
					ID:        userID.(string),
					IPAddress: c.IP(),
				})
			} else {
				scope.SetUser(sentry.User{
					IPAddress: c.IP(),
				})
			}

			// Add custom tags
			scope.SetTag("request_id", GetRequestID(c))
			scope.SetTag("method", c.Method())
			scope.SetTag("path", c.Path())
			scope.SetTag("status_code", string(c.Response().StatusCode()))
		}

		// Capture the error
		sentry.CaptureException(err)
	})
}

// CaptureMessage sends a custom message to Sentry
func CaptureMessage(message string, level sentry.Level) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)
		sentry.CaptureMessage(message)
	})
}

// FlushSentry ensures all events are sent before shutdown
func FlushSentry() {
	sentry.Flush(2 * time.Second)
}

// RecoverWithSentry recovers from panic and sends to Sentry
func RecoverWithSentry(c *fiber.Ctx) {
	if err := recover(); err != nil {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetContext("request", map[string]interface{}{
				"url":    c.OriginalURL(),
				"method": c.Method(),
				"headers": map[string]string{
					"User-Agent":   c.Get("User-Agent"),
					"Content-Type": c.Get("Content-Type"),
					"X-Request-ID": GetRequestID(c),
				},
			})

			scope.SetTag("panic", "true")
			scope.SetTag("request_id", GetRequestID(c))

			sentry.CurrentHub().RecoverWithContext(
				c.Context(),
				err,
			)
		})

		// Re-throw panic
		panic(err)
	}
}

// GetRequestID gets request ID from fiber context
func GetRequestID(c *fiber.Ctx) string {
	requestID := c.Locals("requestid")
	if requestID != nil {
		return requestID.(string)
	}
	return c.Get("X-Request-ID", "unknown")
}
