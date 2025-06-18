package logical

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
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
	case *types.ListExpr:
		elements := make([]types.Value, len(e.Elements))
		for i, elem := range e.Elements {
			val, err := me.Eval(elem)
			if err != nil {
				return nil, err
			}
			elements[i] = val
		}
		return &types.ListValue{Elements: elements}, nil
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
	return valueExpr{value: value}
}

type valueExpr struct {
	value types.Value
}

func (ve valueExpr) String() string {
	return ve.value.String()
}

func TestLogicalPlugin_NewLogicalPlugin(t *testing.T) {
	plugin := NewLogicalPlugin()
	if plugin == nil {
		t.Fatal("NewLogicalPlugin returned nil")
	}

	if plugin.Name() != "logical" {
		t.Errorf("Expected plugin name 'logical', got %s", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", plugin.Version())
	}

	if plugin.Description() != "Logical operations (and, or, not)" {
		t.Errorf("Expected description 'Logical operations (and, or, not)', got %s", plugin.Description())
	}
}

func TestLogicalPlugin_And(t *testing.T) {
	plugin := NewLogicalPlugin()
	eval := newMockEvaluator()

	tests := []struct {
		name     string
		args     []types.Expr
		expected bool
		hasError bool
	}{
		{
			name:     "empty and",
			args:     []types.Expr{},
			expected: true,
		},
		{
			name: "all true",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
				&types.BooleanExpr{Value: true},
			},
			expected: true,
		},
		{
			name: "all false",
			args: []types.Expr{
				&types.BooleanExpr{Value: false},
				&types.BooleanExpr{Value: false},
			},
			expected: false,
		},
		{
			name: "mixed with false",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
				&types.BooleanExpr{Value: false},
			},
			expected: false,
		},
		{
			name: "single true",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
			},
			expected: true,
		},
		{
			name: "truthy numbers",
			args: []types.Expr{
				&types.NumberExpr{Value: 1},
				&types.NumberExpr{Value: 2.5},
			},
			expected: true,
		},
		{
			name: "falsy number",
			args: []types.Expr{
				&types.NumberExpr{Value: 1},
				&types.NumberExpr{Value: 0},
			},
			expected: false,
		},
		{
			name: "truthy strings",
			args: []types.Expr{
				&types.StringExpr{Value: "hello"},
				&types.StringExpr{Value: "world"},
			},
			expected: true,
		},
		{
			name: "falsy string",
			args: []types.Expr{
				&types.StringExpr{Value: "hello"},
				&types.StringExpr{Value: ""},
			},
			expected: false,
		},
		{
			name: "truthy lists",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{&types.NumberExpr{Value: 1}}},
				&types.ListExpr{Elements: []types.Expr{&types.StringExpr{Value: "test"}}},
			},
			expected: true,
		},
		{
			name: "falsy list",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{&types.NumberExpr{Value: 1}}},
				&types.ListExpr{Elements: []types.Expr{}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalAnd(eval, tt.args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Errorf("Expected BooleanValue, got %T", result)
				return
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestLogicalPlugin_Or(t *testing.T) {
	plugin := NewLogicalPlugin()
	eval := newMockEvaluator()

	tests := []struct {
		name     string
		args     []types.Expr
		expected bool
		hasError bool
	}{
		{
			name:     "empty or",
			args:     []types.Expr{},
			expected: false,
		},
		{
			name: "all true",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
				&types.BooleanExpr{Value: true},
			},
			expected: true,
		},
		{
			name: "all false",
			args: []types.Expr{
				&types.BooleanExpr{Value: false},
				&types.BooleanExpr{Value: false},
			},
			expected: false,
		},
		{
			name: "mixed with true",
			args: []types.Expr{
				&types.BooleanExpr{Value: false},
				&types.BooleanExpr{Value: true},
			},
			expected: true,
		},
		{
			name: "single true",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
			},
			expected: true,
		},
		{
			name: "truthy number with falsy",
			args: []types.Expr{
				&types.NumberExpr{Value: 0},
				&types.NumberExpr{Value: 1},
			},
			expected: true,
		},
		{
			name: "all falsy numbers",
			args: []types.Expr{
				&types.NumberExpr{Value: 0},
				&types.NumberExpr{Value: 0},
			},
			expected: false,
		},
		{
			name: "truthy string with falsy",
			args: []types.Expr{
				&types.StringExpr{Value: ""},
				&types.StringExpr{Value: "hello"},
			},
			expected: true,
		},
		{
			name: "all falsy strings",
			args: []types.Expr{
				&types.StringExpr{Value: ""},
				&types.StringExpr{Value: ""},
			},
			expected: false,
		},
		{
			name: "truthy list with falsy",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{}},
				&types.ListExpr{Elements: []types.Expr{&types.NumberExpr{Value: 1}}},
			},
			expected: true,
		},
		{
			name: "all falsy lists",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{}},
				&types.ListExpr{Elements: []types.Expr{}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalOr(eval, tt.args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Errorf("Expected BooleanValue, got %T", result)
				return
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestLogicalPlugin_Not(t *testing.T) {
	plugin := NewLogicalPlugin()
	eval := newMockEvaluator()

	tests := []struct {
		name     string
		args     []types.Expr
		expected bool
		hasError bool
	}{
		{
			name: "not true",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
			},
			expected: false,
		},
		{
			name: "not false",
			args: []types.Expr{
				&types.BooleanExpr{Value: false},
			},
			expected: true,
		},
		{
			name: "not truthy number",
			args: []types.Expr{
				&types.NumberExpr{Value: 1},
			},
			expected: false,
		},
		{
			name: "not falsy number",
			args: []types.Expr{
				&types.NumberExpr{Value: 0},
			},
			expected: true,
		},
		{
			name: "not truthy string",
			args: []types.Expr{
				&types.StringExpr{Value: "hello"},
			},
			expected: false,
		},
		{
			name: "not falsy string",
			args: []types.Expr{
				&types.StringExpr{Value: ""},
			},
			expected: true,
		},
		{
			name: "not truthy list",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{&types.NumberExpr{Value: 1}}},
			},
			expected: false,
		},
		{
			name: "not falsy list",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{}},
			},
			expected: true,
		},
		{
			name:     "no arguments",
			args:     []types.Expr{},
			hasError: true,
		},
		{
			name: "too many arguments",
			args: []types.Expr{
				&types.BooleanExpr{Value: true},
				&types.BooleanExpr{Value: false},
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalNot(eval, tt.args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Errorf("Expected BooleanValue, got %T", result)
				return
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}
