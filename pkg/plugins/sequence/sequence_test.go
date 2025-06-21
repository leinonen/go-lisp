package sequence

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

// Helper function to compare values for equality
func valuesEqual(a, b types.Value) bool {
	switch av := a.(type) {
	case types.NumberValue:
		if bv, ok := b.(types.NumberValue); ok {
			return float64(av) == float64(bv)
		}
	case types.StringValue:
		if bv, ok := b.(types.StringValue); ok {
			return string(av) == string(bv)
		}
	case types.BooleanValue:
		if bv, ok := b.(types.BooleanValue); ok {
			return bool(av) == bool(bv)
		}
	case *types.VectorValue:
		if bv, ok := b.(*types.VectorValue); ok {
			if len(av.Elements) != len(bv.Elements) {
				return false
			}
			for i, elem := range av.Elements {
				if !valuesEqual(elem, bv.Elements[i]) {
					return false
				}
			}
			return true
		}
	case *types.ListValue:
		if bv, ok := b.(*types.ListValue); ok {
			if len(av.Elements) != len(bv.Elements) {
				return false
			}
			for i, elem := range av.Elements {
				if !valuesEqual(elem, bv.Elements[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func TestSequencePlugin_RegisterFunctions(t *testing.T) {
	plugin := NewSequencePluginLegacy()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"vector", "vec", "vector?", "conj"}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestSequencePlugin_EvalVector(t *testing.T) {
	plugin := NewSequencePluginLegacy()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		args     []types.Expr
		expected []types.Value
	}{
		{
			name:     "empty vector",
			args:     []types.Expr{},
			expected: []types.Value{},
		},
		{
			name: "single element vector",
			args: []types.Expr{
				&types.NumberExpr{Value: 1},
			},
			expected: []types.Value{types.NumberValue(1)},
		},
		{
			name: "multiple element vector",
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
			result, err := plugin.evalVector(evaluator, tt.args)
			if err != nil {
				t.Fatalf("evalVector failed: %v", err)
			}

			vector, ok := result.(*types.VectorValue)
			if !ok {
				t.Fatalf("Expected VectorValue, got %T", result)
			}

			if len(vector.Elements) != len(tt.expected) {
				t.Fatalf("Expected %d elements, got %d", len(tt.expected), len(vector.Elements))
			}

			for i, expected := range tt.expected {
				if !valuesEqual(vector.Elements[i], expected) {
					t.Errorf("Element %d: expected %v, got %v", i, expected, vector.Elements[i])
				}
			}
		})
	}
}

func TestSequencePlugin_EvalVec(t *testing.T) {
	plugin := NewSequencePluginLegacy()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    []types.Value
		expectError bool
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:        "too many arguments",
			args:        []types.Expr{&types.NumberExpr{Value: 1}, &types.NumberExpr{Value: 2}},
			expectError: true,
		},
		{
			name: "convert list to vector",
			args: []types.Expr{
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(1),
					types.NumberValue(2),
					types.NumberValue(3),
				}}),
			},
			expected: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			},
		},
		{
			name: "vector to vector (no change)",
			args: []types.Expr{
				wrapValue(types.NewVectorValue([]types.Value{
					types.StringValue("a"),
					types.StringValue("b"),
				})),
			},
			expected: []types.Value{
				types.StringValue("a"),
				types.StringValue("b"),
			},
		},
		{
			name:        "invalid type conversion",
			args:        []types.Expr{&types.NumberExpr{Value: 42}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalVec(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalVec failed: %v", err)
			}

			vector, ok := result.(*types.VectorValue)
			if !ok {
				t.Fatalf("Expected VectorValue, got %T", result)
			}

			if len(vector.Elements) != len(tt.expected) {
				t.Fatalf("Expected %d elements, got %d", len(tt.expected), len(vector.Elements))
			}

			for i, expected := range tt.expected {
				if !valuesEqual(vector.Elements[i], expected) {
					t.Errorf("Element %d: expected %v, got %v", i, expected, vector.Elements[i])
				}
			}
		})
	}
}

func TestSequencePlugin_EvalVectorPredicate(t *testing.T) {
	plugin := NewSequencePluginLegacy()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    bool
		expectError bool
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name: "vector is vector",
			args: []types.Expr{
				wrapValue(types.NewVectorValue([]types.Value{types.NumberValue(1)})),
			},
			expected: true,
		},
		{
			name: "list is not vector",
			args: []types.Expr{
				wrapValue(&types.ListValue{Elements: []types.Value{types.NumberValue(1)}}),
			},
			expected: false,
		},
		{
			name: "number is not vector",
			args: []types.Expr{
				&types.NumberExpr{Value: 42},
			},
			expected: false,
		},
		{
			name: "string is not vector",
			args: []types.Expr{
				&types.StringExpr{Value: "hello"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalVectorPredicate(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalVectorPredicate failed: %v", err)
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

func TestSequencePlugin_EvalConj(t *testing.T) {
	plugin := NewSequencePluginLegacy()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    []types.Value
		expectError bool
	}{
		{
			name:        "no arguments",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name: "conj to vector",
			args: []types.Expr{
				wrapValue(types.NewVectorValue([]types.Value{
					types.NumberValue(1),
					types.NumberValue(2),
				})),
				&types.NumberExpr{Value: 3},
				&types.NumberExpr{Value: 4},
			},
			expected: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
				types.NumberValue(4),
			},
		},
		{
			name: "conj to list (prepends)",
			args: []types.Expr{
				wrapValue(&types.ListValue{Elements: []types.Value{
					types.NumberValue(2),
					types.NumberValue(3),
				}}),
				&types.NumberExpr{Value: 1},
			},
			expected: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			},
		},
		{
			name: "conj single item to empty vector",
			args: []types.Expr{
				wrapValue(types.NewVectorValue([]types.Value{})),
				&types.StringExpr{Value: "hello"},
			},
			expected: []types.Value{
				types.StringValue("hello"),
			},
		},
		{
			name:        "conj to invalid type",
			args:        []types.Expr{&types.NumberExpr{Value: 42}, &types.NumberExpr{Value: 1}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalConj(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("evalConj failed: %v", err)
			}

			// Check if result is vector or list
			switch coll := result.(type) {
			case *types.VectorValue:
				if len(coll.Elements) != len(tt.expected) {
					t.Fatalf("Expected %d elements, got %d", len(tt.expected), len(coll.Elements))
				}

				for i, expected := range tt.expected {
					if !valuesEqual(coll.Elements[i], expected) {
						t.Errorf("Element %d: expected %v, got %v", i, expected, coll.Elements[i])
					}
				}
			case *types.ListValue:
				if len(coll.Elements) != len(tt.expected) {
					t.Fatalf("Expected %d elements, got %d", len(tt.expected), len(coll.Elements))
				}

				for i, expected := range tt.expected {
					if !valuesEqual(coll.Elements[i], expected) {
						t.Errorf("Element %d: expected %v, got %v", i, expected, coll.Elements[i])
					}
				}
			default:
				t.Fatalf("Expected VectorValue or ListValue, got %T", result)
			}
		})
	}
}

func TestSequencePlugin_PluginInfo(t *testing.T) {
	plugin := NewSequencePluginLegacy()

	if plugin.Name() != "sequence" {
		t.Errorf("Expected plugin name 'sequence', got '%s'", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", plugin.Version())
	}

	if plugin.Description() != "Vectors and sequences (vector, vec, seq)" {
		t.Errorf("Expected specific description, got '%s'", plugin.Description())
	}

	deps := plugin.Dependencies()
	expectedDeps := []string{"list"}
	if len(deps) != len(expectedDeps) {
		t.Errorf("Expected %d dependencies, got %d", len(expectedDeps), len(deps))
	}

	for i, expectedDep := range expectedDeps {
		if deps[i] != expectedDep {
			t.Errorf("Expected dependency '%s', got '%s'", expectedDep, deps[i])
		}
	}
}
