# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build and Run
- `make build` - Build the GoLisp interpreter binary to `./bin/golisp`
- `make run` - Build and run the interpreter in REPL mode
- `./bin/golisp` - Run the interpreter directly (REPL mode)
- `./bin/golisp -f script.lisp` - Execute a Lisp file
- `./bin/golisp -e '(+ 1 2 3)'` - Evaluate Lisp code directly

### Testing
- `make test` - Run all tests
- `make test-nocache` - Run all tests without cache (useful for debugging)
- `go test ./pkg/kernel` - Run tests for the kernel package specifically

### Code Quality
- `make fmt` - Format all Go source files
- `go fmt ./...` - Format Go code (alternative)

## Architecture Overview

GoLisp is a minimalist Lisp interpreter written in Go, inspired by Clojure. The codebase follows a clean separation of concerns:

### Core Components

**`pkg/kernel/`** - The heart of the interpreter containing:
- `types.go` - Core data types (Symbol, Keyword, List, Vector, HashMap, etc.) implementing the `Value` interface
- `parser.go` - Lexer and parser for converting text to AST (Token-based parsing with position tracking)
- `eval.go` - Core evaluation engine with context tracking for error reporting
- `env.go` - Environment/scope management for variable bindings
- `repl.go` - Read-Eval-Print-Loop implementation
- `bootstrap.go` - Bootstrapping built-in functions and standard library loading

**`cmd/golisp/main.go`** - CLI entry point supporting:
- Interactive REPL mode (default)
- File execution (`-f` flag)
- Direct code evaluation (`-e` flag)
- Help and usage information

**`lisp/`** - Lisp source files:
- `stdlib.lisp` - Standard library functions (length, map, sum, etc.)
- `self-hosting.lisp` - Self-hosting implementation experiments

### Key Design Patterns

1. **Value Interface**: All Lisp values implement the `Value` interface with a `String()` method
2. **Environment Chain**: Lexical scoping through linked environments
3. **Special Forms**: Core language constructs (if, fn, def, etc.) handled separately from function calls
4. **Error Context**: Rich error reporting with evaluation context and stack traces
5. **Bootstrapped REPL**: The REPL comes pre-loaded with built-ins and standard library

### Data Types Support
- Numbers (integers and floats)
- Strings and symbols
- Keywords (Clojure-style with `:` prefix)
- Lists (linked lists)
- Vectors (indexed collections)
- HashMaps (key-value mappings)
- Functions (built-in and user-defined)

### Testing Strategy
Comprehensive test coverage in `pkg/kernel/*_test.go` files covering:
- Core data types and operations
- Parser functionality
- Evaluation engine
- Environment management
- REPL functionality
- Error handling
- Integration scenarios
- Macro system
- Loop/recur constructs