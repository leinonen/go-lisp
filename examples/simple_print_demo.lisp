; Simple Print Functions Demo
; Demonstrates basic print and println functionality

(println "=== Simple Print Demo ===")

; Basic usage
(println "Hello, World!")
(print "This is print: ")
(println "this continues the line")

; Multiple arguments
(println "Multiple args:" 1 2 3 "hello" #t #f)

; Different data types
(println "String:" "Hello")
(println "Number:" 42)
(println "Boolean:" #t)
(println "List:" (list 1 2 3))
(println "Hash map:" (hash-map "key" "value"))

; Mathematical operations
(println "Math: 2 + 3 =" (+ 2 3))
(println "Modulo: 17 % 5 =" (% 17 5))

; String operations
(define text "Hello World")
(println "Text:" text)
(println "Length:" (string-length text))
(println "Uppercase:" (string-upper text))

; List operations
(define nums (list 1 2 3 4 5))
(println "Numbers:" nums)
(println "First:" (first nums))
(println "Rest:" (rest nums))

; Function definition and use
(define square (lambda (x) (* x x)))
(println "Square of 5:" (square 5))

; Formatting example
(define greeting 
  (lambda (name)
    (string-concat "Hello, " name "!")))

(println (greeting "Alice"))
(println (greeting "Bob"))

(println "=== Demo Complete ===")
