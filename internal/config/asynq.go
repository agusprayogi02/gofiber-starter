package config

import (
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

var (
	// AsynqClient for enqueueing tasks
	AsynqClient *asynq.Client

	// AsynqServer for processing tasks (worker)
	AsynqServer *asynq.Server

	// AsynqScheduler for periodic tasks
	AsynqScheduler *asynq.Scheduler
)

// InitAsynqClient initializes Asynq client for enqueueing tasks
func InitAsynqClient() *asynq.Client {
	redisOpt := asynq.RedisClientOpt{
		Addr:     ENV.REDIS_HOST + ":" + ENV.REDIS_PORT,
		Password: ENV.REDIS_PASSWORD,
		DB:       ENV.REDIS_DB,
	}

	AsynqClient = asynq.NewClient(redisOpt)
	return AsynqClient
}

// InitAsynqServer initializes Asynq server for processing tasks
func InitAsynqServer(concurrency int) *asynq.Server {
	redisOpt := asynq.RedisClientOpt{
		Addr:     ENV.REDIS_HOST + ":" + ENV.REDIS_PORT,
		Password: ENV.REDIS_PASSWORD,
		DB:       ENV.REDIS_DB,
	}

	// Server config
	cfg := asynq.Config{
		Concurrency: concurrency, // Number of concurrent workers
		Queues: map[string]int{
			"critical": 6, // High priority (60%)
			"default":  3, // Medium priority (30%)
			"low":      1, // Low priority (10%)
		},
		// Retry configuration
		RetryDelayFunc: asynq.DefaultRetryDelayFunc,
		IsFailure: func(err error) bool {
			// Custom logic to determine if error needs retry
			return true
		},
	}

	AsynqServer = asynq.NewServer(redisOpt, cfg)
	return AsynqServer
}

// InitAsynqScheduler initializes Asynq scheduler for periodic tasks
func InitAsynqScheduler() *asynq.Scheduler {
	redisOpt := asynq.RedisClientOpt{
		Addr:     ENV.REDIS_HOST + ":" + ENV.REDIS_PORT,
		Password: ENV.REDIS_PASSWORD,
		DB:       ENV.REDIS_DB,
	}

	location := GetTimezone()
	AsynqScheduler = asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{
		Location: location,
	})

	return AsynqScheduler
}

// CloseAsynq closes Asynq connections
func CloseAsynq() {
	if AsynqClient != nil {
		AsynqClient.Close()
	}
	if AsynqScheduler != nil {
		AsynqScheduler.Shutdown()
	}
	// Server shutdown handled separately in main.go
}

// GetAsynqRedisClient returns redis client for Asynq inspector
func GetAsynqRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     ENV.REDIS_HOST + ":" + ENV.REDIS_PORT,
		Password: ENV.REDIS_PASSWORD,
		DB:       ENV.REDIS_DB,
	})
}
