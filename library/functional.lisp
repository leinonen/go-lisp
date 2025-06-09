;; ============================================================================
;; Functional Programming Library
;; ============================================================================
;; This library provides essential functional programming utilities including
;; function combinators, composition operators, and higher-order function tools.

(module functional
  (export 
    ;; Function combinators
    curry curry2 curry3 partial partial2 partial3 complement constantly identity
    ;; Function composition  
    comp comp3 comp4 pipe pipe2 pipe3 pipe4 juxt juxt3 juxt4
    ;; Conditional functions
    if-fn when-fn unless-fn
    ;; Predicate utilities
    every-pred every-pred3 some-pred some-pred3
    ;; Function application
    apply-to thread-first thread-last
    ;; Memoization
    memoize
    ;; Function introspection
    arity
    ;; Higher-order utilities
    fnil fnth map-indexed keep keep-indexed)

  ;; ============================================================================
  ;; Basic Function Combinators
  ;; ============================================================================

  (defun identity (x)
    "Return the argument unchanged"
    x)

  (defun constantly (value)
    "Return a function that always returns the given value"
    (lambda (x) value))

  (defun complement (predicate)
    "Return a function that returns the logical complement of predicate"
    (lambda (x) (not (predicate x))))

  ;; ============================================================================
  ;; Partial Application
  ;; ============================================================================

  (defun partial (fn arg1)
    "Partial application: fix the first argument of a function"
    (lambda (arg2) (fn arg1 arg2)))

  (defun partial2 (fn arg1 arg2)
    "Partial application: fix the first two arguments of a function"
    (lambda (arg3) (fn arg1 arg2 arg3)))

  (defun partial3 (fn arg1 arg2 arg3)
    "Partial application: fix the first three arguments of a function"
    (lambda (arg4) (fn arg1 arg2 arg3 arg4)))

  ;; ============================================================================
  ;; Currying (Transform multi-argument function into chain of single-argument functions)
  ;; ============================================================================

  (defun curry2 (fn)
    "Curry a 2-argument function"
    (lambda (arg1)
      (lambda (arg2)
        (fn arg1 arg2))))

  (defun curry3 (fn)
    "Curry a 3-argument function"
    (lambda (arg1)
      (lambda (arg2)
        (lambda (arg3)
          (fn arg1 arg2 arg3)))))

  ;; Generic curry for common case (2-argument functions)
  (defun curry (fn)
    "Curry a function (assumes 2 arguments)"
    (curry2 fn))

  ;; ============================================================================
  ;; Function Composition
  ;; ============================================================================

  (defun comp (f g)
    "Compose two functions: (comp f g)(x) = f(g(x))"
    (lambda (x) (f (g x))))

  (defun comp3 (f g h)
    "Compose three functions: (comp3 f g h)(x) = f(g(h(x)))"
    (lambda (x) (f (g (h x)))))

  (defun comp4 (f g h i)
    "Compose four functions: (comp4 f g h i)(x) = f(g(h(i(x))))"
    (lambda (x) (f (g (h (i x))))))

  ;; ============================================================================
  ;; Pipeline Operations  
  ;; ============================================================================

  (defun pipe (x f)
    "Apply function f to value x: pipe(x, f) = f(x)"
    (f x))

  (defun pipe2 (x f g)
    "Apply functions in sequence: pipe2(x, f, g) = g(f(x))"
    (g (f x)))

  (defun pipe3 (x f g h)
    "Apply functions in sequence: pipe3(x, f, g, h) = h(g(f(x)))"
    (h (g (f x))))

  (defun pipe4 (x f g h i)
    "Apply functions in sequence: pipe4(x, f, g, h, i) = i(h(g(f(x))))"
    (i (h (g (f x)))))

  ;; ============================================================================
  ;; Juxtaposition (Apply multiple functions to same input)
  ;; ============================================================================

  (defun juxt (f g)
    "Apply multiple functions to same input: juxt(f, g)(x) = [f(x), g(x)]"
    (lambda (x) (list (f x) (g x))))

  (defun juxt3 (f g h)
    "Apply three functions to same input"
    (lambda (x) (list (f x) (g x) (h x))))

  (defun juxt4 (f g h i)
    "Apply four functions to same input"
    (lambda (x) (list (f x) (g x) (h x) (i x))))

  ;; ============================================================================
  ;; Conditional Functions
  ;; ============================================================================

  (defun if-fn (predicate then-fn else-fn)
    "Return a function that conditionally applies then-fn or else-fn"
    (lambda (x)
      (if (predicate x)
          (then-fn x)
          (else-fn x))))

  (defun when-fn (predicate then-fn)
    "Return a function that applies then-fn when predicate is true, else identity"
    (lambda (x)
      (if (predicate x)
          (then-fn x)
          x)))

  (defun unless-fn (predicate else-fn)
    "Return a function that applies else-fn when predicate is false, else identity"
    (lambda (x)
      (if (predicate x)
          x
          (else-fn x))))

  ;; ============================================================================
  ;; Predicate Combinators
  ;; ============================================================================

  (defun every-pred (pred1 pred2)
    "Return a predicate that is true when both predicates are true"
    (lambda (x) (and (pred1 x) (pred2 x))))

  (defun every-pred3 (pred1 pred2 pred3)
    "Return a predicate that is true when all three predicates are true"
    (lambda (x) (and (pred1 x) (pred2 x) (pred3 x))))

  (defun some-pred (pred1 pred2)
    "Return a predicate that is true when either predicate is true"
    (lambda (x) (or (pred1 x) (pred2 x))))

  (defun some-pred3 (pred1 pred2 pred3)
    "Return a predicate that is true when any predicate is true"
    (lambda (x) (or (pred1 x) (pred2 x) (pred3 x))))

  ;; ============================================================================
  ;; Threading/Pipeline Macros (using functions)
  ;; ============================================================================

  (defun thread-first (x f)
    "Thread value x through function f as first argument"
    (f x))

  (defun thread-last (x f)
    "Thread value x through function f as last argument (for single-arg functions, same as thread-first)"
    (f x))

  ;; ============================================================================
  ;; Function Application Utilities
  ;; ============================================================================

  (defun apply-to (value fn)
    "Apply function to value (reverse of normal application)"
    (fn value))

  ;; ============================================================================
  ;; Memoization
  ;; ============================================================================

  (defun memoize (fn)
    "Return a memoized version of function fn (simplified implementation)"
    ;; Note: This is a simplified memoization without actual caching
    ;; In a full implementation, we'd maintain a cache hash-map
    ;; For now, we just return the original function
    fn)

  ;; ============================================================================
  ;; Function Introspection
  ;; ============================================================================

  (defun arity (fn)
    "Return the arity (number of parameters) of a function (placeholder)"
    ;; This would require runtime function introspection
    ;; For now, return a placeholder value
    2)

  ;; ============================================================================
  ;; Higher-Order Utilities
  ;; ============================================================================

  (defun fnil (fn default-val)
    "Return a function that replaces nil arguments with default-val before calling fn"
    (lambda (x)
      (if (= x nil)  ; Direct nil comparison
          (fn default-val)
          (fn x))))

  (defun fnth (n fn)
    "Return a function that applies fn to the nth element of a list"
    (lambda (lst)
      (fn (nth n lst))))

  ;; ============================================================================
  ;; Function Combinators for Lists
  ;; ============================================================================

  (defun map-indexed (fn lst)
    "Map function over list with index as second argument"
    (map-indexed-helper fn lst 0))

  ;; Helper for map-indexed
  (defun map-indexed-helper (fn lst index)
    (if (empty? lst)
        (list)
        (cons (fn index (first lst))
              (map-indexed-helper fn (rest lst) (+ index 1)))))

  (defun keep (fn lst)
    "Apply fn to each element, keep non-nil results"
    (filter (lambda (x) (not (= x nil))) (map fn lst)))

  (defun keep-indexed (fn lst)
    "Apply fn to index and element, keep non-nil results"
    (keep identity (map-indexed fn lst)))

)
