; Higher-Order Function Examples for Lisp Interpreter
; This file demonstrates the newly implemented map, filter, and reduce functions

; Basic map usage - square each number
(map (lambda (x) (* x x)) (list 1 2 3 4 5))
; Expected: (1 4 9 16 25)

; Basic filter usage - keep only positive numbers
(filter (lambda (x) (> x 0)) (list -2 -1 0 1 2 3))
; Expected: (1 2 3)

; Basic reduce usage - sum all numbers
(reduce (lambda (acc x) (+ acc x)) 0 (list 1 2 3 4 5))
; Expected: 15

; Product using reduce
(reduce (lambda (acc x) (* acc x)) 1 (list 2 3 4))
; Expected: 24

; Complex combinations

; 1. Filter positive numbers, then square them
(map (lambda (x) (* x x)) (filter (lambda (x) (> x 0)) (list -2 -1 0 1 2 3)))
; Expected: (1 4 9)

; 2. Square numbers, then sum them
(reduce (lambda (acc x) (+ acc x)) 0 (map (lambda (x) (* x x)) (list 1 2 3)))
; Expected: 14 (1 + 4 + 9)

; 3. Filter, map, then reduce - sum of squares of positive numbers
(reduce (lambda (acc x) (+ acc x)) 0 
  (map (lambda (x) (* x x)) 
    (filter (lambda (x) (> x 0)) (list -3 -2 -1 0 1 2 3 4))))
; Expected: 30 (1 + 4 + 9 + 16)

; 4. Count elements using reduce
(reduce (lambda (acc x) (+ acc 1)) 0 (list "a" "b" "c" "d" "e"))
; Expected: 5

; 5. Find maximum using reduce
(reduce (lambda (acc x) (if (> x acc) x acc)) 0 (list 3 7 2 9 1 8))
; Expected: 9

; 6. Double positive numbers and sum them
(reduce (lambda (acc x) (+ acc x)) 0 
  (map (lambda (x) (* x 2)) 
    (filter (lambda (x) (> x 0)) (list -1 2 -3 4 5))))
; Expected: 22 (4 + 8 + 10)
