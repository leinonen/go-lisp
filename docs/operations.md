# Supported Operations

This document provides a comprehensive reference for all operations supported by the Lisp interpreter.

## Arithmetic

- `(+ 1 2 3)` - Addition with multiple operands
- `(- 10 3)` - Subtraction
- `(* 2 3 4)` - Multiplication with multiple operands
- `(/ 15 3)` - Division

### Big Number Support

The interpreter automatically handles arbitrary precision arithmetic for large integers:

- **Automatic Detection**: Numbers ≥ 10^15 are parsed as big numbers
- **Overflow Protection**: Arithmetic operations automatically promote to big numbers when needed
- **Seamless Integration**: Big numbers work with all arithmetic and comparison operations
- **Readable Output**: Large numbers display in standard decimal format

Examples:
```lisp
; Large multiplication
(* 1000000000000000 1000000000000000)
=> 1000000000000000000000000000000

; Mixed operations
(+ 1000000000000000000 123)
=> 1000000000000000123

; Factorial of large numbers
(defun factorial (n acc)
  (if (= n 0) acc (factorial (- n 1) (* n acc))))

(factorial 50 1)
=> 30414093201713378043612608166064768844377641568960512000000000000
```

## Comparison

- `(= 5 5)` - Equality
- `(< 3 5)` - Less than
- `(> 7 3)` - Greater than
- `(<= 3 5)` - Less than or equal
- `(>= 7 3)` - Greater than or equal

**Note**: All comparison operations work seamlessly with big numbers and mixed number types.

Examples:
```lisp
; Big number comparisons
(= 1000000000000000000 1000000000000000000)
=> #t

(> 1000000000000000001 1000000000000000000)
=> #t

; Mixed type comparisons
(< 999999999999999 1000000000000000000)
=> #t
```

## Logical

- `(and #t #f)` - Logical AND (returns #t if all arguments are #t)
- `(or #t #f)` - Logical OR (returns #t if any argument is #t)
- `(not #t)` - Logical NOT (returns the opposite boolean value)

## Conditional

- `(if condition then-expr else-expr)` - If expression

## Variables

- `(define name value)` - Define a variable with a name and value

## Functions

- `(lambda (params) body)` - Create an anonymous function
- `(defun name (params) body)` - Define a named function (combines define and lambda)
- `(funcname args...)` - Call a user-defined function

## Tail Call Optimization

The interpreter automatically optimizes tail-recursive functions to prevent stack overflow:

- **Tail Calls**: Function calls in tail position (last expression) are optimized
- **Stack Safety**: Large recursive computations won't cause stack overflow
- **Automatic**: No special syntax required - optimization is transparent
- **Preserves Semantics**: Non-tail recursive functions work normally

Examples of tail-recursive functions:
```lisp
; Tail-recursive factorial
(defun fact-tail (n acc)
  (if (= n 0) acc (fact-tail (- n 1) (* n acc))))

; Tail-recursive countdown
(defun countdown (n)
  (if (= n 0) 0 (countdown (- n 1))))
```

## Lists

- `(list)` - Create an empty list
- `(list 1 2 3)` - Create a list with elements
- `(first lst)` - Get the first element of a list
- `(rest lst)` - Get all elements except the first
- `(cons elem lst)` - Prepend an element to a list
- `(length lst)` - Get the number of elements in a list
- `(empty? lst)` - Check if a list is empty
- `(append lst1 lst2)` - Combine two lists into one
- `(reverse lst)` - Reverse the order of elements in a list
- `(nth index lst)` - Get the element at a specific index (0-based)

## Higher-Order Functions

- `(map func lst)` - Apply a function to each element of a list
- `(filter predicate lst)` - Keep only elements that satisfy a predicate
- `(reduce func init lst)` - Reduce a list to a single value using a function

## Comments

- `;` - Line comments (from semicolon to end of line are ignored)
- Comments can appear anywhere in the code
- Useful for documenting code and adding explanations

## Module System

- `(module name (export sym1 sym2...) body...)` - Define a module with exported symbols
- `(import module-name)` - Import all exported symbols from a module into current scope
- `(load "filename.lisp")` - Load and execute a Lisp file
- `module.symbol` - Qualified access to module symbols without importing

## Environment Inspection

- `(env)` - Show all variables and functions in the current environment
- `(modules)` - Show all loaded modules and their exported symbols
- `(builtins)` - Show all available built-in functions and special forms
- `(builtins func-name)` - Get detailed help for a specific built-in function

## Output Functions

- `(print value1 value2 ...)` - Output values to stdout without newline
- `(println value1 value2 ...)` - Output values to stdout with newline

Examples:
```lisp
; Print without newline
(print "Hello" " " "World")  ; Output: Hello World

; Print with newline
(println "Hello World")      ; Output: Hello World\n

; Print multiple values
(println "Result:" (+ 1 2))  ; Output: Result: 3\n

; Print different types
(println "String:" "hello" "Number:" 42 "Boolean:" #t)
; Output: String: hello Number: 42 Boolean: #t\n
```
