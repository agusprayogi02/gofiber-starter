package http

import (
	"fmt"
	"time"

	"starter-gofiber/internal/worker"
	"starter-gofiber/pkg/crypto"
	"starter-gofiber/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// JobHandler handles job-related HTTP requests
type JobHandler struct{}

// NewJobHandler creates new job handler
func NewJobHandler() *JobHandler {
	return &JobHandler{}
}

// TriggerEmailJob endpoint to trigger email job
func (h *JobHandler) TriggerEmailJob(c *fiber.Ctx) error {
	type Request struct {
		To      string `json:"to" validate:"required,email"`
		Subject string `json:"subject" validate:"required"`
		Body    string `json:"body" validate:"required"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Dispatch job to queue
	if err := worker.SendEmailJob(req.To, req.Subject, req.Body, nil); err != nil {
		logger.Error(fmt.Sprintf("Failed to enqueue email job: %v", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to queue email",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Email queued successfully",
		"status":  "processing",
	})
}

// TriggerExportJob endpoint to trigger export job
func (h *JobHandler) TriggerExportJob(c *fiber.Ctx) error {
	type Request struct {
		ExportType string `json:"export_type" validate:"required,oneof=csv excel pdf"`
		Query      string `json:"query"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from JWT token
	userClaims, err := crypto.GetUserFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}
	userID := userClaims.ID

	// Generate filename
	filename := fmt.Sprintf("export_%d_%s.%s", userID, time.Now().Format("20060102_150405"), req.ExportType)

	// Dispatch export job
	if err := worker.ProcessExportJob(userID, req.ExportType, req.Query, filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to queue export",
		})
	}

	// Send notification after 5 seconds (simulate processing time)
	go func() {
		time.Sleep(5 * time.Second)
		worker.SendNotificationJob(
			userID,
			"email",
			"Export Ready",
			fmt.Sprintf("Your %s export is ready for download", req.ExportType),
			map[string]interface{}{
				"filename": filename,
				"type":     req.ExportType,
			},
		)
	}()

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":  "Export job queued successfully",
		"filename": filename,
		"status":   "processing",
	})
}

// GetJobStats returns job queue statistics
func (h *JobHandler) GetJobStats(c *fiber.Ctx) error {
	queues, err := worker.ListAllQueues()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get queue stats",
		})
	}

	stats := make([]fiber.Map, 0)
	for _, q := range queues {
		stats = append(stats, fiber.Map{
			"queue":      q.Queue,
			"pending":    q.Pending,
			"active":     q.Active,
			"scheduled":  q.Scheduled,
			"retry":      q.Retry,
			"archived":   q.Archived,
			"completed":  q.Completed,
			"paused":     q.Paused,
			"size":       q.Size,
			"latency_ms": q.Latency.Milliseconds(),
		})
	}

	return c.JSON(fiber.Map{
		"queues": stats,
	})
}

// PauseQueue pause specific queue
func (h *JobHandler) PauseQueue(c *fiber.Ctx) error {
	queueName := c.Params("queue")
	if queueName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Queue name is required",
		})
	}

	if err := worker.PauseQueue(queueName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to pause queue",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Queue '%s' paused successfully", queueName),
	})
}

// ResumeQueue resume paused queue
func (h *JobHandler) ResumeQueue(c *fiber.Ctx) error {
	queueName := c.Params("queue")
	if queueName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Queue name is required",
		})
	}

	if err := worker.UnpauseQueue(queueName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to resume queue",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Queue '%s' resumed successfully", queueName),
	})
}

// RetryFailedJob retry specific failed job
func (h *JobHandler) RetryFailedJob(c *fiber.Ctx) error {
	type Request struct {
		Queue  string `json:"queue" validate:"required"`
		TaskID string `json:"task_id" validate:"required"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := worker.RetryTask(req.Queue, req.TaskID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retry task",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Task retry initiated successfully",
	})
}

// DeletePendingTasks delete all pending tasks in queue
func (h *JobHandler) DeletePendingTasks(c *fiber.Ctx) error {
	queueName := c.Params("queue")
	if queueName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Queue name is required",
		})
	}

	count, err := worker.DeleteAllPendingTasks(queueName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete pending tasks",
		})
	}

	return c.JSON(fiber.Map{
		"message":       fmt.Sprintf("Deleted %d pending tasks from queue '%s'", count, queueName),
		"deleted_count": count,
	})
}
