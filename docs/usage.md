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

- `(help)` - List all built-in functions  
- `(help <func>)` - Get help for specific function  
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
- Use `(help)` to see all available functions
- Use `(help function-name)` to get help for a specific function
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

## File Operations

The interpreter provides built-in file system operations for reading, writing, and managing files:

```lisp
lisp> (write-file "hello.txt" "Hello, World!")
=> #t
lisp> (file-exists? "hello.txt")
=> #t
lisp> (read-file "hello.txt")
=> "Hello, World!"
```

Common file operation patterns:

```lisp
; Configuration file management
lisp> (if (file-exists? "config.txt")
         (read-file "config.txt")
         "default configuration")

; Data processing pipeline
lisp> (def data (read-file "input.txt"))
lisp> (def processed (string-upper data))
lisp> (write-file "output.txt" processed)
```

For comprehensive file operations documentation, see [File Functions](file_functions.md).
