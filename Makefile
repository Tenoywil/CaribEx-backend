# CaribX Backend Makefile

.PHONY: help build test lint run-dev clean docker-build docker-up docker-down migrate-up migrate-down

# Default target
help:
	@echo "CaribX Backend - Available targets:"
	@echo "  build        - Build the API server binary"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linters"
	@echo "  run-dev      - Run the server in development mode"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up    - Start Docker Compose services"
	@echo "  docker-down  - Stop Docker Compose services"
	@echo "  migrate-up   - Run database migrations up"
	@echo "  migrate-down - Run database migrations down"

# Build the application
build:
	@echo "Building API server..."
	@go build -o bin/api-server ./cmd/api-server

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run the server in development mode
run-dev:
	@echo "Starting server in development mode..."
	@go run ./cmd/api-server

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t caribx-backend:latest .

# Start Docker Compose services
docker-up:
	@echo "Starting Docker Compose services..."
	@docker-compose up -d

# Stop Docker Compose services
docker-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

# Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/caribx?sslmode=disable" up; \
	else \
		echo "migrate not installed. Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
	fi

# Run database migrations down
migrate-down:
	@echo "Running migrations down..."
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/caribx?sslmode=disable" down; \
	else \
		echo "migrate not installed. Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
	fi

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run in watch mode (requires air)
watch:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
	fi
