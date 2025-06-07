# Lisp Interpreter

A feature-rich Lisp interpreter implemented in Go using Test-Driven Development (TDD). This project demonstrates functional programming concepts and provides a complete interactive development environment.

## Quick Start

```bash
# Run the REPL
go run ./cmd/lisp-interpreter

# Or build and run
make build
./lisp-interpreter
```

## Key Features

- **Complete Lisp Implementation**: Full tokenizer, parser, and evaluator
- **Interactive REPL**: Rich development environment with help system
- **Module System**: Organize code with imports and exports
- **Higher-Order Functions**: `map`, `filter`, `reduce`, and more
- **Built-in Help**: Discover functions with `(builtins)` and get detailed help

## Documentation

- **[Features](docs/features.md)** - Complete feature overview and data types
- **[Operations Reference](docs/operations.md)** - Comprehensive guide to all supported operations
- **[Usage Guide](docs/usage.md)** - How to run, build, and use the interpreter
- **[Examples](docs/examples.md)** - Extensive code examples and tutorials
- **[Architecture](docs/architecture.md)** - Technical design and implementation details
- **[Future Enhancements](docs/future.md)** - Planned improvements and roadmap

## Quick Examples

### Basic Operations
```lisp
lisp> (+ 1 2 3)
=> 6

lisp> (define square (lambda (x) (* x x)))
=> #<function([x])>

lisp> (square 5)
=> 25
```

### List Processing
```lisp
lisp> (map (lambda (x) (* x x)) (list 1 2 3 4))
=> (1 4 9 16)

lisp> (filter (lambda (x) (> x 0)) (list -1 2 -3 4))
=> (2 4)
```

### Module System
```lisp
lisp> (module math (export square) (defun square (x) (* x x)))
=> #<module:math>

lisp> (import math)
=> #<module:math>

lisp> (square 5)
=> 25
```

## Building and Testing

### Development Commands
```bash
make build    # Build the interpreter
make run      # Build and run the interpreter  
make test     # Run all tests
```

### Manual Build
```bash
go build -o lisp-interpreter ./cmd/lisp-interpreter
./lisp-interpreter
```

## Project Structure

## Project Structure

Built with Go using clean architecture principles:

```
lisp-interpreter/
├── docs/                        # Comprehensive documentation
├── examples/                    # Example Lisp programs  
├── cmd/lisp-interpreter/        # Main application (REPL)
└── pkg/                         # Core packages
    ├── types/                   # Type definitions
    ├── tokenizer/              # Lexical analysis
    ├── parser/                 # Syntax analysis  
    ├── evaluator/              # Expression evaluation
    └── interpreter/            # High-level API
```

## Contributing

This project welcomes contributions! See the [Architecture](docs/architecture.md) guide for technical details and the [Future Enhancements](docs/future.md) document for planned improvements.

## License

This project is open source. Built as an educational demonstration of Lisp interpreter implementation using test-driven development.
