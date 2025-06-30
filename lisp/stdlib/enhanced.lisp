;; Enhanced Standard Library - Phase 1.1 Implementation
;; Advanced functions implemented in Lisp using core primitives

;; String operations (using core string primitives)
(defn join [sep coll]
  (if (empty? coll)
      ""
      (reduce (fn [acc x] (str acc sep (str x)))
              (str (first coll))
              (rest coll))))

(def split string-split)
(def trim string-trim)
(def replace string-replace)
;; Note: contains? is already defined in core for hash-maps/sets
;; string-contains? remains available for string operations

;; Enhanced collection operations
(defn apply [f coll]
  (if (empty? coll)
      (f)
      (if (= (count coll) 1)
          (f (first coll))
          (if (= (count coll) 2)
              (f (first coll) (nth coll 1))
              (if (= (count coll) 3)
                  (f (first coll) (nth coll 1) (nth coll 2))
                  (f (first coll) (nth coll 1) (nth coll 2) (nth coll 3)))))))

(defn reverse [coll]
  (reduce (fn [acc x] (cons x acc)) () coll))

(defn take [n coll]
  (if (= n 0)
      ()
      (if (empty? coll)
          ()
          (cons (first coll) (take (- n 1) (rest coll))))))

(defn drop [n coll]
  (if (= n 0)
      coll
      (if (empty? coll)
          ()
          (drop (- n 1) (rest coll)))))

(defn concat [coll1 coll2]
  (if (empty? coll1)
      coll2
      (cons (first coll1) (concat (rest coll1) coll2))))

;; Utility functions
(defn inc [x] (+ x 1))
(defn dec [x] (- x 1))
(defn zero? [x] (= x 0))
(defn pos? [x] (> x 0))
(defn neg? [x] (< x 0))
(defn even? [x] (= (% x 2) 0))
(defn odd? [x] (= (% x 2) 1))

;; Boolean operations
(defn and2 [a b]
  (if a b nil))

(defn or2 [a b]
  (if a a b))

;; Enhanced predicates
(defn nil? [x] (= x nil))
(defn some? [x] (not (nil? x)))
(defn true? [x] (= x true))
(defn false? [x] (nil? x))

;; Collection predicates
(defn seq? [x] (or2 (list? x) (vector? x)))
(defn coll? [x] (or2 (list? x) (vector? x)))

;; Functional utilities
(defn comp [f g]
  (fn [x] (f (g x))))

(defn constantly [x]
  (fn [y] x))

(defn identity [x] x)

;; Math utilities
(defn min [& args]
  (reduce (fn [a b] (if (< a b) a b)) (first args) (rest args)))

(defn max [& args]
  (reduce (fn [a b] (if (> a b) a b)) (first args) (rest args)))

(defn abs [x]
  (if (< x 0) (- x) x))

;; Sequence utilities
(defn last [coll]
  (if (empty? (rest coll))
      (first coll)
      (last (rest coll))))

(defn butlast [coll]
  (if (empty? (rest coll))
      ()
      (cons (first coll) (butlast (rest coll)))))

(defn distinct [coll]
  (if (empty? coll)
      ()
      (let [head (first coll)
            rest-distinct (distinct (rest coll))]
        (if (contains-item? head rest-distinct)
            rest-distinct
            (cons head rest-distinct)))))

(defn contains-item? [item coll]
  (if (empty? coll)
      false
      (if (= item (first coll))
          true
          (contains-item? item (rest coll)))))

;; Repeat function
(defn repeat [n x]
  (if (= n 0)
      ()
      (cons x (repeat (- n 1) x))))


;; Enhanced collection functions
(defn sort [coll]
  (if (empty? coll)
      ()
      (if (= (count coll) 1)
          coll
          (let [pivot (first coll)
                rest-coll (rest coll)
                smaller (filter (fn [x] (< x pivot)) rest-coll)
                greater (filter (fn [x] (>= x pivot)) rest-coll)]
            (concat (sort smaller) (cons pivot (sort greater)))))))

;; All function
(defn all? [pred coll]
  (if (empty? coll)
      true
      (if (pred (first coll))
          (all? pred (rest coll))
          false)))

;; Any function  
(defn any? [pred coll]
  (if (empty? coll)
      false
      (if (pred (first coll))
          true
          (any? pred (rest coll)))))

;; Partition function
(defn partition [n coll]
  (if (< (count coll) n)
      ()
      (cons (take n coll) (partition n (drop n coll)))))

;; Interpose
(defn interpose [sep coll]
  (if (empty? coll)
      ()
      (if (empty? (rest coll))
          coll
          (cons (first coll) (cons sep (interpose sep (rest coll)))))))

;; Remove function
(defn remove [pred coll]
  (filter (fn [x] (not (pred x))) coll))

;; Keep function
(defn keep [f coll]
  (if (empty? coll)
      ()
      (let [result (f (first coll))]
        (if (nil? result)
            (keep f (rest coll))
            (cons result (keep f (rest coll)))))))

;; Flatten (simple version)
(defn flatten [coll]
  (if (empty? coll)
      ()
      (if (list? (first coll))
          (concat (flatten (first coll)) (flatten (rest coll)))
          (cons (first coll) (flatten (rest coll))))))

;; Partial function application
(defn partial [f & partial-args]
  (fn [& remaining-args]
    (apply f (concat partial-args remaining-args))))