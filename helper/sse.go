package helper

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SSEClient represents a connected SSE client
type SSEClient struct {
	ID       string
	UserID   uint
	Channel  chan SSEMessage
	Context  context.Context
	Cancel   context.CancelFunc
	LastPing time.Time
}

// SSEMessage represents a message to be sent via SSE
type SSEMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	ID    string      `json:"id,omitempty"`
	Retry int         `json:"retry,omitempty"` // Milliseconds
}

// SSEHub manages all SSE connections
type SSEHub struct {
	clients    map[string]*SSEClient
	userIndex  map[uint][]string // UserID -> ClientIDs
	broadcast  chan SSEMessage
	register   chan *SSEClient
	unregister chan string
	mu         sync.RWMutex
}

var (
	sseHub     *SSEHub
	sseHubOnce sync.Once
)

// GetSSEHub returns the singleton SSE hub instance
func GetSSEHub() *SSEHub {
	sseHubOnce.Do(func() {
		sseHub = &SSEHub{
			clients:    make(map[string]*SSEClient),
			userIndex:  make(map[uint][]string),
			broadcast:  make(chan SSEMessage, 256),
			register:   make(chan *SSEClient),
			unregister: make(chan string),
		}
		go sseHub.run()
	})
	return sseHub
}

// run manages the SSE hub event loop
func (h *SSEHub) run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.userIndex[client.UserID] = append(h.userIndex[client.UserID], client.ID)
			h.mu.Unlock()

			Info("SSE client registered",
				zap.String("client_id", client.ID),
				zap.Uint("user_id", client.UserID),
			)

		case clientID := <-h.unregister:
			h.mu.Lock()
			if client, ok := h.clients[clientID]; ok {
				// Remove from user index
				if clients, exists := h.userIndex[client.UserID]; exists {
					newClients := make([]string, 0, len(clients)-1)
					for _, id := range clients {
						if id != clientID {
							newClients = append(newClients, id)
						}
					}
					if len(newClients) > 0 {
						h.userIndex[client.UserID] = newClients
					} else {
						delete(h.userIndex, client.UserID)
					}
				}

				close(client.Channel)
				delete(h.clients, clientID)

				Info("SSE client unregistered",
					zap.String("client_id", clientID),
					zap.Uint("user_id", client.UserID),
				)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Channel <- message:
				default:
					// Channel full, skip this message for this client
					Warn("SSE message dropped - channel full",
						zap.String("client_id", client.ID),
					)
				}
			}
			h.mu.RUnlock()

		case <-ticker.C:
			// Send ping to all clients to keep connection alive
			h.sendPing()
		}
	}
}

// sendPing sends a ping message to all connected clients
func (h *SSEHub) sendPing() {
	h.mu.RLock()
	defer h.mu.RUnlock()

	pingMessage := SSEMessage{
		Event: "ping",
		Data:  map[string]interface{}{"timestamp": time.Now().Unix()},
	}

	for _, client := range h.clients {
		select {
		case client.Channel <- pingMessage:
			client.LastPing = time.Now()
		default:
			// Skip ping if channel is full
		}
	}
}

// NewClient creates a new SSE client
func (h *SSEHub) NewClient(userID uint, ctx context.Context) *SSEClient {
	clientCtx, cancel := context.WithCancel(ctx)

	client := &SSEClient{
		ID:       uuid.New().String(),
		UserID:   userID,
		Channel:  make(chan SSEMessage, 100),
		Context:  clientCtx,
		Cancel:   cancel,
		LastPing: time.Now(),
	}

	h.register <- client
	return client
}

// RemoveClient removes a client from the hub
func (h *SSEHub) RemoveClient(clientID string) {
	h.unregister <- clientID
}

// Broadcast sends a message to all connected clients
func (h *SSEHub) Broadcast(message SSEMessage) {
	select {
	case h.broadcast <- message:
	default:
		Warn("SSE broadcast dropped - channel full")
	}
}

// SendToUser sends a message to all connections of a specific user
func (h *SSEHub) SendToUser(userID uint, message SSEMessage) {
	h.mu.RLock()
	clientIDs, exists := h.userIndex[userID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clientID := range clientIDs {
		if client, ok := h.clients[clientID]; ok {
			select {
			case client.Channel <- message:
			default:
				Warn("SSE message dropped - channel full",
					zap.String("client_id", clientID),
					zap.Uint("user_id", userID),
				)
			}
		}
	}
}

// SendToClient sends a message to a specific client
func (h *SSEHub) SendToClient(clientID string, message SSEMessage) {
	h.mu.RLock()
	client, ok := h.clients[clientID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	select {
	case client.Channel <- message:
	default:
		Warn("SSE message dropped - channel full",
			zap.String("client_id", clientID),
		)
	}
}

// GetClientCount returns the total number of connected clients
func (h *SSEHub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetUserClientCount returns the number of connections for a specific user
func (h *SSEHub) GetUserClientCount(userID uint) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userIndex[userID])
}

// SSEHandler handles SSE connections
func SSEHandler(c *fiber.Ctx, userID uint) error {
	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("X-Accel-Buffering", "no") // Disable nginx buffering

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		hub := GetSSEHub()
		client := hub.NewClient(userID, c.Context())
		defer hub.RemoveClient(client.ID)

		// Send initial connection success message
		sendSSEMessage(w, SSEMessage{
			Event: "connected",
			Data: map[string]interface{}{
				"client_id": client.ID,
				"timestamp": time.Now().Unix(),
			},
		})

		// Listen for messages
		for {
			select {
			case message, ok := <-client.Channel:
				if !ok {
					return
				}

				if err := sendSSEMessage(w, message); err != nil {
					Error("Failed to send SSE message",
						zap.Error(err),
						zap.String("client_id", client.ID),
					)
					return
				}

				if err := w.Flush(); err != nil {
					Error("Failed to flush SSE message",
						zap.Error(err),
						zap.String("client_id", client.ID),
					)
					return
				}

			case <-client.Context.Done():
				return
			}
		}
	})

	return nil
}

// sendSSEMessage formats and sends an SSE message
func sendSSEMessage(w *bufio.Writer, message SSEMessage) error {
	// Event type
	if message.Event != "" {
		if _, err := fmt.Fprintf(w, "event: %s\n", message.Event); err != nil {
			return err
		}
	}

	// Message ID (for Last-Event-ID support)
	if message.ID != "" {
		if _, err := fmt.Fprintf(w, "id: %s\n", message.ID); err != nil {
			return err
		}
	}

	// Retry interval
	if message.Retry > 0 {
		if _, err := fmt.Fprintf(w, "retry: %d\n", message.Retry); err != nil {
			return err
		}
	}

	// Data (JSON encoded)
	dataJSON, err := json.Marshal(message.Data)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "data: %s\n\n", dataJSON); err != nil {
		return err
	}

	return nil
}

// NotifyUser sends a notification to a specific user
func NotifyUser(userID uint, event string, data interface{}) {
	hub := GetSSEHub()
	hub.SendToUser(userID, SSEMessage{
		Event: event,
		Data:  data,
		ID:    uuid.New().String(),
	})
}

// NotifyAll broadcasts a notification to all connected clients
func NotifyAll(event string, data interface{}) {
	hub := GetSSEHub()
	hub.Broadcast(SSEMessage{
		Event: event,
		Data:  data,
		ID:    uuid.New().String(),
	})
}

// SSEStats returns SSE hub statistics
type SSEStats struct {
	TotalClients int            `json:"total_clients"`
	UserClients  map[uint]int   `json:"user_clients"`
}

// GetSSEStats returns current SSE hub statistics
func GetSSEStats() SSEStats {
	hub := GetSSEHub()
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	userClients := make(map[uint]int)
	for userID, clients := range hub.userIndex {
		userClients[userID] = len(clients)
	}

	return SSEStats{
		TotalClients: len(hub.clients),
		UserClients:  userClients,
	}
}
