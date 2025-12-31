# Server-Sent Events (SSE) Documentation

Dokumentasi lengkap untuk implementasi Server-Sent Events (SSE) - real-time, unidirectional communication dari server ke client.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Setup](#setup)
- [Client Connection](#client-connection)
- [Sending Messages](#sending-messages)
- [Event Types](#event-types)
- [Admin API](#admin-api)
- [Client Implementation](#client-implementation)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## Overview

Server-Sent Events (SSE) adalah teknologi untuk server mengirim real-time updates ke client melalui HTTP. Berbeda dengan WebSocket yang bidirectional, SSE adalah unidirectional (server → client).

### Features

✅ **Real-time Updates** - Server push notifications tanpa polling  
✅ **Auto Reconnect** - Client otomatis reconnect jika terputus  
✅ **Event Types** - Support multiple event types  
✅ **User Targeting** - Send ke specific user atau broadcast ke semua  
✅ **Connection Management** - Track semua active connections  
✅ **Keep-Alive** - Automatic ping untuk maintain connection  

### Use Cases

- 📢 **Notifications** - Real-time user notifications
- 📊 **Live Updates** - Dashboard updates, stock prices, etc.
- 📈 **Progress Tracking** - Long-running task progress
- 💬 **Chat Updates** - New message notifications
- 🔔 **Event Broadcasting** - System-wide announcements

---

## Architecture

### Component Overview

```
┌─────────────┐
│   Client    │
│  Browser    │
└──────┬──────┘
       │ HTTP GET /sse/stream
       │ (with auth token)
       ▼
┌──────────────────┐
│  SSE Handler     │
│  (Fiber)         │
└─────────┬────────┘
          │
          ▼
┌──────────────────┐
│   SSE Hub        │
│  - Manage clients│
│  - Route messages│
│  - Keep-alive    │
└─────────┬────────┘
          │
          ├─────────► Client 1 (User A)
          ├─────────► Client 2 (User A)
          ├─────────► Client 3 (User B)
          └─────────► ...
```

### SSE Hub

Hub singleton yang mengelola semua SSE connections:
- Track semua connected clients
- Index clients by user ID untuk targeted messaging
- Broadcast channel untuk global messages
- Auto ping setiap 30 detik untuk keep-alive

---

## Setup

### 1. Router Configuration

SSE routes sudah terdaftar otomatis di [router/sse.go](../router/sse.go):

```go
// Public SSE stream endpoint (requires auth)
GET /sse/stream

// Admin endpoints
GET  /api/sse/stats           // Connection statistics
POST /api/sse/broadcast       // Broadcast to all
POST /api/sse/send-to-user    // Send to specific user
```

### 2. Dependencies

No additional dependencies required. Menggunakan:
- `github.com/gofiber/fiber/v2` - HTTP framework
- `github.com/google/uuid` - Message IDs
- Built-in Go channels untuk message passing

---

## Client Connection

### Endpoint

```
GET /sse/stream
Authorization: Bearer <jwt_token>
```

### Response Headers

```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
Transfer-Encoding: chunked
X-Accel-Buffering: no
```

### Event Stream Format

```
event: connected
data: {"client_id":"123e4567-e89b-12d3-a456-426614174000","timestamp":1704067200}

event: ping
data: {"timestamp":1704067230}

event: notification
data: {"title":"New message","body":"You have a new message","timestamp":1704067250}
```

---

## Sending Messages

### 1. Notify Specific User

```go
import "starter-gofiber/helper"

// Send notification to user ID 123
helper.NotifyUser(123, "notification", fiber.Map{
    "title":   "New Message",
    "body":    "You have a new message from John",
    "icon":    "/icons/message.png",
    "link":    "/messages/456",
})
```

### 2. Broadcast to All Users

```go
// Send to all connected clients
helper.NotifyAll("announcement", fiber.Map{
    "title":   "System Maintenance",
    "message": "System will be down for 5 minutes",
    "time":    "2024-01-01 02:00:00",
})
```

### 3. Using SSE Hub Directly

```go
hub := helper.GetSSEHub()

// Send to specific client ID
hub.SendToClient("client-id-123", helper.SSEMessage{
    Event: "custom_event",
    Data:  map[string]interface{}{"key": "value"},
    ID:    "msg-001",
    Retry: 3000, // Retry after 3 seconds
})

// Send to all connections of a user
hub.SendToUser(123, helper.SSEMessage{
    Event: "update",
    Data:  map[string]interface{}{"status": "completed"},
})

// Broadcast to everyone
hub.Broadcast(helper.SSEMessage{
    Event: "global_update",
    Data:  map[string]interface{}{"version": "2.0.0"},
})
```

---

## Event Types

### Built-in Events

#### 1. `connected`
Sent immediately after client connects:
```json
{
  "client_id": "123e4567-e89b-12d3-a456-426614174000",
  "timestamp": 1704067200
}
```

#### 2. `ping`
Keep-alive message (every 30 seconds):
```json
{
  "timestamp": 1704067230
}
```

### Custom Events

You can define your own events:

```go
// Notification event
helper.NotifyUser(userID, "notification", fiber.Map{
    "title": "New Message",
    "body":  "...",
    "type":  "message",
})

// Progress update
helper.NotifyUser(userID, "progress", fiber.Map{
    "task_id":  "task-123",
    "progress": 75,
    "status":   "processing",
})

// Real-time data update
helper.NotifyAll("data_update", fiber.Map{
    "entity": "product",
    "action": "created",
    "id":     456,
})
```

### Example Use Cases

#### New Post Created
```go
func (h *PostHandler) Create(c *fiber.Ctx) error {
    // ... create post logic ...
    
    // Notify all users about new post
    helper.NotifyAll("new_post", fiber.Map{
        "post_id":   post.ID,
        "title":     post.Title,
        "author":    post.Author.Name,
        "timestamp": time.Now().Unix(),
    })
    
    return c.JSON(post)
}
```

#### User Message Notification
```go
func (h *MessageHandler) Send(c *fiber.Ctx) error {
    // ... send message logic ...
    
    // Notify recipient
    helper.NotifyUser(message.RecipientID, "new_message", fiber.Map{
        "message_id": message.ID,
        "from":       message.Sender.Name,
        "preview":    message.Body[:50],
        "timestamp":  time.Now().Unix(),
    })
    
    return c.JSON(message)
}
```

#### Export Progress
```go
func (h *ExportHandler) ExportData(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    
    // Start export in background
    go func() {
        total := 1000
        for i := 0; i <= total; i += 100 {
            progress := (i * 100) / total
            
            helper.NotifyUser(userID, "export_progress", fiber.Map{
                "progress": progress,
                "current":  i,
                "total":    total,
                "status":   "processing",
            })
            
            time.Sleep(2 * time.Second) // Simulate work
        }
        
        helper.NotifyUser(userID, "export_complete", fiber.Map{
            "file_url": "/downloads/export.xlsx",
            "status":   "completed",
        })
    }()
    
    return c.JSON(fiber.Map{"message": "Export started"})
}
```

---

## Admin API

### 1. Get Statistics

```bash
GET /api/sse/stats
Authorization: Bearer <admin_token>
```

Response:
```json
{
  "success": true,
  "data": {
    "total_clients": 15,
    "user_clients": {
      "1": 2,
      "5": 1,
      "10": 3,
      "15": 1
    }
  }
}
```

### 2. Broadcast Message

```bash
POST /api/sse/broadcast
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "event": "announcement",
  "data": {
    "title": "System Maintenance",
    "message": "Scheduled maintenance at 2AM",
    "severity": "warning"
  }
}
```

### 3. Send to Specific User

```bash
POST /api/sse/send-to-user
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "user_id": 123,
  "event": "admin_message",
  "data": {
    "title": "Account Verification",
    "message": "Please verify your account"
  }
}
```

---

## Client Implementation

### JavaScript/Browser

#### Basic Connection

```javascript
// Connect to SSE stream
const eventSource = new EventSource('/sse/stream', {
  withCredentials: true, // Include cookies
});

// Listen for connection
eventSource.addEventListener('connected', (e) => {
  const data = JSON.parse(e.data);
  console.log('Connected:', data.client_id);
});

// Listen for ping (keep-alive)
eventSource.addEventListener('ping', (e) => {
  const data = JSON.parse(e.data);
  console.log('Ping received:', new Date(data.timestamp * 1000));
});

// Listen for custom events
eventSource.addEventListener('notification', (e) => {
  const notification = JSON.parse(e.data);
  showNotification(notification.title, notification.body);
});

// Handle errors
eventSource.onerror = (error) => {
  console.error('SSE Error:', error);
  // Browser will automatically try to reconnect
};

// Close connection when done
function disconnect() {
  eventSource.close();
}
```

#### With JWT Token

```javascript
// Since EventSource doesn't support custom headers,
// we need to pass token via query parameter or use fetch API

// Option 1: Query parameter (less secure)
const token = localStorage.getItem('token');
const eventSource = new EventSource(`/sse/stream?token=${token}`);

// Option 2: Use fetch with streaming (recommended)
async function connectSSE() {
  const token = localStorage.getItem('token');
  
  const response = await fetch('/sse/stream', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  const reader = response.body.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    const chunk = decoder.decode(value);
    const lines = chunk.split('\n\n');

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = JSON.parse(line.substring(6));
        handleSSEMessage(data);
      }
    }
  }
}

connectSSE();
```

#### React Example

```jsx
import { useEffect, useState } from 'react';

function useSSE() {
  const [notifications, setNotifications] = useState([]);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const eventSource = new EventSource('/sse/stream');

    eventSource.addEventListener('connected', () => {
      setConnected(true);
    });

    eventSource.addEventListener('notification', (e) => {
      const notification = JSON.parse(e.data);
      setNotifications(prev => [...prev, notification]);
    });

    eventSource.onerror = () => {
      setConnected(false);
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return { notifications, connected };
}

// Usage
function NotificationCenter() {
  const { notifications, connected } = useSSE();

  return (
    <div>
      <div>Status: {connected ? '🟢 Connected' : '🔴 Disconnected'}</div>
      <ul>
        {notifications.map((notif, idx) => (
          <li key={idx}>
            <strong>{notif.title}</strong>
            <p>{notif.body}</p>
          </li>
        ))}
      </ul>
    </div>
  );
}
```

#### Vue.js Example

```vue
<template>
  <div>
    <div>Status: {{ connected ? '🟢 Connected' : '🔴 Disconnected' }}</div>
    <div v-for="(notif, idx) in notifications" :key="idx">
      <h3>{{ notif.title }}</h3>
      <p>{{ notif.body }}</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      eventSource: null,
      connected: false,
      notifications: [],
    };
  },
  
  mounted() {
    this.connectSSE();
  },
  
  beforeUnmount() {
    if (this.eventSource) {
      this.eventSource.close();
    }
  },
  
  methods: {
    connectSSE() {
      this.eventSource = new EventSource('/sse/stream');
      
      this.eventSource.addEventListener('connected', () => {
        this.connected = true;
      });
      
      this.eventSource.addEventListener('notification', (e) => {
        const notification = JSON.parse(e.data);
        this.notifications.push(notification);
      });
      
      this.eventSource.onerror = () => {
        this.connected = false;
      };
    },
  },
};
</script>
```

---

## Best Practices

### 1. Connection Management

✅ **DO:**
- Close connections when user logs out
- Implement reconnection logic with exponential backoff
- Monitor connection status in UI

❌ **DON'T:**
- Keep multiple connections open for same user
- Forget to close connections on page unload

### 2. Message Size

✅ **DO:**
- Keep messages small and focused
- Send only necessary data
- Use message IDs for deduplication

❌ **DON'T:**
- Send large payloads (>1KB) frequently
- Include unnecessary data in events

### 3. Event Types

✅ **DO:**
- Use descriptive event names
- Group related events by prefix (e.g., `user.login`, `user.logout`)
- Document all event types

❌ **DON'T:**
- Use generic event names like `update` or `data`
- Change event schemas without versioning

### 4. Error Handling

```javascript
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;

function connectWithRetry() {
  const eventSource = new EventSource('/sse/stream');
  
  eventSource.addEventListener('connected', () => {
    reconnectAttempts = 0; // Reset on successful connection
  });
  
  eventSource.onerror = () => {
    eventSource.close();
    
    if (reconnectAttempts < maxReconnectAttempts) {
      const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000);
      setTimeout(() => {
        reconnectAttempts++;
        connectWithRetry();
      }, delay);
    } else {
      console.error('Max reconnection attempts reached');
    }
  };
}
```

### 5. Performance

✅ **DO:**
- Use event delegation for multiple listeners
- Debounce UI updates from high-frequency events
- Monitor memory usage for long-running connections

```javascript
// Debounce updates
let updateTimeout;
eventSource.addEventListener('live_data', (e) => {
  clearTimeout(updateTimeout);
  updateTimeout = setTimeout(() => {
    updateUI(JSON.parse(e.data));
  }, 100);
});
```

### 6. Security

✅ **DO:**
- Always require authentication
- Validate user has permission for events
- Sanitize data before sending
- Use HTTPS in production

❌ **DON'T:**
- Send sensitive data without encryption
- Allow unauthenticated connections
- Trust client-provided event data

---

## Troubleshooting

### Issue: Connection keeps dropping

**Causes:**
- Proxy/load balancer timeout
- Network instability
- Server restart

**Solutions:**
```javascript
// Implement auto-reconnect
eventSource.addEventListener('error', () => {
  setTimeout(() => {
    location.reload(); // Or reconnect logic
  }, 3000);
});
```

**Server-side:**
```go
// Increase keep-alive frequency
ticker := time.NewTicker(15 * time.Second) // Instead of 30s
```

### Issue: Messages not received

**Check:**
1. Is client connected? Check `connected` event
2. Are messages being sent? Check server logs
3. Is event name correct?
4. Is user ID correct?

**Debug:**
```javascript
// Listen to all events
eventSource.onmessage = (e) => {
  console.log('Message:', e);
};
```

### Issue: Memory leak

**Cause:** Not closing connections properly

**Solution:**
```javascript
// React
useEffect(() => {
  const es = new EventSource('/sse/stream');
  return () => es.close(); // Cleanup
}, []);

// Vue
beforeUnmount() {
  this.eventSource.close();
}

// Plain JS
window.addEventListener('beforeunload', () => {
  eventSource.close();
});
```

### Issue: CORS errors

**Solution:**
Update CORS configuration in [config/app.go](../config/app.go):

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://yourdomain.com",
    AllowMethods: "GET,POST,OPTIONS",
    AllowHeaders: "Authorization,Content-Type",
    AllowCredentials: true,
}))
```

### Issue: Nginx buffering

**Problem:** Nginx buffers SSE responses

**Solution:**
Add to nginx config:
```nginx
location /sse/ {
    proxy_pass http://backend;
    proxy_buffering off;
    proxy_cache off;
    proxy_set_header Connection '';
    proxy_http_version 1.1;
    chunked_transfer_encoding off;
}
```

---

## Monitoring

### Get Connection Stats

```go
stats := helper.GetSSEStats()
fmt.Printf("Total clients: %d\n", stats.TotalClients)
fmt.Printf("User connections: %v\n", stats.UserClients)
```

### Log Events

```go
helper.Info("SSE message sent",
    zap.Uint("user_id", userID),
    zap.String("event", "notification"),
    zap.Int("client_count", hub.GetUserClientCount(userID)),
)
```

### Prometheus Metrics

Add custom metrics:
```go
var (
    sseConnections = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "sse_active_connections",
        Help: "Number of active SSE connections",
    })
    
    sseMessagesSent = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "sse_messages_sent_total",
        Help: "Total number of SSE messages sent",
    })
)

// Update metrics
sseConnections.Set(float64(hub.GetClientCount()))
sseMessagesSent.Inc()
```

---

## Comparison: SSE vs WebSocket

| Feature | SSE | WebSocket |
|---------|-----|-----------|
| Direction | Unidirectional (Server → Client) | Bidirectional |
| Protocol | HTTP | WebSocket (ws://) |
| Browser Support | ✅ All modern browsers | ✅ All modern browsers |
| Auto Reconnect | ✅ Built-in | ❌ Manual |
| Complexity | Simple | More complex |
| Use Case | Notifications, updates | Chat, gaming |
| Firewall | ✅ Rarely blocked | ⚠️ Sometimes blocked |

**When to use SSE:**
- One-way updates (server → client)
- Notifications, live feeds, dashboards
- Simpler implementation needed

**When to use WebSocket:**
- Two-way communication needed
- Chat applications, gaming
- High-frequency bidirectional data

---

## Summary

✅ **Implemented Features:**
- SSE Hub untuk connection management
- User-targeted dan broadcast messaging
- Auto keep-alive dengan ping
- Admin API untuk monitoring dan messaging
- Complete client examples (JS, React, Vue)

🎯 **Key Benefits:**
- Real-time updates tanpa polling
- Auto-reconnect built-in
- Simple HTTP-based protocol
- Scalable untuk ribuan connections

📚 **Next Steps:**
- Test dengan client favorit Anda
- Customize event types untuk use case Anda
- Add monitoring metrics
- Implement message persistence (optional)

---

**For more information, see:**
- [helper/sse.go](../helper/sse.go) - Core SSE implementation
- [handler/sse.go](../handler/sse.go) - SSE handlers
- [router/sse.go](../router/sse.go) - SSE routes
