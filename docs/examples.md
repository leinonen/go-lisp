# Examples

This document provides comprehensive examples of using the Lisp interpreter, from basic operations to advanced features.

## Basic Function Definition

### Traditional Way (using define + lambda)

```lisp
lisp> (define square (lambda (x) (* x x)))
=> <function>

lisp> (square 5)
=> 25
```

### Convenient Way (using defun)

```lisp
lisp> (defun square (x) (* x x))
=> <function>

lisp> (square 5)
=> 25
```

### Multi-parameter Functions

```lisp
lisp> (defun add (x y) (+ x y))
=> <function>

lisp> (add 3 4)
=> 7
```

## Recursive Functions

```lisp
lisp> (defun factorial (n) 
        (if (= n 0) 
          1 
          (* n (factorial (- n 1)))))
=> <function>

lisp> (factorial 5)
=> 120
```

## Basic Operations

```lisp
lisp> 42
=> 42

lisp> (+ 1 2 3)
=> 6

lisp> (define x 10)
=> 10

lisp> (+ x 5)
=> 15

lisp> (* (+ 2 3) 4)
=> 20

lisp> (if (< 3 5) 100 0)
=> 100

lisp> "hello world"
=> hello world

lisp> (= 5 5)
=> #t
```

## Working with Variables

```lisp
lisp> (define x 10)
=> 10

lisp> x
=> 10

lisp> (define y (* x 3))
=> 30

lisp> (+ x y)
=> 40
```

## Lambda Functions and Closures

```lisp
lisp> (lambda (x) (+ x 1))
=> #<function([x])>

lisp> (define add1 (lambda (x) (+ x 1)))
=> #<function([x])>

lisp> (add1 5)
=> 6

lisp> (define make-adder (lambda (n) (lambda (x) (+ x n))))
=> #<function([n])>

lisp> (define add10 (make-adder 10))
=> #<function([x])>

lisp> (add10 7)
=> 17
```

## Comments in Code

```lisp
; This is a comment - it will be ignored
lisp> (+ 1 2 3) ; Comments can appear at the end of lines
=> 6

; Define a function with comments
lisp> (defun factorial (n) ; Calculate factorial recursively
        (if (= n 0)        ; Base case: 0! = 1
          1 
          (* n (factorial (- n 1))))) ; Recursive case
=> <function>

lisp> (factorial 5) ; Test the function
=> 120
```

## List Operations

### Creating Lists

```lisp
lisp> (list 1 2 3)
=> (1 2 3)

lisp> (list)
=> ()

lisp> (list "hello" 42 #t)
=> (hello 42 #t)

lisp> (define my-list (list 10 20 30))
=> (10 20 30)
```

### Basic List Functions

```lisp
lisp> (first my-list)
=> 10

lisp> (rest my-list)
=> (20 30)

lisp> (length my-list)
=> 3

lisp> (empty? my-list)
=> #f

lisp> (empty? (list))
=> #t

lisp> (cons 5 my-list)
=> (5 10 20 30)
```

### List Manipulation

```lisp
lisp> (append (list 1 2) (list 3 4 5))
=> (1 2 3 4 5)

lisp> (reverse (list 1 2 3 4))
=> (4 3 2 1)

lisp> (nth 0 my-list)
=> 10

lisp> (nth 2 my-list)
=> 30
```

### Complex List Operations

```lisp
lisp> (define sum-list (lambda (lst) (if (empty? lst) 0 (+ (first lst) (sum-list (rest lst))))))
=> #<function([lst])>

lisp> (sum-list (list 1 2 3 4))
=> 10

; Append and reverse operations
lisp> (append (list 1 2 3) (list 4 5 6))
=> (1 2 3 4 5 6)

lisp> (reverse (append (list 1 2) (list 3 4)))
=> (4 3 2 1)

lisp> (nth 1 (reverse (list 10 20 30 40)))
=> 30

lisp> (map (lambda (i) (nth i (list "a" "b" "c" "d"))) (list 0 2 1 3))
=> (a c b d)
```

## Higher-Order Function Examples

### Map

```lisp
; Apply a function to each element
lisp> (map (lambda (x) (* x x)) (list 1 2 3 4 5))
=> (1 4 9 16 25)
```

### Filter

```lisp
; Keep only elements that satisfy a predicate
lisp> (filter (lambda (x) (> x 0)) (list -1 2 -3 4 5))
=> (2 4 5)
```

### Reduce

```lisp
; Combine all elements using a function
lisp> (reduce (lambda (acc x) (+ acc x)) 0 (list 1 2 3 4 5))
=> 15
```

### Complex Combinations

```lisp
; Squares of positive numbers
lisp> (map (lambda (x) (* x x)) (filter (lambda (x) (> x 0)) (list -2 -1 0 1 2 3)))
=> (1 4 9)

; Sum of squares using map and reduce
lisp> (reduce (lambda (acc x) (+ acc x)) 0 (map (lambda (x) (* x x)) (list 1 2 3)))
=> 14

; Count elements using reduce
lisp> (reduce (lambda (acc x) (+ acc 1)) 0 (list "a" "b" "c" "d"))
=> 4
```

## Module System Examples

### Defining a Module

```lisp
; Define a math module with exported functions
lisp> (module math
        (export square cube factorial)
        (defun square (x) (* x x))
        (defun cube (x) (* x x x))
        (defun factorial (n)
          (if (= n 0)
            1
            (* n (factorial (- n 1)))))
        (defun private-helper (x) (+ x 1))) ; not exported
=> #<module:math>
```

### Importing and Using Modules

```lisp
; Import all exported functions from the math module
lisp> (import math)
=> #<module:math>

; Use imported functions directly
lisp> (square 5)
=> 25

lisp> (factorial 4)
=> 24

; Try to use private function (should fail)
lisp> (private-helper 5)
=> Error: undefined symbol: private-helper
```

### Qualified Access

```lisp
; Use qualified access without importing
lisp> (module string-utils
        (export concat reverse-string)
        (defun concat (a b) a) ; simplified
        (defun reverse-string (s) s)) ; simplified
=> #<module:string-utils>

; Access functions using qualified names
lisp> (string-utils.concat "Hello" "World")
=> Hello

lisp> (math.cube 3)
=> 27
```

### Loading Files

```lisp
; Load a file containing module definitions
lisp> (load "examples/math_module.lisp")
=> #<module:math>
```

### Complex Module Example

```lisp
; Create modules with complex functionality
lisp> (module data-processing
        (export process-list filter-numbers)
        (defun process-list (lst)
          (map square (filter positive? lst)))
        (defun filter-numbers (lst)
          (filter (lambda (x) (> x 10)) lst))
        (defun positive? (x) (> x 0))) ; helper function
=> #<module:data-processing>

lisp> (import data-processing)
=> #<module:data-processing>

lisp> (process-list (list -2 3 4 5))
=> (9 16 25)
```

## Environment Inspection Examples

### Basic Environment Inspection

```lisp
; Check what's in the current environment
lisp> (env)
=> ()

; Define some variables and functions
lisp> (define x 42)
=> 42

lisp> (define greeting "Hello, World!")
=> Hello, World!

lisp> (defun square (n) (* n n))
=> #<function([n])>

; Now check the environment again
lisp> (env)
=> ((x 42) (greeting Hello, World!) (square #<function([n])>))
```

### Module Inspection

```lisp
; Check loaded modules (initially empty)
lisp> (modules)
=> ()

; Create a module
lisp> (module math-utils
        (export add multiply)
        (defun add (a b) (+ a b))
        (defun multiply (a b) (* a b)))
=> #<module:math-utils>

; Check modules now
lisp> (modules)
=> ((math-utils (add multiply)))

; Create another module
lisp> (module string-utils
        (export concat length)
        (defun concat (a b) a)
        (defun length (s) 0))
=> #<module:string-utils>

; Check all modules
lisp> (modules)
=> ((math-utils (add multiply)) (string-utils (concat length)))

; Import a module and check environment
lisp> (import math-utils)
=> #<module:math-utils>

lisp> (env)
=> ((x 42) (greeting Hello, World!) (square #<function([n])>) (add #<function([a b])>) (multiply #<function([a b])>))
```

### Built-in Function Discovery

```lisp
; Discover available built-in functions
lisp> (builtins)
=> (+ - * / = < > if define lambda defun list first rest cons length empty? map filter reduce append reverse nth env modules builtins)

; Get help for specific functions
lisp> (builtins reduce)
=> (reduce func init lst)
Reduce a list to a single value using a function.
Example: (reduce (lambda (acc x) (+ acc x)) 0 (list 1 2 3)) => 6

lisp> (builtins map)
=> (map func lst)
Apply a function to each element of a list.
Example: (map (lambda (x) (* x x)) (list 1 2 3)) => (1 4 9)

lisp> (builtins +)
=> (+ num1 num2 ...)
Addition with multiple operands.
Example: (+ 1 2 3) => 6
```
