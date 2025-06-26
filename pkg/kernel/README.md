# Lisp Kernel

Minimal Lisp interpreter with macro system and self-hosting capabilities.

## Features

- **Core data types**: symbols, lists, vectors, hash-maps
- **Lexical environments** with parent chain lookup
- **Special forms**: `if`, `def`, `fn`, `quote`, `do`, `defmacro`
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
; Define a factorial function
(def fact (fn [n] 
  (if (= n 0) 
    1 
    (* n (fact (- n 1))))))

; Create a simple conditional macro
(defmacro when [test body]
  `(if ~test ~body nil))

; Data types examples
(def my-list '(1 2 3))
(def my-vector [1 2 3])
(def my-map (hash-map :name "GoLisp" :version 1))

; Use them
(fact 5)  ; => 120
(when (> 10 5) "true!")  ; => "true!"
(first my-list)  ; => 1
(:name my-map)  ; => "GoLisp"
```