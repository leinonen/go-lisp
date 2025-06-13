# Print Functions Documentation

## Overview

The Lisp interpreter includes built-in `print!` and `println!` functions that enable programs to produce output. These functions are essential for creating interactive programs, debugging, displaying results, and creating user-friendly applications. The exclamation mark indicates that these functions perform side effects (output to stdout).

```lisp
;; Clear and modern function syntax
(def format-greeting (fn [name] 
  (string-concat "Welcome, " name "!")))

;; Nested functions are easier to read
(def process-data (fn [items processor]
  (map (fn [item] (processor item)) items)))
```

## Functions

### `print!`

**Syntax:** `(print! arg1 arg2 ... argN)`

**Description:** Outputs the given arguments to stdout without adding a newline at the end. Multiple arguments are separated by spaces.

**Arguments:** 
- Any number of arguments of any type
- Arguments are evaluated and converted to their string representation

**Returns:** `nil` (to avoid duplicate output in REPL/file execution)

**Examples:**
```lisp
(print! "Hello")                    ; Output: Hello
(print! "Count:" 42)               ; Output: Count: 42
(print! 1 2 3 #t #f)               ; Output: 1 2 3 #t #f
```

### `println!`

**Syntax:** `(println! arg1 arg2 ... argN)`

**Description:** Outputs the given arguments to stdout and adds a newline at the end. Multiple arguments are separated by spaces.

**Arguments:**
- Any number of arguments of any type  
- Arguments are evaluated and converted to their string representation

**Returns:** `nil` (to avoid duplicate output in REPL/file execution)

**Examples:**
```lisp
(println! "Hello World")            ; Output: Hello World\n
(println! "Result:" (+ 2 3))        ; Output: Result: 5\n
(println!)                          ; Output: \n (empty line)
```

## Data Type Support

Both functions support all Lisp data types:

### Basic Types
```lisp
(println! "String:" "Hello")        ; String: Hello
(println! "Number:" 42)             ; Number: 42
(println! "Boolean:" #t)            ; Boolean: #t
(println! "Nil:" nil)               ; Nil: ()
```

### Big Numbers
```lisp
(println! "Big number:" 123456789012345678901234567890)
; Output: Big number: 123456789012345678901234567890
```

### Collections
```lisp
(println! "List:" (list 1 2 3))     ; List: (1 2 3)
(println! "Hash:" (hash-map "a" 1)) ; Hash: {a: 1}
```

### Functions
```lisp
(def square (fn [x] (* x x)))
(println! "Function:" square)       ; Function: #<function([x])>
```

## Usage Patterns

### Basic Output
```lisp
(println! "Program started")
(println! "Processing data...")
(println! "Program completed")
```

### Variable Display
```lisp
(def name "Alice")
(def age 25)
(println! "Name:" name "Age:" age)  ; Name: Alice Age: 25
```

### Mathematical Results
```lisp
(def a 10)
(def b 5)
(println! a "+" b "=" (+ a b))      ; 10 + 5 = 15
(println! a "*" b "=" (* a b))      ; 10 * 5 = 50
```

### List Processing
```lisp
(def numbers (list 1 2 3 4 5))
(println! "Original:" numbers)
(println! "Squared:" (map (fn [x] (* x x)) numbers))
```

### String Operations Integration
```lisp
(def text "Hello World")
(println! "Text:" text)
(println! "Length:" (string-length text))
(println! "Upper:" (string-upper text))
```

### Formatted Output
```lisp
(def format-greeting
  (fn [name]
    (string-concat "Welcome, " name "!")))

(println! (format-greeting "Bob"))  ; Welcome, Bob!
```

### Progress Indication
```lisp
(def show-progress
  (fn [current total]
    (println! "Progress:" current "/" total)))

(show-progress 3 10)               ; Progress: 3 / 10
```

### Menu Systems
```lisp
(def show-menu
  (fn []
    (begin
      (println! "===== MENU =====")
      (println! "1. Option A")
      (println! "2. Option B")
      (println! "================"))));
```

## Difference Between `print!` and `println!`

The key difference is the automatic newline:

```lisp
(print! "Hello")
(print! " ")
(print! "World")
; Output: Hello World

(println! "Hello")
(println! "World")  
; Output: Hello
;         World
```

Combined usage:
```lisp
(print! "Processing")
(print! ".")
(print! ".")
(print! ".")
(println! " Done!")
; Output: Processing... Done!
```

## Implementation Notes

- Both functions return `nil` to prevent duplicate output when used in REPL or file execution
- The interpreter automatically suppresses `nil` values from being printed in command-line and file execution modes
- All argument types are converted to strings using their standard string representation
- Multiple arguments are automatically space-separated
- Empty `(println)` call produces a blank line

## Examples Files

See the following example files for comprehensive demonstrations:

- `examples/simple_print.lisp` - Basic print function usage
- `examples/print_and_strings.lisp` - Integration with string operations
- `examples/string_library.lisp` - Advanced string processing with output

## Integration with Other Features

The print functions work seamlessly with all other interpreter features:

- **String Functions:** Display results of string operations
- **Mathematical Operations:** Show calculation results
- **List Processing:** Display transformed lists
- **Hash Maps:** Show key-value data
- **Function Results:** Display return values
- **Error Handling:** Can be used before calling `(error)` for debugging
- **Module System:** Available in all modules
- **REPL Integration:** Work perfectly in interactive mode

This makes the print functions essential tools for creating complete, user-friendly Lisp programs.
