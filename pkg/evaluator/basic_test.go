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

// TestModulo tests the modulo operation
func TestModulo(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name: "basic modulo",
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
			name: "modulo with negative dividend",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.NumberExpr{Value: -17},
					&types.NumberExpr{Value: 5},
				},
			},
			expected: types.NumberValue(-2),
		},
		{
			name: "modulo with negative divisor",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.NumberExpr{Value: 17},
					&types.NumberExpr{Value: -5},
				},
			},
			expected: types.NumberValue(2),
		},
		{
			name: "modulo with both negative",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.NumberExpr{Value: -17},
					&types.NumberExpr{Value: -5},
				},
			},
			expected: types.NumberValue(-2),
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

// TestModuloErrors tests error cases for modulo operation
func TestModuloErrors(t *testing.T) {
	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "modulo by zero",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.NumberExpr{Value: 17},
					&types.NumberExpr{Value: 0},
				},
			},
		},
		{
			name: "wrong number of arguments - too few",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.NumberExpr{Value: 17},
				},
			},
		},
		{
			name: "wrong number of arguments - too many",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.NumberExpr{Value: 17},
					&types.NumberExpr{Value: 5},
					&types.NumberExpr{Value: 3},
				},
			},
		},
		{
			name: "non-numeric first argument",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.StringExpr{Value: "hello"},
					&types.NumberExpr{Value: 5},
				},
			},
		},
		{
			name: "non-numeric second argument",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "%"},
					&types.NumberExpr{Value: 17},
					&types.BooleanExpr{Value: true},
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

// TestModuloWithBigNumbers tests modulo operation with big numbers
func TestModuloWithBigNumbers(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		checkErr bool
	}{
		{
			name: "big number modulo regular number",
			expr: "(% 123456789012345678901234567890 123)",
		},
		{
			name: "regular number modulo big number",
			expr: "(% 123456 987654321)",
		},
		{
			name: "big number modulo big number",
			expr: "(% 123456789012345678901234567890 987654321)",
		},
		{
			name:     "big number modulo by zero",
			expr:     "(% 123456789012345678901234567890 0)",
			checkErr: true,
		},
		{
			name:     "regular number modulo by zero",
			expr:     "(% 123 0)",
			checkErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAndEval(t, tt.expr)
			if tt.checkErr {
				if err == nil {
					t.Errorf("expected error for modulo by zero")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// For valid operations, just ensure we get a result without error
			if result == nil {
				t.Errorf("expected non-nil result")
			}
		})
	}
}

// TestNilEquality tests that nil values can be compared for equality
func TestNilEquality(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name: "nil equals nil",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.SymbolExpr{Name: "nil"},
					&types.SymbolExpr{Name: "nil"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "nil not equal to number",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.SymbolExpr{Name: "nil"},
					&types.NumberExpr{Value: 42},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "nil not equal to string",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.SymbolExpr{Name: "nil"},
					&types.StringExpr{Value: "hello"},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "nil not equal to boolean",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.SymbolExpr{Name: "nil"},
					&types.BooleanExpr{Value: false},
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

// TestEnhancedEquality tests that all value types can be compared for equality
func TestEnhancedEquality(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name: "string equality",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.StringExpr{Value: "hello"},
					&types.StringExpr{Value: "hello"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "string inequality",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.StringExpr{Value: "hello"},
					&types.StringExpr{Value: "world"},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "boolean equality",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.BooleanExpr{Value: true},
					&types.BooleanExpr{Value: true},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "boolean inequality",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.BooleanExpr{Value: true},
					&types.BooleanExpr{Value: false},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "keyword equality",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.KeywordExpr{Value: "name"},
					&types.KeywordExpr{Value: "name"},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "keyword inequality",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.KeywordExpr{Value: "name"},
					&types.KeywordExpr{Value: "age"},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "cross-type inequality",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.StringExpr{Value: "42"},
					&types.NumberExpr{Value: 42},
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

func TestListEquality(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		expr     *types.ListExpr
		expected types.BooleanValue
	}{
		{
			name: "equal empty lists",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
						},
					},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "equal single element lists",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 42},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 42},
						},
					},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "equal multi-element lists",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.StringExpr{Value: "hello"},
							&types.BooleanExpr{Value: true},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.StringExpr{Value: "hello"},
							&types.BooleanExpr{Value: true},
						},
					},
				},
			},
			expected: types.BooleanValue(true),
		},
		{
			name: "unequal lists - different values",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
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
							&types.NumberExpr{Value: 1},
							&types.NumberExpr{Value: 3},
						},
					},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "unequal lists - different lengths",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.NumberExpr{Value: 2},
						},
					},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "nested list equality",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "="},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "list"},
									&types.NumberExpr{Value: 1},
									&types.NumberExpr{Value: 2},
								},
							},
							&types.NumberExpr{Value: 3},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "list"},
									&types.NumberExpr{Value: 1},
									&types.NumberExpr{Value: 2},
								},
							},
							&types.NumberExpr{Value: 3},
						},
					},
				},
			},
			expected: types.BooleanValue(true),
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
