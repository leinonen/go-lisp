# Core Library Documentation

## Overview
The core library (`core.lisp`) provides a comprehensive collection of mathematical functions, list utilities, and higher-order functions that extend the basic capabilities of the Lisp interpreter.

## Installation and Usage

### Loading the Library
```lisp
(load "library/core.lisp")
```

### Importing Functions
```lisp
; Import all exported functions
(import core)

; Use functions directly
(factorial 10)
(fibonacci 15)

; Or use qualified access without import
(core.factorial 10)
(core.fibonacci 15)
```

## Mathematical Functions

### factorial(n)
Computes the factorial of a non-negative integer using tail recursion.

```lisp
(factorial 0)   ; Returns 1
(factorial 5)   ; Returns 120
(factorial 10)  ; Returns 3628800
(factorial -1)  ; Error: "Factorial not defined for negative numbers"
```

**Implementation**: Uses a private tail-recursive helper `fact-tail(n, acc)` for efficiency.

### fibonacci(n)
Computes the nth Fibonacci number using tail recursion.

```lisp
(fibonacci 0)   ; Returns 0
(fibonacci 1)   ; Returns 1
(fibonacci 10)  ; Returns 55
(fibonacci 20)  ; Returns 6765
(fibonacci -1)  ; Error: "Fibonacci not defined for negative numbers"
```

**Implementation**: Uses a private tail-recursive helper `fib-tail(n, a, b)`.

### gcd(a, b)
Computes the Greatest Common Divisor using the Euclidean algorithm.

```lisp
(gcd 48 18)     ; Returns 6
(gcd 17 13)     ; Returns 1 (coprime)
(gcd 15 0)      ; Returns 15
```

### lcm(a, b)
Computes the Least Common Multiple.

```lisp
(lcm 12 8)      ; Returns 24
(lcm 15 20)     ; Returns 60
(lcm 7 13)      ; Returns 91
```

### abs(x)
Returns the absolute value of a number.

```lisp
(abs 42)        ; Returns 42
(abs -17)       ; Returns 17
(abs 0)         ; Returns 0
```

### min(a, b) and max(a, b)
Return the minimum or maximum of two numbers.

```lisp
(min 10 5)      ; Returns 5
(max 10 5)      ; Returns 10
(min 7 7)       ; Returns 7
```

## List Utility Functions

### length-sq(lst)
Returns the square of the length of a list (useful for complexity analysis).

```lisp
(length-sq (list 1 2 3))        ; Returns 9
(length-sq (list))              ; Returns 0
(length-sq (list 1 2 3 4 5))    ; Returns 25
```

### all(predicate, lst)
Tests if all elements in a list satisfy a predicate.

```lisp
(all (fn [x] (> x 0)) (list 1 2 3))     ; Returns #t
(all (fn [x] (> x 0)) (list 1 -2 3))    ; Returns #f
(all (fn [x] (> x 0)) (list))           ; Returns #t (vacuous truth)
```

### any(predicate, lst)
Tests if any element in a list satisfies a predicate.

```lisp
(any (fn [x] (> x 0)) (list -1 2 -3))   ; Returns #t
(any (fn [x] (> x 0)) (list -1 -2 -3))  ; Returns #f
(any (fn [x] (> x 0)) (list))           ; Returns #f
```

### take(n, lst)
Returns the first n elements of a list.

```lisp
(take 3 (list 1 2 3 4 5))       ; Returns (1 2 3)
(take 0 (list 1 2 3))           ; Returns ()
(take 10 (list 1 2 3))          ; Returns (1 2 3)
```

### drop(n, lst)
Returns the list with the first n elements removed.

```lisp
(drop 2 (list 1 2 3 4 5))       ; Returns (3 4 5)
(drop 0 (list 1 2 3))           ; Returns (1 2 3)
(drop 10 (list 1 2 3))          ; Returns ()
```

## Higher-Order Functions

### compose(f, g)
Returns a function that is the composition of f and g.

```lisp
(def square (fn [x] (* x x)))
(def increment (fn [x] (+ x 1)))
(def square-then-increment (compose increment square))

(square-then-increment 5)       ; Returns 26 (5² + 1)

; Or use directly
((compose square increment) 5)  ; Returns 36 (square(increment(5)))
```

### apply-n(f, n, x)
Applies function f to x exactly n times.

```lisp
(def increment (fn [x] (+ x 1)))

(apply-n increment 0 5)         ; Returns 5
(apply-n increment 3 5)         ; Returns 8
(apply-n square 2 2)            ; Returns 16 (square(square(2)))
```

## Module Structure

### Exported Functions
All public functions are explicitly exported and available after import:
- `factorial`, `fibonacci`, `gcd`, `lcm`, `abs`, `min`, `max`
- `length-sq`, `all`, `any`, `take`, `drop`
- `compose`, `apply-n`

### Private Functions
Helper functions are not exported and remain internal:
- `fact-tail` (factorial helper)
- `fib-tail` (fibonacci helper)

### Access Patterns
```lisp
; After loading but before import - qualified access only
(core.factorial 5)

; After import - direct access
(import core)
(factorial 5)

; Private functions are never accessible
(fact-tail 5 1)        ; Error: undefined symbol
(core.fact-tail 5 1)   ; Error: undefined symbol
```

## Examples

### Mathematical Computations
```lisp
(load "examples/core.lisp")
(import core)

; Compute factorials
(factorial 10)          ; 3628800

; Find GCD and LCM
(gcd 48 18)            ; 6
(lcm 48 18)            ; 144

; Fibonacci sequence
(map fibonacci (list 0 1 2 3 4 5))  ; (0 1 1 2 3 5)
```

### List Processing
```lisp
(def numbers (list 1 2 3 4 5 6 7 8 9 10))

; Check if all are positive
(all (fn [x] (> x 0)) numbers)     ; #t

; Check if any are even
(any (fn [x] (= (% x 2) 0)) numbers) ; #t

; Take first 5, drop first 3
(take 5 numbers)        ; (1 2 3 4 5)
(drop 3 numbers)        ; (4 5 6 7 8 9 10)
```

### Function Composition
```lisp
(def double (fn [x] (* x 2)))
(def add-one (fn [x] (+ x 1)))

; Create composed function
(def double-then-add-one (compose add-one double))
(double-then-add-one 5)  ; Returns 11

; Apply function multiple times
(apply-n double 3 2)     ; Returns 16 (2 → 4 → 8 → 16)
```

## Performance Notes
- All mathematical functions use tail recursion for stack efficiency
- Big number support is automatic when needed
- List functions are implemented efficiently with proper tail recursion
- The module system provides optimal access patterns

## See Also
- [Module System](modules.md)
- [Mathematical Operations](operations.md)
- [List Operations](lists.md)
- [Higher-Order Functions](higher_order.md)
- [Examples](examples.md)
