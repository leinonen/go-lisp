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

// FunctionValue represents a function with parameters and body
type FunctionValue struct {
	Params []string    // parameter names
	Body   Expr        // function body expression
	Env    Environment // captured environment for closures
}

func (f FunctionValue) String() string {
	return fmt.Sprintf("#<function(%v)>", f.Params)
}

// ListValue represents a list of values
type ListValue struct {
	Elements []Value
}

func (l *ListValue) String() string {
	if len(l.Elements) == 0 {
		return "()"
	}
	var elements []string
	for _, elem := range l.Elements {
		elements = append(elements, elem.String())
	}
	result := "("
	for i, elem := range elements {
		if i > 0 {
			result += " "
		}
		result += elem
	}
	result += ")"
	return result
}

// Environment interface for closures
type Environment interface {
	Get(name string) (Value, bool)
	Set(name string, value Value)
	NewChildEnvironment() Environment
}

// ModuleValue represents a module with exported bindings
type ModuleValue struct {
	Name    string
	Exports map[string]Value
	Env     Environment // module's internal environment
}

func (m *ModuleValue) String() string {
	return fmt.Sprintf("#<module:%s>", m.Name)
}

// ModuleExpr represents a module definition expression
type ModuleExpr struct {
	Name    string
	Exports []string
	Body    []Expr
}

func (m *ModuleExpr) String() string {
	return fmt.Sprintf("ModuleExpr(name:%s, exports:%v, body:%v)", m.Name, m.Exports, m.Body)
}

// ImportExpr represents an import expression
type ImportExpr struct {
	ModuleName string
}

func (i *ImportExpr) String() string {
	return fmt.Sprintf("ImportExpr(%s)", i.ModuleName)
}

// LoadExpr represents a load file expression
type LoadExpr struct {
	Filename string
}

func (l *LoadExpr) String() string {
	return fmt.Sprintf("LoadExpr(%s)", l.Filename)
}
