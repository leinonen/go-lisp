package macro

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// Mock evaluator for testing macro features
type mockEvaluator struct {
	callHistory []string
}

func newMockEvaluator() *mockEvaluator {
	return &mockEvaluator{
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
		return types.StringValue(e.Name), nil
	case *types.ListExpr:
		return types.StringValue("list-result"), nil
	default:
		return types.StringValue("unknown"), nil
	}
}

func (me *mockEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	return types.StringValue("function-result"), nil
}

func TestMacroPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewMacroPlugin()
	registry := registry.NewRegistry()

	err := plugin.RegisterFunctions(registry)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	// Check that defmacro function is registered
	if !registry.Has("defmacro") {
		t.Error("defmacro function not registered")
	}

	// Check that macroexpand function is registered
	if !registry.Has("macroexpand") {
		t.Error("macroexpand function not registered")
	}

	// Check that unquote function is registered
	if !registry.Has("unquote") {
		t.Error("unquote function not registered")
	}

	// Get defmacro function and verify properties
	defmacroFunc, exists := registry.Get("defmacro")
	if !exists {
		t.Fatal("defmacro function not found in registry")
	}

	if defmacroFunc.Name() != "defmacro" {
		t.Errorf("Expected function name 'defmacro', got '%s'", defmacroFunc.Name())
	}

	if defmacroFunc.Arity() != 3 {
		t.Errorf("Expected arity 3, got %d", defmacroFunc.Arity())
	}

	// Get macroexpand function and verify properties
	macroexpandFunc, exists := registry.Get("macroexpand")
	if !exists {
		t.Fatal("macroexpand function not found in registry")
	}

	if macroexpandFunc.Name() != "macroexpand" {
		t.Errorf("Expected function name 'macroexpand', got '%s'", macroexpandFunc.Name())
	}

	if macroexpandFunc.Arity() != 1 {
		t.Errorf("Expected arity 1, got %d", macroexpandFunc.Arity())
	}
}

func TestMacroPlugin_EvalDefmacro_BasicMacros(t *testing.T) {
	plugin := NewMacroPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name       string
		macroName  string
		params     []string
		body       types.Expr
		expectType string
	}{
		{
			name:      "simple macro",
			macroName: "when",
			params:    []string{"condition", "body"},
			body: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "if"},
					&types.SymbolExpr{Name: "condition"},
					&types.SymbolExpr{Name: "body"},
				},
			},
			expectType: "*types.MacroValue",
		},
		{
			name:       "macro with no params",
			macroName:  "nil",
			params:     []string{},
			body:       &types.BooleanExpr{Value: false},
			expectType: "*types.MacroValue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build parameter bracket expression
			var paramExprs []types.Expr
			for _, param := range tt.params {
				paramExprs = append(paramExprs, &types.SymbolExpr{Name: param})
			}

			args := []types.Expr{
				&types.SymbolExpr{Name: tt.macroName},
				&types.BracketExpr{Elements: paramExprs},
				tt.body,
			}

			result, err := plugin.evalDefmacro(evaluator, args)
			if err != nil {
				t.Fatalf("evalDefmacro failed: %v", err)
			}

			// Check that result is a MacroValue
			macro, ok := result.(*types.MacroValue)
			if !ok {
				t.Errorf("Expected *types.MacroValue, got %T", result)
				return
			}

			// Verify parameters
			if len(macro.Params) != len(tt.params) {
				t.Errorf("Expected %d params, got %d", len(tt.params), len(macro.Params))
			}

			for i, param := range tt.params {
				if i < len(macro.Params) && macro.Params[i] != param {
					t.Errorf("Expected param %d to be %s, got %s", i, param, macro.Params[i])
				}
			}

			// Verify body is stored
			if macro.Body == nil {
				t.Error("Expected macro body to be stored")
			}
		})
	}
}

func TestMacroPlugin_EvalDefmacro_Errors(t *testing.T) {
	plugin := NewMacroPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expectedErr string
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectedErr: "defmacro requires exactly 3 arguments, got 0",
		},
		{
			name: "non-symbol name",
			args: []types.Expr{
				&types.NumberExpr{Value: 42},
				&types.BracketExpr{Elements: []types.Expr{}},
				&types.BooleanExpr{Value: true},
			},
			expectedErr: "defmacro requires a symbol as first argument, got *types.NumberExpr",
		},
		{
			name: "non-bracket parameters",
			args: []types.Expr{
				&types.SymbolExpr{Name: "mymacro"},
				&types.ListExpr{Elements: []types.Expr{}},
				&types.BooleanExpr{Value: true},
			},
			expectedErr: "defmacro requires a bracket expression for parameters, got *types.ListExpr",
		},
		{
			name: "non-symbol parameter",
			args: []types.Expr{
				&types.SymbolExpr{Name: "mymacro"},
				&types.BracketExpr{Elements: []types.Expr{
					&types.NumberExpr{Value: 42},
				}},
				&types.BooleanExpr{Value: true},
			},
			expectedErr: "defmacro parameter 0 must be a symbol, got *types.NumberExpr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := plugin.evalDefmacro(evaluator, tt.args)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestMacroPlugin_EvalMacroexpand_Basic(t *testing.T) {
	plugin := NewMacroPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		input    types.Expr
		expected string
	}{
		{
			name:     "expand symbol",
			input:    &types.SymbolExpr{Name: "x"},
			expected: "*types.QuotedValue",
		},
		{
			name: "expand list",
			input: &types.ListExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "when"},
				&types.BooleanExpr{Value: true},
				&types.NumberExpr{Value: 42},
			}},
			expected: "*types.QuotedValue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{tt.input}

			result, err := plugin.evalMacroexpand(evaluator, args)
			if err != nil {
				t.Fatalf("evalMacroexpand failed: %v", err)
			}

			// Check that result is a QuotedValue
			quoted, ok := result.(*types.QuotedValue)
			if !ok {
				t.Errorf("Expected *types.QuotedValue, got %T", result)
				return
			}

			// Verify the value is wrapped
			if quoted.Value == nil {
				t.Error("Expected quoted value to contain the input")
			}
		})
	}
}

func TestMacroPlugin_EvalMacroexpand_Errors(t *testing.T) {
	plugin := NewMacroPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expectedErr string
	}{
		{
			name:        "wrong number of arguments - none",
			args:        []types.Expr{},
			expectedErr: "macroexpand requires exactly 1 argument, got 0",
		},
		{
			name: "wrong number of arguments - too many",
			args: []types.Expr{
				&types.SymbolExpr{Name: "a"},
				&types.SymbolExpr{Name: "b"},
			},
			expectedErr: "macroexpand requires exactly 1 argument, got 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := plugin.evalMacroexpand(evaluator, tt.args)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestMacroPlugin_MacroIntegration(t *testing.T) {
	// This test shows how macros would work in practice
	plugin := NewMacroPlugin()
	evaluator := newMockEvaluator()

	// Define a simple "when" macro: (defmacro when [condition body] (if condition body))
	defmacroArgs := []types.Expr{
		&types.SymbolExpr{Name: "when"},
		&types.BracketExpr{Elements: []types.Expr{
			&types.SymbolExpr{Name: "condition"},
			&types.SymbolExpr{Name: "body"},
		}},
		&types.ListExpr{Elements: []types.Expr{
			&types.SymbolExpr{Name: "if"},
			&types.SymbolExpr{Name: "condition"},
			&types.SymbolExpr{Name: "body"},
		}},
	}

	macroResult, err := plugin.evalDefmacro(evaluator, defmacroArgs)
	if err != nil {
		t.Fatalf("Failed to define macro: %v", err)
	}

	// Verify macro was created
	macro, ok := macroResult.(*types.MacroValue)
	if !ok {
		t.Fatalf("Expected MacroValue, got %T", macroResult)
	}

	if len(macro.Params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(macro.Params))
	}

	if macro.Params[0] != "condition" || macro.Params[1] != "body" {
		t.Errorf("Expected params [condition, body], got %v", macro.Params)
	}

	// Test macroexpand on a hypothetical macro call
	macroCall := &types.ListExpr{Elements: []types.Expr{
		&types.SymbolExpr{Name: "when"},
		&types.BooleanExpr{Value: true},
		&types.NumberExpr{Value: 42},
	}}

	expandResult, err := plugin.evalMacroexpand(evaluator, []types.Expr{macroCall})
	if err != nil {
		t.Fatalf("Failed to expand macro: %v", err)
	}

	// Should return a quoted value
	if _, ok := expandResult.(*types.QuotedValue); !ok {
		t.Errorf("Expected QuotedValue, got %T", expandResult)
	}
}

func TestMacroPlugin_EvalUnquote(t *testing.T) {
	plugin := NewMacroPlugin()
	evaluator := newMockEvaluator()

	// Test unquote with a simple expression
	args := []types.Expr{
		&types.NumberExpr{Value: 42},
	}

	result, err := plugin.evalUnquote(evaluator, args)
	if err != nil {
		t.Fatalf("evalUnquote failed: %v", err)
	}

	// Should return the evaluated expression
	expected := types.NumberValue(42)
	if !valuesEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMacroPlugin_ExpandMacro(t *testing.T) {
	plugin := NewMacroPlugin()

	// Create a simple macro: (defmacro inc [x] (+ x 1))
	macro := &types.MacroValue{
		Params: []string{"x"},
		Body: &types.ListExpr{Elements: []types.Expr{
			&types.SymbolExpr{Name: "+"},
			&types.SymbolExpr{Name: "x"},
			&types.NumberExpr{Value: 1},
		}},
	}

	// Expand (inc 5)
	args := []types.Expr{
		&types.NumberExpr{Value: 5},
	}

	expanded, err := plugin.expandMacro(macro, args)
	if err != nil {
		t.Fatalf("expandMacro failed: %v", err)
	}

	// Should expand to (+ 5 1)
	if listExpr, ok := expanded.(*types.ListExpr); ok {
		if len(listExpr.Elements) != 3 {
			t.Errorf("Expected 3 elements in expanded expression, got %d", len(listExpr.Elements))
		}

		// Check first element is +
		if symbol, ok := listExpr.Elements[0].(*types.SymbolExpr); !ok || symbol.Name != "+" {
			t.Errorf("Expected first element to be '+', got %v", listExpr.Elements[0])
		}

		// Check second element is 5
		if num, ok := listExpr.Elements[1].(*types.NumberExpr); !ok || num.Value != 5 {
			t.Errorf("Expected second element to be 5, got %v", listExpr.Elements[1])
		}

		// Check third element is 1
		if num, ok := listExpr.Elements[2].(*types.NumberExpr); !ok || num.Value != 1 {
			t.Errorf("Expected third element to be 1, got %v", listExpr.Elements[2])
		}
	} else {
		t.Errorf("Expected *types.ListExpr, got %T", expanded)
	}
}

func TestMacroPlugin_SubstituteInExpr(t *testing.T) {
	plugin := NewMacroPlugin()

	tests := []struct {
		name     string
		expr     types.Expr
		params   []string
		args     []types.Expr
		expected string
	}{
		{
			name:     "substitute symbol",
			expr:     &types.SymbolExpr{Name: "x"},
			params:   []string{"x"},
			args:     []types.Expr{&types.NumberExpr{Value: 42}},
			expected: "number 42",
		},
		{
			name:     "no substitution needed",
			expr:     &types.SymbolExpr{Name: "y"},
			params:   []string{"x"},
			args:     []types.Expr{&types.NumberExpr{Value: 42}},
			expected: "symbol y",
		},
		{
			name: "substitute in list",
			expr: &types.ListExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "+"},
				&types.SymbolExpr{Name: "x"},
				&types.NumberExpr{Value: 1},
			}},
			params:   []string{"x"},
			args:     []types.Expr{&types.NumberExpr{Value: 5}},
			expected: "list with substitution",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := plugin.substituteInExpr(tt.expr, tt.params, tt.args)

			// Basic verification that substitution occurred
			if result == nil {
				t.Error("Expected result but got nil")
			}

			// For the symbol substitution test, verify the actual substitution
			if tt.name == "substitute symbol" {
				if num, ok := result.(*types.NumberExpr); !ok || num.Value != 42 {
					t.Errorf("Expected NumberExpr with value 42, got %v", result)
				}
			}

			// For the list substitution test, verify structure
			if tt.name == "substitute in list" {
				if listExpr, ok := result.(*types.ListExpr); ok {
					if len(listExpr.Elements) != 3 {
						t.Errorf("Expected 3 elements, got %d", len(listExpr.Elements))
					}
					// Check that x was substituted with 5
					if num, ok := listExpr.Elements[1].(*types.NumberExpr); !ok || num.Value != 5 {
						t.Errorf("Expected second element to be 5, got %v", listExpr.Elements[1])
					}
				} else {
					t.Errorf("Expected *types.ListExpr, got %T", result)
				}
			}
		})
	}
}

func TestMacroPlugin_EvalMacroexpand_BuiltinMacros(t *testing.T) {
	plugin := NewMacroPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		input    types.Expr
		expected string
	}{
		{
			name: "expand when macro",
			input: &types.ListExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "when"},
				&types.BooleanExpr{Value: true},
				&types.NumberExpr{Value: 42},
			}},
			expected: "should expand to if",
		},
		{
			name: "expand unless macro",
			input: &types.ListExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "unless"},
				&types.BooleanExpr{Value: false},
				&types.NumberExpr{Value: 42},
			}},
			expected: "should expand to if with inverted logic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{tt.input}

			result, err := plugin.evalMacroexpand(evaluator, args)
			if err != nil {
				t.Fatalf("evalMacroexpand failed: %v", err)
			}

			// Should return a QuotedValue containing the expanded expression
			quoted, ok := result.(*types.QuotedValue)
			if !ok {
				t.Errorf("Expected *types.QuotedValue, got %T", result)
				return
			}

			// Check that the expanded expression is an if expression
			if listExpr, ok := quoted.Value.(*types.ListExpr); ok {
				if len(listExpr.Elements) == 0 {
					t.Error("Expected non-empty expanded expression")
					return
				}
				if symbol, ok := listExpr.Elements[0].(*types.SymbolExpr); !ok || symbol.Name != "if" {
					t.Errorf("Expected expanded expression to start with 'if', got %v", listExpr.Elements[0])
				}
			} else {
				t.Errorf("Expected expanded expression to be a list, got %T", quoted.Value)
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
	case *types.MacroValue:
		if vb, ok := b.(*types.MacroValue); ok {
			return len(va.Params) == len(vb.Params) // Simplified comparison
		}
	case *types.QuotedValue:
		if vb, ok := b.(*types.QuotedValue); ok {
			return va.Value != nil && vb.Value != nil // Simplified comparison
		}
	}
	return false
}
