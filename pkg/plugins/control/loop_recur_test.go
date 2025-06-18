package control

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

func TestLoopRecur_SimpleLoop(t *testing.T) {
	plugin := NewControlPlugin()
	registry := registry.NewRegistry()

	err := plugin.RegisterFunctions(registry)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	// Check that loop and recur are registered
	if !registry.Has("loop") {
		t.Error("loop function not registered")
	}
	if !registry.Has("recur") {
		t.Error("recur function not registered")
	}
}

func TestLoopRecur_BasicFunctionality(t *testing.T) {
	plugin := NewControlPlugin()
	evaluator := newMockEvaluator()

	// Test a simple loop that counts down from 3 to 0
	// (loop [i 3] (if (= i 0) i (recur (- i 1))))

	// Create binding vector [i 3]
	bindingVector := &types.BracketExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "i"},
			&types.NumberExpr{Value: 3},
		},
	}

	// Create conditional body: (if (= i 0) i (recur (- i 1)))
	// For simplicity, we'll just test that the loop function accepts the right structure
	body := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "if"},
			// ... (simplified for testing)
		},
	}

	args := []types.Expr{bindingVector, body}

	// This test mainly checks that the structure is accepted without error
	// A full test would require a more sophisticated mock evaluator
	_, err := plugin.evalLoop(evaluator, args)

	// We expect this to work without throwing an error about malformed input
	// The actual loop execution would require a more complex setup
	if err != nil && err.Error() != "recur called outside of loop" {
		t.Errorf("Loop evaluation failed with unexpected error: %v", err)
	}
}

func TestRecur_ThrowsException(t *testing.T) {
	plugin := NewControlPlugin()
	evaluator := newMockEvaluator()

	// Test that recur throws a RecurException
	args := []types.Expr{
		&types.NumberExpr{Value: 42},
		&types.StringExpr{Value: "test"},
	}

	_, err := plugin.evalRecur(evaluator, args)

	// Should return a RecurException
	if err == nil {
		t.Error("recur should throw an exception")
	}

	recurErr, ok := err.(*types.RecurException)
	if !ok {
		t.Errorf("Expected RecurException, got %T", err)
	}

	if len(recurErr.Args) != 2 {
		t.Errorf("Expected 2 arguments in RecurException, got %d", len(recurErr.Args))
	}
}

func TestLoop_InvalidBindingVector(t *testing.T) {
	plugin := NewControlPlugin()
	evaluator := newMockEvaluator()

	// Test with invalid binding vector (not a BracketExpr)
	args := []types.Expr{
		&types.ListExpr{Elements: []types.Expr{}}, // Should be BracketExpr
	}

	_, err := plugin.evalLoop(evaluator, args)

	if err == nil || err.Error() != "loop first argument must be a binding vector" {
		t.Errorf("Expected error about binding vector, got: %v", err)
	}
}

func TestLoop_OddBindingElements(t *testing.T) {
	plugin := NewControlPlugin()
	evaluator := newMockEvaluator()

	// Test with odd number of binding elements
	bindingVector := &types.BracketExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "i"},
			// Missing value for 'i'
		},
	}

	args := []types.Expr{bindingVector}

	_, err := plugin.evalLoop(evaluator, args)

	if err == nil || err.Error() != "loop binding vector must have even number of elements" {
		t.Errorf("Expected error about even number of elements, got: %v", err)
	}
}
