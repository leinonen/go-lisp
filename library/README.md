# Lisp Interpreter Library

This directory contains higher-level utilities and extensions built on top of the core interpreter functionality. The library is organized into specialized modules for different programming needs.

## Library Structure

### `core.lisp`
Core mathematical and utility functions that extend the built-in operations:
- **Mathematical Functions**: Advanced math operations and algorithms
- **List Utilities**: Higher-order list processing functions
- **Functional Programming**: Composition and application utilities

### `strings.lisp` 
String manipulation utilities built on top of the built-in string primitives:
- **Convenience Aliases**: Shorter names for common operations
- **Enhanced Validation**: Pattern checking and type validation
- **Text Processing**: Word and line processing utilities
- **Transformation**: Advanced string manipulation and formatting

### `macros.lisp` ⭐ NEW
Comprehensive macro library providing powerful language extensions:
- **Control Flow**: `when`, `unless`, `cond`, short-circuiting operators
- **Variable Binding**: `let1`, `let*`, assignment and modification macros
- **Debugging Tools**: `debug`, `trace`, `assert`, execution timing
- **Iteration**: `dotimes`, `while`, `for-each` loop constructs
- **Utilities**: Multiple values, sequential execution, pattern matching

## Architecture

### Built-in Functions (Go Primitives)
Core operations implemented directly in Go for performance:
- **Arithmetic**: `+`, `-`, `*`, `/`, `%` with big number support
- **Comparison**: `=`, `<`, `>`, `<=`, `>=` for all types
- **String Operations**: 20+ string functions for manipulation and analysis
- **List Operations**: `list`, `cons`, `first`, `rest`, `length`, `empty?`
- **Hash Maps**: `hash-map`, `hash-map-get`, `hash-map-put`, etc.
- **I/O**: `print!`, `load`, module system operations

### Library Functions (Lisp Compositions)
Higher-level utilities that combine the primitives for common tasks and provide more convenient interfaces.

## Usage

```lisp
; Built-in functions are always available
(string-concat "Hello" " " "World")  ; => "Hello World"
(string-length "Hello")              ; => 5
(string-upper "hello")               ; => "HELLO"

; Load the library for higher-level functions
(load "library/strings.lisp")
(import strings)

; Use composed functions
(str-capitalize "hello world")       ; => "Hello world"
(str-title-case "hello world")       ; => "Hello World"
(str-reverse "hello")                ; => "olleh"
(str-numeric? "123")                 ; => true
```

## Usage Examples

### Core Library
```lisp
; Load the core mathematical functions
(load "library/core.lisp")
(import core)

; Use advanced mathematical functions
(factorial 10)                      ; => 3628800
(fibonacci 15)                      ; => 610
(gcd 48 18)                         ; => 6
(lcm 12 8)                          ; => 24

; Higher-order list processing
(all (fn (x) (> x 0)) '(1 2 3)) ; => true
(any (fn (x) (< x 0)) '(1 -2 3)) ; => true
(take 3 '(1 2 3 4 5))               ; => (1 2 3)
(drop 2 '(1 2 3 4 5))               ; => (3 4 5)
```

### String Library
```lisp
; Built-in functions are always available
(string-concat "Hello" " " "World")  ; => "Hello World"
(string-length "Hello")              ; => 5
(string-upper "hello")               ; => "HELLO"

; Load the library for higher-level functions
(load "library/strings.lisp")
(import strings)

; Use composed functions
(str-capitalize "hello world")       ; => "Hello world"
(str-title-case "hello world")       ; => "Hello World"
(str-reverse "hello")                ; => "olleh"
(str-numeric? "123")                 ; => true
(str-words "hello world test")       ; => ("hello" "world" "test")
```

### Macro Library
```lisp
; Load the macro library
(load "library/macros.lisp")

; Control flow macros
(when (> x 5) (print! "x is big"))
(unless (< x 0) (print! "x is not negative"))

; Variable binding
(let* ((x 1) (y (+ x 1)) (z (+ x y))) 
  (+ x y z))                         ; => 5

; Debugging
(debug (+ 1 2 3))                    ; Prints: DEBUG: (+ 1 2 3) => 6

; Iteration
(dotimes i 5 (print! i))              ; Prints: 0 1 2 3 4

; Pattern matching
(match day 
  (1 "Monday") 
  (2 "Tuesday") 
  (_ "Other day"))

; Assignment macros
(def x 10)
(incf x 5)                           ; x is now 15
(decf x 3)                           ; x is now 12
```

## Loading Multiple Libraries

```lisp
; Load all libraries at once
(load "library/core.lisp")
(load "library/strings.lisp") 
(load "library/macros.lisp")

; Now you can use functions and macros from all libraries
(let* ((numbers '(1 2 3 4 5))
       (doubled (map (fn (x) (* x 2)) numbers))
       (result-str (str-join doubled ", ")))
  (debug result-str))                ; DEBUG: result-str => "2, 4, 6, 8, 10"
```

## Benefits of This Architecture

1. **Performance**: Critical operations are implemented in Go for speed
2. **Extensibility**: Higher-level operations can be easily added in Lisp
3. **Modularity**: Users can choose to load only the libraries they need
4. **Readability**: Library functions provide more descriptive names and behaviors
5. **Educational**: Shows how to build complex operations from simple primitives
6. **Language Extension**: Macros allow users to create custom syntax and control structures

The built-in functions provide the essential building blocks, while the library functions demonstrate idiomatic patterns and provide convenient abstractions for common tasks. The macro library enables powerful metaprogramming capabilities for creating domain-specific languages and custom control structures.
