package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// Metrics returns Prometheus metrics
func (h *MetricsHandler) Metrics(c *fiber.Ctx) error {
	// Use fasthttpadaptor to convert http.Handler to fiber.Handler
	handler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	handler(c.Context())
	return nil
}
