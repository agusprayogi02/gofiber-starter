# Starter Template Go Fiber

Backend Golang With Fiber Framework

## Getting Started

### Prerequisites

Sebelum menjalankan project ini, pastikan sistem Anda telah memiliki:

#### Required
- **Go 1.20 or higher** - [Download Go](https://golang.org/dl/)
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
go run main.go migrate up

# Atau menggunakan custom migration command jika ada
# go run scripts/migrate.go up
```

#### 8. (Optional) Seed Database
```bash
# Populate database dengan sample data untuk development
go run scripts/seeder.go
```

#### 9. Run Application

**Development Mode (dengan live reload):**
```bash
air
```

**Standard Mode:**
```bash
go run main.go
```

**Build dan Run Binary:**
```bash
# Build binary
go build -o bin/app

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
// entity/user.go
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
docker build -t starter-gofiber .

# Run with docker-compose
docker-compose up -d

# Check logs
docker-compose logs -f app
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

```
starter-gofiber/
├── assets/                      # Static assets dan files
│   ├── certs/                   # SSL certificates dan private keys (JWT)
│   ├── docker/                  # Docker configurations
│   │   ├── dev/                 # Development docker setup
│   │   └── prod/                # Production docker setup
│   └── rbac/                    # Casbin RBAC policy files
│       ├── model.conf           # RBAC model definition
│       └── policy.csv           # Role-permission mappings
│
├── backups/                     # Database backups (auto-generated)
│
├── config/                      # Application configurations
│   ├── app.go                   # App config (port, name, version)
│   ├── config.go                # Main config loader
│   ├── database.go              # Database connection config
│   ├── permission.go            # Casbin permission init
│   ├── storage.go               # S3/MinIO storage config
│   └── timezone.go              # Timezone configuration
│
├── docs/                        # Documentation files
│   ├── API_AUTH.md              # Authentication API docs
│   ├── API_FEATURES.md          # Advanced API features
│   ├── ARCHITECTURE.md          # Architecture overview
│   ├── BACKGROUND_JOBS.md       # Background jobs & Asynq
│   ├── CACHING.md               # Redis caching guide
│   ├── DATABASE.md              # Database management
│   ├── DEPLOYMENT.md            # Deployment guide
│   ├── EMAIL.md                 # Email system docs
│   ├── FILE_MANAGEMENT.md       # File upload & storage
│   ├── LOGGING_MONITORING.md    # Logging & monitoring
│   ├── METRICS_MONITORING.md    # Metrics & health checks
│   ├── SECURITY.md              # Security features
│   ├── SSE.md                   # Server-Sent Events
│   ├── TESTING.md               # Testing guide
│   └── TEST_COVERAGE.md         # Test coverage reports
│
├── dto/                         # Data Transfer Objects
│   ├── paginate.go              # Pagination DTO
│   ├── post.go                  # Post-related DTOs
│   ├── response.go              # Standard API responses
│   ├── token.go                 # Token DTOs
│   └── user.go                  # User-related DTOs
│
├── entity/                      # Database models (GORM)
│   ├── post.go                  # Post model
│   └── user.go                  # User model + enums
│
├── handler/                     # HTTP request handlers (controllers)
│   ├── auth.go                  # Authentication endpoints
│   └── post.go                  # Post CRUD endpoints
│
├── helper/                      # Utility functions & helpers
│   ├── database.go              # Database helpers (pagination)
│   ├── date.go                  # Date/time utilities
│   ├── email.go                 # Email service
│   ├── email_config.go          # Email configuration
│   ├── error.go                 # Custom error types
│   ├── hash.go                  # Hashing utilities
│   ├── jobs.go                  # Asynq job helpers
│   ├── password.go              # Password utilities
│   ├── response.go              # Response formatters
│   ├── token.go                 # JWT token utilities
│   ├── upload.go                # File upload helpers
│   └── validation.go            # Input validation
│
├── jobs/                        # Background job handlers
│   ├── email.go                 # Email job handlers
│   ├── handlers.go              # Job handler implementations
│   └── jobs.go                  # Job definitions
│
├── middleware/                  # HTTP middlewares
│   ├── auth.go                  # JWT authentication
│   └── authz.go                 # Casbin authorization
│
├── migrations/                  # Database migrations
│   └── *.sql                    # Migration files
│
├── mocks/                       # Mock objects untuk testing
│
├── postman/                     # Postman collections
│   └── *.json                   # API test collections
│
├── public/                      # Public static files
│
├── repository/                  # Database repository layer
│   ├── post.go                  # Post repository
│   ├── repository.go            # Base repository interface
│   └── user.go                  # User repository
│
├── router/                      # Route definitions
│   ├── auth.go                  # Auth routes
│   ├── post.go                  # Post routes
│   └── router.go                # Main router setup
│
├── scripts/                     # Utility scripts
│   ├── seeder.go                # Database seeder
│   └── migrate.go               # Migration runner
│
├── service/                     # Business logic layer
│   ├── auth.go                  # Authentication service
│   └── post.go                  # Post service
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
│   └── constant.go              # App constants
│
├── .air.toml                    # Air live reload config
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
├── CONTRIBUTING.md              # Contributing guidelines
├── docker-compose.yml           # Docker Compose (main)
├── docker-compose.minio.yml     # MinIO S3 storage setup
├── Dockerfile                   # Docker build instructions
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── LICENSE                      # MIT License
├── main.go                      # Application entry point
├── MINIO_SETUP.md              # MinIO setup guide
└── README.MD                    # This file
```

### Architecture Layers

Project ini menggunakan **Clean Architecture** dengan layer separation:

1. **Entity Layer** (`entity/`) - Domain models
2. **Repository Layer** (`repository/`) - Data access
3. **Service Layer** (`service/`) - Business logic
4. **Handler Layer** (`handler/`) - HTTP handlers
5. **Router Layer** (`router/`) - Route definitions

**Data Flow:**
```
Client Request → Router → Middleware → Handler → Service → Repository → Database
                                                    ↓
                                            Helper/Jobs/Config
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
Lihat [docs/API_AUTH.md](docs/API_AUTH.md) untuk dokumentasi lengkap endpoint authentication.

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
- Lihat [Setup Minio](MINIO_SETUP.md) untuk panduan setup MinIO sebagai local S3-compatible storage

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
- [ ] **WebSocket Support** - Real-time bidirectional communication (belum diperlukan)
- [ ] **Pub/Sub System** - Event-driven architecture (belum diperlukan)

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
- [ ] **Makefile** - Common tasks automation
- [ ] **GitHub Actions/GitLab CI** - Automated testing & deployment
- [ ] **Docker Optimization** - Multi-stage build, smaller image size
- [ ] **Docker Compose for Testing** - Separate compose untuk testing
- [ ] **Kubernetes Manifests** - K8s deployment files
- [ ] **Helm Charts** - Kubernetes package manager
- [ ] **Pre-commit Hooks** - Automated checks sebelum commit
- [ ] **Code Linting** - golangci-lint integration
- [ ] **Code Formatting** - Automated formatting dengan gofmt/goimports

### 14. Performance
- [ ] **Database Indexing** - Proper indexes untuk queries
- [ ] **Query Optimization** - Analyze dan optimize slow queries
- [ ] **Connection Pooling** - Optimize database connections
- [ ] **Lazy Loading** - Load relations on demand
- [ ] **Response Pagination** - Limit response size
- [ ] **Request Timeout** - Timeout configuration untuk all requests

### 15. User Management
- [ ] **User Profile** - Get/Update user profile
- [ ] **Avatar Upload** - Profile picture management
- [ ] **User Preferences** - Settings/preferences per user
- [ ] **User Activity Log** - Track user activities
- [ ] **User Roles & Permissions CRUD** - Manage roles via API
- [ ] **User Impersonation** - Admin impersonate user (for support)
- [ ] **Account Deactivation** - Soft delete user accounts

### 16. Notification System
- [ ] **In-app Notifications** - Notification center
- [ ] **Push Notifications** - Firebase/OneSignal integration
- [ ] **Email Notifications** - Automated email notifications
- [ ] **SMS Notifications** - Twilio/Nexmo integration
- [ ] **Notification Preferences** - User can choose notification channels

### 17. Admin Panel Features
- [ ] **Admin Dashboard** - Statistics dan metrics
- [ ] **User Management** - CRUD users from admin
- [ ] **Content Moderation** - Approve/reject user content
- [ ] **System Settings** - Configurable settings via admin
- [ ] **Audit Logs Viewer** - View system audit logs

### 18. Internationalization (i18n)
- [ ] **Multi-language Support** - Support multiple languages
- [ ] **Translation Management** - Easy to add new translations
- [ ] **Locale Detection** - Auto detect user locale
- [ ] **Date/Time Localization** - Format based on locale

### 19. Others
- [ ] **Feature Flags** - Toggle features on/off without deployment
- [ ] **A/B Testing Support** - Framework untuk A/B testing
- [ ] **Analytics Integration** - Google Analytics, Mixpanel, dll
- [ ] **Error Reporting** - Better error messages dengan codes
- [ ] **Request Validation** - Comprehensive input validation
- [ ] **GDPR Compliance** - Data export/delete for privacy compliance
- [ ] **Graceful Shutdown** - Proper shutdown handling

## Prioritas Implementasi (Rekomendasi)

### High Priority (Harus Ada)
1. Unit & Integration Tests
2. API Documentation (Swagger)
3. Refresh Token & Logout
4. Database Migration System
5. Structured Logging
6. Health Check & Metrics
7. Email Verification & Password Reset
8. API Versioning

### Medium Priority (Sangat Direkomendasikan)
1. Redis Caching
2. Background Job Queue
3. File Validation & Cloud Storage
4. Database Seeder & Backup
5. CI/CD Pipeline
6. Error Tracking (Sentry)
7. User Profile Management
8. Advanced Search & Filtering

### Low Priority (Nice to Have)
1. WebSocket/Real-time features
2. GraphQL Support
3. Multi-language Support
4. Social Login
5. Admin Panel Features
6. Feature Flags
7. A/B Testing