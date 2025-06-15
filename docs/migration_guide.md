# Migration Guide: From Monolithic to Modular Architecture

This guide shows how to migrate the existing go-lisp architecture to the new modular plugin system.

## Migration Strategy

### Phase 1: Core Infrastructure
✅ **Completed**
- Function registry system (`pkg/registry`)
- Plugin interface and manager (`pkg/plugins`) 
- Base function utilities (`pkg/functions`)
- Modular evaluator wrapper (`pkg/modular`)

### Phase 2: Core Function Plugins
🔄 **In Progress**
- ✅ Arithmetic plugin (`+`, `-`, `*`, `/`, `%`)
- ✅ Comparison plugin (`=`, `<`, `>`, `<=`, `>=`)
- 🔲 Logical operations plugin (`and`, `or`, `not`)
- 🔲 List operations plugin (`list`, `first`, `rest`, `cons`, etc.)
- 🔲 Control flow plugin (`if`, `do`)

### Phase 3: Advanced Function Plugins
✅ **Completed**
- String operations plugin (string-concat, string-length, string-contains?, etc.)
- Mathematical functions plugin (sqrt, pow, sin, cos, min, max, etc.)
- Hash map operations plugin (hash-map, hash-map-get, hash-map-put, etc.)
- Atom operations plugin (atom, deref, reset!, swap!)
- I/O operations plugin (print!, println!, read-file, write-file, file-exists?)
- HTTP client plugin (http-get, http-post, http-put, http-delete)
- JSON processing plugin (json-parse, json-stringify, json-path, etc.)
- Concurrency plugin (chan, chan-send!, chan-recv!, chan-close!, etc.)

### Phase 4: Integration and Migration
🔲 **Planned**
- Update main evaluator to use modular system
- Migrate existing tests
- Update documentation
- Performance optimization

## Step-by-Step Migration

### 1. Creating Function Category Plugins

Here's how to convert existing built-in functions to plugins:

#### Example: Converting List Operations

**Before (in evaluator.go):**
```go
case "list":
    return e.evalListConstruction(list.Elements[1:])
case "first":
    return e.evalFirst(list.Elements[1:])
case "rest":
    return e.evalRest(list.Elements[1:])
```

**After (as plugin):**
```go
// pkg/plugins/list/list.go
package list

type ListPlugin struct {
    *plugins.BasePlugin
}

func NewListPlugin() *ListPlugin {
    return &ListPlugin{
        BasePlugin: plugins.NewBasePlugin(
            "list",
            "1.0.0", 
            "List manipulation functions",
            []string{},
        ),
    }
}

func (lp *ListPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
    functions := []struct {
        name     string
        arity    int
        help     string
        handler  func(registry.Evaluator, []types.Expr) (types.Value, error)
    }{
        {"list", -1, "Create a list: (list 1 2 3) => (1 2 3)", lp.evalList},
        {"first", 1, "Get first element: (first '(1 2 3)) => 1", lp.evalFirst},
        {"rest", 1, "Get rest of list: (rest '(1 2 3)) => (2 3)", lp.evalRest},
        {"cons", 2, "Prepend element: (cons 0 '(1 2)) => (0 1 2)", lp.evalCons},
        {"length", 1, "Get list length: (length '(1 2 3)) => 3", lp.evalLength},
    }
    
    for _, fn := range functions {
        f := functions.NewFunction(fn.name, registry.CategoryList, fn.arity, fn.help, fn.handler)
        if err := reg.Register(f); err != nil {
            return err
        }
    }
    return nil
}

func (lp *ListPlugin) evalList(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
    values, err := functions.EvalArgs(evaluator, args)
    if err != nil {
        return nil, err
    }
    return &types.ListValue{Elements: values}, nil
}

// ... implement other functions
```

### 2. Migrating Complex Functions

For functions that depend on internal evaluator state:

**Before:**
```go
func (e *Evaluator) evalIf(args []types.Expr) (types.Value, error) {
    if len(args) < 2 || len(args) > 3 {
        return nil, fmt.Errorf("if requires 2 or 3 arguments")
    }
    
    condition, err := e.Eval(args[0])
    if err != nil {
        return nil, err
    }
    
    if isTruthy(condition) {
        return e.Eval(args[1])  // Access to e.Eval
    } else if len(args) == 3 {
        return e.Eval(args[2])
    }
    return types.BooleanValue(false), nil
}
```

**After:**
```go
func (cp *ControlPlugin) evalIf(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
    if len(args) < 2 || len(args) > 3 {
        return nil, fmt.Errorf("if requires 2 or 3 arguments")
    }
    
    condition, err := evaluator.Eval(args[0])  // Use passed evaluator
    if err != nil {
        return nil, err
    }
    
    if functions.IsTruthy(condition) {
        return evaluator.Eval(args[1])
    } else if len(args) == 3 {
        return evaluator.Eval(args[2])
    }
    return types.BooleanValue(false), nil
}
```

### 3. Handling Special Forms

Special forms that need access to unevaluated arguments:

```go
func (fp *FunctionPlugin) evalQuote(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("quote requires exactly 1 argument")
    }
    
    // Don't evaluate the argument - return it as a quoted value
    return &types.QuotedValue{Value: args[0]}, nil
}
```

### 4. Migration Checklist

For each group of functions being migrated:

- [ ] Create new plugin package
- [ ] Implement plugin interface
- [ ] Move function logic to plugin methods
- [ ] Update function signatures to use registry.Evaluator
- [ ] Add proper error handling
- [ ] Create comprehensive tests
- [ ] Update documentation
- [ ] Remove functions from old evaluator switch statement

## Backwards Compatibility Strategy

### 1. Hybrid Evaluator

Create a hybrid evaluator that supports both old and new systems:

```go
type HybridEvaluator struct {
    original      *evaluator.Evaluator
    modular       *modular.ModularEvaluator
    useModular    map[string]bool  // Which functions to use modular for
}

func (he *HybridEvaluator) Eval(expr types.Expr) (types.Value, error) {
    // Try modular first, fall back to original
    if he.shouldUseModular(expr) {
        result, err := he.modular.Eval(expr)
        if err == nil {
            return result, nil
        }
    }
    return he.original.Eval(expr)
}
```

### 2. Gradual Function Migration

Enable specific functions to use the modular system:

```go
// Enable arithmetic functions to use plugins
hybridEval.EnableModular("+", "-", "*", "/", "%")

// Later, enable more functions
hybridEval.EnableModular("=", "<", ">", "<=", ">=")
```

### 3. Configuration-Driven Migration

```go
type EvaluatorConfig struct {
    UsePlugins   []string  `json:"use_plugins"`
    DisablePlugins []string `json:"disable_plugins"`
}

func NewConfiguredEvaluator(config EvaluatorConfig) *HybridEvaluator {
    // Configure which plugins to use based on config
}
```

## Performance Considerations

### 1. Function Lookup Optimization

```go
// Cache frequently used functions
type CachedRegistry struct {
    registry registry.FunctionRegistry
    cache    map[string]registry.BuiltinFunction
    mutex    sync.RWMutex
}

func (cr *CachedRegistry) Get(name string) (registry.BuiltinFunction, bool) {
    cr.mutex.RLock()
    if fn, exists := cr.cache[name]; exists {
        cr.mutex.RUnlock()
        return fn, true
    }
    cr.mutex.RUnlock()
    
    // Cache miss - get from registry and cache it
    fn, exists := cr.registry.Get(name)
    if exists {
        cr.mutex.Lock()
        cr.cache[name] = fn
        cr.mutex.Unlock()
    }
    return fn, exists
}
```

### 2. Plugin Loading Optimization

```go
// Lazy plugin loading
type LazyPluginManager struct {
    plugins     map[string]func() plugins.Plugin
    loaded      map[string]plugins.Plugin
    registry    registry.FunctionRegistry
}

func (lpm *LazyPluginManager) RegisterLazy(name string, factory func() plugins.Plugin) {
    lpm.plugins[name] = factory
}

func (lpm *LazyPluginManager) Get(funcName string) (registry.BuiltinFunction, bool) {
    // Check if function is already loaded
    if fn, exists := lpm.registry.Get(funcName); exists {
        return fn, true
    }
    
    // Try to load plugin that might contain this function
    if plugin := lpm.findPluginForFunction(funcName); plugin != nil {
        lpm.LoadPlugin(plugin)
        return lpm.registry.Get(funcName)
    }
    
    return nil, false
}
```

## Testing Strategy

### 1. Plugin Unit Tests

Each plugin should have comprehensive unit tests:

```go
func TestArithmeticPlugin(t *testing.T) {
    plugin := arithmetic.NewArithmeticPlugin()
    registry := registry.NewRegistry()
    
    err := plugin.RegisterFunctions(registry)
    if err != nil {
        t.Fatal(err)
    }
    
    // Test function registration
    if !registry.Has("+") {
        t.Error("+ function not registered")
    }
    
    // Test function execution
    // ... test each function
}
```

### 2. Integration Tests

Test the complete modular system:

```go
func TestModularIntegration(t *testing.T) {
    env := evaluator.NewEnvironment()
    modularEval, err := modular.NewModularEvaluator(env)
    if err != nil {
        t.Fatal(err)
    }
    
    // Test complex expressions that use multiple plugins
    testCases := []struct {
        input    string
        expected string
    }{
        {"(+ (* 2 3) (- 8 3))", "11"},
        {"(< (+ 1 2) (* 2 3))", "true"},
        // ... more complex cases
    }
    
    for _, tc := range testCases {
        result, err := modularEval.Eval(parseString(tc.input))
        if err != nil {
            t.Errorf("Error evaluating %s: %v", tc.input, err)
        }
        if result.String() != tc.expected {
            t.Errorf("Expected %s, got %s for input %s", tc.expected, result.String(), tc.input)
        }
    }
}
```

### 3. Compatibility Tests

Ensure old code still works:

```go
func TestBackwardsCompatibility(t *testing.T) {
    // Test that existing example files still work
    files := []string{
        "examples/basic_features.lisp",
        "examples/advanced_features.lisp",
        // ... more files
    }
    
    for _, file := range files {
        t.Run(file, func(t *testing.T) {
            // Test with both old and new evaluators
            testWithOldEvaluator(t, file)
            testWithModularEvaluator(t, file)
        })
    }
}
```

## Benefits Realized

After migration, you'll have:

### 1. Easier Extension
```go
// Adding new functionality is now trivial
type DatabasePlugin struct {
    *plugins.BasePlugin
}

func (dp *DatabasePlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
    return reg.Register(functions.NewFunction(
        "db-query", "database", 2,
        "Execute database query: (db-query conn sql)",
        dp.evalQuery,
    ))
}
```

### 2. Better Organization
```
pkg/plugins/
├── arithmetic/     # Math operations
├── comparison/     # Comparison operators  
├── list/          # List operations
├── string/        # String functions
├── io/            # File operations
├── http/          # HTTP client
└── database/      # Database connectivity (3rd party)
```

### 3. Selective Loading
```go
// Load only needed functionality
modularEval := modular.NewMinimalEvaluator()
modularEval.LoadPlugin(arithmetic.NewArithmeticPlugin())
modularEval.LoadPlugin(comparison.NewComparisonPlugin())
// Skip heavy plugins like HTTP, JSON for simple scripts
```

### 4. Community Extensions
```go
// Third-party plugins can be easily integrated
import "github.com/user/lisp-graphics-plugin"

graphicsPlugin := graphics.NewGraphicsPlugin()
modularEval.LoadPlugin(graphicsPlugin)
```

## Next Steps

1. **Complete Core Plugins**: Finish implementing all existing built-in functions as plugins
2. **Update Main Evaluator**: Integrate modular system into main interpreter
3. **Performance Testing**: Benchmark and optimize the plugin system
4. **Documentation**: Update all documentation to reflect new architecture
5. **Community**: Create plugin SDK and examples for third-party developers
