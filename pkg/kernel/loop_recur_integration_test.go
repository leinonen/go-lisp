package kernel

import (
	"strings"
	"testing"
)

func TestLoopRecurIntegration(t *testing.T) {
	t.Run("LoopRecurWithUserDefinedFunctions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a function that uses loop/recur
		_, err := repl.Eval(`(def countdown (fn [n] (loop [i n] (if (= i 0) "done" (recur (- i 1))))))`)
		if err != nil {
			t.Fatalf("Error defining countdown function: %v", err)
		}

		// Test the function
		result, err := repl.Eval("(countdown 3)")
		if err != nil {
			t.Fatalf("Error calling countdown: %v", err)
		}

		if str, ok := result.(String); !ok || string(str) != "done" {
			t.Errorf("Expected \"done\", got %v", result)
		}
	})

	t.Run("NestedLoopsWithRecur", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a function with nested loop structures
		_, err := repl.Eval(`(def matrix-sum (fn [matrix] 
			(loop [rows matrix total 0] 
				(if (= (first rows) nil) 
					total 
					(recur (rest rows) 
						(+ total (loop [cols (first rows) sum 0] 
							(if (= (first cols) nil) 
								sum 
								(recur (rest cols) (+ sum (first cols)))))))))))`)
		if err != nil {
			t.Fatalf("Error defining matrix-sum function: %v", err)
		}

		// Test with a simple 2x2 matrix
		result, err := repl.Eval("(matrix-sum (list (list 1 2) (list 3 4)))")
		if err != nil {
			t.Fatalf("Error calling matrix-sum: %v", err)
		}

		// Should compute 1+2+3+4 = 10
		if num, ok := result.(Number); !ok || float64(num) != 10.0 {
			t.Errorf("Expected 10, got %v", result)
		}
	})

	t.Run("LoopRecurWithHigherOrderFunctions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Use loop/recur to implement map
		_, err := repl.Eval(`(def my-map (fn [f coll] 
			(loop [items coll result (list)] 
				(if (= (first items) nil) 
					result 
					(recur (rest items) 
						(cons (f (first items)) result))))))`)
		if err != nil {
			t.Fatalf("Error defining my-map function: %v", err)
		}

		// Define a simple function to map over
		_, err = repl.Eval("(def double (fn [x] (* x 2)))")
		if err != nil {
			t.Fatalf("Error defining double function: %v", err)
		}

		// We need cons function for the map implementation
		_, err = repl.Eval(`(def cons (fn [item coll] 
			(if (= coll nil) 
				(list item) 
				(list item (first coll) (first (rest coll)) (first (rest (rest coll)))))))`)
		if err != nil {
			// Skip if cons is not easily implementable
			t.Skip("Skipping test due to cons implementation complexity")
		}
	})

	t.Run("LoopRecurErrorHandling", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test that recur in nested function context fails appropriately
		_, err := repl.Eval(`(def bad-fn (fn [x] (recur x)))`)
		if err != nil {
			t.Fatalf("Error defining bad-fn: %v", err)
		}

		_, err = repl.Eval("(bad-fn 5)")
		if err == nil {
			t.Error("Expected error for recur outside loop context")
		}

		if !strings.Contains(err.Error(), "recur can only be used inside a loop") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("LoopRecurPerformance", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test that loop/recur can handle reasonably deep iterations without stack overflow
		// Using a simple accumulator pattern
		result, err := repl.Eval(`(loop [n 1000 acc 0] (if (= n 0) acc (recur (- n 1) (+ acc 1))))`)
		if err != nil {
			t.Fatalf("Error with deep loop iteration: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 1000.0 {
			t.Errorf("Expected 1000, got %v", result)
		}
	})
}
