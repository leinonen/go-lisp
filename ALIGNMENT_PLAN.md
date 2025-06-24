# Alignment Plan: Go-Lisp Codebase → future.md Architecture

## Executive Summary

This plan outlines how to evolve the current Go-Lisp codebase to align with the minimal kernel architecture described in `future.md`. The goal is to create a clean, testable, and extensible Lisp implementation based on a "microkernel" approach.

## ✅ Phase 1: Minimal Kernel Implementation (COMPLETED)

**Status: ✅ COMPLETE**

We have successfully implemented a minimal Lisp kernel in `pkg/minimal/` that demonstrates the core architecture principles:

### Core Components Built:
- **Types System** (`types.go`): Symbol interning, Lists, Numbers, Booleans, Strings, Nil
- **Environment** (`env.go`): Lexical scoping with parent chain lookup
- **Evaluator** (`eval.go`): Core eval/apply logic with special form handling
- **Special Forms**: `quote`, `if`, `fn`, `define`, `do`
- **Bootstrap** (`bootstrap.go`): Higher-level functions implemented in Lisp itself
- **REPL** (`repl.go`): Simple tokenizer, parser, and evaluation loop
- **Tests** (`minimal_test.go`): Comprehensive test coverage

### Key Achievements:
- ✅ Minimal core (<200 lines of eval logic)
- ✅ Self-evaluating expressions
- ✅ Symbol interning
- ✅ Function closures
- ✅ Lexical scoping
- ✅ Bootstrap functions (list, first, rest, arithmetic, comparisons)
- ✅ Working REPL with basic parsing
- ✅ **Clojure-style square bracket `[param]` syntax for function parameters**
- ✅ Factorial recursion working
- ✅ **All tests passing** (fixed test evaluation issue)
- ✅ Test coverage for core functionality

### Demo Output:
```
$ go run cmd/minimal-lisp/main.go examples

=== Minimal Lisp Kernel Examples ===
   (quote hello) => hello
   (define x 42) => defined
   x => 42
   (if true "yes" "no") => "yes"
   (add 3 4) => 7
   (factorial 5) => 120
```

## 🔄 Phase 2: Enhanced Kernel (IN PROGRESS)

### 2.1 Macro System Implementation

**Goal**: Add code-as-data manipulation capabilities

```go
// Add to special forms
case Intern("quasiquote"):
    return specialQuasiquote(args, env)
case Intern("unquote"):
    return specialUnquote(args, env)
case Intern("defmacro"):
    return specialDefmacro(args, env)
```

**Benefits**:
- ✅ Code generation at runtime
- ✅ Language extensibility
- ✅ Syntax transformations
- ✅ Domain-specific languages

### 2.2 Advanced Data Structures

```go
// Vector type (indexed access)
type Vector struct {
    elements []Value
}

// HashMap type (key-value storage)
type HashMap struct {
    elements map[Value]Value
}

// Set type (unique elements)
type Set struct {
    elements map[Value]bool
}
```

### 2.3 Error Handling and Debugging

```go
// Enhanced error types with stack traces
type EvaluationError struct {
    Message string
    StackTrace []string
    SourceLocation Position
}

// Debug information
type DebugInfo struct {
    Breakpoints []Position
    StepMode bool
}
```

## 🔧 Phase 3: Integration with Existing Codebase

### 3.1 Gradual Migration Strategy

**Approach**: Don't replace existing code immediately. Instead, create integration points.

#### Step 1: Adapter Pattern
```go
// pkg/adapters/minimal_adapter.go
type MinimalAdapter struct {
    kernel *minimal.REPL
}

func (a *MinimalAdapter) EvaluateExpression(expr string) (types.Value, error) {
    // Convert existing types.Value to minimal.Value
    minimalExpr := a.convertToMinimal(expr)
    result, err := minimal.Eval(minimalExpr, a.kernel.Env)
    if err != nil {
        return nil, err
    }
    // Convert back to existing types.Value
    return a.convertFromMinimal(result), nil
}
```

#### Step 2: Feature Flags
```go
// pkg/config/config.go
type Config struct {
    UseMinimalKernel bool `env:"USE_MINIMAL_KERNEL" default:"false"`
    FeatureFlags     map[string]bool
}
```

#### Step 3: Side-by-Side Testing
```go
// pkg/testing/compatibility_test.go
func TestCompatibility(t *testing.T) {
    testCases := []string{
        "(+ 1 2 3)",
        "(if true 42 0)",
        "(define factorial (fn [n] ...))",
    }
    
    for _, test := range testCases {
        oldResult := oldEvaluator.Eval(test)
        newResult := minimalKernel.Eval(test)
        assert.Equal(t, oldResult, newResult)
    }
}
```

### 3.2 Plugin System Modernization

**Current Structure**:
```
pkg/plugins/
├── arithmetic/
├── control/
├── core/
├── functional/
└── ...
```

**New Structure**:
```
pkg/kernel/
├── minimal/          # Core kernel
├── stdlib/           # Standard library (implemented in Lisp)
├── extensions/       # Go-based extensions
└── compat/          # Compatibility layer
```

#### Migration Path:
1. **Core Functions** → Move to `stdlib/` as Lisp implementations
2. **Performance-Critical** → Keep as Go extensions
3. **Legacy** → Wrap in compatibility layer

### 3.3 Type System Unification

**Challenge**: Current code uses `pkg/types` extensively

**Solution**: Create type bridges

```go
// pkg/bridges/type_bridge.go
func ConvertValue(v types.Value) minimal.Value {
    switch val := v.(type) {
    case types.NumberValue:
        return minimal.Number(float64(val))
    case types.StringValue:
        return minimal.String(string(val))
    case types.BooleanValue:
        return minimal.Boolean(bool(val))
    case *types.ListValue:
        elements := make([]minimal.Value, len(val.Elements))
        for i, elem := range val.Elements {
            elements[i] = ConvertValue(elem)
        }
        return minimal.NewList(elements...)
    default:
        // Handle unknown types
        return minimal.String(v.String())
    }
}
```

## 🚀 Phase 4: Advanced Features

### 4.1 Module System

```lisp
;; stdlib/math.lisp
(defmodule math
  (export square cube factorial)
  
  (define square (fn [x] (* x x)))
  (define cube (fn [x] (* x x x)))
  (define factorial (fn [n] 
    (if (< n 2) 1 (* n (factorial (- n 1)))))))

;; Usage
(import math)
(math.square 5) ; => 25
```

### 4.2 Concurrency Support

```lisp
;; Built on Go's goroutines and channels
(define async-map (fn [f lst]
  (let [ch (channel (length lst))]
    (for-each (fn [x] 
      (go (send ch (f x)))) lst)
    (collect ch (length lst)))))
```

### 4.3 Performance Optimizations

- **Tail Call Optimization**: Detect and optimize recursive calls
- **Compilation**: JIT compilation for hot code paths
- **Memoization**: Automatic caching for pure functions

## 📊 Testing Strategy

### 4.1 Comprehensive Test Coverage

```go
// pkg/testing/integration_test.go
func TestFullCompatibility(t *testing.T) {
    // Run all existing example files through minimal kernel
    exampleFiles := []string{
        "examples/arithmetic_math.lisp",
        "examples/functional_programming.lisp",
        "examples/control_flow.lisp",
        // ... all examples
    }
    
    for _, file := range exampleFiles {
        t.Run(file, func(t *testing.T) {
            // Test with both old and new evaluators
            testFileCompatibility(t, file)
        })
    }
}
```

### 4.2 Performance Benchmarks

```go
func BenchmarkEvaluation(b *testing.B) {
    tests := []struct{
        name string
        expr string
    }{
        {"simple_arithmetic", "(+ 1 2 3)"},
        {"function_call", "(factorial 10)"},
        {"nested_calls", "(map square (range 100))"},
    }
    
    for _, test := range tests {
        b.Run("old_"+test.name, func(b *testing.B) {
            // Benchmark old evaluator
        })
        b.Run("minimal_"+test.name, func(b *testing.B) {
            // Benchmark minimal kernel
        })
    }
}
```

### 4.3 Regression Testing

```go
func TestRegressions(t *testing.T) {
    // Ensure existing functionality still works
    regressionTests := loadRegressionSuite()
    
    for _, test := range regressionTests {
        result := minimalKernel.Eval(test.Input)
        assert.Equal(t, test.Expected, result, test.Description)
    }
}
```

## 🎯 Benefits of This Architecture

### 1. **Maintainability**
- **Small Core**: Easy to understand and modify
- **Clear Separation**: Core logic separate from library functions
- **Self-Documenting**: Lisp implementations are readable

### 2. **Extensibility**
- **Macro System**: Users can extend the language
- **Module System**: Organized code reuse
- **Plugin Architecture**: Easy to add features

### 3. **Testing**
- **Isolated Components**: Easy to unit test
- **Minimal Dependencies**: Core has no external deps
- **Regression Safety**: Comprehensive test coverage

### 4. **Performance**
- **Optimizable Core**: Small surface area for optimization
- **Selective Compilation**: Compile only hot paths
- **Memory Efficiency**: Minimal object creation

## 📈 Migration Timeline

### Phase 1: ✅ COMPLETE (2 days)
- Minimal kernel implementation
- Basic REPL and examples
- Core test coverage

### Phase 2: 🔄 CURRENT (1 week)
- Macro system
- Advanced data structures
- Enhanced error handling

### Phase 3: 📅 PLANNED (2 weeks)
- Integration adapters
- Type system bridges
- Compatibility testing

### Phase 4: 📅 FUTURE (1 week)
- Module system
- Performance optimizations
- Full migration

## 🏁 Success Metrics

### Functional Goals:
- ✅ All existing examples run on minimal kernel
- ✅ Performance within 10% of current implementation
- ✅ 100% test coverage for core kernel
- ✅ Successful macro system implementation

### Architectural Goals:
- ✅ Core evaluator under 300 lines
- ✅ Clear separation of concerns
- ✅ Extensible design
- ✅ Self-hosting capability

## 🤝 Next Steps

1. **Continue Phase 2**: Implement macro system and advanced data structures
2. **Start Integration**: Begin building adapter layer for existing code
3. **Performance Testing**: Benchmark against current implementation
4. **Community Feedback**: Get input on the architectural direction

This plan provides a clear path forward that respects the existing codebase while moving toward the clean, minimal architecture envisioned in `future.md`.
