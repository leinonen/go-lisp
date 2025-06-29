# GoLisp Self-Hosting Development Guide

## Overview

This guide provides comprehensive instructions for developers working on the GoLisp self-hosting compiler. It covers development workflows, testing strategies, contribution guidelines, and best practices for extending the self-hosting capabilities.

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Understanding the Codebase](#understanding-the-codebase)
3. [Development Workflow](#development-workflow)
4. [Testing and Validation](#testing-and-validation)
5. [Adding New Features](#adding-new-features)
6. [Debugging and Troubleshooting](#debugging-and-troubleshooting)
7. [Performance Optimization](#performance-optimization)
8. [Contributing Guidelines](#contributing-guidelines)

## Development Environment Setup

### Prerequisites

- **Go 1.19+** for the core implementation
- **Make** for build automation
- **Git** for version control
- **Text editor** with Lisp syntax highlighting (recommended: VS Code, Emacs, Vim)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/your-org/go-lisp.git
cd go-lisp

# Build the interpreter
make build

# Run tests to verify setup
make test

# Start REPL with self-hosting compiler
./bin/golisp -f lisp/self-hosting.lisp
```

### Development Commands

```bash
# Build and run REPL
make run

# Run specific test suites
make test-core          # Core Go tests only
make test-nocache       # All tests without cache
make test-core-nocache  # Core tests without cache

# Format code
make fmt

# Build binary
make build

# Clean build artifacts
make clean
```

## Understanding the Codebase

### Architecture Overview

```
go-lisp/
├── pkg/core/           # Go core implementation (3,439 lines)
│   ├── types.go        # Core data types and Value interface
│   ├── reader.go       # Lexer and parser with error reporting
│   ├── eval_*.go       # Modular evaluation engine
│   ├── repl.go         # REPL implementation
│   └── bootstrap.go    # Standard library loader
├── cmd/golisp/         # CLI entry point
├── lisp/               # Self-hosted Lisp code
│   ├── stdlib/         # Standard library in Lisp
│   └── self-hosting.lisp # Self-hosting compiler
└── docs/               # Documentation
```

### Key Design Principles

1. **Minimal Go Core**: Only essential primitives in Go (~50 functions)
2. **Self-Hosted Everything Else**: Higher-level functionality in Lisp
3. **Modular Architecture**: Clean separation of concerns
4. **Comprehensive Testing**: Extensive test coverage for reliability
5. **Context-Driven Compilation**: Threading context through all passes

### Core vs. Self-Hosted Components

#### Go Core Primitives (pkg/core/)
- **Data Types**: Symbol, List, Vector, HashMap, Set, etc.
- **Arithmetic**: `+`, `-`, `*`, `/`, `=`, `<`, `>`, etc.
- **Collections**: `cons`, `first`, `rest`, `nth`, `count`, `conj`
- **I/O**: `slurp`, `spit`, `println`, `file-exists?`
- **Meta**: `eval`, `read-string`, `read-all-string`, `macroexpand`

#### Self-Hosted Components (lisp/)
- **Standard Library**: `map`, `filter`, `reduce`, `sort`, `apply`
- **Compiler**: All compilation logic in `self-hosting.lisp`
- **Optimizations**: Constant folding, dead code elimination
- **Macros**: `when`, `unless`, `cond`

## Development Workflow

### 1. Setting Up Your Development Environment

```bash
# Fork and clone
git clone https://github.com/your-username/go-lisp.git
cd go-lisp

# Create development branch
git checkout -b feature/your-feature-name

# Set up pre-commit hooks (optional)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 2. Making Changes

#### Modifying Go Core
```bash
# Edit core files
vim pkg/core/eval_collections.go

# Run tests frequently
make test-core

# Format code
make fmt

# Verify build
make build
```

#### Modifying Self-Hosted Code
```bash
# Edit Lisp files
vim lisp/self-hosting.lisp

# Test changes interactively
./bin/golisp -f lisp/self-hosting.lisp

# Run full test suite
make test
```

### 3. Interactive Development

#### REPL-Driven Development

```bash
# Start REPL
./bin/golisp

# Load development files
(load-file "lisp/self-hosting.lisp")

# Test individual functions
(make-context)
(compile-expr '(+ 1 2 3) (make-context))

# Experiment with optimizations
(constant-fold-expr '(* (+ 2 3) (- 8 2)))
```

#### Rapid Prototyping

```lisp
;; Add temporary test functions
(defn test-new-feature []
  (println "Testing new feature...")
  ;; Your experimental code here
  )

;; Test immediately
(test-new-feature)
```

### 4. Testing Changes

#### Unit Testing
```bash
# Run specific test files
go test ./pkg/core -run TestNewFeature

# Run with verbose output
go test -v ./pkg/core

# Run with coverage
go test -cover ./pkg/core
```

#### Integration Testing
```bash
# Test self-hosting compiler
./bin/golisp -e "(load-file \"lisp/self-hosting.lisp\") (bootstrap-self-hosting)"

# Test specific compilation scenarios
./bin/golisp -e "(load-file \"lisp/self-hosting.lisp\") (compile-expr '(let [x 1] (+ x 2)) (make-context))"
```

## Testing and Validation

### Test Categories

#### 1. Go Core Tests (`pkg/core/*_test.go`)
- **Unit tests** for individual functions
- **Integration tests** for component interaction
- **Error handling** tests
- **Performance** benchmarks

#### 2. Self-Hosting Compiler Tests
- **Compilation correctness** tests
- **Optimization validation** tests
- **Bootstrap process** tests
- **Error condition** tests

#### 3. End-to-End Tests
- **File compilation** tests
- **Multi-expression** parsing tests
- **REPL integration** tests

### Writing Tests

#### Go Test Example
```go
func TestNewFeature(t *testing.T) {
    // Setup
    env := NewEnvironment()
    
    // Test cases
    testCases := []struct {
        input    string
        expected interface{}
        hasError bool
    }{
        {"(new-feature 1 2)", 3, false},
        {"(new-feature)", nil, true},
    }
    
    for _, tc := range testCases {
        result, err := EvalString(tc.input, env)
        if tc.hasError {
            assert.Error(t, err)
        } else {
            assert.NoError(t, err)
            assert.Equal(t, tc.expected, result)
        }
    }
}
```

#### Lisp Test Example
```lisp
;; Add to self-hosting.lisp or separate test file
(defn test-compilation-feature []
  (let [ctx (make-context)
        input '(your-new-feature 1 2)
        result (compile-expr input ctx)
        expected '(expected-output 1 2)]
    (if (= result expected)
      (println "✓ Test passed")
      (do
        (println "✗ Test failed")
        (println "  Input:" input)
        (println "  Expected:" expected)
        (println "  Got:" result)))))
```

### Test-Driven Development

1. **Write failing test** for new feature
2. **Implement minimal** functionality to pass
3. **Refactor** while keeping tests green
4. **Add edge cases** and error conditions
5. **Document** the new functionality

### Continuous Integration

```bash
# Pre-commit checks
make fmt
make test
make build

# Comprehensive validation
make test-nocache
go vet ./...
golangci-lint run
```

## Adding New Features

### 1. Adding Core Primitives (Go)

#### Step 1: Define the Function
```go
// In appropriate eval_*.go file
func evalNewPrimitive(args []Value, env *Environment) (Value, error) {
    // Validate arguments
    if len(args) != 2 {
        return nil, fmt.Errorf("new-primitive expects 2 arguments, got %d", len(args))
    }
    
    // Implementation
    // ...
    
    return result, nil
}
```

#### Step 2: Register in Evaluator
```go
// In eval_core.go or appropriate file
func init() {
    BuiltinFunctions["new-primitive"] = evalNewPrimitive
}
```

#### Step 3: Add Tests
```go
func TestNewPrimitive(t *testing.T) {
    // Comprehensive test cases
}
```

#### Step 4: Update Documentation
```go
// Add to CLAUDE.md core primitives list
// Update API documentation
```

### 2. Adding Self-Hosted Functions (Lisp)

#### Step 1: Implement in Lisp
```lisp
;; Add to appropriate stdlib file or self-hosting.lisp
(defn new-self-hosted-function [args]
  "Documentation string"
  ;; Implementation using core primitives
  )
```

#### Step 2: Test Interactively
```lisp
;; In REPL
(new-self-hosted-function test-args)
```

#### Step 3: Add Formal Tests
```lisp
;; Add test function
(defn test-new-self-hosted-function []
  ;; Test cases
  )
```

#### Step 4: Integration
```lisp
;; Ensure function is loaded in bootstrap process
;; Add to appropriate load-file chain
```

### 3. Adding Compiler Optimizations

#### Step 1: Define Optimization Function
```lisp
(defn new-optimization-pass [expr]
  (cond
    ;; Pattern matching for optimization opportunities
    (and (list? expr) (= (first expr) 'target-pattern))
    ;; Apply optimization
    (optimize-pattern expr)
    
    ;; Recursive case for complex expressions
    (list? expr)
    (map new-optimization-pass expr)
    
    ;; Base case
    :else expr))
```

#### Step 2: Add Optimization Flag
```lisp
;; In make-context-with-optimizations
(defn make-context-with-optimizations [opt-flags]
  (hash-map :symbols (hash-map)
            :locals '()
            :macros (hash-map)
            :target *compile-target*
            :optimizations (merge {:constant-folding true
                                   :dead-code-elimination true
                                   :new-optimization false} ; Add here
                                  opt-flags)))
```

#### Step 3: Integrate into Pipeline
```lisp
;; In compile-expr
(defn compile-expr [expr ctx]
  (let [;; ... existing passes ...
        ;; New optimization pass
        new-opt-expr (if (optimization-enabled? ctx :new-optimization)
                       (new-optimization-pass compiled-expr)
                       compiled-expr)]
    new-opt-expr))
```

#### Step 4: Add Tests
```lisp
(defn test-new-optimization []
  (let [ctx (make-context-with-optimizations {:new-optimization true})
        input '(pattern-to-optimize 1 2)
        result (compile-expr input ctx)
        expected '(optimized-pattern 1 2)]
    ;; Test logic
    ))
```

### 4. Adding Special Forms

#### Step 1: Add to Core Evaluator
```go
// In eval_special_forms.go
func evalNewSpecialForm(args []Value, env *Environment) (Value, error) {
    // Special form implementation
    // Handle quote/unquote semantics
    // Manage environment bindings
    return result, nil
}
```

#### Step 2: Register Special Form
```go
// In eval_core.go
func init() {
    SpecialForms["new-special-form"] = evalNewSpecialForm
}
```

#### Step 3: Add Compiler Support
```lisp
;; In compile-list function
(defn compile-list [lst ctx]
  ;; ... existing cases ...
  (= head 'new-special-form) (compile-new-special-form args ctx)
  ;; ... rest of cases ...
  )

(defn compile-new-special-form [args ctx]
  ;; Compilation logic for special form
  )
```

## Debugging and Troubleshooting

### 1. Debugging Compilation Issues

#### Enable Verbose Compilation
```lisp
;; Add debug prints to compilation functions
(defn compile-expr [expr ctx]
  (println "Compiling:" expr)
  (let [result (do-compilation expr ctx)]
    (println "Result:" result)
    result))
```

#### Step-by-Step Debugging
```lisp
;; Test individual compilation passes
(def test-expr '(let [x (+ 1 2)] (* x x)))
(def ctx (make-context))

;; Pass 1: Macro expansion
(expand-macros test-expr ctx 0)

;; Pass 2: Constant folding
(constant-fold-expr test-expr)

;; Pass 3: Core compilation
(compile-expr-no-opt test-expr ctx)
```

#### Common Issues and Solutions

**Issue**: Infinite macro expansion
```lisp
;; Solution: Check macro depth limits
(def *max-macro-expansion-depth* 10) ; Reduce for debugging
```

**Issue**: Context corruption
```lisp
;; Solution: Validate context at each step
(defn validate-context [ctx]
  (and (hash-map? ctx)
       (contains? ctx :symbols)
       (contains? ctx :locals)))
```

**Issue**: Optimization breaking semantics
```lisp
;; Solution: Compare optimized vs unoptimized output
(defn debug-optimization [expr]
  (let [ctx-opt (make-context)
        ctx-no-opt (make-context-with-optimizations {})
        opt-result (compile-expr expr ctx-opt)
        no-opt-result (compile-expr expr ctx-no-opt)]
    (println "Original:" expr)
    (println "Optimized:" opt-result)
    (println "Unoptimized:" no-opt-result)
    (= opt-result no-opt-result)))
```

### 2. Debugging Go Core Issues

#### Using Go Debugger
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug specific test
dlv test ./pkg/core -- -test.run TestSpecificFunction

# Debug binary
dlv exec ./bin/golisp
```

#### Adding Debug Prints
```go
func evalFunction(args []Value, env *Environment) (Value, error) {
    fmt.Printf("DEBUG: evaluating function with %d args\n", len(args))
    defer fmt.Printf("DEBUG: function evaluation complete\n")
    
    // Function implementation
    return result, nil
}
```

#### Memory and Performance Profiling
```bash
# CPU profiling
go test -cpuprofile cpu.prof ./pkg/core
go tool pprof cpu.prof

# Memory profiling  
go test -memprofile mem.prof ./pkg/core
go tool pprof mem.prof

# Benchmarking
go test -bench=. ./pkg/core
```

### 3. REPL-Based Debugging

#### Interactive Debugging Session
```lisp
;; Load compiler
(load-file "lisp/self-hosting.lisp")

;; Set up debug environment
(def debug-ctx (make-context))

;; Test problematic expression
(def problem-expr '(complex-expression with many parts))

;; Debug step by step
(expand-macros problem-expr debug-ctx 0)
(constant-fold-expr problem-expr)
(compile-expr problem-expr debug-ctx)

;; Inspect intermediate states
(println "Context locals:" (:locals debug-ctx))
(println "Context symbols:" (:symbols debug-ctx))
```

#### Error Analysis
```lisp
;; Catch and analyze errors
(defn safe-compile [expr ctx]
  (try
    (compile-expr expr ctx)
    (catch Exception e
      (println "Compilation error:")
      (println "  Expression:" expr)
      (println "  Error:" e)
      nil)))
```

## Performance Optimization

### 1. Profiling and Measurement

#### Go Profiling
```bash
# Profile specific operations
go test -bench=BenchmarkCompilation ./pkg/core
go test -cpuprofile compilation.prof -bench=BenchmarkCompilation ./pkg/core
go tool pprof compilation.prof
```

#### Lisp Performance Measurement
```lisp
;; Simple timing function
(defn time-compilation [expr iterations]
  (let [start-time (current-time-millis)]
    (loop [i 0]
      (when (< i iterations)
        (compile-expr expr (make-context))
        (recur (+ i 1))))
    (let [end-time (current-time-millis)]
      (/ (- end-time start-time) iterations))))
```

### 2. Optimization Strategies

#### Compilation Speed Optimizations
- **Context reuse**: Share immutable context components
- **Memoization**: Cache compilation results for repeated expressions
- **Lazy evaluation**: Defer expensive operations
- **Tail recursion**: Use proper tail calls in optimization passes

#### Memory Usage Optimizations
- **Immutable structures**: Minimize data copying
- **Context sharing**: Reuse context components
- **Garbage collection**: Help GC with shorter-lived objects

#### Code Quality Optimizations
- **Constant folding**: More aggressive constant evaluation
- **Dead code elimination**: Better unreachable code detection
- **Function inlining**: Inline small functions
- **Loop optimization**: Better loop/recur patterns

### 3. Performance Testing

#### Benchmark Suite
```lisp
(defn benchmark-compilation []
  (let [test-cases '((+ 1 2 3)
                     (let [x 1] (+ x 2))
                     (if true 42 99)
                     (fn [x] (* x x)))]
    (map (fn [expr]
           (time-compilation expr 1000))
         test-cases)))
```

#### Memory Usage Testing
```go
func BenchmarkCompilationMemory(b *testing.B) {
    env := NewEnvironment()
    expr := "(let [x 1] (+ x 2))"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := EvalString(expr, env)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Contributing Guidelines

### 1. Code Style

#### Go Code Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Add comprehensive documentation
- Include error handling
- Write comprehensive tests

#### Lisp Code Style
- Use consistent indentation (2 spaces)
- Prefer descriptive function names
- Include docstrings for public functions
- Use meaningful variable names
- Follow functional programming patterns

### 2. Commit Guidelines

#### Commit Message Format
```
type(scope): short description

Longer description if needed.

- Specific change 1
- Specific change 2

Closes #issue-number
```

#### Commit Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

### 3. Pull Request Process

#### Before Submitting
1. **Run full test suite**: `make test`
2. **Format code**: `make fmt`
3. **Update documentation**: As needed
4. **Add tests**: For new functionality
5. **Verify build**: `make build`

#### PR Description Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added for new functionality
```

### 4. Review Process

#### Code Review Checklist
- **Functionality**: Does it work as intended?
- **Testing**: Adequate test coverage?
- **Performance**: Any performance implications?
- **Security**: No security vulnerabilities?
- **Documentation**: Well documented?
- **Style**: Follows project conventions?

#### Review Response
- Address all review comments
- Update tests if needed
- Maintain clean commit history
- Be responsive to feedback

## Advanced Topics

### 1. Self-Hosting Bootstrap Process

The bootstrap process demonstrates true self-hosting:

```lisp
(defn full-bootstrap []
  ;; 1. Load minimal core (Go primitives)
  ;; 2. Load self-hosted standard library
  (load-file "lisp/stdlib/core.lisp")
  (load-file "lisp/stdlib/enhanced.lisp")
  
  ;; 3. Load self-hosting compiler
  (load-file "lisp/self-hosting.lisp")
  
  ;; 4. Compile standard library with compiler
  (compile-file "lisp/stdlib/core.lisp" "stdlib-compiled.lisp")
  
  ;; 5. Compile compiler with itself
  (compile-file "lisp/self-hosting.lisp" "compiler-compiled.lisp")
  
  ;; 6. Verify compiled versions work
  (load-file "compiler-compiled.lisp")
  (compile-file "stdlib-compiled.lisp" "stdlib-recompiled.lisp"))
```

### 2. Extending the Language

#### Adding New Data Types
1. **Define in Go**: Add to `types.go`
2. **Add predicates**: Type checking functions
3. **Add operations**: Manipulation functions
4. **Update compiler**: Handle in compilation
5. **Add tests**: Comprehensive coverage

#### Adding New Special Forms
1. **Core implementation**: In `eval_special_forms.go`
2. **Compiler support**: In `self-hosting.lisp`
3. **Macro integration**: If needed
4. **Documentation**: Update guides
5. **Testing**: Edge cases and integration

### 3. Multi-Target Compilation

The compiler architecture supports multiple output targets:

```lisp
;; Future: Multiple compilation targets
(def compilation-targets
  {:lisp compile-to-lisp      ; Current
   :bytecode compile-to-bytecode ; Future
   :go compile-to-go           ; Future
   :js compile-to-js})         ; Future

(defn compile-with-target [expr target]
  (let [compiler (get compilation-targets target)]
    (compiler expr (make-context))))
```

This guide provides a comprehensive foundation for developing and extending the GoLisp self-hosting compiler. The modular architecture, extensive testing, and clean separation between core primitives and self-hosted functionality make it an excellent platform for language experimentation and development.