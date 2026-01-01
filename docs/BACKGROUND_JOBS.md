# Background Jobs & Queue System

Sistem Background Jobs menggunakan **Asynq** - Redis-based distributed task queue yang mirip dengan Laravel Queue. Asynq sangat efisien dan production-ready dengan lebih dari 10k+ stars di GitHub.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Setup](#setup)
- [Basic Usage](#basic-usage)
- [Task Types](#task-types)
- [Scheduled Tasks](#scheduled-tasks)
- [Queue Management](#queue-management)
- [API Endpoints](#api-endpoints)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Features

✅ **Redis-based Queue** - Fast and reliable using Redis  
✅ **Multiple Queues** - Critical, Default, Low priority queues  
✅ **Retry Mechanism** - Automatic retry dengan exponential backoff  
✅ **Scheduled Tasks** - Cron-like periodic tasks (mirip Laravel Scheduler)  
✅ **Delayed Jobs** - Execute job di waktu tertentu  
✅ **Concurrency Control** - Parallel job processing  
✅ **Job Monitoring** - Real-time statistics dan monitoring  
✅ **Laravel-Style API** - Familiar dispatch() pattern  
✅ **Graceful Shutdown** - Safe worker termination  

## Architecture

```
┌─────────────┐
│   Client    │ (dispatch job)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Redis      │ (job queue)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Worker    │ (process job)
│   Server    │
└─────────────┘
```

### Components

**1. Asynq Client** - Enqueue jobs to Redis  
**2. Asynq Server** - Worker yang process jobs  
**3. Asynq Scheduler** - Periodic/scheduled tasks  
**4. Asynq Inspector** - Monitoring dan management  

## Setup

### 1. Configuration

Redis harus enabled di `.env`:

```env
# Redis Configuration
REDIS_ENABLE=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 2. Initialization

Asynq otomatis terinisialisasi saat aplikasi start (jika Redis enabled):

```go
// main.go
if config.ENV.REDIS_ENABLE {
    // Initialize Asynq client
    asynqClient := config.InitAsynqClient()
    helper.SetAsynqClient(asynqClient)
    
    // Initialize scheduler
    config.InitAsynqScheduler()
    
    // Start worker server in background
    go startWorkerServer()
    
    // Start scheduler
    go func() {
        config.AsynqScheduler.Run()
    }()
}
```

## Basic Usage

### Dispatch Job (Laravel Style)

```go
import "starter-gofiber/jobs"

// Send email job
err := jobs.SendEmailJob(
    "user@example.com",
    "Welcome!",
    "Welcome to our platform",
    nil, // optional data
)

// Send verification email
err := jobs.SendVerificationEmailJob(
    "user@example.com",
    "John Doe",
    "verification-token-123",
    "https://example.com",
)

// Process export (background task)
err := jobs.ProcessExportJob(
    userID,
    "excel",
    "SELECT * FROM posts",
    "export_2025.xlsx",
)
```

### Delayed Job

```go
import (
    "time"
    "starter-gofiber/jobs"
)

// Send email after 5 minutes
err := jobs.SendDelayedEmailJob(
    "user@example.com",
    "Reminder",
    "Your session is about to expire",
    5 * time.Minute,
)

// Schedule at specific time
processAt := time.Now().Add(24 * time.Hour)
err := jobs.SendEmailAtJob(
    "admin@example.com",
    "Daily Report",
    "Here is your daily report",
    processAt,
)
```

### Queue Priority

Jobs dapat dikirim ke queue dengan priority berbeda:

```go
import "starter-gofiber/helper"

// Critical queue (highest priority)
helper.EnqueueTaskToQueue(
    helper.TaskSendEmail,
    payload,
    helper.QueueCritical,
)

// Default queue (normal priority)
helper.EnqueueTaskToQueue(
    helper.TaskSendEmail,
    payload,
    helper.QueueDefault,
)

// Low queue (lowest priority)
helper.EnqueueTaskToQueue(
    helper.TaskCleanupOldFiles,
    payload,
    helper.QueueLow,
)
```

**Queue Processing Ratio**:
- Critical: 60% (6 workers)
- Default: 30% (3 workers)
- Low: 10% (1 worker)

### Retry Configuration

```go
import (
    "github.com/hibiken/asynq"
    "starter-gofiber/helper"
)

// Custom retry count
err := helper.EnqueueTaskWithRetry(
    helper.TaskSendEmail,
    payload,
    5, // max 5 retries
)

// Custom retry with delay
opts := []asynq.Option{
    asynq.MaxRetry(3),
    asynq.Timeout(2 * time.Minute),
}
err := helper.EnqueueTask(helper.TaskSendEmail, payload, opts...)
```

## Task Types

### Built-in Tasks

| Task Type | Description | Queue |
|-----------|-------------|-------|
| `email:send` | General email sending | default |
| `email:verification` | Email verification | critical |
| `email:password_reset` | Password reset email | critical |
| `export:process` | Data export (CSV/Excel/PDF) | low |
| `cleanup:old_files` | File cleanup | low |
| `report:generate` | Report generation | default |
| `notification:send` | Notification dispatch | default |

### Create Custom Task

**1. Define Task Type**

```go
// helper/jobs.go
const (
    TaskProcessVideo = "video:process"
)
```

**2. Create Handler**

```go
// jobs/handlers.go
type VideoPayload struct {
    VideoID   uint   `json:"video_id"`
    Quality   string `json:"quality"`
    OutputDir string `json:"output_dir"`
}

func HandleProcessVideo(ctx context.Context, t *asynq.Task) error {
    var payload VideoPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return fmt.Errorf("json.Unmarshal failed: %w", err)
    }

    helper.Logger.Info(fmt.Sprintf("Processing video %d with quality %s", payload.VideoID, payload.Quality))

    // TODO: Implement video processing
    // - Fetch video from storage
    // - Transcode to different quality
    // - Save to output directory
    // - Update database

    helper.Logger.Info("Video processing completed")
    return nil
}
```

**3. Create Job Dispatcher**

```go
// jobs/jobs.go
func ProcessVideoJob(videoID uint, quality, outputDir string) error {
    payload := VideoPayload{
        VideoID:   videoID,
        Quality:   quality,
        OutputDir: outputDir,
    }

    return helper.EnqueueTaskToQueue(
        helper.TaskProcessVideo,
        payload,
        helper.QueueLow, // Heavy processing task
    )
}
```

**4. Register Handler**

```go
// main.go - startWorkerServer()
mux.HandleFunc(helper.TaskProcessVideo, jobs.HandleProcessVideo)
```

## Scheduled Tasks

Scheduled tasks berjalan otomatis sesuai cron schedule (mirip Laravel Task Scheduling).

### Registered Periodic Tasks

```go
// jobs/scheduler.go
func RegisterPeriodicTasks(scheduler *asynq.Scheduler) error {
    // Daily cleanup - every day at 2 AM
    scheduler.Register(
        "@daily 2:00",
        asynq.NewTask(helper.TaskCleanupOldFiles, payload),
    )
    
    // Weekly report - every Monday at 8 AM
    scheduler.Register(
        "0 8 * * MON",
        asynq.NewTask(helper.TaskGenerateReport, payload),
    )
    
    // Hourly health check
    scheduler.Register(
        "@every 1h",
        asynq.NewTask("system:health_check", nil),
    )
}
```

### Cron Expression Format

```go
// Predefined constants
const (
    EveryMinute         = "@every 1m"
    EveryFiveMinutes    = "@every 5m"
    EveryTenMinutes     = "@every 10m"
    EveryFifteenMinutes = "@every 15m"
    EveryThirtyMinutes  = "@every 30m"
    Hourly              = "@hourly"       // "0 * * * *"
    Daily               = "@daily"        // "0 0 * * *"
    Weekly              = "0 0 * * 0"     // Sunday
    Monthly             = "0 0 1 * *"     // 1st day
    Yearly              = "0 0 1 1 *"     // Jan 1st
)

// Helper functions
DailyAt(10, 30)              // Daily at 10:30
WeeklyOn(time.Monday, 9, 0)  // Every Monday at 9:00
MonthlyOn(15, 12, 0)         // 15th of month at 12:00
```

### Dynamic Scheduling

```go
import (
    "time"
    "starter-gofiber/jobs"
    "starter-gofiber/config"
)

// Schedule one-time task
executeAt := time.Now().Add(2 * time.Hour)
err := jobs.ScheduleOneTimeTask(
    helper.TaskSendEmail,
    emailPayload,
    executeAt,
)

// Schedule recurring task
err := jobs.ScheduleRecurringTask(
    config.AsynqScheduler,
    "backup:database",
    backupPayload,
    24 * time.Hour, // Every 24 hours
)

// Unschedule task
err := jobs.UnscheduleTask(config.AsynqScheduler, entryID)
```

### Create Custom Periodic Task

```go
// jobs/scheduler.go
func ScheduleDailyBackup(scheduler *asynq.Scheduler, hour, minute int) error {
    cronSpec := fmt.Sprintf("%d %d * * *", minute, hour)
    
    _, err := scheduler.Register(
        cronSpec,
        asynq.NewTask("backup:database", nil),
        asynq.Queue(helper.QueueCritical),
    )
    
    return err
}

// Usage
jobs.ScheduleDailyBackup(config.AsynqScheduler, 3, 0) // Daily at 03:00
```

## Queue Management

### Monitor Queue Statistics

```go
import "starter-gofiber/helper"

// Get all queues
queues, err := helper.ListAllQueues()
for _, q := range queues {
    fmt.Printf("Queue: %s\n", q.Queue)
    fmt.Printf("  Pending: %d\n", q.Pending)
    fmt.Printf("  Active: %d\n", q.Active)
    fmt.Printf("  Scheduled: %d\n", q.Scheduled)
    fmt.Printf("  Retry: %d\n", q.Retry)
    fmt.Printf("  Archived: %d\n", q.Archived)
}

// Get specific queue stats
queueInfo, err := helper.GetQueueStats("critical")
```

### Queue Control

```go
// Pause queue (stop processing)
err := helper.PauseQueue("low")

// Resume queue
err := helper.UnpauseQueue("low")

// Delete all pending tasks
count, err := helper.DeleteAllPendingTasks("low")

// Archive all pending tasks
count, err := helper.ArchiveAllPendingTasks("low")
```

### Task Management

```go
// Get task info
taskInfo, err := helper.GetTaskInfo("default", "task-id-123")

// Retry failed task
err := helper.RetryTask("default", "task-id-123")

// Delete task
err := helper.DeleteTask("default", "task-id-123")
```

## API Endpoints

### Job Queue Management API

```http
GET    /api/jobs/stats                  # Queue statistics (admin)
POST   /api/jobs/trigger/email           # Trigger email job
POST   /api/jobs/trigger/export          # Trigger export job
PUT    /api/jobs/queue/:queue/pause      # Pause queue (admin)
PUT    /api/jobs/queue/:queue/resume     # Resume queue (admin)
DELETE /api/jobs/queue/:queue/pending    # Delete pending tasks (admin)
POST   /api/jobs/retry                   # Retry failed job (admin)
```

### Example Requests

**Get Queue Statistics**

```bash
curl -X GET http://localhost:3000/api/jobs/stats \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{
  "queues": [
    {
      "queue": "critical",
      "pending": 5,
      "active": 2,
      "scheduled": 10,
      "retry": 1,
      "archived": 100,
      "completed": 1500,
      "paused": false,
      "size": 18,
      "latency_ms": 45
    }
  ]
}
```

**Trigger Email Job**

```bash
curl -X POST http://localhost:3000/api/jobs/trigger/email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "to": "user@example.com",
    "subject": "Test Email",
    "body": "This is a test email from background job"
  }'
```

Response:
```json
{
  "message": "Email queued successfully",
  "status": "processing"
}
```

**Trigger Export Job**

```bash
curl -X POST http://localhost:3000/api/jobs/trigger/export \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "export_type": "excel",
    "query": "SELECT * FROM posts WHERE user_id = 1"
  }'
```

Response:
```json
{
  "message": "Export job queued successfully",
  "filename": "export_1_20250101_120000.excel",
  "status": "processing"
}
```

**Pause Queue**

```bash
curl -X PUT http://localhost:3000/api/jobs/queue/low/pause \
  -H "Authorization: Bearer <admin-token>"
```

**Retry Failed Job**

```bash
curl -X POST http://localhost:3000/api/jobs/retry \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin-token>" \
  -d '{
    "queue": "default",
    "task_id": "abc123"
  }'
```

## Best Practices

### 1. Job Idempotency

Pastikan job bisa dijalankan ulang tanpa side effect:

```go
func HandleProcessOrder(ctx context.Context, t *asynq.Task) error {
    var payload OrderPayload
    json.Unmarshal(t.Payload(), &payload)
    
    // Check if already processed
    order, _ := db.GetOrder(payload.OrderID)
    if order.Status == "processed" {
        return nil // Already done, skip
    }
    
    // Process order
    // ...
    
    return nil
}
```

### 2. Error Handling

Return error untuk retry, atau return nil untuk skip:

```go
func HandleSendEmail(ctx context.Context, t *asynq.Task) error {
    // Permanent error - don't retry
    if invalidEmail {
        helper.Logger.Error("Invalid email format")
        return nil // Skip retry
    }
    
    // Temporary error - retry
    if smtpConnectionFailed {
        return fmt.Errorf("SMTP connection failed: %w", err)
    }
    
    return nil
}
```

### 3. Timeout Configuration

Set timeout untuk long-running tasks:

```go
opts := []asynq.Option{
    asynq.Timeout(10 * time.Minute),
    asynq.MaxRetry(3),
}

helper.EnqueueTask(helper.TaskProcessExport, payload, opts...)
```

### 4. Queue Selection

Pilih queue berdasarkan priority:

- **Critical**: Email verification, password reset, payment processing
- **Default**: Notifications, reports, general tasks
- **Low**: Cleanup, data export, analytics

### 5. Monitoring

Monitor job statistics secara berkala:

```go
// Setup metrics collection
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        queues, _ := helper.ListAllQueues()
        for _, q := range queues {
            // Send metrics to monitoring system
            sendMetric("queue.pending", q.Pending, q.Queue)
            sendMetric("queue.retry", q.Retry, q.Queue)
            sendMetric("queue.latency_ms", q.Latency.Milliseconds(), q.Queue)
        }
    }
}()
```

### 6. Graceful Shutdown

Asynq server sudah handle graceful shutdown otomatis:

```go
// main.go
<-quit // Wait for interrupt signal

// Shutdown Asynq server
if config.AsynqServer != nil {
    helper.Info("Shutting down Asynq worker server...")
    config.AsynqServer.Shutdown() // Wait for running tasks to finish
}
```

## Troubleshooting

### Jobs Not Processing

**Problem**: Jobs enqueued tapi tidak diproses

**Solutions**:
1. Check Redis connection
   ```bash
   redis-cli ping
   ```

2. Verify worker server is running
   ```bash
   # Check logs
   grep "Starting Asynq worker server" logs/app.log
   ```

3. Check queue pause status
   ```bash
   curl http://localhost:3000/api/jobs/stats
   ```

### High Retry Rate

**Problem**: Banyak job yang retry

**Solutions**:
1. Check error logs
   ```go
   helper.Logger.Error("Job failed", zap.Error(err))
   ```

2. Increase timeout
   ```go
   asynq.Timeout(5 * time.Minute)
   ```

3. Review handler logic
   - Cek network issues
   - Cek database connection
   - Cek external API availability

### Memory Issues

**Problem**: High memory usage

**Solutions**:
1. Reduce concurrency
   ```go
   config.InitAsynqServer(5) // Reduce from 10 to 5
   ```

2. Add memory limits
   ```go
   // Process jobs in batches
   for i := 0; i < total; i += 100 {
       processBatch(i, i+100)
       time.Sleep(100 * time.Millisecond) // Breathing room
   }
   ```

3. Monitor job size
   - Avoid large payloads (> 1MB)
   - Use references instead of full data

### Dead Letter Queue

**Problem**: Failed jobs setelah max retry

**Check archived tasks**:
```go
inspector := asynq.NewInspector(redisOpt)
tasks, _ := inspector.ListArchivedTasks("default")

for _, task := range tasks {
    fmt.Printf("Task ID: %s\n", task.ID)
    fmt.Printf("Type: %s\n", task.Type)
    fmt.Printf("Last Error: %s\n", task.LastErr)
    fmt.Printf("Retried: %d times\n", task.Retried)
}
```

**Retry archived tasks**:
```bash
curl -X POST http://localhost:3000/api/jobs/retry \
  -d '{"queue":"default","task_id":"<task-id>"}'
```

## Performance Tuning

### Concurrency

Adjust worker concurrency based on workload:

```go
// Light workload
config.InitAsynqServer(5)

// Medium workload
config.InitAsynqServer(10)

// Heavy workload
config.InitAsynqServer(20)
```

### Redis Configuration

Optimize Redis for job queue:

```conf
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
save ""  # Disable RDB snapshots for faster performance
appendonly yes
appendfsync everysec
```

### Queue Priorities

Adjust queue weights based on needs:

```go
cfg := asynq.Config{
    Queues: map[string]int{
        "critical": 8,  // 80%
        "default":  2,  // 20%
        "low":      0,  // 0% (process only when others empty)
    },
}
```

## Summary

✅ **Asynq** - Production-ready Redis queue  
✅ **Laravel-style API** - Familiar dispatch pattern  
✅ **3 Priority Queues** - Critical, Default, Low  
✅ **Automatic Retry** - Exponential backoff  
✅ **Cron Scheduling** - Periodic tasks  
✅ **Monitoring** - Real-time stats  
✅ **Graceful Shutdown** - Safe termination  

**Next Steps**:
1. Implement email service integration (SMTP/SendGrid)
2. Add job dashboard UI (Asynqmon)
3. Setup alerting for failed jobs
4. Implement job chaining for complex workflows
