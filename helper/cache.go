package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisClient is the global Redis client instance
var RedisClient *redis.Client

// InitRedisClient sets the global Redis client
func InitRedisClient(client *redis.Client) {
	RedisClient = client
}

// CloseRedis closes Redis connection gracefully
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// CacheGet retrieves data from cache
func CacheGet(key string, dest interface{}) error {
	if RedisClient == nil {
		return redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// CacheSet stores data in cache with expiration
func CacheSet(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, data, expiration).Err()
}

// CacheDelete removes data from cache
func CacheDelete(keys ...string) error {
	if RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return RedisClient.Del(ctx, keys...).Err()
}

// CacheExists checks if key exists in cache
func CacheExists(key string) bool {
	if RedisClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	count, err := RedisClient.Exists(ctx, key).Result()
	return err == nil && count > 0
}

// CacheGetOrSet retrieves from cache or executes fn and caches the result
func CacheGetOrSet(key string, dest interface{}, expiration time.Duration, fn func() (interface{}, error)) error {
	// Try to get from cache first
	err := CacheGet(key, dest)
	if err == nil {
		return nil // Cache hit
	}

	// Cache miss - execute function
	result, err := fn()
	if err != nil {
		return err
	}

	// Store in cache
	if err := CacheSet(key, result, expiration); err != nil {
		// Log error but don't fail
		logger, _ := zap.NewProduction()
		logger.Error("Failed to cache result", zap.Error(err), zap.String("key", key))
	}

	// Marshal result to dest
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// CacheInvalidateByPattern removes all keys matching pattern
func CacheInvalidateByPattern(pattern string) error {
	if RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

// CacheIncrement increments a counter in cache
func CacheIncrement(key string, expiration time.Duration) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("redis not enabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	count, err := RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration if this is the first increment
	if count == 1 && expiration > 0 {
		RedisClient.Expire(ctx, key, expiration)
	}

	return count, nil
}

// CacheDecrement decrements a counter in cache
func CacheDecrement(key string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("redis not enabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return RedisClient.Decr(ctx, key).Result()
}

// CacheGetTTL returns remaining time to live for a key
func CacheGetTTL(key string) (time.Duration, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("redis not enabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return RedisClient.TTL(ctx, key).Result()
}

// CacheExpire sets expiration time for existing key
func CacheExpire(key string, expiration time.Duration) error {
	if RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return RedisClient.Expire(ctx, key, expiration).Err()
}

// CacheLock acquires a distributed lock (useful for preventing cache stampede)
func CacheLock(key string, expiration time.Duration) (bool, error) {
	if RedisClient == nil {
		return false, fmt.Errorf("redis not enabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	lockKey := "lock:" + key
	return RedisClient.SetNX(ctx, lockKey, "1", expiration).Result()
}

// CacheUnlock releases a distributed lock
func CacheUnlock(key string) error {
	if RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	lockKey := "lock:" + key
	return RedisClient.Del(ctx, lockKey).Err()
}

// CacheFlush removes all cache entries (use with caution!)
func CacheFlush() error {
	if RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return RedisClient.FlushDB(ctx).Err()
}

// CacheKeys returns all keys matching pattern
func CacheKeys(pattern string) ([]string, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("redis not enabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return RedisClient.Keys(ctx, pattern).Result()
}

// CacheStats returns cache statistics
func CacheStats() map[string]interface{} {
	stats := make(map[string]interface{})

	if RedisClient == nil {
		stats["enabled"] = false
		return stats
	}

	stats["enabled"] = true

	// Get pool stats
	if poolStats := RedisClient.PoolStats(); poolStats != nil {
		stats["hits"] = poolStats.Hits
		stats["misses"] = poolStats.Misses
		stats["timeouts"] = poolStats.Timeouts
		stats["total_conns"] = poolStats.TotalConns
		stats["idle_conns"] = poolStats.IdleConns
		stats["stale_conns"] = poolStats.StaleConns
	}

	// Get database size
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if size, err := RedisClient.DBSize(ctx).Result(); err == nil {
		stats["db_size"] = size
	}

	return stats
}
