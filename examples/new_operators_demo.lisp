; Examples demonstrating new comparison and logical operators

; Comparison operators
(<= 3 5)     ; => #t (3 is less than or equal to 5)
(<= 5 5)     ; => #t (5 is equal to 5)
(<= 7 5)     ; => #f (7 is not less than or equal to 5)

(>= 7 5)     ; => #t (7 is greater than or equal to 5)
(>= 5 5)     ; => #t (5 is equal to 5)
(>= 3 5)     ; => #f (3 is not greater than or equal to 5)

; Logical operators
(and #t #t)  ; => #t (both arguments are true)
(and #t #f)  ; => #f (one argument is false)
(and #f #f)  ; => #f (both arguments are false)

(or #t #t)   ; => #t (at least one argument is true)
(or #t #f)   ; => #t (at least one argument is true)
(or #f #f)   ; => #f (no arguments are true)

(not #t)     ; => #f (opposite of true)
(not #f)     ; => #t (opposite of false)

; Complex expressions combining operators
(and (>= 10 5) (<= 3 7))              ; => #t (both conditions are true)
(or (< 10 5) (> 15 12))               ; => #t (second condition is true)
(not (and (< 5 3) (> 8 10)))          ; => #t (negation of false and)

; Practical examples
(define age 25)
(and (>= age 18) (<= age 65))         ; => #t (age is between 18 and 65)

(define temperature 22)
(or (<= temperature 0) (>= temperature 35))  ; => #f (not extreme temperature)

(define is_weekend #f)
(define is_holiday #t)
(or is_weekend is_holiday)            ; => #t (either weekend or holiday)

; Using in conditional expressions
(if (and (>= age 21) (not is_weekend))
    "Can work today"
    "Cannot work today")               ; => "Can work today"
