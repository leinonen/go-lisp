;;; HTTP Examples
;;; Demonstrates HTTP request functionality

(println! "=== HTTP Request Examples ===")

;; NOTE: These examples use httpbin.org for testing HTTP requests
;; If you don't have internet access, they will fail gracefully

(println! "Testing HTTP GET request...")

;; Simple GET request to httpbin (a testing service)
(def response (http-get "https://httpbin.org/get"))
(def status (hash-map-get response "status"))
(println! "GET request status:" status)

(if (= status 200)
    (do
      (println! "GET request successful!")
      (def body (hash-map-get response "body"))
      (def parsed-body (json-parse body))
      (def url (hash-map-get parsed-body "url"))
      (println! "Requested URL:" url))
    (println! "GET request failed with status:" status))

(println! "Testing HTTP POST request...")

;; POST request with JSON data
(def user-data (hash-map "name" "Alice" "age" 30 "city" "Boston"))
(def json-data (json-stringify user-data))
(def post-response (http-post "https://httpbin.org/post" json-data))
(def post-status (hash-map-get post-response "status"))

(println! "POST request status:" post-status)

(if (= post-status 200)
    (do
      (println! "POST request successful!")
      (def post-body (hash-map-get post-response "body"))
      (def parsed-post (json-parse post-body))
      (def received-data (hash-map-get parsed-post "data"))
      (println! "Server received our data:" received-data))
    (println! "POST request failed with status:" post-status))

(println! "Testing HTTP POST with custom headers...")

;; POST with custom headers
(def headers (hash-map 
  "Content-Type" "application/json"
  "X-API-Key" "test-key-123"
  "User-Agent" "Go-Lisp-Client/1.0"))

(def headers-response (http-post "https://httpbin.org/post" json-data headers))
(def headers-status (hash-map-get headers-response "status"))

(println! "POST with headers status:" headers-status)

(if (= headers-status 200)
    (do
      (println! "POST with headers successful!")
      (def headers-body (hash-map-get headers-response "body"))
      (def parsed-headers (json-parse headers-body))
      (def received-headers (hash-map-get parsed-headers "headers"))
      (def api-key (hash-map-get received-headers "X-Api-Key"))
      (println! "Server received our API key:" api-key))
    (println! "POST with headers failed with status:" headers-status))

(println! "Testing HTTP PUT request...")

;; PUT request
(def updated-data (hash-map "id" 123 "name" "Updated Alice" "status" "active"))
(def put-json (json-stringify updated-data))
(def put-response (http-put "https://httpbin.org/put" put-json))
(def put-status (hash-map-get put-response "status"))

(println! "PUT request status:" put-status)

(println! "Testing HTTP DELETE request...")

;; DELETE request
(def delete-response (http-delete "https://httpbin.org/delete"))
(def delete-status (hash-map-get delete-response "status"))

(println! "DELETE request status:" delete-status)

;; Function to make API calls and handle errors
(defn api-call [method url data headers]
  (def response 
    (if (= method "GET")
        (http-get url)
        (if (= method "POST")
            (if headers
                (http-post url data headers)
                (http-post url data))
            (if (= method "PUT")
                (if headers
                    (http-put url data headers)
                    (http-put url data))
                (if headers
                    (http-delete url headers)
                    (http-delete url))))))
  
  (def status (hash-map-get response "status"))
  (hash-map 
    "success" (and (>= status 200) (< status 300))
    "status" status
    "body" (hash-map-get response "body")
    "headers" (hash-map-get response "headers")))

(println! "Testing reusable API function...")

(def api-result (api-call "GET" "https://httpbin.org/json" nil nil))
(def api-success (hash-map-get api-result "success"))

(if api-success
    (do
      (println! "API call successful!")
      (def api-body (hash-map-get api-result "body"))
      (def parsed-api (json-parse api-body))
      (println! "Received JSON data:" parsed-api))
    (do
      (def api-status (hash-map-get api-result "status"))
      (println! "API call failed with status:" api-status)))

(println! "HTTP examples completed!")
(println! "Note: If you see failures, check your internet connection.")
