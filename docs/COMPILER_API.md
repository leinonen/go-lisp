# GoLisp Self-Hosting Compiler API Documentation

## Overview

This document provides comprehensive API documentation and usage examples for the GoLisp self-hosting compiler. The compiler is implemented entirely in Lisp and provides a rich set of functions for parsing, compiling, and optimizing Lisp code.

## Quick Start

### Loading the Compiler

```bash
# Start GoLisp REPL
./bin/golisp

# Load the self-hosting compiler
(load-file "lisp/self-hosting.lisp")
```

### Basic Compilation Example

```lisp
;; Create a compilation context
(def ctx (make-context))

;; Compile a simple expression
(compile-expr '(+ 1 2 3) ctx)
;; → 6 (constant folding applied)

;; Compile a conditional
(compile-expr '(if true 42 99) ctx)
;; → 42 (dead code elimination applied)
```

## Core API Functions

### 1. Context Management

#### `make-context`
Creates a default compilation context with standard optimizations enabled.

```lisp
(make-context)
;; → {:symbols {} :locals () :macros {} :target eval :optimizations {...}}
```

**Returns:** A hash-map containing:
- `:symbols` - Global symbol table (hash-map)
- `:locals` - Local variable stack (list)  
- `:macros` - Macro definition table (hash-map)
- `:target` - Compilation target (symbol)
- `:optimizations` - Optimization flags (hash-map)

#### `make-context-with-optimizations`
Creates a compilation context with custom optimization settings.

```lisp
;; Disable all optimizations
(make-context-with-optimizations {:constant-folding false :dead-code-elimination false})

;; Enable only constant folding
(make-context-with-optimizations {:constant-folding true :dead-code-elimination false})
```

**Parameters:**
- `opt-flags` (hash-map) - Optimization configuration

**Available optimization flags:**
- `:constant-folding` (boolean) - Enable compile-time constant evaluation
- `:dead-code-elimination` (boolean) - Enable unreachable code removal

#### `optimization-enabled?`
Checks if a specific optimization is enabled in the context.

```lisp
(def ctx (make-context))
(optimization-enabled? ctx :constant-folding)
;; → true

(def no-opt-ctx (make-context-with-optimizations {}))
(optimization-enabled? no-opt-ctx :constant-folding)
;; → nil (false)
```

**Parameters:**
- `ctx` (hash-map) - Compilation context
- `opt-name` (keyword) - Optimization name to check

**Returns:** `true` if enabled, `nil`/`false` if disabled

### 2. Core Compilation Functions

#### `compile-expr`
Main compilation entry point with full optimization pipeline.

```lisp
;; Simple expressions
(compile-expr 42 (make-context))
;; → 42

(compile-expr 'x (make-context))
;; → x

;; Function calls
(compile-expr '(+ 1 2 3) (make-context))
;; → 6 (optimized)

;; Conditionals
(compile-expr '(if (> 5 3) "yes" "no") (make-context))
;; → "yes" (optimized)

;; Function definitions
(compile-expr '(fn [x] (* x x)) (make-context))
;; → (fn [x] (* x x))

;; Let bindings
(compile-expr '(let [x 5 y 10] (+ x y)) (make-context))
;; → (let [x 5 y 10] (+ x y))
```

**Parameters:**
- `expr` - Lisp expression to compile
- `ctx` (hash-map) - Compilation context

**Returns:** Compiled expression (potentially optimized)

**Compilation Pipeline:**
1. **Macro Expansion** - Expand all macros recursively
2. **Constant Folding** - Evaluate constant expressions (if enabled)
3. **Core Compilation** - Transform special forms and function calls
4. **Dead Code Elimination** - Remove unreachable code (if enabled)

#### `compile-expr-no-opt`
Compilation without optimizations (macro expansion only).

```lisp
(compile-expr-no-opt '(+ 1 2 3) (make-context))
;; → (+ 1 2 3) (no constant folding)

(compile-expr-no-opt '(if true 42 99) (make-context))
;; → (if true 42 99) (no dead code elimination)
```

**Use cases:**
- Debugging compilation issues
- Preserving original code structure
- Performance testing

### 3. File Compilation

#### `compile-file`
Compile an entire Lisp source file to an output file.

```lisp
;; Compile standard library
(compile-file "lisp/stdlib/core.lisp" "stdlib-compiled.lisp")

;; Compile application code
(compile-file "my-app.lisp" "my-app-compiled.lisp")
```

**Parameters:**
- `filename` (string) - Input source file path
- `output-filename` (string) - Output compiled file path

**Process:**
1. Read all expressions from source file using `read-all`
2. Compile each expression with shared context
3. Write compiled results to output file with header comment
4. Print compilation status

**Example output file:**
```lisp
;; Compiled from my-app.lisp
(def PI 3.14159)
(defn area [r] (* PI (* r r)))
42
```

#### `read-all`
Parse multiple expressions from a source string.

```lisp
(read-all "(def x 1) (def y 2) (+ x y)")
;; → ((def x 1) (def y 2) (+ x y))

(read-all "(defn square [x] (* x x))\n(square 5)")
;; → ((defn square [x] (* x x)) (square 5))
```

**Parameters:**
- `source` (string) - Source code containing multiple expressions

**Returns:** List of parsed expressions

**Use cases:**
- File compilation
- Multi-expression evaluation
- REPL input processing

### 4. Bootstrap and Self-Compilation

#### `bootstrap-self-hosting`
Complete self-hosting bootstrap process.

```lisp
(bootstrap-self-hosting)
;; Output:
;; === GoLisp Self-Hosting Bootstrap ===
;; 1. Compiling standard library...
;; Compiling lisp/stdlib/core.lisp to stdlib-core-compiled.lisp
;; Compiling lisp/stdlib/enhanced.lisp to stdlib-enhanced-compiled.lisp
;; 2. Compiling self-hosting compiler...
;; Compiling self-hosting.lisp to self-hosting-compiled.lisp
;; 3. Self-hosting bootstrap complete!
```

**Process:**
1. **Compile standard library** - Both core and enhanced modules
2. **Compile compiler itself** - Self-compilation demonstration
3. **Generate output files** - Compiled versions for verification
4. **Report completion** - Summary and next steps

**Generated files:**
- `stdlib-core-compiled.lisp` - Compiled core standard library
- `stdlib-enhanced-compiled.lisp` - Compiled enhanced utilities
- `self-hosting-compiled.lisp` - Compiled compiler

## Optimization Functions

### 1. Constant Folding

#### `constant-fold-expr`
Evaluate constant expressions at compile time.

```lisp
;; Arithmetic folding
(constant-fold-expr '(+ 1 2 3))
;; → 6

(constant-fold-expr '(* 4 (+ 2 3)))
;; → 20

;; Comparison folding
(constant-fold-expr '(< 3 5))
;; → true

(constant-fold-expr '(= "hello" "world"))
;; → false

;; Nested expressions
(constant-fold-expr '(+ (* 2 3) (- 8 2)))
;; → 12

;; Non-constant expressions remain unchanged
(constant-fold-expr '(+ x 1))
;; → (+ x 1)
```

**Supported operations:**
- **Arithmetic**: `+`, `-`, `*`, `/`
- **Comparison**: `=`, `<`, `>`, `<=`, `>=`
- **Nested expressions**: Recursive folding

#### `constant?`
Check if an expression is a compile-time constant.

```lisp
(constant? 42)         ;; → true
(constant? "hello")    ;; → true  
(constant? :keyword)   ;; → true
(constant? true)       ;; → true
(constant? nil)        ;; → true
(constant? 'x)         ;; → false (symbol)
(constant? '(+ 1 2))   ;; → false (expression)
```

**Constant types:**
- Numbers (`42`, `3.14`)
- Strings (`"hello"`)
- Keywords (`:key`)
- Booleans (`true`, `false`)
- Nil (`nil`)

### 2. Dead Code Elimination

#### `eliminate-dead-code`
Remove unreachable code from expressions.

```lisp
;; If statement simplification
(eliminate-dead-code '(if true 42 99))
;; → 42

(eliminate-dead-code '(if false "never" "always"))
;; → "always"

;; Preserve dynamic conditions
(eliminate-dead-code '(if (> x 0) "positive" "not positive"))
;; → (if (> x 0) "positive" "not positive")

;; Nested elimination
(eliminate-dead-code '(if true (if false 1 2) 3))
;; → 2
```

#### `eliminate-dead-if`
Specialized dead code elimination for if expressions.

```lisp
(eliminate-dead-if 'true 42 99)
;; → 42

(eliminate-dead-if 'false 42 99)  
;; → 99

(eliminate-dead-if '(> x 0) 42 99)
;; → (if (> x 0) 42 99)
```

**Parameters:**
- `condition` - If condition expression
- `then-expr` - Then branch expression  
- `else-expr` - Else branch expression

### 3. Macro Expansion

#### `expand-macros`
Recursively expand macros in expressions.

```lisp
;; Built-in macro expansion
(expand-macros '(when true (println "hello")) (make-context) 0)
;; → (if true (do (println "hello")) nil)

(expand-macros '(unless false (println "world")) (make-context) 0)
;; → (if false nil (do (println "world")))

;; Cond macro expansion
(expand-macros '(cond (< x 0) "negative" (> x 0) "positive" :else "zero") (make-context) 0)
;; → (if (< x 0) "negative" (if (> x 0) "positive" "zero"))
```

**Parameters:**
- `expr` - Expression to expand
- `ctx` (hash-map) - Compilation context  
- `depth` (number) - Current expansion depth

**Built-in macros:**
- `when` - Conditional execution with implicit do
- `unless` - Negated conditional execution
- `cond` - Multi-branch conditional

#### `is-macro?`
Check if a symbol refers to a macro.

```lisp
(is-macro? 'when (make-context))
;; → true

(is-macro? 'unless (make-context))
;; → true

(is-macro? '+ (make-context))
;; → false

(is-macro? 'my-custom-macro (make-context))
;; → false (unless defined in context)
```

## Special Form Compilation

### 1. Function Definition

#### `compile-fn`
Compile function definitions with lexical scoping.

```lisp
(compile-fn '([x] (* x x)) (make-context))
;; → (fn [x] (* x x))

(compile-fn '([x y] (+ x y) (- x y)) (make-context))
;; → (fn [x y] (+ x y) (- x y))

;; With optimizations
(compile-fn '([x] (+ 1 2 x)) (make-context))
;; → (fn [x] (+ 3 x)) ; constant folding applied
```

**Features:**
- **Parameter binding** in local context
- **Multiple body expressions** supported
- **Lexical scoping** for nested functions
- **Optimization** applied to function body

### 2. Variable Definition

#### `compile-def`
Compile global variable definitions.

```lisp
(compile-def '(PI 3.14159) (make-context))
;; → (def PI 3.14159)

(compile-def '(square (fn [x] (* x x))) (make-context))
;; → (def square (fn [x] (* x x)))

;; With constant folding
(compile-def '(answer (+ 40 2)) (make-context))
;; → (def answer 42)
```

**Process:**
1. Validate argument count and symbol name
2. Register symbol in global symbol table
3. Compile initialization expression
4. Return compiled def form

### 3. Conditional Compilation

#### `compile-if`
Compile conditional expressions with optimization.

```lisp
(compile-if '((> 5 3) "yes" "no") (make-context))
;; → "yes" (dead code elimination)

(compile-if '(flag result1 result2) (make-context))
;; → (if flag result1 result2)

;; Nested conditionals
(compile-if '((> x 0) (if (< x 10) "small" "big") "negative") (make-context))
;; → (if (> x 0) (if (< x 10) "small" "big") "negative")
```

### 4. Local Bindings

#### `compile-let`
Compile let expressions with local variable tracking.

```lisp
(compile-let '([x 5] (+ x 1)) (make-context))
;; → (let [x 5] (+ x 1))

(compile-let '([x 1 y 2] (+ x y) (* x y)) (make-context))
;; → (let [x 1 y 2] (+ x y) (* x y))

;; With optimizations
(compile-let '([x (+ 1 2)] (* x x)) (make-context))
;; → (let [x 3] (* x x)) ; constant folding in binding
```

## Utility Functions

### 1. Collection Utilities

The compiler includes self-hosted implementations of essential collection functions:

#### `map`
Transform collections with a function.

```lisp
(map (fn [x] (* x x)) '(1 2 3 4))
;; → (1 4 9 16)

(map constant-fold-expr '((+ 1 2) (* 3 4) (- 10 5)))
;; → (3 12 5)
```

#### `filter`
Filter collections by predicate.

```lisp
(filter (fn [x] (> x 0)) '(-2 -1 0 1 2))
;; → (1 2)

(filter constant? '(42 x "hello" (+ 1 2) :key))
;; → (42 "hello" :key)
```

#### `reduce`
Reduce collections with accumulator function.

```lisp
(reduce + 0 '(1 2 3 4))
;; → 10

(reduce (fn [acc item] (cons item acc)) '() '(1 2 3))
;; → (3 2 1) ; reverse
```

#### `any?`
Check if any element matches predicate.

```lisp
(any? (fn [x] (> x 10)) '(1 5 15 3))
;; → true

(any? constant? '(x y (+ 1 2) z))
;; → false
```

### 2. Helper Functions

#### `second`
Get second element of collection.

```lisp
(second '(a b c))
;; → b

(second [10 20 30])
;; → 20
```

#### `length`
Get collection length (alias for `count`).

```lisp
(length '(a b c))
;; → 3

(length [1 2 3 4 5])
;; → 5
```

#### `concat`
Concatenate two collections.

```lisp
(concat '(1 2) '(3 4))
;; → (1 2 3 4)

(concat [] [1 2 3])
;; → [1 2 3]
```

#### `str-join`
Join string representations with separator.

```lisp
(str-join ", " '(1 2 3))
;; → "1, 2, 3"

(str-join "\n" '("line1" "line2" "line3"))
;; → "line1\nline2\nline3"
```

## Advanced Usage Examples

### 1. Custom Optimization Pipeline

```lisp
;; Create context with specific optimizations
(def opt-ctx (make-context-with-optimizations 
               {:constant-folding true 
                :dead-code-elimination false}))

;; Compile with only constant folding
(compile-expr '(if true (+ 1 2 3) (* 4 5)) opt-ctx)
;; → (if true 6 20) ; constants folded, but if not eliminated
```

### 2. Macro Development and Testing

```lisp
;; Define custom macro
(defmacro my-when [condition & body]
  `(if ~condition (do ~@body) nil))

;; Test macro expansion
(expand-macros '(my-when true (println "hello") (println "world")) 
               (make-context) 0)
;; → (if true (do (println "hello") (println "world")) nil)
```

### 3. Multi-File Compilation

```lisp
;; Compile multiple related files
(defn compile-project [files output-dir]
  (map (fn [file]
         (let [output-file (str output-dir "/" file "-compiled.lisp")]
           (compile-file file output-file)
           output-file))
       files))

;; Usage
(compile-project '("utils.lisp" "core.lisp" "main.lisp") "compiled/")
```

### 4. Optimization Analysis

```lisp
;; Compare optimized vs unoptimized output
(defn compare-compilation [expr]
  (let [ctx-opt (make-context)
        ctx-no-opt (make-context-with-optimizations {})]
    {:optimized (compile-expr expr ctx-opt)
     :unoptimized (compile-expr expr ctx-no-opt)}))

;; Example
(compare-compilation '(if true (+ 1 2 3) (* 4 5)))
;; → {:optimized 6, :unoptimized (if true (+ 1 2 3) (* 4 5))}
```

## Error Handling

### Compilation Errors

The compiler provides detailed error messages for various error conditions:

```lisp
;; Arity errors
(compile-def '(x) (make-context))
;; → Error: "def requires exactly 2 arguments"

(compile-if '(condition) (make-context))
;; → Error: "if requires exactly 3 arguments"

;; Type errors
(compile-def '(123 value) (make-context))
;; → Error: "def name must be a symbol"

;; Macro expansion depth
(defmacro infinite-macro [x] `(infinite-macro ~x))
(expand-macros '(infinite-macro foo) (make-context) 0)
;; → Error: "Maximum macro expansion depth exceeded: 20"
```

### Defensive Programming

The compiler includes extensive validation and error checking:

```lisp
;; Nil checking
(compile-expr nil (make-context))
;; → nil (handled gracefully)

;; Empty list handling
(compile-expr '() (make-context))
;; → () (handled gracefully)

;; Invalid optimization flags
(optimization-enabled? nil :constant-folding)
;; → nil (safe failure)
```

## Performance Considerations

### 1. Compilation Speed

- **Optimization overhead**: Multiple passes increase compile time
- **Macro expansion**: Recursive expansion can be expensive
- **Context copying**: Functional approach requires context duplication

### 2. Memory Usage

- **Immutable structures**: All data structures are immutable
- **Context sharing**: Contexts share immutable components
- **Tail recursion**: Optimization passes use tail-recursive patterns

### 3. Output Quality

- **Constant folding**: Reduces runtime computation overhead
- **Dead code elimination**: Produces smaller, cleaner output
- **Semantic preservation**: Optimizations maintain program behavior

## Integration with GoLisp Core

### Required Core Functions

The compiler depends on these Go-implemented primitives:

```lisp
;; Data manipulation
cons first rest nth count empty? conj list vector hash-map

;; Type predicates  
symbol? string? number? list? vector? hash-map? keyword? fn? nil?

;; Meta-programming
eval read-string read-all-string macroexpand

;; I/O operations
slurp spit println prn

;; Arithmetic and comparison
+ - * / = < > <= >=

;; Hash-map operations
get assoc dissoc contains?
```

### Self-Hosted Components

These functions are implemented entirely in Lisp:

```lisp
;; Collection utilities
map filter reduce any? concat reverse flatten1

;; Helper functions
second length not= str-join

;; Optimization functions
constant? constant-fold-expr eliminate-dead-code

;; Compilation functions
compile-expr compile-symbol compile-list compile-file

;; Bootstrap functions
bootstrap-self-hosting read-all
```

This API provides a complete interface for compiling, optimizing, and analyzing Lisp code using the self-hosting compiler. The clean separation between Go core primitives and Lisp-implemented compiler logic demonstrates the power and elegance of the self-hosting approach.