package polymorphic

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// Mock evaluator for testing
type mockEvaluator struct{}

func newMockEvaluator() *mockEvaluator {
	return &mockEvaluator{}
}

func (me *mockEvaluator) Eval(expr types.Expr) (types.Value, error) {
	// Basic implementation for testing
	switch e := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(e.Value), nil
	case *types.StringExpr:
		return types.StringValue(e.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(e.Value), nil
	default:
		return nil, nil
	}
}

func (me *mockEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	return nil, nil
}

func (me *mockEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	return me.Eval(expr)
}

func TestPolymorphicPlugin_RegisterFunctions(t *testing.T) {
	mockEval := newMockEvaluator()
	plugin := NewPolymorphicPlugin(mockEval)
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{
		"first", "rest", "last", "nth", "second", "empty?", "seq",
		"take", "drop", "reverse", "distinct", "sort", "into",
		"seq?", "coll?", "sequential?", "indexed?",
		"identity", "constantly",
	}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestPolymorphicPlugin_ExtractSequence(t *testing.T) {
	mockEval := newMockEvaluator()
	plugin := NewPolymorphicPlugin(mockEval)

	// Test with list
	listValue := &types.ListValue{Elements: []types.Value{
		types.NumberValue(1),
		types.NumberValue(2),
		types.NumberValue(3),
	}}

	elements, err := plugin.extractSequence(listValue)
	if err != nil {
		t.Errorf("extractSequence failed for list: %v", err)
	}
	if len(elements) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(elements))
	}

	// Test with vector
	vectorValue := types.NewVectorValue([]types.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("c"),
	})

	elements, err = plugin.extractSequence(vectorValue)
	if err != nil {
		t.Errorf("extractSequence failed for vector: %v", err)
	}
	if len(elements) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(elements))
	}

	// Test with string
	stringValue := types.StringValue("abc")
	elements, err = plugin.extractSequence(stringValue)
	if err != nil {
		t.Errorf("extractSequence failed for string: %v", err)
	}
	if len(elements) != 3 {
		t.Errorf("Expected 3 elements for string 'abc', got %d", len(elements))
	}

	// Test with nil
	nilValue := &types.NilValue{}
	elements, err = plugin.extractSequence(nilValue)
	if err != nil {
		t.Errorf("extractSequence failed for nil: %v", err)
	}
	if len(elements) != 0 {
		t.Errorf("Expected 0 elements for nil, got %d", len(elements))
	}
}

func TestPolymorphicPlugin_PluginInfo(t *testing.T) {
	mockEval := newMockEvaluator()
	plugin := NewPolymorphicPlugin(mockEval)

	if plugin.Name() != "polymorphic" {
		t.Errorf("Expected plugin name 'polymorphic', got '%s'", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", plugin.Version())
	}
}
