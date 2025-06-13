# File Execution Documentation

## Overview
Go Lisp supports executing Lisp code directly from files, allowing you to run scripts and load libraries without using the interactive REPL.

## Usage
```bash
./lisp <filename.lisp>
```

## Features

### Multi-Expression Support
Files can contain multiple expressions that are executed sequentially:

```lisp
; file: example.lisp
(def x 10)
(def y 20)
(+ x y)
(defn greet [name] (list "Hello" name))
(greet "World")
```

When executed:
```bash
./lisp example.lisp
```

Output:
```
30
(Hello World)
```

### Module Loading
Files can define and export modules:

```lisp
; file: math_utils.lisp
(module math-utils
  (export square cube)
  
  (defn square [x] (* x x))
  (defn cube [x] (* x x x))
  (defn helper [x] (+ x 1))  ; private function
)
```

### Library Usage
You can load and use libraries in scripts:

```lisp
; file: script.lisp
(load "examples/core.lisp")
(import core)

(factorial 10)
(fibonacci 15)
(gcd 48 18)
```

### Error Handling
The interpreter provides clear error messages for:
- File not found
- Syntax errors
- Runtime errors
- Module loading errors

## File Structure Best Practices

### Script Files
For executable scripts:
```lisp
; file: calculator.lisp
; Simple calculator script

(defn calculate [op a b]
  (if (= op "+") (+ a b)
      (if (= op "-") (- a b)
          (if (= op "*") (* a b)
              (if (= op "/") (/ a b)
                  (error "Unknown operation"))))))

; Run calculations
(calculate "+" 10 5)
(calculate "*" 7 8)
```

### Library Files
For reusable libraries:
```lisp
; file: string_utils.lisp
(module string-utils
  (export string-length string-empty? string-reverse)
  
  (defn string-length [s]
    ; Implementation...
    )
    
  (defn string-empty? [s]
    ; Implementation...
    )
    
  (defn string-reverse [s]
    ; Implementation...
    )
)
```

## Command Line Examples

### Basic Execution
```bash
# Execute a simple script
./lisp examples/basic_features.lisp

# Run mathematical calculations
./lisp examples/math_functions.lisp

# Execute core library demonstrations
./lisp examples/core_library.lisp
```

### Integration with Shell Scripts
```bash
#!/bin/bash
# Script to run multiple Lisp files

echo "Running core library tests..."
./lisp examples/core_library.lisp

echo "Running mathematical computations..."
./lisp examples/math_functions.lisp

echo "All tests completed."
```

## Error Messages

### File Not Found
```
Error: open nonexistent.lisp: no such file or directory
```

### Syntax Errors
```
Error in file script.lisp: unexpected token ')' at position 45
```

### Runtime Errors
```
Error in file script.lisp: undefined symbol 'undefined-function'
```

## Performance Considerations
- Files are tokenized and parsed completely before execution
- Large files are handled efficiently with streaming tokenization
- Module loading is cached to avoid re-processing
- Tail call optimization works in file execution mode

## Debugging
- Use the `env` function to inspect the current environment
- Add debug prints with explicit output
- Break large files into smaller, testable modules
- Use the REPL to test individual expressions

## See Also
- [Module System](modules.md)
- [Core Library](core_library.md)
- [REPL Usage](repl.md)
- [Examples](examples.md)
