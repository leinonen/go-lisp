; Core Library Demo
; This file demonstrates the built-in core library functions

; Load the core library
(load "library/core.lisp")

; Import all core functions
(import core)

; Mathematical functions
(factorial 6)                           ; => 720
(fibonacci 10)                          ; => 55
(gcd 48 18)                            ; => 6
(abs -42)                              ; => 42

; List utility functions
(def test-list (list 1 2 3 4 5 6))
(take 3 test-list)                      ; => (1 2 3)
(drop 2 test-list)                      ; => (3 4 5 6)

; Predicate functions
(all (fn [x] (> x 0)) test-list)    ; => true (all positive)
(any (fn [x] (> x 5)) test-list)    ; => true (some > 5)

; Higher-order function utilities
(def double (fn [x] (* x 2)))
(def add-one (fn [x] (+ x 1)))

; Function composition: apply add-one then double
(def double-then-add-one (compose add-one double))
(double-then-add-one 5)                 ; => 11

; Apply a function multiple times
(apply-n add-one 3 5)                   ; => 8
