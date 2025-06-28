package core

import (
	"fmt"
	"testing"
)

func TestEvalBasicValues(t *testing.T) {
	env := NewCoreEnvironment()
	
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
		expr, err := ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}
		
		result, err := Eval(expr, env)
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
	env := NewCoreEnvironment()
	
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
		expr, err := ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}
		
		result, err := Eval(expr, env)
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
	env := NewCoreEnvironment()
	
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
		expr, err := ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}
		
		result, err := Eval(expr, env)
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
	env := NewCoreEnvironment()
	
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
		expr, err := ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}
		
		result, err := Eval(expr, env)
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
	env := NewCoreEnvironment()
	
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
		expr, err := ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}
		
		result, err := Eval(expr, env)
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
	env := NewCoreEnvironment()
	
	// Test quote
	expr, _ := ReadString("(quote x)")
	result, err := Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for quote: %v", err)
	}
	if result.String() != "x" {
		t.Errorf("Expected 'x' for quote, got '%s'", result.String())
	}
	
	// Test shorthand quote
	expr, _ = ReadString("'x")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for shorthand quote: %v", err)
	}
	if result.String() != "x" {
		t.Errorf("Expected 'x' for shorthand quote, got '%s'", result.String())
	}
	
	// Test if - true case
	expr, _ = ReadString("(if true 1 2)")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for if true: %v", err)
	}
	if result.String() != "1" {
		t.Errorf("Expected '1' for if true, got '%s'", result.String())
	}
	
	// Test if - false case
	expr, _ = ReadString("(if nil 1 2)")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for if false: %v", err)
	}
	if result.String() != "2" {
		t.Errorf("Expected '2' for if false, got '%s'", result.String())
	}
	
	// Test if - no else
	expr, _ = ReadString("(if nil 1)")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for if no else: %v", err)
	}
	if result.String() != "nil" {
		t.Errorf("Expected 'nil' for if no else, got '%s'", result.String())
	}
	
	// Test do
	expr, _ = ReadString("(do 1 2 3)")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for do: %v", err)
	}
	if result.String() != "3" {
		t.Errorf("Expected '3' for do, got '%s'", result.String())
	}
}

func TestEvalDefAndSymbolLookup(t *testing.T) {
	env := NewCoreEnvironment()
	
	// Test def
	expr, _ := ReadString("(def x 42)")
	result, err := Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for def: %v", err)
	}
	if result.String() != "x" {
		t.Errorf("Expected 'x' for def, got '%s'", result.String())
	}
	
	// Test symbol lookup
	expr, _ = ReadString("x")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for symbol lookup: %v", err)
	}
	if result.String() != "42" {
		t.Errorf("Expected '42' for symbol lookup, got '%s'", result.String())
	}
	
	// Test undefined symbol
	expr, _ = ReadString("undefined")
	_, err = Eval(expr, env)
	if err == nil {
		t.Error("Expected error for undefined symbol")
	}
}

func TestEvalFunctionDefinitionAndCall(t *testing.T) {
	env := NewCoreEnvironment()
	
	// Define a function
	expr, _ := ReadString("(def add (fn [a b] (+ a b)))")
	_, err := Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for function definition: %v", err)
	}
	
	// Call the function
	expr, _ = ReadString("(add 3 4)")
	result, err := Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for function call: %v", err)
	}
	if result.String() != "7" {
		t.Errorf("Expected '7' for function call, got '%s'", result.String())
	}
	
	// Test function with vector parameters
	expr, _ = ReadString("(def mult (fn [x y] (* x y)))")
	_, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for function definition with vector: %v", err)
	}
	
	expr, _ = ReadString("(mult 5 6)")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for function call with vector params: %v", err)
	}
	if result.String() != "30" {
		t.Errorf("Expected '30' for function call, got '%s'", result.String())
	}
}

func TestEvalRecursiveFunction(t *testing.T) {
	env := NewCoreEnvironment()
	
	// Define factorial function
	expr, _ := ReadString("(def factorial (fn [n] (if (= n 0) 1 (* n (factorial (- n 1))))))")
	_, err := Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for recursive function definition: %v", err)
	}
	
	// Test factorial(5)
	expr, _ = ReadString("(factorial 5)")
	result, err := Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for recursive function call: %v", err)
	}
	if result.String() != "120" {
		t.Errorf("Expected '120' for factorial(5), got '%s'", result.String())
	}
	
	// Test factorial(0)
	expr, _ = ReadString("(factorial 0)")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for factorial(0): %v", err)
	}
	if result.String() != "1" {
		t.Errorf("Expected '1' for factorial(0), got '%s'", result.String())
	}
}

func TestEvalLexicalScoping(t *testing.T) {
	env := NewCoreEnvironment()
	
	// Define outer variable
	expr, _ := ReadString("(def x 10)")
	_, err := Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for outer variable: %v", err)
	}
	
	// Define function that captures outer variable
	expr, _ = ReadString("(def make-adder (fn [y] (fn [z] (+ x y z))))")
	_, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for closure definition: %v", err)
	}
	
	// Create adder function
	expr, _ = ReadString("(def add-5 (make-adder 5))")
	_, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for creating closure: %v", err)
	}
	
	// Call the closure
	expr, _ = ReadString("(add-5 3)")
	result, err := Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for calling closure: %v", err)
	}
	if result.String() != "18" {
		t.Errorf("Expected '18' for closure call, got '%s'", result.String())
	}
}

func TestEvalMetaProgramming(t *testing.T) {
	env := NewCoreEnvironment()
	
	// Test eval
	expr, _ := ReadString("(eval '(+ 1 2))")
	result, err := Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for eval: %v", err)
	}
	if result.String() != "3" {
		t.Errorf("Expected '3' for eval, got '%s'", result.String())
	}
	
	// Test read-string
	expr, _ = ReadString("(read-string \"(+ 1 2)\")")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for read-string: %v", err)
	}
	if result.String() != "(+ 1 2)" {
		t.Errorf("Expected '(+ 1 2)' for read-string, got '%s'", result.String())
	}
	
	// Test eval + read-string
	expr, _ = ReadString("(eval (read-string \"(+ 1 2)\"))")
	result, err = Eval(expr, env)
	if err != nil {
		t.Errorf("Eval error for eval + read-string: %v", err)
	}
	if result.String() != "3" {
		t.Errorf("Expected '3' for eval + read-string, got '%s'", result.String())
	}
}

func TestEvalErrors(t *testing.T) {
	env := NewCoreEnvironment()
	
	tests := []string{
		"(+ 1 \"hello\")",        // Type error
		"(/ 1 0)",                // Division by zero
		"(unknown-function)",     // Unknown function
		"(def)",                  // Wrong number of arguments
		"(fn)",                   // Wrong number of arguments
		"(if)",                   // Wrong number of arguments
		"(quote)",                // Wrong number of arguments
		"(= 1)",                  // Wrong number of arguments for =
		"(< 1)",                  // Wrong number of arguments for <
		"(> 1)",                  // Wrong number of arguments for >
		"(cons 1)",               // Wrong number of arguments for cons
		"(first)",                // Wrong number of arguments for first
		"(rest)",                 // Wrong number of arguments for rest
	}
	
	for _, test := range tests {
		expr, err := ReadString(test)
		if err != nil {
			continue // Skip parse errors for this test
		}
		
		_, err = Eval(expr, env)
		if err == nil {
			t.Errorf("Expected error for '%s', but got none", test)
		}
	}
}

func TestUserFunctionInterface(t *testing.T) {
	env := NewCoreEnvironment()
	
	// Create a user function
	params := NewList(Intern("x"))
	body, _ := ReadString("(+ x 1)")
	userFn := &UserFunction{
		Params: params,
		Body:   body,
		Env:    env,
	}
	
	// Test Call method
	args := []Value{NewNumber(int64(5))}
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
	wrongArgs := []Value{NewNumber(int64(1)), NewNumber(int64(2))}
	_, err = userFn.Call(wrongArgs, env)
	if err == nil {
		t.Error("Expected error for wrong number of arguments")
	}
}

func TestBuiltinFunctionInterface(t *testing.T) {
	// Test a builtin function
	addFn := &BuiltinFunction{
		Name: "test-add",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("test-add expects 2 arguments")
			}
			n1, _ := args[0].(Number)
			n2, _ := args[1].(Number)
			return NewNumber(n1.ToInt() + n2.ToInt()), nil
		},
	}
	
	// Test Call method
	args := []Value{NewNumber(int64(3)), NewNumber(int64(4))}
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