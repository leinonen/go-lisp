package tokenizer

import (
	"reflect"
	"testing"

	"github.com/leinonen/go-lisp/pkg/types"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []types.Token
	}{
		{
			name:     "empty input",
			input:    "",
			expected: []types.Token{},
		},
		{
			name:  "single number",
			input: "42",
			expected: []types.Token{
				{Type: types.NUMBER, Value: "42", Position: types.Position{Line: 1, Column: 1}},
			},
		},
		{
			name:  "single symbol",
			input: "hello",
			expected: []types.Token{
				{Type: types.SYMBOL, Value: "hello", Position: types.Position{Line: 1, Column: 1}},
			},
		},
		{
			name:  "simple expression",
			input: "(+ 1 2)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.SYMBOL, Value: "+", Position: types.Position{Line: 1, Column: 2}},
				{Type: types.NUMBER, Value: "1", Position: types.Position{Line: 1, Column: 4}},
				{Type: types.NUMBER, Value: "2", Position: types.Position{Line: 1, Column: 6}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 7}},
			},
		},
		{
			name:  "nested expression",
			input: "(+ (* 2 3) 4)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.SYMBOL, Value: "+", Position: types.Position{Line: 1, Column: 2}},
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 4}},
				{Type: types.SYMBOL, Value: "*", Position: types.Position{Line: 1, Column: 5}},
				{Type: types.NUMBER, Value: "2", Position: types.Position{Line: 1, Column: 7}},
				{Type: types.NUMBER, Value: "3", Position: types.Position{Line: 1, Column: 9}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 10}},
				{Type: types.NUMBER, Value: "4", Position: types.Position{Line: 1, Column: 12}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 13}},
			},
		},
		{
			name:  "string literal",
			input: `"hello world"`,
			expected: []types.Token{
				{Type: types.STRING, Value: "hello world", Position: types.Position{Line: 1, Column: 1}},
			},
		},
		{
			name:  "boolean values",
			input: "true false",
			expected: []types.Token{
				{Type: types.BOOLEAN, Value: "true", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.BOOLEAN, Value: "false", Position: types.Position{Line: 1, Column: 6}},
			},
		},
		{
			name:  "whitespace handling",
			input: "  (  +   1    2  )  ",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 3}},
				{Type: types.SYMBOL, Value: "+", Position: types.Position{Line: 1, Column: 6}},
				{Type: types.NUMBER, Value: "1", Position: types.Position{Line: 1, Column: 10}},
				{Type: types.NUMBER, Value: "2", Position: types.Position{Line: 1, Column: 15}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 18}},
			},
		},
		{
			name:  "negative numbers",
			input: "(-42 -3.14)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.NUMBER, Value: "-42", Position: types.Position{Line: 1, Column: 2}},
				{Type: types.NUMBER, Value: "-3.14", Position: types.Position{Line: 1, Column: 6}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 11}},
			},
		},
		{
			name:  "comments ignored",
			input: "; this is a comment\n(+ 1 2)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 2, Column: 1}},
				{Type: types.SYMBOL, Value: "+", Position: types.Position{Line: 2, Column: 2}},
				{Type: types.NUMBER, Value: "1", Position: types.Position{Line: 2, Column: 4}},
				{Type: types.NUMBER, Value: "2", Position: types.Position{Line: 2, Column: 6}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 2, Column: 7}},
			},
		},
		{
			name:  "inline comments",
			input: "(+ 1 2) ; this is an inline comment",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.SYMBOL, Value: "+", Position: types.Position{Line: 1, Column: 2}},
				{Type: types.NUMBER, Value: "1", Position: types.Position{Line: 1, Column: 4}},
				{Type: types.NUMBER, Value: "2", Position: types.Position{Line: 1, Column: 6}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 7}},
			},
		},
		{
			name:  "multiple comments",
			input: "; comment 1\n(+ 1 2) ; comment 2\n; comment 3",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 2, Column: 1}},
				{Type: types.SYMBOL, Value: "+", Position: types.Position{Line: 2, Column: 2}},
				{Type: types.NUMBER, Value: "1", Position: types.Position{Line: 2, Column: 4}},
				{Type: types.NUMBER, Value: "2", Position: types.Position{Line: 2, Column: 6}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 2, Column: 7}},
			},
		},
		{
			name:  "comment with special characters",
			input: "; comment with ()#\"symbols\n42",
			expected: []types.Token{
				{Type: types.NUMBER, Value: "42", Position: types.Position{Line: 2, Column: 1}},
			},
		},
		{
			name:  "square brackets",
			input: "[x y]",
			expected: []types.Token{
				{Type: types.LBRACKET, Value: "[", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.SYMBOL, Value: "x", Position: types.Position{Line: 1, Column: 2}},
				{Type: types.SYMBOL, Value: "y", Position: types.Position{Line: 1, Column: 4}},
				{Type: types.RBRACKET, Value: "]", Position: types.Position{Line: 1, Column: 5}},
			},
		},
		{
			name:  "defn with square brackets",
			input: "(defn square [x] (* x x))",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.SYMBOL, Value: "defn", Position: types.Position{Line: 1, Column: 2}},
				{Type: types.SYMBOL, Value: "square", Position: types.Position{Line: 1, Column: 7}},
				{Type: types.LBRACKET, Value: "[", Position: types.Position{Line: 1, Column: 14}},
				{Type: types.SYMBOL, Value: "x", Position: types.Position{Line: 1, Column: 15}},
				{Type: types.RBRACKET, Value: "]", Position: types.Position{Line: 1, Column: 16}},
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 18}},
				{Type: types.SYMBOL, Value: "*", Position: types.Position{Line: 1, Column: 19}},
				{Type: types.SYMBOL, Value: "x", Position: types.Position{Line: 1, Column: 21}},
				{Type: types.SYMBOL, Value: "x", Position: types.Position{Line: 1, Column: 23}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 24}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 25}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tt.input)
			result := tokenizer.Tokenize()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTokenizerError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unterminated string",
			input: `"hello`,
		},
		{
			name:  "invalid character",
			input: "@",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tt.input)
			_, err := tokenizer.TokenizeWithError()

			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
		})
	}
}
