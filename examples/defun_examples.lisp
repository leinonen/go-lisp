; Examples demonstrating the new defun syntax
; which combines define and lambda into a single form

; Basic function definition using defun
(defun square (x) (* x x))
(square 5)

; Function with multiple parameters using defun
(defun add (x y) (+ x y))
(add 3 4)

; More complex function
(defun average (x y) (/ (+ x y) 2))
(average 10 20)

; Recursive function using defun
(defun factorial (n) 
  (if (= n 0) 
    1 
    (* n (factorial (- n 1)))))
(factorial 5)

; Function that returns another function (higher-order function)
(defun make-multiplier (factor) 
  (lambda (x) (* x factor)))
(define double (make-multiplier 2))
(double 7)

; Comparison: Old way vs New way
; Old way: (define old-square (lambda (x) (* x x)))
; New way: (defun new-square (x) (* x x))

; Both create the same functionality
(define old-square (lambda (x) (* x x)))
(defun new-square (x) (* x x))

; Test they work the same
(old-square 4)
(new-square 4)
