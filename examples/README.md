# Examples

Comprehensive examples demonstrating all interpreter features. Each file focuses on specific capabilities with practical, runnable code.

## Core Examples

### `basic_features.lisp`
**Language fundamentals**: arithmetic, functions, lists, closures, recursion with tail-call optimization.

### `advanced_features.lisp` 
**Advanced capabilities**: tail-call optimization, environment introspection, big numbers, error handling.

### `module_system.lisp`
**Module system**: creating modules, exports/imports, qualified access, visibility control.

## Library Examples

### `core_library.lisp`
**Built-in functions**: mathematical operations, list utilities, higher-order functions.

### `functional_library.lisp`
**Functional programming**: currying, composition, partial application, higher-order utilities, data pipelines.

### `string_library.lisp`
**String processing**: 20+ string functions, regex operations, type conversions.

### `macro_library.lisp`
**Macro utilities**: control flow macros, debugging tools, code transformation helpers.

## Data Types

### `hash_maps.lisp`
**Associative arrays**: creation, manipulation, nested structures, immutable operations.

### `keywords.lisp`
**Keyword syntax**: self-evaluating symbols, hash map keys, data structure patterns.

### `atoms.lisp`
**Thread-safe state**: Clojure-style atoms, atomic operations, mutable references, concurrency safety.

### `macro_system.lisp`
**Macro programming**: defmacro, quote forms, code generation, DSL creation.

## Output & Formatting

### `print_functions.lisp`
**Print capabilities**: output formatting, multiple data types, debugging utilities.

### `file_functions.lisp`
**File I/O operations**: reading, writing, checking file existence, data processing pipelines, configuration management.

## Quick Start

```bash
# Run individual examples
./lisp examples/basic_features.lisp
./lisp examples/functional_library.lisp
./lisp examples/macro_system.lisp

# Run all examples
for file in examples/*.lisp; do ./lisp "$file"; done

# Interactive exploration
./lisp
lisp> (load "examples/basic_features.lisp")
```

## Key Features Demonstrated

**Language Core**: Variables, functions, lists, arithmetic, comparisons, comments  
**Modern Data**: Keywords (`:name`), hash maps, big numbers, comprehensive strings  
**Advanced**: Higher-order functions, closures, tail-call optimization, error handling  
**Functional**: Currying, composition, partial application, data transformation pipelines  
**Module System**: Definition, exports, imports, qualified access, visibility  
**Development**: Environment inspection, REPL, file loading, debugging tools

Each example is self-contained and demonstrates practical usage patterns.
