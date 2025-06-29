;; Self-hosting GoLisp Compiler
;; This file demonstrates how GoLisp can compile itself

;; Core compiler data structures
(def *current-env* nil)
(def *compile-target* 'eval) ; 'eval or 'file

;; Symbol table for tracking definitions
(def *symbol-table* (hash-map))

;; Compilation context
(defn make-context []
  {:symbols (hash-map)
   :locals '()
   :target *compile-target*})

;; Core compilation functions
(defn compile-expr [expr ctx]
  (cond
    (symbol? expr) (compile-symbol expr ctx)
    (list? expr) (compile-list expr ctx)
    (vector? expr) (compile-vector expr ctx)
    :else expr)) ; literals

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
      (if (not (symbol? name))
        (throw (str "def name must be a symbol"))
        (do
          ;; Register in symbol table
          (hash-map-put *symbol-table* name :global)
          ;; Compile the value
          (list 'def name (compile-expr value ctx)))))))

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
          fn-ctx (assoc ctx :locals new-locals)]
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
      (let [extract-symbols (fn [bindings-list acc]
                              (if (empty? bindings-list)
                                acc
                                (if (empty? (rest bindings-list))
                                  (throw (str "odd number of binding forms"))
                                  (extract-symbols (rest (rest bindings-list))
                                                 (conj acc (first bindings-list))))))
            binding-symbols (extract-symbols bindings '())
            ;; Create new context with locals  
            new-locals (reduce (fn [acc sym] (conj acc sym)) (:locals ctx) binding-symbols)
            let-ctx (assoc ctx :locals new-locals)
            ;; Compile bindings values with current context
            compile-bindings (fn [bindings-list acc]
                              (if (empty? bindings-list)
                                acc
                                (if (empty? (rest bindings-list))
                                  (throw (str "odd number of binding forms"))
                                  (compile-bindings (rest (rest bindings-list))
                                                  (conj (conj acc (first bindings-list))
                                                        (compile-expr (second bindings-list) ctx))))))
            compiled-bindings (compile-bindings bindings [])]
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
