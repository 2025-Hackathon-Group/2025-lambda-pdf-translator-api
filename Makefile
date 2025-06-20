.DEFAULT_GOAL := run

.PHONY: build test clean run

# Binary name
BINARY_NAME=app
MODULE_PATH=cmd/cli/main.go

# Build the application
build:
	go build -o $(BINARY_NAME) $(MODULE_PATH)

# Run the application
run:
	gow run $(MODULE_PATH)

# Run tests
test:
	go test -v ./...

# Seed the database
seed:
	go run cmd/cli/seed/main.go run

# Clean build files
clean:
	go clean
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	go mod download

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Default target
all: deps build test

migrate:
	go run cmd/migrate/main.go

.PHONY: swag
swag:
	swag init -g cmd/cli/main.go