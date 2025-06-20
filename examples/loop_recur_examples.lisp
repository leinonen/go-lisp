;; Example usage of loop and recur in GoLisp

;; Simple countdown example
(loop [i 5]
  (if (= i 0)
    "done"
    (do
      (println i)
      (recur (- i 1)))))

;; Factorial using loop/recur
(loop [n 5 acc 1]
  (if (<= n 1)
    acc
    (recur (- n 1) (* acc n))))

;; Sum of numbers from 1 to n
(loop [i 1 n 10 sum 0]
  (if (> i n)
    sum
    (recur (+ i 1) n (+ sum i))))

;; The loop construct establishes a recursion point with bindings
;; The recur construct jumps back to the nearest loop with new values
;; This provides efficient tail recursion without growing the call stack
