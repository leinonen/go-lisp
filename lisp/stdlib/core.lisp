;; Core standard library functions implemented in Lisp
;; These complement the Go core primitives

;; Note: defn is implemented as a special form in the core evaluator

;; Note: cond is implemented as a special form in the core evaluator

;; Logical operations
(def not (fn [x] (if x nil true)))

;; Conditional helpers  
(def when
  (fn [condition body]
    (if condition body nil)))

(def unless
  (fn [condition body]
    (if condition nil body)))

;; Collection operations that complement core functions
;; Note: count, empty?, nth, conj are already in core

;; Length alias for count (common in self-hosting compiler)
(def length count)

;; Hash-map mutation (for self-hosting compiler)
;; Note: This is not truly mutable, but works with reassignment
(def hash-map-put 
  (fn [map key value]
    (assoc map key value)))

;; Second, third helpers
(def second (fn [coll] (first (rest coll))))
(def third (fn [coll] (first (rest (rest coll)))))

;; Map function
(def map
  (fn [f coll]
    (if (empty? coll)
        nil
        (cons (f (first coll)) (map f (rest coll))))))

;; Filter function
(def filter
  (fn [pred coll]
    (if (empty? coll)
        nil
        (if (pred (first coll))
            (cons (first coll) (filter pred (rest coll)))
            (filter pred (rest coll))))))

;; Range function (reverse order for simplicity)
(def range
  (fn [n]
    (if (= n 0)
        nil
        (cons (- n 1) (range (- n 1))))))

;; Helper functions for collection operations
(def concat
  (fn [coll1 coll2]
    (if (empty? coll1)
        coll2
        (cons (first coll1) (concat (rest coll1) coll2)))))

;; Any predicate function
(def any?
  (fn [pred coll]
    (if (empty? coll)
        nil
        (if (pred (first coll))
            true
            (any? pred (rest coll))))))

;; Simple hash-map get (placeholder for now)
(def get
  (fn [m key]
    ;; This is a simplified implementation
    ;; Real hash-map operations would need core support
    nil))

;; Simple hash-map assoc (placeholder for now)
(def assoc
  (fn [m key val]
    ;; This is a simplified implementation
    ;; Real hash-map operations would need core support
    m))

;; Reduce function (simplified)
(def reduce
  (fn [f init coll]
    (if (empty? coll)
        init
        (reduce f (f init (first coll)) (rest coll)))))

;; Apply function - apply function to collection as arguments
(def apply
  (fn [f coll]
    (if (empty? coll)
        (f)
        (if (= (count coll) 1)
            (f (first coll))
            (if (= (count coll) 2)
                (f (first coll) (nth coll 1))
                (if (= (count coll) 3)
                    (f (first coll) (nth coll 1) (nth coll 2))
                    (f (first coll) (nth coll 1) (nth coll 2) (nth coll 3))))))))

;; Group-by function - group collection by key function (simplified)
(def group-by
  (fn [f coll]
    ;; Simplified implementation - returns list of (key value) pairs
    (reduce (fn [acc item]
              (let [key (f item)]
                (cons (list key item) acc)))
            nil
            coll)))

;; Improved map function with two collections support
(def map2
  (fn [f coll1 coll2]
    (if (empty? coll1)
        nil
        (if (empty? coll2)
            nil
            (cons (f (first coll1) (first coll2))
                  (map2 f (rest coll1) (rest coll2)))))))

;; Join strings with separator
(def join
  (fn [sep coll]
    (if (empty? coll)
        ""
        (reduce (fn [acc x] (str acc sep (str x)))
                (str (first coll))
                (rest coll)))))

;; Enhanced sort function
(def sort
  (fn [coll]
    (if (empty? coll)
        nil
        (if (= (count coll) 1)
            coll
            (let [pivot (first coll)
                  rest-coll (rest coll)
                  smaller (filter (fn [x] (< x pivot)) rest-coll)
                  greater (filter (fn [x] (>= x pivot)) rest-coll)]
              (concat (sort smaller) (cons pivot (sort greater))))))))