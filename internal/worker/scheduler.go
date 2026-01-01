package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"starter-gofiber/pkg/logger"

	"github.com/hibiken/asynq"
)

// RegisterPeriodicTasks registers scheduled tasks (similar to Laravel scheduler)
func RegisterPeriodicTasks(scheduler *asynq.Scheduler) error {
	// Daily cleanup - every day at 2 AM
	_, err := scheduler.Register(
		"@daily 2:00",
		asynq.NewTask(TaskCleanupOldFiles, []byte(`{"directory":"public/uploads","older_than":30}`)),
		asynq.Queue(QueueLow),
	)
	if err != nil {
		return fmt.Errorf("failed to register daily cleanup: %w", err)
	}

	// Weekly report - every Monday at 8 AM
	_, err = scheduler.Register(
		"0 8 * * MON",
		asynq.NewTask(TaskGenerateReport, []byte(`{"report_type":"weekly","filters":{}}`)),
		asynq.Queue(QueueDefault),
	)
	if err != nil {
		return fmt.Errorf("failed to register weekly report: %w", err)
	}

	// Hourly health check
	_, err = scheduler.Register(
		"@every 1h",
		asynq.NewTask("system:health_check", nil),
		asynq.Queue(QueueCritical),
	)
	if err != nil {
		return fmt.Errorf("failed to register health check: %w", err)
	}

	// Delete expired tokens - every 6 hours
	_, err = scheduler.Register(
		"@every 6h",
		asynq.NewTask("cleanup:expired_tokens", nil),
		asynq.Queue(QueueDefault),
	)
	if err != nil {
		return fmt.Errorf("failed to register token cleanup: %w", err)
	}

	// Monthly archive - first day of month at 3 AM
	_, err = scheduler.Register(
		"0 3 1 * *",
		asynq.NewTask("archive:monthly", nil),
		asynq.Queue(QueueLow),
	)
	if err != nil {
		return fmt.Errorf("failed to register monthly archive: %w", err)
	}

	logger.Info("All periodic tasks registered successfully")
	return nil
}

// Handlers for scheduled tasks

// HandleHealthCheck periodic health check
func HandleHealthCheck(ctx context.Context, t *asynq.Task) error {
	logger.Info("Running scheduled health check")

	// TODO: Implement health check logic
	// - Check database connection
	// - Check Redis connection
	// - Check external services
	// - Send alert if something is down

	logger.Info("Health check completed")
	return nil
}

// HandleCleanupExpiredTokens cleanup expired refresh tokens
func HandleCleanupExpiredTokens(ctx context.Context, t *asynq.Task) error {
	logger.Info("Running scheduled token cleanup")

	// TODO: Implement cleanup logic
	// - Delete expired refresh tokens
	// - Delete expired password reset tokens
	// - Delete expired verification tokens

	logger.Info("Token cleanup completed")
	return nil
}

// HandleMonthlyArchive archive old data
func HandleMonthlyArchive(ctx context.Context, t *asynq.Task) error {
	logger.Info("Running monthly archive")

	// TODO: Implement archive logic
	// - Archive old logs
	// - Archive old audit trails
	// - Compress archived data
	// - Move to cold storage

	logger.Info("Monthly archive completed")
	return nil
}

// Custom periodic task examples

// ScheduleDailyBackup schedule daily database backup
func ScheduleDailyBackup(scheduler *asynq.Scheduler, hour, minute int) error {
	cronSpec := fmt.Sprintf("%d %d * * *", minute, hour)

	_, err := scheduler.Register(
		cronSpec,
		asynq.NewTask("backup:database", nil),
		asynq.Queue(QueueCritical),
	)
	if err != nil {
		return fmt.Errorf("failed to schedule daily backup: %w", err)
	}

	logger.Info(fmt.Sprintf("Daily backup scheduled at %02d:%02d", hour, minute))
	return nil
}

// HandleDatabaseBackup handler for database backup
func HandleDatabaseBackup(ctx context.Context, t *asynq.Task) error {
	logger.Info("Starting database backup")

	// TODO: Implement backup logic
	// - Dump database
	// - Compress backup file
	// - Upload to S3 or storage
	// - Delete old backups

	logger.Info("Database backup completed")
	return nil
}

// ScheduleEmailDigest schedule daily email digest
func ScheduleEmailDigest(scheduler *asynq.Scheduler, hour, minute int, recipients []string) error {
	cronSpec := fmt.Sprintf("%d %d * * *", minute, hour)

	payload := map[string]interface{}{
		"recipients": recipients,
		"type":       "daily_digest",
	}

	_, err := scheduler.Register(
		cronSpec,
		asynq.NewTask("email:daily_digest", toJSON(payload)),
		asynq.Queue(QueueDefault),
	)
	if err != nil {
		return fmt.Errorf("failed to schedule email digest: %w", err)
	}

	return nil
}

// HandleDailyEmailDigest handler for email digest
func HandleDailyEmailDigest(ctx context.Context, t *asynq.Task) error {
	logger.Info("Sending daily email digest")

	// TODO: Implement digest logic
	// - Collect daily statistics
	// - Generate digest content
	// - Send to all recipients

	logger.Info("Daily email digest sent")
	return nil
}

// ScheduleMetricsCollection collect metrics every 5 minutes
func ScheduleMetricsCollection(scheduler *asynq.Scheduler) error {
	_, err := scheduler.Register(
		"@every 5m",
		asynq.NewTask("metrics:collect", nil),
		asynq.Queue(QueueLow),
	)
	if err != nil {
		return fmt.Errorf("failed to schedule metrics collection: %w", err)
	}

	return nil
}

// HandleMetricsCollection handler for metrics collection
func HandleMetricsCollection(ctx context.Context, t *asynq.Task) error {
	logger.Info("Collecting metrics")

	// TODO: Implement metrics collection
	// - Collect system metrics
	// - Store to time-series database
	// - Update dashboards

	logger.Info("Metrics collection completed")
	return nil
}

// toJSON converts payload to JSON bytes
func toJSON(payload interface{}) []byte {
	data, err := json.Marshal(payload)
	if err != nil {
		return []byte("{}")
	}
	return data
}

// Dynamic scheduling examples

// ScheduleOneTimeTask schedule one-time task at specific time
func ScheduleOneTimeTask(taskType string, payload interface{}, executeAt time.Time) error {
	return EnqueueTaskAt(taskType, payload, executeAt)
}

// ScheduleRecurringTask create recurring task with custom interval
func ScheduleRecurringTask(scheduler *asynq.Scheduler, taskType string, payload interface{}, interval time.Duration) error {
	cronSpec := fmt.Sprintf("@every %s", interval.String())

	_, err := scheduler.Register(
		cronSpec,
		asynq.NewTask(taskType, toJSON(payload)),
	)

	return err
}

// UnscheduleTask remove scheduled task
func UnscheduleTask(scheduler *asynq.Scheduler, entryID string) error {
	return scheduler.Unregister(entryID)
}

// ListScheduledTasks get all scheduled tasks (cron entries)
func ListScheduledTasks() []string {
	// Note: Asynq doesn't expose scheduler entries directly
	// This would need custom implementation to track scheduled tasks
	logger.Info("Listing scheduled tasks not yet implemented")
	return []string{}
}

// Cron expression helpers (Laravel-style)

// EveryMinute = "@every 1m"
// EveryFiveMinutes = "@every 5m"
// EveryTenMinutes = "@every 10m"
// EveryFifteenMinutes = "@every 15m"
// EveryThirtyMinutes = "@every 30m"
// Hourly = "@hourly" atau "0 * * * *"
// Daily = "@daily" atau "0 0 * * *"
// DailyAt(hour, minute) = "minute hour * * *"
// Weekly = "0 0 * * 0"
// Monthly = "0 0 1 * *"
// Yearly = "0 0 1 1 *"

const (
	EveryMinute         = "@every 1m"
	EveryFiveMinutes    = "@every 5m"
	EveryTenMinutes     = "@every 10m"
	EveryFifteenMinutes = "@every 15m"
	EveryThirtyMinutes  = "@every 30m"
	Hourly              = "@hourly"
	Daily               = "@daily"
	Weekly              = "0 0 * * 0"
	Monthly             = "0 0 1 * *"
	Yearly              = "0 0 1 1 *"
)

// DailyAt creates cron expression for daily at specific time
func DailyAt(hour, minute int) string {
	return fmt.Sprintf("%d %d * * *", minute, hour)
}

// WeeklyOn creates cron expression for weekly on specific day and time
func WeeklyOn(day time.Weekday, hour, minute int) string {
	return fmt.Sprintf("%d %d * * %d", minute, hour, day)
}

// MonthlyOn creates cron expression for monthly on specific day and time
func MonthlyOn(day, hour, minute int) string {
	return fmt.Sprintf("%d %d %d * *", minute, hour, day)
}
