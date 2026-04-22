# Task Manager App

A production-style task management system with a Go REST API, PostgreSQL persistence, Redis-backed task caching, Swagger documentation, and a lightweight web client served by Caddy.

## Features

- User management (`create`, `list`, `get by id`, `patch`, `delete`)
- Task management (`create`, `list`, `get by id`, `patch`, `delete`)
- Task statistics endpoint with optional filtering by user and date range
- Three-state PATCH semantics (omit field / set value / set `null` where allowed)
- PostgreSQL as source of truth with SQL migrations
- Redis cache layer for task repository operations
- OpenAPI/Swagger docs available from the running service
- Browser UI for managing users, tasks, and statistics
- Containerized local environment (API + Postgres + Redis + Caddy)

## Tech Stack

### Backend

- Go `1.26.2`
- Standard library `net/http` server/mux
- `github.com/go-playground/validator/v10` for request validation
- `github.com/google/uuid` for UUID handling
- `go.uber.org/zap` for structured logging

### Database & Caching

- PostgreSQL `18.1`
- Redis `8.6`
- `github.com/jackc/pgx/v5` (PostgreSQL driver/pool)
- `github.com/redis/go-redis/v9` (Redis client)

### API Documentation

- Swagger/OpenAPI 2.0
- `github.com/swaggo/swag`
- `github.com/swaggo/http-swagger/v2`

### DevOps & Tooling

- Docker + Docker Compose
- Caddy `2.11` as reverse proxy/static server
- `migrate/migrate` for SQL migrations
- Makefile for common workflows

## Installation

### Prerequisites

- Go `1.26.2` (or compatible `1.26.x`)
- Docker Engine + Docker Compose v2
- GNU Make (optional, but recommended)

### 1. Clone the repository

```bash
git clone https://github.com/rod1kutzyy/task-manager-app.git
cd task-manager-app
```

### 2. Configure environment variables

```bash
cp .env.example .env
```

Then fill required values in `.env` (example in [Configuration](#configuration)).

### 3. Start infrastructure (Postgres + Redis)

Using Make:

```bash
make env-up
make env-port-forward
```

Or using Docker Compose directly:

```bash
docker compose up -d notesapp-postgres notesapp-redis port-forwarder
```

### 4. Run DB migrations

```bash
make migrate-up
```

### 5. Run the API locally

```bash
make app-run
```

API will be available at `http://localhost:5050`.

### 6. (Optional) Run web server (Caddy)

```bash
make app-deploy
make web-deploy
```

Then open `http://localhost`.

## Usage

### Common workflows

### A) Local development (Go process + containerized DB/cache)

```bash
make env-up
make env-port-forward
make migrate-up
make app-run
```

### B) Full containerized stack

```bash
docker compose up -d --build notesapp notesapp-postgres notesapp-redis web-server
```

### C) Stop services

```bash
make app-undeploy
make web-undeploy
make env-port-close
make env-down
```

### D) Regenerate Swagger docs

```bash
make swagger-gen
```

## Project Structure

```text
.
|-- cmd/notesapp/                # App entrypoint and API Dockerfile
|-- docs/                        # Generated Swagger docs (yaml/json/go)
|-- internal/
|   |-- core/                    # Shared config, logger, transport, db pools
|   |-- features/
|   |   |-- users/               # Users domain: service/repository/http
|   |   |-- tasks/               # Tasks domain: service/ports/adapters/http
|   |   |-- statistics/          # Statistics domain: service/repository/http
|-- migrations/                  # SQL migrations
|-- web/
|   |-- Caddyfile                # Reverse proxy + static hosting config
|   |-- public/                  # Frontend assets (HTML/CSS/JS)
|-- out/                         # Runtime outputs (logs, db volumes, etc.)
|-- docker-compose.yaml          # Local/dev environment services
|-- Makefile                     # Developer commands
|-- .env.example                 # Environment variable template
```

## Configuration

The app uses environment variables (loaded from `.env` in Docker workflows).

### Example `.env`

```env
CADDY_HOST=http://localhost

HTTP_ADDR=:5050
HTTP_SHUTDOWN_TIMEOUT=30s
HTTP_ALLOWED_ORIGINS=http://localhost:5050,http://127.0.0.1:5050,http://localhost:5500,http://127.0.0.1:5500

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=task_manager
POSTGRES_TIMEOUT=10s

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis
REDIS_DB=0
REDIS_TTL=5m

LOGGER_LEVEL=DEBUG
LOGGER_FOLDER=./out/logs

TIME_ZONE=UTC
```

### Key variables

- `HTTP_ADDR`: API bind address (e.g., `:5050`)
- `HTTP_ALLOWED_ORIGINS`: CORS allowlist (comma-separated)
- `POSTGRES_*`: PostgreSQL connection settings
- `REDIS_*`: Redis connection/cache settings
- `LOGGER_*`: log level and output directory
- `TIME_ZONE`: app timezone (IANA, e.g., `UTC`, `Europe/Moscow`)
- `CADDY_HOST`: host used by Caddy site block

## API Documentation

### Swagger

- Swagger UI: `/swagger/index.html`
- OpenAPI JSON: `/swagger/doc.json`
- Base API path: `/api/v1`

### Endpoint summary

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/users` | List users |
| `POST` | `/api/v1/users` | Create user |
| `GET` | `/api/v1/users/{id}` | Get user by ID |
| `PATCH` | `/api/v1/users/{id}` | Partially update user |
| `DELETE` | `/api/v1/users/{id}` | Delete user |
| `GET` | `/api/v1/tasks` | List tasks |
| `POST` | `/api/v1/tasks` | Create task |
| `GET` | `/api/v1/tasks/{id}` | Get task by ID |
| `PATCH` | `/api/v1/tasks/{id}` | Partially update task |
| `DELETE` | `/api/v1/tasks/{id}` | Delete task |
| `GET` | `/api/v1/statistics` | Task statistics |

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for details.
