package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// HTTPSRedirectMiddleware redirects HTTP to HTTPS
func HTTPSRedirectMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip if already HTTPS
		if c.Protocol() == "https" {
			return c.Next()
		}

		// Skip in development mode
		envType := os.Getenv("ENV_TYPE")
		if envType == "dev" || envType == "test" {
			return c.Next()
		}

		// Skip for health checks and metrics
		path := c.Path()
		if path == "/health" || path == "/health/ready" || path == "/health/live" || path == "/metrics" {
			return c.Next()
		}

		// Redirect to HTTPS
		return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
	}
}

// ForceHTTPSMiddleware enforces HTTPS without exceptions
func ForceHTTPSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Protocol() != "https" {
			return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
		}
		return c.Next()
	}
}
