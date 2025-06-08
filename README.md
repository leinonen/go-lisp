# Lisp Interpreter

A feature-rich Lisp interpreter implemented in Go using Test-Driven Development (TDD). This project demonstrates functional programming concepts and provides a complete interactive development environment.

## Quick Start

```bash
# Run the REPL
go run ./cmd/lisp-interpreter

# Run a Lisp file
go run ./cmd/lisp-interpreter myprogram.lisp

# Or build and run
make build
./lisp-interpreter
./lisp-interpreter myprogram.lisp
```

## Key Features

- **Complete Lisp Implementation**: Full tokenizer, parser, and evaluator
- **Interactive REPL**: Rich development environment with help system
- **File Execution**: Run Lisp programs from files with full multi-expression support
- **Big Number Arithmetic**: Arbitrary precision integers with automatic overflow detection
- **Modulo Operator**: Full modulo support (`%`) with big number compatibility
- **Error Handling**: Built-in `error` function for controlled program termination
- **Core Library**: Rich set of mathematical and utility functions (factorial, abs, gcd, etc.)
- **Tail Call Optimization**: Prevents stack overflow in recursive functions
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

### New Features Documentation

- **[Modulo Operator](docs/modulo_operator.md)** - Complete guide to the `%` operator
- **[Error Function](docs/error_function.md)** - Error handling with the `error` function
- **[File Execution](docs/file_execution.md)** - Running Lisp programs from files
- **[Core Library](docs/core_library.md)** - Mathematical and utility functions

## Quick Examples

### Basic Operations
```lisp
lisp> (+ 1 2 3)
=> 6

lisp> (define square (lambda (x) (* x x)))
=> #<function([x])>

lisp> (square 5)
=> 25

; New modulo operator
lisp> (% 17 5)
=> 2

lisp> (% 1000000000000000001 7)
=> 0
```

### Error Handling
```lisp
lisp> (error "Something went wrong!")
error: Something went wrong!

lisp> (if (< x 0) (error "Negative values not allowed") (sqrt x))
```

### Core Library Functions
```lisp
; Mathematical functions
lisp> (factorial 5)
=> 120

lisp> (abs -42)
=> 42

lisp> (gcd 48 18)
=> 6

; List utilities  
lisp> (all (lambda (x) (> x 0)) (list 1 2 3))
=> #t

lisp> (any (lambda (x) (< x 0)) (list 1 -2 3))
=> #t
```

### File Execution
```bash
# Create a Lisp program file
echo '(define factorial (lambda (n) (if (= n 0) 1 (* n (factorial (- n 1))))))
(factorial 5)
(+ 10 20)' > math.lisp

# Run it
./lisp-interpreter math.lisp
# Output: 120
#         30
```

### List Processing
```lisp
lisp> (map (lambda (x) (* x x)) (list 1 2 3 4))
=> (1 4 9 16)

lisp> (filter (lambda (x) (> x 0)) (list -1 2 -3 4))
=> (2 4)
```

### Tail Call Optimization  
```lisp
; Tail-recursive factorial won't cause stack overflow
lisp> (defun fact-tail (n acc)
        (if (= n 0) acc (fact-tail (- n 1) (* n acc))))
=> #<function([n acc])>

lisp> (fact-tail 1000 1)  ; Handles large recursion efficiently
=> 4023872600770937735...  ; (very large number)
```

### Big Number Arithmetic
```lisp
; Automatic precision handling for large integers
lisp> (* 1000000000000000 1000000000000000)
=> 1000000000000000000000000000000

; Modulo works with big numbers too
lisp> (% 123456789012345678901234567890 7)
=> 4

lisp> (fact-tail 50 1)  ; Large factorials work seamlessly
=> 30414093201713378043612608166064768844377641568960512000000000000
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

### Core Library Access
```lisp
; Core library functions are automatically available
lisp> (factorial 6)
=> 720

lisp> (gcd 24 36)
=> 12

; Load additional functions from examples/core.lisp
lisp> (load "examples/core.lisp")
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

Built with Go using clean architecture principles:

```
lisp-interpreter/
├── docs/                        # Comprehensive documentation
│   ├── modulo_operator.md      # Modulo operator guide
│   ├── error_function.md       # Error handling documentation  
│   ├── file_execution.md       # File execution guide
│   └── core_library.md         # Core library reference
├── examples/                    # Example Lisp programs
│   └── core.lisp              # Core library functions
├── cmd/lisp-interpreter/        # Main application (REPL + file runner)
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
