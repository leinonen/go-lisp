; tail_call_optimization.lisp
; 
; This file demonstrates Tail Call Optimization (TCO) in the Lisp interpreter.
; TCO prevents stack overflow by optimizing tail-recursive functions to use
; constant stack space regardless of recursion depth.

; ==============================================================================
; BASIC TAIL RECURSION
; ==============================================================================

; Traditional factorial (NOT tail-recursive - grows stack)
(defun factorial (n)
  (if (= n 0) 
      1 
      (* n (factorial (- n 1)))))

; Tail-recursive factorial (stack-safe)
(defun fact-tail (n acc)
  (if (= n 0) 
      acc 
      (fact-tail (- n 1) (* n acc))))

; Wrapper function for convenience
(defun factorial-optimized (n)
  (fact-tail n 1))

; Usage examples:
; (factorial 5)           ; => 120 (but limited by stack size)
; (factorial-optimized 5) ; => 120 (no stack limit)
; (factorial-optimized 1000) ; => very large number (no stack overflow!)

; ==============================================================================
; NUMERIC COMPUTATIONS
; ==============================================================================

; Sum from 1 to n using tail recursion
(defun sum-tail (n acc)
  (if (= n 0) 
      acc 
      (sum-tail (- n 1) (+ acc n))))

(defun sum-to-n (n)
  (sum-tail n 0))

; Count down to zero using tail recursion
(defun countdown (n)
  (if (= n 0) 
      0 
      (countdown (- n 1))))

; Fibonacci using tail recursion (much more efficient than naive approach)
(defun fib-tail (n a b)
  (if (= n 0) 
      a 
      (fib-tail (- n 1) b (+ a b))))

(defun fibonacci (n)
  (fib-tail n 0 1))

; ==============================================================================
; LIST PROCESSING WITH TCO
; ==============================================================================

; Reverse a list using tail recursion
(defun reverse-tail (lst acc)
  (if (empty? lst) 
      acc 
      (reverse-tail (rest lst) (cons (first lst) acc))))

(defun reverse-list (lst)
  (reverse-tail lst (list)))

; Length of list using tail recursion
(defun length-tail (lst acc)
  (if (empty? lst) 
      acc 
      (length-tail (rest lst) (+ acc 1))))

(defun list-length (lst)
  (length-tail lst 0))

; Find maximum element in list using tail recursion
(defun max-tail (lst current-max)
  (if (empty? lst) 
      current-max 
      (max-tail (rest lst) 
                (if (> (first lst) current-max) 
                    (first lst) 
                    current-max))))

(defun find-max (lst)
  (if (empty? lst) 
      #f 
      (max-tail (rest lst) (first lst))))

; ==============================================================================
; MUTUALLY RECURSIVE FUNCTIONS
; ==============================================================================

; Even and odd predicates using mutual tail recursion
(defun even? (n)
  (if (= n 0) 
      #t 
      (odd? (- n 1))))

(defun odd? (n)
  (if (= n 0) 
      #f 
      (even? (- n 1))))

; ==============================================================================
; DEMONSTRATION FUNCTIONS
; ==============================================================================

; Function to demonstrate stack safety with large numbers
(defun test-large-recursion ()
  (define large-n 10000)
  
  ; These would cause stack overflow without TCO:
  (define sum-result (sum-to-n large-n))
  (define countdown-result (countdown large-n))
  (define even-result (even? large-n))
  
  ; Create a list with results
  (list sum-result countdown-result even-result))

; Function to compare tail vs non-tail recursion performance
(defun performance-demo ()
  ; Small number - both work
  (define small-n 10)
  (define fact-normal (factorial small-n))
  (define fact-optimized (factorial-optimized small-n))
  
  ; Large number - only optimized version works without stack overflow
  (define large-n 100)
  ; (factorial large-n)           ; Would cause stack overflow
  (define fact-large (factorial-optimized large-n))
  
  (list fact-normal fact-optimized fact-large))

; ==============================================================================
; INTERACTIVE EXAMPLES
; ==============================================================================

; Run these commands in the REPL to see TCO in action:

; Basic examples:
; (factorial-optimized 5)        ; => 120
; (sum-to-n 100)                 ; => 5050
; (fibonacci 10)                 ; => 55

; List processing:
; (reverse-list (list 1 2 3 4))  ; => (4 3 2 1)
; (find-max (list 3 7 2 9 1))    ; => 9

; Large computations (these work thanks to TCO):
; (factorial-optimized 100)      ; => very large number
; (sum-to-n 10000)              ; => 50005000
; (countdown 5000)              ; => 0
; (even? 9999)                  ; => #f

; Performance test:
; (test-large-recursion)         ; => (50005000 0 #t)
; (performance-demo)             ; => (3628800 3628800 very-large-number)

; ==============================================================================
; NOTES ON TAIL CALL OPTIMIZATION
; ==============================================================================

; What makes a call "tail-recursive":
; 1. The recursive call must be the LAST operation in the function
; 2. The result of the recursive call is returned directly (no further computation)
; 3. The call is in "tail position" - nothing happens after it returns

; Examples of TAIL calls:
; (if condition base-case (recursive-call args))  ; recursive call in tail position
; (recursive-call args)                           ; direct tail call

; Examples of NON-TAIL calls:
; (* n (recursive-call args))                     ; multiplication after recursive call
; (+ 1 (recursive-call args))                     ; addition after recursive call
; (cons x (recursive-call args))                  ; cons after recursive call

; The TCO implementation automatically detects tail calls and optimizes them
; to use constant stack space, making it safe to write deeply recursive
; algorithms without worrying about stack overflow.
