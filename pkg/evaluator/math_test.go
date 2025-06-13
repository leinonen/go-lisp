package evaluator

import (
	"math"
	"strings"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestMathematicalFunctions(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Helper function to check if two floats are approximately equal
	approxEqual := func(a, b, tolerance float64) bool {
		return math.Abs(a-b) < tolerance
	}

	tests := []struct {
		name          string
		expr          *types.ListExpr
		expected      float64
		tolerance     float64
		expectError   bool
		errorContains string
	}{
		// sqrt tests
		{
			name: "sqrt of 16",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "sqrt"},
					&types.NumberExpr{Value: 16},
				},
			},
			expected:  4.0,
			tolerance: 1e-10,
		},
		{
			name: "sqrt of 2",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "sqrt"},
					&types.NumberExpr{Value: 2},
				},
			},
			expected:  math.Sqrt(2),
			tolerance: 1e-10,
		},
		{
			name: "sqrt of negative number",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "sqrt"},
					&types.NumberExpr{Value: -4},
				},
			},
			expectError:   true,
			errorContains: "negative number",
		},
		// pow tests
		{
			name: "pow 2^3",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "pow"},
					&types.NumberExpr{Value: 2},
					&types.NumberExpr{Value: 3},
				},
			},
			expected:  8.0,
			tolerance: 1e-10,
		},
		{
			name: "pow 5^0",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "pow"},
					&types.NumberExpr{Value: 5},
					&types.NumberExpr{Value: 0},
				},
			},
			expected:  1.0,
			tolerance: 1e-10,
		},
		// sin tests
		{
			name: "sin of 0",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "sin"},
					&types.NumberExpr{Value: 0},
				},
			},
			expected:  0.0,
			tolerance: 1e-10,
		},
		{
			name: "sin of pi/2",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "sin"},
					&types.NumberExpr{Value: math.Pi / 2},
				},
			},
			expected:  1.0,
			tolerance: 1e-10,
		},
		// cos tests
		{
			name: "cos of 0",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "cos"},
					&types.NumberExpr{Value: 0},
				},
			},
			expected:  1.0,
			tolerance: 1e-10,
		},
		{
			name: "cos of pi",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "cos"},
					&types.NumberExpr{Value: math.Pi},
				},
			},
			expected:  -1.0,
			tolerance: 1e-10,
		},
		// tan tests
		{
			name: "tan of 0",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "tan"},
					&types.NumberExpr{Value: 0},
				},
			},
			expected:  0.0,
			tolerance: 1e-10,
		},
		// log tests
		{
			name: "log of e",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "log"},
					&types.NumberExpr{Value: math.E},
				},
			},
			expected:  1.0,
			tolerance: 1e-10,
		},
		{
			name: "log of 1",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "log"},
					&types.NumberExpr{Value: 1},
				},
			},
			expected:  0.0,
			tolerance: 1e-10,
		},
		{
			name: "log of negative number",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "log"},
					&types.NumberExpr{Value: -1},
				},
			},
			expectError:   true,
			errorContains: "non-positive",
		},
		// exp tests
		{
			name: "exp of 0",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "exp"},
					&types.NumberExpr{Value: 0},
				},
			},
			expected:  1.0,
			tolerance: 1e-10,
		},
		{
			name: "exp of 1",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "exp"},
					&types.NumberExpr{Value: 1},
				},
			},
			expected:  math.E,
			tolerance: 1e-10,
		},
		// floor tests
		{
			name: "floor of 3.7",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "floor"},
					&types.NumberExpr{Value: 3.7},
				},
			},
			expected:  3.0,
			tolerance: 1e-10,
		},
		{
			name: "floor of -2.3",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "floor"},
					&types.NumberExpr{Value: -2.3},
				},
			},
			expected:  -3.0,
			tolerance: 1e-10,
		},
		// ceil tests
		{
			name: "ceil of 3.2",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "ceil"},
					&types.NumberExpr{Value: 3.2},
				},
			},
			expected:  4.0,
			tolerance: 1e-10,
		},
		{
			name: "ceil of -2.7",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "ceil"},
					&types.NumberExpr{Value: -2.7},
				},
			},
			expected:  -2.0,
			tolerance: 1e-10,
		},
		// round tests
		{
			name: "round of 3.4",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "round"},
					&types.NumberExpr{Value: 3.4},
				},
			},
			expected:  3.0,
			tolerance: 1e-10,
		},
		{
			name: "round of 3.6",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "round"},
					&types.NumberExpr{Value: 3.6},
				},
			},
			expected:  4.0,
			tolerance: 1e-10,
		},
		// abs tests
		{
			name: "abs of -5",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "abs"},
					&types.NumberExpr{Value: -5},
				},
			},
			expected:  5.0,
			tolerance: 1e-10,
		},
		{
			name: "abs of 3",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "abs"},
					&types.NumberExpr{Value: 3},
				},
			},
			expected:  3.0,
			tolerance: 1e-10,
		},
		// min tests
		{
			name: "min of 5 and 3",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "min"},
					&types.NumberExpr{Value: 5},
					&types.NumberExpr{Value: 3},
				},
			},
			expected:  3.0,
			tolerance: 1e-10,
		},
		// max tests
		{
			name: "max of 5 and 3",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "max"},
					&types.NumberExpr{Value: 5},
					&types.NumberExpr{Value: 3},
				},
			},
			expected:  5.0,
			tolerance: 1e-10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Eval(tt.expr)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			num, ok := result.(types.NumberValue)
			if !ok {
				t.Fatalf("expected NumberValue, got %T", result)
			}

			if !approxEqual(float64(num), tt.expected, tt.tolerance) {
				t.Errorf("expected %v, got %v", tt.expected, float64(num))
			}
		})
	}
}

func TestMathematicalConstants(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test pi constant
	t.Run("pi constant", func(t *testing.T) {
		result, err := evaluator.Eval(&types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "pi"},
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		num, ok := result.(types.NumberValue)
		if !ok {
			t.Fatalf("expected NumberValue, got %T", result)
		}

		if math.Abs(float64(num)-math.Pi) > 1e-10 {
			t.Errorf("expected %v, got %v", math.Pi, float64(num))
		}
	})

	// Test e constant
	t.Run("e constant", func(t *testing.T) {
		result, err := evaluator.Eval(&types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "e"},
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		num, ok := result.(types.NumberValue)
		if !ok {
			t.Fatalf("expected NumberValue, got %T", result)
		}

		if math.Abs(float64(num)-math.E) > 1e-10 {
			t.Errorf("expected %v, got %v", math.E, float64(num))
		}
	})
}

func TestRandomFunction(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test random with no arguments (0 to 1)
	t.Run("random with no args", func(t *testing.T) {
		result, err := evaluator.Eval(&types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "random"},
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		num, ok := result.(types.NumberValue)
		if !ok {
			t.Fatalf("expected NumberValue, got %T", result)
		}

		value := float64(num)
		if value < 0 || value >= 1 {
			t.Errorf("expected value between 0 and 1, got %v", value)
		}
	})

	// Test random with upper bound
	t.Run("random with upper bound", func(t *testing.T) {
		result, err := evaluator.Eval(&types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "random"},
				&types.NumberExpr{Value: 10},
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		num, ok := result.(types.NumberValue)
		if !ok {
			t.Fatalf("expected NumberValue, got %T", result)
		}

		value := float64(num)
		if value < 0 || value >= 10 || value != float64(int64(value)) {
			t.Errorf("expected integer between 0 and 9, got %v", value)
		}
	})

	// Test random with range
	t.Run("random with range", func(t *testing.T) {
		result, err := evaluator.Eval(&types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "random"},
				&types.NumberExpr{Value: 5},
				&types.NumberExpr{Value: 15},
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		num, ok := result.(types.NumberValue)
		if !ok {
			t.Fatalf("expected NumberValue, got %T", result)
		}

		value := float64(num)
		if value < 5 || value >= 15 || value != float64(int64(value)) {
			t.Errorf("expected integer between 5 and 14, got %v", value)
		}
	})
}

func TestMathematicalFunctionErrors(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Helper function to evaluate a string expression
	evalString := func(expr string) (types.Value, error) {
		tok := tokenizer.NewTokenizer(expr)
		tokens, err := tok.TokenizeWithError()
		if err != nil {
			return nil, err
		}

		p := parser.NewParser(tokens)
		ast, err := p.Parse()
		if err != nil {
			return nil, err
		}
		return evaluator.Eval(ast)
	}

	tests := []struct {
		name     string
		input    string
		errorMsg string
	}{
		{"sqrt with wrong number of args", "(sqrt 4 5)", "sqrt requires exactly 1 argument"},
		{"pow with wrong number of args", "(pow 2)", "pow requires exactly 2 arguments"},
		{"sqrt with non-number", `(sqrt "hello")`, "expected number"},
		{"random with invalid range", "(random 10 5)", "max must be greater than min"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evalString(tt.input)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.input)
			}
			if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("Expected error containing '%s', got: %v", tt.errorMsg, err)
			}
		})
	}
}

// Test additional mathematical functions
func TestAdditionalMathematicalFunctions(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Helper function to evaluate a string expression
	evalString := func(expr string) (types.Value, error) {
		tok := tokenizer.NewTokenizer(expr)
		tokens, err := tok.TokenizeWithError()
		if err != nil {
			return nil, err
		}

		p := parser.NewParser(tokens)
		ast, err := p.Parse()
		if err != nil {
			return nil, err
		}
		return evaluator.Eval(ast)
	}

	tests := []struct {
		name     string
		input    string
		expected float64
		delta    float64
	}{
		// Inverse trigonometric functions
		{"asin of 0.5", "(asin 0.5)", math.Asin(0.5), 1e-10},
		{"acos of 0.5", "(acos 0.5)", math.Acos(0.5), 1e-10},
		{"atan of 1", "(atan 1)", math.Atan(1), 1e-10},
		{"atan2 of (1, 1)", "(atan2 1 1)", math.Atan2(1, 1), 1e-10},

		// Hyperbolic functions
		{"sinh of 0", "(sinh 0)", math.Sinh(0), 1e-10},
		{"cosh of 0", "(cosh 0)", math.Cosh(0), 1e-10},
		{"tanh of 0", "(tanh 0)", math.Tanh(0), 1e-10},
		{"sinh of 1", "(sinh 1)", math.Sinh(1), 1e-10},

		// Angle conversion
		{"degrees of pi", "(degrees (pi))", 180.0, 1e-10},
		{"radians of 180", "(radians 180)", math.Pi, 1e-10},
		{"degrees of pi/2", "(degrees (/ (pi) 2))", 90.0, 1e-10},

		// Logarithm functions
		{"log10 of 100", "(log10 100)", 2.0, 1e-10},
		{"log2 of 8", "(log2 8)", 3.0, 1e-10},
		{"log10 of 1000", "(log10 1000)", 3.0, 1e-10},

		// Utility functions
		{"trunc of 3.7", "(trunc 3.7)", 3.0, 1e-10},
		{"trunc of -3.7", "(trunc -3.7)", -3.0, 1e-10},
		{"sign of 5", "(sign 5)", 1.0, 1e-10},
		{"sign of -5", "(sign -5)", -1.0, 1e-10},
		{"sign of 0", "(sign 0)", 0.0, 1e-10},
		{"mod of 7 and 3", "(mod 7 3)", 1.0, 1e-10},
		{"mod of -7 and 3", "(mod -7 3)", 2.0, 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evalString(tt.input)
			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.input, err)
				return
			}

			num, ok := result.(types.NumberValue)
			if !ok {
				t.Errorf("Expected NumberValue for %s, got %T", tt.input, result)
				return
			}

			if math.Abs(float64(num)-tt.expected) > tt.delta {
				t.Errorf("For %s: expected %f, got %f", tt.input, tt.expected, float64(num))
			}
		})
	}
}

// Test error cases for additional functions
func TestAdditionalMathematicalFunctionErrors(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Helper function to evaluate a string expression
	evalString := func(expr string) (types.Value, error) {
		tok := tokenizer.NewTokenizer(expr)
		tokens, err := tok.TokenizeWithError()
		if err != nil {
			return nil, err
		}

		p := parser.NewParser(tokens)
		ast, err := p.Parse()
		if err != nil {
			return nil, err
		}
		return evaluator.Eval(ast)
	}

	tests := []struct {
		name     string
		input    string
		errorMsg string
	}{
		// Domain errors for inverse trig functions
		{"asin out of range", "(asin 2)", "input must be in range [-1, 1]"},
		{"acos out of range", "(acos -2)", "input must be in range [-1, 1]"},

		// Wrong number of arguments
		{"asin with no args", "(asin)", "asin requires exactly 1 argument"},
		{"atan2 with one arg", "(atan2 1)", "atan2 requires exactly 2 arguments"},
		{"degrees with no args", "(degrees)", "degrees requires exactly 1 argument"},

		// Domain errors for logarithms
		{"log10 of negative", "(log10 -5)", "input must be positive"},
		{"log2 of zero", "(log2 0)", "input must be positive"},

		// Division by zero
		{"mod by zero", "(mod 5 0)", "division by zero"},

		// Non-number arguments
		{"sign of string", `(sign "hello")`, "expected number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evalString(tt.input)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.input)
			}
			if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("Expected error containing '%s', got: %v", tt.errorMsg, err)
			}
		})
	}
}
