package core

import "fmt"

// Value is the core interface for all Lisp values
type Value interface {
	String() string
}

// SourceLocated is an optional interface for values that have source location
type SourceLocated interface {
	GetPosition() Position
	SetPosition(Position)
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

// Call makes keywords callable as functions on hash-maps
func (k Keyword) Call(args []Value, env *Environment) (Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("keyword %s expects 1-2 arguments", k)
	}

	// First argument should be a hash-map
	if hm, ok := args[0].(*HashMap); ok {
		value := hm.Get(k)
		// If value is nil and we have a default value, return the default
		if _, isNil := value.(Nil); isNil && len(args) == 2 {
			return args[1], nil
		}
		return value, nil
	}

	return nil, fmt.Errorf("keyword %s can only be called on hash-maps, got %T", k, args[0])
}

// Number represents both integers and floats
type Number struct {
	Value any // int64 or float64
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

// HashMap represents a key-value mapping
type HashMap struct {
	pairs map[string]Value
	keys  []Value // Maintain insertion order
}

func (h *HashMap) String() string {
	result := "{"
	first := true
	for _, key := range h.keys {
		if !first {
			result += " "
		}
		result += key.String() + " " + h.pairs[h.keyToString(key)].String()
		first = false
	}
	result += "}"
	return result
}

func (h *HashMap) keyToString(key Value) string {
	return key.String()
}

func (h *HashMap) Get(key Value) Value {
	if value, exists := h.pairs[h.keyToString(key)]; exists {
		return value
	}
	return Nil{}
}

func (h *HashMap) Set(key Value, value Value) {
	keyStr := h.keyToString(key)
	if _, exists := h.pairs[keyStr]; !exists {
		h.keys = append(h.keys, key)
	}
	h.pairs[keyStr] = value
}

func (h *HashMap) Count() int {
	return len(h.keys)
}

func (h *HashMap) ContainsKey(key Value) bool {
	_, exists := h.pairs[h.keyToString(key)]
	return exists
}

// Set represents a collection of unique values
type Set struct {
	elements map[string]Value
	order    []Value // Maintain insertion order
}

func (s *Set) String() string {
	result := "#{"
	first := true
	for _, elem := range s.order {
		if !first {
			result += " "
		}
		result += elem.String()
		first = false
	}
	result += "}"
	return result
}

func (s *Set) elemToString(elem Value) string {
	return elem.String()
}

func (s *Set) Add(elem Value) {
	elemStr := s.elemToString(elem)
	if _, exists := s.elements[elemStr]; !exists {
		s.elements[elemStr] = elem
		s.order = append(s.order, elem)
	}
}

func (s *Set) Contains(elem Value) bool {
	_, exists := s.elements[s.elemToString(elem)]
	return exists
}

func (s *Set) Count() int {
	return len(s.order)
}

func (s *Set) Remove(elem Value) {
	elemStr := s.elemToString(elem)
	if _, exists := s.elements[elemStr]; exists {
		delete(s.elements, elemStr)
		// Remove from order slice
		for i, e := range s.order {
			if s.elemToString(e) == elemStr {
				s.order = append(s.order[:i], s.order[i+1:]...)
				break
			}
		}
	}
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

func NewHashMap() *HashMap {
	return &HashMap{
		pairs: make(map[string]Value),
		keys:  make([]Value, 0),
	}
}

func NewHashMapWithPairs(pairs ...Value) *HashMap {
	hm := NewHashMap()
	for i := 0; i < len(pairs)-1; i += 2 {
		hm.Set(pairs[i], pairs[i+1])
	}
	return hm
}

func NewSet() *Set {
	return &Set{
		elements: make(map[string]Value),
		order:    make([]Value, 0),
	}
}

func NewSetWithElements(elements ...Value) *Set {
	s := NewSet()
	for _, elem := range elements {
		s.Add(elem)
	}
	return s
}

func NewNumber(value any) Number {
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
