;;; JSON Examples
;;; Demonstrates JSON parsing and stringification

(println! "=== JSON Processing Examples ===")

;; Creating JSON from Lisp data structures
(def user-data (hash-map 
  "name" "Alice"
  "age" 30
  "email" "alice@example.com"
  "active" true
  "hobbies" (list "reading" "coding" "music")))

(println! "Original Lisp data:" user-data)

;; Convert to JSON
(def user-json (json-stringify user-data))
(println! "As JSON:" user-json)

;; Pretty print JSON
(def pretty-json (json-stringify-pretty user-data))
(println! "Pretty JSON:")
(println! pretty-json)

;; Parse JSON back to Lisp
(def parsed-back (json-parse user-json))
(println! "Parsed back to Lisp:" parsed-back)

;; Access individual fields
(def name (hash-map-get parsed-back "name"))
(def age (hash-map-get parsed-back "age"))
(def hobbies (hash-map-get parsed-back "hobbies"))

(println! "Name:" name)
(println! "Age:" age)
(println! "Hobbies:" hobbies)
(println! "First hobby:" (first hobbies))

;; Working with nested data
(def company-data (hash-map
  "name" "Tech Corp"
  "employees" (list
    (hash-map "id" 1 "name" "Alice" "department" "Engineering")
    (hash-map "id" 2 "name" "Bob" "department" "Marketing")
    (hash-map "id" 3 "name" "Charlie" "department" "Engineering"))
  "metadata" (hash-map
    "founded" 2020
    "employees_count" 3)))

(def company-json (json-stringify-pretty company-data))
(println! "Company data as JSON:")
(println! company-json)

;; Using json-path for easy access
(def company-name (json-path company-json "name"))
(def first-employee-name (json-path company-json "employees.0.name"))
(def employee-count (json-path company-json "metadata.employees_count"))

(println! "Company name via json-path:" company-name)
(println! "First employee name:" first-employee-name)
(println! "Employee count:" employee-count)

;; Data transformation example
(defn transform-employee [employee]
  (hash-map
    "full_name" (hash-map-get employee "name")
    "dept" (hash-map-get employee "department")
    "employee_id" (hash-map-get employee "id")))

(def employees (hash-map-get company-data "employees"))
(def transformed-employees (map transform-employee employees))
(def transformed-json (json-stringify-pretty transformed-employees))

(println! "Transformed employees:")
(println! transformed-json)

;; Round-trip test
(def round-trip-data (json-parse transformed-json))
(def first-transformed (first round-trip-data))
(def full-name (hash-map-get first-transformed "full_name"))

(println! "Round-trip test - first employee full name:" full-name)

(println! "JSON processing examples completed!")
