package core_test

import (
	"fmt"
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
