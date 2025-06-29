;; Self-hosting GoLisp Compiler
;; This file demonstrates how GoLisp can compile itself

;; Helper functions
(defn not= [a b] (if (= a b) false true))

;; Simple any? function for optimization code
(defn any? [pred coll]
  (if (empty? coll)
      false
      (if (pred (first coll))
          true
          (any? pred (rest coll)))))

;; Simple map function for optimization code
(defn map [f coll]
  (if (empty? coll)
      '()
      (cons (f (first coll)) (map f (rest coll)))))

;; Simple reduce function for optimization code  
(defn reduce [f init coll]
  (if (empty? coll)
      init
      (reduce f (f init (first coll)) (rest coll))))

;; Simple filter function for optimization code
(defn filter [pred coll]
  (if (empty? coll)
      '()
      (if (pred (first coll))
          (cons (first coll) (filter pred (rest coll)))
          (filter pred (rest coll)))))

;; Helper functions used in compiler
(defn second [coll] (first (rest coll)))
(defn length [coll] (count coll))

;; Simple concat function  
(defn concat [coll1 coll2]
  (if (empty? coll1)
      coll2
      (cons (first coll1) (concat (rest coll1) coll2))))

;; Flatten one level
(defn flatten1 [colls]
  (reduce concat '() colls))

;; Reverse function
(defn reverse [coll]
  (reduce (fn [acc item] (cons item acc)) '() coll))

;; Optimization functions

;; Check if a value is a constant (number, string, boolean, nil, keyword)
(defn constant? [expr]
  (or (number? expr)
      (string? expr)
      (nil? expr)
      (keyword? expr)
      (= expr 'true)
      (= expr 'false)))

;; Evaluate arithmetic operations on constants
(defn eval-constant-arith [op args]
  (cond
    (= op '+) (reduce + 0 args)
    (= op '-) (if (= (count args) 1) 
                (- (first args))
                (reduce - (first args) (rest args)))
    (= op '*) (reduce * 1 args)
    (= op '/) (if (= (count args) 1)
                (/ 1 (first args))
                (reduce / (first args) (rest args)))
    (= op '=) (if (= (count args) 2) 
                  (if (= (first args) (second args)) 'true 'false)
                  nil)
    (= op '<) (if (= (count args) 2) 
                  (if (< (first args) (second args)) 'true 'false)
                  nil)
    (= op '>) (if (= (count args) 2) 
                  (if (> (first args) (second args)) 'true 'false)
                  nil)
    (= op '<=) (if (= (count args) 2) 
                   (if (<= (first args) (second args)) 'true 'false)
                   nil)
    (= op '>=) (if (= (count args) 2) 
                   (if (>= (first args) (second args)) 'true 'false)
                   nil)
    :else nil))

;; Constant folding optimization
(defn constant-fold-expr [expr]
  (cond
    ;; Already a constant - return as-is
    (constant? expr) expr
    
    ;; List - check for arithmetic operations with all constant args
    (and (list? expr) (if (empty? expr) false true))
    (let [head (first expr)
          args (rest expr)]
      (if (and (symbol? head)
               (any? (fn [op] (= head op)) '(+ - * / = < > <= >=))
               (if (empty? args) false true)
               (reduce (fn [acc arg] (and acc (constant? (constant-fold-expr arg)))) true args))
        ;; All arguments are constants - try to evaluate
        (let [folded-args (map constant-fold-expr args)
              result (eval-constant-arith head folded-args)]
          (if (not= result nil)
            result
            ;; Couldn't fold - return with folded args
            (cons head folded-args)))
        ;; Not all constants or not arithmetic - recursively fold args
        (cons head (map constant-fold-expr args))))
    
    ;; Vector - recursively fold elements
    (vector? expr)
    (vector (map constant-fold-expr expr))
    
    ;; Other types - return as-is
    :else expr))

;; Dead code elimination - remove unused let bindings
(defn find-used-symbols [expr]
  (cond
    (symbol? expr) (list expr)
    (list? expr) (flatten1 (map find-used-symbols expr))
    (vector? expr) (flatten1 (map find-used-symbols expr))
    :else '()))

(defn eliminate-dead-let-bindings [bindings body]
  ;; For now, just return the original bindings to avoid the complex vector/list conversion
  ;; TODO: Implement proper dead code elimination for let bindings
  bindings)

;; Dead code elimination for if expressions  
(defn eliminate-dead-if [condition then-expr else-expr]
  (cond
    ;; If condition is constant true, return then branch
    (= condition 'true) then-expr
    (= condition 'false) else-expr
    ;; Otherwise return full if expression
    :else (list 'if condition then-expr else-expr)))

;; Main dead code elimination function  
(defn eliminate-dead-code [expr]
  (if (list? expr)
    ;; It's a list - check what kind
    (if (empty? expr)
      expr  ; Empty list
      (cond
        ;; Let expression - skip dead code elimination for now due to vector handling complexity
        (= (first expr) 'let)
        (map eliminate-dead-code expr)
        
        ;; If expression - eliminate unreachable branches
        (= (first expr) 'if)
        (if (= (count expr) 4)
          (let [condition (eliminate-dead-code (second expr))
                then-branch (eliminate-dead-code (nth expr 2))
                else-branch (eliminate-dead-code (nth expr 3))]
            (eliminate-dead-if condition then-branch else-branch))
          expr)  ; Malformed if
        
        ;; Other expressions - recursively optimize
        :else (map eliminate-dead-code expr)))
    ;; Non-list (constant or symbol) - return as-is
    expr))

;; Core compiler data structures
(def *current-env* nil)
(def *compile-target* 'eval) ; 'eval or 'file

;; Symbol table for tracking definitions
(def *symbol-table* (hash-map))

;; Compilation context
(defn make-context []
  (make-context-with-optimizations {:constant-folding true :dead-code-elimination true}))

(defn make-context-with-optimizations [opt-flags]
  (hash-map :symbols (hash-map)
            :locals '()
            :macros (hash-map)  ; Track macro definitions
            :target *compile-target*
            :optimizations opt-flags))

;; Helper to check if an optimization is enabled
(defn optimization-enabled? [ctx opt-name]
  (let [opts (:optimizations ctx)]
    (and opts (get opts opt-name))))

;; Macro expansion support
(def *max-macro-expansion-depth* 20)  ; Lower limit for safety

;; Check if a symbol refers to a macro
(defn is-macro? [sym ctx]
  ;; Check compilation context first for user-defined macros
  (if (and ctx (:macros ctx) (contains? (:macros ctx) sym))
    true
    ;; Fall back to built-in macros using any? instead of contains? with set
    (and (symbol? sym)
         (any? (fn [macro-sym] (= macro-sym sym)) 
               '(when unless cond)))))

;; Expand macros recursively with depth limiting
(defn expand-macros [expr ctx depth]
  (if (> depth *max-macro-expansion-depth*)
    (throw (str "Maximum macro expansion depth exceeded: " depth))
    (cond
      ;; Lists - check for macro expansion
      (and (list? expr) (if (empty? expr) false true))
      (let [head (first expr)]
        (if (is-macro? head ctx)
          ;; It's a macro - expand it and recurse
          (let [expanded (macroexpand expr)]
            (expand-macros expanded ctx (+ depth 1)))
          ;; Not a macro - recursively expand elements
          (map (fn [elem] (expand-macros elem ctx depth)) expr)))
      
      ;; Vectors - recursively expand elements
      (vector? expr)
      (vector (map (fn [elem] (expand-macros elem ctx depth)) expr))
      
      ;; Other types - return as-is (symbols, numbers, strings, etc.)
      :else expr)))

;; Core compilation functions
(defn compile-expr [expr ctx]
  ;; Multi-pass compilation with optimizations
  (let [;; Pass 1: Macro expansion
        expanded-expr (expand-macros expr ctx 0)
        ;; Pass 2: Constant folding (if enabled)
        folded-expr (if (optimization-enabled? ctx :constant-folding)
                      (constant-fold-expr expanded-expr)
                      expanded-expr)
        ;; Pass 3: Core compilation
        compiled-expr (cond
                        (symbol? folded-expr) (compile-symbol folded-expr ctx)
                        (list? folded-expr) (compile-list folded-expr ctx)
                        (vector? folded-expr) (compile-vector folded-expr ctx)
                        :else folded-expr)  ; literals
        ;; Pass 4: Dead code elimination (if enabled)
        optimized-expr (if (optimization-enabled? ctx :dead-code-elimination)
                         (eliminate-dead-code compiled-expr)
                         compiled-expr)]
    optimized-expr))

;; Non-optimizing version for when optimizations should be disabled
(defn compile-expr-no-opt [expr ctx]
  ;; Only macro expansion, no other optimizations
  (let [expanded-expr (expand-macros expr ctx 0)]
    (cond
      (symbol? expanded-expr) (compile-symbol expanded-expr ctx)
      (list? expanded-expr) (compile-list expanded-expr ctx)
      (vector? expanded-expr) (compile-vector expanded-expr ctx)
      :else expanded-expr)))

(defn compile-symbol [sym ctx]
  ;; Check if it's a local binding or global using any? for list-based locals
  (let [is-local? (any? (fn [local] (= local sym)) (:locals ctx))]
    (if is-local?
      sym  ; Local reference - keep as-is
      (do
        ;; Global reference - could be optimized
        sym))))

(defn compile-list [lst ctx]
  (if (empty? lst)
    lst
    (let [head (first lst)
          args (rest lst)]
      (cond
        ;; Special forms
        (= head 'def) (compile-def args ctx)
        (= head 'defmacro) (compile-defmacro args ctx)
        (= head 'fn) (compile-fn args ctx)
        (= head 'if) (compile-if args ctx)
        (= head 'quote) (compile-quote args ctx)
        (= head 'do) (compile-do args ctx)
        (= head 'let) (compile-let args ctx)
        
        ;; Function application
        :else (compile-application head args ctx)))))

(defn compile-def [args ctx]
  (if (not= (length args) 2)
    (throw (str "def requires exactly 2 arguments"))
    (let [name (first args)
          value (second args)]
      (if (if (symbol? name) false true)
        (throw (str "def name must be a symbol"))
        (do
          ;; Register in symbol table
          (hash-map-put *symbol-table* name :global)
          ;; Compile the value
          (list 'def name (compile-expr value ctx)))))))

(defn compile-defmacro [args ctx]
  (if (not= (length args) 3)
    (throw (str "defmacro requires exactly 3 arguments (name params body)"))
    (let [name (first args)
          params (second args)
          body (nth args 2)]
      (if (if (symbol? name) false true)
        (throw (str "defmacro name must be a symbol"))
        (do
          ;; For now, just register in symbol table (context modification is complex)
          (hash-map-put *symbol-table* name :macro)
          ;; Return the defmacro form (no compilation of macro body)
          (list 'defmacro name params body))))))

(defn compile-fn [args ctx]
  (if (< (length args) 2)
    (throw (str "fn requires at least 2 arguments"))
    (let [params (first args)
          body (rest args)
          ;; Create new context with local bindings
          new-locals (reduce (fn [acc param] 
                              (conj acc param)) 
                            (:locals ctx) 
                            params)
          fn-ctx (hash-map :symbols (:symbols ctx)
                           :locals new-locals
                           :macros (:macros ctx)
                           :target (:target ctx)
                           :optimizations (:optimizations ctx))]
      ;; Compile function body with new context  
      (cons 'fn (cons params (map (fn [expr] (compile-expr expr fn-ctx)) body))))))

(defn compile-if [args ctx]
  (if (not= (length args) 3)
    (throw (str "if requires exactly 3 arguments"))
    (list 'if 
          (compile-expr (first args) ctx)
          (compile-expr (second args) ctx)
          (compile-expr (nth args 2) ctx))))

(defn compile-quote [args ctx]
  (if (not= (length args) 1)
    (throw (str "quote requires exactly 1 argument"))
    (list 'quote (first args))))

(defn compile-let [args ctx]
  (if (< (length args) 2)
    (throw (str "let requires at least 2 arguments"))
    (let [bindings (first args)
          body (rest args)]
      ;; Extract symbols from bindings for local context
      (let [;; For now, just compile the body without complex binding processing
            ;; TODO: Implement proper local variable tracking for vectors
            compiled-bindings bindings
            let-ctx ctx]
        ;; Return compiled let form
        (cons 'let 
              (cons compiled-bindings
                    (map (fn [expr] (compile-expr expr let-ctx)) body)))))))

(defn compile-do [args ctx]
  (cons 'do (map (fn [expr] (compile-expr expr ctx)) args)))

(defn compile-application [fn args ctx]
  ;; Compile function call
  (cons (compile-expr fn ctx)
        (map (fn [arg] (compile-expr arg ctx)) args)))

(defn compile-vector [vec ctx]
  ;; Compile vector elements
  (vector (map (fn [elem] (compile-expr elem ctx)) vec)))

;; Main compilation entry point
(defn compile-file [filename output-filename]
  (let [source (slurp filename)
        exprs (read-all source)]
    (print (str "Compiling " filename " to " output-filename))
    (let [compiled-exprs (map (fn [expr] 
                                (compile-expr expr (make-context))) 
                              exprs)]
      ;; For now, just write the compiled expressions back as Lisp
      (spit output-filename 
            (str ";; Compiled from " filename "\n"
                 (str-join "\n" (map str compiled-exprs)))))))

;; Read multiple expressions from a string
(defn read-all [source]
  ;; Use the core read-all-string function to parse multiple expressions
  (read-all-string source))

;; Utility for joining strings
(defn str-join [separator coll]
  (if (empty? coll)
    ""
    (reduce (fn [acc item]
              (if (= acc "")
                (str item)
                (str acc separator item)))
            ""
            coll)))

;; Self-compilation bootstrap
(defn bootstrap-self-hosting []
  (print "=== GoLisp Self-Hosting Bootstrap ===")
  
  ;; Compile the standard library
  (print "1. Compiling standard library...")
  (compile-file "lisp/stdlib/core.lisp" "stdlib-core-compiled.lisp")
  (compile-file "lisp/stdlib/enhanced.lisp" "stdlib-enhanced-compiled.lisp")
  
  ;; Compile this compiler itself!
  (print "2. Compiling self-hosting compiler...")
  (compile-file "self-hosting.lisp" "self-hosting-compiled.lisp")
  
  (print "3. Self-hosting bootstrap complete!")
  (print "   - stdlib-core-compiled.lisp: Compiled core standard library")
  (print "   - stdlib-enhanced-compiled.lisp: Compiled enhanced standard library")
  (print "   - self-hosting-compiled.lisp: Compiled compiler")
  (print "")
  (print "Next steps:")
  (print "   - Load compiled versions to verify correctness")
  (print "   - Implement optimizations in the compiler")
  (print "   - Add code generation for other targets"))

;; Example usage:
;; (bootstrap-self-hosting)
