package helper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

var (
	// Global Asynq client - initialized from main
	AsynqClientInstance *asynq.Client
	// Redis connection details for inspector
	RedisAddr     string
	RedisPassword string
	RedisDB       int
)

// SetAsynqClient sets global Asynq client instance
func SetAsynqClient(client *asynq.Client) {
	AsynqClientInstance = client
}

// SetRedisConfig sets Redis config for inspector
func SetRedisConfig(addr, password string, db int) {
	RedisAddr = addr
	RedisPassword = password
	RedisDB = db
}

// Task type constants - similar to Laravel job names
const (
	TaskSendEmail            = "email:send"
	TaskSendVerificationCode = "email:verification"
	TaskSendPasswordReset    = "email:password_reset"
	TaskProcessExport        = "export:process"
	TaskCleanupOldFiles      = "cleanup:old_files"
	TaskGenerateReport       = "report:generate"
	TaskSendNotification     = "notification:send"
)

// Queue names - similar to Laravel queue names
const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

// EnqueueTask adds task to queue (similar to dispatch() in Laravel)
func EnqueueTask(taskType string, payload interface{}, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(taskType, data, opts...)
	info, err := AsynqClientInstance.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	Logger.Info(fmt.Sprintf("Task enqueued: ID=%s Queue=%s Type=%s", info.ID, info.Queue, info.Type))
	return nil
}

// EnqueueTaskWithDelay adds task with delay (similar to dispatch()->delay() in Laravel)
func EnqueueTaskWithDelay(taskType string, payload interface{}, delay time.Duration, opts ...asynq.Option) error {
	opts = append(opts, asynq.ProcessIn(delay))
	return EnqueueTask(taskType, payload, opts...)
}

// EnqueueTaskAt adds task to run at specific time
func EnqueueTaskAt(taskType string, payload interface{}, processAt time.Time, opts ...asynq.Option) error {
	opts = append(opts, asynq.ProcessAt(processAt))
	return EnqueueTask(taskType, payload, opts...)
}

// EnqueueTaskWithRetry adds task with custom retry count
func EnqueueTaskWithRetry(taskType string, payload interface{}, maxRetry int, opts ...asynq.Option) error {
	opts = append(opts, asynq.MaxRetry(maxRetry))
	return EnqueueTask(taskType, payload, opts...)
}

// EnqueueTaskToQueue adds task to specific queue with priority
func EnqueueTaskToQueue(taskType string, payload interface{}, queue string, opts ...asynq.Option) error {
	opts = append(opts, asynq.Queue(queue))
	return EnqueueTask(taskType, payload, opts...)
}

// DispatchNow alias untuk EnqueueTask (Laravel style naming)
func DispatchNow(taskType string, payload interface{}, opts ...asynq.Option) error {
	return EnqueueTask(taskType, payload, opts...)
}

// DispatchLater alias untuk EnqueueTaskWithDelay (Laravel style naming)
func DispatchLater(taskType string, payload interface{}, delay time.Duration, opts ...asynq.Option) error {
	return EnqueueTaskWithDelay(taskType, payload, delay, opts...)
}

// GetTaskInfo gets task information by ID
func GetTaskInfo(queue string, taskID string) (*asynq.TaskInfo, error) {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
	defer inspector.Close()

	return inspector.GetTaskInfo(queue, taskID)
}

// DeleteTask deletes task from queue
func DeleteTask(queue string, taskID string) error {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
	defer inspector.Close()

	return inspector.DeleteTask(queue, taskID)
}

// RetryTask retries failed task
func RetryTask(queue string, taskID string) error {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
	defer inspector.Close()

	return inspector.RunTask(queue, taskID)
}

// GetQueueStats gets queue statistics
func GetQueueStats(queue string) (*asynq.QueueInfo, error) {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
	defer inspector.Close()

	return inspector.GetQueueInfo(queue)
}

// ListAllQueues gets all queues
func ListAllQueues() ([]*asynq.QueueInfo, error) {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
	defer inspector.Close()

	// Get all queue names first
	queueNames, err := inspector.Queues()
	if err != nil {
		return nil, err
	}

	// Get info for each queue
	var queues []*asynq.QueueInfo
	for _, name := range queueNames {
		info, err := inspector.GetQueueInfo(name)
		if err != nil {
			continue
		}
		queues = append(queues, info)
	}

	return queues, nil
}

// PauseQueue pause queue processing
func PauseQueue(queue string) error {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
	defer inspector.Close()

	return inspector.PauseQueue(queue)
}

// UnpauseQueue resume queue processing
func UnpauseQueue(queue string) error {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
	defer inspector.Close()

	return inspector.UnpauseQueue(queue)
}

// DeleteAllPendingTasks deletes all pending tasks in queue
func DeleteAllPendingTasks(queue string) (int, error) {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
	defer inspector.Close()

	return inspector.DeleteAllPendingTasks(queue)
}

// ArchiveAllPendingTasks archives all pending tasks
func ArchiveAllPendingTasks(queue string) (int, error) {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
	defer inspector.Close()

	return inspector.ArchiveAllPendingTasks(queue)
}
