// Test file for new comparison and logical operators
package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestNewComparisonOperators(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name: "less than or equal - true case",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "<="},
					&types.NumberExpr{Value: 3},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "less than or equal - equal case",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "<="},
					&types.NumberExpr{Value: 5},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "less than or equal - false case",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "<="},
					&types.NumberExpr{Value: 7},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "greater than or equal - true case",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: ">="},
					&types.NumberExpr{Value: 7},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "greater than or equal - equal case",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: ">="},
					&types.NumberExpr{Value: 5},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "greater than or equal - false case",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: ">="},
					&types.NumberExpr{Value: 3},
					&types.NumberExpr{Value: 5},
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

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name: "and - all true",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "and"},
					&types.BooleanExpr{Value: true},
					&types.BooleanExpr{Value: true},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "and - one false",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "and"},
					&types.BooleanExpr{Value: true},
					&types.BooleanExpr{Value: false},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "and - multiple arguments all true",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "and"},
					&types.BooleanExpr{Value: true},
					&types.BooleanExpr{Value: true},
					&types.BooleanExpr{Value: true},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "and - empty arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "and"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "or - all false",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "or"},
					&types.BooleanExpr{Value: false},
					&types.BooleanExpr{Value: false},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "or - one true",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "or"},
					&types.BooleanExpr{Value: false},
					&types.BooleanExpr{Value: true},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "or - multiple arguments one true",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "or"},
					&types.BooleanExpr{Value: false},
					&types.BooleanExpr{Value: false},
					&types.BooleanExpr{Value: true},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "or - empty arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "or"},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "not - true",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "not"},
					&types.BooleanExpr{Value: true},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "not - false",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "not"},
					&types.BooleanExpr{Value: false},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "complex logical expression",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "and"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "or"},
							&types.BooleanExpr{Value: true},
							&types.BooleanExpr{Value: false},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "not"},
							&types.BooleanExpr{Value: false},
						},
					},
				},
			},
			expected: types.BooleanValue(true),
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
