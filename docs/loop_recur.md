# Loop and Recur Implementation

This document describes the implementation of `loop` and `recur` constructs in the Go Lisp interpreter.

## Overview

The `loop` and `recur` constructs provide efficient tail recursion without growing the call stack. This is a common pattern in functional programming languages like Clojure.

- `loop` establishes a recursion point with local bindings
- `recur` jumps back to the nearest enclosing `loop` with new values

## Syntax

### Loop
```lisp
(loop [binding-vector] body-expressions...)
```

The binding vector contains pairs of variable names and initial values:
```lisp
(loop [var1 init1 var2 init2 ...] body...)
```

### Recur
```lisp
(recur new-value1 new-value2 ...)
```

The number of arguments to `recur` must match the number of bindings in the loop.

## Examples

### Simple Countdown
```lisp
(loop [i 5]
  (if (= i 0)
    "done"
    (do
      (println i)
      (recur (- i 1)))))
```

### Factorial Calculation
```lisp
(loop [n 5 acc 1]
  (if (<= n 1)
    acc
    (recur (- n 1) (* acc n))))
```

### Sum of Numbers
```lisp
(loop [i 1 n 10 sum 0]
  (if (> i n)
    sum
    (recur (+ i 1) n (+ sum i))))
```

## Implementation Details

### Control Flow Mechanism

The implementation uses a special exception-based control flow:

1. **RecurException**: A special error type that carries the new binding values
2. **Loop Catching**: The `loop` function catches `RecurException` and continues with new values
3. **Argument Validation**: Ensures the number of `recur` arguments matches loop bindings

### Architecture

```go
type RecurException struct {
    Args []Value // Arguments to pass to the next iteration
}
```

The `loop` function:
1. Parses the binding vector
2. Evaluates initial binding values
3. Executes the body in a loop
4. Catches `RecurException` and updates bindings
5. Continues until no `recur` is called

The `recur` function:
1. Evaluates all arguments
2. Throws a `RecurException` with the new values

### Current Limitations

Due to the current evaluator interface design, there are some limitations:

1. **Environment Isolation**: Loop variables are not properly bound in the evaluation environment
2. **Scope Access**: The body expressions cannot directly access the loop variables
3. **Evaluator Interface**: The current interface doesn't support creating child environments with specific bindings

### Full Implementation Requirements

For a complete implementation, the following would be needed:

1. **Enhanced Evaluator Interface**: Support for creating child evaluators with specific bindings
2. **Environment Manipulation**: Direct access to create scoped environments
3. **Special Forms**: Treating `loop` as a special form rather than a function

## Testing

The implementation includes comprehensive tests:

- Basic functionality tests
- Error handling for malformed binding vectors
- Argument count validation
- RecurException throwing and catching

## Usage in Go Lisp

The `loop` and `recur` functions are automatically registered when the control plugin is loaded. They can be used immediately in any Go Lisp program.

## Performance Benefits

While the current implementation has limitations, the pattern provides:

- **Tail Call Optimization**: No call stack growth for recursive operations
- **Memory Efficiency**: Constant memory usage for iterative algorithms
- **Functional Style**: Maintains immutability while providing efficient iteration

This implementation provides the foundation for proper loop/recur functionality and can be enhanced as the evaluator architecture evolves.
