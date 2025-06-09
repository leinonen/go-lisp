# Usage

## Quick Start

```bash
# Interactive REPL
./lisp-interpreter

# Run file
./lisp-interpreter examples/basic_features.lisp

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

```lisp
lisp> (load "examples/math_module.lisp")
=> #<module:math>
```

This executes all expressions in the file and makes any defined modules or functions available in the current session.
