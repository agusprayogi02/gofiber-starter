# Documentation Index

Comprehensive documentation for Go Fiber Starter Kit.

## 📖 Quick Links

### 🚀 Getting Started
- **[../README.MD](../README.MD)** - Main project README with setup instructions
- **[MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)** - Quick start untuk Atlas database migrations (Bahasa Indonesia)
- **[MAKEFILE_COMMANDS.md](MAKEFILE_COMMANDS.md)** - Semua Makefile commands dengan contoh

### 🏗️ Architecture & Database
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Clean Architecture pattern, layers, dan design patterns
- **[DATABASE.md](DATABASE.md)** - Database management lengkap: Atlas migrations, seeding, backup, optimization

### 🔐 Authentication & Security
- **[API_AUTH.md](API_AUTH.md)** - Authentication API: register, login, profile, JWT, refresh token
- **[SECURITY.md](SECURITY.md)** - Security features overview: rate limiting, CSRF, XSS, encryption
- **[SECURITY_IMPLEMENTATION.md](SECURITY_IMPLEMENTATION.md)** - Implementation details untuk security features

### 🚀 API Features
- **[API_FEATURES.md](API_FEATURES.md)** - Advanced features: pagination, bulk operations, export, search, filtering
- **[API_FEATURES_QUICK_REF.md](API_FEATURES_QUICK_REF.md)** - Quick reference untuk API features
- **[PERFORMANCE.md](PERFORMANCE.md)** - Performance optimization: indexing, connection pooling, graceful shutdown

### 💾 Data Storage & Files
- **[FILE_MANAGEMENT.md](FILE_MANAGEMENT.md)** - File upload, validation, cloud storage (S3/MinIO), image processing, versioning
- **[CLOUD_STORAGE_EXAMPLES.md](CLOUD_STORAGE_EXAMPLES.md)** - Cloud storage implementation examples
- **[MINIO_SETUP.md](MINIO_SETUP.md)** - Setup MinIO sebagai local S3-compatible storage
- **[CACHING.md](CACHING.md)** - Redis caching system: strategies, middleware, invalidation

### 🔄 Background Jobs & Real-time
- **[BACKGROUND_JOBS.md](BACKGROUND_JOBS.md)** - Background jobs dengan Asynq: queue, scheduler, email queue, retry
- **[BACKGROUND_JOBS_QUICK_REF.md](BACKGROUND_JOBS_QUICK_REF.md)** - Quick reference untuk background jobs
- **[SSE.md](SSE.md)** - Server-Sent Events untuk real-time push notifications
- **[EMAIL.md](EMAIL.md)** - Email system: templates, SMTP providers, async sending, queue

### 📊 Monitoring & Testing
- **[LOGGING_MONITORING.md](LOGGING_MONITORING.md)** - Structured logging dengan Zap, Sentry integration, request tracking
- **[METRICS_MONITORING.md](METRICS_MONITORING.md)** - Metrics endpoint, health checks, lightweight monitoring
- **[TESTING.md](TESTING.md)** - Testing guidelines, test suite pattern, integration tests
- **[TEST_COVERAGE.md](TEST_COVERAGE.md)** - Test coverage reports, analysis, dan improvement tips

### 🚢 Deployment & DevOps
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Deployment guide: Docker Compose, environment setup, production checklist
- **[DEVOPS.md](DEVOPS.md)** - DevOps setup: CI/CD, Makefile, pre-commit hooks, Docker optimization

## 📚 By Topic

### For New Developers
1. Start with [../README.MD](../README.MD) untuk setup project
2. Baca [ARCHITECTURE.md](ARCHITECTURE.md) untuk understand struktur
3. Follow [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) untuk database setup
4. Explore [API_AUTH.md](API_AUTH.md) untuk understand API endpoints
5. Check [MAKEFILE_COMMANDS.md](MAKEFILE_COMMANDS.md) untuk common tasks

### For API Development
- [API_AUTH.md](API_AUTH.md) - Authentication endpoints
- [API_FEATURES.md](API_FEATURES.md) - Advanced API features
- [API_FEATURES_QUICK_REF.md](API_FEATURES_QUICK_REF.md) - Quick reference
- [CACHING.md](CACHING.md) - Response caching
- [PERFORMANCE.md](PERFORMANCE.md) - API optimization

### For Database Work
- [DATABASE.md](DATABASE.md) - Complete database management
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - Quick start guide (Bahasa Indonesia)
- [PERFORMANCE.md](PERFORMANCE.md) - Database optimization

### For Security Implementation
- [SECURITY.md](SECURITY.md) - Security features overview
- [SECURITY_IMPLEMENTATION.md](SECURITY_IMPLEMENTATION.md) - Implementation details
- [API_AUTH.md](API_AUTH.md) - Authentication security

### For File & Storage
- [FILE_MANAGEMENT.md](FILE_MANAGEMENT.md) - Complete file handling
- [CLOUD_STORAGE_EXAMPLES.md](CLOUD_STORAGE_EXAMPLES.md) - S3 examples
- [MINIO_SETUP.md](MINIO_SETUP.md) - Local S3 setup

### For Background Processing
- [BACKGROUND_JOBS.md](BACKGROUND_JOBS.md) - Complete Asynq guide
- [BACKGROUND_JOBS_QUICK_REF.md](BACKGROUND_JOBS_QUICK_REF.md) - Quick reference
- [EMAIL.md](EMAIL.md) - Email queue system

### For Real-time Features
- [SSE.md](SSE.md) - Server-Sent Events
- [BACKGROUND_JOBS.md](BACKGROUND_JOBS.md) - Background job notifications

### For Production Deployment
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deployment guide
- [DEVOPS.md](DEVOPS.md) - CI/CD setup
- [LOGGING_MONITORING.md](LOGGING_MONITORING.md) - Production monitoring
- [METRICS_MONITORING.md](METRICS_MONITORING.md) - Health checks

### For Testing
- [TESTING.md](TESTING.md) - Testing guidelines
- [TEST_COVERAGE.md](TEST_COVERAGE.md) - Coverage reports

## 🔧 Additional Resources

- **[../cmd/README.md](../cmd/README.md)** - CMD folder conventions (main.go requirement)
- **[../cmd/atlas/README.md](../cmd/atlas/README.md)** - Atlas loader documentation
- **[../CONTRIBUTING.md](../CONTRIBUTING.md)** - Contributing guidelines
- **[../.cursorrules](../.cursorrules)** - AI assistant rules (project conventions)

## 📝 Documentation Standards

When creating new documentation:

1. **Clear Title** - Descriptive and specific
2. **Table of Contents** - For docs > 100 lines
3. **Code Examples** - Working, tested examples
4. **Use Cases** - Real-world scenarios
5. **Troubleshooting** - Common issues & solutions
6. **Best Practices** - Dos and don'ts
7. **References** - Links to related docs

## 🆘 Getting Help

If you can't find what you need:

1. Check [README.MD](../README.MD) first
2. Search in relevant doc by topic (see above)
3. Check [.cursorrules](../.cursorrules) for coding conventions
4. Look at existing code for patterns
5. Check external resources (links in each doc)

## 📊 Documentation Stats

- **Total Docs**: 24 files
- **Categories**: 7 main categories
- **Languages**: English (code), Bahasa Indonesia (MIGRATION_GUIDE)
- **Last Updated**: January 2026

---

**Happy Coding! 🚀**

For project setup, start with [../README.MD](../README.MD)
