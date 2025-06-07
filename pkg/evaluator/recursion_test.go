package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// TestTailCallOptimization tests that tail recursive functions don't cause stack overflow
func TestTailCallOptimization(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	t.Run("tail recursive factorial", func(t *testing.T) {
		// Define a tail-recursive factorial function
		defunExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "fact-tail"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
						&types.SymbolExpr{Name: "acc"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.SymbolExpr{Name: "acc"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "fact-tail"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "-"},
										&types.SymbolExpr{Name: "n"},
										&types.NumberExpr{Value: 1},
									},
								},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "*"},
										&types.SymbolExpr{Name: "n"},
										&types.SymbolExpr{Name: "acc"},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err := evaluator.Eval(defunExpr)
		if err != nil {
			t.Fatalf("failed to define tail-recursive factorial: %v", err)
		}

		// Test small factorial
		callExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "fact-tail"},
				&types.NumberExpr{Value: 5},
				&types.NumberExpr{Value: 1},
			},
		}

		result, err := evaluator.Eval(callExpr)
		if err != nil {
			t.Fatalf("failed to call tail-recursive factorial: %v", err)
		}

		expected := types.NumberValue(120) // 5! = 120
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("large tail recursive computation should not stack overflow", func(t *testing.T) {
		// Define a tail-recursive sum function
		defunExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "sum-tail"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
						&types.SymbolExpr{Name: "acc"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.SymbolExpr{Name: "acc"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "sum-tail"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "-"},
										&types.SymbolExpr{Name: "n"},
										&types.NumberExpr{Value: 1},
									},
								},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "+"},
										&types.SymbolExpr{Name: "acc"},
										&types.SymbolExpr{Name: "n"},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err := evaluator.Eval(defunExpr)
		if err != nil {
			t.Fatalf("failed to define sum-tail function: %v", err)
		}

		// Test with a large number that would normally cause stack overflow
		// Sum from 1 to 1000: 1000 * 1001 / 2 = 500500
		callExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "sum-tail"},
				&types.NumberExpr{Value: 1000},
				&types.NumberExpr{Value: 0},
			},
		}

		result, err := evaluator.Eval(callExpr)
		if err != nil {
			t.Fatalf("failed to call sum-tail with large number: %v", err)
		}

		expected := types.NumberValue(500500)
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

// TestNonTailRecursion tests that non-tail recursive functions still work correctly
func TestNonTailRecursion(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	t.Run("non-tail recursive factorial", func(t *testing.T) {
		// Define a non-tail-recursive factorial function
		defunExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "factorial"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.NumberExpr{Value: 1},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "*"},
								&types.SymbolExpr{Name: "n"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "factorial"},
										&types.ListExpr{
											Elements: []types.Expr{
												&types.SymbolExpr{Name: "-"},
												&types.SymbolExpr{Name: "n"},
												&types.NumberExpr{Value: 1},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err := evaluator.Eval(defunExpr)
		if err != nil {
			t.Fatalf("failed to define non-tail factorial: %v", err)
		}

		// Test small factorial (this is NOT in tail position)
		callExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "factorial"},
				&types.NumberExpr{Value: 5},
			},
		}

		result, err := evaluator.Eval(callExpr)
		if err != nil {
			t.Fatalf("failed to call non-tail factorial: %v", err)
		}

		expected := types.NumberValue(120) // 5! = 120
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

// TestMutualTailRecursion tests mutually tail-recursive functions
func TestMutualTailRecursion(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	t.Run("mutually recursive even/odd", func(t *testing.T) {
		// Define mutually tail-recursive even? function
		evenDefExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "even?"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.BooleanExpr{Value: true},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "odd?"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "-"},
										&types.SymbolExpr{Name: "n"},
										&types.NumberExpr{Value: 1},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err := evaluator.Eval(evenDefExpr)
		if err != nil {
			t.Fatalf("failed to define even? function: %v", err)
		}

		// Define mutually tail-recursive odd? function
		oddDefExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defun"},
				&types.SymbolExpr{Name: "odd?"},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "n"},
					},
				},
				&types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "="},
								&types.SymbolExpr{Name: "n"},
								&types.NumberExpr{Value: 0},
							},
						},
						&types.BooleanExpr{Value: false},
						&types.ListExpr{
							Elements: []types.Expr{
								&types.SymbolExpr{Name: "even?"},
								&types.ListExpr{
									Elements: []types.Expr{
										&types.SymbolExpr{Name: "-"},
										&types.SymbolExpr{Name: "n"},
										&types.NumberExpr{Value: 1},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err = evaluator.Eval(oddDefExpr)
		if err != nil {
			t.Fatalf("failed to define odd? function: %v", err)
		}

		// Test even? with even number
		evenCallExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "even?"},
				&types.NumberExpr{Value: 100},
			},
		}

		result, err := evaluator.Eval(evenCallExpr)
		if err != nil {
			t.Fatalf("failed to call even? function: %v", err)
		}

		expected := types.BooleanValue(true)
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}

		// Test odd? with odd number
		oddCallExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "odd?"},
				&types.NumberExpr{Value: 99},
			},
		}

		result, err = evaluator.Eval(oddCallExpr)
		if err != nil {
			t.Fatalf("failed to call odd? function: %v", err)
		}

		expected = types.BooleanValue(true)
		if !valuesEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}
