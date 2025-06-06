package types

import (
	"testing"
)

func TestTokenTypes(t *testing.T) {
	// Test that token types are properly defined
	tokens := []TokenType{
		LPAREN,
		RPAREN,
		NUMBER,
		SYMBOL,
		STRING,
		BOOLEAN,
	}

	// Just ensure they're distinct values
	for i, token1 := range tokens {
		for j, token2 := range tokens {
			if i != j && token1 == token2 {
				t.Errorf("Token types %d and %d have the same value", i, j)
			}
		}
	}
}

func TestValueTypes(t *testing.T) {
	// Test that value types implement the Value interface
	var _ Value = NumberValue(0)
	var _ Value = StringValue("")
	var _ Value = BooleanValue(false)

	// Test string representations
	if NumberValue(42).String() == "" {
		t.Error("NumberValue should have non-empty string representation")
	}

	if StringValue("hello").String() == "" {
		t.Error("StringValue should have non-empty string representation")
	}

	if BooleanValue(true).String() == "" {
		t.Error("BooleanValue should have non-empty string representation")
	}
}

func TestExprTypes(t *testing.T) {
	// Test that expression types implement the Expr interface
	var _ Expr = &NumberExpr{}
	var _ Expr = &StringExpr{}
	var _ Expr = &BooleanExpr{}
	var _ Expr = &SymbolExpr{}
	var _ Expr = &ListExpr{}
}
