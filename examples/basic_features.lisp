; Basic Language Features Demo
; This file demonstrates the core features of the Lisp interpreter

; Basic arithmetic and variables
(define x 10)
(define y 20)
(+ x y (* 3 4))

; Function definition with defun
(defun square [n] (* n n))
(defun add [a b] (+ a b))

; List operations
(define my-list (list 1 2 3 4 5))
(first my-list)
(rest my-list)
(cons 0 my-list)
(length my-list)

; Higher-order functions
(map square my-list)                    ; => (1 4 9 16 25)
(filter (lambda [x] (> x 3)) my-list)  ; => (4 5)
(reduce (lambda [acc x] (+ acc x)) 0 my-list)  ; => 15

; Closures and function composition
(define make-adder (lambda [n] (lambda [x] (+ x n))))
(define add-five (make-adder 5))
(add-five 10)                           ; => 15

; Recursion with automatic tail call optimization
(defun factorial [n]
  (if (= n 0) 1 (* n (factorial (- n 1)))))
(factorial 10)                          ; => 3628800

; Conditional logic
(if (> (length my-list) 3)
    "List has more than 3 elements"
    "List has 3 or fewer elements")

; Big number support (automatic)
(* 1000000000000000 1000000000000000)

; Error handling
; (error "This would stop execution with a message")

; Environment inspection
(help)  ; Show all built-in functions
(env)       ; Show current variables and functions
