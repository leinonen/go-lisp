package json

import (
	"strings"
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
	return nil, nil // Not needed for JSON tests
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

func TestJSONPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewJSONPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"json-parse", "json-stringify", "json-stringify-pretty", "json-path"}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestJSONPlugin_JsonParse(t *testing.T) {
	plugin := NewJSONPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		json     string
		expected string
	}{
		{
			name:     "simple object",
			json:     `{"name": "John", "age": 30}`,
			expected: "hash-map", // Should return a hash map
		},
		{
			name:     "array",
			json:     `[1, 2, 3]`,
			expected: "list", // Should return a list
		},
		{
			name:     "string",
			json:     `"hello"`,
			expected: "hello",
		},
		{
			name:     "number",
			json:     `42`,
			expected: "42",
		},
		{
			name:     "boolean true",
			json:     `true`,
			expected: "true",
		},
		{
			name:     "boolean false",
			json:     `false`,
			expected: "false",
		},
		{
			name:     "null",
			json:     `null`,
			expected: "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{&types.StringExpr{Value: tt.json}}
			result, err := plugin.evalJsonParse(evaluator, args)
			if err != nil {
				t.Fatalf("evalJsonParse failed: %v", err)
			}

			if result == nil {
				t.Fatal("Result is nil")
			}

			switch tt.expected {
			case "hash-map":
				if _, ok := result.(*types.HashMapValue); !ok {
					t.Errorf("Expected hash map, got %T", result)
				}
			case "list":
				if _, ok := result.(*types.ListValue); !ok {
					t.Errorf("Expected list, got %T", result)
				}
			case "nil":
				if _, ok := result.(*types.NilValue); !ok {
					t.Errorf("Expected nil, got %T", result)
				}
			default:
				if result.String() != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

func TestJSONPlugin_JsonStringify(t *testing.T) {
	plugin := NewJSONPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name     string
		value    types.Value
		expected string
	}{
		{
			name:     "string",
			value:    types.StringValue("hello"),
			expected: `"hello"`,
		},
		{
			name:     "number",
			value:    types.NumberValue(42),
			expected: "42",
		},
		{
			name:     "boolean true",
			value:    types.BooleanValue(true),
			expected: "true",
		},
		{
			name:     "boolean false",
			value:    types.BooleanValue(false),
			expected: "false",
		},
		{
			name:     "nil",
			value:    &types.NilValue{},
			expected: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{wrapValue(tt.value)}
			result, err := plugin.evalJsonStringify(evaluator, args)
			if err != nil {
				t.Fatalf("evalJsonStringify failed: %v", err)
			}

			if result == nil {
				t.Fatal("Result is nil")
			}

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestJSONPlugin_JsonStringifyPretty(t *testing.T) {
	plugin := NewJSONPlugin()
	evaluator := newMockEvaluator()

	// Test with a simple hash map
	hashMap := &types.HashMapValue{Elements: make(map[string]types.Value)}
	hashMap.Elements["name"] = types.StringValue("John")
	hashMap.Elements["age"] = types.NumberValue(30)

	args := []types.Expr{wrapValue(hashMap)}
	result, err := plugin.evalJsonStringifyPretty(evaluator, args)
	if err != nil {
		t.Fatalf("evalJsonStringifyPretty failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	jsonStr := result.String()
	// Pretty JSON should contain newlines and indentation
	if !strings.Contains(jsonStr, "\n") {
		t.Error("Expected pretty JSON to contain newlines")
	}
}

func TestJSONPlugin_JsonPath(t *testing.T) {
	plugin := NewJSONPlugin()
	evaluator := newMockEvaluator()

	// Create a nested JSON structure
	jsonStr := `{"user": {"name": "John", "age": 30, "address": {"city": "New York"}}}`

	tests := []struct {
		name     string
		json     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			json:     jsonStr,
			path:     "user.name",
			expected: "John",
		},
		{
			name:     "nested path",
			json:     jsonStr,
			path:     "user.address.city",
			expected: "New York",
		},
		{
			name:     "number path",
			json:     jsonStr,
			path:     "user.age",
			expected: "30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []types.Expr{
				&types.StringExpr{Value: tt.json},
				&types.StringExpr{Value: tt.path},
			}
			result, err := plugin.evalJsonPath(evaluator, args)
			if err != nil {
				t.Fatalf("evalJsonPath failed: %v", err)
			}

			if result == nil {
				t.Fatal("Result is nil")
			}

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestJSONPlugin_ErrorCases(t *testing.T) {
	plugin := NewJSONPlugin()
	evaluator := newMockEvaluator()

	// Test invalid JSON
	args := []types.Expr{&types.StringExpr{Value: "invalid json"}}
	_, err := plugin.evalJsonParse(evaluator, args)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Test wrong argument count
	_, err = plugin.evalJsonParse(evaluator, []types.Expr{})
	if err == nil {
		t.Error("Expected error for missing arguments")
	}

	// Test wrong argument type
	args = []types.Expr{&types.NumberExpr{Value: 42}}
	_, err = plugin.evalJsonParse(evaluator, args)
	if err == nil {
		t.Error("Expected error for non-string argument")
	}
}

func TestJSONPlugin_ComplexDataStructures(t *testing.T) {
	plugin := NewJSONPlugin()
	evaluator := newMockEvaluator()

	// Test parsing and then stringifying a complex structure
	complexJSON := `{
		"users": [
			{"name": "John", "age": 30},
			{"name": "Jane", "age": 25}
		],
		"active": true,
		"count": 2
	}`

	// Parse the JSON
	args := []types.Expr{&types.StringExpr{Value: complexJSON}}
	result, err := plugin.evalJsonParse(evaluator, args)
	if err != nil {
		t.Fatalf("Failed to parse complex JSON: %v", err)
	}

	// Stringify it back
	args = []types.Expr{wrapValue(result)}
	stringified, err := plugin.evalJsonStringify(evaluator, args)
	if err != nil {
		t.Fatalf("Failed to stringify complex data: %v", err)
	}

	// The result should be valid JSON
	if stringified == nil {
		t.Fatal("Stringified result is nil")
	}

	// Try to parse it again to verify it's valid JSON
	args = []types.Expr{wrapValue(stringified)}
	_, err = plugin.evalJsonParse(evaluator, args)
	if err != nil {
		t.Fatalf("Failed to re-parse stringified JSON: %v", err)
	}
}
