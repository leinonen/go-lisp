; Practical Built-ins Usage Examples
; This file shows practical scenarios for using the builtins function

; Scenario 1: Feature Detection
; Check if advanced features are available before using them

(define advanced-functions (list "map" "filter" "reduce"))

(defun check-advanced-features ()
  (define available (builtins))
  (map (lambda (func) 
         (list func (member? func available)))
       advanced-functions))

; Show which advanced features are available
(check-advanced-features)

; Scenario 2: Dynamic Function Discovery
; Find all functions that match a pattern (e.g., comparison operators)

(defun find-comparison-ops ()
  (filter (lambda (func)
            (or (= func "=")
                (= func "<")
                (= func ">")
                (= func "<=")
                (= func ">=")))
          (builtins)))

(define comparison-ops (find-comparison-ops))
comparison-ops

; Scenario 3: Help System
; Create a simple help system that lists available functions by category

(defun categorize-functions ()
  (define all-functions (builtins))
  (define arithmetic (filter (lambda (f) 
                              (member? f (list "+" "-" "*" "/")))
                            all-functions))
  (define list-ops (filter (lambda (f)
                            (member? f (list "list" "first" "rest" "cons" 
                                           "length" "empty?" "append" "reverse" "nth")))
                          all-functions))
  (define higher-order (filter (lambda (f)
                                (member? f (list "map" "filter" "reduce")))
                              all-functions))
  (define control-flow (filter (lambda (f)
                                (member? f (list "if" "cond")))
                              all-functions))
  
  (list (list "Arithmetic" arithmetic)
        (list "List Operations" list-ops)
        (list "Higher-Order Functions" higher-order)
        (list "Control Flow" control-flow)))

; Show categorized functions
(categorize-functions)

; Scenario 4: Compatibility Checking
; Check if all required functions for a library are available

(defun check-requirements (required-functions)
  (define available (builtins))
  (define missing (filter (lambda (func)
                           (not (member? func available)))
                         required-functions))
  (if (empty? missing)
      "All required functions are available"
      (list "Missing functions:" missing)))

; Test with some required functions
(define my-library-requirements (list "map" "filter" "reduce" "lambda" "+"))
(check-requirements my-library-requirements)

; Test with some non-existent functions
(define invalid-requirements (list "map" "unknown-function" "another-missing"))
(check-requirements invalid-requirements)

; Scenario 5: Interactive Exploration
; Create a function that suggests related functions

(defun suggest-functions (keyword)
  (define all-funcs (builtins))
  (filter (lambda (func)
            ; Simple substring matching (contains keyword)
            (> (length func) 0)) ; Simplified for demo
          all-funcs))

; Show some function suggestions (simplified version)
(define math-functions (filter (lambda (f) 
                                (member? f (list "+" "-" "*" "/" "=" "<" ">")))
                              (builtins)))
(list "Math functions:" math-functions)

; Scenario 6: Runtime Capability Detection
; Determine interpreter capabilities at runtime

(defun get-interpreter-info ()
  (define total-builtins (length (builtins)))
  (define has-modules (member? "modules" (builtins)))
  (define has-env-inspection (member? "env" (builtins)))
  (define has-higher-order (and (member? "map" (builtins))
                               (member? "filter" (builtins))
                               (member? "reduce" (builtins))))
  
  (list (list "Total built-in functions" total-builtins)
        (list "Module system available" has-modules)
        (list "Environment inspection available" has-env-inspection)
        (list "Higher-order functions available" has-higher-order)))

; Show interpreter capabilities
(get-interpreter-info)
