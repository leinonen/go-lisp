; Basic function definition and usage
(define square (lambda (x) (* x x)))
(square 5)

; Function with multiple parameters
(define add (lambda (x y) (+ x y)))
(add 3 4)

; Higher-order functions
(define apply-twice (lambda (f x) (f (f x))))
(define increment (lambda (x) (+ x 1)))
(apply-twice increment 5)

; Closures - functions that capture variables from their environment
(define make-counter (lambda (start) 
  (lambda () (define start (+ start 1)) start)))
(define counter (make-counter 0))

; Recursive functions
(define fibonacci (lambda (n) 
  (if (< n 2) 
    n 
    (+ (fibonacci (- n 1)) (fibonacci (- n 2))))))
(fibonacci 10)

; Factorial using recursion
(define factorial (lambda (n) 
  (if (= n 0) 
    1 
    (* n (factorial (- n 1))))))
(factorial 6)

; Environment inspection
(builtins)  ; List all available built-in functions
(env)       ; Show current environment variables/functions
(modules)   ; Show available modules
