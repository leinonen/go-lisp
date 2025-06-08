// Package evaluator_hashmaps contains hash map functionality
package evaluator

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Hash map operations

// convertKeyToString converts a key value (string or keyword) to its string representation for hash map storage
func convertKeyToString(keyValue types.Value) (string, error) {
	switch key := keyValue.(type) {
	case types.StringValue:
		return string(key), nil
	case types.KeywordValue:
		return ":" + string(key), nil
	default:
		return "", fmt.Errorf("hash map keys must be strings or keywords, got %T", keyValue)
	}
}

// evalHashMap creates a new hash map from key-value pairs
func (e *Evaluator) evalHashMap(args []types.Expr) (types.Value, error) {
	// Check that we have an even number of arguments (key-value pairs)
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("hash-map requires an even number of arguments (key-value pairs), got %d", len(args))
	}

	hashMap := &types.HashMapValue{
		Elements: make(map[string]types.Value),
	}

	// Process key-value pairs
	for i := 0; i < len(args); i += 2 {
		// Evaluate the key
		keyValue, err := e.Eval(args[i])
		if err != nil {
			return nil, err
		}

		// Key must be a string
		keyStr, err := convertKeyToString(keyValue)
		if err != nil {
			return nil, err
		}

		// Evaluate the value
		value, err := e.Eval(args[i+1])
		if err != nil {
			return nil, err
		}

		hashMap.Elements[keyStr] = value
	}

	return hashMap, nil
}

// evalHashMapGet retrieves a value from a hash map by key
func (e *Evaluator) evalHashMapGet(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("hash-map-get requires exactly 2 arguments (hash-map and key), got %d", len(args))
	}

	// Evaluate the hash map
	hashMapValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("hash-map-get first argument must be a hash map, got %T", hashMapValue)
	}

	// Evaluate the key
	keyValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	keyStr, err := convertKeyToString(keyValue)
	if err != nil {
		return nil, err
	}

	// Look up the value
	value, exists := hashMap.Elements[keyStr]
	if !exists {
		return &types.NilValue{}, nil
	}

	return value, nil
}

// evalHashMapPut adds or updates a key-value pair in a hash map (returns new hash map)
func (e *Evaluator) evalHashMapPut(args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("hash-map-put requires exactly 3 arguments (hash-map, key, and value), got %d", len(args))
	}

	// Evaluate the hash map
	hashMapValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("hash-map-put first argument must be a hash map, got %T", hashMapValue)
	}

	// Evaluate the key
	keyValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	keyStr, err := convertKeyToString(keyValue)
	if err != nil {
		return nil, err
	}

	// Evaluate the value
	value, err := e.Eval(args[2])
	if err != nil {
		return nil, err
	}

	// Create a new hash map with the updated key-value pair (immutable)
	newHashMap := &types.HashMapValue{
		Elements: make(map[string]types.Value),
	}

	// Copy existing elements
	for k, v := range hashMap.Elements {
		newHashMap.Elements[k] = v
	}

	// Add or update the new key-value pair
	newHashMap.Elements[keyStr] = value

	return newHashMap, nil
}

// evalHashMapRemove removes a key-value pair from a hash map (returns new hash map)
func (e *Evaluator) evalHashMapRemove(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("hash-map-remove requires exactly 2 arguments (hash-map and key), got %d", len(args))
	}

	// Evaluate the hash map
	hashMapValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("hash-map-remove first argument must be a hash map, got %T", hashMapValue)
	}

	// Evaluate the key
	keyValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	keyStr, err := convertKeyToString(keyValue)
	if err != nil {
		return nil, err
	}

	// Create a new hash map without the specified key (immutable)
	newHashMap := &types.HashMapValue{
		Elements: make(map[string]types.Value),
	}

	// Copy existing elements except the one to remove
	for k, v := range hashMap.Elements {
		if k != keyStr {
			newHashMap.Elements[k] = v
		}
	}

	return newHashMap, nil
}

// evalHashMapContains checks if a hash map contains a specific key
func (e *Evaluator) evalHashMapContains(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("hash-map-contains? requires exactly 2 arguments (hash-map and key), got %d", len(args))
	}

	// Evaluate the hash map
	hashMapValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("hash-map-contains? first argument must be a hash map, got %T", hashMapValue)
	}

	// Evaluate the key
	keyValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	keyStr, err := convertKeyToString(keyValue)
	if err != nil {
		return nil, err
	}

	// Check if key exists
	_, exists := hashMap.Elements[keyStr]
	return types.BooleanValue(exists), nil
}

// evalHashMapKeys returns a list of all keys in a hash map
func (e *Evaluator) evalHashMapKeys(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hash-map-keys requires exactly 1 argument (hash-map), got %d", len(args))
	}

	// Evaluate the hash map
	hashMapValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("hash-map-keys first argument must be a hash map, got %T", hashMapValue)
	}

	// Collect all keys
	keys := make([]types.Value, 0, len(hashMap.Elements))
	for key := range hashMap.Elements {
		keys = append(keys, types.StringValue(key))
	}

	return &types.ListValue{Elements: keys}, nil
}

// evalHashMapValues returns a list of all values in a hash map
func (e *Evaluator) evalHashMapValues(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hash-map-values requires exactly 1 argument (hash-map), got %d", len(args))
	}

	// Evaluate the hash map
	hashMapValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("hash-map-values first argument must be a hash map, got %T", hashMapValue)
	}

	// Collect all values
	values := make([]types.Value, 0, len(hashMap.Elements))
	for _, value := range hashMap.Elements {
		values = append(values, value)
	}

	return &types.ListValue{Elements: values}, nil
}

// evalHashMapSize returns the number of key-value pairs in a hash map
func (e *Evaluator) evalHashMapSize(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hash-map-size requires exactly 1 argument (hash-map), got %d", len(args))
	}

	// Evaluate the hash map
	hashMapValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("hash-map-size first argument must be a hash map, got %T", hashMapValue)
	}

	return types.NumberValue(len(hashMap.Elements)), nil
}

// evalHashMapEmpty checks if a hash map is empty
func (e *Evaluator) evalHashMapEmpty(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hash-map-empty? requires exactly 1 argument (hash-map), got %d", len(args))
	}

	// Evaluate the hash map
	hashMapValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	hashMap, ok := hashMapValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("hash-map-empty? first argument must be a hash map, got %T", hashMapValue)
	}

	return types.BooleanValue(len(hashMap.Elements) == 0), nil
}
