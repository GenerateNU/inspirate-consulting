.PHONY: help test test-verbose test-unit test-db test-coverage test-one test-clean \
        lint lint-fix format format-check \
        mocks \
        db-new db-start db-stop db-reset db-push db-link db-status \
        dev run build clean tidy download verify vendor deps \
        api-preview

# Default target - show help
.DEFAULT_GOAL := help

# Files/packages to exclude from coverage (regex patterns)
COVERAGE_EXCLUDE := main.go|docs.go

# ------------------------
# OS Detection
# ------------------------
OS := $(shell uname 2>/dev/null || echo Windows_NT)
ifeq ($(OS),Windows_NT)
    GREP = findstr
    RM = del /F /Q
    RMDIR = rmdir /S /Q
    MKDIR = mkdir
else
    GREP = grep
    RM = rm -f
    RMDIR = rm -rf
    MKDIR = mkdir -p
endif

# ------------------------
# Colors
# ------------------------
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m
BOLD := \033[1m

# ------------------------
# Configuration
# ------------------------
NO_CACHE := -count=1

# ------------------------
# Help
# ------------------------
help:
	@echo "$(BOLD)Inspirate Consulting Backend - Available Commands$(NC)"
	@echo ""
	@echo "$(BLUE)Testing:$(NC)"
	@echo "  make test              - Run all tests (quiet mode)"
	@echo "  make test-verbose      - Run all tests with verbose output"
	@echo "  make test-unit         - Run unit tests only (skips DB tests)"
	@echo "  make test-db           - Run database tests"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-one TEST=... - Run a specific test"
	@echo "  make test-clean        - Clean test cache"
	@echo ""
	@echo "$(BLUE)Mocks & Code Generation:$(NC)"
	@echo "  make mocks             - Generate mocks with Mockery"
	@echo ""
	@echo "$(BLUE)Code Quality:$(NC)"
	@echo "  make lint              - Run linter"
	@echo "  make lint-fix          - Fix linting issues"
	@echo "  make format            - Format code with gofmt"
	@echo "  make format-check      - Check if code needs formatting"
	@echo ""
	@echo "$(BLUE)Dependencies:$(NC)"
	@echo "  make tidy              - Run go mod tidy"
	@echo "  make download          - Download dependencies"
	@echo "  make verify            - Verify dependencies"
	@echo "  make vendor            - Vendor dependencies"
	@echo "  make deps              - Update all dependencies"
	@echo ""
	@echo "$(BLUE)Database (Supabase):$(NC)"
	@echo "  make db-new NAME=...   - Create new migration"
	@echo "  make db-start          - Start local Supabase (requires Docker)"
	@echo "  make db-stop           - Stop local Supabase"
	@echo "  make db-reset          - Reset local DB and apply migrations"
	@echo "  make db-status         - Show local DB status"
	@echo "  make db-link REF=...   - Link CLI to remote Supabase project"
	@echo "  make db-push           - Push migrations to remote database"
	@echo ""
	@echo "$(BLUE)Development:$(NC)"
	@echo "  make dev               - Run server in development mode"
	@echo "  make run               - Run server"
	@echo "  make build             - Build the application binary"
	@echo "  make clean             - Clean build artifacts and test files"
	@echo ""
	@echo "$(BLUE)API Documentation:$(NC)"
	@echo "  make api-preview       - Open Scalar API documentation in browser"
	@echo ""

# ------------------------
# Testing
# ------------------------
test:
	@echo "$(BOLD)Running all tests...$(NC)"
	@go test ./... $(NO_CACHE) 2>&1 | tee test_output.log
	@echo ""
ifeq ($(OS),Windows_NT)
	@findstr /C:"FAIL" test_output.log >nul 2>&1 && ( \
		echo $(RED)$(BOLD)Tests Failed:$(NC) && \
		echo. && \
		findstr "FAIL" test_output.log && \
		del /F /Q test_output.log && \
		exit 1 \
	) || ( \
		echo $(GREEN)$(BOLD)All Tests Passed!$(NC) && \
		del /F /Q test_output.log \
	)
else
	@if grep -q "FAIL" test_output.log 2>/dev/null; then \
		echo "$(RED)$(BOLD)Tests Failed:$(NC)"; \
		echo ""; \
		grep "FAIL" test_output.log; \
		echo ""; \
		echo "$(YELLOW)Tip: Run 'make test-verbose' for detailed output$(NC)"; \
		rm -f test_output.log; \
		exit 1; \
	else \
		echo "$(GREEN)$(BOLD)✓ All Tests Passed!$(NC)"; \
		rm -f test_output.log; \
	fi
endif

test-verbose:
	@echo "$(BOLD)Running all tests (verbose)...$(NC)"
	@go test ./... -v $(NO_CACHE)

test-unit:
	@echo "$(BOLD)Running unit tests (short mode)...$(NC)"
	@go test ./... $(NO_CACHE) -short

test-db:
	@echo "$(BOLD)Running database tests...$(NC)"
	@go test ./internal/data/postgres/... $(NO_CACHE)

# Run tests with coverage
test-coverage:
	@echo "$(BOLD)Running tests with coverage...$(NC)"
	@go test ./... -coverprofile=coverage.out $(NO_CACHE)
ifeq ($(OS),Windows_NT)
	@type coverage.out | findstr /V /R "$(COVERAGE_EXCLUDE)" > coverage_filtered.out 2>nul || copy coverage.out coverage_filtered.out >nul
	@go tool cover -html=coverage_filtered.out -o coverage.html
	@echo $(GREEN)Coverage report: coverage.html$(NC)
	@go tool cover -func=coverage_filtered.out | findstr "total"
else
	@grep -vE '$(COVERAGE_EXCLUDE)' coverage.out > coverage_filtered.out || cp coverage.out coverage_filtered.out
	@go tool cover -html=coverage_filtered.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"
	@go tool cover -func=coverage_filtered.out | grep total
endif
	@echo "$(YELLOW)Excluded patterns: $(COVERAGE_EXCLUDE)$(NC)"

coverage-percentage:
	@go tool cover -func=coverage.out | $(GREP) total | awk '{print substr($$3, 1, length($$3)-1)}'

test-one:
	@if [ -z "$(TEST)" ]; then \
		echo "$(RED)Please specify a test: make test-one TEST=TestName$(NC)"; \
		exit 1; \
	fi
	@echo "$(BOLD)Running test: $(TEST)$(NC)"
	@go test ./... -v -run "$(TEST)" $(NO_CACHE)

test-clean:
	@go clean -testcache
	@$(RM) coverage.out coverage.html coverage_filtered.out test_output.log 2>/dev/null || true
	@echo "$(GREEN)Test cache cleaned$(NC)"

# ------------------------
# Mocks & Code Generation
# ------------------------
mocks:
	@echo "$(BOLD)Generating repository mocks with Mockery...$(NC)"
	@mockery --name=GreetingRepository --dir=internal/data --output=internal/data/repo-mocks --outpkg=mocks --filename=greetingMock.go
	@echo "$(GREEN)Mocks generated in internal/data/repo-mocks/$(NC)"

# ------------------------
# Code Quality
# ------------------------
lint:
	@echo "$(BOLD)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install: https://golangci-lint.run/usage/install/$(NC)"; \
		exit 1; \
	fi

lint-fix:
	@echo "$(BOLD)Fixing linting issues...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix; \
		echo "$(GREEN)Linting issues fixed$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install: https://golangci-lint.run/usage/install/$(NC)"; \
		exit 1; \
	fi

format:
	@echo "$(BOLD)Formatting code...$(NC)"
	@gofmt -w .
	@echo "$(GREEN)Code formatted$(NC)"

format-check:
	@echo "$(BOLD)Checking code formatting...$(NC)"
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "$(RED)The following files need formatting:$(NC)"; \
		gofmt -l .; \
		exit 1; \
	else \
		echo "$(GREEN)All files are properly formatted$(NC)"; \
	fi

# ------------------------
# Dependencies
# ------------------------
tidy:
	@echo "$(BOLD)Running go mod tidy...$(NC)"
	@go mod tidy
	@echo "$(GREEN)Dependencies tidied$(NC)"

download:
	@echo "$(BOLD)Downloading dependencies...$(NC)"
	@go mod download
	@echo "$(GREEN)Dependencies downloaded$(NC)"

verify:
	@echo "$(BOLD)Verifying dependencies...$(NC)"
	@go mod verify
	@echo "$(GREEN)Dependencies verified$(NC)"

vendor:
	@echo "$(BOLD)Vendoring dependencies...$(NC)"
	@go mod vendor
	@echo "$(GREEN)Dependencies vendored to vendor/$(NC)"

deps: tidy download verify
	@echo "$(GREEN)All dependency operations complete$(NC)"

# ------------------------
# Database (Supabase)
# ------------------------
db-new:
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Please specify migration name: make db-new NAME=add_users_table$(NC)"; \
		exit 1; \
	fi
	@echo "$(BOLD)Creating new migration: $(NAME)$(NC)"
	@cd internal/supabase && supabase migration new $(NAME)
	@echo "$(GREEN)Migration created in internal/supabase/migrations/$(NC)"
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Add your SQL to the migration file"
	@echo "  2. Run 'make db-reset' to apply locally"
	@echo "  3. Test your changes at http://localhost:54323"

db-start:
	@echo "$(BOLD)Starting local Supabase (requires Docker)...$(NC)"
ifeq ($(OS),Windows_NT)
	@docker info >nul 2>&1 || ( \
		echo $(RED)Docker is not running. Please start Docker first.$(NC) && \
		exit 1 \
	)
else
	@if ! docker info >/dev/null 2>&1; then \
		echo "$(RED)Docker is not running. Please start Docker first.$(NC)"; \
		exit 1; \
	fi
endif
	@cd internal/supabase && supabase start
	@echo ""
	@echo "$(GREEN)Supabase is running!$(NC)"
	@echo "Dashboard: $(BLUE)http://localhost:54323$(NC)"
	@echo "DB URL: $(BLUE)postgresql://postgres:postgres@127.0.0.1:54322/postgres$(NC)"

db-stop:
	@echo "$(BOLD)Stopping local Supabase...$(NC)"
	@cd internal/supabase && supabase stop
	@echo "$(GREEN)Supabase stopped$(NC)"

db-reset:
	@echo "$(BOLD)Resetting local database and applying migrations...$(NC)"
	@cd internal/supabase && supabase db reset
	@echo "$(GREEN)Database reset complete$(NC)"
	@echo "View at: $(BLUE)http://localhost:54323$(NC)"

db-status:
	@echo "$(BOLD)Local Supabase status:$(NC)"
	@cd internal/supabase && supabase status

db-link:
	@if [ -z "$(REF)" ]; then \
		echo "$(RED)Please provide project ref: make db-link REF=your-project-ref$(NC)"; \
		echo "$(YELLOW)Find your project ref in Supabase UI after /project/$(NC)"; \
		exit 1; \
	fi
	@echo "$(BOLD)Linking to Supabase project...$(NC)"
	@cd internal/supabase && supabase link --project-ref $(REF)
	@echo "$(GREEN)Successfully linked to project$(NC)"

db-push:
	@echo "$(BOLD)Pushing migrations to remote database...$(NC)"
	@echo "$(YELLOW)⚠️  This will modify the production database!$(NC)"
	@echo "Make sure your migrations are tested and approved by a TL."
ifeq ($(OS),Windows_NT)
	@set /p CONFIRM="Continue? [y/N]: " && if /i "!CONFIRM!" neq "y" ( \
		echo $(YELLOW)Cancelled$(NC) && exit /b 1 \
	) else ( \
		cd internal/supabase && supabase db push && echo $(GREEN)Migrations pushed successfully$(NC) \
	)
else
	@read -p "Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		cd internal/supabase && supabase db push; \
		echo "$(GREEN)Migrations pushed successfully$(NC)"; \
	else \
		echo "$(YELLOW)Cancelled$(NC)"; \
	fi
endif

# ------------------------
# Development
# ------------------------
dev:
	@echo "$(BOLD)Starting development server...$(NC)"
	@if [ ! -r .env ] && [ ! -r ../.env ]; then \
		echo "$(RED)Missing required .env file (looked in . and ..)$(NC)"; \
		exit 1; \
	fi; \
	if [ -r .env ]; then \
		set -a; . ./.env; set +a; \
	else \
		set -a; . ../.env; set +a; \
	fi; \
	go run cmd/main.go

run:
	@echo "$(BOLD)Starting server...$(NC)"
	@if [ ! -r .env ] && [ ! -r ../.env ]; then \
		echo "$(RED)Missing required .env file (looked in . and ..)$(NC)"; \
		exit 1; \
	fi; \
	if [ -r .env ]; then \
		set -a; . ./.env; set +a; \
	else \
		set -a; . ../.env; set +a; \
	fi; \
	go run cmd/main.go

build:
	@echo "$(BOLD)Building application binary...$(NC)"
	@$(MKDIR) bin 2>/dev/null || true
	@go build -o bin/server cmd/main.go
	@echo "$(GREEN)Build complete: bin/server$(NC)"

# ------------------------
# Cleanup
# ------------------------
clean:
	@echo "$(BOLD)Cleaning build artifacts...$(NC)"
	@$(RMDIR) bin 2>/dev/null || true
	@$(RM) coverage.out coverage.html coverage_filtered.out test_output.log 2>/dev/null || true
	@go clean -testcache
	@echo "$(GREEN)Cleaned$(NC)"

# ------------------------
# API Documentation
# ------------------------
api-preview:
	@echo "$(BOLD)Opening Scalar API documentation...$(NC)"
ifeq ($(OS),Windows_NT)
	@start http://localhost:8080/docs
else ifeq ($(OS),Darwin)
	@open http://localhost:8080/docs
else
	@xdg-open http://localhost:8080/docs 2>/dev/null || echo "$(GREEN)Visit: http://localhost:8080/docs$(NC)"
endif
	@echo "$(YELLOW)Note: Make sure your dev server is running (make dev)$(NC)"
