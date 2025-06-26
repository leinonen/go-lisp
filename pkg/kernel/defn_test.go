package kernel

import "testing"

// TestDefnMacro tests the defn macro functionality
func TestDefnMacro(t *testing.T) {
	t.Run("BasicDefn", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a function using defn
		// (defn square [x] (* x x))
		defnExpr := NewList(
			Intern("defn"),
			Intern("square"),
			NewVector(Intern("x")),
			NewList(Intern("*"), Intern("x"), Intern("x")),
		)

		result, err := Eval(defnExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining function with defn: %v", err)
		}

		// Should return DefinedValue
		if _, ok := result.(DefinedValue); !ok {
			t.Errorf("Expected DefinedValue from defn, got %T", result)
		}

		// Verify the function is defined and callable
		square, err := repl.Env.Get(Intern("square"))
		if err != nil {
			t.Fatalf("Error getting defined function: %v", err)
		}

		// Should be a UserFunction
		if _, ok := square.(*UserFunction); !ok {
			t.Errorf("Expected UserFunction, got %T", square)
		}

		// Test calling the function
		callExpr := NewList(Intern("square"), Number(5))
		result, err = Eval(callExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling square function: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 25.0 {
			t.Errorf("Expected square(5) = 25, got %v", result)
		}
	})

	t.Run("DefnWithMultipleParams", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a function with multiple parameters
		// (defn add [x y] (+ x y))
		defnExpr := NewList(
			Intern("defn"),
			Intern("add"),
			NewVector(Intern("x"), Intern("y")),
			NewList(Intern("+"), Intern("x"), Intern("y")),
		)

		_, err := Eval(defnExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining multi-param function with defn: %v", err)
		}

		// Test calling the function
		callExpr := NewList(Intern("add"), Number(3), Number(7))
		result, err := Eval(callExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling add function: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 10.0 {
			t.Errorf("Expected add(3, 7) = 10, got %v", result)
		}
	})

	t.Run("DefnWithComplexBody", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a function with conditional logic
		// (defn abs [x] (if (< x 0) (- x) x))
		defnExpr := NewList(
			Intern("defn"),
			Intern("abs"),
			NewVector(Intern("x")),
			NewList(
				Intern("if"),
				NewList(Intern("<"), Intern("x"), Number(0)),
				NewList(Intern("-"), Intern("x")),
				Intern("x"),
			),
		)

		_, err := Eval(defnExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining abs function with defn: %v", err)
		}

		// Test with negative number
		callExpr1 := NewList(Intern("abs"), Number(-5))
		result, err := Eval(callExpr1, repl.Env)
		if err != nil {
			t.Fatalf("Error calling abs(-5): %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 5.0 {
			t.Errorf("Expected abs(-5) = 5, got %v", result)
		}

		// Test with positive number
		callExpr2 := NewList(Intern("abs"), Number(3))
		result, err = Eval(callExpr2, repl.Env)
		if err != nil {
			t.Fatalf("Error calling abs(3): %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 3.0 {
			t.Errorf("Expected abs(3) = 3, got %v", result)
		}
	})

	t.Run("DefnErrorHandling", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test defn with wrong number of arguments
		errorCases := []*List{
			NewList(Intern("defn")),                             // no args
			NewList(Intern("defn"), Intern("foo")),              // only name
			NewList(Intern("defn"), Intern("foo"), NewVector()), // no body
		}

		for _, tc := range errorCases {
			_, err := Eval(tc, repl.Env)
			if err == nil {
				t.Errorf("Expected error for defn with wrong args: %s", tc.String())
			}
		}

		// Test defn with non-symbol name
		nonSymbolName := NewList(
			Intern("defn"),
			Number(42), // not a symbol
			NewVector(),
			Number(1),
		)

		_, err := Eval(nonSymbolName, repl.Env)
		if err == nil {
			t.Error("Expected error for defn with non-symbol name")
		}
	})

	t.Run("DefnRecursiveFunction", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a recursive factorial function
		// (defn fact [n] (if (= n 0) 1 (* n (fact (- n 1)))))
		defnExpr := NewList(
			Intern("defn"),
			Intern("fact"),
			NewVector(Intern("n")),
			NewList(
				Intern("if"),
				NewList(Intern("="), Intern("n"), Number(0)),
				Number(1),
				NewList(
					Intern("*"),
					Intern("n"),
					NewList(
						Intern("fact"),
						NewList(Intern("-"), Intern("n"), Number(1)),
					),
				),
			),
		)

		_, err := Eval(defnExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining recursive function with defn: %v", err)
		}

		// Test factorial calculation
		callExpr := NewList(Intern("fact"), Number(5))
		result, err := Eval(callExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling fact(5): %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 120.0 {
			t.Errorf("Expected fact(5) = 120, got %v", result)
		}
	})
}
