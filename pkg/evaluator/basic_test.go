package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestEvaluator(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name:     "number literal",
			expr:     &types.NumberExpr{Value: 42},
			expected: types.NumberValue(42),
		},
		{
			name:     "string literal",
			expr:     &types.StringExpr{Value: "hello"},
			expected: types.StringValue("hello"),
		},
		{
			name:     "boolean true",
			expr:     &types.BooleanExpr{Value: true},
			expected: types.BooleanValue(true),
		},
		{
			name:     "boolean false",
			expr:     &types.BooleanExpr{Value: false},
			expected: types.BooleanValue(false),
		},
		{
			name: "addition",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 2},
				},
			},
			expected: types.NumberValue(3),
		},
		{
			name: "subtraction",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "-"},
					&types.NumberExpr{Value: 10},
					&types.NumberExpr{Value: 3},
				},
			},
			expected: types.NumberValue(7),
		},
		{
			name: "multiplication",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "*"},
					&types.NumberExpr{Value: 6},
					&types.NumberExpr{Value: 7},
				},
			},
			expected: types.NumberValue(42),
		},
		{
			name: "division",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "/"},
					&types.NumberExpr{Value: 15},
					&types.NumberExpr{Value: 3},
				},
			},
			expected: types.NumberValue(5),
		},
		{
			name: "nested arithmetic",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "*"},
							&types.NumberExpr{Value: 2},
							&types.NumberExpr{Value: 3},
						},
					},
					&types.NumberExpr{Value: 4},
				},
			},
			expected: types.NumberValue(10),
		},
		{
			name: "multiple operands addition",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 2},
					&types.NumberExpr{Value: 3},
					&types.NumberExpr{Value: 4},
				},
			},
			expected: types.NumberValue(10),
		},
		{
			name: "equality comparison",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.NumberExpr{Value: 5},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "inequality comparison",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.NumberExpr{Value: 5},
					&types.NumberExpr{Value: 3},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "less than",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "<"},
					&types.NumberExpr{Value: 3},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "greater than",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: ">"},
					&types.NumberExpr{Value: 7},
					&types.NumberExpr{Value: 3},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "less than or equal - true",
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
			name: "less than or equal - equal",
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
			name: "greater than or equal - true",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: ">="},
					&types.NumberExpr{Value: 7},
					&types.NumberExpr{Value: 3},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "and - true",
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
			name: "and - false",
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
			name: "or - true",
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
			name: "not - true to false",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "not"},
					&types.BooleanExpr{Value: true},
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

func TestEvaluatorVariables(t *testing.T) {
	env := NewEnvironment()
	env.Set("x", types.NumberValue(10))
	env.Set("y", types.NumberValue(20))

	evaluator := NewEvaluator(env)

	// Test variable lookup
	result, err := evaluator.Eval(&types.SymbolExpr{Name: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valuesEqual(result, types.NumberValue(10)) {
		t.Errorf("expected 10, got %v", result)
	}

	// Test expression with variables
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "+"},
			&types.SymbolExpr{Name: "x"},
			&types.SymbolExpr{Name: "y"},
		},
	}

	result, err = evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valuesEqual(result, types.NumberValue(30)) {
		t.Errorf("expected 30, got %v", result)
	}
}

func TestEvaluatorError(t *testing.T) {
	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "undefined symbol",
			expr: &types.SymbolExpr{Name: "undefined"},
		},
		{
			name: "division by zero",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "/"},
					&types.NumberExpr{Value: 10},
					&types.NumberExpr{Value: 0},
				},
			},
		},
		{
			name: "invalid function",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "nonexistent"},
					&types.NumberExpr{Value: 1},
				},
			},
		},
		{
			name: "wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
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
