# Lisp Interpreter

A basic Lisp interpreter implemented in Go using Test-Driven Development (TDD).

## Features

- **Tokenizer/Lexer**: Converts input text into tokens
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens  
- **Evaluator**: Evaluates the AST to produce results
- **REPL**: Interactive Read-Eval-Print Loop
- **Comments**: Support for line comments using semicolons

## Supported Operations

### Arithmetic
- `(+ 1 2 3)` - Addition with multiple operands
- `(- 10 3)` - Subtraction
- `(* 2 3 4)` - Multiplication with multiple operands
- `(/ 15 3)` - Division

### Comparison
- `(= 5 5)` - Equality
- `(< 3 5)` - Less than
- `(> 7 3)` - Greater than

### Conditional
- `(if condition then-expr else-expr)` - If expression

### Variables
- `(define name value)` - Define a variable with a name and value

### Functions
- `(lambda (params) body)` - Create an anonymous function
- `(defun name (params) body)` - Define a named function (combines define and lambda)
- `(funcname args...)` - Call a user-defined function

### Lists
- `(list)` - Create an empty list
- `(list 1 2 3)` - Create a list with elements
- `(first lst)` - Get the first element of a list
- `(rest lst)` - Get all elements except the first
- `(cons elem lst)` - Prepend an element to a list
- `(length lst)` - Get the number of elements in a list
- `(empty? lst)` - Check if a list is empty
- `(append lst1 lst2)` - Combine two lists into one
- `(reverse lst)` - Reverse the order of elements in a list
- `(nth index lst)` - Get the element at a specific index (0-based)

### Higher-Order Functions
- `(map func lst)` - Apply a function to each element of a list
- `(filter predicate lst)` - Keep only elements that satisfy a predicate
- `(reduce func init lst)` - Reduce a list to a single value using a function

### Comments
- `;` - Line comments (from semicolon to end of line are ignored)
- Comments can appear anywhere in the code
- Useful for documenting code and adding explanations

### Module System
- `(module name (export sym1 sym2...) body...)` - Define a module with exported symbols
- `(import module-name)` - Import all exported symbols from a module into current scope
- `(load "filename.lisp")` - Load and execute a Lisp file
- `module.symbol` - Qualified access to module symbols without importing

### Advanced Function Features
- **First-class functions**: Functions can be stored in variables, passed as arguments, and returned from other functions
- **Closures**: Functions capture and remember variables from their creation environment
- **Recursion**: Functions can call themselves for recursive algorithms
- **Higher-order functions**: Functions that take other functions as arguments or return functions

### Data Types
- Numbers: `42`, `-3.14`
- Strings: `"hello world"`
- Booleans: `#t`, `#f`
- Lists: `(1 2 3)`, `("a" "b" "c")`, `()`
- Symbols: `+`, `-`, `x`, `my-var`
- Functions: `#<function([param1 param2])>`

## Usage

### Running the REPL
```bash
go run ./cmd/lisp-interpreter
```

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o lisp-interpreter ./cmd/lisp-interpreter
./lisp-interpreter
```

### Using the Makefile
```bash
make build    # Build the interpreter
make run      # Build and run the interpreter
make test     # Run all tests
```

## Examples

### Basic Function Definition

Traditional way (using define + lambda):
```lisp
lisp> (define square (lambda (x) (* x x)))
=> <function>

lisp> (square 5)
=> 25
```

New convenient way (using defun):
```lisp
lisp> (defun square (x) (* x x))
=> <function>

lisp> (square 5)
=> 25
```

### Multi-parameter Functions
```lisp
lisp> (defun add (x y) (+ x y))
=> <function>

lisp> (add 3 4)
=> 7
```

### Recursive Functions
```lisp
lisp> (defun factorial (n) 
        (if (= n 0) 
          1 
          (* n (factorial (- n 1)))))
=> <function>

lisp> (factorial 5)
=> 120
```

### Basic Operations
```lisp
lisp> 42
=> 42

lisp> (+ 1 2 3)
=> 6

lisp> (define x 10)
=> 10

lisp> (+ x 5)
=> 15
```

### Examples with Comments
```lisp
; This is a comment - it will be ignored
lisp> (+ 1 2 3) ; Comments can appear at the end of lines
=> 6

; Define a function with comments
lisp> (defun factorial (n) ; Calculate factorial recursively
        (if (= n 0)        ; Base case: 0! = 1
          1 
          (* n (factorial (- n 1))))) ; Recursive case
=> <function>

lisp> (factorial 5) ; Test the function
=> 120
```

### More Examples

```lisp
lisp> (* (+ 2 3) 4)
=> 20

lisp> (if (< 3 5) 100 0)
=> 100

lisp> "hello world"
=> hello world

lisp> (= 5 5)
=> #t

lisp> (define x 10)
=> 10

lisp> x
=> 10

lisp> (define y (* x 3))
=> 30

lisp> (+ x y)
=> 40

lisp> (lambda (x) (+ x 1))
=> #<function([x])>

lisp> (define add1 (lambda (x) (+ x 1)))
=> #<function([x])>

lisp> (add1 5)
=> 6

lisp> (define factorial (lambda (n) (if (= n 0) 1 (* n (factorial (- n 1))))))
=> #<function([n])>

lisp> (factorial 5)
=> 120

lisp> (define make-adder (lambda (n) (lambda (x) (+ x n))))
=> #<function([n])>

lisp> (define add10 (make-adder 10))
=> #<function([x])>

lisp> (add10 7)
=> 17

lisp> (list 1 2 3)
=> (1 2 3)

lisp> (list)
=> ()

lisp> (list "hello" 42 #t)
=> (hello 42 #t)

lisp> (define my-list (list 10 20 30))
=> (10 20 30)

lisp> (first my-list)
=> 10

lisp> (rest my-list)
=> (20 30)

lisp> (length my-list)
=> 3

lisp> (empty? my-list)
=> #f

lisp> (empty? (list))
=> #t

lisp> (cons 5 my-list)
=> (5 10 20 30)

lisp> (append (list 1 2) (list 3 4 5))
=> (1 2 3 4 5)

lisp> (reverse (list 1 2 3 4))
=> (4 3 2 1)

lisp> (nth 0 my-list)
=> 10

lisp> (nth 2 my-list)
=> 30

lisp> (define sum-list (lambda (lst) (if (empty? lst) 0 (+ (first lst) (sum-list (rest lst))))))
=> #<function([lst])>

lisp> (sum-list (list 1 2 3 4))
=> 10
```

### Module System Examples

```lisp
; Define a math module with exported functions
lisp> (module math
        (export square cube factorial)
        (defun square (x) (* x x))
        (defun cube (x) (* x x x))
        (defun factorial (n)
          (if (= n 0)
            1
            (* n (factorial (- n 1)))))
        (defun private-helper (x) (+ x 1))) ; not exported
=> #<module:math>

; Import all exported functions from the math module
lisp> (import math)
=> #<module:math>

; Use imported functions directly
lisp> (square 5)
=> 25

lisp> (factorial 4)
=> 24

; Try to use private function (should fail)
lisp> (private-helper 5)
=> Error: undefined symbol: private-helper

; Use qualified access without importing
lisp> (module string-utils
        (export concat reverse-string)
        (defun concat (a b) a) ; simplified
        (defun reverse-string (s) s)) ; simplified
=> #<module:string-utils>

; Access functions using qualified names
lisp> (string-utils.concat "Hello" "World")
=> Hello

lisp> (math.cube 3)
=> 27

; Load a file containing module definitions
lisp> (load "examples/math_module.lisp")
=> #<module:math>

; Create modules with complex functionality
lisp> (module data-processing
        (export process-list filter-numbers)
        (defun process-list (lst)
          (map square (filter positive? lst)))
        (defun filter-numbers (lst)
          (filter (lambda (x) (> x 10)) lst))
        (defun positive? (x) (> x 0))) ; helper function
=> #<module:data-processing>

lisp> (import data-processing)
=> #<module:data-processing>

lisp> (process-list (list -2 3 4 5))
=> (9 16 25)
```

### Higher-Order Function Examples
```lisp
; Map: apply a function to each element
lisp> (map (lambda (x) (* x x)) (list 1 2 3 4 5))
=> (1 4 9 16 25)

; Filter: keep only elements that satisfy a predicate
lisp> (filter (lambda (x) (> x 0)) (list -1 2 -3 4 5))
=> (2 4 5)

; Reduce: combine all elements using a function
lisp> (reduce (lambda (acc x) (+ acc x)) 0 (list 1 2 3 4 5))
=> 15

; Complex combination: squares of positive numbers
lisp> (map (lambda (x) (* x x)) (filter (lambda (x) (> x 0)) (list -2 -1 0 1 2 3)))
=> (1 4 9)

; Sum of squares using map and reduce
lisp> (reduce (lambda (acc x) (+ acc x)) 0 (map (lambda (x) (* x x)) (list 1 2 3)))
=> 14

; Count elements using reduce
lisp> (reduce (lambda (acc x) (+ acc 1)) 0 (list "a" "b" "c" "d"))
=> 4
```

### Additional List Operations Examples

```lisp
; Append: combine two lists
lisp> (append (list 1 2 3) (list 4 5 6))
=> (1 2 3 4 5 6)

lisp> (append (list) (list 1 2 3))
=> (1 2 3)

lisp> (append (list "hello" "world") (list "from" "lisp"))
=> (hello world from lisp)

; Reverse: reverse the order of elements
lisp> (reverse (list 1 2 3 4 5))
=> (5 4 3 2 1)

lisp> (reverse (list "a" "b" "c"))
=> (c b a)

lisp> (reverse (list))
=> ()

; Nth: get element at specific index (0-based)
lisp> (define colors (list "red" "green" "blue" "yellow"))
=> (red green blue yellow)

lisp> (nth 0 colors)
=> red

lisp> (nth 2 colors)
=> blue

lisp> (nth 3 colors)
=> yellow

; Combining operations for complex data manipulation
lisp> (reverse (append (list 1 2) (list 3 4)))
=> (4 3 2 1)

lisp> (nth 1 (reverse (list 10 20 30 40)))
=> 30

lisp> (map (lambda (i) (nth i (list "a" "b" "c" "d"))) (list 0 2 1 3))
=> (a c b d)
```

## Architecture

The interpreter follows a traditional three-phase design:

1. **Tokenization**: Convert input string into tokens
2. **Parsing**: Build AST from tokens
3. **Evaluation**: Evaluate AST in an environment

### Test-Driven Development

This project was built using TDD principles:
- Tests were written first for each component
- Implementation followed to make tests pass
- Tests serve as documentation and ensure correctness

### File Structure

The project follows Go standard layout with modular architecture:

```
lisp-interpreter/
├── go.mod                       # Go module definition
├── README.md                    # Project documentation
├── Makefile                     # Build automation
├── cmd/
│   └── lisp-interpreter/
│       └── main.go              # REPL and main program
└── pkg/
    ├── types/
    │   ├── types.go             # Core types and interfaces
    │   └── types_test.go        # Type system tests
    ├── tokenizer/
    │   ├── tokenizer.go         # Lexical analysis
    │   └── tokenizer_test.go    # Tokenizer tests
    ├── parser/
    │   ├── parser.go            # Syntax analysis
    │   └── parser_test.go       # Parser tests
    ├── evaluator/
    │   ├── evaluator.go         # Expression evaluation
    │   └── evaluator_test.go    # Evaluator tests
    └── interpreter/
        ├── interpreter.go       # High-level interpreter API
        └── interpreter_test.go  # Integration tests
```

### Package Overview

- **`pkg/types`** - Core type definitions (Token, Expr, Value interfaces)
- **`pkg/tokenizer`** - Lexical analysis (string → tokens)
- **`pkg/parser`** - Syntax analysis (tokens → AST)
- **`pkg/evaluator`** - Expression evaluation (AST → values)
- **`pkg/interpreter`** - High-level API combining all components
- **`cmd/lisp-interpreter`** - Main application with REPL

## Future Enhancements

- ✅ List data structures and operations (implemented: `list`, `first`, `rest`, `cons`, `length`, `empty?`)
- ✅ More built-in list functions (implemented: `map`, `filter`, `reduce`)
- ✅ Additional list operations (implemented: `append`, `reverse`, `nth`)
- Error recovery in parser
- Better error messages with line numbers
- ✅ Support for comments
- Tail call optimization for recursive functions
- ✅ Module system for code organization (implemented: `module`, `import`, `load`, qualified access)
- Macro system for code transformation
