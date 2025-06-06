# Lisp Interpreter

A basic Lisp interpreter implemented in Go using Test-Driven Development (TDD).

## Features

- **Tokenizer/Lexer**: Converts input text into tokens
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens
- **Evaluator**: Evaluates the AST to produce results
- **REPL**: Interactive Read-Eval-Print Loop

## Supported Operations

### Arithmetic
- `(+ 1 2 3)` - Addition with multiple operands
- `(- 10 3)` - Subtraction
- `(* 2 3 4)` - Multiplication with multiple operands
- `(/ 15 3)` - Division

### Comparison
- `(= 5 5)` - Equality
- `(< 3 5)` - Less than
- `(> 7 3)` - Greater than

### Conditional
- `(if condition then-expr else-expr)` - If expression

### Variables
- `(define name value)` - Define a variable with a name and value

### Data Types
- Numbers: `42`, `-3.14`
- Strings: `"hello world"`
- Booleans: `#t`, `#f`
- Symbols: `+`, `-`, `x`, `my-var`

## Usage

### Running the REPL
```bash
go run ./cmd/lisp-interpreter
```

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o lisp-interpreter ./cmd/lisp-interpreter
./lisp-interpreter
```

### Using the Makefile
```bash
make build    # Build the interpreter
make run      # Build and run the interpreter
make test     # Run all tests
```

## Examples

```lisp
lisp> 42
=> 42

lisp> (+ 1 2 3)
=> 6

lisp> (* (+ 2 3) 4)
=> 20

lisp> (if (< 3 5) 100 0)
=> 100

lisp> "hello world"
=> hello world

lisp> (= 5 5)
=> #t

lisp> (define x 10)
=> 10

lisp> x
=> 10

lisp> (define y (* x 3))
=> 30

lisp> (+ x y)
=> 40
```

## Architecture

The interpreter follows a traditional three-phase design:

1. **Tokenization**: Convert input string into tokens
2. **Parsing**: Build AST from tokens
3. **Evaluation**: Evaluate AST in an environment

### Test-Driven Development

This project was built using TDD principles:
- Tests were written first for each component
- Implementation followed to make tests pass
- Tests serve as documentation and ensure correctness

### File Structure

The project follows Go standard layout with modular architecture:

```
lisp-interpreter/
├── go.mod                       # Go module definition
├── README.md                    # Project documentation
├── Makefile                     # Build automation
├── cmd/
│   └── lisp-interpreter/
│       └── main.go              # REPL and main program
└── pkg/
    ├── types/
    │   ├── types.go             # Core types and interfaces
    │   └── types_test.go        # Type system tests
    ├── tokenizer/
    │   ├── tokenizer.go         # Lexical analysis
    │   └── tokenizer_test.go    # Tokenizer tests
    ├── parser/
    │   ├── parser.go            # Syntax analysis
    │   └── parser_test.go       # Parser tests
    ├── evaluator/
    │   ├── evaluator.go         # Expression evaluation
    │   └── evaluator_test.go    # Evaluator tests
    └── interpreter/
        ├── interpreter.go       # High-level interpreter API
        └── interpreter_test.go  # Integration tests
```

### Package Overview

- **`pkg/types`** - Core type definitions (Token, Expr, Value interfaces)
- **`pkg/tokenizer`** - Lexical analysis (string → tokens)
- **`pkg/parser`** - Syntax analysis (tokens → AST)
- **`pkg/evaluator`** - Expression evaluation (AST → values)
- **`pkg/interpreter`** - High-level API combining all components
- **`cmd/lisp-interpreter`** - Main application with REPL

## Future Enhancements

- Function definitions with `(lambda (args) body)`
- More built-in functions
- Error recovery in parser
- Better error messages
- Support for comments
