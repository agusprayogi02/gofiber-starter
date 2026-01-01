package cache

import (
	"fmt"
	"time"
)

// CacheStrategy defines different cache invalidation patterns
type CacheStrategy string

const (
	// StrategyTTL - Time-To-Live based invalidation (automatic)
	StrategyTTL CacheStrategy = "ttl"

	// StrategyWriteThrough - Invalidate on write operations
	StrategyWriteThrough CacheStrategy = "write_through"

	// StrategyWriteBack - Delay invalidation (eventual consistency)
	StrategyWriteBack CacheStrategy = "write_back"

	// StrategyLRU - Least Recently Used (handled by Redis)
	StrategyLRU CacheStrategy = "lru"
)

// CacheKeyBuilder helps build consistent cache keys
type CacheKeyBuilder struct {
	prefix string
	parts  []string
}

// NewCacheKeyBuilder creates a new cache key builder
func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{
		prefix: prefix,
		parts:  make([]string, 0),
	}
}

// Add adds a part to the cache key
func (b *CacheKeyBuilder) Add(part string) *CacheKeyBuilder {
	b.parts = append(b.parts, part)
	return b
}

// AddInt adds an integer part to the cache key
func (b *CacheKeyBuilder) AddInt(num int) *CacheKeyBuilder {
	b.parts = append(b.parts, fmt.Sprintf("%d", num))
	return b
}

// AddUint adds an uint part to the cache key
func (b *CacheKeyBuilder) AddUint(num uint) *CacheKeyBuilder {
	b.parts = append(b.parts, fmt.Sprintf("%d", num))
	return b
}

// Build constructs the final cache key
func (b *CacheKeyBuilder) Build() string {
	key := b.prefix
	for _, part := range b.parts {
		key += ":" + part
	}
	return key
}

// InvalidateRelated invalidates all cache keys related to a resource
func InvalidateRelated(resourceType string, resourceID interface{}) error {
	// Build pattern to match all related keys
	pattern := fmt.Sprintf("*:%s:%v:*", resourceType, resourceID)
	return CacheInvalidateByPattern(pattern)
}

// InvalidateCollection invalidates cache for a collection/list
func InvalidateCollection(collectionName string) error {
	pattern := fmt.Sprintf("*:%s:*", collectionName)
	return CacheInvalidateByPattern(pattern)
}

// InvalidateUser invalidates all cache for a specific user
func InvalidateUser(userID uint) error {
	pattern := fmt.Sprintf("*:user:%d:*", userID)
	return CacheInvalidateByPattern(pattern)
}

// WriteThrough performs write operation and invalidates cache
func WriteThrough(key string, writeFn func() error) error {
	// Execute write operation
	if err := writeFn(); err != nil {
		return err
	}

	// Invalidate cache after successful write
	return CacheDelete(key)
}

// WriteThroughPattern performs write and invalidates cache by pattern
func WriteThroughPattern(pattern string, writeFn func() error) error {
	// Execute write operation
	if err := writeFn(); err != nil {
		return err
	}

	// Invalidate cache after successful write
	return CacheInvalidateByPattern(pattern)
}

// CacheAside implements cache-aside (lazy loading) pattern
func CacheAside(key string, dest interface{}, expiration time.Duration, loadFn func() (interface{}, error)) error {
	return CacheGetOrSet(key, dest, expiration, loadFn)
}

// RefreshCache refreshes cache by deleting and reloading
func RefreshCache(key string, dest interface{}, expiration time.Duration, loadFn func() (interface{}, error)) error {
	// Delete existing cache
	if err := CacheDelete(key); err != nil {
		return err
	}

	// Reload from source
	return CacheGetOrSet(key, dest, expiration, loadFn)
}

// BulkInvalidate invalidates multiple cache keys
func BulkInvalidate(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return CacheDelete(keys...)
}

// TaggedCacheKey creates a cache key with tags for grouped invalidation
func TaggedCacheKey(key string, tags ...string) string {
	for _, tag := range tags {
		key = tag + ":" + key
	}
	return key
}

// InvalidateByTag invalidates all cache keys with a specific tag
func InvalidateByTag(tag string) error {
	pattern := fmt.Sprintf("%s:*", tag)
	return CacheInvalidateByPattern(pattern)
}

// CacheWarming pre-loads cache with frequently accessed data
func CacheWarming(items map[string]interface{}, expiration time.Duration) error {
	for key, value := range items {
		if err := CacheSet(key, value, expiration); err != nil {
			return fmt.Errorf("failed to warm cache for key %s: %w", key, err)
		}
	}
	return nil
}

// StaleWhileRevalidate serves stale cache while refreshing in background
func StaleWhileRevalidate(key string, dest interface{}, expiration time.Duration, loadFn func() (interface{}, error)) error {
	// Try to get from cache
	err := CacheGet(key, dest)
	if err == nil {
		// Cache hit - check if it's close to expiration
		ttl, _ := CacheGetTTL(key)
		if ttl < expiration/10 { // If less than 10% TTL remaining
			// Refresh in background
			go func() {
				data, err := loadFn()
				if err == nil {
					CacheSet(key, data, expiration)
				}
			}()
		}
		return nil
	}

	// Cache miss - load synchronously
	return CacheGetOrSet(key, dest, expiration, loadFn)
}

// CacheStampedePrevention prevents cache stampede using distributed lock
func CacheStampedePrevention(key string, dest interface{}, expiration time.Duration, loadFn func() (interface{}, error)) error {
	// Try to get from cache first
	err := CacheGet(key, dest)
	if err == nil {
		return nil // Cache hit
	}

	// Cache miss - try to acquire lock
	locked, err := CacheLock(key, 10*time.Second)
	if err != nil {
		return err
	}

	if !locked {
		// Another process is loading - wait and retry
		time.Sleep(100 * time.Millisecond)
		return CacheGet(key, dest)
	}

	// We got the lock - load data
	defer CacheUnlock(key)

	// Double check cache (might have been loaded by another process)
	err = CacheGet(key, dest)
	if err == nil {
		return nil
	}

	// Load from source
	result, err := loadFn()
	if err != nil {
		return err
	}

	// Store in cache
	if err := CacheSet(key, result, expiration); err != nil {
		return err
	}

	// Marshal to dest
	return CacheGet(key, dest)
}

// GetCachePattern returns standard cache key patterns
type CachePattern struct{}

var Pattern = CachePattern{}

// User returns user-related cache pattern
func (p CachePattern) User(userID uint) string {
	return fmt.Sprintf("user:%d", userID)
}

// Post returns post-related cache pattern
func (p CachePattern) Post(postID uint) string {
	return fmt.Sprintf("post:%d", postID)
}

// PostList returns post list cache pattern
func (p CachePattern) PostList(page int) string {
	return fmt.Sprintf("posts:list:page:%d", page)
}

// UserPosts returns user's posts cache pattern
func (p CachePattern) UserPosts(userID uint) string {
	return fmt.Sprintf("user:%d:posts", userID)
}

// All returns pattern to match all cache keys
func (p CachePattern) All() string {
	return "*"
}
