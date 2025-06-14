.PHONY: build test clean run

# Binary name
BINARY_NAME=app

# Build the application
build:
	go build -o $(BINARY_NAME) ./...

# Run the application
run:
	go run ./...

# Run tests
test:
	go test -v ./...

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
