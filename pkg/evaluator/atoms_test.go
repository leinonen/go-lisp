package evaluator

import (
	"sync"
	"testing"
	"time"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Helper function to evaluate string expressions for testing
func evalExprStr(evaluator *Evaluator, input string) (types.Value, error) {
	tok := tokenizer.NewTokenizer(input)
	tokens, err := tok.TokenizeWithError()
	if err != nil {
		return nil, err
	}

	parser := parser.NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return evaluator.Eval(expr)
}

func TestAtomCreation(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name        string
		input       string
		expectValue types.Value
		expectError bool
	}{
		{
			name:        "create atom with number",
			input:       "(atom 42)",
			expectValue: types.NumberValue(42),
			expectError: false,
		},
		{
			name:        "create atom with string",
			input:       "(atom \"hello\")",
			expectValue: types.StringValue("hello"),
			expectError: false,
		},
		{
			name:        "create atom with boolean",
			input:       "(atom true)",
			expectValue: types.BooleanValue(true),
			expectError: false,
		},
		{
			name:        "create atom with list",
			input:       "(atom (list 1 2 3))",
			expectValue: &types.ListValue{Elements: []types.Value{types.NumberValue(1), types.NumberValue(2), types.NumberValue(3)}},
			expectError: false,
		},
		{
			name:        "atom with wrong number of arguments - none",
			input:       "(atom)",
			expectError: true,
		},
		{
			name:        "atom with wrong number of arguments - too many",
			input:       "(atom 1 2)",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := evalExprStr(evaluator, test.input)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check that result is an AtomValue
			atom, ok := result.(*types.AtomValue)
			if !ok {
				t.Fatalf("expected AtomValue, got %T", result)
			}

			// Check the initial value of the atom
			actualValue := atom.Value()
			if !valuesEqual(actualValue, test.expectValue) {
				t.Errorf("expected atom value %v, got %v", test.expectValue, actualValue)
			}

			// Verify atom string representation contains expected value
			atomStr := atom.String()
			if atomStr == "" {
				t.Errorf("atom string representation should not be empty")
			}
		})
	}
}

func TestDeref(t *testing.T) {
	tests := []struct {
		name        string
		setup       string
		input       string
		expectValue types.Value
		expectError bool
	}{
		{
			name:        "deref number atom",
			setup:       "(def a (atom 42))",
			input:       "(deref a)",
			expectValue: types.NumberValue(42),
			expectError: false,
		},
		{
			name:        "deref string atom",
			setup:       "(def a (atom \"hello\"))",
			input:       "(deref a)",
			expectValue: types.StringValue("hello"),
			expectError: false,
		},
		{
			name:        "deref boolean atom",
			setup:       "(def a (atom false))",
			input:       "(deref a)",
			expectValue: types.BooleanValue(false),
			expectError: false,
		},
		{
			name:        "deref list atom",
			setup:       "(def a (atom (list 1 2 3)))",
			input:       "(deref a)",
			expectValue: &types.ListValue{Elements: []types.Value{types.NumberValue(1), types.NumberValue(2), types.NumberValue(3)}},
			expectError: false,
		},
		{
			name:        "deref non-atom value",
			setup:       "(def a 42)",
			input:       "(deref a)",
			expectError: true,
		},
		{
			name:        "deref with wrong number of arguments - none",
			input:       "(deref)",
			expectError: true,
		},
		{
			name:        "deref with wrong number of arguments - too many",
			setup:       "(def a (atom 42))",
			input:       "(deref a 1)",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Fresh environment for each test
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			if test.setup != "" {
				_, err := evalExprStr(evaluator, test.setup)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			result, err := evalExprStr(evaluator, test.input)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, test.expectValue) {
				t.Errorf("expected %v, got %v", test.expectValue, result)
			}
		})
	}
}

func TestSwap(t *testing.T) {
	tests := []struct {
		name        string
		setup       string
		input       string
		expectValue types.Value
		expectError bool
	}{
		{
			name:        "swap with increment function",
			setup:       "(def a (atom 42))",
			input:       "(swap! a (fn [x] (+ x 1)))",
			expectValue: types.NumberValue(43),
			expectError: false,
		},
		{
			name:        "swap with double function",
			setup:       "(def a (atom 10))",
			input:       "(swap! a (fn [x] (* x 2)))",
			expectValue: types.NumberValue(20),
			expectError: false,
		},
		{
			name:        "swap with string concat",
			setup:       "(def a (atom \"hello\"))",
			input:       "(swap! a (fn [x] (string-concat x \" world\")))",
			expectValue: types.StringValue("hello world"),
			expectError: false,
		},
		{
			name:        "swap with list append",
			setup:       "(def a (atom (list 1 2)))",
			input:       "(swap! a (fn [x] (cons 3 x)))",
			expectValue: &types.ListValue{Elements: []types.Value{types.NumberValue(3), types.NumberValue(1), types.NumberValue(2)}},
			expectError: false,
		},
		{
			name:        "swap non-atom value",
			setup:       "(def a 42)",
			input:       "(swap! a (fn [x] (+ x 1)))",
			expectError: true,
		},
		{
			name:        "swap with non-function",
			setup:       "(def a (atom 42))",
			input:       "(swap! a 1)",
			expectError: true,
		},
		{
			name:        "swap with wrong number of arguments - not enough",
			setup:       "(def a (atom 42))",
			input:       "(swap! a)",
			expectError: true,
		},
		{
			name:        "swap with wrong number of arguments - too many",
			setup:       "(def a (atom 42))",
			input:       "(swap! a (fn [x] (+ x 1)) 1)",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Fresh environment for each test
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			if test.setup != "" {
				_, err := evalExprStr(evaluator, test.setup)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			result, err := evalExprStr(evaluator, test.input)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, test.expectValue) {
				t.Errorf("expected %v, got %v", test.expectValue, result)
			}
		})
	}
}

func TestReset(t *testing.T) {
	tests := []struct {
		name        string
		setup       string
		input       string
		expectValue types.Value
		expectError bool
	}{
		{
			name:        "reset number atom",
			setup:       "(def a (atom 42))",
			input:       "(reset! a 100)",
			expectValue: types.NumberValue(100),
			expectError: false,
		},
		{
			name:        "reset string atom",
			setup:       "(def a (atom \"hello\"))",
			input:       "(reset! a \"goodbye\")",
			expectValue: types.StringValue("goodbye"),
			expectError: false,
		},
		{
			name:        "reset boolean atom",
			setup:       "(def a (atom true))",
			input:       "(reset! a false)",
			expectValue: types.BooleanValue(false),
			expectError: false,
		},
		{
			name:        "reset list atom",
			setup:       "(def a (atom (list 1 2)))",
			input:       "(reset! a (list 3 4 5))",
			expectValue: &types.ListValue{Elements: []types.Value{types.NumberValue(3), types.NumberValue(4), types.NumberValue(5)}},
			expectError: false,
		},
		{
			name:        "reset non-atom value",
			setup:       "(def a 42)",
			input:       "(reset! a 100)",
			expectError: true,
		},
		{
			name:        "reset with wrong number of arguments - not enough",
			setup:       "(def a (atom 42))",
			input:       "(reset! a)",
			expectError: true,
		},
		{
			name:        "reset with wrong number of arguments - too many",
			setup:       "(def a (atom 42))",
			input:       "(reset! a 100 200)",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Fresh environment for each test
			env := NewEnvironment()
			evaluator := NewEvaluator(env)

			if test.setup != "" {
				_, err := evalExprStr(evaluator, test.setup)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			result, err := evalExprStr(evaluator, test.input)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, test.expectValue) {
				t.Errorf("expected %v, got %v", test.expectValue, result)
			}
		})
	}
}

func TestAtomStateConsistency(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Setup an atom and verify state changes persist
	_, err := evalExprStr(evaluator, "(def counter (atom 0))")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Initial value should be 0
	result, err := evalExprStr(evaluator, "(deref counter)")
	if err != nil {
		t.Fatalf("deref failed: %v", err)
	}
	if !valuesEqual(result, types.NumberValue(0)) {
		t.Errorf("expected 0, got %v", result)
	}

	// Increment with swap!
	result, err = evalExprStr(evaluator, "(swap! counter (fn [x] (+ x 1)))")
	if err != nil {
		t.Fatalf("swap failed: %v", err)
	}
	if !valuesEqual(result, types.NumberValue(1)) {
		t.Errorf("expected 1, got %v", result)
	}

	// Verify the value persisted
	result, err = evalExprStr(evaluator, "(deref counter)")
	if err != nil {
		t.Fatalf("deref failed: %v", err)
	}
	if !valuesEqual(result, types.NumberValue(1)) {
		t.Errorf("expected 1, got %v", result)
	}

	// Reset to a new value
	result, err = evalExprStr(evaluator, "(reset! counter 100)")
	if err != nil {
		t.Fatalf("reset failed: %v", err)
	}
	if !valuesEqual(result, types.NumberValue(100)) {
		t.Errorf("expected 100, got %v", result)
	}

	// Verify the reset persisted
	result, err = evalExprStr(evaluator, "(deref counter)")
	if err != nil {
		t.Fatalf("deref failed: %v", err)
	}
	if !valuesEqual(result, types.NumberValue(100)) {
		t.Errorf("expected 100, got %v", result)
	}
}

func TestAtomConcurrency(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Setup atom with initial value
	_, err := evalExprStr(evaluator, "(def counter (atom 0))")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Get the atom
	atomValue, ok := env.Get("counter")
	if !ok {
		t.Fatalf("counter atom not found")
	}
	atom, ok := atomValue.(*types.AtomValue)
	if !ok {
		t.Fatalf("expected AtomValue, got %T", atomValue)
	}

	// Test concurrent access
	const numGoroutines = 50
	const incrementsPerGoroutine = 5

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch goroutines that increment the counter
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				// Create a fresh environment and evaluator for this goroutine
				localEnv := NewEnvironment()
				localEvaluator := NewEvaluator(localEnv)
				localEnv.Set("counter", atom)

				// Use swap! to safely increment
				_, err := evalExprStr(localEvaluator, "(swap! counter (fn [x] (+ x 1)))")
				if err != nil {
					t.Errorf("swap failed: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	// Check final value
	finalValue := atom.Value()
	expectedValue := types.NumberValue(numGoroutines * incrementsPerGoroutine)

	if !valuesEqual(finalValue, expectedValue) {
		t.Errorf("expected final value %v, got %v", expectedValue, finalValue)
	}
}

func TestAtomConcurrentReadsAndWrites(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Setup atom
	_, err := evalExprStr(evaluator, "(def shared (atom 0))")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Get the atom
	atomValue, ok := env.Get("shared")
	if !ok {
		t.Fatalf("shared atom not found")
	}
	atom, ok := atomValue.(*types.AtomValue)
	if !ok {
		t.Fatalf("expected AtomValue, got %T", atomValue)
	}

	const numReaders = 10
	const numWriters = 3
	const operationsPerGoroutine = 5

	var wg sync.WaitGroup
	wg.Add(numReaders + numWriters)

	// Launch readers
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				localEnv := NewEnvironment()
				localEvaluator := NewEvaluator(localEnv)
				localEnv.Set("shared", atom)

				_, err := evalExprStr(localEvaluator, "(deref shared)")
				if err != nil {
					t.Errorf("deref failed: %v", err)
					return
				}
				time.Sleep(time.Microsecond) // Small delay to increase chance of race conditions
			}
		}()
	}

	// Launch writers
	for i := 0; i < numWriters; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				localEnv := NewEnvironment()
				localEvaluator := NewEvaluator(localEnv)
				localEnv.Set("shared", atom)

				_, err := evalExprStr(localEvaluator, "(swap! shared (fn [x] (+ x 1)))")
				if err != nil {
					t.Errorf("swap failed: %v", err)
					return
				}
				time.Sleep(time.Microsecond) // Small delay to increase chance of race conditions
			}
		}()
	}

	wg.Wait()

	// Verify final value
	result, err := evalExprStr(evaluator, "(deref shared)")
	if err != nil {
		t.Fatalf("final deref failed: %v", err)
	}

	expectedValue := types.NumberValue(numWriters * operationsPerGoroutine)
	if !valuesEqual(result, expectedValue) {
		t.Errorf("expected final value %v, got %v", expectedValue, result)
	}
}
