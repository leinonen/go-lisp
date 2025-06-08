;; Core Library Demo
;; This file demonstrates the use of the core library functions

;; First, load the core library
(load "library/core.lisp")

;; Import all core functions into the current namespace
(import core)

;; Demonstrate mathematical functions
(factorial 5)      ; Should return 120
(factorial 0)      ; Should return 1
(fibonacci 10)     ; Should return 55
(gcd 48 18)        ; Should return 6
(lcm 12 8)         ; Should return 24
(abs -42)          ; Should return 42
(abs 17)           ; Should return 17
(min 10 5)         ; Should return 5
(max 10 5)         ; Should return 10

;; Demonstrate list utility functions
(length-sq (list 1 2 3))                    ; Should return 9 (3^2)
(take 3 (list 1 2 3 4 5 6))                ; Should return (1 2 3)
(drop 2 (list 1 2 3 4 5))                  ; Should return (3 4 5)

;; Test edge cases
(take 0 (list 1 2 3))                      ; Should return ()
(drop 0 (list 1 2 3))                      ; Should return (1 2 3)
(take 10 (list 1 2 3))                     ; Should return (1 2 3)
(drop 10 (list 1 2 3))                     ; Should return ()

;; Demonstrate qualified names (alternative to import)
(core.factorial 6)                          ; Should return 720
(core.fibonacci 7)                          ; Should return 13
