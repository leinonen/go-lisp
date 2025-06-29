;; Core standard library functions implemented in Lisp
;; These complement the Go core primitives

;; Note: defn is implemented as a special form in the core evaluator

;; Note: cond is implemented as a special form in the core evaluator

;; Logical operations
(defn not [x] (if x nil true))

;; Conditional helpers  
(defmacro when [condition & body]
  (list 'if condition (cons 'do body) nil))

(defmacro unless [condition & body]
  (list 'if condition nil (cons 'do body)))

(defmacro cond [& clauses]
  (if (empty? clauses)
    nil
    (let [condition (first clauses)
          then-expr (second clauses)
          rest-clauses (rest (rest clauses))]
      (if (= condition :else)
        then-expr
        (list 'if condition then-expr (cons 'cond rest-clauses))))))

;; Collection operations that complement core functions
;; Note: count, empty?, nth, conj are already in core

;; Length alias for count (common in self-hosting compiler)
(def length count)

;; Hash-map mutation (for self-hosting compiler)
;; Note: This is not truly mutable, but works with reassignment
(defn hash-map-put [map key value]
  (assoc map key value))

;; Second, third helpers
(defn second [coll] (first (rest coll)))
(defn third [coll] (first (rest (rest coll))))

;; Map function
(defn map [f coll]
  (if (empty? coll)
      ()
      (cons (f (first coll)) (map f (rest coll)))))

;; Filter function
(defn filter [pred coll]
  (if (empty? coll)
      ()
      (if (pred (first coll))
          (cons (first coll) (filter pred (rest coll)))
          (filter pred (rest coll)))))

;; Range function (reverse order for simplicity)
(defn range [n]
  (if (= n 0)
      ()
      (cons (- n 1) (range (- n 1)))))

;; Reduce function (simplified)
(defn reduce [f init coll]
  (if (empty? coll)
      init
      (reduce f (f init (first coll)) (rest coll))))

;; Group-by function - group collection by key function (simplified)
(defn group-by [f coll]
  ;; Simplified implementation - returns list of (key value) pairs
  (reduce (fn [acc item]
            (let [key (f item)]
              (cons (list key item) acc)))
          nil
          coll))

;; Improved map function with two collections support
(defn map2 [f coll1 coll2]
  (if (empty? coll1)
      nil
      (if (empty? coll2)
          nil
          (cons (f (first coll1) (first coll2))
                (map2 f (rest coll1) (rest coll2))))))