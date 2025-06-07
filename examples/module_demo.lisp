; Example demonstrating module usage
; Load the math module
(load "examples/math_module.lisp")

; Import all exported functions from math module
(import math)

; Now we can use the imported functions directly
(define result1 (square 5))    ; Should be 25
(define result2 (cube 3))      ; Should be 27
(define result3 (add-squares 3 4)) ; Should be 25 (9 + 16)

; We can also use qualified access without importing
(load "examples/string_module.lisp")

; Use qualified access to call functions from string-utils module
(define result4 (string-utils.concat-three "Hello" " " "World"))

; Test that private functions are not accessible
; This should fail:
; (private-helper 5)  ; Uncomment to see error

; Display results
result1
result2  
result3
result4
