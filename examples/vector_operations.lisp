;; Vector Operations Examples in go-lisp
;; This file demonstrates vector functionality

;; Basic vector creation
(println "=== Vector Creation ===")
(println "Vector literal:" [1 2 3 4 5])
(println "Using vector function:" (vector 1 2 3 4 5))
(println "Converting list to vector:" (vec (list 1 2 3 4 5)))

;; Vector predicates
(println "\n=== Vector Predicates ===")
(println "Is [1 2 3] a vector?" (vector? [1 2 3]))
(println "Is (list 1 2 3) a vector?" (vector? (list 1 2 3)))

;; Vector operations
(println "\n=== Vector Operations ===")
(println "Count elements in [1 2 3 4 5]:" (count [1 2 3 4 5]))
(println "Get 2nd element (index 1) of [10 20 30]:" (nth [10 20 30] 1))
(println "Conjoin element 4 to [1 2 3]:" (conj [1 2 3] 4))
(println "Conjoin multiple elements:" (conj [1 2] 3 4 5))

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
(println "Access nested element:" (nth (nth [[1 2] [3 4] [5 6]] 1) 0))
