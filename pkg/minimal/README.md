# Minimal Lisp Kernel Implementation

This package implements a minimal Lisp kernel following the clean architecture principles outlined in `future.md`.

## 1. ✅ Minimal Core Interpreter (Kernel)
Built with a tiny, trusted set of primitive operations:

- **Symbols and interning** (`types.go`) - Unique symbol representation
- **Lists and Vectors** (`types.go`) - Core data structures with Clojure-style `[param]` syntax
- **Environments** (`env.go`) - Lexical scope with parent chain lookup
- **Eval/apply logic** (`eval.go`) - Core evaluation engine
- **Special forms**: `if`, `define`, `fn`, `quote`, `do`, `quasiquote`, `unquote`, `defmacro`
- **Minimal REPL** (`repl.go`) - Interactive development environment

💡 **Design achieved**: The core is small enough to hold in your head (~366 lines in eval.go). This is our "Lisp microkernel".

## 2. ✅ Bootstrap Language in Itself
Higher-level constructs implemented using the language itself:

- **Built-in functions** (`bootstrap.go`) - `list`, `first`, `rest`, arithmetic, comparisons
- **User-defined functions** - Closures with lexical scoping
- **Macro system** - Code generation and language extension
- **Control structures** - Built using macros (`when`, `unless`, etc.)

This makes the language self-hosting and keeps the core clean.

## 3. ✅ Macro System Implementation
Full metaprogramming capabilities achieved:

- **Code-as-data manipulation** - Quasiquote (`` ` ``) and unquote (`~`) syntax
- **Macro definitions** - `defmacro` with Clojure-style square bracket `[param]` syntax
- **Language extensibility** - Users can define new control structures
- **Template system** - Selective evaluation within code templates

**Example Usage:**
```lisp
(defmacro when [condition body] `(if ~condition ~body nil))
(when true 42)  ; => 42

(defmacro unless [condition body] `(if ~condition nil ~body))
(unless false 99)  ; => 99
```

This makes the language infinitely extensible through metaprogramming.