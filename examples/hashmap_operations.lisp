;; Hash map examples in Go Lisp

;; Creating hash maps
(hash-map "name" "John" "age" 30 "city" "Boston")

;; Getting values
(def person (hash-map "name" "Alice" "age" 25 "job" "Engineer"))
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
(hash-map-empty? (hash-map))                      ; => true
(hash-map-contains? person "name")                ; => true
(hash-map-contains? person "salary")              ; => false

;; Getting keys and values
(hash-map-keys person)                            ; => ("name" "age" "job")
(hash-map-values person)                          ; => ("Alice" 25 "Engineer")

;; Nested hash maps
(def company 
  (hash-map 
    "name" "TechCorp"
    "employees" (hash-map
                  "alice" (hash-map "role" "dev" "salary" 80000)
                  "bob" (hash-map "role" "manager" "salary" 95000))))

(hash-map-get (hash-map-get (hash-map-get company "employees") "alice") "role")  ; => "dev"
