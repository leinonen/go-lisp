; String Library Module
; Provides additional string manipulation utilities built on top of the core string functions

(module strings
  (export 
    ; Basic string utilities
    str-concat str-empty str-blank? str-non-empty?
    ; String analysis
    str-words str-lines str-char-count
    ; String transformation
    str-reverse str-capitalize str-title-case
    ; String validation
    str-numeric? str-alpha? str-alnum?
    ; String formatting
    str-pad-left str-pad-right str-center)

  ; Basic string utilities
  (defun str-concat (strings)
    "Concatenate a list of strings"
    (string-join strings ""))

  (defun str-empty (str)
    "Check if string is empty"
    (string-empty? str))

  (defun str-blank? (str)
    "Check if string is empty or contains only whitespace"
    (string-empty? (string-trim str)))

  (defun str-non-empty? (str)
    "Check if string is not empty"
    (not (string-empty? str)))

  ; String analysis functions
  (defun str-words (str)
    "Split string into words (by whitespace)"
    (filter (lambda (word) (not (string-empty? word)))
            (string-split (string-trim str) " ")))

  (defun str-lines (str)
    "Split string into lines"
    (string-split str "\n"))

  (defun str-char-count (str char)
    "Count occurrences of a character in string"
    (length (string-regex-find-all str char)))

  ; String transformation functions
  (defun str-reverse (str)
    "Reverse a string"
    (string-join (reverse (map (lambda (i) (string-char-at str i))
                              (range 0 (string-length str)))) ""))

  ; Helper function to create a range of numbers
  (defun range (start end)
    "Create a list of numbers from start to end-1"
    (if (>= start end)
        (list)
        (cons start (range (+ start 1) end))))

  (defun str-capitalize (str)
    "Capitalize first character of string"
    (if (string-empty? str)
        str
        (string-concat (string-upper (string-substring str 0 1))
                      (string-lower (string-substring str 1 (string-length str))))))

  (defun str-title-case (str)
    "Convert string to title case (capitalize each word)"
    (string-join (map str-capitalize (str-words str)) " "))

  ; String validation functions
  (defun str-numeric? (str)
    "Check if string contains only numeric characters"
    (string-regex-match? str "^[0-9]+$"))

  (defun str-alpha? (str)
    "Check if string contains only alphabetic characters"
    (string-regex-match? str "^[a-zA-Z]+$"))

  (defun str-alnum? (str)
    "Check if string contains only alphanumeric characters"
    (string-regex-match? str "^[a-zA-Z0-9]+$"))

  ; String formatting functions
  (defun str-pad-left (str width pad-char)
    "Pad string on the left to specified width"
    (let ((current-len (string-length str))
          (needed (- width current-len)))
      (if (<= needed 0)
          str
          (string-concat (string-repeat pad-char needed) str))))

  (defun str-pad-right (str width pad-char)
    "Pad string on the right to specified width"
    (let ((current-len (string-length str))
          (needed (- width current-len)))
      (if (<= needed 0)
          str
          (string-concat str (string-repeat pad-char needed)))))

  (defun str-center (str width pad-char)
    "Center string within specified width"
    (let ((current-len (string-length str))
          (total-padding (- width current-len)))
      (if (<= total-padding 0)
          str
          (let ((left-padding (/ total-padding 2))
                (right-padding (- total-padding left-padding)))
            (string-concat (string-repeat pad-char left-padding)
                          str
                          (string-repeat pad-char right-padding))))))

  ; Private helper function (not exported)
  (defun let (bindings body)
    "Simple let implementation for local bindings"
    ; This is a simplified let - in a full implementation this would be a special form
    ; For now, we'll implement it as a lambda application
    ((lambda (names vals)
       (apply (lambda names body) vals))
     (map first bindings)
     (map (lambda (binding) (first (rest binding))) bindings)))
)
