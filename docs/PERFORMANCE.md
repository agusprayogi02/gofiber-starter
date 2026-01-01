# Performance Optimization Guide

Dokumentasi lengkap untuk performance optimization termasuk database indexing, request timeout, dan graceful shutdown.

## Daftar Isi

- [Database Indexing](#database-indexing)
- [Request Timeout](#request-timeout)
- [Graceful Shutdown](#graceful-shutdown)
- [Best Practices](#best-practices)

---

## Database Indexing

Database indexing adalah salah satu cara terpenting untuk meningkatkan performa query. Index membantu database menemukan data lebih cepat tanpa harus scan seluruh table.

### Index Utilities

Package `pkg/database/indexing.go` menyediakan utilities untuk membuat dan mengelola indexes:

```go
import "starter-gofiber/pkg/database"

// Create single index
index := database.IndexDefinition{
    Name:    "idx_users_email",
    Table:   "users",
    Columns: []string{"email"},
    Unique:  true,
}

err := database.CreateIndex(db, index)
```

### Index Types

**Single Column Index:**
```go
index := database.IndexDefinition{
    Name:    "idx_users_created_at",
    Table:   "users",
    Columns: []string{"created_at"},
}
```

**Composite Index (Multiple Columns):**
```go
index := database.IndexDefinition{
    Name:    "idx_posts_status_created",
    Table:   "posts",
    Columns: []string{"status", "created_at"},
}
```

**Unique Index:**
```go
index := database.IndexDefinition{
    Name:    "idx_users_email",
    Table:   "users",
    Columns: []string{"email"},
    Unique:  true,
}
```

**Partial Index (PostgreSQL only):**
```go
index := database.IndexDefinition{
    Name:    "idx_posts_active",
    Table:   "posts",
    Columns: []string{"created_at"},
    Where:   "status = 'active'",
}
```

**Concurrent Index (PostgreSQL only):**
```go
index := database.IndexDefinition{
    Name:       "idx_users_email",
    Table:      "users",
    Columns:    []string{"email"},
    Concurrent: true, // Doesn't lock table during creation
}
```

### Recommended Indexes

Package menyediakan recommended indexes untuk common tables:

```go
recommended := database.RecommendedIndexes()

// Create all recommended indexes for users table
if indexes, ok := recommended["users"]; ok {
    err := database.CreateIndexes(db, indexes)
}
```

### Check Index Existence

```go
exists, err := database.IndexExists(db, "users", "idx_users_email")
if err != nil {
    return err
}

if !exists {
    // Create index
    err := database.CreateIndex(db, index)
}
```

### List Table Indexes

```go
indexes, err := database.GetTableIndexes(db, "users")
if err != nil {
    return err
}

for _, indexName := range indexes {
    fmt.Println(indexName)
}
```

### Drop Index

```go
err := database.DropIndex(db, "users", "idx_users_email")
```

### Migration Example

```go
// In migration file
func up(db *gorm.DB) error {
    indexes := []database.IndexDefinition{
        {Name: "idx_users_email", Table: "users", Columns: []string{"email"}, Unique: true},
        {Name: "idx_users_created_at", Table: "users", Columns: []string{"created_at"}},
        {Name: "idx_posts_user_id", Table: "posts", Columns: []string{"user_id"}},
        {Name: "idx_posts_status_created", Table: "posts", Columns: []string{"status", "created_at"}},
    }
    
    return database.CreateIndexes(db, indexes)
}

func down(db *gorm.DB) error {
    indexes := []string{
        "idx_users_email",
        "idx_users_created_at",
        "idx_posts_user_id",
        "idx_posts_status_created",
    }
    
    for _, indexName := range indexes {
        if err := database.DropIndex(db, "users", indexName); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Index Best Practices

1. **Index Frequently Queried Columns**
   - Columns used in WHERE clauses
   - Columns used in JOIN conditions
   - Columns used in ORDER BY

2. **Composite Index Order Matters**
   ```sql
   -- Good: matches WHERE user_id = ? ORDER BY created_at
   CREATE INDEX idx_posts_user_created ON posts(user_id, created_at);
   
   -- Bad: doesn't help with WHERE user_id = ?
   CREATE INDEX idx_posts_created_user ON posts(created_at, user_id);
   ```

3. **Don't Over-Index**
   - Each index slows down INSERT/UPDATE/DELETE
   - Only index columns that are frequently queried

4. **Use Partial Indexes (PostgreSQL)**
   - Index only rows that match a condition
   - Smaller index = faster queries

5. **Monitor Index Usage**
   ```sql
   -- PostgreSQL
   SELECT * FROM pg_stat_user_indexes WHERE idx_scan = 0;
   
   -- MySQL
   SELECT * FROM sys.schema_unused_indexes;
   ```

---

## Request Timeout

Request timeout mencegah request yang terlalu lama menghabiskan resources server.

### Configuration

Set di `.env`:
```env
REQUEST_TIMEOUT=30  # Timeout in seconds (default: 30)
```

### How It Works

Request timeout dikonfigurasi di Fiber app:
- **ReadTimeout**: Maximum time to read request body
- **WriteTimeout**: Maximum time to write response
- **IdleTimeout**: Keep-alive timeout (default: 120 seconds)

```go
conf := fiber.Config{
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

### Timeout Behavior

- **ReadTimeout**: Jika client tidak mengirim request body dalam waktu yang ditentukan, connection akan ditutup
- **WriteTimeout**: Jika server tidak bisa mengirim response dalam waktu yang ditentukan, connection akan ditutup
- **IdleTimeout**: Jika tidak ada activity dalam waktu yang ditentukan, connection akan ditutup

### Best Practices

1. **Set Appropriate Timeout**
   - API endpoints: 30-60 seconds
   - File uploads: 300+ seconds
   - Long-running operations: Use background jobs instead

2. **Different Timeouts for Different Routes**
   ```go
   // Long timeout for file upload
   upload := app.Group("/upload")
   upload.Use(func(c *fiber.Ctx) error {
       c.Context().SetReadTimeout(5 * time.Minute)
       return c.Next()
   })
   ```

3. **Handle Timeout Errors**
   ```go
   if err := c.BodyParser(&data); err != nil {
       if err == context.DeadlineExceeded {
           return c.Status(408).JSON(fiber.Map{
               "error": "Request timeout",
           })
       }
       return err
   }
   ```

---

## Graceful Shutdown

Graceful shutdown memastikan aplikasi menutup semua resources dengan benar sebelum exit, mencegah data loss dan connection leaks.

### How It Works

1. **Signal Handling**: Aplikasi mendengarkan SIGTERM dan SIGINT
2. **Shutdown Sequence**: Menutup resources dalam urutan yang benar
3. **Timeout Protection**: Force exit jika shutdown terlalu lama

### Shutdown Sequence

1. Stop accepting new requests
2. Wait for active requests to complete
3. Shutdown Fiber server
4. Shutdown Asynq scheduler
5. Shutdown Asynq worker server
6. Close database connections
7. Close Redis connections
8. Flush logs

### Configuration

Set di `.env`:
```env
SHUTDOWN_TIMEOUT=10  # Timeout in seconds (default: 10)
```

### Implementation

**API Server (`cmd/api/main.go`):**
```go
// Wait for interrupt signal
<-quit
logger.Info("Shutting down server...")

// Create shutdown context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Graceful shutdown sequence
go func() {
    app.Shutdown()
    config.AsynqScheduler.Shutdown()
    config.AsynqServer.Shutdown()
    // Close connections...
}()

// Wait for completion or timeout
select {
case <-shutdownDone:
    logger.Info("Shutdown complete")
case <-ctx.Done():
    logger.Warn("Shutdown timeout exceeded")
}
```

**Worker Server (`cmd/worker/main.go`):**
```go
// Similar shutdown sequence for worker server
```

### Testing Graceful Shutdown

**Send SIGTERM:**
```bash
# Find process ID
ps aux | grep starter-gofiber

# Send SIGTERM
kill -TERM <PID>

# Or use killall
killall -TERM starter-gofiber
```

**Expected Behavior:**
1. Server stops accepting new requests
2. Active requests complete
3. Resources closed gracefully
4. Server exits cleanly

### Docker/Kubernetes Integration

**Docker:**
```yaml
services:
  app:
    stop_grace_period: 30s  # Give time for graceful shutdown
```

**Kubernetes:**
```yaml
spec:
  terminationGracePeriodSeconds: 30
```

### Troubleshooting

**Shutdown takes too long:**
- Increase `SHUTDOWN_TIMEOUT`
- Check for hanging database queries
- Check for long-running background jobs

**Connection leaks:**
- Ensure all connections are closed in shutdown sequence
- Use defer statements for cleanup
- Monitor connection pool stats

**Data loss:**
- Ensure transactions are committed before shutdown
- Use background jobs for long operations
- Implement request queuing

---

## Best Practices

### 1. Database Indexing

- ✅ Index foreign keys
- ✅ Index columns used in WHERE clauses
- ✅ Index columns used in ORDER BY
- ✅ Use composite indexes for multi-column queries
- ✅ Monitor index usage and remove unused indexes
- ❌ Don't index every column
- ❌ Don't create indexes on small tables (< 1000 rows)

### 2. Request Timeout

- ✅ Set appropriate timeout based on endpoint type
- ✅ Use background jobs for long operations
- ✅ Handle timeout errors gracefully
- ✅ Log timeout occurrences for monitoring
- ❌ Don't set timeout too high (wastes resources)
- ❌ Don't set timeout too low (causes false timeouts)

### 3. Graceful Shutdown

- ✅ Always implement graceful shutdown
- ✅ Set appropriate shutdown timeout
- ✅ Close all connections properly
- ✅ Wait for active requests to complete
- ✅ Log shutdown progress
- ❌ Don't force exit immediately
- ❌ Don't ignore shutdown errors

### 4. Performance Monitoring

```go
// Monitor query performance
db.Debug().Where("...").Find(&results)

// Monitor index usage
indexes, _ := database.GetTableIndexes(db, "users")

// Monitor connection pool
sqlDB, _ := db.DB()
stats := sqlDB.Stats()
logger.Info("DB stats",
    zap.Int("open_connections", stats.OpenConnections),
    zap.Int("in_use", stats.InUse),
    zap.Int("idle", stats.Idle),
)
```

---

## Examples

### Complete Setup Example

```go
// 1. Create indexes during migration
func up(db *gorm.DB) error {
    indexes := database.RecommendedIndexes()["users"]
    return database.CreateIndexes(db, indexes)
}

// 2. Configure timeouts
conf := fiber.Config{
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
}

// 3. Implement graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Shutdown sequence...
```

---

## References

- [PostgreSQL Indexing](https://www.postgresql.org/docs/current/indexes.html)
- [MySQL Indexing](https://dev.mysql.com/doc/refman/8.0/en/mysql-indexes.html)
- [Fiber Timeouts](https://docs.gofiber.io/api/app#readtimeout)
- [Go Context Timeout](https://golang.org/pkg/context/#WithTimeout)

---

**Last Updated**: January 2026

