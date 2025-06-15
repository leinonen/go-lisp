package string

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
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

func wrapValue(value types.Value) types.Expr {
	return valueExpr{value}
}

type valueExpr struct {
	value types.Value
}

func (ve valueExpr) String() string {
	return ve.value.String()
}

func TestStringPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewStringPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{
		"string-concat", "string-length", "string-substring", "string-char-at",
		"string-upper", "string-lower", "string-trim", "string-split", "string-join",
		"string-contains?", "string-starts-with?", "string-ends-with?", "string-replace",
		"string-index-of", "string->number", "number->string", "string-regex-match?",
		"string-regex-find-all", "string-repeat", "string?", "string-empty?",
	}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestStringPlugin_BasicStringOperations(t *testing.T) {
	plugin := NewStringPlugin()
	evaluator := newMockEvaluator()

	// Test string-concat
	concatArgs := []types.Expr{
		&types.StringExpr{Value: "Hello"},
		&types.StringExpr{Value: " "},
		&types.StringExpr{Value: "World"},
	}
	result, err := plugin.evalStringConcat(evaluator, concatArgs)
	if err != nil {
		t.Fatalf("evalStringConcat failed: %v", err)
	}
	if result != types.StringValue("Hello World") {
		t.Errorf("Expected 'Hello World', got %v", result)
	}

	// Test string-length
	lengthArgs := []types.Expr{&types.StringExpr{Value: "Hello"}}
	result, err = plugin.evalStringLength(evaluator, lengthArgs)
	if err != nil {
		t.Fatalf("evalStringLength failed: %v", err)
	}
	if result != types.NumberValue(5) {
		t.Errorf("Expected 5, got %v", result)
	}

	// Test string-substring
	substringArgs := []types.Expr{
		&types.StringExpr{Value: "Hello World"},
		&types.NumberExpr{Value: 6},
		&types.NumberExpr{Value: 11},
	}
	result, err = plugin.evalStringSubstring(evaluator, substringArgs)
	if err != nil {
		t.Fatalf("evalStringSubstring failed: %v", err)
	}
	if result != types.StringValue("World") {
		t.Errorf("Expected 'World', got %v", result)
	}

	// Test string-char-at
	charAtArgs := []types.Expr{
		&types.StringExpr{Value: "Hello"},
		&types.NumberExpr{Value: 1},
	}
	result, err = plugin.evalStringCharAt(evaluator, charAtArgs)
	if err != nil {
		t.Fatalf("evalStringCharAt failed: %v", err)
	}
	if result != types.StringValue("e") {
		t.Errorf("Expected 'e', got %v", result)
	}

	// Test string-upper
	upperArgs := []types.Expr{&types.StringExpr{Value: "hello"}}
	result, err = plugin.evalStringUpper(evaluator, upperArgs)
	if err != nil {
		t.Fatalf("evalStringUpper failed: %v", err)
	}
	if result != types.StringValue("HELLO") {
		t.Errorf("Expected 'HELLO', got %v", result)
	}

	// Test string-lower
	lowerArgs := []types.Expr{&types.StringExpr{Value: "HELLO"}}
	result, err = plugin.evalStringLower(evaluator, lowerArgs)
	if err != nil {
		t.Fatalf("evalStringLower failed: %v", err)
	}
	if result != types.StringValue("hello") {
		t.Errorf("Expected 'hello', got %v", result)
	}

	// Test string-trim
	trimArgs := []types.Expr{&types.StringExpr{Value: "  hello  "}}
	result, err = plugin.evalStringTrim(evaluator, trimArgs)
	if err != nil {
		t.Fatalf("evalStringTrim failed: %v", err)
	}
	if result != types.StringValue("hello") {
		t.Errorf("Expected 'hello', got %v", result)
	}
}

func TestStringPlugin_StringPredicates(t *testing.T) {
	plugin := NewStringPlugin()
	evaluator := newMockEvaluator()

	// Test string-contains?
	containsArgs := []types.Expr{
		&types.StringExpr{Value: "Hello World"},
		&types.StringExpr{Value: "World"},
	}
	result, err := plugin.evalStringContains(evaluator, containsArgs)
	if err != nil {
		t.Fatalf("evalStringContains failed: %v", err)
	}
	if result != types.BooleanValue(true) {
		t.Errorf("Expected true, got %v", result)
	}

	// Test string-starts-with?
	startsWithArgs := []types.Expr{
		&types.StringExpr{Value: "Hello World"},
		&types.StringExpr{Value: "Hello"},
	}
	result, err = plugin.evalStringStartsWith(evaluator, startsWithArgs)
	if err != nil {
		t.Fatalf("evalStringStartsWith failed: %v", err)
	}
	if result != types.BooleanValue(true) {
		t.Errorf("Expected true, got %v", result)
	}

	// Test string-ends-with?
	endsWithArgs := []types.Expr{
		&types.StringExpr{Value: "Hello World"},
		&types.StringExpr{Value: "World"},
	}
	result, err = plugin.evalStringEndsWith(evaluator, endsWithArgs)
	if err != nil {
		t.Fatalf("evalStringEndsWith failed: %v", err)
	}
	if result != types.BooleanValue(true) {
		t.Errorf("Expected true, got %v", result)
	}

	// Test string?
	stringPredArgs := []types.Expr{&types.StringExpr{Value: "hello"}}
	result, err = plugin.evalStringPredicate(evaluator, stringPredArgs)
	if err != nil {
		t.Fatalf("evalStringPredicate failed: %v", err)
	}
	if result != types.BooleanValue(true) {
		t.Errorf("Expected true, got %v", result)
	}

	// Test string? with non-string
	stringPredArgs = []types.Expr{&types.NumberExpr{Value: 42}}
	result, err = plugin.evalStringPredicate(evaluator, stringPredArgs)
	if err != nil {
		t.Fatalf("evalStringPredicate failed: %v", err)
	}
	if result != types.BooleanValue(false) {
		t.Errorf("Expected false, got %v", result)
	}

	// Test string-empty?
	emptyArgs := []types.Expr{&types.StringExpr{Value: ""}}
	result, err = plugin.evalStringEmpty(evaluator, emptyArgs)
	if err != nil {
		t.Fatalf("evalStringEmpty failed: %v", err)
	}
	if result != types.BooleanValue(true) {
		t.Errorf("Expected true, got %v", result)
	}

	// Test string-empty? with non-empty string
	emptyArgs = []types.Expr{&types.StringExpr{Value: "hello"}}
	result, err = plugin.evalStringEmpty(evaluator, emptyArgs)
	if err != nil {
		t.Fatalf("evalStringEmpty failed: %v", err)
	}
	if result != types.BooleanValue(false) {
		t.Errorf("Expected false, got %v", result)
	}
}

func TestStringPlugin_StringConversion(t *testing.T) {
	plugin := NewStringPlugin()
	evaluator := newMockEvaluator()

	// Test string->number
	stringToNumArgs := []types.Expr{&types.StringExpr{Value: "42"}}
	result, err := plugin.evalStringToNumber(evaluator, stringToNumArgs)
	if err != nil {
		t.Fatalf("evalStringToNumber failed: %v", err)
	}
	if result != types.NumberValue(42) {
		t.Errorf("Expected 42, got %v", result)
	}

	// Test number->string
	numToStringArgs := []types.Expr{&types.NumberExpr{Value: 42}}
	result, err = plugin.evalNumberToString(evaluator, numToStringArgs)
	if err != nil {
		t.Fatalf("evalNumberToString failed: %v", err)
	}
	if result != types.StringValue("42") {
		t.Errorf("Expected '42', got %v", result)
	}
}

func TestStringPlugin_StringManipulation(t *testing.T) {
	plugin := NewStringPlugin()
	evaluator := newMockEvaluator()

	// Test string-replace
	replaceArgs := []types.Expr{
		&types.StringExpr{Value: "Hello World"},
		&types.StringExpr{Value: "World"},
		&types.StringExpr{Value: "Go"},
	}
	result, err := plugin.evalStringReplace(evaluator, replaceArgs)
	if err != nil {
		t.Fatalf("evalStringReplace failed: %v", err)
	}
	if result != types.StringValue("Hello Go") {
		t.Errorf("Expected 'Hello Go', got %v", result)
	}

	// Test string-index-of
	indexOfArgs := []types.Expr{
		&types.StringExpr{Value: "Hello World"},
		&types.StringExpr{Value: "World"},
	}
	result, err = plugin.evalStringIndexOf(evaluator, indexOfArgs)
	if err != nil {
		t.Fatalf("evalStringIndexOf failed: %v", err)
	}
	if result != types.NumberValue(6) {
		t.Errorf("Expected 6, got %v", result)
	}

	// Test string-repeat
	repeatArgs := []types.Expr{
		&types.StringExpr{Value: "Hi"},
		&types.NumberExpr{Value: 3},
	}
	result, err = plugin.evalStringRepeat(evaluator, repeatArgs)
	if err != nil {
		t.Fatalf("evalStringRepeat failed: %v", err)
	}
	if result != types.StringValue("HiHiHi") {
		t.Errorf("Expected 'HiHiHi', got %v", result)
	}
}

func TestStringPlugin_ErrorCases(t *testing.T) {
	plugin := NewStringPlugin()
	evaluator := newMockEvaluator()

	// Test string-substring with invalid indices
	_, err := plugin.evalStringSubstring(evaluator, []types.Expr{
		&types.StringExpr{Value: "Hello"},
		&types.NumberExpr{Value: 10}, // Index out of bounds
		&types.NumberExpr{Value: 15},
	})
	if err == nil {
		t.Error("Expected error for out of bounds substring")
	}

	// Test string-char-at with invalid index
	_, err = plugin.evalStringCharAt(evaluator, []types.Expr{
		&types.StringExpr{Value: "Hello"},
		&types.NumberExpr{Value: 10}, // Index out of bounds
	})
	if err == nil {
		t.Error("Expected error for out of bounds char-at")
	}

	// Test string->number with invalid string
	_, err = plugin.evalStringToNumber(evaluator, []types.Expr{
		&types.StringExpr{Value: "not-a-number"},
	})
	if err == nil {
		t.Error("Expected error for invalid number string")
	}

	// Test string functions with non-string arguments
	nonString := types.NumberValue(42)
	_, err = plugin.evalStringLength(evaluator, []types.Expr{
		wrapValue(nonString),
	})
	if err == nil {
		t.Error("Expected error for non-string argument")
	}
}
