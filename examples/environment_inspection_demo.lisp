; Environment Inspection Demo
; This file demonstrates the new environment inspection functions

; Start with empty environment
(env)      ; Should show empty list: ()
(modules)  ; Should show empty list: ()

; Define some variables
(define pi 3.14159)
(define name "Lisp Interpreter")
(define version 1.0)

; Define some functions
(defun circle-area (radius)
  (* pi radius radius))

(defun greet (name)
  (+ "Hello, " name "!"))

; Check environment after definitions
(env)      ; Should show our variables and functions

; Create a math utilities module
(module math-utils
  (export square cube factorial power)
  
  (defun square (x) (* x x))
  (defun cube (x) (* x x x))
  (defun factorial (n)
    (if (= n 0)
      1
      (* n (factorial (- n 1)))))
  (defun power (base exp)
    (if (= exp 0)
      1
      (* base (power base (- exp 1)))))
  
  ; Private helper function (not exported)
  (defun helper (x) (+ x 1)))

; Create a string utilities module  
(module string-utils
  (export reverse-words count-chars)
  
  (defun reverse-words (text) text)  ; simplified
  (defun count-chars (text) 0))      ; simplified

; Check modules after creation
(modules)  ; Should show both modules with their exports

; Import math-utils and check environment
(import math-utils)
(env)      ; Should now include imported functions

; Test the functionality
(square 5)           ; => 25
(factorial 4)        ; => 24
(circle-area 2)      ; => 12.566...

; Demonstrate qualified access
(string-utils.reverse-words "Hello World")

; Final environment and module inspection
(env)
(modules)
