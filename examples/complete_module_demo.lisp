; Comprehensive Module System Demo
; ================================

; Step 1: Create a math utilities module
(module math-utils
  (export square cube add-squares factorial)
  
  ; Define exported functions
  (defun square (x) 
    (* x x))
    
  (defun cube (x) 
    (* x x x))
    
  (defun add-squares (x y) 
    (+ (square x) (square y)))
    
  (defun factorial (n)
    (if (= n 0)
      1
      (* n (factorial (- n 1)))))
      
  ; Define a private helper function (not exported)
  (defun private-helper (x) 
    (+ x 42)))

; Step 2: Create a list utilities module  
(module list-utils
  (export sum-list product-list)
  
  (defun sum-list (lst)
    (if (empty? lst)
      0
      (+ (first lst) (sum-list (rest lst)))))
      
  (defun product-list (lst)
    (if (empty? lst)
      1
      (* (first lst) (product-list (rest lst))))))

; Step 3: Import math-utils and test imported functions
(import math-utils)

; Test imported functions
(define test1 (square 5))      ; Should be 25
(define test2 (cube 3))        ; Should be 27  
(define test3 (add-squares 3 4)) ; Should be 25 (9 + 16)
(define test4 (factorial 5))   ; Should be 120

; Step 4: Test qualified access without importing
(define test5 (list-utils.sum-list (list 1 2 3 4)))     ; Should be 10
(define test6 (list-utils.product-list (list 2 3 4)))   ; Should be 24

; Step 5: Test that qualified access works for imported modules too
(define test7 (math-utils.square 6))  ; Should be 36

; Display all results
test1
test2  
test3
test4
test5
test6
test7

; Note: The following would cause an error because private-helper is not exported:
; (private-helper 5)  ; Uncomment to see error
