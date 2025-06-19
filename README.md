# Go Lisp

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

Go Lisp is a **functional Lisp interpreter** with a modern plugin architecture under active development:

- Core Language Support: Basic Lisp constructs implemented (arithmetic, lists, functions)
- Plugin Architecture: Modular design with 21+ plugin categories
- Interactive REPL: Basic REPL functionality for interactive development
- Data Types: Numbers, strings, booleans, lists, hash maps, and functions
- Hash Map Literals: Clojure-style syntax for easy hash map creation using `{:key "value"}`
- File Execution: Support for running Lisp programs from files
- Clean Architecture: Well-structured codebase with comprehensive testing

**Go Compatibility**: Go 1.24.2+  
**Platform Support**: Linux, macOS, Windows

**Note**: This is an actively developed project. Some advanced features mentioned in the documentation may be in various stages of implementation.

## Plugin Architecture

Go Lisp features a modular plugin system with the following categories:

- **arithmetic** - Basic arithmetic operations
- **atom** - Atomic operations and type checking  
- **comparison** - Comparison and equality operations
- **concurrency** - Concurrent programming primitives
- **control** - Control flow constructs
- **core** - Core language functionality (def, fn, quote, help, etc.)
- **functional** - Functional programming utilities
- **hashmap** - Hash map operations and utilities
- **http** - HTTP client functionality
- **io** - Input/output operations
- **json** - JSON parsing and manipulation
- **keyword** - Keyword support
- **list** - List processing and manipulation
- **logical** - Logical operations
- **math** - Mathematical functions and constants
- **string** - String processing and manipulation
- **binding** - Variable binding and scoping
- **environment** - Environment management
- **macro** - Macro system
- **sequence** - Sequence operations
- **advanced** - Advanced language features

*Note: Plugin functionality is in various stages of implementation and testing.*

## Key Features

- **Core Lisp Implementation**: Tokenizer, parser, and evaluator with plugin architecture
- **Interactive REPL**: Development environment for interactive programming
- **File Execution**: Run Lisp programs from files with multi-expression support
- **Basic Arithmetic**: Support for fundamental mathematical operations
- **List Processing**: Core list manipulation and processing functions
- **Plugin System**: Modular architecture with 21+ plugin categories
- **Modern Syntax**: Clean, readable Lisp syntax
- **Error Handling**: Clear error messages for debugging
- **Development Tools**: Command-line interface with multiple execution modes
- **Extensible Design**: Plugin-based architecture for easy feature addition

*Note: Advanced features like big number arithmetic, comprehensive string processing, HTTP clients, and other specialized functionality are in development.*

## REPL Features

The interactive REPL provides a development environment:

- **Interactive Execution**: Read-eval-print loop for immediate feedback
- **Multi-line Input**: Support for complex expressions
- **Basic Error Handling**: Clear error messages for common issues
- **Expression Evaluation**: Direct evaluation of Lisp expressions

*Note: Advanced REPL features like tab completion, syntax highlighting, and built-in help system are in development.*

## Architecture

Go Lisp is built with a modern, modular plugin architecture that provides clean separation of concerns and extensibility. The core components include:

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
- **List**: Collection operations (list creation, manipulation)
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
- Modern Go Lisp implementation with plugin architecture  
- Test-driven development methodologies
- Clean architecture principles in Go
- Modular design patterns
- Software engineering best practices

This project showcases both educational value and practical Lisp programming capabilities through a well-structured, extensible codebase.

## License

This project is open source and available under the MIT License. Built as both an educational resource and a practical tool for Lisp programming with modern software engineering practices.