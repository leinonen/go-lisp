package types

import (
	"fmt"
	"strconv"
)

// TokenType represents the type of a token
type TokenType int

const (
	LPAREN TokenType = iota
	RPAREN
	NUMBER
	SYMBOL
	STRING
	BOOLEAN
)

// Token represents a single token in the input
type Token struct {
	Type  TokenType
	Value string
}

// Expr represents an expression in the AST
type Expr interface {
	String() string
}

// Value represents a value that can be returned from evaluation
type Value interface {
	String() string
}

// Expression types
type NumberExpr struct {
	Value float64
}

func (n *NumberExpr) String() string {
	return fmt.Sprintf("NumberExpr(%g)", n.Value)
}

type StringExpr struct {
	Value string
}

func (s *StringExpr) String() string {
	return fmt.Sprintf("StringExpr(%q)", s.Value)
}

type BooleanExpr struct {
	Value bool
}

func (b *BooleanExpr) String() string {
	return fmt.Sprintf("BooleanExpr(%t)", b.Value)
}

type SymbolExpr struct {
	Name string
}

func (s *SymbolExpr) String() string {
	return fmt.Sprintf("SymbolExpr(%s)", s.Name)
}

type ListExpr struct {
	Elements []Expr
}

func (l *ListExpr) String() string {
	return fmt.Sprintf("ListExpr(%v)", l.Elements)
}

// Value types
type NumberValue float64

func (n NumberValue) String() string {
	return strconv.FormatFloat(float64(n), 'g', -1, 64)
}

type StringValue string

func (s StringValue) String() string {
	return string(s)
}

type BooleanValue bool

func (b BooleanValue) String() string {
	if b {
		return "#t"
	}
	return "#f"
}
