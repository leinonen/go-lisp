;; Hash map examples in GoLisp

;; Creating hash maps with the new literal syntax (Clojure-style)
{:name "John" :age 30 :city "Boston"}

;; Creating hash maps with function syntax (still supported)
(hash-map "name" "John" "age" 30 "city" "Boston")

;; Creating hash maps with mixed key types
{:name "Alice" "job" "Engineer" :age 25}

;; Getting values
(def person {:name "Alice" :age 25 :job "Engineer"})
(hash-map-get person "name")                      ; => "Alice"
(hash-map-get person "age")                       ; => 25

;; Adding/updating values
(def updated-person (hash-map-put person "age" 26))
(def with-hobby (hash-map-put updated-person "hobby" "reading"))

;; Removing values
(def without-job (hash-map-remove person "job"))

;; Hash map properties
(hash-map-size person)                            ; => 3
(hash-map-empty? person)                          ; => false
(hash-map-empty? {})                              ; => true
(hash-map-contains? person "name")                ; => true
(hash-map-contains? person "salary")              ; => false

;; Getting keys and values
(hash-map-keys person)                            ; => ("name" "age" "job")
(hash-map-values person)                          ; => ("Alice" 25 "Engineer")

;; Nested hash maps with literal syntax
(def company 
  {:name "TechCorp"
   :employees {:alice {:role "dev" :salary 80000}
               :bob {:role "manager" :salary 95000}}})

(hash-map-get (hash-map-get (hash-map-get company "employees") "alice") "role")  ; => "dev"

;; Empty hash map
{}

;; Comparison with the old syntax
;; Old way:
(hash-map "name" "Charlie" "age" 40)

;; New way (much cleaner!):
{:name "Charlie" :age 40}
