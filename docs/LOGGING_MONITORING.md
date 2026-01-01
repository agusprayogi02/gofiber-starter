# Logging & Monitoring Guide

Dokumentasi lengkap untuk sistem logging dan monitoring pada aplikasi starter-gofiber.

## Table of Contents

- [Structured Logging](#structured-logging)
- [Request Tracking](#request-tracking)
- [Error Tracking dengan Sentry](#error-tracking-dengan-sentry)
- [Metrics dengan Prometheus](#metrics-dengan-prometheus)
- [Health Check Endpoints](#health-check-endpoints)
- [Database Query Logging](#database-query-logging)

---

## Structured Logging

Aplikasi menggunakan **Zap** sebagai structured logging library dengan **fiberzap** middleware untuk Fiber HTTP logging integration. Format JSON untuk production dan colored console untuk development.

### Konfigurasi

Logging dikonfigurasi otomatis berdasarkan environment (`ENV_TYPE`):

**Development Mode:**
- Output: Console dengan warna
- Level: Debug
- Format: Human-readable
- File logging: Disabled

**Production Mode:**
- Output: File (`logs/app.log`)
- Level: Info
- Format: JSON
- Log rotation: 100MB, max 7 backup files

### Penggunaan

```go
import "starter-gofiber/pkg/apierror"
import "go.uber.org/zap"

// Basic logging
logger.Info("User logged in")
helper.Warn("Cache miss")
apierror.Error("Database error", zap.Error(err))
helper.Debug("Variable value", zap.String("key", value))

// Fatal (logs and exits)
helper.Fatal("Cannot start server", zap.Error(err))

// Request logging (automatic via middleware)
helper.LogRequest(c, duration, statusCode)

// Error logging with context
helper.LogError(err, map[string]interface{}{
    "user_id": userID,
    "action": "create_post",
})

// Database query logging (automatic via GORM logger)
helper.LogDBQuery(sql, duration, rows, err)
```

### HTTP Request Logging with fiberzap

HTTP requests are automatically logged using the **fiberzap** middleware (official Fiber-Zap integration):

```go
// middleware/logger.go
func RequestLogger() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Generate request ID before fiberzap logging
        requestID := uuid.New().String()
        c.Locals("requestID", requestID)
        c.Set("X-Request-ID", requestID)

        // Use fiberzap middleware for structured HTTP logging
        return fiberzap.New(fiberzap.Config{
            Logger: helper.Logger,
            Fields: []string{"ip", "latency", "status", "method", "url", "error"},
        })(c)
    }
}
```

**Benefits of fiberzap:**
- ✅ Official Fiber integration for Zap
- ✅ Consistent structured logging format
- ✅ Automatic performance metrics (latency)
- ✅ Request/response details captured
- ✅ Error logging with stack traces

### Log Fields

Setiap log entry otomatis menyertakan:
- `timestamp`: ISO8601 format
- `level`: debug, info, warn, error, fatal
- `msg`: Pesan log
- `caller`: File dan line number
- Custom fields yang ditambahkan

---

## Request Tracking

Setiap HTTP request mendapatkan unique Request ID untuk tracing.

### Features

- **UUID Request ID**: Unique identifier per request
- **Request ID Header**: `X-Request-ID` dikirim dalam response
- **Context Storage**: Request ID tersimpan di Fiber context
- **Automatic Logging**: Semua request dicatat dengan detail lengkap

### Log Format

```json
{
  "level": "info",
  "timestamp": "2025-06-15T10:30:45+07:00",
  "caller": "middleware/logger.go:45",
  "msg": "HTTP Request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/posts",
  "status": 201,
  "duration_ms": 45.23,
  "ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "body_size": 1024
}
```

### Menggunakan Request ID

```go
// Di handler/service
func (h *PostHandler) Create(c *fiber.Ctx) error {
    requestID := middleware.GetRequestID(c)
    logger.Info("Creating post", 
        zap.String("request_id", requestID),
        zap.String("user_id", userID),
    )
    // ...
}
```

---

## Error Tracking dengan Sentry

Integrasi dengan Sentry untuk tracking error production.

### Konfigurasi

Tambahkan DSN Sentry di `.env`:

```env
SENTRY_DSN=https://your-dsn@sentry.io/project-id
```

Jika `SENTRY_DSN` kosong, error tracking akan di-skip (tidak error).

### Features

- **Automatic Error Capture**: Semua error dari handlers otomatis tercatat
- **Request Context**: URL, method, headers, user info tersimpan
- **Panic Recovery**: Panic ditangkap dan dikirim ke Sentry
- **Custom Tags**: Request ID, method, path, status code
- **User Tracking**: User ID dan IP address
- **Environment Separation**: dev/staging/production environments

### Manual Error Capture

```go
import "starter-gofiber/pkg/apierror"

// Capture error dengan context
func (h *Handler) SomeMethod(c *fiber.Ctx) error {
    err := someOperation()
    if err != nil {
        helper.CaptureError(err, c)
        return err
    }
}

// Capture custom message
helper.CaptureMessage("Important event occurred", sentry.LevelInfo)
```

### Sentry Dashboard

Error di Sentry akan menyertakan:
- Full stack trace
- Request details (URL, method, headers)
- User information (ID, IP)
- Request ID untuk cross-reference dengan logs
- Custom tags untuk filtering

---

## Metrics dengan Prometheus

Prometheus metrics tersedia di `/metrics` endpoint.

### Available Metrics

#### HTTP Metrics

```prometheus
# Total requests
http_requests_total{method="POST", path="/api/posts", status="201"}

# Request duration (histogram)
http_request_duration_seconds{method="GET", path="/api/posts"}

# Request size (histogram)
http_request_size_bytes{method="POST", path="/api/posts"}

# Response size (histogram)
http_response_size_bytes{method="GET", path="/api/posts"}

# In-flight requests (gauge)
http_requests_in_flight
```

#### Database Metrics

```prometheus
# Connections in use
db_connections_in_use

# Idle connections
db_connections_idle
```

#### Authentication Metrics

```prometheus
# Auth attempts
auth_attempts_total{result="success"}
auth_attempts_total{result="failure"}

# Tokens created
auth_tokens_created_total
```

### Prometheus Configuration

Tambahkan job di `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'gofiber-app'
    static_configs:
      - targets: ['localhost:3000']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Grafana Dashboard

Contoh query untuk dashboard:

```promql
# Request rate
rate(http_requests_total[5m])

# Error rate
rate(http_requests_total{status=~"5.."}[5m])

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Database connections
db_connections_in_use + db_connections_idle

# Auth success rate
rate(auth_attempts_total{result="success"}[5m]) / 
rate(auth_attempts_total[5m])
```

### Custom Metrics

Tambahkan custom metrics di `middleware/prometheus.go`:

```go
var customCounter = promauto.NewCounter(
    prometheus.CounterOpts{
        Name: "my_custom_metric_total",
        Help: "Description of metric",
    },
)

// Increment di code
customCounter.Inc()
```

---

## Health Check Endpoints

Tiga health check endpoints untuk monitoring dan orchestration.

### Endpoints

#### 1. `/health` - Comprehensive Health Check

Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-06-15T10:30:45+07:00",
  "uptime": "2h30m15s",
  "services": {
    "database": {
      "status": "healthy",
      "message": "Connected"
    }
  }
}
```

Status codes:
- `200 OK`: Semua services healthy
- `503 Service Unavailable`: Ada service yang down

#### 2. `/health/ready` - Readiness Probe

Untuk Kubernetes readiness probe. Memeriksa apakah aplikasi siap menerima traffic.

Response:
```json
{
  "status": "ready",
  "timestamp": "2025-06-15T10:30:45+07:00"
}
```

#### 3. `/health/live` - Liveness Probe

Untuk Kubernetes liveness probe. Memeriksa apakah aplikasi masih running.

Response:
```json
{
  "status": "alive",
  "timestamp": "2025-06-15T10:30:45+07:00"
}
```

### Kubernetes Configuration

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    image: gofiber-app:latest
    livenessProbe:
      httpGet:
        path: /health/live
        port: 3000
      initialDelaySeconds: 30
      periodSeconds: 10
    readinessProbe:
      httpGet:
        path: /health/ready
        port: 3000
      initialDelaySeconds: 5
      periodSeconds: 5
```

### Docker Compose Healthcheck

```yaml
services:
  app:
    image: gofiber-app:latest
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

---

## Database Query Logging

GORM queries otomatis di-log dengan custom logger yang terintegrasi dengan Zap.

### Features

- **Slow Query Detection**: Queries lebih lambat dari threshold dicatat sebagai warning
- **SQL Logging**: Full SQL dengan parameters
- **Execution Time**: Duration setiap query
- **Affected Rows**: Jumlah rows yang terpengaruh
- **Error Logging**: Database errors dengan full context

### Konfigurasi

Threshold slow query berbeda per environment:

**Development:**
- Threshold: 200ms
- Log level: Info (semua queries)

**Production:**
- Threshold: 1 second
- Log level: Warn (hanya slow queries dan errors)

### Log Format

```json
{
  "level": "warn",
  "timestamp": "2025-06-15T10:30:45+07:00",
  "msg": "Slow SQL Query",
  "duration_ms": 1250.5,
  "rows": 100,
  "sql": "SELECT * FROM users WHERE created_at > ? ORDER BY id LIMIT 100",
  "source": "repository/user.go:45"
}
```

### Best Practices

1. **Monitor Slow Queries**: Gunakan log untuk identifikasi queries yang perlu optimization
2. **Add Indexes**: Jika query slow secara konsisten, tambahkan index
3. **Pagination**: Selalu gunakan LIMIT/OFFSET untuk large datasets
4. **Eager Loading**: Gunakan Preload untuk menghindari N+1 queries
5. **Review Logs**: Periksa logs secara berkala untuk pattern

### Query Optimization

Contoh slow query dan optimisasi:

```go
// ❌ Slow: N+1 problem
users := []user.User{}
db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts)
}

// ✅ Fast: Eager loading
users := []user.User{}
db.Preload("Posts").Find(&users)
```

---

## Monitoring Stack Setup

### Complete Monitoring Stack

```yaml
version: '3.8'

services:
  # Application
  app:
    image: gofiber-app:latest
    ports:
      - "3000:3000"
    environment:
      - SENTRY_DSN=${SENTRY_DSN}
      - ENV_TYPE=prod

  # Prometheus
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus

  # Grafana
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana
    depends_on:
      - prometheus

volumes:
  prometheus-data:
  grafana-data:
```

### Alert Rules

Contoh Prometheus alert rules (`alerts.yml`):

```yaml
groups:
  - name: app_alerts
    interval: 30s
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          
      - alert: SlowRequests
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "95th percentile latency > 1s"
          
      - alert: DatabaseConnectionsHigh
        expr: db_connections_in_use > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Database connections usage high"
```

---

## Troubleshooting

### Logs tidak muncul

1. Pastikan `ENV_TYPE` di `.env` sudah benar
2. Check file permissions untuk folder `logs/`
3. Verify logger initialization di `main.go`

### Sentry tidak menerima errors

1. Verify `SENTRY_DSN` di `.env`
2. Check network connectivity ke Sentry
3. Pastikan middleware SentryMiddleware terpasang

### Prometheus metrics tidak update

1. Check `/metrics` endpoint accessible
2. Verify Prometheus scrape config
3. Check middleware PrometheusMiddleware terpasang

### Health check returns unhealthy

1. Check database connection
2. Verify database credentials
3. Check database server status

---

## Performance Tips

1. **Log Level Production**: Gunakan Info/Warn level di production, hindari Debug
2. **Sampling**: Jika traffic tinggi, pertimbangkan sampling untuk Sentry
3. **Metrics Cardinality**: Jangan tambahkan label dengan high cardinality (user_id, request_id)
4. **Log Retention**: Atur retention policy untuk logs (default 7 hari)
5. **Database Connection Pool**: Monitor dan tune sesuai load

---

## Resources

- [Zap Documentation](https://github.com/uber-go/zap)
- [Sentry Go SDK](https://docs.sentry.io/platforms/go/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
