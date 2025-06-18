package tokenizer

import (
	"reflect"
	"testing"

	"github.com/leinonen/go-lisp/pkg/types"
)

func TestKeywordTokenization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []types.Token
	}{
		{
			name:  "simple keyword",
			input: ":name",
			expected: []types.Token{
				{Type: types.KEYWORD, Value: "name"},
			},
		},
		{
			name:  "keyword with dash",
			input: ":first-name",
			expected: []types.Token{
				{Type: types.KEYWORD, Value: "first-name"},
			},
		},
		{
			name:  "keyword with numbers",
			input: ":item2",
			expected: []types.Token{
				{Type: types.KEYWORD, Value: "item2"},
			},
		},
		{
			name:  "keyword in expression",
			input: "(:name \"Alice\")",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.KEYWORD, Value: "name"},
				{Type: types.STRING, Value: "Alice"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:  "hash map with keywords",
			input: "(hash-map :name \"Alice\" :age 30)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "hash-map"},
				{Type: types.KEYWORD, Value: "name"},
				{Type: types.STRING, Value: "Alice"},
				{Type: types.KEYWORD, Value: "age"},
				{Type: types.NUMBER, Value: "30"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:  "multiple keywords",
			input: ":a :b :c",
			expected: []types.Token{
				{Type: types.KEYWORD, Value: "a"},
				{Type: types.KEYWORD, Value: "b"},
				{Type: types.KEYWORD, Value: "c"},
			},
		},
		{
			name:  "keywords with special characters",
			input: ":test? :check! :count+",
			expected: []types.Token{
				{Type: types.KEYWORD, Value: "test?"},
				{Type: types.KEYWORD, Value: "check!"},
				{Type: types.KEYWORD, Value: "count+"},
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

func TestKeywordError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "lone colon",
			input: ":",
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
