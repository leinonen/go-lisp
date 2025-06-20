# GoLisp Function Reference

This is a comprehensive reference of all functions available in GoLisp, organized by category.

## Core Functions

### `def`
**Definition:** Define a variable: `(def name value)`
**Example:**
```lisp
(def x 42)
(def name "John")
```

### `fn`
**Definition:** Create a function: `(fn [params] body)`
**Example:**
```lisp
(def add (fn [x y] (+ x y)))
(add 3 4) ; => 7
```

### `defn`
**Definition:** Define a named function: `(defn name [params] body)`
**Example:**
```lisp
(defn square [x] (* x x))
(square 5) ; => 25
```

### `quote`
**Definition:** Quote an expression: `(quote expr)`
**Example:**
```lisp
(quote (+ 1 2)) ; => (+ 1 2)
'(a b c) ; => (a b c)
```

### `help`
**Definition:** Show help: `(help)` for all functions, `(help function-name)` for specific function
**Example:**
```lisp
(help)
(help +)
```

### `env`
**Definition:** Show environment variables: `(env)`
**Example:**
```lisp
(env) ; Shows all variables and functions
```

### `plugins`
**Definition:** Show loaded plugins: `(plugins)`
**Example:**
```lisp
(plugins) ; Shows all loaded plugin categories
```

### `count`
**Definition:** Get count of elements in collection: `(count coll)` - works on lists, vectors, hash-maps, strings
**Examples:**
```lisp
(count '(1 2 3)) ; => 3
(count [1 2 3 4]) ; => 4
(count {:a 1 :b 2}) ; => 2
(count "hello") ; => 5
```

## Arithmetic Functions

### `+`
**Definition:** Add numbers: `(+ 1 2 3)` => 6
**Example:**
```lisp
(+ 1 2 3) ; => 6
(+ 3.14 2.86) ; => 6.0
```

### `-`
**Definition:** Subtract numbers: `(- 10 3 2)` => 5, `(- 5)` => -5
**Examples:**
```lisp
(- 10 3 2) ; => 5
(- 5) ; => -5 (negation)
```

### `*`
**Definition:** Multiply numbers: `(* 2 3 4)` => 24
**Example:**
```lisp
(* 2 3 4) ; => 24
(* 3.14 2) ; => 6.28
```

### `/`
**Definition:** Divide numbers: `(/ 12 3 2)` => 2
**Example:**
```lisp
(/ 12 3 2) ; => 2
(/ 10 3) ; => 3.333...
```

### `%`
**Definition:** Modulo operation: `(% 10 3)` => 1
**Example:**
```lisp
(% 10 3) ; => 1
(% 7 2) ; => 1
```

### `inc`
**Definition:** Increment by 1: `(inc 5)` => 6
**Example:**
```lisp
(inc 5) ; => 6
(inc -1) ; => 0
```

### `dec`
**Definition:** Decrement by 1: `(dec 5)` => 4
**Example:**
```lisp
(dec 5) ; => 4
(dec 1) ; => 0
```

## Comparison Functions

### `=`
**Definition:** Test equality: `(= 1 1 1)` => true, `(= 1 2)` => false
**Examples:**
```lisp
(= 1 1) ; => true
(= 1 2) ; => false
(= "hello" "hello") ; => true
```

### `<`
**Definition:** Test less than: `(< 1 2 3)` => true, `(< 1 3 2)` => false
**Example:**
```lisp
(< 1 2 3) ; => true
(< 5 3) ; => false
```

### `>`
**Definition:** Test greater than: `(> 3 2 1)` => true, `(> 3 1 2)` => false
**Example:**
```lisp
(> 3 2 1) ; => true
(> 1 2) ; => false
```

### `<=`
**Definition:** Test less than or equal: `(<= 1 2 2)` => true, `(<= 2 1)` => false
**Example:**
```lisp
(<= 1 2 2) ; => true
(<= 5 5) ; => true
```

### `>=`
**Definition:** Test greater than or equal: `(>= 3 2 2)` => true, `(>= 2 3)` => false
**Example:**
```lisp
(>= 3 2 2) ; => true
(>= 5 5) ; => true
```

## Logical Functions

### `and`
**Definition:** Logical AND: `(and true true false)` => false, `(and true true)` => true
**Example:**
```lisp
(and true true) ; => true
(and true false) ; => false
(and 1 2 3) ; => true (all truthy)
```

### `or`
**Definition:** Logical OR: `(or false false true)` => true, `(or false false)` => false
**Example:**
```lisp
(or false true) ; => true
(or false false) ; => false
(or 0 1) ; => true
```

### `not`
**Definition:** Logical NOT: `(not true)` => false, `(not false)` => true
**Example:**
```lisp
(not true) ; => false
(not false) ; => true
(not 0) ; => true (0 is falsy)
```

## Math Functions

### `sqrt`
**Definition:** Square root: `(sqrt 16)` => 4
**Example:**
```lisp
(sqrt 16) ; => 4
(sqrt 2) ; => 1.414...
```

### `pow`
**Definition:** Power: `(pow 2 3)` => 8
**Example:**
```lisp
(pow 2 3) ; => 8
(pow 10 2) ; => 100
```

### `abs`
**Definition:** Absolute value: `(abs -7)` => 7
**Example:**
```lisp
(abs -7) ; => 7
(abs 7) ; => 7
```

### `floor`
**Definition:** Floor (largest integer ≤ x): `(floor 3.7)` => 3
**Example:**
```lisp
(floor 3.7) ; => 3
(floor -3.7) ; => -4
```

### `ceil`
**Definition:** Ceiling (smallest integer ≥ x): `(ceil 3.2)` => 4
**Example:**
```lisp
(ceil 3.2) ; => 4
(ceil -3.2) ; => -3
```

### `round`
**Definition:** Round to nearest integer: `(round 3.6)` => 4
**Example:**
```lisp
(round 3.6) ; => 4
(round 3.4) ; => 3
```

### `trunc`
**Definition:** Truncate towards zero: `(trunc 3.7)` => 3
**Example:**
```lisp
(trunc 3.7) ; => 3
(trunc -3.7) ; => -3
```

### Trigonometric Functions

### `sin`
**Definition:** Sine: `(sin 0)` => 0
**Example:**
```lisp
(sin 0) ; => 0
(sin (/ pi 2)) ; => 1
```

### `cos`
**Definition:** Cosine: `(cos 0)` => 1
**Example:**
```lisp
(cos 0) ; => 1
(cos pi) ; => -1
```

### `tan`
**Definition:** Tangent: `(tan 0)` => 0
**Example:**
```lisp
(tan 0) ; => 0
(tan (/ pi 4)) ; => 1
```

### `asin`
**Definition:** Arc sine: `(asin 0.5)` => 0.5236 (π/6)
**Example:**
```lisp
(asin 0.5) ; => 0.5236
(asin 1) ; => 1.5708 (π/2)
```

### `acos`
**Definition:** Arc cosine: `(acos 0.5)` => 1.0472 (π/3)
**Example:**
```lisp
(acos 0.5) ; => 1.0472
(acos 1) ; => 0
```

### `atan`
**Definition:** Arc tangent: `(atan 1)` => 0.7854 (π/4)
**Example:**
```lisp
(atan 1) ; => 0.7854
(atan 0) ; => 0
```

### `atan2`
**Definition:** Two-argument arc tangent: `(atan2 1 1)` => 0.7854 (π/4)
**Example:**
```lisp
(atan2 1 1) ; => 0.7854
(atan2 1 0) ; => 1.5708
```

### Hyperbolic Functions

### `sinh`
**Definition:** Hyperbolic sine: `(sinh 0)` => 0
**Example:**
```lisp
(sinh 0) ; => 0
(sinh 1) ; => 1.1752
```

### `cosh`
**Definition:** Hyperbolic cosine: `(cosh 0)` => 1
**Example:**
```lisp
(cosh 0) ; => 1
(cosh 1) ; => 1.5431
```

### `tanh`
**Definition:** Hyperbolic tangent: `(tanh 0)` => 0
**Example:**
```lisp
(tanh 0) ; => 0
(tanh 1) ; => 0.7616
```

### Logarithmic and Exponential Functions

### `log`
**Definition:** Natural logarithm: `(log (e))` => 1
**Example:**
```lisp
(log (e)) ; => 1
(log 10) ; => 2.3026
```

### `exp`
**Definition:** Exponential (e^x): `(exp 1)` => 2.7183
**Example:**
```lisp
(exp 1) ; => 2.7183
(exp 0) ; => 1
```

### `log10`
**Definition:** Base-10 logarithm: `(log10 100)` => 2
**Example:**
```lisp
(log10 100) ; => 2
(log10 10) ; => 1
```

### `log2`
**Definition:** Base-2 logarithm: `(log2 8)` => 3
**Example:**
```lisp
(log2 8) ; => 3
(log2 2) ; => 1
```

### Angle Conversion

### `degrees`
**Definition:** Convert radians to degrees: `(degrees (pi))` => 180
**Example:**
```lisp
(degrees pi) ; => 180
(degrees (/ pi 2)) ; => 90
```

### `radians`
**Definition:** Convert degrees to radians: `(radians 180)` => 3.1416
**Example:**
```lisp
(radians 180) ; => 3.1416
(radians 90) ; => 1.5708
```

### Statistical Functions

### `min`
**Definition:** Minimum of numbers: `(min 5 2 8 1)` => 1
**Example:**
```lisp
(min 5 2 8 1) ; => 1
(min -3 -1 -5) ; => -5
```

### `max`
**Definition:** Maximum of numbers: `(max 5 2 8 1)` => 8
**Example:**
```lisp
(max 5 2 8 1) ; => 8
(max -3 -1 -5) ; => -1
```

### Utility Functions

### `sign`
**Definition:** Sign of number (-1, 0, or 1): `(sign -5)` => -1
**Example:**
```lisp
(sign -5) ; => -1
(sign 0) ; => 0
(sign 5) ; => 1
```

### `mod`
**Definition:** Mathematical modulo: `(mod 7 3)` => 1
**Example:**
```lisp
(mod 7 3) ; => 1
(mod -7 3) ; => 2
```

### Mathematical Constants

### `pi`
**Definition:** Pi constant: `(pi)` => 3.1416
**Example:**
```lisp
(pi) ; => 3.1416
(* 2 (pi)) ; => 6.2832
```

### `e`
**Definition:** Euler's number: `(e)` => 2.7183
**Example:**
```lisp
(e) ; => 2.7183
(exp 1) ; => 2.7183
```

### `random`
**Definition:** Random number: `(random)` => 0-1, `(random 10)` => 0-9, `(random 5 15)` => 5-14
**Examples:**
```lisp
(random) ; => 0.7234 (random float 0-1)
(random 10) ; => 7 (random int 0-9)
(random 5 15) ; => 12 (random int 5-14)
```

## Control Flow Functions

### `if`
**Definition:** Conditional: `(if condition then-expr else-expr?)` => evaluated result
**Examples:**
```lisp
(if true "yes" "no") ; => "yes"
(if false "yes" "no") ; => "no"
(if (> 5 3) "bigger") ; => "bigger"
```

### `do`
**Definition:** Sequential evaluation: `(do expr1 expr2 expr3)` => result of last expr
**Example:**
```lisp
(do 
  (def x 5)
  (def y 10)
  (+ x y)) ; => 15
```

### `cond`
**Definition:** Multi-branch conditional: `(cond test1 expr1 test2 expr2 :else default)`
**Example:**
```lisp
(cond 
  (< x 0) "negative"
  (= x 0) "zero"
  :else "positive")
```

### `when`
**Definition:** Conditional execution: `(when test expr1 expr2 ...)`
**Example:**
```lisp
(when (> x 0)
  (print "positive")
  (* x 2))
```

### `when-not`
**Definition:** Negated conditional execution: `(when-not test expr1 expr2 ...)`
**Example:**
```lisp
(when-not (= x 0)
  (print "not zero")
  (/ 1 x))
```

### `loop`
**Definition:** Loop establishes a recursion point with local bindings: `(loop [var1 init1 var2 init2 ...] body...)`
**Example:**
```lisp
(loop [i 0 sum 0]
  (if (< i 5)
    (recur (inc i) (+ sum i))
    sum)) ; => 10
```

### `recur`
**Definition:** Recur jumps back to the nearest loop with new values
**Example:** See `loop` example above.

## List Functions

### `list`
**Definition:** Create a list: `(list 1 2 3)` => (1 2 3)
**Example:**
```lisp
(list 1 2 3) ; => (1 2 3)
(list "a" "b" "c") ; => ("a" "b" "c")
```

### `cons`
**Definition:** Prepend element: `(cons 0 '(1 2))` => (0 1 2)
**Example:**
```lisp
(cons 0 '(1 2)) ; => (0 1 2)
(cons "x" '("y" "z")) ; => ("x" "y" "z")
```

### `length`
**Definition:** Get list length: `(length '(1 2 3))` => 3
**Example:**
```lisp
(length '(1 2 3)) ; => 3
(length '()) ; => 0
```

### `append`
**Definition:** Append lists: `(append '(1 2) '(3 4))` => (1 2 3 4)
**Example:**
```lisp
(append '(1 2) '(3 4)) ; => (1 2 3 4)
(append '(a) '(b) '(c)) ; => (a b c)
```

### `concat`
**Definition:** Concatenate lists: `(concat '(1 2) '(3 4))` => (1 2 3 4)
**Example:**
```lisp
(concat '(1 2) '(3 4)) ; => (1 2 3 4)
```

## Vector Functions

### `vector`
**Definition:** Create a vector: `(vector 1 2 3)` => [1 2 3]
**Example:**
```lisp
(vector 1 2 3) ; => [1 2 3]
(vector "a" "b") ; => ["a" "b"]
```

### `vec`
**Definition:** Convert to vector: `(vec '(1 2 3))` => [1 2 3]
**Example:**
```lisp
(vec '(1 2 3)) ; => [1 2 3]
(vec "abc") ; => ["a" "b" "c"]
```

### `vector?`
**Definition:** Check if value is a vector: `(vector? [1 2 3])` => true
**Examples:**
```lisp
(vector? [1 2 3]) ; => true
(vector? '(1 2 3)) ; => false
```

### `conj`
**Definition:** Add elements to collection: `(conj [1 2] 3 4)` => [1 2 3 4]
**Examples:**
```lisp
(conj [1 2] 3 4) ; => [1 2 3 4] (vectors append)
(conj '(1 2) 3 4) ; => (4 3 1 2) (lists prepend)
```

## Polymorphic Sequence Functions

These functions work on multiple collection types (lists, vectors, strings):

### `first`
**Definition:** Get first element of sequence: `(first coll)` - works on lists, vectors, strings
**Examples:**
```lisp
(first '(1 2 3)) ; => 1
(first [1 2 3]) ; => 1
(first "hello") ; => "h"
```

### `rest`
**Definition:** Get rest of sequence: `(rest coll)` - works on lists, vectors, strings
**Examples:**
```lisp
(rest '(1 2 3)) ; => (2 3)
(rest [1 2 3]) ; => (2 3)
(rest "hello") ; => ("e" "l" "l" "o")
```

### `last`
**Definition:** Get last element of sequence: `(last coll)` - works on lists, vectors, strings
**Examples:**
```lisp
(last '(1 2 3)) ; => 3
(last [1 2 3]) ; => 3
(last "hello") ; => "o"
```

### `nth`
**Definition:** Get nth element of sequence: `(nth coll n)` - works on lists, vectors, strings
**Examples:**
```lisp
(nth '(1 2 3) 1) ; => 2
(nth [1 2 3] 0) ; => 1
(nth "hello" 1) ; => "e"
```

### `second`
**Definition:** Get second element of sequence: `(second coll)` - works on lists, vectors, strings
**Examples:**
```lisp
(second '(1 2 3)) ; => 2
(second [1 2 3]) ; => 2
(second "hello") ; => "e"
```

### `empty?`
**Definition:** Check if collection is empty: `(empty? coll)` - works on all collections
**Examples:**
```lisp
(empty? '()) ; => true
(empty? []) ; => true
(empty? "") ; => true
(empty? '(1)) ; => false
```

### `seq`
**Definition:** Convert to sequence: `(seq coll)` - works on lists, vectors, strings
**Examples:**
```lisp
(seq [1 2 3]) ; => (1 2 3)
(seq "abc") ; => ("a" "b" "c")
```

### `take`
**Definition:** Take first n elements: `(take n coll)` - works on all sequences
**Examples:**
```lisp
(take 2 '(1 2 3 4)) ; => (1 2)
(take 3 [1 2 3 4 5]) ; => (1 2 3)
(take 2 "hello") ; => ("h" "e")
```

### `drop`
**Definition:** Drop first n elements: `(drop n coll)` - works on all sequences
**Examples:**
```lisp
(drop 2 '(1 2 3 4)) ; => (3 4)
(drop 1 [1 2 3]) ; => (2 3)
(drop 2 "hello") ; => ("l" "l" "o")
```

### `reverse`
**Definition:** Reverse sequence: `(reverse coll)` - works on lists, vectors, strings
**Examples:**
```lisp
(reverse '(1 2 3)) ; => (3 2 1)
(reverse [1 2 3]) ; => (3 2 1)
(reverse "hello") ; => ("o" "l" "l" "e" "h")
```

### `distinct`
**Definition:** Remove duplicates: `(distinct coll)` - works on all sequences
**Examples:**
```lisp
(distinct '(1 2 2 3 3 3)) ; => (1 2 3)
(distinct ["a" "b" "a" "c" "b"]) ; => ("a" "b" "c")
(distinct "hello") ; => ("h" "e" "l" "o")
```

### `sort`
**Definition:** Sort sequence: `(sort coll)` - works on all sequences
**Examples:**
```lisp
(sort '(3 1 4 1 5)) ; => (1 1 3 4 5)
(sort ["c" "a" "b"]) ; => ("a" "b" "c")
(sort "hello") ; => ("e" "h" "l" "l" "o")
```

### `into`
**Definition:** Merge collections: `(into to from)` - works on all collections
**Examples:**
```lisp
(into [] '(1 2 3)) ; => [1 2 3]
(into '() [1 2 3]) ; => (3 2 1)
```

## Predicate Functions

### `seq?`
**Definition:** Check if sequential: `(seq? x)`
**Examples:**
```lisp
(seq? '(1 2 3)) ; => true
(seq? [1 2 3]) ; => true
(seq? "abc") ; => true
(seq? {:a 1}) ; => false
```

### `coll?`
**Definition:** Check if collection: `(coll? x)`
**Examples:**
```lisp
(coll? '(1 2 3)) ; => true
(coll? [1 2 3]) ; => true
(coll? {:a 1}) ; => true
(coll? "abc") ; => true
```

### `sequential?`
**Definition:** Check if sequential: `(sequential? x)`
**Examples:**
```lisp
(sequential? '(1 2 3)) ; => true
(sequential? [1 2 3]) ; => true
(sequential? {:a 1}) ; => false
```

### `indexed?`
**Definition:** Check if indexed: `(indexed? x)`
**Examples:**
```lisp
(indexed? [1 2 3]) ; => true
(indexed? "abc") ; => true
(indexed? '(1 2 3)) ; => false
```

## Utility Functions

### `identity`
**Definition:** Return argument unchanged: `(identity x)`
**Example:**
```lisp
(identity 42) ; => 42
(identity "hello") ; => "hello"
```

### `constantly`
**Definition:** Return constant function: `(constantly x)`
**Example:**
```lisp
(def always-5 (constantly 5))
(always-5 1 2 3) ; => 5
```

## Functional Programming Functions

### `map`
**Definition:** Apply a function to each element of a list: `(map fn list)`
**Example:**
```lisp
(map inc '(1 2 3)) ; => (2 3 4)
(map (fn [x] (* x x)) [1 2 3 4]) ; => (1 4 9 16)
```

### `filter`
**Definition:** Filter elements of a list using a predicate: `(filter pred list)`
**Example:**
```lisp
(filter (fn [x] (> x 2)) '(1 2 3 4 5)) ; => (3 4 5)
(filter even? [1 2 3 4 5 6]) ; => (2 4 6)
```

### `reduce`
**Definition:** Reduce a list to a single value: `(reduce fn init list)`
**Example:**
```lisp
(reduce + 0 '(1 2 3 4)) ; => 10
(reduce * 1 [1 2 3 4]) ; => 24
```

### `apply`
**Definition:** Apply a function to a list of arguments: `(apply fn args)`
**Example:**
```lisp
(apply + '(1 2 3 4)) ; => 10
(apply max [5 2 8 1]) ; => 8
```

## Hash Map Functions

### `hash-map`
**Definition:** Create a hash map: `(hash-map :key1 val1 :key2 val2)`
**Example:**
```lisp
(hash-map :name "John" :age 30) ; => {:name "John" :age 30}
```

### `get`
**Definition:** Get value from hash map: `(get map key)` or `(get map key default)`
**Examples:**
```lisp
(get {:a 1 :b 2} :a) ; => 1
(get {:a 1 :b 2} :c 0) ; => 0 (default value)
```

### `assoc`
**Definition:** Associate key-value pairs: `(assoc map key val)`
**Example:**
```lisp
(assoc {:a 1} :b 2) ; => {:a 1 :b 2}
(assoc {} :name "John" :age 30) ; => {:name "John" :age 30}
```

### `dissoc`
**Definition:** Dissociate keys: `(dissoc map key1 key2 ...)`
**Example:**
```lisp
(dissoc {:a 1 :b 2 :c 3} :b) ; => {:a 1 :c 3}
(dissoc {:a 1 :b 2 :c 3} :a :c) ; => {:b 2}
```

### `contains?`
**Definition:** Check if hash map contains key: `(contains? map key)`
**Examples:**
```lisp
(contains? {:a 1 :b 2} :a) ; => true
(contains? {:a 1 :b 2} :c) ; => false
```

### `keys`
**Definition:** Get all keys: `(keys map)`
**Example:**
```lisp
(keys {:a 1 :b 2 :c 3}) ; => (:a :b :c)
```

### `vals`
**Definition:** Get all values: `(vals map)`
**Example:**
```lisp
(vals {:a 1 :b 2 :c 3}) ; => (1 2 3)
```

## String Functions

### `str` (alias: `string-concat`)
**Definition:** Concatenate strings: `(str "Hello" " " "World")`
**Example:**
```lisp
(str "Hello" " " "World") ; => "Hello World"
(str "Number: " 42) ; => "Number: 42"
```

### `string-length`
**Definition:** Get string length: `(string-length "hello")` => 5
**Example:**
```lisp
(string-length "hello") ; => 5
(string-length "") ; => 0
```

### `subs` (alias: `string-substring`)
**Definition:** Get substring: `(subs "hello" 1 3)` => "el"
**Example:**
```lisp
(subs "hello" 1 3) ; => "el"
(subs "hello" 2) ; => "llo"
```

### `string-char-at`
**Definition:** Get character at index: `(string-char-at "hello" 1)` => "e"
**Example:**
```lisp
(string-char-at "hello" 1) ; => "e"
(string-char-at "hello" 0) ; => "h"
```

### `string-upper`
**Definition:** Convert to uppercase: `(string-upper "hello")` => "HELLO"
**Example:**
```lisp
(string-upper "hello") ; => "HELLO"
(string-upper "Hello World") ; => "HELLO WORLD"
```

### `string-lower`
**Definition:** Convert to lowercase: `(string-lower "HELLO")` => "hello"
**Example:**
```lisp
(string-lower "HELLO") ; => "hello"
(string-lower "Hello World") ; => "hello world"
```

### `string-trim`
**Definition:** Trim whitespace: `(string-trim " hello ")` => "hello"
**Example:**
```lisp
(string-trim " hello ") ; => "hello"
(string-trim "\n\ttest\n") ; => "test"
```

### `string-split`
**Definition:** Split string: `(string-split "a,b,c" ",")` => ("a" "b" "c")
**Example:**
```lisp
(string-split "a,b,c" ",") ; => ("a" "b" "c")
(string-split "hello world" " ") ; => ("hello" "world")
```

### `string-join`
**Definition:** Join strings: `(string-join ["a" "b" "c"] ",")` => "a,b,c"
**Example:**
```lisp
(string-join ["a" "b" "c"] ",") ; => "a,b,c"
(string-join '("hello" "world") " ") ; => "hello world"
```

### `string-contains?`
**Definition:** Check if string contains substring: `(string-contains? "hello" "ell")` => true
**Example:**
```lisp
(string-contains? "hello" "ell") ; => true
(string-contains? "hello" "xyz") ; => false
```

### `string-starts-with?`
**Definition:** Check if string starts with prefix: `(string-starts-with? "hello" "he")` => true
**Example:**
```lisp
(string-starts-with? "hello" "he") ; => true
(string-starts-with? "hello" "lo") ; => false
```

### `string-ends-with?`
**Definition:** Check if string ends with suffix: `(string-ends-with? "hello" "lo")` => true
**Example:**
```lisp
(string-ends-with? "hello" "lo") ; => true
(string-ends-with? "hello" "he") ; => false
```

### `string-replace`
**Definition:** Replace all occurrences: `(string-replace "hello" "l" "x")` => "hexxo"
**Example:**
```lisp
(string-replace "hello" "l" "x") ; => "hexxo"
(string-replace "hello world" "o" "0") ; => "hell0 w0rld"
```

### `string-index-of`
**Definition:** Find index of substring: `(string-index-of "hello" "ll")` => 2
**Example:**
```lisp
(string-index-of "hello" "ll") ; => 2
(string-index-of "hello" "x") ; => -1 (not found)
```

### `string->number`
**Definition:** Convert string to number: `(string->number "42")` => 42
**Example:**
```lisp
(string->number "42") ; => 42
(string->number "3.14") ; => 3.14
```

### `number->string`
**Definition:** Convert number to string: `(number->string 42)` => "42"
**Example:**
```lisp
(number->string 42) ; => "42"
(number->string 3.14) ; => "3.14"
```

### `string-regex-match?`
**Definition:** Check if string matches regex: `(string-regex-match? "hello123" "[0-9]+")` => true
**Example:**
```lisp
(string-regex-match? "hello123" "[0-9]+") ; => true
(string-regex-match? "hello" "[0-9]+") ; => false
```

### `string-regex-find-all`
**Definition:** Find all regex matches: `(string-regex-find-all "abc123def456" "[0-9]+")` => ("123" "456")
**Example:**
```lisp
(string-regex-find-all "abc123def456" "[0-9]+") ; => ("123" "456")
```

### `string-repeat`
**Definition:** Repeat string: `(string-repeat "Hi" 3)` => "HiHiHi"
**Example:**
```lisp
(string-repeat "Hi" 3) ; => "HiHiHi"
(string-repeat "-" 5) ; => "-----"
```

### `string?`
**Definition:** Check if value is string: `(string? "hello")` => true
**Examples:**
```lisp
(string? "hello") ; => true
(string? 42) ; => false
```

### `string-empty?`
**Definition:** Check if string is empty: `(string-empty? "")` => true
**Examples:**
```lisp
(string-empty? "") ; => true
(string-empty? "hello") ; => false
```

## I/O Functions

### `print`
**Definition:** Print values: `(print "Hello" "World")`
**Example:**
```lisp
(print "Hello World") ; Prints: Hello World
(print "Value:" 42) ; Prints: Value: 42
```

### `println`
**Definition:** Print values with newline: `(println "Hello" "World")`
**Example:**
```lisp
(println "Hello") ; Prints: Hello\n
(println "Line 1") ; Prints: Line 1\n
```

### `read-file`
**Definition:** Read file contents: `(read-file "filename.txt")`
**Example:**
```lisp
(read-file "data.txt") ; => "file contents"
```

### `write-file`
**Definition:** Write to file: `(write-file "filename.txt" "content")`
**Example:**
```lisp
(write-file "output.txt" "Hello World")
```

### `file-exists?`
**Definition:** Check if file exists: `(file-exists? "filename.txt")`
**Example:**
```lisp
(file-exists? "data.txt") ; => true or false
```

## JSON Functions

### `json-parse`
**Definition:** Parse JSON string: `(json-parse "{\"key\": \"value\"}")`
**Example:**
```lisp
(json-parse "{\"name\": \"John\", \"age\": 30}")
; => {:name "John" :age 30}
```

### `json-stringify`
**Definition:** Convert to JSON string: `(json-stringify {:key "value"})`
**Example:**
```lisp
(json-stringify {:name "John" :age 30})
; => "{\"name\":\"John\",\"age\":30}"
```

### `json-stringify-pretty`
**Definition:** Convert to pretty JSON: `(json-stringify-pretty {:key "value"})`
**Example:**
```lisp
(json-stringify-pretty {:name "John" :age 30})
; => "{\n  \"name\": \"John\",\n  \"age\": 30\n}"
```

### `json-path`
**Definition:** Extract value using JSON path: `(json-path data "$.key")`
**Example:**
```lisp
(json-path {:user {:name "John"}} "$.user.name") ; => "John"
```

## Keyword Functions

### `keyword`
**Definition:** Create a keyword: `(keyword "test")` => :test
**Example:**
```lisp
(keyword "test") ; => :test
(keyword "my-key") ; => :my-key
```

### `keyword?`
**Definition:** Check if value is keyword: `(keyword? :test)` => true
**Examples:**
```lisp
(keyword? :test) ; => true
(keyword? "test") ; => false
```

## Atom Functions (Concurrency)

### `atom`
**Definition:** Create an atom: `(atom initial-value)`
**Example:**
```lisp
(def counter (atom 0))
(def state (atom {:count 0}))
```

### `deref` (alias: `@`)
**Definition:** Dereference atom: `(deref atom)` or `@atom`
**Example:**
```lisp
@counter ; => 0
(deref state) ; => {:count 0}
```

### `swap!`
**Definition:** Swap atom value: `(swap! atom function & args)`
**Example:**
```lisp
(swap! counter inc) ; increments counter
(swap! counter + 5) ; adds 5 to counter
```

### `reset!`
**Definition:** Reset atom value: `(reset! atom new-value)`
**Example:**
```lisp
(reset! counter 10) ; sets counter to 10
(reset! state {:count 5}) ; resets state
```

## Binding Functions

### `let`
**Definition:** Local bindings: `(let [var1 val1 var2 val2] body)`
**Example:**
```lisp
(let [x 10 y 20]
  (+ x y)) ; => 30
```

### `let*`
**Definition:** Sequential bindings: `(let* [var1 val1 var2 (f var1)] body)`
**Example:**
```lisp
(let* [x 10 y (+ x 5)]
  (* x y)) ; => 150
```

### `letfn`
**Definition:** Local function bindings: `(letfn [(f [x] body)] body)`
**Example:**
```lisp
(letfn [(factorial [n]
          (if (<= n 1) 1 (* n (factorial (dec n)))))]
  (factorial 5)) ; => 120
```

## Advanced Binding Functions

### `let-destructure`
**Definition:** Destructuring let: `(let-destructure [[a b] [1 2]] (+ a b))`
**Example:**
```lisp
(let-destructure [[x y] [10 20]]
  (+ x y)) ; => 30
```

### `fn-destructure`
**Definition:** Destructuring function: `(fn-destructure [[a b]] (+ a b))`
**Example:**
```lisp
(def add-pair (fn-destructure [[x y]] (+ x y)))
(add-pair [3 4]) ; => 7
```

## Macro Functions

### `defmacro`
**Definition:** Define a macro: `(defmacro name [params] body)`
**Example:**
```lisp
(defmacro when-positive [x & body]
  `(when (> ~x 0) ~@body))
```

### `macroexpand`
**Definition:** Expand macro: `(macroexpand '(when-positive x (print "positive")))`
**Example:**
```lisp
(macroexpand '(when-positive 5 (print "yes")))
```

### `unquote`
**Definition:** Unquote in macro: Used within macro definitions
**Example:** See `defmacro` example above.

## Concurrency Functions

### `go`
**Definition:** Start goroutine: `(go expression)`
**Example:**
```lisp
(go (print "Hello from goroutine"))
```

### `go-wait`
**Definition:** Wait for goroutine: `(go-wait goroutine-id)`
**Example:**
```lisp
(def g (go (+ 1 2)))
(go-wait g) ; => 3
```

### `go-wait-all`
**Definition:** Wait for all goroutines: `(go-wait-all [id1 id2 ...])`
**Example:**
```lisp
(def g1 (go (* 2 3)))
(def g2 (go (+ 4 5)))
(go-wait-all [g1 g2]) ; => [6 9]
```

### `chan`
**Definition:** Create channel: `(chan)` or `(chan buffer-size)`
**Example:**
```lisp
(def ch (chan))
(def buffered-ch (chan 10))
```

### `chan-send!`
**Definition:** Send to channel: `(chan-send! channel value)`
**Example:**
```lisp
(chan-send! ch "hello")
```

### `chan-recv!`
**Definition:** Receive from channel: `(chan-recv! channel)`
**Example:**
```lisp
(chan-recv! ch) ; => "hello"
```

### `chan-try-recv!`
**Definition:** Try receive from channel: `(chan-try-recv! channel)`
**Example:**
```lisp
(chan-try-recv! ch) ; => value or nil if empty
```

### `chan-close!`
**Definition:** Close channel: `(chan-close! channel)`
**Example:**
```lisp
(chan-close! ch)
```

### `chan-closed?`
**Definition:** Check if channel is closed: `(chan-closed? channel)`
**Example:**
```lisp
(chan-closed? ch) ; => true or false
```

## Utility Functions (Collection Operations)

### `frequencies`
**Definition:** Count occurrences: `(frequencies coll)`
**Example:**
```lisp
(frequencies '(1 2 2 3 3 3)) ; => {1 1, 2 2, 3 3}
(frequencies "hello") ; => {"h" 1, "e" 1, "l" 2, "o" 1}
```

### `group-by`
**Definition:** Group by function result: `(group-by fn coll)`
**Example:**
```lisp
(group-by (fn [x] (mod x 2)) [1 2 3 4 5 6])
; => {0 [2 4 6], 1 [1 3 5]}
```

### `partition`
**Definition:** Partition into chunks: `(partition n coll)`
**Example:**
```lisp
(partition 2 [1 2 3 4 5 6]) ; => ((1 2) (3 4) (5 6))
(partition 3 "abcdef") ; => (("a" "b" "c") ("d" "e" "f"))
```

### `interleave`
**Definition:** Interleave sequences: `(interleave seq1 seq2 ...)`
**Example:**
```lisp
(interleave [1 2 3] [4 5 6]) ; => (1 4 2 5 3 6)
(interleave "abc" "def") ; => ("a" "d" "b" "e" "c" "f")
```

### `interpose`
**Definition:** Insert separator: `(interpose sep coll)`
**Example:**
```lisp
(interpose "," ["a" "b" "c"]) ; => ("a" "," "b" "," "c")
(interpose 0 [1 2 3]) ; => (1 0 2 0 3)
```

### `flatten`
**Definition:** Flatten nested sequences: `(flatten coll)`
**Example:**
```lisp
(flatten [[1 2] [3 4] [5]]) ; => (1 2 3 4 5)
(flatten [1 [2 [3 4]] 5]) ; => (1 2 3 4 5)
```

### `shuffle`
**Definition:** Randomize order: `(shuffle coll)`
**Example:**
```lisp
(shuffle [1 2 3 4 5]) ; => [3 1 5 2 4] (random order)
(shuffle "hello") ; => ("o" "l" "h" "l" "e") (random order)
```

### `remove`
**Definition:** Remove elements: `(remove pred coll)`
**Example:**
```lisp
(remove (fn [x] (< x 3)) [1 2 3 4 5]) ; => (3 4 5)
(remove string-empty? ["a" "" "b" ""]) ; => ("a" "b")
```

### `keep`
**Definition:** Keep non-nil results: `(keep fn coll)`
**Example:**
```lisp
(keep (fn [x] (when (> x 2) (* x x))) [1 2 3 4])
; => (9 16)
```

### `mapcat`
**Definition:** Map then concatenate: `(mapcat fn coll)`
**Example:**
```lisp
(mapcat (fn [x] [x (* x 2)]) [1 2 3])
; => (1 2 2 4 3 6)
```

### `take-while`
**Definition:** Take while predicate true: `(take-while pred coll)`
**Example:**
```lisp
(take-while (fn [x] (< x 5)) [1 2 3 4 5 6 7])
; => (1 2 3 4)
```

### `drop-while`
**Definition:** Drop while predicate true: `(drop-while pred coll)`
**Example:**
```lisp
(drop-while (fn [x] (< x 5)) [1 2 3 4 5 6 7])
; => (5 6 7)
```

### `split-at`
**Definition:** Split at index: `(split-at n coll)`
**Example:**
```lisp
(split-at 3 [1 2 3 4 5 6]) ; => [(1 2 3) (4 5 6)]
```

### `split-with`
**Definition:** Split with predicate: `(split-with pred coll)`
**Example:**
```lisp
(split-with (fn [x] (< x 5)) [1 2 3 4 5 6 7])
; => [(1 2 3 4) (5 6 7)]
```

## Function Utilities

### `comp`
**Definition:** Function composition: `(comp f g h)`
**Example:**
```lisp
(def add-then-square (comp (fn [x] (* x x)) (fn [x] (+ x 1))))
(add-then-square 3) ; => 16 ((3+1)^2)
```

### `partial`
**Definition:** Partial application: `(partial f arg1 arg2)`
**Example:**
```lisp
(def add5 (partial + 5))
(add5 3) ; => 8
```

### `complement`
**Definition:** Complement predicate: `(complement pred)`
**Example:**
```lisp
(def not-empty? (complement empty?))
(not-empty? [1 2 3]) ; => true
```

### `juxt`
**Definition:** Apply multiple functions: `(juxt f g h)`
**Example:**
```lisp
(def stats (juxt count first last))
(stats [1 2 3 4 5]) ; => [5 1 5]
```

## Set Operations

### `union`
**Definition:** Union of collections: `(union coll1 coll2 ...)`
**Example:**
```lisp
(union [1 2 3] [3 4 5]) ; => (1 2 3 4 5)
(union "abc" "bcd") ; => ("a" "b" "c" "d")
```

### `intersection`
**Definition:** Intersection of collections: `(intersection coll1 coll2 ...)`
**Example:**
```lisp
(intersection [1 2 3] [2 3 4]) ; => (2 3)
(intersection "abc" "bcd") ; => ("b" "c")
```

### `difference`
**Definition:** Difference of collections: `(difference coll1 coll2)`
**Example:**
```lisp
(difference [1 2 3 4] [2 4]) ; => (1 3)
(difference "abcd" "bd") ; => ("a" "c")
```

---

*This documentation covers all the built-in functions available in GoLisp. Functions are organized by category and include both their formal definitions and practical examples. Many functions are polymorphic, meaning they work across different data types (lists, vectors, strings, hash-maps) as shown in the examples.*
