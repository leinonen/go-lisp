package interpreter

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:     "simple number",
			input:    "42",
			expected: types.NumberValue(42),
		},
		{
			name:     "simple addition",
			input:    "(+ 1 2)",
			expected: types.NumberValue(3),
		},
		{
			name:     "nested expression",
			input:    "(+ (* 2 3) 4)",
			expected: types.NumberValue(10),
		},
		{
			name:     "boolean expression",
			input:    "(< 3 5)",
			expected: types.BooleanValue(true),
		},
		{
			name:     "string literal",
			input:    `"hello world"`,
			expected: types.StringValue("hello world"),
		},
		{
			name:     "if expression",
			input:    "(if (< 3 5) 42 0)",
			expected: types.NumberValue(42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()
			result, err := interpreter.Interpret(tt.input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInterpreterDefine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:     "simple variable definition",
			input:    "(define x 42)",
			expected: types.NumberValue(42),
		},
		{
			name:     "define with expression",
			input:    "(define y (+ 10 20))",
			expected: types.NumberValue(30),
		},
		{
			name:     "define string variable",
			input:    `(define greeting "hello world")`,
			expected: types.StringValue("hello world"),
		},
		{
			name:     "define boolean variable",
			input:    "(define flag #t)",
			expected: types.BooleanValue(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()
			result, err := interpreter.Interpret(tt.input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInterpreterDefineAndUse(t *testing.T) {
	interpreter := NewInterpreter()

	// Define a variable
	_, err := interpreter.Interpret("(define x 10)")
	if err != nil {
		t.Fatalf("unexpected error defining x: %v", err)
	}

	// Use the variable
	result, err := interpreter.Interpret("x")
	if err != nil {
		t.Fatalf("unexpected error accessing x: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(10)) {
		t.Errorf("expected 10, got %v", result)
	}

	// Use the variable in an expression
	result, err = interpreter.Interpret("(+ x 5)")
	if err != nil {
		t.Fatalf("unexpected error in expression: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(15)) {
		t.Errorf("expected 15, got %v", result)
	}

	// Define another variable using the first
	result, err = interpreter.Interpret("(define y (* x 3))")
	if err != nil {
		t.Fatalf("unexpected error defining y: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(30)) {
		t.Errorf("expected 30, got %v", result)
	}

	// Use both variables
	result, err = interpreter.Interpret("(+ x y)")
	if err != nil {
		t.Fatalf("unexpected error in final expression: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(40)) {
		t.Errorf("expected 40, got %v", result)
	}
}

func TestInterpreterFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:  "lambda expression",
			input: "(lambda (x) (+ x 1))",
			expected: types.FunctionValue{
				Params: []string{"x"},
				Body: &types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "+"},
						&types.SymbolExpr{Name: "x"},
						&types.NumberExpr{Value: 1},
					},
				},
			},
		},
		{
			name:  "lambda with multiple parameters",
			input: "(lambda (x y) (* x y))",
			expected: types.FunctionValue{
				Params: []string{"x", "y"},
				Body: &types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "*"},
						&types.SymbolExpr{Name: "x"},
						&types.SymbolExpr{Name: "y"},
					},
				},
			},
		},
		{
			name:  "lambda with no parameters",
			input: "(lambda () 42)",
			expected: types.FunctionValue{
				Params: []string{},
				Body:   &types.NumberExpr{Value: 42},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()
			result, err := interpreter.Interpret(tt.input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// For function values, we need a special comparison
			if !functionsEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInterpreterFunctionCalls(t *testing.T) {
	tests := []struct {
		name     string
		setup    []string // expressions to run before the test
		input    string
		expected types.Value
	}{
		{
			name:     "simple function call",
			setup:    []string{"(define add1 (lambda (x) (+ x 1)))"},
			input:    "(add1 5)",
			expected: types.NumberValue(6),
		},
		{
			name:     "function with multiple parameters",
			setup:    []string{"(define multiply (lambda (x y) (* x y)))"},
			input:    "(multiply 3 4)",
			expected: types.NumberValue(12),
		},
		{
			name:     "function with no parameters",
			setup:    []string{"(define get-answer (lambda () 42))"},
			input:    "(get-answer)",
			expected: types.NumberValue(42),
		},
		{
			name: "nested function calls",
			setup: []string{
				"(define add1 (lambda (x) (+ x 1)))",
				"(define double (lambda (x) (* x 2)))",
			},
			input:    "(double (add1 5))",
			expected: types.NumberValue(12),
		},
		{
			name: "function using outer variables",
			setup: []string{
				"(define y 10)",
				"(define add-y (lambda (x) (+ x y)))",
			},
			input:    "(add-y 5)",
			expected: types.NumberValue(15),
		},
		{
			name:     "recursive function",
			setup:    []string{"(define factorial (lambda (n) (if (= n 0) 1 (* n (factorial (- n 1))))))"},
			input:    "(factorial 5)",
			expected: types.NumberValue(120),
		},
		{
			name: "higher-order function",
			setup: []string{
				"(define apply-twice (lambda (f x) (f (f x))))",
				"(define add1 (lambda (x) (+ x 1)))",
			},
			input:    "(apply-twice add1 5)",
			expected: types.NumberValue(7),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()

			// Run setup expressions
			for _, setupExpr := range tt.setup {
				_, err := interpreter.Interpret(setupExpr)
				if err != nil {
					t.Fatalf("unexpected error in setup %q: %v", setupExpr, err)
				}
			}

			// Run the test expression
			result, err := interpreter.Interpret(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInterpreterFunctionErrors(t *testing.T) {
	tests := []struct {
		name  string
		setup []string
		input string
	}{
		{
			name:  "lambda with wrong argument count",
			input: "(lambda)",
		},
		{
			name:  "lambda with non-list parameters",
			input: "(lambda x (+ x 1))",
		},
		{
			name:  "lambda with non-symbol in parameter list",
			input: "(lambda (x 42) (+ x 1))",
		},
		{
			name:  "function call with wrong number of arguments - too few",
			setup: []string{"(define add (lambda (x y) (+ x y)))"},
			input: "(add 5)",
		},
		{
			name:  "function call with wrong number of arguments - too many",
			setup: []string{"(define add1 (lambda (x) (+ x 1)))"},
			input: "(add1 5 6)",
		},
		{
			name:  "calling non-function",
			setup: []string{"(define x 42)"},
			input: "(x)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()

			// Run setup expressions
			for _, setupExpr := range tt.setup {
				_, err := interpreter.Interpret(setupExpr)
				if err != nil {
					t.Fatalf("unexpected error in setup %q: %v", setupExpr, err)
				}
			}

			// Run the test expression - should produce an error
			_, err := interpreter.Interpret(tt.input)
			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
		})
	}
}

func TestInterpreterClosure(t *testing.T) {
	interpreter := NewInterpreter()

	// Define a function that returns a closure
	_, err := interpreter.Interpret("(define make-adder (lambda (n) (lambda (x) (+ x n))))")
	if err != nil {
		t.Fatalf("unexpected error defining make-adder: %v", err)
	}

	// Create an adder function
	_, err = interpreter.Interpret("(define add5 (make-adder 5))")
	if err != nil {
		t.Fatalf("unexpected error creating add5: %v", err)
	}

	// Use the closure
	result, err := interpreter.Interpret("(add5 10)")
	if err != nil {
		t.Fatalf("unexpected error calling add5: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(15)) {
		t.Errorf("expected 15, got %v", result)
	}

	// Create another adder with different captured value
	_, err = interpreter.Interpret("(define add10 (make-adder 10))")
	if err != nil {
		t.Fatalf("unexpected error creating add10: %v", err)
	}

	result, err = interpreter.Interpret("(add10 7)")
	if err != nil {
		t.Fatalf("unexpected error calling add10: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(17)) {
		t.Errorf("expected 17, got %v", result)
	}
}

// Helper function to compare values
func valuesEqual(a, b types.Value) bool {
	switch va := a.(type) {
	case types.NumberValue:
		if vb, ok := b.(types.NumberValue); ok {
			return va == vb
		}
	case types.StringValue:
		if vb, ok := b.(types.StringValue); ok {
			return va == vb
		}
	case types.BooleanValue:
		if vb, ok := b.(types.BooleanValue); ok {
			return va == vb
		}
	}
	return false
}

// Helper function to compare function values
func functionsEqual(a, b types.Value) bool {
	fa, ok1 := a.(types.FunctionValue)
	fb, ok2 := b.(types.FunctionValue)

	if !ok1 || !ok2 {
		return valuesEqual(a, b)
	}

	// Compare parameters count
	if len(fa.Params) != len(fb.Params) {
		return false
	}

	// Compare parameter names
	for i, param := range fa.Params {
		if param != fb.Params[i] {
			return false
		}
	}

	// For simplicity, we'll assume bodies are equal if they have the same string representation
	// In a more sophisticated implementation, we'd do structural comparison
	return fa.Body.String() == fb.Body.String()
}

func TestInterpreterComplexFunctionExample(t *testing.T) {
	interpreter := NewInterpreter()

	// Test a complex example with multiple function features
	expressions := []struct {
		input    string
		expected types.Value
	}{
		// Define helper functions
		{"(define square (lambda (x) (* x x)))", types.NumberValue(0)}, // Function value comparison not straightforward
		{"(define add (lambda (x y) (+ x y)))", types.NumberValue(0)},

		// Test basic function calls
		{"(square 4)", types.NumberValue(16)},
		{"(add 3 7)", types.NumberValue(10)},

		// Test higher-order functions
		{"(define apply-twice (lambda (f x) (f (f x))))", types.NumberValue(0)},
		{"(define increment (lambda (x) (+ x 1)))", types.NumberValue(0)},
		{"(apply-twice increment 5)", types.NumberValue(7)},
		{"(apply-twice square 2)", types.NumberValue(16)}, // square(square(2)) = square(4) = 16

		// Test closures
		{"(define make-multiplier (lambda (n) (lambda (x) (* x n))))", types.NumberValue(0)},
		{"(define double (make-multiplier 2))", types.NumberValue(0)},
		{"(define triple (make-multiplier 3))", types.NumberValue(0)},
		{"(double 5)", types.NumberValue(10)},
		{"(triple 4)", types.NumberValue(12)},

		// Test recursion
		{"(define sum-to (lambda (n) (if (= n 0) 0 (+ n (sum-to (- n 1))))))", types.NumberValue(0)},
		{"(sum-to 5)", types.NumberValue(15)}, // 1+2+3+4+5 = 15
	}

	for i, expr := range expressions {
		result, err := interpreter.Interpret(expr.input)
		if err != nil {
			t.Fatalf("step %d: unexpected error for %q: %v", i+1, expr.input, err)
		}

		// For function definitions, we can't easily compare the result, so we skip detailed comparison
		if _, ok := result.(types.FunctionValue); ok {
			continue
		}

		if !valuesEqual(result, expr.expected) {
			t.Errorf("step %d: for input %q, expected %v, got %v", i+1, expr.input, expr.expected, result)
		}
	}
}
