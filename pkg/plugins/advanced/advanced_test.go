package advanced

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// Mock evaluator for testing advanced features
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
				case "let":
					// Mock let implementation - just return a test value
					return types.NumberValue(42), nil
				case "fn":
					// Mock function creation
					return types.StringValue("function"), nil
				case "nth":
					// Mock nth function - return the index as a number
					if len(e.Elements) == 3 {
						if idx, ok := e.Elements[2].(*types.NumberExpr); ok {
							return types.NumberValue(idx.Value + 10), nil // Add 10 to make it identifiable
						}
					}
					return types.NumberValue(0), nil
				case "+":
					// Mock addition
					if len(e.Elements) == 3 {
						left, _ := me.Eval(e.Elements[1])
						right, _ := me.Eval(e.Elements[2])
						if l, ok := left.(types.NumberValue); ok {
							if r, ok := right.(types.NumberValue); ok {
								return types.NumberValue(l + r), nil
							}
						}
					}
					return types.NumberValue(0), nil
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

func TestAdvancedPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewAdvancedPlugin()
	registry := registry.NewRegistry()

	err := plugin.RegisterFunctions(registry)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	// Check that let-destructure function is registered
	if !registry.Has("let-destructure") {
		t.Error("let-destructure function not registered")
	}

	// Check that fn-destructure function is registered
	if !registry.Has("fn-destructure") {
		t.Error("fn-destructure function not registered")
	}

	// Verify function properties
	letDestructureFunc, exists := registry.Get("let-destructure")
	if !exists {
		t.Fatal("let-destructure function not found in registry")
	}

	if letDestructureFunc.Name() != "let-destructure" {
		t.Errorf("Expected function name 'let-destructure', got '%s'", letDestructureFunc.Name())
	}

	if letDestructureFunc.Arity() != 2 {
		t.Errorf("Expected arity 2, got %d", letDestructureFunc.Arity())
	}
}

func TestAdvancedPlugin_ExpandPattern_SimpleBinding(t *testing.T) {
	plugin := NewAdvancedPlugin()

	// Test simple symbol binding (no destructuring)
	pattern := &types.SymbolExpr{Name: "x"}
	value := &types.NumberExpr{Value: 42}

	bindings, err := plugin.expandPattern(pattern, value)
	if err != nil {
		t.Fatalf("expandPattern failed: %v", err)
	}

	if len(bindings) != 2 {
		t.Errorf("Expected 2 bindings, got %d", len(bindings))
	}

	if symbol, ok := bindings[0].(*types.SymbolExpr); !ok || symbol.Name != "x" {
		t.Errorf("Expected first binding to be symbol 'x', got %v", bindings[0])
	}

	if num, ok := bindings[1].(*types.NumberExpr); !ok || num.Value != 42 {
		t.Errorf("Expected second binding to be number 42, got %v", bindings[1])
	}
}

func TestAdvancedPlugin_ExpandPattern_VectorDestructuring(t *testing.T) {
	plugin := NewAdvancedPlugin()

	// Test vector destructuring: [a b c]
	pattern := &types.BracketExpr{Elements: []types.Expr{
		&types.SymbolExpr{Name: "a"},
		&types.SymbolExpr{Name: "b"},
		&types.SymbolExpr{Name: "c"},
	}}
	value := &types.SymbolExpr{Name: "my-vec"}

	bindings, err := plugin.expandPattern(pattern, value)
	if err != nil {
		t.Fatalf("expandPattern failed: %v", err)
	}

	// Should generate 6 bindings: a (nth my-vec 0) b (nth my-vec 1) c (nth my-vec 2)
	if len(bindings) != 6 {
		t.Errorf("Expected 6 bindings, got %d", len(bindings))
	}

	// Check first binding: a
	if symbol, ok := bindings[0].(*types.SymbolExpr); !ok || symbol.Name != "a" {
		t.Errorf("Expected first binding to be symbol 'a', got %v", bindings[0])
	}

	// Check second binding: (nth my-vec 0)
	if nthExpr, ok := bindings[1].(*types.ListExpr); ok {
		if len(nthExpr.Elements) != 3 {
			t.Errorf("Expected nth expression with 3 elements, got %d", len(nthExpr.Elements))
		}
		if symbol, ok := nthExpr.Elements[0].(*types.SymbolExpr); !ok || symbol.Name != "nth" {
			t.Errorf("Expected nth function call, got %v", nthExpr.Elements[0])
		}
		if num, ok := nthExpr.Elements[2].(*types.NumberExpr); !ok || num.Value != 0 {
			t.Errorf("Expected index 0, got %v", nthExpr.Elements[2])
		}
	} else {
		t.Errorf("Expected nth expression, got %T", bindings[1])
	}
}

func TestAdvancedPlugin_ExpandPattern_ListDestructuring(t *testing.T) {
	plugin := NewAdvancedPlugin()

	// Test list destructuring: (a b)
	pattern := &types.ListExpr{Elements: []types.Expr{
		&types.SymbolExpr{Name: "first"},
		&types.SymbolExpr{Name: "second"},
	}}
	value := &types.SymbolExpr{Name: "my-list"}

	bindings, err := plugin.expandPattern(pattern, value)
	if err != nil {
		t.Fatalf("expandPattern failed: %v", err)
	}

	// Should generate 4 bindings: first (nth my-list 0) second (nth my-list 1)
	if len(bindings) != 4 {
		t.Errorf("Expected 4 bindings, got %d", len(bindings))
	}
}

func TestAdvancedPlugin_ExpandPattern_Errors(t *testing.T) {
	plugin := NewAdvancedPlugin()

	tests := []struct {
		name     string
		pattern  types.Expr
		value    types.Expr
		errorMsg string
	}{
		{
			name:     "unsupported pattern type",
			pattern:  &types.NumberExpr{Value: 42},
			value:    &types.SymbolExpr{Name: "x"},
			errorMsg: "unsupported destructuring pattern type: *types.NumberExpr",
		},
		{
			name: "non-symbol in vector pattern",
			pattern: &types.BracketExpr{Elements: []types.Expr{
				&types.NumberExpr{Value: 1},
			}},
			value:    &types.SymbolExpr{Name: "x"},
			errorMsg: "destructuring pattern elements must be symbols, got *types.NumberExpr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := plugin.expandPattern(tt.pattern, tt.value)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if err.Error() != tt.errorMsg {
				t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestAdvancedPlugin_EvalLetDestructure_Basic(t *testing.T) {
	plugin := NewAdvancedPlugin()
	evaluator := newMockEvaluator()

	// Test: (let-destructure [[a b] my-vec] (+ a b))
	args := []types.Expr{
		&types.BracketExpr{Elements: []types.Expr{
			&types.BracketExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "a"},
				&types.SymbolExpr{Name: "b"},
			}},
			&types.SymbolExpr{Name: "my-vec"},
		}},
		&types.ListExpr{Elements: []types.Expr{
			&types.SymbolExpr{Name: "+"},
			&types.SymbolExpr{Name: "a"},
			&types.SymbolExpr{Name: "b"},
		}},
	}

	result, err := plugin.evalLetDestructure(evaluator, args)
	if err != nil {
		t.Fatalf("evalLetDestructure failed: %v", err)
	}

	// Mock let returns 42
	expected := types.NumberValue(42)
	if !valuesEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestAdvancedPlugin_EvalLetDestructure_Errors(t *testing.T) {
	plugin := NewAdvancedPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expectedErr string
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectedErr: "let-destructure requires exactly 2 arguments, got 0",
		},
		{
			name: "non-bracket bindings",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{}},
				&types.NumberExpr{Value: 1},
			},
			expectedErr: "let-destructure requires a bracket expression for bindings, got *types.ListExpr",
		},
		{
			name: "odd number of bindings",
			args: []types.Expr{
				&types.BracketExpr{Elements: []types.Expr{
					&types.SymbolExpr{Name: "x"},
				}},
				&types.NumberExpr{Value: 1},
			},
			expectedErr: "let-destructure bindings must come in pairs, got 1 elements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := plugin.evalLetDestructure(evaluator, tt.args)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestAdvancedPlugin_EvalFnDestructure_Basic(t *testing.T) {
	plugin := NewAdvancedPlugin()
	evaluator := newMockEvaluator()

	// Test: (fn-destructure [[a b]] (+ a b))
	args := []types.Expr{
		&types.BracketExpr{Elements: []types.Expr{
			&types.BracketExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "a"},
				&types.SymbolExpr{Name: "b"},
			}},
		}},
		&types.ListExpr{Elements: []types.Expr{
			&types.SymbolExpr{Name: "+"},
			&types.SymbolExpr{Name: "a"},
			&types.SymbolExpr{Name: "b"},
		}},
	}

	result, err := plugin.evalFnDestructure(evaluator, args)
	if err != nil {
		t.Fatalf("evalFnDestructure failed: %v", err)
	}

	// Mock fn returns "function"
	expected := types.StringValue("function")
	if !valuesEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestAdvancedPlugin_EvalFnDestructure_Errors(t *testing.T) {
	plugin := NewAdvancedPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expectedErr string
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectedErr: "fn-destructure requires exactly 2 arguments, got 0",
		},
		{
			name: "non-bracket parameters",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{}},
				&types.NumberExpr{Value: 1},
			},
			expectedErr: "fn-destructure requires a bracket expression for parameters, got *types.ListExpr",
		},
		{
			name: "multiple parameters not supported",
			args: []types.Expr{
				&types.BracketExpr{Elements: []types.Expr{
					&types.SymbolExpr{Name: "a"},
					&types.SymbolExpr{Name: "b"},
				}},
				&types.NumberExpr{Value: 1},
			},
			expectedErr: "fn-destructure currently supports only one destructuring parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := plugin.evalFnDestructure(evaluator, tt.args)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
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
