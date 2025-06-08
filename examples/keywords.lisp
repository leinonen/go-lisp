;;; Keywords Examples
;;; Demonstrates the keyword data type and its use with hash maps

;; Basic keyword usage - keywords are self-evaluating
(define example-keywords (list :name :first-name :user-id))
example-keywords

;; Keywords are self-evaluating
(define my-keyword :status)
my-keyword

;; Keywords in lists
(list :name :age :email)

;; Hash maps with keyword keys (more idiomatic)
(define person (hash-map
  :name "Alice"
  :age 30
  :city "Boston"))

person

;; Accessing with keyword keys
(hash-map-get person :name)
(hash-map-get person :age)
(hash-map-get person :missing)  ; Returns nil

;; Adding and updating with keywords
(define updated-person (hash-map-put person :email "alice@example.com"))
updated-person

;; Mixed string and keyword keys
(define mixed-map (hash-map
  "string-key" "string-value"
  :keyword-key "keyword-value"))

mixed-map

;; Keywords vs strings in hash maps
(hash-map-get mixed-map "string-key")    ; Regular string key
(hash-map-get mixed-map :keyword-key)    ; Keyword key

;; Querying hash maps with keywords
(hash-map-contains? person :name)
(hash-map-contains? person :missing)

;; Removing with keywords
(define minimal-person (hash-map-remove person :city))
minimal-person

;; Keywords are great for configuration
(define config (hash-map
  :debug #t
  :port 8080
  :host "localhost"
  :database-url "localhost:5432"
  :max-connections 100))

config

;; Using keywords for function parameters (more readable)
(define create-user (lambda (name email age)
  (hash-map
    :name name
    :email email
    :age age
    :created-at "2024-01-15"
    :active #t)))

(define user1 (create-user "Charlie" "charlie@example.com" 28))
(define user2 (create-user "Diana" "diana@example.com" 32))

user1
user2

;; Helper function for safe access with defaults
(define get-with-default (lambda (map key default)
  (if (hash-map-contains? map key)
      (hash-map-get map key)
      default)))

(get-with-default user1 :phone "no phone")
(get-with-default user1 :name "unknown")

;; Example: Application state with keywords
(define app-state (hash-map
  :current-user nil
  :logged-in #f
  :theme :dark
  :language :en
  :notifications (hash-map
    :email #t
    :push #f
    :desktop #t)))

app-state

;; Accessing nested keyword-based hash maps
(hash-map-get (hash-map-get app-state :notifications) :email)
