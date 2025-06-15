;;; HTTP Examples
;;; Demonstrates HTTP request functionality

;; Basic HTTP GET request
(def response (http-get "https://httpbin.org/get"))
(println! "GET Response:")
(println! (:status response))
(println! (:body response))

;; HTTP GET with JSON parsing
(def json-response (http-get "https://httpbin.org/json"))
(def json-body (:body json-response))
(def parsed-json (json-parse json-body))
(println! "Parsed JSON:")
(println! parsed-json)

;; HTTP POST with JSON data
(def user-data (hash-map :name "Alice" :age 30 :city "Boston"))
(def json-data (json-stringify user-data))
(def post-response (http-post "https://httpbin.org/post" json-data))
(println! "POST Response Status:" (:status post-response))

;; HTTP POST with custom headers
(def custom-headers (hash-map 
  "Content-Type" "application/json"
  "X-API-Key" "your-api-key"
  "User-Agent" "Go-Lisp/1.0"))
(def headers-response (http-post "https://httpbin.org/post" json-data custom-headers))
(println! "POST with headers response:" (:status headers-response))

;; HTTP PUT request
(def updated-data (json-stringify (hash-map :id 1 :name "Updated Name")))
(def put-response (http-put "https://httpbin.org/put" updated-data))
(println! "PUT Response Status:" (:status put-response))

;; HTTP DELETE request
(def delete-response (http-delete "https://httpbin.org/delete"))
(println! "DELETE Response Status:" (:status delete-response))

;; HTTP DELETE with headers
(def auth-headers (hash-map "Authorization" "Bearer your-token"))
(def auth-delete-response (http-delete "https://httpbin.org/delete" auth-headers))
(println! "DELETE with auth response:" (:status auth-delete-response))

;; Working with response headers
(def header-response (http-get "https://httpbin.org/response-headers?Content-Type=application/json&X-Custom=test"))
(def response-headers (:headers header-response))
(println! "Response headers:")
(println! (hash-map-keys response-headers))
(println! "Content-Type:" (hash-map-get response-headers "Content-Type"))

;; Error handling example
(def error-response (http-get "https://httpbin.org/status/404"))
(def error-status (:status error-response))
(if (= error-status 404)
    (println! "Got expected 404 error")
    (println! "Unexpected status:" error-status))

;; Practical example: REST API interaction
(defn create-user [name email]
  (def user-data (hash-map :name name :email email :created-at "2024-01-15"))
  (def json-body (json-stringify user-data))
  (def response (http-post "https://httpbin.org/post" json-body))
  (def response-body (:body response))
  (json-parse response-body))

(def new-user (create-user "Bob Smith" "bob@example.com"))
(println! "Created user:" new-user)

;; Fetching and processing JSON data
(defn get-user-info [user-id]
  (def url (string-concat "https://httpbin.org/json?user_id=" (number->string user-id)))
  (def response (http-get url))
  (def json-body (:body response))
  (json-parse json-body))

(def user-info (get-user-info 123))
(println! "User info:" user-info)
