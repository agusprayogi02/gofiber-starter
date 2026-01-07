package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"starter-gofiber/internal/infrastructure/cache"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// CacheConfig defines cache middleware configuration
type CacheConfig struct {
	// Expiration time for cached responses
	Expiration time.Duration

	// Skip caching for these paths
	ExcludePaths []string

	// Skip caching for these methods (default: only GET is cached)
	ExcludeMethods []string

	// Cache key prefix
	KeyPrefix string

	// Skip cache if response status is in this list
	ExcludeStatuses []int
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		Expiration:      5 * time.Minute,
		ExcludePaths:    []string{"/health", "/metrics", "/admin", "/api/auth/*"},
		ExcludeMethods:  []string{"POST", "PUT", "PATCH", "DELETE"},
		KeyPrefix:       "cache:",
		ExcludeStatuses: []int{500, 502, 503, 504},
	}
}

// CacheMiddleware creates a new cache middleware with custom config
func CacheMiddleware(cfg ...CacheConfig) fiber.Handler {
	// Use default config if none provided
	var cacheConfig CacheConfig
	if len(cfg) > 0 {
		cacheConfig = cfg[0]
	} else {
		cacheConfig = DefaultCacheConfig()
	}

	return func(c *fiber.Ctx) error {
		// Skip if Redis is not enabled
		if cache.RedisClient == nil {
			return c.Next()
		}

		// Skip non-GET requests
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		// Skip excluded paths
		for _, path := range cacheConfig.ExcludePaths {
			if c.Path() == path {
				return c.Next()
			}
		}

		// Generate cache key from request
		cacheKey := generateCacheKey(c, cacheConfig.KeyPrefix)

		// Try to get from cache
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		cached, err := cache.RedisClient.Get(ctx, cacheKey).Bytes()
		if err == nil && len(cached) > 0 {
			// Cache hit - return cached response
			c.Set("X-Cache", "HIT")
			c.Set("Content-Type", "application/json")
			return c.Send(cached)
		}

		// Cache miss - capture response
		c.Set("X-Cache", "MISS")

		// Create a custom response writer to capture body
		originalBody := c.Response().Body()
		buf := new(bytes.Buffer)

		// Execute handler
		if err := c.Next(); err != nil {
			return err
		}

		// Get response status and body
		status := c.Response().StatusCode()
		responseBody := c.Response().Body()

		// Check if we should cache this response
		shouldCache := true
		for _, excludeStatus := range cacheConfig.ExcludeStatuses {
			if status == excludeStatus {
				shouldCache = false
				break
			}
		}

		// Cache the response if applicable
		if shouldCache && status >= 200 && status < 300 {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			if err := cache.RedisClient.Set(ctx, cacheKey, responseBody, cacheConfig.Expiration).Err(); err != nil {
				// Log error but don't fail the request
				logger, _ := zap.NewProduction()
				logger.Error("Failed to cache response",
					zap.Error(err),
					zap.String("key", cacheKey),
				)
			}
		}

		// Restore original body
		c.Response().SetBody(originalBody)
		buf.Write(responseBody)

		return nil
	}
}

// SimpleCacheMiddleware creates a cache middleware with 5 minute expiration
func SimpleCacheMiddleware() fiber.Handler {
	return CacheMiddleware(DefaultCacheConfig())
}

// generateCacheKey creates a unique cache key from request
func generateCacheKey(c *fiber.Ctx, prefix string) string {
	// Include method, path, query params, and user ID (if authenticated)
	keyData := fmt.Sprintf("%s:%s:%s",
		c.Method(),
		c.Path(),
		string(c.Request().URI().QueryString()),
	)

	// Include user ID if authenticated
	if userID := c.Locals("user_id"); userID != nil {
		keyData += fmt.Sprintf(":user:%v", userID)
	}

	// Hash the key data
	hash := sha256.Sum256([]byte(keyData))
	return prefix + hex.EncodeToString(hash[:])
}

// InvalidateCache removes a specific cache entry
func InvalidateCache(c *fiber.Ctx, paths ...string) error {
	if cache.RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for _, path := range paths {
		// Generate cache key pattern
		pattern := fmt.Sprintf("cache:*%s*", path)

		// Find matching keys
		keys, err := cache.RedisClient.Keys(ctx, pattern).Result()
		if err != nil {
			return err
		}

		// Delete all matching keys
		if len(keys) > 0 {
			if err := cache.RedisClient.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
	}

	return nil
}

// InvalidateCacheByPattern removes cache entries matching a pattern
func InvalidateCacheByPattern(pattern string) error {
	if cache.RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	keys, err := cache.RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return cache.RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

// ClearAllCache removes all cached responses
func ClearAllCache() error {
	if cache.RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Delete all keys with cache prefix
	keys, err := cache.RedisClient.Keys(ctx, "cache:*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return cache.RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}
