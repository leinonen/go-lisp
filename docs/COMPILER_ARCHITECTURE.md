# GoLisp Self-Hosting Compiler Architecture

## Overview

The GoLisp self-hosting compiler is a multi-pass compiler written entirely in Lisp that compiles Lisp source code. It demonstrates true self-hosting capability where the language is sophisticated enough to implement its own compiler.

## Architecture Principles

### 1. Minimal Core Foundation
- **Go Core (3,439 lines)**: Essential primitives only (~50 functions)
- **Self-Hosted Everything Else**: All higher-level functionality implemented in Lisp
- **Modular Design**: Clean separation between core primitives and self-hosted components

### 2. Multi-Pass Compilation Pipeline
```
Source Code → Macro Expansion → Constant Folding → Core Compilation → Dead Code Elimination → Output
```

### 3. Context-Driven Compilation
- **Compilation Context**: Tracks symbols, locals, macros, and optimization flags
- **Lexical Scoping**: Proper variable resolution through context chain
- **Optimization Control**: Fine-grained control over optimization passes

## Core Components

### 1. Data Structures (`pkg/core/types.go`)

#### Value Interface
All Lisp values implement the `Value` interface:
- `String()` method for uniform string representation
- Type-specific behavior through interface methods

#### Core Types
- **Symbol**: Named identifiers (`x`, `foo`, `+`)
- **Number**: Integers and floats (`42`, `3.14`)
- **String**: Text literals (`"hello"`)
- **Keyword**: Clojure-style keywords (`:key`)
- **List**: Linked lists (`(1 2 3)`)
- **Vector**: Indexed collections (`[1 2 3]`)
- **HashMap**: Key-value mappings (`{:a 1 :b 2}`)
- **Set**: Unique collections (`#{1 2 3}`)
- **Function**: Built-in and user-defined functions
- **Nil/Boolean**: `nil`, `true`, `false`

### 2. Parser (`pkg/core/reader.go`)

#### Lexer
- **Token-based parsing** with position tracking
- **Error reporting** with line/column information
- **Multi-expression support** via `read-all-string`

#### Parser Features
- **Recursive descent parsing** for nested structures
- **Quote syntax**: `'expr`, `` `expr ``, `~expr`, `~@expr`
- **Literal syntax**: Numbers, strings, keywords, collections
- **Error recovery** with descriptive messages

### 3. Evaluation Engine (Modular `eval_*.go`)

#### Core Evaluation (`eval_core.go`)
- **Environment chain** for lexical scoping
- **Special form dispatch** to specialized handlers
- **Context tracking** for error reporting

#### Specialized Modules
- **Arithmetic** (`eval_arithmetic.go`): `+`, `-`, `*`, `/`, comparisons
- **Collections** (`eval_collections.go`): `cons`, `first`, `rest`, collection operations
- **Strings** (`eval_strings.go`): String manipulation functions
- **I/O** (`eval_io.go`): File operations, printing
- **Meta-programming** (`eval_meta.go`): `eval`, `read-string`, `macroexpand`
- **Special Forms** (`eval_special_forms.go`): `if`, `fn`, `def`, `let`, etc.

## Self-Hosting Compiler (`lisp/self-hosting.lisp`)

### 1. Compilation Context

#### Context Structure
```lisp
{:symbols (hash-map)        ; Global symbol table
 :locals '()               ; Local variable stack
 :macros (hash-map)        ; Macro definitions
 :target 'eval             ; Compilation target
 :optimizations {...}}     ; Optimization flags
```

#### Context Functions
- `make-context` - Create default compilation context
- `make-context-with-optimizations` - Create context with custom optimizations
- `optimization-enabled?` - Check if optimization is enabled

### 2. Multi-Pass Compilation Pipeline

#### Pass 1: Macro Expansion
```lisp
(defn expand-macros [expr ctx depth]
  ;; Recursively expand macros with depth limiting
  ;; Supports built-in macros: when, unless, cond
  ;; User-defined macros tracked in context
  )
```

**Features:**
- **Depth limiting** prevents infinite recursion
- **Built-in macro support** for `when`, `unless`, `cond`
- **User-defined macros** tracked in compilation context
- **Recursive expansion** through all data structures

#### Pass 2: Constant Folding Optimization
```lisp
(defn constant-fold-expr [expr]
  ;; Evaluate compile-time constants
  ;; Fold arithmetic: (+ 1 2 3) → 6
  ;; Fold comparisons: (< 3 5) → true
  )
```

**Optimizations:**
- **Arithmetic folding**: `(+ 1 2 3)` → `6`
- **Comparison folding**: `(< 3 5)` → `true`
- **Nested expressions**: `(+ (* 2 3) (- 8 2))` → `12`
- **Type safety** ensures only valid operations are folded

#### Pass 3: Core Compilation
```lisp
(defn compile-expr [expr ctx]
  ;; Dispatch to specialized compilers based on expression type
  ;; Symbol resolution through context
  ;; Special form handling
  )
```

**Compilation Strategies:**
- **Symbol compilation**: Local vs. global resolution
- **Special forms**: `def`, `fn`, `if`, `let`, `do`, `quote`
- **Function application**: Head + arguments compilation
- **Data structures**: Vectors, lists with recursive compilation

#### Pass 4: Dead Code Elimination
```lisp
(defn eliminate-dead-code [expr]
  ;; Remove unreachable code branches
  ;; Simplify conditional expressions
  )
```

**Optimizations:**
- **If simplification**: `(if true x y)` → `x`
- **Branch elimination**: Remove unreachable else branches
- **Let optimization**: Remove unused bindings (future)

### 3. Special Form Compilation

#### Function Definition (`fn`)
```lisp
(defn compile-fn [args ctx]
  ;; Create new lexical scope
  ;; Track parameter bindings in context
  ;; Compile function body with extended context
  )
```

#### Variable Definition (`def`)
```lisp
(defn compile-def [args ctx]
  ;; Register symbol in global symbol table
  ;; Compile initialization expression
  ;; Track definition for optimization
  )
```

#### Conditional (`if`)
```lisp
(defn compile-if [args ctx]
  ;; Compile condition, then, else branches
  ;; Apply dead code elimination if enabled
  )
```

#### Local Bindings (`let`)
```lisp
(defn compile-let [args ctx]
  ;; Extract binding symbols
  ;; Create extended context with local bindings
  ;; Compile body with local context
  )
```

### 4. Optimization System

#### Optimization Flags
```lisp
{:constant-folding true
 :dead-code-elimination true
 :tail-call-optimization false}  ; Future
```

#### Configurable Optimizations
- **Per-compilation control** via context flags
- **Individual pass enabling/disabling**
- **Extensible architecture** for new optimizations

#### Current Optimizations
1. **Constant Folding**: Compile-time evaluation of constant expressions
2. **Dead Code Elimination**: Removal of unreachable code branches
3. **Macro Pre-expansion**: Early macro expansion before compilation

#### Future Optimizations
- **Tail Call Recognition**: Identify tail-recursive patterns
- **Function Inlining**: Inline small functions at call sites
- **Loop Optimization**: Enhance loop/recur patterns

## Bootstrap Process

### 1. Self-Compilation Workflow
```lisp
(defn bootstrap-self-hosting []
  ;; 1. Compile standard library with compiler
  ;; 2. Compile compiler with itself
  ;; 3. Verify compiled versions work identically
  )
```

### 2. File Compilation
```lisp
(defn compile-file [input-file output-file]
  ;; 1. Read all expressions from source file
  ;; 2. Compile each expression with shared context
  ;; 3. Write compiled output to target file
  )
```

### 3. Multi-Expression Support
```lisp
(defn read-all [source]
  ;; Parse multiple top-level expressions
  ;; Handle realistic Lisp programs
  ;; Support for libraries and modules
  )
```

## Error Handling and Diagnostics

### 1. Compilation Errors
- **Type checking** during compilation
- **Arity validation** for special forms
- **Symbol resolution** errors
- **Macro expansion** depth limits

### 2. Source Location Tracking
- **Parse errors** with exact line/column
- **Visual error indicators** showing error position
- **Stack trace support** for nested compilation

### 3. Defensive Programming
- **Nil checking** throughout compilation pipeline
- **Input validation** for all compilation functions
- **Graceful degradation** when optimizations fail

## Integration with Go Core

### 1. Core Primitive Dependencies
The compiler requires these Go-implemented primitives:
- **Data manipulation**: `cons`, `first`, `rest`, `nth`, `count`
- **Type predicates**: `symbol?`, `list?`, `vector?`, etc.
- **Meta-programming**: `eval`, `read-string`, `read-all-string`
- **I/O operations**: `slurp`, `spit`, `println`

### 2. Self-Hosted Functions
These functions are implemented in Lisp within the compiler:
- **Collection utilities**: `map`, `filter`, `reduce`, `any?`
- **Helper functions**: `second`, `length`, `concat`, `reverse`
- **Optimization algorithms**: All optimization passes
- **Compilation logic**: All compilation functions

### 3. Macro System Integration
- **Built-in macros**: `when`, `unless`, `cond` (implemented in stdlib)
- **User macros**: `defmacro` compilation and expansion
- **Macro expansion**: Recursive expansion with depth limits

## Design Patterns and Principles

### 1. Functional Programming
- **Immutable data structures** throughout compilation
- **Pure functions** for most compilation logic
- **Higher-order functions** for optimization passes

### 2. Context-Driven Design
- **Compilation context** threading through all passes
- **Lexical scope tracking** via context chain
- **Optimization control** through context flags

### 3. Modular Architecture
- **Separation of concerns** between passes
- **Pluggable optimizations** via flag system
- **Extensible design** for new language features

### 4. Error-First Programming
- **Early validation** of inputs and invariants
- **Descriptive error messages** with context
- **Graceful failure** modes for robustness

## Performance Characteristics

### 1. Compilation Speed
- **Multi-pass overhead** balanced by optimization benefits
- **Macro expansion caching** (future enhancement)
- **Incremental compilation** support (future)

### 2. Output Quality
- **Constant folding** reduces runtime computation
- **Dead code elimination** produces smaller output
- **Semantic preservation** ensures correctness

### 3. Memory Usage
- **Immutable structures** with GC-friendly patterns
- **Context sharing** reduces memory duplication
- **Tail recursion** in optimization passes

## Extensibility and Future Development

### 1. Adding New Optimizations
1. **Define optimization function** following existing patterns
2. **Add optimization flag** to context system
3. **Integrate into compilation pipeline**
4. **Add comprehensive tests**

### 2. Supporting New Language Features
1. **Extend core primitives** (if needed)
2. **Add special form compilation**
3. **Update macro expansion** (if applicable)
4. **Enhance error handling**

### 3. Code Generation Targets
- **Current**: Lisp → Lisp (AST transformation)
- **Future**: Lisp → Bytecode, Lisp → Go, Lisp → JavaScript

This architecture demonstrates that a relatively small core (3,439 lines of Go) can support a sophisticated, self-hosting compiler entirely implemented in the target language. The clean separation between minimal core primitives and self-hosted higher-level functionality enables rapid language evolution and optimization.