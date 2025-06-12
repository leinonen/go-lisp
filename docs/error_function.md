# Error Function Documentation

## Overview
The `error` function allows you to raise custom errors with specified messages from within Lisp code. This is useful for input validation, error propagation, and creating robust programs.

## Syntax
```lisp
(error message)
```

## Parameters
- `message`: The error message to display. Can be any value that will be converted to a string representation.

## Return Value
This function never returns normally - it always raises an error with the specified message.

## Examples

### String Messages
```lisp
(error "Something went wrong!")
; Raises: Something went wrong!

(error "Invalid input: expected positive number")
; Raises: Invalid input: expected positive number
```

### Non-String Messages
```lisp
(error 404)
; Raises: 404

(error #t)
; Raises: #t

(error (list 1 2 3))
; Raises: (1 2 3)
```

### Conditional Error Handling
```lisp
(defun divide-safe [a b]
  (if (= b 0)
      (error "Division by zero!")
      (/ a b)))

(divide-safe 10 2)  ; Returns 5
(divide-safe 10 0)  ; Raises: Division by zero!
```

### Input Validation
```lisp
(defun factorial [n]
  (if (< n 0)
      (error "Factorial not defined for negative numbers")
      (if (= n 0)
          1
          (* n (factorial (- n 1))))))

(factorial 5)   ; Returns 120
(factorial -1)  ; Raises: Factorial not defined for negative numbers
```

### Error Propagation
```lisp
(defun validate-age [age]
  (if (< age 0)
      (error "Age cannot be negative")
      (if (> age 150)
          (error "Age seems unrealistic")
          age)))

(defun create-person [name age]
  (list name (validate-age age)))

(create-person "Alice" 30)   ; Returns ("Alice" 30)
(create-person "Bob" -5)     ; Raises: Age cannot be negative
(create-person "Eve" 200)    ; Raises: Age seems unrealistic
```

## Error Conditions
- **Wrong Number of Arguments**: The function requires exactly 1 argument
- **Expression Evaluation Error**: If the message expression cannot be evaluated, that error is raised instead

## Implementation Notes
- The message argument is evaluated before creating the error
- Any Lisp value can be used as an error message
- Non-string values are converted to their string representation
- This function integrates with the interpreter's error handling system

## Common Patterns

### Guard Clauses
```lisp
(defun sqrt-positive [x]
  (if (< x 0) (error "Square root of negative number"))
  (sqrt x))
```

### Assertion-Style Checks
```lisp
(defun process-list [lst]
  (if (not (list? lst)) (error "Expected a list"))
  (if (empty? lst) (error "List cannot be empty"))
  ; Process the list...
  )
```

### Custom Error Types (Using Conventions)
```lisp
(defun file-operation [filename]
  (if (not (file-exists? filename))
      (error (list "FileNotFound" filename))
      ; Process file...
      ))
```

## See Also
- [Control Flow](control_flow.md)
- [Function Definition](functions.md)
- [Input Validation Patterns](patterns.md)
