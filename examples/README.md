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

### `functional_library_demo.lisp` ⭐ NEW
**Functional programming**: currying, composition, partial application, higher-order utilities, data pipelines.

### `string_library_demo.lisp`
**String processing**: 20+ string functions, regex operations, type conversions.

### `macro_library_demo.lisp`
**Macro utilities**: control flow macros, debugging tools, code transformation helpers.

## Data Types

### `hash_maps.lisp`
**Associative arrays**: creation, manipulation, nested structures, immutable operations.

### `keywords.lisp`
**Keyword syntax**: self-evaluating symbols, hash map keys, data structure patterns.

### `macro_system.lisp`
**Macro programming**: defmacro, quote forms, code generation, DSL creation.

## Output & Formatting

### `print_functions_demo.lisp`
**Print capabilities**: output formatting, multiple data types, debugging utilities.

## Quick Start

```bash
# Run individual examples
./lisp-interpreter examples/basic_features.lisp
./lisp-interpreter examples/functional_library_demo.lisp
./lisp-interpreter examples/macro_system.lisp

# Run all examples
for file in examples/*.lisp; do ./lisp-interpreter "$file"; done

# Interactive exploration
./lisp-interpreter
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
