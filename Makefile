# GoLisp Makefile

# Variables
BINARY_NAME=golisp
CORE_BINARY_NAME=golisp-core

.PHONY: build run test fmt build-core run-core

# Default target
all: build

build: ## Build the interpreter binary
	go build -o ./bin/$(BINARY_NAME) ./cmd/golisp

build-core: ## Build the minimal core interpreter
	go build -o ./bin/$(CORE_BINARY_NAME) ./cmd/golisp-core

run: build ## Build and run the interpreter
	./bin/$(BINARY_NAME)

run-core: build-core ## Build and run the minimal core interpreter
	./bin/$(CORE_BINARY_NAME)

test: ## Run all tests
	go test ./...

test-core: ## Run minimal core tests only
	go test ./pkg/core/...

test-nocache: ## Run all tests without cache
	go test -count=1 ./...

test-core-nocache: ## Run minimal core tests without cache
	go test -count=1 ./pkg/core/...

fmt: ## Format all Go source files
	go fmt ./...
