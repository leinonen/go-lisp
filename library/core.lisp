; Core Library Module
; This module provides a collection of useful functions for common operations.
; It includes mathematical functions, list utilities, and other helpful tools.

(module core
  (export factorial fibonacci gcd lcm abs min max
          length-sq all any take drop compose apply-n)
  
  ; Mathematical Functions
  
  ; Factorial function - computes n! using tail recursion
  ; Public interface that delegates to private tail-recursive implementation
  (defun factorial [n]
    (if (< n 0)
        (error "factorial argument must be non-negative")
        (fact-tail n 1)))
  
  ; Private tail-recursive factorial helper
  ; This function is not exported and only used internally
  (defun fact-tail [n acc]
    (if (= n 0)
        acc
        (fact-tail (- n 1) (* n acc))))
  
  ; Fibonacci sequence using tail recursion
  (defun fibonacci [n]
    (if (< n 0)
        (error "fibonacci argument must be non-negative") 
        (if (< n 2)
            n
            (fib-tail n 0 1))))
            
  ; Private tail-recursive fibonacci helper
  (defun fib-tail [n a b]
    (if (= n 0)
        a
        (fib-tail (- n 1) b (+ a b))))
  
  ; Greatest Common Divisor using Euclidean algorithm
  (defun gcd [a b]
    (if (= b 0)
        a
        (gcd b (% a b))))
        
  ; Least Common Multiple 
  (defun lcm [a b]
    (/ (* a b) (gcd a b)))
  
  ; Absolute value
  (defun abs [x]
    (if (< x 0) (- 0 x) x))
  
  ; Minimum of two numbers
  (defun min [a b]
    (if (< a b) a b))
    
  ; Maximum of two numbers  
  (defun max [a b]
    (if (> a b) a b))
  
  ; List Utility Functions
  
  ; Square of the length of a list (useful for complexity analysis)
  (defun length-sq [lst]
    (* (length lst) (length lst)))
      
  ; Check if all elements in a list satisfy a predicate
  (defun all [predicate lst]
    (if (empty? lst)
        #t
        (if (predicate (first lst))
            (all predicate (rest lst))
            #f)))
            
  ; Check if any element in a list satisfies a predicate
  (defun any [predicate lst]
    (if (empty? lst)
        #f
        (if (predicate (first lst))
            #t
            (any predicate (rest lst)))))
            
  ; Take first n elements from a list
  (defun take [n lst]
    (if (or (= n 0) (empty? lst))
        (list)
        (cons (first lst) (take (- n 1) (rest lst)))))
        
  ; Drop first n elements from a list  
  (defun drop [n lst]
    (if (or (= n 0) (empty? lst))
        lst
        (drop (- n 1) (rest lst))))
  
  ; Higher-Order Utility Functions  
  
  ; Function composition - (compose f g) returns a function that applies g then f
  (defun compose [f g]
    (lambda [x] (f (g x))))
    
  ; Apply a function n times - (apply-n f n x) = f(f(...f(x)...))
  (defun apply-n [f n x]
    (if (= n 0)
        x
        (apply-n f (- n 1) (f x))))
)
