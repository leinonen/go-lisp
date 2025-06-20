# Modern Polymorphic Functions in Go-Lisp

This document describes the modern function names and polymorphic behavior that have been added to make Go-Lisp more compatible with contemporary Lisp syntax and conventions, particularly those inspired by Clojure.

## Polymorphic Sequence Functions

Go-Lisp now features true polymorphic functions that work across all sequence and collection types, providing a unified interface similar to Clojure. These functions automatically adapt their behavior based on the input type.

### Core Polymorphic Functions

| Function | Description | Works On |
|----------|-------------|----------|
| `first` | Get first element | Lists, vectors, strings |
| `rest` | Get all but first element | Lists, vectors, strings |
| `last` | Get last element | Lists, vectors, strings |
| `nth` | Get nth element (0-indexed) | Lists, vectors, strings |
| `second` | Get second element | Lists, vectors, strings |
| `empty?` | Check if collection is empty | All collections |
| `count` | Get number of elements | All collections |
| `get` | Get element by key/index | Hashmaps, vectors, lists, strings |
| `contains?` | Check if contains key/index | Hashmaps, vectors, lists, strings |

### Polymorphic Sequence Functions Examples

#### `first` - Get first element of any sequence
```lisp
;; Lists
(first '(1 2 3))        ; => 1
(first '())             ; => nil

;; Vectors  
(first ["a" "b" "c"])   ; => "a"
(first [])              ; => nil

;; Strings
(first "hello")         ; => "h"
(first "")              ; => nil
```

#### `rest` - Get all elements except the first
```lisp
;; Lists
(rest '(1 2 3))         ; => (2 3)
(rest '(1))             ; => ()
(rest '())              ; => ()

;; Vectors
(rest ["a" "b" "c"])    ; => ("b" "c")  ; Always returns a list
(rest ["a"])            ; => ()

;; Strings  
(rest "hello")          ; => ("e" "l" "l" "o")  ; Returns list of characters
(rest "a")              ; => ()
```

#### `last` - Get last element of any sequence
```lisp
;; Lists
(last '(1 2 3))         ; => 3
(last '())              ; => nil

;; Vectors
(last ["a" "b" "c"])    ; => "c"
(last [])               ; => nil

;; Strings
(last "hello")          ; => "o"
(last "")               ; => nil
```

#### `nth` - Get nth element (0-indexed)
```lisp
;; Lists
(nth '("a" "b" "c" "d") 2)    ; => "c"
(nth '(1 2 3) 0)              ; => 1

;; Vectors
(nth [10 20 30 40] 1)         ; => 20
(nth [] 0)                    ; => nil (out of bounds)

;; Strings
(nth "hello" 1)               ; => "e"
(nth "hello" 4)               ; => "o"
```

#### `second` - Get second element of any sequence
```lisp
;; Lists
(second '(1 2 3))       ; => 2
(second '(1))           ; => nil

;; Vectors
(second ["a" "b" "c"])  ; => "b"
(second [])             ; => nil

;; Strings
(second "hello")        ; => "e"
(second "h")            ; => nil
```

### Polymorphic Collection Functions

#### `empty?` - Check if any collection is empty
```lisp
;; Lists
(empty? '())            ; => true
(empty? '(1 2 3))       ; => false

;; Vectors
(empty? [])             ; => true
(empty? [1 2 3])        ; => false

;; Hashmaps
(empty? {})             ; => true
(empty? {:a 1})         ; => false

;; Strings
(empty? "")             ; => true
(empty? "hello")        ; => false

;; Nil
(empty? nil)            ; => true
```

#### `count` - Get count of elements in any collection
```lisp
;; Lists
(count '(1 2 3 4))        ; => 4
(count '())               ; => 0

;; Vectors
(count [a b c])           ; => 3
(count [])                ; => 0

;; Hash maps (counts key-value pairs)
(count {:a 1 :b 2 :c 3})  ; => 3
(count {})                ; => 0

;; Strings (counts characters, not bytes)
(count "hello")           ; => 5
(count "")                ; => 0
(count "héllo")           ; => 5  ; Unicode-aware

;; Nil
(count nil)               ; => 0
```

### Polymorphic Access Functions

#### `get` - Get element by key or index from any collection
```lisp
;; Hashmaps
(def my-map {:name "Alice" :age 30})
(get my-map :name)        ; => "Alice"
(get my-map :city)        ; => nil

;; Vectors (by index)
(def my-vec [10 20 30 40])
(get my-vec 0)            ; => 10
(get my-vec 2)            ; => 30
(get my-vec 10)           ; => nil

;; Lists (by index)  
(def my-list '("a" "b" "c" "d"))
(get my-list 1)           ; => "b"
(get my-list 3)           ; => "d"

;; Strings (by index)
(get "hello" 1)           ; => "e"
(get "hello" 4)           ; => "o"
```

#### `contains?` - Check if collection contains key/index or substring
```lisp
;; Hashmaps (check for key)
(contains? my-map :name)      ; => true
(contains? my-map :city)      ; => false

;; Vectors (check for valid index)
(contains? my-vec 2)          ; => true (index 2 exists)
(contains? my-vec 10)         ; => false (index 10 doesn't exist)

;; Strings (check for index or substring)
(contains? "hello" 2)         ; => true (index 2 exists)
(contains? "hello" "ell")     ; => true (substring exists)
(contains? "hello" "xyz")     ; => false
```

### Polymorphic Transformation Functions

#### `take` - Take first n elements from any sequence
```lisp
;; Lists
(take 2 '(1 2 3 4 5))         ; => (1 2)
(take 0 '(1 2 3))             ; => ()

;; Vectors
(take 3 ["a" "b" "c" "d" "e"]) ; => ("a" "b" "c")

;; Strings
(take 2 "hello")              ; => ("h" "e")
```

#### `drop` - Drop first n elements from any sequence
```lisp
;; Lists
(drop 2 '(1 2 3 4 5))         ; => (3 4 5)
(drop 0 '(1 2 3))             ; => (1 2 3)

;; Vectors
(drop 1 ["a" "b" "c" "d"])    ; => ("b" "c" "d")

;; Strings
(drop 2 "hello")              ; => ("l" "l" "o")
```

#### `reverse` - Reverse any sequence
```lisp
;; Lists
(reverse '(1 2 3))            ; => (3 2 1)

;; Vectors
(reverse ["a" "b" "c"])       ; => ["c" "b" "a"]

;; Strings
(reverse "hello")             ; => "olleh"
```

### Type Predicate Functions

#### Collection Type Predicates
```lisp
;; seq? - Check if value is a sequence (lists, vectors, strings)
(seq? '(1 2 3))               ; => true
(seq? [1 2 3])                ; => true
(seq? "abc")                  ; => true
(seq? {:a 1})                 ; => false

;; coll? - Check if value is a collection (lists, vectors, hashmaps)
(coll? '(1 2 3))              ; => true
(coll? [1 2 3])               ; => true
(coll? {:a 1})                ; => true
(coll? "abc")                 ; => false

;; sequential? - Check if value is sequential (lists, vectors)
(sequential? '(1 2 3))        ; => true
(sequential? [1 2 3])         ; => true
(sequential? {:a 1})          ; => false

;; indexed? - Check if value supports indexed access (vectors, strings)
(indexed? [1 2 3])            ; => true
(indexed? "abc")              ; => true
(indexed? '(1 2 3))           ; => false
```

## Hash Map Functions

### Core Functions with Modern Aliases

| Go-Lisp Function | Modern Alias | Description |
|------------------|---------------|-------------|
| `hash-map-get` | `get` | Get value from hash map by key (now polymorphic) |
| `hash-map-contains?` | `contains?` | Check if hash map contains key (now polymorphic) |
| `hash-map-keys` | `keys` | Get all keys from hash map |
| `hash-map-values` | `vals` | Get all values from hash map |

### New Clojure-Style Functions

#### `assoc` - Associate key-value pairs
```lisp
;; Basic usage
(def my-map (hash-map "name" "Alice" "age" 30))
(assoc my-map "city" "New York")  ; => {"name" "Alice" "age" 30 "city" "New York"}

;; Multiple key-value pairs
(assoc my-map "city" "Boston" "country" "USA")
; => {"name" "Alice" "age" 30 "city" "Boston" "country" "USA"}

;; With keywords
(def person {:name "Bob" :age 25})
(assoc person :job "Developer" :salary 75000)
```

#### `dissoc` - Dissociate keys from hash map
```lisp
(def user {:name "Charlie" :age 35 :email "charlie@email.com" :phone "555-1234"})

;; Remove single key
(dissoc user :phone)  ; => {:name "Charlie" :age 35 :email "charlie@email.com"}

;; Remove multiple keys
(dissoc user :email :phone)  ; => {:name "Charlie" :age 35}
```

#### Examples using polymorphic aliases
```lisp
;; Using polymorphic get - works on hashmaps, vectors, lists, strings
(def data {:users [{:name "Alice"} {:name "Bob"}] :count 2})
(get data :count)                     ; => 2
(get (get data :users) 0)            ; => {:name "Alice"}
(get (get (get data :users) 0) :name) ; => "Alice"

;; Using polymorphic contains? - works on all collection types
(contains? data :users)  ; => true (hashmap key)
(contains? data :admin)  ; => false

(contains? [1 2 3] 1)    ; => true (vector index)
(contains? [1 2 3] 5)    ; => false

(contains? "hello" 2)    ; => true (string index)
(contains? "hello" "ell") ; => true (substring)

;; Using keys and vals (hashmap-specific)
(keys data)   ; => ("users" "count")
(vals data)   ; => ([{:name "Alice"} {:name "Bob"}] 2)
```

## Arithmetic Functions

### New Increment/Decrement Functions

#### `inc` - Increment a number
```lisp
(inc 5)      ; => 6
(inc -3)     ; => -2
(inc 0)      ; => 1
(inc 3.14)   ; => 4.14

;; Works with big numbers too
(inc 999999999999999999999)  ; => 1000000000000000000000
```

#### `dec` - Decrement a number
```lisp
(dec 5)      ; => 4
(dec -3)     ; => -4
(dec 0)      ; => -1
(dec 3.14)   ; => 2.14

;; Works with big numbers too
(dec 1000000000000000000000)  ; => 999999999999999999999
```

#### Practical examples
```lisp
;; Counter pattern
(def counter 0)
(def counter (inc counter))  ; => 1
(def counter (inc counter))  ; => 2

;; Loop with increment
(loop [i 0]
  (if (< i 5)
    (do
      (println i)
      (recur (inc i)))
    "done"))
```

## List Functions

### New Clojure-Style Functions

#### `second` - Get the second element
```lisp
(second '(a b c d))     ; => b
(second [1 2 3])        ; => 2
(second '(only-one))    ; => nil
(second '())            ; => nil
```

#### `concat` - Concatenate sequences (alias for `append`)
```lisp
(concat '(1 2) '(3 4))           ; => (1 2 3 4)
(concat '(a) '(b c) '(d e))      ; => (a b c d e)
(concat '() '(1 2 3))            ; => (1 2 3)
(concat)                         ; => ()
```

## Polymorphic Count Function

The `count` function has been made polymorphic and moved to the core plugin, making it work with all collection types just like in Clojure.

### `count` - Get count of elements in any collection
```lisp
;; Lists
(count '(1 2 3 4))        ; => 4
(count '())               ; => 0

;; Vectors
(count [a b c])           ; => 3
(count [])                ; => 0

;; Hash maps (counts key-value pairs)
(count {:a 1 :b 2 :c 3})  ; => 3
(count {})                ; => 0

;; Strings (counts characters, not bytes)
(count "hello")           ; => 5
(count "")                ; => 0
(count "héllo")           ; => 5  ; Unicode-aware

;; Nil
(count nil)               ; => 0
```

## String Functions

### Clojure-Style String Functions

| Function | Description | Example |
|----------|-------------|---------|
| `str` | String concatenation | `(str "Hello" " " "World")` → `"Hello World"` |
| `subs` | Substring extraction | `(subs "Hello" 1 4)` → `"ell"` |
| `split` | Split string | `(split "a,b,c" ",")` → `["a" "b" "c"]` |
| `join` | Join strings | `(join "," ["a" "b" "c"])` → `"a,b,c"` |
| `replace` | Replace substring | `(replace "hello" "l" "x")` → `"hexxo"` |
| `trim` | Trim whitespace | `(trim "  hello  ")` → `"hello"` |
| `upper-case` | Convert to uppercase | `(upper-case "hello")` → `"HELLO"` |
| `lower-case` | Convert to lowercase | `(lower-case "HELLO")` → `"hello"` |

### String Function Examples
```lisp
;; Building strings
(str "User: " "Alice" " (ID: " 123 ")")  ; => "User: Alice (ID: 123)"

;; String manipulation pipeline
(-> "  Hello, World!  "
    trim
    lower-case
    (replace "world" "clojure"))  ; => "hello, clojure!"

;; Working with collections of strings
(join " | " (map upper-case ["apple" "banana" "cherry"]))
; => "APPLE | BANANA | CHERRY"
```

## Migration Guide

### From Go-Lisp style to Modern Polymorphic style

```lisp
;; Old style (type-specific functions)
(list-first my-list)
(list-rest my-list)
(list-length my-list)
(hash-map-get my-map "key")
(hash-map-size my-map)

;; New Polymorphic style (works on all types)
(first my-list)        ; Works on lists, vectors, strings
(rest my-list)         ; Works on lists, vectors, strings
(count my-list)        ; Works on all collections
(get my-map "key")     ; Works on hashmaps, vectors, lists, strings
(count my-map)         ; Works on all collections

;; Cross-type compatibility
(first "hello")        ; => "h"
(rest [1 2 3])         ; => (2 3)
(last {:a 1 :b 2})     ; Not applicable, but count works
(count "hello")        ; => 5
(get [1 2 3] 1)        ; => 2
(contains? "hello" "ell") ; => true
```

## Best Practices

1. **Use polymorphic functions** for maximum code reusability and Clojure compatibility
2. **Prefer `first`, `rest`, `last`, `nth`** over type-specific alternatives
3. **Use `count`** over type-specific length functions for all collections
4. **Use polymorphic `get` and `contains?`** for unified data access patterns
5. **Use `assoc`/`dissoc`** for hash map operations instead of put/remove
6. **Use `inc`/`dec`** for simple increment/decrement operations
7. **Use keywords** (`:key`) as hash map keys when possible
8. **Leverage type predicates** (`seq?`, `coll?`, etc.) for type-aware code

## Real-World Examples

### Data Processing Pipeline
```lisp
;; Works with any sequence type
(defn process-collection [coll]
  (if (empty? coll)
    "No data"
    (str "First: " (first coll) 
         ", Last: " (last coll)
         ", Count: " (count coll))))

(process-collection '(1 2 3))      ; => "First: 1, Last: 3, Count: 3"
(process-collection [1 2 3])       ; => "First: 1, Last: 3, Count: 3"  
(process-collection "hello")       ; => "First: h, Last: o, Count: 5"
```

### Unified Data Access
```lisp
;; Single function works on all data structures
(defn safe-get [coll key-or-index]
  (if (contains? coll key-or-index)
    (get coll key-or-index)
    "Not found"))

(safe-get {:name "Alice"} :name)   ; => "Alice"
(safe-get [1 2 3] 1)               ; => 2
(safe-get "hello" 1)               ; => "e"
(safe-get [1 2] 5)                 ; => "Not found"
```

## Performance Notes

- All polymorphic functions maintain excellent performance through efficient type dispatch
- The polymorphic functions have minimal overhead for type checking
- Functions preserve input types where possible (e.g., `reverse` on vectors returns vectors)
- `rest` always returns a list for consistency with Clojure behavior
- `inc` and `dec` functions handle both regular and big numbers efficiently
- String functions use Unicode-aware operations where appropriate
- Cross-type operations are optimized for the most common use cases
