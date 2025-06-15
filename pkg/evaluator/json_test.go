package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestJsonParse(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		json     string
		expected types.Value
	}{
		{
			name:     "parse null",
			json:     "null",
			expected: &types.NilValue{},
		},
		{
			name:     "parse boolean true",
			json:     "true",
			expected: types.BooleanValue(true),
		},
		{
			name:     "parse boolean false",
			json:     "false",
			expected: types.BooleanValue(false),
		},
		{
			name:     "parse number",
			json:     "42",
			expected: types.NumberValue(42),
		},
		{
			name:     "parse string",
			json:     `"hello"`,
			expected: types.StringValue("hello"),
		},
		{
			name: "parse array",
			json: `[1, 2, 3]`,
			expected: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			}},
		},
		{
			name: "parse object",
			json: `{"name": "Alice", "age": 30}`,
			expected: &types.HashMapValue{Elements: map[string]types.Value{
				"name": types.StringValue("Alice"),
				"age":  types.NumberValue(30),
			}},
		},
		{
			name: "parse nested object",
			json: `{"user": {"name": "Bob", "active": true}, "count": 5}`,
			expected: &types.HashMapValue{Elements: map[string]types.Value{
				"user": &types.HashMapValue{Elements: map[string]types.Value{
					"name":   types.StringValue("Bob"),
					"active": types.BooleanValue(true),
				}},
				"count": types.NumberValue(5),
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-parse"},
					&types.StringExpr{Value: tt.json},
				},
			}

			result, err := evaluator.Eval(expr)
			if err != nil {
				t.Fatalf("json-parse failed: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestJsonStringify(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		value    types.Expr
		expected string
	}{
		{
			name:     "stringify null",
			value:    &types.SymbolExpr{Name: "nil"},
			expected: "null",
		},
		{
			name:     "stringify boolean",
			value:    &types.BooleanExpr{Value: true},
			expected: "true",
		},
		{
			name:     "stringify number",
			value:    &types.NumberExpr{Value: 42},
			expected: "42",
		},
		{
			name:     "stringify string",
			value:    &types.StringExpr{Value: "hello"},
			expected: `"hello"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-stringify"},
					tt.value,
				},
			}

			result, err := evaluator.Eval(expr)
			if err != nil {
				t.Fatalf("json-stringify failed: %v", err)
			}

			resultStr, ok := result.(types.StringValue)
			if !ok {
				t.Fatalf("expected StringValue, got %T", result)
			}

			if string(resultStr) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(resultStr))
			}
		})
	}
}

func TestJsonStringifyComplex(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test with hash map
	hashMapExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "hash-map"},
			&types.StringExpr{Value: "name"},
			&types.StringExpr{Value: "Alice"},
			&types.StringExpr{Value: "age"},
			&types.NumberExpr{Value: 30},
		},
	}

	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "json-stringify"},
			hashMapExpr,
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("json-stringify failed: %v", err)
	}

	resultStr, ok := result.(types.StringValue)
	if !ok {
		t.Fatalf("expected StringValue, got %T", result)
	}

	// Parse it back to verify it's valid JSON
	parseExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "json-parse"},
			&types.StringExpr{Value: string(resultStr)},
		},
	}

	parsed, err := evaluator.Eval(parseExpr)
	if err != nil {
		t.Fatalf("failed to parse back stringified JSON: %v", err)
	}

	// Should be a hash map with the expected values
	parsedMap, ok := parsed.(*types.HashMapValue)
	if !ok {
		t.Fatalf("expected HashMapValue, got %T", parsed)
	}

	if name, exists := parsedMap.Elements["name"]; !exists || name != types.StringValue("Alice") {
		t.Errorf("expected name=Alice, got %v (exists: %v)", name, exists)
	}

	if age, exists := parsedMap.Elements["age"]; !exists || age != types.NumberValue(30) {
		t.Errorf("expected age=30, got %v (exists: %v)", age, exists)
	}
}

func TestJsonStringifyPretty(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create a hash map
	hashMapExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "hash-map"},
			&types.StringExpr{Value: "name"},
			&types.StringExpr{Value: "Alice"},
			&types.StringExpr{Value: "age"},
			&types.NumberExpr{Value: 30},
		},
	}

	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "json-stringify-pretty"},
			hashMapExpr,
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("json-stringify-pretty failed: %v", err)
	}

	resultStr, ok := result.(types.StringValue)
	if !ok {
		t.Fatalf("expected StringValue, got %T", result)
	}

	// Should contain newlines and proper indentation
	jsonStr := string(resultStr)
	if len(jsonStr) < 10 { // Should be longer than compact JSON
		t.Errorf("pretty JSON seems too short: %s", jsonStr)
	}

	// Should contain newlines
	containsNewline := false
	for _, char := range jsonStr {
		if char == '\n' {
			containsNewline = true
			break
		}
	}
	if !containsNewline {
		t.Errorf("pretty JSON should contain newlines: %s", jsonStr)
	}
}

func TestJsonPath(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	jsonStr := `{
		"user": {
			"name": "Alice",
			"age": 30,
			"addresses": [
				{"type": "home", "city": "Boston"},
				{"type": "work", "city": "New York"}
			]
		},
		"active": true
	}`

	tests := []struct {
		name     string
		path     string
		expected types.Value
	}{
		{
			name:     "get user name",
			path:     "user.name",
			expected: types.StringValue("Alice"),
		},
		{
			name:     "get user age",
			path:     "user.age",
			expected: types.NumberValue(30),
		},
		{
			name:     "get active status",
			path:     "active",
			expected: types.BooleanValue(true),
		},
		{
			name:     "get first address",
			path:     "user.addresses.0.city",
			expected: types.StringValue("Boston"),
		},
		{
			name:     "get second address type",
			path:     "user.addresses.1.type",
			expected: types.StringValue("work"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-path"},
					&types.StringExpr{Value: jsonStr},
					&types.StringExpr{Value: tt.path},
				},
			}

			result, err := evaluator.Eval(expr)
			if err != nil {
				t.Fatalf("json-path failed: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestJsonErrors(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "json-parse with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-parse"},
				},
			},
		},
		{
			name: "json-parse with non-string argument",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-parse"},
					&types.NumberExpr{Value: 123},
				},
			},
		},
		{
			name: "json-parse with invalid JSON",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-parse"},
					&types.StringExpr{Value: "{invalid json}"},
				},
			},
		},
		{
			name: "json-stringify with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-stringify"},
				},
			},
		},
		{
			name: "json-path with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-path"},
					&types.StringExpr{Value: "{}"},
				},
			},
		},
		{
			name: "json-path with invalid JSON",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-path"},
					&types.StringExpr{Value: "{invalid}"},
					&types.StringExpr{Value: "key"},
				},
			},
		},
		{
			name: "json-path with invalid path",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "json-path"},
					&types.StringExpr{Value: `{"name": "Alice"}`},
					&types.StringExpr{Value: "nonexistent"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluator.Eval(tt.expr)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}
