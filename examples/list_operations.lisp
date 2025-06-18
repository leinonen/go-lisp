;; List operations examples in Go Lisp

;; Creating lists
(list 1 2 3 4 5)                    ; => (1 2 3 4 5)
(list "hello" "world" "!")           ; => ("hello" "world" "!")

;; Getting elements
(first '(1 2 3))                     ; => 1
(rest '(1 2 3))                      ; => (2 3)
(last '(1 2 3))                      ; => 3
(nth 1 '(a b c d))                   ; => b (0-indexed)

;; Adding elements
(cons 0 '(1 2 3))                    ; => (0 1 2 3)
(append '(1 2) '(3 4) '(5 6))        ; => (1 2 3 4 5 6)

;; List properties
(length '(1 2 3 4))                  ; => 4
(empty? '())                         ; => true
(empty? '(1))                        ; => false

;; List transformations
(reverse '(1 2 3 4))                 ; => (4 3 2 1)

;; Building a list with cons
(cons 1 (cons 2 (cons 3 '())))       ; => (1 2 3)

;; Nested lists
(list '(1 2) '(3 4) '(5 6))          ; => ((1 2) (3 4) (5 6))
(first (first '((a b) (c d))))       ; => a
