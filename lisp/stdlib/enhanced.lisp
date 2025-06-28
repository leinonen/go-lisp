;; Enhanced Standard Library - Phase 1.1 Implementation
;; Advanced functions implemented in Lisp using core primitives

;; String operations (using core string primitives)
(def join
  (fn [sep coll]
    (if (empty? coll)
        ""
        (reduce (fn [acc x] (str acc sep (str x)))
                (str (first coll))
                (rest coll)))))

(def split string-split)
(def trim string-trim)
(def replace string-replace)
;; Note: contains? is already defined in core for hash-maps/sets
;; string-contains? remains available for string operations

;; Enhanced collection operations
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

(def reverse
  (fn [coll]
    (reduce (fn [acc x] (cons x acc)) nil coll)))

(def take
  (fn [n coll]
    (if (= n 0)
        nil
        (if (empty? coll)
            nil
            (cons (first coll) (take (- n 1) (rest coll)))))))

(def drop
  (fn [n coll]
    (if (= n 0)
        coll
        (if (empty? coll)
            nil
            (drop (- n 1) (rest coll))))))

(def concat
  (fn [coll1 coll2]
    (if (empty? coll1)
        coll2
        (cons (first coll1) (concat (rest coll1) coll2)))))

;; Utility functions
(def inc (fn [x] (+ x 1)))
(def dec (fn [x] (- x 1)))
(def zero? (fn [x] (= x 0)))
(def pos? (fn [x] (> x 0)))
(def neg? (fn [x] (< x 0)))
(def even? (fn [x] (= (% x 2) 0)))
(def odd? (fn [x] (= (% x 2) 1)))

;; Boolean operations
(def and2
  (fn [a b]
    (if a b nil)))

(def or2
  (fn [a b]
    (if a a b)))

;; Enhanced predicates
(def nil? (fn [x] (= x nil)))
(def some? (fn [x] (not (nil? x))))
(def true? (fn [x] (= x true)))
(def false? (fn [x] (nil? x)))

;; Collection predicates
(def seq? (fn [x] (or2 (list? x) (vector? x))))
(def coll? (fn [x] (or2 (list? x) (vector? x))))

;; Functional utilities
(def comp
  (fn [f g]
    (fn [x] (f (g x)))))

(def constantly
  (fn [x]
    (fn [y] x)))

(def identity (fn [x] x))

;; Math utilities
(def min
  (fn [a b]
    (if (< a b) a b)))

(def max
  (fn [a b]
    (if (> a b) a b)))

(def abs
  (fn [x]
    (if (< x 0) (- x) x)))

;; Sequence utilities
(def last
  (fn [coll]
    (if (empty? (rest coll))
        (first coll)
        (last (rest coll)))))

(def butlast
  (fn [coll]
    (if (empty? (rest coll))
        nil
        (cons (first coll) (butlast (rest coll))))))

(def distinct
  (fn [coll]
    (if (empty? coll)
        nil
        (let [head (first coll)
              rest-distinct (distinct (rest coll))]
          (if (contains-item? head rest-distinct)
              rest-distinct
              (cons head rest-distinct))))))

(def contains-item?
  (fn [item coll]
    (if (empty? coll)
        nil
        (if (= item (first coll))
            true
            (contains-item? item (rest coll))))))

;; Repeat function
(def repeat
  (fn [n x]
    (if (= n 0)
        nil
        (cons x (repeat (- n 1) x)))))

;; Some useful conditional helpers
(def when
  (fn [condition body]
    (if condition body nil)))

(def unless
  (fn [condition body]
    (if condition nil body)))

;; Enhanced collection functions
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

;; All function
(def all?
  (fn [pred coll]
    (if (empty? coll)
        true
        (if (pred (first coll))
            (all? pred (rest coll))
            nil))))

;; Any function  
(def any?
  (fn [pred coll]
    (if (empty? coll)
        nil
        (if (pred (first coll))
            true
            (any? pred (rest coll))))))

;; Partition function
(def partition
  (fn [n coll]
    (if (< (count coll) n)
        nil
        (cons (take n coll) (partition n (drop n coll))))))

;; Interpose
(def interpose
  (fn [sep coll]
    (if (empty? coll)
        nil
        (if (empty? (rest coll))
            coll
            (cons (first coll) (cons sep (interpose sep (rest coll))))))))

;; Remove function
(def remove
  (fn [pred coll]
    (filter (fn [x] (not (pred x))) coll)))

;; Keep function
(def keep
  (fn [f coll]
    (if (empty? coll)
        nil
        (let [result (f (first coll))]
          (if (nil? result)
              (keep f (rest coll))
              (cons result (keep f (rest coll))))))))

;; Flatten (simple version)
(def flatten
  (fn [coll]
    (if (empty? coll)
        nil
        (if (list? (first coll))
            (concat (flatten (first coll)) (flatten (rest coll)))
            (cons (first coll) (flatten (rest coll)))))))