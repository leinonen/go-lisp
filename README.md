# GoLisp

![GoLisp logo](./docs/img/golisp-logo.png)

A minimalist, self-hosting Lisp interpreter written in Go, inspired by Clojure. Features a modular core with self-hosted standard library.

## Features

- **Minimal Core**: ~2,700 lines of focused Go code providing essential primitives
- **Self-Hosting**: Standard library functions implemented in Lisp
- **Rich Data Types**: Numbers, strings, symbols, keywords, lists, vectors, hash-maps, sets
- **Functional Programming**: First-class functions, closures, and higher-order functions
- **Advanced Language Features**: `defn`, `defmacro`, `cond`, multiple body expressions
- **Macro System**: Full macro expansion with `defmacro` and quasiquote support
- **Meta-Programming**: Full `eval`/`read-string` capabilities with macro system
- **Self-Hosting Compiler**: Integrated compiler for self-compilation capabilities
- **REPL**: Interactive development environment
- **File Operations**: Load and execute Lisp files
- **Comprehensive Testing**: 4,300+ lines of test coverage

## Quick Start

### Installation

```bash
git clone https://github.com/leinonen/go-lisp
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
(defn square [x] (* x x))            ; define function (using defn)
(square 5)                           ; 25

(def numbers [1 2 3 4 5])            ; vector
(def person {:name "Alice" :age 30}) ; hash-map
```

### Advanced Language Features
```lisp
;; Conditional expressions
(cond
  (< x 0) "negative"
  (= x 0) "zero"
  :else   "positive")

;; Macros (variadic parameters with &)
(defmacro when [condition & body]
  (list 'if condition (cons 'do body) nil))

(defmacro unless [condition & body]
  (list 'if condition nil (cons 'do body)))

;; Alternative quasiquote syntax
(defmacro when [condition & body]
  `(if ~condition (do ~@body) nil))

;; Quasiquote for template construction
(def x 42)
(def lst (list 1 2 3))
`(a ~x c)                           ; (a 42 c) - unquote substitution
`(a ~@lst d)                        ; (a 1 2 3 d) - unquote-splicing
`{:value ~x :type "number"}         ; {:value 42 :type "number"}

;; Multiple body expressions
(defn complex-function [x]
  (println "Processing" x)
  (def result (* x x))
  (println "Result:" result)
  result)
```

### Higher-Order Functions (Self-Hosted)
```lisp
(map (fn [x] (* x 2)) [1 2 3 4])      ; (2 4 6 8)
(filter (fn [x] (> x 2)) [1 2 3 4 5]) ; (3 4 5)
(reduce + 0 [1 2 3 4 5])              ; 15
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

;; Macro expansion
(macroexpand '(when true (println "hello")))
;; => (if true (do (println "hello")) nil)
```

### Self-Hosting Compiler
```lisp
;; Load the self-hosting compiler
;; The compiler can compile Lisp code to Lisp code
(def ctx (make-context))
(compile-expr '(+ 1 2) ctx)        ; Compile arithmetic expression
(compile-expr 'my-var ctx)         ; Compile symbol reference
```

## Architecture

GoLisp uses a modular architecture with clear separation between the Go core and Lisp standard library:

### Core (Go Implementation)
- **Types & Parser**: Essential data types and parsing with macro support
- **Evaluator**: Modular evaluation engine (~60 core primitives including special forms
- **Special Forms**: `def`, `fn`, `defn`, `defmacro`, `cond`, `if`, `let`, `do`, `quote`
- **Macro System**: Full macro expansion with `defmacro` and macro call evaluation
- **REPL**: Interactive environment with file loading

### Standard Library (Lisp Implementation)
- **Collections**: `map`, `filter`, `reduce`, `sort`, `apply`, `length`
- **Logic**: `not`, `when`, `unless`, `cond` (enhanced)
- **Utilities**: `range`, `join`, `group-by`, `hash-map-put`
- **Error Handling**: `throw` for runtime error generation

### Self-Hosting Compiler (Lisp Implementation)
- **Compilation Context**: Environment and symbol table management
- **Expression Compilation**: Handles symbols, lists, vectors with context
- **Special Form Compilation**: `def`, `fn`, `if`, `quote`, `do`, `let`
- **Code Generation**: Lisp-to-Lisp compilation with optimization hooks

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

lisp/
├── stdlib/         # Self-hosted standard library
│   ├── core.lisp   # Core functions in Lisp
│   └── enhanced.lisp # Enhanced utilities
└── self-hosting.lisp # Self-hosting compiler

cmd/golisp/         # CLI entry point
```

## Self-Hosting

GoLisp is designed for self-hosting with:
- Minimal Go core providing essential primitives (~60 functions)
- Standard library implemented in Lisp using core primitives
- Advanced language features: `defn`, `defmacro`, `cond`, multiple body expressions
- Full macro system with `defmacro` and macro expansion
- Meta-programming capabilities (`eval`, `read-string`, `macroexpand`)
- Self-hosting compiler in `lisp/self-hosting.lisp` (functional and integrated)

### Current Self-Hosting Status
- ✅ **Phase 1**: Meta-programming core complete
- ✅ **Phase 2**: Enhanced standard library complete  
- ✅ **Phase 3.1**: Basic self-hosting compiler integration complete
- 🚧 **Phase 3.2**: Advanced compiler features in progress

The self-hosting compiler can currently:
- Load and compile basic Lisp expressions
- Handle compilation contexts and symbol tables
- Compile special forms (`def`, `fn`, `if`, `quote`, `do`)
- Support for basic optimization hooks

This architecture enables language evolution in Lisp rather than Go, making it easy to extend and modify the language itself.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `make test`
5. Submit a pull request

## License

MIT
