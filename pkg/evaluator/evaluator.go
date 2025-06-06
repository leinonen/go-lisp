package evaluator

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Environment represents a variable binding environment
type Environment struct {
	bindings map[string]types.Value
	parent   *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		bindings: make(map[string]types.Value),
		parent:   nil,
	}
}

func (e *Environment) Set(name string, value types.Value) {
	e.bindings[name] = value
}

func (e *Environment) Get(name string) (types.Value, bool) {
	if value, ok := e.bindings[name]; ok {
		return value, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

// NewChildEnvironment creates a new environment with this environment as parent
func (e *Environment) NewChildEnvironment() types.Environment {
	return &Environment{
		bindings: make(map[string]types.Value),
		parent:   e,
	}
}

// Evaluator evaluates expressions
type Evaluator struct {
	env *Environment
}

func NewEvaluator(env *Environment) *Evaluator {
	return &Evaluator{env: env}
}

func (e *Evaluator) Eval(expr types.Expr) (types.Value, error) {
	switch ex := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(ex.Value), nil
	case *types.StringExpr:
		return types.StringValue(ex.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(ex.Value), nil
	case *types.SymbolExpr:
		value, ok := e.env.Get(ex.Name)
		if !ok {
			return nil, fmt.Errorf("undefined symbol: %s", ex.Name)
		}
		return value, nil
	case *types.ListExpr:
		return e.evalList(ex)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (e *Evaluator) evalList(list *types.ListExpr) (types.Value, error) {
	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("empty list cannot be evaluated")
	}

	// The first element could be a special form, built-in function, or user-defined function
	firstExpr := list.Elements[0]

	// Check if it's a symbol (special form or function name)
	if symbolExpr, ok := firstExpr.(*types.SymbolExpr); ok {
		switch symbolExpr.Name {
		case "+":
			return e.evalArithmetic(list.Elements[1:], func(a, b float64) float64 { return a + b })
		case "-":
			return e.evalArithmetic(list.Elements[1:], func(a, b float64) float64 { return a - b })
		case "*":
			return e.evalArithmetic(list.Elements[1:], func(a, b float64) float64 { return a * b })
		case "/":
			return e.evalDivision(list.Elements[1:])
		case "=":
			return e.evalEquality(list.Elements[1:])
		case "<":
			return e.evalComparison(list.Elements[1:], func(a, b float64) bool { return a < b })
		case ">":
			return e.evalComparison(list.Elements[1:], func(a, b float64) bool { return a > b })
		case "if":
			return e.evalIf(list.Elements[1:])
		case "define":
			return e.evalDefine(list.Elements[1:])
		case "lambda":
			return e.evalLambda(list.Elements[1:])
		case "defun":
			return e.evalDefun(list.Elements[1:])
		case "list":
			return e.evalListConstruction(list.Elements[1:])
		case "first":
			return e.evalFirst(list.Elements[1:])
		case "rest":
			return e.evalRest(list.Elements[1:])
		case "cons":
			return e.evalCons(list.Elements[1:])
		case "length":
			return e.evalLength(list.Elements[1:])
		case "empty?":
			return e.evalEmpty(list.Elements[1:])
		default:
			// Try to call it as a user-defined function
			return e.evalFunctionCall(symbolExpr.Name, list.Elements[1:])
		}
	}

	// If first element is not a symbol, evaluate it (could be a lambda expression)
	funcValue, err := e.Eval(firstExpr)
	if err != nil {
		return nil, err
	}

	// Call the function
	return e.callFunction(funcValue, list.Elements[1:])
}

func (e *Evaluator) evalArithmetic(args []types.Expr, op func(float64, float64) float64) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("arithmetic operation requires at least one argument")
	}

	first, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	firstNum, ok := first.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("arithmetic operation requires numbers")
	}

	result := float64(firstNum)

	for i := 1; i < len(args); i++ {
		val, err := e.Eval(args[i])
		if err != nil {
			return nil, err
		}

		num, ok := val.(types.NumberValue)
		if !ok {
			return nil, fmt.Errorf("arithmetic operation requires numbers")
		}

		result = op(result, float64(num))
	}

	return types.NumberValue(result), nil
}

func (e *Evaluator) evalDivision(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("division requires exactly 2 arguments")
	}

	first, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	second, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	firstNum, ok := first.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("division requires numbers")
	}

	secondNum, ok := second.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("division requires numbers")
	}

	if secondNum == 0 {
		return nil, fmt.Errorf("division by zero")
	}

	return types.NumberValue(float64(firstNum) / float64(secondNum)), nil
}

func (e *Evaluator) evalEquality(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("equality requires exactly 2 arguments")
	}

	first, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	second, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// For simplicity, only compare numbers for now
	firstNum, ok1 := first.(types.NumberValue)
	secondNum, ok2 := second.(types.NumberValue)

	if ok1 && ok2 {
		return types.BooleanValue(firstNum == secondNum), nil
	}

	return types.BooleanValue(false), nil
}

func (e *Evaluator) evalComparison(args []types.Expr, op func(float64, float64) bool) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("comparison requires exactly 2 arguments")
	}

	first, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	second, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	firstNum, ok := first.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("comparison requires numbers")
	}

	secondNum, ok := second.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("comparison requires numbers")
	}

	return types.BooleanValue(op(float64(firstNum), float64(secondNum))), nil
}

func (e *Evaluator) evalIf(args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("if requires exactly 3 arguments: condition, then, else")
	}

	condition, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	condBool, ok := condition.(types.BooleanValue)
	if !ok {
		return nil, fmt.Errorf("if condition must be a boolean")
	}

	if condBool {
		return e.Eval(args[1])
	} else {
		return e.Eval(args[2])
	}
}

func (e *Evaluator) evalDefine(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("define requires exactly 2 arguments: name and value")
	}

	// First argument must be a symbol (variable name)
	nameExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("define first argument must be a symbol")
	}

	// Evaluate the second argument (the value)
	value, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Set the variable in the environment
	e.env.Set(nameExpr.Name, value)

	// Return the value that was defined
	return value, nil
}

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

func (e *Evaluator) evalFunctionCall(funcName string, args []types.Expr) (types.Value, error) {
	// Look up the function in the environment
	funcValue, ok := e.env.Get(funcName)
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", funcName)
	}

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

	// Create a new environment for the function call, extending the captured environment
	var funcEnv types.Environment
	if function.Env != nil {
		// Use the captured environment as the parent (for closures)
		funcEnv = function.Env.NewChildEnvironment()
	} else {
		// Fallback to current environment as parent
		funcEnv = e.env.NewChildEnvironment()
	}

	// Evaluate arguments and bind them to parameters
	for i, arg := range args {
		argValue, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}
		funcEnv.Set(function.Params[i], argValue)
	}

	// Create a new evaluator with the function environment
	// We need to convert back to concrete type for the evaluator
	concreteEnv, ok := funcEnv.(*Environment)
	if !ok {
		return nil, fmt.Errorf("internal error: environment type conversion failed")
	}
	funcEvaluator := NewEvaluator(concreteEnv)

	// Evaluate the function body
	return funcEvaluator.Eval(function.Body)
}

// List operation methods

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
