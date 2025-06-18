;; Arithmetic and math examples in Go Lisp

;; Basic arithmetic
(+ 1 2 3 4 5)                                     ; => 15
(- 10 3 2)                                        ; => 5
(- 5)                                             ; => -5 (negation)
(* 2 3 4)                                         ; => 24
(/ 20 4)                                          ; => 5
(/ 22 7)                                          ; => 3.142857...
(% 17 5)                                          ; => 2 (modulo)

;; Comparison operations
(= 5 5)                                           ; => true
(= 5 6)                                           ; => false
(< 3 5)                                           ; => true
(<= 5 5)                                          ; => true
(> 7 3)                                           ; => true
(>= 5 5)                                          ; => true

;; Math functions
(abs -5)                                          ; => 5
(max 1 5 3 9 2)                                   ; => 9
(min 1 5 3 9 2)                                   ; => 1
(sqrt 16)                                         ; => 4
(pow 2 3)                                         ; => 8
(exp 1)                                           ; => 2.718...
(log 10)                                          ; => 2.302...
(sin 0)                                           ; => 0
(cos 0)                                           ; => 1
(tan 0)                                           ; => 0

;; Rounding functions
(floor 3.7)                                       ; => 3
(ceil 3.2)                                        ; => 4
(round 3.6)                                       ; => 4

;; Mathematical calculations
;; Calculate area of a circle
(defn circle-area [radius]
  (* 3.14159 (pow radius 2)))

(circle-area 5)                                   ; => 78.53975

;; Calculate compound interest
(defn compound-interest [principal rate time compounds-per-year]
  (* principal (pow (+ 1 (/ rate compounds-per-year)) (* compounds-per-year time))))

(compound-interest 1000 0.05 2 12)               ; => 1104.89...

;; Fibonacci with arithmetic
(defn fib [n]
  (if (<= n 1)
    n
    (+ (fib (- n 1)) (fib (- n 2)))))

(map fib '(0 1 2 3 4 5 6 7 8))                   ; => (0 1 1 2 3 5 8 13 21)
