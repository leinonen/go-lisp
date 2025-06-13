# Mathematical Functions Reference

Complete reference for all built-in mathematical functions in Go Lisp.

## Overview

Go Lisp provides a comprehensive set of mathematical functions covering:
- Basic mathematical operations (sqrt, pow, abs)
- Trigonometric functions (sin, cos, tan and their inverses)
- Hyperbolic functions (sinh, cosh, tanh)
- Logarithmic and exponential functions
- Rounding and truncation functions
- Statistical functions (min, max)
- Utility functions (sign, mod)
- Mathematical constants (pi, e)
- Random number generation

All mathematical functions handle both regular numbers and big numbers automatically.

## Basic Mathematical Functions

### Square Root
```lisp
(sqrt x)           ; Square root of x
```

**Examples:**
```lisp
(sqrt 16)          ; => 4
(sqrt 2)           ; => 1.4142135623730951
(sqrt 0)           ; => 0
```

**Notes:**
- Returns error for negative numbers
- Returns exact integer when possible

### Power
```lisp
(pow base exponent) ; base raised to the power of exponent
```

**Examples:**
```lisp
(pow 2 3)          ; => 8
(pow 5 0)          ; => 1
(pow 10 -2)        ; => 0.01
(pow 2.5 2)        ; => 6.25
```

### Absolute Value
```lisp
(abs x)            ; Absolute value of x
```

**Examples:**
```lisp
(abs -7)           ; => 7
(abs 3)            ; => 3
(abs 0)            ; => 0
```

## Trigonometric Functions

### Basic Trigonometric Functions
```lisp
(sin x)            ; Sine of x (x in radians)
(cos x)            ; Cosine of x (x in radians)
(tan x)            ; Tangent of x (x in radians)
```

**Examples:**
```lisp
(sin 0)            ; => 0
(cos 0)            ; => 1
(tan 0)            ; => 0
(sin (/ (pi) 2))   ; => 1 (sin of π/2)
(cos (pi))         ; => -1 (cos of π)
```

### Inverse Trigonometric Functions
```lisp
(asin x)           ; Arc sine of x (returns radians)
(acos x)           ; Arc cosine of x (returns radians)
(atan x)           ; Arc tangent of x (returns radians)
(atan2 y x)        ; Two-argument arc tangent of y/x
```

**Examples:**
```lisp
(asin 0.5)         ; => 0.5235987755982989 (π/6)
(acos 0.5)         ; => 1.0471975511965979 (π/3)
(atan 1)           ; => 0.7853981633974483 (π/4)
(atan2 1 1)        ; => 0.7853981633974483 (π/4)
```

**Notes:**
- `asin` and `acos` require input in range [-1, 1]
- `atan2` handles quadrant correctly and avoids division by zero

## Hyperbolic Functions

```lisp
(sinh x)           ; Hyperbolic sine of x
(cosh x)           ; Hyperbolic cosine of x
(tanh x)           ; Hyperbolic tangent of x
```

**Examples:**
```lisp
(sinh 0)           ; => 0
(cosh 0)           ; => 1
(tanh 0)           ; => 0
(sinh 1)           ; => 1.1752011936438014
```

## Angle Conversion

```lisp
(degrees radians)  ; Convert radians to degrees
(radians degrees)  ; Convert degrees to radians
```

**Examples:**
```lisp
(degrees (pi))     ; => 180
(radians 180)      ; => 3.141592653589793
(degrees (/ (pi) 2)) ; => 90
(radians 45)       ; => 0.7853981633974483
```

## Logarithmic and Exponential Functions

### Natural Logarithm and Exponential
```lisp
(log x)            ; Natural logarithm (base e)
(exp x)            ; Exponential function (e^x)
```

**Examples:**
```lisp
(log (e))          ; => 1
(log 1)            ; => 0
(exp 0)            ; => 1
(exp 1)            ; => 2.718281828459045 (e)
```

### Other Logarithm Functions
```lisp
(log10 x)          ; Base-10 logarithm
(log2 x)           ; Base-2 logarithm
```

**Examples:**
```lisp
(log10 100)        ; => 2
(log10 1000)       ; => 3
(log2 8)           ; => 3
(log2 1024)        ; => 10
```

**Notes:**
- All logarithm functions require positive input
- Returns error for zero or negative numbers

## Rounding and Truncation Functions

```lisp
(floor x)          ; Largest integer ≤ x
(ceil x)           ; Smallest integer ≥ x
(round x)          ; Round to nearest integer
(trunc x)          ; Truncate towards zero
```

**Examples:**
```lisp
(floor 3.7)        ; => 3
(floor -2.3)       ; => -3
(ceil 3.2)         ; => 4
(ceil -2.7)        ; => -2
(round 3.4)        ; => 3
(round 3.6)        ; => 4
(trunc 3.7)        ; => 3
(trunc -3.7)       ; => -3
```

**Notes:**
- `floor` always rounds down (towards negative infinity)
- `ceil` always rounds up (towards positive infinity)
- `round` rounds to nearest integer (0.5 rounds up)
- `trunc` removes fractional part (rounds towards zero)

## Statistical Functions

```lisp
(min a b)          ; Minimum of two numbers
(max a b)          ; Maximum of two numbers
```

**Examples:**
```lisp
(min 10 5)         ; => 5
(max 10 5)         ; => 10
(min -3 2)         ; => -3
(max -3 2)         ; => 2
```

## Utility Functions

### Sign Function
```lisp
(sign x)           ; Sign of x (-1, 0, or 1)
```

**Examples:**
```lisp
(sign 5)           ; => 1
(sign -3)          ; => -1
(sign 0)           ; => 0
```

### Mathematical Modulo
```lisp
(mod x y)          ; Mathematical modulo (proper handling of negatives)
```

**Examples:**
```lisp
(mod 7 3)          ; => 1
(mod -7 3)         ; => 2  (always non-negative when divisor is positive)
(mod 17 5)         ; => 2
```

**Notes:**
- Unlike the `%` operator, `mod` ensures the result has the same sign as the divisor
- For positive divisor, result is always non-negative

## Mathematical Constants

```lisp
(pi)               ; π (pi) ≈ 3.141592653589793
(e)                ; e (Euler's number) ≈ 2.718281828459045
```

**Examples:**
```lisp
(pi)               ; => 3.141592653589793
(e)                ; => 2.718281828459045
(* 2 (pi))         ; => 6.283185307179586 (2π)
(pow (e) 2)        ; => 7.3890560989306504 (e²)
```

## Random Number Generation

```lisp
(random)           ; Random float between 0 and 1
(random n)         ; Random integer between 0 and n-1
(random min max)   ; Random integer between min and max-1
```

**Examples:**
```lisp
(random)           ; => 0.6234567890123456 (example)
(random 10)        ; => 7 (example, 0-9)
(random 5 15)      ; => 12 (example, 5-14)
```

**Notes:**
- Random numbers are generated using a seeded PRNG
- For integer ranges, both bounds must be integers
- Upper bound is exclusive (not included in range)

## Complex Mathematical Expressions

### Distance Formula
```lisp
(defn distance [x1 y1 x2 y2]
  (sqrt (+ (pow (- x2 x1) 2) (pow (- y2 y1) 2))))

(distance 0 0 3 4)  ; => 5
```

### Quadratic Formula
```lisp
(defn quadratic-roots [a b c]
  (let [discriminant (- (pow b 2) (* 4 a c))]
    (if (< discriminant 0)
      "No real roots"
      (list (/ (+ (- b) (sqrt discriminant)) (* 2 a))
            (/ (- (- b) (sqrt discriminant)) (* 2 a))))))

(quadratic-roots 1 -5 6)  ; => (-2 -3)
```

### Compound Interest
```lisp
(defn compound-interest [principal rate times years]
  (* principal (pow (+ 1 (/ rate times)) (* times years))))

(compound-interest 1000 0.05 12 10)  ; => ~1643.62
```

### Area and Circumference of Circle
```lisp
(defn circle-area [radius]
  (* (pi) (pow radius 2)))

(defn circle-circumference [radius]
  (* 2 (pi) radius))

(circle-area 5)           ; => 78.53981633974483
(circle-circumference 5)  ; => 31.41592653589793
```

## Error Handling

Mathematical functions provide clear error messages for invalid inputs:

```lisp
(sqrt -1)          ; Error: "cannot compute square root of negative number"
(log 0)            ; Error: "cannot compute logarithm of non-positive number"
(asin 2)           ; Error: "input must be in range [-1, 1]"
(mod 5 0)          ; Error: "division by zero"
```

## Performance Notes

- All functions handle both regular and big numbers automatically
- Functions use Go's `math` package for optimal performance
- Domain checking is performed to prevent invalid operations
- Results are returned as the most appropriate number type

## See Also

- [Basic Operations](operations.md) - Arithmetic and comparison operations
- [Core Library](core_library.md) - Higher-level mathematical functions
- [Examples](examples.md) - More complex mathematical examples
