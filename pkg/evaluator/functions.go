// Package evaluator_functions contains function creation, calling, and closure functionality
package evaluator

import (
	"fmt"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Function definition

// Helper function to extract parameter names from BracketExpr only
func extractParameters(paramsExpr types.Expr) ([]string, error) {
	bracketExpr, ok := paramsExpr.(*types.BracketExpr)
	if !ok {
		return nil, fmt.Errorf("parameter list must use square brackets [], got %T", paramsExpr)
	}

	params := make([]string, len(bracketExpr.Elements))
	for i, paramExpr := range bracketExpr.Elements {
		symbolExpr, ok := paramExpr.(*types.SymbolExpr)
		if !ok {
			return nil, fmt.Errorf("parameter must be a symbol, got %T", paramExpr)
		}
		params[i] = symbolExpr.Name
	}

	return params, nil
}

func (e *Evaluator) evalLambda(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("lambda requires exactly 2 arguments: parameters and body")
	}

	// Extract parameter names from square brackets
	params, err := extractParameters(args[0])
	if err != nil {
		return nil, err
	}

	// Create the function value with captured environment
	return types.FunctionValue{
		Params: params,
		Body:   args[1],
		Env:    e.env, // capture current environment for closures
	}, nil
}

func (e *Evaluator) evalDefun(args []types.Expr) (types.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("defn requires at least 3 arguments: name, parameters, and body")
	}

	// First argument must be a symbol (function name)
	nameExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("defn first argument must be a symbol")
	}

	// Second argument must be a list of parameter names in square brackets
	params, err := extractParameters(args[1])
	if err != nil {
		return nil, err
	}

	// If there's only one body expression, use it directly
	// If there are multiple, wrap them in a 'do' form (we'll need to implement this)
	var body types.Expr
	if len(args) == 3 {
		body = args[2]
	} else {
		// For now, we'll just use the last expression as the body
		// In a more complete implementation, we'd want to evaluate all expressions
		// and return the last one (similar to a 'do' or 'progn' form)
		body = args[len(args)-1]
	}

	// Create the function value with captured environment
	function := types.FunctionValue{
		Params: params,
		Body:   body,
		Env:    e.env, // capture current environment for closures
	}

	// Set the function in the environment
	e.env.Set(nameExpr.Name, function)

	// Return the function that was defined
	return function, nil
}

// Function calling

func (e *Evaluator) evalFunctionCall(funcName string, args []types.Expr) (types.Value, error) {
	// Check if this is a qualified symbol (module.function)
	if strings.Contains(funcName, ".") {
		funcValue, err := e.evalQualifiedSymbol(funcName)
		if err != nil {
			return nil, err
		}
		return e.callFunctionWithTailCheck(funcValue, args)
	}

	// Look up the function in the environment
	funcValue, ok := e.env.Get(funcName)
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", funcName)
	}

	return e.callFunctionWithTailCheck(funcValue, args)
}

// callFunctionWithTailCheck checks if we can optimize this call as a tail call
func (e *Evaluator) callFunctionWithTailCheck(funcValue types.Value, args []types.Expr) (types.Value, error) {
	// If tail calls are enabled and this is a tail call, set up the tail call info
	if e.tailCallOK {
		function, ok := funcValue.(types.FunctionValue)
		if ok {
			// Evaluate arguments
			argValues := make([]types.Value, len(args))
			for i, arg := range args {
				argValue, err := e.Eval(arg)
				if err != nil {
					return nil, err
				}
				argValues[i] = argValue
			}

			// Set up tail call instead of making the call
			e.tailCall = &TailCallInfo{
				Function: function,
				Args:     argValues,
			}
			// Return a placeholder - this won't be used since tail call will be detected
			return nil, nil
		}

		// Keywords cannot be tail-call optimized, so handle them normally
		if _, ok := funcValue.(types.KeywordValue); ok {
			return e.callFunction(funcValue, args)
		}

		// Arithmetic operations cannot be tail-call optimized, so handle them normally
		if _, ok := funcValue.(*types.ArithmeticFunctionValue); ok {
			return e.callFunction(funcValue, args)
		}
	}

	// Regular function call
	return e.callFunction(funcValue, args)
}

func (e *Evaluator) callFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	// Check if the value is a keyword being used as a function
	if keyword, ok := funcValue.(types.KeywordValue); ok {
		return e.evalKeywordAsFunction(keyword, args)
	}

	// Check if the value is an arithmetic operation
	if arithFunc, ok := funcValue.(*types.ArithmeticFunctionValue); ok {
		return e.callArithmeticFunction(arithFunc, args)
	}

	// Check if the value is a built-in function
	if builtinFunc, ok := funcValue.(*types.BuiltinFunctionValue); ok {
		return e.callBuiltinFunction(builtinFunc, args)
	}

	function, ok := funcValue.(types.FunctionValue)
	if !ok {
		return nil, fmt.Errorf("value is not a function: %T", funcValue)
	}

	// Check argument count
	if len(args) != len(function.Params) {
		return nil, fmt.Errorf("function expects %d arguments, got %d", len(function.Params), len(args))
	}

	// Evaluate arguments first
	argValues := make([]types.Value, len(args))
	for i, arg := range args {
		argValue, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}
		argValues[i] = argValue
	}

	// Check if we should use tail call optimization
	if e.tailCallOK {
		// Tail call optimization: iterative execution
		currentFunc := function
		currentArgs := argValues

		for {
			// Create a new environment for the function call, extending the captured environment
			var funcEnv types.Environment
			if currentFunc.Env != nil {
				// Use the captured environment as the parent (for closures)
				funcEnv = currentFunc.Env.NewChildEnvironment()
			} else {
				// Fallback to current environment as parent
				funcEnv = e.env.NewChildEnvironment()
			}

			// Bind arguments to parameters
			for i, param := range currentFunc.Params {
				funcEnv.Set(param, currentArgs[i])
			}

			// Create a new evaluator with the function environment
			concreteEnv, ok := funcEnv.(*Environment)
			if !ok {
				return nil, fmt.Errorf("internal error: environment type conversion failed")
			}
			funcEvaluator := NewEvaluator(concreteEnv)
			funcEvaluator.tailCallOK = true // Enable tail call detection

			// Evaluate the function body
			result, err := funcEvaluator.Eval(currentFunc.Body)
			if err != nil {
				return nil, err
			}

			// Check if a tail call was detected
			if funcEvaluator.tailCall != nil {
				// Continue with the tail call instead of returning
				currentFunc = funcEvaluator.tailCall.Function
				currentArgs = funcEvaluator.tailCall.Args
				continue
			}

			// No tail call, return the result
			return result, nil
		}
	} else {
		// Regular function call without tail call optimization
		return e.callFunctionRegular(function, argValues)
	}
}

func (e *Evaluator) callFunctionRegular(function types.FunctionValue, argValues []types.Value) (types.Value, error) {
	// Create a new environment for the function call, extending the captured environment
	var funcEnv types.Environment
	if function.Env != nil {
		// Use the captured environment as the parent (for closures)
		funcEnv = function.Env.NewChildEnvironment()
	} else {
		// Fallback to current environment as parent
		funcEnv = e.env.NewChildEnvironment()
	}

	// Bind arguments to parameters
	for i, param := range function.Params {
		funcEnv.Set(param, argValues[i])
	}

	// Create a new evaluator with the function environment
	concreteEnv, ok := funcEnv.(*Environment)
	if !ok {
		return nil, fmt.Errorf("internal error: environment type conversion failed")
	}
	funcEvaluator := NewEvaluator(concreteEnv)
	// Do not enable tail call detection for regular calls
	funcEvaluator.tailCallOK = false

	// Evaluate the function body
	return funcEvaluator.Eval(function.Body)
}

// evalKeywordAsFunction handles keywords used as functions to get values from hash maps
func (e *Evaluator) evalKeywordAsFunction(keyword types.KeywordValue, args []types.Expr) (types.Value, error) {
	// Keywords as functions can take 1 or 2 arguments:
	// 1 argument: (:key hash-map) -> value or nil
	// 2 arguments: (:key hash-map default) -> value or default
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("keyword function requires 1 or 2 arguments (hash-map and optional default), got %d", len(args))
	}

	// Evaluate the hash map argument
	hashMapValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("keyword function first argument must be a hash map, got %T", hashMapValue)
	}

	// Convert keyword to the string key format used in hash maps
	keyStr := ":" + string(keyword)

	// Look up the value
	value, exists := hashMap.Elements[keyStr]
	if !exists {
		// If no default value provided, return nil
		if len(args) == 1 {
			return &types.NilValue{}, nil
		}

		// Evaluate and return the default value
		defaultValue, err := e.Eval(args[1])
		if err != nil {
			return nil, err
		}
		return defaultValue, nil
	}

	return value, nil
}

// callArithmeticFunction handles calling arithmetic operations as functions
func (e *Evaluator) callArithmeticFunction(arithFunc *types.ArithmeticFunctionValue, args []types.Expr) (types.Value, error) {
	switch arithFunc.Operation {
	case "+":
		return e.evalArithmetic(args, func(a, b float64) float64 { return a + b })
	case "-":
		return e.evalSubtraction(args)
	case "*":
		return e.evalMultiplication(args)
	case "/":
		return e.evalDivision(args)
	case "%":
		return e.evalModulo(args)
	default:
		return nil, fmt.Errorf("unknown arithmetic operation: %s", arithFunc.Operation)
	}
}

// callBuiltinFunction handles calling built-in functions as callable values
func (e *Evaluator) callBuiltinFunction(builtinFunc *types.BuiltinFunctionValue, args []types.Expr) (types.Value, error) {
	// Simply call the built-in function by name through the existing evaluation mechanism
	return e.evalFunctionCall(builtinFunc.Name, args)
}
