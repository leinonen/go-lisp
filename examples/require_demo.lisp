; Require Function Demo
; This file demonstrates the basic require function for simplified module loading

; The require function combines load and import in a single operation

(print! "=== Require Function Demo ===")

; Traditional approach - two separate commands
(print! "Loading core library with require:")
(require "library/core.lisp")

(print! "Functions now available:")
(print! "Factorial of 5:")
(print! (factorial 5))

(print! "GCD of 48 and 18:")
(print! (gcd 48 18))

(print! "Fibonacci of 10:")
(print! (fibonacci 10))

(print! "")

; Load functional library
(print! "Loading functional library:")
(require "library/functional.lisp")

(print! "Using functional library functions:")
(print! "Map with double function:")
(print! (map (fn [x] (* x 2)) (list 1 2 3 4 5)))

(print! "Filter even numbers:")
(print! (filter (fn [x] (= (% x 2) 0)) (list 1 2 3 4 5 6 7 8)))

(print! "Compose function example:")
(def add-one (fn [x] (+ x 1)))
(def double (fn [x] (* x 2)))
(def add-one-then-double (comp double add-one))
(print! "Add 1 then double 5:")
(print! (add-one-then-double 5))

(print! "")

; Show that require prevents re-loading
(print! "=== File Caching Demo ===")
(print! "First require of core.lisp already done above")
(print! "Second require of core.lisp (should be cached):")
(require "library/core.lisp")
(print! "Functions still available after second require:")
(print! (abs -42))

(print! "")

; Show module introspection
(print! "=== Module Introspection ===")
(print! "Available modules:")
(modules)

(print! "")
(print! "Require demo completed successfully!")
