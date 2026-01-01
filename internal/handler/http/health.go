package http

import (
	"time"

	"starter-gofiber/internal/config"
	"starter-gofiber/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Services  map[string]ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

var startTime = time.Now()

// Health returns the health status of the application
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	uptime := time.Since(startTime)

	services := make(map[string]ServiceInfo)

	// Check database connection
	dbStatus := "healthy"
	dbMessage := ""
	if config.DB != nil {
		sqlDB, err := config.DB.DB()
		if err != nil {
			dbStatus = "unhealthy"
			dbMessage = err.Error()
			logger.Error("Database health check failed", zap.Error(err))
		} else {
			err = sqlDB.Ping()
			if err != nil {
				dbStatus = "unhealthy"
				dbMessage = err.Error()
				logger.Error("Database ping failed", zap.Error(err))
			}
		}
	} else {
		dbStatus = "unavailable"
		dbMessage = "Database not initialized"
	}

	services["database"] = ServiceInfo{
		Status:  dbStatus,
		Message: dbMessage,
	}

	// Overall status
	overallStatus := "healthy"
	if dbStatus != "healthy" {
		overallStatus = "degraded"
	}

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    uptime.String(),
		Services:  services,
	}

	statusCode := fiber.StatusOK
	if overallStatus == "degraded" {
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(response)
}

// Ready returns readiness probe (for Kubernetes)
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	// Check if all critical services are ready
	if config.DB == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "not ready",
			"message": "Database not initialized",
		})
	}

	sqlDB, err := config.DB.DB()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "not ready",
			"message": "Database connection error",
		})
	}

	err = sqlDB.Ping()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "not ready",
			"message": "Database ping failed",
		})
	}

	return c.JSON(fiber.Map{
		"status": "ready",
	})
}

// Live returns liveness probe (for Kubernetes)
func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "alive",
	})
}
