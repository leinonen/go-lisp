package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestStringConcat(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "concatenate two strings",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-concat"},
					&types.StringExpr{Value: "Hello"},
					&types.StringExpr{Value: " World"},
				},
			},
			expected: types.StringValue("Hello World"),
		},
		{
			name: "concatenate multiple strings",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-concat"},
					&types.StringExpr{Value: "A"},
					&types.StringExpr{Value: "B"},
					&types.StringExpr{Value: "C"},
				},
			},
			expected: types.StringValue("ABC"),
		},
		{
			name: "concatenate mixed types",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-concat"},
					&types.StringExpr{Value: "Number: "},
					&types.NumberExpr{Value: 42},
				},
			},
			expected: types.StringValue("Number: 42"),
		},
		{
			name: "empty concatenation",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-concat"},
				},
			},
			expected: types.StringValue(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringLength(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "length of normal string",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-length"},
					&types.StringExpr{Value: "hello"},
				},
			},
			expected: types.NumberValue(5),
		},
		{
			name: "length of empty string",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-length"},
					&types.StringExpr{Value: ""},
				},
			},
			expected: types.NumberValue(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringSubstring(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "substring middle",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-substring"},
					&types.StringExpr{Value: "hello"},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 4},
				},
			},
			expected: types.StringValue("ell"),
		},
		{
			name: "substring from beginning",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-substring"},
					&types.StringExpr{Value: "hello"},
					&types.NumberExpr{Value: 0},
					&types.NumberExpr{Value: 2},
				},
			},
			expected: types.StringValue("he"),
		},
		{
			name: "substring to end",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-substring"},
					&types.StringExpr{Value: "hello"},
					&types.NumberExpr{Value: 3},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.StringValue("lo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringCharAt(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "char at index 0",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-char-at"},
					&types.StringExpr{Value: "hello"},
					&types.NumberExpr{Value: 0},
				},
			},
			expected: types.StringValue("h"),
		},
		{
			name: "char at middle index",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-char-at"},
					&types.StringExpr{Value: "hello"},
					&types.NumberExpr{Value: 2},
				},
			},
			expected: types.StringValue("l"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringCase(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "uppercase",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-upper"},
					&types.StringExpr{Value: "hello"},
				},
			},
			expected: types.StringValue("HELLO"),
		},
		{
			name: "lowercase",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-lower"},
					&types.StringExpr{Value: "HELLO"},
				},
			},
			expected: types.StringValue("hello"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringTrim(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "string-trim"},
			&types.StringExpr{Value: "  hello world  "},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := types.StringValue("hello world")
	if !valuesEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestStringSplit(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "string-split"},
			&types.StringExpr{Value: "a,b,c"},
			&types.StringExpr{Value: ","},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := &types.ListValue{
		Elements: []types.Value{
			types.StringValue("a"),
			types.StringValue("b"),
			types.StringValue("c"),
		},
	}

	if !valuesEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestStringJoin(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "string-join"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "list"},
					&types.StringExpr{Value: "a"},
					&types.StringExpr{Value: "b"},
					&types.StringExpr{Value: "c"},
				},
			},
			&types.StringExpr{Value: ","},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := types.StringValue("a,b,c")
	if !valuesEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestStringPredicates(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "string contains substring",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-contains?"},
					&types.StringExpr{Value: "hello world"},
					&types.StringExpr{Value: "world"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "string does not contain substring",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-contains?"},
					&types.StringExpr{Value: "hello world"},
					&types.StringExpr{Value: "xyz"},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "string starts with prefix",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-starts-with?"},
					&types.StringExpr{Value: "hello world"},
					&types.StringExpr{Value: "hello"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "string does not start with prefix",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-starts-with?"},
					&types.StringExpr{Value: "hello world"},
					&types.StringExpr{Value: "world"},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "string ends with suffix",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-ends-with?"},
					&types.StringExpr{Value: "hello world"},
					&types.StringExpr{Value: "world"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "string does not end with suffix",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-ends-with?"},
					&types.StringExpr{Value: "hello world"},
					&types.StringExpr{Value: "hello"},
				},
			},
			expected: types.BooleanValue(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringReplace(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "string-replace"},
			&types.StringExpr{Value: "hello world"},
			&types.StringExpr{Value: "world"},
			&types.StringExpr{Value: "universe"},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := types.StringValue("hello universe")
	if !valuesEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestStringIndexOf(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "found substring",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-index-of"},
					&types.StringExpr{Value: "hello world"},
					&types.StringExpr{Value: "world"},
				},
			},
			expected: types.NumberValue(6),
		},
		{
			name: "not found substring",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-index-of"},
					&types.StringExpr{Value: "hello world"},
					&types.StringExpr{Value: "xyz"},
				},
			},
			expected: types.NumberValue(-1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringNumberConversion(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "string to number",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string->number"},
					&types.StringExpr{Value: "42.5"},
				},
			},
			expected: types.NumberValue(42.5),
		},
		{
			name: "number to string",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "number->string"},
					&types.NumberExpr{Value: 42.5},
				},
			},
			expected: types.StringValue("42.5"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringRegex(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "regex match success",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-regex-match?"},
					&types.StringExpr{Value: "hello123"},
					&types.StringExpr{Value: "[a-z]+[0-9]+"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "regex match failure",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-regex-match?"},
					&types.StringExpr{Value: "hello"},
					&types.StringExpr{Value: "[0-9]+"},
				},
			},
			expected: types.BooleanValue(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringRegexFindAll(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "string-regex-find-all"},
			&types.StringExpr{Value: "abc123def456"},
			&types.StringExpr{Value: "[0-9]+"},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := &types.ListValue{
		Elements: []types.Value{
			types.StringValue("123"),
			types.StringValue("456"),
		},
	}

	if !valuesEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestStringRepeat(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "string-repeat"},
			&types.StringExpr{Value: "ha"},
			&types.NumberExpr{Value: 3},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := types.StringValue("hahaha")
	if !valuesEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestStringTypeChecks(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.Value
	}{
		{
			name: "string? with string",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string?"},
					&types.StringExpr{Value: "hello"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "string? with number",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string?"},
					&types.NumberExpr{Value: 42},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "string-empty? with empty string",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-empty?"},
					&types.StringExpr{Value: ""},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "string-empty? with non-empty string",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-empty?"},
					&types.StringExpr{Value: "hello"},
				},
			},
			expected: types.BooleanValue(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringErrors(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name string
		expr *types.ListExpr
	}{
		{
			name: "string-length with wrong type",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-length"},
					&types.NumberExpr{Value: 42},
				},
			},
		},
		{
			name: "string-substring out of bounds",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string-substring"},
					&types.StringExpr{Value: "hello"},
					&types.NumberExpr{Value: 10},
					&types.NumberExpr{Value: 15},
				},
			},
		},
		{
			name: "string->number with invalid format",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "string->number"},
					&types.StringExpr{Value: "not-a-number"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluator.Eval(tt.expr)
			if err == nil {
				t.Error("expected error but got none")
			}
		})
	}
}
