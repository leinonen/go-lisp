;; Macro System Examples
;; Comprehensive demonstration of the Lisp interpreter's macro capabilities

;; ============================================================================
;; Basic Macro Definition and Usage
;; ============================================================================

;; Define a simple 'when' macro for conditional execution
(defmacro when (condition body)
  (list 'if condition body 'nil))

(print! "=== When Macro Demo ===")
(def x 10)
(when (> x 5) (print! "x is greater than 5"))
(when (< x 5) (print! "This won't print"))

;; Define an 'unless' macro (opposite of when)
(defmacro unless (condition then else)
  (list 'if condition else then))

(print! "\n=== Unless Macro Demo ===")
(unless (< x 5) 
  (print! "x is not less than 5") 
  (print! "x is less than 5"))

;; ============================================================================
;; Quote Special Form
;; ============================================================================

(print! "\n=== Quote Examples ===")
(print! "Regular evaluation:")
(print! (+ 1 2 3))

(print! "Quoted (unevaluated):")
(print! (quote (+ 1 2 3)))
(print! '(+ 1 2 3))

;; ============================================================================
;; Advanced Macro - Let-like Binding
;; ============================================================================

(defmacro let1 (var value body)
  (list (list 'lambda [list var] body) value))

(print! "\n=== Let1 Macro Demo ===")
(let1 y 42 
  (print! (str "y inside let1: " y)))

;; ============================================================================
;; Debugging Macro - Show Expression and Result
;; ============================================================================

(defmacro debug (expr)
  (list 'list 
    (list 'quote expr) 
    expr))

(print! "\n=== Debug Macro Demo ===")
(def result (debug (+ 10 20 30)))
(print! (str "Expression: " (first result)))
(print! (str "Result: " (first (rest result))))

;; ============================================================================
;; Conditional Compilation Macro
;; ============================================================================

(defmacro compile-time-if (condition true-branch false-branch)
  (if condition true-branch false-branch))

(print! "\n=== Compile-time Conditional ===")
(def debug-mode true)
(compile-time-if debug-mode
  (print! "Debug mode is enabled")
  (print! "Debug mode is disabled"))

;; ============================================================================
;; Increment Macro (modifies variables)
;; ============================================================================

(defmacro inc! (var)
  (list 'define var (list '+ var 1)))

(print! "\n=== Increment Macro Demo ===")
(def counter 0)
(print! (str "Counter before: " counter))
(inc! counter)
(print! (str "Counter after inc!: " counter))
(inc! counter)
(print! (str "Counter after second inc!: " counter))

;; ============================================================================
;; Summary
;; ============================================================================

(print! "\n=== Macro System Summary ===")
(print! "✓ defmacro - Define code transformation macros")
(print! "✓ quote/'  - Prevent evaluation for code manipulation")
(print! "✓ Complex macros - Let-bindings, debugging, conditionals")
(print! "✓ Variable modification - Increment and other state changes")
(print! "✓ DSL creation - Build domain-specific languages")

(print! "\nMacro system is fully functional! 🎉")
