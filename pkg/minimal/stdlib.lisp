;; Minimal Standard Library for Lisp Kernel
;; Essential higher-level functions built on the core primitives
;; Demonstrates building the language in itself

;; =============================================================================
;; LIST UTILITIES
;; =============================================================================

(define length 
  (fn [lst]
    ;; Count elements in a list
    (if (= lst ())
      0
      (+ 1 (length (rest lst))))))

(define append 
  (fn [lst item]
    ;; Add item to end of list (inefficient but simple)
    (if (= lst ())
      (list item)
      (cons (first lst) (append (rest lst) item)))))

(define reverse 
  (fn [lst]
    ;; Reverse a list
    (if (= lst ())
      ()
      (append (reverse (rest lst)) (first lst)))))

;; =============================================================================
;; FUNCTIONAL OPERATIONS
;; =============================================================================

(define map 
  (fn [f lst]
    ;; Apply function to each element
    (if (= lst ())
      ()
      (cons (f (first lst)) (map f (rest lst))))))

(define filter 
  (fn [pred lst]
    ;; Keep elements that satisfy predicate
    (if (= lst ())
      ()
      (if (pred (first lst))
        (cons (first lst) (filter pred (rest lst)))
        (filter pred (rest lst))))))

(define reduce 
  (fn [f init lst]
    ;; Reduce list with function and initial value
    (if (= lst ())
      init
      (reduce f (f init (first lst)) (rest lst)))))

;; =============================================================================
;; MATHEMATICAL UTILITIES
;; =============================================================================

(define sum 
  (fn [lst]
    ;; Sum all numbers in list
    (reduce + 0 lst)))

(define product 
  (fn [lst]
    ;; Multiply all numbers in list
    (reduce * 1 lst)))

(define range 
  (fn [n]
    ;; Create list of numbers from 0 to n-1
    (if (= n 0)
      ()
      (cons (- n 1) (range (- n 1))))))

;; =============================================================================
;; EXAMPLES
;; =============================================================================

(define demo 
  (fn []
    ;; Demonstrate standard library
    (print "=== Standard Library Demo ===")
    (define nums (list 1 2 3 4 5))
    (print "Numbers:" nums)
    (print "Length:" (length nums))
    (print "Sum:" (sum nums))
    (print "Doubled:" (map (fn [x] (* x 2)) nums))
    (print "Filter > 3:" (filter (fn [x] (< 3 x)) nums))
    (print "Range 5:" (range 5))
    (print "=== Demo Complete ===")))

;; Load notification
(print "Minimal standard library loaded!")
(print "Call (demo) to see examples.")

;; Load notification
(print "Minimal standard library loaded!")
(print "Call (demo) to see examples.")
