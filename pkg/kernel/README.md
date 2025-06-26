# Lisp Kernel

Minimal Lisp interpreter with macro system and self-hosting capabilities.

## Features

- **Core data types**: symbols, lists, vectors, hash-maps
- **Lexical environments** with parent chain lookup
- **Special forms**: `if`, `def`, `fn`, `quote`, `do`, `defmacro`, `loop`, `recur`
- **Tail call optimization** with loop/recur for efficient iteration
- **Macro system** with quasiquote/unquote templates
- **Self-hosting**: higher-level features built in Lisp itself
- **Interactive REPL** for development

## Files

- `types.go` - Data types and structures
- `env.go` - Environment management
- `eval.go` - Core evaluator (~475 lines)
- `bootstrap.go` - Built-in functions
- `repl.go` - Interactive shell

## Example

```lisp
; Define a factorial function using recursion
(def fact (fn [n] 
  (if (= n 0) 
    1 
    (* n (fact (- n 1))))))

; Define factorial using loop/recur for tail call optimization
(def fact-loop (fn [n]
  (loop [i n acc 1]
    (if (= i 0)
      acc
      (recur (- i 1) (* acc i))))))

; Create a simple conditional macro
(defmacro when [test body]
  `(if ~test ~body nil))

; Data types examples
(def my-list '(1 2 3))
(def my-vector [1 2 3])
(def my-map (hash-map :name "GoLisp" :version 1))

; Use them
(fact 5)  ; => 120
(fact-loop 5)  ; => 120 (tail-call optimized)
(when (> 10 5) "true!")  ; => "true!"
(first my-list)  ; => 1
(:name my-map)  ; => "GoLisp"

; Loop examples
(loop [x 3] (if (= x 0) x (recur (- x 1))))  ; => 0
(loop [i 5 sum 0] (if (= i 0) sum (recur (- i 1) (+ sum i))))  ; => 15

; Built-in collection functions
(range 5)        ; => (0 1 2 3 4)
(range 2 8)      ; => (2 3 4 5 6 7)
(range 0 10 2)   ; => (0 2 4 6 8)

; Reduce function for collection processing
(reduce + (range 5))           ; => 10 (sum of 0+1+2+3+4)
(reduce * (range 1 6))         ; => 120 (factorial: 1*2*3*4*5)
(reduce + (range 5) 100)       ; => 110 (sum with initial value)
```