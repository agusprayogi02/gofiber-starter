package router

import (
	"starter-gofiber/internal/handler/http"
	"starter-gofiber/internal/handler/middleware"

	"github.com/gofiber/fiber/v2"
)

func SSERouter(app *fiber.App) {
	sseHandler := http.NewSSEHandler()

	// SSE endpoint (requires authentication)
	app.Get("/sse/stream", middleware.AuthMiddleware(), sseHandler.Connect)

	// Admin endpoints for sending messages
	authz := middleware.LoadAuthzMiddleware()
	admin := app.Group("/api/sse", middleware.AuthMiddleware(), authz.RequiresPermissions([]string{"admin:write"}))
	{
		admin.Get("/stats", sseHandler.Stats)
		admin.Post("/broadcast", sseHandler.BroadcastMessage)
		admin.Post("/send-to-user", sseHandler.SendToUser)
	}
}
