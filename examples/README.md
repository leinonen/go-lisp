# GoLisp Examples

This directory contains focused examples demonstrating various features of the GoLisp interpreter.

## Available Examples

### `polymorphic_functions.lisp`
Demonstrates polymorphic functions that work across all data types:
- Sequence functions: `first`, `rest`, `last`, `nth`, `second` on lists, vectors, strings
- Collection functions: `count`, `empty?` on all collection types
- Access functions: `get`, `contains?` on hashmaps, vectors, lists, strings
- Transformation functions: `take`, `drop`, `reverse` on all sequences
- Type predicate functions: `seq?`, `coll?`, `sequential?`, `indexed?`

### `polymorphic_demo.lisp`
Interactive demonstration of polymorphic functions working across different data types:
- Side-by-side comparison of the same functions on lists, vectors, strings, and hashmaps
- Practical examples of unified data processing
- Complete overview of type predicates and their behavior

### `loop_recur_examples.lisp`
Demonstrates the `loop` and `recur` constructs for efficient tail recursion:
- Simple countdown
- Factorial calculation
- Sum of numbers

### `list_operations.lisp`
Shows list manipulation functions:
- Creating lists with `list`
- Accessing elements with polymorphic functions (`first`, `rest`, `last`, `nth`)
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

### `vector_operations.lisp`
Demonstrates vector data structure operations:
- Creating vectors with `[...]` literals and `vector` function
- Vector type checking with `vector?`
- Vector operations: `count`, `nth`, `conj`
- Functional operations preserving vector type: `map`, `filter`, `reduce`
- Converting between vectors and lists: `vec`, `seq`
- Working with nested vectors

## Running Examples

To run these examples, use the GoLisp REPL or interpreter:

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

- **Polymorphic Functions**: Type-aware functions that work across lists, vectors, strings, and hashmaps
- **Cross-Type Operations**: Unified interface for all collection types with consistent behavior
- **Data Structures**: Lists, Hash Maps, Strings, Numbers, Vectors
- **Functional Programming**: Map, Filter, Reduce, Higher-order functions
- **Control Flow**: Conditionals, Function definitions, Recursion
- **String Processing**: Manipulation, Regular expressions, Formatting
- **Mathematical Operations**: Arithmetic, Trigonometry, Statistics
- **Loop Constructs**: Efficient tail recursion with loop/recur

Each example file is self-contained and includes comments explaining the functionality being demonstrated.
