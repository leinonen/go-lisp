;; Core standard library functions implemented in Lisp
;; These complement the Go core primitives

;; Note: defn is implemented as a special form in the core evaluator

;; Note: cond is implemented as a special form in the core evaluator

;; Logical operations
(def not (fn [x] (if x nil true)))

;; Conditional helpers  
(defmacro when [condition & body]
  (list 'if condition (cons 'do body) nil))

(defmacro unless [condition & body]
  (list 'if condition nil (cons 'do body)))

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

;; Reduce function (simplified)
(def reduce
  (fn [f init coll]
    (if (empty? coll)
        init
        (reduce f (f init (first coll)) (rest coll)))))

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