package kernel

import (
	"strings"
	"testing"
)

func TestLoopRecur(t *testing.T) {
	t.Run("BasicLoop", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Simple loop that counts down from 3 to 0
		// (loop [x 3] (if (= x 0) x (recur (- x 1))))
		result, err := repl.Eval(`(loop [x 3] (if (= x 0) x (recur (- x 1))))`)
		if err != nil {
			t.Fatalf("Error evaluating basic loop: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 0.0 {
			t.Errorf("Expected 0, got %v", result)
		}
	})

	t.Run("LoopWithMultipleBindings", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Loop with multiple bindings - compute sum
		// (loop [i 5 sum 0] (if (= i 0) sum (recur (- i 1) (+ sum i))))
		result, err := repl.Eval(`(loop [i 5 sum 0] (if (= i 0) sum (recur (- i 1) (+ sum i))))`)
		if err != nil {
			t.Fatalf("Error evaluating loop with multiple bindings: %v", err)
		}

		// Should compute 5+4+3+2+1 = 15
		if num, ok := result.(Number); !ok || float64(num) != 15.0 {
			t.Errorf("Expected 15, got %v", result)
		}
	})

	t.Run("NestedLoops", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Nested loops - compute factorial using loops
		// (loop [n 5 result 1]
		//   (if (= n 0)
		//     result
		//     (recur (- n 1) (* result n))))
		result, err := repl.Eval(`(loop [n 5 result 1] (if (= n 0) result (recur (- n 1) (* result n))))`)
		if err != nil {
			t.Fatalf("Error evaluating factorial loop: %v", err)
		}

		// Should compute 5! = 120
		if num, ok := result.(Number); !ok || float64(num) != 120.0 {
			t.Errorf("Expected 120, got %v", result)
		}
	})

	t.Run("RecurOutsideLoopError", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// recur outside of loop should error
		_, err := repl.Eval(`(recur 1 2)`)
		if err == nil {
			t.Error("Expected error for recur outside loop")
		}
		if !strings.Contains(err.Error(), "recur can only be used inside a loop") {
			t.Errorf("Expected specific error message to contain 'recur can only be used inside a loop', got: %v", err)
		}
	})

	t.Run("RecurWrongArityError", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// recur with wrong number of arguments should error
		_, err := repl.Eval(`(loop [x 1] (recur))`)
		if err == nil {
			t.Error("Expected error for recur with wrong arity")
		}
	})

	t.Run("LoopBodyMultipleExpressions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Loop with multiple expressions in body
		// (loop [x 3]
		//   (def temp (* x 2))  ; side effect
		//   (if (= x 0) temp (recur (- x 1))))
		result, err := repl.Eval(`(loop [x 3] (def temp (* x 2)) (if (= x 0) temp (recur (- x 1))))`)
		if err != nil {
			t.Fatalf("Error evaluating loop with multiple body expressions: %v", err)
		}

		// Should compute (* 0 2) = 0
		if num, ok := result.(Number); !ok || float64(num) != 0.0 {
			t.Errorf("Expected 0, got %v", result)
		}
	})

	t.Run("LoopInvalidBindings", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Invalid bindings - odd number of elements
		_, err := repl.Eval(`(loop [x] x)`)
		if err == nil {
			t.Error("Expected error for invalid bindings")
		}
	})

	t.Run("LoopNonVectorBindings", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Invalid bindings - not a vector
		_, err := repl.Eval(`(loop (x 1) x)`)
		if err == nil {
			t.Error("Expected error for non-vector bindings")
		}
	})

	t.Run("LoopWithNonSymbolBinding", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Invalid bindings - binding name is not a symbol
		_, err := repl.Eval(`(loop [1 2] x)`)
		if err == nil {
			t.Error("Expected error for non-symbol binding name")
		}
	})
}
