# Database Management

Dokumentasi lengkap untuk database management, migrations, seeding, backup, dan optimasi.

**Supported Databases**: MySQL & PostgreSQL

## Daftar Isi

- [Migration System](#migration-system)
- [Database Seeder](#database-seeder)
- [Backup & Restore](#backup--restore)
- [Connection Pooling](#connection-pooling)
- [Read Replica](#read-replica)
- [Best Practices](#best-practices)

---

## Migration System

Menggunakan **golang-migrate** untuk database schema versioning. Support untuk MySQL dan PostgreSQL.

### Database Type Detection

Migration system otomatis mendeteksi database type dari `DB_TYPE` environment variable:
- `mysql` - Menggunakan MySQL migration driver
- `postgres` - Menggunakan PostgreSQL migration driver

### Setup

Migration files disimpan di folder `migrations/`:

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_posts_table.up.sql
└── 000002_create_posts_table.down.sql
```

### Creating Migrations

**Menggunakan Helper Function**:
```go
import "starter-gofiber/helper"

// Creates timestamped migration files
helper.CreateMigration("create_comments_table")
```

Output:
```
Created migration files:
  - migrations/1704153600_create_comments_table.up.sql
  - migrations/1704153600_create_comments_table.down.sql
```

**Menggunakan CLI**:
```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

### Writing Migrations

**UP Migration** (`*_up.sql`):

**MySQL**:
```sql
-- Migration: create_users_table
-- Created at: 2024-01-01T10:00:00Z

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_email (email),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**PostgreSQL**:
```sql
-- Migration: create_users_table
-- Created at: 2024-01-01T10:00:00Z

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_deleted_at ON users(deleted_at);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

**DOWN Migration** (`*_down.sql`):
```sql
-- Migration: create_users_table
-- Created at: 2024-01-01T10:00:00Z

DROP TABLE IF EXISTS users;
```

### Running Migrations

**Programmatically**:
```go
import "starter-gofiber/helper"

// Run all pending migrations
if err := helper.RunMigrations(config.DB); err != nil {
    log.Fatal(err)
}
```

**Using CLI**:

**MySQL**:
```bash
# Run all up migrations
migrate -path ./migrations -database "mysql://user:pass@tcp(localhost:3306)/dbname" up

# Rollback last migration
migrate -path ./migrations -database "mysql://user:pass@tcp(localhost:3306)/dbname" down 1

# Go to specific version
migrate -path ./migrations -database "mysql://user:pass@tcp(localhost:3306)/dbname" goto 2
```

**PostgreSQL**:
```bash
# Run all up migrations
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# Rollback last migration
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down 1

# Go to specific version
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" goto 2
```

### Migration Commands

```go
// Get current version
version, dirty, err := helper.GetMigrationVersion(config.DB)

// Rollback migrations
err := helper.RollbackMigration(config.DB, 1) // Rollback 1 step

// Force version (fix dirty state)
err := helper.ForceMigrationVersion(config.DB, 2)
```

---

## Database Seeder

Populate database dengan sample data untuk development/testing.

### Available Seeders

Default seeders registered:
- `users` - Creates 3 sample users
- `posts` - Creates 5 sample posts

### Running Seeders

**Run All Seeders**:
```go
import "starter-gofiber/helper"

if err := helper.RunAllSeeders(config.DB); err != nil {
    log.Fatal(err)
}
```

**Run Specific Seeder**:
```go
// Run only user seeder
if err := helper.RunSeeder(config.DB, "users"); err != nil {
    log.Fatal(err)
}
```

### Creating Custom Seeders

**Register New Seeder**:
```go
// In your init() or main
helper.RegisterSeeder("categories", func(db *gorm.DB) error {
    categories := []entity.Category{
        {Name: "Technology"},
        {Name: "Business"},
        {Name: "Lifestyle"},
    }
    
    for _, cat := range categories {
        var existing entity.Category
        if db.Where("name = ?", cat.Name).First(&existing).Error == gorm.ErrRecordNotFound {
            if err := db.Create(&cat).Error; err != nil {
                return err
            }
        }
    }
    
    return nil
})
```

### Default Seeded Data

**Users** (password: `password123`):
- admin@example.com - Admin User
- test@example.com - Test User
- demo@example.com - Demo User

**Posts**:
- Getting Started with Go
- Building REST APIs with Fiber
- Database Best Practices
- Testing in Go
- Microservices Architecture

### Resetting Database

⚠️ **Development Only**:
```go
// Truncates all tables and re-runs seeders
if err := helper.ResetDatabase(config.DB); err != nil {
    log.Fatal(err)
}
```

---

## Backup & Restore

Automated database backup scripts dengan compression dan retention policy.

**Available Scripts**:
- MySQL: `backup_mysql.sh`, `restore_mysql.sh`, `auto_backup_mysql.sh`
- PostgreSQL: `backup_postgres.sh`, `restore_postgres.sh`, `auto_backup_postgres.sh`

### Backup Script

**Manual Backup**:

**MySQL**:
```bash
./scripts/backup/backup_mysql.sh
```

**PostgreSQL**:
```bash
./scripts/backup/backup_postgres.sh
```

Output:
```
🔄 Starting database backup...
Database: myapp
Backup file: ./backups/backup_myapp_20240101_100000.sql
📦 Compressing backup...
✅ Backup completed successfully!
File: ./backups/backup_myapp_20240101_100000.sql.gz
Size: 1.2M
```

**Features**:
- Single transaction/snapshot (consistent backup)
- Includes routines, triggers, events
- GZIP compression
- Auto cleanup (keeps last 7 days)

### Restore Script

**Restore from Backup**:

**MySQL**:
```bash
./scripts/backup/restore_mysql.sh ./backups/backup_myapp_20240101_100000.sql.gz
```

**PostgreSQL**:
```bash
./scripts/backup/restore_postgres.sh ./backups/backup_myapp_20240101_100000.sql.gz
```

**Interactive Confirmation**:
```
⚠️  WARNING: This will REPLACE all data in database: myapp
Backup file: ./backups/backup_myapp_20240101_100000.sql.gz
Are you sure you want to continue? (yes/no):
```

### Automated Daily Backup

**Setup Cron Job**:

**MySQL**:
```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /path/to/project/scripts/backup/auto_backup_mysql.sh
```

**PostgreSQL**:
```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /path/to/project/scripts/backup/auto_backup_postgres.sh
```

**Or Use Auto Backup Script**:

**MySQL**: `./scripts/backup/auto_backup_mysql.sh`

**PostgreSQL**: `./scripts/backup/auto_backup_postgres.sh`

**Features**:
- Logging to `logs/backup.log`
- Configurable retention (default: 30 days)
- Backup statistics
- Optional webhook notifications

### Backup Configuration

```bash
# Environment variables
BACKUP_DIR=./backups          # Backup directory
RETENTION_DAYS=30             # Keep backups for N days
LOG_FILE=./logs/backup.log    # Log file location
```

### Backup Best Practices

1. **Regular Schedule** - Daily backups minimum
2. **Off-site Storage** - Copy to S3/Cloud Storage
3. **Test Restores** - Regularly test backup restoration
4. **Monitor Disk Space** - Ensure enough space for backups
5. **Encrypt Backups** - Encrypt sensitive data backups

---

## Connection Pooling

Optimized database connection pool configuration.

### Configuration

Connection pool settings auto-adjust based on environment:

**Production** (`ENV_TYPE=prod`):
```go
MaxIdleConns:     25           // Minimum idle connections
MaxOpenConns:     200          // Maximum open connections
ConnMaxLifetime:  1 hour       // Connection reuse lifetime
ConnMaxIdleTime:  10 minutes   // Idle connection timeout
```

**Development** (`ENV_TYPE=dev`):
```go
MaxIdleConns:     10
MaxOpenConns:     50
ConnMaxLifetime:  30 minutes
ConnMaxIdleTime:  5 minutes
```

### Tuning Guidelines

**MaxIdleConns** - Minimum idle connections:
- Set to average concurrent queries
- Too low → Frequent connection creation
- Too high → Wasted resources

**MaxOpenConns** - Maximum total connections:
- Consider database server limits
- MySQL default: 151 connections
- Formula: `(RAM in GB) * 100` as starting point
- Monitor with `SHOW PROCESSLIST`

**ConnMaxLifetime** - Connection age limit:
- Prevents stale connections
- Recommended: 30min - 2 hours
- Must be < database server timeout

**ConnMaxIdleTime** - Idle timeout:
- Close unused connections
- Recommended: 5-15 minutes
- Frees resources during low traffic

### Monitoring Pool

```go
import "starter-gofiber/config"

sqlDB, _ := config.DB.DB()
stats := sqlDB.Stats()

log.Printf("Open Connections: %d", stats.OpenConnections)
log.Printf("In Use: %d", stats.InUse)
log.Printf("Idle: %d", stats.Idle)
log.Printf("Wait Count: %d", stats.WaitCount)
log.Printf("Wait Duration: %v", stats.WaitDuration)
log.Printf("Max Idle Closed: %d", stats.MaxIdleClosed)
log.Printf("Max Lifetime Closed: %d", stats.MaxLifetimeClosed)
```

---

## Read Replica

Load balancing dengan read replicas untuk scale-out reads.

### Setup

**Environment Configuration**:

**MySQL**:
```bash
# .env
# Primary Database (Write)
DB_TYPE=mysql
DB_URL=localhost:3306
DB_USER=root
DB_PASS=password
DB_NAME=myapp

# Read Replica (Optional)
DB_READ_HOST=replica.example.com
DB_READ_PORT=3306
DB_READ_USER=readonly
DB_READ_PASS=password
DB_READ_NAME=myapp
```

**PostgreSQL**:
```bash
# .env
# Primary Database (Write)
DB_TYPE=postgres
DB_URL=localhost:5432
DB_USER=postgres
DB_PASS=password
DB_NAME=myapp

# Read Replica (Optional)
DB_READ_HOST=replica.example.com
DB_READ_PORT=5432
DB_READ_USER=readonly
DB_READ_PASS=password
DB_READ_NAME=myapp
```

**Initialize Read Replica**:
```go
// main.go
config.LoadDB()  // Load primary DB

// Load read replica (falls back to primary if not configured)
config.LoadReadReplica()
```

### Usage

**Automatic Routing**:
```go
import "starter-gofiber/config"

// Use read replica for queries
users := []entity.User{}
config.UseReadReplica().Find(&users)

// Use primary DB for writes
config.UseWriteDB().Create(&user)
```

**DB Resolver**:
```go
resolver := config.NewDBResolver()

// Read operations
resolver.GetRead().Where("status = ?", "active").Find(&posts)

// Write operations
resolver.GetWrite().Create(&newPost)
```

**In Services**:
```go
type UserService struct {
    dbRead  *gorm.DB
    dbWrite *gorm.DB
}

func NewUserService() *UserService {
    return &UserService{
        dbRead:  config.UseReadReplica(),
        dbWrite: config.UseWriteDB(),
    }
}

func (s *UserService) GetAll() ([]entity.User, error) {
    var users []entity.User
    // Use read replica
    err := s.dbRead.Find(&users).Error
    return users, err
}

func (s *UserService) Create(user *entity.User) error {
    // Use primary DB
    return s.dbWrite.Create(user).Error
}
```

### Read Replica Configuration

**Connection Pool for Read Replica**:

Production:
```go
MaxIdleConns:     50   // Higher for read-heavy workload
MaxOpenConns:     300  // More connections for reads
ConnMaxLifetime:  1 hour
ConnMaxIdleTime:  10 minutes
```

### Best Practices

1. **Read Operations Only** - Never write to replicas
2. **Replication Lag** - Be aware of eventual consistency
3. **Critical Reads** - Use primary for immediately consistent reads
4. **Monitor Lag** - Watch replication lag metrics
5. **Failover** - Handle replica unavailability gracefully

### Replication Lag

**Check Lag**:
```sql
-- On MySQL replica
SHOW SLAVE STATUS\G
-- Look for: Seconds_Behind_Master
```

**Handle in Code**:
```go
// For critical reads, use primary DB
user := entity.User{}
config.UseWriteDB().Where("id = ?", userID).First(&user)
```

---

## Best Practices

### 1. Schema Design

- Use appropriate data types
- Add indexes for frequently queried columns
- Use foreign keys for referential integrity
- Normalize to reduce redundancy
- Denormalize for performance (where needed)

### 2. Indexing

```sql
-- Single column index
CREATE INDEX idx_users_email ON users(email);

-- Composite index (order matters!)
CREATE INDEX idx_posts_user_status ON posts(user_id, status);

-- Covering index
CREATE INDEX idx_users_lookup ON users(email, name, created_at);

-- Full-text index
CREATE FULLTEXT INDEX idx_posts_content ON posts(content);
```

### 3. Query Optimization

**Use EXPLAIN**:
```go
config.DB.Debug().
    Where("status = ?", "active").
    Find(&posts)
// Check EXPLAIN output for query plan
```

**Avoid N+1 Queries**:
```go
// ❌ Bad: N+1 queries
var users []entity.User
db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts) // N queries
}

// ✅ Good: Use Preload
var users []entity.User
db.Preload("Posts").Find(&users) // 2 queries total
```

### 4. Transactions

```go
// Use transaction for multiple operations
err := config.DB.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err // Rollback
    }
    
    if err := tx.Create(&profile).Error; err != nil {
        return err // Rollback
    }
    
    return nil // Commit
})
```

### 5. Batch Operations

```go
// ❌ Bad: Loop creates
for _, user := range users {
    db.Create(&user)
}

// ✅ Good: Batch create
db.CreateInBatches(users, 100) // Batch size: 100
```

### 6. Soft Deletes

```go
// Already included in gorm.Model
type User struct {
    gorm.Model // Includes DeletedAt
}

// Soft delete
db.Delete(&user) // Sets DeletedAt

// Include deleted records
db.Unscoped().Find(&users)

// Permanent delete
db.Unscoped().Delete(&user)
```

### 7. Database Maintenance

```sql
-- Analyze tables
ANALYZE TABLE users, posts;

-- Optimize tables
OPTIMIZE TABLE users;

-- Check table status
SHOW TABLE STATUS LIKE 'users';

-- Monitor slow queries
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2; -- Queries > 2 seconds
```

---

## Troubleshooting

### Connection Pool Exhausted

**Symptoms**:
- Timeouts waiting for connection
- High `WaitCount` in pool stats

**Solutions**:
```go
// Increase MaxOpenConns
sqlDB.SetMaxOpenConns(300)

// Or reduce connection lifetime
sqlDB.SetConnMaxLifetime(30 * time.Minute)

// Check for connection leaks
// Always close rows:
rows, _ := db.Raw("SELECT ...").Rows()
defer rows.Close()
```

### Slow Queries

**Debug Query**:
```go
config.DB.Debug().Where("...").Find(&results)
```

**Add Indexes**:
```sql
CREATE INDEX idx_column ON table(column);
```

**Use Query Cache**:
```go
// Cache query results
helper.CacheGetOrSet("users:all", &users, 5*time.Minute, func() (interface{}, error) {
    return repository.GetAll()
})
```

### Migration Failed (Dirty State)

**Check Version**:
```go
version, dirty, _ := helper.GetMigrationVersion(config.DB)
fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
```

**Force Clean**:
```go
// Force to version (use with caution!)
helper.ForceMigrationVersion(config.DB, 5)
```

### Backup Failed

**Check Permissions**:
```bash
# Ensure backup directory is writable
chmod 755 ./backups
```

**Check Disk Space**:
```bash
df -h
```

**Manual mysqldump/pg_dump**:

**MySQL**:
```bash
mysqldump -u user -p database > backup.sql
```

**PostgreSQL**:
```bash
pg_dump -U user -d database > backup.sql
```

---

## References

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [GORM Documentation](https://gorm.io/)
- [MySQL Connection Pool](https://dev.mysql.com/doc/connector-j/8.0/en/connector-j-usagenotes-j2ee-concepts-connection-pooling.html)
- [PostgreSQL Connection Pool](https://www.postgresql.org/docs/current/runtime-config-connection.html)
- [Database Indexing Best Practices](https://use-the-index-luke.com/)
- [PostgreSQL Performance Tuning](https://wiki.postgresql.org/wiki/Performance_Optimization)

---

**Last Updated**: January 2026
