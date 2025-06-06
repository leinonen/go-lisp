package parser

import (
	"reflect"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []types.Token
		expected types.Expr
	}{
		{
			name:     "single number",
			tokens:   []types.Token{{Type: types.NUMBER, Value: "42"}},
			expected: &types.NumberExpr{Value: 42},
		},
		{
			name:     "single symbol",
			tokens:   []types.Token{{Type: types.SYMBOL, Value: "x"}},
			expected: &types.SymbolExpr{Name: "x"},
		},
		{
			name:     "boolean true",
			tokens:   []types.Token{{Type: types.BOOLEAN, Value: "#t"}},
			expected: &types.BooleanExpr{Value: true},
		},
		{
			name:     "boolean false",
			tokens:   []types.Token{{Type: types.BOOLEAN, Value: "#f"}},
			expected: &types.BooleanExpr{Value: false},
		},
		{
			name:     "string literal",
			tokens:   []types.Token{{Type: types.STRING, Value: "hello"}},
			expected: &types.StringExpr{Value: "hello"},
		},
		{
			name: "simple addition",
			tokens: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "+"},
				{Type: types.NUMBER, Value: "1"},
				{Type: types.NUMBER, Value: "2"},
				{Type: types.RPAREN, Value: ")"},
			},
			expected: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.NumberExpr{Value: 1},
					&types.NumberExpr{Value: 2},
				},
			},
		},
		{
			name: "nested expression",
			tokens: []types.Token{
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
			expected: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "*"},
							&types.NumberExpr{Value: 2},
							&types.NumberExpr{Value: 3},
						},
					},
					&types.NumberExpr{Value: 4},
				},
			},
		},
		{
			name: "empty list",
			tokens: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.RPAREN, Value: ")"},
			},
			expected: &types.ListExpr{Elements: []types.Expr{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.tokens)
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

func TestParserError(t *testing.T) {
	tests := []struct {
		name   string
		tokens []types.Token
	}{
		{
			name: "unmatched opening paren",
			tokens: []types.Token{
				{Type: types.LPAREN, Value: "("},
				{Type: types.SYMBOL, Value: "+"},
				{Type: types.NUMBER, Value: "1"},
			},
		},
		{
			name: "unmatched closing paren",
			tokens: []types.Token{
				{Type: types.SYMBOL, Value: "+"},
				{Type: types.NUMBER, Value: "1"},
				{Type: types.RPAREN, Value: ")"},
			},
		},
		{
			name:   "empty input",
			tokens: []types.Token{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.tokens)
			_, err := parser.Parse()

			if err == nil {
				t.Errorf("expected error for tokens %v", tt.tokens)
			}
		})
	}
}
