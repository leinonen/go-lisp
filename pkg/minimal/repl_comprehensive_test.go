package minimal

import (
	"strings"
	"testing"
)

// Helper function to evaluate a string expression in a REPL
func (r *REPL) Eval(input string) (Value, error) {
	parsed, err := r.Parse(input)
	if err != nil {
		return nil, err
	}
	return Eval(parsed, r.Env)
}

// Helper function to create a bootstrapped REPL
func newBootstrappedREPL() *REPL {
	repl := NewREPL()
	Bootstrap(repl.Env)
	return repl
}

// TestREPLComprehensive provides comprehensive testing of REPL functionality
func TestREPLComprehensive(t *testing.T) {
	t.Run("BasicREPLOperations", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test that REPL has arithmetic functions
		result, err := repl.Eval("(+ 1 2 3)")
		if err != nil {
			t.Fatalf("Error evaluating arithmetic: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 6.0 {
			t.Errorf("Expected 6, got %v", result)
		}

		// Test that REPL maintains state between evaluations
		_, err = repl.Eval("(def x 42)")
		if err != nil {
			t.Fatalf("Error defining variable: %v", err)
		}

		result, err = repl.Eval("x")
		if err != nil {
			t.Fatalf("Error accessing defined variable: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected x to be 42, got %v", result)
		}
	})

	t.Run("ParseMethod", func(t *testing.T) {
		repl := NewREPL()

		// Test parsing numbers
		result, err := repl.Parse("42")
		if err != nil {
			t.Fatalf("Error parsing number: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected Number 42, got %v", result)
		}

		// Test parsing strings
		result, err = repl.Parse(`"hello"`)
		if err != nil {
			t.Fatalf("Error parsing string: %v", err)
		}

		if str, ok := result.(String); !ok || string(str) != "hello" {
			t.Errorf("Expected String 'hello', got %v", result)
		}

		// Test parsing lists
		result, err = repl.Parse("(+ 1 2)")
		if err != nil {
			t.Fatalf("Error parsing list: %v", err)
		}

		if list, ok := result.(*List); !ok || list.Length() != 3 {
			t.Errorf("Expected List of length 3, got %v", result)
		}

		// Test parsing error
		_, err = repl.Parse("(unclosed")
		if err == nil {
			t.Error("Expected parse error for unclosed expression")
		}
	})

	t.Run("EvalMethod", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test basic evaluation
		result, err := repl.Eval("(* 6 7)")
		if err != nil {
			t.Fatalf("Error evaluating: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected 42, got %v", result)
		}

		// Test evaluation error
		_, err = repl.Eval("(+ 1 \"not-a-number\")")
		if err == nil {
			t.Error("Expected evaluation error for type mismatch")
		}

		// Test undefined symbol
		_, err = repl.Eval("undefined-symbol")
		if err == nil {
			t.Error("Expected error for undefined symbol")
		}
	})

	t.Run("IsBalanced", func(t *testing.T) {
		repl := NewREPL()

		testCases := []struct {
			input    string
			expected bool
		}{
			{"", true},
			{"42", true},
			{"(+ 1 2)", true},
			{"(+ 1 2", false},
			{"+ 1 2)", false},
			{"(+ (* 2 3) (- 5 1))", true},
			{"(+ (* 2 3) (- 5 1)", false},
			{"[1 2 3]", true},
			{"[1 2 3", false},
			{"1 2 3]", false},
			{"(+ [1 2] 3)", true},
			{"(+ [1 2 3)", false},
			{`"hello"`, true},
			{`"hello`, false},
			{`"hello (with parens)"`, true},
			{`"hello (with parens`, false},
			{`(print "hello")`, true},
			{`(print "hello"`, false},
			{"((()))", true},
			{"((()", false},
			{"[[[]]]", true},
			{"[[[", false},
			{"(list [1 2] [3 4])", true},
			{"(list [1 2] [3 4)", false},
		}

		for _, tc := range testCases {
			result := repl.isBalanced(tc.input)
			if result != tc.expected {
				t.Errorf("isBalanced(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		}
	})

	t.Run("StringHandlingInBalance", func(t *testing.T) {
		repl := NewREPL()

		// Test that parentheses inside strings don't affect balance
		stringCases := []struct {
			input    string
			expected bool
		}{
			{`"("`, true},
			{`")"`, true},
			{`"()"`, true},
			{`"((("`, true},
			{`"))))"`, true},
			{`(print "hello (world)")`, true},
			{`(print "hello (world)"`, false},
			{`(print "hello )world(")`, true},
			{`"with \" escaped quote"`, true},
			{`"with \" escaped quote`, false},
		}

		for _, tc := range stringCases {
			result := repl.isBalanced(tc.input)
			if result != tc.expected {
				t.Errorf("isBalanced(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		}
	})

	t.Run("NewBootstrappedREPL", func(t *testing.T) {
		repl := newBootstrappedREPL()

		// Should have all basic functions
		basicFunctions := []string{"+", "-", "*", "list", "first", "rest", "=", "<"}

		for _, fn := range basicFunctions {
			_, err := repl.Env.Get(Intern(fn))
			if err != nil {
				t.Errorf("Bootstrapped REPL missing function: %s", fn)
			}
		}

		// Should be able to evaluate complex expressions
		result, err := repl.Eval("(first (list 1 2 3))")
		if err != nil {
			t.Fatalf("Error evaluating in bootstrapped REPL: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected 1 from bootstrapped REPL, got %v", result)
		}
	})

	t.Run("ErrorReporting", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test various error conditions
		errorCases := []struct {
			input     string
			shouldErr bool
			checkErr  func(error) bool
		}{
			{
				input:     "(+ 1 \"string\")",
				shouldErr: true,
				checkErr: func(err error) bool {
					return strings.Contains(strings.ToLower(err.Error()), "number")
				},
			},
			{
				input:     "undefined-variable",
				shouldErr: true,
				checkErr: func(err error) bool {
					return strings.Contains(err.Error(), "undefined")
				},
			},
			{
				input:     "(42 1 2)", // trying to call number as function
				shouldErr: true,
				checkErr: func(err error) bool {
					return strings.Contains(err.Error(), "function")
				},
			},
			{
				input:     "(first)", // wrong number of args
				shouldErr: true,
				checkErr: func(err error) bool {
					return strings.Contains(err.Error(), "argument")
				},
			},
		}

		for _, tc := range errorCases {
			result, err := repl.Eval(tc.input)

			if tc.shouldErr {
				if err == nil {
					t.Errorf("Expected error for input %q, but got result: %v", tc.input, result)
				} else if !tc.checkErr(err) {
					t.Errorf("Error check failed for input %q: %v", tc.input, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tc.input, err)
			}
		}
	})

	t.Run("MultilineInput", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test that multiline input works correctly
		multilineInput := `(def factorial
  (fn [n]
    (if (= n 0)
        1
        (* n (factorial (- n 1))))))`

		result, err := repl.Eval(multilineInput)
		if err != nil {
			t.Fatalf("Error evaluating multiline input: %v", err)
		}

		if _, ok := result.(DefinedValue); !ok {
			t.Errorf("Expected DefinedValue, got %T", result)
		}

		// Test that the function works
		factResult, err := repl.Eval("(factorial 5)")
		if err != nil {
			t.Fatalf("Error calling defined function: %v", err)
		}

		if num, ok := factResult.(Number); !ok || float64(num) != 120.0 {
			t.Errorf("Expected factorial(5) = 120, got %v", factResult)
		}
	})

	t.Run("ComplexNestedExpressions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test deeply nested expressions - use a simpler nested expression since we don't have division
		simpleNested := `(+ (* 2 3) (- 10 5) (* 1 3))`

		result, err := repl.Eval(simpleNested)
		if err != nil {
			t.Fatalf("Error evaluating nested expression: %v", err)
		}

		// (+ 6 5 3) = 14
		if num, ok := result.(Number); !ok || float64(num) != 14.0 {
			t.Errorf("Expected 14, got %v", result)
		}
	})

	t.Run("EnvironmentIsolation", func(t *testing.T) {
		// Test that different REPL instances have isolated environments
		repl1 := NewREPL()
		repl2 := NewREPL()
		Bootstrap(repl1.Env)
		Bootstrap(repl2.Env)

		// Define variable in repl1
		_, err := repl1.Eval("(def test-var 42)")
		if err != nil {
			t.Fatalf("Error defining variable in repl1: %v", err)
		}

		// repl2 should not see this variable
		_, err = repl2.Eval("test-var")
		if err == nil {
			t.Error("repl2 should not see variable defined in repl1")
		}

		// Define different value in repl2
		_, err = repl2.Eval("(def test-var 100)")
		if err != nil {
			t.Fatalf("Error defining variable in repl2: %v", err)
		}

		// Check that repl1 still has its value
		result1, err := repl1.Eval("test-var")
		if err != nil {
			t.Fatalf("Error accessing variable in repl1: %v", err)
		}

		if num, ok := result1.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected repl1 test-var to be 42, got %v", result1)
		}

		// Check that repl2 has its value
		result2, err := repl2.Eval("test-var")
		if err != nil {
			t.Fatalf("Error accessing variable in repl2: %v", err)
		}

		if num, ok := result2.(Number); !ok || float64(num) != 100.0 {
			t.Errorf("Expected repl2 test-var to be 100, got %v", result2)
		}
	})
}
