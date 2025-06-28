package core

import "fmt"

// Value is the core interface for all Lisp values
type Value interface {
	String() string
}

// Symbol represents an interned symbol
type Symbol string

func (s Symbol) String() string {
	return string(s)
}

// Keyword represents an interned keyword (like Clojure keywords)
type Keyword string

func (k Keyword) String() string {
	return ":" + string(k)
}

// Number represents both integers and floats
type Number struct {
	Value interface{} // int64 or float64
}

func (n Number) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func (n Number) IsInteger() bool {
	_, ok := n.Value.(int64)
	return ok
}

func (n Number) IsFloat() bool {
	_, ok := n.Value.(float64)
	return ok
}

func (n Number) ToInt() int64 {
	if i, ok := n.Value.(int64); ok {
		return i
	}
	if f, ok := n.Value.(float64); ok {
		return int64(f)
	}
	return 0
}

func (n Number) ToFloat() float64 {
	if f, ok := n.Value.(float64); ok {
		return f
	}
	if i, ok := n.Value.(int64); ok {
		return float64(i)
	}
	return 0.0
}

// String represents a string value
type String string

func (s String) String() string {
	return fmt.Sprintf("%q", string(s))
}

// Nil represents the nil/null value
type Nil struct{}

func (n Nil) String() string {
	return "nil"
}

// List represents a linked list
type List struct {
	head Value
	tail *List
}

func (l *List) String() string {
	if l == nil {
		return "()"
	}

	result := "("
	current := l
	first := true

	for current != nil {
		if !first {
			result += " "
		}
		if current.head != nil {
			result += current.head.String()
		} else {
			result += "nil"
		}
		current = current.tail
		first = false
	}

	result += ")"
	return result
}

func (l *List) IsEmpty() bool {
	return l == nil
}

func (l *List) First() Value {
	if l == nil {
		return Nil{}
	}
	return l.head
}

func (l *List) Rest() *List {
	if l == nil {
		return nil
	}
	return l.tail
}

// Vector represents an indexed collection
type Vector struct {
	elements []Value
}

func (v *Vector) String() string {
	result := "["
	for i, elem := range v.elements {
		if i > 0 {
			result += " "
		}
		result += elem.String()
	}
	result += "]"
	return result
}

func (v *Vector) Get(index int) Value {
	if index < 0 || index >= len(v.elements) {
		return Nil{}
	}
	return v.elements[index]
}

func (v *Vector) Count() int {
	return len(v.elements)
}

// Environment represents a lexical environment for variable bindings
type Environment struct {
	bindings map[Symbol]Value
	parent   *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		bindings: make(map[Symbol]Value),
		parent:   parent,
	}
}

func (env *Environment) Get(sym Symbol) (Value, error) {
	if value, exists := env.bindings[sym]; exists {
		return value, nil
	}

	if env.parent != nil {
		return env.parent.Get(sym)
	}

	return nil, fmt.Errorf("undefined symbol: %s", sym)
}

func (env *Environment) Set(sym Symbol, value Value) {
	env.bindings[sym] = value
}

// Constructors
func NewList(elements ...Value) *List {
	if len(elements) == 0 {
		return nil
	}

	var result *List
	for i := len(elements) - 1; i >= 0; i-- {
		result = &List{head: elements[i], tail: result}
	}
	return result
}

func NewVector(elements ...Value) *Vector {
	return &Vector{elements: elements}
}

func NewNumber(value interface{}) Number {
	return Number{Value: value}
}

// Intern table for symbols
var internTable = make(map[string]Symbol)

// Intern ensures symbol uniqueness
func Intern(name string) Symbol {
	if sym, exists := internTable[name]; exists {
		return sym
	}
	sym := Symbol(name)
	internTable[name] = sym
	return sym
}

// Keyword intern table
var keywordInternTable = make(map[string]Keyword)

// InternKeyword ensures keyword uniqueness
func InternKeyword(name string) Keyword {
	if kw, exists := keywordInternTable[name]; exists {
		return kw
	}
	kw := Keyword(name)
	keywordInternTable[name] = kw
	return kw
}
