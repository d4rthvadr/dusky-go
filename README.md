# Motivation

## Technologies Used

- **Go 1.24+** - Programming language
- **Chi Router** - Lightweight HTTP router
- **PostgreSQL** - Relational database
- **Docker** - Database containerization
- **godotenv** - Environment variable management
- **lib/pq** - PostgreSQL driver

## Project Structure

```
.
├── cmd/
│   ├── api/           # Application entrypoint and HTTP handlers
│   └── migrations/    # Database migrations
├── internal/
│   ├── config/        # Application configuration
│   ├── db/            # Database connection setup
│   ├── env/           # Environment variable helpers
│   ├── models/        # Data models
│   └── store/         # Data access layer
├── docker-compose.yml # PostgreSQL setup
└── go.mod            # Go module dependencies
```

## Setup

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose

### Installation

1. Clone the repository

```bash
git clone <repository-url>
cd dusky-go
```

2. Install dependencies

```bash
go mod download
```

3. Set up environment variables

```bash
cp .env.example .env
```

4. Start PostgreSQL with Docker

```bash
docker-compose up -d
```

5. Run database migrations

```bash
make migrate-up
```

6. Build and run the application

```bash
go run ./cmd/api
```

The API will be available at `http://localhost:8082`

## Dev Container Onboarding

Use this section when onboarding a new engineer with VS Code Dev Containers.

### Prerequisites (host machine)

- Docker Desktop / Docker Engine running
- VS Code with the Dev Containers extension

### What comes out of the box in the dev container

- Go toolchain (from `mcr.microsoft.com/devcontainers/go`)
- VS Code extensions:
	- `golang.go`
	- `ms-azuretools.vscode-docker`
	- `redhat.vscode-yaml`
- `postgres` and `redis` started automatically by devcontainer compose
- Runtime environment variables for DB and Redis preconfigured in devcontainer
- Go module/build cache paths preconfigured for the non-root `vscode` user

### Open and start the project

1. Open the `dusky-go` folder in VS Code and choose **Reopen in Container**.
	- The container opens directly at the module root, so Explorer shows only project files (not parent workspace folders).
2. Wait for container setup to complete (`postCreateCommand` runs `go mod download`).
3. In the container terminal:

```bash
cp .env.example .env
```

4. Start the API:

```bash
go run ./cmd/api
```

The API will be available at `http://localhost:8082`.

### Optional: build and run binary

```bash
go build -o bin/main ./cmd/api
./bin/main
```

### Optional: migrations

If the `migrate` CLI is installed in the container:

```bash
make migrate-up
```

### Optional: hot reload (`make dev`)

`make dev` uses `air`. If `air` is missing, install it in the container:

```bash
go install github.com/air-verse/air@latest
```

## Development

### Running the server

```bash
go run ./cmd/api
```

### Database Migrations

View all available make commands:

```bash
make help
```

Create a new migration:

```bash
make migrate name=create_posts
make migrate name=add_index_to_users
```

Run all pending migrations:

```bash
make migrate-up
```

Rollback the last migration:

```bash
make migrate-down
```

Check current migration version:

```bash
make migrate-version
```

Force migration version (use when database is in dirty state):

```bash
make migrate-force version=1
```

### Building for production

```bash
go build -o bin/main ./cmd/api
./bin/main
```

### Stopping the database

```bash
docker-compose down
```

## Configuration

Environment variables can be configured in `.env`:

- `ADDR` - Server port (default in this repo: 8082)
- `DB_ADDR` - PostgreSQL connection string
- `DB_MAX_OPEN_CONNS` - Maximum open database connections (default: 30)
- `DB_MAX_IDLE_CONNS` - Maximum idle database connections (default: 30)
- `DB_MAX_IDLE_TIME` - Maximum idle time for connections (default: 15m)
