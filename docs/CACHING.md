# Caching with Redis

Dokumentasi lengkap sistem caching menggunakan Redis untuk meningkatkan performa aplikasi.

## Daftar Isi

- [Overview](#overview)
- [Setup & Configuration](#setup--configuration)
- [Cache Middleware](#cache-middleware)
- [Query Result Caching](#query-result-caching)
- [Cache Invalidation](#cache-invalidation)
- [Advanced Patterns](#advanced-patterns)
- [Best Practices](#best-practices)
- [Monitoring](#monitoring)

---

## Overview

Sistem caching ini menyediakan:

- **Response Caching** - Cache HTTP responses untuk GET requests
- **Query Result Caching** - Cache hasil database queries
- **Distributed Lock** - Prevent cache stampede dengan distributed locking
- **Pattern-based Invalidation** - Invalidate cache secara terstruktur
- **Multiple Strategies** - TTL, Write-Through, LRU, Stale-While-Revalidate

### Benefits

- **Reduced Database Load** - Kurangi query repetitive ke database
- **Faster Response Time** - Serve dari memory instead of disk
- **Better Scalability** - Handle lebih banyak requests dengan resource yang sama
- **Cost Savings** - Kurangi biaya infrastructure

---

## Setup & Configuration

### 1. Install Redis

**Docker (Recommended)**:
```bash
docker run -d --name redis \
  -p 6379:6379 \
  redis:alpine redis-server --appendonly yes
```

**Ubuntu/Debian**:
```bash
sudo apt update
sudo apt install redis-server
sudo systemctl enable redis-server
sudo systemctl start redis-server
```

**macOS**:
```bash
brew install redis
brew services start redis
```

### 2. Environment Configuration

```bash
# .env
REDIS_ENABLE=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=        # Leave empty if no password
REDIS_DB=0             # Database number (0-15)
```

### 3. Connection Settings

Konfigurasi connection pool di [config/redis.go](config/redis.go):

```go
redis.Options{
    Addr:         "localhost:6379",
    Password:     "",          // No password
    DB:           0,           // Default DB
    DialTimeout:  5 * time.Second,
    ReadTimeout:  3 * time.Second,
    WriteTimeout: 3 * time.Second,
    PoolSize:     10,          // Max connections
    MinIdleConns: 5,           // Min idle connections
}
```

---

## Cache Middleware

Middleware untuk cache HTTP responses secara otomatis.

### Basic Usage

```go
// config/app.go
import "starter-gofiber/middleware"

// Cache all GET requests dengan default 5 menit TTL
app.Use(middleware.SimpleCacheMiddleware())
```

### Custom Configuration

```go
cacheConfig := middleware.CacheConfig{
    Expiration:      10 * time.Minute,
    ExcludePaths:    []string{"/health", "/metrics", "/admin"},
    ExcludeMethods:  []string{"POST", "PUT", "PATCH", "DELETE"},
    KeyPrefix:       "api:cache:",
    ExcludeStatuses: []int{500, 502, 503, 504},
}

app.Use(middleware.CacheMiddleware(cacheConfig))
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Expiration` | `time.Duration` | 5 minutes | Cache TTL |
| `ExcludePaths` | `[]string` | `/health`, `/metrics`, `/admin` | Skip caching |
| `ExcludeMethods` | `[]string` | POST, PUT, PATCH, DELETE | Only GET cached |
| `KeyPrefix` | `string` | `cache:` | Cache key prefix |
| `ExcludeStatuses` | `[]int` | 500, 502, 503, 504 | Don't cache errors |

### Response Headers

```
X-Cache: HIT      # Served from cache
X-Cache: MISS     # Served from origin
```

### Per-Route Caching

```go
// Cache specific route group
api := app.Group("/api")
api.Use(middleware.SimpleCacheMiddleware())

// Don't cache admin routes
admin := app.Group("/admin")
// No cache middleware
```

---

## Query Result Caching

Helper functions untuk cache hasil database queries.

### Simple Get/Set

```go
import "starter-gofiber/helper"

// Set cache
err := helper.CacheSet("user:123", userData, 10*time.Minute)

// Get cache
var user entity.User
err := helper.CacheGet("user:123", &user)
if err == redis.Nil {
    // Cache miss - load from database
}
```

### Cache-Aside Pattern (Lazy Loading)

```go
func (s *PostService) GetByID(id uint) (*entity.Post, error) {
    var post entity.Post
    cacheKey := helper.Pattern.Post(id)
    
    // Try cache first, load from DB if miss
    err := helper.CacheGetOrSet(cacheKey, &post, 5*time.Minute, func() (interface{}, error) {
        return s.repo.GetByID(id)
    })
    
    return &post, err
}
```

### Preventing Cache Stampede

```go
func (s *PostService) GetPopularPosts() ([]entity.Post, error) {
    var posts []entity.Post
    cacheKey := "posts:popular"
    
    // Use distributed lock to prevent stampede
    err := helper.CacheStampedePrevention(cacheKey, &posts, 10*time.Minute, func() (interface{}, error) {
        return s.repo.GetPopularPosts()
    })
    
    return posts, err
}
```

### Stale-While-Revalidate

```go
func (s *UserService) GetProfile(userID uint) (*entity.User, error) {
    var user entity.User
    cacheKey := helper.Pattern.User(userID)
    
    // Serve stale cache while refreshing in background
    err := helper.StaleWhileRevalidate(cacheKey, &user, 30*time.Minute, func() (interface{}, error) {
        return s.repo.GetByID(userID)
    })
    
    return &user, err
}
```

---

## Cache Invalidation

Strategi untuk invalidate cache saat data berubah.

### Individual Key

```go
// Delete specific cache
helper.CacheDelete("user:123")

// Delete multiple keys
helper.BulkInvalidate("user:123", "user:456", "posts:list")
```

### Pattern-Based Invalidation

```go
// Invalidate all user-related cache
helper.InvalidateUser(userID)

// Invalidate all posts
helper.InvalidateCollection("posts")

// Invalidate related resources
helper.InvalidateRelated("post", postID)

// Custom pattern
helper.InvalidateByTag("trending")
```

### Write-Through Pattern

```go
func (h *PostHandler) Update(c *fiber.Ctx) error {
    // ... validation ...
    
    // Update and invalidate cache atomically
    err := helper.WriteThroughPattern("post:*", func() error {
        return h.service.Update(&post, id)
    })
    
    if err != nil {
        return err
    }
    
    return helper.Response(...)
}
```

### Cache Strategies

#### 1. **Time-To-Live (TTL)**

```go
// Auto-expire after TTL
helper.CacheSet("data", value, 10*time.Minute)
```

#### 2. **Write-Through**

```go
// Invalidate immediately on write
helper.WriteThrough("user:123", func() error {
    return updateDatabase()
})
```

#### 3. **Write-Back** (Eventual Consistency)

```go
// Update cache first, DB later
helper.CacheSet("user:123", newData, time.Hour)
go func() {
    time.Sleep(5 * time.Second)
    updateDatabase()
}()
```

---

## Advanced Patterns

### 1. Cache Key Builder

```go
import "starter-gofiber/helper"

// Build structured cache key
key := helper.NewCacheKeyBuilder("posts").
    Add("list").
    AddInt(page).
    AddUint(userID).
    Build()
// Result: "posts:list:1:42"
```

### 2. Tagged Cache

```go
// Tag cache for grouped invalidation
key := helper.TaggedCacheKey("user:123", "trending", "featured")
helper.CacheSet(key, data, time.Hour)

// Invalidate all "trending" cache
helper.InvalidateByTag("trending")
```

### 3. Cache Warming

```go
// Pre-load frequently accessed data
warmingData := map[string]interface{}{
    "stats:total_users": userCount,
    "stats:total_posts": postCount,
    "config:settings": siteSettings,
}

helper.CacheWarming(warmingData, 24*time.Hour)
```

### 4. Distributed Lock

```go
// Acquire lock
locked, err := helper.CacheLock("process:import", 10*time.Second)
if !locked {
    return errors.New("process already running")
}
defer helper.CacheUnlock("process:import")

// ... critical section ...
```

### 5. Counter Operations

```go
// Increment counter
count, err := helper.CacheIncrement("views:post:123", 24*time.Hour)

// Decrement counter
count, err := helper.CacheDecrement("stock:product:456")
```

---

## Best Practices

### 1. Cache Key Naming Convention

```go
// Good: Hierarchical, descriptive
"user:123:profile"
"post:456:comments"
"stats:daily:2024-01-01"

// Bad: Flat, unclear
"u123"
"data456"
"x"
```

### 2. Set Appropriate TTL

```go
// Frequently changing data - short TTL
helper.CacheSet("trending:posts", data, 5*time.Minute)

// Rarely changing data - long TTL
helper.CacheSet("categories", data, 24*time.Hour)

// Static data - very long TTL
helper.CacheSet("site:config", data, 7*24*time.Hour)
```

### 3. Handle Cache Misses Gracefully

```go
var user entity.User
err := helper.CacheGet("user:123", &user)

if err == redis.Nil {
    // Cache miss - load from DB
    user, err = loadFromDatabase(123)
    if err != nil {
        return err
    }
    
    // Repopulate cache
    helper.CacheSet("user:123", user, 10*time.Minute)
} else if err != nil {
    // Redis error - fall back to DB
    return loadFromDatabase(123)
}

return user
```

### 4. Invalidate Proactively

```go
func (h *PostHandler) Create(c *fiber.Ctx) error {
    post, err := h.service.Create(&req)
    if err != nil {
        return err
    }
    
    // Invalidate related cache
    helper.InvalidateCollection("posts")           // List cache
    helper.InvalidateUser(req.UserID)             // User's posts
    helper.InvalidateByTag("trending")             // Trending posts
    
    return helper.Response(...)
}
```

### 5. Monitor Cache Performance

```go
// Get cache statistics
stats := helper.CacheStats()
log.Printf("Cache Stats: %+v", stats)
// Output: {enabled:true hits:1234 misses:56 total_conns:10 idle_conns:5}
```

### 6. Use Compression for Large Data

```go
import "compress/gzip"

// Before caching
compressed := compressData(largeData)
helper.CacheSet("large:data", compressed, time.Hour)

// After retrieving
data, _ := helper.CacheGet("large:data", &compressed)
original := decompressData(compressed)
```

### 7. Avoid Caching Sensitive Data

```go
// ❌ Don't cache
helper.CacheSet("user:password", hashedPassword, time.Hour)

// ✅ Use encrypted storage or skip caching
// Store passwords only in database with encryption
```

---

## Monitoring

### Cache Statistics

```go
// Handler untuk monitoring
func CacheStatsHandler(c *fiber.Ctx) error {
    stats := helper.CacheStats()
    return c.JSON(stats)
}

// Response:
{
    "enabled": true,
    "hits": 12500,
    "misses": 340,
    "timeouts": 2,
    "total_conns": 10,
    "idle_conns": 7,
    "stale_conns": 0,
    "db_size": 1523
}
```

### Health Check

```go
func RedisHealthCheck(c *fiber.Ctx) error {
    if helper.RedisClient == nil {
        return c.Status(503).JSON(fiber.Map{
            "status": "unavailable",
            "cache": "disabled",
        })
    }
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    if err := helper.RedisClient.Ping(ctx).Err(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status": "unhealthy",
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "status": "healthy",
        "cache": "enabled",
    })
}
```

### Cache Hit Rate

```go
hitRate := float64(stats["hits"].(uint32)) / float64(stats["hits"].(uint32) + stats["misses"].(uint32)) * 100
log.Printf("Cache Hit Rate: %.2f%%", hitRate)
```

### Redis CLI Monitoring

```bash
# Monitor commands in real-time
redis-cli monitor

# Get server info
redis-cli info

# Check key count
redis-cli dbsize

# List all keys (use with caution in production!)
redis-cli keys "*"

# Get key TTL
redis-cli ttl "user:123"
```

---

## Common Use Cases

### 1. User Profile Caching

```go
func (s *UserService) GetProfile(userID uint) (*dto.UserProfile, error) {
    var profile dto.UserProfile
    key := helper.Pattern.User(userID)
    
    err := helper.CacheGetOrSet(key, &profile, 15*time.Minute, func() (interface{}, error) {
        // Load from database
        user, err := s.repo.GetByID(userID)
        if err != nil {
            return nil, err
        }
        
        // Transform to DTO
        return &dto.UserProfile{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            // ... more fields
        }, nil
    })
    
    return &profile, err
}
```

### 2. Paginated List Caching

```go
func (s *PostService) GetAll(params *dto.Pagination) ([]entity.Post, error) {
    var posts []entity.Post
    key := helper.Pattern.PostList(params.Page)
    
    err := helper.CacheGetOrSet(key, &posts, 5*time.Minute, func() (interface{}, error) {
        return s.repo.GetAll(params)
    })
    
    return posts, err
}
```

### 3. API Rate Limiting with Redis

```go
func CheckRateLimit(userID uint) (bool, error) {
    key := fmt.Sprintf("ratelimit:user:%d", userID)
    
    count, err := helper.CacheIncrement(key, 1*time.Minute)
    if err != nil {
        return false, err
    }
    
    if count > 100 {
        return false, errors.New("rate limit exceeded")
    }
    
    return true, nil
}
```

### 4. Session Storage

```go
func StoreSession(sessionID string, data interface{}) error {
    key := fmt.Sprintf("session:%s", sessionID)
    return helper.CacheSet(key, data, 30*time.Minute)
}

func GetSession(sessionID string, dest interface{}) error {
    key := fmt.Sprintf("session:%s", sessionID)
    return helper.CacheGet(key, dest)
}
```

---

## Troubleshooting

### Cache Not Working

**Check Redis connection**:
```bash
redis-cli ping
# Should return: PONG
```

**Check environment**:
```bash
echo $REDIS_ENABLE
echo $REDIS_HOST
```

**Check logs**:
```bash
grep "Redis" logs/app.log
```

### High Cache Miss Rate

- Increase TTL untuk data yang jarang berubah
- Review cache invalidation strategy
- Check if cache keys are consistent
- Monitor popular queries dan warm cache

### Memory Issues

**Check memory usage**:
```bash
redis-cli info memory
```

**Set max memory**:
```bash
# redis.conf
maxmemory 256mb
maxmemory-policy allkeys-lru
```

**Clear cache**:
```go
helper.CacheFlush() // Use with caution!
```

---

## Performance Tips

1. **Use pipelining** untuk multiple operations
2. **Batch invalidation** instead of one-by-one
3. **Compress large objects** before caching
4. **Use appropriate data structures** (strings, hashes, sets)
5. **Monitor and tune** connection pool size
6. **Set realistic TTLs** - not too short, not too long
7. **Use Redis clustering** untuk high-availability production

---

## References

- [Redis Documentation](https://redis.io/documentation)
- [go-redis Client](https://github.com/redis/go-redis)
- [Caching Best Practices](https://aws.amazon.com/caching/best-practices/)
- [Cache Stampede Prevention](https://en.wikipedia.org/wiki/Cache_stampede)

---

**Last Updated**: January 2026
