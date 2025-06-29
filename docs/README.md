# GoLisp Documentation

This directory contains comprehensive documentation for the GoLisp self-hosting compiler project.

## Documentation Overview

### Core Documentation

- **[COMPILER_ARCHITECTURE.md](COMPILER_ARCHITECTURE.md)** - Complete architecture documentation covering design patterns, data structures, compilation pipeline, and integration with the Go core
- **[COMPILER_API.md](COMPILER_API.md)** - Full API reference with usage examples, function signatures, and practical development scenarios
- **[SELF_HOSTING_GUIDE.md](SELF_HOSTING_GUIDE.md)** - Developer guide covering setup, workflows, testing, debugging, performance optimization, and contribution guidelines
- **[REPL.md](REPL.md)** - Comprehensive guide to the enhanced REPL with multi-line support, dynamic autocomplete, and interactive development features

### Quick Navigation

#### For New Developers
1. Start with **[REPL.md](REPL.md)** to get hands-on with the interactive environment
2. Read **[COMPILER_ARCHITECTURE.md](COMPILER_ARCHITECTURE.md)** to understand the overall design
3. Follow **[SELF_HOSTING_GUIDE.md](SELF_HOSTING_GUIDE.md)** for development setup and workflows
4. Reference **[COMPILER_API.md](COMPILER_API.md)** for specific function usage

#### For Contributors
1. Follow the setup instructions in **[SELF_HOSTING_GUIDE.md](SELF_HOSTING_GUIDE.md)**
2. Review the contribution guidelines in the same document
3. Use **[COMPILER_API.md](COMPILER_API.md)** for API examples and patterns

#### For Advanced Users
1. Study the compilation pipeline in **[COMPILER_ARCHITECTURE.md](COMPILER_ARCHITECTURE.md)**
2. Explore optimization techniques in **[COMPILER_API.md](COMPILER_API.md)**
3. Reference performance optimization in **[SELF_HOSTING_GUIDE.md](SELF_HOSTING_GUIDE.md)**

## Documentation Features

### Comprehensive Coverage
- **Architecture**: Multi-pass compilation pipeline, context system, optimization framework
- **API**: Complete function reference with 50+ documented functions and examples
- **Development**: Setup, testing, debugging, and contribution workflows

### Practical Examples
- **Code samples**: Hundreds of working examples throughout all documentation
- **Use cases**: Real-world scenarios for compiler development and usage
- **Troubleshooting**: Common issues and step-by-step solutions

### Professional Quality
- **Detailed explanations**: In-depth coverage of design decisions and trade-offs
- **Best practices**: Coding standards, testing approaches, and performance optimization
- **Future roadmap**: Clear guidance for extending the compiler with new features

## Key Topics Covered

### Architecture & Design
- Minimal Go core vs. self-hosted components
- Multi-pass compilation pipeline
- Context-driven compilation
- Modular evaluation engine
- Error handling and diagnostics

### API & Usage
- Core compilation functions
- Optimization system
- File compilation and bootstrap process
- Utility functions and helpers
- Error handling patterns

### Development & Contributing
- Environment setup and development tools
- Testing strategies and validation
- Debugging techniques and troubleshooting
- Performance optimization and profiling
- Code style and contribution guidelines

## Getting Started

```bash
# Clone and build
git clone https://github.com/your-org/go-lisp.git
cd go-lisp
make build

# Start with documentation
open docs/COMPILER_ARCHITECTURE.md

# Load the compiler
./bin/golisp -f lisp/self-hosting.lisp

# Try basic compilation
(compile-expr '(+ 1 2 3) (make-context))
```

## Documentation Maintenance

This documentation is actively maintained and updated with each major release. When contributing to the project:

1. **Update documentation** for any API changes
2. **Add examples** for new features
3. **Maintain consistency** with existing documentation style
4. **Test all examples** to ensure they work correctly

## Related Files

### Project Root
- **[CLAUDE.md](../CLAUDE.md)** - Project instructions and architecture overview
- **[ROADMAP.md](../ROADMAP.md)** - Development roadmap and progress tracking
- **[Makefile](../Makefile)** - Build targets and development commands

### Source Code
- **[pkg/core/](../pkg/core/)** - Go core implementation
- **[lisp/](../lisp/)** - Self-hosted Lisp code
- **[cmd/golisp/](../cmd/golisp/)** - CLI entry point

This documentation suite provides everything needed to understand, use, and contribute to the GoLisp self-hosting compiler project.