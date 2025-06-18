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
			expected: &types.KeywordExpr{Value: "name", Position: types.Position{Line: 1, Column: 1}},
		},
		{
			name:     "keyword with dash",
			input:    ":first-name",
			expected: &types.KeywordExpr{Value: "first-name", Position: types.Position{Line: 1, Column: 1}},
		},
		{
			name:  "keyword in list",
			input: "(:name \"Alice\")",
			expected: &types.ListExpr{
				Position: types.Position{Line: 1, Column: 1},
				Elements: []types.Expr{
					&types.KeywordExpr{Value: "name", Position: types.Position{Line: 1, Column: 2}},
					&types.StringExpr{Value: "Alice", Position: types.Position{Line: 1, Column: 8}},
				},
			},
		},
		{
			name:  "hash map with keywords",
			input: "(hash-map :name \"Alice\" :age 30)",
			expected: &types.ListExpr{
				Position: types.Position{Line: 1, Column: 1},
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "hash-map", Position: types.Position{Line: 1, Column: 2}},
					&types.KeywordExpr{Value: "name", Position: types.Position{Line: 1, Column: 11}},
					&types.StringExpr{Value: "Alice", Position: types.Position{Line: 1, Column: 17}},
					&types.KeywordExpr{Value: "age", Position: types.Position{Line: 1, Column: 25}},
					&types.NumberExpr{Value: 30, Position: types.Position{Line: 1, Column: 30}},
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
