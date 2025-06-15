# Modular Architecture Usage Examples

This document demonstrates how to use and extend the new modular architecture in go-lisp.

## Basic Usage

### Creating a Modular Evaluator

```go
package main

import (
    "fmt"
    "github.com/leinonen/lisp-interpreter/pkg/evaluator"
    "github.com/leinonen/lisp-interpreter/pkg/modular"
)

func main() {
    // Create environment
    env := evaluator.NewEnvironment()
    
    // Create modular evaluator (automatically loads core plugins)
    modularEval, err := modular.NewModularEvaluator(env)
    if err != nil {
        panic(err)
    }
    
    // Use arithmetic operations from plugins
    result, err := modularEval.Eval(parseString("(+ 1 2 3)"))
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result) // Output: 6
}
```

### Loading Additional Plugins

```go
// Load comparison plugin
comparisonPlugin := comparison.NewComparisonPlugin()
err := modularEval.LoadPlugin(comparisonPlugin)
if err != nil {
    panic(err)
}

// Now you can use comparison operations
result, err := modularEval.Eval(parseString("(< 1 2 3)"))
fmt.Println(result) // Output: true
```

## Plugin Management

### Listing Loaded Plugins

```go
plugins := modularEval.ListPlugins()
for _, plugin := range plugins {
    fmt.Printf("Plugin: %s v%s - %s\n", 
        plugin.Name, plugin.Version, plugin.Description)
    fmt.Printf("  Functions: %v\n", plugin.Functions)
}
```

### Listing Functions by Category

```go
// List all arithmetic functions
arithmeticFuncs := modularEval.ListFunctionsByCategory("arithmetic")
fmt.Println("Arithmetic functions:", arithmeticFuncs)

// List all comparison functions
comparisonFuncs := modularEval.ListFunctionsByCategory("comparison")
fmt.Println("Comparison functions:", comparisonFuncs)
```

### Getting Function Help

```go
help, exists := modularEval.GetFunctionHelp("+")
if exists {
    fmt.Println("Help for +:", help)
}
```

## Creating Custom Plugins

### Simple Plugin Example

```go
package mymath

import (
    "math"
    "github.com/leinonen/lisp-interpreter/pkg/functions"
    "github.com/leinonen/lisp-interpreter/pkg/plugins"
    "github.com/leinonen/lisp-interpreter/pkg/registry"
    "github.com/leinonen/lisp-interpreter/pkg/types"
)

type MyMathPlugin struct {
    *plugins.BasePlugin
}

func NewMyMathPlugin() *MyMathPlugin {
    return &MyMathPlugin{
        BasePlugin: plugins.NewBasePlugin(
            "mymath",
            "1.0.0",
            "Custom math functions",
            []string{"arithmetic"}, // Depends on arithmetic plugin
        ),
    }
}

func (mp *MyMathPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
    // Square function
    squareFunc := functions.NewFunction(
        "square",
        "math",
        1, // Exactly 1 argument
        "Square a number: (square 5) => 25",
        mp.evalSquare,
    )
    
    return reg.Register(squareFunc)
}

func (mp *MyMathPlugin) evalSquare(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
    values, err := functions.EvalArgs(evaluator, args)
    if err != nil {
        return nil, err
    }
    
    num, err := functions.ExtractFloat64(values[0])
    if err != nil {
        return nil, err
    }
    
    result := num * num
    return types.NumberValue(result), nil
}
```

### Plugin with Dependencies

```go
type AdvancedMathPlugin struct {
    *plugins.BasePlugin
}

func NewAdvancedMathPlugin() *AdvancedMathPlugin {
    return &AdvancedMathPlugin{
        BasePlugin: plugins.NewBasePlugin(
            "advancedmath",
            "1.0.0",
            "Advanced mathematical functions",
            []string{"arithmetic", "comparison", "mymath"}, // Multiple dependencies
        ),
    }
}

func (amp *AdvancedMathPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
    // Factorial function
    factFunc := functions.NewFunction(
        "factorial",
        "math",
        1,
        "Calculate factorial: (factorial 5) => 120",
        amp.evalFactorial,
    )
    
    return reg.Register(factFunc)
}

func (amp *AdvancedMathPlugin) evalFactorial(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
    values, err := functions.EvalArgs(evaluator, args)
    if err != nil {
        return nil, err
    }
    
    num, err := functions.ExtractFloat64(values[0])
    if err != nil {
        return nil, err
    }
    
    if num < 0 {
        return nil, fmt.Errorf("factorial of negative number")
    }
    
    result := 1.0
    for i := 2; i <= int(num); i++ {
        result *= float64(i)
    }
    
    return types.NumberValue(result), nil
}
```

## Plugin Lifecycle Management

### Loading and Unloading Plugins

```go
// Load a plugin
myPlugin := mymath.NewMyMathPlugin()
err := modularEval.LoadPlugin(myPlugin)
if err != nil {
    fmt.Printf("Failed to load plugin: %v\n", err)
}

// Use plugin functions
result, err := modularEval.Eval(parseString("(square 5)"))
fmt.Println(result) // Output: 25

// Unload the plugin
err = modularEval.UnloadPlugin("mymath")
if err != nil {
    fmt.Printf("Failed to unload plugin: %v\n", err)
}
```

### Plugin Dependency Resolution

The plugin system automatically checks dependencies:

```go
// This will fail if arithmetic plugin is not loaded
advancedPlugin := NewAdvancedMathPlugin()
err := modularEval.LoadPlugin(advancedPlugin)
if err != nil {
    fmt.Printf("Dependency error: %v\n", err)
}
```

## Integration with Existing Code

### Fallback to Original Evaluator

The modular evaluator automatically falls back to the original evaluator for functions that aren't registered as plugins:

```go
// This will use the original evaluator's implementation
result, err := modularEval.Eval(parseString("(defn my-func [x] (+ x 1))"))
```

### Gradual Migration

You can gradually migrate existing built-in functions to plugins:

1. Create a plugin for a category of functions (e.g., string operations)
2. Load the plugin alongside the original evaluator
3. Registered plugin functions take precedence
4. Unregistered functions fall back to the original evaluator

## Best Practices

### Plugin Design

1. **Single Responsibility**: Each plugin should focus on one area of functionality
2. **Clear Dependencies**: Explicitly declare plugin dependencies
3. **Good Documentation**: Provide helpful function descriptions
4. **Error Handling**: Return clear error messages
5. **Type Safety**: Use the helper functions for type extraction

### Function Implementation

```go
func (p *MyPlugin) myFunction(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
    // 1. Validate arguments
    if err := functions.ValidateExactArity("my-func", 2, args); err != nil {
        return nil, err
    }
    
    // 2. Evaluate arguments
    values, err := functions.EvalArgs(evaluator, args)
    if err != nil {
        return nil, err
    }
    
    // 3. Extract typed values
    str, err := functions.ExtractString(values[0])
    if err != nil {
        return nil, fmt.Errorf("my-func first argument: %v", err)
    }
    
    num, err := functions.ExtractFloat64(values[1])
    if err != nil {
        return nil, fmt.Errorf("my-func second argument: %v", err)
    }
    
    // 4. Perform operation
    result := fmt.Sprintf("%s: %.2f", str, num)
    
    // 5. Return result
    return types.StringValue(result), nil
}
```

### Testing Plugins

```go
func TestMyPlugin(t *testing.T) {
    env := evaluator.NewEnvironment()
    modularEval, err := modular.NewModularEvaluator(env)
    if err != nil {
        t.Fatal(err)
    }
    
    // Load your plugin
    plugin := NewMyPlugin()
    err = modularEval.LoadPlugin(plugin)
    if err != nil {
        t.Fatal(err)
    }
    
    // Test plugin functions
    result, err := modularEval.Eval(parseString("(my-func \"test\" 42)"))
    if err != nil {
        t.Fatal(err)
    }
    
    if result.String() != "test: 42.00" {
        t.Errorf("Expected 'test: 42.00', got %s", result.String())
    }
}
```

## Benefits of the Modular Architecture

1. **Easy Extension**: Add new functionality without modifying core code
2. **Hot-Pluggable**: Load/unload features at runtime
3. **Dependency Management**: Automatic dependency resolution
4. **Better Organization**: Related functions grouped together
5. **Testing**: Test plugins in isolation
6. **Performance**: Load only needed functionality
7. **Community**: Enable third-party plugin development
8. **Backwards Compatibility**: Existing code continues to work
