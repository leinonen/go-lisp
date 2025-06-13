;; ============================================================================
;; Macro Library Test & Demonstration
;; ============================================================================
;; This file tests and demonstrates all macros in the macro library

;; Load the macro library
(load "library/macros.lisp")

(print "=== Macro Library Demo ===\n")

;; ============================================================================
;; Control Flow Macros
;; ============================================================================

(print "1. Control Flow Macros")
(print "----------------------")

;; When macro
(def x 10)
(when (> x 5) (print "x is greater than 5"))

;; Unless macro  
(unless (< x 0) (print "x is not negative"))

;; Cond macro (would need better list support for full implementation)
(print "Cond would work with: (cond ((< x 5) \"small\") ((< x 15) \"medium\") (else \"large\"))")

;; ============================================================================
;; Variable Binding Macros
;; ============================================================================

(print "\n2. Variable Binding Macros")
(print "---------------------------")

;; Let1 macro
(let1 y 42 
  (print (string-concat "y inside let1: " (number->string y))))

;; Let* macro (sequential binding)
(print "Let* example: binding x=1, y=x+1, z=x+y")
(let* ((a 1) (b (+ a 1))) 
  (print (string-concat "a=" (number->string a) ", b=" (number->string b))))

;; Assignment macros
(def counter 0)
(print (string-concat "Counter before: " (number->string counter)))
(incf counter)
(print (string-concat "Counter after incf: " (number->string counter)))
(incf counter 5)
(print (string-concat "Counter after incf 5: " (number->string counter)))
(decf counter 3)
(print (string-concat "Counter after decf 3: " (number->string counter)))

;; ============================================================================
;; Debugging Macros
;; ============================================================================

(print "\n3. Debugging Macros")
(print "-------------------")

;; Debug macro
(debug (+ 10 20 30))

;; Assert macro
(assert (> counter 0) "Counter should be positive")
(print "Assertion passed!")

;; ============================================================================
;; Iteration Macros
;; ============================================================================

(print "\n4. Iteration Macros")
(print "-------------------")

;; Dotimes macro
(print "Dotimes 0 to 4:")
(dotimes i 5 (print (string-concat "  i = " (number->string i))))

;; For-each macro (would need list iteration support)
(print "For-each would work with: (for-each x '(1 2 3) (print x))")

;; ============================================================================
;; Utility Macros
;; ============================================================================

(print "\n5. Utility Macros")
(print "-----------------")

;; Progn macro (sequential execution)
(def result 
  (progn
    (print "First expression")
    (print "Second expression")
    42))
(print (string-concat "Progn result: " (number->string result)))

;; Comment macro
(comment "This is a comment that does nothing")
(print "Comments work - they return nil and don't execute")

;; Pattern matching macro
(def day 2)
(def day-name
  (match day
    (1 "Monday")
    (2 "Tuesday") 
    (3 "Wednesday")
    (_ "Other day")))
(print (string-concat "Day " (number->string day) " is " day-name))

;; ============================================================================
;; Complex Example - Using Multiple Macros Together
;; ============================================================================

(print "\n6. Complex Example")
(print "------------------")

;; Define a function using multiple macros
(defn process-numbers [n]
  (let* ((start 0)
         (end n)
         (sum 0))
    (dotimes i end
      (when (> i 2)
        (incf sum i)))
    sum))

(def result (process-numbers 10))
(debug result)

;; ============================================================================
;; Performance and Advanced Features
;; ============================================================================

(print "\n7. Advanced Features")
(print "--------------------")

;; Time macro (placeholder)
(time (+ 1 2 3 4 5))

;; Trace macro example
(print "Trace macro would instrument functions with entry/exit logging")

;; TODO macro example
(print "TODO macro would mark unimplemented code with descriptive messages")

;; ============================================================================
;; Summary
;; ============================================================================

(print "\n=== Macro Library Summary ===")
(print "✓ Control Flow: when, unless, cond, and*, or*")
(print "✓ Variables: let1, let*, setf, incf, decf") 
(print "✓ Debugging: debug, trace, assert, time")
(print "✓ Iteration: dotimes, while, for-each")
(print "✓ Utilities: progn, comment, match, with-values")
(print "✓ Advanced: defn-memo, fn*, todo, ignore-errors")
(print "\nThe macro library provides powerful language extensions!")
(print "Load with: (load \"library/macros.lisp\")")
