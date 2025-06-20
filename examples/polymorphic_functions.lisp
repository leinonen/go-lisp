;; Polymorphic Functions Examples in Go Lisp
;; These examples demonstrate how functions work across different data types

;; ========== SEQUENCE FUNCTIONS ==========

;; first - works on lists, vectors, strings
(first '(1 2 3))              ; => 1
(first ["a" "b" "c"])         ; => "a"  
(first "hello")               ; => "h"
(first nil)                   ; => nil

;; rest - works on lists, vectors, strings  
(rest '(1 2 3))               ; => (2 3)
(rest ["a" "b" "c"])          ; => ("b" "c")
(rest "hello")                ; => ("e" "l" "l" "o")
(rest nil)                    ; => ()

;; last - works on lists, vectors, strings
(last '(1 2 3))               ; => 3
(last ["a" "b" "c"])          ; => "c"
(last "hello")                ; => "o"
(last nil)                    ; => nil

;; nth - works on lists, vectors, strings
(nth '("a" "b" "c" "d") 2)    ; => "c"
(nth [10 20 30 40] 1)         ; => 20
(nth "hello" 1)               ; => "e"

;; second - works on lists, vectors, strings
(second '(1 2 3))             ; => 2
(second ["a" "b" "c"])        ; => "b"
(second "hello")              ; => "e"

;; count - works on all collections (already polymorphic)
(count '(1 2 3))              ; => 3
(count ["a" "b" "c"])         ; => 3
(count {:a 1 :b 2})           ; => 2
(count "hello")               ; => 5
(count nil)                   ; => 0

;; empty? - works on all collections
(empty? '())                  ; => true
(empty? [])                   ; => true
(empty? {})                   ; => true
(empty? "")                   ; => true
(empty? nil)                  ; => true

;; ========== POLYMORPHIC GET AND CONTAINS ==========

;; get - works on hashmaps, vectors, lists, strings
(def my-map {:name "Alice" :age 30})
(get my-map :name)            ; => "Alice"
(get my-map :city)            ; => nil

(def my-vec [10 20 30 40])
(get my-vec 0)                ; => 10
(get my-vec 2)                ; => 30
(get my-vec 10)               ; => nil

(def my-list '("a" "b" "c" "d"))
(get my-list 1)               ; => "b"
(get my-list 3)               ; => "d"

(get "hello" 1)               ; => "e"
(get "hello" 4)               ; => "o"

;; contains? - works on hashmaps, vectors, lists, strings
(contains? my-map :name)      ; => true
(contains? my-map :city)      ; => false

(contains? my-vec 2)          ; => true (index 2 exists)
(contains? my-vec 10)         ; => false (index 10 doesn't exist)

(contains? "hello" 2)         ; => true (index 2 exists)
(contains? "hello" "ell")     ; => true (substring exists)
(contains? "hello" "xyz")     ; => false

;; ========== TRANSFORMATION FUNCTIONS ==========

;; take - take first n elements from any sequence
(take 2 '(1 2 3 4 5))         ; => (1 2)
(take 3 ["a" "b" "c" "d" "e"]) ; => ("a" "b" "c")
(take 2 "hello")              ; => ("h" "e")

;; drop - drop first n elements from any sequence
(drop 2 '(1 2 3 4 5))         ; => (3 4 5)
(drop 1 ["a" "b" "c" "d"])    ; => ("b" "c" "d")
(drop 2 "hello")              ; => ("l" "l" "o")

;; reverse - reverse any sequence
(reverse '(1 2 3))            ; => (3 2 1)
(reverse ["a" "b" "c"])       ; => ["c" "b" "a"]
(reverse "hello")             ; => "olleh"

;; seq - convert to sequence (list)
(seq [1 2 3])                 ; => (1 2 3)
(seq "abc")                   ; => ("a" "b" "c")
(seq nil)                     ; => nil

;; into - merge collections
(into [] '(1 2 3))            ; => [1 2 3]
(into '() ["a" "b" "c"])      ; => ("a" "b" "c")
(into [1 2] [3 4])            ; => [1 2 3 4]

;; distinct - remove duplicates
(distinct '(1 2 2 3 3 3))     ; => (1 2 3)
(distinct ["a" "b" "a" "c" "b"]) ; => ("a" "b" "c")
(distinct "hello")            ; => ("h" "e" "l" "o")

;; sort - sort any sequence
(sort '(3 1 4 1 5))           ; => (1 1 3 4 5)
(sort ["c" "a" "b"])          ; => ("a" "b" "c")
(sort "hello")                ; => ("e" "h" "l" "l" "o")

;; ========== PREDICATE FUNCTIONS ==========

;; Type checking predicates
(seq? '(1 2 3))               ; => true
(seq? [1 2 3])                ; => true
(seq? "abc")                  ; => true
(seq? {:a 1})                 ; => false

(coll? '(1 2 3))              ; => true
(coll? [1 2 3])               ; => true
(coll? {:a 1})                ; => true
(coll? "abc")                 ; => false

(sequential? '(1 2 3))        ; => true
(sequential? [1 2 3])         ; => true
(sequential? {:a 1})          ; => false

(indexed? [1 2 3])            ; => true
(indexed? "abc")              ; => true
(indexed? '(1 2 3))           ; => false

;; ========== UTILITY FUNCTIONS ==========

;; identity - return argument unchanged
(identity 42)                 ; => 42
(identity "hello")            ; => "hello"
(map identity '(1 2 3))       ; => (1 2 3)

;; ========== PRACTICAL EXAMPLES ==========

;; Data transformation pipeline
(def data [1 2 3 4 5 6 7 8 9 10])
(->> data
     (filter (fn [x] (> x 5)))
     (map (fn [x] (* x x)))
     (take 3)
     (reverse))               ; => (100 81 64)

;; Working with mixed data types
(def collections 
  ['(1 2 3) [4 5 6] "789"])

(map count collections)        ; => (3 3 3)
(map first collections)        ; => (1 4 "7")
(map last collections)         ; => (3 6 "9")

;; Polymorphic get with different data structures
(def mixed-data {:users [{"name" "Alice"} {"name" "Bob"}]
                 :count 2})

(get mixed-data :count)                    ; => 2
(get (get mixed-data :users) 0)           ; => {"name" "Alice"}
(get (get (get mixed-data :users) 0) "name") ; => "Alice"

;; Using contains? for validation
(defn valid-index? [coll idx]
  (contains? coll idx))

(valid-index? [1 2 3] 1)      ; => true
(valid-index? [1 2 3] 5)      ; => false
(valid-index? "hello" 2)      ; => true
