package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// TestBasicGoroutine tests basic goroutine creation and waiting
func TestBasicGoroutine(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test (go (+ 1 2 3))
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "go"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 2},
					&types.NumberExpr{Value: 3},
				},
			},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	future, ok := result.(*types.FutureValue)
	if !ok {
		t.Fatalf("Expected FutureValue, got: %T", result)
	}

	// Wait for the result
	waitResult, waitErr := future.Wait()
	if waitErr != nil {
		t.Fatalf("Expected no error waiting for future, got: %v", waitErr)
	}

	if num, ok := waitResult.(types.NumberValue); !ok || num != 6 {
		t.Errorf("Expected 6, got: %v", waitResult)
	}
}

// TestGoroutineWithFunction tests goroutine with user-defined function
func TestGoroutineWithFunction(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// First define a function: (defn square [x] (* x x))
	defnExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "defn"},
			&types.SymbolExpr{Name: "square"},
			&types.BracketExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "x"},
				},
			},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "*"},
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "x"},
				},
			},
		},
	}

	_, err := evaluator.Eval(defnExpr)
	if err != nil {
		t.Fatalf("Expected no error defining function, got: %v", err)
	}

	// Now test (go (square 5))
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "go"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "square"},
					&types.NumberExpr{Value: 5},
				},
			},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	future, ok := result.(*types.FutureValue)
	if !ok {
		t.Fatalf("Expected FutureValue, got: %T", result)
	}

	// Wait for the result
	waitResult, waitErr := future.Wait()
	if waitErr != nil {
		t.Fatalf("Expected no error waiting for future, got: %v", waitErr)
	}

	if num, ok := waitResult.(types.NumberValue); !ok || num != 25 {
		t.Errorf("Expected 25, got: %v", waitResult)
	}
}

// TestGoroutineWait tests the go-wait function
func TestGoroutineWait(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create a future manually for testing
	future := types.NewFuture()
	future.SetResult(types.NumberValue(42))

	// Test (go-wait future)
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "go-wait"},
			&types.SymbolExpr{Name: "test-future"},
		},
	}

	// Set the future in the environment
	env.Set("test-future", future)

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if num, ok := result.(types.NumberValue); !ok || num != 42 {
		t.Errorf("Expected 42, got: %v", result)
	}
}

// TestMultipleGoroutines tests multiple goroutines running concurrently
func TestMultipleGoroutines(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	var futures []*types.FutureValue

	// Create multiple goroutines
	for i := 1; i <= 5; i++ {
		expr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "go"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "*"},
						&types.NumberExpr{Value: float64(i)},
						&types.NumberExpr{Value: float64(i)},
					},
				},
			},
		}

		result, err := evaluator.Eval(expr)
		if err != nil {
			t.Fatalf("Expected no error for goroutine %d, got: %v", i, err)
		}

		future, ok := result.(*types.FutureValue)
		if !ok {
			t.Fatalf("Expected FutureValue for goroutine %d, got: %T", i, result)
		}

		futures = append(futures, future)
	}

	// Wait for all results
	expectedResults := []float64{1, 4, 9, 16, 25}
	for i, future := range futures {
		waitResult, waitErr := future.Wait()
		if waitErr != nil {
			t.Fatalf("Expected no error waiting for future %d, got: %v", i, waitErr)
		}

		if num, ok := waitResult.(types.NumberValue); !ok || num != types.NumberValue(expectedResults[i]) {
			t.Errorf("Expected %v for goroutine %d, got: %v", expectedResults[i], i, waitResult)
		}
	}
}

// TestGoroutineWithError tests goroutine error handling
func TestGoroutineWithError(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test (go (/ 1 0)) - division by zero
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "go"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "/"},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 0},
				},
			},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("Expected no error creating goroutine, got: %v", err)
	}

	future, ok := result.(*types.FutureValue)
	if !ok {
		t.Fatalf("Expected FutureValue, got: %T", result)
	}

	// Wait for the result - should get an error
	_, waitErr := future.Wait()
	if waitErr == nil {
		t.Fatalf("Expected an error from division by zero, got none")
	}
}

// TestGoroutineWithAtoms tests goroutines with shared atoms
func TestGoroutineWithAtoms(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create a shared atom: (def counter (atom 0))
	atomExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "def"},
			&types.SymbolExpr{Name: "counter"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "atom"},
					&types.NumberExpr{Value: 0},
				},
			},
		},
	}

	_, err := evaluator.Eval(atomExpr)
	if err != nil {
		t.Fatalf("Expected no error creating atom, got: %v", err)
	}

	// Start multiple goroutines that increment the counter
	var futures []*types.FutureValue
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		expr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "go"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "swap!"},
						&types.SymbolExpr{Name: "counter"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "fn"},
								&types.BracketExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "x"},
									},
								},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "+"},
										&types.SymbolExpr{Name: "x"},
										&types.NumberExpr{Value: 1},
									},
								},
							},
						},
					},
				},
			},
		}

		result, err := evaluator.Eval(expr)
		if err != nil {
			t.Fatalf("Expected no error for goroutine %d, got: %v", i, err)
		}

		future, ok := result.(*types.FutureValue)
		if !ok {
			t.Fatalf("Expected FutureValue for goroutine %d, got: %T", i, result)
		}

		futures = append(futures, future)
	}

	// Wait for all goroutines to complete
	for i, future := range futures {
		_, waitErr := future.Wait()
		if waitErr != nil {
			t.Fatalf("Expected no error waiting for future %d, got: %v", i, waitErr)
		}
	}

	// Check final value of counter
	counterValue, ok := env.Get("counter")
	if !ok {
		t.Fatal("Counter not found in environment")
	}

	atom, ok := counterValue.(*types.AtomValue)
	if !ok {
		t.Fatalf("Expected AtomValue, got: %T", counterValue)
	}

	finalValue := atom.Value()
	if num, ok := finalValue.(types.NumberValue); !ok || num != types.NumberValue(numGoroutines) {
		t.Errorf("Expected %d, got: %v", numGoroutines, finalValue)
	}
}

// TestGoroutineWaitAll tests waiting for multiple goroutines
func TestGoroutineWaitAll(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create multiple goroutines
	var futures []types.Value
	for i := 1; i <= 3; i++ {
		goExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "go"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "*"},
						&types.NumberExpr{Value: float64(i)},
						&types.NumberExpr{Value: 10},
					},
				},
			},
		}

		future, err := evaluator.Eval(goExpr)
		if err != nil {
			t.Fatalf("Expected no error creating goroutine %d, got: %v", i, err)
		}

		futures = append(futures, future)
	}

	// Create a list of futures
	listValue := &types.ListValue{Elements: futures}
	env.Set("future-list", listValue)

	// Test (go-wait-all future-list)
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "go-wait-all"},
			&types.SymbolExpr{Name: "future-list"},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	resultList, ok := result.(*types.ListValue)
	if !ok {
		t.Fatalf("Expected ListValue, got: %T", result)
	}

	expectedResults := []types.NumberValue{10, 20, 30}
	if len(resultList.Elements) != len(expectedResults) {
		t.Fatalf("Expected %d results, got %d", len(expectedResults), len(resultList.Elements))
	}

	for i, expected := range expectedResults {
		if num, ok := resultList.Elements[i].(types.NumberValue); !ok || num != expected {
			t.Errorf("Expected %v at index %d, got: %v", expected, i, resultList.Elements[i])
		}
	}
}
