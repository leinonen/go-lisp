// Package functional provides higher-order functions for the Lisp interpreter
package functional

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// FunctionalPlugin provides higher-order functions
type FunctionalPlugin struct {
	*plugins.BasePlugin
}

// NewFunctionalPlugin creates a new functional plugin
func NewFunctionalPlugin() *FunctionalPlugin {
	return &FunctionalPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"functional",
			"1.0.0",
			"Higher-order functions (map, filter, reduce, etc.)",
			[]string{}, // No dependencies
		),
	}
}

// RegisterFunctions registers higher-order functions
func (p *FunctionalPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// map function
	mapFunc := functions.NewFunction(
		"map",
		registry.CategoryFunctional,
		2,
		"Apply a function to each element of a list: (map fn list)",
		p.mapFunc,
	)
	if err := reg.Register(mapFunc); err != nil {
		return err
	}

	// filter function
	filterFunc := functions.NewFunction(
		"filter",
		registry.CategoryFunctional,
		2,
		"Filter elements of a list using a predicate: (filter pred list)",
		p.filterFunc,
	)
	if err := reg.Register(filterFunc); err != nil {
		return err
	}

	// reduce function
	reduceFunc := functions.NewFunction(
		"reduce",
		registry.CategoryFunctional,
		3,
		"Reduce a list to a single value: (reduce fn init list)",
		p.reduceFunc,
	)
	if err := reg.Register(reduceFunc); err != nil {
		return err
	}

	// apply function
	applyFunc := functions.NewFunction(
		"apply",
		registry.CategoryFunctional,
		2,
		"Apply a function to a list of arguments: (apply fn args)",
		p.applyFunc,
	)
	if err := reg.Register(applyFunc); err != nil {
		return err
	}

	return nil
}

// mapFunc implements the map higher-order function
func (p *FunctionalPlugin) mapFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the function
	fnValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the list
	listValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Check if it's a list or vector
	var elements []types.Value
	var isVector bool

	switch v := listValue.(type) {
	case *types.ListValue:
		elements = v.Elements
		isVector = false
	case *types.VectorValue:
		elements = v.Elements
		isVector = true
	default:
		return nil, fmt.Errorf("map requires a list or vector as second argument, got %T", listValue)
	}

	// Apply function to each element
	var results []types.Value
	for _, elem := range elements {
		// Create argument list for function call
		argExpr := p.valueToExpr(elem)
		result, err := p.callFunction(evaluator, fnValue, []types.Expr{argExpr})
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	// Return same type as input
	if isVector {
		return types.NewVectorValue(results), nil
	}
	return &types.ListValue{Elements: results}, nil
}

// filterFunc implements the filter higher-order function
func (p *FunctionalPlugin) filterFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("filter requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the predicate function
	predValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the list
	listValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Check if it's a list or vector
	var elements []types.Value
	var isVector bool

	switch v := listValue.(type) {
	case *types.ListValue:
		elements = v.Elements
		isVector = false
	case *types.VectorValue:
		elements = v.Elements
		isVector = true
	default:
		return nil, fmt.Errorf("filter requires a list or vector as second argument, got %T", listValue)
	}

	// Filter elements based on predicate
	var results []types.Value
	for _, elem := range elements {
		// Create argument list for predicate call
		argExpr := p.valueToExpr(elem)
		result, err := p.callFunction(evaluator, predValue, []types.Expr{argExpr})
		if err != nil {
			return nil, err
		}

		// Check if result is truthy
		if p.isTruthy(result) {
			results = append(results, elem)
		}
	}

	// Return same type as input
	if isVector {
		return types.NewVectorValue(results), nil
	}
	return &types.ListValue{Elements: results}, nil
}

// reduceFunc implements the reduce higher-order function
func (p *FunctionalPlugin) reduceFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("reduce requires exactly 3 arguments, got %d", len(args))
	}

	// Evaluate the function
	fnValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the initial value
	accumulator, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Evaluate the list
	listValue, err := evaluator.Eval(args[2])
	if err != nil {
		return nil, err
	}

	// Check if it's a list or vector
	var elements []types.Value

	switch v := listValue.(type) {
	case *types.ListValue:
		elements = v.Elements
	case *types.VectorValue:
		elements = v.Elements
	default:
		return nil, fmt.Errorf("reduce requires a list or vector as third argument, got %T", listValue)
	}

	// Reduce the collection
	for _, elem := range elements {
		// Create argument list for function call (accumulator, current element)
		accExpr := p.valueToExpr(accumulator)
		elemExpr := p.valueToExpr(elem)
		result, err := p.callFunction(evaluator, fnValue, []types.Expr{accExpr, elemExpr})
		if err != nil {
			return nil, err
		}
		accumulator = result
	}

	return accumulator, nil
}

// applyFunc implements the apply higher-order function
func (p *FunctionalPlugin) applyFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("apply requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the function
	fnValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the arguments list
	argsValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Check if it's a list
	argsList, ok := argsValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("apply requires a list as second argument, got %T", argsValue)
	}

	// Convert values to expressions
	var argExprs []types.Expr
	for _, val := range argsList.Elements {
		argExprs = append(argExprs, p.valueToExpr(val))
	}

	// Call the function with the arguments
	return p.callFunction(evaluator, fnValue, argExprs)
}

// callFunction calls a function with given arguments
func (p *FunctionalPlugin) callFunction(evaluator registry.Evaluator, fnValue types.Value, args []types.Expr) (types.Value, error) {
	// Use the evaluator's CallFunction method
	return evaluator.CallFunction(fnValue, args)
}

// valueToExpr converts a value back to an expression for function calls
func (p *FunctionalPlugin) valueToExpr(value types.Value) types.Expr {
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
		var elements []types.Expr
		for _, elem := range v.Elements {
			elements = append(elements, p.valueToExpr(elem))
		}
		return &types.ListExpr{Elements: elements}
	case *types.VectorValue:
		var elements []types.Expr
		for _, elem := range v.Elements {
			elements = append(elements, p.valueToExpr(elem))
		}
		return &types.BracketExpr{Elements: elements}
	default:
		// For complex types, create a symbol that would represent the value
		// This is a fallback and might not work perfectly in all cases
		return &types.SymbolExpr{Name: fmt.Sprintf("#<%T>", value)}
	}
}

// isTruthy determines if a value is truthy
func (p *FunctionalPlugin) isTruthy(value types.Value) bool {
	switch v := value.(type) {
	case types.BooleanValue:
		return bool(v)
	case *types.NilValue:
		return false
	case types.NumberValue:
		return float64(v) != 0
	case types.StringValue:
		return string(v) != ""
	case *types.ListValue:
		return len(v.Elements) > 0
	case *types.VectorValue:
		return len(v.Elements) > 0
	default:
		return true // Non-nil values are generally truthy
	}
}
