;; Control flow and variable binding examples in GoLisp

;; Variable definitions
(def x 10)
(def name "Alice")
(def numbers '(1 2 3 4 5))

;; Local bindings with let
(let [a 5
      b 10
      sum (+ a b)]
  (* sum 2))                                      ; => 30

;; Conditional expressions
(if (> x 5) 
  "greater than 5" 
  "not greater than 5")                          ; => "greater than 5"

;; Multiple conditions with cond
(def grade 85)
(cond 
  (>= grade 90) "A"
  (>= grade 80) "B"
  (>= grade 70) "C"
  (>= grade 60) "D"
  :else "F")                                     ; => "B"

;; Function definitions
(defn square [x]
  (* x x))

(defn factorial [n]
  (if (<= n 1)
    1
    (* n (factorial (- n 1)))))

(defn greet [name]
  (string-concat "Hello, " name "!"))

;; Using functions
(square 5)                                       ; => 25
(factorial 5)                                    ; => 120
(greet "World")                                  ; => "Hello, World!"

;; Anonymous functions (lambdas)
(def add-ten (fn [x] (+ x 10)))
(add-ten 5)                                      ; => 15

;; Higher-order functions
(defn apply-twice [f x]
  (f (f x)))

(apply-twice square 3)                           ; => 81 (3^2)^2

;; Multiple parameter functions
(defn calculate-area [length width]
  (* length width))

(defn full-name [first last]
  (string-concat first " " last))

;; Practical examples
;; Check if a number is prime
(defn prime? [n]
  (cond 
    (< n 2) false
    (= n 2) true
    (= (% n 2) 0) false
    :else (not (reduce (fn [acc x] 
                        (or acc (= (% n x) 0)))
                      false
                      (range 3 (+ (sqrt n) 1) 2)))))

;; Find the maximum value in a list
(defn list-max [lst]
  (if (empty? lst)
    nil
    (reduce (fn [acc x] (if (> x acc) x acc))
            (first lst)
            (rest lst))))

;; Count occurrences of an element
(defn count-occurrences [item lst]
  (reduce (fn [acc x] 
            (if (= x item) 
              (+ acc 1) 
              acc))
          0 
          lst))

(count-occurrences 2 '(1 2 3 2 4 2 5))          ; => 3
