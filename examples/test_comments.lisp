; This is a comment at the beginning
; Let's test basic arithmetic with comments

(+ 1 2) ; This should return 3

; Define a simple function
(defun square (x) ; x is the parameter
  (* x x)) ; multiply x by itself

; Test the function
(square 5) ; Should return 25

; Multiple expressions with comments
(+ 1 2 3) ; Addition
(* 2 3 4) ; Multiplication
(- 10 3)  ; Subtraction

; Comments with special characters ()[]{}!@#$%^&*
; Even quotes "like this" and symbols + - * / should work in comments
42 ; Final result
