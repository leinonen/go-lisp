package minimal

import (
	"testing"
)

// TestMinimalKernelIntegration provides comprehensive integration testing
func TestMinimalKernelIntegration(t *testing.T) {
	t.Run("CompleteFactorialExample", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define factorial using recursion
		factorialDef := `(def factorial
			(fn [n]
				(if (= n 0)
					1
					(* n (factorial (- n 1))))))`

		_, err := repl.Eval(factorialDef)
		if err != nil {
			t.Fatalf("Error defining factorial: %v", err)
		}

		// Test various factorial values
		testCases := []struct {
			input    int
			expected int
		}{
			{0, 1},
			{1, 1},
			{2, 2},
			{3, 6},
			{4, 24},
			{5, 120},
		}

		for _, tc := range testCases {
			// Convert to proper expression
			factExpr := `(factorial ` + []string{"0", "1", "2", "3", "4", "5"}[tc.input] + `)`

			result, err := repl.Eval(factExpr)
			if err != nil {
				t.Fatalf("Error evaluating factorial(%d): %v", tc.input, err)
			}

			if num, ok := result.(Number); !ok || int(float64(num)) != tc.expected {
				t.Errorf("factorial(%d): expected %d, got %v", tc.input, tc.expected, result)
			}
		}
	})

	t.Run("FibonacciSequence", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define fibonacci using recursion
		fibDef := `(def fib
			(fn [n]
				(if (< n 2)
					n
					(+ (fib (- n 1)) (fib (- n 2))))))`

		_, err := repl.Eval(fibDef)
		if err != nil {
			t.Fatalf("Error defining fibonacci: %v", err)
		}

		// Test fibonacci sequence
		fibTests := []struct {
			n        int
			expected int
		}{
			{0, 0},
			{1, 1},
			{2, 1},
			{3, 2},
			{4, 3},
			{5, 5},
			{6, 8},
		}

		for _, tc := range fibTests {
			fibExpr := `(fib ` + []string{"0", "1", "2", "3", "4", "5", "6"}[tc.n] + `)`

			result, err := repl.Eval(fibExpr)
			if err != nil {
				t.Fatalf("Error evaluating fib(%d): %v", tc.n, err)
			}

			if num, ok := result.(Number); !ok || int(float64(num)) != tc.expected {
				t.Errorf("fib(%d): expected %d, got %v", tc.n, tc.expected, result)
			}
		}
	})

	t.Run("HigherOrderFunctions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define map function
		mapDef := `(def map-list
			(fn [f lst]
				(if (= (first lst) nil)
					(list)
					(cons (f (first lst)) (map-list f (rest lst))))))`

		_, err := repl.Eval(mapDef)
		if err != nil {
			t.Fatalf("Error defining map-list: %v", err)
		}

		// Define a simple doubling function
		doubleDef := `(def double (fn [x] (* x 2)))`
		_, err = repl.Eval(doubleDef)
		if err != nil {
			t.Fatalf("Error defining double: %v", err)
		}

		// Test mapping double over a list
		mapExpr := `(map-list double (list 1 2 3 4))`
		result, err := repl.Eval(mapExpr)
		if err != nil {
			t.Fatalf("Error evaluating map: %v", err)
		}

		list, ok := result.(*List)
		if !ok {
			t.Fatalf("Expected List result, got %T", result)
		}

		// Should be (2 4 6 8)
		expected := []float64{2, 4, 6, 8}
		if list.Length() != len(expected) {
			t.Errorf("Expected list length %d, got %d", len(expected), list.Length())
		}

		current := list
		for i, exp := range expected {
			if current.IsEmpty() {
				t.Fatalf("List too short at element %d", i)
			}

			if num, ok := current.First().(Number); !ok || float64(num) != exp {
				t.Errorf("Element %d: expected %f, got %v", i, exp, current.First())
			}

			current = current.Rest()
		}
	})

	t.Run("MacroExpansionChain", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a chain of macros
		// First macro: when
		whenDef := `(defmacro when [condition body]
			(quasiquote (if (unquote condition) (unquote body) nil)))`

		_, err := repl.Eval(whenDef)
		if err != nil {
			t.Fatalf("Error defining when macro: %v", err)
		}

		// Second macro: unless (uses when)
		unlessDef := `(defmacro unless [condition body]
			(quasiquote (when (= (unquote condition) false) (unquote body))))`

		_, err = repl.Eval(unlessDef)
		if err != nil {
			t.Fatalf("Error defining unless macro: %v", err)
		}

		// Test the macro chain
		unlessExpr := `(unless false 42)`
		result, err := repl.Eval(unlessExpr)
		if err != nil {
			t.Fatalf("Error evaluating unless: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected 42 from unless, got %v", result)
		}

		// Test with true condition
		unlessExpr2 := `(unless true 42)`
		result, err = repl.Eval(unlessExpr2)
		if err != nil {
			t.Fatalf("Error evaluating unless with true: %v", err)
		}

		if _, ok := result.(Nil); !ok {
			t.Errorf("Expected nil from unless with true condition, got %v", result)
		}
	})

	t.Run("ComplexDataStructureManipulation", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Create a complex nested data structure
		structDef := `(def data
			(list 
				(hash-map "name" "Alice" "age" 30)
				(hash-map "name" "Bob" "age" 25)
				(vector 1 2 3)))`

		_, err := repl.Eval(structDef)
		if err != nil {
			t.Fatalf("Error defining data structure: %v", err)
		}

		// Access nested data
		accessExpr := `(hash-map-get (first data) "name")`
		result, err := repl.Eval(accessExpr)
		if err != nil {
			t.Fatalf("Error accessing nested data: %v", err)
		}

		if str, ok := result.(String); !ok || string(str) != "Alice" {
			t.Errorf("Expected 'Alice', got %v", result)
		}

		// Access vector in list
		vectorAccessExpr := `(vector-get (first (rest (rest data))) 1)`
		result, err = repl.Eval(vectorAccessExpr)
		if err != nil {
			t.Fatalf("Error accessing vector in nested structure: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 2.0 {
			t.Errorf("Expected 2 from vector access, got %v", result)
		}
	})

	t.Run("EnvironmentScoping", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test lexical scoping with nested functions
		// Note: This test demonstrates the current scoping limitations
		// A proper closure would maintain mutable state, but our implementation
		// uses immutable environments
		scopeDef := `(def make-adder
			(fn [x]
				(fn [y]
					(+ x y))))`

		_, err := repl.Eval(scopeDef)
		if err != nil {
			t.Fatalf("Error defining make-adder: %v", err)
		}

		// Create an adder
		adderDef := `(def add5 (make-adder 5))`
		_, err = repl.Eval(adderDef)
		if err != nil {
			t.Fatalf("Error creating adder: %v", err)
		}

		// Use the adder multiple times
		result1, err := repl.Eval("(add5 3)")
		if err != nil {
			t.Fatalf("Error calling adder first time: %v", err)
		}

		if num, ok := result1.(Number); !ok || float64(num) != 8.0 {
			t.Errorf("Expected 8 from first adder call, got %v", result1)
		}

		result2, err := repl.Eval("(add5 5)")
		if err != nil {
			t.Fatalf("Error calling adder second time: %v", err)
		}

		if num, ok := result2.(Number); !ok || float64(num) != 10.0 {
			t.Errorf("Expected 10 from second adder call, got %v", result2)
		}
	})

	t.Run("ErrorRecovery", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a variable successfully
		_, err := repl.Eval("(def good-var 42)")
		if err != nil {
			t.Fatalf("Error defining good variable: %v", err)
		}

		// Try something that should fail
		_, err = repl.Eval("(+ 1 \"bad\")")
		if err == nil {
			t.Error("Expected error for type mismatch")
		}

		// Verify that the REPL still works after the error
		result, err := repl.Eval("good-var")
		if err != nil {
			t.Fatalf("Error accessing variable after error: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected good-var to still be 42 after error, got %v", result)
		}

		// Should be able to define new variables
		_, err = repl.Eval("(def another-var 100)")
		if err != nil {
			t.Fatalf("Error defining new variable after error: %v", err)
		}

		result, err = repl.Eval("(+ good-var another-var)")
		if err != nil {
			t.Fatalf("Error evaluating after recovery: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 142.0 {
			t.Errorf("Expected 142 after recovery, got %v", result)
		}
	})

	t.Run("QuasiquoteComplexity", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test complex quasiquote patterns
		complexQuasiDef := `(defmacro let-values [bindings body]
			(quasiquote 
				((fn [(unquote (first bindings))] (unquote body)) 
				 (unquote (first (rest bindings))))))`

		_, err := repl.Eval(complexQuasiDef)
		if err != nil {
			t.Fatalf("Error defining complex macro: %v", err)
		}

		// Use the macro
		letExpr := `(let-values (x 10) (* x x))`
		result, err := repl.Eval(letExpr)
		if err != nil {
			t.Fatalf("Error using complex macro: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 100.0 {
			t.Errorf("Expected 100 from complex macro, got %v", result)
		}
	})

	t.Run("StressTestRecursion", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test moderate recursion depth (not too deep to avoid stack overflow)
		sumDef := `(def sum-to-n
			(fn [n]
				(if (= n 0)
					0
					(+ n (sum-to-n (- n 1))))))`

		_, err := repl.Eval(sumDef)
		if err != nil {
			t.Fatalf("Error defining sum-to-n: %v", err)
		}

		// Test sum-to-n(10) = 55
		result, err := repl.Eval("(sum-to-n 10)")
		if err != nil {
			t.Fatalf("Error evaluating sum-to-n: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 55.0 {
			t.Errorf("Expected 55 from sum-to-n(10), got %v", result)
		}
	})
}
