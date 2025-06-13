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
  (defn str-concat [strings]
    "Concatenate a list of strings"
    (string-join strings ""))

  (defn str-empty [str]
    "Check if string is empty"
    (string-empty? str))

  (defn str-blank? [str]
    "Check if string is empty or contains only whitespace"
    (string-empty? (string-trim str)))

  (defn str-non-empty? [str]
    "Check if string is not empty"
    (not (string-empty? str)))

  ; String analysis functions
  (defn str-words [str]
    "Split string into words (by whitespace)"
    (filter (fn [word] (not (string-empty? word)))
            (string-split (string-trim str) " ")))

  (defn str-lines [str]
    "Split string into lines"
    (string-split str "\n"))

  (defn str-char-count [str char]
    "Count occurrences of a character in string"
    (length (string-regex-find-all str char)))

  ; String transformation functions
  (defn str-reverse [str]
    "Reverse a string"
    (string-join (reverse (map (fn [i] (string-char-at str i))
                              (range 0 (string-length str)))) ""))

  ; Helper function to create a range of numbers
  (defn range [start end]
    "Create a list of numbers from start to end-1"
    (if (>= start end)
        (list)
        (cons start (range (+ start 1) end))))

  (defn str-capitalize [str]
    "Capitalize first character of string"
    (if (string-empty? str)
        str
        (string-concat (string-upper (string-substring str 0 1))
                      (string-lower (string-substring str 1 (string-length str))))))

  (defn str-title-case [str]
    "Convert string to title case (capitalize each word)"
    (string-join (map str-capitalize (str-words str)) " "))

  ; String validation functions
  (defn str-numeric? [str]
    "Check if string contains only numeric characters"
    (string-regex-match? str "^[0-9]+$"))

  (defn str-alpha? [str]
    "Check if string contains only alphabetic characters"
    (string-regex-match? str "^[a-zA-Z]+$"))

  (defn str-alnum? [str]
    "Check if string contains only alphanumeric characters"
    (string-regex-match? str "^[a-zA-Z0-9]+$"))

  ; String formatting functions
  (defn str-pad-left [str width pad-char]
    "Pad string on the left to specified width"
    (let ((current-len (string-length str))
          (needed (- width current-len)))
      (if (<= needed 0)
          str
          (string-concat (string-repeat pad-char needed) str))))

  (defn str-pad-right [str width pad-char]
    "Pad string on the right to specified width"
    (let ((current-len (string-length str))
          (needed (- width current-len)))
      (if (<= needed 0)
          str
          (string-concat str (string-repeat pad-char needed)))))

  (defn str-center [str width pad-char]
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
  (defn let [bindings body]
    "Simple let implementation for local bindings"
    ; This is a simplified let - in a full implementation this would be a special form
    ; For now, we'll implement it as a fn application
    ((fn [names vals]
       (apply (fn names body) vals))
     (map first bindings)
     (map (fn [binding] (first (rest binding))) bindings)))
)
