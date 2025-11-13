# Binary name
binary_name := "mw-cli"

# Default recipe - show help
default:
    @just --list

# Build the application
build:
    @echo "Building {{binary_name}}..."
    go build -o {{binary_name}} -v main.go
    @echo "Build complete: ./{{binary_name}}"

# Install the binary to GOPATH/bin
install:
    @echo "Installing {{binary_name}}..."
    go install
    @echo "Installed to $(go env GOPATH)/bin/{{binary_name}}"

# Run the application with a word
# Usage: just run serendipity
run word:
    @go build -o {{binary_name}} main.go
    @./{{binary_name}} {{word}}

# Run tests
test:
    @echo "Running tests..."
    go test -v ./...

# Run tests with coverage
test-coverage:
    @echo "Running tests with coverage..."
    go test -v -cover -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
    @echo "Cleaning..."
    go clean
    rm -f {{binary_name}}
    rm -f coverage.out coverage.html
    @echo "Clean complete"

# Format code
fmt:
    @echo "Formatting code..."
    go fmt ./...
    @echo "Format complete"

# Run linter
lint:
    @echo "Running linter..."
    golangci-lint run
    @echo "Lint complete"

# Run go vet
vet:
    @echo "Running go vet..."
    go vet ./...
    @echo "Vet complete"

# Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
    @echo "All checks passed!"

# Download and tidy dependencies
deps:
    @echo "Downloading dependencies..."
    go get -v ./...
    go mod tidy
    @echo "Dependencies updated"

# Build and run in one command
# Usage: just dev serendipity
dev word: build
    @./{{binary_name}} {{word}}
