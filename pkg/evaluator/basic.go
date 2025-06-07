// Package evaluator_basic contains basic arithmetic, comparison, and conditional operations
package evaluator

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Basic arithmetic operations

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

// Comparison operations

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

// Conditional operations

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

	// Evaluate the appropriate branch
	// Both branches can be in tail position, so preserve tail call context
	if condBool {
		result, err := e.Eval(args[1])
		if err != nil {
			return nil, err
		}
		// If this was a tail call, propagate it
		return result, nil
	} else {
		result, err := e.Eval(args[2])
		if err != nil {
			return nil, err
		}
		// If this was a tail call, propagate it
		return result, nil
	}
}
