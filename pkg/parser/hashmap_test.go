package parser

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/tokenizer"
)

func TestParseHashMapLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty hash map",
			input:    "{}",
			expected: "HashMapExpr([])",
		},
		{
			name:     "simple hash map",
			input:    "{:name \"John\"}",
			expected: "HashMapExpr([KeywordExpr(:name) StringExpr(\"John\")])",
		},
		{
			name:     "hash map with multiple pairs",
			input:    "{:name \"John\" :age 30}",
			expected: "HashMapExpr([KeywordExpr(:name) StringExpr(\"John\") KeywordExpr(:age) NumberExpr(30)])",
		},
		{
			name:     "hash map with string keys",
			input:    "{\"name\" \"Alice\" \"age\" 25}",
			expected: "HashMapExpr([StringExpr(\"name\") StringExpr(\"Alice\") StringExpr(\"age\") NumberExpr(25)])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := tokenizer.NewTokenizer(tt.input)
			tokens, err := tokenizer.TokenizeWithError()
			if err != nil {
				t.Fatalf("tokenization failed: %v", err)
			}

			parser := NewParser(tokens)
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("parsing failed: %v", err)
			}

			if expr.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}
		})
	}
}

func TestParseHashMapLiteralErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "odd number of elements",
			input: "{:name \"John\" :age}",
		},
		{
			name:  "unmatched opening brace",
			input: "{:name \"John\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := tokenizer.NewTokenizer(tt.input)
			tokens, err := tokenizer.TokenizeWithError()
			if err != nil {
				t.Fatalf("tokenization failed: %v", err)
			}

			parser := NewParser(tokens)
			_, err = parser.Parse()
			if err == nil {
				t.Errorf("expected parsing to fail for input %s", tt.input)
			}
		})
	}
}
