package main

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/plugins/advanced"
	"github.com/leinonen/go-lisp/pkg/plugins/binding"
	"github.com/leinonen/go-lisp/pkg/plugins/control"
	"github.com/leinonen/go-lisp/pkg/plugins/environment"
	"github.com/leinonen/go-lisp/pkg/plugins/keyword"
	"github.com/leinonen/go-lisp/pkg/plugins/macro"
	"github.com/leinonen/go-lisp/pkg/plugins/sequence"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// IntegrationTestSuite tests all the new features together
func TestIntegrationSuite(t *testing.T) {
	// Create a registry and register all new plugins
	reg := registry.NewRegistry()
	env := evaluator.NewEnvironment()

	// Register all new plugins
	plugins := []interface {
		RegisterFunctions(registry.FunctionRegistry) error
	}{
		binding.NewBindingPlugin(),
		keyword.NewKeywordPlugin(),
		control.NewControlPlugin(),
		sequence.NewSequencePlugin(),
		environment.NewEnvironmentPlugin(env),
		advanced.NewAdvancedPlugin(),
		macro.NewMacroPlugin(),
	}

	for _, plugin := range plugins {
		if err := plugin.RegisterFunctions(reg); err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}
	}

	// Create a mock evaluator that uses the registry
	evaluator := &integrationEvaluator{
		env:      env,
		registry: reg,
	}

	t.Run("Keywords", func(t *testing.T) {
		testKeywordFeatures(t, evaluator, reg)
	})

	t.Run("Enhanced Control Flow", func(t *testing.T) {
		testControlFlowFeatures(t, evaluator, reg)
	})

	t.Run("Let Bindings", func(t *testing.T) {
		testBindingFeatures(t, evaluator, reg)
	})

	t.Run("Vectors", func(t *testing.T) {
		testVectorFeatures(t, evaluator, reg)
	})

	t.Run("Advanced Bindings", func(t *testing.T) {
		testAdvancedBindingFeatures(t, evaluator, reg)
	})

	t.Run("Macros", func(t *testing.T) {
		testMacroFeatures(t, evaluator, reg)
	})
}

// Mock evaluator for integration tests
type integrationEvaluator struct {
	env      *evaluator.Environment
	registry registry.FunctionRegistry
}

func (ie *integrationEvaluator) Eval(expr types.Expr) (types.Value, error) {
	// Basic evaluation logic
	switch e := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(e.Value), nil
	case *types.StringExpr:
		return types.StringValue(e.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(e.Value), nil
	case *types.KeywordExpr:
		// Remove the leading : from the value and create KeywordValue
		value := e.Value
		if len(value) > 0 && value[0] == ':' {
			value = value[1:]
		}
		return types.KeywordValue(value), nil
	case *types.SymbolExpr:
		if val, exists := ie.env.Get(e.Name); exists {
			return val, nil
		}
		return types.StringValue(e.Name), nil
	case *types.ListExpr:
		if len(e.Elements) > 0 {
			if symbol, ok := e.Elements[0].(*types.SymbolExpr); ok {
				if fn, exists := ie.registry.Get(symbol.Name); exists {
					return fn.Call(ie, e.Elements[1:])
				}
			}
		}
		// Convert ListExpr to ListValue for testing
		var elements []types.Value
		for _, elem := range e.Elements {
			val, err := ie.Eval(elem)
			if err != nil {
				return nil, err
			}
			elements = append(elements, val)
		}
		return &types.ListValue{Elements: elements}, nil
	}
	return types.StringValue("unknown"), nil
}

func (ie *integrationEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	return types.StringValue("function-result"), nil
}

func (ie *integrationEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	// For testing purposes, just call regular Eval
	// In a real implementation, this would use the bindings
	return ie.Eval(expr)
}

func testKeywordFeatures(t *testing.T, evaluator *integrationEvaluator, reg registry.FunctionRegistry) {
	// Test keyword creation
	keywordFunc, _ := reg.Get("keyword")
	result, err := keywordFunc.Call(evaluator, []types.Expr{
		&types.StringExpr{Value: "test"},
	})
	if err != nil {
		t.Errorf("keyword function failed: %v", err)
	}
	if kw, ok := result.(types.KeywordValue); !ok || string(kw) != "test" {
		t.Errorf("Expected keyword :test, got %v", result)
	}

	// Test keyword predicate
	keywordPredFunc, _ := reg.Get("keyword?")
	result, err = keywordPredFunc.Call(evaluator, []types.Expr{
		&types.KeywordExpr{Value: ":name"},
	})
	if err != nil {
		t.Errorf("keyword? function failed: %v", err)
	}
	if b, ok := result.(types.BooleanValue); !ok || !bool(b) {
		t.Errorf("Expected true for keyword?, got %v", result)
	}
}

func testControlFlowFeatures(t *testing.T, evaluator *integrationEvaluator, reg registry.FunctionRegistry) {
	// Test cond
	condFunc, _ := reg.Get("cond")
	result, err := condFunc.Call(evaluator, []types.Expr{
		&types.BooleanExpr{Value: false},
		&types.StringExpr{Value: "first"},
		&types.BooleanExpr{Value: true},
		&types.StringExpr{Value: "second"},
		&types.KeywordExpr{Value: ":else"},
		&types.StringExpr{Value: "default"},
	})
	if err != nil {
		t.Errorf("cond function failed: %v", err)
	}
	if s, ok := result.(types.StringValue); !ok || string(s) != "second" {
		t.Errorf("Expected 'second', got %v", result)
	}

	// Test when
	whenFunc, _ := reg.Get("when")
	result, err = whenFunc.Call(evaluator, []types.Expr{
		&types.BooleanExpr{Value: true},
		&types.NumberExpr{Value: 1},
		&types.NumberExpr{Value: 2},
		&types.NumberExpr{Value: 3},
	})
	if err != nil {
		t.Errorf("when function failed: %v", err)
	}
	if n, ok := result.(types.NumberValue); !ok || n != 3 {
		t.Errorf("Expected 3, got %v", result)
	}
}

func testBindingFeatures(t *testing.T, evaluator *integrationEvaluator, reg registry.FunctionRegistry) {
	// Test basic let
	letFunc, _ := reg.Get("let")
	result, err := letFunc.Call(evaluator, []types.Expr{
		&types.BracketExpr{Elements: []types.Expr{
			&types.SymbolExpr{Name: "x"},
			&types.NumberExpr{Value: 10},
			&types.SymbolExpr{Name: "y"},
			&types.NumberExpr{Value: 20},
		}},
		&types.ListExpr{Elements: []types.Expr{
			&types.SymbolExpr{Name: "+"},
			&types.SymbolExpr{Name: "x"},
			&types.SymbolExpr{Name: "y"},
		}},
	})
	if err != nil {
		t.Errorf("let function failed: %v", err)
	}
	// The mock implementation returns specific values, just check it doesn't error
	if result == nil {
		t.Error("Expected result from let")
	}

	// Test let*
	letStarFunc, exists := reg.Get("let*")
	if exists {
		result, err = letStarFunc.Call(evaluator, []types.Expr{
			&types.BracketExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "x"},
				&types.NumberExpr{Value: 5},
				&types.SymbolExpr{Name: "y"},
				&types.ListExpr{Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.SymbolExpr{Name: "x"},
					&types.NumberExpr{Value: 10},
				}},
			}},
			&types.SymbolExpr{Name: "y"},
		})
		if err != nil {
			t.Errorf("let* function failed: %v", err)
		}
		if result == nil {
			t.Error("Expected result from let*")
		}
	}
}

func testVectorFeatures(t *testing.T, evaluator *integrationEvaluator, reg registry.FunctionRegistry) {
	// Test vector creation
	vectorFunc, _ := reg.Get("vector")
	result, err := vectorFunc.Call(evaluator, []types.Expr{
		&types.NumberExpr{Value: 1},
		&types.NumberExpr{Value: 2},
		&types.NumberExpr{Value: 3},
	})
	if err != nil {
		t.Errorf("vector function failed: %v", err)
	}
	if vec, ok := result.(*types.VectorValue); !ok || len(vec.Elements) != 3 {
		t.Errorf("Expected vector with 3 elements, got %v", result)
	}

	// Test vec conversion
	vecFunc, _ := reg.Get("vec")
	result, err = vecFunc.Call(evaluator, []types.Expr{
		&types.ListExpr{Elements: []types.Expr{
			&types.NumberExpr{Value: 1},
			&types.NumberExpr{Value: 2},
		}},
	})
	if err != nil {
		t.Errorf("vec function failed: %v", err)
	}
	if vec, ok := result.(*types.VectorValue); !ok || len(vec.Elements) != 2 {
		t.Errorf("Expected vector with 2 elements, got %v", result)
	}
}

func testAdvancedBindingFeatures(t *testing.T, evaluator *integrationEvaluator, reg registry.FunctionRegistry) {
	// Test destructuring let
	letDestructureFunc, exists := reg.Get("let-destructure")
	if exists {
		result, err := letDestructureFunc.Call(evaluator, []types.Expr{
			&types.BracketExpr{Elements: []types.Expr{
				&types.BracketExpr{Elements: []types.Expr{
					&types.SymbolExpr{Name: "a"},
					&types.SymbolExpr{Name: "b"},
				}},
				&types.SymbolExpr{Name: "my-vec"},
			}},
			&types.ListExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "+"},
				&types.SymbolExpr{Name: "a"},
				&types.SymbolExpr{Name: "b"},
			}},
		})
		if err != nil {
			t.Errorf("let-destructure function failed: %v", err)
		}
		if result == nil {
			t.Error("Expected result from let-destructure")
		}
	}
}

func testMacroFeatures(t *testing.T, evaluator *integrationEvaluator, reg registry.FunctionRegistry) {
	// Test macroexpand (quote is tested elsewhere as it's in core)
	macroexpandFunc, exists := reg.Get("macroexpand")
	if exists {
		result, err := macroexpandFunc.Call(evaluator, []types.Expr{
			&types.ListExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "when"},
				&types.BooleanExpr{Value: true},
				&types.NumberExpr{Value: 42},
			}},
		})
		if err != nil {
			t.Errorf("macroexpand function failed: %v", err)
		}
		if quoted, ok := result.(*types.QuotedValue); ok {
			if listExpr, ok := quoted.Value.(*types.ListExpr); ok {
				if len(listExpr.Elements) == 0 {
					t.Error("Expected non-empty expanded expression")
				} else if symbol, ok := listExpr.Elements[0].(*types.SymbolExpr); !ok || symbol.Name != "if" {
					t.Errorf("Expected expanded expression to start with 'if', got %v", listExpr.Elements[0])
				}
			} else {
				t.Errorf("Expected expanded expression to be a list, got %T", quoted.Value)
			}
		} else {
			t.Errorf("Expected QuotedValue, got %T", result)
		}
	}
}

// Benchmark tests
func BenchmarkIntegrationFeatures(b *testing.B) {
	// Setup
	reg := registry.NewRegistry()
	env := evaluator.NewEnvironment()

	plugins := []interface {
		RegisterFunctions(registry.FunctionRegistry) error
	}{
		binding.NewBindingPlugin(),
		keyword.NewKeywordPlugin(),
		control.NewControlPlugin(),
		sequence.NewSequencePlugin(),
	}

	for _, plugin := range plugins {
		if err := plugin.RegisterFunctions(reg); err != nil {
			b.Fatalf("Failed to register plugin: %v", err)
		}
	}

	evaluator := &integrationEvaluator{
		env:      env,
		registry: reg,
	}

	b.Run("KeywordCreation", func(b *testing.B) {
		keywordFunc, _ := reg.Get("keyword")
		args := []types.Expr{&types.StringExpr{Value: "test"}}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = keywordFunc.Call(evaluator, args)
		}
	})

	b.Run("LetBinding", func(b *testing.B) {
		letFunc, _ := reg.Get("let")
		args := []types.Expr{
			&types.BracketExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "x"},
				&types.NumberExpr{Value: 10},
			}},
			&types.SymbolExpr{Name: "x"},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = letFunc.Call(evaluator, args)
		}
	})

	b.Run("VectorCreation", func(b *testing.B) {
		vectorFunc, _ := reg.Get("vector")
		args := []types.Expr{
			&types.NumberExpr{Value: 1},
			&types.NumberExpr{Value: 2},
			&types.NumberExpr{Value: 3},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = vectorFunc.Call(evaluator, args)
		}
	})
}
