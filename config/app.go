package config

import (
	"time"

	"starter-gofiber/dto"
	"starter-gofiber/middleware"
	"starter-gofiber/variables"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func App(app *fiber.App) {
	// CORS middleware
	app.Use(cors.New())

	// Structured logging middleware (replaces default logger)
	app.Use(middleware.RequestLogger())

	// Sentry error tracking middleware
	app.Use(middleware.SentryMiddleware())

	// Prometheus metrics middleware
	app.Use(middleware.PrometheusMiddleware())

	// Rate limiter
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        30,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			msg := "C-LimitReached"
			err := dto.ErrorResponse{
				Code:      fiber.StatusTooManyRequests,
				Order:     &msg,
				Message:   "Too many requests",
				Timestamp: time.Now().Format(variables.FORMAT_TIME),
			}
			return c.Status(fiber.StatusTooManyRequests).JSON(err)
		},
		Storage: STORAGE,
	}))
}
