// Package evaluator_atoms contains atom functionality
package evaluator

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// evalAtom creates a new atom with the given initial value
func (e *Evaluator) evalAtom(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("atom requires exactly 1 argument")
	}

	// Evaluate the initial value
	initialValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Create a new atom with the initial value
	return types.NewAtom(initialValue), nil
}

// evalDeref returns the current value of an atom
func (e *Evaluator) evalDeref(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("deref requires exactly 1 argument")
	}

	// Evaluate the atom argument
	atomValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Check if it's an atom
	atom, ok := atomValue.(*types.AtomValue)
	if !ok {
		return nil, fmt.Errorf("deref requires an atom, got %T", atomValue)
	}

	// Return the current value of the atom
	return atom.Value(), nil
}

// evalSwap applies a function to the current value of an atom and updates it
func (e *Evaluator) evalSwap(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("swap! requires exactly 2 arguments")
	}

	// Evaluate the atom argument
	atomValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Check if it's an atom
	atom, ok := atomValue.(*types.AtomValue)
	if !ok {
		return nil, fmt.Errorf("swap! requires an atom as first argument, got %T", atomValue)
	}

	// Evaluate the function argument
	funcValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Check if it's a function
	function, ok := funcValue.(types.FunctionValue)
	if !ok {
		return nil, fmt.Errorf("swap! requires a function as second argument, got %T", funcValue)
	}

	// The function must take exactly 1 parameter
	if len(function.Params) != 1 {
		return nil, fmt.Errorf("swap! function must take exactly 1 parameter, got %d", len(function.Params))
	}

	// Apply the function to the current atom value and update
	return atom.SwapValue(func(currentValue types.Value) types.Value {
		// Create a new environment for the function call
		var funcEnv types.Environment
		if function.Env != nil {
			funcEnv = function.Env.NewChildEnvironment()
		} else {
			funcEnv = e.env.NewChildEnvironment()
		}
		funcEnv.Set(function.Params[0], currentValue)

		// Create evaluator with function environment
		concreteEnv, ok := funcEnv.(*Environment)
		if !ok {
			// This should not happen in normal operation
			panic("internal error: environment type conversion failed")
		}
		funcEvaluator := NewEvaluator(concreteEnv)

		// Evaluate function body
		result, err := funcEvaluator.Eval(function.Body)
		if err != nil {
			// In a real implementation, we'd need better error handling here
			// For now, we'll panic on evaluation errors during swap
			panic(fmt.Sprintf("swap! function evaluation failed: %v", err))
		}

		return result
	}), nil
}

// evalReset sets an atom to a new value directly
func (e *Evaluator) evalReset(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("reset! requires exactly 2 arguments")
	}

	// Evaluate the atom argument
	atomValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Check if it's an atom
	atom, ok := atomValue.(*types.AtomValue)
	if !ok {
		return nil, fmt.Errorf("reset! requires an atom as first argument, got %T", atomValue)
	}

	// Evaluate the new value
	newValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Set the atom to the new value
	return atom.SetValue(newValue), nil
}
