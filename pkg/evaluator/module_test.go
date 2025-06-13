package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestModuleSystem(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple module definition",
			input: `(module math 
				(export add multiply) 
				(defn add [x y] (+ x y))
				(defn multiply [x y] (* x y)))`,
			expected: "#<module:math>",
		},
		{
			name:     "module import and usage",
			input:    `(import math)`,
			expected: "#<module:math>",
		},
		{
			name:     "use imported function",
			input:    `(add 3 4)`,
			expected: "7",
		},
		{
			name:     "qualified access",
			input:    `(math.multiply 3 4)`,
			expected: "12",
		},
	}

	// First create the module
	moduleInput := `(module math 
		(export add multiply) 
		(defn add [x y] (+ x y))
		(defn multiply [x y] (* x y)))`

	moduleExpr := parseString(t, moduleInput)
	_, err := evaluator.Eval(moduleExpr)
	if err != nil {
		t.Fatalf("Failed to create module: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "simple module definition" {
				// Skip - already created above
				return
			}

			expr := parseString(t, tt.input)
			result, err := evaluator.Eval(expr)
			if err != nil {
				t.Fatalf("Evaluation error: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestModuleExports(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create a module with some private and public functions
	moduleInput := `(module utils 
		(export double) 
		(defn double [x] (* x 2))
		(defn private-helper [x] (+ x 1)))`

	moduleExpr := parseString(t, moduleInput)
	_, err := evaluator.Eval(moduleExpr)
	if err != nil {
		t.Fatalf("Failed to create module: %v", err)
	}

	// Import the module
	importExpr := parseString(t, `(import utils)`)
	_, err = evaluator.Eval(importExpr)
	if err != nil {
		t.Fatalf("Failed to import module: %v", err)
	}

	// Test that exported function is accessible
	doubleExpr := parseString(t, `(double 5)`)
	result, err := evaluator.Eval(doubleExpr)
	if err != nil {
		t.Fatalf("Failed to call exported function: %v", err)
	}

	if result.String() != "10" {
		t.Errorf("Expected 10, got %s", result.String())
	}

	// Test that private function is not accessible
	privateExpr := parseString(t, `(private-helper 5)`)
	_, err = evaluator.Eval(privateExpr)
	if err == nil {
		t.Error("Expected error when calling private function, but got none")
	}
}

func TestQualifiedAccess(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create multiple modules
	mathModule := `(module math 
		(export add) 
		(defn add [x y] (+ x y)))`

	stringModule := `(module strings 
		(export concat) 
		(defn concat [x y] x))` // simplified concat for testing

	// Create both modules
	mathExpr := parseString(t, mathModule)
	_, err := evaluator.Eval(mathExpr)
	if err != nil {
		t.Fatalf("Failed to create math module: %v", err)
	}

	stringExpr := parseString(t, stringModule)
	_, err = evaluator.Eval(stringExpr)
	if err != nil {
		t.Fatalf("Failed to create strings module: %v", err)
	}

	// Test qualified access without importing
	qualifiedExpr := parseString(t, `(math.add 2 3)`)
	result, err := evaluator.Eval(qualifiedExpr)
	if err != nil {
		t.Fatalf("Failed qualified access: %v", err)
	}

	if result.String() != "5" {
		t.Errorf("Expected 5, got %s", result.String())
	}
}

func parseString(t *testing.T, input string) types.Expr {
	tok := tokenizer.NewTokenizer(input)
	tokens, err := tok.TokenizeWithError()
	if err != nil {
		t.Fatalf("Tokenization error: %v", err)
	}

	p := parser.NewParser(tokens)
	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	return expr
}
