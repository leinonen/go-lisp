package functional

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
	case *types.KeywordExpr:
		return types.KeywordValue(e.Value), nil
	case *types.SymbolExpr:
		// Try to resolve from environment
		if value, exists := me.env.Get(e.Name); exists {
			return value, nil
		}
		// Return the symbol name as a function placeholder for testing
		return &mockFunction{name: e.Name}, nil
	case *types.ListExpr:
		// Convert to list value
		var elements []types.Value
		for _, elem := range e.Elements {
			val, err := me.Eval(elem)
			if err != nil {
				return nil, err
			}
			elements = append(elements, val)
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
	// Mock function calls for testing
	if mockFn, ok := funcValue.(*mockFunction); ok {
		switch mockFn.name {
		case "double":
			// Double function: multiply by 2
			if len(args) != 1 {
				return nil, nil
			}
			val, err := me.Eval(args[0])
			if err != nil {
				return nil, err
			}
			if num, ok := val.(types.NumberValue); ok {
				return types.NumberValue(float64(num) * 2), nil
			}
		case "isEven":
			// isEven predicate: check if number is even
			if len(args) != 1 {
				return nil, nil
			}
			val, err := me.Eval(args[0])
			if err != nil {
				return nil, err
			}
			if num, ok := val.(types.NumberValue); ok {
				return types.BooleanValue(int(num)%2 == 0), nil
			}
		case "add":
			// Add function: add two numbers
			if len(args) != 2 {
				return nil, nil
			}
			val1, err := me.Eval(args[0])
			if err != nil {
				return nil, err
			}
			val2, err := me.Eval(args[1])
			if err != nil {
				return nil, err
			}
			if num1, ok1 := val1.(types.NumberValue); ok1 {
				if num2, ok2 := val2.(types.NumberValue); ok2 {
					return types.NumberValue(float64(num1) + float64(num2)), nil
				}
			}
		case "concat":
			// Concat function: concatenate two strings
			if len(args) != 2 {
				return nil, nil
			}
			val1, err := me.Eval(args[0])
			if err != nil {
				return nil, err
			}
			val2, err := me.Eval(args[1])
			if err != nil {
				return nil, err
			}
			if str1, ok1 := val1.(types.StringValue); ok1 {
				if str2, ok2 := val2.(types.StringValue); ok2 {
					return types.StringValue(string(str1) + string(str2)), nil
				}
			}
		}
	}
	return nil, nil
}

// Mock function type for testing
type mockFunction struct {
	name string
}

func (mf *mockFunction) String() string {
	return mf.name
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

func TestFunctionalPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewFunctionalPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"map", "filter", "reduce", "apply"}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestFunctionalPlugin_MapFunc(t *testing.T) {
	plugin := NewFunctionalPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    []float64
		expectError bool
	}{
		{
			name:        "wrong number of arguments (0)",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (1)",
			args:        []types.Expr{&types.SymbolExpr{Name: "double"}},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (3)",
			args:        []types.Expr{&types.SymbolExpr{Name: "double"}, wrapValue(&types.ListValue{Elements: []types.Value{}}), &types.NumberExpr{Value: 1}},
			expectError: true,
		},
		{
			name:        "second argument not a list",
			args:        []types.Expr{&types.SymbolExpr{Name: "double"}, &types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name: "map double over empty list",
			args: []types.Expr{
				&types.SymbolExpr{Name: "double"},
				wrapValue(&types.ListValue{Elements: []types.Value{}}),
			},
			expected: []float64{},
		},
		{
			name: "map double over single element",
			args: []types.Expr{
				&types.SymbolExpr{Name: "double"},
				wrapValue(&types.ListValue{Elements: []types.Value{types.NumberValue(5)}}),
			},
			expected: []float64{10},
		},
		{
			name: "map double over multiple elements",
			args: []types.Expr{
				&types.SymbolExpr{Name: "double"},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(1),
					types.NumberValue(2),
					types.NumberValue(3),
				}}),
			},
			expected: []float64{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.mapFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("mapFunc failed: %v", err)
			}

			listResult, ok := result.(*types.ListValue)
			if !ok {
				t.Fatalf("Expected ListValue, got %T", result)
			}

			if len(listResult.Elements) != len(tt.expected) {
				t.Errorf("Expected list length %d, got %d", len(tt.expected), len(listResult.Elements))
			}

			for i, expectedVal := range tt.expected {
				if i >= len(listResult.Elements) {
					t.Errorf("Missing element at index %d", i)
					continue
				}

				numVal, ok := listResult.Elements[i].(types.NumberValue)
				if !ok {
					t.Errorf("Expected NumberValue at index %d, got %T", i, listResult.Elements[i])
					continue
				}

				if float64(numVal) != expectedVal {
					t.Errorf("Expected %f at index %d, got %f", expectedVal, i, float64(numVal))
				}
			}
		})
	}
}

func TestFunctionalPlugin_FilterFunc(t *testing.T) {
	plugin := NewFunctionalPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    []float64
		expectError bool
	}{
		{
			name:        "wrong number of arguments (0)",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (1)",
			args:        []types.Expr{&types.SymbolExpr{Name: "isEven"}},
			expectError: true,
		},
		{
			name:        "second argument not a list",
			args:        []types.Expr{&types.SymbolExpr{Name: "isEven"}, &types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name: "filter isEven over empty list",
			args: []types.Expr{
				&types.SymbolExpr{Name: "isEven"},
				wrapValue(&types.ListValue{Elements: []types.Value{}}),
			},
			expected: []float64{},
		},
		{
			name: "filter isEven over list with only odd numbers",
			args: []types.Expr{
				&types.SymbolExpr{Name: "isEven"},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(1),
					types.NumberValue(3),
					types.NumberValue(5),
				}}),
			},
			expected: []float64{},
		},
		{
			name: "filter isEven over list with only even numbers",
			args: []types.Expr{
				&types.SymbolExpr{Name: "isEven"},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(2),
					types.NumberValue(4),
					types.NumberValue(6),
				}}),
			},
			expected: []float64{2, 4, 6},
		},
		{
			name: "filter isEven over mixed list",
			args: []types.Expr{
				&types.SymbolExpr{Name: "isEven"},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(1),
					types.NumberValue(2),
					types.NumberValue(3),
					types.NumberValue(4),
					types.NumberValue(5),
				}}),
			},
			expected: []float64{2, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.filterFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("filterFunc failed: %v", err)
			}

			listResult, ok := result.(*types.ListValue)
			if !ok {
				t.Fatalf("Expected ListValue, got %T", result)
			}

			if len(listResult.Elements) != len(tt.expected) {
				t.Errorf("Expected list length %d, got %d", len(tt.expected), len(listResult.Elements))
			}

			for i, expectedVal := range tt.expected {
				if i >= len(listResult.Elements) {
					t.Errorf("Missing element at index %d", i)
					continue
				}

				numVal, ok := listResult.Elements[i].(types.NumberValue)
				if !ok {
					t.Errorf("Expected NumberValue at index %d, got %T", i, listResult.Elements[i])
					continue
				}

				if float64(numVal) != expectedVal {
					t.Errorf("Expected %f at index %d, got %f", expectedVal, i, float64(numVal))
				}
			}
		})
	}
}

func TestFunctionalPlugin_ReduceFunc(t *testing.T) {
	plugin := NewFunctionalPlugin()
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
			name:        "wrong number of arguments (2)",
			args:        []types.Expr{&types.SymbolExpr{Name: "add"}, &types.NumberExpr{Value: 0}},
			expectError: true,
		},
		{
			name:        "third argument not a list",
			args:        []types.Expr{&types.SymbolExpr{Name: "add"}, &types.NumberExpr{Value: 0}, &types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name: "reduce add over empty list",
			args: []types.Expr{
				&types.SymbolExpr{Name: "add"},
				&types.NumberExpr{Value: 10},
				wrapValue(&types.ListValue{Elements: []types.Value{}}),
			},
			expected: 10, // Initial value should be returned
		},
		{
			name: "reduce add over single element",
			args: []types.Expr{
				&types.SymbolExpr{Name: "add"},
				&types.NumberExpr{Value: 0},
				wrapValue(&types.ListValue{Elements: []types.Value{types.NumberValue(5)}}),
			},
			expected: 5, // 0 + 5 = 5
		},
		{
			name: "reduce add over multiple elements",
			args: []types.Expr{
				&types.SymbolExpr{Name: "add"},
				&types.NumberExpr{Value: 0},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(1),
					types.NumberValue(2),
					types.NumberValue(3),
					types.NumberValue(4),
				}}),
			},
			expected: 10, // 0 + 1 + 2 + 3 + 4 = 10
		},
		{
			name: "reduce add with non-zero initial value",
			args: []types.Expr{
				&types.SymbolExpr{Name: "add"},
				&types.NumberExpr{Value: 100},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(1),
					types.NumberValue(2),
					types.NumberValue(3),
				}}),
			},
			expected: 106, // 100 + 1 + 2 + 3 = 106
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.reduceFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("reduceFunc failed: %v", err)
			}

			numResult, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("Expected NumberValue, got %T", result)
			}

			if float64(numResult) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(numResult))
			}
		})
	}
}

func TestFunctionalPlugin_ApplyFunc(t *testing.T) {
	plugin := NewFunctionalPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    interface{}
		expectError bool
	}{
		{
			name:        "wrong number of arguments (0)",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (1)",
			args:        []types.Expr{&types.SymbolExpr{Name: "add"}},
			expectError: true,
		},
		{
			name:        "second argument not a list",
			args:        []types.Expr{&types.SymbolExpr{Name: "add"}, &types.NumberExpr{Value: 5}},
			expectError: true,
		},
		{
			name: "apply add to two numbers",
			args: []types.Expr{
				&types.SymbolExpr{Name: "add"},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(3),
					types.NumberValue(7),
				}}),
			},
			expected: 10.0,
		},
		{
			name: "apply concat to two strings",
			args: []types.Expr{
				&types.SymbolExpr{Name: "concat"},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.StringValue("hello"),
					types.StringValue("world"),
				}}),
			},
			expected: "helloworld",
		},
		{
			name: "apply function to empty argument list",
			args: []types.Expr{
				&types.SymbolExpr{Name: "add"},
				wrapValue(&types.ListValue{Elements: []types.Value{}}),
			},
			expected: nil, // Function will return nil for empty args
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.applyFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("applyFunc failed: %v", err)
			}

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}

			switch expectedVal := tt.expected.(type) {
			case float64:
				numResult, ok := result.(types.NumberValue)
				if !ok {
					t.Fatalf("Expected NumberValue, got %T", result)
				}
				if float64(numResult) != expectedVal {
					t.Errorf("Expected %f, got %f", expectedVal, float64(numResult))
				}
			case string:
				strResult, ok := result.(types.StringValue)
				if !ok {
					t.Fatalf("Expected StringValue, got %T", result)
				}
				if string(strResult) != expectedVal {
					t.Errorf("Expected %s, got %s", expectedVal, string(strResult))
				}
			}
		})
	}
}

func TestFunctionalPlugin_ValueToExpr(t *testing.T) {
	plugin := NewFunctionalPlugin()

	tests := []struct {
		name     string
		value    types.Value
		expected string
	}{
		{
			name:     "number value",
			value:    types.NumberValue(42),
			expected: "NumberExpr(42)",
		},
		{
			name:     "string value",
			value:    types.StringValue("hello"),
			expected: "StringExpr(\"hello\")",
		},
		{
			name:     "boolean value",
			value:    types.BooleanValue(true),
			expected: "BooleanExpr(true)",
		},
		{
			name:     "keyword value",
			value:    types.KeywordValue("test"),
			expected: "KeywordExpr(:test)",
		},
		{
			name: "list value",
			value: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
			}},
			expected: "ListExpr([NumberExpr(1) NumberExpr(2)])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := plugin.valueToExpr(tt.value)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFunctionalPlugin_IsTruthy(t *testing.T) {
	plugin := NewFunctionalPlugin()

	tests := []struct {
		name     string
		value    types.Value
		expected bool
	}{
		{
			name:     "boolean true",
			value:    types.BooleanValue(true),
			expected: true,
		},
		{
			name:     "boolean false",
			value:    types.BooleanValue(false),
			expected: false,
		},
		{
			name:     "nil value",
			value:    &types.NilValue{},
			expected: false,
		},
		{
			name:     "non-zero number",
			value:    types.NumberValue(42),
			expected: true,
		},
		{
			name:     "zero number",
			value:    types.NumberValue(0),
			expected: false,
		},
		{
			name:     "non-empty string",
			value:    types.StringValue("hello"),
			expected: true,
		},
		{
			name:     "empty string",
			value:    types.StringValue(""),
			expected: false,
		},
		{
			name:     "non-empty list",
			value:    &types.ListValue{Elements: []types.Value{types.NumberValue(1)}},
			expected: true,
		},
		{
			name:     "empty list",
			value:    &types.ListValue{Elements: []types.Value{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := plugin.isTruthy(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFunctionalPlugin_PluginInfo(t *testing.T) {
	plugin := NewFunctionalPlugin()

	if plugin.Name() != "functional" {
		t.Errorf("Expected plugin name 'functional', got %s", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected plugin version '1.0.0', got %s", plugin.Version())
	}

	if plugin.Description() != "Higher-order functions (map, filter, reduce, etc.)" {
		t.Errorf("Expected specific description, got %s", plugin.Description())
	}

	deps := plugin.Dependencies()
	if len(deps) != 0 {
		t.Errorf("Expected no dependencies, got %v", deps)
	}
}
