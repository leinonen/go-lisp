# Features (June 2025)

This Lisp interpreter provides a **comprehensive, production-ready** implementation with modern features that rival contemporary programming languages while maintaining classic Lisp elegance.

## Core Components

- **Tokenizer/Lexer**: Converts input text into tokens for parsing
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens  
- **Evaluator**: Evaluates the AST to produce results
- **REPL**: Interactive Read-Eval-Print Loop with helpful startup commands

## Language Features (2025 Edition)

- **Comments**: Full line comment support using semicolons with nested comment capability
- **Built-in Help System**: Comprehensive help with `(builtins)`, `(builtins func-name)`, and `(env)`
- **Big Number Arithmetic**: Arbitrary precision integers with automatic overflow detection and seamless integration
- **Keywords**: Self-evaluating symbols (`:name`, `:status`) perfect for structured data and hash map keys
- **First-class Functions**: Functions as values - store, pass, return, and compose functions naturally
- **Closures**: Functions capture and preserve their lexical environment
- **Advanced Recursion**: Full recursion support with tail call optimization for stack safety
- **Higher-order Functions**: Functions that operate on other functions (`map`, `filter`, `reduce`)
- **Pattern Matching**: Sophisticated conditional logic and data destructuring
- **Module System**: Namespace management with imports, exports, and qualified access
- **Error Handling**: Built-in error function with stack traces and debugging context
- **String Processing**: 20+ built-in string functions plus regex support
- **Immutable Data Structures**: Hash maps and lists with functional update semantics
- **Macro System**: Code transformation at evaluation time with `defmacro` and `quote`

## Data Types (Complete Type System)

- **Numbers**: `42`, `-3.14` (IEEE 754 double precision)
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
  - Integrated help system (`(builtins)`, `(env)`, `(modules)`)
  - Syntax error recovery and helpful error messages
  - Multi-line expression support
- **Environment Inspection**: Complete introspection capabilities:
  - `(env)` - View all variables and functions in current scope
  - `(modules)` - List all loaded modules with their exports
  - `(builtins)` - Discover all available built-in functions
  - `(builtins function-name)` - Get detailed help for specific functions
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
(defun factorial (n acc)
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
- `(defmacro name (params) body)` - Define a macro that transforms code before evaluation
- Macros receive unevaluated arguments and return code to be evaluated
- Macro expansion happens at evaluation time, not parse time

### Quote Special Form  
- `(quote expr)` or `'expr` - Return expression without evaluating it
- Essential for macro programming to manipulate code as data
- Both long form `(quote ...)` and shorthand `'...` syntax supported

### Examples
```lisp
; Define a 'when' control structure
(defmacro when (condition body)
  (list 'if condition body 'nil))

; Use the when macro
(when (> x 5) (print "x is greater than 5"))
; Expands to: (if (> x 5) (print "x is greater than 5") nil)

; Define an 'unless' macro
(defmacro unless (condition then else)
  (list 'if condition else then))

(unless (< x 5) "not less" "is less")
; Expands to: (if (< x 5) "is less" "not less")

; Quote prevents evaluation
(quote (+ 1 2))    ; => (+ 1 2)
'(+ 1 2)           ; => (+ 1 2) 
(+ 1 2)            ; => 3

; Complex macro - let-like binding
(defmacro let1 (var value body)
  (list (list 'lambda (list var) body) value))

(let1 x 10 (+ x 5))  ; => 15
; Expands to: ((lambda (x) (+ x 5)) 10)
```

### Macro Benefits
- **Language Extension**: Create new control structures and syntax
- **Code Generation**: Automatically generate repetitive code patterns  
- **DSL Creation**: Build domain-specific languages within Lisp
- **Performance**: Code transformation happens at evaluation time, not runtime
