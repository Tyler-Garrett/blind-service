.PHONY: help build run test clean docker-build docker-run lint fmt vet tidy seed deploy

# Variables
SERVICE_NAME=blind-service
DOCKER_IMAGE=blind-service
DOCKER_TAG=latest
GO_VERSION=1.26

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(SERVICE_NAME)..."
	@go build -o $(SERVICE_NAME) ./cmd/api

run: ## Run the application
	@echo "Running $(SERVICE_NAME)..."
	@go run ./cmd/api

dev: ## Run with hot reload (requires air)
	@echo "Starting development server with hot reload..."
	@air -c .air.toml

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(SERVICE_NAME)
	@rm -f coverage.txt coverage.html
	@rm -rf tmp/

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f deployments/docker/Dockerfile .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Start services with docker-compose
	@echo "Starting services with docker-compose..."
	@cd deployments/docker && docker-compose up

docker-compose-down: ## Stop services with docker-compose
	@echo "Stopping services with docker-compose..."
	@cd deployments/docker && docker-compose down

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy

seed: ## Seed database with blind data
	@echo "Seeding database..."
	@go run scripts/seed-data.go

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

# CI/CD targets
ci: lint vet test ## Run CI checks

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest

# Azure deployment
deploy-azure: ## Deploy to Azure Container Apps
	@echo "Deploying to Azure..."
	@bash scripts/deploy.sh