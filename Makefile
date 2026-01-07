.PHONY: help build run test test-coverage clean migrate migrate-up migrate-down migrate-create lint fmt vet docker-build docker-run docker-up docker-down install-tools

# Variables
APP_NAME=starter-gofiber
API_BINARY=bin/api
WORKER_BINARY=bin/worker
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_REGISTRY?=
COVERAGE_THRESHOLD=60

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# Build commands
build: ## Build API and Worker binaries
	@echo "$(GREEN)Building binaries...$(NC)"
	@mkdir -p bin
	@go build -o $(API_BINARY) ./cmd/api/main.go
	@go build -o $(WORKER_BINARY) ./cmd/worker/main.go
	@echo "$(GREEN)✓ Build complete$(NC)"

build-api: ## Build API binary only
	@echo "$(GREEN)Building API binary...$(NC)"
	@mkdir -p bin
	@go build -o $(API_BINARY) ./cmd/api/main.go
	@echo "$(GREEN)✓ API build complete$(NC)"

build-worker: ## Build Worker binary only
	@echo "$(GREEN)Building Worker binary...$(NC)"
	@mkdir -p bin
	@go build -o $(WORKER_BINARY) ./cmd/worker/main.go
	@echo "$(GREEN)✓ Worker build complete$(NC)"

# Run commands
run: ## Run API server (development)
	@echo "$(GREEN)Starting API server...$(NC)"
	@go run ./cmd/api/main.go

run-worker: ## Run Worker server
	@echo "$(GREEN)Starting Worker server...$(NC)"
	@go run ./cmd/worker/main.go

run-air: ## Run API server with Air (live reload)
	@echo "$(GREEN)Starting API server with Air...$(NC)"
	@air

# Test commands
test: ## Run all tests
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v ./tests/...

test-short: ## Run tests in short mode
	@echo "$(GREEN)Running tests (short mode)...$(NC)"
	@go test -short -v ./tests/...

test-coverage: ## Run tests with coverage report
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@bash scripts/test-coverage.sh

test-coverage-html: test-coverage ## Generate HTML coverage report
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

# Database commands (Atlas)
atlas-install: ## Install Atlas CLI
	@echo "$(GREEN)Installing Atlas CLI...$(NC)"
	@curl -sSf https://atlasgo.sh | sh
	@echo "$(GREEN)✓ Atlas CLI installed$(NC)"
	@atlas version

atlas-diff: ## Generate migration from GORM schema changes
	@echo "$(GREEN)Generating migration from GORM models...$(NC)"
	@ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas migrate diff --env dev
	@echo "$(GREEN)✓ Migration generated$(NC)"

atlas-diff-name: ## Generate migration with custom name (usage: make atlas-diff-name NAME=add_user_bio)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required$(NC)"; \
		echo "Usage: make atlas-diff-name NAME=add_user_bio"; \
		exit 1; \
	fi
	@echo "$(GREEN)Generating migration: $(NAME)$(NC)"
	@ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas migrate diff $(NAME) --env dev

atlas-apply: ## Apply pending migrations to database
	@echo "$(GREEN)Applying migrations...$(NC)"
	@ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas migrate apply --env dev
	@echo "$(GREEN)✓ Migrations applied$(NC)"

atlas-apply-prod: ## Apply migrations to production (with confirmation)
	@echo "$(RED)⚠️  WARNING: Applying migrations to PRODUCTION!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas migrate apply --env prod; \
		echo "$(GREEN)✓ Production migrations applied$(NC)"; \
	else \
		echo "$(YELLOW)Migration cancelled$(NC)"; \
	fi

atlas-status: ## Show migration status
	@echo "$(GREEN)Migration status:$(NC)"
	@ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas migrate status --env dev

atlas-status-prod: ## Show production migration status
	@echo "$(GREEN)Production migration status:$(NC)"
	@ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas migrate status --env prod

atlas-validate: ## Validate migration files
	@echo "$(GREEN)Validating migrations...$(NC)"
	@ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas migrate validate --env dev
	@echo "$(GREEN)✓ Migrations valid$(NC)"

atlas-inspect: ## Inspect current database schema
	@echo "$(GREEN)Inspecting database schema...$(NC)"
	@ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas schema inspect --env dev

atlas-lint: ## Lint migration files
	@echo "$(GREEN)Linting migrations...$(NC)"
	@ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas migrate lint --env dev

atlas-test: ## Test migrations with dry-run
	@echo "$(GREEN)Testing migrations (dry-run)...$(NC)"
	@ATLAS_DB_URL=$$(bash scripts/get-atlas-url.sh) atlas migrate apply --env dev --dry-run

atlas-migrate: ## Full workflow: generate, validate, and apply
	@echo "$(GREEN)Running full migration workflow...$(NC)"
	@make atlas-diff
	@make atlas-validate
	@make atlas-apply
	@echo "$(GREEN)✓ Migration workflow complete!$(NC)"

atlas-clean: ## Clean all migration files (DANGEROUS!)
	@echo "$(RED)⚠️  WARNING: This will DELETE all migration files!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		rm -rf migrations/*.sql; \
		echo "$(GREEN)Migration files cleaned$(NC)"; \
	else \
		echo "$(YELLOW)Clean cancelled$(NC)"; \
	fi

# Legacy migrate commands (deprecated - use atlas-* commands instead)
migrate-up: atlas-apply ## Alias for atlas-apply (deprecated)
	@echo "$(YELLOW)Note: Use 'make atlas-apply' instead$(NC)"

migrate-down: ## Rollback not supported with Atlas (use version control)
	@echo "$(RED)Atlas doesn't support down migrations$(NC)"
	@echo "$(YELLOW)Use git to revert migration files and regenerate$(NC)"

migrate-status: atlas-status ## Alias for atlas-status (deprecated)
	@echo "$(YELLOW)Note: Use 'make atlas-status' instead$(NC)"

# Code quality commands
lint: ## Run golangci-lint
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

fmt: ## Format code with gofmt
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

fmt-check: ## Check code formatting without modifying
	@echo "$(GREEN)Checking code format...$(NC)"
	@if [ $$(gofmt -l . | wc -l) -ne 0 ]; then \
		echo "$(RED)Code is not formatted. Run 'make fmt' to fix.$(NC)"; \
		gofmt -d .; \
		exit 1; \
	else \
		echo "$(GREEN)✓ Code is properly formatted$(NC)"; \
	fi

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Go vet passed$(NC)"

mod-tidy: ## Tidy go.mod and go.sum
	@echo "$(GREEN)Tidying dependencies...$(NC)"
	@go mod tidy
	@go mod verify
	@echo "$(GREEN)✓ Dependencies tidied$(NC)"

mod-download: ## Download dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@go mod download
	@go mod verify
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

# Docker commands
docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -f assets/docker/Dockerfile -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)✓ Docker image built: $(DOCKER_IMAGE)$(NC)"

docker-build-optimized: ## Build optimized Docker image (multi-stage)
	@echo "$(GREEN)Building optimized Docker image...$(NC)"
	@docker build -f assets/docker/Dockerfile.optimized -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)✓ Optimized Docker image built: $(DOCKER_IMAGE)$(NC)"

docker-run: ## Run Docker container
	@echo "$(GREEN)Running Docker container...$(NC)"
	@docker run -p 3000:3000 --env-file .env $(DOCKER_IMAGE)

docker-up: ## Start services with docker-compose
	@echo "$(GREEN)Starting services...$(NC)"
	@docker-compose -f assets/docker/docker-compose.yml up -d
	@echo "$(GREEN)✓ Services started$(NC)"

docker-down: ## Stop services with docker-compose
	@echo "$(YELLOW)Stopping services...$(NC)"
	@docker-compose -f assets/docker/docker-compose.yml down
	@echo "$(GREEN)✓ Services stopped$(NC)"

docker-logs: ## Show docker-compose logs
	@docker-compose -f assets/docker/docker-compose.yml logs -f

docker-clean: ## Clean Docker images and containers
	@echo "$(YELLOW)Cleaning Docker resources...$(NC)"
	@docker-compose -f assets/docker/docker-compose.yml down -v
	@docker system prune -f
	@echo "$(GREEN)✓ Docker cleaned$(NC)"

# Installation commands
install-tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(NC)"
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"
	@echo "$(YELLOW)Note: Install Atlas CLI separately with 'make atlas-install'$(NC)"

# Cleanup commands
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f coverage.out coverage.html
	@rm -f build-errors.log
	@echo "$(GREEN)✓ Clean complete$(NC)"

clean-all: clean ## Clean everything including test cache
	@echo "$(YELLOW)Cleaning test cache...$(NC)"
	@go clean -testcache
	@echo "$(GREEN)✓ All clean$(NC)"

# CI/CD helpers
ci-test: fmt-check vet test ## Run all CI checks (format, vet, test)
	@echo "$(GREEN)✓ All CI checks passed$(NC)"

ci-build: mod-download build ## Build for CI
	@echo "$(GREEN)✓ CI build complete$(NC)"

# Development workflow
dev-setup: install-tools mod-download atlas-install ## Setup development environment
	@echo "$(GREEN)✓ Development environment ready$(NC)"
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Copy .env to configure your environment"
	@echo "  2. Set DB_GEN=true for development (uses AutoMigrate)"
	@echo "  3. Or use 'make atlas-diff' + 'make atlas-apply' for Atlas migrations"
	@echo "  4. Run 'make run-air' to start development server with live reload"

# Default target
.DEFAULT_GOAL := help

