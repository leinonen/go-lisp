package pure

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/parser"
	"github.com/leinonen/go-lisp/pkg/tokenizer"
)

func TestEvalHashMapLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty hash map",
			input:    "{}",
			expected: "{}",
		},
		{
			name:     "simple hash map with keyword keys",
			input:    "{:name \"John\" :age 30}",
			expected: "{\"age\" 30 \"name\" John}",
		},
		{
			name:     "hash map with string keys",
			input:    "{\"name\" \"Alice\" \"job\" \"Engineer\"}",
			expected: "{\"job\" Engineer \"name\" Alice}",
		},
		{
			name:     "mixed key types",
			input:    "{:name \"Bob\" \"age\" 35}",
			expected: "{\"age\" 35 \"name\" Bob}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create environment and evaluator
			env := evaluator.NewEnvironment()
			eval, err := NewPureEvaluator(env)
			if err != nil {
				t.Fatalf("failed to create evaluator: %v", err)
			}

			// Tokenize, parse, and evaluate
			tok := tokenizer.NewTokenizer(tt.input)
			tokens, err := tok.TokenizeWithError()
			if err != nil {
				t.Fatalf("tokenization failed: %v", err)
			}

			parser := parser.NewParser(tokens)
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("parsing failed: %v", err)
			}

			result, err := eval.Eval(expr)
			if err != nil {
				t.Fatalf("evaluation failed: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestEvalHashMapLiteralErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid key type",
			input: "{123 \"value\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create environment and evaluator
			env := evaluator.NewEnvironment()
			eval, err := NewPureEvaluator(env)
			if err != nil {
				t.Fatalf("failed to create evaluator: %v", err)
			}

			// Tokenize, parse, and evaluate
			tok := tokenizer.NewTokenizer(tt.input)
			tokens, err := tok.TokenizeWithError()
			if err != nil {
				t.Fatalf("tokenization failed: %v", err)
			}

			parser := parser.NewParser(tokens)
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("parsing failed: %v", err)
			}

			_, err = eval.Eval(expr)
			if err == nil {
				t.Errorf("expected evaluation to fail for input %s", tt.input)
			}
		})
	}
}
