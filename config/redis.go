package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

// InitRedis initializes Redis connection
func InitRedis() (*redis.Client, error) {
	if !ENV.REDIS_ENABLE {
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", ENV.REDIS_HOST, ENV.REDIS_PORT),
		Password:     ENV.REDIS_PASSWORD,
		DB:           ENV.REDIS_DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Log success
	logger, _ := zap.NewProduction()
	logger.Info("Redis connected successfully",
		zap.String("host", ENV.REDIS_HOST),
		zap.String("port", ENV.REDIS_PORT),
		zap.Int("db", ENV.REDIS_DB),
	)

	return client, nil
}

// GetRedisStats returns Redis connection pool statistics
func GetRedisStats(client *redis.Client) *redis.PoolStats {
	if client != nil {
		return client.PoolStats()
	}
	return nil
}
