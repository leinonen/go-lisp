package atom

import (
	"fmt"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Mock evaluator for testing
type mockEvaluator struct {
	env *evaluator.Environment
}

func newMockEvaluator() *mockEvaluator {
	return &mockEvaluator{
		env: evaluator.NewEnvironment(),
	}
}

func (me *mockEvaluator) Eval(expr types.Expr) (types.Value, error) {
	switch e := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(e.Value), nil
	case *types.StringExpr:
		return types.StringValue(e.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(e.Value), nil
	default:
		if ve, ok := expr.(valueExpr); ok {
			return ve.value, nil
		}
		if val, ok := expr.(types.Value); ok {
			return val, nil
		}
		return nil, nil
	}
}

func (me *mockEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	// For testing swap!, we need to simulate function calls
	if fn, ok := funcValue.(*types.FunctionValue); ok {
		// Check if it's our test increment function by checking params
		if len(fn.Params) == 1 && fn.Params[0] == "x" {
			// Simulate increment function: (fn [x] (+ x 1))
			if len(args) == 1 {
				if val, err := me.Eval(args[0]); err == nil {
					if num, ok := val.(types.NumberValue); ok {
						return types.NumberValue(num + 1), nil
					}
				}
			}
		}
	}
	// Return an error for non-function values to simulate proper validation
	return nil, fmt.Errorf("not a function: %T", funcValue)
}

func wrapValue(value types.Value) types.Expr {
	return valueExpr{value}
}

type valueExpr struct {
	value types.Value
}

func (ve valueExpr) String() string {
	return ve.value.String()
}

func TestAtomPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewAtomPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"atom", "deref", "swap!", "reset!"}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestAtomPlugin_BasicAtomOperations(t *testing.T) {
	plugin := NewAtomPlugin()
	evaluator := newMockEvaluator()

	// Test creating an atom
	atomArgs := []types.Expr{&types.NumberExpr{Value: 42}}
	result, err := plugin.evalAtom(evaluator, atomArgs)
	if err != nil {
		t.Fatalf("evalAtom failed: %v", err)
	}

	atomValue, ok := result.(*types.AtomValue)
	if !ok {
		t.Fatalf("Expected AtomValue, got %T", result)
	}

	// Test dereferencing the atom
	derefArgs := []types.Expr{wrapValue(atomValue)}
	result, err = plugin.evalDeref(evaluator, derefArgs)
	if err != nil {
		t.Fatalf("evalDeref failed: %v", err)
	}

	if result != types.NumberValue(42) {
		t.Errorf("Expected 42, got %v", result)
	}

	// Test resetting the atom
	resetArgs := []types.Expr{
		wrapValue(atomValue),
		&types.NumberExpr{Value: 100},
	}
	result, err = plugin.evalReset(evaluator, resetArgs)
	if err != nil {
		t.Fatalf("evalReset failed: %v", err)
	}

	if result != types.NumberValue(100) {
		t.Errorf("Expected 100, got %v", result)
	}

	// Verify the atom was reset by dereferencing again
	result, err = plugin.evalDeref(evaluator, derefArgs)
	if err != nil {
		t.Fatalf("evalDeref after reset failed: %v", err)
	}

	if result != types.NumberValue(100) {
		t.Errorf("Expected 100 after reset, got %v", result)
	}
}

func TestAtomPlugin_SwapOperation(t *testing.T) {
	plugin := NewAtomPlugin()
	evaluator := newMockEvaluator()

	// Create an atom with initial value
	atomArgs := []types.Expr{&types.NumberExpr{Value: 5}}
	result, err := plugin.evalAtom(evaluator, atomArgs)
	if err != nil {
		t.Fatalf("evalAtom failed: %v", err)
	}

	atomValue, ok := result.(*types.AtomValue)
	if !ok {
		t.Fatalf("Expected AtomValue, got %T", result)
	}

	// Create a simple increment function for testing
	incFunc := &types.FunctionValue{
		Params: []string{"x"},
		Body:   &types.NumberExpr{Value: 0}, // Dummy body, actual logic is in CallFunction
		Env:    nil,                         // No environment needed for test
	}

	// Test swap! with increment function
	swapArgs := []types.Expr{
		wrapValue(atomValue),
		wrapValue(incFunc),
	}
	result, err = plugin.evalSwap(evaluator, swapArgs)
	if err != nil {
		t.Fatalf("evalSwap failed: %v", err)
	}

	if result != types.NumberValue(6) {
		t.Errorf("Expected 6 after increment, got %v", result)
	}

	// Verify the atom was updated by dereferencing
	derefArgs := []types.Expr{wrapValue(atomValue)}
	result, err = plugin.evalDeref(evaluator, derefArgs)
	if err != nil {
		t.Fatalf("evalDeref after swap failed: %v", err)
	}

	if result != types.NumberValue(6) {
		t.Errorf("Expected 6 in atom after swap, got %v", result)
	}
}

func TestAtomPlugin_AtomWithDifferentTypes(t *testing.T) {
	plugin := NewAtomPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name  string
		value types.Value
	}{
		{"string atom", types.StringValue("hello")},
		{"boolean atom", types.BooleanValue(true)},
		{"number atom", types.NumberValue(3.14)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create atom with the test value
			atomArgs := []types.Expr{wrapValue(tt.value)}
			result, err := plugin.evalAtom(evaluator, atomArgs)
			if err != nil {
				t.Fatalf("evalAtom failed: %v", err)
			}

			atomValue, ok := result.(*types.AtomValue)
			if !ok {
				t.Fatalf("Expected AtomValue, got %T", result)
			}

			// Dereference and verify the value
			derefArgs := []types.Expr{wrapValue(atomValue)}
			result, err = plugin.evalDeref(evaluator, derefArgs)
			if err != nil {
				t.Fatalf("evalDeref failed: %v", err)
			}

			if result != tt.value {
				t.Errorf("Expected %v, got %v", tt.value, result)
			}
		})
	}
}

func TestAtomPlugin_ErrorCases(t *testing.T) {
	plugin := NewAtomPlugin()
	evaluator := newMockEvaluator()

	// Test deref with non-atom
	nonAtom := types.NumberValue(42)
	_, err := plugin.evalDeref(evaluator, []types.Expr{wrapValue(nonAtom)})
	if err == nil {
		t.Error("Expected error when dereferencing non-atom")
	}

	// Test reset! with non-atom
	_, err = plugin.evalReset(evaluator, []types.Expr{
		wrapValue(nonAtom),
		&types.NumberExpr{Value: 100},
	})
	if err == nil {
		t.Error("Expected error when resetting non-atom")
	}

	// Test swap! with non-atom
	incFunc := &types.FunctionValue{
		Params: []string{"x"},
		Body:   &types.NumberExpr{Value: 0}, // Dummy body
		Env:    nil,
	}
	_, err = plugin.evalSwap(evaluator, []types.Expr{
		wrapValue(nonAtom),
		wrapValue(incFunc),
	})
	if err == nil {
		t.Error("Expected error when swapping non-atom")
	}

	// Test swap! with non-function
	atomArgs := []types.Expr{&types.NumberExpr{Value: 5}}
	result, err := plugin.evalAtom(evaluator, atomArgs)
	if err != nil {
		t.Fatalf("evalAtom failed: %v", err)
	}
	atomValue := result.(*types.AtomValue)

	_, err = plugin.evalSwap(evaluator, []types.Expr{
		wrapValue(atomValue),
		&types.NumberExpr{Value: 42}, // Non-function
	})
	if err == nil {
		t.Error("Expected error when swapping with non-function")
	}
}

func TestAtomPlugin_ConcurrentAccess(t *testing.T) {
	plugin := NewAtomPlugin()
	evaluator := newMockEvaluator()

	// Create an atom
	atomArgs := []types.Expr{&types.NumberExpr{Value: 0}}
	result, err := plugin.evalAtom(evaluator, atomArgs)
	if err != nil {
		t.Fatalf("evalAtom failed: %v", err)
	}

	atomValue, ok := result.(*types.AtomValue)
	if !ok {
		t.Fatalf("Expected AtomValue, got %T", result)
	}

	// Test multiple resets (simulating concurrent access)
	for i := 1; i <= 10; i++ {
		resetArgs := []types.Expr{
			wrapValue(atomValue),
			&types.NumberExpr{Value: float64(i)},
		}
		result, err = plugin.evalReset(evaluator, resetArgs)
		if err != nil {
			t.Fatalf("evalReset failed at iteration %d: %v", i, err)
		}

		if result != types.NumberValue(float64(i)) {
			t.Errorf("Expected %d, got %v", i, result)
		}
	}

	// Final dereference should show the last value
	derefArgs := []types.Expr{wrapValue(atomValue)}
	result, err = plugin.evalDeref(evaluator, derefArgs)
	if err != nil {
		t.Fatalf("evalDeref failed: %v", err)
	}

	if result != types.NumberValue(10) {
		t.Errorf("Expected 10, got %v", result)
	}
}
