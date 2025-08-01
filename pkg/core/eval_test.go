package core_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/leinonen/go-lisp/pkg/core"
)

func TestEvalBasicValues(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		{"42", "42"},
		{"-42", "-42"},
		{"3.14", "3.14"},
		{"\"hello\"", "\"hello\""},
		{":keyword", ":keyword"},
		{"nil", "nil"},
		{"true", "true"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestEvalArithmetic(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		{"(+ 1 2)", "3"},
		{"(+ 1 2 3)", "6"},
		{"(+)", "0"},
		{"(- 5 3)", "2"},
		{"(- 10)", "-10"},
		{"(* 2 3)", "6"},
		{"(* 2 3 4)", "24"},
		{"(*)", "1"},
		{"(/ 6 2)", "3"},
		{"(/ 10 2 2)", "2.5"},
		{"(+ 1.5 2.5)", "4"},
		{"(* 2.5 4)", "10"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestEvalComparison(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		{"(= 1 1)", "true"},
		{"(= 1 2)", "nil"},
		{"(= 1 1 1)", "true"},
		{"(= 1 1 2)", "nil"},
		{"(< 1 2)", "true"},
		{"(< 2 1)", "nil"},
		{"(> 2 1)", "true"},
		{"(> 1 2)", "nil"},
		{"(= \"hello\" \"hello\")", "true"},
		{"(= \"hello\" \"world\")", "nil"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestEvalListOperations(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		{"(cons 1 nil)", "(1 nil)"},
		{"(cons 1 (cons 2 nil))", "(1 2 nil)"},
		{"(first (cons 1 nil))", "1"},
		{"(first nil)", "nil"},
		{"(rest (cons 1 (cons 2 nil)))", "(2 nil)"},
		{"(rest nil)", "()"},
		{"(first [1 2 3])", "1"},
		{"(first [])", "nil"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestEvalTypePredicates(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		{"(symbol? 'x)", "true"},
		{"(symbol? 42)", "nil"},
		{"(string? \"hello\")", "true"},
		{"(string? 42)", "nil"},
		{"(number? 42)", "true"},
		{"(number? \"hello\")", "nil"},
		{"(list? '(1 2 3))", "true"},
		{"(list? [1 2 3])", "nil"},
		{"(vector? [1 2 3])", "true"},
		{"(vector? '(1 2 3))", "nil"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestEvalSpecialForms(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Test quote
	expr, _ := core.ReadString("(quote x)")
	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for quote: %v", err)
	}
	if result.String() != "x" {
		t.Errorf("Expected 'x' for quote, got '%s'", result.String())
	}

	// Test shorthand quote
	expr, _ = core.ReadString("'x")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for shorthand quote: %v", err)
	}
	if result.String() != "x" {
		t.Errorf("Expected 'x' for shorthand quote, got '%s'", result.String())
	}

	// Test if - true case
	expr, _ = core.ReadString("(if true 1 2)")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for if true: %v", err)
	}
	if result.String() != "1" {
		t.Errorf("Expected '1' for if true, got '%s'", result.String())
	}

	// Test if - false case
	expr, _ = core.ReadString("(if nil 1 2)")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for if false: %v", err)
	}
	if result.String() != "2" {
		t.Errorf("Expected '2' for if false, got '%s'", result.String())
	}

	// Test if - no else
	expr, _ = core.ReadString("(if nil 1)")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for if no else: %v", err)
	}
	if result.String() != "nil" {
		t.Errorf("Expected 'nil' for if no else, got '%s'", result.String())
	}

	// Test do
	expr, _ = core.ReadString("(do 1 2 3)")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for do: %v", err)
	}
	if result.String() != "3" {
		t.Errorf("Expected '3' for do, got '%s'", result.String())
	}
}

func TestEvalDefAndSymbolLookup(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Test def
	expr, _ := core.ReadString("(def x 42)")
	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for def: %v", err)
	}
	if result.String() != "x" {
		t.Errorf("Expected 'x' for def, got '%s'", result.String())
	}

	// Test symbol lookup
	expr, _ = core.ReadString("x")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for symbol lookup: %v", err)
	}
	if result.String() != "42" {
		t.Errorf("Expected '42' for symbol lookup, got '%s'", result.String())
	}

	// Test undefined symbol
	expr, _ = core.ReadString("undefined")
	_, err = core.Eval(expr, env)
	if err == nil {
		t.Error("Expected error for undefined symbol")
	}
}

func TestEvalFunctionDefinitionAndCall(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Define a function
	expr, _ := core.ReadString("(def add (fn [a b] (+ a b)))")
	_, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for function definition: %v", err)
	}

	// Call the function
	expr, _ = core.ReadString("(add 3 4)")
	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for function call: %v", err)
	}
	if result.String() != "7" {
		t.Errorf("Expected '7' for function call, got '%s'", result.String())
	}

	// Test function with vector parameters
	expr, _ = core.ReadString("(def mult (fn [x y] (* x y)))")
	_, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for function definition with vector: %v", err)
	}

	expr, _ = core.ReadString("(mult 5 6)")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for function call with vector params: %v", err)
	}
	if result.String() != "30" {
		t.Errorf("Expected '30' for function call, got '%s'", result.String())
	}
}

func TestEvalRecursiveFunction(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Define factorial function
	expr, _ := core.ReadString("(def factorial (fn [n] (if (= n 0) 1 (* n (factorial (- n 1))))))")
	_, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for recursive function definition: %v", err)
	}

	// Test factorial(5)
	expr, _ = core.ReadString("(factorial 5)")
	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for recursive function call: %v", err)
	}
	if result.String() != "120" {
		t.Errorf("Expected '120' for factorial(5), got '%s'", result.String())
	}

	// Test factorial(0)
	expr, _ = core.ReadString("(factorial 0)")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for factorial(0): %v", err)
	}
	if result.String() != "1" {
		t.Errorf("Expected '1' for factorial(0), got '%s'", result.String())
	}
}

func TestEvalLexicalScoping(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Define outer variable
	expr, _ := core.ReadString("(def x 10)")
	_, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for outer variable: %v", err)
	}

	// Define function that captures outer variable
	expr, _ = core.ReadString("(def make-adder (fn [y] (fn [z] (+ x y z))))")
	_, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for closure definition: %v", err)
	}

	// Create adder function
	expr, _ = core.ReadString("(def add-5 (make-adder 5))")
	_, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for creating closure: %v", err)
	}

	// Call the closure
	expr, _ = core.ReadString("(add-5 3)")
	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for calling closure: %v", err)
	}
	if result.String() != "18" {
		t.Errorf("Expected '18' for closure call, got '%s'", result.String())
	}
}

func TestEvalMetaProgramming(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Test eval
	expr, _ := core.ReadString("(eval '(+ 1 2))")
	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for eval: %v", err)
	}
	if result.String() != "3" {
		t.Errorf("Expected '3' for eval, got '%s'", result.String())
	}

	// Test read-string
	expr, _ = core.ReadString("(read-string \"(+ 1 2)\")")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for read-string: %v", err)
	}
	if result.String() != "(+ 1 2)" {
		t.Errorf("Expected '(+ 1 2)' for read-string, got '%s'", result.String())
	}

	// Test eval + read-string
	expr, _ = core.ReadString("(eval (read-string \"(+ 1 2)\"))")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for eval + read-string: %v", err)
	}
	if result.String() != "3" {
		t.Errorf("Expected '3' for eval + read-string, got '%s'", result.String())
	}

	// Test gensym with default prefix
	expr, _ = core.ReadString("(gensym)")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for gensym: %v", err)
	}
	if result.String() != "G__1" {
		t.Errorf("Expected 'G__1' for first gensym, got '%s'", result.String())
	}

	// Test gensym with custom prefix
	expr, _ = core.ReadString("(gensym \"my-prefix\")")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for gensym with prefix: %v", err)
	}
	if result.String() != "my-prefix2" {
		t.Errorf("Expected 'my-prefix2' for gensym with prefix, got '%s'", result.String())
	}

	// Test gensym uniqueness
	expr, _ = core.ReadString("(gensym)")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for third gensym: %v", err)
	}
	if result.String() != "G__3" {
		t.Errorf("Expected 'G__3' for third gensym, got '%s'", result.String())
	}
}

func TestEvalMacroExpand(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Define a simple macro
	expr, _ := core.ReadString("(defmacro when [condition body] `(if ~condition ~body nil))")
	_, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Error defining when macro: %v", err)
	}

	// Test macroexpand on the macro
	expr, _ = core.ReadString("(macroexpand '(when true (+ 1 2)))")
	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for macroexpand: %v", err)
	}
	expected := "(if true (+ 1 2) nil)"
	if result.String() != expected {
		t.Errorf("Expected '%s' for macroexpand, got '%s'", expected, result.String())
	}

	// Test macroexpand on non-macro (should return unchanged)
	expr, _ = core.ReadString("(macroexpand '(+ 1 2))")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for macroexpand on non-macro: %v", err)
	}
	if result.String() != "(+ 1 2)" {
		t.Errorf("Expected '(+ 1 2)' for macroexpand on non-macro, got '%s'", result.String())
	}
}

func TestEvalVariadicFunctions(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Test variadic function with no extra args
	expr, _ := core.ReadString("(def vfn (fn [x & rest] (list x rest)))")
	_, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Error defining variadic function: %v", err)
	}

	// Call with minimum args
	expr, _ = core.ReadString("(vfn 1)")
	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for variadic function with min args: %v", err)
	}
	if result.String() != "(1 ())" {
		t.Errorf("Expected '(1 ())' for variadic function with min args, got '%s'", result.String())
	}

	// Call with extra args
	expr, _ = core.ReadString("(vfn 1 2 3 4)")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for variadic function with extra args: %v", err)
	}
	if result.String() != "(1 (2 3 4))" {
		t.Errorf("Expected '(1 (2 3 4))' for variadic function with extra args, got '%s'", result.String())
	}

	// Test variadic macro
	expr, _ = core.ReadString("(defmacro when [condition & body] `(if ~condition (do ~@body) nil))")
	_, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Error defining when macro: %v", err)
	}

	// Test macroexpand on when macro
	expr, _ = core.ReadString("(macroexpand '(when true (println \"hello\") (+ 1 2)))")
	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for when macro expansion: %v", err)
	}
	expected := "(if true (do (println \"hello\") (+ 1 2)) nil)"
	if result.String() != expected {
		t.Errorf("Expected '%s' for when macro expansion, got '%s'", expected, result.String())
	}

	// Test unless macro (should be available from stdlib)
	env2, err := core.CreateBootstrappedEnvironment()
	if err != nil {
		t.Errorf("Error creating bootstrapped environment: %v", err)
	}

	// Test unless macro expansion
	expr, _ = core.ReadString("(macroexpand '(unless false (println \"hello\") (+ 1 2)))")
	result, err = core.Eval(expr, env2)
	if err != nil {
		t.Errorf("Eval error for unless macro expansion: %v", err)
	}
	expected = "(if false nil (do (println \"hello\") (+ 1 2)))"
	if result.String() != expected {
		t.Errorf("Expected '%s' for unless macro expansion, got '%s'", expected, result.String())
	}
}

func TestEvalErrorReporting(t *testing.T) {
	// Test parse errors with location information
	parseErrorTests := []struct {
		input    string
		contains []string // Substrings that should be in the error message
	}{
		{"(+ 1 2", []string{"ParseError", "line 1", "column"}},
		{"(def x \"unterminated", []string{"unterminated string", "line 1", "column"}},
		{")", []string{"unexpected token", "line 1", "column"}},
	}

	for _, test := range parseErrorTests {
		_, err := core.ReadString(test.input)
		if err == nil {
			t.Errorf("Expected parse error for input '%s', but got none", test.input)
			continue
		}

		errorMsg := err.Error()
		for _, substr := range test.contains {
			if !strings.Contains(errorMsg, substr) {
				t.Errorf("Expected error message to contain '%s' for input '%s', got: %s", substr, test.input, errorMsg)
			}
		}
	}
}

func TestEvalErrors(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []string{
		"(+ 1 \"hello\")",    // Type error
		"(/ 1 0)",            // Division by zero
		"(unknown-function)", // Unknown function
		"(def)",              // Wrong number of arguments
		"(fn)",               // Wrong number of arguments
		"(if)",               // Wrong number of arguments
		"(quote)",            // Wrong number of arguments
		"(= 1)",              // Wrong number of arguments for =
		"(< 1)",              // Wrong number of arguments for <
		"(> 1)",              // Wrong number of arguments for >
		"(cons 1)",           // Wrong number of arguments for cons
		"(first)",            // Wrong number of arguments for first
		"(rest)",             // Wrong number of arguments for rest
	}

	for _, test := range tests {
		expr, err := core.ReadString(test)
		if err != nil {
			continue // Skip parse errors for this test
		}

		_, err = core.Eval(expr, env)
		if err == nil {
			t.Errorf("Expected error for '%s', but got none", test)
		}
	}
}

func TestUserFunctionInterface(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Create a user function
	params := core.NewList(core.Intern("x"))
	body, _ := core.ReadString("(+ x 1)")
	userFn := &core.UserFunction{
		Params: params,
		Body:   body,
		Env:    env,
	}

	// Test Call method
	args := []core.Value{core.NewNumber(int64(5))}
	result, err := userFn.Call(args, env)
	if err != nil {
		t.Errorf("Unexpected error calling user function: %v", err)
	}
	if result.String() != "6" {
		t.Errorf("Expected '6', got '%s'", result.String())
	}

	// Test String method
	if userFn.String() != "#<function>" {
		t.Errorf("Expected '#<function>', got '%s'", userFn.String())
	}

	// Test wrong number of arguments
	wrongArgs := []core.Value{core.NewNumber(int64(1)), core.NewNumber(int64(2))}
	_, err = userFn.Call(wrongArgs, env)
	if err == nil {
		t.Error("Expected error for wrong number of arguments")
	}
}

func TestBuiltinFunctionInterface(t *testing.T) {
	// Test a builtin function
	addFn := &core.BuiltinFunction{
		Name: "test-add",
		Fn: func(args []core.Value, env *core.Environment) (core.Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("test-add expects 2 arguments")
			}
			n1, _ := args[0].(core.Number)
			n2, _ := args[1].(core.Number)
			return core.NewNumber(n1.ToInt() + n2.ToInt()), nil
		},
	}

	// Test Call method
	args := []core.Value{core.NewNumber(int64(3)), core.NewNumber(int64(4))}
	result, err := addFn.Call(args, nil)
	if err != nil {
		t.Errorf("Unexpected error calling builtin function: %v", err)
	}
	if result.String() != "7" {
		t.Errorf("Expected '7', got '%s'", result.String())
	}

	// Test String method
	if addFn.String() != "#<builtin:test-add>" {
		t.Errorf("Expected '#<builtin:test-add>', got '%s'", addFn.String())
	}
}

func TestNewCoreFunctions(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// Test list function
		{"(list)", "()"},
		{"(list 1)", "(1)"},
		{"(list 1 2 3)", "(1 2 3)"},
		{"(list \"a\" \"b\" \"c\")", "(\"a\" \"b\" \"c\")"},

		// Test count function
		{"(count (list 1 2 3))", "3"},
		{"(count [])", "0"},
		{"(count [1 2 3 4])", "4"},
		{"(count nil)", "0"},

		// Test empty? function
		{"(empty? (list))", "true"},
		{"(empty? [])", "true"},
		{"(empty? (list 1))", "nil"},
		{"(empty? [1])", "nil"},
		{"(empty? nil)", "true"},

		// Test nth function
		{"(nth [1 2 3] 0)", "1"},
		{"(nth [1 2 3] 1)", "2"},
		{"(nth [1 2 3] 2)", "3"},
		{"(nth (list 1 2 3) 1)", "2"},

		// Test conj function
		{"(conj [1 2] 3)", "[1 2 3]"},
		{"(conj [] 1)", "[1]"},
		{"(conj (list 1 2) 3)", "(3 1 2)"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestStringOperations(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// Test str function
		{"(str \"hello\")", "\"hello\""},
		{"(str \"hello\" \" \" \"world\")", "\"hello world\""},
		{"(str 1 2 3)", "\"123\""},
		{"(str)", "\"\""},

		// Test substring function
		{"(substring \"hello\" 1 4)", "\"ell\""},
		{"(substring \"world\" 0 5)", "\"world\""},
		{"(substring \"test\" 2 4)", "\"st\""},

		// Test string-split function
		{"(string-split \"a,b,c\" \",\")", "[\"a\" \"b\" \"c\"]"},
		{"(string-split \"hello world\" \" \")", "[\"hello\" \"world\"]"},
		{"(string-split \"test\" \",\")", "[\"test\"]"},

		// Test string-trim function
		{"(string-trim \"  hello  \")", "\"hello\""},
		{"(string-trim \"\\n\\ttest\\n\")", "\"test\""},
		{"(string-trim \"normal\")", "\"normal\""},

		// Test string-replace function
		{"(string-replace \"hello world\" \"world\" \"universe\")", "\"hello universe\""},
		{"(string-replace \"test test\" \"test\" \"demo\")", "\"demo demo\""},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestLetSpecialForm(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// Basic let bindings
		{"(let [x 1] x)", "1"},
		{"(let [x 1 y 2] (+ x y))", "3"},
		{"(let [x 10] (let [y 20] (+ x y)))", "30"},

		// Let with function calls
		{"(let [x (+ 1 2)] (* x 3))", "9"},

		// Let with multiple expressions in body
		{"(let [x 1] x x)", "1"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestFileSystemOperations(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Test file-exists? with existing file (from project root)
	expr, err := core.ReadString("(file-exists? \"../../README.md\")")
	if err != nil {
		t.Errorf("Parse error for file-exists?: %v", err)
		return
	}

	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for file-exists?: %v", err)
		return
	}

	if result.String() != "true" {
		t.Errorf("Expected 'true' for existing file, got '%s'", result.String())
	}

	// Test file-exists? with non-existing file
	expr, err = core.ReadString("(file-exists? \"nonexistent-file.txt\")")
	if err != nil {
		t.Errorf("Parse error for file-exists? non-existing: %v", err)
		return
	}

	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for file-exists? non-existing: %v", err)
		return
	}

	if result.String() != "nil" {
		t.Errorf("Expected 'nil' for non-existing file, got '%s'", result.String())
	}

	// Test list-dir with current directory
	expr, err = core.ReadString("(list-dir \"../..\")")
	if err != nil {
		t.Errorf("Parse error for list-dir: %v", err)
		return
	}

	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for list-dir: %v", err)
		return
	}

	// Check that result is a vector and contains expected files
	if vector, ok := result.(*core.Vector); ok {
		found := false
		for i := 0; i < vector.Count(); i++ {
			item := vector.Get(i)
			if str, ok := item.(core.String); ok && string(str) == "README.md" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected README.md to be in directory listing")
		}
	} else {
		t.Errorf("Expected vector result from list-dir, got %T", result)
	}
}

func TestCoreFunctionsOnly(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Test only core functions that don't require stdlib loading
	tests := []struct {
		input    string
		expected string
	}{
		// Test that core environment is working
		{"(+ 1 2 3)", "6"},
		{"(* 2 3 4)", "24"},
		{"(str \"hello\" \" \" \"world\")", "\"hello world\""},

		// Test that string operations work
		{"(string-split \"a,b,c\" \",\")", "[\"a\" \"b\" \"c\"]"},
		{"(string-trim \"  hello  \")", "\"hello\""},
		{"(string-replace \"hello world\" \"world\" \"test\")", "\"hello test\""},

		// Test type predicates
		{"(string? \"hello\")", "true"},
		{"(number? 42)", "true"},
		{"(symbol? 'test)", "true"},

		// Test collection operations
		{"(cons 1 '(2 3))", "(1 2 3)"},
		{"(first '(1 2 3))", "1"},
		{"(rest '(1 2 3))", "(2 3)"},
		{"(count [1 2 3 4])", "4"},
		{"(empty? [])", "true"},
		{"(empty? [1])", "nil"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestIOOperations(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Test println - capture output would require more complex setup
	// For now just test that it doesn't error
	expr, err := core.ReadString("(println \"test\")")
	if err != nil {
		t.Errorf("Parse error for println: %v", err)
		return
	}

	result, err := core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for println: %v", err)
		return
	}

	// println should return nil
	if result.String() != "nil" {
		t.Errorf("Expected 'nil' for println return, got '%s'", result.String())
	}

	// Test prn similarly
	expr, err = core.ReadString("(prn \"test\")")
	if err != nil {
		t.Errorf("Parse error for prn: %v", err)
		return
	}

	result, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for prn: %v", err)
		return
	}

	if result.String() != "nil" {
		t.Errorf("Expected 'nil' for prn return, got '%s'", result.String())
	}
}

func TestHashMapOperations(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// Test hash-map constructor
		{"(hash-map)", "{}"},
		{"(hash-map :name \"Alice\")", "{:name \"Alice\"}"},
		{"(hash-map :name \"Alice\" :age 30)", "{:name \"Alice\" :age 30}"},

		// Test hash-map literal syntax
		{"{}", "{}"},
		{"{:name \"Bob\"}", "{:name \"Bob\"}"},
		{"{:name \"Bob\" :age 25}", "{:name \"Bob\" :age 25}"},

		// Test get function
		{"(get {:name \"Alice\" :age 30} :name)", "\"Alice\""},
		{"(get {:name \"Alice\" :age 30} :age)", "30"},
		{"(get {:name \"Alice\"} :nonexistent)", "nil"},
		{"(get {:name \"Alice\"} :nonexistent \"default\")", "\"default\""},

		// Test assoc function
		{"(assoc {} :key \"value\")", "{:key \"value\"}"},
		{"(assoc {:a 1} :b 2)", "{:a 1 :b 2}"},
		{"(assoc {:a 1} :a 2)", "{:a 2}"},

		// Test dissoc function
		{"(dissoc {:a 1 :b 2} :a)", "{:b 2}"},
		{"(dissoc {:a 1 :b 2 :c 3} :b)", "{:a 1 :c 3}"},
		{"(dissoc {:a 1} :nonexistent)", "{:a 1}"},

		// Test contains? function
		{"(contains? {:name \"Alice\"} :name)", "true"},
		{"(contains? {:name \"Alice\"} :age)", "nil"},

		// Test hash-map? predicate
		{"(hash-map? {})", "true"},
		{"(hash-map? {:a 1})", "true"},
		{"(hash-map? [])", "nil"},
		{"(hash-map? \"test\")", "nil"},

		// Test count with hash-map
		{"(count {})", "0"},
		{"(count {:a 1})", "1"},
		{"(count {:a 1 :b 2 :c 3})", "3"},

		// Test empty? with hash-map
		{"(empty? {})", "true"},
		{"(empty? {:a 1})", "nil"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestSetOperations(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// Test set constructor
		{"(set)", "#{}"},
		{"(set 1)", "#{1}"},
		{"(set 1 2 3)", "#{1 2 3}"},
		{"(set 1 2 3 2 1)", "#{1 2 3}"},

		// Test set literal syntax
		{"#{}", "#{}"},
		{"#{1}", "#{1}"},
		{"#{1 2 3}", "#{1 2 3}"},

		// Test contains? function with sets
		{"(contains? #{1 2 3} 1)", "true"},
		{"(contains? #{1 2 3} 2)", "true"},
		{"(contains? #{1 2 3} 4)", "nil"},

		// Test set? predicate
		{"(set? #{})", "true"},
		{"(set? #{1 2 3})", "true"},
		{"(set? [])", "nil"},
		{"(set? {})", "nil"},

		// Test count with set
		{"(count #{})", "0"},
		{"(count #{1})", "1"},
		{"(count #{1 2 3})", "3"},

		// Test empty? with set
		{"(empty? #{})", "true"},
		{"(empty? #{1})", "nil"},

		// Test union operation
		{"(union #{1 2} #{3 4})", "#{1 2 3 4}"},
		{"(union #{1 2} #{2 3})", "#{1 2 3}"},
		{"(union #{} #{1 2})", "#{1 2}"},
		{"(union #{1 2} #{})", "#{1 2}"},
		{"(union #{1 2} #{2 3} #{3 4})", "#{1 2 3 4}"},

		// Test intersection operation
		{"(intersection #{1 2 3} #{2 3 4})", "#{2 3}"},
		{"(intersection #{1 2} #{3 4})", "#{}"},
		{"(intersection #{1 2 3} #{1 2 3})", "#{1 2 3}"},
		{"(intersection #{1 2 3} #{2} #{2 4})", "#{2}"},
		{"(intersection #{1 2 3} #{2 3} #{2 3 4})", "#{2 3}"},

		// Test difference operation
		{"(difference #{1 2 3} #{2})", "#{1 3}"},
		{"(difference #{1 2 3} #{2 3})", "#{1}"},
		{"(difference #{1 2 3} #{4 5})", "#{1 2 3}"},
		{"(difference #{1 2 3} #{1 2 3})", "#{}"},
		{"(difference #{1 2 3 4} #{2} #{3})", "#{1 4}"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestKeywordAsFunction(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// Test keyword function calls on hash-maps
		{"(:name {:name \"Alice\" :age 30})", "\"Alice\""},
		{"(:age {:name \"Alice\" :age 30})", "30"},
		{"(:nonexistent {:name \"Alice\"})", "nil"},
		{"(:nonexistent {:name \"Alice\"} \"default\")", "\"default\""},

		// Test with different data types as keys
		{"(:key {:key 42})", "42"},
		{"(:flag {:flag true})", "true"},

		// Test with nested structures
		{"(:user {:user {:name \"Bob\"}})", "{:name \"Bob\"}"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestReadAllString(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
		count    int
		desc     string
	}{
		{
			input:    "(read-all-string \"(+ 1 2)\")",
			expected: "((+ 1 2))",
			count:    1,
			desc:     "single expression",
		},
		{
			input:    "(read-all-string \"(+ 1 2) (* 3 4)\")",
			expected: "((+ 1 2) (* 3 4))",
			count:    2,
			desc:     "two expressions",
		},
		{
			input:    "(read-all-string \"(def x 10) (def y 20) (+ x y)\")",
			expected: "((def x 10) (def y 20) (+ x y))",
			count:    3,
			desc:     "three expressions with definitions",
		},
		{
			input:    "(read-all-string \"[1 2 3] {:a 1} #{1 2}\")",
			expected: "([1 2 3] {:a 1} #{1 2})",
			count:    3,
			desc:     "data structures",
		},
		{
			input:    "(read-all-string \"(defn f [x] (* x x)) (f 5)\")",
			expected: "((defn f [x] (* x x)) (f 5))",
			count:    2,
			desc:     "function definition and call",
		},
		{
			input:    "(read-all-string \"\")",
			expected: "()",
			count:    0,
			desc:     "empty string",
		},
	}

	for _, test := range tests {
		// Test the result
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for %s '%s': %v", test.desc, test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for %s '%s': %v", test.desc, test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for %s, got '%s'", test.expected, test.desc, result.String())
		}

		// Test the count
		countExpr, err := core.ReadString(fmt.Sprintf("(count %s)", test.input))
		if err != nil {
			t.Errorf("Parse error for count test %s: %v", test.desc, err)
			continue
		}

		countResult, err := core.Eval(countExpr, env)
		if err != nil {
			t.Errorf("Eval error for count test %s: %v", test.desc, err)
			continue
		}

		expectedCount := fmt.Sprintf("%d", test.count)
		if countResult.String() != expectedCount {
			t.Errorf("Expected count %s for %s, got '%s'", expectedCount, test.desc, countResult.String())
		}
	}
}

func TestReadAllStringErrors(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input string
		desc  string
	}{
		{"(read-all-string)", "no arguments"},
		{"(read-all-string \"(+ 1 2)\" \"extra\")", "too many arguments"},
		{"(read-all-string 123)", "non-string argument"},
		{"(read-all-string \"(unclosed list\")", "syntax error"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			// This is expected for some syntax errors
			continue
		}

		_, err = core.Eval(expr, env)
		if err == nil {
			t.Errorf("Expected error for %s '%s', but got none", test.desc, test.input)
		}
	}
}

func TestLoadFile(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Create a temporary test file
	testContent := `(def test-var 42)
(def test-fn (fn [x] (* x 2)))
(def test-result (test-fn test-var))`

	// Write test file
	testFile := "/tmp/test-load-file.lisp"
	expr, err := core.ReadString(fmt.Sprintf("(spit \"%s\" \"%s\")", testFile, testContent))
	if err != nil {
		t.Errorf("Parse error creating test file: %v", err)
		return
	}
	_, err = core.Eval(expr, env)
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
		return
	}

	// Test load-file
	loadExpr, err := core.ReadString(fmt.Sprintf("(load-file \"%s\")", testFile))
	if err != nil {
		t.Errorf("Parse error for load-file: %v", err)
		return
	}

	result, err := core.Eval(loadExpr, env)
	if err != nil {
		t.Errorf("Error loading file: %v", err)
		return
	}

	// The result should be the last expression's symbol (test-result)
	if result.String() != "test-result" {
		t.Errorf("Expected 'test-result' from loaded file, got '%s'", result.String())
	}

	// Test that the variables are now defined in the environment
	tests := []struct {
		input    string
		expected string
		desc     string
	}{
		{"test-var", "42", "loaded variable"},
		{"(test-fn 10)", "20", "loaded function"},
		{"test-result", "84", "loaded computation result"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for %s: %v", test.desc, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for %s: %v", test.desc, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for %s, got '%s'", test.expected, test.desc, result.String())
		}
	}
}

func TestLoadFileErrors(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input string
		desc  string
	}{
		{"(load-file)", "no arguments"},
		{"(load-file \"file1\" \"file2\")", "too many arguments"},
		{"(load-file 123)", "non-string argument"},
		{"(load-file \"/nonexistent/file.lisp\")", "nonexistent file"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			continue
		}

		_, err = core.Eval(expr, env)
		if err == nil {
			t.Errorf("Expected error for %s '%s', but got none", test.desc, test.input)
		}
	}
}

func TestLogicalOperations(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// Test and with no arguments
		{"(and)", "true"},
		// Test and with single argument
		{"(and true)", "true"},
		{"(and nil)", "nil"},
		{"(and 42)", "42"},
		// Test and with multiple arguments - returns last value if all truthy
		{"(and true true)", "true"},
		{"(and 1 2 3)", "3"},
		{"(and \"hello\" \"world\")", "\"world\""},
		// Test and with falsy values - returns first falsy value
		{"(and true nil)", "nil"},
		{"(and 1 nil 3)", "nil"},
		{"(and nil false)", "nil"},

		// Test or with no arguments
		{"(or)", "nil"},
		// Test or with single argument
		{"(or true)", "true"},
		{"(or nil)", "nil"},
		{"(or 42)", "42"},
		// Test or with multiple arguments - returns first truthy value
		{"(or true false)", "true"},
		{"(or nil 42)", "42"},
		{"(or nil false \"hello\")", "\"hello\""},
		// Test or with all falsy values - returns last value
		{"(or nil false)", "nil"},
		{"(or false nil)", "nil"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestEvalLoopRecur(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// Basic loop with no recur - just returns final expression
		{"(loop [x 5] x)", "5"},
		{"(loop [x 1 y 2] (+ x y))", "3"},
		
		// Simple factorial using loop/recur
		{"(loop [n 5 acc 1] (if (= n 0) acc (recur (- n 1) (* acc n))))", "120"},
		
		// Countdown to zero
		{"(loop [i 3] (if (= i 0) \"done\" (recur (- i 1))))", "\"done\""},
		
		// Sum from 1 to n
		{"(loop [n 5 sum 0] (if (= n 0) sum (recur (- n 1) (+ sum n))))", "15"},
		
		// Fibonacci using loop/recur
		{"(loop [n 6 a 0 b 1] (if (= n 0) a (recur (- n 1) b (+ a b))))", "8"},
		
		// Loop with multiple body expressions
		{"(loop [x 10] (def temp x) (if (= temp 0) \"zero\" (recur (- temp 1))))", "\"zero\""},
		
		// Empty loop body
		{"(loop [] 42)", "42"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestEvalLoopRecurErrors(t *testing.T) {
	env := core.NewCoreEnvironment()

	errorTests := []struct {
		input       string
		expectedErr string
	}{
		// Loop with wrong binding format
		{"(loop [x] x)", "loop bindings must be even number of forms"},
		{"(loop [x 1 y] x)", "loop bindings must be even number of forms"},
		{"(loop [1 2] 3)", "loop binding names must be symbols"},
		
		// Loop with wrong number of arguments
		{"(loop)", "loop expects at least 2 arguments"},
		{"(loop [x 1])", "loop expects at least 2 arguments"},
		
		// Loop with wrong binding types
		{"(loop 5 x)", "loop expects vector or list for bindings"},
		{"(loop \"bindings\" x)", "loop expects vector or list for bindings"},
		
		// Recur with wrong arity
		{"(loop [x 1] (recur 1 2))", "recur expects 1 arguments, got 2"},
		{"(loop [x 1 y 2] (recur 1))", "recur expects 2 arguments, got 1"},
		
		// Recur outside of loop context (should still work but will be caught by function)
		{"(recur 1)", "#<recur>"},  // This should return the RecurValue since no enclosing loop
	}

	for _, test := range errorTests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if test.expectedErr == "#<recur>" {
			// Special case - expecting RecurValue return
			if err != nil {
				t.Errorf("Expected RecurValue for '%s', got error: %v", test.input, err)
				continue
			}
			if result.String() != test.expectedErr {
				t.Errorf("Expected '%s' for '%s', got '%s'", test.expectedErr, test.input, result.String())
			}
		} else {
			// Expecting an error
			if err == nil {
				t.Errorf("Expected error for '%s', got result: %s", test.input, result.String())
				continue
			}
			if !strings.Contains(err.Error(), test.expectedErr) {
				t.Errorf("Expected error containing '%s' for '%s', got: %v", test.expectedErr, test.input, err)
			}
		}
	}
}

func TestEvalRecurInFunction(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Test recur in user-defined functions
	tests := []struct {
		setup    string
		input    string
		expected string
	}{
		// Factorial function using recur
		{
			"(defn factorial [n] (if (= n 0) 1 (* n (factorial (- n 1)))))",
			"(factorial 5)",
			"120",
		},
		// Tail-recursive factorial using recur
		{
			"(defn factorial-recur [n acc] (if (= n 0) acc (recur (- n 1) (* acc n))))",
			"(factorial-recur 5 1)",
			"120",
		},
		// Countdown function with recur
		{
			"(defn countdown [n] (if (= n 0) \"done\" (recur (- n 1))))",
			"(countdown 3)",
			"\"done\"",
		},
		// Fibonacci with recur
		{
			"(defn fib-helper [n a b] (if (= n 0) a (recur (- n 1) b (+ a b))))",
			"(fib-helper 6 0 1)",
			"8",
		},
	}

	for _, test := range tests {
		// Setup function definition
		if test.setup != "" {
			setupExpr, err := core.ReadString(test.setup)
			if err != nil {
				t.Errorf("Parse error for setup '%s': %v", test.setup, err)
				continue
			}
			_, err = core.Eval(setupExpr, env)
			if err != nil {
				t.Errorf("Setup error for '%s': %v", test.setup, err)
				continue
			}
		}

		// Test the function call
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestMapOperations(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// keys function
		{"(keys {})", "()"},
		{"(keys {:a 1})", "(:a)"},
		{"(keys {:a 1 :b 2 :c 3})", "(:a :b :c)"},
		
		// vals function
		{"(vals {})", "()"},
		{"(vals {:a 1})", "(1)"},
		{"(vals {:a 1 :b 2 :c 3})", "(1 2 3)"},
		
		// zipmap function
		{"(zipmap [] [])", "{}"},
		{"(zipmap [:a] [1])", "{:a 1}"},
		{"(zipmap [:a :b :c] [1 2 3])", "{:a 1 :b 2 :c 3}"},
		{"(zipmap [:a :b :c] [1 2])", "{:a 1 :b 2}"},
		{"(zipmap [:a :b] [1 2 3])", "{:a 1 :b 2}"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestMetaProgrammingConstructors(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// symbol function
		{"(symbol \"test\")", "test"},
		{"(symbol \"foo-bar\")", "foo-bar"},
		{"(symbol 'existing)", "existing"},
		
		// keyword function
		{"(keyword \"test\")", ":test"},
		{"(keyword \"foo-bar\")", ":foo-bar"},
		{"(keyword \":already\")", ":already"},
		{"(keyword 'sym)", ":sym"},
		{"(keyword :existing)", ":existing"},
		
		// name function
		{"(name 'test)", "\"test\""},
		{"(name :keyword)", "\"keyword\""},
		{"(name \":prefixed\")", "\"prefixed\""},
		{"(name \"string\")", "\"string\""},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestSetSubsetSuperset(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input    string
		expected string
	}{
		// subset? function
		{"(subset? #{} #{})", "true"},
		{"(subset? #{} #{1 2})", "true"},
		{"(subset? #{1} #{1 2})", "true"},
		{"(subset? #{1 2} #{1 2 3})", "true"},
		{"(subset? #{1 2} #{1 2})", "true"},
		{"(subset? #{1 3} #{1 2})", "nil"},
		{"(subset? #{1 2 3} #{1 2})", "nil"},
		
		// superset? function
		{"(superset? #{} #{})", "true"},
		{"(superset? #{1 2} #{})", "true"},
		{"(superset? #{1 2} #{1})", "true"},
		{"(superset? #{1 2 3} #{1 2})", "true"},
		{"(superset? #{1 2} #{1 2})", "true"},
		{"(superset? #{1 2} #{1 3})", "nil"},
		{"(superset? #{1 2} #{1 2 3})", "nil"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestNewFunctionErrors(t *testing.T) {
	env := core.NewCoreEnvironment()

	tests := []struct {
		input       string
		shouldError bool
		errorMatch  string
	}{
		// keys errors
		{"(keys)", true, "keys expects 1 argument"},
		{"(keys 1 2)", true, "keys expects 1 argument"},
		{"(keys \"not-a-map\")", true, "keys expects hash-map"},
		
		// vals errors
		{"(vals)", true, "vals expects 1 argument"},
		{"(vals 1 2)", true, "vals expects 1 argument"},
		{"(vals [1 2 3])", true, "vals expects hash-map"},
		
		// zipmap errors
		{"(zipmap)", true, "zipmap expects 2 arguments"},
		{"(zipmap [])", true, "zipmap expects 2 arguments"},
		{"(zipmap [] [] [])", true, "zipmap expects 2 arguments"},
		{"(zipmap \"not-collection\" [])", true, "expected collection"},
		
		// symbol errors
		{"(symbol)", true, "symbol expects 1 argument"},
		{"(symbol 1)", true, "symbol expects string or symbol"},
		
		// keyword errors
		{"(keyword)", true, "keyword expects 1 argument"},
		{"(keyword 123)", true, "keyword expects string, symbol, or keyword"},
		
		// name errors
		{"(name)", true, "name expects 1 argument"},
		{"(name 123)", true, "name expects symbol, keyword, or string"},
		
		// subset?/superset? errors
		{"(subset?)", true, "subset? expects 2 arguments"},
		{"(subset? #{})", true, "subset? expects 2 arguments"},
		{"(subset? #{} \"not-set\")", true, "subset? expects two sets"},
		{"(superset? #{})", true, "superset? expects 2 arguments"},
		{"(superset? #{} [1 2])", true, "superset? expects two sets"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			if !test.shouldError {
				t.Errorf("Unexpected parse error for '%s': %v", test.input, err)
			}
			continue
		}

		result, err := core.Eval(expr, env)
		if test.shouldError {
			if err == nil {
				t.Errorf("Expected error for '%s', but got result: %v", test.input, result)
			} else if !strings.Contains(err.Error(), test.errorMatch) {
				t.Errorf("Expected error containing '%s' for '%s', got: %v", test.errorMatch, test.input, err)
			}
		} else if err != nil {
			t.Errorf("Unexpected error for '%s': %v", test.input, err)
		}
	}
}
