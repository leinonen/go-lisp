;; Hash Map Literal Syntax Demo
;; Now you can create hash maps using Clojure-style syntax!

;; Define a person using the new literal syntax
(def person {:name "Alice" :age 25 :job "Engineer" :active true})

;; Access values (same as before)
(hash-map-get person "name")      ; => "Alice"
(hash-map-get person "age")       ; => 25

;; Create nested hash maps easily
(def company
  {:name "TechCorp"
   :location "San Francisco"
   :employees {:alice {:role "Frontend Developer" :salary 80000}
               :bob {:role "Backend Developer" :salary 85000}
               :charlie {:role "Product Manager" :salary 95000}}
   :founded 2010})

;; Access nested data
(hash-map-get company "name")
(hash-map-get (hash-map-get company "employees") "alice")

;; Mix different key types
(def mixed-keys {:keyword-key "value1" "string-key" "value2" :another-key 42})

;; Empty hash map
(def empty-map {})
(hash-map-empty? empty-map)       ; => true

;; Compare old vs new syntax:
;; Old way:
(def old-way (hash-map "name" "John" "age" 30 "city" "Boston"))

;; New way (much cleaner!):
(def new-way {:name "John" :age 30 :city "Boston"})

;; Both create the same result
old-way
new-way
