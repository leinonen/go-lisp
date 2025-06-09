# Lisp Interpreter Makefile

# Variables
BINARY_NAME=lisp

.PHONY: build run test

# Default target
all: build

build: ## Build the interpreter binary
	go build -o $(BINARY_NAME) ./cmd/lisp-interpreter

run: build ## Build and run the interpreter
	./$(BINARY_NAME)

test: ## Run all tests
	go test ./...
