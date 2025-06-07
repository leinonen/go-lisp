; Math module - provides basic mathematical operations
(module math
  (export square cube add-squares)
  
  ; Define a function to square a number
  (defun square (x) (* x x))
  
  ; Define a function to cube a number  
  (defun cube (x) (* x x x))
  
  ; Define a function that adds the squares of two numbers
  (defun add-squares (x y) (+ (square x) (square y)))
  
  ; This is a private helper function - not exported
  (defun private-helper (x) (+ x 1))
)
