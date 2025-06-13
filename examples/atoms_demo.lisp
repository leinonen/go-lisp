;; Atoms - Thread-safe Mutable References
;; This example demonstrates the use of Clojure-style atoms for managing mutable state.

; Create an atom with an initial value
(def counter (atom 0))
(println! "Created counter atom with initial value:" (deref counter))

; Use deref to get the current value
(println! "Current counter value:" (deref counter))

; Use swap! to atomically update the value by applying a function
(swap! counter (fn [x] (+ x 1)))
(println! "After incrementing:" (deref counter))

; swap! can use any function that transforms the current value
(swap! counter (fn [x] (* x 10)))
(println! "After multiplying by 10:" (deref counter))

; Use reset! to set a completely new value
(reset! counter 42)
(println! "After reset to 42:" (deref counter))

; Atoms work with any data type
(def message (atom "Hello"))
(println! "String atom:" (deref message))

(swap! message (fn [s] (string-concat s ", World!")))
(println! "After string concatenation:" (deref message))

; Atoms with lists
(def items (atom (list)))
(println! "Empty list atom:" (deref items))

(swap! items (fn [lst] (cons "first" lst)))
(swap! items (fn [lst] (cons "second" lst)))
(println! "List after adding items:" (deref items))

; Atoms are thread-safe - perfect for shared state
; This is especially useful in concurrent scenarios
(def shared-counter (atom 0))

; Function to increment the shared counter safely
(defn increment-counter []
  (swap! shared-counter (fn [x] (+ x 1))))

; Multiple calls are safe and atomic
(increment-counter)
(increment-counter)
(increment-counter)
(println! "Shared counter after 3 increments:" (deref shared-counter))

(println! "Atoms provide thread-safe, atomic updates to shared state!")
