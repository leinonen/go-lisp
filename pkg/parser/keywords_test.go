package parser

import (
	"reflect"
	"testing"

	"github.com/leinonen/go-lisp/pkg/tokenizer"
	"github.com/leinonen/go-lisp/pkg/types"
)

func TestKeywordParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Expr
	}{
		{
			name:     "simple keyword",
			input:    ":name",
			expected: &types.KeywordExpr{Value: "name"},
		},
		{
			name:     "keyword with dash",
			input:    ":first-name",
			expected: &types.KeywordExpr{Value: "first-name"},
		},
		{
			name:  "keyword in list",
			input: "(:name \"Alice\")",
			expected: &types.ListExpr{
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "name"},
					&types.StringExpr{Value: "Alice"},
				},
			},
		},
		{
			name:  "hash map with keywords",
			input: "(hash-map :name \"Alice\" :age 30)",
			expected: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map"},
					&types.KeywordExpr{Value: "name"},
					&types.StringExpr{Value: "Alice"},
					&types.KeywordExpr{Value: "age"},
					&types.NumberExpr{Value: 30},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := tokenizer.NewTokenizer(tt.input)
			tokens := tokenizer.Tokenize()
			parser := NewParser(tokens)

			result, err := parser.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
