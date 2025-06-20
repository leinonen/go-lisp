# GoLisp

![GoLisp logo](./docs/img/golisp-logo.png)

A modern, production-ready Lisp dialect implemented in Go with a modular plugin architecture. Features comprehensive language support, advanced data types, functional programming utilities, and a powerful plugin system.

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

GoLisp is a **production-ready Lisp interpreter** with a modern dependency injection plugin architecture:

- **Core Language Support**: Complete Lisp constructs with comprehensive functionality
- **Modern Plugin Architecture**: Dependency injection-based modular design with 21+ plugin categories
- **Interactive REPL**: Full-featured REPL with tab completion and error handling
- **Advanced Data Types**: Numbers, big numbers, strings, booleans, lists, vectors, hash maps, atoms, and functions
- **Hash Map Literals**: Clojure-style syntax for easy hash map creation using `{:key "value"}`
- **File Execution**: Robust support for running Lisp programs from files
- **Clean Architecture**: Well-structured codebase with comprehensive testing and dependency injection
- **Production Ready**: Clean separation of concerns, proper error handling, and extensive test coverage

**Go Compatibility**: Go 1.24.2+  
**Platform Support**: Linux, macOS, Windows

**Recent Updates**: ✅ **Refactoring Complete** - Successfully implemented dependency injection throughout the entire plugin system, removing all legacy constructors and enforcing modern architectural patterns. All tests pass and the codebase now uses clean, dependency injection-based design patterns exclusively.

## Plugin Architecture

GoLisp features a modern **dependency injection-based plugin system** with the following categories:

### Core Plugins (Fully Refactored)
- **arithmetic** - Basic arithmetic operations with proper dependency injection
- **atom** - Atomic operations and thread-safe mutable state  
- **comparison** - Comparison and equality operations
- **concurrency** - Concurrent programming primitives (goroutines, channels)
- **control** - Control flow constructs (if, do, cond, when, loop/recur)
- **core** - Core language functionality (def, fn, quote, help, etc.)
- **hashmap** - Hash map operations and utilities with polymorphic support
- **io** - Input/output operations and file handling
- **keyword** - Keyword support and utilities
- **list** - List processing and manipulation
- **math** - Mathematical functions and constants
- **binding** - Variable binding and scoping (let)
- **environment** - Enhanced environment management (let*, letfn)
- **macro** - Macro system (defmacro, macroexpand)
- **polymorphic** - Cross-type polymorphic functions
- **utils** - Utility functions (frequencies, partition, shuffle, etc.)

### Additional Plugins
- **functional** - Functional programming utilities
- **http** - HTTP client functionality
- **json** - JSON parsing and manipulation
- **logical** - Logical operations
- **string** - String processing and manipulation
- **sequence** - Sequence operations
- **advanced** - Advanced language features

**Architecture**: ✅ **Modernization Complete** - All core plugins now use modern dependency injection patterns exclusively. Legacy constructors have been completely removed from the codebase. The plugin system enforces proper dependency management, enables comprehensive testing with mock evaluators, and provides clean separation of concerns. All 16 core plugins have been refactored and verified with passing tests.

### Testing Strategy
The refactoring employed comprehensive testing methodologies:
- **Mock Evaluators**: All plugin tests use dependency injection with mock evaluators for isolated testing
- **Integration Testing**: Full end-to-end testing ensures plugins work together correctly
- **Regression Testing**: All existing functionality verified to work after architectural changes
- **Test Coverage**: Each refactored plugin maintains or improves test coverage
- **Continuous Verification**: Tests run after each plugin refactoring to catch issues early

## Key Features

- **Core Lisp Implementation**: Tokenizer, parser, and evaluator with plugin architecture
- **Interactive REPL**: Development environment for interactive programming
- **File Execution**: Run Lisp programs from files with multi-expression support
- **Clojure-Style Functions**: Modern Clojure-compatible function names and aliases
- **Polymorphic Functions**: Type-aware functions that work across different data types (lists, vectors, strings, hashmaps)
- **Unified Sequence Operations**: Functions like `first`, `rest`, `last`, `nth`, `count`, `empty?` work on all collection types
- **Cross-Type Compatibility**: `get` and `contains?` work on hashmaps, vectors, lists, and strings
- **Basic Arithmetic**: Support for fundamental mathematical operations
- **List Processing**: Core list manipulation and processing functions
- **Plugin System**: Modular architecture with 21+ plugin categories
- **Modern Syntax**: Clean, readable Lisp syntax with Clojure-style literals
- **Error Handling**: Clear error messages for debugging
- **Development Tools**: Command-line interface with multiple execution modes
- **Extensible Design**: Plugin-based architecture for easy feature addition

*Note: Advanced features like big number arithmetic, comprehensive string processing, HTTP clients, and other specialized functionality are in development.*

## Language Design

GoLisp draws major inspiration from **Clojure**, incorporating modern Lisp design principles and idiomatic function names. The language features:

- **Clojure-Compatible Functions**: Familiar function names like `get`, `assoc`, `dissoc`, `contains?`, `keys`, `vals`
- **Polymorphic Operations**: Functions like `first`, `rest`, `last`, `nth`, `count`, `empty?` work across all collection types (lists, vectors, strings, hashmaps)
- **Cross-Type Compatibility**: Advanced functions like `get` and `contains?` work seamlessly on hashmaps, vectors, lists, and strings
- **Unified Sequence Interface**: All collections can be treated as sequences with consistent behavior
- **Modern Syntax**: Clean, readable syntax with support for hash map literals `{:key "value"}`
- **Functional Paradigm**: Emphasis on immutable data structures and functional programming patterns
- **Practical Design**: Balance between academic correctness and real-world usability

While maintaining its own identity, GoLisp aims to provide a familiar experience for developers coming from Clojure while offering the performance and deployment advantages of Go.

## REPL Features

The interactive REPL provides a development environment:

- **Interactive Execution**: Read-eval-print loop for immediate feedback
- **Multi-line Input**: Support for complex expressions
- **Basic Error Handling**: Clear error messages for common issues
- **Expression Evaluation**: Direct evaluation of Lisp expressions

*Note: Advanced REPL features like tab completion, syntax highlighting, and built-in help system are in development.*

## Architecture

GoLisp is built with a modern, modular plugin architecture that provides clean separation of concerns and extensibility. The core components include:

### Core Components
- **Tokenizer**: Lexical analysis converting source code to tokens
- **Parser**: Builds abstract syntax trees from token streams  
- **Evaluator**: Executes parsed expressions with environment management
- **Interpreter**: Coordinates parsing and evaluation
- **REPL**: Interactive read-eval-print loop with advanced features
- **Executor**: Handles file execution and command-line evaluation

### Plugin System
The plugin architecture allows for modular functionality:
- **Plugin Registry**: Central registration and discovery system
- **Function Registry**: Type-safe function registration with metadata
- **Category System**: Logical grouping of related functionality
- **Environment Integration**: Seamless integration with the core evaluator

All functionality is organized into focused plugins, making the codebase maintainable and extensible while providing a comprehensive set of built-in capabilities.

## Function Categories

The built-in functions are organized into logical categories:

- **Core Functions**: Variable definition, function creation, quoting
- **Arithmetic**: Basic mathematical operations (+, -, *, /, %)
- **Comparison**: Equality and ordering operations
- **Logical**: Boolean logic operations
- **Polymorphic Functions**: Type-aware sequence operations (`first`, `rest`, `last`, `nth`, `second`, `empty?`, `take`, `drop`, `reverse`)
- **Collection Operations**: Cross-type functions (`get`, `contains?`, `count`) that work on all data types
- **List**: List-specific operations (list creation, cons, append)
- **Control**: Flow control constructs (if, cond)

*Note: Additional categories like advanced math, string processing, HTTP clients, JSON handling, and other specialized functions are implemented as plugins in various stages of development.*

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

## Project Structure

The project follows a clean, modular architecture:

```
go-lisp/
├── bin/golisp             # Built binary
├── cmd/go-lisp/           # Main application entry point
├── pkg/                   # Core packages
│   ├── evaluator/         # Expression evaluation and environment
│   ├── executor/          # File execution coordination  
│   ├── functions/         # Function definition and metadata
│   ├── interpreter/       # Main interpreter logic
│   ├── parser/           # Syntax analysis and AST generation
│   ├── plugins/          # Plugin system and implementations
│   │   ├── arithmetic/   # Basic math operations
│   │   ├── atom/         # Atomic operations
│   │   ├── comparison/   # Comparison operations
│   │   ├── concurrency/  # Concurrent programming
│   │   ├── control/      # Control flow
│   │   ├── core/         # Core language features
│   │   ├── functional/   # Functional programming
│   │   ├── hashmap/      # Hash map operations
│   │   ├── http/         # HTTP client
│   │   ├── io/           # Input/output operations
│   │   ├── json/         # JSON processing
│   │   ├── list/         # List operations
│   │   ├── logical/      # Boolean logic
│   │   ├── math/         # Mathematical functions
│   │   └── string/       # String processing
│   ├── registry/         # Function and plugin registration
│   ├── repl/             # Interactive REPL with completion
│   ├── tokenizer/        # Lexical analysis
│   └── types/            # Core data types and expressions
└── Makefile              # Build automation
```

## Contributing

This project features a modern, plugin-based architecture built with clean coding standards. We welcome contributions in several areas:

### Areas for Contribution
- **Bug Reports**: Help us maintain reliability by reporting issues
- **Feature Implementation**: Complete plugin functionality and add new capabilities
- **Documentation**: Improve guides, add tutorials, or enhance examples  
- **Testing**: Add test coverage and edge case validation
- **Performance**: Optimize critical paths and memory usage
- **Tool Integration**: IDE plugins, syntax highlighting, or language servers
- **Examples**: Real-world applications and algorithm implementations

### Development Standards
- **Test-Driven Development**: All features require comprehensive tests
- **Plugin Architecture**: New functionality should follow the modular plugin pattern
- **Documentation**: New features need documentation and examples
- **Clean Architecture**: Maintain separation of concerns and modularity
- **Go Best Practices**: Follow Go idioms and conventions

The plugin system makes it easy to extend functionality while maintaining clean separation of concerns.

## Recognition

**Built in 2025** as a demonstration of:
- Modern GoLisp implementation with dependency injection plugin architecture  
- Test-driven development methodologies with comprehensive refactoring
- Clean architecture principles and dependency injection patterns in Go
- Modular design patterns with proper separation of concerns
- Software engineering best practices including legacy code modernization

This project showcases both educational value and practical Lisp programming capabilities through a well-structured, fully modernized, and extensively tested codebase.

## License

This project is open source and available under the MIT License. Built as both an educational resource and a practical tool for Lisp programming with modern software engineering practices.