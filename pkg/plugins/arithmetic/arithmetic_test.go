package arithmetic

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

func TestArithmeticPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewArithmeticPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"+", "-", "*", "/", "%"}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestArithmeticPlugin_Addition(t *testing.T) {
	plugin := NewArithmeticPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		args     []types.Expr
		expected float64
	}{
		{
			name:     "no arguments",
			args:     []types.Expr{},
			expected: 0,
		},
		{
			name:     "single argument",
			args:     []types.Expr{&types.NumberExpr{Value: 5}},
			expected: 5,
		},
		{
			name:     "two arguments",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 7}},
			expected: 10,
		},
		{
			name:     "multiple arguments",
			args:     []types.Expr{&types.NumberExpr{Value: 1}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 4}},
			expected: 10,
		},
		{
			name:     "negative numbers",
			args:     []types.Expr{&types.NumberExpr{Value: -5}, &types.NumberExpr{Value: 3}},
			expected: -2,
		},
		{
			name:     "decimal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 1.5}, &types.NumberExpr{Value: 2.5}},
			expected: 4.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalAdd(evaluator, tt.args)
			if err != nil {
				t.Fatalf("evalAdd failed: %v", err)
			}

			numResult, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(numResult), tt.expected, 1e-10) {
				t.Errorf("Expected %f, got %f", tt.expected, float64(numResult))
			}
		})
	}
}

func TestArithmeticPlugin_AdditionWithBigNumbers(t *testing.T) {
	plugin := NewArithmeticPlugin()
	evaluator := newMockEvaluator()

	// Test with big numbers
	bigNum := types.NewBigNumberFromInt64(1000000000000000000)
	regularNum := types.NumberValue(5)

	args := []types.Expr{
		wrapValue(bigNum),
		wrapValue(regularNum),
	}

	result, err := plugin.evalAdd(evaluator, args)
	if err != nil {
		t.Fatalf("evalAdd with big numbers failed: %v", err)
	}

	bigResult, ok := result.(*types.BigNumberValue)
	if !ok {
		t.Fatalf("Expected BigNumberValue, got %T", result)
	}

	// Check that result is greater than the original big number
	if bigResult.Value.Cmp(bigNum.Value) <= 0 {
		t.Errorf("Big number addition failed: result should be greater than original")
	}
}

func TestArithmeticPlugin_Subtraction(t *testing.T) {
	plugin := NewArithmeticPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    float64
		expectError bool
	}{
		{
			name:        "no arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:     "unary minus",
			args:     []types.Expr{&types.NumberExpr{Value: 5}},
			expected: -5,
		},
		{
			name:     "two arguments",
			args:     []types.Expr{&types.NumberExpr{Value: 10}, &types.NumberExpr{Value: 3}},
			expected: 7,
		},
		{
			name:     "multiple arguments",
			args:     []types.Expr{&types.NumberExpr{Value: 20}, &types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 3}},
			expected: 12,
		},
		{
			name:     "negative result",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 10}},
			expected: -7,
		},
		{
			name:     "decimal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 5.5}, &types.NumberExpr{Value: 2.3}},
			expected: 3.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalSubtract(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalSubtract failed: %v", err)
			}

			numResult, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(numResult), tt.expected, 1e-10) {
				t.Errorf("Expected %f, got %f", tt.expected, float64(numResult))
			}
		})
	}
}

func TestArithmeticPlugin_SubtractionWithBigNumbers(t *testing.T) {
	plugin := NewArithmeticPlugin()
	evaluator := newMockEvaluator()

	// Test unary minus with big number
	bigNum := types.NewBigNumberFromInt64(1000000000000000000)
	args := []types.Expr{wrapValue(bigNum)}

	result, err := plugin.evalSubtract(evaluator, args)
	if err != nil {
		t.Fatalf("evalSubtract with big number failed: %v", err)
	}

	bigResult, ok := result.(*types.BigNumberValue)
	if !ok {
		t.Fatalf("Expected BigNumberValue, got %T", result)
	}

	// Check that result is negative
	if bigResult.Value.Sign() != -1 {
		t.Errorf("Big number unary minus failed: result should be negative")
	}
}

func TestArithmeticPlugin_Multiplication(t *testing.T) {
	plugin := NewArithmeticPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		args     []types.Expr
		expected float64
	}{
		{
			name:     "no arguments",
			args:     []types.Expr{},
			expected: 1,
		},
		{
			name:     "single argument",
			args:     []types.Expr{&types.NumberExpr{Value: 5}},
			expected: 5,
		},
		{
			name:     "two arguments",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 4}},
			expected: 12,
		},
		{
			name:     "multiple arguments",
			args:     []types.Expr{&types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 4}},
			expected: 24,
		},
		{
			name:     "with zero",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 0}},
			expected: 0,
		},
		{
			name:     "negative numbers",
			args:     []types.Expr{&types.NumberExpr{Value: -2}, &types.NumberExpr{Value: 3}},
			expected: -6,
		},
		{
			name:     "decimal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 2.5}, &types.NumberExpr{Value: 4}},
			expected: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalMultiply(evaluator, tt.args)
			if err != nil {
				t.Fatalf("evalMultiply failed: %v", err)
			}

			numResult, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(numResult), tt.expected, 1e-10) {
				t.Errorf("Expected %f, got %f", tt.expected, float64(numResult))
			}
		})
	}
}

func TestArithmeticPlugin_MultiplicationWithBigNumbers(t *testing.T) {
	plugin := NewArithmeticPlugin()
	evaluator := newMockEvaluator()

	// Test with big numbers
	bigNum := types.NewBigNumberFromInt64(1000000000)
	regularNum := types.NumberValue(2)

	args := []types.Expr{
		wrapValue(bigNum),
		wrapValue(regularNum),
	}

	result, err := plugin.evalMultiply(evaluator, args)
	if err != nil {
		t.Fatalf("evalMultiply with big numbers failed: %v", err)
	}

	bigResult, ok := result.(*types.BigNumberValue)
	if !ok {
		t.Fatalf("Expected BigNumberValue, got %T", result)
	}

	// Check that result is twice the original big number
	expected := types.NewBigNumberFromInt64(2000000000)
	if bigResult.Value.Cmp(expected.Value) != 0 {
		t.Errorf("Big number multiplication failed: expected %s, got %s", expected.Value, bigResult.Value)
	}
}

func TestArithmeticPlugin_Division(t *testing.T) {
	plugin := NewArithmeticPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    float64
		expectError bool
	}{
		{
			name:        "no arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:     "single argument (reciprocal)",
			args:     []types.Expr{&types.NumberExpr{Value: 4}},
			expected: 0.25,
		},
		{
			name:     "two arguments",
			args:     []types.Expr{&types.NumberExpr{Value: 12}, &types.NumberExpr{Value: 3}},
			expected: 4,
		},
		{
			name:     "multiple arguments",
			args:     []types.Expr{&types.NumberExpr{Value: 24}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 3}},
			expected: 4,
		},
		{
			name:     "decimal result",
			args:     []types.Expr{&types.NumberExpr{Value: 10}, &types.NumberExpr{Value: 4}},
			expected: 2.5,
		},
		{
			name:     "negative numbers",
			args:     []types.Expr{&types.NumberExpr{Value: -12}, &types.NumberExpr{Value: 3}},
			expected: -4,
		},
		{
			name:        "division by zero",
			args:        []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 0}},
			expectError: true,
		},
		{
			name:        "reciprocal of zero",
			args:        []types.Expr{&types.NumberExpr{Value: 0}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalDivide(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalDivide failed: %v", err)
			}

			numResult, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(numResult), tt.expected, 1e-10) {
				t.Errorf("Expected %f, got %f", tt.expected, float64(numResult))
			}
		})
	}
}

func TestArithmeticPlugin_Modulo(t *testing.T) {
	plugin := NewArithmeticPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    float64
		expectError bool
	}{
		{
			name:        "wrong number of arguments (0)",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (1)",
			args:        []types.Expr{&types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (3)",
			args:        []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 2}},
			expectError: true,
		},
		{
			name:     "basic modulo",
			args:     []types.Expr{&types.NumberExpr{Value: 10}, &types.NumberExpr{Value: 3}},
			expected: 1,
		},
		{
			name:     "exact division",
			args:     []types.Expr{&types.NumberExpr{Value: 15}, &types.NumberExpr{Value: 5}},
			expected: 0,
		},
		{
			name:     "negative dividend",
			args:     []types.Expr{&types.NumberExpr{Value: -10}, &types.NumberExpr{Value: 3}},
			expected: -1,
		},
		{
			name:     "negative divisor",
			args:     []types.Expr{&types.NumberExpr{Value: 10}, &types.NumberExpr{Value: -3}},
			expected: 1,
		},
		{
			name:     "decimal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 5.5}, &types.NumberExpr{Value: 2.0}},
			expected: 1.5,
		},
		{
			name:        "modulo by zero",
			args:        []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 0}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalModulo(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalModulo failed: %v", err)
			}

			numResult, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if !floatEqual(float64(numResult), tt.expected, 1e-10) {
				t.Errorf("Expected %f, got %f", tt.expected, float64(numResult))
			}
		})
	}
}

func TestArithmeticPlugin_ErrorHandling(t *testing.T) {
	plugin := NewArithmeticPlugin()
	evaluator := newMockEvaluator()

	// Test with invalid argument types
	stringArg := &types.StringExpr{Value: "not a number"}

	t.Run("addition with string", func(t *testing.T) {
		_, err := plugin.evalAdd(evaluator, []types.Expr{stringArg})
		if err == nil {
			t.Error("Expected error when adding string, but got none")
		}
	})

	t.Run("subtraction with string", func(t *testing.T) {
		_, err := plugin.evalSubtract(evaluator, []types.Expr{stringArg})
		if err == nil {
			t.Error("Expected error when subtracting string, but got none")
		}
	})

	t.Run("multiplication with string", func(t *testing.T) {
		_, err := plugin.evalMultiply(evaluator, []types.Expr{stringArg})
		if err == nil {
			t.Error("Expected error when multiplying string, but got none")
		}
	})

	t.Run("division with string", func(t *testing.T) {
		_, err := plugin.evalDivide(evaluator, []types.Expr{stringArg})
		if err == nil {
			t.Error("Expected error when dividing string, but got none")
		}
	})

	t.Run("modulo with string", func(t *testing.T) {
		_, err := plugin.evalModulo(evaluator, []types.Expr{stringArg, &types.NumberExpr{Value: 2}})
		if err == nil {
			t.Error("Expected error when modulo with string, but got none")
		}
	})
}

func TestArithmeticPlugin_PluginInfo(t *testing.T) {
	plugin := NewArithmeticPlugin()

	if plugin.Name() != "arithmetic" {
		t.Errorf("Expected plugin name 'arithmetic', got %s", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected plugin version '1.0.0', got %s", plugin.Version())
	}

	if plugin.Description() != "Basic arithmetic operations (+, -, *, /, %)" {
		t.Errorf("Expected specific description, got %s", plugin.Description())
	}

	deps := plugin.Dependencies()
	if len(deps) != 0 {
		t.Errorf("Expected no dependencies, got %v", deps)
	}
}
