# Makefile Commands Reference

Quick reference untuk semua available commands di Makefile.

## 📋 View All Commands

```bash
make help
```

## 🗄️ Database & Migration (Atlas)

### Setup

```bash
# Install Atlas CLI
make atlas-install
```

### Generate Migrations

```bash
# Auto-generate migration from GORM schema changes
make atlas-diff

# Generate with custom name
make atlas-diff-name NAME=add_user_bio
```

### Apply Migrations

```bash
# Apply to development
make atlas-apply

# Apply to production (with confirmation)
make atlas-apply-prod
```

### Check Status

```bash
# Development status
make atlas-status

# Production status
make atlas-status-prod
```

### Validate & Test

```bash
# Validate migration files
make atlas-validate

# Test with dry-run
make atlas-test

# Lint migrations
make atlas-lint

# Inspect database schema
make atlas-inspect
```

### Workflows

```bash
# Full workflow: diff + validate + apply
make atlas-migrate

# Clean migrations (DANGEROUS!)
make atlas-clean
```

## 🏗️ Build Commands

```bash
# Build both API and Worker
make build

# Build API only
make build-api

# Build Worker only
make build-worker
```

## 🚀 Run Commands

```bash
# Run API server
make run

# Run with Air (live reload)
make run-air

# Run Worker
make run-worker
```

## 🧪 Test Commands

```bash
# Run all tests
make test

# Run tests (short mode)
make test-short

# Run with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

## 🔍 Code Quality

```bash
# Format code
make fmt

# Check format without modifying
make fmt-check

# Run linter
make lint

# Run go vet
make vet

# Tidy dependencies
make mod-tidy

# Download dependencies
make mod-download
```

## 🐳 Docker Commands

```bash
# Build Docker image
make docker-build

# Build optimized image
make docker-build-optimized

# Run container
make docker-run

# Start services with docker-compose
make docker-up

# Stop services
make docker-down

# View logs
make docker-logs

# Clean Docker resources
make docker-clean
```

## 🛠️ Development Setup

```bash
# Full development setup (installs tools, Atlas, downloads deps)
make dev-setup

# Install development tools only
make install-tools
```

## 🧹 Cleanup Commands

```bash
# Clean build artifacts
make clean

# Clean everything including test cache
make clean-all
```

## 🤖 CI/CD Helpers

```bash
# Run all CI checks (format, vet, test)
make ci-test

# Build for CI
make ci-build
```

## 💡 Common Workflows

### First Time Setup

```bash
make dev-setup           # Install tools and dependencies
# Configure .env file
# Set DB_GEN=true for development
make run-air             # Start development server
```

### Adding New Feature with Migration

```bash
# 1. Create GORM entity
# 2. Add to GetModelsForMigration() in internal/config/database.go
# 3. Generate migration
make atlas-diff-name NAME=add_categories_table

# 4. Review generated SQL
cat migrations/[timestamp]_add_categories_table.sql

# 5. Apply migration
make atlas-apply

# 6. Verify
make atlas-status
```

### Before Committing Code

```bash
make fmt                 # Format code
make vet                 # Check for issues
make test                # Run tests
make lint                # Run linter (if installed)
```

### Production Deployment

```bash
# On production server
git pull origin main
make atlas-apply-prod    # Apply migrations (with confirmation)
make build               # Build binaries
# Restart service
```

## 🔗 Related Documentation

- **Atlas Migration Guide**: `MIGRATION_GUIDE.md`
- **Database Documentation**: `docs/DATABASE.md`
- **Migration README**: `migrations/README.md`

---

**Pro Tip**: Run `make help` anytime untuk melihat semua available commands dengan deskripsi!
