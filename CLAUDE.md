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
- `types.go` - Core data types (Symbol, Keyword, List, Vector, HashMap, etc.) implementing the `Value` interface, plus comprehensive error handling system
- `reader.go` - Lexer and parser for converting text to AST (Token-based parsing with position tracking)
- `eval_*.go` - Modular evaluation engine split across specialized files:
  - `eval_core.go` - Core evaluation logic, special forms, and context-aware evaluation with stack tracking
  - `eval_arithmetic.go` - Arithmetic operations (+, -, *, /, =, <, >)
  - `eval_collections.go` - Collection operations (cons, first, rest, nth, count, etc.)
  - `eval_strings.go` - String operations (string-split, substring, string-trim, etc.)
  - `eval_io.go` - I/O operations (slurp, spit, println, file-exists?, etc.)
  - `eval_meta.go` - Meta-programming (eval, read-string, macroexpand, gensym, throw, etc.)
  - `eval_special_forms.go` - Special forms (if, fn, def, quote, quasiquote, do, loop, recur, etc.)
- `repl.go` - Enhanced Read-Eval-Print-Loop implementation with multi-line support, dynamic autocomplete, history navigation, and context-aware error reporting
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
3. **Special Forms**: Core language constructs (if, fn, def, quote, quasiquote, loop, recur, etc.) handled separately from function calls
4. **Modular Evaluation**: Core primitives split into focused modules for maintainability
5. **Self-Hosting**: Standard library functions implemented in Lisp using core primitives
6. **Enhanced Error Handling**: Professional-grade error reporting with categorized errors, stack traces, and source context
7. **Bootstrapped REPL**: REPL automatically loads self-hosted standard library
8. **Tail-Call Optimization**: Loop/recur provides efficient tail recursion without stack growth

### Enhanced Error Handling System

GoLisp implements a comprehensive error handling system designed for professional development:

#### Error Categories
- **`ParseError`**: Syntax errors during lexing/parsing with source location and context
- **`TypeError`**: Type mismatches (e.g., `(+ 1 "string")`)
- **`ArityError`**: Wrong number of function arguments (e.g., `(+ 1)`, `(def)`)
- **`NameError`**: Undefined symbols/variables (e.g., `undefined-variable`)
- **`RuntimeError`**: General runtime errors and exceptions
- **`IOError`**: File system and I/O related errors

#### Error Information
- **Source Location**: File name, line number, column number, and character offset
- **Source Context**: Visual display of the error location with pointer (`^`)
- **Stack Traces**: Call stack information for debugging nested function calls
- **Error Chaining**: Support for wrapped/caused-by error relationships
- **Position Tracking**: Unified position tracking across parser and evaluator

#### Error Message Examples

**Parse Error with Visual Context:**
```
ParseError: unexpected end of input at line 3, column 7
(+ 1 2
      ^
```

**Type Error:**
```
TypeError: + expects numbers, got core.String
```

**Arity Error with Location:**
```
ArityError: def expects 2 arguments, got 0 at script.lisp:5:1
```

**Name Error:**
```
NameError: undefined symbol: unknown-function
```

#### Context-Aware Evaluation
- **`EvaluationContext`**: Tracks call stack and source information during evaluation
- **`EvalWithContext`**: Enhanced evaluation function with context tracking
- **Stack Frame Tracking**: Automatic call stack construction for debugging
- **File Context**: REPL and file loading maintain filename context for errors

#### Error System Architecture
- **`LispError`**: Primary error type with rich metadata and formatting
- **Error Preservation**: Error types maintained through evaluation chain
- **Backward Compatibility**: Existing `Eval` function preserved for compatibility
- **Enhanced REPL**: Context-aware evaluation in interactive mode and file loading

#### Developer Experience
- **Clear Error Messages**: Descriptive messages following consistent patterns
- **Visual Error Location**: Source code display with error position indicators
- **Professional Output**: Production-quality error reporting for debugging
- **Non-Crashing REPL**: Errors don't terminate interactive sessions

### Data Types Support
- Numbers (integers and floats)
- Strings and symbols
- Keywords (Clojure-style with `:` prefix)
- Lists (linked lists)
- Vectors (indexed collections)
- HashMaps (key-value mappings with insertion order preservation)
- Sets (unique collections)
- Functions (built-in and user-defined)
- Nil and boolean values

#### HashMap Features
- **Literal syntax**: `{}` for empty, `{:key "value" :other 42}` for populated
- **Constructor**: `(hash-map key1 value1 key2 value2 ...)`
- **Access**: `(get map key)`, `(get map key default)`, or `(:key map)` (keyword as function)
- **Immutable operations**: `(assoc map :key value)`, `(dissoc map :key)`
- **Predicates**: `(hash-map? value)`, `(contains? map key)`, `(empty? map)`
- **Flexible keys**: Any value type can be used as a key
- **Insertion order**: Keys maintain their insertion order

#### Quasiquote System
GoLisp implements a complete quasiquote system compatible with Clojure:
- **Syntax**: `` `expr `` (quasiquote), `~expr` (unquote), `~@expr` (unquote-splicing)
- **Basic quasiquote**: `` `(a b c) `` acts like `'(a b c)` - prevents evaluation
- **Unquote**: `` `(a ~x c) `` evaluates `x` and substitutes its value
- **Unquote-splicing**: `` `(a ~@lst c) `` evaluates `lst` and splices sequence elements
- **Data structure support**: Works with lists, vectors, and hash maps
- **Nested evaluation**: Supports complex expressions like `` `(+ 1 ~(* 2 3)) ``

Examples:
```lisp
(def x 42)
(def lst (list 1 2 3))

`(a b c)        ; → (a b c)
`(a ~x c)       ; → (a 42 c)
`(a ~@lst d)    ; → (a 1 2 3 d)
`[a ~x c]       ; → [a 42 c]
`{:a ~x :b 2}   ; → {:a 42 :b 2}
```

#### Loop/Recur System
GoLisp implements efficient tail-call optimization through `loop` and `recur` special forms:
- **Syntax**: `(loop [binding-vector] body...)` and `(recur new-values...)`
- **Tail optimization**: `recur` jumps back to enclosing `loop` or function without stack growth
- **Arity validation**: `recur` must provide exact number of values as loop bindings
- **Function support**: `recur` works in both `loop` forms and user-defined functions
- **Error handling**: Clear validation of binding syntax and argument arity

Examples:
```lisp
;; Factorial using loop/recur
(loop [n 5 acc 1]
  (if (= n 0)
    acc
    (recur (- n 1) (* acc n))))  ; → 120

;; Tail-recursive function
(defn factorial [n acc]
  (if (= n 0) acc (recur (- n 1) (* acc n))))
(factorial 6 1)  ; → 720

;; Countdown example
(loop [i 3]
  (if (= i 0)
    "done"
    (recur (- i 1))))  ; → "done"
```

### Core Primitives (Go Implementation)
The minimal core provides ~50 essential primitives:

**Arithmetic**: `+`, `-`, `*`, `/`, `=`, `<`, `>`, `<=`, `>=`
**Collections**: `cons`, `first`, `rest`, `nth`, `count`, `empty?`, `conj`, `list`, `vector`, `hash-map`, `set`
**HashMap**: `get`, `assoc`, `dissoc`, `contains?`
**Types**: `symbol?`, `string?`, `number?`, `list?`, `vector?`, `hash-map?`, `set?`, `keyword?`, `fn?`, `nil?`
**Strings**: `str`, `string-split`, `substring`, `string-trim`, `string-replace`
**I/O**: `slurp`, `spit`, `println`, `prn`, `file-exists?`, `list-dir`, `load-file`
**Meta**: `eval`, `read-string`, `read-all-string`, `macroexpand`, `gensym`, `throw`
**Special**: `symbol`, `keyword`, `name`, `throw`
**Quasiquote**: `quasiquote` (`` ` ``), `unquote` (`~`), `unquote-splicing` (`~@`)
**Control Flow**: `loop`, `recur` (tail-call optimization)

### Self-Hosted Standard Library
Higher-level functions implemented in Lisp:

**Logical**: `not`, `when`, `unless`
**Collections**: `map`, `filter`, `reduce`, `apply`, `sort`, `concat`, `any?`, `second`, `third`
**Utilities**: `range`, `join`, `group-by`

### Multi-Expression Support
The GoLisp interpreter provides comprehensive support for handling multiple expressions in source files:

#### File Loading Functions
- **`load-file`**: Loads and evaluates all expressions from a Lisp file in the current environment
  - Usage: `(load-file "filename.lisp")`
  - Returns the value of the last expression in the file
  - All definitions and side effects are applied to the current environment
  - Supports relative and absolute file paths

#### Multi-Expression Parsing
- **`read-all-string`**: Parses multiple expressions from a string
  - Usage: `(read-all-string "(def x 1) (def y 2) (+ x y)")`
  - Returns a list of parsed expressions: `((def x 1) (def y 2) (+ x y))`
  - Handles all data types: lists, vectors, hash-maps, sets, literals
  - Used internally by the self-hosting compiler for source file processing

#### Self-Hosting Compiler Integration
- **`read-all`**: Lisp function that wraps `read-all-string` for compiler use
  - Defined in `lisp/self-hosting.lisp` as `(defn read-all [source] (read-all-string source))`
  - Used by `compile-file` to parse source files with multiple top-level forms
  - Enables the self-hosting compiler to handle realistic Lisp programs

### Enhanced REPL Features

GoLisp provides a sophisticated Read-Eval-Print-Loop with modern interactive development features:

#### Multi-line Expression Support
- **Smart Parentheses Balancing**: Automatically detects incomplete expressions and continues multi-line input
- **Dynamic Prompt Changes**: `GoLisp> ` for new expressions, `      > ` for continuation lines
- **Comment Handling**: Full support for Lisp comments (`;` to end of line) in multi-line expressions
- **String Literal Support**: Proper handling of strings with embedded parentheses and escape sequences
- **Force Evaluation**: Type `)` on an empty line during multi-line input to force evaluation with automatic closing parentheses

#### Dynamic Auto-completion
- **Environment-Aware**: Completion suggestions based on actually loaded functions and symbols
- **117+ Symbols**: Includes all built-in functions and standard library functions
- **Smart Categorization**: Functions get `(functionname` prefix, literals like `nil` get no parentheses
- **Real-time Updates**: Completion list updates as new symbols are defined (framework ready)
- **Comprehensive Coverage**: Includes arithmetic (`+`, `-`), collections (`map`, `filter`), I/O (`println`), meta-programming (`eval`), and more

#### History and Navigation
- **Command History**: Up/down arrow keys to navigate through previous commands
- **Cursor Movement**: Left/right arrow keys for in-line cursor positioning  
- **Session Persistence**: History maintained during REPL session
- **Readline Integration**: Professional terminal handling via `chzyer/readline` library

#### Error Handling and Feedback
- **Context-Aware Errors**: Professional error reporting with source location and context
- **Non-Crashing REPL**: Errors don't terminate the interactive session
- **Input Validation**: Clear feedback for invalid input (e.g., unexpected closing parentheses)
- **Graceful Cancellation**: Ctrl+C cancels multi-line input without exiting REPL

#### Usage Examples
```lisp
GoLisp> (defn factorial [n]
      >   (if (= n 0)
      >       1
      >       (* n (factorial (- n 1)))))
#<function>

GoLisp> (map +
      >      (list 1 2 3)
      >      (list 4 5 6))
[5 7 9]

GoLisp> (factorial 5)
120

GoLisp> ; Tab completion demo: type "(ma" + Tab
GoLisp> (map    ; Auto-completes with opening parenthesis
```

#### Technical Implementation
- **Balanced Expression Detection**: Sophisticated algorithm handling nested structures, strings, and comments
- **Content Analysis**: Distinguishes between meaningful expressions and whitespace/comments
- **Force Evaluation Logic**: Smart bracket counting for safe completion of incomplete expressions
- **Comprehensive Testing**: 79 unit tests covering all REPL functionality and edge cases

### Testing Strategy
Comprehensive test coverage with 4,500+ lines of tests:
- Core data types and operations (`types_test.go`)
- Parser functionality (`reader_test.go`)
- Evaluation engine with new multi-expression functions (`eval_test.go`)
- **Enhanced REPL functionality (`repl_test.go`)**: 79 unit tests covering:
  - Parentheses balancing with 33 test cases (strings, comments, nesting)
  - Content detection with 19 test cases (whitespace, comments, literals)
  - Bracket counting for force evaluation (11 test cases)
  - Evaluation logic with 16 test cases (complete vs incomplete expressions)
  - REPL integration tests (creation, environment access, evaluation)
- Standard library functions (`stdlib_test.go`)
- Integration scenarios (`integration_test.go`)
- Self-hosting compiler integration (`self_hosting_test.go`)
- Multi-expression parsing and file loading tests
- Enhanced error handling system with categorized error testing
- Error context preservation and stack trace functionality
- Parse error reporting with source location and visual context
- Loop/recur tail-call optimization with comprehensive test coverage
- Self-hosting capabilities

### Self-Hosting Features
- **Minimal Core**: 2,719 lines of Go code providing essential primitives
- **Self-Hosted Stdlib**: Standard library functions written in Lisp
- **Compiler Foundation**: Self-hosting compiler in `lisp/self-hosting.lisp`
- **Multi-Expression Support**: Complete support for parsing and compiling multi-expression source files
- **File Loading System**: `load-file` function for loading and evaluating Lisp files
- **Bootstrap Process**: Automatic loading of Lisp-based standard library
- **Meta-Programming**: Full eval/read-string/read-all-string capabilities for self-modification