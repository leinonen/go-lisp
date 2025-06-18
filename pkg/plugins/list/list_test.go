package list

import (
	"reflect"
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

// valuesEqual compares two values for equality
func valuesEqual(a, b types.Value) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	switch va := a.(type) {
	case types.NumberValue:
		vb := b.(types.NumberValue)
		return va == vb
	case types.StringValue:
		vb := b.(types.StringValue)
		return va == vb
	case types.BooleanValue:
		vb := b.(types.BooleanValue)
		return va == vb
	case *types.ListValue:
		vb := b.(*types.ListValue)
		if len(va.Elements) != len(vb.Elements) {
			return false
		}
		for i, elem := range va.Elements {
			if !valuesEqual(elem, vb.Elements[i]) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a, b)
	}
}

func TestListPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewListPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{
		"list", "first", "rest", "cons", "length", "empty?",
		"append", "reverse", "nth", "last",
	}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestListPlugin_evalList(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		args     []types.Expr
		expected []types.Value
	}{
		{
			name:     "empty list",
			args:     []types.Expr{},
			expected: []types.Value{},
		},
		{
			name: "single element",
			args: []types.Expr{
				&types.NumberExpr{Value: 1},
			},
			expected: []types.Value{types.NumberValue(1)},
		},
		{
			name: "multiple elements",
			args: []types.Expr{
				&types.NumberExpr{Value: 1},
				&types.NumberExpr{Value: 2},
				&types.NumberExpr{Value: 3},
			},
			expected: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			},
		},
		{
			name: "mixed types",
			args: []types.Expr{
				&types.NumberExpr{Value: 1},
				&types.StringExpr{Value: "hello"},
				&types.BooleanExpr{Value: true},
			},
			expected: []types.Value{
				types.NumberValue(1),
				types.StringValue("hello"),
				types.BooleanValue(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalList(evaluator, tt.args)
			if err != nil {
				t.Fatalf("evalList failed: %v", err)
			}

			list, ok := result.(*types.ListValue)
			if !ok {
				t.Fatalf("Expected ListValue, got %T", result)
			}

			if len(list.Elements) != len(tt.expected) {
				t.Fatalf("Expected %d elements, got %d", len(tt.expected), len(list.Elements))
			}

			for i, expected := range tt.expected {
				if !valuesEqual(list.Elements[i], expected) {
					t.Errorf("Element %d: expected %v, got %v", i, expected, list.Elements[i])
				}
			}
		})
	}
}

func TestListPlugin_evalFirst(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		list        *types.ListValue
		expected    types.Value
		expectError bool
	}{
		{
			name: "non-empty list",
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			}},
			expected: types.NumberValue(1),
		},
		{
			name: "single element list",
			list: &types.ListValue{Elements: []types.Value{
				types.StringValue("hello"),
			}},
			expected: types.StringValue("hello"),
		},
		{
			name:        "empty list",
			list:        &types.ListValue{Elements: []types.Value{}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{wrapValue(tt.list)}
			result, err := plugin.evalFirst(evaluator, args)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalFirst failed: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestListPlugin_evalRest(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		list     *types.ListValue
		expected []types.Value
	}{
		{
			name: "multiple elements",
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			}},
			expected: []types.Value{
				types.NumberValue(2),
				types.NumberValue(3),
			},
		},
		{
			name: "single element",
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
			}},
			expected: []types.Value{},
		},
		{
			name:     "empty list",
			list:     &types.ListValue{Elements: []types.Value{}},
			expected: []types.Value{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{wrapValue(tt.list)}
			result, err := plugin.evalRest(evaluator, args)
			if err != nil {
				t.Fatalf("evalRest failed: %v", err)
			}

			list, ok := result.(*types.ListValue)
			if !ok {
				t.Fatalf("Expected ListValue, got %T", result)
			}

			if len(list.Elements) != len(tt.expected) {
				t.Fatalf("Expected %d elements, got %d", len(tt.expected), len(list.Elements))
			}

			for i, expected := range tt.expected {
				if !valuesEqual(list.Elements[i], expected) {
					t.Errorf("Element %d: expected %v, got %v", i, expected, list.Elements[i])
				}
			}
		})
	}
}

func TestListPlugin_evalCons(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		element  types.Value
		list     *types.ListValue
		expected []types.Value
	}{
		{
			name:    "prepend to non-empty list",
			element: types.NumberValue(0),
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
			}},
			expected: []types.Value{
				types.NumberValue(0),
				types.NumberValue(1),
				types.NumberValue(2),
			},
		},
		{
			name:     "prepend to empty list",
			element:  types.StringValue("hello"),
			list:     &types.ListValue{Elements: []types.Value{}},
			expected: []types.Value{types.StringValue("hello")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{
				wrapValue(tt.element),
				wrapValue(tt.list),
			}
			result, err := plugin.evalCons(evaluator, args)
			if err != nil {
				t.Fatalf("evalCons failed: %v", err)
			}

			list, ok := result.(*types.ListValue)
			if !ok {
				t.Fatalf("Expected ListValue, got %T", result)
			}

			if len(list.Elements) != len(tt.expected) {
				t.Fatalf("Expected %d elements, got %d", len(tt.expected), len(list.Elements))
			}

			for i, expected := range tt.expected {
				if !valuesEqual(list.Elements[i], expected) {
					t.Errorf("Element %d: expected %v, got %v", i, expected, list.Elements[i])
				}
			}
		})
	}
}

func TestListPlugin_evalLength(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		list     *types.ListValue
		expected int
	}{
		{
			name:     "empty list",
			list:     &types.ListValue{Elements: []types.Value{}},
			expected: 0,
		},
		{
			name: "single element",
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
			}},
			expected: 1,
		},
		{
			name: "multiple elements",
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			}},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{wrapValue(tt.list)}
			result, err := plugin.evalLength(evaluator, args)
			if err != nil {
				t.Fatalf("evalLength failed: %v", err)
			}

			length, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if int(length) != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, int(length))
			}
		})
	}
}

func TestListPlugin_evalEmpty(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		list     *types.ListValue
		expected bool
	}{
		{
			name:     "empty list",
			list:     &types.ListValue{Elements: []types.Value{}},
			expected: true,
		},
		{
			name: "non-empty list",
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
			}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{wrapValue(tt.list)}
			result, err := plugin.evalEmpty(evaluator, args)
			if err != nil {
				t.Fatalf("evalEmpty failed: %v", err)
			}

			isEmpty, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}

			if bool(isEmpty) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(isEmpty))
			}
		})
	}
}

func TestListPlugin_evalAppend(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		lists    []*types.ListValue
		expected []types.Value
	}{
		{
			name:     "no lists",
			lists:    []*types.ListValue{},
			expected: []types.Value{},
		},
		{
			name: "single list",
			lists: []*types.ListValue{
				{Elements: []types.Value{types.NumberValue(1), types.NumberValue(2)}},
			},
			expected: []types.Value{types.NumberValue(1), types.NumberValue(2)},
		},
		{
			name: "multiple lists",
			lists: []*types.ListValue{
				{Elements: []types.Value{types.NumberValue(1), types.NumberValue(2)}},
				{Elements: []types.Value{types.NumberValue(3), types.NumberValue(4)}},
				{Elements: []types.Value{types.NumberValue(5)}},
			},
			expected: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
				types.NumberValue(4),
				types.NumberValue(5),
			},
		},
		{
			name: "with empty lists",
			lists: []*types.ListValue{
				{Elements: []types.Value{types.NumberValue(1)}},
				{Elements: []types.Value{}},
				{Elements: []types.Value{types.NumberValue(2)}},
			},
			expected: []types.Value{types.NumberValue(1), types.NumberValue(2)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]types.Expr, len(tt.lists))
			for i, list := range tt.lists {
				args[i] = wrapValue(list)
			}

			result, err := plugin.evalAppend(evaluator, args)
			if err != nil {
				t.Fatalf("evalAppend failed: %v", err)
			}

			list, ok := result.(*types.ListValue)
			if !ok {
				t.Fatalf("Expected ListValue, got %T", result)
			}

			if len(list.Elements) != len(tt.expected) {
				t.Fatalf("Expected %d elements, got %d", len(tt.expected), len(list.Elements))
			}

			for i, expected := range tt.expected {
				if !valuesEqual(list.Elements[i], expected) {
					t.Errorf("Element %d: expected %v, got %v", i, expected, list.Elements[i])
				}
			}
		})
	}
}

func TestListPlugin_evalReverse(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		list     *types.ListValue
		expected []types.Value
	}{
		{
			name:     "empty list",
			list:     &types.ListValue{Elements: []types.Value{}},
			expected: []types.Value{},
		},
		{
			name: "single element",
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
			}},
			expected: []types.Value{types.NumberValue(1)},
		},
		{
			name: "multiple elements",
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			}},
			expected: []types.Value{
				types.NumberValue(3),
				types.NumberValue(2),
				types.NumberValue(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{wrapValue(tt.list)}
			result, err := plugin.evalReverse(evaluator, args)
			if err != nil {
				t.Fatalf("evalReverse failed: %v", err)
			}

			list, ok := result.(*types.ListValue)
			if !ok {
				t.Fatalf("Expected ListValue, got %T", result)
			}

			if len(list.Elements) != len(tt.expected) {
				t.Fatalf("Expected %d elements, got %d", len(tt.expected), len(list.Elements))
			}

			for i, expected := range tt.expected {
				if !valuesEqual(list.Elements[i], expected) {
					t.Errorf("Element %d: expected %v, got %v", i, expected, list.Elements[i])
				}
			}
		})
	}
}

func TestListPlugin_evalNth(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	list := &types.ListValue{Elements: []types.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("c"),
	}}

	tests := []struct {
		name        string
		index       int
		expected    types.Value
		expectError bool
	}{
		{
			name:     "valid index 0",
			index:    0,
			expected: types.StringValue("a"),
		},
		{
			name:     "valid index 1",
			index:    1,
			expected: types.StringValue("b"),
		},
		{
			name:     "valid index 2",
			index:    2,
			expected: types.StringValue("c"),
		},
		{
			name:        "negative index",
			index:       -1,
			expectError: true,
		},
		{
			name:        "index out of bounds",
			index:       3,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{
				wrapValue(list),
				wrapValue(types.NumberValue(tt.index)),
			}
			result, err := plugin.evalNth(evaluator, args)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalNth failed: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestListPlugin_evalLast(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		list        *types.ListValue
		expected    types.Value
		expectError bool
	}{
		{
			name: "non-empty list",
			list: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			}},
			expected: types.NumberValue(3),
		},
		{
			name: "single element list",
			list: &types.ListValue{Elements: []types.Value{
				types.StringValue("hello"),
			}},
			expected: types.StringValue("hello"),
		},
		{
			name:        "empty list",
			list:        &types.ListValue{Elements: []types.Value{}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{wrapValue(tt.list)}
			result, err := plugin.evalLast(evaluator, args)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalLast failed: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestListPlugin_InvalidArguments(t *testing.T) {
	plugin := NewListPlugin()
	evaluator := newMockEvaluator()

	// Test with non-list arguments
	nonList := types.NumberValue(42)

	tests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		args     []types.Expr
	}{
		{
			name:     "first with non-list",
			function: plugin.evalFirst,
			args:     []types.Expr{wrapValue(nonList)},
		},
		{
			name:     "rest with non-list",
			function: plugin.evalRest,
			args:     []types.Expr{wrapValue(nonList)},
		},
		{
			name:     "length with non-list",
			function: plugin.evalLength,
			args:     []types.Expr{wrapValue(nonList)},
		},
		{
			name:     "empty? with non-list",
			function: plugin.evalEmpty,
			args:     []types.Expr{wrapValue(nonList)},
		},
		{
			name:     "reverse with non-list",
			function: plugin.evalReverse,
			args:     []types.Expr{wrapValue(nonList)},
		},
		{
			name:     "last with non-list",
			function: plugin.evalLast,
			args:     []types.Expr{wrapValue(nonList)},
		},
		{
			name:     "cons with non-list second arg",
			function: plugin.evalCons,
			args: []types.Expr{
				wrapValue(types.NumberValue(1)),
				wrapValue(nonList),
			},
		},
		{
			name:     "nth with non-number first arg",
			function: plugin.evalNth,
			args: []types.Expr{
				wrapValue(types.StringValue("not-a-number")),
				wrapValue(&types.ListValue{Elements: []types.Value{types.NumberValue(1)}}),
			},
		},
		{
			name:     "nth with non-list second arg",
			function: plugin.evalNth,
			args: []types.Expr{
				wrapValue(types.NumberValue(0)),
				wrapValue(nonList),
			},
		},
		{
			name:     "append with non-list",
			function: plugin.evalAppend,
			args:     []types.Expr{wrapValue(nonList)},
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
