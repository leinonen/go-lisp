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
				{Type: types.KEYWORD, Value: "name", Position: types.Position{Line: 1, Column: 1}},
			},
		},
		{
			name:  "keyword with dash",
			input: ":first-name",
			expected: []types.Token{
				{Type: types.KEYWORD, Value: "first-name", Position: types.Position{Line: 1, Column: 1}},
			},
		},
		{
			name:  "keyword with numbers",
			input: ":item2",
			expected: []types.Token{
				{Type: types.KEYWORD, Value: "item2", Position: types.Position{Line: 1, Column: 1}},
			},
		},
		{
			name:  "keyword in expression",
			input: "(:name \"Alice\")",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.KEYWORD, Value: "name", Position: types.Position{Line: 1, Column: 2}},
				{Type: types.STRING, Value: "Alice", Position: types.Position{Line: 1, Column: 8}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 15}},
			},
		},
		{
			name:  "hash map with keywords",
			input: "(hash-map :name \"Alice\" :age 30)",
			expected: []types.Token{
				{Type: types.LPAREN, Value: "(", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.SYMBOL, Value: "hash-map", Position: types.Position{Line: 1, Column: 2}},
				{Type: types.KEYWORD, Value: "name", Position: types.Position{Line: 1, Column: 11}},
				{Type: types.STRING, Value: "Alice", Position: types.Position{Line: 1, Column: 17}},
				{Type: types.KEYWORD, Value: "age", Position: types.Position{Line: 1, Column: 25}},
				{Type: types.NUMBER, Value: "30", Position: types.Position{Line: 1, Column: 30}},
				{Type: types.RPAREN, Value: ")", Position: types.Position{Line: 1, Column: 32}},
			},
		},
		{
			name:  "multiple keywords",
			input: ":a :b :c",
			expected: []types.Token{
				{Type: types.KEYWORD, Value: "a", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.KEYWORD, Value: "b", Position: types.Position{Line: 1, Column: 4}},
				{Type: types.KEYWORD, Value: "c", Position: types.Position{Line: 1, Column: 7}},
			},
		},
		{
			name:  "keywords with special characters",
			input: ":test? :check! :count+",
			expected: []types.Token{
				{Type: types.KEYWORD, Value: "test?", Position: types.Position{Line: 1, Column: 1}},
				{Type: types.KEYWORD, Value: "check!", Position: types.Position{Line: 1, Column: 8}},
				{Type: types.KEYWORD, Value: "count+", Position: types.Position{Line: 1, Column: 16}},
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
