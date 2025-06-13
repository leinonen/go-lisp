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

  (defn identity [x]
    "Return the argument unchanged"
    x)

  (defn constantly [value]
    "Return a function that always returns the given value"
    (fn [x] value))

  (defn complement [predicate]
    "Return a function that returns the logical complement of predicate"
    (fn [x] (not (predicate x))))

  ;; ============================================================================
  ;; Partial Application
  ;; ============================================================================

  (defn partial [f arg1]
    "Partial application: fix the first argument of a function"
    (fn [arg2] (f arg1 arg2)))

  (defn partial2 [f arg1 arg2]
    "Partial application: fix the first two arguments of a function"
    (fn [arg3] (f arg1 arg2 arg3)))

  (defn partial3 [f arg1 arg2 arg3]
    "Partial application: fix the first three arguments of a function"
    (fn [arg4] (f arg1 arg2 arg3 arg4)))

  ;; ============================================================================
  ;; Currying (Transform multi-argument function into chain of single-argument functions)
  ;; ============================================================================

  (defn curry2 [f]
    "Curry a 2-argument function"
    (fn [arg1]
      (fn [arg2]
        (f arg1 arg2))))

  (defn curry3 [f]
    "Curry a 3-argument function"
    (fn [arg1]
      (fn [arg2]
        (fn [arg3]
          (f arg1 arg2 arg3)))))

  ;; Generic curry for common case (2-argument functions)
  (defn curry [f]
    "Curry a function (assumes 2 arguments)"
    (curry2 f))

  ;; ============================================================================
  ;; Function Composition
  ;; ============================================================================

  (defn comp [f g]
    "Compose two functions: (comp f g)(x) = f(g(x))"
    (fn [x] (f (g x))))

  (defn comp3 [f g h]
    "Compose three functions: (comp3 f g h)(x) = f(g(h(x)))"
    (fn [x] (f (g (h x)))))

  (defn comp4 [f g h i]
    "Compose four functions: (comp4 f g h i)(x) = f(g(h(i(x))))"
    (fn [x] (f (g (h (i x))))))

  ;; ============================================================================
  ;; Pipeline Operations  
  ;; ============================================================================

  (defn pipe [x f]
    "Apply function f to value x: pipe(x, f) = f(x)"
    (f x))

  (defn pipe2 [x f g]
    "Apply functions in sequence: pipe2(x, f, g) = g(f(x))"
    (g (f x)))

  (defn pipe3 [x f g h]
    "Apply functions in sequence: pipe3(x, f, g, h) = h(g(f(x)))"
    (h (g (f x))))

  (defn pipe4 [x f g h i]
    "Apply functions in sequence: pipe4(x, f, g, h, i) = i(h(g(f(x))))"
    (i (h (g (f x)))))

  ;; ============================================================================
  ;; Juxtaposition (Apply multiple functions to same input)
  ;; ============================================================================

  (defn juxt [f g]
    "Apply multiple functions to same input: juxt(f, g)(x) = [f(x), g(x)]"
    (fn [x] (list (f x) (g x))))

  (defn juxt3 [f g h]
    "Apply three functions to same input"
    (fn [x] (list (f x) (g x) (h x))))

  (defn juxt4 [f g h i]
    "Apply four functions to same input"
    (fn [x] (list (f x) (g x) (h x) (i x))))

  ;; ============================================================================
  ;; Conditional Functions
  ;; ============================================================================

  (defn if-fn [predicate then-fn else-fn]
    "Return a function that conditionally applies then-fn or else-fn"
    (fn [x]
      (if (predicate x)
          (then-fn x)
          (else-fn x))))

  (defn when-fn [predicate then-fn]
    "Return a function that applies then-fn when predicate is true, else identity"
    (fn [x]
      (if (predicate x)
          (then-fn x)
          x)))

  (defn unless-fn [predicate else-fn]
    "Return a function that applies else-fn when predicate is false, else identity"
    (fn [x]
      (if (predicate x)
          x
          (else-fn x))))

  ;; ============================================================================
  ;; Predicate Combinators
  ;; ============================================================================

  (defn every-pred [pred1 pred2]
    "Return a predicate that is true when both predicates are true"
    (fn [x] (and (pred1 x) (pred2 x))))

  (defn every-pred3 [pred1 pred2 pred3]
    "Return a predicate that is true when all three predicates are true"
    (fn [x] (and (pred1 x) (pred2 x) (pred3 x))))

  (defn some-pred [pred1 pred2]
    "Return a predicate that is true when either predicate is true"
    (fn [x] (or (pred1 x) (pred2 x))))

  (defn some-pred3 [pred1 pred2 pred3]
    "Return a predicate that is true when any predicate is true"
    (fn [x] (or (pred1 x) (pred2 x) (pred3 x))))

  ;; ============================================================================
  ;; Threading/Pipeline Macros (using functions)
  ;; ============================================================================

  (defn thread-first [x f]
    "Thread value x through function f as first argument"
    (f x))

  (defn thread-last [x f]
    "Thread value x through function f as last argument (for single-arg functions, same as thread-first)"
    (f x))

  ;; ============================================================================
  ;; Function Application Utilities
  ;; ============================================================================

  (defn apply-to [value f]
    "Apply function to value (reverse of normal application)"
    (f value))

  ;; ============================================================================
  ;; Memoization
  ;; ============================================================================

  (defn memoize [f]
    "Return a memoized version of function f (simplified implementation)"
    ;; Note: This is a simplified memoization without actual caching
    ;; In a full implementation, we'd maintain a cache hash-map
    ;; For now, we just return the original function
    f)

  ;; ============================================================================
  ;; Function Introspection
  ;; ============================================================================

  (defn arity [f]
    "Return the arity (number of parameters) of a function (placeholder)"
    ;; This would require runtime function introspection
    ;; For now, return a placeholder value
    2)

  ;; ============================================================================
  ;; Higher-Order Utilities
  ;; ============================================================================

  (defn fnil [f default-val]
    "Return a function that replaces nil arguments with default-val before calling f"
    (fn [x]
      (if (= x nil)  ; Direct nil comparison
          (f default-val)
          (f x))))

  (defn fnth [n f]
    "Return a function that applies f to the nth element of a list"
    (fn [lst]
      (f (nth n lst))))

  ;; ============================================================================
  ;; Function Combinators for Lists
  ;; ============================================================================

  (defn map-indexed [f lst]
    "Map function over list with index as second argument"
    (map-indexed-helper f lst 0))

  ;; Helper for map-indexed
  (defn map-indexed-helper [f lst index]
    (if (empty? lst)
        (list)
        (cons (f index (first lst))
              (map-indexed-helper f (rest lst) (+ index 1)))))

  (defn keep [f lst]
    "Apply f to each element, keep non-nil results"
    (filter (fn [x] (not (= x nil))) (map f lst)))

  (defn keep-indexed [f lst]
    "Apply f to index and element, keep non-nil results"
    (keep identity (map-indexed f lst)))

)
