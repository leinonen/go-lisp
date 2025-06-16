# Go Lisp

A modern, production-ready Lisp dialect implemented in Go with a modular plugin architecture. Features comprehensive language support, advanced data types, functional programming utilities, and a powerful plugin system.

## Quick Start

```bash
# Build go-lisp
make build

# Show help and available options
./lisp -help

# Start interactive REPL
./lisp

# Evaluate code directly from command line
./lisp -e "(+ 1 2 3)"

# Execute a Lisp file (explicit flag)
./lisp -f myprogram.lisp

# Execute a Lisp file (legacy positional argument)
./lisp myprogram.lisp

# Development mode (without building)
go run ./cmd/go-lisp -e "(* 6 7)"
go run ./cmd/go-lisp -f examples/basic_features.lisp
```

## Command Line Usage

Go Lisp supports multiple modes of operation through command line parameters:

### Interactive REPL Mode
```bash
./lisp
# Starts the interactive Read-Eval-Print Loop
```

### Direct Code Evaluation
```bash
# Evaluate expressions directly from the command line
./lisp -e "(+ 1 2 3)"           # => 6
./lisp -e "(list 1 2 3 4 5)"    # => (1 2 3 4 5)
./lisp -e "(* 6 7)"             # => 42

# Modern square bracket function syntax
./lisp -e "(defn square [x] (* x x))"
./lisp -e "((fn [x y] (+ x y)) 10 20)"  # => 30

# Perfect for quick calculations and one-liners
./lisp -e "(% 1000000000000000001 7)"  # => 0
```

### File Execution
```bash
# Execute Lisp programs from files
./lisp -f script.lisp           # Explicit file flag
./lisp script.lisp              # Legacy positional argument (still supported)
```

### Help and Information
```bash
./lisp -help                    # Show all available options
```

### Command Line Options

| Option | Description | Example |
|--------|-------------|---------|
| `-help` | Show help message and usage examples | `./lisp -help` |
| `-e <code>` | Evaluate Lisp code directly | `./lisp -e "(+ 1 2)"` |
| `-f <file>` | Execute a Lisp file | `./lisp -f program.lisp` |
| (none) | Start interactive REPL | `./lisp` |

**Note**: The interpreter maintains backward compatibility - you can still use `./lisp filename.lisp` without the `-f` flag.

## Current Status

Go Lisp is **feature-complete** and **production-ready** with a modern plugin architecture:

- Complete Language Support: All core Lisp constructs implemented
- Plugin Architecture: Modular design with 15 plugin categories providing 122+ functions
- Modern Tooling: REPL with help system, tab completion, and color-coded error messages
- Advanced Data Types: Big numbers, hash maps, keywords, strings, lists, and functions
- Performance Optimized: Efficient evaluation with proper tail call handling
- High Code Quality: Clean architecture, comprehensive testing, TDD development

**Go Compatibility**: Go 1.24.2+  
**Platform Support**: Linux, macOS, Windows

## Plugin Architecture

Go Lisp features a modular plugin system with the following categories:

- **arithmetic** (5 functions) - Basic arithmetic operations
- **atom** (4 functions) - Atomic operations and type checking  
- **comparison** (5 functions) - Comparison and equality operations
- **concurrency** (9 functions) - Concurrent programming primitives
- **control** (2 functions) - Control flow constructs
- **core** (7 functions) - Core language functionality (def, fn, quote, help, etc.)
- **functional** (4 functions) - Functional programming utilities
- **hashmap** (9 functions) - Hash map operations and utilities
- **http** (4 functions) - HTTP client functionality
- **io** (5 functions) - Input/output operations
- **json** (4 functions) - JSON parsing and manipulation
- **list** (10 functions) - List processing and manipulation
- **logical** (3 functions) - Logical operations
- **math** (30 functions) - Mathematical functions and constants
- **string** (21 functions) - String processing and manipulation

Use `(plugins)` in the REPL to see all loaded plugins, or `(help)` to see all available functions.

## Key Features

- **Complete Lisp Implementation**: Full tokenizer, parser, and evaluator with modern plugin architecture
- **Interactive REPL**: Rich development environment with built-in help system and tab completion
- **File Execution**: Run Lisp programs from files with full multi-expression support
- **Big Number Arithmetic**: Arbitrary precision integers with automatic overflow detection
- **Advanced Data Types**: Lists, hash maps, keywords, strings, and functions as first-class citizens
- **Performance Optimized**: Efficient evaluation with proper tail call handling
- **Plugin System**: Modular architecture with 15 plugin categories providing 122+ functions
- **String Processing**: Comprehensive string manipulation with 21+ built-in functions
- **Error Handling**: Color-coded error messages with helpful suggestions
- **Development Tools**: Environment inspection, debugging helpers, and comprehensive help system
- **Mathematical Functions**: 30+ built-in math functions (trigonometry, logarithms, statistics, constants)
- **Functional Programming**: Complete functional programming utilities with higher-order functions
- **HTTP Client**: Full REST API support with GET, POST, PUT, DELETE operations
- **JSON Processing**: Parse, stringify, and extract data with comprehensive JSON support
- **Concurrency Support**: 9 concurrent programming primitives for parallel processing
- **Hash Map Operations**: 9 functions for associative data structure manipulation
- **Keywords Support**: Self-evaluating symbols perfect for hash map keys
- **Modern Syntax**: Square bracket function parameters for improved readability

## REPL Features

The interactive REPL provides a rich development environment:

- **Tab Completion**: Press TAB to see available functions and complete names
- **Color-coded Errors**: Different error types displayed with helpful suggestions
- **Multi-line Input**: Automatic detection of balanced parentheses
- **Built-in Help**: Use `(help)` for all functions, `(help function-name)` for specific help
- **Environment Inspection**: Use `(env)` to see current variables and functions
- **Plugin Information**: Use `(plugins)` to see all loaded plugin categories

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

The 122+ built-in functions are organized into logical categories:

- **Core Functions**: Variable definition, function creation, quoting, help system
- **Arithmetic**: Basic mathematical operations (+, -, *, /, %)
- **Comparison**: Equality and ordering operations (=, <, >, <=, >=)
- **Logical**: Boolean logic operations (and, or, not)
- **Math**: Advanced mathematical functions (trigonometry, logarithms, statistics)
- **String**: Text processing and manipulation (length, substring, split, join)
- **List**: Collection operations (map, filter, reduce, first, rest)
- **Hash Map**: Associative data structure operations (get, put, keys, values)
- **I/O**: Input/output operations (print, println, read-file, write-file)
- **HTTP**: Web client functionality (GET, POST, PUT, DELETE requests)
- **JSON**: Data serialization (parse, stringify, extract)
- **Atom**: Atomic operations for concurrent programming
- **Control**: Flow control constructs (if, cond)
- **Functional**: Higher-order programming utilities (curry, compose, partial)
- **Concurrency**: Parallel processing primitives (spawn, await, channel operations)

Use `(help)` in the REPL to explore all available functions, or `(help category-name)` for category-specific functions.

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

## Building and Testing

### Quick Start
```bash
# Clone and build
git clone https://github.com/leinonen/lisp-interpreter.git
cd go-lisp
make build

# Start the REPL
./lisp

# Try some examples
./lisp -e "(* 1000000000000000000 1000000000000000000)"
./lisp -e "(plugins)"  # See all available plugin categories
./lisp -e "(help)"     # See all available functions
```

### Development Commands
```bash
make build    # Build the interpreter binary
make run      # Build and run the interpreter (REPL mode)  
make test     # Run comprehensive test suite

# Development mode
go run ./cmd/go-lisp                    # Start REPL without building
go run ./cmd/go-lisp -e "(+ 1 2 3)"    # Evaluate expression
go run ./cmd/go-lisp -f script.lisp    # Execute file

# Plugin exploration
./lisp -e "(plugins)"                  # List all plugin categories
./lisp -e "(help math)"                # Math functions (if implemented)
./lisp -e "(help string)"              # String functions (if implemented)
```

### Requirements and Compatibility
```bash
# System requirements  
go version     # Requires Go 1.24.2 or later
make --version # GNU Make for build automation

# Platform support
# Linux (primary development platform)
# macOS (full compatibility)  
# Windows (cross-compiled)
```

### Manual Build
```bash
go build -o lisp ./cmd/go-lisp
./lisp
```

## Project Structure

The project follows a clean, modular architecture:

```
go-lisp/
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

This project features a modern, plugin-based architecture built with production-quality standards. We welcome contributions in several areas:

### Areas for Contribution
- **Bug Reports**: Help us maintain reliability by reporting issues
- **Documentation**: Improve guides, add tutorials, or enhance examples  
- **Performance**: Optimize critical paths and memory usage
- **Plugin Development**: Add new function categories or extend existing ones
- **Tool Integration**: IDE plugins, syntax highlighting, or language servers
- **Examples**: Real-world applications and algorithm implementations
- **Testing**: Edge cases, stress testing, and platform validation

### Development Standards
- **Test-Driven Development**: All features require comprehensive tests
- **Plugin Architecture**: New functionality should follow the modular plugin pattern
- **Documentation First**: New features need documentation and examples
- **Clean Architecture**: Maintain separation of concerns and modularity
- **Go Best Practices**: Follow Go idioms and conventions
- **Backward Compatibility**: Preserve existing functionality

The plugin system makes it easy to extend functionality while maintaining clean separation of concerns.

## Recognition

**Built in 2025** as a comprehensive demonstration of:
- Modern Go Lisp implementation with plugin architecture
- Test-driven development methodologies  
- Clean architecture principles in Go
- Modular design patterns
- Production-quality software engineering practices

## License

This project is open source and available under the MIT License. Built as both an educational resource and a practical tool for Lisp programming with modern software engineering practices.
