# Makefile Guide

This repository contains two Makefiles to simplify development workflows:
1. **Root `Makefile`**: Manages Docker orchestration, multi-service lifecycle (frontend + backend), and repo-wide utilities.
2. **Backend `backend/Makefile`**: Manages Go testing, linting, formatting, Mockery mock generation, and Supabase database migrations.

---

## 1. Root Commands (`/Makefile`)

Run these commands from the repository root (`/`):

### Service Lifecycle & Docker
| Command | Description |
|---|---|
| `make up` | Start all services (frontend + backend) in Docker with hot-reload |
| `make down` | Stop and remove all Docker containers |
| `make restart` | Restart all containers |
| `make logs` | Stream logs from all running services |
| `make build` | Build all service images |
| `make stop` | Stop containers without removing them |
| `make ps` | Show running Docker containers |

### Individual Services
| Command | Description |
|---|---|
| `make up-backend` | Start only the backend service with hot-reload |
| `make up-frontend` | Start only the frontend service with hot-reload |
| `make logs-backend` | Stream backend logs |
| `make logs-frontend` | Stream frontend logs |
| `make build-backend` | Build backend container |
| `make build-frontend` | Build frontend container |
| `make shell-backend` | Open a shell inside the running backend container |
| `make shell-frontend` | Open a shell inside the running frontend container |

### Utilities & Code Generation
| Command | Description |
|---|---|
| `make generate-api` | Run Orval to regenerate TypeScript API clients & Zod schemas from OpenAPI |
| `make format-frontend` | Run frontend linter with auto-fix |
| `make clean-node-modules` | Remove all `node_modules` folders and reinstall fresh dependencies |
| `make clean` | Remove all Docker containers and volumes |
| `make prune` | Remove all unused Docker system resources and build cache |

---

## 2. Backend Commands (`backend/Makefile`)

Run these commands from the `backend/` directory:

### Testing & Quality
| Command | Description |
|---|---|
| `make test` | Run all tests (quiet mode with summary) |
| `make test-verbose` | Run all tests with verbose output (`-v`) |
| `make test-unit` | Run fast unit & route mock tests (skips slow DB integration tests with `-short`) |
| `make test-db` | Run database integration tests in `internal/data/postgres/` |
| `make test-coverage` | Run tests and generate HTML coverage report (`coverage.html`) |
| `make test-one TEST=...` | Run a specific test function by name (e.g., `make test-one TEST=TestRoute_CreateGreeting`) |
| `make test-clean` | Clear Go test cache and temporary test logs/coverage files |
| `make lint` | Run `golangci-lint` |
| `make lint-fix` | Run `golangci-lint` with auto-fix |
| `make format` | Format Go code with `gofmt` |
| `make format-check` | Check if code is formatted without making changes |

### Mock Generation
| Command | Description |
|---|---|
| `make mocks` | Regenerate repository mocks using Mockery into `internal/data/repo-mocks/` |

### Database & Supabase
| Command | Description |
|---|---|
| `make db-start` | Start local Supabase instance (requires Docker) |
| `make db-stop` | Stop local Supabase instance |
| `make db-status` | Display local Supabase service status and credentials |
| `make db-new NAME=...` | Create a new Supabase SQL migration file |
| `make db-reset` | Reset local database and run all pending migrations |
| `make db-link REF=...` | Link local CLI to remote Supabase project |
| `make db-push` | Push migrations to remote Supabase production database |

### Development & Build
| Command | Description |
|---|---|
| `make dev` | Run backend server locally with `.env` loaded |
| `make build` | Build the Go binary to `bin/server` |
| `make clean` | Clean compiled binaries and test cache |
| `make tidy` | Run `go mod tidy` |
| `make deps` | Run tidy, download, and verify dependencies |
| `make api-preview` | Open the interactive Scalar API documentation (`/docs`) in your browser |
