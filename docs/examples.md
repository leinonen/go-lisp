# GoLisp Standard Library Examples

This document provides working examples for all functions, variables, and macros defined in the GoLisp standard library.

## Table of Contents

- [Core Library Functions](#core-library-functions) (from `lisp/stdlib/core.lisp`)
- [Enhanced Library Functions](#enhanced-library-functions) (from `lisp/stdlib/enhanced.lisp`)
- [Control Flow](#control-flow) (loop/recur tail-call optimization)
- [Running Examples](#running-examples)

## Core Library Functions

### Logical Operations

#### `not`
Returns logical negation - `nil` if x is truthy, `true` if x is falsy.

```lisp
(not true)   ; => nil
(not false)  ; => true
(not nil)    ; => true
(not 42)     ; => nil
```

### Conditional Macros

#### `when`
Conditional execution: executes body when condition is truthy.

```lisp
(when true (println "This will print"))
; => This will print

(when false (println "This won't print"))
; => nil
```

#### `unless`
Conditional execution: executes body when condition is falsy.

```lisp
(unless false (println "This will print"))
; => This will print

(unless true (println "This won't print"))
; => nil
```

### Control Flow

#### `loop` and `recur`
Efficient tail-call optimization for recursive algorithms without stack growth.

```lisp
;; Basic loop with binding vector
(loop [x 5] x)
; => 5

;; Factorial using loop/recur
(loop [n 5 acc 1]
  (if (= n 0)
    acc
    (recur (- n 1) (* acc n))))
; => 120

;; Countdown example
(loop [i 3]
  (if (= i 0)
    "done"
    (recur (- i 1))))
; => "done"

;; Sum from 1 to n
(loop [n 5 sum 0]
  (if (= n 0)
    sum
    (recur (- n 1) (+ sum n))))
; => 15

;; Fibonacci sequence
(loop [n 6 a 0 b 1]
  (if (= n 0)
    a
    (recur (- n 1) b (+ a b))))
; => 8

;; Recur in user-defined functions
(defn factorial [n acc]
  (if (= n 0)
    acc
    (recur (- n 1) (* acc n))))

(factorial 6 1)
; => 720

;; Tail-recursive function for fibonacci
(defn fib-helper [n a b]
  (if (= n 0)
    a
    (recur (- n 1) b (+ a b))))

(fib-helper 7 0 1)
; => 13
```

### Collection Operations

#### `second`
Returns the second element of a collection.

```lisp
(second (list 1 2 3 4))  ; => 2
(second [10 20 30])      ; => 20
```

#### `third`
Returns the third element of a collection.

```lisp
(third (list 1 2 3 4 5))  ; => 3
(third [10 20 30 40])     ; => 30
```

#### `map`
Applies function to each element of collection, returning a new collection.

```lisp
(map inc (list 1 2 3))           ; => (2 3 4)
(map (fn [x] (* x 2)) [1 2 3])  ; => (2 4 6)
```

#### `filter`
Returns elements from collection where predicate function returns truthy.

```lisp
(filter even? (list 1 2 3 4 5))  ; => (2 4)
(filter pos? [-1 0 1 2])         ; => (1 2)
```

#### `range`
Returns a collection of integers from 0 to n-1 (in reverse order).

```lisp
(range 5)  ; => (4 3 2 1 0)
(range 3)  ; => (2 1 0)
```

#### `reduce`
Reduces collection to single value using function and initial value.

```lisp
(reduce + 0 (list 1 2 3 4))      ; => 10
(reduce * 1 (list 2 3 4))        ; => 24
(reduce max 0 (list 5 2 8 1))    ; => 8
```

#### `group-by`
Groups collection elements by key function, returns list of (key value) pairs.

```lisp
(group-by even? (list 1 2 3 4))  ; => ((true 4) (false 3) (true 2) (false 1))
(group-by (fn [x] (> x 5)) (list 1 7 3 9 2))  ; => ((false 2) (true 9) (false 3) (true 7) (false 1))
```

#### `map2`
Maps function over two collections simultaneously.

```lisp
(map2 + (list 1 2 3) (list 4 5 6))  ; => (5 7 9)
(map2 * [1 2 3] [2 3 4])            ; => (2 6 12)
```

### Hash Map Operations

#### `hash-map-put`
Associates a key-value pair in a hash map (wrapper around `assoc`).

```lisp
(hash-map-put {:a 1 :b 2} :c 3)  ; => {:a 1 :b 2 :c 3}
(hash-map-put {} :name "Alice")   ; => {:name "Alice"}
```

### Utility Variables

#### `length`
Alias for the `count` function.

```lisp
(length (list 1 2 3 4))  ; => 4
(length [1 2 3])         ; => 3
(length "hello")         ; => 5
```

## Enhanced Library Functions

### String Operations

#### `join`
Joins collection elements with separator string.

```lisp
(join ", " (list "apple" "banana" "cherry"))  ; => "apple, banana, cherry"
(join "-" [1 2 3 4])                          ; => "1-2-3-4"
```

#### String Aliases
Aliases for core string functions:

```lisp
; split - alias for string-split
(split "hello,world" ",")  ; => ("hello" "world")

; trim - alias for string-trim  
(trim "  hello  ")         ; => "hello"

; replace - alias for string-replace
(replace "hello world" "world" "GoLisp")  ; => "hello GoLisp"
```

### Function Application

#### `apply`
Applies function to arguments taken from collection (limited to 4 args).

```lisp
(apply + (list 1 2 3))      ; => 6
(apply max [5 2 8 1])       ; => 8
(apply str ["Hello" " " "World"])  ; => "Hello World"
```

### Collection Manipulation

#### `reverse`
Returns collection in reverse order.

```lisp
(reverse (list 1 2 3 4))  ; => (4 3 2 1)
(reverse [1 2 3])         ; => (3 2 1)
```

#### `take`
Returns first n elements from collection.

```lisp
(take 3 (list 1 2 3 4 5))  ; => (1 2 3)
(take 2 [10 20 30 40])     ; => (10 20)
```

#### `drop`
Returns collection with first n elements removed.

```lisp
(drop 2 (list 1 2 3 4 5))  ; => (3 4 5)
(drop 1 [10 20 30])        ; => (20 30)
```

#### `concat`
Concatenates two collections.

```lisp
(concat (list 1 2) (list 3 4))  ; => (1 2 3 4)
(concat [1 2] [3 4 5])          ; => (1 2 3 4 5)
```

#### `last`
Returns last element of collection.

```lisp
(last (list 1 2 3 4))  ; => 4
(last [10 20 30])      ; => 30
```

#### `butlast`
Returns collection with last element removed.

```lisp
(butlast (list 1 2 3 4))  ; => (1 2 3)
(butlast [10 20 30])      ; => (10 20)
```

#### `distinct`
Returns collection with duplicate elements removed.

```lisp
(distinct (list 1 2 2 3 1 4))  ; => (4 3 2 1)
(distinct [1 1 2 3 2 4])       ; => (4 3 2 1)
```

#### `contains-item?`
Tests if collection contains specific item.

```lisp
(contains-item? 3 (list 1 2 3 4))  ; => true
(contains-item? 5 [1 2 3 4])       ; => nil
```

#### `repeat`
Returns collection of n copies of value.

```lisp
(repeat 4 "hello")  ; => ("hello" "hello" "hello" "hello")
(repeat 3 42)       ; => (42 42 42)
```

#### `sort`
Returns collection sorted in ascending order.

```lisp
(sort (list 3 1 4 1 5))  ; => (1 1 3 4 5)
(sort [5 2 8 1 9])       ; => (1 2 5 8 9)
```

#### `partition`
Partitions collection into chunks of size n.

```lisp
(partition 2 (list 1 2 3 4 5 6))  ; => ((1 2) (3 4) (5 6))
(partition 3 [1 2 3 4 5 6 7 8])   ; => ((1 2 3) (4 5 6))
```

#### `interpose`
Inserts separator between each element of collection.

```lisp
(interpose ", " (list "a" "b" "c"))  ; => ("a" ", " "b" ", " "c")
(interpose 0 [1 2 3])                ; => (1 0 2 0 3)
```

#### `remove`
Returns elements from collection where predicate returns falsy.

```lisp
(remove even? (list 1 2 3 4 5))  ; => (1 3 5)
(remove neg? [-1 0 1 2])         ; => (0 1 2)
```

#### `keep`
Applies function to each element, keeping non-nil results.

```lisp
(keep (fn [x] (if (even? x) (* x 2))) (list 1 2 3 4 5))  ; => (4 8)
(keep identity [1 nil 2 nil 3])                          ; => (1 2 3)
```

#### `flatten`
Flattens nested collections into single-level collection.

```lisp
(flatten (list (list 1 2) (list 3 4)))  ; => (1 2 3 4)
(flatten [1 [2 3] [4 [5 6]]])           ; => (1 2 3 4 5 6)
```

### Arithmetic Operations

#### `inc`
Increments number by 1.

```lisp
(inc 5)    ; => 6
(inc -2)   ; => -1
```

#### `dec`
Decrements number by 1.

```lisp
(dec 5)    ; => 4
(dec 0)    ; => -1
```

#### `abs`
Returns absolute value of number.

```lisp
(abs -5)   ; => 5
(abs 3)    ; => 3
(abs 0)    ; => 0
```

#### `min`
Returns smallest of multiple numbers.

```lisp
(min 5 3)       ; => 3
(min -1 2 7)    ; => -1
(min 10 5 8 1)  ; => 1
```

#### `max`
Returns largest of multiple numbers.

```lisp
(max 5 3)       ; => 5
(max -1 2 7)    ; => 7
(max 10 5 8 1)  ; => 10
```

### Predicates

#### `zero?`
Tests if number equals zero.

```lisp
(zero? 0)   ; => true
(zero? 5)   ; => nil
(zero? -3)  ; => nil
```

#### `pos?`
Tests if number is positive.

```lisp
(pos? 5)    ; => true
(pos? -3)   ; => nil
(pos? 0)    ; => nil
```

#### `neg?`
Tests if number is negative.

```lisp
(neg? -5)   ; => true
(neg? 3)    ; => nil
(neg? 0)    ; => nil
```

#### `even?`
Tests if number is even.

```lisp
(even? 4)   ; => true
(even? 3)   ; => nil
(even? 0)   ; => true
```

#### `odd?`
Tests if number is odd.

```lisp
(odd? 3)    ; => true
(odd? 4)    ; => nil
(odd? 0)    ; => nil
```

#### `nil?`
Tests if value is nil.

```lisp
(nil? nil)    ; => true
(nil? false)  ; => nil
(nil? 0)      ; => nil
```

#### `some?`
Tests if value is not nil.

```lisp
(some? 42)    ; => true
(some? false) ; => true
(some? nil)   ; => nil
```

#### `true?`
Tests if value equals true.

```lisp
(true? true)   ; => true
(true? false)  ; => nil
(true? 1)      ; => nil
```

#### `false?`
Tests if value is falsy (nil).

```lisp
(false? nil)    ; => true
(false? false)  ; => true
(false? true)   ; => nil
```

#### `seq?`
Tests if value is a sequence (list or vector).

```lisp
(seq? (list 1 2 3))  ; => true
(seq? [1 2 3])       ; => true
(seq? {:a 1})        ; => nil
```

#### `coll?`
Tests if value is a collection (list or vector).

```lisp
(coll? (list 1 2 3))  ; => true
(coll? [1 2 3])       ; => true
(coll? "hello")       ; => nil
```

#### `all?`
Tests if predicate is true for all elements in collection.

```lisp
(all? even? (list 2 4 6))    ; => true
(all? pos? [1 2 3])          ; => true
(all? even? (list 1 2 3))    ; => nil
```

#### `any?`
Tests if predicate is true for any element in collection.

```lisp
(any? even? (list 1 2 3))    ; => true
(any? neg? [1 2 3])          ; => nil
(any? zero? [1 0 3])         ; => true
```

### Logical Operations

#### `and2`
Logical AND operation for two values.

```lisp
(and2 true true)    ; => true
(and2 true false)   ; => nil
(and2 42 "hello")   ; => "hello"
```

#### `or2`
Logical OR operation for two values.

```lisp
(or2 false true)    ; => true
(or2 nil "hello")   ; => "hello"
(or2 false nil)     ; => nil
```

### Functional Programming

#### `comp`
Returns function composition of f and g.

```lisp
(def add-one-then-double (comp (fn [x] (* x 2)) inc))
(add-one-then-double 5)  ; => 12
```

#### `constantly`
Returns function that always returns x regardless of arguments.

```lisp
(def always-42 (constantly 42))
(always-42)           ; => 42
(always-42 1 2 3)     ; => 42
```

#### `identity`
Returns its argument unchanged.

```lisp
(identity 42)      ; => 42
(identity "hello") ; => "hello"
(map identity [1 2 3])  ; => (1 2 3)
```

## Running Examples

To run these examples, you can:

1. **Interactive REPL**: Start the GoLisp REPL and type the examples directly:
   ```bash
   make run
   ```

2. **File execution**: Save examples to a `.lisp` file and run:
   ```bash
   ./bin/golisp -f examples.lisp
   ```

3. **Direct evaluation**: Run single examples from command line:
   ```bash
   ./bin/golisp -e "(map inc (list 1 2 3))"
   ```

**Note**: The standard library is automatically loaded when starting GoLisp, so all these functions are immediately available.

## Function Categories Summary

- **Collection Operations**: `map`, `filter`, `reduce`, `reverse`, `take`, `drop`, `concat`, `sort`, etc.
- **Predicates**: `nil?`, `even?`, `odd?`, `zero?`, `pos?`, `neg?`, `all?`, `any?`, etc.  
- **String Operations**: `join`, `split`, `trim`, `replace`
- **Math Operations**: `inc`, `dec`, `abs`, `min`, `max`
- **Logical Operations**: `not`, `and2`, `or2`
- **Functional Programming**: `apply`, `comp`, `constantly`, `identity`
- **Conditional Macros**: `when`, `unless`

For more information about the GoLisp language and core primitives, see the main [README.md](../README.md).