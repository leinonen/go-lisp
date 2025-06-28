package core

import (
	"fmt"
)

// Function interface for callable values
type Function interface {
	Call(args []Value, env *Environment) (Value, error)
}

// BuiltinFunction represents a built-in function
type BuiltinFunction struct {
	Name string
	Fn   func(args []Value, env *Environment) (Value, error)
}

func (bf *BuiltinFunction) Call(args []Value, env *Environment) (Value, error) {
	return bf.Fn(args, env)
}

func (bf *BuiltinFunction) String() string {
	return fmt.Sprintf("#<builtin:%s>", bf.Name)
}

// UserFunction represents a user-defined function
type UserFunction struct {
	Params *List
	Body   Value
	Env    *Environment
}

func (uf *UserFunction) Call(args []Value, env *Environment) (Value, error) {
	// Create new environment for function execution
	fnEnv := NewEnvironment(uf.Env)

	// Bind parameters to arguments
	err := bindParams(uf.Params, args, fnEnv)
	if err != nil {
		return nil, err
	}

	// Evaluate function body
	return Eval(uf.Body, fnEnv)
}

func (uf *UserFunction) String() string {
	return "#<function>"
}

// bindParams binds function parameters to arguments
func bindParams(params *List, args []Value, env *Environment) error {
	paramList := listToSlice(params)

	if len(paramList) != len(args) {
		return fmt.Errorf("function expects %d arguments, got %d", len(paramList), len(args))
	}

	for i, param := range paramList {
		if sym, ok := param.(Symbol); ok {
			env.Set(sym, args[i])
		} else {
			return fmt.Errorf("parameter must be a symbol, got %T", param)
		}
	}

	return nil
}

// listToSlice converts a List to a slice of Values
func listToSlice(list *List) []Value {
	var result []Value
	current := list

	for current != nil {
		result = append(result, current.First())
		current = current.Rest()
	}

	return result
}

// Eval evaluates a Lisp expression
func Eval(expr Value, env *Environment) (Value, error) {
	switch v := expr.(type) {
	case Symbol:
		// Look up symbol in environment
		return env.Get(v)

	case *List:
		if v.IsEmpty() {
			return v, nil // Empty list evaluates to itself
		}

		// Check if first element is a special form
		if sym, ok := v.First().(Symbol); ok {
			result, err := evalSpecialForm(sym, v.Rest(), env)
			if err == nil {
				return result, nil
			}
			// If it's a recognized special form but had an error, return the error
			if isSpecialForm(sym) {
				return nil, err
			}
		}

		// Regular function call
		return evalFunctionCall(v, env)

	case Number, String, Keyword, *Vector:
		// These evaluate to themselves
		return expr, nil

	default:
		return expr, nil
	}
}

// evalFunctionCall evaluates a function call
func evalFunctionCall(list *List, env *Environment) (Value, error) {
	// Evaluate the function
	fn, err := Eval(list.First(), env)
	if err != nil {
		return nil, err
	}

	// Check if it's callable
	callable, ok := fn.(Function)
	if !ok {
		return nil, fmt.Errorf("cannot call non-function: %T", fn)
	}

	// Evaluate arguments
	var args []Value
	current := list.Rest()

	for current != nil {
		arg, err := Eval(current.First(), env)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		current = current.Rest()
	}

	// Call the function
	return callable.Call(args, env)
}

// isTruthy determines if a value is truthy
func isTruthy(v Value) bool {
	switch val := v.(type) {
	case Nil:
		return false
	case Number:
		if val.IsInteger() {
			return val.ToInt() != 0
		}
		return val.ToFloat() != 0.0
	case String:
		return string(val) != ""
	default:
		return true
	}
}

// valuesEqual compares two values for equality
func valuesEqual(a, b Value) bool {
	switch va := a.(type) {
	case Symbol:
		if vb, ok := b.(Symbol); ok {
			return va == vb
		}
	case String:
		if vb, ok := b.(String); ok {
			return va == vb
		}
	case Number:
		if vb, ok := b.(Number); ok {
			return va.ToFloat() == vb.ToFloat()
		}
	case Keyword:
		if vb, ok := b.(Keyword); ok {
			return va == vb
		}
	case Nil:
		_, ok := b.(Nil)
		return ok
	}
	return false
}

// NewCoreEnvironment creates an environment with core primitives
// This function coordinates the setup from all specialized modules
func NewCoreEnvironment() *Environment {
	env := NewEnvironment(nil)

	// Set up different categories of operations
	setupArithmeticOperations(env) // +, -, *, /, %, =, <, >, >=, <=
	setupCollectionOperations(env) // count, empty?, nth, conj, cons, first, rest, list, list?, vector?
	setupStringOperations(env)     // str, substring, string-split, string-replace, string-contains?, string-trim, string?
	setupIOOperations(env)         // println, prn, slurp, spit, file-exists?, list-dir
	setupMetaProgramming(env)      // eval, read-string, symbol?, number?, keyword?, nil?, fn?

	return env
}
