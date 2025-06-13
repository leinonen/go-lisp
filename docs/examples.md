# Examples

For comprehensive, runnable examples, see the `examples/` directory. This document provides quick reference snippets.

## Basic Usage

```lisp
; Functions
(defn square [x] (* x x))
(square 5) ; => 25

; Lists and higher-order functions
(map square (list 1 2 3 4)) ; => (1 4 9 16)
(filter (fn [x] (> x 5)) (list 1 8 3 10)) ; => (8 10)

; Data structures
(def user (hash-map :name "Alice" :age 30))
(hash-map-get user :name) ; => "Alice"

; Sequential evaluation
(do 
  (def x 10)
  (def y 20)
  (+ x y)) ; => 30
```

## Advanced Features

```lisp
; Modules
(module math (export factorial))
(defn factorial [n] (if (<= n 1) 1 (* n (factorial (- n 1)))))

; Functional programming (load "library/functional.lisp")
(def add-then-double (comp (partial * 2) (partial + 1)))
(add-then-double 5) ; => 12

; Macros
(defmacro when [condition &rest body]
  `(if ~condition (do ~@body)))
```

For complete, working examples with detailed explanations, run:
```bash
./lisp examples/basic_features.lisp
./lisp examples/functional_library.lisp
./lisp examples/test_do.lisp                  # Sequential evaluation with do
./lisp examples/goroutines.lisp          # Concurrency with goroutines and channels
```
