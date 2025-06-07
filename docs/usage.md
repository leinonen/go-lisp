# Usage Guide

This guide covers how to run, build, and work with the Lisp interpreter.

## Running the REPL

Start the interactive Read-Eval-Print Loop:

```bash
go run ./cmd/lisp-interpreter
```

The interpreter starts with a helpful message showing key commands:

```
Welcome to the Lisp Interpreter!
Type expressions to evaluate them, or 'quit' to exit.

Helpful commands:
  (builtins)        - List all available built-in functions
  (builtins <func>) - Get help for a specific function
  (env)             - Show current environment variables
  (modules)         - Show loaded modules
```

## Building the Interpreter

Create a standalone executable:

```bash
go build -o lisp-interpreter ./cmd/lisp-interpreter
./lisp-interpreter
```

## Running Tests

Execute the test suite:

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
