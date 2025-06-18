// Package math implements mathematical functions as a plugin
package math

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// MathPlugin implements mathematical functions
type MathPlugin struct {
	*plugins.BasePlugin
}

// Random number generator (initialize once)
var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// NewMathPlugin creates a new math plugin
func NewMathPlugin() *MathPlugin {
	return &MathPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"math",
			"1.0.0",
			"Mathematical functions (sqrt, pow, sin, cos, log, etc.)",
			[]string{}, // No dependencies
		),
	}
}

// Functions returns the list of functions provided by this plugin
func (p *MathPlugin) Functions() []string {
	return []string{
		"sqrt", "pow", "abs", "floor", "ceil", "round", "trunc",
		"sin", "cos", "tan", "asin", "acos", "atan", "atan2",
		"sinh", "cosh", "tanh", "log", "exp", "log10", "log2",
		"degrees", "radians", "min", "max", "sign", "mod",
		"pi", "e", "random",
	}
}

// RegisterFunctions registers all math functions with the registry
func (p *MathPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// Basic mathematical functions
	sqrtFunc := functions.NewFunction(
		"sqrt",
		registry.CategoryMath,
		1,
		"Square root: (sqrt 16) => 4",
		p.evalSqrt,
	)
	if err := reg.Register(sqrtFunc); err != nil {
		return err
	}

	powFunc := functions.NewFunction(
		"pow",
		registry.CategoryMath,
		2,
		"Power: (pow 2 3) => 8",
		p.evalPow,
	)
	if err := reg.Register(powFunc); err != nil {
		return err
	}

	absFunc := functions.NewFunction(
		"abs",
		registry.CategoryMath,
		1,
		"Absolute value: (abs -7) => 7",
		p.evalAbs,
	)
	if err := reg.Register(absFunc); err != nil {
		return err
	}

	// Rounding functions
	floorFunc := functions.NewFunction(
		"floor",
		registry.CategoryMath,
		1,
		"Floor (largest integer ≤ x): (floor 3.7) => 3",
		p.evalFloor,
	)
	if err := reg.Register(floorFunc); err != nil {
		return err
	}

	ceilFunc := functions.NewFunction(
		"ceil",
		registry.CategoryMath,
		1,
		"Ceiling (smallest integer ≥ x): (ceil 3.2) => 4",
		p.evalCeil,
	)
	if err := reg.Register(ceilFunc); err != nil {
		return err
	}

	roundFunc := functions.NewFunction(
		"round",
		registry.CategoryMath,
		1,
		"Round to nearest integer: (round 3.6) => 4",
		p.evalRound,
	)
	if err := reg.Register(roundFunc); err != nil {
		return err
	}

	truncFunc := functions.NewFunction(
		"trunc",
		registry.CategoryMath,
		1,
		"Truncate towards zero: (trunc 3.7) => 3",
		p.evalTrunc,
	)
	if err := reg.Register(truncFunc); err != nil {
		return err
	}

	// Trigonometric functions
	sinFunc := functions.NewFunction(
		"sin",
		registry.CategoryMath,
		1,
		"Sine: (sin 0) => 0",
		p.evalSin,
	)
	if err := reg.Register(sinFunc); err != nil {
		return err
	}

	cosFunc := functions.NewFunction(
		"cos",
		registry.CategoryMath,
		1,
		"Cosine: (cos 0) => 1",
		p.evalCos,
	)
	if err := reg.Register(cosFunc); err != nil {
		return err
	}

	tanFunc := functions.NewFunction(
		"tan",
		registry.CategoryMath,
		1,
		"Tangent: (tan 0) => 0",
		p.evalTan,
	)
	if err := reg.Register(tanFunc); err != nil {
		return err
	}

	// Inverse trigonometric functions
	asinFunc := functions.NewFunction(
		"asin",
		registry.CategoryMath,
		1,
		"Arc sine: (asin 0.5) => 0.5236 (π/6)",
		p.evalAsin,
	)
	if err := reg.Register(asinFunc); err != nil {
		return err
	}

	acosFunc := functions.NewFunction(
		"acos",
		registry.CategoryMath,
		1,
		"Arc cosine: (acos 0.5) => 1.0472 (π/3)",
		p.evalAcos,
	)
	if err := reg.Register(acosFunc); err != nil {
		return err
	}

	atanFunc := functions.NewFunction(
		"atan",
		registry.CategoryMath,
		1,
		"Arc tangent: (atan 1) => 0.7854 (π/4)",
		p.evalAtan,
	)
	if err := reg.Register(atanFunc); err != nil {
		return err
	}

	atan2Func := functions.NewFunction(
		"atan2",
		registry.CategoryMath,
		2,
		"Two-argument arc tangent: (atan2 1 1) => 0.7854 (π/4)",
		p.evalAtan2,
	)
	if err := reg.Register(atan2Func); err != nil {
		return err
	}

	// Hyperbolic functions
	sinhFunc := functions.NewFunction(
		"sinh",
		registry.CategoryMath,
		1,
		"Hyperbolic sine: (sinh 0) => 0",
		p.evalSinh,
	)
	if err := reg.Register(sinhFunc); err != nil {
		return err
	}

	coshFunc := functions.NewFunction(
		"cosh",
		registry.CategoryMath,
		1,
		"Hyperbolic cosine: (cosh 0) => 1",
		p.evalCosh,
	)
	if err := reg.Register(coshFunc); err != nil {
		return err
	}

	tanhFunc := functions.NewFunction(
		"tanh",
		registry.CategoryMath,
		1,
		"Hyperbolic tangent: (tanh 0) => 0",
		p.evalTanh,
	)
	if err := reg.Register(tanhFunc); err != nil {
		return err
	}

	// Logarithmic and exponential functions
	logFunc := functions.NewFunction(
		"log",
		registry.CategoryMath,
		1,
		"Natural logarithm: (log (e)) => 1",
		p.evalLog,
	)
	if err := reg.Register(logFunc); err != nil {
		return err
	}

	expFunc := functions.NewFunction(
		"exp",
		registry.CategoryMath,
		1,
		"Exponential (e^x): (exp 1) => 2.7183",
		p.evalExp,
	)
	if err := reg.Register(expFunc); err != nil {
		return err
	}

	log10Func := functions.NewFunction(
		"log10",
		registry.CategoryMath,
		1,
		"Base-10 logarithm: (log10 100) => 2",
		p.evalLog10,
	)
	if err := reg.Register(log10Func); err != nil {
		return err
	}

	log2Func := functions.NewFunction(
		"log2",
		registry.CategoryMath,
		1,
		"Base-2 logarithm: (log2 8) => 3",
		p.evalLog2,
	)
	if err := reg.Register(log2Func); err != nil {
		return err
	}

	// Angle conversion
	degreesFunc := functions.NewFunction(
		"degrees",
		registry.CategoryMath,
		1,
		"Convert radians to degrees: (degrees (pi)) => 180",
		p.evalDegrees,
	)
	if err := reg.Register(degreesFunc); err != nil {
		return err
	}

	radiansFunc := functions.NewFunction(
		"radians",
		registry.CategoryMath,
		1,
		"Convert degrees to radians: (radians 180) => 3.1416",
		p.evalRadians,
	)
	if err := reg.Register(radiansFunc); err != nil {
		return err
	}

	// Statistical functions
	minFunc := functions.NewFunction(
		"min",
		registry.CategoryMath,
		-1,
		"Minimum of numbers: (min 5 2 8 1) => 1",
		p.evalMin,
	)
	if err := reg.Register(minFunc); err != nil {
		return err
	}

	maxFunc := functions.NewFunction(
		"max",
		registry.CategoryMath,
		-1,
		"Maximum of numbers: (max 5 2 8 1) => 8",
		p.evalMax,
	)
	if err := reg.Register(maxFunc); err != nil {
		return err
	}

	// Utility functions
	signFunc := functions.NewFunction(
		"sign",
		registry.CategoryMath,
		1,
		"Sign of number (-1, 0, or 1): (sign -5) => -1",
		p.evalSign,
	)
	if err := reg.Register(signFunc); err != nil {
		return err
	}

	modFunc := functions.NewFunction(
		"mod",
		registry.CategoryMath,
		2,
		"Mathematical modulo: (mod 7 3) => 1",
		p.evalMod,
	)
	if err := reg.Register(modFunc); err != nil {
		return err
	}

	// Mathematical constants
	piFunc := functions.NewFunction(
		"pi",
		registry.CategoryMath,
		0,
		"Pi constant: (pi) => 3.1416",
		p.evalPi,
	)
	if err := reg.Register(piFunc); err != nil {
		return err
	}

	eFunc := functions.NewFunction(
		"e",
		registry.CategoryMath,
		0,
		"Euler's number: (e) => 2.7183",
		p.evalE,
	)
	if err := reg.Register(eFunc); err != nil {
		return err
	}

	// Random number generation
	randomFunc := functions.NewFunction(
		"random",
		registry.CategoryMath,
		-1, // Variadic: 0, 1, or 2 arguments
		"Random number: (random) => 0-1, (random 10) => 0-9, (random 5 15) => 5-14",
		p.evalRandom,
	)
	if err := reg.Register(randomFunc); err != nil {
		return err
	}

	return nil
}

// Helper function to extract numeric value from any number type
func (p *MathPlugin) extractFloat64(val types.Value) (float64, error) {
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

// Helper function to create appropriate number type
func (p *MathPlugin) createNumber(f float64) types.Value {
	// For simplicity, always return regular NumberValue
	// In a full implementation, we might check for overflow and use BigNumber
	return types.NumberValue(f)
}

// Square root
func (p *MathPlugin) evalSqrt(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sqrt requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("sqrt: %v", err)
	}

	if num < 0 {
		return nil, fmt.Errorf("sqrt: cannot compute square root of negative number: %v", num)
	}

	result := math.Sqrt(num)
	return p.createNumber(result), nil
}

// Power
func (p *MathPlugin) evalPow(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("pow requires exactly 2 arguments, got %d", len(args))
	}

	baseVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	expVal, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	base, err := p.extractFloat64(baseVal)
	if err != nil {
		return nil, fmt.Errorf("pow base: %v", err)
	}

	exp, err := p.extractFloat64(expVal)
	if err != nil {
		return nil, fmt.Errorf("pow exponent: %v", err)
	}

	result := math.Pow(base, exp)
	return p.createNumber(result), nil
}

// Absolute value
func (p *MathPlugin) evalAbs(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("abs requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("abs: %v", err)
	}

	result := math.Abs(num)
	return p.createNumber(result), nil
}

// Floor
func (p *MathPlugin) evalFloor(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("floor requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("floor: %v", err)
	}

	result := math.Floor(num)
	return p.createNumber(result), nil
}

// Ceiling
func (p *MathPlugin) evalCeil(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ceil requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("ceil: %v", err)
	}

	result := math.Ceil(num)
	return p.createNumber(result), nil
}

// Round
func (p *MathPlugin) evalRound(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("round requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("round: %v", err)
	}

	result := math.Round(num)
	return p.createNumber(result), nil
}

// Truncate
func (p *MathPlugin) evalTrunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trunc requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("trunc: %v", err)
	}

	result := math.Trunc(num)
	return p.createNumber(result), nil
}

// Sine
func (p *MathPlugin) evalSin(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sin requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("sin: %v", err)
	}

	result := math.Sin(num)
	return p.createNumber(result), nil
}

// Cosine
func (p *MathPlugin) evalCos(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cos requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("cos: %v", err)
	}

	result := math.Cos(num)
	return p.createNumber(result), nil
}

// Tangent
func (p *MathPlugin) evalTan(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("tan requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("tan: %v", err)
	}

	result := math.Tan(num)

	// Check for very large results (near vertical asymptotes)
	if math.Abs(result) > 1e15 {
		return nil, fmt.Errorf("tan: result too large (near asymptote)")
	}

	return p.createNumber(result), nil
}

// Arc sine
func (p *MathPlugin) evalAsin(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("asin requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("asin: %v", err)
	}

	if num < -1 || num > 1 {
		return nil, fmt.Errorf("asin: input must be in range [-1, 1], got %v", num)
	}

	result := math.Asin(num)
	return p.createNumber(result), nil
}

// Arc cosine
func (p *MathPlugin) evalAcos(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("acos requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("acos: %v", err)
	}

	if num < -1 || num > 1 {
		return nil, fmt.Errorf("acos: input must be in range [-1, 1], got %v", num)
	}

	result := math.Acos(num)
	return p.createNumber(result), nil
}

// Arc tangent
func (p *MathPlugin) evalAtan(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("atan requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("atan: %v", err)
	}

	result := math.Atan(num)
	return p.createNumber(result), nil
}

// Two-argument arc tangent
func (p *MathPlugin) evalAtan2(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("atan2 requires exactly 2 arguments, got %d", len(args))
	}

	yVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	xVal, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	y, err := p.extractFloat64(yVal)
	if err != nil {
		return nil, fmt.Errorf("atan2 y: %v", err)
	}

	x, err := p.extractFloat64(xVal)
	if err != nil {
		return nil, fmt.Errorf("atan2 x: %v", err)
	}

	result := math.Atan2(y, x)
	return p.createNumber(result), nil
}

// Hyperbolic sine
func (p *MathPlugin) evalSinh(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sinh requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("sinh: %v", err)
	}

	result := math.Sinh(num)
	return p.createNumber(result), nil
}

// Hyperbolic cosine
func (p *MathPlugin) evalCosh(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cosh requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("cosh: %v", err)
	}

	result := math.Cosh(num)
	return p.createNumber(result), nil
}

// Hyperbolic tangent
func (p *MathPlugin) evalTanh(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("tanh requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("tanh: %v", err)
	}

	result := math.Tanh(num)
	return p.createNumber(result), nil
}

// Natural logarithm
func (p *MathPlugin) evalLog(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("log: %v", err)
	}

	if num <= 0 {
		return nil, fmt.Errorf("log: cannot compute logarithm of non-positive number: %v", num)
	}

	result := math.Log(num)
	return p.createNumber(result), nil
}

// Exponential
func (p *MathPlugin) evalExp(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("exp requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("exp: %v", err)
	}

	result := math.Exp(num)
	return p.createNumber(result), nil
}

// Base-10 logarithm
func (p *MathPlugin) evalLog10(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log10 requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("log10: %v", err)
	}

	if num <= 0 {
		return nil, fmt.Errorf("log10: input must be positive, got %v", num)
	}

	result := math.Log10(num)
	return p.createNumber(result), nil
}

// Base-2 logarithm
func (p *MathPlugin) evalLog2(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log2 requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("log2: %v", err)
	}

	if num <= 0 {
		return nil, fmt.Errorf("log2: input must be positive, got %v", num)
	}

	result := math.Log2(num)
	return p.createNumber(result), nil
}

// Convert radians to degrees
func (p *MathPlugin) evalDegrees(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("degrees requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("degrees: %v", err)
	}

	result := num * 180.0 / math.Pi
	return p.createNumber(result), nil
}

// Convert degrees to radians
func (p *MathPlugin) evalRadians(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("radians requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("radians: %v", err)
	}

	result := num * math.Pi / 180.0
	return p.createNumber(result), nil
}

// Minimum of multiple numbers
func (p *MathPlugin) evalMin(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("min requires at least 1 argument, got %d", len(args))
	}

	// Evaluate first argument
	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	result, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("min first argument: %v", err)
	}

	// If more than one argument, find minimum
	for i := 1; i < len(args); i++ {
		val, err := evaluator.Eval(args[i])
		if err != nil {
			return nil, err
		}

		num, err := p.extractFloat64(val)
		if err != nil {
			return nil, fmt.Errorf("min argument %d: %v", i+1, err)
		}

		if num < result {
			result = num
		}
	}

	return p.createNumber(result), nil
}

// Maximum of multiple numbers
func (p *MathPlugin) evalMax(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("max requires at least 1 argument, got %d", len(args))
	}

	// Evaluate first argument
	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	result, err := p.extractFloat64(val)
	if err != nil {
		return nil, fmt.Errorf("max first argument: %v", err)
	}

	// If more than one argument, find maximum
	for i := 1; i < len(args); i++ {
		val, err := evaluator.Eval(args[i])
		if err != nil {
			return nil, err
		}

		num, err := p.extractFloat64(val)
		if err != nil {
			return nil, fmt.Errorf("max argument %d: %v", i+1, err)
		}

		if num > result {
			result = num
		}
	}

	return p.createNumber(result), nil
}

// Sign of number
func (p *MathPlugin) evalSign(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sign requires exactly 1 argument, got %d", len(args))
	}

	val, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	num, err := p.extractFloat64(val)
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

	return p.createNumber(result), nil
}

// Mathematical modulo
func (p *MathPlugin) evalMod(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("mod requires exactly 2 arguments, got %d", len(args))
	}

	val1, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	val2, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	num1, err := p.extractFloat64(val1)
	if err != nil {
		return nil, fmt.Errorf("mod first argument: %v", err)
	}

	num2, err := p.extractFloat64(val2)
	if err != nil {
		return nil, fmt.Errorf("mod second argument: %v", err)
	}

	if num2 == 0 {
		return nil, fmt.Errorf("mod: division by zero")
	}

	result := math.Mod(num1, num2)
	return p.createNumber(result), nil
}

// Pi constant
func (p *MathPlugin) evalPi(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("pi requires no arguments, got %d", len(args))
	}
	return p.createNumber(math.Pi), nil
}

// Euler's number
func (p *MathPlugin) evalE(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("e requires no arguments, got %d", len(args))
	}
	return p.createNumber(math.E), nil
}

// Random number generation
func (p *MathPlugin) evalRandom(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	switch len(args) {
	case 0:
		// Random float between 0 and 1
		result := rng.Float64()
		return p.createNumber(result), nil

	case 1:
		// Random integer between 0 and n-1
		maxVal, err := evaluator.Eval(args[0])
		if err != nil {
			return nil, err
		}

		max, err := p.extractFloat64(maxVal)
		if err != nil {
			return nil, fmt.Errorf("random max: %v", err)
		}

		if max <= 0 {
			return nil, fmt.Errorf("random: max must be positive, got %v", max)
		}

		result := float64(rng.Intn(int(max)))
		return p.createNumber(result), nil

	case 2:
		// Random integer between min and max-1
		minVal, err := evaluator.Eval(args[0])
		if err != nil {
			return nil, err
		}

		maxVal, err := evaluator.Eval(args[1])
		if err != nil {
			return nil, err
		}

		min, err := p.extractFloat64(minVal)
		if err != nil {
			return nil, fmt.Errorf("random min: %v", err)
		}

		max, err := p.extractFloat64(maxVal)
		if err != nil {
			return nil, fmt.Errorf("random max: %v", err)
		}

		if min >= max {
			return nil, fmt.Errorf("random: min must be less than max, got min=%v, max=%v", min, max)
		}

		result := min + float64(rng.Intn(int(max-min)))
		return p.createNumber(result), nil

	default:
		return nil, fmt.Errorf("random requires 0, 1, or 2 arguments, got %d", len(args))
	}
}
