// Package json implements JSON functions as a plugin
package json

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// JSONPlugin implements JSON functions
type JSONPlugin struct {
	*plugins.BasePlugin
}

// NewJSONPlugin creates a new JSON plugin
func NewJSONPlugin() *JSONPlugin {
	return &JSONPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"json",
			"1.0.0",
			"JSON processing functions (json-parse, json-stringify, json-path, etc.)",
			[]string{}, // No dependencies
		),
	}
}

// Functions returns the list of functions provided by this plugin
func (p *JSONPlugin) Functions() []string {
	return []string{
		"json-parse", "json-stringify", "json-stringify-pretty", "json-path",
	}
}

// RegisterFunctions registers all JSON functions with the registry
func (p *JSONPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// json-parse function
	jsonParseFunc := functions.NewFunction(
		"json-parse",
		registry.CategoryJSON,
		1,
		"Parse JSON string to Lisp data: (json-parse \"{\\\"key\\\": \\\"value\\\"}\")",
		p.evalJsonParse,
	)
	if err := reg.Register(jsonParseFunc); err != nil {
		return err
	}

	// json-stringify function
	jsonStringifyFunc := functions.NewFunction(
		"json-stringify",
		registry.CategoryJSON,
		1,
		"Convert Lisp data to JSON string: (json-stringify (hash-map \"key\" \"value\"))",
		p.evalJsonStringify,
	)
	if err := reg.Register(jsonStringifyFunc); err != nil {
		return err
	}

	// json-stringify-pretty function
	jsonStringifyPrettyFunc := functions.NewFunction(
		"json-stringify-pretty",
		registry.CategoryJSON,
		1,
		"Convert Lisp data to pretty JSON string: (json-stringify-pretty data)",
		p.evalJsonStringifyPretty,
	)
	if err := reg.Register(jsonStringifyPrettyFunc); err != nil {
		return err
	}

	// json-path function
	jsonPathFunc := functions.NewFunction(
		"json-path",
		registry.CategoryJSON,
		2,
		"Extract value using JSON path: (json-path json-string \"path.to.field\")",
		p.evalJsonPath,
	)
	if err := reg.Register(jsonPathFunc); err != nil {
		return err
	}

	return nil
}

// Helper function to convert JSON interface{} to Lisp Value
func (p *JSONPlugin) jsonToLispValue(jsonValue interface{}) types.Value {
	if jsonValue == nil {
		return &types.NilValue{}
	}

	switch v := jsonValue.(type) {
	case bool:
		return types.BooleanValue(v)
	case float64:
		return types.NumberValue(v)
	case string:
		return types.StringValue(v)
	case []interface{}:
		var elements []types.Value
		for _, elem := range v {
			elements = append(elements, p.jsonToLispValue(elem))
		}
		return &types.ListValue{Elements: elements}
	case map[string]interface{}:
		elements := make(map[string]types.Value)
		for key, value := range v {
			elements[key] = p.jsonToLispValue(value)
		}
		return &types.HashMapValue{Elements: elements}
	default:
		// Fallback to string representation
		return types.StringValue(fmt.Sprintf("%v", v))
	}
}

// Helper function to convert Lisp Value to JSON-compatible interface{}
func (p *JSONPlugin) lispValueToJSON(value types.Value) interface{} {
	switch v := value.(type) {
	case types.BooleanValue:
		return bool(v)
	case types.NumberValue:
		return float64(v)
	case types.StringValue:
		return string(v)
	case types.KeywordValue:
		return string(v)
	case *types.NilValue:
		return nil
	case *types.ListValue:
		var result []interface{}
		for _, elem := range v.Elements {
			result = append(result, p.lispValueToJSON(elem))
		}
		return result
	case *types.HashMapValue:
		result := make(map[string]interface{})
		for key, val := range v.Elements {
			result[key] = p.lispValueToJSON(val)
		}
		return result
	default:
		// Fallback to string representation
		return value.String()
	}
}

// evalJsonParse parses a JSON string into Lisp data structures
func (p *JSONPlugin) evalJsonParse(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("json-parse requires exactly 1 argument, got %d", len(args))
	}

	jsonStringValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	jsonString, ok := jsonStringValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("json-parse requires a string argument, got %T", jsonStringValue)
	}

	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return p.jsonToLispValue(jsonData), nil
}

// evalJsonStringify converts Lisp data to a JSON string
func (p *JSONPlugin) evalJsonStringify(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("json-stringify requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	jsonData := p.lispValueToJSON(value)
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to stringify JSON: %v", err)
	}

	return types.StringValue(string(jsonBytes)), nil
}

// evalJsonStringifyPretty converts Lisp data to a pretty-printed JSON string
func (p *JSONPlugin) evalJsonStringifyPretty(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("json-stringify-pretty requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	jsonData := p.lispValueToJSON(value)
	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to stringify pretty JSON: %v", err)
	}

	return types.StringValue(string(jsonBytes)), nil
}

// evalJsonPath extracts a value from JSON using a dot-separated path
func (p *JSONPlugin) evalJsonPath(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("json-path requires exactly 2 arguments, got %d", len(args))
	}

	jsonStringValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	jsonString, ok := jsonStringValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("json-path first argument must be a string, got %T", jsonStringValue)
	}

	pathValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	pathString, ok := pathValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("json-path second argument must be a string, got %T", pathValue)
	}

	// Parse JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Navigate the path
	result, err := p.navigateJSONPath(jsonData, string(pathString))
	if err != nil {
		return nil, err
	}

	return p.jsonToLispValue(result), nil
}

// Helper function to navigate a JSON path
func (p *JSONPlugin) navigateJSONPath(data interface{}, path string) (interface{}, error) {
	if path == "" {
		return data, nil
	}

	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Check if this is an array index
		if idx, err := strconv.Atoi(part); err == nil {
			// Array access
			switch arr := current.(type) {
			case []interface{}:
				if idx < 0 || idx >= len(arr) {
					return nil, fmt.Errorf("array index %d out of bounds", idx)
				}
				current = arr[idx]
			default:
				return nil, fmt.Errorf("cannot index non-array with %d", idx)
			}
		} else {
			// Object property access
			switch obj := current.(type) {
			case map[string]interface{}:
				value, exists := obj[part]
				if !exists {
					return nil, fmt.Errorf("property %s not found", part)
				}
				current = value
			default:
				return nil, fmt.Errorf("cannot access property %s on non-object", part)
			}
		}
	}

	return current, nil
}
