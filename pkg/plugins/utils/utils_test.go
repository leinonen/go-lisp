package utils

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

func TestUtilsPlugin_RegisterFunctions(t *testing.T) {
	mockEval := &mockEvaluator{env: evaluator.NewEnvironment()}
	plugin := NewUtilsPlugin(mockEval)
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{
		"frequencies", "group-by", "partition", "interleave", "interpose",
		"flatten", "shuffle", "remove", "keep", "mapcat",
		"take-while", "drop-while", "split-at", "split-with",
		"comp", "partial", "complement", "juxt",
		"union", "intersection", "difference",
	}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestUtilsPlugin_PluginInfo(t *testing.T) {
	mockEval := &mockEvaluator{env: evaluator.NewEnvironment()}
	plugin := NewUtilsPlugin(mockEval)

	if plugin.Name() != "utils" {
		t.Errorf("Expected plugin name 'utils', got '%s'", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", plugin.Version())
	}
}

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
		// Return the symbol name as a string for testing
		return types.StringValue(e.Name), nil
	case valueExpr:
		return e.value, nil
	default:
		return types.StringValue("mock"), nil
	}
}

func (me *mockEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	// Simple mock: for testing purposes, simulate common functions
	if funcStr, ok := funcValue.(types.StringValue); ok {
		switch string(funcStr) {
		case "inc":
			if len(args) == 1 {
				if val, err := me.Eval(args[0]); err == nil {
					if numVal, ok := val.(types.NumberValue); ok {
						return types.NumberValue(numVal + 1), nil
					}
				}
			}
		case "even?":
			if len(args) == 1 {
				if val, err := me.Eval(args[0]); err == nil {
					if numVal, ok := val.(types.NumberValue); ok {
						return types.BooleanValue(int(numVal)%2 == 0), nil
					}
				}
			}
		}
	}
	return types.StringValue("function"), nil
}

func (me *mockEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	return me.Eval(expr)
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

func (ve valueExpr) GetPosition() types.Position {
	return types.Position{}
}

func TestUtilsPlugin_EvalFrequencies(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewUtilsPlugin(evaluator)

	tests := []struct {
		name        string
		args        []types.Expr
		expectError bool
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name: "frequencies of list",
			args: []types.Expr{
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.StringValue("a"),
					types.StringValue("b"),
					types.StringValue("a"),
					types.StringValue("c"),
					types.StringValue("a"),
				}}),
			},
			expectError: false,
		},
		{
			name:        "frequencies of non-collection",
			args:        []types.Expr{&types.NumberExpr{Value: 42}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalFrequencies(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalFrequencies failed: %v", err)
			}

			// Should return a hash map
			if _, ok := result.(*types.HashMapValue); !ok {
				t.Errorf("Expected HashMapValue, got %T", result)
			}
		})
	}
}

func TestUtilsPlugin_EvalPartition(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewUtilsPlugin(evaluator)

	tests := []struct {
		name        string
		args        []types.Expr
		expectError bool
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name: "partition list into chunks of 2",
			args: []types.Expr{
				&types.NumberExpr{Value: 2},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(1),
					types.NumberValue(2),
					types.NumberValue(3),
					types.NumberValue(4),
					types.NumberValue(5),
				}}),
			},
			expectError: false,
		},
		{
			name:        "invalid chunk size",
			args:        []types.Expr{&types.NumberExpr{Value: 0}, wrapValue(&types.ListValue{Elements: []types.Value{}})},
			expectError: true,
		},
		{
			name:        "partition non-collection",
			args:        []types.Expr{&types.NumberExpr{Value: 2}, &types.NumberExpr{Value: 42}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalPartition(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalPartition failed: %v", err)
			}

			// Should return a list of lists
			if _, ok := result.(*types.ListValue); !ok {
				t.Errorf("Expected ListValue, got %T", result)
			}
		})
	}
}

func TestUtilsPlugin_EvalFlatten(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewUtilsPlugin(evaluator)

	tests := []struct {
		name        string
		args        []types.Expr
		expectError bool
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name: "flatten nested list",
			args: []types.Expr{
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(1),
					&types.ListValue{Elements: []types.Value{
						types.NumberValue(2),
						types.NumberValue(3),
					}},
					types.NumberValue(4),
				}}),
			},
			expectError: false,
		},
		{
			name:        "flatten single value",
			args:        []types.Expr{&types.NumberExpr{Value: 42}},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalFlatten(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalFlatten failed: %v", err)
			}

			// Should return a flattened list
			if _, ok := result.(*types.ListValue); !ok {
				t.Errorf("Expected ListValue, got %T", result)
			}
		})
	}
}

func TestUtilsPlugin_EvalInterpose(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewUtilsPlugin(evaluator)

	tests := []struct {
		name        string
		args        []types.Expr
		expectError bool
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name: "interpose separator between elements",
			args: []types.Expr{
				&types.StringExpr{Value: ","},
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.StringValue("a"),
					types.StringValue("b"),
					types.StringValue("c"),
				}}),
			},
			expectError: false,
		},
		{
			name:        "interpose in non-collection",
			args:        []types.Expr{&types.StringExpr{Value: ","}, &types.NumberExpr{Value: 42}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalInterpose(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalInterpose failed: %v", err)
			}

			// Should return a list with interposed elements
			if _, ok := result.(*types.ListValue); !ok {
				t.Errorf("Expected ListValue, got %T", result)
			}
		})
	}
}

func TestUtilsPlugin_EvalShuffle(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewUtilsPlugin(evaluator)

	tests := []struct {
		name        string
		args        []types.Expr
		expectError bool
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name: "shuffle list",
			args: []types.Expr{
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(1),
					types.NumberValue(2),
					types.NumberValue(3),
					types.NumberValue(4),
					types.NumberValue(5),
				}}),
			},
			expectError: false,
		},
		{
			name:        "shuffle non-collection",
			args:        []types.Expr{&types.NumberExpr{Value: 42}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalShuffle(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalShuffle failed: %v", err)
			}

			// Should return a shuffled list
			if _, ok := result.(*types.ListValue); !ok {
				t.Errorf("Expected ListValue, got %T", result)
			}
		})
	}
}

func TestUtilsPlugin_ErrorHandling(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewUtilsPlugin(evaluator)

	// Test various error conditions
	errorTests := []struct {
		name     string
		function func(registry.Evaluator, []types.Expr) (types.Value, error)
		args     []types.Expr
	}{
		{
			name:     "frequencies with no args",
			function: plugin.evalFrequencies,
			args:     []types.Expr{},
		},
		{
			name:     "partition with no args",
			function: plugin.evalPartition,
			args:     []types.Expr{},
		},
		{
			name:     "flatten with no args",
			function: plugin.evalFlatten,
			args:     []types.Expr{},
		},
		{
			name:     "interpose with no args",
			function: plugin.evalInterpose,
			args:     []types.Expr{},
		},
		{
			name:     "shuffle with no args",
			function: plugin.evalShuffle,
			args:     []types.Expr{},
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.function(evaluator, tt.args)
			if err == nil {
				t.Errorf("Expected error for %s, but got none", tt.name)
			}
		})
	}
}
