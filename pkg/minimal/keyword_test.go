package minimal

import (
	"testing"
)

func TestKeywordBasics(t *testing.T) {
	// Test keyword creation and string representation
	keyword := InternKeyword("test")
	expected := ":test"
	if keyword.String() != expected {
		t.Errorf("Expected keyword string %s, got %s", expected, keyword.String())
	}

	// Test keyword interning - same keyword should be identical
	keyword2 := InternKeyword("test")
	if keyword != keyword2 {
		t.Error("Keywords with same name should be identical after interning")
	}

	// Test different keywords are different
	keyword3 := InternKeyword("other")
	if keyword == keyword3 {
		t.Error("Different keywords should not be identical")
	}
}

func TestKeywordParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{":test", ":test"},
		{":name", ":name"},
		{":foo-bar", ":foo-bar"},
		{":123", ":123"},
		{":a", ":a"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			value, err := Parse(test.input)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test.input, err)
			}

			keyword, ok := value.(Keyword)
			if !ok {
				t.Fatalf("Expected Keyword, got %T", value)
			}

			if keyword.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, keyword.String())
			}
		})
	}
}

func TestKeywordInvalidParsing(t *testing.T) {
	tests := []string{
		":",  // empty keyword
		": ", // keyword with space
		":(", // keyword with special character
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			_, err := Parse(test)
			if err == nil {
				t.Errorf("Expected error parsing invalid keyword %s", test)
			}
		})
	}
}

func TestKeywordAsFunction(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)

	// Create a map with keyword keys
	hm := NewHashMap().
		PutByValue(InternKeyword("name"), String("Alice")).
		PutByValue(InternKeyword("age"), Number(30))

	// Test keyword function access
	keyword := InternKeyword("name")

	// Test with one argument (map only)
	result, err := keyword.Call([]Value{hm}, env)
	if err != nil {
		t.Fatalf("Error calling keyword as function: %v", err)
	}

	expected := String("Alice")
	if result.String() != expected.String() {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test with missing key
	missingKeyword := InternKeyword("missing")
	result, err = missingKeyword.Call([]Value{hm}, env)
	if err != nil {
		t.Fatalf("Error calling keyword with missing key: %v", err)
	}

	if _, ok := result.(Nil); !ok {
		t.Errorf("Expected Nil for missing key, got %T", result)
	}

	// Test with default value
	result, err = missingKeyword.Call([]Value{hm, String("default")}, env)
	if err != nil {
		t.Fatalf("Error calling keyword with default value: %v", err)
	}

	expected = String("default")
	if result.String() != expected.String() {
		t.Errorf("Expected default value %s, got %s", expected, result)
	}
}

func TestKeywordAccessOnNil(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)
	keyword := InternKeyword("test")

	// Test keyword access on nil should return default or nil
	result, err := keyword.Call([]Value{Nil{}}, env)
	if err != nil {
		t.Fatalf("Error calling keyword on nil: %v", err)
	}

	if _, ok := result.(Nil); !ok {
		t.Errorf("Expected Nil when accessing nil map, got %T", result)
	}

	// Test with default value
	result, err = keyword.Call([]Value{Nil{}, String("default")}, env)
	if err != nil {
		t.Fatalf("Error calling keyword on nil with default: %v", err)
	}

	expected := String("default")
	if result.String() != expected.String() {
		t.Errorf("Expected default value %s, got %s", expected, result)
	}
}

func TestKeywordFunctionErrors(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)
	keyword := InternKeyword("test")

	// Test with no arguments
	_, err := keyword.Call([]Value{}, env)
	if err == nil {
		t.Error("Expected error when calling keyword with no arguments")
	}

	// Test with too many arguments
	_, err = keyword.Call([]Value{Nil{}, String("default"), String("extra")}, env)
	if err == nil {
		t.Error("Expected error when calling keyword with too many arguments")
	}

	// Test with non-map argument
	_, err = keyword.Call([]Value{String("not-a-map")}, env)
	if err == nil {
		t.Error("Expected error when calling keyword on non-map")
	}
}

func TestKeywordInLispExpressions(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)

	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "keyword access on map",
			expr:     "(:name (hash-map :name \"Alice\" :age 30))",
			expected: "\"Alice\"",
		},
		{
			name:     "keyword access with default",
			expr:     "(:missing (hash-map :name \"Alice\") \"default\")",
			expected: "\"default\"",
		},
		{
			name:     "keyword access on nil",
			expr:     "(:test nil)",
			expected: "nil",
		},
		{
			name:     "nested keyword access",
			expr:     "(:city (:address (hash-map :address (hash-map :city \"NYC\"))))",
			expected: "\"NYC\"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := Parse(test.expr)
			if err != nil {
				t.Fatalf("Failed to parse expression %s: %v", test.expr, err)
			}

			result, err := Eval(expr, env)
			if err != nil {
				t.Fatalf("Failed to evaluate expression %s: %v", test.expr, err)
			}

			if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}

func TestKeywordHashMapIntegration(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)

	// Test hash-map creation with keywords
	expr := "(hash-map :name \"Alice\" :age 30 :city \"NYC\")"
	parsedExpr, err := Parse(expr)
	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	result, err := Eval(parsedExpr, env)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	hm, ok := result.(*HashMap)
	if !ok {
		t.Fatalf("Expected HashMap, got %T", result)
	}

	// Test keyword access
	nameKeyword := InternKeyword("name")
	name, err := nameKeyword.Call([]Value{hm}, env)
	if err != nil {
		t.Fatalf("Error accessing name with keyword: %v", err)
	}

	if name.String() != "\"Alice\"" {
		t.Errorf("Expected \"Alice\", got %s", name.String())
	}

	// Test hash-map-get with keyword
	getExpr := "(hash-map-get (hash-map :name \"Bob\") :name)"
	parsedGetExpr, err := Parse(getExpr)
	if err != nil {
		t.Fatalf("Failed to parse get expression: %v", err)
	}

	getResult, err := Eval(parsedGetExpr, env)
	if err != nil {
		t.Fatalf("Failed to evaluate get expression: %v", err)
	}

	if getResult.String() != "\"Bob\"" {
		t.Errorf("Expected \"Bob\", got %s", getResult.String())
	}
}

func TestKeywordAsMapKey(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)

	// Test putting a keyword as a key
	expr := "(hash-map-put (hash-map) :new-key \"new-value\")"
	parsedExpr, err := Parse(expr)
	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	result, err := Eval(parsedExpr, env)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	hm, ok := result.(*HashMap)
	if !ok {
		t.Fatalf("Expected HashMap, got %T", result)
	}

	// Test accessing the value with the keyword
	keyword := InternKeyword("new-key")
	value, err := keyword.Call([]Value{hm}, env)
	if err != nil {
		t.Fatalf("Error accessing value with keyword: %v", err)
	}

	if value.String() != "\"new-value\"" {
		t.Errorf("Expected \"new-value\", got %s", value.String())
	}
}
