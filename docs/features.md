# Features

This Lisp interpreter provides a comprehensive set of features for functional programming and interactive development.

## Core Components

- **Tokenizer/Lexer**: Converts input text into tokens for parsing
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens  
- **Evaluator**: Evaluates the AST to produce results
- **REPL**: Interactive Read-Eval-Print Loop with helpful startup commands

## Language Features

- **Comments**: Support for line comments using semicolons
- **Built-in Help System**: Discover functions with `(builtins)` and get help with `(builtins func-name)`
- **First-class Functions**: Functions can be stored in variables, passed as arguments, and returned from other functions
- **Closures**: Functions capture and remember variables from their creation environment
- **Recursion**: Functions can call themselves for recursive algorithms
- **Tail Call Optimization**: Prevents stack overflow in tail-recursive functions by eliminating stack growth
- **Higher-order Functions**: Functions that take other functions as arguments or return functions

## Data Types

- **Numbers**: `42`, `-3.14`
- **Strings**: `"hello world"`
- **Booleans**: `#t`, `#f`
- **Lists**: `(1 2 3)`, `("a" "b" "c")`, `()`
- **Symbols**: `+`, `-`, `x`, `my-var`
- **Functions**: `#<function([param1 param2])>`

## Development Features

- **Interactive Environment**: Full REPL with command history and helpful startup messages
- **Environment Inspection**: View current variables, functions, and loaded modules
- **Module System**: Organize code into modules with explicit exports and imports
- **File Loading**: Load and execute Lisp files from disk
- **Error Handling**: Clear error messages for debugging
