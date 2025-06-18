// Package functions provides utilities for creating built-in functions
package functions

import (
	"fmt"
	"math/big"

	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// BaseFunction provides a basic implementation of BuiltinFunction
type BaseFunction struct {
	name     string
	category string
	arity    int
	help     string
	callFunc func(evaluator registry.Evaluator, args []types.Expr) (types.Value, error)
}

// NewFunction creates a new base function
func NewFunction(name, category string, arity int, help string,
	callFunc func(evaluator registry.Evaluator, args []types.Expr) (types.Value, error)) *BaseFunction {
	return &BaseFunction{
		name:     name,
		category: category,
		arity:    arity,
		help:     help,
		callFunc: callFunc,
	}
}

// Name returns the function name
func (bf *BaseFunction) Name() string {
	return bf.name
}

// Category returns the function category
func (bf *BaseFunction) Category() string {
	return bf.category
}

// Arity returns the function arity
func (bf *BaseFunction) Arity() int {
	return bf.arity
}

// Help returns the function help text
func (bf *BaseFunction) Help() string {
	return bf.help
}

// Call executes the function
func (bf *BaseFunction) Call(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	// Validate arity if not variadic
	if bf.arity >= 0 && len(args) != bf.arity {
		return nil, fmt.Errorf("%s requires exactly %d arguments, got %d", bf.name, bf.arity, len(args))
	}

	return bf.callFunc(evaluator, args)
}

// Helper functions for common operations

// EvalArgs evaluates all arguments
func EvalArgs(evaluator registry.Evaluator, args []types.Expr) ([]types.Value, error) {
	values := make([]types.Value, len(args))
	for i, arg := range args {
		val, err := evaluator.Eval(arg)
		if err != nil {
			return nil, err
		}
		values[i] = val
	}
	return values, nil
}

// ExtractFloat64 extracts a float64 from a value
func ExtractFloat64(value types.Value) (float64, error) {
	switch v := value.(type) {
	case types.NumberValue:
		return float64(v), nil
	case *types.BigNumberValue:
		// Convert big number to float64 (may lose precision)
		f := new(big.Float).SetInt(v.Value)
		result, _ := f.Float64()
		return result, nil
	default:
		return 0, fmt.Errorf("expected number, got %T", value)
	}
}

// ExtractString extracts a string from a value
func ExtractString(value types.Value) (string, error) {
	switch v := value.(type) {
	case types.StringValue:
		return string(v), nil
	default:
		return "", fmt.Errorf("expected string, got %T", value)
	}
}

// ExtractBool extracts a boolean from a value
func ExtractBool(value types.Value) (bool, error) {
	switch v := value.(type) {
	case types.BooleanValue:
		return bool(v), nil
	default:
		return false, fmt.Errorf("expected boolean, got %T", value)
	}
}

// ExtractList extracts a list from a value
func ExtractList(value types.Value) (*types.ListValue, error) {
	switch v := value.(type) {
	case *types.ListValue:
		return v, nil
	default:
		return nil, fmt.Errorf("expected list, got %T", value)
	}
}

// ValidateArity validates argument count for variadic functions
func ValidateArity(funcName string, minArgs int, args []types.Expr) error {
	if len(args) < minArgs {
		return fmt.Errorf("%s requires at least %d arguments, got %d", funcName, minArgs, len(args))
	}
	return nil
}

// ValidateExactArity validates argument count for fixed-arity functions
func ValidateExactArity(funcName string, expectedArgs int, args []types.Expr) error {
	if len(args) != expectedArgs {
		return fmt.Errorf("%s requires exactly %d arguments, got %d", funcName, expectedArgs, len(args))
	}
	return nil
}

// IsTruthy returns whether a value is considered true in Lisp
func IsTruthy(value types.Value) bool {
	switch v := value.(type) {
	case types.BooleanValue:
		return bool(v)
	case types.NumberValue:
		return float64(v) != 0
	case types.StringValue:
		return string(v) != ""
	case *types.ListValue:
		return len(v.Elements) > 0
	default:
		return value != nil
	}
}
