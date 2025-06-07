// Package evaluator_lists contains list operations and manipulation functionality
package evaluator

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Basic list operations

func (e *Evaluator) evalListConstruction(args []types.Expr) (types.Value, error) {
	elements := make([]types.Value, len(args))
	for i, arg := range args {
		value, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}
		elements[i] = value
	}
	return &types.ListValue{Elements: elements}, nil
}

func (e *Evaluator) evalFirst(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("first requires exactly 1 argument")
	}

	listValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("first requires a list, got %T", listValue)
	}

	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("first: list is empty")
	}

	return list.Elements[0], nil
}

func (e *Evaluator) evalRest(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("rest requires exactly 1 argument")
	}

	listValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("rest requires a list, got %T", listValue)
	}

	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("rest: list is empty")
	}

	restElements := make([]types.Value, len(list.Elements)-1)
	copy(restElements, list.Elements[1:])
	return &types.ListValue{Elements: restElements}, nil
}

func (e *Evaluator) evalCons(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("cons requires exactly 2 arguments")
	}

	elementValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	listValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("cons second argument must be a list, got %T", listValue)
	}

	newElements := make([]types.Value, len(list.Elements)+1)
	newElements[0] = elementValue
	copy(newElements[1:], list.Elements)
	return &types.ListValue{Elements: newElements}, nil
}

func (e *Evaluator) evalLength(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("length requires exactly 1 argument")
	}

	listValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("length requires a list, got %T", listValue)
	}

	return types.NumberValue(float64(len(list.Elements))), nil
}

func (e *Evaluator) evalEmpty(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("empty? requires exactly 1 argument")
	}

	listValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("empty? requires a list, got %T", listValue)
	}

	return types.BooleanValue(len(list.Elements) == 0), nil
}

// List manipulation functions

func (e *Evaluator) evalAppend(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("append requires exactly 2 arguments")
	}

	// Evaluate the first list
	firstValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the second list
	secondValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	firstList, ok := firstValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("append first argument must be a list, got %T", firstValue)
	}

	secondList, ok := secondValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("append second argument must be a list, got %T", secondValue)
	}

	// Create a new list with combined elements
	resultElements := make([]types.Value, 0, len(firstList.Elements)+len(secondList.Elements))
	resultElements = append(resultElements, firstList.Elements...)
	resultElements = append(resultElements, secondList.Elements...)

	return &types.ListValue{Elements: resultElements}, nil
}

func (e *Evaluator) evalReverse(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("reverse requires exactly 1 argument")
	}

	// Evaluate the list
	listValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("reverse argument must be a list, got %T", listValue)
	}

	// Create a new list with reversed elements
	resultElements := make([]types.Value, len(list.Elements))
	for i, elem := range list.Elements {
		resultElements[len(list.Elements)-1-i] = elem
	}

	return &types.ListValue{Elements: resultElements}, nil
}

func (e *Evaluator) evalNth(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("nth requires exactly 2 arguments")
	}

	// Evaluate the index
	indexValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the list
	listValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	index, ok := indexValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("nth first argument must be a number, got %T", indexValue)
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("nth second argument must be a list, got %T", listValue)
	}

	// Check bounds
	idx := int(index)
	if idx < 0 {
		return nil, fmt.Errorf("nth index cannot be negative: %d", idx)
	}

	if idx >= len(list.Elements) {
		return nil, fmt.Errorf("nth index %d out of bounds for list of length %d", idx, len(list.Elements))
	}

	return list.Elements[idx], nil
}

// Higher-order functions

func (e *Evaluator) evalMap(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map requires exactly 2 arguments")
	}

	// Evaluate the function
	funcValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the list
	listValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("map second argument must be a list, got %T", listValue)
	}

	function, ok := funcValue.(types.FunctionValue)
	if !ok {
		return nil, fmt.Errorf("map first argument must be a function, got %T", funcValue)
	}

	if len(function.Params) != 1 {
		return nil, fmt.Errorf("map function must take exactly 1 parameter, got %d", len(function.Params))
	}

	// Apply function to each element
	resultElements := make([]types.Value, len(list.Elements))
	for i, elem := range list.Elements {
		// Create a new environment for the function call
		var funcEnv types.Environment
		if function.Env != nil {
			funcEnv = function.Env.NewChildEnvironment()
		} else {
			funcEnv = e.env.NewChildEnvironment()
		}
		funcEnv.Set(function.Params[0], elem)

		// Create evaluator with function environment
		concreteEnv, ok := funcEnv.(*Environment)
		if !ok {
			return nil, fmt.Errorf("internal error: environment type conversion failed")
		}
		funcEvaluator := NewEvaluator(concreteEnv)

		// Evaluate function body
		result, err := funcEvaluator.Eval(function.Body)
		if err != nil {
			return nil, err
		}
		resultElements[i] = result
	}

	return &types.ListValue{Elements: resultElements}, nil
}

func (e *Evaluator) evalFilter(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("filter requires exactly 2 arguments")
	}

	// Evaluate the predicate function
	funcValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the list
	listValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("filter second argument must be a list, got %T", listValue)
	}

	function, ok := funcValue.(types.FunctionValue)
	if !ok {
		return nil, fmt.Errorf("filter first argument must be a function, got %T", funcValue)
	}

	if len(function.Params) != 1 {
		return nil, fmt.Errorf("filter function must take exactly 1 parameter, got %d", len(function.Params))
	}

	// Filter elements based on predicate
	var resultElements []types.Value
	for _, elem := range list.Elements {
		// Create a new environment for the function call
		var funcEnv types.Environment
		if function.Env != nil {
			funcEnv = function.Env.NewChildEnvironment()
		} else {
			funcEnv = e.env.NewChildEnvironment()
		}
		funcEnv.Set(function.Params[0], elem)

		// Create evaluator with function environment
		concreteEnv, ok := funcEnv.(*Environment)
		if !ok {
			return nil, fmt.Errorf("internal error: environment type conversion failed")
		}
		funcEvaluator := NewEvaluator(concreteEnv)

		// Evaluate predicate function
		result, err := funcEvaluator.Eval(function.Body)
		if err != nil {
			return nil, err
		}

		// Check if result is truthy
		if isTruthy(result) {
			resultElements = append(resultElements, elem)
		}
	}

	return &types.ListValue{Elements: resultElements}, nil
}

func (e *Evaluator) evalReduce(args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("reduce requires exactly 3 arguments")
	}

	// Evaluate the function
	funcValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the initial value
	accumulator, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Evaluate the list
	listValue, err := e.Eval(args[2])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("reduce third argument must be a list, got %T", listValue)
	}

	function, ok := funcValue.(types.FunctionValue)
	if !ok {
		return nil, fmt.Errorf("reduce first argument must be a function, got %T", funcValue)
	}

	if len(function.Params) != 2 {
		return nil, fmt.Errorf("reduce function must take exactly 2 parameters, got %d", len(function.Params))
	}

	// Reduce over the list
	for _, elem := range list.Elements {
		// Create a new environment for the function call
		var funcEnv types.Environment
		if function.Env != nil {
			funcEnv = function.Env.NewChildEnvironment()
		} else {
			funcEnv = e.env.NewChildEnvironment()
		}
		funcEnv.Set(function.Params[0], accumulator)
		funcEnv.Set(function.Params[1], elem)

		// Create evaluator with function environment
		concreteEnv, ok := funcEnv.(*Environment)
		if !ok {
			return nil, fmt.Errorf("internal error: environment type conversion failed")
		}
		funcEvaluator := NewEvaluator(concreteEnv)

		// Evaluate function body
		result, err := funcEvaluator.Eval(function.Body)
		if err != nil {
			return nil, err
		}
		accumulator = result
	}

	return accumulator, nil
}

// Helper function to check if a value is truthy
func isTruthy(value types.Value) bool {
	switch v := value.(type) {
	case types.BooleanValue:
		return bool(v)
	case types.NumberValue:
		return v != 0
	case types.StringValue:
		return v != ""
	case *types.ListValue:
		return len(v.Elements) > 0
	default:
		return true // Other values are considered truthy
	}
}
