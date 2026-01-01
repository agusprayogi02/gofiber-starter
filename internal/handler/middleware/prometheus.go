package middleware

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Lightweight in-memory metrics (no external dependencies, ringan!)
var (
	totalRequests   uint64
	successRequests uint64
	errorRequests   uint64
	totalDuration   uint64 // nanoseconds
	inFlightReqs    int32
	startTime       = time.Now()
)

// GetMetrics returns current system metrics in JSON format
func GetMetrics() fiber.Map {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	total := atomic.LoadUint64(&totalRequests)
	success := atomic.LoadUint64(&successRequests)
	errors := atomic.LoadUint64(&errorRequests)
	duration := atomic.LoadUint64(&totalDuration)
	inFlight := atomic.LoadInt32(&inFlightReqs)

	// Calculate averages
	avgLatency := float64(0)
	if total > 0 {
		avgLatency = float64(duration) / float64(total) / 1e6 // Convert to milliseconds
	}

	successRate := float64(0)
	if total > 0 {
		successRate = float64(success) / float64(total) * 100
	}

	return fiber.Map{
		"uptime_seconds": time.Since(startTime).Seconds(),
		"requests": fiber.Map{
			"total":        total,
			"success":      success,
			"errors":       errors,
			"in_flight":    inFlight,
			"success_rate": successRate,
		},
		"performance": fiber.Map{
			"avg_latency_ms":         avgLatency,
			"total_duration_seconds": float64(duration) / 1e9,
		},
		"system": fiber.Map{
			"memory_alloc_mb": float64(m.Alloc) / 1024 / 1024,
			"memory_total_mb": float64(m.TotalAlloc) / 1024 / 1024,
			"memory_sys_mb":   float64(m.Sys) / 1024 / 1024,
			"gc_runs":         m.NumGC,
			"goroutines":      runtime.NumGoroutine(),
		},
	}
}

// RecordRequest updates request metrics (lightweight, thread-safe)
func RecordRequest(c *fiber.Ctx, start time.Time) {
	// Skip metrics endpoint
	if c.Path() == "/metrics" || c.Path() == "/health" {
		return
	}

	duration := time.Since(start)
	atomic.AddUint64(&totalRequests, 1)
	atomic.AddUint64(&totalDuration, uint64(duration.Nanoseconds()))

	if c.Response().StatusCode() < 400 {
		atomic.AddUint64(&successRequests, 1)
	} else {
		atomic.AddUint64(&errorRequests, 1)
	}
}

// MetricsMiddleware tracks HTTP requests (no external dependencies, super ringan!)
func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		atomic.AddInt32(&inFlightReqs, 1)
		defer atomic.AddInt32(&inFlightReqs, -1)

		start := time.Now()
		err := c.Next()

		RecordRequest(c, start)
		return err
	}
}
