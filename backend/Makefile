# Space Notifications Backend Makefile

.PHONY: help build run clean utils-help utils-cleanup-events test fmt vet

# Default target
help:
	@echo "Space Notifications Backend"
	@echo ""
	@echo "Available targets:"
	@echo "  help                 Show this help message"
	@echo "  build               Build the main backend application"
	@echo "  run                 Run the main backend application"
	@echo "  build-utils         Build the utility commands"
	@echo "  utils-help          Show utility commands help"
	@echo "  utils-cleanup-events Run the duplicate events cleanup utility"
	@echo "  test                Run all tests"
	@echo "  fmt                 Format Go code"
	@echo "  vet                 Run Go vet"
	@echo "  clean               Clean build artifacts"
	@echo ""
	@echo "Environment variables needed:"
	@echo "  PB_URL              PocketBase URL (e.g., http://localhost:8080)"
	@echo "  PB_ADMIN_EMAIL      PocketBase admin email"
	@echo "  PB_ADMIN_PASSWORD   PocketBase admin password"

# Build targets
build:
	@echo "🔨 Building main backend application..."
	go build -o bin/space-notifications main.go

build-utils:
	@echo "🔨 Building utility commands..."
	go build -o bin/space-utils cmd/utils-main.go

# Run targets
run: build
	@echo "🚀 Starting main backend application..."
	./bin/space-notifications

utils-help: build-utils
	@echo "🔧 Running utility help..."
	./bin/space-utils -help

utils-cleanup-events: build-utils
	@echo "🧹 Running duplicate events cleanup..."
	./bin/space-utils -cleanup-events

# Alternative run commands (without building binaries)
run-dev:
	@echo "🚀 Running main backend in development mode..."
	go run main.go

utils-cleanup-events-dev:
	@echo "🧹 Running cleanup utility in development mode..."
	go run cmd/utils-main.go -cleanup-events

utils-help-dev:
	@echo "🔧 Running utility help in development mode..."
	go run cmd/utils-main.go -help

# Development targets
test:
	@echo "🧪 Running tests..."
	go test -v ./...

fmt:
	@echo "📝 Formatting Go code..."
	go fmt ./...

vet:
	@echo "🔍 Running Go vet..."
	go vet ./...

# Clean up
clean:
	@echo "🧽 Cleaning up..."
	rm -rf bin/
	go clean

# Create bin directory
bin:
	mkdir -p bin
