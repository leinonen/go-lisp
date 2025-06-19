;; Assoc and Dissoc Examples in Go Lisp
;; These functions provide traditional Lisp-style hashmap manipulation

;; Create initial hashmap using literal syntax
(def person {:name "Alice" :age 25 :job "Engineer"})

;; Using assoc to add/update key-value pairs
;; assoc returns a new hashmap (immutable)

;; Add a single key-value pair
(def person-with-city (assoc person :city "Boston"))
person-with-city                                  ; => {"name" Alice "age" 25 "job" Engineer "city" Boston}

;; Add multiple key-value pairs at once
(def full-person (assoc person :city "Boston" :salary 80000 :active true))
full-person                                       ; => {"name" Alice "age" 25 "job" Engineer "city" Boston "salary" 80000 "active" true}

;; Update existing key
(def older-person (assoc person :age 26))
older-person                                      ; => {"name" Alice "age" 26 "job" Engineer}

;; Original hashmap remains unchanged (immutable)
person                                            ; => {"name" Alice "age" 25 "job" Engineer}

;; Using dissoc to remove keys
;; dissoc returns a new hashmap without specified keys

;; Remove a single key
(def person-no-job (dissoc person :job))
person-no-job                                     ; => {"name" Alice "age" 25}

;; Remove multiple keys at once
(def minimal-person (dissoc full-person :salary :active :city))
minimal-person                                    ; => {"name" Alice "age" 25 "job" Engineer}

;; Removing non-existing key doesn't cause error
(def same-person (dissoc person :non-existing-key))
same-person                                       ; => {"name" Alice "age" 25 "job" Engineer}

;; Chaining operations
(def transformed-person 
  (-> person
      (assoc :age 30 :city "New York")
      (dissoc :job)
      (assoc :status "available")))
transformed-person                                ; => {"name" Alice "age" 30 "city" New York "status" available}

;; Working with empty hashmaps
(def empty-map {})
(def populated (assoc empty-map :first "value1" :second "value2"))
populated                                         ; => {"first" value1 "second" value2}

(def back-to-empty (dissoc populated :first :second))
back-to-empty                                     ; => {}

;; Comparison with hash-map-put and hash-map-remove
;; These are equivalent:
(assoc person :city "Boston")                     ; New style
(hash-map-put person :city "Boston")              ; Old style

(dissoc person :job)                              ; New style
(hash-map-remove person :job)                     ; Old style

;; assoc and dissoc are more idiomatic and support multiple operations
;; hash-map-put and hash-map-remove only work with single key-value pairs

;; Practical example: updating user preferences
(def user-prefs {:theme "dark" :notifications true :language "en"})

;; Update multiple preferences at once
(def new-prefs (assoc user-prefs 
                      :theme "light" 
                      :timezone "UTC" 
                      :auto-save true))

;; Remove unwanted preferences
(def clean-prefs (dissoc new-prefs :auto-save :notifications))

clean-prefs                                       ; => {"theme" light "language" en "timezone" UTC}
