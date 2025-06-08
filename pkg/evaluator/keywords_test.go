package evaluator

import (
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
