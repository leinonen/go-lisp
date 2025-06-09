package types

import (
	"fmt"
	"math/big"
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
	KEYWORD
	QUOTE
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

type BigNumberExpr struct {
	Value string // Store as string to preserve precision during parsing
}

func (b *BigNumberExpr) String() string {
	return fmt.Sprintf("BigNumberExpr(%s)", b.Value)
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

type KeywordExpr struct {
	Value string
}

func (k *KeywordExpr) String() string {
	return fmt.Sprintf("KeywordExpr(:%s)", k.Value)
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
	// For very large numbers that exceed float64 precision, show scientific notation
	f := float64(n)
	if f > 1e15 || f < -1e15 {
		return fmt.Sprintf("%.6e", f)
	}
	// Check if it's an integer
	if f == float64(int64(f)) && f >= -1e15 && f <= 1e15 {
		return fmt.Sprintf("%.0f", f)
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// BigNumberValue represents large integers using math/big
type BigNumberValue struct {
	Value *big.Int
}

func (b *BigNumberValue) String() string {
	return b.Value.String()
}

// Helper function to create a BigNumberValue
func NewBigNumberValue(x *big.Int) *BigNumberValue {
	return &BigNumberValue{Value: new(big.Int).Set(x)}
}

// Helper function to create a BigNumberValue from int64
func NewBigNumberFromInt64(x int64) *BigNumberValue {
	return &BigNumberValue{Value: big.NewInt(x)}
}

// Helper function to create a BigNumberValue from string
func NewBigNumberFromString(s string) (*BigNumberValue, bool) {
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(s, 10)
	if !ok {
		return nil, false
	}
	return &BigNumberValue{Value: bigInt}, true
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

type KeywordValue string

func (k KeywordValue) String() string {
	return ":" + string(k)
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

// MacroValue represents a macro with parameters and body
type MacroValue struct {
	Params []string    // parameter names
	Body   Expr        // macro body expression
	Env    Environment // captured environment for closures
}

func (m MacroValue) String() string {
	return fmt.Sprintf("#<macro(%v)>", m.Params)
}

// QuotedValue represents a quoted expression that should not be evaluated
type QuotedValue struct {
	Value Expr
}

func (q QuotedValue) String() string {
	// Special handling for quoted symbols - return just the symbol name
	if symbolExpr, ok := q.Value.(*SymbolExpr); ok {
		return symbolExpr.Name
	}
	return q.Value.String()
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

// HashMapValue represents a hash map with string keys and arbitrary values
type HashMapValue struct {
	Elements map[string]Value
}

func (h *HashMapValue) String() string {
	if len(h.Elements) == 0 {
		return "{}"
	}
	result := "{"
	first := true
	for key, value := range h.Elements {
		if !first {
			result += " "
		}
		result += fmt.Sprintf("%q %s", key, value.String())
		first = false
	}
	result += "}"
	return result
}

// NilValue represents the absence of a value (nil/null)
type NilValue struct{}

func (n *NilValue) String() string {
	return "nil"
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

// ArithmeticFunctionValue represents a built-in arithmetic operation as a callable function
type ArithmeticFunctionValue struct {
	Operation string // "+", "-", "*", "/", "%"
}

func (a ArithmeticFunctionValue) String() string {
	return fmt.Sprintf("#<built-in:%s>", a.Operation)
}
