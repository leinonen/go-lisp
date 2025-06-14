;;; Hash Map Examples
;;; Demonstrates the hash map data structure and operations

;; Basic hash map creation
(def empty-map (hash-map()))
(def person (hash-map "name" "Alice" "age" 30 "city" "Boston"))

;; Display the hash maps
empty-map
person

;; Accessing values
(hash-map-get person "name")
(hash-map-get person "age")
(hash-map-get person "missing")  ; Returns nil

;; Adding and updating values (immutable operations)
(def updated-person (hash-map-put person "email" "alice@example.com"))
updated-person

;; Original hash map is unchanged
person

;; Update existing key
(def renamed-person (hash-map-put person "name" "Alice Smith"))
renamed-person

;; Removing keys
(def minimal-person (hash-map-remove person "city"))
minimal-person

;; Querying hash maps
(hash-map-contains? person "name")
(hash-map-contains? person "missing")
(hash-map-size person)
(hash-map-empty? empty-map)
(hash-map-empty? person)

;; Getting all keys and values
(hash-map-keys person)
(hash-map-values person)

;; Practical example: Configuration management
(def config (hash-map 
  "database-url" "localhost:5432"
  "debug" true
  "max-connections" 100
  "timeout" 30))

(def production-config 
  (hash-map-put 
    (hash-map-put config "debug" false)
    "database-url" "prod-server:5432"))

config
production-config

;; Example: Inventory management
(def inventory (hash-map 
  "apples" 50
  "oranges" 30
  "bananas" 25
  "grapes" 15))

;; Find all fruit names
(hash-map-keys inventory)

;; Update inventory after sale
(def after-sale 
  (hash-map-put 
    (hash-map-put inventory "apples" 45)
    "oranges" 28))

after-sale

;; Example: Nested hash maps
(def user-profile (hash-map
  "personal" (hash-map "name" "Bob" "age" 35)
  "preferences" (hash-map "theme" "dark" "language" "en")
  "stats" (hash-map "login-count" 42 "last-login" "2024-01-15")))

user-profile

;; Access nested values
(hash-map-get (hash-map-get user-profile "personal") "name")
(hash-map-get (hash-map-get user-profile "stats") "login-count")

;; Helper function for safe access with defaults
(def get-with-default (fn [map key default]
  (if (hash-map-contains? map key)
      (hash-map-get map key)
      default)))

(get-with-default person "phone" "no phone")
(get-with-default person "name" "unknown")

;; Example: Building a simple database record
(def create-user (fn [name email age]
  (hash-map 
    "name" name 
    "email" email 
    "age" age
    "created-at" "2024-01-15"
    "active" true)))

(def user1 (create-user "Charlie" "charlie@example.com" 28))
(def user2 (create-user "Diana" "diana@example.com" 32))

user1
user2

;; Example: Hash map as a lookup table
(def color-codes (hash-map
  "red" "falseF0000"
  "green" "#00FF00"
  "blue" "#0000FF"
  "yellow" "falseFFF00"
  "purple" "#800080"))

(def get-color-code (fn [color]
  (get-with-default color-codes color "unknown color")))

(get-color-code "red")
(get-color-code "blue")
(get-color-code "orange")
