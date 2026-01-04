.PHONY: help build test test-unit test-integration test-coverage lint run clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/api cmd/api/main.go

test: ## Run all tests
	go test -v ./test/...

test-unit: ## Run unit tests only
	go test -v -short ./test/...

test-integration: ## Run integration tests only
	go test -v -run Integration ./test/integration/...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	golangci-lint run ./...

run: ## Run the application
	go run cmd/api/main.go

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
