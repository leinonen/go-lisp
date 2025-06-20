;; Vector Operations Examples in go-lisp
;; This file demonstrates vector functionality and polymorphic operations

;; Basic vector creation
(println "=== Vector Creation ===")
(println "Vector literal:" [1 2 3 4 5])
(println "Using vector function:" (vector 1 2 3 4 5))
(println "Converting list to vector:" (vec (list 1 2 3 4 5)))

;; Vector predicates
(println "\n=== Vector Predicates ===")
(println "Is [1 2 3] a vector?" (vector? [1 2 3]))
(println "Is (list 1 2 3) a vector?" (vector? (list 1 2 3)))

;; Polymorphic operations - work on vectors, lists, and strings
(println "\n=== Polymorphic Operations ===")
(println "Count elements in [1 2 3 4 5]:" (count [1 2 3 4 5]))
(println "First element of [10 20 30]:" (first [10 20 30]))
(println "Rest of [10 20 30]:" (rest [10 20 30]))  ; Returns list
(println "Last element of [10 20 30]:" (last [10 20 30]))
(println "Get 2nd element (index 1) of [10 20 30]:" (nth [10 20 30] 1))
(println "Second element of [10 20 30]:" (second [10 20 30]))

;; Polymorphic get and contains? functions
(println "\n=== Polymorphic Access ===")
(println "Get index 2 from [a b c d]:" (get ["a" "b" "c" "d"] 2))
(println "Contains index 1 in [a b c]:" (contains? ["a" "b" "c"] 1))
(println "Contains index 5 in [a b c]:" (contains? ["a" "b" "c"] 5))

;; Vector-specific operations
(println "\n=== Vector-Specific Operations ===")
(println "Conjoin element 4 to [1 2 3]:" (conj [1 2 3] 4))
(println "Conjoin multiple elements:" (conj [1 2] 3 4 5))

;; Polymorphic transformations
(println "\n=== Polymorphic Transformations ===")
(println "Reverse [1 2 3 4]:" (reverse [1 2 3 4]))  ; Returns vector
(println "Take 3 from [1 2 3 4 5 6]:" (take 3 [1 2 3 4 5 6]))
(println "Drop 2 from [1 2 3 4 5]:" (drop 2 [1 2 3 4 5]))

;; Cross-type polymorphic examples
(println "\n=== Cross-Type Examples ===")
(println "First of vector:" (first [1 2 3]))     ; => 1
(println "First of list:" (first '(1 2 3)))      ; => 1  
(println "First of string:" (first "hello"))     ; => "h"

(println "Count of vector:" (count [1 2 3]))     ; => 3
(println "Count of list:" (count '(1 2 3)))      ; => 3
(println "Count of string:" (count "hello"))     ; => 5

;; Functional operations on vectors
(println "\n=== Functional Operations ===")
(println "Map double over [1 2 3 4]:" (map (fn [x] (* x 2)) [1 2 3 4]))
(println "Filter even numbers from [1 2 3 4 5 6]:" (filter (fn [x] (= (% x 2) 0)) [1 2 3 4 5 6]))
(println "Sum all elements in [1 2 3 4 5]:" (reduce + 0 [1 2 3 4 5]))

;; Conversion between vectors and lists
(println "\n=== Conversions ===")
(println "Convert vector to list:" (seq [1 2 3]))
(println "Convert list to vector:" (vec (list 1 2 3)))

;; Nested vectors
(println "\n=== Nested Vectors ===")
(println "Nested vector:" [[1 2] [3 4] [5 6]])
(println "Access nested element:" (get (get [[1 2] [3 4] [5 6]] 1) 0))
