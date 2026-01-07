# Database Management

Dokumentasi lengkap untuk database management, migrations, seeding, backup, dan optimasi.

**Supported Databases**: MySQL, PostgreSQL & SQL Server

## Daftar Isi

- [Migration System (Atlas)](#migration-system-atlas)
- [Database Seeder](#database-seeder)
- [Backup & Restore](#backup--restore)
- [Connection Pooling](#connection-pooling)
- [Read Replica](#read-replica)
- [Soft Delete](#soft-delete)
- [Audit Log](#audit-log)
- [Best Practices](#best-practices)

---

## Migration System (Atlas)

Menggunakan **Atlas with GORM Provider** untuk automatic schema migrations dari GORM entities. Atlas generate migrations otomatis dari perubahan GORM models, tidak perlu menulis SQL manual!

### Why Atlas?

✅ **Auto-generate migrations** dari GORM models
✅ **Type-safe** - Migrations dari Go code, bukan SQL manual
✅ **Version control friendly** - Migration files di-commit ke git
✅ **Database agnostic** - Support MySQL, PostgreSQL, SQL Server
✅ **Smart diffing** - Hanya generate perubahan yang diperlukan
✅ **Rollback support** - Automatic down migrations
✅ **CI/CD ready** - CLI commands untuk automation

### Installation

**Install Atlas CLI**:

```bash
# Using curl (Linux/macOS)
curl -sSf https://atlasgo.sh | sh

# Using Homebrew (macOS)
brew install ariga/tap/atlas

# Using Go
go install ariga.io/atlas@latest

# Verify installation
atlas version
```

**Or using Makefile**:
```bash
make atlas-install
```

### Setup

Project sudah dikonfigurasi dengan:
- `atlas.hcl` - Atlas configuration file
- `cmd/atlas/` - GORM schema loader (uses models from config)
- `internal/config/database.go` - Single source of truth for models
- `Makefile` - Atlas commands integrated in main Makefile
- `migrations/` - Migration files directory

### Environment Configuration

Set database URL di `.env`:

```bash
# Development
ATLAS_DEV_DB_URL=mysql://user:pass@localhost:3306/myapp_dev

# Production
ATLAS_DB_URL=mysql://user:pass@localhost:3306/myapp
```

**PostgreSQL**:
```bash
ATLAS_DEV_DB_URL=postgres://user:pass@localhost:5432/myapp_dev?sslmode=disable
ATLAS_DB_URL=postgres://user:pass@localhost:5432/myapp?sslmode=disable
```

### Quick Start Workflow

**1. Modify GORM Entity**:

```go
// internal/domain/product/entity.go
package product

type Product struct {
    ID          uint    `gorm:"primaryKey"`
    Name        string  `gorm:"type:varchar(200);not null"`
    Price       float64 `gorm:"type:decimal(10,2);not null"`
    Description string  `gorm:"type:text"`
    gorm.Model
}
```

**2. Register Model** in `internal/config/database.go`:

```go
func GetModelsForMigration() []interface{} {
    models := []interface{}{
        &user.User{},
        &post.Post{},
        &product.Product{}, // Add new model here
        // ... other models
    }
    return models
}
```

**Note**: Single source of truth! Used by both AutoMigrate and Atlas.

**3. Generate Migration**:

```bash
# Generate migration from schema changes
atlas migrate diff create_products_table --env dev
```

Output:
```
Analyzing GORM models...
Comparing with current database schema...
Generated: migrations/20260107120000_create_products_table.sql
```

**4. Review Migration**:

```sql
-- migrations/20260107120000_create_products_table.sql
-- Create "products" table
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `description` text,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_products_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

**5. Apply Migration**:

```bash
# Apply to development database
atlas migrate apply --env dev

# Or apply to production
atlas migrate apply --env prod
```

### Atlas Commands

**Generate Migrations**:

```bash
# Generate with auto-name
atlas migrate diff --env dev

# Generate with custom name
atlas migrate diff add_user_bio --env dev

# Generate for specific environment
atlas migrate diff --env prod
```

**Apply Migrations**:

```bash
# Apply all pending migrations
atlas migrate apply --env dev

# Apply specific number of migrations
atlas migrate apply --env dev --to-version 20260107120000

# Dry run (preview without applying)
atlas migrate apply --env dev --dry-run
```

**Migration Status**:

```bash
# Check current status
atlas migrate status --env dev

# Output:
# Migration Status: OK
# Current Version: 20260107120000
# Next Version: 20260107130000
# Pending Migrations: 1
```

**Inspect Schema**:

```bash
# Inspect current database schema
atlas schema inspect --env dev

# Inspect and save to file
atlas schema inspect --env dev > schema.sql
```

**Validate Migrations**:

```bash
# Validate migration files
atlas migrate validate --env dev

# Lint migrations
atlas migrate lint --env dev
```

### Makefile Commands

All Atlas commands are integrated in main `Makefile`:

```bash
# Install Atlas CLI
make atlas-install

# Generate migration from GORM models
make atlas-diff

# Generate with custom name
make atlas-diff-name NAME=add_user_bio

# Apply migrations
make atlas-apply

# Check status
make atlas-status

# Validate migrations
make atlas-validate

# Inspect database
make atlas-inspect

# Full workflow (diff + validate + apply)
make atlas-migrate

# Test with dry-run
make atlas-test

# View all available commands
make help
```

### Development Workflow

**Option 1: AutoMigrate (Quick Prototyping)**

Set `DB_GEN=true` di `.env` untuk development:

```bash
DB_GEN=true  # Use GORM AutoMigrate
```

Server akan otomatis sync schema saat startup menggunakan `db.AutoMigrate()`.

⚠️ **Only for local development!** Jangan gunakan di production.

**Option 2: Atlas Migrations (Recommended)**

Set `DB_GEN=false` dan gunakan Atlas:

```bash
DB_GEN=false  # Use Atlas migrations
```

1. Modify GORM entities
2. Generate migration: `make atlas-diff`
3. Review generated SQL
4. Apply migration: `make atlas-apply`

### Production Workflow

**1. Generate Migration (Local)**:

```bash
# Di environment development
make atlas-diff-dev
```

**2. Commit Migration Files**:

```bash
git add migrations/
git commit -m "Add products table migration"
git push
```

**3. Apply in Production**:

```bash
# Deploy ke production server
# Set production database URL
export ATLAS_DB_URL="postgres://user:pass@prod-db:5432/myapp"

# Apply migrations
make atlas-apply-prod

# Or using CLI
atlas migrate apply --env prod
```

### Advanced Features

**Multiple Databases**:

```hcl
// atlas.hcl
env "db1" {
  url = var.db1_url
  migration {
    dir = "file://migrations/db1"
  }
}

env "db2" {
  url = var.db2_url
  migration {
    dir = "file://migrations/db2"
  }
}
```

**Migration Policies**:

```hcl
// Prevent destructive changes in production
env "prod" {
  lint {
    destructive {
      error = true  // Block DROP operations
    }
  }
  
  diff {
    skip {
      drop_schema = true
      drop_table  = true
    }
  }
}
```

**Custom Migration**:

Jika perlu custom migration (data migration, complex logic):

```bash
# Generate empty migration
atlas migrate new custom_data_migration --env dev
```

Edit file dan tambahkan SQL manual:

```sql
-- migrations/20260107140000_custom_data_migration.sql
-- Update existing records
UPDATE users SET role = 'user' WHERE role IS NULL;

-- Add constraint
ALTER TABLE users ADD CONSTRAINT chk_role CHECK (role IN ('admin', 'user'));
```

### Troubleshooting

**"Failed to connect to database"**:

```bash
# Check database URL
echo $ATLAS_DB_URL

# Test connection
atlas schema inspect --env dev
```

**"Dirty database state"**:

```bash
# Check migration status
atlas migrate status --env dev

# Force to specific version (use with caution!)
atlas migrate set 20260107120000 --env dev
```

**"Schema drift detected"**:

Database schema berbeda dengan migrations:

```bash
# Inspect current state
atlas schema inspect --env dev

# Generate migration to fix drift
atlas migrate diff fix_schema_drift --env dev
```

**Schema Not Updating**:

Ensure model is registered in `internal/config/database.go`:

```go
func GetModelsForMigration() []interface{} {
    models := []interface{}{
        &user.User{},
        &product.Product{}, // Make sure new model is here
        // ... other models
    }
    return models
}
```

This is the single source of truth used by both AutoMigrate and Atlas.

### Migration Best Practices

1. **Always Review Generated Migrations** - Check SQL before applying
2. **Test Locally First** - Apply to dev database before production
3. **Version Control** - Commit migration files to git
4. **Backward Compatible** - Design migrations that won't break running app
5. **Data Migrations** - Run separately from schema migrations
6. **Rollback Plan** - Test down migrations before deploying
7. **CI/CD Integration** - Automate migration testing in pipeline

### CI/CD Integration

**GitHub Actions Example**:

```yaml
name: Database Migration

on:
  push:
    branches: [main]
    paths:
      - 'migrations/**'

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Atlas
        run: curl -sSf https://atlasgo.sh | sh
      
      - name: Validate Migrations
        run: atlas migrate validate --env dev
        env:
          ATLAS_DEV_DB_URL: ${{ secrets.DB_URL }}
      
      - name: Apply Migrations
        run: atlas migrate apply --env prod
        env:
          ATLAS_DB_URL: ${{ secrets.PROD_DB_URL }}
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
import "starter-gofiber/pkg/apierror"

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
import "starter-gofiber/internal/config"

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
import "starter-gofiber/internal/config"

// Use read replica for queries
users := []user.User{}
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

func (s *UserService) GetAll() ([]user.User, error) {
    var users []user.User
    // Use read replica
    err := s.dbRead.Find(&users).Error
    return users, err
}

func (s *UserService) Create(user *user.User) error {
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
user := user.User{}
config.UseWriteDB().Where("id = ?", userID).First(&user)
```

---

## Soft Delete

Implementasi soft delete menggunakan GORM's built-in soft delete dengan helper functions tambahan.

### How It Works

Semua entities yang menggunakan `gorm.Model` otomatis mendapat soft delete support:
```go
type User struct {
    gorm.Model  // Includes: ID, CreatedAt, UpdatedAt, DeletedAt
    Name  string
    Email string
}
```

### Basic Usage

**Soft Delete** (sets DeletedAt):
```go
// Delete user (soft delete)
db.Delete(&user)

// Delete by ID
db.Delete(&user.User{}, userID)
```

**Query Behavior**:
```go
// Normal queries exclude soft deleted records
var users []user.User
db.Find(&users)  // Only returns non-deleted users
```

### Helper Functions

**Include Soft Deleted Records**:
```go
import "starter-gofiber/pkg/apierror"

// Include deleted records
var users []user.User
helper.WithTrashed(db).Find(&users)

// Or using Scopes
db.Scopes(helper.SoftDeleteScope()).Find(&users)
```

**Only Deleted Records**:
```go
// Get only soft deleted records
var deletedUsers []user.User
helper.OnlyTrashed(db).Find(&deletedUsers)
```

**Restore Deleted Records**:
```go
// Restore a record
helper.Restore(db, &user)

// Restore by ID
helper.RestoreByID(db, &user.User{}, userID)

// Restore all deleted records of a model
helper.RestoreAll(db, &user.User{})
```

**Force Delete** (permanent):
```go
// Permanently delete a record
helper.ForceDelete(db, &user)

// Permanently delete by ID
helper.ForceDeleteByID(db, &user.User{}, userID)
```

**Utility Functions**:
```go
// Count soft deleted records
count, err := helper.CountTrashed(db, &user.User{})

// Check if a record is soft deleted
isTrashed, err := helper.IsTrashed(db, &user.User{}, userID)
```

### Example Usage

```go
// In your service
func (s *UserService) SoftDeleteUser(id uint) error {
    var user user.User
    if err := s.db.First(&user, id).Error; err != nil {
        return err
    }
    
    // Soft delete
    return s.db.Delete(&user).Error
}

func (s *UserService) RestoreUser(id uint) error {
    // Restore soft deleted user
    return helper.RestoreByID(s.db, &user.User{}, id)
}

func (s *UserService) GetDeletedUsers() ([]user.User, error) {
    var users []user.User
    err := helper.OnlyTrashed(s.db).Find(&users).Error
    return users, err
}

func (s *UserService) PermanentlyDeleteUser(id uint) error {
    // This cannot be undone!
    return helper.ForceDeleteByID(s.db, &user.User{}, id)
}
```

### Best Practices

1. **Default to Soft Delete** - Use soft delete for most user data
2. **Cleanup Strategy** - Schedule periodic cleanup of old soft deleted records
3. **Audit Before Force Delete** - Log before permanent deletions
4. **User Notification** - Inform users before permanent data deletion
5. **Cascade Deletes** - Handle related records appropriately

---

## Audit Log

Automatic tracking of all data changes (who, when, what) dengan GORM callbacks.

### Features

- **Automatic Logging** - Tracks CREATE, UPDATE, DELETE, RESTORE operations
- **User Tracking** - Records user ID, username, IP address, user agent
- **Change Tracking** - Stores old and new values in JSON format
- **Request Tracing** - Links changes to request ID for debugging
- **Query Support** - Filter logs by entity, user, date, action

### How It Works

Audit logging otomatis aktif untuk semua GORM operations. Sistem menggunakan GORM callbacks untuk intercept dan log setiap perubahan data.

### Audit Log Schema

```go
type AuditLog struct {
    ID          uint
    EntityType  string       // e.g., "users", "posts"
    EntityID    uint         // ID of affected record
    Action      AuditAction  // CREATE, UPDATE, DELETE, RESTORE
    Description string       // Human-readable description
    OldValues   string       // JSON of old values
    NewValues   string       // JSON of new values
    UserID      *uint        // User who performed action
    Username    string       // Cached username
    IPAddress   string       // IPv4 or IPv6
    UserAgent   string       // Browser/client info
    RequestID   string       // Trace request chain
    CreatedAt   time.Time
}
```

### Manual Logging

**With User Context**:
```go
import "starter-gofiber/pkg/apierror"

// Create audit logger
logger := helper.NewAuditLogger(config.DB).
    WithUser(userID, username).
    WithRequest(ipAddress, userAgent, requestID)

// Log create
logger.LogCreate("users", user.ID, user)

// Log update
logger.LogUpdate("users", user.ID, oldUser, newUser)

// Log delete
logger.LogDelete("users", user.ID, user, true) // true = soft delete

// Log restore
logger.LogRestore("users", user.ID, user)
```

### Automatic Logging

GORM callbacks sudah terdaftar dan akan otomatis log semua operations. Untuk memberikan user context, gunakan context dalam DB statement:

```go
import "context"

// Create context with user info
ctx := context.Background()
ctx = context.WithValue(ctx, "user_id", uint(123))
ctx = context.WithValue(ctx, "username", "john@example.com")
ctx = context.WithValue(ctx, "ip_address", "192.168.1.1")
ctx = context.WithValue(ctx, "user_agent", "Mozilla/5.0...")
ctx = context.WithValue(ctx, "request_id", "req-abc-123")

// Use context in query
db.WithContext(ctx).Create(&user)
db.WithContext(ctx).Updates(&user)
db.WithContext(ctx).Delete(&user)
```

### Query Audit Logs

**Get All Logs**:
```go
logs, total, err := helper.GetAuditLogs(config.DB, entity.AuditLogFilter{}, 1, 50)
```

**Filter by Entity**:
```go
filter := entity.AuditLogFilter{
    EntityType: "users",
    EntityID:   &userID,
}
logs, total, err := helper.GetAuditLogs(config.DB, filter, 1, 50)
```

**Filter by User**:
```go
filter := entity.AuditLogFilter{
    UserID: &adminID,
}
logs, total, err := helper.GetAuditLogs(config.DB, filter, 1, 50)
```

**Filter by Date Range**:
```go
startDate := time.Now().AddDate(0, 0, -7) // Last 7 days
endDate := time.Now()

filter := entity.AuditLogFilter{
    StartDate: &startDate,
    EndDate:   &endDate,
}
logs, total, err := helper.GetAuditLogs(config.DB, filter, 1, 50)
```

**Get Entity History**:
```go
// Get full history of a user
history, err := helper.GetEntityAuditHistory(config.DB, "users", userID)

// History shows all changes chronologically
for _, log := range history {
    fmt.Printf("%s - %s: %s\n", log.CreatedAt, log.Action, log.Description)
}
```

**Get User Activity**:
```go
startDate := time.Now().AddDate(0, 0, -30) // Last 30 days
activities, err := helper.GetUserActivity(config.DB, userID, &startDate, nil)
```

### Cleanup Old Logs

```go
// Delete audit logs older than 90 days
err := helper.CleanupOldAuditLogs(config.DB, 90)
```

### Example in Handler

```go
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    username := c.Locals("username").(string)
    requestID := c.Locals("request_id").(string)
    
    var updateData dto.UpdateUserRequest
    if err := c.BodyParser(&updateData); err != nil {
        return err
    }
    
    // Create context with user info
    ctx := context.Background()
    ctx = context.WithValue(ctx, "user_id", userID)
    ctx = context.WithValue(ctx, "username", username)
    ctx = context.WithValue(ctx, "ip_address", c.IP())
    ctx = context.WithValue(ctx, "user_agent", c.Get("User-Agent"))
    ctx = context.WithValue(ctx, "request_id", requestID)
    
    // Update will be automatically logged
    var user user.User
    if err := config.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
        return err
    }
    
    if err := config.DB.WithContext(ctx).Model(&user).Updates(updateData).Error; err != nil {
        return err
    }
    
    return c.JSON(fiber.Map{"message": "User updated successfully"})
}
```

### Best Practices

1. **Always Provide Context** - Include user, IP, request ID for complete audit trail
2. **Retention Policy** - Clean up old logs periodically (e.g., keep 1 year)
3. **Sensitive Data** - Don't log passwords or sensitive fields
4. **Performance** - Index entity_type, entity_id, user_id, created_at
5. **Monitoring** - Alert on unusual patterns (e.g., mass deletes)
6. **Compliance** - Essential for GDPR, SOC2, and other regulations

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
var users []user.User
db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts) // N queries
}

// ✅ Good: Use Preload
var users []user.User
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
