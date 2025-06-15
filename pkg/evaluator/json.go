// Package evaluator_json contains JSON conversion functionality
package evaluator

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// evalJsonParse converts a JSON string to a Lisp value (hash map, list, etc.)
func (e *Evaluator) evalJsonParse(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("json-parse requires exactly 1 argument (JSON string), got %d", len(args))
	}

	// Evaluate the JSON string
	jsonValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	jsonStr, ok := jsonValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("json-parse argument must be a string, got %T", jsonValue)
	}

	// Parse the JSON
	var jsonData interface{}
	err = json.Unmarshal([]byte(string(jsonStr)), &jsonData)
	if err != nil {
		return nil, fmt.Errorf("json-parse failed: %v", err)
	}

	// Convert JSON data to Lisp value
	return convertJsonToLisp(jsonData)
}

// evalJsonStringify converts a Lisp value to a JSON string
func (e *Evaluator) evalJsonStringify(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("json-stringify requires exactly 1 argument, got %d", len(args))
	}

	// Evaluate the argument
	value, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Convert Lisp value to JSON-compatible Go value
	jsonData, err := convertLispToJson(value)
	if err != nil {
		return nil, fmt.Errorf("json-stringify failed: %v", err)
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("json-stringify marshal failed: %v", err)
	}

	return types.StringValue(string(jsonBytes)), nil
}

// evalJsonStringifyPretty converts a Lisp value to a pretty-printed JSON string
func (e *Evaluator) evalJsonStringifyPretty(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("json-stringify-pretty requires exactly 1 argument, got %d", len(args))
	}

	// Evaluate the argument
	value, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Convert Lisp value to JSON-compatible Go value
	jsonData, err := convertLispToJson(value)
	if err != nil {
		return nil, fmt.Errorf("json-stringify-pretty failed: %v", err)
	}

	// Marshal to pretty JSON
	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json-stringify-pretty marshal failed: %v", err)
	}

	return types.StringValue(string(jsonBytes)), nil
}

// convertJsonToLisp converts Go JSON data to Lisp values
func convertJsonToLisp(data interface{}) (types.Value, error) {
	switch v := data.(type) {
	case nil:
		return &types.NilValue{}, nil
	case bool:
		return types.BooleanValue(v), nil
	case float64:
		return types.NumberValue(v), nil
	case string:
		return types.StringValue(v), nil
	case []interface{}:
		// Convert JSON array to Lisp list
		elements := make([]types.Value, len(v))
		for i, item := range v {
			elem, err := convertJsonToLisp(item)
			if err != nil {
				return nil, err
			}
			elements[i] = elem
		}
		return &types.ListValue{Elements: elements}, nil
	case map[string]interface{}:
		// Convert JSON object to Lisp hash map
		hashMap := &types.HashMapValue{
			Elements: make(map[string]types.Value),
		}
		for key, value := range v {
			lispValue, err := convertJsonToLisp(value)
			if err != nil {
				return nil, err
			}
			hashMap.Elements[key] = lispValue
		}
		return hashMap, nil
	default:
		return nil, fmt.Errorf("unsupported JSON type: %T", data)
	}
}

// convertLispToJson converts Lisp values to Go JSON-compatible values
func convertLispToJson(value types.Value) (interface{}, error) {
	switch v := value.(type) {
	case *types.NilValue:
		return nil, nil
	case types.BooleanValue:
		return bool(v), nil
	case types.NumberValue:
		return float64(v), nil
	case types.StringValue:
		return string(v), nil
	case types.KeywordValue:
		// Convert keywords to strings in JSON
		return ":" + string(v), nil
	case *types.ListValue:
		// Convert Lisp list to JSON array
		result := make([]interface{}, len(v.Elements))
		for i, elem := range v.Elements {
			jsonElem, err := convertLispToJson(elem)
			if err != nil {
				return nil, err
			}
			result[i] = jsonElem
		}
		return result, nil
	case *types.HashMapValue:
		// Convert Lisp hash map to JSON object
		result := make(map[string]interface{})
		for key, val := range v.Elements {
			jsonVal, err := convertLispToJson(val)
			if err != nil {
				return nil, err
			}
			result[key] = jsonVal
		}
		return result, nil
	case *types.BigNumberValue:
		// Convert big numbers to strings in JSON to preserve precision
		return v.String(), nil
	default:
		// For unsupported types, convert to string representation
		return value.String(), nil
	}
}

// evalJsonPath extracts a value from JSON using a simple path notation
func (e *Evaluator) evalJsonPath(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("json-path requires exactly 2 arguments (JSON string, path), got %d", len(args))
	}

	// Evaluate the JSON string
	jsonValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	jsonStr, ok := jsonValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("json-path first argument must be a string, got %T", jsonValue)
	}

	// Evaluate the path
	pathValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	pathStr, ok := pathValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("json-path second argument must be a string, got %T", pathValue)
	}

	// Parse the JSON
	var jsonData interface{}
	err = json.Unmarshal([]byte(string(jsonStr)), &jsonData)
	if err != nil {
		return nil, fmt.Errorf("json-path failed to parse JSON: %v", err)
	}

	// Navigate the path (simple dot notation like "user.name" or "users.0.name")
	result, err := navigateJsonPath(jsonData, string(pathStr))
	if err != nil {
		return nil, fmt.Errorf("json-path navigation failed: %v", err)
	}

	// Convert result to Lisp value
	return convertJsonToLisp(result)
}

// navigateJsonPath navigates a JSON structure using dot notation
func navigateJsonPath(data interface{}, path string) (interface{}, error) {
	if path == "" {
		return data, nil
	}

	// Split path by dots
	parts := []string{}
	current := ""
	for _, char := range path {
		if char == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}

	// Navigate through the path
	current_data := data
	for _, part := range parts {
		switch v := current_data.(type) {
		case map[string]interface{}:
			// Object access
			if value, exists := v[part]; exists {
				current_data = value
			} else {
				return nil, fmt.Errorf("key '%s' not found", part)
			}
		case []interface{}:
			// Array access
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid array index '%s'", part)
			}
			if index < 0 || index >= len(v) {
				return nil, fmt.Errorf("array index %d out of bounds", index)
			}
			current_data = v[index]
		default:
			return nil, fmt.Errorf("cannot navigate path '%s' through non-object/array", part)
		}
	}

	return current_data, nil
}
