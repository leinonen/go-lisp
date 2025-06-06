package evaluator

import (
	"math"
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

func TestEvaluatorLambdaErrors(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "lambda with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "lambda"},
					&types.ListExpr{Elements: []types.Expr{}}, // empty params
					// missing body
				},
			},
		},
		{
			name: "lambda with non-list parameters",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "lambda"},
					&types.SymbolExpr{Name: "x"}, // should be a list
					&types.NumberExpr{Value: 42},
				},
			},
		},
		{
			name: "lambda with non-symbol parameter",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "lambda"},
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

func TestEvaluatorFunctionCallErrors(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Define a function that expects 2 arguments
	function := types.FunctionValue{
		Params: []string{"x", "y"},
		Body: &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "+"},
				&types.SymbolExpr{Name: "x"},
				&types.SymbolExpr{Name: "y"},
			},
		},
		Env: env,
	}
	env.Set("add2", function)

	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "too few arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "add2"},
					&types.NumberExpr{Value: 5},
					// missing second argument
				},
			},
		},
		{
			name: "too many arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "add2"},
					&types.NumberExpr{Value: 5},
					&types.NumberExpr{Value: 10},
					&types.NumberExpr{Value: 15}, // extra argument
				},
			},
		},
		{
			name: "calling undefined function",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "undefined"},
					&types.NumberExpr{Value: 5},
				},
			},
		},
		{
			name: "calling non-function value",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "notafunc"},
					&types.NumberExpr{Value: 5},
				},
			},
		},
	}

	// Set a non-function value for the last test
	env.Set("notafunc", types.NumberValue(42))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluator.Eval(tt.expr)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}
