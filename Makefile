# GoLisp Makefile

# Variables
BINARY_NAME=golisp

.PHONY: build run test fmt

# Default target
all: build

build: ## Build the interpreter binary
	go build -o ./bin/$(BINARY_NAME) ./cmd/golisp

run: build ## Build and run the interpreter
	./bin/$(BINARY_NAME)

test: ## Run all tests
	go test ./...

test-core: ## Run core tests only
	go test ./pkg/core/...

test-nocache: ## Run all tests without cache
	go test -count=1 ./...

test-core-nocache: ## Run core tests without cache
	go test -count=1 ./pkg/core/...

fmt: ## Format all Go source files
	go fmt ./...
