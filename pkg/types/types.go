package types

import (
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
)

// TokenType represents the type of a token
type TokenType int

const (
	LPAREN TokenType = iota
	RPAREN
	LBRACKET
	RBRACKET
	NUMBER
	SYMBOL
	STRING
	BOOLEAN
	KEYWORD
	QUOTE
)

// Position represents the location of a token in the source code
type Position struct {
	Line   int // 1-based line number
	Column int // 1-based column number
}

// Token represents a single token in the input
type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

// Expr represents an expression in the AST
type Expr interface {
	String() string
	GetPosition() Position
}

// Value represents a value that can be returned from evaluation
type Value interface {
	String() string
}

// Expression types
type NumberExpr struct {
	Value    float64
	Position Position
}

func (n *NumberExpr) String() string {
	return fmt.Sprintf("NumberExpr(%g)", n.Value)
}

func (n *NumberExpr) GetPosition() Position {
	return n.Position
}

type BigNumberExpr struct {
	Value    string // Store as string to preserve precision during parsing
	Position Position
}

func (b *BigNumberExpr) String() string {
	return fmt.Sprintf("BigNumberExpr(%s)", b.Value)
}

func (b *BigNumberExpr) GetPosition() Position {
	return b.Position
}

type StringExpr struct {
	Value    string
	Position Position
}

func (s *StringExpr) String() string {
	return fmt.Sprintf("StringExpr(%q)", s.Value)
}

func (s *StringExpr) GetPosition() Position {
	return s.Position
}

type BooleanExpr struct {
	Value    bool
	Position Position
}

func (b *BooleanExpr) String() string {
	return fmt.Sprintf("BooleanExpr(%t)", b.Value)
}

func (b *BooleanExpr) GetPosition() Position {
	return b.Position
}

type SymbolExpr struct {
	Name     string
	Position Position
}

func (s *SymbolExpr) String() string {
	return fmt.Sprintf("SymbolExpr(%s)", s.Name)
}

func (s *SymbolExpr) GetPosition() Position {
	return s.Position
}

type KeywordExpr struct {
	Value    string
	Position Position
}

func (k *KeywordExpr) String() string {
	return fmt.Sprintf("KeywordExpr(:%s)", k.Value)
}

func (k *KeywordExpr) GetPosition() Position {
	return k.Position
}

type ListExpr struct {
	Elements []Expr
	Position Position
}

func (l *ListExpr) String() string {
	return fmt.Sprintf("ListExpr(%v)", l.Elements)
}

func (l *ListExpr) GetPosition() Position {
	return l.Position
}

type BracketExpr struct {
	Elements []Expr
	Position Position
}

func (b *BracketExpr) String() string {
	return fmt.Sprintf("BracketExpr(%v)", b.Elements)
}

func (b *BracketExpr) GetPosition() Position {
	return b.Position
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
		return "true"
	}
	return "false"
}

// KeywordValue represents a keyword (self-evaluating symbol starting with :)
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
		if elem == nil {
			elements = append(elements, "nil")
		} else {
			elements = append(elements, elem.String())
		}
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
	Name     string
	Exports  []string
	Body     []Expr
	Position Position
}

func (m *ModuleExpr) String() string {
	return fmt.Sprintf("ModuleExpr(name:%s, exports:%v, body:%v)", m.Name, m.Exports, m.Body)
}

func (m *ModuleExpr) GetPosition() Position {
	return m.Position
}

// ImportExpr represents an import expression
type ImportExpr struct {
	ModuleName string
	Position   Position
}

func (i *ImportExpr) String() string {
	return fmt.Sprintf("ImportExpr(%s)", i.ModuleName)
}

func (i *ImportExpr) GetPosition() Position {
	return i.Position
}

// LoadExpr represents a load file expression
type LoadExpr struct {
	Filename string
	Position Position
}

func (l *LoadExpr) String() string {
	return fmt.Sprintf("LoadExpr(%s)", l.Filename)
}

func (l *LoadExpr) GetPosition() Position {
	return l.Position
}

// RequireExpr represents a require expression that combines load and import
type RequireExpr struct {
	Filename string
	AsAlias  string   // For :as syntax
	OnlyList []string // For :only syntax
	Position Position
}

func (r *RequireExpr) String() string {
	if r.AsAlias != "" {
		return fmt.Sprintf("RequireExpr(%s :as %s)", r.Filename, r.AsAlias)
	}
	if len(r.OnlyList) > 0 {
		return fmt.Sprintf("RequireExpr(%s :only %v)", r.Filename, r.OnlyList)
	}
	return fmt.Sprintf("RequireExpr(%s)", r.Filename)
}

func (r *RequireExpr) GetPosition() Position {
	return r.Position
}

// ArithmeticFunctionValue represents a built-in arithmetic operation as a callable function
type ArithmeticFunctionValue struct {
	Operation string // "+", "-", "*", "/", "%"
}

func (a ArithmeticFunctionValue) String() string {
	return fmt.Sprintf("#<built-in:%s>", a.Operation)
}

// BuiltinFunctionValue represents a built-in function as a callable value
type BuiltinFunctionValue struct {
	Name string // name of the builtin function
}

func (b BuiltinFunctionValue) String() string {
	return fmt.Sprintf("#<built-in:%s>", b.Name)
}

// AtomValue represents a mutable reference to a value, protected by a mutex
type AtomValue struct {
	value Value
	mutex sync.RWMutex
}

// NewAtom creates a new atom with the given initial value
func NewAtom(initialValue Value) *AtomValue {
	return &AtomValue{
		value: initialValue,
	}
}

// Value returns the current value of the atom (thread-safe read)
func (a *AtomValue) Value() Value {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.value
}

// SetValue sets the atom's value (thread-safe write)
func (a *AtomValue) SetValue(newValue Value) Value {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.value = newValue
	return newValue
}

// SwapValue applies a function to the current value and updates the atom (thread-safe)
func (a *AtomValue) SwapValue(fn func(Value) Value) Value {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.value = fn(a.value)
	return a.value
}

func (a *AtomValue) String() string {
	currentValue := a.Value()
	if currentValue == nil {
		return "#<atom:nil>"
	}
	return fmt.Sprintf("#<atom:%s>", currentValue.String())
}

// FutureValue represents a future result from a goroutine
type FutureValue struct {
	result chan Value
	err    chan error
	done   atomic.Bool
}

// NewFuture creates a new future
func NewFuture() *FutureValue {
	return &FutureValue{
		result: make(chan Value, 1),
		err:    make(chan error, 1),
	}
}

// SetResult sets the result of the future
func (f *FutureValue) SetResult(value Value) {
	if f.done.CompareAndSwap(false, true) {
		f.result <- value
		// Don't close the channels immediately, let Wait() handle it
	}
}

// SetError sets the error of the future
func (f *FutureValue) SetError(err error) {
	if f.done.CompareAndSwap(false, true) {
		f.err <- err
		// Don't close the channels immediately, let Wait() handle it
	}
}

// Wait waits for the future to complete and returns the result or error
func (f *FutureValue) Wait() (Value, error) {
	select {
	case result, ok := <-f.result:
		if !ok {
			// Channel was closed without a value, should not happen
			return nil, fmt.Errorf("future result channel closed unexpectedly")
		}
		return result, nil
	case err, ok := <-f.err:
		if !ok {
			// Channel was closed without an error, should not happen
			return nil, fmt.Errorf("future error channel closed unexpectedly")
		}
		return nil, err
	}
}

// IsDone returns true if the future is complete
func (f *FutureValue) IsDone() bool {
	return f.done.Load()
}

func (f *FutureValue) String() string {
	if f.IsDone() {
		return "#<future:done>"
	}
	return "#<future:pending>"
}

// ChannelValue represents a channel for goroutine communication
type ChannelValue struct {
	ch     chan Value
	mutex  sync.RWMutex
	size   int
	closed bool
}

// NewChannel creates a new channel with the given buffer size
func NewChannel(size int) *ChannelValue {
	return &ChannelValue{
		ch:   make(chan Value, size),
		size: size,
	}
}

// Send sends a value to the channel (blocking)
func (c *ChannelValue) Send(value Value) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.closed {
		return fmt.Errorf("cannot send to closed channel")
	}

	c.ch <- value
	return nil
}

// Receive receives a value from the channel (blocking)
func (c *ChannelValue) Receive() (Value, bool) {
	value, ok := <-c.ch
	return value, ok
}

// TryReceive tries to receive a value from the channel (non-blocking)
func (c *ChannelValue) TryReceive() (Value, bool) {
	select {
	case value, ok := <-c.ch:
		return value, ok
	default:
		return nil, false
	}
}

// Close closes the channel
func (c *ChannelValue) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.closed {
		close(c.ch)
		c.closed = true
	}
}

// IsClosed returns true if the channel is closed
func (c *ChannelValue) IsClosed() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.closed
}

func (c *ChannelValue) String() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.closed {
		return fmt.Sprintf("#<channel:closed:size=%d>", c.size)
	}
	return fmt.Sprintf("#<channel:open:size=%d>", c.size)
}

// WaitGroupValue represents a wait group for coordinating goroutines
type WaitGroupValue struct {
	wg *sync.WaitGroup
}

// NewWaitGroup creates a new wait group
func NewWaitGroup() *WaitGroupValue {
	return &WaitGroupValue{
		wg: &sync.WaitGroup{},
	}
}

// Add adds delta to the wait group counter
func (w *WaitGroupValue) Add(delta int) {
	w.wg.Add(delta)
}

// Done decrements the wait group counter
func (w *WaitGroupValue) Done() {
	w.wg.Done()
}

// Wait waits for the wait group counter to go to zero
func (w *WaitGroupValue) Wait() {
	w.wg.Wait()
}

func (w *WaitGroupValue) String() string {
	return "#<wait-group>"
}

// RecurException represents a recur call that should jump back to the nearest loop
type RecurException struct {
	Args []Value // Arguments to pass to the next iteration
}

func (r *RecurException) Error() string {
	return "recur called outside of loop"
}

func (r *RecurException) String() string {
	return fmt.Sprintf("RecurException(%v)", r.Args)
}

// Helper function to create a RecurException
func NewRecurException(args []Value) *RecurException {
	return &RecurException{Args: args}
}

// PositionalError wraps an error with position information
type PositionalError struct {
	Message  string
	Position Position
	Cause    error
}

func (pe *PositionalError) Error() string {
	if pe.Position.Line > 0 {
		return fmt.Sprintf("line %d, column %d: %s", pe.Position.Line, pe.Position.Column, pe.Message)
	}
	return pe.Message
}

func (pe *PositionalError) Unwrap() error {
	return pe.Cause
}

// NewPositionalError creates a new positional error
func NewPositionalError(message string, pos Position, cause error) *PositionalError {
	return &PositionalError{
		Message:  message,
		Position: pos,
		Cause:    cause,
	}
}
