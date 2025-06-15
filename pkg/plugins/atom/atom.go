// Package atom implements atom functions as a plugin
package atom

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// AtomPlugin implements atom functions
type AtomPlugin struct {
	*plugins.BasePlugin
}

// NewAtomPlugin creates a new atom plugin
func NewAtomPlugin() *AtomPlugin {
	return &AtomPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"atom",
			"1.0.0",
			"Atom functions for thread-safe mutable state (atom, deref, swap!, reset!)",
			[]string{}, // No dependencies
		),
	}
}

// Functions returns the list of functions provided by this plugin
func (p *AtomPlugin) Functions() []string {
	return []string{
		"atom", "deref", "swap!", "reset!",
	}
}

// RegisterFunctions registers all atom functions with the registry
func (p *AtomPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// atom function
	atomFunc := functions.NewFunction(
		"atom",
		registry.CategoryAtom,
		1,
		"Create an atom with initial value: (atom 42)",
		p.evalAtom,
	)
	if err := reg.Register(atomFunc); err != nil {
		return err
	}

	// deref function
	derefFunc := functions.NewFunction(
		"deref",
		registry.CategoryAtom,
		1,
		"Get current value of atom: (deref my-atom)",
		p.evalDeref,
	)
	if err := reg.Register(derefFunc); err != nil {
		return err
	}

	// swap! function
	swapFunc := functions.NewFunction(
		"swap!",
		registry.CategoryAtom,
		-1, // Variable arguments (atom, function, optional args)
		"Atomically update atom by applying function: (swap! my-atom inc)",
		p.evalSwap,
	)
	if err := reg.Register(swapFunc); err != nil {
		return err
	}

	// reset! function
	resetFunc := functions.NewFunction(
		"reset!",
		registry.CategoryAtom,
		2,
		"Set atom to new value: (reset! my-atom 100)",
		p.evalReset,
	)
	if err := reg.Register(resetFunc); err != nil {
		return err
	}

	return nil
}

// evalAtom creates a new atom with the given initial value
func (p *AtomPlugin) evalAtom(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("atom requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	return types.NewAtom(value), nil
}

// evalDeref gets the current value of an atom
func (p *AtomPlugin) evalDeref(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("deref requires exactly 1 argument, got %d", len(args))
	}

	atomValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	atom, ok := atomValue.(*types.AtomValue)
	if !ok {
		return nil, fmt.Errorf("deref requires an atom, got %T", atomValue)
	}

	return atom.Value(), nil
}

// evalSwap atomically updates an atom by applying a function to its current value
func (p *AtomPlugin) evalSwap(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("swap! requires at least 2 arguments (atom and function), got %d", len(args))
	}

	// Evaluate the atom
	atomValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	atom, ok := atomValue.(*types.AtomValue)
	if !ok {
		return nil, fmt.Errorf("swap! first argument must be an atom, got %T", atomValue)
	}

	// Evaluate the function
	fnValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Validate that the second argument is actually a function
	if _, ok := fnValue.(*types.FunctionValue); !ok {
		return nil, fmt.Errorf("swap! second argument must be a function, got %T", fnValue)
	}

	// Prepare additional arguments if any
	additionalArgs := args[2:]

	// Apply the function atomically
	newValue := atom.SwapValue(func(currentValue types.Value) types.Value {
		// Create argument list: current value + additional arguments
		allArgs := make([]types.Expr, len(additionalArgs)+1)
		// First argument is the current value (we need to convert it to an expression)
		allArgs[0] = p.valueToExpr(currentValue)

		// Add additional arguments
		copy(allArgs[1:], additionalArgs)

		// Call the function with the current value and additional arguments
		result, err := evaluator.CallFunction(fnValue, allArgs)
		if err != nil {
			// In case of error, return the current value unchanged
			// Note: This is a limitation - we can't propagate errors from within SwapValue
			return currentValue
		}
		return result
	})

	return newValue, nil
}

// Helper function to convert a value back to an expression
// This is needed for the swap function to work properly
func (p *AtomPlugin) valueToExpr(value types.Value) types.Expr {
	switch v := value.(type) {
	case types.NumberValue:
		return &types.NumberExpr{Value: float64(v)}
	case types.StringValue:
		return &types.StringExpr{Value: string(v)}
	case types.BooleanValue:
		return &types.BooleanExpr{Value: bool(v)}
	case types.KeywordValue:
		return &types.KeywordExpr{Value: string(v)}
	case *types.ListValue:
		elements := make([]types.Expr, len(v.Elements))
		for i, elem := range v.Elements {
			elements[i] = p.valueToExpr(elem)
		}
		return &types.ListExpr{Elements: elements}
	default:
		// For complex types that don't have a direct expression representation,
		// we'll need to create a special expression type or handle differently
		// For now, we'll create a symbol expression that would need special handling
		return &types.SymbolExpr{Name: fmt.Sprintf("#<value:%T>", value)}
	}
}

// evalReset sets an atom to a new value
func (p *AtomPlugin) evalReset(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("reset! requires exactly 2 arguments, got %d", len(args))
	}

	atomValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	atom, ok := atomValue.(*types.AtomValue)
	if !ok {
		return nil, fmt.Errorf("reset! first argument must be an atom, got %T", atomValue)
	}

	newValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	return atom.SetValue(newValue), nil
}
