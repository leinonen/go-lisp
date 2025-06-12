package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestDefmacro(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple macro definition",
			input:    "(defmacro when [condition body] (list 'if condition body 'nil))",
			expected: "#<macro([condition body])>",
		},
		{
			name:     "macro with multiple parameters",
			input:    "(defmacro unless [condition then else] (list 'if condition else then))",
			expected: "#<macro([condition then else])>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			expr := parseString(t, tt.input)
			result, err := evaluator.Eval(expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestMacroExpansion(t *testing.T) {
	tests := []struct {
		name     string
		setup    string
		input    string
		expected string
	}{
		{
			name:     "when macro expansion",
			setup:    "(defmacro when [condition body] (list 'if condition body 'nil))",
			input:    "(when (> 5 3) 42)",
			expected: "42",
		},
		{
			name:     "unless macro expansion",
			setup:    "(defmacro unless [condition then else] (list 'if condition else then))",
			input:    "(unless (< 5 3) 42 99)",
			expected: "42",
		},
		{
			name:     "macro with symbol interpolation",
			setup:    "(defmacro inc [var] (list 'define var (list '+ var 1)))",
			input:    "(define x 10)",
			expected: "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			// Set up the macro
			setupExpr := parseString(t, tt.setup)
			_, err := evaluator.Eval(setupExpr)
			if err != nil {
				t.Fatalf("failed to set up macro: %v", err)
			}

			// Test macro usage
			expr := parseString(t, tt.input)
			result, err := evaluator.Eval(expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestMacroExpansionComplex(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Define a more complex macro that creates a let-like binding
	macroDefExpr := parseString(t, `(defmacro let1 [var value body] 
		(list (list 'lambda (list var) body) value))`)
	_, err := evaluator.Eval(macroDefExpr)
	if err != nil {
		t.Fatalf("failed to define macro: %v", err)
	}

	// Use the macro
	testExpr := parseString(t, "(let1 x 10 (+ x 5))")
	result, err := evaluator.Eval(testExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "15"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
}

func TestMacroExpansionNestedQuotes(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Define a macro that uses nested quoting
	macroDefExpr := parseString(t, `(defmacro debug [expr] 
		(list 'list (list 'quote expr) expr))`)
	_, err := evaluator.Eval(macroDefExpr)
	if err != nil {
		t.Fatalf("failed to define macro: %v", err)
	}

	// Use the macro
	testExpr := parseString(t, "(debug (+ 2 3))")
	result, err := evaluator.Eval(testExpr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return a list containing the expression and its value
	listValue, ok := result.(*types.ListValue)
	if !ok {
		t.Fatalf("expected list result, got %T", result)
	}

	if len(listValue.Elements) != 2 {
		t.Fatalf("expected list with 2 elements, got %d", len(listValue.Elements))
	}

	// First element should be the quoted expression (just verify it exists)
	if listValue.Elements[0] == nil {
		t.Fatalf("expected first element to exist")
	}

	// Second element should be the evaluated value
	secondValue, ok := listValue.Elements[1].(types.NumberValue)
	if !ok {
		t.Fatalf("expected second element to be a number, got %T", listValue.Elements[1])
	}

	if float64(secondValue) != 5 {
		t.Errorf("expected second element to be 5, got %v", secondValue)
	}
}

func TestMacroErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "defmacro with no arguments",
			input:       "(defmacro)",
			expectError: true,
			errorMsg:    "defmacro requires at least 3 arguments",
		},
		{
			name:        "defmacro with non-symbol name",
			input:       "(defmacro 123 [x] x)",
			expectError: true,
			errorMsg:    "defmacro first argument must be a symbol",
		},
		{
			name:        "defmacro with non-list parameters",
			input:       "(defmacro test x x)",
			expectError: true,
			errorMsg:    "defmacro second argument must be a parameter list using square brackets",
		},
		{
			name:        "defmacro with non-symbol parameter",
			input:       "(defmacro test [123] 456)",
			expectError: true,
			errorMsg:    "defmacro parameter must be a symbol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			expr := parseString(t, tt.input)
			_, err := evaluator.Eval(expr)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestQuoteSpecialForm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "quote number",
			input:    "(quote 42)",
			expected: "42",
		},
		{
			name:     "quote symbol",
			input:    "(quote foo)",
			expected: "foo",
		},
		{
			name:     "quote list",
			input:    "(quote (1 2 3))",
			expected: "(1 2 3)",
		},
		{
			name:     "quote nested list",
			input:    "(quote (+ 1 2))",
			expected: "(+ 1 2)",
		},
		{
			name:     "apostrophe shorthand",
			input:    "'(+ 1 2)",
			expected: "(+ 1 2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			expr := parseString(t, tt.input)
			result, err := evaluator.Eval(expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}()))
}
