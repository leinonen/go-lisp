package minimal

// Core data types for the minimal Lisp kernel

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

// List represents a Lisp list
type List struct {
	elements []Value
}

// NewList creates a new list
func NewList(values ...Value) *List {
	return &List{elements: values}
}

func (l *List) String() string {
	if len(l.elements) == 0 {
		return "()"
	}
	result := "("
	for i, v := range l.elements {
		if i > 0 {
			result += " "
		}
		result += v.String()
	}
	result += ")"
	return result
}

// First returns the first element
func (l *List) First() Value {
	if len(l.elements) == 0 {
		return nil
	}
	return l.elements[0]
}

// Rest returns a new list with all but the first element
func (l *List) Rest() *List {
	if len(l.elements) <= 1 {
		return NewList()
	}
	return NewList(l.elements[1:]...)
}

// IsEmpty returns true if the list is empty
func (l *List) IsEmpty() bool {
	return len(l.elements) == 0
}

// Length returns the number of elements
func (l *List) Length() int {
	return len(l.elements)
}

// Number represents a numeric value
type Number float64

func (n Number) String() string {
	if n == Number(int64(n)) {
		return fmt.Sprintf("%.0f", float64(n))
	}
	return fmt.Sprintf("%g", float64(n))
}

// Boolean represents a boolean value
type Boolean bool

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}

// Nil represents the nil value
type Nil struct{}

func (n Nil) String() string {
	return "nil"
}

// String represents a string value
type String string

func (s String) String() string {
	return fmt.Sprintf("\"%s\"", string(s))
}

// Vector represents a Clojure-style vector with square bracket notation
type Vector struct {
	elements []Value
}

// NewVector creates a new vector
func NewVector(values ...Value) *Vector {
	return &Vector{elements: values}
}

func (v *Vector) String() string {
	if len(v.elements) == 0 {
		return "[]"
	}
	result := "["
	for i, val := range v.elements {
		if i > 0 {
			result += " "
		}
		result += val.String()
	}
	result += "]"
	return result
}

// First returns the first element
func (v *Vector) First() Value {
	if len(v.elements) == 0 {
		return nil
	}
	return v.elements[0]
}

// Rest returns a new vector with all but the first element
func (v *Vector) Rest() *Vector {
	if len(v.elements) <= 1 {
		return NewVector()
	}
	return NewVector(v.elements[1:]...)
}

// IsEmpty returns true if the vector is empty
func (v *Vector) IsEmpty() bool {
	return len(v.elements) == 0
}

// Length returns the number of elements
func (v *Vector) Length() int {
	return len(v.elements)
}

// ToList converts the vector to a list for compatibility
func (v *Vector) ToList() *List {
	return NewList(v.elements...)
}
