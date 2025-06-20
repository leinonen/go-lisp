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

func (ve valueExpr) GetPosition() types.Position {
	return types.Position{Line: 1, Column: 1}
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
		"hash-map-size", "hash-map-empty?", "assoc", "dissoc",
		// Clojure-style aliases
		"get", "contains?", "keys", "vals",
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

func TestHashMapPlugin_Assoc(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Create initial hash map
	initialHashMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			"key1": types.StringValue("value1"),
			"key2": types.NumberValue(42),
		},
	}

	// Test assoc with one key-value pair
	assocArgs := []types.Expr{
		wrapValue(initialHashMap),
		&types.StringExpr{Value: "key3"},
		&types.StringExpr{Value: "value3"},
	}

	result, err := plugin.evalAssoc(evaluator, assocArgs)
	if err != nil {
		t.Fatalf("evalAssoc failed: %v", err)
	}

	newHashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected HashMapValue, got %T", result)
	}

	// Check that the new key was added
	if newHashMap.Elements["key3"] != types.StringValue("value3") {
		t.Errorf("Expected 'value3' for key3, got %v", newHashMap.Elements["key3"])
	}

	// Check that original keys are still present
	if newHashMap.Elements["key1"] != types.StringValue("value1") {
		t.Errorf("Expected original key1 to be preserved")
	}

	// Test assoc with multiple key-value pairs
	multiAssocArgs := []types.Expr{
		wrapValue(initialHashMap),
		&types.StringExpr{Value: "key3"},
		&types.StringExpr{Value: "value3"},
		&types.StringExpr{Value: "key4"},
		&types.NumberExpr{Value: 100},
	}

	result, err = plugin.evalAssoc(evaluator, multiAssocArgs)
	if err != nil {
		t.Fatalf("evalAssoc with multiple pairs failed: %v", err)
	}

	multiHashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected HashMapValue, got %T", result)
	}

	// Check that both new keys were added
	if multiHashMap.Elements["key3"] != types.StringValue("value3") {
		t.Errorf("Expected 'value3' for key3, got %v", multiHashMap.Elements["key3"])
	}
	if multiHashMap.Elements["key4"] != types.NumberValue(100) {
		t.Errorf("Expected 100 for key4, got %v", multiHashMap.Elements["key4"])
	}

	// Test assoc with existing key (should update)
	updateArgs := []types.Expr{
		wrapValue(initialHashMap),
		&types.StringExpr{Value: "key1"},
		&types.StringExpr{Value: "updated_value"},
	}

	result, err = plugin.evalAssoc(evaluator, updateArgs)
	if err != nil {
		t.Fatalf("evalAssoc update failed: %v", err)
	}

	updatedHashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected HashMapValue, got %T", result)
	}

	if updatedHashMap.Elements["key1"] != types.StringValue("updated_value") {
		t.Errorf("Expected 'updated_value' for key1, got %v", updatedHashMap.Elements["key1"])
	}
}

func TestHashMapPlugin_AssocErrors(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Test assoc with too few arguments
	_, err := plugin.evalAssoc(evaluator, []types.Expr{
		wrapValue(&types.HashMapValue{Elements: map[string]types.Value{}}),
		&types.StringExpr{Value: "key"},
		// Missing value
	})
	if err == nil {
		t.Error("Expected error for too few arguments")
	}

	// Test assoc with odd number of key-value pairs
	_, err = plugin.evalAssoc(evaluator, []types.Expr{
		wrapValue(&types.HashMapValue{Elements: map[string]types.Value{}}),
		&types.StringExpr{Value: "key1"},
		&types.StringExpr{Value: "value1"},
		&types.StringExpr{Value: "key2"},
		// Missing value for key2
	})
	if err == nil {
		t.Error("Expected error for odd number of key-value pairs")
	}

	// Test assoc with non-hashmap
	_, err = plugin.evalAssoc(evaluator, []types.Expr{
		wrapValue(types.NumberValue(42)),
		&types.StringExpr{Value: "key"},
		&types.StringExpr{Value: "value"},
	})
	if err == nil {
		t.Error("Expected error for non-hashmap argument")
	}
}

func TestHashMapPlugin_Dissoc(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Create initial hash map
	initialHashMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			"key1": types.StringValue("value1"),
			"key2": types.NumberValue(42),
			"key3": types.StringValue("value3"),
		},
	}

	// Test dissoc with one key
	dissocArgs := []types.Expr{
		wrapValue(initialHashMap),
		&types.StringExpr{Value: "key2"},
	}

	result, err := plugin.evalDissoc(evaluator, dissocArgs)
	if err != nil {
		t.Fatalf("evalDissoc failed: %v", err)
	}

	newHashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected HashMapValue, got %T", result)
	}

	// Check that the key was removed
	if _, exists := newHashMap.Elements["key2"]; exists {
		t.Errorf("Expected key2 to be removed")
	}

	// Check that other keys are still present
	if newHashMap.Elements["key1"] != types.StringValue("value1") {
		t.Errorf("Expected key1 to be preserved")
	}
	if newHashMap.Elements["key3"] != types.StringValue("value3") {
		t.Errorf("Expected key3 to be preserved")
	}

	// Test dissoc with multiple keys
	multiDissocArgs := []types.Expr{
		wrapValue(initialHashMap),
		&types.StringExpr{Value: "key1"},
		&types.StringExpr{Value: "key3"},
	}

	result, err = plugin.evalDissoc(evaluator, multiDissocArgs)
	if err != nil {
		t.Fatalf("evalDissoc with multiple keys failed: %v", err)
	}

	multiHashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected HashMapValue, got %T", result)
	}

	// Check that both keys were removed
	if _, exists := multiHashMap.Elements["key1"]; exists {
		t.Errorf("Expected key1 to be removed")
	}
	if _, exists := multiHashMap.Elements["key3"]; exists {
		t.Errorf("Expected key3 to be removed")
	}

	// Check that key2 is still present
	if multiHashMap.Elements["key2"] != types.NumberValue(42) {
		t.Errorf("Expected key2 to be preserved")
	}

	// Test dissoc with non-existing key (should not error)
	nonExistingArgs := []types.Expr{
		wrapValue(initialHashMap),
		&types.StringExpr{Value: "non_existing_key"},
	}

	result, err = plugin.evalDissoc(evaluator, nonExistingArgs)
	if err != nil {
		t.Fatalf("evalDissoc with non-existing key failed: %v", err)
	}

	resultHashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected HashMapValue, got %T", result)
	}

	// Should return same size since key didn't exist
	if len(resultHashMap.Elements) != len(initialHashMap.Elements) {
		t.Errorf("Expected same size after dissoc with non-existing key")
	}
}

func TestHashMapPlugin_DissocErrors(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Test dissoc with too few arguments
	_, err := plugin.evalDissoc(evaluator, []types.Expr{
		wrapValue(&types.HashMapValue{Elements: map[string]types.Value{}}),
		// Missing key
	})
	if err == nil {
		t.Error("Expected error for too few arguments")
	}

	// Test dissoc with non-hashmap
	_, err = plugin.evalDissoc(evaluator, []types.Expr{
		wrapValue(types.NumberValue(42)),
		&types.StringExpr{Value: "key"},
	})
	if err == nil {
		t.Error("Expected error for non-hashmap argument")
	}
}

// Tests for Clojure-style aliases
func TestHashMapPlugin_ClojureGet(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Create a hash map
	hashMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			"name": types.StringValue("Alice"),
			"age":  types.NumberValue(30),
		},
	}

	// Test get with existing key
	result, err := plugin.evalHashMapGet(evaluator, []types.Expr{
		wrapValue(hashMap),
		&types.StringExpr{Value: "name"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	strVal, ok := result.(types.StringValue)
	if !ok {
		t.Fatalf("Expected StringValue, got %T", result)
	}

	if string(strVal) != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", string(strVal))
	}

	// Test get with non-existing key
	result, err = plugin.evalHashMapGet(evaluator, []types.Expr{
		wrapValue(hashMap),
		&types.StringExpr{Value: "missing"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if _, ok := result.(*types.NilValue); !ok {
		t.Errorf("Expected NilValue for missing key, got %T", result)
	}
}

func TestHashMapPlugin_ClojureContains(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Create a hash map
	hashMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			"name": types.StringValue("Alice"),
			"age":  types.NumberValue(30),
		},
	}

	// Test contains? with existing key
	result, err := plugin.evalHashMapContains(evaluator, []types.Expr{
		wrapValue(hashMap),
		&types.StringExpr{Value: "name"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	boolVal, ok := result.(types.BooleanValue)
	if !ok {
		t.Fatalf("Expected BooleanValue, got %T", result)
	}

	if !bool(boolVal) {
		t.Error("Expected true for existing key")
	}

	// Test contains? with non-existing key
	result, err = plugin.evalHashMapContains(evaluator, []types.Expr{
		wrapValue(hashMap),
		&types.StringExpr{Value: "missing"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	boolVal, ok = result.(types.BooleanValue)
	if !ok {
		t.Fatalf("Expected BooleanValue, got %T", result)
	}

	if bool(boolVal) {
		t.Error("Expected false for non-existing key")
	}
}

func TestHashMapPlugin_ClojureKeys(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Create a hash map
	hashMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			"name": types.StringValue("Alice"),
			"age":  types.NumberValue(30),
		},
	}

	// Test keys
	result, err := plugin.evalHashMapKeys(evaluator, []types.Expr{
		wrapValue(hashMap),
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	listVal, ok := result.(*types.ListValue)
	if !ok {
		t.Fatalf("Expected ListValue, got %T", result)
	}

	if len(listVal.Elements) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(listVal.Elements))
	}

	// Check that we have the expected keys (order may vary)
	keyStrings := make(map[string]bool)
	for _, key := range listVal.Elements {
		if strKey, ok := key.(types.StringValue); ok {
			keyStrings[string(strKey)] = true
		} else {
			t.Errorf("Expected string key, got %T", key)
		}
	}

	if !keyStrings["name"] {
		t.Error("Missing 'name' key")
	}
	if !keyStrings["age"] {
		t.Error("Missing 'age' key")
	}
}

func TestHashMapPlugin_ClojureVals(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Create a hash map
	hashMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			"name": types.StringValue("Alice"),
			"age":  types.NumberValue(30),
		},
	}

	// Test vals
	result, err := plugin.evalHashMapValues(evaluator, []types.Expr{
		wrapValue(hashMap),
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	listVal, ok := result.(*types.ListValue)
	if !ok {
		t.Fatalf("Expected ListValue, got %T", result)
	}

	if len(listVal.Elements) != 2 {
		t.Errorf("Expected 2 values, got %d", len(listVal.Elements))
	}

	// Check that we have the expected values (order may vary)
	hasAlice := false
	hasAge := false
	for _, value := range listVal.Elements {
		if strVal, ok := value.(types.StringValue); ok && string(strVal) == "Alice" {
			hasAlice = true
		}
		if numVal, ok := value.(types.NumberValue); ok && float64(numVal) == 30 {
			hasAge = true
		}
	}

	if !hasAlice {
		t.Error("Missing 'Alice' value")
	}
	if !hasAge {
		t.Error("Missing age value 30")
	}
}

func TestHashMapPlugin_ClojureCount(t *testing.T) {
	plugin := NewHashMapPlugin()
	evaluator := newMockEvaluator()

	// Create a hash map
	hashMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			"name": types.StringValue("Alice"),
			"age":  types.NumberValue(30),
		},
	}

	// Test count
	result, err := plugin.evalHashMapSize(evaluator, []types.Expr{
		wrapValue(hashMap),
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	numVal, ok := result.(types.NumberValue)
	if !ok {
		t.Fatalf("Expected NumberValue, got %T", result)
	}

	if float64(numVal) != 2 {
		t.Errorf("Expected count 2, got %f", float64(numVal))
	}

	// Test count on empty hash map
	emptyMap := &types.HashMapValue{Elements: make(map[string]types.Value)}
	result, err = plugin.evalHashMapSize(evaluator, []types.Expr{
		wrapValue(emptyMap),
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	numVal, ok = result.(types.NumberValue)
	if !ok {
		t.Fatalf("Expected NumberValue, got %T", result)
	}

	if float64(numVal) != 0 {
		t.Errorf("Expected count 0, got %f", float64(numVal))
	}
}
