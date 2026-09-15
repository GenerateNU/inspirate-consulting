# Inspirate Consulting Setup

This guide starts the PostgreSQL database, Go backend, and React frontend on
your computer. Run commands from the repository root unless a step says
otherwise.

## 1. Install the Required Tools

Install these tools before downloading project dependencies:

| Tool | Why it is needed | Check your installation |
| --- | --- | --- |
| [Git](https://git-scm.com/downloads) | Downloads and manages the repository | `git --version` |
| [Go 1.25+](https://go.dev/doc/install) | Builds and runs the backend | `go version` |
| [Bun](https://bun.sh/docs/installation) | Installs and runs the frontend | `bun --version` |
| [Docker](https://docs.docker.com/get-docker/) | Runs the local Supabase services | `docker --version` |
| [Docker Compose](https://docs.docker.com/compose/install/) | Runs the app in containers when desired | `docker compose version` |
| [Supabase CLI](https://supabase.com/docs/guides/local-development/cli/getting-started) | Manages the local database and migrations | `supabase --version` |

Start Docker before running any Supabase command.

Huma, Fiber, pgx, React, Vite, Orval, and the other code libraries do **not**
need to be installed globally. The next step downloads the versions used in

## 2. Install Project Dependencies

Install the backend dependencies from `backend/go.mod` and `backend/go.sum`:

```sh
cd backend
go mod download
cd ..
```

Install the frontend dependencies from `frontend/package.json` and
`frontend/bun.lock`:

```sh
cd frontend
bun install --frozen-lockfile
cd ..
```

## 3. Configure Environment Variables

Create your local environment file:

```sh
cp .env.example .env
```

The backend requires values for:

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME`
- `SUPABASE_URL`, `SUPABASE_ANON_KEY`, and `SUPABASE_SERVICE_ROLE_KEY`

The example file also contains optional application URLs and database settings.
Never commit `.env`; it contains credentials.

For local development, the next step's `supabase status` command prints the
database URL and Supabase keys. Use those values to fill in `.env`. A local
database normally uses port `54322` and `DB_SSLMODE=disable`.

The Go application reads variables from the process environment; it does not
load `.env` automatically. Before running the backend directly, export the file
from the repository root:

```sh
set -a
source .env
set +a
```

Set `TEST_MODE=true` locally if you need to call protected endpoints without a
JWT. Do not enable test mode in a deployed environment.

## 4. Start the Local Database

Supabase runs PostgreSQL and its supporting services in Docker:

```sh
cd backend/internal/supabase
supabase start
supabase status
supabase db reset
cd ../../..
```

`supabase db reset` deletes local database data and applies every migration. It
does not reset the hosted production database. Open Supabase Studio at
[http://127.0.0.1:54323](http://127.0.0.1:54323) to inspect local tables.

See the [Supabase guide](backend/internal/supabase/README.md) before creating or
deploying migrations.

## 5. Run the Application

Use two terminal windows so the backend and frontend can keep running together.

### Terminal 1: Backend

Load `.env`, enter the Go `main` package, and run it:

```sh
cd backend
set -a
source .env
set +a
cd backend/cmd
go run .
```

The API listens on [http://localhost:8080](http://localhost:8080) by default.
Huma's interactive API documentation is available at
[http://localhost:8080/docs](http://localhost:8080/docs).

You can also run the backend from the `backend` directory with `go run ./cmd`.

### Terminal 2: Frontend

```sh
cd frontend
bun run dev
```

Vite prints the frontend URL, normally
[http://localhost:5173](http://localhost:5173).

Stop either development server with `Ctrl+C`. Stop local Supabase from its
directory with `supabase stop`.

## Useful Commands

Run backend commands from `backend`:

```sh
go test ./...       # Run all Go tests
go fmt ./...        # Format Go code
go mod tidy         # Update dependency metadata after imports change
```

Run frontend commands from `frontend`:

```sh
bun run dev         # Start the Vite development server
bun run build       # Type-check and create a production build
bun run lint        # Run ESLint
bunx orval          # Regenerate API types, request functions, and mocks
```

Run the full application in development containers from the repository root:

```sh
docker compose up --build
docker compose down
```

The Compose backend reads the root `.env` file. The Supabase database must
already be running and `DB_HOST` must be reachable from the backend container.

## Common Problems

- **The backend reports a missing environment variable:** load `.env` in the
	same terminal before `go run .`.
- **The backend cannot connect to PostgreSQL:** make sure Docker is running,
	check `supabase status`, and compare its values with `.env`.
- **No tables appear in Studio:** run `supabase db reset`, check the command for
	migration errors, and refresh Studio.
- **Supabase rejects keys in `config.toml`:** update the Supabase CLI. Older CLI
	versions may not understand a config generated by a newer version.
- **A port is already in use:** stop the process using ports `8080`, `5173`, or
	the Supabase ports shown by `supabase status`.