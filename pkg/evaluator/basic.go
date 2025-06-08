// Package evaluator_basic contains basic arithmetic, comparison, and conditional operations
package evaluator

import (
	"fmt"
	"math"
	"math/big"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Basic arithmetic operations

// Helper function to convert any numeric value to big.Int for computation
func toBigInt(val types.Value) (*big.Int, error) {
	switch v := val.(type) {
	case types.NumberValue:
		f := float64(v)
		// Check if it's a whole number
		if f != float64(int64(f)) {
			return nil, fmt.Errorf("cannot convert non-integer to big integer: %v", f)
		}
		return big.NewInt(int64(f)), nil
	case *types.BigNumberValue:
		return new(big.Int).Set(v.Value), nil
	default:
		return nil, fmt.Errorf("not a number: %T", val)
	}
}

// Helper function to determine if we should use big integers
func shouldUseBigInt(args []types.Value) bool {
	for _, val := range args {
		if _, ok := val.(*types.BigNumberValue); ok {
			return true
		}
		if num, ok := val.(types.NumberValue); ok {
			f := float64(num)
			// Use big int for large integers to avoid precision loss
			if f >= 1e15 || f <= -1e15 {
				return true
			}
		}
	}
	return false
}

// Helper function to check if multiplication might overflow
func mightOverflowMultiplication(args []types.Value) bool {
	if len(args) < 2 {
		return false
	}

	// Calculate product estimate to see if it might exceed safe float64 range
	product := 1.0
	for _, val := range args {
		switch v := val.(type) {
		case types.NumberValue:
			f := float64(v)
			if f == 0 {
				return false // multiplication by zero won't overflow
			}
			product *= math.Abs(f)
			// If intermediate product already exceeds 1e15, use big arithmetic
			if product > 1e15 {
				return true
			}
		case *types.BigNumberValue:
			return true // already a big number
		}
	}
	return false
}

// Enhanced arithmetic evaluation that supports both regular and big numbers
func (e *Evaluator) evalArithmetic(args []types.Expr, op func(float64, float64) float64) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("arithmetic operation requires at least one argument")
	}

	// Evaluate all arguments first
	values := make([]types.Value, len(args))
	for i, arg := range args {
		val, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}
		values[i] = val
	}

	// Check if we should use big integers
	if shouldUseBigInt(values) {
		return e.evalBigArithmetic(values, op)
	}

	// Use regular float64 arithmetic
	firstNum, ok := values[0].(types.NumberValue)
	if !ok {
		// Check if it's a big number that we need to convert
		if bigNum, ok := values[0].(*types.BigNumberValue); ok {
			f, _ := bigNum.Value.Float64()
			firstNum = types.NumberValue(f)
		} else {
			return nil, fmt.Errorf("arithmetic operation requires numbers")
		}
	}

	result := float64(firstNum)
	for i := 1; i < len(values); i++ {
		var num types.NumberValue
		var ok bool

		if num, ok = values[i].(types.NumberValue); !ok {
			if bigNum, ok := values[i].(*types.BigNumberValue); ok {
				f, _ := bigNum.Value.Float64()
				num = types.NumberValue(f)
			} else {
				return nil, fmt.Errorf("arithmetic operation requires numbers")
			}
		}
		result = op(result, float64(num))
	}

	// Check if result is too large for safe float64 integer representation
	if result > 1e15 || result < -1e15 {
		// Check if this looks like it should be an integer (for multiplication)
		if result == float64(int64(result)) {
			// Convert to big integer for better precision
			return types.NewBigNumberFromInt64(int64(result)), nil
		}
	}

	return types.NumberValue(result), nil
}

// Big integer arithmetic for operations that might overflow
func (e *Evaluator) evalBigArithmetic(values []types.Value, op func(float64, float64) float64) (types.Value, error) {
	// For big integer operations, we need to handle each operation type separately
	// This is a simplified version that works for multiplication

	first, err := toBigInt(values[0])
	if err != nil {
		return nil, err
	}

	result := new(big.Int).Set(first)

	for i := 1; i < len(values); i++ {
		val, err := toBigInt(values[i])
		if err != nil {
			return nil, err
		}

		// We need to determine the operation type from the function pointer
		// This is a bit tricky, so for now we'll handle multiplication specifically
		// Check if this looks like multiplication by testing with small numbers
		testResult := op(2.0, 3.0)
		if testResult == 6.0 {
			// Multiplication
			result.Mul(result, val)
		} else if testResult == 5.0 {
			// Addition
			result.Add(result, val)
		} else if testResult == -1.0 {
			// Subtraction
			result.Sub(result, val)
		} else {
			// For division and other operations, fall back to float64
			f1, _ := result.Float64()
			f2, _ := val.Float64()
			finalResult := op(f1, f2)
			return types.NumberValue(finalResult), nil
		}
	}

	return types.NewBigNumberValue(result), nil
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

	// Handle big numbers
	if _, ok := first.(*types.BigNumberValue); ok {
		firstBig, err := toBigInt(first)
		if err != nil {
			return nil, err
		}
		secondBig, err := toBigInt(second)
		if err != nil {
			return nil, err
		}
		if secondBig.Sign() == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		// For division, convert to float64 since result might not be integer
		f1, _ := firstBig.Float64()
		f2, _ := secondBig.Float64()
		return types.NumberValue(f1 / f2), nil
	}

	if _, ok := second.(*types.BigNumberValue); ok {
		firstBig, err := toBigInt(first)
		if err != nil {
			return nil, err
		}
		secondBig, err := toBigInt(second)
		if err != nil {
			return nil, err
		}
		if secondBig.Sign() == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		f1, _ := firstBig.Float64()
		f2, _ := secondBig.Float64()
		return types.NumberValue(f1 / f2), nil
	}

	// Regular number division
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

	// Handle big number comparisons
	firstBig, isBig1 := first.(*types.BigNumberValue)
	secondBig, isBig2 := second.(*types.BigNumberValue)

	if isBig1 && isBig2 {
		return types.BooleanValue(firstBig.Value.Cmp(secondBig.Value) == 0), nil
	}

	if isBig1 || isBig2 {
		// One is big, one is regular - convert both to big for comparison
		big1, err := toBigInt(first)
		if err != nil {
			return types.BooleanValue(false), nil
		}
		big2, err := toBigInt(second)
		if err != nil {
			return types.BooleanValue(false), nil
		}
		return types.BooleanValue(big1.Cmp(big2) == 0), nil
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

	// Handle big number comparisons
	_, isBig1 := first.(*types.BigNumberValue)
	_, isBig2 := second.(*types.BigNumberValue)

	if isBig1 || isBig2 {
		// At least one is a big number - convert both to big for comparison
		big1, err := toBigInt(first)
		if err != nil {
			return nil, fmt.Errorf("comparison requires numbers")
		}
		big2, err := toBigInt(second)
		if err != nil {
			return nil, fmt.Errorf("comparison requires numbers")
		}

		// Use big.Int.Cmp for precise comparison
		cmp := big1.Cmp(big2)

		// Determine the result based on the comparison function
		// We need to figure out which operation this is by testing with known values
		testResult := op(1.0, 2.0) // Test with 1 < 2

		if testResult { // This is < or <=
			secondTest := op(2.0, 2.0) // Test with 2 == 2
			if secondTest { // This is <=
				return types.BooleanValue(cmp <= 0), nil
			} else { // This is <
				return types.BooleanValue(cmp < 0), nil
			}
		} else { // This is > or >=
			secondTest := op(2.0, 2.0) // Test with 2 == 2
			if secondTest { // This is >=
				return types.BooleanValue(cmp >= 0), nil
			} else { // This is >
				return types.BooleanValue(cmp > 0), nil
			}
		}
	}

	// Both are regular numbers
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

// Logical operations

func (e *Evaluator) evalAnd(args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return types.BooleanValue(true), nil // empty and is true
	}

	for _, arg := range args {
		val, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}

		boolVal, ok := val.(types.BooleanValue)
		if !ok {
			return nil, fmt.Errorf("and requires boolean arguments")
		}

		if !boolVal {
			return types.BooleanValue(false), nil
		}
	}

	return types.BooleanValue(true), nil
}

func (e *Evaluator) evalOr(args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return types.BooleanValue(false), nil // empty or is false
	}

	for _, arg := range args {
		val, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}

		boolVal, ok := val.(types.BooleanValue)
		if !ok {
			return nil, fmt.Errorf("or requires boolean arguments")
		}

		if boolVal {
			return types.BooleanValue(true), nil
		}
	}

	return types.BooleanValue(false), nil
}

func (e *Evaluator) evalNot(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("not requires exactly 1 argument")
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	boolVal, ok := val.(types.BooleanValue)
	if !ok {
		return nil, fmt.Errorf("not requires a boolean argument")
	}

	return types.BooleanValue(!boolVal), nil
}

// Enhanced multiplication evaluation that detects potential overflow
func (e *Evaluator) evalMultiplication(args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("multiplication requires at least one argument")
	}

	// Evaluate all arguments first
	values := make([]types.Value, len(args))
	for i, arg := range args {
		val, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}
		values[i] = val
	}

	// Check if we should use big integers (either operands are big or result might be big)
	if shouldUseBigInt(values) || mightOverflowMultiplication(values) {
		return e.evalBigArithmetic(values, func(a, b float64) float64 { return a * b })
	}

	// Use regular float64 arithmetic for small numbers
	firstNum, ok := values[0].(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("multiplication requires numbers")
	}

	result := float64(firstNum)
	for i := 1; i < len(values); i++ {
		num, ok := values[i].(types.NumberValue)
		if !ok {
			return nil, fmt.Errorf("multiplication requires numbers")
		}
		result *= float64(num)
	}

	// Check if result is too large for safe float64 integer representation
	if result > 1e15 || result < -1e15 {
		// Check if this looks like it should be an integer
		if result == float64(int64(result)) {
			// Convert to big integer for better precision
			return types.NewBigNumberFromInt64(int64(result)), nil
		}
	}

	return types.NumberValue(result), nil
}

func (e *Evaluator) evalModulo(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("modulo requires exactly 2 arguments")
	}

	first, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	second, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Handle big numbers
	if _, ok := first.(*types.BigNumberValue); ok {
		firstBig, err := toBigInt(first)
		if err != nil {
			return nil, err
		}
		secondBig, err := toBigInt(second)
		if err != nil {
			return nil, err
		}
		if secondBig.Sign() == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		result := new(big.Int)
		result.Mod(firstBig, secondBig)
		return types.NewBigNumberValue(result), nil
	}

	if _, ok := second.(*types.BigNumberValue); ok {
		firstBig, err := toBigInt(first)
		if err != nil {
			return nil, err
		}
		secondBig, err := toBigInt(second)
		if err != nil {
			return nil, err
		}
		if secondBig.Sign() == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		result := new(big.Int)
		result.Mod(firstBig, secondBig)
		return types.NewBigNumberValue(result), nil
	}

	// Regular number modulo
	firstNum, ok := first.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("modulo requires numbers")
	}

	secondNum, ok := second.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("modulo requires numbers")
	}

	if secondNum == 0 {
		return nil, fmt.Errorf("modulo by zero")
	}

	// Convert to integers for modulo operation
	result := int64(firstNum) % int64(secondNum)
	return types.NumberValue(float64(result)), nil
}
