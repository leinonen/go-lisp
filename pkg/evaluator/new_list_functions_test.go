package evaluator

import (
	"reflect"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestEvaluatorNewListFunctions(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
		hasError bool
	}{
		// last function tests
		{
			name: "last of non-empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "last"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.NumberExpr{Value: 2},
							&types.NumberExpr{Value: 3},
						},
					},
				},
			},
			expected: types.NumberValue(3),
			hasError: false,
		},
		{
			name: "last of single element list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "last"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 42},
						},
					},
				},
			},
			expected: types.NumberValue(42),
			hasError: false,
		},

		// butlast function tests
		{
			name: "butlast of non-empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "butlast"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.NumberExpr{Value: 2},
							&types.NumberExpr{Value: 3},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(1), types.NumberValue(2)}},
			hasError: false,
		},
		{
			name: "butlast of single element list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "butlast"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 42},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{}},
			hasError: false,
		},
		{
			name: "butlast of empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "butlast"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{}},
			hasError: false,
		},

		// flatten function tests
		{
			name: "flatten nested list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "flatten"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "list"},
									&types.NumberExpr{Value: 2},
									&types.NumberExpr{Value: 3},
								},
							},
							&types.NumberExpr{Value: 4},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
				types.NumberValue(4),
			}},
			hasError: false,
		},

		// distinct function tests
		{
			name: "distinct removes duplicates",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "distinct"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.NumberExpr{Value: 2},
							&types.NumberExpr{Value: 1},
							&types.NumberExpr{Value: 3},
							&types.NumberExpr{Value: 2},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			}},
			hasError: false,
		},

		// concat function tests
		{
			name: "concat multiple lists",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "concat"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.NumberExpr{Value: 2},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 3},
							&types.NumberExpr{Value: 4},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 5},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
				types.NumberValue(4),
				types.NumberValue(5),
			}},
			hasError: false,
		},

		// partition function tests
		{
			name: "partition list into chunks",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "partition"},
					&types.NumberExpr{Value: 2},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.NumberExpr{Value: 2},
							&types.NumberExpr{Value: 3},
							&types.NumberExpr{Value: 4},
							&types.NumberExpr{Value: 5},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{
				&types.ListValue{Elements: []types.Value{types.NumberValue(1), types.NumberValue(2)}},
				&types.ListValue{Elements: []types.Value{types.NumberValue(3), types.NumberValue(4)}},
				&types.ListValue{Elements: []types.Value{types.NumberValue(5)}},
			}},
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			result, err := evaluator.Eval(test.expr)

			if test.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}
