# Examples

This directory contains focused examples that demonstrate the key features of the Lisp interpreter. The examples have been consolidated to provide clear demonstrations without redundancy.

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

## Key Features Demonstrated

### Language Core
- ✅ Variables and functions
- ✅ Lists and list operations
- ✅ Arithmetic and comparisons
- ✅ Conditional logic
- ✅ Comments

### Advanced Features
- ✅ Higher-order functions (map, filter, reduce)
- ✅ Closures and lexical scoping
- ✅ Recursion with tail call optimization
- ✅ Big number arithmetic
- ✅ Error handling with `error` function

### Module System
- ✅ Module definition and exports
- ✅ Import system
- ✅ Qualified access (module.function)
- ✅ Private vs. public functions

### Development Tools
- ✅ Environment inspection (`env`, `modules`, `builtins`)
- ✅ Interactive REPL
- ✅ File execution
- ✅ Core library integration

### Built-in Functions
- ✅ Mathematical: factorial, fibonacci, gcd, lcm, abs, min, max
- ✅ List utilities: take, drop, all, any, length-sq
- ✅ Higher-order: compose, apply-n
- ✅ Meta: builtins, env, modules

## Running Examples

### Individual Files
```bash
# Run any example file
./lisp-interpreter examples/basic_features.lisp
./lisp-interpreter examples/module_system.lisp
./lisp-interpreter examples/core_library.lisp
./lisp-interpreter examples/advanced_features.lisp
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
