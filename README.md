# Motivation

## Technologies Used

- **Go 1.23.2** - Programming language
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

- Go 1.23.2 or higher
- Docker and Docker Compose

### Installation

1. Clone the repository

```bash
git clone <repository-url>
cd go-api-tutorial
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

5. Run database migrations (if applicable)

```bash
# Add migration commands here
```

6. Build and run the application

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:3000`

## Development

### Running the server

```bash
go run cmd/api/main.go
```

### Building for production

```bash
go build -o bin/main cmd/api/main.go
./bin/main
```

### Stopping the database

```bash
docker-compose down
```

## Configuration

Environment variables can be configured in `.env`:

- `ADDR` - Server port (default: 3000)
- `DB_ADDR` - PostgreSQL connection string
- `DB_MAX_OPEN_CONNS` - Maximum open database connections (default: 30)
- `DB_MAX_IDLE_CONNS` - Maximum idle database connections (default: 30)
- `DB_MAX_IDLE_TIME` - Maximum idle time for connections (default: 15m)
