# GoLisp

![GoLisp logo](./docs/img/golisp-logo.png)

A modern Lisp dialect built on a minimal kernel architecture, combining the elegance of Clojure with the performance and simplicity of Go. Features a clean core interpreter with extensible modules for building powerful, fun programming experiences without the Java bloat.

## Quick Start

```bash
# Build go-lisp
make build

# Show help and available options
./bin/golisp -help

# Start interactive REPL
./bin/golisp

# Evaluate code directly from command line
./bin/golisp -e "(+ 1 2 3)"

# Execute a Lisp file (explicit flag)
./bin/golisp -f myprogram.lisp

```

### Command Line Options

| Option | Description | Example |
|--------|-------------|---------|
| `-help` | Show help message and usage examples | `./bin/golisp -help` |
| `-e <code>` | Evaluate Lisp code directly | `./bin/golisp -e "(+ 1 2)"` |
| `-f <file>` | Execute a Lisp file | `./bin/golisp -f program.lisp` |
| (none) | Start interactive REPL | `./bin/golisp` |

## Current Status

GoLisp is built on a **minimal Lisp kernel** with a modular architecture designed for extensibility and fun:

- **Minimal Core**: Clean, trusted kernel with essential primitives (~475 lines)
- **Self-Hosting**: Higher-level constructs implemented in Lisp itself
- **Macro System**: Full metaprogramming with quasiquote/unquote syntax
- **Modular Design**: Separate modules for different functionality areas
- **Interactive REPL**: Development environment with comprehensive examples
- **Data Structures**: Core types plus hash maps, vectors, and sets
- **File Loading**: Module system for organizing code

**Philosophy**: Like Clojure's elegance without Java's complexity, leveraging Go's strengths to create something new and enjoyable.

**Go Compatibility**: Go 1.24.2+  
**Platform Support**: Linux, macOS, Windows

## Architecture Philosophy

GoLisp follows a **minimal kernel** approach inspired by the best of both Clojure and Go:

### Minimal Lisp Kernel (`pkg/minimal/`)
The core is a tiny, trusted set of primitive operations that forms our "Lisp microkernel":

- **Essential Primitives**: Symbols, lists, vectors, environments, eval/apply
- **Special Forms**: `if`, `def`, `fn`, `quote`, `do`, `quasiquote`, `unquote`, `defmacro`
- **Macro System**: Full metaprogramming with code-as-data manipulation
- **Self-Hosting**: Higher-level constructs implemented in Lisp itself
- **Clean Separation**: ~475 lines maintaining clear concerns

### Modular Extension System
Built on the kernel, separate modules provide specialized functionality:

- **Domain-Specific Modules**: HTTP, JSON, math, string processing, concurrency
- **Go Integration**: Leverage Go's strengths (performance, concurrency, tooling)
- **Incremental Loading**: Add functionality as needed
- **Clean Interfaces**: Well-defined boundaries between kernel and modules

### Design Goals
- **Clojure's Elegance**: Modern Lisp design without Java's complexity
- **Go's Strengths**: Performance, simple deployment, excellent tooling
- **Minimal Core**: Keep the kernel small and trusted
- **Maximum Fun**: Create something new and enjoyable to use

## Key Features

### Minimal Kernel Design
- **Trusted Core**: Small, verifiable kernel with essential primitives
- **Self-Hosting**: Language constructs implemented in Lisp itself
- **Macro System**: Full metaprogramming capabilities with quasiquote/unquote
- **Clean Architecture**: Clear separation between kernel and extensions

### Clojure-Inspired Experience
- **Modern Syntax**: Square bracket `[param]` syntax for function parameters
- **Data Structures**: Vectors, hash maps, sets with Clojure-style literals
- **Functional Programming**: Immutable data structures and functional patterns
- **Polymorphic Functions**: Operations that work across different collection types

### Go Integration Benefits
- **Performance**: Native compilation without JVM overhead
- **Simple Deployment**: Single binary distribution
- **Concurrency**: Leverage Go's goroutines and channels
- **Ecosystem Access**: Integrate with Go libraries and tools

### Development Experience
- **Interactive REPL**: Rich development environment with examples
- **File Loading**: Module system for code organization
- **Standard Library**: Higher-level functions implemented in Lisp
- **Extensible Design**: Add new modules without touching the kernel

*Note: This project focuses on creating something new and fun by combining the best aspects of Clojure and Go.*

## Language Design

GoLisp combines the best of **Clojure** and **Go** to create something new and delightful:

### From Clojure
- **Modern Lisp Design**: Clean syntax with bracket notation `[param]` for parameters
- **Rich Data Structures**: Vectors, hash maps, sets with literal syntax
- **Functional Paradigm**: Immutable data and functional programming patterns
- **Macro System**: Full metaprogramming with quasiquote/unquote
- **Polymorphic Operations**: Functions that work across collection types

### From Go
- **Minimal Core**: Small, trusted kernel without unnecessary complexity
- **Performance**: Native compilation and efficient execution
- **Simplicity**: Easy deployment and straightforward tooling
- **Concurrency**: Future integration with Go's concurrency primitives

### New & Fun
- **Kernel Architecture**: Minimal, verifiable core with modular extensions
- **Self-Hosting**: Higher-level features implemented in the language itself
- **Module System**: Clean separation between core and specialized functionality
- **Development Joy**: Focus on creating something enjoyable to use

**Goal**: Capture Clojure's elegance without Java's bloat, while leveraging Go's strengths to build something both practical and fun.

## REPL Features

The interactive REPL provides a rich development environment built on the minimal kernel:

- **Interactive Execution**: Read-eval-print loop for immediate feedback
- **Macro Exploration**: Test metaprogramming features interactively
- **Standard Library**: Access to all bootstrap functions and higher-level constructs
- **File Loading**: Load and experiment with modular code
- **Examples**: Built-in examples demonstrating kernel capabilities

*Example REPL session:*
```lisp
user> (defmacro when [condition body] `(if ~condition ~body nil))
#<macro when>
user> (when true 42)
42
user> (load "stdlib.lisp")
; Standard library functions now available
```

## Architecture

GoLisp follows a **minimal kernel** approach with clean modular extensions:

### Minimal Lisp Kernel (`pkg/minimal/`)
The core implements a tiny, trusted set of primitives that form our "Lisp microkernel":

- **Core Types**: Symbols, lists, vectors, hash maps with interning
- **Environments**: Lexical scoping with parent chain lookup
- **Evaluator**: Essential eval/apply logic (~475 lines)
- **Special Forms**: `if`, `def`, `fn`, `quote`, `do`, `quasiquote`, `unquote`, `defmacro`
- **Macro System**: Full metaprogramming capabilities
- **Bootstrap**: Built-in functions implemented in Lisp itself
- **File Loading**: Module system for code organization

### Extension Modules (Future)
Specialized functionality built on the kernel foundation:
- **HTTP Module**: Web client and server capabilities
- **JSON Module**: Data serialization and parsing
- **Math Module**: Advanced mathematical functions
- **Concurrency Module**: Go-style concurrent programming
- **String Module**: Text processing utilities
- **File System Module**: File and directory operations

### Design Principles
- **Minimal Core**: Keep the kernel small and verifiable
- **Self-Hosting**: Implement features in Lisp when possible
- **Clean Interfaces**: Well-defined boundaries between components
- **Go Integration**: Leverage Go's strengths where appropriate

## Kernel Functions

The minimal kernel provides essential functionality organized into logical categories:

### Core Language
- **Definition**: `def` for variables, `fn` for functions, `defmacro` for macros
- **Control Flow**: `if`, `do`, `quote` for essential program structure
- **Metaprogramming**: Quasiquote (`` ` ``) and unquote (`~`) for macro templates

### Data Structures
- **Lists**: `list`, `first`, `rest`, `cons` for list processing
- **Vectors**: Square bracket notation `[1 2 3]` with indexing
- **Hash Maps**: `hash-map`, `hash-map-get`, `hash-map-put`, `hash-map-keys`
- **Type Checking**: `list?`, `vector?`, `map?`, `nil?`, `keyword?`

### Arithmetic & Comparison
- **Basic Math**: `+`, `-`, `*`, `/`, `%` for arithmetic operations
- **Comparisons**: `=`, `<`, `>`, `<=`, `>=` for ordering and equality
- **Logic**: `and`, `or`, `not` for boolean operations

### Self-Hosted Features
Built using the kernel primitives and implemented in Lisp:
- **Control Macros**: `when`, `unless` for conditional execution
- **Higher-Order Functions**: Functional programming utilities
- **Standard Library**: Collection of useful functions in `stdlib.lisp`

*The beauty of this approach: most functionality is implemented in Lisp itself, keeping the kernel minimal while maximizing expressiveness.*

## Documentation

The project includes comprehensive documentation (though some may reference the previous architecture):

### Core Documentation
- Complete feature overview and data types
- Guide to all supported operations  
- Mathematical function reference
- File system operations and I/O
- HTTP client and JSON processing
- Code examples and tutorials
- Usage and building guide

### Modern Language Features
- **[Modern Functions Guide](docs/modern_functions.md)** - Complete reference for modern function names and aliases
- **[Quick Reference](docs/quick_reference.md)** - Cheat sheet for modern syntax and functions
- **[Examples](examples/modern_functions_demo.lisp)** - Practical examples demonstrating modern features

### Development Documentation
- Implementation design and plugin architecture
- Future enhancements and roadmap
- Contributing guidelines and standards

## Building and testing

### Quick Start
```bash
# Clone and build
git clone https://github.com/leinonen/go-lisp.git
cd go-lisp
make build

# Testing
make test
```

## Contributing

GoLisp is built on a **minimal kernel architecture** that makes contributions focused and impactful:

### Kernel Development
- **Core Improvements**: Enhance the minimal kernel (performance, features, correctness)
- **Macro System**: Expand metaprogramming capabilities
- **Standard Library**: Add useful functions implemented in Lisp
- **Testing**: Comprehensive test coverage for kernel reliability

### Module Development
- **New Modules**: Create specialized functionality (HTTP, JSON, concurrency, etc.)
- **Go Integration**: Bridge Lisp and Go ecosystems effectively
- **Domain Expertise**: Contribute modules in your area of expertise
- **Clean Interfaces**: Maintain clear boundaries between kernel and modules

### Documentation & Examples
- **Kernel Guide**: Document the minimal core and its philosophy
- **Module Examples**: Show how to build on the kernel foundation
- **Real-World Apps**: Demonstrate practical uses of the language
- **Learning Materials**: Help others understand Lisp and metaprogramming

### Development Philosophy
- **Minimal Kernel**: Keep the core small, trusted, and verifiable
- **Self-Hosting**: Prefer implementing features in Lisp over Go when possible
- **Clean Architecture**: Maintain clear separation between components
- **Fun First**: Create something enjoyable to use and extend

The kernel approach makes it easy to contribute meaningfully while maintaining the project's core vision of combining Clojure's elegance with Go's strengths.

## Recognition

**Built in 2025** as a demonstration of:
- **Minimal Kernel Architecture**: Small, trusted core with modular extensions
- **Self-Hosting Language Design**: Higher-level features implemented in Lisp itself  
- **Modern Lisp Innovation**: Combining Clojure's elegance with Go's practical benefits
- **Clean Software Engineering**: Separation of concerns and testable components
- **Fun Programming Experience**: Creating something new and enjoyable to use

This project showcases how to build a modern Lisp that's both educationally valuable and practically useful, proving that you can have Clojure's expressiveness without Java's complexity.

## License

This project is open source and available under the MIT License. Built as both an educational resource and a practical tool for exploring modern Lisp design, demonstrating how to combine the best of Clojure and Go to create something new and fun.