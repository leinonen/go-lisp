package evaluator

import (
	"math"
	"strings"
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
			name: "if expression - true condition",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "if"},
					&types.BooleanExpr{Value: true},
					&types.NumberExpr{Value: 42},
					&types.NumberExpr{Value: 0},
				},
			},
			expected: types.NumberValue(42),
		},
		{
			name: "if expression - false condition",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "if"},
					&types.BooleanExpr{Value: false},
					&types.NumberExpr{Value: 42},
					&types.NumberExpr{Value: 0},
				},
			},
			expected: types.NumberValue(0),
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

func TestEvaluatorDefine(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test basic variable definition
	defineExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "define"},
			&types.SymbolExpr{Name: "x"},
			&types.NumberExpr{Value: 42},
		},
	}

	result, err := evaluator.Eval(defineExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// define should return the value that was defined
	if !valuesEqual(result, types.NumberValue(42)) {
		t.Errorf("expected 42, got %v", result)
	}

	// Test that the variable was actually defined
	symbolExpr := &types.SymbolExpr{Name: "x"}
	result, err = evaluator.Eval(symbolExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valuesEqual(result, types.NumberValue(42)) {
		t.Errorf("expected 42, got %v", result)
	}
}

func TestEvaluatorDefineWithExpression(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Define a variable with a computed expression
	defineExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "define"},
			&types.SymbolExpr{Name: "y"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.NumberExpr{Value: 10},
					&types.NumberExpr{Value: 20},
				},
			},
		},
	}

	result, err := evaluator.Eval(defineExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(30)) {
		t.Errorf("expected 30, got %v", result)
	}

	// Test that the variable can be used in other expressions
	useExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "*"},
			&types.SymbolExpr{Name: "y"},
			&types.NumberExpr{Value: 2},
		},
	}

	result, err = evaluator.Eval(useExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valuesEqual(result, types.NumberValue(60)) {
		t.Errorf("expected 60, got %v", result)
	}
}

func TestEvaluatorDefineOverwrite(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Define a variable
	defineExpr1 := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "define"},
			&types.SymbolExpr{Name: "z"},
			&types.NumberExpr{Value: 100},
		},
	}

	_, err := evaluator.Eval(defineExpr1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Redefine the same variable with a different value
	defineExpr2 := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "define"},
			&types.SymbolExpr{Name: "z"},
			&types.StringExpr{Value: "hello"},
		},
	}

	result, err := evaluator.Eval(defineExpr2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !valuesEqual(result, types.StringValue("hello")) {
		t.Errorf("expected 'hello', got %v", result)
	}

	// Check that the variable now has the new value
	symbolExpr := &types.SymbolExpr{Name: "z"}
	result, err = evaluator.Eval(symbolExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valuesEqual(result, types.StringValue("hello")) {
		t.Errorf("expected 'hello', got %v", result)
	}
}

func TestEvaluatorDefineErrors(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "too few arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "define"},
					&types.SymbolExpr{Name: "x"},
				},
			},
		},
		{
			name: "too many arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "define"},
					&types.SymbolExpr{Name: "x"},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 2},
				},
			},
		},
		{
			name: "first argument not a symbol",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "define"},
					&types.NumberExpr{Value: 42}, // should be a symbol
					&types.NumberExpr{Value: 1},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluator.Eval(tt.expr)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}

func TestEvaluatorLambda(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test lambda creation
	lambdaExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "lambda"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "y"},
				},
			},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "y"},
				},
			},
		},
	}

	result, err := evaluator.Eval(lambdaExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	function, ok := result.(types.FunctionValue)
	if !ok {
		t.Fatalf("expected FunctionValue, got %T", result)
	}

	if len(function.Params) != 2 {
		t.Errorf("expected 2 parameters, got %d", len(function.Params))
	}

	if function.Params[0] != "x" || function.Params[1] != "y" {
		t.Errorf("unexpected parameters: %v", function.Params)
	}
}

func TestEvaluatorFunctionCall(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Define a function
	function := types.FunctionValue{
		Params: []string{"x"},
		Body: &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "+"},
				&types.SymbolExpr{Name: "x"},
				&types.NumberExpr{Value: 1},
			},
		},
		Env: env,
	}
	env.Set("add1", function)

	// Call the function
	callExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "add1"},
			&types.NumberExpr{Value: 5},
		},
	}

	result, err := evaluator.Eval(callExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(6)) {
		t.Errorf("expected 6, got %v", result)
	}
}

func TestEvaluatorClosure(t *testing.T) {
	env := NewEnvironment()
	env.Set("n", types.NumberValue(10))

	// Create a closure that captures 'n'
	function := types.FunctionValue{
		Params: []string{"x"},
		Body: &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "+"},
				&types.SymbolExpr{Name: "x"},
				&types.SymbolExpr{Name: "n"}, // captured variable
			},
		},
		Env: env, // captured environment
	}

	// Call the function in a new environment that doesn't have 'n'
	newEnv := NewEnvironment()
	newEvaluator := NewEvaluator(newEnv)

	result, err := newEvaluator.callFunction(function, []types.Expr{&types.NumberExpr{Value: 5}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(15)) {
		t.Errorf("expected 15, got %v", result)
	}
}

func TestEvaluatorDefun(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test defun creation
	defunExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "defun"},
			&types.SymbolExpr{Name: "square"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "x"},
				},
			},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "*"},
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "x"},
				},
			},
		},
	}

	result, err := evaluator.Eval(defunExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that the function was created correctly
	function, ok := result.(types.FunctionValue)
	if !ok {
		t.Fatalf("expected FunctionValue, got %T", result)
	}

	if len(function.Params) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(function.Params))
	}

	if function.Params[0] != "x" {
		t.Errorf("unexpected parameter: %v", function.Params[0])
	}

	// Test that the function was defined in the environment
	funcValue, ok := env.Get("square")
	if !ok {
		t.Error("function 'square' not found in environment")
	}

	if !valuesEqual(funcValue, function) {
		t.Error("function in environment doesn't match returned function")
	}

	// Test calling the defined function
	callExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "square"},
			&types.NumberExpr{Value: 5},
		},
	}

	callResult, err := evaluator.Eval(callExpr)
	if err != nil {
		t.Fatalf("unexpected error calling function: %v", err)
	}

	if !valuesEqual(callResult, types.NumberValue(25)) {
		t.Errorf("expected 25, got %v", callResult)
	}
}

func TestEvaluatorDefunMultipleParams(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test defun with multiple parameters
	defunExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "defun"},
			&types.SymbolExpr{Name: "add"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "y"},
				},
			},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "y"},
				},
			},
		},
	}

	_, err := evaluator.Eval(defunExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test calling the function
	callExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "add"},
			&types.NumberExpr{Value: 3},
			&types.NumberExpr{Value: 4},
		},
	}

	result, err := evaluator.Eval(callExpr)
	if err != nil {
		t.Fatalf("unexpected error calling function: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(7)) {
		t.Errorf("expected 7, got %v", result)
	}
}

func TestEvaluatorDefunErrors(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "defun with too few arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "defun"},
					&types.SymbolExpr{Name: "foo"},
					// missing parameters and body
				},
			},
		},
		{
			name: "defun with non-symbol function name",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "defun"},
					&types.NumberExpr{Value: 42}, // should be symbol
					&types.ListExpr{Elements: []types.Expr{}},
					&types.NumberExpr{Value: 1},
				},
			},
		},
		{
			name: "defun with non-list parameters",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "defun"},
					&types.SymbolExpr{Name: "foo"},
					&types.SymbolExpr{Name: "x"}, // should be a list
					&types.NumberExpr{Value: 42},
				},
			},
		},
		{
			name: "defun with non-symbol parameter",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "defun"},
					&types.SymbolExpr{Name: "foo"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.NumberExpr{Value: 42}, // should be symbol
						},
					},
					&types.NumberExpr{Value: 42},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluator.Eval(tt.expr)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}

// TestTailCallOptimization tests that tail recursive functions don't cause stack overflow
func TestTailCallOptimization(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	t.Run("tail recursive factorial", func(t *testing.T) {
		// Define a tail-recursive factorial function
		defunExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "fact-tail"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
						&types.SymbolExpr{Name: "acc"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.SymbolExpr{Name: "acc"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "fact-tail"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "-"},
										&types.SymbolExpr{Name: "n"},
										&types.NumberExpr{Value: 1},
									},
								},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "*"},
										&types.SymbolExpr{Name: "n"},
										&types.SymbolExpr{Name: "acc"},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err := evaluator.Eval(defunExpr)
		if err != nil {
			t.Fatalf("failed to define tail-recursive factorial: %v", err)
		}

		// Test small factorial
		callExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "fact-tail"},
				&types.NumberExpr{Value: 5},
				&types.NumberExpr{Value: 1},
			},
		}

		result, err := evaluator.Eval(callExpr)
		if err != nil {
			t.Fatalf("failed to call tail-recursive factorial: %v", err)
		}

		expected := types.NumberValue(120) // 5! = 120
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("large tail recursive computation should not stack overflow", func(t *testing.T) {
		// Define a tail-recursive sum function
		defunExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "sum-tail"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
						&types.SymbolExpr{Name: "acc"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.SymbolExpr{Name: "acc"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "sum-tail"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "-"},
										&types.SymbolExpr{Name: "n"},
										&types.NumberExpr{Value: 1},
									},
								},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "+"},
										&types.SymbolExpr{Name: "acc"},
										&types.SymbolExpr{Name: "n"},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err := evaluator.Eval(defunExpr)
		if err != nil {
			t.Fatalf("failed to define sum-tail function: %v", err)
		}

		// Test with a large number that would normally cause stack overflow
		// Sum from 1 to 1000: 1000 * 1001 / 2 = 500500
		callExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "sum-tail"},
				&types.NumberExpr{Value: 1000},
				&types.NumberExpr{Value: 0},
			},
		}

		result, err := evaluator.Eval(callExpr)
		if err != nil {
			t.Fatalf("failed to call sum-tail with large number: %v", err)
		}

		expected := types.NumberValue(500500)
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

// TestNonTailRecursion tests that non-tail recursive functions still work correctly
func TestNonTailRecursion(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	t.Run("non-tail recursive factorial", func(t *testing.T) {
		// Define a non-tail-recursive factorial function
		defunExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "factorial"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.NumberExpr{Value: 1},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "*"},
								&types.SymbolExpr{Name: "n"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "factorial"},
										&types.ListExpr{
											Elements: []types.Expr{
												&types.SymbolExpr{Name: "-"},
												&types.SymbolExpr{Name: "n"},
												&types.NumberExpr{Value: 1},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err := evaluator.Eval(defunExpr)
		if err != nil {
			t.Fatalf("failed to define non-tail factorial: %v", err)
		}

		// Test small factorial (this is NOT in tail position)
		callExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "factorial"},
				&types.NumberExpr{Value: 5},
			},
		}

		result, err := evaluator.Eval(callExpr)
		if err != nil {
			t.Fatalf("failed to call non-tail factorial: %v", err)
		}

		expected := types.NumberValue(120) // 5! = 120
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

// TestMutualTailRecursion tests mutually tail-recursive functions
func TestMutualTailRecursion(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	t.Run("mutually recursive even/odd", func(t *testing.T) {
		// Define mutually tail-recursive even? function
		evenDefExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "even?"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.BooleanExpr{Value: true},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "odd?"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "-"},
										&types.SymbolExpr{Name: "n"},
										&types.NumberExpr{Value: 1},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err := evaluator.Eval(evenDefExpr)
		if err != nil {
			t.Fatalf("failed to define even? function: %v", err)
		}

		// Define mutually tail-recursive odd? function
		oddDefExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "odd?"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.BooleanExpr{Value: false},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "even?"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "-"},
										&types.SymbolExpr{Name: "n"},
										&types.NumberExpr{Value: 1},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err = evaluator.Eval(oddDefExpr)
		if err != nil {
			t.Fatalf("failed to define odd? function: %v", err)
		}

		// Test even? with even number
		evenCallExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "even?"},
				&types.NumberExpr{Value: 100},
			},
		}

		result, err := evaluator.Eval(evenCallExpr)
		if err != nil {
			t.Fatalf("failed to call even? function: %v", err)
		}

		expected := types.BooleanValue(true)
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}

		// Test odd? with odd number
		oddCallExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "odd?"},
				&types.NumberExpr{Value: 99},
			},
		}

		result, err = evaluator.Eval(oddCallExpr)
		if err != nil {
			t.Fatalf("failed to call odd? function: %v", err)
		}

		expected = types.BooleanValue(true)
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

// Helper function to compare values
func valuesEqual(a, b types.Value) bool {
	switch va := a.(type) {
	case types.NumberValue:
		if vb, ok := b.(types.NumberValue); ok {
			return math.Abs(float64(va-vb)) < 1e-9
		}
	case types.StringValue:
		if vb, ok := b.(types.StringValue); ok {
			return va == vb
		}
	case types.BooleanValue:
		if vb, ok := b.(types.BooleanValue); ok {
			return va == vb
		}
	case *types.ListValue:
		if vb, ok := b.(*types.ListValue); ok {
			if len(va.Elements) != len(vb.Elements) {
				return false
			}
			for i, elem := range va.Elements {
				if !valuesEqual(elem, vb.Elements[i]) {
					return false
				}
			}
			return true
		}
	case types.FunctionValue:
		if vb, ok := b.(types.FunctionValue); ok {
			// For functions, compare parameter lists and body string representation
			if len(va.Params) != len(vb.Params) {
				return false
			}
			for i, param := range va.Params {
				if param != vb.Params[i] {
					return false
				}
			}
			return va.Body.String() == vb.Body.String()
		}
	}
	return false
}

func TestEvaluatorListOperations(t *testing.T) {
	tests := []struct {
		name     string
		expr     types.Expr
		expected types.Value
	}{
		{
			name: "list creation with no arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "list"},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{}},
		},
		{
			name: "list creation with single element",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "list"},
					&types.NumberExpr{Value: 42},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(42)}},
		},
		{
			name: "list creation with multiple elements",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "list"},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 2},
					&types.NumberExpr{Value: 3},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			}},
		},
		{
			name: "list with mixed types",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "list"},
					&types.NumberExpr{Value: 42},
					&types.StringExpr{Value: "hello"},
					&types.BooleanExpr{Value: true},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{
				types.NumberValue(42),
				types.StringValue("hello"),
				types.BooleanValue(true),
			}},
		},
		{
			name: "empty? on empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "empty?"},
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
			name: "empty? on non-empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "empty?"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
						},
					},
				},
			},
			expected: types.BooleanValue(false),
		},
		{
			name: "length of empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "length"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
						},
					},
				},
			},
			expected: types.NumberValue(0),
		},
		{
			name: "length of non-empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "length"},
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
		},
		{
			name: "first of non-empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "first"},
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
			expected: types.NumberValue(1),
		},
		{
			name: "rest of non-empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "rest"},
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
			expected: &types.ListValue{Elements: []types.Value{
				types.NumberValue(2),
				types.NumberValue(3),
			}},
		},
		{
			name: "rest of single element list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "rest"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{}},
		},
		{
			name: "cons element to list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "cons"},
					&types.NumberExpr{Value: 0},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
							&types.NumberExpr{Value: 2},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{
				types.NumberValue(0),
				types.NumberValue(1),
				types.NumberValue(2),
			}},
		},
		{
			name: "cons element to empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "cons"},
					&types.NumberExpr{Value: 42},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
						},
					},
				},
			},
			expected: &types.ListValue{Elements: []types.Value{
				types.NumberValue(42),
			}},
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

func TestEvaluatorListOperationErrors(t *testing.T) {
	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "first on empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "first"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
						},
					},
				},
			},
		},
		{
			name: "rest on empty list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "rest"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
						},
					},
				},
			},
		},
		{
			name: "first with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "first"},
				},
			},
		},
		{
			name: "rest with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "rest"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "list"},
							&types.NumberExpr{Value: 1},
						},
					},
					&types.NumberExpr{Value: 2}, // extra argument
				},
			},
		},
		{
			name: "cons with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "cons"},
					&types.NumberExpr{Value: 1},
				},
			},
		},
		{
			name: "length with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "length"},
				},
			},
		},
		{
			name: "empty? with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "empty?"},
				},
			},
		},
		{
			name: "first on non-list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "first"},
					&types.NumberExpr{Value: 42},
				},
			},
		},
		{
			name: "rest on non-list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "rest"},
					&types.NumberExpr{Value: 42},
				},
			},
		},
		{
			name: "cons with non-list second argument",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "cons"},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 42},
				},
			},
		},
		{
			name: "length on non-list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "length"},
					&types.NumberExpr{Value: 42},
				},
			},
		},
		{
			name: "empty? on non-list",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "empty?"},
					&types.NumberExpr{Value: 42},
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
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}

func TestEnvironmentInspection(t *testing.T) {
	t.Run("env function shows environment variables", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Set some variables in the environment
		env.Set("x", types.NumberValue(42))
		env.Set("greeting", types.StringValue("hello"))
		env.Set("flag", types.BooleanValue(true))

		// Call (env) to inspect environment
		envExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "env"},
			},
		}

		result, err := evaluator.Eval(envExpr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Result should be a list of (name value) pairs
		listResult, ok := result.(*types.ListValue)
		if !ok {
			t.Fatalf("expected ListValue, got %T", result)
		}

		// Should have 3 pairs (x, greeting, flag) plus built-in functions
		if len(listResult.Elements) < 3 {
			t.Errorf("expected at least 3 environment entries, got %d", len(listResult.Elements))
		}

		// Check that our variables are included
		found := make(map[string]bool)
		for _, elem := range listResult.Elements {
			pair, ok := elem.(*types.ListValue)
			if ok && len(pair.Elements) == 2 {
				if name, ok := pair.Elements[0].(types.StringValue); ok {
					found[string(name)] = true
				}
			}
		}

		if !found["x"] {
			t.Error("variable 'x' not found in environment listing")
		}
		if !found["greeting"] {
			t.Error("variable 'greeting' not found in environment listing")
		}
		if !found["flag"] {
			t.Error("variable 'flag' not found in environment listing")
		}
	})

	t.Run("env function works with empty environment", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		envExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "env"},
			},
		}

		result, err := evaluator.Eval(envExpr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Result should be a list (possibly empty or with only built-ins)
		listResult, ok := result.(*types.ListValue)
		if !ok {
			t.Fatalf("expected ListValue, got %T", result)
		}

		// Should at least return a list structure (even if empty)
		if listResult == nil {
			t.Error("expected non-nil list result")
		}
	})

	t.Run("modules function shows loaded modules", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Create a mock module
		module := &types.ModuleValue{
			Name: "test-module",
			Exports: map[string]types.Value{
				"test-func": &types.FunctionValue{
					Params: []string{"x"},
					Body:   &types.NumberExpr{Value: 42},
					Env:    env,
				},
			},
		}
		env.SetModule("test-module", module)

		// Call (modules) to inspect modules
		modulesExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "modules"},
			},
		}

		result, err := evaluator.Eval(modulesExpr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Result should be a list of (module-name exports) pairs
		listResult, ok := result.(*types.ListValue)
		if !ok {
			t.Fatalf("expected ListValue, got %T", result)
		}

		// Should have at least 1 module entry
		if len(listResult.Elements) < 1 {
			t.Errorf("expected at least 1 module entry, got %d", len(listResult.Elements))
		}

		// Check that our test module is included
		found := false
		for _, elem := range listResult.Elements {
			pair, ok := elem.(*types.ListValue)
			if ok && len(pair.Elements) == 2 {
				if name, ok := pair.Elements[0].(types.StringValue); ok {
					if string(name) == "test-module" {
						found = true
						break
					}
				}
			}
		}

		if !found {
			t.Error("module 'test-module' not found in modules listing")
		}
	})

	t.Run("modules function works with no modules", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		modulesExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "modules"},
			},
		}

		result, err := evaluator.Eval(modulesExpr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Result should be an empty list
		listResult, ok := result.(*types.ListValue)
		if !ok {
			t.Fatalf("expected ListValue, got %T", result)
		}

		if len(listResult.Elements) != 0 {
			t.Errorf("expected empty list for no modules, got %d elements", len(listResult.Elements))
		}
	})

	t.Run("env function shows functions defined with defun", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Define a function using defun
		defunExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "square"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "x"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "*"},
						&types.SymbolExpr{Name: "x"},
						&types.SymbolExpr{Name: "x"},
					},
				},
			},
		}

		_, err := evaluator.Eval(defunExpr)
		if err != nil {
			t.Fatalf("unexpected error defining function: %v", err)
		}

		// Now call (env) to see the defined function
		envExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "env"},
			},
		}

		result, err := evaluator.Eval(envExpr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		listResult, ok := result.(*types.ListValue)
		if !ok {
			t.Fatalf("expected ListValue, got %T", result)
		}

		// Check that the square function is listed
		found := false
		for _, elem := range listResult.Elements {
			pair, ok := elem.(*types.ListValue)
			if ok && len(pair.Elements) == 2 {
				if name, ok := pair.Elements[0].(types.StringValue); ok {
					if string(name) == "square" {
						found = true
						break
					}
				}
			}
		}

		if !found {
			t.Error("function 'square' not found in environment listing")
		}
	})

	t.Run("builtins function shows all built-in functions", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Call (builtins) to get list of built-in functions
		builtinsExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "builtins"},
			},
		}

		result, err := evaluator.Eval(builtinsExpr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Result should be a list of string values
		listResult, ok := result.(*types.ListValue)
		if !ok {
			t.Fatalf("expected ListValue, got %T", result)
		}

		// Check that we have a reasonable number of built-in functions
		if len(listResult.Elements) < 20 {
			t.Errorf("expected at least 20 built-in functions, got %d", len(listResult.Elements))
		}

		// Check that some essential built-in functions are included
		builtinMap := make(map[string]bool)
		for _, elem := range listResult.Elements {
			if name, ok := elem.(types.StringValue); ok {
				builtinMap[string(name)] = true
			}
		}

		essentialBuiltins := []string{"+", "-", "*", "/", "if", "define", "lambda", "defun", "list", "env", "modules"}
		for _, builtin := range essentialBuiltins {
			if !builtinMap[builtin] {
				t.Errorf("essential built-in function '%s' not found in builtins listing", builtin)
			}
		}

		// Check that 'builtins' itself is included (meta!)
		if !builtinMap["builtins"] {
			t.Error("'builtins' function should include itself in the listing")
		}
	})

	t.Run("builtins function shows help for specific functions", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Test help for reduce function
		builtinsHelpExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "builtins"},
				&types.SymbolExpr{Name: "reduce"},
			},
		}

		result, err := evaluator.Eval(builtinsHelpExpr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Result should be a string with help text
		helpText, ok := result.(types.StringValue)
		if !ok {
			t.Fatalf("expected StringValue, got %T", result)
		}

		helpStr := string(helpText)
		// Check that the help contains key information
		if !strings.Contains(helpStr, "reduce") {
			t.Error("help text should contain 'reduce'")
		}
		if !strings.Contains(helpStr, "func") {
			t.Error("help text should contain 'func'")
		}
		if !strings.Contains(helpStr, "Example") {
			t.Error("help text should contain an example")
		}

		// Test help for a simple arithmetic function
		plusHelpExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "builtins"},
				&types.SymbolExpr{Name: "+"},
			},
		}

		result, err = evaluator.Eval(plusHelpExpr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		helpText, ok = result.(types.StringValue)
		if !ok {
			t.Fatalf("expected StringValue, got %T", result)
		}

		helpStr = string(helpText)
		if !strings.Contains(helpStr, "+") {
			t.Error("help text should contain '+'")
		}
		if !strings.Contains(helpStr, "Addition") {
			t.Error("help text should contain 'Addition'")
		}
	})

	t.Run("builtins function fails for unknown functions", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Test help for non-existent function
		unknownHelpExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "builtins"},
				&types.SymbolExpr{Name: "unknown-function"},
			},
		}

		_, err := evaluator.Eval(unknownHelpExpr)
		if err == nil {
			t.Error("expected error for unknown function, but got none")
		}

		expectedError := "no help available for 'unknown-function'"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got: %v", expectedError, err)
		}
	})

	t.Run("builtins function fails with too many arguments", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Test with too many arguments
		tooManyArgsExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "builtins"},
				&types.SymbolExpr{Name: "reduce"},
				&types.SymbolExpr{Name: "extra"},
			},
		}

		_, err := evaluator.Eval(tooManyArgsExpr)
		if err == nil {
			t.Error("expected error for too many arguments, but got none")
		}

		expectedError := "builtins requires 0 or 1 arguments"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got: %v", expectedError, err)
		}
	})
}
