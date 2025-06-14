# Operations Reference

Complete guide to all operations and built-in functions in Go Lisp.

## Function Definition

### Function Definition with `defn`
```lisp
(defn function-name [param1 param2 ...] body)
```

**Examples:**
```lisp
(defn square [x] (* x x))
(defn add [a b] (+ a b))
(defn greet [name] (string-concat "Hello, " name "!"))
```

### Lambda Functions
```lisp
(fn [param1 param2 ...] body)
```

**Examples:**
```lisp
(fn [x] (* x x))
(fn [x y] (+ x y))
(map (fn [x] (* x 2)) (list 1 2 3 4))  ; => (2 4 6 8)
```

## Arithmetic Operations

### Basic Arithmetic
```lisp
(+ a b c ...)      ; Addition (multiple operands)
(- a b)            ; Subtraction (two operands)
(* a b c ...)      ; Multiplication (multiple operands)
(/ a b)            ; Division (two operands)
(% a b)            ; Modulo (remainder)
```

**Examples:**
```lisp
(+ 1 2 3 4)        ; => 10
(- 10 3)           ; => 7
(* 2 3 4)          ; => 24
(/ 15 3)           ; => 5
(% 17 5)           ; => 2
```

### Big Number Support
```lisp
; Automatic big number handling for large integers
(* 1000000000000000 1000000000000000)
; => 1000000000000000000000000000000

; Factorial with big numbers
(defn factorial [n]
  (if (= n 0) 1 (* n (factorial (- n 1)))))
(factorial 50)  ; Handles arbitrarily large results
```

## Comparison Operations

```lisp
(= a b)            ; Equality
(< a b)            ; Less than
(> a b)            ; Greater than
(<= a b)           ; Less than or equal
(>= a b)           ; Greater than or equal
```

**Examples:**
```lisp
(= 5 5)            ; => true
(< 3 5)            ; => true
(> 10 7)           ; => true
(<= 3 3)           ; => true
(>= 8 5)           ; => true
```

## Logical Operations

```lisp
(and expr1 expr2 ...)   ; Logical AND (short-circuiting)
(or expr1 expr2 ...)    ; Logical OR (short-circuiting)
(not expr)              ; Logical NOT
```

**Examples:**
```lisp
(and true true)        ; => true
(and true false)        ; => false
(or false true)         ; => true
(not true)           ; => false
```

## Control Flow

### Conditional Expressions
```lisp
(if condition then-expr else-expr)
```

**Examples:**
```lisp
(if (> x 0) "positive" "not positive")
(if (empty? my-list) "empty" "has elements")
```

### Sequential Evaluation
```lisp
(do expr1 expr2 ... exprN)
```

Evaluates expressions in sequence and returns the result of the last expression.

**Examples:**
```lisp
(do (def x 5) (def y 10) (+ x y))  ; => 15
(do (println! "Step 1") (println! "Step 2") "done")  ; => "done"
(do)  ; => nil
```

### Variable Definition
```lisp
(def name value)
```

**Examples:**
```lisp
(def x 42)
(def my-list (list 1 2 3))
(def square-fn (fn [x] (* x x)))
```

## List Operations

### List Creation
```lisp
(list elem1 elem2 ...)
```

**Examples:**
```lisp
(list 1 2 3 4 5)           ; => (1 2 3 4 5)
(list "a" "b" "c")         ; => ("a" "b" "c")
(list)                     ; => ()
```

### List Access
```lisp
(first lst)        ; Get first element
(rest lst)         ; Get all elements except first
(nth n lst)        ; Get nth element (0-indexed)
(length lst)       ; Get number of elements
(empty? lst)       ; Check if list is empty
```

**Examples:**
```lisp
(first (list 1 2 3))       ; => 1
(rest (list 1 2 3))        ; => (2 3)
(nth 1 (list "a" "b" "c")) ; => "b"
(length (list 1 2 3 4))    ; => 4
(empty? (list))            ; => true
```

### List Manipulation
```lisp
(cons elem lst)           ; Prepend element to list
(append lst1 lst2)        ; Concatenate lists
(reverse lst)             ; Reverse list order
(last lst)                ; Get last element
(butlast lst)             ; Get all but last element
(flatten lst)             ; Flatten nested lists
(distinct lst)            ; Remove duplicates
(concat lst1 lst2 ...)    ; Concatenate multiple lists
(partition n lst)         ; Split list into chunks of size n
(zip lst1 lst2 ...)       ; Combine lists by pairing elements
(sort lst [comparator])   ; Sort list (optionally with custom comparator)
```

**Examples:**
```lisp
(cons 0 (list 1 2 3))     ; => (0 1 2 3)
(append (list 1 2) (list 3 4))  ; => (1 2 3 4)
(reverse (list 1 2 3))    ; => (3 2 1)
(last (list 1 2 3 4))     ; => 4
(butlast (list 1 2 3 4))  ; => (1 2 3)
(flatten (list 1 (list 2 3) 4))  ; => (1 2 3 4)
(distinct (list 1 2 2 3 1))      ; => (1 2 3)
(concat (list 1 2) (list 3 4) (list 5))  ; => (1 2 3 4 5)
(partition 2 (list 1 2 3 4 5))   ; => ((1 2) (3 4) (5))
(zip (list 1 2 3) (list "a" "b" "c"))  ; => ((1 "a") (2 "b") (3 "c"))
(sort (list 5 2 8 1 9 3)) ; => (1 2 3 5 8 9)
```

## Higher-Order Functions

### Map, Filter, Reduce
```lisp
(map func lst)                    ; Apply function to each element
(filter predicate lst)            ; Keep elements matching predicate
(reduce func init lst)            ; Reduce list to single value
```

**Examples:**
```lisp
(map (fn [x] (* x x)) (list 1 2 3 4))
; => (1 4 9 16)

(filter (fn [x] (> x 0)) (list -1 2 -3 4))
; => (2 4)

(reduce (fn [acc x] (+ acc x)) 0 (list 1 2 3 4))
; => 10
```

## Hash Map Operations

### Hash Map Creation
```lisp
(hash-map key1 val1 key2 val2 ...)
```

**Examples:**
```lisp
(hash-map :name "Alice" :age 30)
(hash-map "x" 10 "y" 20)
```

### Hash Map Access
```lisp
(hash-map-get hm key)             ; Get value by key
(hash-map-put hm key val)         ; Add/update key-value pair
(hash-map-remove hm key)          ; Remove key-value pair
(hash-map-contains? hm key)       ; Check if key exists
(hash-map-keys hm)                ; Get all keys
(hash-map-values hm)              ; Get all values
(hash-map-size hm)                ; Get number of entries
(hash-map-empty? hm)              ; Check if empty
```

**Examples:**
```lisp
(def person (hash-map :name "Bob" :age 25))
(hash-map-get person :name)       ; => "Bob"
(hash-map-contains? person :age)  ; => true
(hash-map-keys person)            ; => (:name :age)
```

## String Operations

### String Functions
```lisp
(string-concat str1 str2 ...)     ; Concatenate strings
(string-length str)               ; Get string length
(string-substring str start end)  ; Extract substring
(string-char-at str index)        ; Get character at index
(string-upper str)                ; Convert to uppercase
(string-lower str)                ; Convert to lowercase
(string-trim str)                 ; Remove whitespace
```

**Examples:**
```lisp
(string-concat "Hello" " " "World")    ; => "Hello World"
(string-length "Hello")                ; => 5
(string-substring "Hello" 1 4)         ; => "ell"
(string-upper "hello")                 ; => "HELLO"
```

### String Predicates
```lisp
(string? value)                   ; Check if value is string
(string-empty? str)               ; Check if string is empty
(string-contains? str substr)     ; Check if contains substring
(string-starts-with? str prefix)  ; Check if starts with prefix
(string-ends-with? str suffix)    ; Check if ends with suffix
```

### String Conversion
```lisp
(string->number str)              ; Convert string to number
(number->string num)              ; Convert number to string
```

**Examples:**
```lisp
(string->number "42")             ; => 42
(number->string 42)               ; => "42"
```

## Macro System

### Macro Definition
```lisp
(defmacro name [param1 param2 ...] template)
```

**Examples:**
```lisp
(defmacro when [condition body]
  (list 'if condition body 'nil))

(defmacro unless [condition then else]
  (list 'if condition else then))
```

### Quote Special Form
```lisp
(quote expr)                      ; Prevent evaluation
'expr                             ; Shorthand for quote
```

**Examples:**
```lisp
(quote (+ 1 2))                   ; => (+ 1 2)
'(+ 1 2)                          ; => (+ 1 2)
(+ 1 2)                           ; => 3
```

## Module System

### Module Definition
```lisp
(module name
  (export symbol1 symbol2 ...)
  definitions...)
```

**Example:**
```lisp
(module math-utils
  (export square cube)
  
  (defn square [x] (* x x))
  (defn cube [x] (* x x x)))
```

### Module Loading and Import
```lisp
(load "filename.lisp")            ; Load file
(import module-name)              ; Import module exports
(require "filename.lisp")         ; Load and import in one step
```

**Examples:**
```lisp
(load "library/core.lisp")
(import core)
(factorial 10)                    ; Use imported function

; Or use qualified access
(core.factorial 10)

; Simplified with require
(require "library/core.lisp")
(factorial 10)
```

## Environment Inspection

### Development Tools
```lisp
(env)                             ; Show current environment
(modules)                         ; Show loaded modules
(help)                            ; Show all built-in functions
(help function-name)              ; Get help for specific function
```

**Examples:**
```lisp
(env)                             ; List all variables and functions
(modules)                         ; Show available modules
(help +)                          ; Get help for + function
```

## Concurrency Operations

### Goroutines (Parallel Execution)
```lisp
(go expression)                   ; Start goroutine, returns future
(go-wait future)                  ; Wait for single goroutine completion
(go-wait-all futures)             ; Wait for multiple goroutines
```

**Examples:**
```lisp
; Basic goroutine
(def result (go (* 5 5)))
(go-wait result)                  ; => 25

; Multiple goroutines
(def futures (list (go (+ 1 2)) (go (* 3 4))))
(go-wait-all futures)             ; => (3 12)
```

### Channels (Communication)
```lisp
(chan)                            ; Create unbuffered channel
(chan size)                       ; Create buffered channel
(chan-send! channel value)        ; Send value to channel
(chan-recv! channel)              ; Receive value from channel
(chan-close! channel)             ; Close channel
```

**Examples:**
```lisp
(def ch (chan 2))                 ; Buffered channel
(chan-send! ch "hello")
(chan-send! ch "world")
(chan-recv! ch)                   ; => "hello"
(chan-recv! ch)                   ; => "world"
```

## I/O Operations

### Print Functions
```lisp
(print! value1 value2 ...)         ; Print values separated by spaces
(println! value1 value2 ...)       ; Print values with newline
```

**Examples:**
```lisp
(print! "Hello" "World")           ; Prints: Hello World
(println! "Line 1")                ; Prints: Line 1\n
(println! "Line 2")                ; Prints: Line 2\n
```

### File Operations
```lisp
(read-file filename)              ; Read file contents as string
(write-file filename content)     ; Write string content to file
(file-exists? filename)           ; Check if file exists
```

**Examples:**
```lisp
; Read a file
(read-file "data.txt")            ; => "file contents..."

; Write to a file  
(write-file "output.txt" "Hello, World!")  ; => true

; Check file existence
(if (file-exists? "config.txt")
    (read-file "config.txt")
    "default config")

; Data processing pipeline
(def data (read-file "input.txt"))
(def processed (string-upper data))
(write-file "output.txt" processed)
```

### Error Handling
```lisp
(error message)                   ; Raise an error with message
```

**Example:**
```lisp
(defn safe-divide [a b]
  (if (= b 0)
      (error "Division by zero")
      (/ a b)))
```

## Type Predicates

```lisp
(string? value)                   ; Check if string
(number? value)                   ; Check if number  
(boolean? value)                  ; Check if boolean
(list? value)                     ; Check if list
(function? value)                 ; Check if function
(keyword? value)                  ; Check if keyword
(nil? value)                      ; Check if nil
```

## Keywords

Keywords are self-evaluating symbols perfect for hash map keys:

```lisp
:name                             ; => :name
:status                           ; => :status
:id                               ; => :id

; Common usage in hash maps
(hash-map :name "Alice" :age 30)
(:name person-map)                ; Keywords as accessor functions
```

## Constants

```lisp
nil                               ; Null/empty value (falsy)
true                                ; Boolean true
false                                ; Boolean false
```

This reference covers all the core operations available in Go Lisp. For more examples and advanced usage patterns, see the [Examples Guide](examples.md), [Mathematical Functions](mathematical_functions.md), and the various library documentation files.
