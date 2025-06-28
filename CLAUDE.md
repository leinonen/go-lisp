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
  - `eval_special_forms.go` - Special forms (if, fn, def, quote, quasiquote, do, etc.)
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
3. **Special Forms**: Core language constructs (if, fn, def, quote, quasiquote, etc.) handled separately from function calls
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

### Core Primitives (Go Implementation)
The minimal core provides ~50 essential primitives:

**Arithmetic**: `+`, `-`, `*`, `/`, `=`, `<`, `>`, `<=`, `>=`
**Collections**: `cons`, `first`, `rest`, `nth`, `count`, `empty?`, `conj`, `list`, `vector`, `hash-map`, `set`
**HashMap**: `get`, `assoc`, `dissoc`, `contains?`
**Types**: `symbol?`, `string?`, `number?`, `list?`, `vector?`, `hash-map?`, `set?`, `keyword?`, `fn?`, `nil?`
**Strings**: `str`, `string-split`, `substring`, `string-trim`, `string-replace`
**I/O**: `slurp`, `spit`, `println`, `prn`, `file-exists?`, `list-dir`, `load-file`
**Meta**: `eval`, `read-string`, `read-all-string`, `macroexpand`, `gensym`
**Special**: `symbol`, `keyword`, `name`, `throw`
**Quasiquote**: `quasiquote` (`` ` ``), `unquote` (`~`), `unquote-splicing` (`~@`)

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

### Testing Strategy
Comprehensive test coverage with 4,500+ lines of tests:
- Core data types and operations (`types_test.go`)
- Parser functionality (`reader_test.go`)
- Evaluation engine with new multi-expression functions (`eval_test.go`)
- Standard library functions (`stdlib_test.go`)
- Integration scenarios (`integration_test.go`)
- Self-hosting compiler integration (`self_hosting_test.go`)
- Multi-expression parsing and file loading tests
- Error handling and edge cases
- Self-hosting capabilities

### Self-Hosting Features
- **Minimal Core**: 2,719 lines of Go code providing essential primitives
- **Self-Hosted Stdlib**: Standard library functions written in Lisp
- **Compiler Foundation**: Self-hosting compiler in `lisp/self-hosting.lisp`
- **Multi-Expression Support**: Complete support for parsing and compiling multi-expression source files
- **File Loading System**: `load-file` function for loading and evaluating Lisp files
- **Bootstrap Process**: Automatic loading of Lisp-based standard library
- **Meta-Programming**: Full eval/read-string/read-all-string capabilities for self-modification