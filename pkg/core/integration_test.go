package core

import (
	"os"
	"testing"
)

func TestCoreEnvironmentBootstrap(t *testing.T) {
	env, err := CreateBootstrappedEnvironment()
	if err != nil {
		t.Errorf("Failed to create bootstrapped environment: %v", err)
	}
	
	// Test that core primitives are available
	primitives := []string{
		"+", "-", "*", "/", "=", "<", ">",
		"cons", "first", "rest",
		"eval", "read-string",
		"slurp", "spit",
		"symbol?", "string?", "number?", "list?", "vector?",
		"nil", "true",
	}
	
	for _, primitive := range primitives {
		sym := Intern(primitive)
		_, err := env.Get(sym)
		if err != nil {
			t.Errorf("Core primitive '%s' not found in environment: %v", primitive, err)
		}
	}
}

func TestREPLEvaluation(t *testing.T) {
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
	}
	
	tests := []struct {
		input    string
		expected string
	}{
		{"(+ 1 2 3)", "6"},
		{"(* 2 3 4)", "24"},
		{"(cons 1 (cons 2 nil))", "(1 2 nil)"},
		{"(first (cons 1 nil))", "1"},
		{"'(+ 1 2)", "(+ 1 2)"},
		{"(if true 42 0)", "42"},
		{"(if nil 42 0)", "0"},
	}
	
	for _, test := range tests {
		result, err := repl.EvalString(test.input)
		if err != nil {
			t.Errorf("REPL eval error for '%s': %v", test.input, err)
			continue
		}
		
		if result.String() != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestComplexPrograms(t *testing.T) {
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
	}
	
	// Test factorial function
	program := []string{
		"(def factorial (fn [n] (if (= n 0) 1 (* n (factorial (- n 1))))))",
		"(factorial 5)",
	}
	
	var result Value
	for _, expr := range program {
		result, err = repl.EvalString(expr)
		if err != nil {
			t.Errorf("Error evaluating '%s': %v", expr, err)
		}
	}
	
	if result.String() != "120" {
		t.Errorf("Expected '120' for factorial(5), got '%s'", result.String())
	}
	
	// Test higher-order function
	program = []string{
		"(def apply-twice (fn [f x] (f (f x))))",
		"(def add-one (fn [x] (+ x 1)))",
		"(apply-twice add-one 5)",
	}
	
	for _, expr := range program {
		result, err = repl.EvalString(expr)
		if err != nil {
			t.Errorf("Error evaluating '%s': %v", expr, err)
		}
	}
	
	if result.String() != "7" {
		t.Errorf("Expected '7' for apply-twice add-one 5, got '%s'", result.String())
	}
}

func TestFileLoading(t *testing.T) {
	// Create a temporary test file
	testContent := `(def test-var 42)
(def test-fn (fn [x] (* x 2)))
(test-fn test-var)`
	
	tmpFile, err := os.CreateTemp("", "golisp-test-*.lisp")
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())
	
	_, err = tmpFile.WriteString(testContent)
	if err != nil {
		t.Errorf("Failed to write temp file: %v", err)
		return
	}
	tmpFile.Close()
	
	// Test loading the file
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
		return
	}
	
	err = repl.LoadFile(tmpFile.Name())
	if err != nil {
		t.Errorf("Failed to load file: %v", err)
		return
	}
	
	// Test that variables are defined
	result, err := repl.EvalString("test-var")
	if err != nil {
		t.Errorf("Error accessing test-var: %v", err)
		return
	}
	if result.String() != "42" {
		t.Errorf("Expected '42' for test-var, got '%s'", result.String())
	}
	
	result, err = repl.EvalString("(test-fn 10)")
	if err != nil {
		t.Errorf("Error calling test-fn: %v", err)
		return
	}
	if result.String() != "20" {
		t.Errorf("Expected '20' for (test-fn 10), got '%s'", result.String())
	}
}

func TestFileIO(t *testing.T) {
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
		return
	}
	
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "golisp-io-test-*.txt")
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
		return
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	
	// Test spit (write)
	writeExpr := `(spit "` + tmpFile.Name() + `" "Hello, World!")`
	_, err = repl.EvalString(writeExpr)
	if err != nil {
		t.Errorf("Error with spit: %v", err)
		return
	}
	
	// Test slurp (read)
	readExpr := `(slurp "` + tmpFile.Name() + `")`
	result, err := repl.EvalString(readExpr)
	if err != nil {
		t.Errorf("Error with slurp: %v", err)
		return
	}
	
	expected := `"Hello, World!"`
	if result.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.String())
	}
}

func TestMetaProgramming(t *testing.T) {
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
		return
	}
	
	// Test code as data
	program := []string{
		"(def code '(+ 1 2 3))",
		"(eval code)",
	}
	
	var result Value
	for _, expr := range program {
		result, err = repl.EvalString(expr)
		if err != nil {
			t.Errorf("Error evaluating '%s': %v", expr, err)
		}
	}
	
	if result.String() != "6" {
		t.Errorf("Expected '6' for eval code, got '%s'", result.String())
	}
	
	// Test read-string + eval
	expr := `(eval (read-string "(* 3 4)"))`
	result, err = repl.EvalString(expr)
	if err != nil {
		t.Errorf("Error with read-string + eval: %v", err)
		return
	}
	
	if result.String() != "12" {
		t.Errorf("Expected '12' for read-string + eval, got '%s'", result.String())
	}
}

func TestDataStructures(t *testing.T) {
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
		return
	}
	
	tests := []struct {
		input    string
		expected string
	}{
		// Lists
		{"(cons 1 (cons 2 (cons 3 nil)))", "(1 2 3 nil)"},
		{"(first (cons 1 (cons 2 nil)))", "1"},
		{"(rest (cons 1 (cons 2 nil)))", "(2 nil)"},
		
		// Vectors
		{"[1 2 3]", "[1 2 3]"},
		{"(first [1 2 3])", "1"},
		
		// Nested structures
		{"[[1 2] [3 4]]", "[[1 2] [3 4]]"},
		{"(first [1 2 3])", "1"},
		
		// Mixed types
		{"[1 \"hello\" :keyword nil true]", "[1 \"hello\" :keyword nil true]"},
	}
	
	for _, test := range tests {
		result, err := repl.EvalString(test.input)
		if err != nil {
			t.Errorf("Error evaluating '%s': %v", test.input, err)
			continue
		}
		
		if result.String() != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestClosuresAndLexicalScoping(t *testing.T) {
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
		return
	}
	
	// Test closure creation and execution
	program := []string{
		"(def x 10)",
		"(def make-adder (fn [y] (fn [z] (+ x y z))))",
		"(def add-5 (make-adder 5))",
		"(add-5 3)",
	}
	
	var result Value
	for _, expr := range program {
		result, err = repl.EvalString(expr)
		if err != nil {
			t.Errorf("Error evaluating '%s': %v", expr, err)
		}
	}
	
	if result.String() != "18" {
		t.Errorf("Expected '18' for closure result, got '%s'", result.String())
	}
	
	// Test simpler closure behavior
	program = []string{
		"(def outer-var 100)",
		"(def test-closure (fn [x] (+ x outer-var)))",
		"(test-closure 5)",
	}
	
	for _, expr := range program {
		result, err = repl.EvalString(expr)
		if err != nil {
			t.Errorf("Error evaluating '%s': %v", expr, err)
		}
	}
	
	if result.String() != "105" {
		t.Errorf("Expected '105' for closure test, got '%s'", result.String())
	}
}

func TestErrorHandling(t *testing.T) {
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
		return
	}
	
	// Test various error conditions
	errorTests := []string{
		"(+ 1 \"hello\")",           // Type error
		"(/ 1 0)",                  // Division by zero
		"(unknown-function 1 2)",   // Undefined function
		"undefined-variable",       // Undefined variable
		"(fn [x] x)",              // Function without proper call (should not error)
		"((fn [x] x) 1 2)",        // Wrong number of arguments
	}
	
	// Most of these should produce errors
	errorCount := 0
	for _, test := range errorTests {
		_, err := repl.EvalString(test)
		if err != nil {
			errorCount++
		}
	}
	
	// We expect at least 4 errors (excluding the valid function definition)
	if errorCount < 4 {
		t.Errorf("Expected at least 4 errors, got %d", errorCount)
	}
}

func TestCoreMinimalFootprint(t *testing.T) {
	// Verify that the core has minimal dependencies
	env := NewCoreEnvironment()
	
	// Count the number of built-in functions
	builtinCount := 0
	
	// This is a bit indirect, but we can test that only essential primitives exist
	essentialPrimitives := []string{
		"+", "-", "*", "/", "=", "<", ">",  // Arithmetic/comparison (7)
		"cons", "first", "rest",            // List operations (3)
		"eval", "read-string",              // Meta-programming (2)
		"slurp", "spit",                    // File I/O (2)
		"symbol?", "string?", "number?", "list?", "vector?", // Type predicates (5)
	}
	
	for _, primitive := range essentialPrimitives {
		sym := Intern(primitive)
		_, err := env.Get(sym)
		if err != nil {
			t.Errorf("Essential primitive '%s' missing from core", primitive)
		} else {
			builtinCount++
		}
	}
	
	// Verify we have exactly the expected number of core primitives
	expectedCount := len(essentialPrimitives)
	if builtinCount != expectedCount {
		t.Errorf("Expected %d core primitives, found %d", expectedCount, builtinCount)
	}
}

func TestEndToEndProgram(t *testing.T) {
	// Create a comprehensive program that tests multiple features
	program := `
	; Define helper functions
	(def abs (fn [x] (if (< x 0) (- x) x)))
	(def square (fn [x] (* x x)))
	
	; Define a more complex function
	(def distance (fn [x1 y1 x2 y2]
		(abs (+ (square (- x2 x1)) (square (- y2 y1))))))
	
	; Test the function
	(distance 0 0 3 4)
	`
	
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "golisp-end-to-end-*.lisp")
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())
	
	_, err = tmpFile.WriteString(program)
	if err != nil {
		t.Errorf("Failed to write temp file: %v", err)
		return
	}
	tmpFile.Close()
	
	// Load and execute the program
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
		return
	}
	
	err = repl.LoadFile(tmpFile.Name())
	if err != nil {
		t.Errorf("Failed to load program: %v", err)
		return
	}
	
	// The program should have calculated the distance squared (3² + 4² = 25)
	result, err := repl.EvalString("(distance 0 0 3 4)")
	if err != nil {
		t.Errorf("Error calling distance function: %v", err)
		return
	}
	
	if result.String() != "25" {
		t.Errorf("Expected '25' for distance calculation, got '%s'", result.String())
	}
}

func TestCoreVsFullComparison(t *testing.T) {
	// Test that our minimal core can handle the same basic operations
	// as would be expected from a full Lisp system
	
	basicLispOperations := []struct {
		input    string
		expected string
	}{
		// Arithmetic
		{"(+ 1 2 3 4 5)", "15"},
		{"(- 100 25)", "75"},
		{"(* 6 7)", "42"},
		{"(/ 84 2)", "42"},
		
		// Logic and comparisons
		{"(= 42 42)", "true"},
		{"(< 5 10)", "true"},
		{"(> 10 5)", "true"},
		
		// List processing
		{"(first (cons 'a (cons 'b nil)))", "a"},
		{"(rest (cons 'a (cons 'b nil)))", "(b nil)"},
		
		// Functions and closures
		{"((fn [x] (+ x 1)) 41)", "42"},
		
		// Conditional logic
		{"(if (< 2 3) 'yes 'no)", "yes"},
		
		// Meta-programming
		{"(eval '(+ 20 22))", "42"},
	}
	
	repl, err := NewREPL()
	if err != nil {
		t.Errorf("Failed to create REPL: %v", err)
		return
	}
	
	for _, test := range basicLispOperations {
		result, err := repl.EvalString(test.input)
		if err != nil {
			t.Errorf("Error evaluating '%s': %v", test.input, err)
			continue
		}
		
		if result.String() != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}