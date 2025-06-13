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

func TestListValue(t *testing.T) {
	// Test that ListValue implements the Value interface
	var _ Value = &ListValue{}

	// Test empty list
	emptyList := &ListValue{Elements: []Value{}}
	if emptyList.String() != "()" {
		t.Errorf("Empty list should be '()', got %q", emptyList.String())
	}

	// Test list with single element
	singleList := &ListValue{Elements: []Value{NumberValue(42)}}
	expected := "(42)"
	if singleList.String() != expected {
		t.Errorf("Single element list should be %q, got %q", expected, singleList.String())
	}

	// Test list with multiple elements
	multiList := &ListValue{Elements: []Value{
		NumberValue(1),
		NumberValue(2),
		NumberValue(3),
	}}
	expected = "(1 2 3)"
	if multiList.String() != expected {
		t.Errorf("Multi element list should be %q, got %q", expected, multiList.String())
	}

	// Test mixed type list
	mixedList := &ListValue{Elements: []Value{
		NumberValue(42),
		StringValue("hello"),
		BooleanValue(true),
	}}
	expected = "(42 hello true)"
	if mixedList.String() != expected {
		t.Errorf("Mixed type list should be %q, got %q", expected, mixedList.String())
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
