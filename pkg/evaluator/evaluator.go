package evaluator

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Environment represents a variable binding environment
type Environment struct {
	bindings map[string]types.Value
	parent   *Environment
	modules  map[string]*types.ModuleValue // module registry
}

func NewEnvironment() *Environment {
	return &Environment{
		bindings: make(map[string]types.Value),
		parent:   nil,
		modules:  make(map[string]*types.ModuleValue),
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
		modules:  e.modules, // share module registry with parent
	}
}

// Module-related methods
func (e *Environment) GetModule(name string) (*types.ModuleValue, bool) {
	if module, ok := e.modules[name]; ok {
		return module, true
	}
	if e.parent != nil {
		return e.parent.GetModule(name)
	}
	return nil, false
}

func (e *Environment) SetModule(name string, module *types.ModuleValue) {
	if e.modules == nil {
		e.modules = make(map[string]*types.ModuleValue)
	}
	e.modules[name] = module
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
		// Check for qualified module access (module.symbol)
		if strings.Contains(ex.Name, ".") {
			return e.evalQualifiedSymbol(ex.Name)
		}
		value, ok := e.env.Get(ex.Name)
		if !ok {
			return nil, fmt.Errorf("undefined symbol: %s", ex.Name)
		}
		return value, nil
	case *types.ListExpr:
		return e.evalList(ex)
	case *types.ModuleExpr:
		return e.evalModule(ex)
	case *types.ImportExpr:
		return e.evalImport(ex)
	case *types.LoadExpr:
		return e.evalLoad(ex)
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
		case "map":
			return e.evalMap(list.Elements[1:])
		case "filter":
			return e.evalFilter(list.Elements[1:])
		case "reduce":
			return e.evalReduce(list.Elements[1:])
		case "append":
			return e.evalAppend(list.Elements[1:])
		case "reverse":
			return e.evalReverse(list.Elements[1:])
		case "nth":
			return e.evalNth(list.Elements[1:])
		case "env":
			return e.evalEnv(list.Elements[1:])
		case "modules":
			return e.evalModules(list.Elements[1:])
		case "builtins":
			return e.evalBuiltins(list.Elements[1:])
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
	// Check if this is a qualified symbol (module.function)
	if strings.Contains(funcName, ".") {
		funcValue, err := e.evalQualifiedSymbol(funcName)
		if err != nil {
			return nil, err
		}
		return e.callFunction(funcValue, args)
	}

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

// Environment inspection methods

func (e *Evaluator) evalEnv(args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("env requires no arguments")
	}

	// Create a list of (name value) pairs for all bindings in the environment
	var elements []types.Value

	for name, value := range e.env.bindings {
		// Create a pair (name value)
		pair := &types.ListValue{
			Elements: []types.Value{
				types.StringValue(name),
				value,
			},
		}
		elements = append(elements, pair)
	}

	return &types.ListValue{Elements: elements}, nil
}

func (e *Evaluator) evalModules(args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("modules requires no arguments")
	}

	// Create a list of (module-name exports-list) pairs for all modules
	var elements []types.Value

	for name, module := range e.env.modules {
		// Create a list of exported symbol names
		var exports []types.Value
		for exportName := range module.Exports {
			exports = append(exports, types.StringValue(exportName))
		}

		// Create a pair (module-name exports-list)
		pair := &types.ListValue{
			Elements: []types.Value{
				types.StringValue(name),
				&types.ListValue{Elements: exports},
			},
		}
		elements = append(elements, pair)
	}

	return &types.ListValue{Elements: elements}, nil
}

// Built-in functions

func (e *Evaluator) evalBuiltins(args []types.Expr) (types.Value, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("builtins requires 0 or 1 arguments")
	}

	// If no arguments, list all built-in functions
	if len(args) == 0 {
		// List all built-in functions and special forms
		builtinNames := []string{
			// Arithmetic operations
			"+", "-", "*", "/",
			// Comparison operations
			"=", "<", ">",
			// Control flow
			"if",
			// Variable and function definition
			"define", "lambda", "defun",
			// List operations
			"list", "first", "rest", "cons", "length", "empty?",
			// Higher-order functions
			"map", "filter", "reduce",
			// List manipulation
			"append", "reverse", "nth",
			// Environment inspection
			"env", "modules", "builtins",
		}

		// Convert to list of string values
		var elements []types.Value
		for _, name := range builtinNames {
			elements = append(elements, types.StringValue(name))
		}

		return &types.ListValue{Elements: elements}, nil
	}

	// If one argument, show help for that function
	funcNameExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("builtins argument must be a symbol")
	}

	funcName := funcNameExpr.Name
	helpText := e.getBuiltinHelp(funcName)
	if helpText == "" {
		return nil, fmt.Errorf("no help available for '%s' (not a built-in function)", funcName)
	}

	return types.StringValue(helpText), nil
}

func (e *Evaluator) getBuiltinHelp(funcName string) string {
	helpMap := map[string]string{
		// Arithmetic operations
		"+": "(+ num1 num2 ...)\nAddition with multiple operands.\nExample: (+ 1 2 3) => 6",
		"-": "(- num1 num2)\nSubtraction of two numbers.\nExample: (- 10 3) => 7",
		"*": "(* num1 num2 ...)\nMultiplication with multiple operands.\nExample: (* 2 3 4) => 24",
		"/": "(/ num1 num2)\nDivision of two numbers.\nExample: (/ 15 3) => 5",

		// Comparison operations
		"=": "(= val1 val2)\nEquality comparison.\nExample: (= 5 5) => #t",
		"<": "(< num1 num2)\nLess than comparison.\nExample: (< 3 5) => #t",
		">": "(> num1 num2)\nGreater than comparison.\nExample: (> 7 3) => #t",

		// Control flow
		"if": "(if condition then-expr else-expr)\nConditional expression.\nExample: (if (< 3 5) \"yes\" \"no\") => \"yes\"",

		// Variable and function definition
		"define": "(define name value)\nDefine a variable with a name and value.\nExample: (define x 42)",
		"lambda": "(lambda (params) body)\nCreate an anonymous function.\nExample: (lambda (x) (+ x 1))",
		"defun":  "(defun name (params) body)\nDefine a named function.\nExample: (defun square (x) (* x x))",

		// List operations
		"list":   "(list elem1 elem2 ...)\nCreate a list with the given elements.\nExample: (list 1 2 3) => (1 2 3)",
		"first":  "(first lst)\nGet the first element of a list.\nExample: (first (list 1 2 3)) => 1",
		"rest":   "(rest lst)\nGet all elements except the first.\nExample: (rest (list 1 2 3)) => (2 3)",
		"cons":   "(cons elem lst)\nPrepend an element to a list.\nExample: (cons 0 (list 1 2)) => (0 1 2)",
		"length": "(length lst)\nGet the number of elements in a list.\nExample: (length (list 1 2 3)) => 3",
		"empty?": "(empty? lst)\nCheck if a list is empty.\nExample: (empty? (list)) => #t",

		// Higher-order functions
		"map":    "(map func lst)\nApply a function to each element of a list.\nExample: (map (lambda (x) (* x x)) (list 1 2 3)) => (1 4 9)",
		"filter": "(filter predicate lst)\nKeep only elements that satisfy a predicate.\nExample: (filter (lambda (x) (> x 0)) (list -1 2 -3 4)) => (2 4)",
		"reduce": "(reduce func init lst)\nReduce a list to a single value using a function.\nExample: (reduce (lambda (acc x) (+ acc x)) 0 (list 1 2 3)) => 6",

		// List manipulation
		"append":  "(append lst1 lst2)\nCombine two lists into one.\nExample: (append (list 1 2) (list 3 4)) => (1 2 3 4)",
		"reverse": "(reverse lst)\nReverse the order of elements in a list.\nExample: (reverse (list 1 2 3)) => (3 2 1)",
		"nth":     "(nth index lst)\nGet the element at a specific index (0-based).\nExample: (nth 1 (list \"a\" \"b\" \"c\")) => \"b\"",

		// Environment inspection
		"env":      "(env)\nShow all variables and functions in the current environment.\nExample: (env) => ((x 42) (square #<function([x])>))",
		"modules":  "(modules)\nShow all loaded modules and their exported symbols.\nExample: (modules) => ((math (square cube)))",
		"builtins": "(builtins) or (builtins func-name)\nShow all built-in functions or help for a specific function.\nExample: (builtins) => (+ - * / ...) or (builtins reduce) => help for reduce",
	}

	return helpMap[funcName]
}

// Module system methods

func (e *Evaluator) evalQualifiedSymbol(name string) (types.Value, error) {
	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid qualified symbol: %s", name)
	}

	moduleName := parts[0]
	symbolName := parts[1]

	module, ok := e.env.GetModule(moduleName)
	if !ok {
		return nil, fmt.Errorf("module not found: %s", moduleName)
	}

	value, ok := module.Exports[symbolName]
	if !ok {
		return nil, fmt.Errorf("symbol %s not exported by module %s", symbolName, moduleName)
	}

	return value, nil
}

func (e *Evaluator) evalModule(moduleExpr *types.ModuleExpr) (types.Value, error) {
	// Create a new environment for the module
	moduleEnv := e.env.NewChildEnvironment()
	moduleEvaluator := NewEvaluator(moduleEnv.(*Environment))

	// Evaluate the module body
	for _, expr := range moduleExpr.Body {
		_, err := moduleEvaluator.Eval(expr)
		if err != nil {
			return nil, fmt.Errorf("error in module %s: %v", moduleExpr.Name, err)
		}
	}

	// Create the module value with exported bindings
	module := &types.ModuleValue{
		Name:    moduleExpr.Name,
		Exports: make(map[string]types.Value),
		Env:     moduleEnv,
	}

	// Collect exported symbols
	for _, exportName := range moduleExpr.Exports {
		value, ok := moduleEnv.Get(exportName)
		if !ok {
			return nil, fmt.Errorf("exported symbol %s not found in module %s", exportName, moduleExpr.Name)
		}
		module.Exports[exportName] = value
	}

	// Register the module in the global environment
	e.env.SetModule(moduleExpr.Name, module)

	return module, nil
}

func (e *Evaluator) evalImport(importExpr *types.ImportExpr) (types.Value, error) {
	module, ok := e.env.GetModule(importExpr.ModuleName)
	if !ok {
		return nil, fmt.Errorf("module not found: %s", importExpr.ModuleName)
	}

	// Import all exported symbols into current environment
	for name, value := range module.Exports {
		e.env.Set(name, value)
	}

	return module, nil
}

func (e *Evaluator) evalLoad(loadExpr *types.LoadExpr) (types.Value, error) {
	// Read the file
	content, err := ioutil.ReadFile(loadExpr.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", loadExpr.Filename, err)
	}

	// Tokenize
	tokenizer := tokenizer.NewTokenizer(string(content))
	tokens, err := tokenizer.TokenizeWithError()
	if err != nil {
		return nil, fmt.Errorf("tokenization error in %s: %v", loadExpr.Filename, err)
	}

	// Parse and evaluate each expression in the file
	var lastValue types.Value = types.BooleanValue(true) // default return value

	// Parse expressions one by one until we reach the end
	i := 0
	for i < len(tokens) {
		// Find the end of the current expression
		if tokens[i].Type == types.TokenType(-1) { // EOF
			break
		}

		// Extract tokens for this expression
		exprTokens, newIndex := e.extractExpression(tokens, i)
		if len(exprTokens) == 0 {
			break
		}

		// Parse the expression
		parser := parser.NewParser(exprTokens)
		expr, err := parser.Parse()
		if err != nil {
			return nil, fmt.Errorf("parse error in %s: %v", loadExpr.Filename, err)
		}

		// Evaluate the expression
		value, err := e.Eval(expr)
		if err != nil {
			return nil, fmt.Errorf("evaluation error in %s: %v", loadExpr.Filename, err)
		}

		lastValue = value
		i = newIndex
	}

	return lastValue, nil
}

// Helper function to extract a complete expression from tokens
func (e *Evaluator) extractExpression(tokens []types.Token, start int) ([]types.Token, int) {
	if start >= len(tokens) {
		return nil, start
	}

	// Handle single token expressions (numbers, strings, booleans, symbols)
	if tokens[start].Type != types.LPAREN {
		return tokens[start : start+1], start + 1
	}

	// Handle list expressions - find matching closing paren
	parenCount := 0
	end := start
	for end < len(tokens) {
		if tokens[end].Type == types.LPAREN {
			parenCount++
		} else if tokens[end].Type == types.RPAREN {
			parenCount--
			if parenCount == 0 {
				end++
				break
			}
		}
		end++
	}

	if parenCount != 0 {
		// Unmatched parentheses - return what we have
		return tokens[start:end], end
	}

	return tokens[start:end], end
}
