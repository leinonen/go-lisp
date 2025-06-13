// Package evaluator_print contains print functionality for the Lisp interpreter
package evaluator

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Print! function - outputs values to stdout without newline (side effect)
func (e *Evaluator) evalPrint(args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return &types.NilValue{}, nil
	}

	// Evaluate all arguments
	for i, arg := range args {
		value, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(e.valueToString(value))
	}

	return &types.NilValue{}, nil
}

// Println! function - outputs values to stdout with newline (side effect)
func (e *Evaluator) evalPrintln(args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		fmt.Println()
		return &types.NilValue{}, nil
	}

	// Evaluate all arguments
	for i, arg := range args {
		value, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(e.valueToString(value))
	}
	fmt.Println()

	return &types.NilValue{}, nil
}

// Helper function to convert a value to its string representation for printing
func (e *Evaluator) valueToString(value types.Value) string {
	switch v := value.(type) {
	case *types.NilValue:
		return "nil"
	case types.StringValue:
		return string(v)
	case types.NumberValue:
		return fmt.Sprintf("%.0f", float64(v))
	case types.BooleanValue:
		if v {
			return "#t"
		}
		return "#f"
	case *types.ListValue:
		result := "("
		for i, elem := range v.Elements {
			if i > 0 {
				result += " "
			}
			result += e.valueToString(elem)
		}
		result += ")"
		return result
	case *types.FunctionValue:
		paramNames := make([]string, len(v.Params))
		for i, param := range v.Params {
			paramNames[i] = param
		}
		return fmt.Sprintf("#<function([%s])>", fmt.Sprintf("%v", paramNames))
	case *types.HashMapValue:
		result := "{"
		first := true
		for key, val := range v.Elements {
			if !first {
				result += ", "
			}
			result += fmt.Sprintf("%s: %s", key, e.valueToString(val))
			first = false
		}
		result += "}"
		return result
	case *types.BigNumberValue:
		return v.String()
	default:
		return fmt.Sprintf("%v", value)
	}
}
