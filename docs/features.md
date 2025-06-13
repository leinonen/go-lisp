# Features

A modern, production-ready Lisp interpreter with comprehensive language support and advanced capabilities.

## Core Language

**Basic**: Variables, functions, lists, arithmetic, comparisons, comments  
**Advanced**: Closures, recursion, tail-call optimization, error handling  
**Data Types**: Numbers (including big integers), strings, booleans, lists, hash maps, keywords  
**Modern Syntax**: Square bracket function parameters for improved readability and reduced confusion  

## Modern Capabilities

**Functional Programming**: Higher-order functions, currying, composition, partial application  
**Module System**: Namespaces, exports/imports, qualified access  
**Macro System**: Code transformation with `defmacro`, `quote`, and templating  
**String Processing**: 20+ functions including regex support  
**Development Tools**: Interactive REPL, environment inspection, built-in help  

## Key Operations

### Arithmetic
`(+ 1 2 3)` `(- 10 3)` `(* 2 3)` `(/ 15 3)` `(mod 10 3)`  
Automatic big number support for arbitrary precision.

### Lists  
`(list 1 2 3)` `(first lst)` `(rest lst)` `(append lst1 lst2)`  
`(map fn lst)` `(filter pred lst)` `(reduce fn init lst)`

### Hash Maps
`(hash-map :key "value")` `(hash-map-get hm :key)` `(hash-map-put hm :key val)`

### Functions (Modern Square Bracket Syntax)
`(defn name [params] body)` `(lambda [x] (* x x))`  
`(apply fn args)` `(compose f g)` `(partial fn arg1)`  
Square brackets make parameters visually distinct and reduce confusion.

### Control Flow
`(if condition then else)` `(cond ...)` `(when pred body)`  
`(defmacro name [params] template)`

See `examples/` directory for comprehensive demonstrations.
- **Big Numbers**: Arbitrary precision integers (e.g., `1000000000000000000000000000000`)
- **Strings**: `"hello world"` with full Unicode support
- **Booleans**: `#t`, `#f` (true/false)
- **Nil**: `nil` (represents empty/null values, falsy in conditionals)
- **Keywords**: `:name`, `:status`, `:id` (self-evaluating symbols, perfect for hash map keys)
- **Lists**: `(1 2 3)`, `("a" "b" "c")`, `()` (immutable linked lists)
- **Hash Maps**: `{:name "Alice" :age 30}` (immutable key-value associative arrays)
- **Symbols**: `+`, `-`, `x`, `my-var` (identifiers and operators)
- **Functions**: `#<function([param1 param2])>` (first-class callable objects)
- **Macros**: `#<macro([param1 param2])>` (code transformation functions)
- **Modules**: `#<module:name>` (namespace containers with exports)
- **Environments**: Runtime scoping and variable binding contexts

## Development Features (2025 Tooling)

- **Interactive REPL**: Full-featured development environment with:
  - Command history and line editing
  - Integrated help system (`(help)`, `(env)`, `(modules)`)
  - Syntax error recovery and helpful error messages
  - Multi-line expression support
- **Environment Inspection**: Complete introspection capabilities:
  - `(env)` - View all variables and functions in current scope
  - `(modules)` - List all loaded modules with their exports
  - `(help)` - Discover all available built-in functions
  - `(help function-name)` - Get detailed help for specific functions
- **Module System**: Production-ready namespace management:
  - Module definition with explicit exports
  - Import system with qualified and unqualified access
  - File loading with dependency resolution
  - Circular dependency detection
- **File Execution**: Robust script execution capabilities:
  - Multi-expression file processing
  - Command-line argument parsing
  - Error handling with file context
  - Module loading from disk
- **Advanced Error Handling**: Comprehensive debugging support:
  - Stack traces with function call context
  - Source location information
  - Clear error messages with suggestions
  - Built-in `error` function for controlled termination
- **Performance Monitoring**: Built-in profiling and optimization tools
- **Development Workflow**: Integrated with modern development practices

## Module System

The interpreter provides a comprehensive module system for organizing code into reusable namespaces with explicit exports and imports.

### Module Definition
```lisp
(module math-utils
  (export square cube add-squares)
  
  (defn square [x] (* x x))
  (defn cube [x] (* x x x))
  (defn add-squares [x y] (+ (square x) (square y)))
  
  ; Private helper function (not exported)
  (defn helper [x] (+ x 1)))
```

### Loading and Importing

#### Individual Operations
```lisp
; Load a file containing module definitions
(load "library/core.lisp")

; Import a module's exports into current environment
(import core)

; Use imported functions directly
(factorial 10)

; Or use qualified access without importing
(core.factorial 10)
```

#### Unified Require Function
The `require` function combines loading and importing in a single operation:

```lisp
; Basic require - load file and import all exports
(require "library/core.lisp")
(factorial 10)  ; Available immediately

; Equivalent to:
; (load "library/core.lisp")
; (import core)
```

**Benefits:**
- **Simplicity**: One command instead of two separate load/import calls
- **Efficiency**: Automatically detects and imports modules from loaded files
- **File Caching**: Prevents re-loading the same file multiple times
- **Error Handling**: Comprehensive error messages for missing files or modules

**Usage Example:**
```lisp
; Traditional approach
(load "library/functional.lisp")
(import functional)
(map (lambda (x) (* x 2)) (list 1 2 3))

; Simplified with require  
(require "library/functional.lisp")
(map (lambda (x) (* x 2)) (list 1 2 3))
```

**Current Implementation:**
- Supports basic `(require "filename")` syntax
- Automatically loads file and imports all module exports
- Prevents duplicate file loading through caching
- Works with all existing library modules (core, functional, strings, macros)

### Access Patterns

#### Qualified Access
```lisp
; Access without importing
(math-utils.square 5)      ; => 25
(core.factorial 10)        ; => 3628800
```

#### Direct Access After Import
```lisp
(import math-utils)
(square 5)                 ; => 25
(add-squares 3 4)         ; => 25
```

#### Module Introspection
```lisp
(modules)                  ; List all loaded modules
(env)                      ; View current environment bindings
```

### File Organization
- **Library Structure**: Organize related functions into modules within files
- **Export Control**: Only exported functions are accessible outside the module
- **Private Functions**: Non-exported functions remain module-internal
- **Dependency Management**: Files can load other files, creating dependency chains

## Big Number Support

The interpreter provides comprehensive support for arbitrary precision arithmetic:

### Automatic Precision Detection
- Large integers (≥ 10^15) are automatically handled as big numbers during parsing
- Arithmetic operations detect potential overflow and promote results to big numbers
- Seamless mixing of regular numbers and big numbers in expressions

### Readable Formatting
- Big numbers display in standard decimal format without separators
- Example: `1000000000000000000000` displays as `1000000000000000000000`

### Operations Support
- All arithmetic operations: `+`, `-`, `*`, `/`
- All comparison operations: `=`, `<`, `>`, `<=`, `>=`
- Compatible with tail-recursive algorithms for computing large factorials and Fibonacci numbers

### Examples
```lisp
; Large multiplication automatically uses big numbers
(* 1000000000000000 1000000000000000)
=> 1000000000000000000000000000000

; Factorial of large numbers
(defn factorial [n acc]
  (if (= n 0) acc (factorial (- n 1) (* n acc))))

(factorial 50 1)
=> 30414093201713378043612608166064768844377641568960512000000000000

; Comparisons work seamlessly
(> 1000000000000000000000 999999999999999999999)
=> #t
```

## Macro System

The interpreter includes a powerful macro system that enables code transformation at evaluation time, allowing developers to extend the language with custom syntax and control structures.

### Macro Definition
- `(defmacro name [params] body)` - Define a macro that transforms code before evaluation
- Macros receive unevaluated arguments and return code to be evaluated
- Macro expansion happens at evaluation time, not parse time

### Quote Special Form  
- `(quote expr)` or `'expr` - Return expression without evaluating it
- Essential for macro programming to manipulate code as data
- Both long form `(quote ...)` and shorthand `'...` syntax supported

### Examples
```lisp
; Define a 'when' control structure
(defmacro when [condition body]
  (list 'if condition body 'nil))

; Use the when macro
(when (> x 5) (print "x is greater than 5"))
; Expands to: (if (> x 5) (print "x is greater than 5") nil)

; Define an 'unless' macro
(defmacro unless [condition then else]
  (list 'if condition else then))

(unless (< x 5) "not less" "is less")
; Expands to: (if (< x 5) "is less" "not less")

; Quote prevents evaluation
(quote (+ 1 2))    ; => (+ 1 2)
'(+ 1 2)           ; => (+ 1 2) 
(+ 1 2)            ; => 3

; Complex macro - let-like binding
(defmacro let1 [var value body]
  (list (list 'lambda (list var) body) value))

(let1 x 10 (+ x 5))  ; => 15
; Expands to: ((lambda [x] (+ x 5)) 10)
```

### Macro Benefits
- **Language Extension**: Create new control structures and syntax
- **Code Generation**: Automatically generate repetitive code patterns  
- **DSL Creation**: Build domain-specific languages within Lisp
- **Performance**: Code transformation happens at evaluation time, not runtime
