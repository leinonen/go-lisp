# Go Lisp Examples

This directory contains focused examples demonstrating various features of the Go Lisp interpreter.

## Available Examples

### `loop_recur_examples.lisp`
Demonstrates the `loop` and `recur` constructs for efficient tail recursion:
- Simple countdown
- Factorial calculation
- Sum of numbers

### `list_operations.lisp`
Shows list manipulation functions:
- Creating lists with `list`
- Accessing elements with `first`, `rest`, `last`, `nth`
- Adding elements with `cons`, `append`
- List properties: `length`, `empty?`
- List transformations: `reverse`

### `functional_programming.lisp`
Higher-order functions for functional programming:
- `map` - apply function to each element
- `filter` - select elements matching predicate
- `reduce` - combine elements with function
- Function composition examples

### `hashmap_operations.lisp`
Hash map data structure operations:
- Creating hash maps
- Getting and setting values
- Hash map properties and predicates
- Working with keys and values
- Nested hash map examples

### `string_manipulation.lisp`
Comprehensive string processing:
- Basic operations: concat, length, substring
- Case conversion and trimming
- String searching and replacement
- Regular expression matching
- Number/string conversion

### `arithmetic_math.lisp`
Mathematical operations and functions:
- Basic arithmetic: `+`, `-`, `*`, `/`, `%`
- Comparison operations
- Math functions: `sqrt`, `pow`, `sin`, `cos`, etc.
- Number predicates and properties
- Practical calculation examples

### `control_flow.lisp`
Control structures and function definitions:
- Variable binding with `def` and `let`
- Conditional expressions: `if`, `cond`
- Function definitions: `defn`
- Anonymous functions
- Higher-order function examples

### `comprehensive_example.lisp`
Real-world example combining multiple features:
- Employee data processing pipeline
- Data filtering and transformation
- Statistical calculations
- Report generation
- Complex nested data navigation

## Running Examples

To run these examples, use the Go Lisp REPL or interpreter:

```bash
# Build the interpreter
make build

# Run individual examples
./bin/golisp < examples/list_operations.lisp

# Or start the REPL and load examples
./bin/golisp
> (load "examples/functional_programming.lisp")
```

## Features Demonstrated

- **Data Structures**: Lists, Hash Maps, Strings, Numbers
- **Functional Programming**: Map, Filter, Reduce, Higher-order functions
- **Control Flow**: Conditionals, Function definitions, Recursion
- **String Processing**: Manipulation, Regular expressions, Formatting
- **Mathematical Operations**: Arithmetic, Trigonometry, Statistics
- **Loop Constructs**: Efficient tail recursion with loop/recur

Each example file is self-contained and includes comments explaining the functionality being demonstrated.
