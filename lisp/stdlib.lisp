;; Minimal working stdlib

(def length 
  (fn [lst]
    (if (= lst ())
      0
      (+ 1 (length (rest lst))))))

(def map 
  (fn [f lst]
    (if (= lst ())
      ()
      (cons (f (first lst)) (map f (rest lst))))))

(def sum
  (fn [lst]
    (if (= lst ())
      0
      (+ (first lst) (sum (rest lst))))))

(def demo 
  (fn []
    (print "=== Standard Library Demo ===")
    (print "Numbers:" (list 1 2 3 4 5))
    (print "Length:" (length (list 1 2 3 4 5)))
    (print "Doubled:" (map (fn [x] (* x 2)) (list 1 2 3 4 5)))
    (print "=== Demo Complete ===")))

;; (print "Minimal standard library loaded!")