# Usage

## Quick Start

```bash
# Interactive REPL
./lisp

# Run file
./lisp examples/basic_features.lisp

# Build from source
make build
# or: go build -o lisp-interpreter ./cmd/lisp-interpreter
```

## REPL Commands

- `(builtins)` - List all built-in functions  
- `(builtins <func>)` - Get help for specific function  
- `(env)` - Show current environment  
- `(modules)` - Show loaded modules  
- `quit` - Exit interpreter

## Testing

```bash
go test ./...
```

This runs all tests across all packages in the project.

## Using the Makefile

The project includes a Makefile with convenient commands:

```bash
make build    # Build the interpreter
make run      # Build and run the interpreter
make test     # Run all tests
```

## REPL Commands

Once in the REPL, you can:

- Enter any Lisp expression to evaluate it
- Type `quit` to exit the interpreter
- Use `(builtins)` to see all available functions
- Use `(builtins function-name)` to get help for a specific function
- Use `(env)` to inspect the current environment
- Use `(modules)` to see loaded modules

## Loading Files

You can load Lisp files from within the REPL:

### Traditional Approach
```lisp
lisp> (load "examples/math_module.lisp")
=> #<module:math>
lisp> (import math)
=> #<module:math>
lisp> (factorial 5)
=> 120
```

### Simplified with Require
```lisp
lisp> (require "examples/math_module.lisp")
=> #<module:math>
lisp> (factorial 5)
=> 120
```

The `require` function combines loading and importing in a single operation:
- **Loads** the file and executes all expressions
- **Automatically detects** any modules defined in the file
- **Imports** all exported functions into the current environment
- **Prevents** re-loading the same file multiple times

This makes it easier to work with library files and reduces the number of commands needed to use external modules.
