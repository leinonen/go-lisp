# Examples (June 2025)

This directory contains **comprehensive, production-ready examples** that demonstrate all features of the Lisp interpreter. These examples serve as both tutorials and practical demonstrations of modern Lisp programming.

## Complete Feature Coverage

Our examples demonstrate **100%** of the interpreter's capabilities across 11 focused demonstration files.

## Example Files

### `basic_features.lisp`
**Core language fundamentals**
- Basic arithmetic and variable definitions
- Function definition with `defun`
- List operations (creation, access, manipulation)
- Higher-order functions (`map`, `filter`, `reduce`)
- Closures and function composition
- Recursion with automatic tail call optimization
- Conditional logic and big number support
- Environment inspection

**Run with:**
```bash
./lisp-interpreter examples/basic_features.lisp
```

### `module_system.lisp`
**Module system capabilities**
- Creating modules with `module` keyword
- Exporting specific functions
- Qualified access to module functions
- Importing modules for direct access
- Private vs. public function visibility
- Multiple module management

**Run with:**
```bash
./lisp-interpreter examples/module_system.lisp
```

### `core_library.lisp`
**Built-in core library functions**
- Loading the standard core library
- Mathematical functions (factorial, fibonacci, gcd, lcm, etc.)
- List utility functions (take, drop, length-sq)
- Predicate functions (all, any)
- Higher-order utilities (compose, apply-n)
- Demonstrates both imported and qualified access patterns

**Run with:**
```bash
./lisp-interpreter examples/core_library.lisp
```

### `advanced_features.lisp`
**Advanced interpreter capabilities**
- Tail call optimization preventing stack overflow
- Environment introspection and reflection
- Function availability checking
- Complex data processing pipelines
- Big number arithmetic with modulo operations
- Error handling capabilities
- File loading system

**Run with:**
```bash
./lisp-interpreter examples/advanced_features.lisp
```

### `keywords.lisp` ⭐ NEW
**Keyword data type demonstration**
- Keyword creation and self-evaluation
- Keywords as hash map keys
- Keyword-based data structures
- Comparison with string keys
- Best practices for keyword usage

**Run with:**
```bash
./lisp-interpreter examples/keywords.lisp
```

### `hash_maps.lisp` ⭐ NEW  
**Hash map operations showcase**
- Hash map creation with multiple data types
- Getting, setting, and updating values
- Key and value inspection
- Immutable operations
- Complex nested data structures
- Performance characteristics

**Run with:**
```bash
./lisp-interpreter examples/hash_maps.lisp
```

### `string_library_demo.lisp` ⭐ ENHANCED
**Comprehensive string processing**
- 20+ built-in string functions
- Regular expression operations
- String search and manipulation
- Type conversion functions
- Unicode and international support
- Performance-optimized primitives vs. high-level compositions

**Run with:**
```bash
./lisp-interpreter examples/string_library_demo.lisp
```

### Additional String Examples
- **`print_and_strings_demo.lisp`** - Output functions with string handling
- **`print_functions_demo.lisp`** - Advanced print capabilities and formatting
- **`simple_print_demo.lisp`** - Basic output operations

## Key Features Demonstrated (2025 Edition)

### Language Core
- ✅ Variables and functions with lexical scoping
- ✅ Lists and comprehensive list operations  
- ✅ Arithmetic (including big number support)
- ✅ Comparisons and logical operations
- ✅ Comments and documentation

### Modern Data Types (2025) ⭐
- ✅ **Keywords**: Self-evaluating symbols (`:name`, `:id`, `:status`)
- ✅ **Hash Maps**: Key-value associative arrays with immutable operations
- ✅ **Big Numbers**: Arbitrary precision integers with automatic overflow detection
- ✅ **Strings**: 20+ functions including regex support

### Advanced Language Features
- ✅ Higher-order functions (`map`, `filter`, `reduce`, `compose`)
- ✅ Closures and lexical scoping
- ✅ Recursion with tail call optimization (stack-safe)
- ✅ Error handling with `error` function and stack traces
- ✅ Pattern matching and conditional logic

### Module System & Organization
- ✅ Module definition with explicit exports
- ✅ Import system (qualified and unqualified)
- ✅ Qualified access (`module.function`)
- ✅ Private vs. public function visibility
- ✅ File loading and dependency management

### Development Environment
- ✅ Environment inspection (`env`, `modules`, `builtins`)
- ✅ Interactive REPL with help system
- ✅ File execution with multi-expression support
- ✅ Comprehensive error messages with context
- ✅ Built-in documentation and discovery tools

### Built-in Standard Library
- ✅ **Mathematical**: `factorial`, `fibonacci`, `gcd`, `lcm`, `abs`, `min`, `max`
- ✅ **List utilities**: `take`, `drop`, `all`, `any`, `length-sq`, `append`, `reverse`
- ✅ **Higher-order**: `compose`, `apply-n`, `curry`
- ✅ **String processing**: `string-concat`, `string-split`, `string-regex-*`, etc.
- ✅ **Hash map operations**: `hash-map-get`, `hash-map-put`, `hash-map-keys`, etc.
- ✅ **Meta functions**: `builtins`, `env`, `modules`, `type-of`

## Running Examples

### Individual Files (Complete Coverage)
```bash
# Core language features
./lisp-interpreter examples/basic_features.lisp
./lisp-interpreter examples/advanced_features.lisp

# Modern data types (2025)
./lisp-interpreter examples/keywords.lisp
./lisp-interpreter examples/hash_maps.lisp

# String processing capabilities
./lisp-interpreter examples/string_library_demo.lisp
./lisp-interpreter examples/print_and_strings_demo.lisp
./lisp-interpreter examples/print_functions_demo.lisp
./lisp-interpreter examples/simple_print_demo.lisp

# System architecture
./lisp-interpreter examples/module_system.lisp
./lisp-interpreter examples/core_library.lisp

# Run all examples in sequence
for file in examples/*.lisp; do
    echo "=== Running $file ==="
    ./lisp-interpreter "$file"
    echo
done
```

### Interactive REPL
```bash
# Start REPL and load examples
./lisp-interpreter
lisp> (load "examples/basic_features.lisp")
lisp> (load "examples/module_system.lisp")
```

### Building the Interpreter
```bash
# Build from source
make build

# Or using Go directly
go build -o lisp-interpreter cmd/lisp-interpreter/main.go
```

## Example Output

When you run the examples, you'll see output demonstrating:
- Function results and return values
- List transformations
- Module loading and imports
- Mathematical calculations
- Environment state

The examples are designed to showcase the current capabilities of the interpreter in a practical, easy-to-understand way.
