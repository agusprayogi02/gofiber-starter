package handler

import (
	"starter-gofiber/helper"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SSEHandler struct{}

func NewSSEHandler() *SSEHandler {
	return &SSEHandler{}
}

// Connect handles SSE connection endpoint
func (h *SSEHandler) Connect(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return &helper.UnauthorizedError{
			Message: "User not authenticated",
			Order:   "H-SSE-Connect-1",
		}
	}

	helper.Info("SSE connection established",
		zap.Uint("user_id", userID),
		zap.String("ip", c.IP()),
	)

	return helper.SSEHandler(c, userID)
}

// Stats returns SSE hub statistics (admin only)
func (h *SSEHandler) Stats(c *fiber.Ctx) error {
	stats := helper.GetSSEStats()
	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// BroadcastMessage broadcasts a message to all connected clients (admin only)
func (h *SSEHandler) BroadcastMessage(c *fiber.Ctx) error {
	var req struct {
		Event string      `json:"event"`
		Data  interface{} `json:"data"`
	}

	if err := c.BodyParser(&req); err != nil {
		return &helper.BadRequestError{
			Message: "Invalid request body",
			Order:   "H-SSE-Broadcast-1",
		}
	}

	if req.Event == "" || req.Data == nil {
		return &helper.BadRequestError{
			Message: "Event and data are required",
			Order:   "H-SSE-Broadcast-2",
		}
	}

	helper.NotifyAll(req.Event, req.Data)

	helper.Info("SSE broadcast sent",
		zap.String("event", req.Event),
	)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Broadcast sent to all clients",
	})
}

// SendToUser sends a message to a specific user (admin only)
func (h *SSEHandler) SendToUser(c *fiber.Ctx) error {
	var req struct {
		UserID uint        `json:"user_id"`
		Event  string      `json:"event"`
		Data   interface{} `json:"data"`
	}

	if err := c.BodyParser(&req); err != nil {
		return &helper.BadRequestError{
			Message: "Invalid request body",
			Order:   "H-SSE-SendToUser-1",
		}
	}

	if req.UserID == 0 || req.Event == "" || req.Data == nil {
		return &helper.BadRequestError{
			Message: "UserID, event, and data are required",
			Order:   "H-SSE-SendToUser-2",
		}
	}

	helper.NotifyUser(req.UserID, req.Event, req.Data)

	helper.Info("SSE message sent to user",
		zap.Uint("user_id", req.UserID),
		zap.String("event", req.Event),
	)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Message sent to user",
	})
}

// Example: Send notification when new post is created
func (h *SSEHandler) NotifyNewPost(userID uint, post interface{}) {
	helper.NotifyAll("new_post", fiber.Map{
		"message":   "New post created",
		"post":      post,
		"timestamp": time.Now().Unix(),
	})
}

// Example: Send notification to specific user
func (h *SSEHandler) NotifyUserMessage(userID uint, message interface{}) {
	helper.NotifyUser(userID, "new_message", fiber.Map{
		"message":   "You have a new message",
		"data":      message,
		"timestamp": time.Now().Unix(),
	})
}

// Example: Progress update notification
func (h *SSEHandler) NotifyProgress(userID uint, progress int, status string) {
	helper.NotifyUser(userID, "progress_update", fiber.Map{
		"progress":  progress,
		"status":    status,
		"timestamp": time.Now().Unix(),
	})
}
