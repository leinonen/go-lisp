package math

import (
	"math"
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// Mock evaluator for testing
type mockEvaluator struct {
	env *evaluator.Environment
}

func newMockEvaluator() *mockEvaluator {
	return &mockEvaluator{
		env: evaluator.NewEnvironment(),
	}
}

func (me *mockEvaluator) Eval(expr types.Expr) (types.Value, error) {
	// For simple literals, convert them directly
	switch e := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(e.Value), nil
	case *types.StringExpr:
		return types.StringValue(e.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(e.Value), nil
	default:
		// For values wrapped in valueExpr, return them as-is
		if ve, ok := expr.(valueExpr); ok {
			return ve.value, nil
		}
		// For values, return them as-is
		if val, ok := expr.(types.Value); ok {
			return val, nil
		}
		return nil, nil
	}
}

func (me *mockEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	return nil, nil // Not used in these tests
}

func (me *mockEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	// For testing purposes, just call regular Eval
	// In a real implementation, this would use the bindings
	return me.Eval(expr)
}

// Helper function to wrap values as expressions
func wrapValue(value types.Value) types.Expr {
	return valueExpr{value}
}

type valueExpr struct {
	value types.Value
}

func (ve valueExpr) String() string {
	return ve.value.String()
}

func (ve valueExpr) GetPosition() types.Position {
	return types.Position{Line: 1, Column: 1}
}

// Helper function for floating point comparison
func floatEqual(a, b float64, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestMathPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewMathPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{
		"sqrt", "pow", "abs", "floor", "ceil", "round", "trunc",
		"sin", "cos", "tan", "asin", "acos", "atan", "atan2",
		"sinh", "cosh", "tanh", "log", "exp", "log10", "log2",
		"degrees", "radians", "min", "max", "sign", "mod",
		"pi", "e", "random",
	}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestMathPlugin_BasicFunctions(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		args     []types.Expr
		expected float64
	}{
		{
			name:     "sqrt of 16",
			function: plugin.evalSqrt,
			args:     []types.Expr{&types.NumberExpr{Value: 16}},
			expected: 4,
		},
		{
			name:     "sqrt of 2",
			function: plugin.evalSqrt,
			args:     []types.Expr{&types.NumberExpr{Value: 2}},
			expected: math.Sqrt(2),
		},
		{
			name:     "pow 2^3",
			function: plugin.evalPow,
			args: []types.Expr{
				&types.NumberExpr{Value: 2},
				&types.NumberExpr{Value: 3},
			},
			expected: 8,
		},
		{
			name:     "abs of -7",
			function: plugin.evalAbs,
			args:     []types.Expr{&types.NumberExpr{Value: -7}},
			expected: 7,
		},
		{
			name:     "abs of 7",
			function: plugin.evalAbs,
			args:     []types.Expr{&types.NumberExpr{Value: 7}},
			expected: 7,
		},
		{
			name:     "floor of 3.7",
			function: plugin.evalFloor,
			args:     []types.Expr{&types.NumberExpr{Value: 3.7}},
			expected: 3,
		},
		{
			name:     "floor of -3.7",
			function: plugin.evalFloor,
			args:     []types.Expr{&types.NumberExpr{Value: -3.7}},
			expected: -4,
		},
		{
			name:     "ceil of 3.2",
			function: plugin.evalCeil,
			args:     []types.Expr{&types.NumberExpr{Value: 3.2}},
			expected: 4,
		},
		{
			name:     "ceil of -3.2",
			function: plugin.evalCeil,
			args:     []types.Expr{&types.NumberExpr{Value: -3.2}},
			expected: -3,
		},
		{
			name:     "round of 3.5",
			function: plugin.evalRound,
			args:     []types.Expr{&types.NumberExpr{Value: 3.5}},
			expected: 4,
		},
		{
			name:     "round of 3.4",
			function: plugin.evalRound,
			args:     []types.Expr{&types.NumberExpr{Value: 3.4}},
			expected: 3,
		},
		{
			name:     "trunc of 3.9",
			function: plugin.evalTrunc,
			args:     []types.Expr{&types.NumberExpr{Value: 3.9}},
			expected: 3,
		},
		{
			name:     "trunc of -3.9",
			function: plugin.evalTrunc,
			args:     []types.Expr{&types.NumberExpr{Value: -3.9}},
			expected: -3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.function(evaluator, tt.args)
			if err != nil {
				t.Fatalf("Function failed: %v", err)
			}

			number, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(number), tt.expected, 1e-10) {
				t.Errorf("Expected %v, got %v", tt.expected, float64(number))
			}
		})
	}
}

func TestMathPlugin_TrigonometricFunctions(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		input    float64
		expected float64
	}{
		{
			name:     "sin(0)",
			function: plugin.evalSin,
			input:    0,
			expected: 0,
		},
		{
			name:     "sin(π/2)",
			function: plugin.evalSin,
			input:    math.Pi / 2,
			expected: 1,
		},
		{
			name:     "cos(0)",
			function: plugin.evalCos,
			input:    0,
			expected: 1,
		},
		{
			name:     "cos(π)",
			function: plugin.evalCos,
			input:    math.Pi,
			expected: -1,
		},
		{
			name:     "tan(0)",
			function: plugin.evalTan,
			input:    0,
			expected: 0,
		},
		{
			name:     "tan(π/4)",
			function: plugin.evalTan,
			input:    math.Pi / 4,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{&types.NumberExpr{Value: tt.input}}
			result, err := tt.function(evaluator, args)
			if err != nil {
				t.Fatalf("Function failed: %v", err)
			}

			number, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(number), tt.expected, 1e-10) {
				t.Errorf("Expected %v, got %v", tt.expected, float64(number))
			}
		})
	}
}

func TestMathPlugin_InverseTrigonometricFunctions(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		input    float64
		expected float64
	}{
		{
			name:     "asin(0)",
			function: plugin.evalAsin,
			input:    0,
			expected: 0,
		},
		{
			name:     "asin(1)",
			function: plugin.evalAsin,
			input:    1,
			expected: math.Pi / 2,
		},
		{
			name:     "acos(1)",
			function: plugin.evalAcos,
			input:    1,
			expected: 0,
		},
		{
			name:     "acos(0)",
			function: plugin.evalAcos,
			input:    0,
			expected: math.Pi / 2,
		},
		{
			name:     "atan(0)",
			function: plugin.evalAtan,
			input:    0,
			expected: 0,
		},
		{
			name:     "atan(1)",
			function: plugin.evalAtan,
			input:    1,
			expected: math.Pi / 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{&types.NumberExpr{Value: tt.input}}
			result, err := tt.function(evaluator, args)
			if err != nil {
				t.Fatalf("Function failed: %v", err)
			}

			number, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(number), tt.expected, 1e-10) {
				t.Errorf("Expected %v, got %v", tt.expected, float64(number))
			}
		})
	}
}

func TestMathPlugin_LogarithmicFunctions(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		input    float64
		expected float64
	}{
		{
			name:     "log(e)",
			function: plugin.evalLog,
			input:    math.E,
			expected: 1,
		},
		{
			name:     "log(1)",
			function: plugin.evalLog,
			input:    1,
			expected: 0,
		},
		{
			name:     "exp(0)",
			function: plugin.evalExp,
			input:    0,
			expected: 1,
		},
		{
			name:     "exp(1)",
			function: plugin.evalExp,
			input:    1,
			expected: math.E,
		},
		{
			name:     "log10(10)",
			function: plugin.evalLog10,
			input:    10,
			expected: 1,
		},
		{
			name:     "log10(100)",
			function: plugin.evalLog10,
			input:    100,
			expected: 2,
		},
		{
			name:     "log2(2)",
			function: plugin.evalLog2,
			input:    2,
			expected: 1,
		},
		{
			name:     "log2(8)",
			function: plugin.evalLog2,
			input:    8,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{&types.NumberExpr{Value: tt.input}}
			result, err := tt.function(evaluator, args)
			if err != nil {
				t.Fatalf("Function failed: %v", err)
			}

			number, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(number), tt.expected, 1e-10) {
				t.Errorf("Expected %v, got %v", tt.expected, float64(number))
			}
		})
	}
}

func TestMathPlugin_MinMax(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		args     []types.Expr
		expected float64
	}{
		{
			name:     "min of two numbers",
			function: plugin.evalMin,
			args: []types.Expr{
				&types.NumberExpr{Value: 5},
				&types.NumberExpr{Value: 3},
			},
			expected: 3,
		},
		{
			name:     "min of multiple numbers",
			function: plugin.evalMin,
			args: []types.Expr{
				&types.NumberExpr{Value: 5},
				&types.NumberExpr{Value: 3},
				&types.NumberExpr{Value: 8},
				&types.NumberExpr{Value: 1},
			},
			expected: 1,
		},
		{
			name:     "max of two numbers",
			function: plugin.evalMax,
			args: []types.Expr{
				&types.NumberExpr{Value: 5},
				&types.NumberExpr{Value: 3},
			},
			expected: 5,
		},
		{
			name:     "max of multiple numbers",
			function: plugin.evalMax,
			args: []types.Expr{
				&types.NumberExpr{Value: 5},
				&types.NumberExpr{Value: 3},
				&types.NumberExpr{Value: 8},
				&types.NumberExpr{Value: 1},
			},
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.function(evaluator, tt.args)
			if err != nil {
				t.Fatalf("Function failed: %v", err)
			}

			number, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(number), tt.expected, 1e-10) {
				t.Errorf("Expected %v, got %v", tt.expected, float64(number))
			}
		})
	}
}

func TestMathPlugin_AngleConversion(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		input    float64
		expected float64
	}{
		{
			name:     "degrees(π)",
			function: plugin.evalDegrees,
			input:    math.Pi,
			expected: 180,
		},
		{
			name:     "degrees(π/2)",
			function: plugin.evalDegrees,
			input:    math.Pi / 2,
			expected: 90,
		},
		{
			name:     "radians(180)",
			function: plugin.evalRadians,
			input:    180,
			expected: math.Pi,
		},
		{
			name:     "radians(90)",
			function: plugin.evalRadians,
			input:    90,
			expected: math.Pi / 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{&types.NumberExpr{Value: tt.input}}
			result, err := tt.function(evaluator, args)
			if err != nil {
				t.Fatalf("Function failed: %v", err)
			}

			number, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(number), tt.expected, 1e-10) {
				t.Errorf("Expected %v, got %v", tt.expected, float64(number))
			}
		})
	}
}

func TestMathPlugin_SpecialFunctions(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	// Test sign function
	signTests := []struct {
		input    float64
		expected float64
	}{
		{5, 1},
		{-5, -1},
		{0, 0},
	}

	for _, tt := range signTests {
		t.Run("sign", func(t *testing.T) {
			args := []types.Expr{&types.NumberExpr{Value: tt.input}}
			result, err := plugin.evalSign(evaluator, args)
			if err != nil {
				t.Fatalf("evalSign failed: %v", err)
			}

			number, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if float64(number) != tt.expected {
				t.Errorf("sign(%v): expected %v, got %v", tt.input, tt.expected, float64(number))
			}
		})
	}

	// Test mod function
	args := []types.Expr{
		&types.NumberExpr{Value: 10},
		&types.NumberExpr{Value: 3},
	}
	result, err := plugin.evalMod(evaluator, args)
	if err != nil {
		t.Fatalf("evalMod failed: %v", err)
	}

	number, ok := result.(types.NumberValue)
	if !ok {
		t.Fatalf("Expected NumberValue, got %T", result)
	}

	if float64(number) != 1 {
		t.Errorf("mod(10, 3): expected 1, got %v", float64(number))
	}

	// Test pi constant
	result, err = plugin.evalPi(evaluator, []types.Expr{})
	if err != nil {
		t.Fatalf("evalPi failed: %v", err)
	}

	number, ok = result.(types.NumberValue)
	if !ok {
		t.Fatalf("Expected NumberValue, got %T", result)
	}

	if !floatEqual(float64(number), math.Pi, 1e-10) {
		t.Errorf("pi: expected %v, got %v", math.Pi, float64(number))
	}

	// Test e constant
	result, err = plugin.evalE(evaluator, []types.Expr{})
	if err != nil {
		t.Fatalf("evalE failed: %v", err)
	}

	number, ok = result.(types.NumberValue)
	if !ok {
		t.Fatalf("Expected NumberValue, got %T", result)
	}

	if !floatEqual(float64(number), math.E, 1e-10) {
		t.Errorf("e: expected %v, got %v", math.E, float64(number))
	}
}

func TestMathPlugin_Random(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	// Test random function returns a value between 0 and 1
	result, err := plugin.evalRandom(evaluator, []types.Expr{})
	if err != nil {
		t.Fatalf("evalRandom failed: %v", err)
	}

	number, ok := result.(types.NumberValue)
	if !ok {
		t.Fatalf("Expected NumberValue, got %T", result)
	}

	val := float64(number)
	if val < 0 || val >= 1 {
		t.Errorf("random(): expected value in [0,1), got %v", val)
	}

	// Test multiple calls return different values (with high probability)
	result2, err := plugin.evalRandom(evaluator, []types.Expr{})
	if err != nil {
		t.Fatalf("evalRandom failed: %v", err)
	}

	number2, ok := result2.(types.NumberValue)
	if !ok {
		t.Fatalf("Expected NumberValue, got %T", result2)
	}

	// Very unlikely to be equal
	if float64(number) == float64(number2) {
		t.Logf("Warning: Two random() calls returned the same value: %v", float64(number))
	}
}

func TestMathPlugin_ErrorCases(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		args     []types.Expr
	}{
		{
			name:     "sqrt of negative number",
			function: plugin.evalSqrt,
			args:     []types.Expr{&types.NumberExpr{Value: -1}},
		},
		{
			name:     "log of zero",
			function: plugin.evalLog,
			args:     []types.Expr{&types.NumberExpr{Value: 0}},
		},
		{
			name:     "log of negative number",
			function: plugin.evalLog,
			args:     []types.Expr{&types.NumberExpr{Value: -1}},
		},
		{
			name:     "asin of value > 1",
			function: plugin.evalAsin,
			args:     []types.Expr{&types.NumberExpr{Value: 2}},
		},
		{
			name:     "acos of value > 1",
			function: plugin.evalAcos,
			args:     []types.Expr{&types.NumberExpr{Value: 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.function(evaluator, tt.args)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
		})
	}
}

func TestMathPlugin_InvalidArguments(t *testing.T) {
	plugin := NewMathPlugin()
	evaluator := newMockEvaluator()

	// Test with non-number arguments
	nonNumber := types.StringValue("not-a-number")

	tests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		args     []types.Expr
	}{
		{
			name:     "sqrt with non-number",
			function: plugin.evalSqrt,
			args:     []types.Expr{wrapValue(nonNumber)},
		},
		{
			name:     "pow with non-number first arg",
			function: plugin.evalPow,
			args: []types.Expr{
				wrapValue(nonNumber),
				&types.NumberExpr{Value: 2},
			},
		},
		{
			name:     "pow with non-number second arg",
			function: plugin.evalPow,
			args: []types.Expr{
				&types.NumberExpr{Value: 2},
				wrapValue(nonNumber),
			},
		},
		{
			name:     "min with no arguments",
			function: plugin.evalMin,
			args:     []types.Expr{},
		},
		{
			name:     "max with no arguments",
			function: plugin.evalMax,
			args:     []types.Expr{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.function(evaluator, tt.args)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
		})
	}
}
