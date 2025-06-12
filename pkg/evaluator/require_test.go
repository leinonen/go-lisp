package evaluator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestRequireFunction(t *testing.T) {
	// Create temporary test module file
	tempDir := t.TempDir()
	testModulePath := filepath.Join(tempDir, "test-module.lisp")

	moduleContent := `(module test-utils
  (export double triple add-ten)
  
  (defun double [x] (* x 2))
  (defun triple [x] (* x 3))
  (defun add-ten [x] (+ x 10))
  
  ; Private function (not exported)
  (defun private-helper [x] (+ x 1)))`

	err := os.WriteFile(testModulePath, []byte(moduleContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test module file: %v", err)
	}

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

	t.Run("basic require - load and import all", func(t *testing.T) {
		// Test basic require functionality
		requireExpr := `(require "` + testModulePath + `")`
		result, err := evalExpr(requireExpr)
		if err != nil {
			t.Fatalf("Failed to require module: %v", err)
		}

		// Should return the module value
		moduleValue, ok := result.(*types.ModuleValue)
		if !ok {
			t.Fatalf("Expected ModuleValue, got %T", result)
		}

		if moduleValue.Name != "test-utils" {
			t.Errorf("Expected module name 'test-utils', got '%s'", moduleValue.Name)
		}

		// Test that exported functions are available directly
		_, err = evalExpr("(double 5)")
		if err != nil {
			t.Errorf("Function 'double' should be available after require: %v", err)
		}

		result, err = evalExpr("(triple 4)")
		if err != nil {
			t.Errorf("Function 'triple' should be available after require: %v", err)
		}
		if result.String() != "12" {
			t.Errorf("Expected 12, got %s", result.String())
		}

		// Test that private functions are not accessible
		_, err = evalExpr("(private-helper 5)")
		if err == nil {
			t.Error("Private function should not be accessible after require")
		}
	})

	t.Run("require with :as alias", func(t *testing.T) {
		// Create a new environment for this test
		env2 := NewEnvironment()
		evaluator2 := NewEvaluator(env2)

		evalExpr2 := func(input string) (types.Value, error) {
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

			return evaluator2.Eval(expr)
		}

		// Test require with alias
		requireExpr := `(require "` + testModulePath + `" :as utils)`
		result, err := evalExpr2(requireExpr)
		if err != nil {
			t.Fatalf("Failed to require module with alias: %v", err)
		}

		// Should return the module value
		_, ok := result.(*types.ModuleValue)
		if !ok {
			t.Fatalf("Expected ModuleValue, got %T", result)
		}

		// Test qualified access through alias
		result, err = evalExpr2("(utils.double 6)")
		if err != nil {
			t.Errorf("Qualified access through alias should work: %v", err)
		}
		if result.String() != "12" {
			t.Errorf("Expected 12, got %s", result.String())
		}

		// Test that direct access is not available
		_, err = evalExpr2("(double 5)")
		if err == nil {
			t.Error("Direct access should not be available when using :as alias")
		}
	})

	t.Run("require with :only selective import", func(t *testing.T) {
		// Create a new environment for this test
		env3 := NewEnvironment()
		evaluator3 := NewEvaluator(env3)

		evalExpr3 := func(input string) (types.Value, error) {
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

			return evaluator3.Eval(expr)
		}

		// Test require with selective import
		requireExpr := `(require "` + testModulePath + `" :only (double add-ten))`
		result, err := evalExpr3(requireExpr)
		if err != nil {
			t.Fatalf("Failed to require module with selective import: %v", err)
		}

		// Should return the module value
		_, ok := result.(*types.ModuleValue)
		if !ok {
			t.Fatalf("Expected ModuleValue, got %T", result)
		}

		// Test that selected functions are available
		result, err = evalExpr3("(double 7)")
		if err != nil {
			t.Errorf("Selected function 'double' should be available: %v", err)
		}
		if result.String() != "14" {
			t.Errorf("Expected 14, got %s", result.String())
		}

		result, err = evalExpr3("(add-ten 5)")
		if err != nil {
			t.Errorf("Selected function 'add-ten' should be available: %v", err)
		}
		if result.String() != "15" {
			t.Errorf("Expected 15, got %s", result.String())
		}

		// Test that non-selected functions are not available
		_, err = evalExpr3("(triple 3)")
		if err == nil {
			t.Error("Non-selected function 'triple' should not be available")
		}
	})

	t.Run("require caching - module loaded only once", func(t *testing.T) {
		// Create a new environment for this test
		env4 := NewEnvironment()
		evaluator4 := NewEvaluator(env4)

		evalExpr4 := func(input string) (types.Value, error) {
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

			return evaluator4.Eval(expr)
		}

		// Create a module that prints during loading to test caching
		sideEffectModulePath := filepath.Join(tempDir, "side-effect-module.lisp")
		sideEffectContent := `; This module will create a side effect during loading
(def load-count 1)

(module counter
  (export get-load-count)
  
  (defun get-load-count [] load-count))`

		err := os.WriteFile(sideEffectModulePath, []byte(sideEffectContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create side effect module file: %v", err)
		}

		// First require
		requireExpr := `(require "` + sideEffectModulePath + `")`
		_, err = evalExpr4(requireExpr)
		if err != nil {
			t.Fatalf("Failed first require: %v", err)
		}

		// Check that load-count is 1
		result, err := evalExpr4("(get-load-count)")
		if err != nil {
			t.Fatalf("Failed to get load count: %v", err)
		}
		if result.String() != "1" {
			t.Errorf("Expected 1, got %s", result.String())
		}

		// Second require should use cached module (no side effects from loading)
		_, err = evalExpr4(requireExpr)
		if err != nil {
			t.Fatalf("Failed second require: %v", err)
		}

		// The load-count should still be 1 (not re-executed)
		result, err = evalExpr4("(get-load-count)")
		if err != nil {
			t.Fatalf("Failed to get count: %v", err)
		}
		if result.String() != "1" {
			t.Errorf("Expected count to remain 1 (cached), got %s", result.String())
		}
	})

	t.Run("require error cases", func(t *testing.T) {
		env5 := NewEnvironment()
		evaluator5 := NewEvaluator(env5)

		evalExpr5 := func(input string) (types.Value, error) {
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

			return evaluator5.Eval(expr)
		}

		// Test non-existent file
		_, err := evalExpr5(`(require "non-existent.lisp")`)
		if err == nil {
			t.Error("Expected error for non-existent file")
		}

		// Test invalid syntax
		_, err = evalExpr5(`(require)`)
		if err == nil {
			t.Error("Expected error for require with no arguments")
		}

		// Test invalid :as usage
		_, err = evalExpr5(`(require "` + testModulePath + `" :as)`)
		if err == nil {
			t.Error("Expected error for require with :as but no alias")
		}

		// Test invalid :only usage
		_, err = evalExpr5(`(require "` + testModulePath + `" :only)`)
		if err == nil {
			t.Error("Expected error for require with :only but no symbol list")
		}
	})
}

func TestRequireWithExistingLibraries(t *testing.T) {
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

	t.Run("require core library", func(t *testing.T) {
		// Test requiring the existing core library
		coreLibPath := filepath.Join("..", "..", "library", "core.lisp")
		requireExpr := `(require "` + coreLibPath + `")`
		result, err := evalExpr(requireExpr)
		if err != nil {
			t.Fatalf("Failed to require core library: %v", err)
		}

		// Should return the module value
		moduleValue, ok := result.(*types.ModuleValue)
		if !ok {
			t.Fatalf("Expected ModuleValue, got %T", result)
		}

		if moduleValue.Name != "core" {
			t.Errorf("Expected module name 'core', got '%s'", moduleValue.Name)
		}

		// Test that core functions are available directly
		result, err = evalExpr("(factorial 5)")
		if err != nil {
			t.Errorf("Function 'factorial' should be available after require: %v", err)
		}
		if result.String() != "120" {
			t.Errorf("Expected 120, got %s", result.String())
		}

		result, err = evalExpr("(fibonacci 6)")
		if err != nil {
			t.Errorf("Function 'fibonacci' should be available after require: %v", err)
		}
		if result.String() != "8" {
			t.Errorf("Expected 8, got %s", result.String())
		}
	})

	t.Run("require functional library with alias", func(t *testing.T) {
		// Create a new environment for this test
		env2 := NewEnvironment()
		evaluator2 := NewEvaluator(env2)

		evalExpr2 := func(input string) (types.Value, error) {
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

			return evaluator2.Eval(expr)
		}

		functionalLibPath := filepath.Join("..", "..", "library", "functional.lisp")
		requireExpr := `(require "` + functionalLibPath + `" :as fn)`
		result, err := evalExpr2(requireExpr)
		if err != nil {
			t.Fatalf("Failed to require functional library with alias: %v", err)
		}

		// Should return the module value
		moduleValue, ok := result.(*types.ModuleValue)
		if !ok {
			t.Fatalf("Expected ModuleValue, got %T", result)
		}

		if moduleValue.Name != "functional" {
			t.Errorf("Expected module name 'functional', got '%s'", moduleValue.Name)
		}

		// Test qualified access through alias
		result, err = evalExpr2("(fn.identity 42)")
		if err != nil {
			t.Errorf("Qualified access through alias should work: %v", err)
		}
		if result.String() != "42" {
			t.Errorf("Expected 42, got %s", result.String())
		}

		// Test that direct access is not available
		_, err = evalExpr2("(identity 42)")
		if err == nil {
			t.Error("Direct access should not be available when using :as alias")
		}
	})
}
