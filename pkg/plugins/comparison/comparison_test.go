package comparison

import (
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

func TestComparisonPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewComparisonPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"=", "<", ">", "<=", ">="}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestComparisonPlugin_Equality(t *testing.T) {
	plugin := NewComparisonPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    bool
		expectError bool
	}{
		{
			name:        "too few arguments",
			args:        []types.Expr{&types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name:     "two equal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 5}},
			expected: true,
		},
		{
			name:     "two different numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 3}},
			expected: false,
		},
		{
			name:     "three equal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 2}},
			expected: true,
		},
		{
			name:     "three numbers with one different",
			args:     []types.Expr{&types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 3}},
			expected: false,
		},
		{
			name:     "equal strings",
			args:     []types.Expr{&types.StringExpr{Value: "hello"}, &types.StringExpr{Value: "hello"}},
			expected: true,
		},
		{
			name:     "different strings",
			args:     []types.Expr{&types.StringExpr{Value: "hello"}, &types.StringExpr{Value: "world"}},
			expected: false,
		},
		{
			name:     "equal booleans",
			args:     []types.Expr{&types.BooleanExpr{Value: true}, &types.BooleanExpr{Value: true}},
			expected: true,
		},
		{
			name:     "different booleans",
			args:     []types.Expr{&types.BooleanExpr{Value: true}, &types.BooleanExpr{Value: false}},
			expected: false,
		},
		{
			name:     "different types (number and string)",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.StringExpr{Value: "5"}},
			expected: false,
		},
		{
			name:     "decimal numbers equal",
			args:     []types.Expr{&types.NumberExpr{Value: 3.14}, &types.NumberExpr{Value: 3.14}},
			expected: true,
		},
		{
			name:     "decimal numbers different",
			args:     []types.Expr{&types.NumberExpr{Value: 3.14}, &types.NumberExpr{Value: 3.15}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalEquality(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalEquality failed: %v", err)
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestComparisonPlugin_EqualityWithBigNumbers(t *testing.T) {
	plugin := NewComparisonPlugin()
	evaluator := newMockEvaluator()

	// Test equality with big numbers
	bigNum1 := types.NewBigNumberFromInt64(1000000000000000000)
	bigNum2 := types.NewBigNumberFromInt64(1000000000000000000)
	bigNum3 := types.NewBigNumberFromInt64(1000000000000000001)

	tests := []struct {
		name     string
		args     []types.Expr
		expected bool
	}{
		{
			name:     "equal big numbers",
			args:     []types.Expr{wrapValue(bigNum1), wrapValue(bigNum2)},
			expected: true,
		},
		{
			name:     "different big numbers",
			args:     []types.Expr{wrapValue(bigNum1), wrapValue(bigNum3)},
			expected: false,
		},
		{
			name:     "big number equal to regular number",
			args:     []types.Expr{wrapValue(types.NewBigNumberFromInt64(5)), wrapValue(types.NumberValue(5))},
			expected: true,
		},
		{
			name:     "big number not equal to regular number",
			args:     []types.Expr{wrapValue(types.NewBigNumberFromInt64(5)), wrapValue(types.NumberValue(6))},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalEquality(evaluator, tt.args)
			if err != nil {
				t.Fatalf("evalEquality failed: %v", err)
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestComparisonPlugin_LessThan(t *testing.T) {
	plugin := NewComparisonPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    bool
		expectError bool
	}{
		{
			name:        "too few arguments",
			args:        []types.Expr{&types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name:     "two numbers - true",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 5}},
			expected: true,
		},
		{
			name:     "two numbers - false",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 3}},
			expected: false,
		},
		{
			name:     "equal numbers - false",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 5}},
			expected: false,
		},
		{
			name:     "three numbers ascending - true",
			args:     []types.Expr{&types.NumberExpr{Value: 1}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 3}},
			expected: true,
		},
		{
			name:     "three numbers not ascending - false",
			args:     []types.Expr{&types.NumberExpr{Value: 1}, &types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 2}},
			expected: false,
		},
		{
			name:     "negative numbers",
			args:     []types.Expr{&types.NumberExpr{Value: -5}, &types.NumberExpr{Value: -3}},
			expected: true,
		},
		{
			name:     "decimal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 1.5}, &types.NumberExpr{Value: 2.5}},
			expected: true,
		},
		{
			name:     "non-numeric values - false",
			args:     []types.Expr{&types.StringExpr{Value: "a"}, &types.StringExpr{Value: "b"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalLessThan(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalLessThan failed: %v", err)
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestComparisonPlugin_GreaterThan(t *testing.T) {
	plugin := NewComparisonPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    bool
		expectError bool
	}{
		{
			name:        "too few arguments",
			args:        []types.Expr{&types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name:     "two numbers - true",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 3}},
			expected: true,
		},
		{
			name:     "two numbers - false",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 5}},
			expected: false,
		},
		{
			name:     "equal numbers - false",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 5}},
			expected: false,
		},
		{
			name:     "three numbers descending - true",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 1}},
			expected: true,
		},
		{
			name:     "three numbers not descending - false",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 1}, &types.NumberExpr{Value: 2}},
			expected: false,
		},
		{
			name:     "negative numbers",
			args:     []types.Expr{&types.NumberExpr{Value: -3}, &types.NumberExpr{Value: -5}},
			expected: true,
		},
		{
			name:     "decimal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 2.5}, &types.NumberExpr{Value: 1.5}},
			expected: true,
		},
		{
			name:     "non-numeric values - false",
			args:     []types.Expr{&types.StringExpr{Value: "b"}, &types.StringExpr{Value: "a"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalGreaterThan(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalGreaterThan failed: %v", err)
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestComparisonPlugin_LessThanOrEqual(t *testing.T) {
	plugin := NewComparisonPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    bool
		expectError bool
	}{
		{
			name:        "too few arguments",
			args:        []types.Expr{&types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name:     "two numbers - less than",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 5}},
			expected: true,
		},
		{
			name:     "two numbers - equal",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 5}},
			expected: true,
		},
		{
			name:     "two numbers - greater than",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 3}},
			expected: false,
		},
		{
			name:     "three numbers ascending",
			args:     []types.Expr{&types.NumberExpr{Value: 1}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 3}},
			expected: true,
		},
		{
			name:     "three numbers with equal values",
			args:     []types.Expr{&types.NumberExpr{Value: 1}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 2}},
			expected: true,
		},
		{
			name:     "three numbers with descending part",
			args:     []types.Expr{&types.NumberExpr{Value: 1}, &types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 2}},
			expected: false,
		},
		{
			name:     "negative numbers",
			args:     []types.Expr{&types.NumberExpr{Value: -5}, &types.NumberExpr{Value: -3}},
			expected: true,
		},
		{
			name:     "decimal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 1.5}, &types.NumberExpr{Value: 2.5}},
			expected: true,
		},
		{
			name:     "non-numeric values - true (since neither is greater than the other)",
			args:     []types.Expr{&types.StringExpr{Value: "a"}, &types.StringExpr{Value: "b"}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalLessThanOrEqual(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalLessThanOrEqual failed: %v", err)
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestComparisonPlugin_GreaterThanOrEqual(t *testing.T) {
	plugin := NewComparisonPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    bool
		expectError bool
	}{
		{
			name:        "too few arguments",
			args:        []types.Expr{&types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name:     "two numbers - greater than",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 3}},
			expected: true,
		},
		{
			name:     "two numbers - equal",
			args:     []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 5}},
			expected: true,
		},
		{
			name:     "two numbers - less than",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 5}},
			expected: false,
		},
		{
			name:     "three numbers descending",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 1}},
			expected: true,
		},
		{
			name:     "three numbers with equal values",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 2}},
			expected: true,
		},
		{
			name:     "three numbers with ascending part",
			args:     []types.Expr{&types.NumberExpr{Value: 3}, &types.NumberExpr{Value: 1}, &types.NumberExpr{Value: 2}},
			expected: false,
		},
		{
			name:     "negative numbers",
			args:     []types.Expr{&types.NumberExpr{Value: -3}, &types.NumberExpr{Value: -5}},
			expected: true,
		},
		{
			name:     "decimal numbers",
			args:     []types.Expr{&types.NumberExpr{Value: 2.5}, &types.NumberExpr{Value: 1.5}},
			expected: true,
		},
		{
			name:     "non-numeric values - true (since neither is less than the other)",
			args:     []types.Expr{&types.StringExpr{Value: "b"}, &types.StringExpr{Value: "a"}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalGreaterThanOrEqual(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalGreaterThanOrEqual failed: %v", err)
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestComparisonPlugin_ValueComparison(t *testing.T) {
	plugin := NewComparisonPlugin()

	// Test valuesEqual method
	t.Run("valuesEqual", func(t *testing.T) {
		tests := []struct {
			name     string
			a, b     types.Value
			expected bool
		}{
			{
				name:     "equal numbers",
				a:        types.NumberValue(5),
				b:        types.NumberValue(5),
				expected: true,
			},
			{
				name:     "different numbers",
				a:        types.NumberValue(5),
				b:        types.NumberValue(3),
				expected: false,
			},
			{
				name:     "equal strings",
				a:        types.StringValue("hello"),
				b:        types.StringValue("hello"),
				expected: true,
			},
			{
				name:     "different strings",
				a:        types.StringValue("hello"),
				b:        types.StringValue("world"),
				expected: false,
			},
			{
				name:     "equal booleans",
				a:        types.BooleanValue(true),
				b:        types.BooleanValue(true),
				expected: true,
			},
			{
				name:     "different booleans",
				a:        types.BooleanValue(true),
				b:        types.BooleanValue(false),
				expected: false,
			},
			{
				name:     "equal keywords",
				a:        types.KeywordValue("test"),
				b:        types.KeywordValue("test"),
				expected: true,
			},
			{
				name:     "different keywords",
				a:        types.KeywordValue("test"),
				b:        types.KeywordValue("other"),
				expected: false,
			},
			{
				name:     "different types",
				a:        types.NumberValue(5),
				b:        types.StringValue("5"),
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := plugin.valuesEqual(tt.a, tt.b)
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			})
		}
	})

	// Test valueLessThan method
	t.Run("valueLessThan", func(t *testing.T) {
		tests := []struct {
			name     string
			a, b     types.Value
			expected bool
		}{
			{
				name:     "less than",
				a:        types.NumberValue(3),
				b:        types.NumberValue(5),
				expected: true,
			},
			{
				name:     "greater than",
				a:        types.NumberValue(5),
				b:        types.NumberValue(3),
				expected: false,
			},
			{
				name:     "equal",
				a:        types.NumberValue(5),
				b:        types.NumberValue(5),
				expected: false,
			},
			{
				name:     "non-numeric values",
				a:        types.StringValue("a"),
				b:        types.StringValue("b"),
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := plugin.valueLessThan(tt.a, tt.b)
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			})
		}
	})

	// Test valueGreaterThan method
	t.Run("valueGreaterThan", func(t *testing.T) {
		tests := []struct {
			name     string
			a, b     types.Value
			expected bool
		}{
			{
				name:     "greater than",
				a:        types.NumberValue(5),
				b:        types.NumberValue(3),
				expected: true,
			},
			{
				name:     "less than",
				a:        types.NumberValue(3),
				b:        types.NumberValue(5),
				expected: false,
			},
			{
				name:     "equal",
				a:        types.NumberValue(5),
				b:        types.NumberValue(5),
				expected: false,
			},
			{
				name:     "non-numeric values",
				a:        types.StringValue("b"),
				b:        types.StringValue("a"),
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := plugin.valueGreaterThan(tt.a, tt.b)
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			})
		}
	})
}

func TestComparisonPlugin_PluginInfo(t *testing.T) {
	plugin := NewComparisonPlugin()

	if plugin.Name() != "comparison" {
		t.Errorf("Expected plugin name 'comparison', got %s", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected plugin version '1.0.0', got %s", plugin.Version())
	}

	if plugin.Description() != "Comparison operations (=, <, >, <=, >=)" {
		t.Errorf("Expected specific description, got %s", plugin.Description())
	}

	deps := plugin.Dependencies()
	if len(deps) != 0 {
		t.Errorf("Expected no dependencies, got %v", deps)
	}
}
