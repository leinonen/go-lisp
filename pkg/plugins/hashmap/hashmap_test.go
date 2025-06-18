package hashmap

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
	default:
		if ve, ok := expr.(valueExpr); ok {
			return ve.value, nil
		}
		if val, ok := expr.(types.Value); ok {
			return val, nil
		}
		return nil, nil
	}
}

func (me *mockEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	return nil, nil
}

func (me *mockEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	// For testing purposes, just call regular Eval
	// In a real implementation, this would use the bindings
	return me.Eval(expr)
}

func wrapValue(value types.Value) types.Expr {
	return valueExpr{value}
}

type valueExpr struct {
	value types.Value
}

func (ve valueExpr) String() string {
	return ve.value.String()
}

func TestHashMapPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewHashMapPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{
		"hash-map", "hash-map-get", "hash-map-put", "hash-map-remove",
		"hash-map-contains?", "hash-map-keys", "hash-map-values",
		"hash-map-size", "hash-map-empty?",
	}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestHashMapPlugin_BasicOperations(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Test creating a hash map
	hashMapArgs := []types.Expr{
		&types.StringExpr{Value: "key1"},
		&types.StringExpr{Value: "value1"},
		&types.StringExpr{Value: "key2"},
		&types.NumberExpr{Value: 42},
	}

	result, err := plugin.evalHashMap(evaluator, hashMapArgs)
	if err != nil {
		t.Fatalf("evalHashMap failed: %v", err)
	}

	hashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected HashMapValue, got %T", result)
	}

	// Test getting values
	getArgs := []types.Expr{
		wrapValue(hashMap),
		&types.StringExpr{Value: "key1"},
	}

	result, err = plugin.evalHashMapGet(evaluator, getArgs)
	if err != nil {
		t.Fatalf("evalHashMapGet failed: %v", err)
	}

	if result != types.StringValue("value1") {
		t.Errorf("Expected 'value1', got %v", result)
	}

	// Test putting new value
	putArgs := []types.Expr{
		wrapValue(hashMap),
		&types.StringExpr{Value: "key3"},
		&types.StringExpr{Value: "value3"},
	}

	result, err = plugin.evalHashMapPut(evaluator, putArgs)
	if err != nil {
		t.Fatalf("evalHashMapPut failed: %v", err)
	}

	newHashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected HashMapValue, got %T", result)
	}

	// Test contains
	containsArgs := []types.Expr{
		wrapValue(newHashMap),
		&types.StringExpr{Value: "key3"},
	}

	result, err = plugin.evalHashMapContains(evaluator, containsArgs)
	if err != nil {
		t.Fatalf("evalHashMapContains failed: %v", err)
	}

	if result != types.BooleanValue(true) {
		t.Errorf("Expected true, got %v", result)
	}

	// Test size
	sizeArgs := []types.Expr{wrapValue(newHashMap)}
	result, err = plugin.evalHashMapSize(evaluator, sizeArgs)
	if err != nil {
		t.Fatalf("evalHashMapSize failed: %v", err)
	}

	if result != types.NumberValue(3) {
		t.Errorf("Expected 3, got %v", result)
	}

	// Test empty?
	emptyArgs := []types.Expr{wrapValue(newHashMap)}
	result, err = plugin.evalHashMapEmpty(evaluator, emptyArgs)
	if err != nil {
		t.Fatalf("evalHashMapEmpty failed: %v", err)
	}

	if result != types.BooleanValue(false) {
		t.Errorf("Expected false, got %v", result)
	}
}

func TestHashMapPlugin_EmptyHashMap(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Test creating empty hash map
	result, err := plugin.evalHashMap(evaluator, []types.Expr{})
	if err != nil {
		t.Fatalf("evalHashMap failed: %v", err)
	}

	hashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected HashMapValue, got %T", result)
	}

	// Test empty? on empty map
	emptyArgs := []types.Expr{wrapValue(hashMap)}
	result, err = plugin.evalHashMapEmpty(evaluator, emptyArgs)
	if err != nil {
		t.Fatalf("evalHashMapEmpty failed: %v", err)
	}

	if result != types.BooleanValue(true) {
		t.Errorf("Expected true, got %v", result)
	}

	// Test size on empty map
	sizeArgs := []types.Expr{wrapValue(hashMap)}
	result, err = plugin.evalHashMapSize(evaluator, sizeArgs)
	if err != nil {
		t.Fatalf("evalHashMapSize failed: %v", err)
	}

	if result != types.NumberValue(0) {
		t.Errorf("Expected 0, got %v", result)
	}
}

func TestHashMapPlugin_ErrorCases(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Test hash-map with odd number of arguments
	_, err := plugin.evalHashMap(evaluator, []types.Expr{
		&types.StringExpr{Value: "key1"},
		&types.StringExpr{Value: "value1"},
		&types.StringExpr{Value: "key2"},
		// Missing value for key2
	})
	if err == nil {
		t.Error("Expected error for odd number of arguments")
	}

	// Test get with non-hashmap
	nonHashMap := types.NumberValue(42)
	_, err = plugin.evalHashMapGet(evaluator, []types.Expr{
		wrapValue(nonHashMap),
		&types.StringExpr{Value: "key"},
	})
	if err == nil {
		t.Error("Expected error for non-hashmap argument")
	}
}
