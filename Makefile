.PHONY: build run test clean deps lint

# Build the application
build:
	go build -o bin/loan-service main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run linter
lint:
	golangci-lint run

# Run linter with fix
lint-fix:
	golangci-lint run --fix

# Generate API documentation (if using swagger)
docs:
	swag init -g main.go

# Run in development mode
dev:
	go run main.go

# Database migrations
migrate:
	go run main.go migrate

# Rollback migrations
migrate-down:
	go run main.go migrate-down

# Create new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql

# Generate mocks
generate-mocks:
	mockery --dir internal/service --name LoanServiceInterface --output mocks --outpkg mocks
	mockery --dir internal/repository --name LoanRepositoryInterface --output mocks --outpkg mocks
	mockery --dir internal/client --name NotificationClientInterface --output mocks --outpkg mocks

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/loan-service main.go

# Docker build
docker-build:
	docker build -t loan-service .

# Docker run
docker-run:
	docker run -p 8080:8080 loan-service
