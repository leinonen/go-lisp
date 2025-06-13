package tokenizer

import (
	"reflect"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
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
				{Type: types.NUMBER, Value: "42"},
			},
		},
		{
			name:  "single symbol",
			input: "hello",
			expected: []types.Token{
				{Type: types.SYMBOL, Value: "hello"},
			},
		},
		{
			name:  "simple expression",
			input: "(+ 1 2)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "+"},
				{Type: types.NUMBER, Value: "1"},
				{Type: types.NUMBER, Value: "2"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:  "nested expression",
			input: "(+ (* 2 3) 4)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "+"},
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "*"},
				{Type: types.NUMBER, Value: "2"},
				{Type: types.NUMBER, Value: "3"},
				{Type: types.RPAREN, Value: ")"},
				{Type: types.NUMBER, Value: "4"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:  "string literal",
			input: `"hello world"`,
			expected: []types.Token{
				{Type: types.STRING, Value: "hello world"},
			},
		},
		{
			name:  "boolean values",
			input: "true false",
			expected: []types.Token{
				{Type: types.BOOLEAN, Value: "true"},
				{Type: types.BOOLEAN, Value: "false"},
			},
		},
		{
			name:  "whitespace handling",
			input: "  (  +   1    2  )  ",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "+"},
				{Type: types.NUMBER, Value: "1"},
				{Type: types.NUMBER, Value: "2"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:  "negative numbers",
			input: "(-42 -3.14)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.NUMBER, Value: "-42"},
				{Type: types.NUMBER, Value: "-3.14"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:  "comments ignored",
			input: "; this is a comment\n(+ 1 2)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "+"},
				{Type: types.NUMBER, Value: "1"},
				{Type: types.NUMBER, Value: "2"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:  "inline comments",
			input: "(+ 1 2) ; this is an inline comment",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "+"},
				{Type: types.NUMBER, Value: "1"},
				{Type: types.NUMBER, Value: "2"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:  "multiple comments",
			input: "; comment 1\n(+ 1 2) ; comment 2\n; comment 3",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "+"},
				{Type: types.NUMBER, Value: "1"},
				{Type: types.NUMBER, Value: "2"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:  "comment with special characters",
			input: "; comment with ()#\"symbols\n42",
			expected: []types.Token{
				{Type: types.NUMBER, Value: "42"},
			},
		},
		{
			name:  "square brackets",
			input: "[x y]",
			expected: []types.Token{
				{Type: types.LBRACKET, Value: "["},
				{Type: types.SYMBOL, Value: "x"},
				{Type: types.SYMBOL, Value: "y"},
				{Type: types.RBRACKET, Value: "]"},
			},
		},
		{
			name:  "defn with square brackets",
			input: "(defn square [x] (* x x))",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "defn"},
				{Type: types.SYMBOL, Value: "square"},
				{Type: types.LBRACKET, Value: "["},
				{Type: types.SYMBOL, Value: "x"},
				{Type: types.RBRACKET, Value: "]"},
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "*"},
				{Type: types.SYMBOL, Value: "x"},
				{Type: types.SYMBOL, Value: "x"},
				{Type: types.RPAREN, Value: ")"},
				{Type: types.RPAREN, Value: ")"},
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
