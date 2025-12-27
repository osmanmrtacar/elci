.PHONY: help run build test clean install dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy

run: ## Run the application
	@echo "Starting server..."
	go run cmd/server/main.go

dev: ## Run in development mode with auto-reload (requires air)
	@echo "Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/air-verse/air@latest"; \
		echo "Running without auto-reload..."; \
		go run cmd/server/main.go; \
	fi

build: ## Build the application
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

format: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"; \
	fi

db-reset: ## Reset database (delete and recreate)
	@echo "Resetting database..."
	rm -f data/*.db data/*.db-*
	@echo "Database reset. It will be recreated on next server start."
