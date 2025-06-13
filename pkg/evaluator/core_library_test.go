package evaluator

import (
	"path/filepath"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// TestCoreLibraryFunctions tests the core library functions by loading the core.lisp file
// and testing each exported function
func TestCoreLibraryFunctions(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Helper function to evaluate expressions
	evalExpr := func(input string) (types.Value, error) {
		tok := tokenizer.NewTokenizer(input)
		tokens, err := tok.TokenizeWithError()
		if err != nil {
			return nil, err
		}

		p := parser.NewParser(tokens)
		expr, err := p.Parse()
		if err != nil {
			return nil, err
		}

		return evaluator.Eval(expr)
	}

	// Load the core library
	coreLibPath := filepath.Join("..", "..", "library", "core.lisp")
	loadExpr := `(load "` + coreLibPath + `")`
	_, err := evalExpr(loadExpr)
	if err != nil {
		t.Fatalf("Failed to load core library: %v", err)
	}

	// Import the core module
	_, err = evalExpr("(import core)")
	if err != nil {
		t.Fatalf("Failed to import core module: %v", err)
	}

	t.Run("factorial function", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"factorial of 0", "(factorial 0)", "1"},
			{"factorial of 1", "(factorial 1)", "1"},
			{"factorial of 5", "(factorial 5)", "120"},
			{"factorial of 10", "(factorial 10)", "3628800"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := evalExpr(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			})
		}

		// Test error case for negative input
		t.Run("factorial negative error", func(t *testing.T) {
			_, err := evalExpr("(factorial -1)")
			if err == nil {
				t.Error("expected error for negative factorial input")
			}
		})
	})

	t.Run("fibonacci function", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"fibonacci of 0", "(fibonacci 0)", "0"},
			{"fibonacci of 1", "(fibonacci 1)", "1"},
			{"fibonacci of 5", "(fibonacci 5)", "5"},
			{"fibonacci of 10", "(fibonacci 10)", "55"},
			{"fibonacci of 15", "(fibonacci 15)", "610"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := evalExpr(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			})
		}

		// Test error case for negative input
		t.Run("fibonacci negative error", func(t *testing.T) {
			_, err := evalExpr("(fibonacci -1)")
			if err == nil {
				t.Error("expected error for negative fibonacci input")
			}
		})
	})

	t.Run("gcd function", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"gcd of 48 and 18", "(gcd 48 18)", "6"},
			{"gcd of 56 and 42", "(gcd 56 42)", "14"},
			{"gcd of 17 and 13", "(gcd 17 13)", "1"},
			{"gcd with zero", "(gcd 15 0)", "15"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := evalExpr(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			})
		}
	})

	t.Run("lcm function", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"lcm of 12 and 8", "(lcm 12 8)", "24"},
			{"lcm of 15 and 20", "(lcm 15 20)", "60"},
			{"lcm of 7 and 13", "(lcm 7 13)", "91"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := evalExpr(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			})
		}
	})

	t.Run("abs function", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"abs of positive number", "(abs 42)", "42"},
			{"abs of negative number", "(abs -17)", "17"},
			{"abs of zero", "(abs 0)", "0"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := evalExpr(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			})
		}
	})

	t.Run("min and max functions", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"min of two numbers", "(min 10 5)", "5"},
			{"min with equal numbers", "(min 7 7)", "7"},
			{"max of two numbers", "(max 10 5)", "10"},
			{"max with equal numbers", "(max 7 7)", "7"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := evalExpr(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			})
		}
	})

	t.Run("list utility functions", func(t *testing.T) {
		// Test length-sq
		t.Run("length-sq", func(t *testing.T) {
			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"length-sq of 3-element list", "(length-sq (list 1 2 3))", "9"},
				{"length-sq of empty list", "(length-sq (list))", "0"},
				{"length-sq of 5-element list", "(length-sq (list 1 2 3 4 5))", "25"},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					result, err := evalExpr(tt.input)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if result.String() != tt.expected {
						t.Errorf("expected %s, got %s", tt.expected, result.String())
					}
				})
			}
		})

		// Test all predicate
		t.Run("all predicate", func(t *testing.T) {
			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"all positive", "(all (fn [x] (> x 0)) (list 1 2 3))", "#t"},
				{"not all positive", "(all (fn [x] (> x 0)) (list 1 -2 3))", "#f"},
				{"all on empty list", "(all (fn [x] (> x 0)) (list))", "#t"},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					result, err := evalExpr(tt.input)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if result.String() != tt.expected {
						t.Errorf("expected %s, got %s", tt.expected, result.String())
					}
				})
			}
		})

		// Test any predicate
		t.Run("any predicate", func(t *testing.T) {
			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"any positive", "(any (fn [x] (> x 0)) (list -1 2 -3))", "#t"},
				{"no positive", "(any (fn [x] (> x 0)) (list -1 -2 -3))", "#f"},
				{"any on empty list", "(any (fn [x] (> x 0)) (list))", "#f"},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					result, err := evalExpr(tt.input)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if result.String() != tt.expected {
						t.Errorf("expected %s, got %s", tt.expected, result.String())
					}
				})
			}
		})

		// Test take function
		t.Run("take function", func(t *testing.T) {
			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"take 3 from list", "(take 3 (list 1 2 3 4 5))", "(1 2 3)"},
				{"take 0 elements", "(take 0 (list 1 2 3))", "()"},
				{"take more than length", "(take 10 (list 1 2 3))", "(1 2 3)"},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					result, err := evalExpr(tt.input)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if result.String() != tt.expected {
						t.Errorf("expected %s, got %s", tt.expected, result.String())
					}
				})
			}
		})

		// Test drop function
		t.Run("drop function", func(t *testing.T) {
			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"drop 2 from list", "(drop 2 (list 1 2 3 4 5))", "(3 4 5)"},
				{"drop 0 elements", "(drop 0 (list 1 2 3))", "(1 2 3)"},
				{"drop more than length", "(drop 10 (list 1 2 3))", "()"},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					result, err := evalExpr(tt.input)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if result.String() != tt.expected {
						t.Errorf("expected %s, got %s", tt.expected, result.String())
					}
				})
			}
		})
	})

	t.Run("higher-order functions", func(t *testing.T) {
		// Test compose function
		t.Run("compose function", func(t *testing.T) {
			// Define helper functions for testing
			_, err := evalExpr("(def square (fn [x] (* x x)))")
			if err != nil {
				t.Fatalf("failed to define square: %v", err)
			}
			_, err = evalExpr("(def increment (fn [x] (+ x 1)))")
			if err != nil {
				t.Fatalf("failed to define increment: %v", err)
			}

			// Test composition
			result, err := evalExpr("((compose square increment) 5)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// (compose square increment)(5) = square(increment(5)) = square(6) = 36
			if result.String() != "36" {
				t.Errorf("expected 36, got %s", result.String())
			}
		})

		// Test apply-n function
		t.Run("apply-n function", func(t *testing.T) {
			// Define increment function if not already defined
			_, err := evalExpr("(def increment (fn [x] (+ x 1)))")
			if err != nil {
				t.Fatalf("failed to define increment: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"apply-n 0 times", "(apply-n increment 0 5)", "5"},
				{"apply-n 3 times", "(apply-n increment 3 5)", "8"},
				{"apply-n with square", "(apply-n square 2 2)", "16"}, // square(square(2)) = square(4) = 16
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					result, err := evalExpr(tt.input)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if result.String() != tt.expected {
						t.Errorf("expected %s, got %s", tt.expected, result.String())
					}
				})
			}
		})
	})
}

// TestCoreLibraryModuleSystem tests that the core library module system works correctly
func TestCoreLibraryModuleSystem(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Helper function to evaluate expressions
	evalExpr := func(input string) (types.Value, error) {
		tok := tokenizer.NewTokenizer(input)
		tokens, err := tok.TokenizeWithError()
		if err != nil {
			return nil, err
		}

		p := parser.NewParser(tokens)
		expr, err := p.Parse()
		if err != nil {
			return nil, err
		}

		return evaluator.Eval(expr)
	}

	// Load the core library
	coreLibPath := filepath.Join("..", "..", "library", "core.lisp")
	loadExpr := `(load "` + coreLibPath + `")`
	_, err := evalExpr(loadExpr)
	if err != nil {
		t.Fatalf("Failed to load core library: %v", err)
	}

	t.Run("qualified access without import", func(t *testing.T) {
		// Test qualified access to core functions without importing
		result, err := evalExpr("(core.factorial 5)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.String() != "120" {
			t.Errorf("expected 120, got %s", result.String())
		}

		// Test another function
		result, err = evalExpr("(core.fibonacci 7)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.String() != "13" {
			t.Errorf("expected 13, got %s", result.String())
		}
	})

	t.Run("private functions not accessible", func(t *testing.T) {
		// Attempt to access private fact-tail function should fail
		_, err := evalExpr("(fact-tail 5 1)")
		if err == nil {
			t.Error("expected error when accessing private function without import")
		}

		// Even with qualified access, private functions should not be accessible
		_, err = evalExpr("(core.fact-tail 5 1)")
		if err == nil {
			t.Error("expected error when accessing private function with qualified access")
		}
	})

	t.Run("exported functions list", func(t *testing.T) {
		// Check that the module exports the expected functions
		_, err := evalExpr("(import core)")
		if err != nil {
			t.Fatalf("Failed to import core module: %v", err)
		}

		// All these functions should be available after import
		exportedFunctions := []string{"factorial", "fibonacci", "gcd", "lcm", "abs", "min", "max", "length-sq", "all", "any", "take", "drop", "compose", "apply-n"}

		for _, funcName := range exportedFunctions {
			// Test that each function is callable
			switch funcName {
			case "factorial", "fibonacci":
				_, err := evalExpr("(" + funcName + " 3)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "gcd", "lcm", "min", "max":
				_, err := evalExpr("(" + funcName + " 12 8)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "abs":
				_, err := evalExpr("(" + funcName + " -5)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "length-sq":
				_, err := evalExpr("(" + funcName + " (list 1 2 3))")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			}
		}
	})
}
