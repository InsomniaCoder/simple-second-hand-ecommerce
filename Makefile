.PHONY: help build run test fmt vet coverage clean check install

# Default target
help:
	@echo "Available targets:"
	@echo "  make build      - Build the API binary"
	@echo "  make run        - Run the API server"
	@echo "  make test       - Run all tests"
	@echo "  make fmt        - Format all Go files"
	@echo "  make vet        - Run go vet"
	@echo "  make coverage   - Generate coverage report"
	@echo "  make check      - Run fmt, vet, and test"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make install    - Install development tools"

# Build the binary
build:
	@echo "Building..."
	go build -o bin/api cmd/api/main.go

# Run the server
run:
	@echo "Starting server..."
	go run cmd/api/main.go

# Run all tests
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	@bash scripts/test-coverage.sh

# Run all checks (format, vet, test)
check: fmt vet test
	@echo "✓ All checks passed"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Install development tools
install:
	@echo "Installing development tools..."
	@if ! command -v goimports >/dev/null 2>&1; then \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@echo "✓ Development tools installed"
