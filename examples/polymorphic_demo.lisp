;; Polymorphic Functions Demo
;; This example demonstrates how the same functions work across different data types

(println "=== POLYMORPHIC FUNCTIONS DEMO ===")
(println)

;; Define sample data of different types
(def my-list '(1 2 3 4 5))
(def my-vector [10 20 30 40 50])
(def my-string "hello")
(def my-hashmap {:a 1 :b 2 :c 3})

(println "Sample data:")
(println "List:" my-list)
(println "Vector:" my-vector)
(println "String:" my-string)
(println "Hashmap:" my-hashmap)
(println)

;; Sequence functions work on lists, vectors, strings
(println "=== SEQUENCE FUNCTIONS ===")

(println "first function:")
(println "  (first list):" (first my-list))
(println "  (first vector):" (first my-vector))
(println "  (first string):" (first my-string))
(println)

(println "rest function:")
(println "  (rest list):" (rest my-list))
(println "  (rest vector):" (rest my-vector))
(println "  (rest string):" (rest my-string))
(println)

(println "last function:")
(println "  (last list):" (last my-list))
(println "  (last vector):" (last my-vector))
(println "  (last string):" (last my-string))
(println)

(println "nth function:")
(println "  (nth list 2):" (nth my-list 2))
(println "  (nth vector 2):" (nth my-vector 2))
(println "  (nth string 2):" (nth my-string 2))
(println)

(println "second function:")
(println "  (second list):" (second my-list))
(println "  (second vector):" (second my-vector))
(println "  (second string):" (second my-string))
(println)

;; Collection functions work on all collection types
(println "=== COLLECTION FUNCTIONS ===")

(println "count function:")
(println "  (count list):" (count my-list))
(println "  (count vector):" (count my-vector))
(println "  (count string):" (count my-string))
(println "  (count hashmap):" (count my-hashmap))
(println)

(println "empty? function:")
(println "  (empty? list):" (empty? my-list))
(println "  (empty? vector):" (empty? my-vector))
(println "  (empty? string):" (empty? my-string))
(println "  (empty? hashmap):" (empty? my-hashmap))
(println "  (empty? []):" (empty? []))
(println)

;; Access functions work on different collection types
(println "=== ACCESS FUNCTIONS ===")

(println "get function:")
(println "  (get list 1):" (get my-list 1))
(println "  (get vector 1):" (get my-vector 1))
(println "  (get string 1):" (get my-string 1))
(println "  (get hashmap :b):" (get my-hashmap :b))
(println)

(println "contains? function:")
(println "  (contains? list 2):" (contains? my-list 2))     ; index exists
(println "  (contains? vector 3):" (contains? my-vector 3)) ; index exists
(println "  (contains? string 2):" (contains? my-string 2)) ; index exists
(println "  (contains? hashmap :b):" (contains? my-hashmap :b)) ; key exists
(println)

;; Transformation functions
(println "=== TRANSFORMATION FUNCTIONS ===")

(println "reverse function:")
(println "  (reverse list):" (reverse my-list))
(println "  (reverse vector):" (reverse my-vector))
(println "  (reverse string):" (reverse my-string))
(println)

(println "take function:")
(println "  (take 3 list):" (take 3 my-list))
(println "  (take 3 vector):" (take 3 my-vector))
(println "  (take 3 string):" (take 3 my-string))
(println)

(println "drop function:")
(println "  (drop 2 list):" (drop 2 my-list))
(println "  (drop 2 vector):" (drop 2 my-vector))
(println "  (drop 2 string):" (drop 2 my-string))
(println)

;; Type predicates
(println "=== TYPE PREDICATES ===")

(println "seq? (sequences: lists, vectors, strings):")
(println "  (seq? list):" (seq? my-list))
(println "  (seq? vector):" (seq? my-vector))
(println "  (seq? string):" (seq? my-string))
(println "  (seq? hashmap):" (seq? my-hashmap))
(println)

(println "coll? (collections: lists, vectors, hashmaps):")
(println "  (coll? list):" (coll? my-list))
(println "  (coll? vector):" (coll? my-vector))
(println "  (coll? string):" (coll? my-string))
(println "  (coll? hashmap):" (coll? my-hashmap))
(println)

(println "indexed? (supports index access: vectors, strings):")
(println "  (indexed? list):" (indexed? my-list))
(println "  (indexed? vector):" (indexed? my-vector))
(println "  (indexed? string):" (indexed? my-string))
(println "  (indexed? hashmap):" (indexed? my-hashmap))
(println)

(println "=== POLYMORPHIC DEMO COMPLETE ===")
