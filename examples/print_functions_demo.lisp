; Print Functions Demonstration
; This example showcases the print and println functions with various data types

(println "=== Print Functions Demo ===")

; Basic print vs println difference
(println "\n--- Basic Usage ---")
(print "This is 'print' - no automatic newline")
(print " - continues on same line")
(println " - 'println' adds newline")
(println "Next line after println")

; Multiple arguments
(println "\n--- Multiple Arguments ---")
(println "Multiple args:" "arg1" "arg2" "arg3" 42 #t)
(print "Print multiple:" 1 2 3 "hello" #f)
(println)  ; Empty println for newline

; Different data types
(println "\n--- Data Types ---")
(println "String:" "Hello World")
(println "Number:" 42)
(println "Big number:" 123456789012345678901234567890)
(println "Boolean true:" #t)
(println "Boolean false:" #f)
(println "List:" (list 1 2 3 4 5))
(println "Hash map:" (hash-map "name" "Alice" "age" 30))
(println "Function:" (lambda [x] (* x x)))
(println "Nil value:" (list))  ; Empty list evaluates to nil in display

; Mathematical expressions
(println "\n--- Mathematical Results ---")
(println "Addition: 2 + 3 =" (+ 2 3))
(println "Multiplication: 6 * 7 =" (* 6 7))
(println "Subtraction: 10 - 3 =" (- 10 3))
(println "Division: 15 / 3 =" (/ 15 3))
(println "Modulo 17 % 5 =" (% 17 5))

; Define our own simple factorial for demo
(define factorial 
  (lambda [n]
    (if (= n 0)
        1
        (* n (factorial (- n 1))))))

(println "Factorial of 5:" (factorial 5))

; String operations
(println "\n--- String Operations ---")
(define text "Hello World")
(println "Original text:" text)
(println "Length:" (string-length text))
(println "Uppercase:" (string-upper text))
(println "Substring (0,5):" (string-substring text 0 5))
(println "Contains 'World':" (string-contains? text "World"))

; List operations with output
(println "\n--- List Operations ---")
(define numbers (list 1 2 3 4 5))
(println "Numbers:" numbers)
(println "Squared:" (map (lambda [x] (* x x)) numbers))
(println "Even only:" (filter (lambda [x] (= (% x 2) 0)) numbers))
(println "Manual sum: 1+2+3+4+5 =" (+ 1 2 3 4 5))

; Formatting examples
(println "\n--- Formatting Examples ---")
(define format-currency
  (lambda [amount]
    (string-concat "$" (number->string amount))))

(println "Price:" (format-currency 19.99))

(define create-greeting
  (lambda [name time]
    (string-concat "Good " time ", " name "!")))

(println (create-greeting "Alice" "morning"))
(println (create-greeting "Bob" "evening"))

; Table-like output
(println "\n--- Tabular Output ---")
(println "Name        Age    City")
(println "------------------------")
(println "Alice       25     Boston")
(println "Bob         30     Seattle")
(println "Charlie     35     Austin")

; Progress indication simulation
(println "\n--- Progress Simulation ---")
(define show-progress
  (lambda [current total]
    (let ((percent (* (/ current total) 100)))
      (println "Progress:" current "/" total "(" percent "%)"))))

(show-progress 1 5)
(show-progress 3 5)
(show-progress 5 5)

; Error demonstration (commented out to not stop execution)
; (println "\n--- Error Example (commented) ---")
; (println "This would show an error:" (error "Demo error message"))

; Complex nested data
(println "\n--- Complex Data Structures ---")
(define person (hash-map 
  "name" "John Doe"
  "age" 30))

(println "Person data:" person)
(println "Name:" (hash-map-get person "name"))
(println "Age:" (hash-map-get person "age"))

; Function composition with output
(println "\n--- Function Composition ---")
(define compose
  (lambda [f g]
    (lambda [x] (f (g x)))))

(define add-one (lambda [x] (+ x 1)))
(define double (lambda [x] (* x 2)))
(define add-one-then-double (compose double add-one))

(println "Compose example: (double (add-one 5)) =" (add-one-then-double 5))

; Recursive function output
(println "\n--- Simple Recursion ---")
(define simple-countdown
  (lambda [n]
    (if (= n 0)
        (println "Done!")
        (begin
          (println "Count:" n)
          (simple-countdown (- n 1))))))

(simple-countdown 3)

(println "\n=== Print Demo Complete ===")
