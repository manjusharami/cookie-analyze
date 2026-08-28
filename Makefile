 .PHONY: build run test test-unit test-integration test-coverage clean help

BINARY_NAME=cookie_analyzer
MAIN_PATH=./cmd/main.go

# Default parameters for 'make run' if none are supplied
FILE ?= cookie_log.csv
DATE ?= 2018-12-09

build:
	@echo "Building binary..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	./$(BINARY_NAME) -f $(FILE) -d $(DATE)

test:
	@echo "Running all tests..."
	go test -v ./...

test-unit:
	@echo "Running unit tests..."
	go test -v ./internal/... ./cmd/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/...

test-coverage:
	@echo "Running tests with statement coverage..."
	go test -coverpkg=./cmd/...,./internal/... -coverprofile=coverage.out ./cmd/... ./internal/...
	go tool cover -func=coverage.out

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) coverage.out