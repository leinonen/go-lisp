;; Core standard library functions implemented in Lisp
;; These complement the Go core primitives

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