# DevOps & CI/CD Guide

Dokumentasi lengkap untuk setup DevOps dan CI/CD pipeline.

## Daftar Isi

- [Makefile](#makefile)
- [GitHub Actions CI/CD](#github-actions-cicd)
- [Docker Optimization](#docker-optimization)
- [Pre-commit Hooks](#pre-commit-hooks)
- [Development Workflow](#development-workflow)

---

## Makefile

Makefile menyediakan automation untuk common tasks development dan deployment.

### Installation

Tidak perlu install, Make sudah tersedia di sebagian besar sistem Unix/Linux.

### Available Commands

**Build Commands:**
```bash
make build              # Build API dan Worker binaries
make build-api          # Build API binary saja
make build-worker       # Build Worker binary saja
```

**Run Commands:**
```bash
make run                # Run API server (development)
make run-worker         # Run Worker server
make run-air            # Run dengan Air (live reload)
```

**Test Commands:**
```bash
make test               # Run semua tests
make test-short         # Run tests dalam short mode
make test-coverage      # Run tests dengan coverage report
make test-coverage-html # Generate HTML coverage report
```

**Database Commands:**
```bash
make migrate-up         # Run database migrations
make migrate-down       # Rollback last migration
make migrate-create NAME=create_users_table  # Create new migration
make migrate-status     # Show migration status
```

**Code Quality:**
```bash
make lint               # Run golangci-lint
make fmt                # Format code dengan gofmt
make fmt-check          # Check formatting tanpa modify
make vet                # Run go vet
make mod-tidy           # Tidy go.mod dan go.sum
```

**Docker Commands:**
```bash
make docker-build              # Build Docker image
make docker-build-optimized    # Build optimized Docker image
make docker-run                # Run Docker container
make docker-up                 # Start services dengan docker-compose
make docker-down               # Stop services
make docker-logs               # Show docker-compose logs
make docker-clean              # Clean Docker resources
```

**Installation:**
```bash
make install-tools      # Install development tools (Air, golangci-lint, migrate)
```

**Cleanup:**
```bash
make clean              # Clean build artifacts
make clean-all          # Clean everything termasuk test cache
```

**CI/CD Helpers:**
```bash
make ci-test            # Run all CI checks (format, vet, test)
make ci-build           # Build for CI
```

**Development Setup:**
```bash
make dev-setup          # Setup development environment lengkap
```

### Examples

**Setup Development Environment:**
```bash
make dev-setup          # Install tools dan download dependencies
cp .env.example .env    # Setup environment variables
make migrate-up         # Run migrations
make run-air            # Start development server
```

**Before Committing:**
```bash
make fmt                # Format code
make lint               # Check linting
make test               # Run tests
```

**Production Build:**
```bash
make build              # Build binaries
make docker-build-optimized  # Build optimized Docker image
```

---

## GitHub Actions CI/CD

CI/CD pipeline otomatis untuk testing, linting, dan building.

### Setup

1. **Push workflow file** ke repository:
   ```bash
   git add .github/workflows/ci.yml
   git commit -m "Add CI/CD pipeline"
   git push
   ```

2. **Configure secrets** (jika diperlukan):
   - Go to repository Settings → Secrets
   - Add secrets untuk deployment (jika ada)

### Workflow Jobs

Pipeline terdiri dari 4 jobs yang berjalan secara sequential:

1. **Lint Job**
   - Checkout code
   - Setup Go environment
   - Run golangci-lint
   - Cache Go modules

2. **Test Job**
   - Setup PostgreSQL dan Redis services
   - Run tests dengan coverage
   - Check coverage threshold (60%)
   - Upload coverage ke Codecov

3. **Build Job**
   - Build API binary
   - Build Worker binary
   - Upload binaries sebagai artifacts

4. **Docker Build Job**
   - Build optimized Docker image
   - Cache Docker layers untuk faster builds

### Trigger Events

Workflow akan trigger pada:
- Push ke `main` atau `develop` branch
- Pull request ke `main` atau `develop` branch

### Environment Variables

Workflow menggunakan environment variables:
- `GO_VERSION`: Go version (default: 1.25)
- `COVERAGE_THRESHOLD`: Minimum coverage percentage (default: 60)

### Customization

**Change Go Version:**
```yaml
env:
  GO_VERSION: '1.26'  # Update version
```

**Change Coverage Threshold:**
```yaml
env:
  COVERAGE_THRESHOLD: 70  # Increase threshold
```

**Add Deployment Step:**
```yaml
deploy:
  name: Deploy
  runs-on: ubuntu-latest
  needs: docker-build
  if: github.ref == 'refs/heads/main'
  steps:
    - name: Deploy to production
      # Add deployment steps here
```

---

## Docker Optimization

Optimized Dockerfile dengan multi-stage build untuk smaller image size.

### Features

- **Multi-stage build**: Separate build dan runtime stages
- **Smaller image**: Hanya include runtime dependencies
- **Security**: Non-root user untuk running container
- **Health check**: Built-in health check endpoint
- **Static binaries**: Fully static binaries untuk portability

### Build Optimized Image

```bash
# Build optimized image
make docker-build-optimized

# Atau langsung dengan docker
docker build -f assets/docker/Dockerfile.optimized -t starter-gofiber:latest .
```

### Image Size Comparison

- **Standard Dockerfile**: ~300-400 MB
- **Optimized Dockerfile**: ~20-30 MB (90% reduction!)

### Usage

**Run Container:**
```bash
docker run -p 3000:3000 --env-file .env starter-gofiber:latest
```

**Run with docker-compose:**
```yaml
services:
  app:
    build:
      context: .
      dockerfile: assets/docker/Dockerfile.optimized
    ports:
      - "3000:3000"
    env_file:
      - .env
```

### Optimization Techniques

1. **Multi-stage build**: Build di satu stage, copy binaries ke minimal runtime stage
2. **Alpine Linux**: Minimal base image (~5MB)
3. **Static linking**: CGO_ENABLED=0 untuk fully static binaries
4. **Strip symbols**: `-ldflags='-w -s'` untuk remove debug symbols
5. **Layer caching**: Copy go.mod/go.sum first untuk better caching

---

## Pre-commit Hooks

Automated checks sebelum commit untuk menjaga code quality.

### Installation

**Install pre-commit:**
```bash
# Install pre-commit (Python required)
pip install pre-commit

# Install hooks
pre-commit install

# Install golangci-lint (required)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Available Hooks

**General File Checks:**
- Trailing whitespace removal
- End of file fixer
- YAML/JSON validation
- Large file detection
- Merge conflict detection
- Private key detection

**Go Specific:**
- `go-fmt`: Format Go code
- `go-vet`: Run go vet
- `go-mod-tidy`: Tidy go.mod
- `golangci-lint`: Run linter

**Other:**
- Makefile linting (checkmake)
- Dockerfile linting (hadolint)
- Markdown linting

### Usage

**Run manually:**
```bash
# Run on all files
pre-commit run --all-files

# Run on staged files only (default)
pre-commit run
```

**Skip hooks (not recommended):**
```bash
git commit --no-verify -m "message"
```

### Configuration

Edit `.pre-commit-config.yaml` untuk customize hooks:

```yaml
# Disable specific hook
- id: golangci-lint
  exclude: ^(tests/|mocks/)  # Skip linting on tests/mocks
```

---

## Development Workflow

### Initial Setup

```bash
# 1. Clone repository
git clone <repo-url>
cd starter-gofiber

# 2. Setup development environment
make dev-setup

# 3. Setup environment variables
cp .env.example .env
# Edit .env dengan konfigurasi yang sesuai

# 4. Setup database
make migrate-up

# 5. Start development server
make run-air
```

### Daily Workflow

**Before Starting Work:**
```bash
git pull origin main
make mod-download  # Update dependencies jika ada perubahan
```

**During Development:**
```bash
make run-air  # Start dengan live reload
# Make changes...
# Server akan auto-reload
```

**Before Committing:**
```bash
make fmt          # Format code
make lint          # Check linting
make test          # Run tests
make test-coverage # Check coverage
```

**Commit:**
```bash
git add .
git commit -m "feat: add new feature"
# Pre-commit hooks akan run otomatis
```

**Push:**
```bash
git push origin feature-branch
# CI/CD pipeline akan run otomatis
```

### Code Review Checklist

Sebelum membuat PR, pastikan:
- [ ] Code sudah di-format (`make fmt`)
- [ ] Linting passed (`make lint`)
- [ ] Tests passed (`make test`)
- [ ] Coverage meets threshold (`make test-coverage`)
- [ ] Documentation updated (jika diperlukan)
- [ ] Migration files created (jika ada perubahan DB)

---

## Troubleshooting

### Makefile Issues

**Command not found:**
```bash
# Install Make (Ubuntu/Debian)
sudo apt-get install make

# Install Make (macOS)
xcode-select --install
```

**Go command not found:**
```bash
# Install Go
# https://golang.org/dl/
```

### Pre-commit Issues

**Pre-commit not found:**
```bash
pip install pre-commit
```

**golangci-lint not found:**
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Hooks too slow:**
```bash
# Skip golangci-lint untuk faster commits (not recommended)
# Edit .pre-commit-config.yaml dan comment out golangci-lint hook
```

### CI/CD Issues

**Tests failing:**
- Check database connection settings
- Verify Redis is running
- Check environment variables

**Coverage below threshold:**
- Write more tests
- Increase coverage threshold di workflow (tidak recommended)
- Check coverage report untuk identify uncovered code

**Docker build failing:**
- Check Dockerfile syntax
- Verify all dependencies available
- Check build logs untuk detailed errors

---

## Best Practices

1. **Always run `make fmt` before committing**
2. **Run `make test` sebelum push**
3. **Keep coverage above threshold**
4. **Use meaningful commit messages**
5. **Review CI/CD logs jika build fails**
6. **Update documentation jika ada perubahan**
7. **Test Docker image locally sebelum deploy**

---

## Additional Resources

- [Makefile Documentation](https://www.gnu.org/software/make/manual/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Pre-commit Documentation](https://pre-commit.com/)
- [golangci-lint Documentation](https://golangci-lint.run/)

