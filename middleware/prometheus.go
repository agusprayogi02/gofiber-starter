package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request counter
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP request duration histogram
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// HTTP request size histogram
	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path"},
	)

	// HTTP response size histogram
	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path"},
	)

	// Active requests gauge
	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)
)

// PrometheusMiddleware records HTTP metrics
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Increment in-flight requests
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Get request size
		reqSize := len(c.Request().Body())

		// Continue processing
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()

		// Skip metrics endpoint itself
		if path == "/metrics" {
			return err
		}

		// Record request count
		httpRequestsTotal.WithLabelValues(
			method,
			path,
			strconv.Itoa(status),
		).Inc()

		// Record duration
		httpRequestDuration.WithLabelValues(
			method,
			path,
		).Observe(duration)

		// Record request size
		httpRequestSize.WithLabelValues(
			method,
			path,
		).Observe(float64(reqSize))

		// Record response size
		respSize := len(c.Response().Body())
		httpResponseSize.WithLabelValues(
			method,
			path,
		).Observe(float64(respSize))

		return err
	}
}
