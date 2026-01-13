# News Fetcher Makefile

.PHONY: help build run test clean docker-build docker-run docker-stop deps lint

# Default target
help:
	@echo "Available targets:"
	@echo "  deps        - Download Go dependencies"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application locally"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linter"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run with Docker Compose"
	@echo "  docker-stop - Stop Docker Compose services"
	@echo "  setup       - Setup development environment"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build:
	go build -o bin/news-fetcher cmd/worker/main.go

# Run the application locally
run:
	go run cmd/worker/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	docker build -t news-fetcher:latest .

# Run with Docker Compose
docker-run:
	docker-compose up --build

# Run in background
docker-run-bg:
	docker-compose up -d --build

# Stop Docker Compose services
docker-stop:
	docker-compose down

# View logs
docker-logs:
	docker-compose logs -f

# Setup development environment
setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp env.example .env; echo "Created .env file from template"; fi
	@echo "Please edit .env file with your API keys"
	@echo "Run 'make docker-run' to start Kafka and the application"

# Development helpers
dev-start-kafka:
	docker-compose up -d zookeeper kafka

dev-stop-kafka:
	docker-compose stop zookeeper kafka

# Check if required tools are installed
check-tools:
	@command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Aborting." >&2; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Aborting." >&2; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose is required but not installed. Aborting." >&2; exit 1; }

