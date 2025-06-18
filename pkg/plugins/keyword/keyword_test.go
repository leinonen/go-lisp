package keyword

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
	case *types.StringExpr:
		return types.StringValue(e.Value), nil
	case *types.KeywordExpr:
		return types.KeywordValue(e.Value), nil
	default:
		return types.StringValue("test"), nil
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

func TestKeywordPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewKeywordPlugin()
	registry := registry.NewRegistry()

	err := plugin.RegisterFunctions(registry)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"keyword", "keyword?"}
	for _, funcName := range expectedFunctions {
		if !registry.Has(funcName) {
			t.Errorf("Function %s not registered", funcName)
		}
	}
}

func TestKeywordPlugin_EvalKeyword(t *testing.T) {
	plugin := NewKeywordPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		args     []types.Expr
		expected types.Value
	}{
		{
			name:     "string to keyword",
			args:     []types.Expr{&types.StringExpr{Value: "test"}},
			expected: types.KeywordValue("test"),
		},
		{
			name:     "keyword to keyword (idempotent)",
			args:     []types.Expr{&types.KeywordExpr{Value: "existing"}},
			expected: types.KeywordValue("existing"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalKeyword(evaluator, tt.args)
			if err != nil {
				t.Fatalf("evalKeyword failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestKeywordPlugin_EvalKeyword_Errors(t *testing.T) {
	plugin := NewKeywordPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expectedErr string
	}{
		{
			name:        "wrong number of arguments",
			args:        []types.Expr{},
			expectedErr: "keyword requires exactly 1 argument, got 0",
		},
		{
			name:        "too many arguments",
			args:        []types.Expr{&types.StringExpr{Value: "a"}, &types.StringExpr{Value: "b"}},
			expectedErr: "keyword requires exactly 1 argument, got 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := plugin.evalKeyword(evaluator, tt.args)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestKeywordPlugin_EvalKeywordPredicate(t *testing.T) {
	plugin := NewKeywordPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		args     []types.Expr
		expected types.Value
	}{
		{
			name:     "keyword is keyword",
			args:     []types.Expr{&types.KeywordExpr{Value: "test"}},
			expected: types.BooleanValue(true),
		},
		{
			name:     "string is not keyword",
			args:     []types.Expr{&types.StringExpr{Value: "test"}},
			expected: types.BooleanValue(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.evalKeywordPredicate(evaluator, tt.args)
			if err != nil {
				t.Fatalf("evalKeywordPredicate failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
