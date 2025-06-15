// Package hashmap implements hash map functions as a plugin
package hashmap

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// HashMapPlugin implements hash map functions
type HashMapPlugin struct {
	*plugins.BasePlugin
}

// NewHashMapPlugin creates a new hash map plugin
func NewHashMapPlugin() *HashMapPlugin {
	return &HashMapPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"hashmap",
			"1.0.0",
			"Hash map functions (hash-map, hash-map-get, hash-map-put, etc.)",
			[]string{}, // No dependencies
		),
	}
}

// Functions returns the list of functions provided by this plugin
func (p *HashMapPlugin) Functions() []string {
	return []string{
		"hash-map", "hash-map-get", "hash-map-put", "hash-map-remove",
		"hash-map-contains?", "hash-map-keys", "hash-map-values",
		"hash-map-size", "hash-map-empty?",
	}
}

// RegisterFunctions registers all hash map functions with the registry
func (p *HashMapPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// hash-map function
	hashMapFunc := functions.NewFunction(
		"hash-map",
		registry.CategoryHashMap,
		-1, // Variable arguments
		"Create a hash map from key-value pairs: (hash-map \"key1\" \"value1\" \"key2\" \"value2\")",
		p.evalHashMap,
	)
	if err := reg.Register(hashMapFunc); err != nil {
		return err
	}

	// hash-map-get function
	hashMapGetFunc := functions.NewFunction(
		"hash-map-get",
		registry.CategoryHashMap,
		2,
		"Get value from hash map by key: (hash-map-get map \"key\")",
		p.evalHashMapGet,
	)
	if err := reg.Register(hashMapGetFunc); err != nil {
		return err
	}

	// hash-map-put function
	hashMapPutFunc := functions.NewFunction(
		"hash-map-put",
		registry.CategoryHashMap,
		3,
		"Add/update key-value pair in hash map: (hash-map-put map \"key\" \"value\")",
		p.evalHashMapPut,
	)
	if err := reg.Register(hashMapPutFunc); err != nil {
		return err
	}

	// hash-map-remove function
	hashMapRemoveFunc := functions.NewFunction(
		"hash-map-remove",
		registry.CategoryHashMap,
		2,
		"Remove key from hash map: (hash-map-remove map \"key\")",
		p.evalHashMapRemove,
	)
	if err := reg.Register(hashMapRemoveFunc); err != nil {
		return err
	}

	// hash-map-contains? function
	hashMapContainsFunc := functions.NewFunction(
		"hash-map-contains?",
		registry.CategoryHashMap,
		2,
		"Check if hash map contains key: (hash-map-contains? map \"key\")",
		p.evalHashMapContains,
	)
	if err := reg.Register(hashMapContainsFunc); err != nil {
		return err
	}

	// hash-map-keys function
	hashMapKeysFunc := functions.NewFunction(
		"hash-map-keys",
		registry.CategoryHashMap,
		1,
		"Get all keys from hash map: (hash-map-keys map)",
		p.evalHashMapKeys,
	)
	if err := reg.Register(hashMapKeysFunc); err != nil {
		return err
	}

	// hash-map-values function
	hashMapValuesFunc := functions.NewFunction(
		"hash-map-values",
		registry.CategoryHashMap,
		1,
		"Get all values from hash map: (hash-map-values map)",
		p.evalHashMapValues,
	)
	if err := reg.Register(hashMapValuesFunc); err != nil {
		return err
	}

	// hash-map-size function
	hashMapSizeFunc := functions.NewFunction(
		"hash-map-size",
		registry.CategoryHashMap,
		1,
		"Get size of hash map: (hash-map-size map)",
		p.evalHashMapSize,
	)
	if err := reg.Register(hashMapSizeFunc); err != nil {
		return err
	}

	// hash-map-empty? function
	hashMapEmptyFunc := functions.NewFunction(
		"hash-map-empty?",
		registry.CategoryHashMap,
		1,
		"Check if hash map is empty: (hash-map-empty? map)",
		p.evalHashMapEmpty,
	)
	if err := reg.Register(hashMapEmptyFunc); err != nil {
		return err
	}

	return nil
}

// Helper function to convert a key value (string or keyword) to its string representation for hash map storage
func (p *HashMapPlugin) convertKeyToString(keyValue types.Value) (string, error) {
	switch kv := keyValue.(type) {
	case types.StringValue:
		return string(kv), nil
	case types.KeywordValue:
		return string(kv), nil
	default:
		return "", fmt.Errorf("hash map keys must be strings or keywords, got %T", keyValue)
	}
}

// evalHashMap creates a new hash map from key-value pairs
func (p *HashMapPlugin) evalHashMap(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("hash-map requires an even number of arguments")
	}

	elements := make(map[string]types.Value)

	for i := 0; i < len(args); i += 2 {
		keyExpr := args[i]
		valueExpr := args[i+1]

		keyValue, err := evaluator.Eval(keyExpr)
		if err != nil {
			return nil, err
		}

		valueValue, err := evaluator.Eval(valueExpr)
		if err != nil {
			return nil, err
		}

		keyStr, err := p.convertKeyToString(keyValue)
		if err != nil {
			return nil, err
		}

		elements[keyStr] = valueValue
	}

	return &types.HashMapValue{Elements: elements}, nil
}

// evalHashMapGet retrieves a value from a hash map by key
func (p *HashMapPlugin) evalHashMapGet(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("hash-map-get requires exactly 2 arguments, got %d", len(args))
	}

	hashMapValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	keyValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("expected hash-map as first argument to hash-map-get, got %T", hashMapValue)
	}

	keyStr, err := p.convertKeyToString(keyValue)
	if err != nil {
		return nil, err
	}

	value, exists := hashMap.Elements[keyStr]
	if !exists {
		return &types.NilValue{}, nil
	}

	return value, nil
}

// evalHashMapPut adds or updates a key-value pair in a hash map (returns new hash map)
func (p *HashMapPlugin) evalHashMapPut(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("hash-map-put requires exactly 3 arguments, got %d", len(args))
	}

	hashMapValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	keyValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	valueValue, err := evaluator.Eval(args[2])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("expected hash-map as first argument to hash-map-put, got %T", hashMapValue)
	}

	keyStr, err := p.convertKeyToString(keyValue)
	if err != nil {
		return nil, err
	}

	// Create a new hash map with the updated value
	newElements := make(map[string]types.Value)
	for k, v := range hashMap.Elements {
		newElements[k] = v
	}
	newElements[keyStr] = valueValue

	return &types.HashMapValue{Elements: newElements}, nil
}

// evalHashMapRemove removes a key-value pair from a hash map (returns new hash map)
func (p *HashMapPlugin) evalHashMapRemove(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("hash-map-remove requires exactly 2 arguments, got %d", len(args))
	}

	hashMapValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	keyValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("expected hash-map as first argument to hash-map-remove, got %T", hashMapValue)
	}

	keyStr, err := p.convertKeyToString(keyValue)
	if err != nil {
		return nil, err
	}

	// Create a new hash map without the specified key
	newElements := make(map[string]types.Value)
	for k, v := range hashMap.Elements {
		if k != keyStr {
			newElements[k] = v
		}
	}

	return &types.HashMapValue{Elements: newElements}, nil
}

// evalHashMapContains checks if a hash map contains a specific key
func (p *HashMapPlugin) evalHashMapContains(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("hash-map-contains? requires exactly 2 arguments, got %d", len(args))
	}

	hashMapValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	keyValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("expected hash-map as first argument to hash-map-contains?, got %T", hashMapValue)
	}

	keyStr, err := p.convertKeyToString(keyValue)
	if err != nil {
		return nil, err
	}

	_, exists := hashMap.Elements[keyStr]
	return types.BooleanValue(exists), nil
}

// evalHashMapKeys returns a list of all keys in a hash map
func (p *HashMapPlugin) evalHashMapKeys(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hash-map-keys requires exactly 1 argument, got %d", len(args))
	}

	hashMapValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("expected hash-map as argument to hash-map-keys, got %T", hashMapValue)
	}

	var keys []types.Value
	for key := range hashMap.Elements {
		keys = append(keys, types.StringValue(key))
	}

	return &types.ListValue{Elements: keys}, nil
}

// evalHashMapValues returns a list of all values in a hash map
func (p *HashMapPlugin) evalHashMapValues(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hash-map-values requires exactly 1 argument, got %d", len(args))
	}

	hashMapValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("expected hash-map as argument to hash-map-values, got %T", hashMapValue)
	}

	var values []types.Value
	for _, value := range hashMap.Elements {
		values = append(values, value)
	}

	return &types.ListValue{Elements: values}, nil
}

// evalHashMapSize returns the number of key-value pairs in a hash map
func (p *HashMapPlugin) evalHashMapSize(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hash-map-size requires exactly 1 argument, got %d", len(args))
	}

	hashMapValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("expected hash-map as argument to hash-map-size, got %T", hashMapValue)
	}

	return types.NumberValue(float64(len(hashMap.Elements))), nil
}

// evalHashMapEmpty checks if a hash map is empty
func (p *HashMapPlugin) evalHashMapEmpty(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hash-map-empty? requires exactly 1 argument, got %d", len(args))
	}

	hashMapValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("expected hash-map as argument to hash-map-empty?, got %T", hashMapValue)
	}

	return types.BooleanValue(len(hashMap.Elements) == 0), nil
}
