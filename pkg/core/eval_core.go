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

// Macro represents a macro
type Macro struct {
	Name   Symbol
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

func (m *Macro) String() string {
	return fmt.Sprintf("#<macro:%s>", m.Name)
}

// bindParams binds function parameters to arguments, supporting variadic parameters
func bindParams(params *List, args []Value, env *Environment) error {
	paramList := listToSlice(params)

	// Check for variadic parameters (& rest-param)
	var restParamIndex = -1
	for i, param := range paramList {
		if sym, ok := param.(Symbol); ok && sym == "&" {
			if i == len(paramList)-1 {
				return fmt.Errorf("& must be followed by a parameter name")
			}
			if i != len(paramList)-2 {
				return fmt.Errorf("& parameter must be the last parameter")
			}
			restParamIndex = i
			break
		}
	}

	if restParamIndex >= 0 {
		// Variadic function
		minArgs := restParamIndex
		if len(args) < minArgs {
			return fmt.Errorf("function expects at least %d arguments, got %d", minArgs, len(args))
		}

		// Bind regular parameters
		for i := 0; i < restParamIndex; i++ {
			if sym, ok := paramList[i].(Symbol); ok {
				env.Set(sym, args[i])
			} else {
				return fmt.Errorf("parameter must be a symbol, got %T", paramList[i])
			}
		}

		// Bind rest parameter as a list
		restParamName, ok := paramList[restParamIndex+1].(Symbol)
		if !ok {
			return NewTypeError("rest parameter must be a symbol, got %T", paramList[restParamIndex+1])
		}

		// Collect remaining arguments into a list
		restArgs := args[restParamIndex:]
		if len(restArgs) == 0 {
			env.Set(restParamName, NewList()) // Empty list
		} else {
			env.Set(restParamName, NewList(restArgs...))
		}
	} else {
		// Non-variadic function - exact parameter count required
		if len(paramList) != len(args) {
			return NewArityError("function expects %d arguments, got %d", len(paramList), len(args))
		}

		for i, param := range paramList {
			if sym, ok := param.(Symbol); ok {
				env.Set(sym, args[i])
			} else {
				return NewTypeError("parameter must be a symbol, got %T", param)
			}
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

// EvalWithContext evaluates a Lisp expression with context tracking
func EvalWithContext(expr Value, env *Environment, ctx *EvaluationContext) (Value, error) {
	return evalWithContext(expr, env, ctx)
}

// Eval evaluates a Lisp expression (backward compatibility)
func Eval(expr Value, env *Environment) (Value, error) {
	ctx := NewEvaluationContext()
	return evalWithContext(expr, env, ctx)
}

// evalWithContext is the internal evaluation function with context tracking
func evalWithContext(expr Value, env *Environment, ctx *EvaluationContext) (Value, error) {
	switch v := expr.(type) {
	case Symbol:
		// Look up symbol in environment
		result, err := env.Get(v)
		if err != nil {
			return nil, ctx.EnhanceError(err)
		}
		return result, nil

	case *List:
		if v.IsEmpty() {
			return v, nil // Empty list evaluates to itself
		}

		// Check if first element is a special form
		if sym, ok := v.First().(Symbol); ok {
			ctx.PushFrame(string(sym), Position{})
			result, err := evalSpecialFormWithContext(sym, v.Rest(), env, ctx)
			ctx.PopFrame()
			if err == nil {
				return result, nil
			}
			// If it's a recognized special form but had an error, return the error
			if isSpecialForm(sym) {
				return nil, ctx.EnhanceError(err)
			}
		}

		// Regular function call
		return evalFunctionCallWithContext(v, env, ctx)

	case Number, String, Keyword, *Vector:
		// These evaluate to themselves
		return expr, nil

	default:
		return expr, nil
	}
}

// evalFunctionCallWithContext evaluates a function call with context tracking
func evalFunctionCallWithContext(list *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	// Evaluate the function
	fn, err := evalWithContext(list.First(), env, ctx)
	if err != nil {
		return nil, err
	}

	// Get function name for stack trace
	fnName := "anonymous"
	if sym, ok := list.First().(Symbol); ok {
		fnName = string(sym)
	}

	// Check if it's a macro - macros are expanded without evaluating arguments
	if macro, ok := fn.(*Macro); ok {
		ctx.PushFrame(fmt.Sprintf("macro %s", fnName), Position{})
		result, err := expandMacroWithContext(macro, list.Rest(), env, ctx)
		ctx.PopFrame()
		if err != nil {
			return nil, ctx.EnhanceError(err)
		}
		return result, nil
	}

	// Check if it's callable
	callable, ok := fn.(Function)
	if !ok {
		return nil, ctx.EnhanceError(NewTypeError("cannot call non-function: %T", fn))
	}

	// Evaluate arguments
	var args []Value
	current := list.Rest()
	
	for current != nil {
		arg, err := evalWithContext(current.First(), env, ctx)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		current = current.Rest()
	}

	// Call the function with context tracking
	ctx.PushFrame(fnName, Position{})
	result, err := callable.Call(args, env)
	ctx.PopFrame()
	
	if err != nil {
		return nil, ctx.EnhanceError(err)
	}
	
	return result, nil
}


// evalSpecialFormWithContext handles special forms with context tracking
func evalSpecialFormWithContext(sym Symbol, args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	// For now, just use the regular evalSpecialForm
	// TODO: Enhance special forms to use context for better error reporting
	_ = ctx // Suppress unused parameter warning
	return evalSpecialForm(sym, args, env)
}

// expandMacroWithContext expands a macro with context tracking  
func expandMacroWithContext(macro *Macro, args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	// For now, just use the regular expandMacro
	// TODO: Enhance macro expansion to use context for better error reporting
	_ = ctx // Suppress unused parameter warning
	return expandMacro(macro, args, env)
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

// expandMacro expands a macro call
func expandMacro(macro *Macro, args *List, env *Environment) (Value, error) {
	// Create new environment for macro expansion
	macroEnv := NewEnvironment(macro.Env)

	// Bind macro parameters to arguments (unevaluated)
	err := bindParams(macro.Params, listToSlice(args), macroEnv)
	if err != nil {
		return nil, err
	}

	// Evaluate macro body to get expansion
	expansion, err := Eval(macro.Body, macroEnv)
	if err != nil {
		return nil, err
	}

	// Evaluate the expansion
	return Eval(expansion, env)
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
