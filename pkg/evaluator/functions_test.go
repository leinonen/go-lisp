package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

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
