package binding

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// Mock evaluator for testing
type mockEvaluator struct {
	env         *evaluator.Environment
	callHistory []string
}

func newMockEvaluator() *mockEvaluator {
	return &mockEvaluator{
		env:         evaluator.NewEnvironment(),
		callHistory: make([]string, 0),
	}
}

func (me *mockEvaluator) Eval(expr types.Expr) (types.Value, error) {
	me.callHistory = append(me.callHistory, expr.String())

	// Handle basic expressions for testing
	switch e := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(e.Value), nil
	case *types.StringExpr:
		return types.StringValue(e.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(e.Value), nil
	case *types.SymbolExpr:
		// Simulate variable lookup
		if val, exists := me.env.Get(e.Name); exists {
			return val, nil
		}
		// Return the symbol name as string for testing
		return types.StringValue(e.Name), nil
	case *types.ListExpr:
		// Simple simulation of evaluation
		if len(e.Elements) > 0 {
			if symbol, ok := e.Elements[0].(*types.SymbolExpr); ok {
				switch symbol.Name {
				case "def":
					// Simulate def operation
					if len(e.Elements) == 3 {
						if varSymbol, ok := e.Elements[1].(*types.SymbolExpr); ok {
							value, _ := me.Eval(e.Elements[2])
							me.env.Set(varSymbol.Name, value)
							return value, nil
						}
					}
				case "do":
					// Simulate do operation - evaluate all, return last
					var result types.Value = types.BooleanValue(false)
					for i := 1; i < len(e.Elements); i++ {
						result, _ = me.Eval(e.Elements[i])
					}
					return result, nil
				case "+":
					// Simulate addition
					if len(e.Elements) == 3 {
						left, _ := me.Eval(e.Elements[1])
						right, _ := me.Eval(e.Elements[2])
						if l, ok := left.(types.NumberValue); ok {
							if r, ok := right.(types.NumberValue); ok {
								return types.NumberValue(l + r), nil
							}
						}
					}
				}
			}
		}
		return types.StringValue("list-result"), nil
	default:
		return types.StringValue("unknown"), nil
	}
}

func (me *mockEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	return types.StringValue("function-result"), nil
}

func (me *mockEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	// For testing purposes, just call regular Eval
	// In a real implementation, this would use the bindings
	return me.Eval(expr)
}

func TestBindingPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewBindingPlugin()
	registry := registry.NewRegistry()

	err := plugin.RegisterFunctions(registry)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	// Check that let function is registered
	if !registry.Has("let") {
		t.Error("let function not registered")
	}

	// Get function and verify properties
	letFunc, exists := registry.Get("let")
	if !exists {
		t.Fatal("let function not found in registry")
	}

	if letFunc.Name() != "let" {
		t.Errorf("Expected function name 'let', got '%s'", letFunc.Name())
	}

	if letFunc.Arity() != 2 {
		t.Errorf("Expected arity 2, got %d", letFunc.Arity())
	}
}

func TestBindingPlugin_EvalLet_BasicBindings(t *testing.T) {
	plugin := NewBindingPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		bindings []types.Expr
		body     types.Expr
		expected types.Value
	}{
		{
			name: "single binding",
			bindings: []types.Expr{
				&types.SymbolExpr{Name: "x"},
				&types.NumberExpr{Value: 42},
			},
			body:     &types.SymbolExpr{Name: "x"},
			expected: types.NumberValue(42),
		},
		{
			name: "multiple bindings",
			bindings: []types.Expr{
				&types.SymbolExpr{Name: "x"},
				&types.NumberExpr{Value: 10},
				&types.SymbolExpr{Name: "y"},
				&types.NumberExpr{Value: 20},
			},
			body: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "y"},
				},
			},
			expected: types.NumberValue(30),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{
				&types.BracketExpr{Elements: tt.bindings},
				tt.body,
			}

			result, err := plugin.evalLet(evaluator, args)
			if err != nil {
				t.Fatalf("evalLet failed: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBindingPlugin_EvalLet_Errors(t *testing.T) {
	plugin := NewBindingPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expectedErr string
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectedErr: "let requires exactly 2 arguments, got 0",
		},
		{
			name: "non-bracket bindings",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{}},
				&types.NumberExpr{Value: 1},
			},
			expectedErr: "let requires a bracket expression for bindings, got *types.ListExpr",
		},
		{
			name: "odd number of bindings",
			args: []types.Expr{
				&types.BracketExpr{Elements: []types.Expr{
					&types.SymbolExpr{Name: "x"},
				}},
				&types.NumberExpr{Value: 1},
			},
			expectedErr: "let bindings must come in pairs, got 1 elements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := plugin.evalLet(evaluator, tt.args)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestBindingPlugin_EvalLet_SequentialBinding(t *testing.T) {
	// Test that later bindings can reference earlier ones
	plugin := NewBindingPlugin()
	evaluator := newMockEvaluator()

	args := []types.Expr{
		&types.BracketExpr{Elements: []types.Expr{
			&types.SymbolExpr{Name: "x"},
			&types.NumberExpr{Value: 5},
			&types.SymbolExpr{Name: "y"},
			&types.ListExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "+"},
				&types.SymbolExpr{Name: "x"},
				&types.NumberExpr{Value: 10},
			}},
		}},
		&types.SymbolExpr{Name: "y"},
	}

	result, err := plugin.evalLet(evaluator, args)
	if err != nil {
		t.Fatalf("evalLet failed: %v", err)
	}

	expected := types.NumberValue(15)
	if !valuesEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// Helper function to compare values
func valuesEqual(a, b types.Value) bool {
	switch va := a.(type) {
	case types.NumberValue:
		if vb, ok := b.(types.NumberValue); ok {
			return va == vb
		}
	case types.StringValue:
		if vb, ok := b.(types.StringValue); ok {
			return va == vb
		}
	case types.BooleanValue:
		if vb, ok := b.(types.BooleanValue); ok {
			return va == vb
		}
	case types.KeywordValue:
		if vb, ok := b.(types.KeywordValue); ok {
			return va == vb
		}
	}
	return false
}
