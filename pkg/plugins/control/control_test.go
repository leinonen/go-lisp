package control

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
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
	case *types.ListExpr:
		values := make([]types.Value, len(e.Elements))
		for i, elem := range e.Elements {
			val, err := me.Eval(elem)
			if err != nil {
				return nil, err
			}
			values[i] = val
		}
		return &types.ListValue{Elements: values}, nil
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

func TestControlPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewControlPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"if", "do"}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestControlPlugin_evalIf(t *testing.T) {
	plugin := NewControlPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    types.Value
		expectError bool
	}{
		{
			name: "true condition with then clause",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
				&types.NumberExpr{Value: 10},
			},
			expected: types.NumberValue(10),
		},
		{
			name: "false condition with then clause - no else",
			args: []types.Expr{
				&types.BooleanExpr{Value: false},
				&types.NumberExpr{Value: 10},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "true condition with then and else clause",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
				&types.StringExpr{Value: "then"},
				&types.StringExpr{Value: "else"},
			},
			expected: types.StringValue("then"),
		},
		{
			name: "false condition with then and else clause",
			args: []types.Expr{
				&types.BooleanExpr{Value: false},
				&types.StringExpr{Value: "then"},
				&types.StringExpr{Value: "else"},
			},
			expected: types.StringValue("else"),
		},
		{
			name: "truthy number condition",
			args: []types.Expr{
				&types.NumberExpr{Value: 5},
				&types.StringExpr{Value: "positive"},
				&types.StringExpr{Value: "not positive"},
			},
			expected: types.StringValue("positive"),
		},
		{
			name: "falsy zero condition",
			args: []types.Expr{
				&types.NumberExpr{Value: 0},
				&types.StringExpr{Value: "positive"},
				&types.StringExpr{Value: "zero"},
			},
			expected: types.StringValue("zero"),
		},
		{
			name: "truthy string condition",
			args: []types.Expr{
				&types.StringExpr{Value: "hello"},
				&types.NumberExpr{Value: 1},
				&types.NumberExpr{Value: 2},
			},
			expected: types.NumberValue(1),
		},
		{
			name: "falsy empty string condition",
			args: []types.Expr{
				&types.StringExpr{Value: ""},
				&types.NumberExpr{Value: 1},
				&types.NumberExpr{Value: 2},
			},
			expected: types.NumberValue(2),
		},
		{
			name: "too few arguments",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
			},
			expectError: true,
		},
		{
			name: "too many arguments",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
				&types.NumberExpr{Value: 1},
				&types.NumberExpr{Value: 2},
				&types.NumberExpr{Value: 3},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalIf(evaluator, tt.args)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalIf failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestControlPlugin_evalDo(t *testing.T) {
	plugin := NewControlPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    types.Value
		expectError bool
	}{
		{
			name: "single expression",
			args: []types.Expr{
				&types.NumberExpr{Value: 42},
			},
			expected: types.NumberValue(42),
		},
		{
			name: "multiple expressions - returns last",
			args: []types.Expr{
				&types.NumberExpr{Value: 1},
				&types.StringExpr{Value: "hello"},
				&types.BooleanExpr{Value: true},
				&types.NumberExpr{Value: 99},
			},
			expected: types.NumberValue(99),
		},
		{
			name: "mixed types",
			args: []types.Expr{
				&types.StringExpr{Value: "first"},
				&types.BooleanExpr{Value: false},
				&types.StringExpr{Value: "last"},
			},
			expected: types.StringValue("last"),
		},
		{
			name:        "no arguments",
			args:        []types.Expr{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalDo(evaluator, tt.args)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalDo failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestControlPlugin_IfWithComplexConditions(t *testing.T) {
	plugin := NewControlPlugin()
	evaluator := newMockEvaluator()

	// Test with non-empty list (should be truthy)
	listValue := &types.ListValue{Elements: []types.Value{types.NumberValue(1)}}
	args := []types.Expr{
		wrapValue(listValue),
		&types.StringExpr{Value: "list not empty"},
		&types.StringExpr{Value: "list empty"},
	}

	result, err := plugin.evalIf(evaluator, args)
	if err != nil {
		t.Fatalf("evalIf failed: %v", err)
	}

	expected := types.StringValue("list not empty")
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test with empty list (should be falsy)
	emptyListValue := &types.ListValue{Elements: []types.Value{}}
	args = []types.Expr{
		wrapValue(emptyListValue),
		&types.StringExpr{Value: "list not empty"},
		&types.StringExpr{Value: "list empty"},
	}

	result, err = plugin.evalIf(evaluator, args)
	if err != nil {
		t.Fatalf("evalIf failed: %v", err)
	}

	expected = types.StringValue("list empty")
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestControlPlugin_DoSequentialEvaluation(t *testing.T) {
	plugin := NewControlPlugin()
	evaluator := newMockEvaluator()

	// Create a more complex test where we simulate side effects
	// by having expressions that would typically modify state
	args := []types.Expr{
		&types.StringExpr{Value: "step1"},
		&types.StringExpr{Value: "step2"},
		&types.StringExpr{Value: "step3"},
		&types.NumberExpr{Value: 42}, // This should be the final result
	}

	result, err := plugin.evalDo(evaluator, args)
	if err != nil {
		t.Fatalf("evalDo failed: %v", err)
	}

	expected := types.NumberValue(42)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
