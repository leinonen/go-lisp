// Package evaluator_functions contains function creation, calling, and closure functionality
package evaluator

import (
	"fmt"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Function definition

func (e *Evaluator) evalLambda(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("lambda requires exactly 2 arguments: parameters and body")
	}

	// First argument must be a list of parameter names
	paramsExpr, ok := args[0].(*types.ListExpr)
	if !ok {
		return nil, fmt.Errorf("lambda first argument must be a parameter list")
	}

	// Extract parameter names
	params := make([]string, len(paramsExpr.Elements))
	for i, paramExpr := range paramsExpr.Elements {
		symbolExpr, ok := paramExpr.(*types.SymbolExpr)
		if !ok {
			return nil, fmt.Errorf("lambda parameter must be a symbol, got %T", paramExpr)
		}
		params[i] = symbolExpr.Name
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
		return nil, fmt.Errorf("defun requires at least 3 arguments: name, parameters, and body")
	}

	// First argument must be a symbol (function name)
	nameExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("defun first argument must be a symbol")
	}

	// Second argument must be a list of parameter names
	paramsExpr, ok := args[1].(*types.ListExpr)
	if !ok {
		return nil, fmt.Errorf("defun second argument must be a parameter list")
	}

	// Extract parameter names
	params := make([]string, len(paramsExpr.Elements))
	for i, paramExpr := range paramsExpr.Elements {
		symbolExpr, ok := paramExpr.(*types.SymbolExpr)
		if !ok {
			return nil, fmt.Errorf("defun parameter must be a symbol, got %T", paramExpr)
		}
		params[i] = symbolExpr.Name
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
	}

	// Regular function call
	return e.callFunction(funcValue, args)
}

func (e *Evaluator) callFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
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
