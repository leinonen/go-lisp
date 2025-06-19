package tokenizer

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/types"
)

func TestTokenizeBraces(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []types.TokenType
	}{
		{
			name:     "empty braces",
			input:    "{}",
			expected: []types.TokenType{types.LBRACE, types.RBRACE},
		},
		{
			name:     "simple hash map",
			input:    "{:name \"John\"}",
			expected: []types.TokenType{types.LBRACE, types.KEYWORD, types.STRING, types.RBRACE},
		},
		{
			name:     "nested braces",
			input:    "{{}}",
			expected: []types.TokenType{types.LBRACE, types.LBRACE, types.RBRACE, types.RBRACE},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tt.input)
			tokens, err := tokenizer.TokenizeWithError()
			if err != nil {
				t.Fatalf("tokenization failed: %v", err)
			}

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, token := range tokens {
				if token.Type != tt.expected[i] {
					t.Errorf("token %d: expected %v, got %v", i, tt.expected[i], token.Type)
				}
			}
		})
	}
}
