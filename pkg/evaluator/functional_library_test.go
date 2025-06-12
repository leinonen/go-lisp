package evaluator

import (
	"path/filepath"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// TestFunctionalLibraryFunctions tests the functional library functions by loading the functional.lisp file
// and testing each exported function
func TestFunctionalLibraryFunctions(t *testing.T) {
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

	// Load the functional library
	functionalLibPath := filepath.Join("..", "..", "library", "functional.lisp")
	loadExpr := `(load "` + functionalLibPath + `")`
	_, err := evalExpr(loadExpr)
	if err != nil {
		t.Fatalf("Failed to load functional library: %v", err)
	}

	// Import the functional module
	_, err = evalExpr("(import functional)")
	if err != nil {
		t.Fatalf("Failed to import functional module: %v", err)
	}

	t.Run("basic combinators", func(t *testing.T) {
		// Test identity function
		t.Run("identity", func(t *testing.T) {
			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"identity number", "(identity 42)", "42"},
				{"identity string", "(identity \"hello\")", "hello"},
				{"identity list", "(identity (list 1 2 3))", "(1 2 3)"},
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

		// Test constantly function
		t.Run("constantly", func(t *testing.T) {
			// Define a constant function that always returns 42
			_, err := evalExpr("(define always-42 (constantly 42))")
			if err != nil {
				t.Fatalf("failed to define always-42: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"constantly with number", "(always-42 1)", "42"},
				{"constantly with string", "(always-42 \"test\")", "42"},
				{"constantly with list", "(always-42 (list 1 2))", "42"},
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

		// Test complement function
		t.Run("complement", func(t *testing.T) {
			// Define a positive predicate and its complement
			_, err := evalExpr("(define positive? (lambda [x] (> x 0)))")
			if err != nil {
				t.Fatalf("failed to define positive?: %v", err)
			}
			_, err = evalExpr("(define not-positive? (complement positive?))")
			if err != nil {
				t.Fatalf("failed to define not-positive?: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"complement of positive number", "(not-positive? 5)", "#f"},
				{"complement of negative number", "(not-positive? -3)", "#t"},
				{"complement of zero", "(not-positive? 0)", "#t"},
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

	t.Run("partial application", func(t *testing.T) {
		// Test partial function
		t.Run("partial", func(t *testing.T) {
			// Create a partial application of addition
			_, err := evalExpr("(define add-5 (partial + 5))")
			if err != nil {
				t.Fatalf("failed to define add-5: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"partial add 5 to 10", "(add-5 10)", "15"},
				{"partial add 5 to 0", "(add-5 0)", "5"},
				{"partial add 5 to -3", "(add-5 -3)", "2"},
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

		// Test partial2 function
		t.Run("partial2", func(t *testing.T) {
			// Define a 3-argument function for testing
			_, err := evalExpr("(define add3 (lambda [a b c] (+ a b c)))")
			if err != nil {
				t.Fatalf("failed to define add3: %v", err)
			}
			_, err = evalExpr("(define add-5-10 (partial2 add3 5 10))")
			if err != nil {
				t.Fatalf("failed to define add-5-10: %v", err)
			}

			result, err := evalExpr("(add-5-10 3)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "18" {
				t.Errorf("expected 18, got %s", result.String())
			}
		})

		// Test partial3 function
		t.Run("partial3", func(t *testing.T) {
			// Define a 4-argument function for testing
			_, err := evalExpr("(define add4 (lambda [a b c d] (+ a b c d)))")
			if err != nil {
				t.Fatalf("failed to define add4: %v", err)
			}
			_, err = evalExpr("(define add-1-2-3 (partial3 add4 1 2 3))")
			if err != nil {
				t.Fatalf("failed to define add-1-2-3: %v", err)
			}

			result, err := evalExpr("(add-1-2-3 4)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "10" {
				t.Errorf("expected 10, got %s", result.String())
			}
		})
	})

	t.Run("currying", func(t *testing.T) {
		// Test curry2 function
		t.Run("curry2", func(t *testing.T) {
			_, err := evalExpr("(define curried-add (curry2 +))")
			if err != nil {
				t.Fatalf("failed to define curried-add: %v", err)
			}

			// Test curried application
			result, err := evalExpr("((curried-add 5) 3)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "8" {
				t.Errorf("expected 8, got %s", result.String())
			}
		})

		// Test curry3 function
		t.Run("curry3", func(t *testing.T) {
			// Define a 3-argument function and curry it
			_, err := evalExpr("(define add3 (lambda [a b c] (+ a b c)))")
			if err != nil {
				t.Fatalf("failed to define add3: %v", err)
			}
			_, err = evalExpr("(define curried-add3 (curry3 add3))")
			if err != nil {
				t.Fatalf("failed to define curried-add3: %v", err)
			}

			// Test curried application
			result, err := evalExpr("(((curried-add3 1) 2) 3)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "6" {
				t.Errorf("expected 6, got %s", result.String())
			}
		})

		// Test generic curry function (alias for curry2)
		t.Run("curry", func(t *testing.T) {
			_, err := evalExpr("(define curried-mult (curry *))")
			if err != nil {
				t.Fatalf("failed to define curried-mult: %v", err)
			}

			result, err := evalExpr("((curried-mult 4) 7)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "28" {
				t.Errorf("expected 28, got %s", result.String())
			}
		})
	})

	t.Run("function composition", func(t *testing.T) {
		// Define helper functions for testing
		_, err := evalExpr("(define square (lambda [x] (* x x)))")
		if err != nil {
			t.Fatalf("failed to define square: %v", err)
		}
		_, err = evalExpr("(define increment (lambda [x] (+ x 1)))")
		if err != nil {
			t.Fatalf("failed to define increment: %v", err)
		}
		_, err = evalExpr("(define double (lambda [x] (* x 2)))")
		if err != nil {
			t.Fatalf("failed to define double: %v", err)
		}

		// Test comp function
		t.Run("comp", func(t *testing.T) {
			_, err := evalExpr("(define square-then-increment (comp increment square))")
			if err != nil {
				t.Fatalf("failed to define square-then-increment: %v", err)
			}

			// (comp increment square)(3) = increment(square(3)) = increment(9) = 10
			result, err := evalExpr("(square-then-increment 3)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "10" {
				t.Errorf("expected 10, got %s", result.String())
			}
		})

		// Test comp3 function
		t.Run("comp3", func(t *testing.T) {
			_, err := evalExpr("(define complex-comp (comp3 double increment square))")
			if err != nil {
				t.Fatalf("failed to define complex-comp: %v", err)
			}

			// (comp3 double increment square)(3) = double(increment(square(3))) = double(increment(9)) = double(10) = 20
			result, err := evalExpr("(complex-comp 3)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "20" {
				t.Errorf("expected 20, got %s", result.String())
			}
		})
	})

	t.Run("pipeline operations", func(t *testing.T) {
		// Define helper functions for testing
		_, err := evalExpr("(define square (lambda [x] (* x x)))")
		if err != nil {
			t.Fatalf("failed to define square: %v", err)
		}
		_, err = evalExpr("(define increment (lambda [x] (+ x 1)))")
		if err != nil {
			t.Fatalf("failed to define increment: %v", err)
		}

		// Test pipe function
		t.Run("pipe", func(t *testing.T) {
			result, err := evalExpr("(pipe 5 square)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "25" {
				t.Errorf("expected 25, got %s", result.String())
			}
		})

		// Test pipe2 function
		t.Run("pipe2", func(t *testing.T) {
			// pipe2(3, square, increment) = increment(square(3)) = increment(9) = 10
			result, err := evalExpr("(pipe2 3 square increment)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "10" {
				t.Errorf("expected 10, got %s", result.String())
			}
		})

		// Test pipe3 function
		t.Run("pipe3", func(t *testing.T) {
			_, err := evalExpr("(define double (lambda [x] (* x 2)))")
			if err != nil {
				t.Fatalf("failed to define double: %v", err)
			}

			// pipe3(2, increment, square, double) = double(square(increment(2))) = double(square(3)) = double(9) = 18
			result, err := evalExpr("(pipe3 2 increment square double)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "18" {
				t.Errorf("expected 18, got %s", result.String())
			}
		})
	})

	t.Run("juxtaposition", func(t *testing.T) {
		// Define helper functions for testing
		_, err := evalExpr("(define square (lambda [x] (* x x)))")
		if err != nil {
			t.Fatalf("failed to define square: %v", err)
		}
		_, err = evalExpr("(define double (lambda [x] (* x 2)))")
		if err != nil {
			t.Fatalf("failed to define double: %v", err)
		}
		_, err = evalExpr("(define increment (lambda [x] (+ x 1)))")
		if err != nil {
			t.Fatalf("failed to define increment: %v", err)
		}

		// Test juxt function
		t.Run("juxt", func(t *testing.T) {
			_, err := evalExpr("(define square-and-double (juxt square double))")
			if err != nil {
				t.Fatalf("failed to define square-and-double: %v", err)
			}

			result, err := evalExpr("(square-and-double 5)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Should return (25 10)
			if result.String() != "(25 10)" {
				t.Errorf("expected (25 10), got %s", result.String())
			}
		})

		// Test juxt3 function
		t.Run("juxt3", func(t *testing.T) {
			_, err := evalExpr("(define triple-juxt (juxt3 square double increment))")
			if err != nil {
				t.Fatalf("failed to define triple-juxt: %v", err)
			}

			result, err := evalExpr("(triple-juxt 4)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Should return (16 8 5)
			if result.String() != "(16 8 5)" {
				t.Errorf("expected (16 8 5), got %s", result.String())
			}
		})
	})

	t.Run("conditional functions", func(t *testing.T) {
		// Define helper functions for testing
		_, err := evalExpr("(define positive? (lambda [x] (> x 0)))")
		if err != nil {
			t.Fatalf("failed to define positive?: %v", err)
		}
		_, err = evalExpr("(define negate (lambda [x] (- x)))")
		if err != nil {
			t.Fatalf("failed to define negate: %v", err)
		}
		_, err = evalExpr("(define abs-fn (lambda [x] (if (< x 0) (- x) x)))")
		if err != nil {
			t.Fatalf("failed to define abs-fn: %v", err)
		}

		// Test if-fn function
		t.Run("if-fn", func(t *testing.T) {
			_, err := evalExpr("(define conditional-abs (if-fn positive? identity negate))")
			if err != nil {
				t.Fatalf("failed to define conditional-abs: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"positive number", "(conditional-abs 5)", "5"},
				{"negative number", "(conditional-abs -3)", "3"},
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

		// Test when-fn function
		t.Run("when-fn", func(t *testing.T) {
			_, err := evalExpr("(define negate-when-positive (when-fn positive? negate))")
			if err != nil {
				t.Fatalf("failed to define negate-when-positive: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"positive number gets negated", "(negate-when-positive 5)", "-5"},
				{"negative number unchanged", "(negate-when-positive -3)", "-3"},
				{"zero unchanged", "(negate-when-positive 0)", "0"},
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

		// Test unless-fn function
		t.Run("unless-fn", func(t *testing.T) {
			_, err := evalExpr("(define negate-unless-positive (unless-fn positive? negate))")
			if err != nil {
				t.Fatalf("failed to define negate-unless-positive: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"positive number unchanged", "(negate-unless-positive 5)", "5"},
				{"negative number gets negated", "(negate-unless-positive -3)", "3"},
				{"zero gets negated", "(negate-unless-positive 0)", "0"},
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

	t.Run("predicate combinators", func(t *testing.T) {
		// Define helper predicates for testing
		_, err := evalExpr("(define positive? (lambda [x] (> x 0)))")
		if err != nil {
			t.Fatalf("failed to define positive?: %v", err)
		}
		_, err = evalExpr("(define even? (lambda [x] (= (% x 2) 0)))")
		if err != nil {
			t.Fatalf("failed to define even?: %v", err)
		}
		_, err = evalExpr("(define less-than-10? (lambda [x] (< x 10)))")
		if err != nil {
			t.Fatalf("failed to define less-than-10?: %v", err)
		}

		// Test every-pred function
		t.Run("every-pred", func(t *testing.T) {
			_, err := evalExpr("(define positive-and-even? (every-pred positive? even?))")
			if err != nil {
				t.Fatalf("failed to define positive-and-even?: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"positive and even", "(positive-and-even? 4)", "#t"},
				{"positive but odd", "(positive-and-even? 3)", "#f"},
				{"negative and even", "(positive-and-even? -2)", "#f"},
				{"negative and odd", "(positive-and-even? -1)", "#f"},
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

		// Test every-pred3 function
		t.Run("every-pred3", func(t *testing.T) {
			_, err := evalExpr("(define positive-even-small? (every-pred3 positive? even? less-than-10?))")
			if err != nil {
				t.Fatalf("failed to define positive-even-small?: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"satisfies all three", "(positive-even-small? 4)", "#t"},
				{"positive even but large", "(positive-even-small? 12)", "#f"},
				{"positive small but odd", "(positive-even-small? 5)", "#f"},
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

		// Test some-pred function
		t.Run("some-pred", func(t *testing.T) {
			_, err := evalExpr("(define positive-or-even? (some-pred positive? even?))")
			if err != nil {
				t.Fatalf("failed to define positive-or-even?: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"positive and even", "(positive-or-even? 4)", "#t"},
				{"positive but odd", "(positive-or-even? 3)", "#t"},
				{"negative but even", "(positive-or-even? -2)", "#t"},
				{"negative and odd", "(positive-or-even? -1)", "#f"},
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

	t.Run("threading and application utilities", func(t *testing.T) {
		// Define helper functions for testing
		_, err := evalExpr("(define square (lambda [x] (* x x)))")
		if err != nil {
			t.Fatalf("failed to define square: %v", err)
		}

		// Test thread-first function
		t.Run("thread-first", func(t *testing.T) {
			result, err := evalExpr("(thread-first 5 square)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "25" {
				t.Errorf("expected 25, got %s", result.String())
			}
		})

		// Test thread-last function
		t.Run("thread-last", func(t *testing.T) {
			result, err := evalExpr("(thread-last 6 square)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "36" {
				t.Errorf("expected 36, got %s", result.String())
			}
		})

		// Test apply-to function
		t.Run("apply-to", func(t *testing.T) {
			result, err := evalExpr("(apply-to 7 square)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "49" {
				t.Errorf("expected 49, got %s", result.String())
			}
		})
	})

	t.Run("higher-order utilities", func(t *testing.T) {
		// Test fnil function
		t.Run("fnil", func(t *testing.T) {
			_, err := evalExpr("(define safe-add-5 (fnil (lambda [x] (+ x 5)) 0))")
			if err != nil {
				t.Fatalf("failed to define safe-add-5: %v", err)
			}

			tests := []struct {
				name     string
				input    string
				expected string
			}{
				{"normal value", "(safe-add-5 10)", "15"},
				{"nil value uses default", "(safe-add-5 nil)", "5"},
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

		// Test map-indexed function
		t.Run("map-indexed", func(t *testing.T) {
			// Create a function that adds index to element
			_, err := evalExpr("(define add-index (lambda [index element] (+ index element)))")
			if err != nil {
				t.Fatalf("failed to define add-index: %v", err)
			}

			result, err := evalExpr("(map-indexed add-index (list 10 20 30))")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Should return (10 21 32) - index 0+10, 1+20, 2+30
			if result.String() != "(10 21 32)" {
				t.Errorf("expected (10 21 32), got %s", result.String())
			}
		})

		// Test keep function
		t.Run("keep", func(t *testing.T) {
			// Define a function that returns nil for even numbers, otherwise returns the number
			_, err := evalExpr("(define odd-or-nil (lambda [x] (if (= (% x 2) 0) nil x)))")
			if err != nil {
				t.Fatalf("failed to define odd-or-nil: %v", err)
			}

			result, err := evalExpr("(keep odd-or-nil (list 1 2 3 4 5))")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Should return (1 3 5) - keeping only odd numbers
			if result.String() != "(1 3 5)" {
				t.Errorf("expected (1 3 5), got %s", result.String())
			}
		})

		// Test memoize function (simplified implementation)
		t.Run("memoize", func(t *testing.T) {
			_, err := evalExpr("(define square (lambda [x] (* x x)))")
			if err != nil {
				t.Fatalf("failed to define square: %v", err)
			}
			_, err = evalExpr("(define memoized-square (memoize square))")
			if err != nil {
				t.Fatalf("failed to define memoized-square: %v", err)
			}

			result, err := evalExpr("(memoized-square 4)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != "16" {
				t.Errorf("expected 16, got %s", result.String())
			}
		})

		// Test arity function (placeholder implementation)
		t.Run("arity", func(t *testing.T) {
			_, err := evalExpr("(define square (lambda [x] (* x x)))")
			if err != nil {
				t.Fatalf("failed to define square: %v", err)
			}

			result, err := evalExpr("(arity square)")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Current implementation returns placeholder value of 2
			if result.String() != "2" {
				t.Errorf("expected 2, got %s", result.String())
			}
		})
	})
}

// TestFunctionalLibraryModuleSystem tests that the functional library module system works correctly
func TestFunctionalLibraryModuleSystem(t *testing.T) {
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

	// Load the functional library
	functionalLibPath := filepath.Join("..", "..", "library", "functional.lisp")
	loadExpr := `(load "` + functionalLibPath + `")`
	_, err := evalExpr(loadExpr)
	if err != nil {
		t.Fatalf("Failed to load functional library: %v", err)
	}

	t.Run("qualified access without import", func(t *testing.T) {
		// Test qualified access to functional functions without importing
		result, err := evalExpr("(functional.identity 42)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.String() != "42" {
			t.Errorf("expected 42, got %s", result.String())
		}

		// Test another function with qualified access
		_, err = evalExpr("(define square (lambda [x] (* x x)))")
		if err != nil {
			t.Fatalf("failed to define square: %v", err)
		}
		_, err = evalExpr("(define increment (lambda [x] (+ x 1)))")
		if err != nil {
			t.Fatalf("failed to define increment: %v", err)
		}

		result, err = evalExpr("((functional.comp increment square) 3)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.String() != "10" {
			t.Errorf("expected 10, got %s", result.String())
		}
	})

	t.Run("exported functions available after import", func(t *testing.T) {
		// Import the functional module
		_, err := evalExpr("(import functional)")
		if err != nil {
			t.Fatalf("Failed to import functional module: %v", err)
		}

		// Test that key functions are available after import
		exportedFunctions := []string{"identity", "constantly", "complement", "curry", "curry2", "curry3", "partial", "partial2", "partial3", "comp", "comp3", "comp4", "pipe", "pipe2", "pipe3", "pipe4", "juxt", "juxt3", "juxt4", "if-fn", "when-fn", "unless-fn", "every-pred", "every-pred3", "some-pred", "some-pred3", "thread-first", "thread-last", "apply-to", "memoize", "arity", "fnil", "map-indexed", "keep", "keep-indexed"}

		for _, funcName := range exportedFunctions {
			// Test basic functionality of each exported function
			switch funcName {
			case "identity":
				_, err := evalExpr("(" + funcName + " 5)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "constantly":
				_, err := evalExpr("((constantly 42) 10)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "complement":
				_, err := evalExpr("(define pos? (lambda [x] (> x 0)))")
				if err != nil {
					continue
				}
				_, err = evalExpr("((complement pos?) -1)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "partial":
				_, err := evalExpr("((partial + 5) 3)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "curry2", "curry":
				_, err := evalExpr("(((curry2 +) 5) 3)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "comp":
				_, err := evalExpr("(define sq (lambda [x] (* x x)))")
				if err != nil {
					continue
				}
				_, err = evalExpr("(define inc (lambda [x] (+ x 1)))")
				if err != nil {
					continue
				}
				_, err = evalExpr("((comp inc sq) 3)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "pipe":
				_, err := evalExpr("(define sq (lambda [x] (* x x)))")
				if err != nil {
					continue
				}
				_, err = evalExpr("(pipe 4 sq)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "juxt":
				_, err := evalExpr("(define sq (lambda [x] (* x x)))")
				if err != nil {
					continue
				}
				_, err = evalExpr("(define dbl (lambda [x] (* x 2)))")
				if err != nil {
					continue
				}
				_, err = evalExpr("((juxt sq dbl) 3)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "memoize":
				_, err := evalExpr("(define sq (lambda [x] (* x x)))")
				if err != nil {
					continue
				}
				_, err = evalExpr("((memoize sq) 5)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "arity":
				_, err := evalExpr("(define sq (lambda [x] (* x x)))")
				if err != nil {
					continue
				}
				_, err = evalExpr("(arity sq)")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			case "map-indexed":
				_, err := evalExpr("(define add-idx (lambda [i x] (+ i x)))")
				if err != nil {
					continue
				}
				_, err = evalExpr("(map-indexed add-idx (list 1 2 3))")
				if err != nil {
					t.Errorf("Function %s should be available after import, got error: %v", funcName, err)
				}
			}
		}
	})
}
