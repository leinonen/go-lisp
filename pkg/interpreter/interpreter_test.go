package interpreter

import (
	"math"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

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

// Helper function to compare functions specifically
func functionsEqual(a, b types.Value) bool {
	return valuesEqual(a, b)
}

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

func TestInterpreterListOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:     "empty list creation",
			input:    "(list)",
			expected: &types.ListValue{Elements: []types.Value{}},
		},
		{
			name:     "single element list",
			input:    "(list 42)",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(42)}},
		},
		{
			name:     "multi-element list",
			input:    "(list 1 2 3)",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(1), types.NumberValue(2), types.NumberValue(3)}},
		},
		{
			name:     "mixed type list",
			input:    `(list 42 "hello" #t)`,
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(42), types.StringValue("hello"), types.BooleanValue(true)}},
		},
		{
			name:     "empty check on empty list",
			input:    "(empty? (list))",
			expected: types.BooleanValue(true),
		},
		{
			name:     "empty check on non-empty list",
			input:    "(empty? (list 1))",
			expected: types.BooleanValue(false),
		},
		{
			name:     "length of empty list",
			input:    "(length (list))",
			expected: types.NumberValue(0),
		},
		{
			name:     "length of non-empty list",
			input:    "(length (list 1 2 3))",
			expected: types.NumberValue(3),
		},
		{
			name:     "first of list",
			input:    "(first (list 1 2 3))",
			expected: types.NumberValue(1),
		},
		{
			name:     "rest of list",
			input:    "(rest (list 1 2 3))",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(2), types.NumberValue(3)}},
		},
		{
			name:     "rest of single element list",
			input:    "(rest (list 42))",
			expected: &types.ListValue{Elements: []types.Value{}},
		},
		{
			name:     "cons to list",
			input:    "(cons 0 (list 1 2))",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(0), types.NumberValue(1), types.NumberValue(2)}},
		},
		{
			name:     "cons to empty list",
			input:    "(cons 42 (list))",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(42)}},
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

func TestInterpreterListOperationsComplex(t *testing.T) {
	interpreter := NewInterpreter()

	// Test complex list operations with variables and functions
	expressions := []struct {
		input    string
		expected types.Value
	}{
		// Define a list
		{"(define my-list (list 1 2 3))", &types.ListValue{Elements: []types.Value{types.NumberValue(1), types.NumberValue(2), types.NumberValue(3)}}},

		// Test operations on the defined list
		{"(length my-list)", types.NumberValue(3)},
		{"(first my-list)", types.NumberValue(1)},
		{"(rest my-list)", &types.ListValue{Elements: []types.Value{types.NumberValue(2), types.NumberValue(3)}}},

		// Build a new list using cons
		{"(define extended-list (cons 0 my-list))", &types.ListValue{Elements: []types.Value{types.NumberValue(0), types.NumberValue(1), types.NumberValue(2), types.NumberValue(3)}}},
		{"(length extended-list)", types.NumberValue(4)},

		// Nested list operations
		{"(first (rest extended-list))", types.NumberValue(1)},
		{"(rest (rest extended-list))", &types.ListValue{Elements: []types.Value{types.NumberValue(2), types.NumberValue(3)}}},

		// List with expressions
		{"(list (+ 1 2) (* 3 4) (if #t 5 6))", &types.ListValue{Elements: []types.Value{types.NumberValue(3), types.NumberValue(12), types.NumberValue(5)}}},
	}

	for i, expr := range expressions {
		result, err := interpreter.Interpret(expr.input)
		if err != nil {
			t.Fatalf("step %d: unexpected error for %q: %v", i+1, expr.input, err)
		}

		if !valuesEqual(result, expr.expected) {
			t.Errorf("step %d: for input %q, expected %v, got %v", i+1, expr.input, expr.expected, result)
		}
	}
}

func TestInterpreterListsWithFunctions(t *testing.T) {
	interpreter := NewInterpreter()

	// Test lists with functions
	expressions := []struct {
		input    string
		expected types.Value
	}{
		// Define a function that works with lists
		{"(define list-add1 (lambda (lst) (cons (+ (first lst) 1) (rest lst))))", types.NumberValue(0)}, // Function definition

		// Test the function
		{"(list-add1 (list 5 10 15))", &types.ListValue{Elements: []types.Value{types.NumberValue(6), types.NumberValue(10), types.NumberValue(15)}}},

		// Define a function that creates lists
		{"(define make-range (lambda (n) (if (= n 0) (list) (cons n (make-range (- n 1))))))", types.NumberValue(0)},
		{"(make-range 3)", &types.ListValue{Elements: []types.Value{types.NumberValue(3), types.NumberValue(2), types.NumberValue(1)}}},

		// Function that processes list recursively
		{"(define sum-list (lambda (lst) (if (empty? lst) 0 (+ (first lst) (sum-list (rest lst))))))", types.NumberValue(0)},
		{"(sum-list (list 1 2 3 4))", types.NumberValue(10)},
	}

	for i, expr := range expressions {
		result, err := interpreter.Interpret(expr.input)
		if err != nil {
			t.Fatalf("step %d: unexpected error for %q: %v", i+1, expr.input, err)
		}

		// Skip function value comparison
		if _, ok := result.(types.FunctionValue); ok {
			continue
		}

		if !valuesEqual(result, expr.expected) {
			t.Errorf("step %d: for input %q, expected %v, got %v", i+1, expr.input, expr.expected, result)
		}
	}
}

func TestMapFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:     "map with square function",
			input:    "(map (lambda (x) (* x x)) (list 1 2 3 4))",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(1), types.NumberValue(4), types.NumberValue(9), types.NumberValue(16)}},
		},
		{
			name:     "map with add one function",
			input:    "(map (lambda (x) (+ x 1)) (list 10 20 30))",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(11), types.NumberValue(21), types.NumberValue(31)}},
		},
		{
			name:     "map with empty list",
			input:    "(map (lambda (x) (* x 2)) (list))",
			expected: &types.ListValue{Elements: []types.Value{}},
		},
		{
			name:     "map with string transformation",
			input:    "(map (lambda (s) s) (list \"a\" \"b\" \"c\"))",
			expected: &types.ListValue{Elements: []types.Value{types.StringValue("a"), types.StringValue("b"), types.StringValue("c")}},
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

func TestFilterFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:     "filter positive numbers",
			input:    "(filter (lambda (x) (> x 0)) (list -1 2 -3 4 5))",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(2), types.NumberValue(4), types.NumberValue(5)}},
		},
		{
			name:     "filter numbers greater than 3",
			input:    "(filter (lambda (x) (> x 3)) (list 1 2 3 4 5 6))",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(4), types.NumberValue(5), types.NumberValue(6)}},
		},
		{
			name:     "filter with empty list",
			input:    "(filter (lambda (x) (> x 0)) (list))",
			expected: &types.ListValue{Elements: []types.Value{}},
		},
		{
			name:     "filter all elements match",
			input:    "(filter (lambda (x) (> x 0)) (list 1 2 3))",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(1), types.NumberValue(2), types.NumberValue(3)}},
		},
		{
			name:     "filter no elements match",
			input:    "(filter (lambda (x) (< x 0)) (list 1 2 3))",
			expected: &types.ListValue{Elements: []types.Value{}},
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

func TestReduceFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:     "reduce sum with initial value",
			input:    "(reduce (lambda (acc x) (+ acc x)) 0 (list 1 2 3 4))",
			expected: types.NumberValue(10),
		},
		{
			name:     "reduce product with initial value",
			input:    "(reduce (lambda (acc x) (* acc x)) 1 (list 2 3 4))",
			expected: types.NumberValue(24),
		},
		{
			name:     "reduce with empty list",
			input:    "(reduce (lambda (acc x) (+ acc x)) 0 (list))",
			expected: types.NumberValue(0),
		},
		{
			name:     "reduce with single element",
			input:    "(reduce (lambda (acc x) (+ acc x)) 10 (list 5))",
			expected: types.NumberValue(15),
		},
		{
			name:     "reduce with lambda function",
			input:    "(reduce (lambda (acc x) (+ acc (* x x))) 0 (list 1 2 3))",
			expected: types.NumberValue(14), // 0 + 1² + 2² + 3² = 14
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

func TestHigherOrderFunctionCombinations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:     "map then reduce",
			input:    "(reduce (lambda (acc x) (+ acc x)) 0 (map (lambda (x) (* x x)) (list 1 2 3)))",
			expected: types.NumberValue(14), // sum of squares: 1 + 4 + 9 = 14
		},
		{
			name:     "filter then map",
			input:    "(map (lambda (x) (* x 2)) (filter (lambda (x) (> x 0)) (list -1 2 -3 4)))",
			expected: &types.ListValue{Elements: []types.Value{types.NumberValue(4), types.NumberValue(8)}},
		},
		{
			name:     "filter then reduce",
			input:    "(reduce (lambda (acc x) (+ acc x)) 0 (filter (lambda (x) (> x 0)) (list -1 2 -3 4 5)))",
			expected: types.NumberValue(11), // 2 + 4 + 5 = 11
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
