package concurrency

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
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
	return nil, nil // Not used in concurrency tests
}

func (me *mockEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	// For testing purposes, just call regular Eval
	// In a real implementation, this would use the bindings
	return me.Eval(expr)
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

func (ve valueExpr) GetPosition() types.Position {
	return types.Position{Line: 1, Column: 1}
}

func TestConcurrencyPlugin_RegisterFunctions(t *testing.T) {
	mockEval := newMockEvaluator()
	plugin := NewConcurrencyPlugin(mockEval)
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{
		"go", "go-wait", "go-wait-all", "chan", "chan-send!", "chan-recv!",
		"chan-try-recv!", "chan-close!", "chan-closed?",
	}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestConcurrencyPlugin_GoFunction(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewConcurrencyPlugin(evaluator)

	// Test go function with a simple expression
	args := []types.Expr{&types.NumberExpr{Value: 42}}
	result, err := plugin.goFunc(evaluator, args)
	if err != nil {
		t.Fatalf("goFunc failed: %v", err)
	}

	// Result should be a future
	future, ok := result.(*types.FutureValue)
	if !ok {
		t.Fatalf("Expected future, got %T", result)
	}

	// Wait for the future to complete
	resultValue, err := future.Wait()
	if err != nil {
		t.Fatalf("Future failed: %v", err)
	}

	// Check the result
	if resultValue.String() != "42" {
		t.Errorf("Expected result 42, got %s", resultValue.String())
	}
}

func TestConcurrencyPlugin_GoWait(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewConcurrencyPlugin(evaluator)

	// Create a future first
	goArgs := []types.Expr{&types.NumberExpr{Value: 100}}
	futureResult, err := plugin.goFunc(evaluator, goArgs)
	if err != nil {
		t.Fatalf("goFunc failed: %v", err)
	}

	// Wait for the future
	waitArgs := []types.Expr{wrapValue(futureResult)}
	result, err := plugin.goWait(evaluator, waitArgs)
	if err != nil {
		t.Fatalf("goWait failed: %v", err)
	}

	// Result should be the value from the future
	if result.String() != "100" {
		t.Errorf("Expected result 100, got %s", result.String())
	}
}

func TestConcurrencyPlugin_ChannelCreation(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewConcurrencyPlugin(evaluator)

	// Create an unbuffered channel
	chanResult, err := plugin.chan_(evaluator, []types.Expr{})
	if err != nil {
		t.Fatalf("chan failed: %v", err)
	}

	_, ok := chanResult.(*types.ChannelValue)
	if !ok {
		t.Fatalf("Expected channel, got %T", chanResult)
	}

	// Create a buffered channel
	bufferedChanResult, err := plugin.chan_(evaluator, []types.Expr{&types.NumberExpr{Value: 5}})
	if err != nil {
		t.Fatalf("chan with buffer failed: %v", err)
	}

	_, ok = bufferedChanResult.(*types.ChannelValue)
	if !ok {
		t.Fatalf("Expected buffered channel, got %T", bufferedChanResult)
	}
}

func TestConcurrencyPlugin_ErrorCases(t *testing.T) {
	evaluator := newMockEvaluator()
	plugin := NewConcurrencyPlugin(evaluator)

	// Test go function with no arguments
	_, err := plugin.goFunc(evaluator, []types.Expr{})
	if err == nil {
		t.Error("Expected error for go function with no arguments")
	}

	// Test go-wait with wrong argument type
	_, err = plugin.goWait(evaluator, []types.Expr{&types.NumberExpr{Value: 42}})
	if err == nil {
		t.Error("Expected error for go-wait with non-future argument")
	}

	// Test chan with negative buffer size
	_, err = plugin.chan_(evaluator, []types.Expr{&types.NumberExpr{Value: -1}})
	if err == nil {
		t.Error("Expected error for negative buffer size")
	}

	// Test chan-send with wrong argument count
	_, err = plugin.chanSend(evaluator, []types.Expr{&types.NumberExpr{Value: 42}})
	if err == nil {
		t.Error("Expected error for chan-send with wrong argument count")
	}

	// Test chan-recv with wrong argument type
	_, err = plugin.chanRecv(evaluator, []types.Expr{&types.NumberExpr{Value: 42}})
	if err == nil {
		t.Error("Expected error for chan-recv with non-channel argument")
	}
}
