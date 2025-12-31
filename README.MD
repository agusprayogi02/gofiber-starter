# Starter Template Go Fiber

Backend Golang With Fiber Framework

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Docker (optional)
- Open ssl -> untuk generate token jwt

### Installation
1. Install `air` for live reload,
```bash
go install github.com/air-verse/air@latest
```
2. Copy `.env.example` to `.env` and fill the values.
3. Membuat certificate menggunakan open ssl
	- Masuk ke folder `assets/certs` 
	- buat private key
		```bash
		openssl genpkey -algorithm RSA -out certificate.pem -pkeyopt rsa_keygen_bits:4096
		```
	- Masukkan kedalam file `.env` dengan format:
		```bash
		LOCATION_CERT="assets/certs/certificate.pem"
		```
3. Change config on .env file, like database connection string.
4. For `Sql Server` change Enum Data Type on `entity/user.go` like: 
```go
type User struct {
	ID       uint     `gorm:"primaryKey;autoIncrement"`
	Name     string   `gorm:"type:varchar(200);not null"`
	Email    string   `gorm:"type:varchar(200);uniqueIndex;not null"`
	Password string   `gorm:"type:varchar(150);not null"`
	Role     UserRole `gorm:"type:varchar(10);default:user"` // this
	gorm.Model
}
```
5. Run `air` to start the server with live reload.
6. if you want to run in the docker container, run `docker-compose up -d` and then `docker ps` to check the container is running.

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
- [x] **Structured Logging** - Implementasi structured logging (zap) menggantikan log.Print
- [x] **Request ID Tracking** - Trace request dengan unique ID (UUID)
- [x] **Access Logs** - Logging untuk semua incoming requests dengan detail lengkap
- [x] **Error Tracking** - Integration dengan Sentry untuk production error tracking
- [x] **Performance Monitoring** - Prometheus metrics untuk monitoring performa
- [x] **Health Check Endpoint** - `/health`, `/health/ready`, `/health/live` untuk orchestration
- [x] **Metrics Endpoint** - Prometheus metrics di `/metrics`
- [x] **Database Query Logging** - Logging slow queries dengan custom GORM logger

### Dokumentasi
Lihat [docs/LOGGING_MONITORING.md](docs/LOGGING_MONITORING.md) untuk dokumentasi lengkap logging dan monitoring

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
- [ ] **Redis Integration** - Caching layer dengan Redis
- [ ] **Cache Middleware** - Response caching untuk GET requests
- [ ] **Query Result Caching** - Cache database query results
- [ ] **Cache Invalidation Strategy** - Proper cache invalidation mechanism

### 7. Database
- [ ] **Database Migration System** - Proper migration tool (golang-migrate, gorm/migrator)
- [ ] **Seeder** - Database seeder untuk development/testing
- [ ] **Database Backup Scripts** - Automated backup mechanism
- [ ] **Connection Pooling Configuration** - Optimal connection pool settings
- [ ] **Read Replica Support** - Load balancing untuk read operations
- [ ] **Soft Delete Support** - Soft delete implementation untuk semua entities
- [ ] **Audit Log** - Track semua perubahan data (who, when, what)

### 8. File Management
- [ ] **File Validation** - Validate file type, size, dan content
- [ ] **Multiple File Upload** - Support upload multiple files sekaligus
- [ ] **Cloud Storage Integration** - AWS S3/Google Cloud Storage integration
- [ ] **Image Processing** - Resize, compress, thumbnail generation
- [ ] **CDN Integration** - Serve static files via CDN
- [ ] **File Versioning** - Keep track of file versions

### 9. API Features
- [ ] **Response Compression** - Gzip compression untuk response
- [ ] **Pagination Cursor-based** - Alternative untuk offset pagination
- [ ] **GraphQL Support** - GraphQL endpoint (optional)
- [ ] **Webhook System** - Trigger webhooks untuk events tertentu
- [ ] **Bulk Operations** - Bulk create/update/delete endpoints
- [ ] **Export Data** - Export ke CSV/Excel/PDF
- [ ] **Search & Filter** - Advanced search dengan multiple filters
- [ ] **Sorting** - Multiple field sorting support

### 10. Real-time Features
- [ ] **WebSocket Support** - Real-time communication
- [ ] **Server-Sent Events (SSE)** - Push notifications ke client
- [ ] **Pub/Sub System** - Event-driven architecture

### 11. Background Jobs
- [ ] **Job Queue** - Background job processing (Asynq, Machinery)
- [ ] **Scheduled Tasks** - Cron jobs untuk periodic tasks
- [ ] **Email Queue** - Async email sending
- [ ] **Retry Mechanism** - Auto retry failed jobs
- [ ] **Job Dashboard** - Monitor job status dan history

### 12. Email System
- [ ] **Email Templates** - HTML email templates
- [ ] **Email Service Integration** - SendGrid/Mailgun/AWS SES
- [ ] **Email Queue** - Async email processing
- [ ] **Email Verification** - Template untuk verify email
- [ ] **Password Reset Email** - Template untuk reset password
- [ ] **Newsletter System** - Bulk email sending

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