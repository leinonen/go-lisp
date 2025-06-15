# HTTP and JSON Functions

This document describes the HTTP request and JSON processing functions available in Go Lisp.

## HTTP Functions

### `http-get`
Performs an HTTP GET request.

**Syntax:**
```lisp
(http-get url)
```

**Parameters:**
- `url`: String - The URL to request

**Returns:**
A hash map with the following keys:
- `:status`: Number - HTTP status code (e.g., 200, 404)
- `:status-text`: String - HTTP status text (e.g., "200 OK")
- `:body`: String - Response body content
- `:headers`: Hash map - Response headers

**Examples:**
```lisp
(def response (http-get "https://api.example.com/users"))
(def status (:status response))
(def body (:body response))

; Alternative syntax using hash-map-get
(def status (hash-map-get response ":status"))
(def body (hash-map-get response ":body"))
```

### `http-post`
Performs an HTTP POST request.

**Syntax:**
```lisp
(http-post url body)
(http-post url body headers)
```

**Parameters:**
- `url`: String - The URL to post to
- `body`: String - Request body (typically JSON)
- `headers`: Hash map (optional) - Custom headers

**Returns:**
Same format as `http-get`.

**Examples:**
```lisp
; Simple POST
(def data (json-stringify (hash-map "name" "Alice" "age" 30)))
(def response (http-post "https://api.example.com/users" data))
(def status (:status response))

; POST with custom headers
(def headers (hash-map "Authorization" "Bearer token123"))
(def response (http-post "https://api.example.com/users" data headers))
(def body (:body response))
```

### `http-put`
Performs an HTTP PUT request.

**Syntax:**
```lisp
(http-put url body)
(http-put url body headers)
```

**Parameters:**
- `url`: String - The URL to put to
- `body`: String - Request body
- `headers`: Hash map (optional) - Custom headers

**Returns:**
Same format as `http-get`.

**Examples:**
```lisp
(def data (json-stringify (hash-map "id" 1 "name" "Updated")))
(def response (http-put "https://api.example.com/users/1" data))
(def status (:status response))
```

### `http-delete`
Performs an HTTP DELETE request.

**Syntax:**
```lisp
(http-delete url)
(http-delete url headers)
```

**Parameters:**
- `url`: String - The URL to delete
- `headers`: Hash map (optional) - Custom headers

**Returns:**
Same format as `http-get`.

**Examples:**
```lisp
(def response (http-delete "https://api.example.com/users/1"))
(def status (:status response))

; DELETE with authentication
(def auth-headers (hash-map "Authorization" "Bearer token123"))
(def response (http-delete "https://api.example.com/users/1" auth-headers))
(def success (= (:status response) 204))
```

## HTTP Response Access

All HTTP functions return hash maps with keyword keys, allowing for convenient access using keyword syntax:

### Keyword Function Syntax
```lisp
(def response (http-get "https://api.example.com/data"))

; Use keywords as functions (recommended)
(:status response)      ; Returns status code
(:status-text response) ; Returns status text
(:body response)        ; Returns response body
(:headers response)     ; Returns headers hash map
```

### Traditional Hash Map Access
```lisp
; Alternative using hash-map-get
(hash-map-get response ":status")      ; Returns status code
(hash-map-get response ":status-text") ; Returns status text
(hash-map-get response ":body")        ; Returns response body
(hash-map-get response ":headers")     ; Returns headers hash map
```

### Working with Headers
Response headers are stored as a hash map with string keys:
```lisp
(def response (http-get "https://api.example.com/data"))
(def headers (:headers response))
(def content-type (hash-map-get headers "Content-Type"))
```

### Complete Example
```lisp
(def response (http-get "https://httpbin.org/json"))
(println! "Status:" (:status response))
(println! "Content-Type:" (hash-map-get (:headers response) "Content-Type"))

; Parse JSON body if successful
(if (= (:status response) 200)
    (def parsed-data (json-parse (:body response)))
    (println! "Request failed with status:" (:status response)))
```

## JSON Functions

### `json-parse`
Converts a JSON string to Lisp data structures.

**Syntax:**
```lisp
(json-parse json-string)
```

**Parameters:**
- `json-string`: String - Valid JSON text

**Returns:**
Lisp value corresponding to the JSON:
- JSON objects → Hash maps
- JSON arrays → Lists  
- JSON strings → Strings
- JSON numbers → Numbers
- JSON booleans → Booleans
- JSON null → nil

**Examples:**
```lisp
(def data (json-parse "{\"name\": \"Alice\", \"age\": 30}"))
; Returns: {name: Alice, age: 30}

(def array (json-parse "[1, 2, 3]"))
; Returns: (1 2 3)
```

### `json-stringify`
Converts Lisp data structures to JSON string.

**Syntax:**
```lisp
(json-stringify value)
```

**Parameters:**
- `value`: Any Lisp value

**Returns:**
String containing compact JSON representation.

**Examples:**
```lisp
(def data (hash-map "name" "Alice" "age" 30))
(def json (json-stringify data))
; Returns: "{\"age\":30,\"name\":\"Alice\"}"
```

### `json-stringify-pretty`
Converts Lisp data structures to pretty-printed JSON string.

**Syntax:**
```lisp
(json-stringify-pretty value)
```

**Parameters:**
- `value`: Any Lisp value

**Returns:**
String containing formatted JSON with indentation and newlines.

**Examples:**
```lisp
(def data (hash-map "name" "Alice" "age" 30))
(def pretty (json-stringify-pretty data))
; Returns:
; {
;   "age": 30,
;   "name": "Alice"
; }
```

### `json-path`
Extracts a value from JSON using simple path notation.

**Syntax:**
```lisp
(json-path json-string path)
```

**Parameters:**
- `json-string`: String - Valid JSON text
- `path`: String - Dot-separated path (e.g., "user.name" or "items.0.title")

**Returns:**
The value at the specified path, or error if path doesn't exist.

**Examples:**
```lisp
(def json "{\"user\": {\"name\": \"Alice\"}, \"items\": [\"a\", \"b\"]}")
(def name (json-path json "user.name"))        ; Returns: Alice
(def first-item (json-path json "items.0"))    ; Returns: a
```

## Data Type Conversions

### Lisp to JSON
- Hash maps → JSON objects
- Lists → JSON arrays
- Strings → JSON strings
- Numbers → JSON numbers
- Booleans → JSON booleans
- nil → JSON null
- Keywords → JSON strings (with colon prefix)

### JSON to Lisp
- JSON objects → Hash maps
- JSON arrays → Lists
- JSON strings → Strings
- JSON numbers → Numbers
- JSON booleans → Booleans
- JSON null → nil

## Practical Examples

### REST API Client
```lisp
(defn get-user [id]
  (def url (string-concat "https://api.example.com/users/" (number->string id)))
  (def response (http-get url))
  (if (= (:status response) 200)
      (json-parse (:body response))
      nil))

(defn create-user [name email]
  (def user-data (hash-map "name" name "email" email))
  (def json-data (json-stringify user-data))
  (def response (http-post "https://api.example.com/users" json-data))
  (json-parse (:body response)))
```

### JSON Data Processing
```lisp
(def api-response "{\"users\": [{\"id\": 1, \"name\": \"Alice\"}, {\"id\": 2, \"name\": \"Bob\"}]}")
(def data (json-parse api-response))
(def users (hash-map-get data "users"))
(def names (map (fn [user] (hash-map-get user "name")) users))
; Returns: ("Alice" "Bob")
```

### Configuration Management
```lisp
(def config (hash-map
  "database" (hash-map "host" "localhost" "port" 5432)
  "api" (hash-map "timeout" 30 "retries" 3)))

(def config-json (json-stringify-pretty config))
; Save to file or send via HTTP

(def parsed-config (json-parse config-json))
(def db-host (json-path config-json "database.host"))
```

## Error Handling

All HTTP and JSON functions can throw errors:

- Invalid URLs, network timeouts, or connection failures
- Invalid JSON syntax
- Type mismatches (e.g., non-string URLs)
- Invalid JSON paths

Always check HTTP status codes and handle potential parsing errors in production code.
