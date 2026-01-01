package http

import (
	"starter-gofiber/internal/handler/middleware"

	"github.com/gofiber/fiber/v2"
)

type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// Metrics returns lightweight system metrics in JSON (no external dependencies, super ringan!)
func (h *MetricsHandler) Metrics(c *fiber.Ctx) error {
	return c.JSON(middleware.GetMetrics())
}
