// Package evaluator_math contains mathematical functions
package evaluator

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Mathematical functions

// Helper function to extract numeric value from any number type
func extractFloat64(val types.Value) (float64, error) {
	switch v := val.(type) {
	case types.NumberValue:
		return float64(v), nil
	case *types.BigNumberValue:
		f, _ := v.Value.Float64()
		return f, nil
	default:
		return 0, fmt.Errorf("expected number, got %T", val)
	}
}

// Helper function to create appropriate number type (regular or big)
func createNumber(f float64) types.Value {
	// If the number is an integer and within reasonable bounds, return as NumberValue
	if f == float64(int64(f)) && f >= -1e15 && f <= 1e15 {
		return types.NumberValue(f)
	}
	// For very large numbers or fractional results, use NumberValue
	return types.NumberValue(f)
}

// evalSqrt computes square root
func (e *Evaluator) evalSqrt(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sqrt requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("sqrt: %v", err)
	}

	if num < 0 {
		return nil, fmt.Errorf("sqrt: cannot compute square root of negative number: %v", num)
	}

	result := math.Sqrt(num)
	return createNumber(result), nil
}

// evalPow computes power (base^exponent)
func (e *Evaluator) evalPow(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("pow requires exactly 2 arguments (base exponent), got %d", len(args))
	}

	baseVal, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	expVal, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	base, err := extractFloat64(baseVal)
	if err != nil {
		return nil, fmt.Errorf("pow base: %v", err)
	}

	exp, err := extractFloat64(expVal)
	if err != nil {
		return nil, fmt.Errorf("pow exponent: %v", err)
	}

	result := math.Pow(base, exp)

	// Check for invalid results
	if math.IsNaN(result) {
		return nil, fmt.Errorf("pow: result is not a number")
	}
	if math.IsInf(result, 0) {
		return nil, fmt.Errorf("pow: result is infinite")
	}

	return createNumber(result), nil
}

// evalSin computes sine
func (e *Evaluator) evalSin(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sin requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("sin: %v", err)
	}

	result := math.Sin(num)
	return createNumber(result), nil
}

// evalCos computes cosine
func (e *Evaluator) evalCos(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cos requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("cos: %v", err)
	}

	result := math.Cos(num)
	return createNumber(result), nil
}

// evalTan computes tangent
func (e *Evaluator) evalTan(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("tan requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("tan: %v", err)
	}

	result := math.Tan(num)

	// Check for very large results (near vertical asymptotes)
	if math.Abs(result) > 1e15 {
		return nil, fmt.Errorf("tan: result too large (near asymptote)")
	}

	return createNumber(result), nil
}

// evalLog computes natural logarithm
func (e *Evaluator) evalLog(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("log: %v", err)
	}

	if num <= 0 {
		return nil, fmt.Errorf("log: cannot compute logarithm of non-positive number: %v", num)
	}

	result := math.Log(num)
	return createNumber(result), nil
}

// evalExp computes e^x
func (e *Evaluator) evalExp(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("exp requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("exp: %v", err)
	}

	result := math.Exp(num)

	// Check for overflow
	if math.IsInf(result, 0) {
		return nil, fmt.Errorf("exp: result is infinite (overflow)")
	}

	return createNumber(result), nil
}

// evalFloor computes floor (largest integer ≤ x)
func (e *Evaluator) evalFloor(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("floor requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("floor: %v", err)
	}

	result := math.Floor(num)
	return createNumber(result), nil
}

// evalCeil computes ceiling (smallest integer ≥ x)
func (e *Evaluator) evalCeil(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ceil requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("ceil: %v", err)
	}

	result := math.Ceil(num)
	return createNumber(result), nil
}

// evalRound computes round to nearest integer
func (e *Evaluator) evalRound(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("round requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("round: %v", err)
	}

	result := math.Round(num)
	return createNumber(result), nil
}

// evalAbs computes absolute value
func (e *Evaluator) evalAbs(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("abs requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("abs: %v", err)
	}

	result := math.Abs(num)
	return createNumber(result), nil
}

// evalMin computes minimum of two numbers
func (e *Evaluator) evalMin(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("min requires exactly 2 arguments, got %d", len(args))
	}

	val1, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	val2, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	num1, err := extractFloat64(val1)
	if err != nil {
		return nil, fmt.Errorf("min first argument: %v", err)
	}

	num2, err := extractFloat64(val2)
	if err != nil {
		return nil, fmt.Errorf("min second argument: %v", err)
	}

	result := math.Min(num1, num2)
	return createNumber(result), nil
}

// evalMax computes maximum of two numbers
func (e *Evaluator) evalMax(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("max requires exactly 2 arguments, got %d", len(args))
	}

	val1, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	val2, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	num1, err := extractFloat64(val1)
	if err != nil {
		return nil, fmt.Errorf("max first argument: %v", err)
	}

	num2, err := extractFloat64(val2)
	if err != nil {
		return nil, fmt.Errorf("max second argument: %v", err)
	}

	result := math.Max(num1, num2)
	return createNumber(result), nil
}

// Random number generator (initialize once)
var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// evalRandom generates random number
func (e *Evaluator) evalRandom(args []types.Expr) (types.Value, error) {
	switch len(args) {
	case 0:
		// Random float between 0 and 1
		result := rng.Float64()
		return createNumber(result), nil
	case 1:
		// Random integer between 0 and n-1
		val, err := e.Eval(args[0])
		if err != nil {
			return nil, err
		}

		num, err := extractFloat64(val)
		if err != nil {
			return nil, fmt.Errorf("random: %v", err)
		}

		if num <= 0 {
			return nil, fmt.Errorf("random: upper bound must be positive, got %v", num)
		}

		if num != float64(int64(num)) {
			return nil, fmt.Errorf("random: upper bound must be an integer, got %v", num)
		}

		n := int64(num)
		result := rng.Int63n(n)
		return createNumber(float64(result)), nil
	case 2:
		// Random integer between min and max-1
		minVal, err := e.Eval(args[0])
		if err != nil {
			return nil, err
		}

		maxVal, err := e.Eval(args[1])
		if err != nil {
			return nil, err
		}

		minNum, err := extractFloat64(minVal)
		if err != nil {
			return nil, fmt.Errorf("random min: %v", err)
		}

		maxNum, err := extractFloat64(maxVal)
		if err != nil {
			return nil, fmt.Errorf("random max: %v", err)
		}

		if minNum != float64(int64(minNum)) {
			return nil, fmt.Errorf("random: min must be an integer, got %v", minNum)
		}

		if maxNum != float64(int64(maxNum)) {
			return nil, fmt.Errorf("random: max must be an integer, got %v", maxNum)
		}

		if maxNum <= minNum {
			return nil, fmt.Errorf("random: max must be greater than min, got min=%v max=%v", minNum, maxNum)
		}

		min := int64(minNum)
		max := int64(maxNum)
		result := rng.Int63n(max-min) + min
		return createNumber(float64(result)), nil
	default:
		return nil, fmt.Errorf("random requires 0, 1, or 2 arguments, got %d", len(args))
	}
}

// Mathematical constants
func (e *Evaluator) evalPi(args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("pi requires no arguments, got %d", len(args))
	}
	return createNumber(math.Pi), nil
}

func (e *Evaluator) evalE(args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("e requires no arguments, got %d", len(args))
	}
	return createNumber(math.E), nil
}

// Additional trigonometric functions

// evalAsin computes arcsine (inverse sine)
func (e *Evaluator) evalAsin(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("asin requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("asin: %v", err)
	}

	if num < -1 || num > 1 {
		return nil, fmt.Errorf("asin: input must be in range [-1, 1], got %v", num)
	}

	result := math.Asin(num)
	return createNumber(result), nil
}

// evalAcos computes arccosine (inverse cosine)
func (e *Evaluator) evalAcos(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("acos requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("acos: %v", err)
	}

	if num < -1 || num > 1 {
		return nil, fmt.Errorf("acos: input must be in range [-1, 1], got %v", num)
	}

	result := math.Acos(num)
	return createNumber(result), nil
}

// evalAtan computes arctangent (inverse tangent)
func (e *Evaluator) evalAtan(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("atan requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("atan: %v", err)
	}

	result := math.Atan(num)
	return createNumber(result), nil
}

// evalAtan2 computes atan2(y, x) - arctangent of y/x, handling all quadrants
func (e *Evaluator) evalAtan2(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("atan2 requires exactly 2 arguments, got %d", len(args))
	}

	yVal, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	xVal, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	y, err := extractFloat64(yVal)
	if err != nil {
		return nil, fmt.Errorf("atan2 y: %v", err)
	}

	x, err := extractFloat64(xVal)
	if err != nil {
		return nil, fmt.Errorf("atan2 x: %v", err)
	}

	result := math.Atan2(y, x)
	return createNumber(result), nil
}

// Hyperbolic functions

// evalSinh computes hyperbolic sine
func (e *Evaluator) evalSinh(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sinh requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("sinh: %v", err)
	}

	result := math.Sinh(num)
	return createNumber(result), nil
}

// evalCosh computes hyperbolic cosine
func (e *Evaluator) evalCosh(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cosh requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("cosh: %v", err)
	}

	result := math.Cosh(num)
	return createNumber(result), nil
}

// evalTanh computes hyperbolic tangent
func (e *Evaluator) evalTanh(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("tanh requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("tanh: %v", err)
	}

	result := math.Tanh(num)
	return createNumber(result), nil
}

// Angle conversion functions

// evalDegrees converts radians to degrees
func (e *Evaluator) evalDegrees(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("degrees requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("degrees: %v", err)
	}

	result := num * 180.0 / math.Pi
	return createNumber(result), nil
}

// evalRadians converts degrees to radians
func (e *Evaluator) evalRadians(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("radians requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("radians: %v", err)
	}

	result := num * math.Pi / 180.0
	return createNumber(result), nil
}

// Additional logarithm functions

// evalLog10 computes logarithm base 10
func (e *Evaluator) evalLog10(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log10 requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("log10: %v", err)
	}

	if num <= 0 {
		return nil, fmt.Errorf("log10: input must be positive, got %v", num)
	}

	result := math.Log10(num)
	return createNumber(result), nil
}

// evalLog2 computes logarithm base 2
func (e *Evaluator) evalLog2(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log2 requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("log2: %v", err)
	}

	if num <= 0 {
		return nil, fmt.Errorf("log2: input must be positive, got %v", num)
	}

	result := math.Log2(num)
	return createNumber(result), nil
}

// Additional utility functions

// evalTrunc truncates towards zero
func (e *Evaluator) evalTrunc(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trunc requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("trunc: %v", err)
	}

	result := math.Trunc(num)
	return createNumber(result), nil
}

// evalSign returns the sign of the number (-1, 0, or 1)
func (e *Evaluator) evalSign(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sign requires exactly 1 argument, got %d", len(args))
	}

	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("sign: %v", err)
	}

	var result float64
	if num > 0 {
		result = 1
	} else if num < 0 {
		result = -1
	} else {
		result = 0
	}

	return createNumber(result), nil
}

// evalMod computes modulo (like % but handles negative numbers properly)
func (e *Evaluator) evalMod(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("mod requires exactly 2 arguments, got %d", len(args))
	}

	xVal, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	yVal, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	x, err := extractFloat64(xVal)
	if err != nil {
		return nil, fmt.Errorf("mod x: %v", err)
	}

	y, err := extractFloat64(yVal)
	if err != nil {
		return nil, fmt.Errorf("mod y: %v", err)
	}

	if y == 0 {
		return nil, fmt.Errorf("mod: division by zero")
	}

	// Implement proper mathematical modulo that always returns a non-negative result
	// when the divisor is positive
	result := math.Mod(x, y)
	if result < 0 && y > 0 {
		result += y
	} else if result > 0 && y < 0 {
		result += y
	}

	return createNumber(result), nil
}
