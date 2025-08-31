# Loan Service

A RESTful API service for managing loans built with Go, gorilla/mux, and MySQL.

## Features

- RESTful API endpoints for loan management
- MySQL database with GORM ORM
- Request validation and error handling
- Docker support for easy deployment
- Comprehensive logging

## Project Structure

```
loan-service/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── Dockerfile             # Docker configuration
├── docker-compose.yml     # Docker Compose setup
├── Makefile              # Build and development commands
├── env.example           # Environment variables template
├── README.md             # Project documentation
├── migrations/            # Database migration files
└── internal/
    ├── config/           # Configuration management
    ├── database/        # Database connection and setup
    ├── handlers/        # HTTP request handlers
    ├── models/          # Data models and structs
    ├── repository/      # Data access layer
    └── service/         # Business logic layer
```

## Prerequisites

- Go 1.23 or higher
- MySQL 8.0 or higher
- Docker and Docker Compose (optional)

## Quick Start

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd loan-service
```

2. Copy environment file:
```bash
cp env.example .env
```

3. Start the application with Docker Compose:
```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

### Manual Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Set up environment variables:
```bash
cp env.example .env
# Edit .env with your database credentials
```

3. Start MySQL database

4. Run the application:
```bash
go run main.go
```

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Loans
- `GET /v1/loans` - Get all loans
- `POST /v1/loans` - Create new loan
- `GET /v1/loans/{id}` - Get loan by ID
- `PUT /v1/loans/{id}` - Update loan
- `DELETE /v1/loans/{id}` - Delete loan

## Development

### Available Make Commands

```bash
make build        # Build the application
make run          # Run the application
make test         # Run tests
make test-coverage # Run tests with coverage
make clean        # Clean build artifacts
make deps         # Install dependencies
make lint         # Run linter
make dev          # Run in development mode
```

### Running Tests

```bash
go test ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Database Migrations

The application uses Goose for database migrations. Migrations are SQL files located in the `migrations/` directory.

#### Running Migrations

```bash
# Run all pending migrations
make migrate

# Or directly
go run main.go migrate

# Rollback last migration
make migrate-down

# Or directly
go run main.go migrate-down

# Create new migration
make migrate-create
```

#### Migration Files

Migration files follow the naming convention: `{timestamp}_{description}.sql`

Example: `20240101000000_create_loans_table.sql`

The timestamp format is `YYYYMMDDHHMMSS` which ensures proper ordering and avoids conflicts when multiple developers create migrations simultaneously.

## Configuration

The application uses environment variables for configuration. See `env.example` for available options:

- `SERVER_PORT` - HTTP server port (default: 8080)
- `SERVER_HOST` - HTTP server host (default: localhost)
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 3306)
- `DB_USER` - Database user (default: root)
- `DB_PASSWORD` - Database password (default: password)
- `DB_NAME` - Database name (default: loan_service)


## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License.
