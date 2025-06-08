; String Library Demonstration
; This example showcases both built-in string functions and higher-level library functions
; Built-in functions are primitives implemented in Go for performance
; Library functions are Lisp compositions built on top of the primitives

(println "=== String Library Demo ===")

; Test built-in string functions (implemented as Go primitives)
(println "\n--- Built-in String Functions (Go Primitives) ---")
(println "String concatenation:" (string-concat "Hello" " " "World"))
(println "String length:" (string-length "Hello World"))
(println "Substring:" (string-substring "Hello World" 0 5))
(println "Character at index 6:" (string-char-at "Hello World" 6))
(println "Uppercase:" (string-upper "hello world"))
(println "Lowercase:" (string-lower "HELLO WORLD"))
(println "Trimmed:" (string-trim "  Hello World  "))

; Test string search and manipulation
(println "\n--- String Search & Manipulation ---")
(println "Contains 'World':" (string-contains? "Hello World" "World"))
(println "Starts with 'Hello':" (string-starts-with? "Hello World" "Hello"))
(println "Ends with 'World':" (string-ends-with? "Hello World" "World"))
(println "Replace 'World' with 'Universe':" (string-replace "Hello World" "World" "Universe"))
(println "Index of 'World':" (string-index-of "Hello World" "World"))
(println "Repeat 'Hi' 3 times:" (string-repeat "Hi" 3))

; Test string splitting and joining
(println "\n--- String Splitting & Joining ---")
(define words (string-split "apple,banana,cherry" ","))
(println "Split by comma:" words)
(println "Join with ' | ':" (string-join words " | "))

; Test string validation
(println "\n--- String Validation ---")
(println "Is '123' a string?:" (string? "123"))
(println "Is '123' empty?:" (string-empty? "123"))
(println "Is '' empty?:" (string-empty? ""))

; Test type conversion
(println "\n--- Type Conversion ---")
(println "String to number '42.5':" (string->number "42.5"))
(println "Number to string 42.5:" (number->string 42.5))

; Test regex functions
(println "\n--- Regular Expressions ---")
(println "Regex match 'H.*d' in 'Hello World':" (string-regex-match? "Hello World" "H.*d"))
(println "Find all digits in 'abc123def456':" (string-regex-find-all "abc123def456" "[0-9]+"))

; Load and test the high-level string library
(println "\n--- High-Level String Library (Lisp Compositions) ---")
(println "These functions are built on top of the primitives for convenience and readability")

; Note: In a real implementation, we would load the module like this:
; (load "library/strings.lisp")

; For now, let's demonstrate the concept with inline implementations
; These show how to build higher-level functions from the primitive built-ins

(define str-reverse 
  (lambda (str)
    "Reverse a string by converting to list and back"
    (let ((chars (map (lambda (i) (string-char-at str i))
                      (range 0 (string-length str)))))
      (string-join (reverse chars) ""))))

(define range
  (lambda (start end)
    "Create a list of numbers from start to end-1"
    (if (>= start end)
        (list)
        (cons start (range (+ start 1) end)))))

(define str-capitalize
  (lambda (str)
    "Capitalize first character of string - built from primitives"
    (if (string-empty? str)
        str
        (string-concat (string-upper (string-substring str 0 1))
                      (string-lower (string-substring str 1 (string-length str)))))))

(define str-title-case
  (lambda (str)
    "Convert to title case - capitalize each word"
    (let ((words (string-split (string-trim str) " ")))
      (string-join (map str-capitalize words) " "))))

(define str-blank?
  (lambda (str)
    "Check if string is empty or only whitespace"
    (string-empty? (string-trim str))))

(define str-numeric?
  (lambda (str)
    "Check if string contains only numeric characters"
    (string-regex-match? str "^[0-9]+$")))

; Test our higher-level utility functions
(println "String reverse 'Hello':" (str-reverse "Hello"))
(println "Capitalize 'hello world':" (str-capitalize "hello world"))
(println "Title case 'hello world example':" (str-title-case "hello world example"))
(println "Is '  ' blank?:" (str-blank? "  "))
(println "Is '123' numeric?:" (str-numeric? "123"))
(println "Is 'abc' numeric?:" (str-numeric? "abc"))

; Demonstrate interactive string processing
(println "\n--- Interactive String Processing ---")
(define sample-text "  The Quick Brown Fox Jumps Over The Lazy Dog  ")
(println "Original text:" sample-text)
(println "Trimmed:" (string-trim sample-text))
(println "Uppercase:" (string-upper (string-trim sample-text)))
(println "Lowercase:" (string-lower (string-trim sample-text)))
(println "Capitalized:" (str-capitalize (string-lower (string-trim sample-text))))

; Demonstrate string analysis with output
(println "\n--- String Analysis ---")
(define analysis-text "Hello World 123!")
(println "Analyzing text:" analysis-text)
(println "Length:" (string-length analysis-text))
(println "Contains digits:" (string-regex-match? analysis-text "[0-9]"))
(println "Contains letters:" (string-regex-match? analysis-text "[a-zA-Z]"))
(println "First 5 characters:" (string-substring analysis-text 0 5))
(println "Last 4 characters:" (string-substring analysis-text (- (string-length analysis-text) 4) (string-length analysis-text)))

; Demonstrate practical use cases
(println "\n--- Practical Examples ---")

; Example 1: Text formatting
(define format-name 
  (lambda (first last)
    (string-concat (str-capitalize first) " " (str-capitalize last))))

(println "Formatted name:" (format-name "john" "doe"))

; Example 2: Data processing with output
(define csv-line "apple,banana,cherry,date")
(println "CSV data:" csv-line)
(define fruits (string-split csv-line ","))
(println "Parsed fruits:" fruits)
(println "Fruit count:" (length fruits))

; Print each fruit with formatting
(define print-fruit-list
  (lambda (fruits)
    (if (null? fruits)
        "Done!"
        (begin
          (println "- " (str-capitalize (first fruits)))
          (print-fruit-list (rest fruits))))))

(println "Formatted fruit list:")
(print-fruit-list fruits)

(println "\n=== Demo Complete ===")

; Example usage showing print vs println difference
(println "\n--- Print vs Println Demonstration ---")
(print "This is print: no newline")
(print " - continues on same line")
(println " - println ends the line")
(println "New line after println")

; Multiple arguments
(println "Multiple arguments:" "arg1" "arg2" "arg3")
(print "Print multiple:" "arg1" "arg2" "arg3")
(println)  ; Empty println for newline
; (use strings)

; For this demo, we'll define some of the functions inline to test them
(defun str-words (str)
  "Split string into words (by whitespace)"
  (filter (lambda (word) (not (string-empty? word)))
          (string-split (string-trim str) " ")))

(defun str-capitalize (str)
  "Capitalize first character of string"
  (if (string-empty? str)
      str
      (string-concat (string-upper (string-substring str 0 1))
                    (string-lower (string-substring str 1 (string-length str))))))

(defun str-title-case (str)
  "Convert string to title case (capitalize each word)"
  (string-join (map str-capitalize (str-words str)) " "))

(defun str-blank? (str)
  "Check if string is empty or contains only whitespace"
  (string-empty? (string-trim str)))

(defun str-pad-left (str width pad-char)
  "Pad string on the left to specified width"
  (let ((current-len (string-length str))
        (needed (- width current-len)))
    (if (<= needed 0)
        str
        (string-concat (string-repeat pad-char needed) str))))

; Test the high-level functions
(println "\n--- Testing High-Level Functions ---")
(println "Words in 'hello world from lisp':" (str-words "  hello   world   from   lisp  "))
(println "Capitalize 'hello':" (str-capitalize "hello"))
(println "Title case 'hello world':" (str-title-case "hello world"))
(println "Is '   ' blank?:" (str-blank? "   "))
(println "Is 'text' blank?:" (str-blank? "text"))
(println "Pad left 'Hi' to width 10 with '*':" (str-pad-left "Hi" 10 "*"))

; Test string processing pipeline
(println "\n--- String Processing Pipeline ---")
(define text "  HELLO world from LISP interpreter  ")
(println "Original:" text)
(println "Trimmed & title case:" (str-title-case (string-trim text)))

; Test with different data types
(println "\n--- Edge Cases & Error Handling ---")
(println "Empty string operations:")
(println "  Length of '':" (string-length ""))
(println "  Trim '':" (string-trim ""))
(println "  Upper '':" (string-upper ""))

(println "\nString library demonstration complete!")
