# Architecture

The Lisp interpreter follows a traditional three-phase design with a modular, test-driven architecture.

## Core Design

The interpreter implements a classic three-phase pipeline:

1. **Tokenization**: Convert input string into tokens
2. **Parsing**: Build Abstract Syntax Tree (AST) from tokens
3. **Evaluation**: Evaluate AST in an environment context

## Test-Driven Development

This project was built using TDD principles:

- **Tests First**: Tests were written before implementation for each component
- **Red-Green-Refactor**: Implementation followed to make tests pass, then refactored
- **Living Documentation**: Tests serve as executable documentation and ensure correctness
- **Regression Prevention**: Comprehensive test suite prevents breaking changes

## File Structure

The project follows Go standard layout with clear separation of concerns:

```
lisp-interpreter/
├── go.mod                       # Go module definition
├── README.md                    # Project overview and quick start
├── Makefile                     # Build automation
├── docs/                        # Detailed documentation
│   ├── features.md              # Feature overview
│   ├── operations.md            # Operation reference
│   ├── usage.md                 # Usage guide
│   ├── examples.md              # Code examples
│   ├── architecture.md          # This file
│   └── future.md                # Planned enhancements
├── examples/                    # Example Lisp programs
│   ├── math_module.lisp         # Module system examples
│   ├── higher_order_functions.lisp # HOF examples
│   └── ...                     # Additional examples
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
    │   ├── module_test.go       # Module system tests
    │   └── evaluator_test.go    # Evaluator tests
    └── interpreter/
        ├── interpreter.go       # High-level interpreter API
        └── interpreter_test.go  # Integration tests
```

## Package Overview

### Core Packages

- **`pkg/types`** - Core type definitions and interfaces
  - Defines `Token`, `Expr`, and `Value` interfaces
  - Provides type safety and abstraction boundaries
  - Contains shared constants and enums

- **`pkg/tokenizer`** - Lexical analysis (string → tokens)
  - Converts raw input into structured tokens
  - Handles comments, strings, numbers, symbols
  - Provides error reporting for invalid input

- **`pkg/parser`** - Syntax analysis (tokens → AST)
  - Builds Abstract Syntax Tree from token stream
  - Implements recursive descent parsing
  - Handles operator precedence and associativity

- **`pkg/evaluator`** - Expression evaluation (AST → values)
  - Evaluates expressions in environment contexts
  - Implements built-in functions and special forms
  - Manages variable scoping and function calls
  - Handles module system and imports
  - **Tail Call Optimization**: Automatically optimizes tail-recursive calls to prevent stack overflow

- **`pkg/interpreter`** - High-level API
  - Combines all components into unified interface
  - Provides convenient methods for evaluation
  - Manages REPL state and error handling

### Application Layer

- **`cmd/lisp-interpreter`** - Main application
  - Implements interactive REPL
  - Handles command-line arguments
  - Provides user-friendly error messages
  - Manages session state and history

## Design Principles

### Separation of Concerns

Each package has a single, well-defined responsibility:
- Tokenizer only handles lexical analysis
- Parser only handles syntax analysis  
- Evaluator only handles semantic evaluation
- Clear interfaces between components

### Interface-Driven Design

Key abstractions are defined as interfaces:
- `Token` interface for different token types
- `Expr` interface for AST nodes
- `Value` interface for runtime values
- Enables extensibility and testing

### Immutable Data Structures

Most data structures are immutable:
- Tokens are read-only after creation
- AST nodes don't change after parsing
- Values are immutable (except environments)
- Reduces bugs and enables safe concurrency

### Error Handling

Comprehensive error handling throughout:
- Lexical errors (invalid characters)
- Syntax errors (malformed expressions)
- Runtime errors (undefined variables, type mismatches)
- Clear error messages with context

## Testing Strategy

### Unit Tests

Each package has comprehensive unit tests:
- **Tokenizer tests**: Verify correct token generation
- **Parser tests**: Verify correct AST construction
- **Evaluator tests**: Verify correct evaluation results
- **Type tests**: Verify interface implementations

### Integration Tests

The interpreter package provides integration tests:
- End-to-end evaluation of complex expressions
- REPL functionality testing
- Module system integration
- Error handling across components

### Test Coverage

Tests cover:
- Happy path scenarios
- Edge cases and error conditions
- Performance characteristics
- Regression prevention

## Performance Considerations

### Memory Management

- Minimal allocations in hot paths
- Reuse of common structures where possible
- Garbage collection friendly design
- **Big Numbers**: Use `math/big` package for arbitrary precision arithmetic

### Parsing Efficiency

- Single-pass tokenization
- Recursive descent parsing (O(n) complexity)
- Minimal backtracking
- **Big Number Detection**: Automatic detection of large integers during parsing

### Evaluation Optimization

- Direct AST evaluation (no intermediate compilation)
- Efficient environment lookup
- **Tail Call Optimization**: Implemented to prevent stack overflow in recursive functions
- **Overflow Detection**: Automatic promotion to big integers when operations would overflow

## Big Number Support

### Technical Design

The interpreter provides comprehensive support for arbitrary precision integers using Go's `math/big` package:

1. **Automatic Detection**: Large integers (≥1e15) are automatically detected during parsing
2. **Overflow Protection**: Arithmetic operations detect potential overflow and promote to big integers
3. **Seamless Integration**: Big numbers work transparently with regular integers in all operations
4. **Display Formatting**: Large numbers are formatted in standard decimal notation

### Implementation Details

```go
type BigNumberValue struct {
    value *big.Int
}

type BigNumberExpr struct {
    Value *big.Int
}
```

### Key Components

- **`BigNumberValue`**: Runtime value type for arbitrary precision integers
- **`BigNumberExpr`**: AST node for big number literals
- **Overflow Detection**: `mightOverflowMultiplication()` checks for potential overflow
- **Type Conversion**: Automatic promotion between regular and big integers
- **Arithmetic Enhancement**: All operations (`+`, `-`, `*`, `/`, `%`) support big numbers
- **Comparison Support**: All comparison operations (`=`, `<`, `>`, `<=`, `>=`) work with big numbers

### Parsing Strategy

The parser automatically detects large numbers using multiple criteria:
- Scientific notation (contains 'e' or 'E')
- String length (>15 characters)
- Magnitude threshold (≥1e15)

### Benefits

- **Precision**: No loss of precision for large integer calculations
- **Transparency**: Existing code continues to work without changes
- **Performance**: Regular integers used for small values to maintain performance
- **Usability**: Automatic promotion means users don't need to think about number types

## Tail Call Optimization Implementation

### Technical Design

The TCO implementation uses an iterative approach to eliminate stack growth for tail-recursive functions:

1. **Tail Call Detection**: The evaluator identifies when a function call is in tail position
2. **Iterative Execution**: Instead of recursive calls, uses a loop to reuse the same stack frame
3. **Argument Evaluation**: Arguments are evaluated once and reused in the optimization loop
4. **Semantic Preservation**: Non-tail recursive functions maintain normal call semantics

### Implementation Details

```go
type TailCallInfo struct {
    Function types.Value    // The function to call
    Args     []types.Value  // Pre-evaluated arguments
}

type Evaluator struct {
    // ...existing fields...
    tailCall   *TailCallInfo  // Current tail call information
    tailCallOK bool           // Whether tail calls are enabled
}
```

### Key Components

- **`callFunction`**: Main function call handler with iterative tail call loop
- **`callFunctionWithTailCheck`**: Used for calls in tail position
- **`callFunctionRegular`**: Used for calls in non-tail position  
- **`evalFunctionCall`**: Intelligently chooses appropriate call method
- **`evalIf`**: Properly handles tail calls in conditional branches

### Benefits

- **Stack Safety**: Large recursive computations don't cause stack overflow
- **Performance**: Linear space complexity instead of exponential for tail calls
- **Transparency**: No special syntax required - optimization is automatic
- **Correctness**: Preserves all language semantics for non-tail recursive code

## Extensibility Points

The architecture supports future extensions:

### New Data Types
- Add new `Value` implementations
- Extend tokenizer for new literals
- Update parser for new syntax

### New Built-in Functions
- Add functions to evaluator's built-in registry
- Automatic help system integration
- Consistent error handling

### Language Features
- Module system provides namespace management
- Macro system foundation exists
- Meta-programming capabilities possible

### Performance Improvements
- Bytecode compilation layer can be added
- JIT compilation possible
- Parallel evaluation for pure functions
