package evaluator

import (
	"strings"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestEvaluatorDef(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test basic variable definition
	defineExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "def"},
			&types.SymbolExpr{Name: "x"},
			&types.NumberExpr{Value: 42},
		},
	}

	result, err := evaluator.Eval(defineExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// def should return the value that was defined
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

func TestEvaluatorDefWithExpression(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Define a variable with a computed expression
	defineExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "def"},
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

func TestEvaluatorDefOverwrite(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Define a variable
	defineExpr1 := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "def"},
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
			&types.SymbolExpr{Name: "def"},
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

func TestEvaluatorDefErrors(t *testing.T) {
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
					&types.SymbolExpr{Name: "def"},
					&types.SymbolExpr{Name: "x"},
				},
			},
		},
		{
			name: "too many arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "def"},
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
					&types.SymbolExpr{Name: "def"},
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

		// Define a function using defn
		defnExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defn"},
				&types.SymbolExpr{Name: "square"},
				&types.BracketExpr{
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

		_, err := evaluator.Eval(defnExpr)
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

	t.Run("help function shows all built-in functions", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Call (help) to get list of built-in functions
		helpExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "help"},
			},
		}

		result, err := evaluator.Eval(helpExpr)
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

		essentialBuiltins := []string{"+", "-", "*", "/", "if", "def", "fn", "defn", "list", "env", "modules"}
		for _, builtin := range essentialBuiltins {
			if !builtinMap[builtin] {
				t.Errorf("essential built-in function '%s' not found in help listing", builtin)
			}
		}

		// Check that 'help' itself is included (meta!)
		if !builtinMap["help"] {
			t.Error("'help' function should include itself in the listing")
		}
	})

	t.Run("help function shows help for specific functions", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Test help for reduce function
		helpHelpExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "help"},
				&types.SymbolExpr{Name: "reduce"},
			},
		}

		result, err := evaluator.Eval(helpHelpExpr)
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
				&types.SymbolExpr{Name: "help"},
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

	t.Run("help function fails for unknown functions", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Test help for non-existent function
		unknownHelpExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "help"},
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

	t.Run("help function fails with too many arguments", func(t *testing.T) {
		env := NewEnvironment()
		evaluator := NewEvaluator(env)

		// Test with too many arguments
		tooManyArgsExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "help"},
				&types.SymbolExpr{Name: "reduce"},
				&types.SymbolExpr{Name: "extra"},
			},
		}

		_, err := evaluator.Eval(tooManyArgsExpr)
		if err == nil {
			t.Error("expected error for too many arguments, but got none")
		}

		expectedError := "help requires 0 or 1 arguments"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got: %v", expectedError, err)
		}
	})
}

// TestEvalError tests the error function
func TestEvalError(t *testing.T) {
	tests := []struct {
		name        string
		expr        types.Expr
		expectError bool
		errorMsg    string
	}{
		{
			name: "error with string message",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "error"},
					&types.StringExpr{Value: "This is a test error"},
				},
			},
			expectError: true,
			errorMsg:    "This is a test error",
		},
		{
			name: "error with number message",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "error"},
					&types.NumberExpr{Value: 42},
				},
			},
			expectError: true,
			errorMsg:    "42",
		},
		{
			name: "error with boolean message",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "error"},
					&types.BooleanExpr{Value: true},
				},
			},
			expectError: true,
			errorMsg:    "#t",
		},
		{
			name: "error with string concatenation that works",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "error"},
					&types.StringExpr{Value: "Error code: 404"},
				},
			},
			expectError: true,
			errorMsg:    "Error code: 404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			_, err := evaluator.Eval(tt.expr)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestEvalErrorEdgeCases tests error cases for the error function
func TestEvalErrorEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "error with no arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "error"},
				},
			},
		},
		{
			name: "error with too many arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "error"},
					&types.StringExpr{Value: "message1"},
					&types.StringExpr{Value: "message2"},
				},
			},
		},
		{
			name: "error with invalid expression argument",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "error"},
					&types.SymbolExpr{Name: "undefined_symbol"},
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
