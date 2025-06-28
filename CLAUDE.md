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
- `make test-core` - Run core package tests only
- `make test-nocache` - Run all tests without cache (useful for debugging)
- `make test-core-nocache` - Run core tests without cache

### Code Quality
- `make fmt` - Format all Go source files
- `go fmt ./...` - Format Go code (alternative)

## Architecture Overview

GoLisp is a minimalist, self-hosting Lisp interpreter written in Go, inspired by Clojure. The architecture features a minimal core with self-hosted standard library.

### Core Components

**`pkg/core/`** - Minimal kernel (2,719 lines total):
- `types.go` - Core data types (Symbol, Keyword, List, Vector, HashMap, etc.) implementing the `Value` interface
- `reader.go` - Lexer and parser for converting text to AST (Token-based parsing with position tracking)
- `eval_*.go` - Modular evaluation engine split across specialized files:
  - `eval_core.go` - Core evaluation logic and special forms
  - `eval_arithmetic.go` - Arithmetic operations (+, -, *, /, =, <, >)
  - `eval_collections.go` - Collection operations (cons, first, rest, nth, count, etc.)
  - `eval_strings.go` - String operations (string-split, substring, string-trim, etc.)
  - `eval_io.go` - I/O operations (slurp, spit, println, file-exists?, etc.)
  - `eval_meta.go` - Meta-programming (eval, read-string, macroexpand, etc.)
  - `eval_special_forms.go` - Special forms (if, fn, def, quote, do, etc.)
- `repl.go` - Read-Eval-Print-Loop implementation
- `bootstrap.go` - Standard library loader and environment initialization

**`cmd/golisp/main.go`** - CLI entry point supporting:
- Interactive REPL mode (default)
- File execution (`-f` flag)
- Direct code evaluation (`-e` flag)
- Help and usage information

**`lisp/`** - Self-hosted Lisp source files:
- `stdlib.lisp` - Legacy minimal standard library
- `stdlib/core.lisp` - Self-hosted standard library (map, filter, reduce, etc.)
- `stdlib/enhanced.lisp` - Enhanced collection operations and utilities
- `self-hosting.lisp` - Self-hosting compiler implementation

### Key Design Patterns

1. **Value Interface**: All Lisp values implement the `Value` interface with a `String()` method
2. **Environment Chain**: Lexical scoping through linked environments
3. **Special Forms**: Core language constructs (if, fn, def, etc.) handled separately from function calls
4. **Modular Evaluation**: Core primitives split into focused modules for maintainability
5. **Self-Hosting**: Standard library functions implemented in Lisp using core primitives
6. **Error Context**: Rich error reporting with evaluation context and stack traces
7. **Bootstrapped REPL**: REPL automatically loads self-hosted standard library

### Data Types Support
- Numbers (integers and floats)
- Strings and symbols
- Keywords (Clojure-style with `:` prefix)
- Lists (linked lists)
- Vectors (indexed collections)
- HashMaps (key-value mappings)
- Sets (unique collections)
- Functions (built-in and user-defined)
- Nil and boolean values

### Core Primitives (Go Implementation)
The minimal core provides ~50 essential primitives:

**Arithmetic**: `+`, `-`, `*`, `/`, `=`, `<`, `>`, `<=`, `>=`
**Collections**: `cons`, `first`, `rest`, `nth`, `count`, `empty?`, `conj`, `list`, `vector`, `hash-map`, `set`
**Types**: `symbol?`, `string?`, `number?`, `list?`, `vector?`, `hash-map?`, `set?`, `keyword?`, `fn?`, `nil?`
**Strings**: `str`, `string-split`, `substring`, `string-trim`, `string-replace`
**I/O**: `slurp`, `spit`, `println`, `prn`, `file-exists?`, `list-dir`
**Meta**: `eval`, `read-string`, `macroexpand`, `gensym`
**Special**: `symbol`, `keyword`, `name`, `throw`

### Self-Hosted Standard Library
Higher-level functions implemented in Lisp:

**Logical**: `not`, `when`, `unless`
**Collections**: `map`, `filter`, `reduce`, `apply`, `sort`, `concat`, `any?`, `second`, `third`
**Utilities**: `range`, `join`, `group-by`

### Testing Strategy
Comprehensive test coverage with 4,319 lines of tests:
- Core data types and operations (`types_test.go`)
- Parser functionality (`reader_test.go`)
- Evaluation engine (`eval_test.go`)
- Standard library functions (`stdlib_test.go`)
- Integration scenarios (`integration_test.go`)
- Error handling and edge cases
- Self-hosting capabilities

### Self-Hosting Features
- **Minimal Core**: 2,719 lines of Go code providing essential primitives
- **Self-Hosted Stdlib**: Standard library functions written in Lisp
- **Compiler Foundation**: Self-hosting compiler in `lisp/self-hosting.lisp`
- **Bootstrap Process**: Automatic loading of Lisp-based standard library
- **Meta-Programming**: Full eval/read-string capabilities for self-modification