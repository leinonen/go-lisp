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

// DefinedValue represents the result of a successful define operation
type DefinedValue struct{}

func (d DefinedValue) String() string {
	return "defined"
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

// Enhanced Vector operations
func (v *Vector) Get(index int) Value {
	if index < 0 || index >= len(v.elements) {
		return Nil{}
	}
	return v.elements[index]
}

func (v *Vector) Append(val Value) *Vector {
	newElements := make([]Value, len(v.elements)+1)
	copy(newElements, v.elements)
	newElements[len(v.elements)] = val
	return &Vector{elements: newElements}
}

func (v *Vector) Update(index int, val Value) *Vector {
	if index < 0 || index >= len(v.elements) {
		return v // Return unchanged if index out of bounds
	}
	newElements := make([]Value, len(v.elements))
	copy(newElements, v.elements)
	newElements[index] = val
	return &Vector{elements: newElements}
}

// HashMap implementation
type HashMap struct {
	elements map[string]Value
}

func NewHashMap() *HashMap {
	return &HashMap{elements: make(map[string]Value)}
}

func (h *HashMap) String() string {
	if len(h.elements) == 0 {
		return "{}"
	}

	result := "{"
	first := true
	for k, v := range h.elements {
		if !first {
			result += " "
		}
		result += k + " " + v.String()
		first = false
	}
	result += "}"
	return result
}

func (h *HashMap) Get(key string) Value {
	if val, exists := h.elements[key]; exists {
		return val
	}
	return Nil{}
}

func (h *HashMap) Put(key string, val Value) *HashMap {
	newElements := make(map[string]Value)
	for k, v := range h.elements {
		newElements[k] = v
	}
	newElements[key] = val
	return &HashMap{elements: newElements}
}

func (h *HashMap) Keys() *Vector {
	keys := make([]Value, 0, len(h.elements))
	for k := range h.elements {
		keys = append(keys, String(k))
	}
	return &Vector{elements: keys}
}

func (h *HashMap) Values() *Vector {
	vals := make([]Value, 0, len(h.elements))
	for _, v := range h.elements {
		vals = append(vals, v)
	}
	return &Vector{elements: vals}
}

func (h *HashMap) Length() int {
	return len(h.elements)
}

// Set implementation (built on HashMap)
type Set struct {
	elements *HashMap
}

func NewSet() *Set {
	return &Set{elements: NewHashMap()}
}

func (s *Set) String() string {
	if s.elements.Length() == 0 {
		return "#{}"
	}

	result := "#{"
	keys := s.elements.Keys()
	for i := 0; i < keys.Length(); i++ {
		if i > 0 {
			result += " "
		}
		result += keys.elements[i].String()
	}
	result += "}"
	return result
}

func (s *Set) Add(val Value) *Set {
	key := val.String() // Simple string-based hashing for now
	newElements := s.elements.Put(key, Boolean(true))
	return &Set{elements: newElements}
}

func (s *Set) Contains(val Value) bool {
	key := val.String()
	result := s.elements.Get(key)
	_, isNil := result.(Nil)
	return !isNil
}

func (s *Set) Remove(val Value) *Set {
	newElements := make(map[string]Value)
	key := val.String()
	for k, v := range s.elements.elements {
		if k != key {
			newElements[k] = v
		}
	}
	return &Set{elements: &HashMap{elements: newElements}}
}

func (s *Set) Length() int {
	return s.elements.Length()
}
