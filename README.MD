# Starter Template Go Fiber

Backend Golang With Fiber Framework

## Getting Started

### Prerequisites

Sebelum menjalankan project ini, pastikan sistem Anda telah memiliki:

#### Required or higher** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 14+** atau **MySQL 8+** atau **SQL Server 2019+** - Database utama
- **Redis 6+** - Untuk caching dan background job queue (Asynq)
- **OpenSSL** - Untuk generate RSA key pair (JWT signing)

#### Optional (Recommended)
- **Docker & Docker Compose** - Untuk menjalankan dependencies (Redis, Database, Mailpit)
- **Make** - Untuk automation tasks (jika menggunakan Makefile)
- **Git** - Version control
- **Mailpit** - Local SMTP server untuk testing email (development)
  ```bash
  docker run -d -p 1025:1025 -p 8025:8025 mailpit/mailpit
  ```

#### Development Tools
- **Air** - Live reload untuk development
- **golangci-lint** - Code linting (optional)
- **Postman** atau **Insomnia** - API testing

### Installation

#### 1. Clone Repository
```bash
git clone <repository-url>
cd starter-gofiber
```

#### 2. Install Dependencies
```bash
go mod download
go mod verify
```

#### 3. Install Development Tools
```bash
# Install Air for live reload
go install github.com/air-verse/air@latest

# (Optional) Install golangci-lint for code linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### 4. Setup Environment Variables
```bash
# Copy example environment file
cp .env.example .env

# Edit .env file dengan text editor favorit
nano .env  # atau vim, code, dll
```

**Konfigurasi minimal yang harus diisi:**
```env
# Application
ENV_TYPE=dev
APP_URL=http://localhost:3000

# Database (pilih salah satu)
# PostgreSQL
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=starter_gofiber
DB_USER=postgres
DB_PASSWORD=your_password

# Redis (required untuk caching dan background jobs)
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# SMTP (untuk email - gunakan Mailpit untuk development)
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_FROM=noreply@localhost
SMTP_FROM_NAME=Starter GoFiber
SMTP_USE_TLS=false

# Sentry (optional - untuk error tracking)
SENTRY_DSN=
```

#### 5. Generate RSA Key Pair untuk JWT
```bash
# Buat directory untuk certificates
mkdir -p assets/certs

# Generate private key (4096 bit RSA)
openssl genpkey -algorithm RSA -out assets/certs/certificate.pem -pkeyopt rsa_keygen_bits:4096

# Set permission (recommended)
chmod 600 assets/certs/certificate.pem
```

**Update .env:**
```env
LOCATION_CERT=assets/certs/certificate.pem
```

#### 6. Setup Database

**Opsi A: Menggunakan Docker Compose**
```bash
# Start PostgreSQL + Redis + Mailpit
docker-compose up -d

# Check status
docker-compose ps
```

**Opsi B: Manual Setup**
```bash
# Buat database (PostgreSQL example)
createdb starter_gofiber

# Atau menggunakan psql
psql -U postgres -c "CREATE DATABASE starter_gofiber;"
```

#### 7. Run Database Migrations
```bash
# Buat table dan structure database
go run ./cmd/api/main.go migrate up

# Atau menggunakan custom migration command jika ada
# go run scripts/migrate.go up
```

#### 8. (Optional) Seed Database
```bash
# Seeders akan dijalankan secara otomatis jika dikonfigurasi
# Atau bisa menggunakan pkg/database seeder secara manual
```

#### 9. Run Application

**Development Mode (dengan Air - Live Reload):**
```bash
# Air akan otomatis build dan run dari cmd/api/main.go
air
```

**Standard Mode:**
```bash
# Run langsung tanpa build
go run ./cmd/api/main.go
```

**Build dan Run Binary:**
```bash
# Build binary
go build -o bin/app ./cmd/api

# Run binary
./bin/app
```

#### 10. Verify Installation

Application akan berjalan di `http://localhost:3000` (atau port yang dikonfigurasi di `.env`)

**Health Check:**
```bash
curl http://localhost:3000/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "timestamp": "2026-01-01T12:00:00Z",
  "version": "1.0.0"
}
```

**Access Points:**
- API: `http://localhost:3000/api`
- Health: `http://localhost:3000/health`
- Metrics: `http://localhost:3000/metrics`
- Mailpit UI (email testing): `http://localhost:8025`

### Database Configuration Notes

#### PostgreSQL
Default dan recommended. Tidak perlu konfigurasi khusus.

#### MySQL
Pastikan menggunakan charset `utf8mb4`:
```env
DB_DRIVER=mysql
DB_CHARSET=utf8mb4
```

#### SQL Server
Untuk SQL Server, perlu adjust enum data type di model:
```go
// internal/domain/user/entity.go
type User struct {
    ID       uint     `gorm:"primaryKey;autoIncrement"`
    Name     string   `gorm:"type:varchar(200);not null"`
    Email    string   `gorm:"type:varchar(200);uniqueIndex;not null"`
    Password string   `gorm:"type:varchar(150);not null"`
    Role     UserRole `gorm:"type:varchar(10);default:user"` // Gunakan varchar untuk SQL Server
    gorm.Model
}
```

### Docker Deployment

```bash
# Build image
docker build -f assets/docker/Dockerfile -t starter-gofiber .

# Run with docker-compose
docker-compose -f assets/docker/docker-compose.yml up -d

# Check logs
docker-compose -f assets/docker/docker-compose.yml logs -f app
```

### Troubleshooting

#### Port Already in Use
```bash
# Find process using port 3000
lsof -ti:3000

# Kill process
kill -9 $(lsof -ti:3000)
```

#### Redis Connection Failed
```bash
# Check Redis is running
redis-cli ping
# Expected: PONG

# Start Redis (if not running)
redis-server
# Or with Docker:
docker run -d -p 6379:6379 redis:7-alpine
```

#### Database Connection Failed
- Check database credentials di `.env`
- Verify database server is running
- Check firewall rules
- Ensure database exists: `CREATE DATABASE starter_gofiber;`

## Project Structure

Project ini menggunakan **Go Standard Layout** dengan **Domain-Driven Design** dan **Clean Architecture**:

```
starter-gofiber/
├── assets/                      # Static assets dan files
│   ├── certs/                   # SSL certificates dan private keys (JWT)
│   ├── docker/                  # Docker configurations
│   │   ├── Dockerfile           # Default Dockerfile
│   │   ├── Dockerfile.optimized # Optimized multi-stage Dockerfile
│   │   ├── docker-compose.yml   # Default docker-compose
│   │   ├── docker-compose.minio.yml  # MinIO S3 storage setup
│   │   ├── dev/                 # Development docker setup
│   │   │   ├── docker-compose.yml
│   │   │   └── Dockerfile
│   │   └── prod/                # Production docker setup
│   │       ├── docker-compose.yml
│   │       └── Dockerfile
│   └── rbac/                    # Casbin RBAC policy files
│       ├── model.conf           # RBAC model definition
│       └── policy.csv           # Role-permission mappings
│
├── backups/                     # Database backups (auto-generated)
│
├── cmd/                         # Application entry points
│   ├── api/                     # Main API server
│   │   └── main.go              # API server entry point
│   └── worker/                  # Background worker server
│       └── main.go              # Worker server entry point
│
├── migrations/                  # Database migration files
│   ├── README.md                # Migration documentation
│   ├── 000003_add_user_profile_fields.up.sql
│   └── 000003_add_user_profile_fields.down.sql
│
├── migrations/                  # Database migration files
│   ├── 000003_add_user_profile_fields.up.sql
│   └── 000003_add_user_profile_fields.down.sql
│
├── cmd/                         # Application entry points
│   ├── api/                     # Main API server
│   │   └── main.go              # API server entry point
│   └── worker/                  # Background worker server
│       └── main.go              # Worker server entry point (optional)
│
├── docs/                        # Documentation files
│   ├── API_AUTH.md              # Authentication API docs
│   ├── API_FEATURES.md          # Advanced API features
│   ├── API_FEATURES_QUICK_REF.md # API features quick reference
│   ├── ARCHITECTURE.md          # Architecture overview
│   ├── BACKGROUND_JOBS.md       # Background jobs & Asynq
│   ├── BACKGROUND_JOBS_QUICK_REF.md # Background jobs quick reference
│   ├── CACHING.md               # Redis caching guide
│   ├── CLOUD_STORAGE_EXAMPLES.md # Cloud storage examples
│   ├── DATABASE.md              # Database management & Atlas migrations
│   ├── DEPLOYMENT.md            # Deployment guide
│   ├── DEVOPS.md                # DevOps & CI/CD guide
│   ├── EMAIL.md                 # Email system docs
│   ├── FILE_MANAGEMENT.md       # File upload & storage
│   ├── LOGGING_MONITORING.md    # Logging & monitoring
│   ├── MAKEFILE_COMMANDS.md     # Makefile commands reference
│   ├── METRICS_MONITORING.md    # Metrics & health checks
│   ├── MIGRATION_GUIDE.md       # Atlas migration quick start
│   ├── MINIO_SETUP.md           # MinIO S3 storage setup guide
│   ├── PERFORMANCE.md           # Performance optimization guide
│   ├── SECURITY.md              # Security features
│   ├── SECURITY_IMPLEMENTATION.md # Security implementation details
│   ├── SSE.md                   # Server-Sent Events
│   ├── TESTING.md               # Testing guide
│   └── TEST_COVERAGE.md         # Test coverage reports
│
├── internal/                    # Private application code (tidak bisa diimport dari luar)
│   ├── config/                  # Application configurations
│   │   ├── app.go               # App config (middleware setup)
│   │   ├── asynq.go             # Asynq background jobs config
│   │   ├── cloud_storage.go     # S3/MinIO storage config
│   │   ├── config.go            # Main config loader
│   │   ├── database.go          # Database connection config
│   │   ├── gorm_logger.go       # GORM logger with Zap
│   │   ├── permission.go        # Casbin permission init
│   │   ├── permissions.go       # Permission definitions
│   │   ├── read_replica.go      # Read replica config
│   │   ├── redis.go             # Redis connection config
│   │   ├── storage.go            # Storage config
│   │   └── timezone.go          # Timezone configuration
│   │
│   ├── domain/                  # Domain layer (business entities & interfaces)
│   │   ├── post/                # Post domain
│   │   │   ├── dto.go           # Post DTOs (requests/responses)
│   │   │   ├── entity.go        # Post entity (GORM model)
│   │   │   ├── repository.go    # Post repository interface
│   │   │   └── service.go       # Post service interface
│   │   └── user/                # User domain
│   │       ├── apikey.go        # API Key entity
│   │       ├── dto.go           # User DTOs (requests/responses)
│   │       ├── email_verification.go  # Email verification entity
│   │       ├── entity.go        # User entity (GORM model)
│   │       ├── password_reset.go      # Password reset entity
│   │       ├── preferences.go   # User preferences entity
│   │       ├── refresh_token.go       # Refresh token entity
│   │       ├── repository.go    # User repository interface
│   │       ├── service.go       # User service interface
│   │       └── token_dto.go     # JWT token DTOs
│   │
│   ├── handler/                 # HTTP handlers (presentation layer)
│   │   ├── http/                # HTTP handlers
│   │   │   ├── auth.go          # Authentication & profile endpoints
│   │   │   ├── health.go        # Health check endpoints
│   │   │   ├── jobs.go          # Job management endpoints
│   │   │   ├── metrics.go       # Metrics endpoint
│   │   │   ├── post.go          # Post CRUD endpoints
│   │   │   └── sse.go           # Server-Sent Events handler
│   │   └── middleware/          # HTTP middlewares
│   │       ├── apikey.go        # API Key authentication
│   │       ├── auth.go          # JWT authentication
│   │       ├── authz.go         # Casbin authorization
│   │       ├── cache.go         # Response caching
│   │       ├── compression.go   # Response compression
│   │       ├── csrf.go          # CSRF protection
│   │       ├── https.go         # HTTPS redirect
│   │       ├── ip_filter.go     # IP filtering
│   │       ├── logger.go         # Request logging
│   │       ├── prometheus.go    # Prometheus metrics
│   │       ├── rate_limiter.go   # Rate limiting
│   │       ├── security.go      # Security headers
│   │       └── sentry.go        # Sentry error tracking
│   │
│   ├── infrastructure/          # Infrastructure layer (external services)
│   │   ├── auth/                # Authentication infrastructure
│   │   │   └── apikey.go        # API Key management
│   │   ├── cache/               # Caching infrastructure
│   │   │   ├── redis.go         # Redis client wrapper
│   │   │   └── strategy.go      # Cache strategy patterns
│   │   ├── email/               # Email infrastructure
│   │   │   ├── config.go        # Email configuration
│   │   │   └── sender.go       # Email sender implementation
│   │   └── storage/             # Storage infrastructure
│   │       ├── handler.go       # File upload handler
│   │       ├── image.go         # Image processing
│   │       ├── s3.go            # S3 client wrapper
│   │       ├── upload.go        # File upload utilities
│   │       ├── validation.go    # File validation
│   │       └── versioning.go     # File versioning
│   │
│   ├── repository/               # Data access layer
│   │   └── postgres/            # PostgreSQL implementation
│   │       ├── email_verification.go
│   │       ├── password_reset.go
│   │       ├── post.go          # Post repository implementation
│   │       ├── refresh_token.go
│   │       └── user.go          # User repository implementation (includes preferences)
│   │
│   ├── service/                 # Business logic layer
│   │   ├── auth/                # Authentication service
│   │   │   └── service.go       # Auth service implementation
│   │   └── post/                # Post service
│   │       └── service.go       # Post service implementation
│   │
│   └── worker/                  # Background worker (Asynq)
│       ├── client.go            # Asynq client wrapper
│       ├── email.go             # Email job handlers
│       ├── handlers.go          # Job handler implementations
│       ├── jobs.go              # Job definitions & enqueue functions
│       └── scheduler.go         # Periodic task scheduler
│
├── pkg/                         # Public library code (bisa diimport dari luar)
│   ├── apierror/                # API error types
│   │   └── error.go            # Custom error definitions
│   ├── crypto/                  # Cryptographic utilities
│   │   ├── encryption.go       # AES encryption
│   │   ├── hash.go             # Hashing utilities
│   │   ├── jwt.go              # JWT token utilities
│   │   └── password.go         # Password hashing (bcrypt)
│   ├── database/                # Database utilities
│   │   ├── audit.go            # Audit logging
│   │   ├── bulk.go             # Bulk operations
│   │   ├── indexing.go         # Database indexing utilities
│   │   ├── migration.go        # Migration utilities
│   │   ├── pagination.go       # Pagination helpers
│   │   ├── seeder.go           # Database seeder
│   │   └── soft_delete.go      # Soft delete utilities
│   ├── dto/                     # Shared Data Transfer Objects
│   │   ├── paginate.go         # Pagination DTO
│   │   ├── response.go         # Standard API responses
│   │   ├── token.go            # Token DTOs
│   ├── logger/                  # Logging utilities
│   │   ├── logger.go           # Zap logger wrapper
│   │   └── sentry.go           # Sentry integration
│   ├── pagination/              # Pagination utilities
│   │   └── cursor.go           # Cursor-based pagination
│   ├── response/                # Response utilities
│   │   └── response.go         # Response formatters
│   ├── utils/                   # General utilities
│   │   ├── date.go             # Date/time utilities
│   │   ├── export.go           # Data export utilities
│   │   ├── filter.go           # Filtering utilities
│   │   ├── metrics.go          # Metrics utilities
│   │   ├── sort.go             # Sorting utilities
│   │   └── sse.go              # SSE utilities
│   └── validator/               # Validation utilities
│       ├── sanitize.go         # Input sanitization
│       └── validator.go        # Input validation
│
├── router/                      # Route definitions
│   ├── auth.go                  # Auth routes
│   ├── post.go                  # Post routes
│   ├── router.go                # Main router setup
│   └── sse.go                   # SSE routes
│
├── templates/                   # HTML templates
│   └── email/                   # Email templates
│       ├── welcome.html         # Welcome email HTML
│       ├── welcome.txt          # Welcome email text
│       ├── reset-password.html  # Password reset HTML
│       ├── reset-password.txt   # Password reset text
│       ├── verify-email.html    # Email verification HTML
│       └── verify-email.txt     # Email verification text
│
├── tests/                       # Test files
│   ├── auth_test.go             # Auth endpoint tests
│   ├── email_test.go            # Email system tests
│   ├── post_test.go             # Post endpoint tests
│   └── setup_test.go            # Test setup & helpers
│
├── variables/                   # Global constants & variables
│   └── constant.go              # App constants (includes AVATAR_PATH)
│
├── .air.toml                    # Air live reload config
├── .cursorrules                 # AI assistant rules & best practices guide
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
├── .golangci.yml                # golangci-lint configuration
├── .pre-commit-config.yaml      # Pre-commit hooks configuration
├── CONTRIBUTING.md              # Contributing guidelines
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── LICENSE                      # MIT License
├── Makefile                     # Automation tasks (build, test, migrate, etc.)
├── README.MD                    # This file
└── .github/                     # GitHub configuration
    └── workflows/               # GitHub Actions workflows
        └── ci.yml               # CI/CD pipeline
```

**Note:** Semua code baru harus menggunakan struktur di `internal/` dan `pkg/`. Shared DTOs berada di `pkg/dto/` dan dapat digunakan oleh semua domain.

### Architecture Layers

Project ini menggunakan **Clean Architecture** dengan **Domain-Driven Design**:

1. **Domain Layer** (`internal/domain/`) - Business entities, DTOs, dan interfaces
   - Setiap domain memiliki: `entity.go`, `dto.go`, `repository.go`, `service.go`
   - Domain tidak bergantung pada layer lain (pure business logic)

2. **Repository Layer** (`internal/repository/`) - Data access implementations
   - Implementasi dari interface yang didefinisikan di domain
   - Menggunakan GORM untuk database operations

3. **Service Layer** (`internal/service/`) - Business logic implementations
   - Implementasi dari interface yang didefinisikan di domain
   - Menggunakan repository untuk data access

4. **Handler Layer** (`internal/handler/`) - HTTP request handlers
   - HTTP handlers (`internal/handler/http/`)
   - Middleware (`internal/handler/middleware/`)
   - Menggunakan service untuk business logic

5. **Infrastructure Layer** (`internal/infrastructure/`) - External services
   - Email, Cache, Storage, Auth infrastructure
   - Implementasi untuk external dependencies

6. **Worker Layer** (`internal/worker/`) - Background jobs
   - Asynq job handlers dan schedulers
   - Background task processing

**Data Flow:**
```
Client Request → Router → Middleware → Handler → Service → Repository → Database
                                                    ↓
                                    Infrastructure (Email, Cache, Storage)
                                                    ↓
                                            Worker (Background Jobs)
```

## Warning

- Untuk `Order` pada helper digunakan untuk menandai urutan dari error yang muncul.
	- Keterangan kode huruf pertama: 
		- S: Service
		- H: Handler
		- R: Repository
		- U: Util
	- Contoh:
		- S1: Ada pada service dengan error handling pertama

## Aturan
- jika ada error pada `service`, `util`, `helper` wajib mengembalikan error menggunakan `helper` error handling contoh: 
```go
return nil, &helper.BadRequestError{
	Message: err.Error(),
	Order:   "S1",
}
```
- Agar pada `handler` tinggal return error, contoh:
```go
if err != nil {
	return err
}
```

## Rilis Production dengan docker
- pada file `.env` ganti `ENV_TYPE="prod"`
- ganti file `docker-composer.yml` dan `Dockerfile` dengan yang ada didalam folder `assets/docker/prod`, jangan lupa menyesuaikan configurasi yang akan digunakan 


## Authz
- Authz dibuat dengan `casbin` library.
	- https://github.com/gofiber/contrib/blob/main/casbin/casbin_test.go

## Fitur yang Masih Kurang

### 1. Authentication & Authorization
- [x] **Refresh Token** - Implementasi refresh token untuk memperpanjang sesi pengguna tanpa harus login ulang
- [x] **Logout** - Endpoint untuk logout dan invalidasi token
- [x] **Forgot Password** - Fitur reset password via email (backend ready, perlu email service)
- [x] **Email Verification** - Verifikasi email setelah registrasi (backend ready, perlu email service)
- [x] **Change Password** - Endpoint untuk ubah password user yang sudah login
- [x] **Session Management** - Manage active sessions pengguna

### Dokumentasi
- [docs/API_AUTH.md](docs/API_AUTH.md) - Dokumentasi lengkap endpoint authentication

### 2. Testing
- [x] **Integration Tests** - Testing end-to-end untuk API endpoints
- [x] **Test Coverage Report** - Setup coverage report dengan minimal threshold (Current: 23.2%)
- [x] **Mock Database** - Setup mock untuk testing database operations
- [x] **API Testing Collection** - Postman collection untuk manual testing

### Dokumentasi
Lihat [docs/TESTING.md](docs/TESTING.md) untuk dokumentasi lengkap dari testing
Lihat [docs/TEST_COVERAGE.md](docs/TEST_COVERAGE.md) untuk dokumentasi lengkap test coverage report.

### 3. Documentation
- [x] **Postman Collection** - Export collection untuk API testing
- [x] **Architecture Documentation** - Dokumentasi arsitektur dan design pattern yang digunakan
- [x] **Deployment Guide** - Panduan deploy Docker Compose dengan production setup
- [x] **Contributing Guide** - Guidelines untuk kontributor dengan coding standards
- [x] **LICENSE File** - MIT License untuk open source project

### Dokumentasi
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) - Clean Architecture pattern dan design patterns
- [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) - Deployment guide untuk Docker Compose
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contributing guidelines dan coding standards

### 4. Logging & Monitoring
- [x] **Structured Logging** - Zap logger with Fiber integration via fiberzap middleware
- [x] **Request ID Tracking** - Trace request dengan unique ID (UUID)
- [x] **Access Logs** - Logging untuk semua incoming requests dengan detail lengkap
- [x] **Error Tracking** - Integration dengan Sentry untuk production error tracking
- [x] **Health Check Endpoint** - `/health`, `/health/ready`, `/health/live` untuk orchestration
- [x] **Lightweight Metrics** - Zero-dependency metrics di `/metrics` (tidak pakai Prometheus, super ringan!)
- [x] **Database Query Logging** - Logging slow queries dengan custom GORM logger

### Dokumentasi
Lihat [docs/LOGGING_MONITORING.md](docs/LOGGING_MONITORING.md) dan [docs/METRICS_MONITORING.md](docs/METRICS_MONITORING.md) untuk dokumentasi lengkap

### 5. Security
- [x] **API Rate Limiting per User** - Rate limit berdasarkan user ID/IP dengan in-memory store
- [x] **CSRF Protection** - Cross-Site Request Forgery protection dengan token validation
- [x] **Content Security Policy** - CSP headers configuration dengan custom policy support
- [x] **SQL Injection Prevention** - Prepared statements validation dan pattern detection
- [x] **XSS Protection** - Input sanitization dengan bluemonday
- [x] **HTTPS Redirect** - Auto redirect HTTP ke HTTPS di production
- [x] **Helmet Middleware** - Security headers (X-Frame-Options, HSTS, dll)
- [x] **API Key Authentication** - Alternative auth method untuk machine-to-machine access
- [x] **IP Whitelist/Blacklist** - Restrict access by IP dengan wildcard support
- [x] **Encryption at Rest** - Sensitive data encryption dengan AES-256-GCM

### Dokumentasi
Lihat [docs/SECURITY.md](docs/SECURITY.md) untuk dokumentasi lengkap fitur keamanan

### 6. Caching
- [x] **Redis Integration** - Caching layer dengan Redis client (go-redis v9)
- [x] **Cache Middleware** - Response caching untuk GET requests dengan configurable TTL
- [x] **Query Result Caching** - Helper functions untuk cache database query results
- [x] **Cache Invalidation Strategy** - Write-Through, TTL, pattern-based, dan cache stampede prevention

### Dokumentasi
Lihat [docs/CACHING.md](docs/CACHING.md) untuk dokumentasi lengkap caching system

### 7. Database
- [x] **Database Migration System** - Proper migration tool (golang-migrate) dengan helper functions
- [x] **Seeder** - Database seeder untuk development/testing dengan sample data
- [x] **Database Backup Scripts** - Automated backup mechanism dengan compression dan retention
- [x] **Connection Pooling Configuration** - Optimal connection pool settings (prod vs dev)
- [x] **Read Replica Support** - Load balancing untuk read operations dengan DBResolver
- [x] **Soft Delete Support** - Soft delete implementation dengan helper functions (WithTrashed, OnlyTrashed, Restore, ForceDelete)
- [x] **Audit Log** - Track semua perubahan data (who, when, what) dengan automatic GORM callbacks

### Dokumentasi
Lihat [docs/DATABASE.md](docs/DATABASE.md) untuk dokumentasi lengkap database management

### 8. File Management
- [x] **File Validation** - Validate file type, size, dan content dengan magic number checking
- [x] **Multiple File Upload** - Support upload multiple files sekaligus dengan concurrent processing
- [x] **Cloud Storage Integration** - AWS S3 integration dengan presigned URLs dan S3-compatible services
- [x] **Image Processing** - Resize, compress, thumbnail generation, rotate, flip, blur, sharpen, brightness/contrast adjustment
- [x] **File Versioning** - Keep track of file versions dengan checksum verification dan restore capabilities

### Dokumentasi
- Lihat [docs/FILE_MANAGEMENT.md](docs/FILE_MANAGEMENT.md) untuk dokumentasi lengkap file management system
- Lihat [Setup Minio](docs/MINIO_SETUP.md) untuk panduan setup MinIO sebagai local S3-compatible storage

### 9. API Features
- [x] **Response Compression** - Gzip compression untuk response dengan 3 compression levels
- [x] **Pagination Cursor-based** - Alternative untuk offset pagination, O(1) performance untuk large datasets
- [x] **Bulk Operations** - Bulk create/update/delete/restore/upsert dengan error tracking per record
- [x] **Export Data** - Export ke CSV/Excel/PDF dengan styling dan pagination otomatis
- [x] **Search & Filter** - Advanced search dengan 15+ operators (eq, gt, like, in, between, dll) dan AND/OR logic
- [x] **Sorting** - Multiple field sorting support dengan validation dan SQL injection protection

### Dokumentasi
Lihat [docs/API_FEATURES.md](docs/API_FEATURES.md) untuk dokumentasi lengkap API features dengan contoh penggunaan

### 10. Real-time Features
- [x] **Server-Sent Events (SSE)** - Real-time push notifications dari server ke client dengan event types, user targeting, broadcast, dan auto keep-alive

### Dokumentasi
Lihat [docs/SSE.md](docs/SSE.md) untuk dokumentasi lengkap SSE dengan contoh client implementation (JS, React, Vue)

### 11. Background Jobs
- [x] **Job Queue** - Asynq (Redis-based) untuk background task processing
- [x] **Scheduled Tasks** - Cron-like periodic tasks (mirip Laravel Scheduler)
- [x] **Email Queue** - Async email sending dengan priority queues
- [x] **Retry Mechanism** - Auto retry failed jobs dengan exponential backoff
- [x] **Job Dashboard** - Queue statistics, pause/resume, task management

### Dokumentasi
Lihat [docs/BACKGROUND_JOBS.md](docs/BACKGROUND_JOBS.md) untuk dokumentasi lengkap background jobs dan queue system

### 12. Email System
- [x] **Email Templates** - HTML + text email templates dengan responsive design
- [x] **Email Service Integration** - Multi-provider SMTP support (Mailpit, Mailgun, SendGrid, AWS SES, Gmail)
- [x] **Email Queue** - Async email processing via Asynq dengan auto-retry
- [x] **Email Verification** - Template untuk verify email dengan expiration
- [x] **Password Reset Email** - Template untuk reset password dengan security warnings
- [x] **Welcome Email** - Professional welcome email untuk new users
- [x] **Custom Email Support** - Send custom emails dengan HTML/text dan attachments
- [x] **Template System** - Go html/template engine dengan data binding

### Dokumentasi
Lihat [docs/EMAIL.md](docs/EMAIL.md) untuk dokumentasi lengkap email notification system dengan contoh konfigurasi SMTP providers dan template customization

### 13. DevOps & CI/CD
- [x] **Makefile** - Automation untuk common tasks (build, test, migrate)
- [x] **CI/CD Pipeline** - GitHub Actions untuk automated testing & deployment
- [x] **Docker Optimization** - Multi-stage build untuk smaller image size (~90% reduction)
- [x] **Pre-commit Hooks** - Automated checks sebelum commit (linting, formatting)

### Dokumentasi
Lihat [docs/DEVOPS.md](docs/DEVOPS.md) untuk dokumentasi lengkap DevOps & CI/CD setup

### 14. Performance
- [x] **Database Indexing** - Proper indexes untuk optimize queries dengan utilities dan recommended indexes
- [x] **Request Timeout** - Timeout configuration untuk all requests (ReadTimeout, WriteTimeout, IdleTimeout)
- [x] **Graceful Shutdown** - Proper shutdown handling dengan timeout protection untuk zero-downtime deployment

### Dokumentasi
Lihat [docs/PERFORMANCE.md](docs/PERFORMANCE.md) untuk dokumentasi lengkap performance optimization

### 15. User Management
- [x] **User Profile** - Get/Update user profile dengan avatar upload
- [x] **User Preferences** - Settings/preferences per user (notifications, privacy, UI settings)

### Dokumentasi
- [docs/API_AUTH.md](docs/API_AUTH.md) - Dokumentasi lengkap authentication endpoints

## 📚 Dokumentasi Lengkap

> **📖 Complete Index**: See [docs/README.md](docs/README.md) for organized documentation index by topic



### 🏗️ Architecture & Setup
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - Clean Architecture pattern, design patterns, dan layer responsibilities
- **[docs/MIGRATION_GUIDE.md](docs/MIGRATION_GUIDE.md)** - Quick start guide untuk Atlas database migrations (Bahasa Indonesia)
- **[docs/DATABASE.md](docs/DATABASE.md)** - Database management lengkap: Atlas migrations, seeding, backup, optimization
- **[docs/MAKEFILE_COMMANDS.md](docs/MAKEFILE_COMMANDS.md)** - Reference lengkap untuk semua Makefile commands
- **[cmd/README.md](cmd/README.md)** - Convention untuk cmd/ folder dan cara run programs

### 🔐 Security & Authentication
- **[docs/API_AUTH.md](docs/API_AUTH.md)** - Authentication endpoints: register, login, profile, password reset
- **[docs/SECURITY.md](docs/SECURITY.md)** - Security features: rate limiting, CSRF, XSS, encryption, IP filtering
- **[docs/SECURITY_IMPLEMENTATION.md](docs/SECURITY_IMPLEMENTATION.md)** - Implementation details untuk security features

### 🚀 API Features & Optimization
- **[docs/API_FEATURES.md](docs/API_FEATURES.md)** - Advanced API features: pagination, bulk ops, export, search, filtering
- **[docs/API_FEATURES_QUICK_REF.md](docs/API_FEATURES_QUICK_REF.md)** - Quick reference untuk API features
- **[docs/PERFORMANCE.md](docs/PERFORMANCE.md)** - Performance optimization: indexing, timeout, graceful shutdown
- **[docs/CACHING.md](docs/CACHING.md)** - Redis caching system dengan strategies dan best practices

### 💾 Data & Storage
- **[docs/FILE_MANAGEMENT.md](docs/FILE_MANAGEMENT.md)** - File upload, validation, cloud storage (S3/MinIO), image processing
- **[docs/CLOUD_STORAGE_EXAMPLES.md](docs/CLOUD_STORAGE_EXAMPLES.md)** - Cloud storage implementation examples
- **[docs/MINIO_SETUP.md](docs/MINIO_SETUP.md)** - Setup MinIO sebagai local S3-compatible storage

### 🔄 Background Jobs & Real-time
- **[docs/BACKGROUND_JOBS.md](docs/BACKGROUND_JOBS.md)** - Background jobs dengan Asynq: queue, scheduler, email queue
- **[docs/BACKGROUND_JOBS_QUICK_REF.md](docs/BACKGROUND_JOBS_QUICK_REF.md)** - Quick reference untuk background jobs
- **[docs/SSE.md](docs/SSE.md)** - Server-Sent Events untuk real-time notifications
- **[docs/EMAIL.md](docs/EMAIL.md)** - Email system: templates, SMTP providers, async sending

### 📊 Monitoring & Testing
- **[docs/LOGGING_MONITORING.md](docs/LOGGING_MONITORING.md)** - Structured logging dengan Zap dan Sentry integration
- **[docs/METRICS_MONITORING.md](docs/METRICS_MONITORING.md)** - Metrics endpoint dan health checks
- **[docs/TESTING.md](docs/TESTING.md)** - Testing guidelines dan best practices
- **[docs/TEST_COVERAGE.md](docs/TEST_COVERAGE.md)** - Test coverage reports dan analysis

### 🚢 Deployment & DevOps
- **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)** - Deployment guide untuk Docker Compose
- **[docs/DEVOPS.md](docs/DEVOPS.md)** - DevOps setup: CI/CD, Makefile, pre-commit hooks, Docker optimization

### 📖 Additional Resources
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contributing guidelines dan coding standards
- **[LICENSE](LICENSE)** - MIT License
- **[.cursorrules](.cursorrules)** - AI assistant rules dan project conventions

---

## Prioritas Implementasi (Rekomendasi)

### High Priority
1. **API Documentation (Swagger)** - Auto-generated API docs
2. **Graceful Shutdown** - Zero-downtime deployment support
3. **Database Indexing** - Query performance optimization

### Medium Priority
1. **CI/CD Pipeline** - Automated testing & deployment
2. **Makefile** - Development automation
3. **User Profile Management** - Basic user profile features

### Low Priority
1. **API Versioning** - Backward compatibility management
2. **Request Timeout** - Request timeout configuration