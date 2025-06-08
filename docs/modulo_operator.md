# Modulo Operator Documentation

## Overview
The modulo operator (`%`) computes the remainder of integer division. It supports both regular numbers and arbitrary-precision big numbers.

## Syntax
```lisp
(% dividend divisor)
```

## Parameters
- `dividend`: The number to be divided (can be a regular number or big number)
- `divisor`: The number to divide by (can be a regular number or big number, must not be zero)

## Return Value
Returns the remainder of the division. The result type depends on the input types:
- If both inputs are regular numbers that fit in float64, returns a regular number
- If either input is a big number or the result requires big number precision, returns a big number

## Examples

### Basic Usage
```lisp
(% 17 5)     ; Returns 2
(% 10 3)     ; Returns 1
(% 15 5)     ; Returns 0 (no remainder)
```

### Negative Numbers
```lisp
(% -17 5)    ; Returns -2
(% 17 -5)    ; Returns 2
(% -17 -5)   ; Returns -2
```

### Big Numbers
```lisp
(% 123456789012345678901234567890 123)  ; Works with very large numbers
(% 999999999999999999 7)                ; Automatic big number conversion
```

## Error Conditions
- **Division by Zero**: `(% 5 0)` raises an error
- **Wrong Number of Arguments**: The operator requires exactly 2 arguments
- **Non-Numeric Arguments**: Both arguments must be numbers

## Implementation Notes
- The modulo operation preserves the sign of the dividend
- Automatic conversion to big numbers when necessary to maintain precision
- Zero-division protection for both regular and big numbers

## See Also
- [Basic Arithmetic Operations](operations.md)
- [Big Number Support](big_numbers.md)
- [Error Handling](error_handling.md)
