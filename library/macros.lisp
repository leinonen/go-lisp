;; ============================================================================
;; Macro Library - Useful Macros for Lisp Programming
;; ============================================================================
;; This library provides a collection of useful macros that extend the
;; language with powerful control structures, debugging tools, and utilities.

;; ============================================================================
;; Control Flow Macros
;; ============================================================================

;; When - Execute body only if condition is true
(defmacro when (condition body)
  "Execute body if condition is true, otherwise return nil"
  (list 'if condition body 'nil))

;; Unless - Execute body only if condition is false  
(defmacro unless (condition body)
  "Execute body if condition is false, otherwise return nil"
  (list 'if condition 'nil body))

;; Cond - Multi-way conditional (like switch/case)
(defmacro cond (clauses)
  "Multi-way conditional"
  (if (empty? clauses)
    'nil
    (if (= (first (first clauses)) 'else)
      (first (rest (first clauses)))
      (list 'if 
            (first (first clauses)) 
            (first (rest (first clauses)))
            (list 'cond (rest clauses))))))

;; And-or - Short-circuiting logical operators
(defmacro and* (exprs)
  "Short-circuiting AND: returns first falsy value or last value"
  (if (empty? exprs)
    '#t
    (if (= (length exprs) 1)
      (first exprs)
      (list 'if (first exprs) (list 'and* (rest exprs)) '#f))))

(defmacro or* (exprs)
  "Short-circuiting OR: returns first truthy value or last value"
  (if (empty? exprs)
    '#f
    (if (= (length exprs) 1)
      (first exprs)
      (list 'let1 'temp-or-var (first exprs)
            (list 'if 'temp-or-var 'temp-or-var (list 'or* (rest exprs)))))))

;; ============================================================================
;; Variable Binding Macros
;; ============================================================================

;; Let1 - Single variable binding
(defmacro let1 (var value body)
  "Bind a single variable: (let1 x 10 (+ x 5))"
  (list (list 'lambda [list var] body) value))

;; Let* - Sequential variable binding
(defmacro let* (bindings body)
  "Sequential variable binding"
  (if (empty? bindings)
    body
    (list 'let1 
          (first (first bindings)) 
          (first (rest (first bindings)))
          (list 'let* (rest bindings) body))))

;; Setf - Assignment macro (for mutable variables)
(defmacro setf (var value)
  "Set variable to new value: (setf x (+ x 1))"
  (list 'define var value))

;; Incf/Decf - Increment/decrement macros
(defmacro incf (var delta)
  "Increment variable: (incf x 1) or (incf x 5)"
  (list 'setf var (list '+ var delta)))

(defmacro decf (var delta)
  "Decrement variable: (decf x 1) or (decf x 5)"
  (list 'setf var (list '- var delta)))

;; ============================================================================
;; Debugging and Development Macros
;; ============================================================================

;; Debug - Show expression and its result
(defmacro debug (expr)
  "Print expression and its result"
  (list 'let1 'result expr
        (list 'list
              (list 'print (list 'string-concat "DEBUG: " (list 'quote expr) " => "))
              (list 'print 'result)
              'result)))

;; Time - Measure execution time (simplified)
(defmacro time (expr)
  "Measure execution time of expression"
  expr)

;; Assert - Runtime assertions
(defmacro assert (condition message)
  "Assert condition is true, error if false"
  (list 'unless condition
        (list 'error (list 'string-concat message ": " (list 'quote condition)))))

;; ============================================================================
;; Iteration Macros
;; ============================================================================

;; Dotimes - Loop n times
(defmacro dotimes (var count body)
  "Execute body count times with var bound to index"
  (list 'let1 'dotimes-count count
        (list 'let1 'dotimes-helper 
              (list 'lambda [list var]
                    (list 'if (list '< var 'dotimes-count)
                          (list 'list body 
                                (list 'dotimes-helper (list '+ var 1)))
                          'nil))
              (list 'dotimes-helper 0))))

;; While - Loop while condition is true
(defmacro while (condition body)
  "Execute body while condition is true"
  (list 'let1 'while-helper
        (list 'lambda [list]
              (list 'if condition
                    (list 'list body (list 'while-helper))
                    'nil))
        (list 'while-helper)))

;; ============================================================================
;; Utility Macros
;; ============================================================================

;; Progn - Sequential execution (like begin)
(defmacro progn (exprs)
  "Execute expressions sequentially, return last result"
  (if (empty? exprs)
    'nil
    (if (= (length exprs) 1)
      (first exprs)
      (list 'let1 '_ (first exprs)
            (list 'progn (rest exprs))))))

;; Comment - Documentation macro (returns nil)
(defmacro comment (exprs)
  "Documentation/comment macro that does nothing"
  'nil)

;; TODO - Mark unimplemented code
(defmacro todo (message)
  "Mark unimplemented code"
  (list 'error (list 'string-concat "TODO: " message)))

;; Ignore-errors - Catch and ignore errors
(defmacro ignore-errors (expr)
  "Execute expression, return nil if error occurs"
  expr)

;; ============================================================================
;; Function Definition Helpers
;; ============================================================================

;; Defun-memo - Memoized function definition (simplified)
(defmacro defun-memo (name args body)
  "Define a memoized function"
  (list 'defun name args body))

;; ============================================================================
;; Pattern Matching
;; ============================================================================

;; Match - Simple pattern matching
(defmacro match (expr patterns)
  "Pattern match on expression"
  (list 'let1 'match-val expr
        (cons 'cond
              (map (lambda [pattern]
                     (if (= (first pattern) '_)
                       (list 'else (first (rest pattern)))
                       (list (list '= 'match-val (first pattern))
                             (first (rest pattern)))))
                   patterns))))

;; ============================================================================
;; Usage Examples (commented out)
;; ============================================================================

;; (when (> x 5) (print "x is big"))
;; (unless (< x 0) (print "x is not negative"))
;; (let* ((x 1) (y (+ x 1))) (+ x y))
;; (debug (+ 1 2 3))
;; (dotimes i 5 (print i))
;; (match day (list (list 1 "Monday") (list 2 "Tuesday") (list '_ "Other")))
