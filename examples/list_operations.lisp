; List Operations Examples
; This file demonstrates the additional list operations: append, reverse, and nth

; Define some example lists
(define list1 (list 1 2 3))
(define list2 (list 4 5 6))
(define mixed-list (list "hello" 42 #t 3.14))

; APPEND - Combines two lists into one
(append list1 list2)                    ; => (1 2 3 4 5 6)
(append (list) list1)                   ; => (1 2 3)
(append list1 (list))                   ; => (1 2 3)
(append (list "a" "b") (list "c" "d"))  ; => ("a" "b" "c" "d")

; REVERSE - Reverses the order of elements in a list
(reverse list1)                         ; => (3 2 1)
(reverse (list))                        ; => ()
(reverse (list 42))                     ; => (42)
(reverse mixed-list)                    ; => (3.14 #t 42 "hello")

; NTH - Gets the element at a specific index (0-based)
(nth 0 list1)                          ; => 1 (first element)
(nth 1 list1)                          ; => 2 (second element)
(nth 2 list1)                          ; => 3 (third element)
(nth 0 mixed-list)                     ; => "hello"
(nth 2 mixed-list)                     ; => #t

; Combining operations
(reverse (append list1 list2))         ; => (6 5 4 3 2 1)
(nth 0 (reverse list1))                ; => 3 (first element of reversed list)
(append (reverse list1) (reverse list2)) ; => (3 2 1 6 5 4)

; Working with single elements
(nth 0 (list "only"))                  ; => "only"
(reverse (append (list 1) (list 2)))   ; => (2 1)

; Error cases (these would cause errors):
; (nth 10 list1)                       ; Error: index out of bounds
; (nth -1 list1)                       ; Error: negative index
; (nth 0 (list))                       ; Error: empty list

; Complex example: Processing a list
(define numbers (list 1 2 3 4 5))
(define doubled (map (lambda (x) (* x 2)) numbers))
(define reversed-doubled (reverse doubled))
(define first-three (list (nth 0 reversed-doubled) 
                         (nth 1 reversed-doubled) 
                         (nth 2 reversed-doubled)))
first-three                            ; => (10 8 6)

; Another example: Building lists incrementally
(define step1 (list 1))
(define step2 (append step1 (list 2)))
(define step3 (append step2 (list 3)))
(define final (reverse step3))
final                                  ; => (3 2 1)
