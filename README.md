# GoLisp

![GoLisp logo](./docs/img/golisp-logo.png)

A minimalist, self-hosting Lisp interpreter written in Go, inspired by Clojure. Features a modular core with self-hosted standard library.

## Features

- **Minimal Core**: ~2,700 lines of focused Go code providing essential primitives
- **Self-Hosting**: Standard library functions implemented in Lisp
- **Rich Data Types**: Numbers, strings, symbols, keywords, lists, vectors, hash-maps, sets
- **Functional Programming**: First-class functions, closures, and higher-order functions
- **Meta-Programming**: Full `eval`/`read-string` capabilities with macro system
- **REPL**: Interactive development environment
- **File Operations**: Load and execute Lisp files
- **Comprehensive Testing**: 4,300+ lines of test coverage

## Quick Start

### Installation

```bash
git clone https://github.com/yourusername/go-lisp
cd go-lisp
make build
```

### Usage

```bash
# Interactive REPL
./bin/golisp

# Execute a file
./bin/golisp -f script.lisp

# Evaluate expression directly
./bin/golisp -e '(+ 1 2 3)'
```

## Examples

### Basic Operations
```lisp
(+ 1 2 3)                          ; 6
(* 2 3 4)                          ; 24
(= 5 (+ 2 3))                      ; true
```

### Functions and Variables
```lisp
(def square (fn [x] (* x x)))      ; define function
(square 5)                         ; 25

(def numbers [1 2 3 4 5])          ; vector
(def person {:name "Alice" :age 30}) ; hash-map
```

### Higher-Order Functions (Self-Hosted)
```lisp
(map (fn [x] (* x 2)) [1 2 3 4])   ; (2 4 6 8)
(filter (fn [x] (> x 2)) [1 2 3 4 5]) ; (3 4 5)
(reduce + 0 [1 2 3 4 5])           ; 15
```

### Collections
```lisp
[1 2 3]                            ; vectors
'(1 2 3)                           ; lists  
{:name "Bob" :age 25}              ; hash-maps
#{1 2 3}                           ; sets
```

### Meta-Programming
```lisp
(eval '(+ 1 2 3))                  ; 6
(read-string "(+ 1 2)")            ; (+ 1 2)
```

## Architecture

GoLisp uses a modular architecture with clear separation between the Go core and Lisp standard library:

### Core (Go Implementation)
- **Types & Parser**: Essential data types and parsing
- **Evaluator**: Modular evaluation engine (~50 core primitives)
- **REPL**: Interactive environment with file loading

### Standard Library (Lisp Implementation)
- **Collections**: `map`, `filter`, `reduce`, `sort`, `apply`
- **Logic**: `not`, `when`, `unless`
- **Utilities**: `range`, `join`, `group-by`

## Development

### Building and Testing
```bash
make build          # Build interpreter
make test           # Run all tests
make test-core      # Test core package only
make fmt            # Format Go code
```

### Project Structure
```
pkg/core/           # Minimal Go core (~2,700 lines)
├── types.go        # Data types and environment
├── reader.go       # Parser and lexer
├── eval_*.go       # Modular evaluator
├── repl.go         # Interactive REPL
└── bootstrap.go    # Standard library loader

lisp/stdlib/        # Self-hosted standard library
├── core.lisp       # Core functions in Lisp
└── enhanced.lisp   # Enhanced utilities

cmd/golisp/         # CLI entry point
```

## Self-Hosting

GoLisp is designed for self-hosting with:
- Minimal Go core providing essential primitives
- Standard library implemented in Lisp using core primitives
- Meta-programming capabilities (`eval`, `read-string`, macros)
- Self-hosting compiler foundation in `lisp/self-hosting.lisp`

This architecture enables language evolution in Lisp rather than Go, making it easy to extend and modify the language itself.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `make test`
5. Submit a pull request

## License

MIT
