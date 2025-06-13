; Advanced Features Demo
; This file demonstrates advanced interpreter capabilities

; Tail Call Optimization - prevents stack overflow
(defn tail-factorial [n acc]
  "Tail-recursive factorial with accumulator"
  (if (= n 0) acc (tail-factorial (- n 1) (* n acc))))

(defn factorial [n] [tail-factorial n 1])

; This can handle very large numbers without stack overflow
(factorial 20)                          ; => 2432902008176640000

; Tail-recursive fibonacci
(defn fib-tail [n a b]
  (if (= n 0) a (fib-tail (- n 1) b (+ a b))))

(defn fibonacci [n] [fib-tail n 0 1])
(fibonacci 30)                          ; => 832040

; Environment inspection and introspection
(def my-var 42)
(defn my-function [x] [+ x 1])

; Check what's in our environment
(env)                                   ; Shows all defined variables and functions

; Check available built-in functions
(length (help))                         ; Count of built-in functions

; Environment introspection - check what's available in current scope
(def my-functions (env))
(length my-functions)                   ; Show number of defined symbols

; Advanced list operations with higher-order functions
(def numbers (list 1 2 3 4 5 6 7 8 9 10))

; Chain operations together
(def evens (filter (lambda [x] (= (% x 2) 0)) numbers))
(def doubled-evens (map (lambda [x] (* x 2)) evens))
(def sum-doubled-evens (reduce (lambda [acc x] (+ acc x)) 0 doubled-evens))

evens                                   ; => (2 4 6 8 10)
doubled-evens                           ; => (4 8 12 16 20)
sum-doubled-evens                       ; => 60

; Big number arithmetic
(def big-num 123456789012345678901234567890)
(% big-num 7)                          ; => 4 (modulo works with big numbers)

; Comments and documentation
; The interpreter ignores everything after semicolons
(+ 1 2)  ; This calculates 1 + 2 = 3

; Error handling demonstration
; Uncomment the line below to see error handling:
; (error "This is a custom error message!")

; File loading capability (load other .lisp files)
; (load "examples/basic_features.lisp")
