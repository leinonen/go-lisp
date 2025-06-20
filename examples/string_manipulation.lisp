;; String manipulation examples in GoLisp

;; Basic string operations
(string-concat "Hello" " " "World" "!")           ; => "Hello World!"
(string-length "Hello")                           ; => 5
(string-substring "Hello World" 0 5)              ; => "Hello"
(string-char-at "Hello" 1)                        ; => "e"

;; Case conversion
(string-upper "hello world")                      ; => "HELLO WORLD"
(string-lower "HELLO WORLD")                      ; => "hello world"
(string-trim "  hello world  ")                   ; => "hello world"

;; String searching
(string-contains? "Hello World" "World")          ; => true
(string-starts-with? "Hello World" "Hello")       ; => true
(string-ends-with? "Hello World" "World")         ; => true
(string-index-of "Hello World" "o")               ; => 4

;; String splitting and joining
(string-split "apple,banana,cherry" ",")          ; => ("apple" "banana" "cherry")
(string-join '("apple" "banana" "cherry") ", ")   ; => "apple, banana, cherry"

;; String replacement
(string-replace "Hello World" "World" "Universe") ; => "Hello Universe"

;; String predicates
(string? "hello")                                 ; => true
(string? 123)                                     ; => false
(string-empty? "")                                ; => true
(string-empty? "hello")                           ; => false

;; String repetition
(string-repeat "Hi! " 3)                          ; => "Hi! Hi! Hi! "

;; Number/string conversion
(number->string 42)                               ; => "42"
(string->number "42")                             ; => 42
(string->number "3.14")                           ; => 3.14

;; Regular expressions
(string-regex-match? "hello123" "^[a-z]+[0-9]+$") ; => true
(string-regex-find-all "The numbers are 123 and 456" "[0-9]+")  ; => ("123" "456")

;; Practical examples
;; Email validation
(defn valid-email? [email]
  (string-regex-match? email "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"))

;; Word count
(defn word-count [text]
  (length (string-split (string-trim text) " ")))

(word-count "Hello world from GoLisp")           ; => 5
