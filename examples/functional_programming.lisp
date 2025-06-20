;; Functional programming examples in GoLisp

;; Map - apply a function to each element
(map (fn [x] (* x x)) '(1 2 3 4 5))              ; => (1 4 9 16 25)
(map (fn [x] (+ x 10)) '(1 2 3))                 ; => (11 12 13)
(map string-upper '("hello" "world"))            ; => ("HELLO" "WORLD")

;; ;; Filter - keep elements that satisfy a predicate
(filter (fn [x] (> x 3)) '(1 2 3 4 5 6))         ; => (4 5 6)
(filter (fn [x] (= (% x 2) 0)) '(1 2 3 4 5 6))   ; => (2 4 6) - even numbers
(filter (fn [s] (> (string-length s) 3)) '("hi" "hello" "world" "a"))  ; => ("hello" "world")

;; ;; Reduce - combine all elements using a function
(reduce + 0 '(1 2 3 4 5))                        ; => 15
(reduce * 1 '(1 2 3 4))                          ; => 24
(reduce max 0 '(3 7 2 9 1))                      ; => 9

;; ;; Combining functional operations
;; ;; Find sum of squares of even numbers
(reduce + 0 
  (map (fn [x] (* x x))
    (filter (fn [x] (= (% x 2) 0)) '(1 2 3 4 5 6))))  ; => 56

;; ;; Count words longer than 4 characters
(reduce + 0
  (map (fn [s] (if (> (string-length s) 4) 1 0))
    '("cat" "hello" "world" "hi" "programming")))      ; => 3
