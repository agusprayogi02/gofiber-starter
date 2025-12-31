package router

import (
	"starter-gofiber/config"
	"starter-gofiber/dto"
	"starter-gofiber/handler"
	"starter-gofiber/helper"
	"starter-gofiber/middleware"
	"starter-gofiber/variables"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func AppRouter(app *fiber.App) {
	// Only initialize enforcer if not already set (e.g., from tests)
	if config.Enforcer == nil {
		enforcer, err := casbin.NewEnforcer("./assets/rbac/model.conf", "assets/rbac/policy.csv")
		if err != nil {
			panic(err)
		}
		err = config.InitializePermission(enforcer)
		if err != nil {
			panic(err)
		}
	}

	// Recover middleware for production
	if config.ENV.ENV_TYPE != "dev" {
		app.Use(recover.New())
	}

	// Health check endpoints
	healthHandler := handler.NewHealthHandler()
	app.Get("/health", healthHandler.Health)
	app.Get("/health/ready", healthHandler.Ready)
	app.Get("/health/live", healthHandler.Live)

	// Metrics endpoint for Prometheus
	metricsHandler := handler.NewMetricsHandler()
	app.Get("/metrics", metricsHandler.Metrics)

	// Ping endpoint (legacy)
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
			Code:      fiber.StatusOK,
			Message:   "pong",
			Timestamp: helper.TimeNow(),
		})
	})

	// Static files
	static := app.Group(variables.STATIC_PATH, middleware.AuthMiddleware())
	authz := middleware.LoadAuthzMiddleware()
	static.Use(authz.RequiresPermissions([]string{"files:read"}))
	static.Static("/", "./public")
	app.Static("/favicon.ico", "./public/favicon.ico")

	// API routes
	api := app.Group("/api")
	auth := api.Group("/auth")
	NewAuthentication(auth, config.Enforcer)
	NewPostRouter(api)

	// SSE routes
	SSERouter(app)
}
