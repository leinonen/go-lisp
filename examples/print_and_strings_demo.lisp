; Complete Print Functions and String Operations Demo
; Demonstrates the integration of print/println with string manipulation

(println "=== Print Functions & String Operations Demo ===")

; Basic print functionality
(println "\n--- Basic Print Functions ---")
(println "This is println - adds newline automatically")
(print "This is print - no newline")
(print " - continues on same line")
(println " - println ends here")

; Multiple arguments
(println "\n--- Multiple Arguments ---")
(println "Multiple values:" 42 "hello" #t (list 1 2 3))
(print "Print also supports multiple args:" 1 2 3)
(println) ; Empty println for newline

; String operations with output
(println "\n--- String Operations ---")
(def sample-text "Hello, Lisp World!")
(println "Original:" sample-text)
(println "Length:" (string-length sample-text))
(println "Uppercase:" (string-upper sample-text))
(println "Lowercase:" (string-lower sample-text))
(println "First 5 chars:" (string-substring sample-text 0 5))
(println "Contains 'Lisp':" (string-contains? sample-text "Lisp"))

; String building and formatting
(println "\n--- String Building ---")
(def build-greeting
  (lambda [name time]
    (string-concat "Good " time ", " name "!")))

(println (build-greeting "Alice" "morning"))
(println (build-greeting "Bob" "evening"))

; Data formatting with print
(println "\n--- Data Formatting ---")
(def format-person
  (lambda [name age]
    (string-concat name " is " (number->string age) " years old")))

(println (format-person "Charlie" 25))
(println (format-person "Diana" 30))

; List processing with output
(println "\n--- List Processing ---")
(def numbers (list 1 2 3 4 5))
(println "Original numbers:" numbers)

; Process and display each number
(def print-squares
  (lambda [nums]
    (if (= (length nums) 0)
        (println "Done processing squares")
        (begin
          (def n (first nums))
          (println "Square of" n "is" (* n n))
          (print-squares (rest nums))))))

(print-squares numbers)

; Mathematical expressions with output
(println "\n--- Mathematical Results ---")
(def show-calculation
  (lambda [a b operation]
    (cond
      ((= operation 1) (println a "+" b "=" (+ a b)))
      ((= operation 2) (println a "*" b "=" (* a b)))
      ((= operation 3) (println a "%" b "=" (% a b))))))

(show-calculation 15 4 1)  ; Addition
(show-calculation 15 4 2)  ; Multiplication  
(show-calculation 15 4 3)  ; Modulo

; Hash map operations with output
(println "\n--- Hash Map Display ---")
(def person (hash-map "name" "John" "age" 28 "city" "Boston"))
(println "Person record:" person)
(println "Name field:" (hash-map-get person "name"))
(println "Age field:" (hash-map-get person "age"))
(println "City field:" (hash-map-get person "city"))

; Text processing pipeline
(println "\n--- Text Processing Pipeline ---")
(def process-text
  (lambda [text]
    (begin
      (println "Input text:" text)
      (def trimmed (string-trim text))
      (println "After trim:" trimmed)
      (def upper (string-upper trimmed))
      (println "Uppercase:" upper)
      (def length (string-length upper))
      (println "Final length:" length)
      (println "Processing complete!")
      upper)))

(def result (process-text "  welcome to lisp programming  "))
(println "Final result:" result)

; Interactive menu simulation
(println "\n--- Menu Simulation ---")
(def show-menu
  (lambda []
    (begin
      (println "===== LISP CALCULATOR =====")
      (println "1. Addition")
      (println "2. Multiplication") 
      (println "3. Modulo")
      (println "========================="))))

(show-menu)
(println "Selected option 1: Addition")
(println "Result: 10 + 5 =" (+ 10 5))

; Progress display
(println "\n--- Progress Display ---")
(def show-progress
  (lambda [step total task]
    (println "Step" step "of" total ":" task)))

(show-progress 1 4 "Initializing")
(show-progress 2 4 "Processing data") 
(show-progress 3 4 "Computing results")
(show-progress 4 4 "Complete!")

(println "\n=== Demo Complete ===")
(println "Print and println functions are working perfectly!")
(print "Thanks for trying the demo")
