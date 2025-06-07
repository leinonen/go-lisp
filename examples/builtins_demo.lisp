; Built-in Functions Demo
; This file demonstrates how to inspect and use the builtins function

; Get list of all available built-in functions
(define available-functions (builtins))

; Show all built-in functions
available-functions

; Count how many built-in functions are available
(length available-functions)

; Check if specific functions are available
(define has-plus (member? "+" available-functions))
(define has-map (member? "map" available-functions))
(define has-filter (member? "filter" available-functions))

; Show results
has-plus    ; Should be true
has-map     ; Should be true  
has-filter  ; Should be true

; Demonstrate using some of the built-in functions
; Arithmetic operations
(+ 5 3)     ; Addition
(* 4 6)     ; Multiplication
(/ 15 3)    ; Division

; List operations
(define my-list (list 1 2 3 4 5))
(first my-list)    ; Get first element
(rest my-list)     ; Get rest of list
(length my-list)   ; Get list length

; Higher-order functions
(map (lambda (x) (* x 2)) my-list)        ; Double each element
(filter (lambda (x) (> x 3)) my-list)     ; Filter elements > 3
(reduce + my-list 0)                      ; Sum all elements

; Control flow
(if (> (length my-list) 3)
    "List has more than 3 elements"
    "List has 3 or fewer elements")

; Function definition
(defun square (x) (* x x))
(square 7)

; Environment inspection
(env)      ; Show current environment
(modules)  ; Show current modules
(builtins) ; Show built-in functions again

; Helper function to check if a function exists in builtins
(defun function-exists? (func-name)
  (member? func-name (builtins)))

; Test the helper function
(function-exists? "cons")     ; Should be true
(function-exists? "unknown")  ; Should be false

; Print summary
(+ "Total built-in functions available: " (length (builtins)))
