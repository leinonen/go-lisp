// Example of how plugins can now use the new dependency injection interfaces
package interfaces

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/types"
)

// ExamplePlugin demonstrates the new dependency injection pattern
type ExamplePlugin struct {
	evaluator   CoreEvaluator
	funcCaller  FunctionCaller
	envProvider EnvironmentProvider
}

// NewExamplePlugin creates a plugin that depends only on what it needs
func NewExamplePlugin(
	evaluator CoreEvaluator,
	funcCaller FunctionCaller,
	envProvider EnvironmentProvider,
) *ExamplePlugin {
	return &ExamplePlugin{
		evaluator:   evaluator,
		funcCaller:  funcCaller,
		envProvider: envProvider,
	}
}

// ExampleFunction shows how a plugin function would use the injected dependencies
func (ep *ExamplePlugin) ExampleFunction(evaluator Evaluator, args []types.Expr) (types.Value, error) {
	// This plugin only needs to evaluate expressions and doesn't care about
	// environment management details - much cleaner than before!

	if len(args) != 1 {
		return nil, fmt.Errorf("example-func expects 1 argument, got %d", len(args))
	}

	// Use the injected evaluator
	result, err := ep.evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Could also call functions if needed
	// funcResult, err := ep.funcCaller.CallFunction(someFunc, someArgs)

	// Could access environment if needed
	// env := ep.envProvider.GetEnvironment()
	// value, exists := env.Get("some-var")

	return result, nil
}

// RegisterFunctions shows how the plugin would register its functions
// Note: In practice, this would be implemented by the actual plugin
// and would use the registry package, but we can't import it here due to cycles
func (ep *ExamplePlugin) RegisterFunctions(reg interface{}) error {
	// Function registration remains the same
	// The key difference is that the plugin constructor receives
	// only the interfaces it actually needs
	return nil
}
