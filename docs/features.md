# Features

This Lisp interpreter provides a comprehensive set of features for functional programming and interactive development.

## Core Components

- **Tokenizer/Lexer**: Converts input text into tokens for parsing
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens  
- **Evaluator**: Evaluates the AST to produce results
- **REPL**: Interactive Read-Eval-Print Loop with helpful startup commands

## Language Features

- **Comments**: Support for line comments using semicolons
- **Built-in Help System**: Discover functions with `(builtins)` and get help with `(builtins func-name)`
- **Big Number Arithmetic**: Arbitrary precision integers with automatic overflow detection and readable formatting
- **First-class Functions**: Functions can be stored in variables, passed as arguments, and returned from other functions
- **Closures**: Functions capture and remember variables from their creation environment
- **Recursion**: Functions can call themselves for recursive algorithms
- **Tail Call Optimization**: Prevents stack overflow in tail-recursive functions by eliminating stack growth
- **Higher-order Functions**: Functions that take other functions as arguments or return functions

## Data Types

- **Numbers**: `42`, `-3.14`
- **Big Numbers**: Large integers with arbitrary precision (e.g., `1000000000000000000`)
- **Strings**: `"hello world"`
- **Booleans**: `#t`, `#f`
- **Lists**: `(1 2 3)`, `("a" "b" "c")`, `()`
- **Symbols**: `+`, `-`, `x`, `my-var`
- **Functions**: `#<function([param1 param2])>`

## Development Features

- **Interactive Environment**: Full REPL with command history and helpful startup messages
- **Environment Inspection**: View current variables, functions, and loaded modules
- **Module System**: Organize code into modules with explicit exports and imports
- **File Loading**: Load and execute Lisp files from disk
- **Error Handling**: Clear error messages for debugging

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
