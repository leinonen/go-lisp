;;; JSON Examples
;;; Demonstrates JSON parsing and stringification

;;; JSON Examples
;;; Demonstrates JSON parsing and stringification

;; Create JSON from Lisp data structures (easier than escaping)
(def user-data (hash-map 
  "name" "Alice"
  "age" 30
  "active" true))

(def user-json (json-stringify user-data))
(println! "User as JSON:" user-json)

;; Parse it back
(def parsed-data (json-parse user-json))
(println! "Parsed JSON:" parsed-data)

;; Accessing parsed JSON data
(def name (hash-map-get parsed-data "name"))
(def age (hash-map-get parsed-data "age"))
(println! "Name:" name)
(println! "Age:" age)

;; JSON arrays from Lisp lists
(def array-data (list 1 2 3 "four" true))
(def array-json (json-stringify array-data))
(println! "Array JSON:" array-json)

(def parsed-array (json-parse array-json))
(println! "Parsed array:" parsed-array)
(println! "First element:" (first parsed-array))
(println! "Length:" (length parsed-array))

;; Nested JSON objects - create from Lisp structures
(def nested-data (hash-map
  "user" (hash-map
    "personal" (hash-map
      "name" "Bob"
      "age" 25)
    "preferences" (hash-map
      "theme" "dark"
      "notifications" true))
  "posts" (list
    (hash-map "id" 1 "title" "First Post")
    (hash-map "id" 2 "title" "Second Post"))))

(def nested-json (json-stringify-pretty nested-data))
(println! "Nested JSON:")
(println! nested-json)

;; Using json-path to extract specific values
(def user-name (json-path nested-json "user.personal.name"))
(def user-age (json-path nested-json "user.personal.age"))
(def theme (json-path nested-json "user.preferences.theme"))
(def first-post-title (json-path nested-json "posts.0.title"))
(def second-post-id (json-path nested-json "posts.1.id"))

(println! "User name:" user-name)
(println! "User age:" user-age)
(println! "Theme:" theme)
(println! "First post title:" first-post-title)
(println! "Second post ID:" second-post-id)

;; Creating JSON from Lisp data structures
(def user-data (hash-map 
  "name" "Charlie"
  "age" 28
  "email" "charlie@example.com"
  "active" true
  "tags" (list "developer" "golang" "lisp")))

(def user-json (json-stringify user-data))
(println! "User as JSON:" user-json)

;; Pretty printing JSON
(def pretty-json (json-stringify-pretty user-data))
(println! "Pretty JSON:")
(println! pretty-json)

;; Working with lists and nested structures
(def complex-data (hash-map
  :users (list
    (hash-map :id 1 :name "Alice" :role "admin")
    (hash-map :id 2 :name "Bob" :role "user")
    (hash-map :id 3 :name "Charlie" :role "moderator"))
  :metadata (hash-map 
    :total 3
    :page 1
    :created-at "2024-01-15")))

(def complex-json (json-stringify-pretty complex-data))
(println! "Complex data as pretty JSON:")
(println! complex-json)

;; Parsing it back
(def parsed-back (json-parse complex-json))
(println! "Parsed back:" parsed-back)

;; Round-trip verification
(def users (hash-map-get parsed-back ":users"))
(def first-user (first users))
(def first-user-name (hash-map-get first-user ":name"))
(println! "First user name from round-trip:" first-user-name)

;; Working with API responses
(defn process-api-response [json-response]
  (def data (json-parse json-response))
  (def status (hash-map-get data "status"))
  (def message (hash-map-get data "message"))
  (if (= status "success")
      (hash-map :result "ok" :data (hash-map-get data "data"))
      (hash-map :result "error" :message message)))

(def api-response-1 `{"status": "success", "data": {"id": 123, "name": "Test"}, "message": "OK"}`)
(def api-response-2 `{"status": "error", "message": "Not found", "data": null}`)

(def processed-1 (process-api-response api-response-1))
(def processed-2 (process-api-response api-response-2))

(println! "Processed success response:" processed-1)
(println! "Processed error response:" processed-2)

;; JSON transformation example
(defn transform-user-data [json-string]
  (def data (json-parse json-string))
  (def transformed (hash-map
    :full-name (hash-map-get data "name")
    :years-old (hash-map-get data "age")
    :is-active (hash-map-get data "active")
    :contact-email (hash-map-get data "email")))
  (json-stringify-pretty transformed))

(def original-user `{"name": "David", "age": 35, "active": true, "email": "david@example.com"}`)
(def transformed-user (transform-user-data original-user))
(println! "Transformed user data:")
(println! transformed-user)

;; Handling arrays of objects
(def employees-json `[
  {"id": 1, "name": "Emma", "department": "Engineering", "salary": 75000},
  {"id": 2, "name": "Frank", "department": "Marketing", "salary": 65000},
  {"id": 3, "name": "Grace", "department": "Engineering", "salary": 80000}
]`)

(def employees (json-parse employees-json))
(println! "Number of employees:" (length employees))

;; Filter engineers
(def engineers (filter 
  (fn [employee] 
    (= (hash-map-get employee "department") "Engineering"))
  employees))

(println! "Engineers:" (json-stringify-pretty engineers))

;; Calculate average engineering salary
(def engineer-salaries (map 
  (fn [engineer] (hash-map-get engineer "salary"))
  engineers))

(def avg-salary (/ (reduce + 0 engineer-salaries) (length engineer-salaries)))
(println! "Average engineering salary:" avg-salary)
