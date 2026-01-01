# Background Jobs - Quick Reference

## Quick Start

### 1. Basic Job Dispatch

```go
import "starter-gofiber/jobs"

// Send email
jobs.SendEmailJob("user@example.com", "Subject", "Body", nil)

// Process export
jobs.ProcessExportJob(userID, "excel", "query", "file.xlsx")

// Send notification
jobs.SendNotificationJob(userID, "email", "Title", "Message", nil)
```

### 2. Delayed Jobs

```go
import "time"

// 5 minutes delay
jobs.SendDelayedEmailJob(email, subject, body, 5*time.Minute)

// Specific time
jobs.SendEmailAtJob(email, subject, body, time.Now().Add(24*time.Hour))
```

### 3. Queue Priority

```go
import "starter-gofiber/helper"

// Critical (60% processing power)
helper.EnqueueTaskToQueue(helper.TaskSendEmail, payload, helper.QueueCritical)

// Default (30% processing power)
helper.EnqueueTaskToQueue(helper.TaskSendEmail, payload, helper.QueueDefault)

// Low (10% processing power)
helper.EnqueueTaskToQueue(helper.TaskCleanupOldFiles, payload, helper.QueueLow)
```

## Task Types

| Constant | Description | Default Queue |
|----------|-------------|---------------|
| `TaskSendEmail` | General email | default |
| `TaskSendVerificationCode` | Verification email | critical |
| `TaskSendPasswordReset` | Password reset | critical |
| `TaskProcessExport` | Data export | low |
| `TaskCleanupOldFiles` | File cleanup | low |
| `TaskGenerateReport` | Report generation | default |
| `TaskSendNotification` | Notifications | default |

## Scheduled Tasks (Cron)

```go
// Predefined schedules
EveryMinute         = "@every 1m"
EveryFiveMinutes    = "@every 5m"
Hourly              = "@hourly"
Daily               = "@daily"
Weekly              = "0 0 * * 0"
Monthly             = "0 0 1 * *"

// Custom schedules
DailyAt(10, 30)              // Daily at 10:30
WeeklyOn(time.Monday, 9, 0)  // Every Monday at 9:00
MonthlyOn(15, 12, 0)         // 15th of month at 12:00
```

## Job Management

```go
// Get queue stats
queues, _ := helper.ListAllQueues()
queueInfo, _ := helper.GetQueueStats("default")

// Pause/Resume queue
helper.PauseQueue("low")
helper.UnpauseQueue("low")

// Task management
helper.GetTaskInfo("default", "task-id")
helper.RetryTask("default", "task-id")
helper.DeleteTask("default", "task-id")

// Bulk operations
helper.DeleteAllPendingTasks("low")
helper.ArchiveAllPendingTasks("low")
```

## API Endpoints

```
GET    /api/jobs/stats                  # Queue stats (admin)
POST   /api/jobs/trigger/email           # Trigger email
POST   /api/jobs/trigger/export          # Trigger export
PUT    /api/jobs/queue/:queue/pause      # Pause queue (admin)
PUT    /api/jobs/queue/:queue/resume     # Resume queue (admin)
DELETE /api/jobs/queue/:queue/pending    # Delete pending (admin)
POST   /api/jobs/retry                   # Retry failed (admin)
```

## cURL Examples

```bash
# Get stats
curl http://localhost:3000/api/jobs/stats \
  -H "Authorization: Bearer <token>"

# Trigger email
curl -X POST http://localhost:3000/api/jobs/trigger/email \
  -H "Content-Type: application/json" \
  -d '{"to":"user@example.com","subject":"Test","body":"Body"}'

# Pause queue
curl -X PUT http://localhost:3000/api/jobs/queue/low/pause \
  -H "Authorization: Bearer <admin-token>"

# Retry failed job
curl -X POST http://localhost:3000/api/jobs/retry \
  -H "Content-Type: application/json" \
  -d '{"queue":"default","task_id":"abc123"}'
```

## Create Custom Job

**1. Define constant** (`helper/jobs.go`):
```go
const TaskProcessVideo = "video:process"
```

**2. Create payload & handler** (`jobs/handlers.go`):
```go
type VideoPayload struct {
    VideoID uint `json:"video_id"`
}

func HandleProcessVideo(ctx context.Context, t *asynq.Task) error {
    var payload VideoPayload
    json.Unmarshal(t.Payload(), &payload)
    // Process video...
    return nil
}
```

**3. Create dispatcher** (`jobs/jobs.go`):
```go
func ProcessVideoJob(videoID uint) error {
    return helper.EnqueueTask(helper.TaskProcessVideo, VideoPayload{VideoID: videoID})
}
```

**4. Register** (`main.go`):
```go
mux.HandleFunc(helper.TaskProcessVideo, jobs.HandleProcessVideo)
```

## Error Handling

```go
func HandleSendEmail(ctx context.Context, t *asynq.Task) error {
    // Permanent error - skip retry
    if invalidEmail {
        return nil
    }
    
    // Temporary error - retry
    if networkError {
        return fmt.Errorf("network error: %w", err)
    }
    
    return nil
}
```

## Configuration

```env
# .env
REDIS_ENABLE=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Monitoring

```go
// Periodic stats collection
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        queues, _ := helper.ListAllQueues()
        for _, q := range queues {
            log.Printf("%s: pending=%d active=%d retry=%d", 
                q.Queue, q.Pending, q.Active, q.Retry)
        }
    }
}()
```

## Best Practices

✅ Make jobs idempotent (safe to run multiple times)  
✅ Return error to retry, nil to skip  
✅ Set appropriate timeout for long tasks  
✅ Use correct queue based on priority  
✅ Monitor queue stats regularly  
✅ Handle graceful shutdown properly  

## Common Issues

**Jobs not processing?**
- Check Redis connection
- Verify worker server is running
- Check if queue is paused

**High retry rate?**
- Check error logs
- Increase timeout
- Review handler logic

**Memory issues?**
- Reduce concurrency
- Add memory limits
- Avoid large payloads

For detailed documentation, see [BACKGROUND_JOBS.md](BACKGROUND_JOBS.md)
