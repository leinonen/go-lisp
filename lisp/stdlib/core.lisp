;; Core standard library functions implemented in Lisp
;; These replace the Go implementations in bootstrap.go

;; Logical operations
(def not (fn [x] (if x nil true)))

;; Conditional helpers  
(def when
  (fn [condition body]
    (if condition body nil)))

(def unless
  (fn [condition body]
    (if condition nil body)))

;; Collection operations
(def nth
  (fn [coll n]
    (if (= n 0)
        (first coll)
        (nth (rest coll) (- n 1)))))

;; Second, third helpers
(def second (fn [coll] (first (rest coll))))
(def third (fn [coll] (first (rest (rest coll)))))

;; Count function
(def count
  (fn [coll]
    (if (= coll nil)
        0
        (+ 1 (count (rest coll))))))

;; Empty? predicate
(def empty?
  (fn [coll]
    (= coll nil)))

;; Map function
(def map
  (fn [f coll]
    (if (= coll nil)
        nil
        (cons (f (first coll)) (map f (rest coll))))))

;; Filter function
(def filter
  (fn [pred coll]
    (if (= coll nil)
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
    (if (= coll nil)
        init
        (reduce f (f init (first coll)) (rest coll)))))