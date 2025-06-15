// Package arithmetic provides basic arithmetic operations as a plugin
package arithmetic

import (
	"fmt"
	"math"

	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// ArithmeticPlugin provides basic arithmetic operations
type ArithmeticPlugin struct {
	*plugins.BasePlugin
}

// NewArithmeticPlugin creates a new arithmetic plugin
func NewArithmeticPlugin() *ArithmeticPlugin {
	return &ArithmeticPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"arithmetic",
			"1.0.0",
			"Basic arithmetic operations (+, -, *, /, %)",
			[]string{}, // No dependencies
		),
	}
}

// RegisterFunctions registers arithmetic functions
func (ap *ArithmeticPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// Addition
	addFunc := functions.NewFunction(
		"+",
		registry.CategoryArithmetic,
		-1, // Variadic
		"Add numbers: (+ 1 2 3) => 6",
		ap.evalAdd,
	)
	if err := reg.Register(addFunc); err != nil {
		return err
	}

	// Subtraction
	subFunc := functions.NewFunction(
		"-",
		registry.CategoryArithmetic,
		-1, // Variadic
		"Subtract numbers: (- 10 3 2) => 5, (- 5) => -5",
		ap.evalSubtract,
	)
	if err := reg.Register(subFunc); err != nil {
		return err
	}

	// Multiplication
	mulFunc := functions.NewFunction(
		"*",
		registry.CategoryArithmetic,
		-1, // Variadic
		"Multiply numbers: (* 2 3 4) => 24",
		ap.evalMultiply,
	)
	if err := reg.Register(mulFunc); err != nil {
		return err
	}

	// Division
	divFunc := functions.NewFunction(
		"/",
		registry.CategoryArithmetic,
		-1, // Variadic
		"Divide numbers: (/ 12 3 2) => 2",
		ap.evalDivide,
	)
	if err := reg.Register(divFunc); err != nil {
		return err
	}

	// Modulo
	modFunc := functions.NewFunction(
		"%",
		registry.CategoryArithmetic,
		2, // Exactly 2 arguments
		"Modulo operation: (% 10 3) => 1",
		ap.evalModulo,
	)
	return reg.Register(modFunc)
}

// evalAdd implements addition
func (ap *ArithmeticPlugin) evalAdd(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return types.NumberValue(0), nil
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	// Check if any values are big numbers
	hasBig := false
	for _, val := range values {
		if _, ok := val.(*types.BigNumberValue); ok {
			hasBig = true
			break
		}
	}

	if hasBig {
		return ap.addBig(values)
	}
	return ap.addRegular(values)
}

// addRegular handles regular number addition
func (ap *ArithmeticPlugin) addRegular(values []types.Value) (types.Value, error) {
	result := 0.0
	for _, val := range values {
		num, err := functions.ExtractFloat64(val)
		if err != nil {
			return nil, err
		}
		result += num
	}
	return types.NumberValue(result), nil
}

// addBig handles big number addition
func (ap *ArithmeticPlugin) addBig(values []types.Value) (types.Value, error) {
	result := types.NewBigNumberFromInt64(0)
	for _, val := range values {
		switch v := val.(type) {
		case types.NumberValue:
			// Convert regular number to big number and add
			num := types.NewBigNumberFromInt64(int64(v))
			result.Value.Add(result.Value, num.Value)
		case *types.BigNumberValue:
			result.Value.Add(result.Value, v.Value)
		default:
			return nil, fmt.Errorf("expected number, got %T", val)
		}
	}
	return result, nil
}

// evalSubtract implements subtraction
func (ap *ArithmeticPlugin) evalSubtract(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("- requires at least 1 argument")
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	// Unary minus
	if len(values) == 1 {
		switch v := values[0].(type) {
		case types.NumberValue:
			return types.NumberValue(-float64(v)), nil
		case *types.BigNumberValue:
			result := types.NewBigNumberValue(v.Value)
			result.Value.Neg(result.Value)
			return result, nil
		default:
			return nil, fmt.Errorf("expected number, got %T", v)
		}
	}

	// Check if any values are big numbers
	hasBig := false
	for _, val := range values {
		if _, ok := val.(*types.BigNumberValue); ok {
			hasBig = true
			break
		}
	}

	if hasBig {
		return ap.subtractBig(values)
	}
	return ap.subtractRegular(values)
}

// subtractRegular handles regular number subtraction
func (ap *ArithmeticPlugin) subtractRegular(values []types.Value) (types.Value, error) {
	first, err := functions.ExtractFloat64(values[0])
	if err != nil {
		return nil, err
	}

	result := first
	for i := 1; i < len(values); i++ {
		num, err := functions.ExtractFloat64(values[i])
		if err != nil {
			return nil, err
		}
		result -= num
	}
	return types.NumberValue(result), nil
}

// subtractBig handles big number subtraction
func (ap *ArithmeticPlugin) subtractBig(values []types.Value) (types.Value, error) {
	// Convert first value to big number
	var result *types.BigNumberValue
	switch v := values[0].(type) {
	case types.NumberValue:
		result = types.NewBigNumberFromInt64(int64(v))
	case *types.BigNumberValue:
		result = types.NewBigNumberValue(v.Value)
	default:
		return nil, fmt.Errorf("expected number, got %T", v)
	}

	// Subtract remaining values
	for i := 1; i < len(values); i++ {
		switch v := values[i].(type) {
		case types.NumberValue:
			num := types.NewBigNumberFromInt64(int64(v))
			result.Value.Sub(result.Value, num.Value)
		case *types.BigNumberValue:
			result.Value.Sub(result.Value, v.Value)
		default:
			return nil, fmt.Errorf("expected number, got %T", v)
		}
	}
	return result, nil
}

// evalMultiply implements multiplication
func (ap *ArithmeticPlugin) evalMultiply(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return types.NumberValue(1), nil
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	// Check if any values are big numbers
	hasBig := false
	for _, val := range values {
		if _, ok := val.(*types.BigNumberValue); ok {
			hasBig = true
			break
		}
	}

	if hasBig {
		return ap.multiplyBig(values)
	}
	return ap.multiplyRegular(values)
}

// multiplyRegular handles regular number multiplication
func (ap *ArithmeticPlugin) multiplyRegular(values []types.Value) (types.Value, error) {
	result := 1.0
	for _, val := range values {
		num, err := functions.ExtractFloat64(val)
		if err != nil {
			return nil, err
		}
		result *= num
	}
	return types.NumberValue(result), nil
}

// multiplyBig handles big number multiplication
func (ap *ArithmeticPlugin) multiplyBig(values []types.Value) (types.Value, error) {
	result := types.NewBigNumberFromInt64(1)
	for _, val := range values {
		switch v := val.(type) {
		case types.NumberValue:
			num := types.NewBigNumberFromInt64(int64(v))
			result.Value.Mul(result.Value, num.Value)
		case *types.BigNumberValue:
			result.Value.Mul(result.Value, v.Value)
		default:
			return nil, fmt.Errorf("expected number, got %T", val)
		}
	}
	return result, nil
}

// evalDivide implements division
func (ap *ArithmeticPlugin) evalDivide(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("/ requires at least 1 argument")
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	// For now, always use regular division (returns float)
	// TODO: Could implement integer division for big numbers
	first, err := functions.ExtractFloat64(values[0])
	if err != nil {
		return nil, err
	}

	if len(values) == 1 {
		// (/ x) means 1/x
		if first == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return types.NumberValue(1.0 / first), nil
	}

	result := first
	for i := 1; i < len(values); i++ {
		num, err := functions.ExtractFloat64(values[i])
		if err != nil {
			return nil, err
		}
		if num == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result /= num
	}
	return types.NumberValue(result), nil
}

// evalModulo implements modulo operation
func (ap *ArithmeticPlugin) evalModulo(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("modulo requires exactly 2 arguments, got %d", len(args))
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	x, err := functions.ExtractFloat64(values[0])
	if err != nil {
		return nil, err
	}

	y, err := functions.ExtractFloat64(values[1])
	if err != nil {
		return nil, err
	}

	if y == 0 {
		return nil, fmt.Errorf("modulo by zero")
	}

	// Use Go's built-in modulo behavior
	result := math.Mod(x, y)
	return types.NumberValue(result), nil
}
