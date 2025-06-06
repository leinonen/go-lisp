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

	// The first element should be a function name
	firstExpr := list.Elements[0]
	symbolExpr, ok := firstExpr.(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("first element must be a symbol")
	}

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
	default:
		return nil, fmt.Errorf("unknown function: %s", symbolExpr.Name)
	}
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
