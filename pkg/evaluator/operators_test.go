package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestAllOperators(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		// Literals
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
		
		// Arithmetic operators
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
			name: "modulo basic",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.NumberExpr{Value: 17},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.NumberValue(2),
		},
		{
			name: "modulo with zero remainder",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.NumberExpr{Value: 10},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.NumberValue(0),
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
		
		// Comparison operators
		{
			name: "equality comparison true",
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
			name: "equality comparison false",
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
			name: "less than true",
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
			name: "less than false",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "<"},
					&types.NumberExpr{Value: 7},
					&types.NumberExpr{Value: 3},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "greater than true",
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
			name: "greater than false",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: ">"},
					&types.NumberExpr{Value: 2},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.BooleanValue(false),
		},
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
		
		// Logical operators
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
		
		// Conditional operators
		{
			name: "if with true condition",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "if"},
					&types.BooleanExpr{Value: true},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 2},
				},
			},
			expected: types.NumberValue(1),
		},
		{
			name: "if with false condition",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "if"},
					&types.BooleanExpr{Value: false},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 2},
				},
			},
			expected: types.NumberValue(2),
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

// Keep essential non-operator tests
func TestEvaluatorVariables(t *testing.T) {
	env := NewEnvironment()
	env.Set("x", types.NumberValue(10))
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name:     "variable lookup",
			expr:     &types.SymbolExpr{Name: "x"},
			expected: types.NumberValue(10),
		},
		{
			name: "variable in expression",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.SymbolExpr{Name: "x"},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.NumberValue(15),
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
