package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestHashMapOperations(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name: "hash-map creation with no arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map"},
				},
			},
			expected: &types.HashMapValue{Elements: map[string]types.Value{}},
		},
		{
			name: "hash-map creation with key-value pairs",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map"},
					&types.StringExpr{Value: "name"},
					&types.StringExpr{Value: "Alice"},
					&types.StringExpr{Value: "age"},
					&types.NumberExpr{Value: 30},
				},
			},
			expected: &types.HashMapValue{
				Elements: map[string]types.Value{
					"name": types.StringValue("Alice"),
					"age":  types.NumberValue(30),
				},
			},
		},
		{
			name: "hash-map-get with existing key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-get"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "x"},
							&types.NumberExpr{Value: 42},
						},
					},
					&types.StringExpr{Value: "x"},
				},
			},
			expected: types.NumberValue(42),
		},
		{
			name: "hash-map-get with non-existing key returns nil",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-get"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "x"},
							&types.NumberExpr{Value: 42},
						},
					},
					&types.StringExpr{Value: "y"},
				},
			},
			expected: &types.NilValue{},
		},
		{
			name: "hash-map-put adds new key-value pair",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-put"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "x"},
							&types.NumberExpr{Value: 42},
						},
					},
					&types.StringExpr{Value: "y"},
					&types.NumberExpr{Value: 99},
				},
			},
			expected: &types.HashMapValue{
				Elements: map[string]types.Value{
					"x": types.NumberValue(42),
					"y": types.NumberValue(99),
				},
			},
		},
		{
			name: "hash-map-put updates existing key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-put"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "x"},
							&types.NumberExpr{Value: 42},
						},
					},
					&types.StringExpr{Value: "x"},
					&types.NumberExpr{Value: 100},
				},
			},
			expected: &types.HashMapValue{
				Elements: map[string]types.Value{
					"x": types.NumberValue(100),
				},
			},
		},
		{
			name: "hash-map-remove removes existing key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-remove"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "x"},
							&types.NumberExpr{Value: 42},
							&types.StringExpr{Value: "y"},
							&types.NumberExpr{Value: 99},
						},
					},
					&types.StringExpr{Value: "x"},
				},
			},
			expected: &types.HashMapValue{
				Elements: map[string]types.Value{
					"y": types.NumberValue(99),
				},
			},
		},
		{
			name: "hash-map-contains? returns true for existing key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-contains?"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "x"},
							&types.NumberExpr{Value: 42},
						},
					},
					&types.StringExpr{Value: "x"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "hash-map-contains? returns false for non-existing key",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-contains?"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "x"},
							&types.NumberExpr{Value: 42},
						},
					},
					&types.StringExpr{Value: "y"},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "hash-map-keys returns list of all keys",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-keys"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "a"},
							&types.NumberExpr{Value: 1},
							&types.StringExpr{Value: "b"},
							&types.NumberExpr{Value: 2},
						},
					},
				},
			},
			// Note: order is not guaranteed in hash maps, so we'll test both possible orders
			expected: &types.ListValue{
				Elements: []types.Value{
					types.StringValue("a"),
					types.StringValue("b"),
				},
			},
		},
		{
			name: "hash-map-values returns list of all values",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-values"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "a"},
							&types.NumberExpr{Value: 1},
							&types.StringExpr{Value: "b"},
							&types.NumberExpr{Value: 2},
						},
					},
				},
			},
			// Note: order is not guaranteed in hash maps, so we'll test both possible orders
			expected: &types.ListValue{
				Elements: []types.Value{
					types.NumberValue(1),
					types.NumberValue(2),
				},
			},
		},
		{
			name: "hash-map-size returns number of key-value pairs",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-size"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "a"},
							&types.NumberExpr{Value: 1},
							&types.StringExpr{Value: "b"},
							&types.NumberExpr{Value: 2},
						},
					},
				},
			},
			expected: types.NumberValue(2),
		},
		{
			name: "hash-map-empty? returns true for empty hash map",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-empty?"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
						},
					},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "hash-map-empty? returns false for non-empty hash map",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-empty?"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
							&types.StringExpr{Value: "x"},
							&types.NumberExpr{Value: 42},
						},
					},
				},
			},
			expected: types.BooleanValue(false),
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

			// Special handling for hash-map-keys and hash-map-values since order is not guaranteed
			if tt.name == "hash-map-keys returns list of all keys" || tt.name == "hash-map-values returns list of all values" {
				resultList, ok := result.(*types.ListValue)
				if !ok {
					t.Fatalf("expected ListValue, got %T", result)
				}
				expectedList, ok := tt.expected.(*types.ListValue)
				if !ok {
					t.Fatalf("expected test case to have ListValue, got %T", tt.expected)
				}

				// Check that both lists have same length
				if len(resultList.Elements) != len(expectedList.Elements) {
					t.Errorf("expected %d elements, got %d", len(expectedList.Elements), len(resultList.Elements))
					return
				}

				// Check that all expected elements are present (order doesn't matter)
				for _, expectedElem := range expectedList.Elements {
					found := false
					for _, resultElem := range resultList.Elements {
						if valuesEqual(expectedElem, resultElem) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected element %v not found in result %v", expectedElem, resultList.Elements)
					}
				}
				return
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHashMapErrors(t *testing.T) {
	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "hash-map with odd number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map"},
					&types.StringExpr{Value: "key1"},
					// missing value for key1
				},
			},
		},
		{
			name: "hash-map-get with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-get"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
						},
					},
					// missing key argument
				},
			},
		},
		{
			name: "hash-map-get with non-hash-map",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-get"},
					&types.NumberExpr{Value: 42}, // not a hash map
					&types.StringExpr{Value: "key"},
				},
			},
		},
		{
			name: "hash-map-put with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-put"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "hash-map"},
						},
					},
					&types.StringExpr{Value: "key"},
					// missing value argument
				},
			},
		},
		{
			name: "hash-map-contains? with non-hash-map",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map-contains?"},
					&types.StringExpr{Value: "not-a-hash-map"},
					&types.StringExpr{Value: "key"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			_, err := evaluator.Eval(tt.expr)
			if err == nil {
				t.Errorf("expected error for expression %v", tt.expr)
			}
		})
	}
}
