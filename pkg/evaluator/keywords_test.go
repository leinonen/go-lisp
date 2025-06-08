package evaluator

import (
	"strings"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestKeywordBasics(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name:     "simple keyword evaluation",
			expr:     &types.KeywordExpr{Value: "name"},
			expected: types.KeywordValue("name"),
		},
		{
			name:     "keyword with dash",
			expr:     &types.KeywordExpr{Value: "first-name"},
			expected: types.KeywordValue("first-name"),
		},
		{
			name: "keyword in list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "list"},
					&types.KeywordExpr{Value: "key"},
					&types.StringExpr{Value: "value"},
				},
			},
			expected: &types.ListValue{
				Elements: []types.Value{
					types.KeywordValue("key"),
					types.StringValue("value"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

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

func TestKeywordHashMaps(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name: "hash-map with keyword keys",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map"},
					&types.KeywordExpr{Value: "name"},
					&types.StringExpr{Value: "Alice"},
					&types.KeywordExpr{Value: "age"},
					&types.NumberExpr{Value: 30},
				},
			},
			expected: &types.HashMapValue{
				Elements: map[string]types.Value{
					":name": types.StringValue("Alice"),
					":age":  types.NumberValue(30),
				},
			},
		},
		{
			name: "hash-map-get with keyword key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-get"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "name"},
							&types.StringExpr{Value: "Alice"},
						},
					},
					&types.KeywordExpr{Value: "name"},
				},
			},
			expected: types.StringValue("Alice"),
		},
		{
			name: "hash-map-put with keyword key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-put"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "name"},
							&types.StringExpr{Value: "Alice"},
						},
					},
					&types.KeywordExpr{Value: "age"},
					&types.NumberExpr{Value: 25},
				},
			},
			expected: &types.HashMapValue{
				Elements: map[string]types.Value{
					":name": types.StringValue("Alice"),
					":age":  types.NumberValue(25),
				},
			},
		},
		{
			name: "hash-map-contains? with keyword key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-contains?"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "name"},
							&types.StringExpr{Value: "Alice"},
						},
					},
					&types.KeywordExpr{Value: "name"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "hash-map-remove with keyword key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-remove"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "name"},
							&types.StringExpr{Value: "Alice"},
							&types.KeywordExpr{Value: "age"},
							&types.NumberExpr{Value: 30},
						},
					},
					&types.KeywordExpr{Value: "name"},
				},
			},
			expected: &types.HashMapValue{
				Elements: map[string]types.Value{
					":age": types.NumberValue(30),
				},
			},
		},
		{
			name: "mixed string and keyword keys",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map"},
					&types.StringExpr{Value: "string-key"},
					&types.StringExpr{Value: "string-value"},
					&types.KeywordExpr{Value: "keyword-key"},
					&types.StringExpr{Value: "keyword-value"},
				},
			},
			expected: &types.HashMapValue{
				Elements: map[string]types.Value{
					"string-key":   types.StringValue("string-value"),
					":keyword-key": types.StringValue("keyword-value"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

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

func TestKeywordErrors(t *testing.T) {
	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "hash-map-get with non-existent keyword key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-get"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "name"},
							&types.StringExpr{Value: "Alice"},
						},
					},
					&types.KeywordExpr{Value: "missing"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			result, err := evaluator.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Non-existent keys should return nil
			if result == nil || result.String() != "nil" {
				t.Errorf("expected nil for non-existent key, got %v", result)
			}
		})
	}
}

func TestKeywordAsFunction(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name: "keyword as function - simple get",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "name"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "name"},
							&types.StringExpr{Value: "Alice"},
							&types.KeywordExpr{Value: "age"},
							&types.NumberExpr{Value: 30},
						},
					},
				},
			},
			expected: types.StringValue("Alice"),
		},
		{
			name: "keyword as function - get number value",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "age"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "name"},
							&types.StringExpr{Value: "Bob"},
							&types.KeywordExpr{Value: "age"},
							&types.NumberExpr{Value: 25},
						},
					},
				},
			},
			expected: types.NumberValue(25),
		},
		{
			name: "keyword as function - non-existent key returns nil",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "missing"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "name"},
							&types.StringExpr{Value: "Charlie"},
						},
					},
				},
			},
			expected: &types.NilValue{},
		},
		{
			name: "keyword as function - empty hash map returns nil",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "key"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
						},
					},
				},
			},
			expected: &types.NilValue{},
		},
		{
			name: "keyword as function - mixed key types",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "keyword-key"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "string-key"},
							&types.StringExpr{Value: "string-value"},
							&types.KeywordExpr{Value: "keyword-key"},
							&types.StringExpr{Value: "keyword-value"},
						},
					},
				},
			},
			expected: types.StringValue("keyword-value"),
		},
		{
			name: "keyword as function - with default value",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "missing"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "existing"},
							&types.StringExpr{Value: "value"},
						},
					},
					&types.StringExpr{Value: "default"},
				},
			},
			expected: types.StringValue("default"),
		},
		{
			name: "keyword as function - with default value when key exists",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "existing"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "existing"},
							&types.StringExpr{Value: "found"},
						},
					},
					&types.StringExpr{Value: "default"},
				},
			},
			expected: types.StringValue("found"),
		},
		{
			name: "keyword as function - nested hash maps",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "user"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.KeywordExpr{Value: "user"},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "hash-map"},
									&types.KeywordExpr{Value: "name"},
									&types.StringExpr{Value: "Diana"},
								},
							},
						},
					},
				},
			},
			expected: &types.HashMapValue{
				Elements: map[string]types.Value{
					":name": types.StringValue("Diana"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

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

func TestKeywordAsFunctionErrors(t *testing.T) {
	tests := []struct {
		name        string
		expr        types.Expr
		expectError bool
		errorMsg    string
	}{
		{
			name: "keyword as function - wrong number of arguments (no args)",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "key"},
				},
			},
			expectError: true,
			errorMsg:    "keyword function requires 1 or 2 arguments",
		},
		{
			name: "keyword as function - too many arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "key"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
						},
					},
					&types.StringExpr{Value: "default"},
					&types.StringExpr{Value: "extra"},
				},
			},
			expectError: true,
			errorMsg:    "keyword function requires 1 or 2 arguments",
		},
		{
			name: "keyword as function - non-hash-map argument",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "key"},
					&types.StringExpr{Value: "not-a-hash-map"},
				},
			},
			expectError: true,
			errorMsg:    "keyword function first argument must be a hash map",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			_, err := evaluator.Eval(tt.expr)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
