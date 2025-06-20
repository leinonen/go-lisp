;; List operations examples in Go Lisp
;; Note: Many functions are now polymorphic and work on lists, vectors, and strings

;; Creating lists
(list 1 2 3 4 5)                    ; => (1 2 3 4 5)
(list "hello" "world" "!")           ; => ("hello" "world" "!")

;; Getting elements - polymorphic functions work on lists, vectors, strings
(first '(1 2 3))                     ; => 1
(first ["a" "b" "c"])                ; => "a"  
(first "hello")                      ; => "h"

(rest '(1 2 3))                      ; => (2 3)
(rest ["a" "b" "c"])                 ; => ("b" "c")
(rest "hello")                       ; => ("e" "l" "l" "o")

(last '(1 2 3))                      ; => 3
(last ["a" "b" "c"])                 ; => "c"
(last "hello")                       ; => "o"

(nth '(a b c d) 1)                   ; => b (0-indexed)
(nth ["x" "y" "z"] 2)                ; => "z"
(nth "hello" 1)                      ; => "e"

;; Adding elements (list-specific)
(cons 0 '(1 2 3))                    ; => (0 1 2 3)
(append '(1 2) '(3 4) '(5 6))        ; => (1 2 3 4 5 6)

;; Collection properties - polymorphic functions
(count '(1 2 3 4))                   ; => 4 (polymorphic count)
(count ["a" "b" "c"])                ; => 3
(count "hello")                      ; => 5
(length '(1 2 3 4))                  ; => 4 (list-specific, but count is preferred)

(empty? '())                         ; => true
(empty? [])                          ; => true  
(empty? "")                          ; => true
(empty? '(1))                        ; => false

;; Transformations - polymorphic functions
(reverse '(1 2 3 4))                 ; => (4 3 2 1)
(reverse ["a" "b" "c"])              ; => ["c" "b" "a"]
(reverse "hello")                    ; => "olleh"

;; Building a list with cons
(cons 1 (cons 2 (cons 3 '())))       ; => (1 2 3)

;; Nested lists
(list '(1 2) '(3 4) '(5 6))          ; => ((1 2) (3 4) (5 6))
(first (first '((a b) (c d))))       ; => a
