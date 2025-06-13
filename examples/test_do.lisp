;; Test do construct
(println! "Testing do construct...")
(def result (do (def x 5) (def y 10) (+ x y)))
(println! "Result:" result)
