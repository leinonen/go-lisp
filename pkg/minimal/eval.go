package minimal

// Core evaluation logic for the minimal Lisp kernel

import "fmt"

// Eval evaluates a Lisp expression in the given environment
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
			// Not a special form - check if the error is our specific "not a special form" error
			if err.Error() != "not a special form" {
				return nil, err // Return actual errors from special forms
			}
			// Fall through to function application
		}

		// Regular function application
		return Apply(v, env)

	default:
		// Self-evaluating expressions (numbers, strings, booleans, etc.)
		return expr, nil
	}
}

// Apply applies a function to arguments
func Apply(list *List, env *Environment) (Value, error) {
	if list.IsEmpty() {
		return nil, fmt.Errorf("cannot apply empty list")
	}

	// Evaluate the function
	fn, err := Eval(list.First(), env)
	if err != nil {
		return nil, err
	}

	function, ok := fn.(Function)
	if !ok {
		return nil, fmt.Errorf("%v is not a function", fn)
	}

	// Evaluate all arguments
	args := make([]Value, 0)
	for current := list.Rest(); !current.IsEmpty(); current = current.Rest() {
		evaluated, err := Eval(current.First(), env)
		if err != nil {
			return nil, err
		}
		args = append(args, evaluated)
	}

	// Call the function
	return function.Call(args, env)
}

// Function interface for all callable functions
type Function interface {
	Value
	Call(args []Value, env *Environment) (Value, error)
}

// evalSpecialForm handles evaluation of special forms
// Returns (nil, nil) if the symbol is not a special form
func evalSpecialForm(name Symbol, args *List, env *Environment) (Value, error) {
	switch name {
	case Intern("quote"):
		return specialQuote(args, env)
	case Intern("if"):
		return specialIf(args, env)
	case Intern("fn"):
		return specialFn(args, env)
	case Intern("define"):
		return specialDefine(args, env)
	case Intern("do"):
		return specialDo(args, env)
	default:
		return nil, fmt.Errorf("not a special form") // Use a specific error
	}
}

func specialQuote(args *List, env *Environment) (Value, error) {
	if args.Length() != 1 {
		return nil, fmt.Errorf("quote requires exactly 1 argument")
	}
	return args.First(), nil
}

func specialIf(args *List, env *Environment) (Value, error) {
	if args.Length() < 2 || args.Length() > 3 {
		return nil, fmt.Errorf("if requires 2 or 3 arguments")
	}

	// Evaluate condition
	condition, err := Eval(args.First(), env)
	if err != nil {
		return nil, err
	}

	// Check if condition is truthy
	if isTruthy(condition) {
		return Eval(args.Rest().First(), env)
	} else if args.Length() == 3 {
		return Eval(args.Rest().Rest().First(), env)
	}

	// Return nil if no else branch
	return Nil{}, nil
}

func isTruthy(v Value) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(Boolean); ok && !bool(b) {
		return false
	}
	if _, ok := v.(Nil); ok {
		return false
	}
	return true
}

func specialFn(args *List, env *Environment) (Value, error) {
	if args.Length() != 2 {
		return nil, fmt.Errorf("fn requires exactly 2 arguments: (fn [params] body)")
	}

	// Get parameter list - accept either vector or list
	var params *List
	switch paramArg := args.First().(type) {
	case *Vector:
		// Convert vector to list for internal use
		params = paramArg.ToList()
	case *List:
		params = paramArg
	default:
		return nil, fmt.Errorf("fn parameter list must be a vector [params] or list (params)")
	}

	// Get function body
	body := args.Rest().First()

	// Create user function (closure)
	return &UserFunction{
		Params: params,
		Body:   body,
		Env:    env, // Capture current environment
	}, nil
}

func specialDefine(args *List, env *Environment) (Value, error) {
	if args.Length() != 2 {
		return nil, fmt.Errorf("define requires exactly 2 arguments")
	}

	// Get symbol name
	name, ok := args.First().(Symbol)
	if !ok {
		return nil, fmt.Errorf("define first argument must be a symbol")
	}

	// Evaluate value
	value, err := Eval(args.Rest().First(), env)
	if err != nil {
		return nil, err
	}

	// Define in current environment
	env.Set(name, value)

	return value, nil
}

func specialDo(args *List, env *Environment) (Value, error) {
	var result Value = Nil{}
	var err error

	// Evaluate each expression in sequence
	for current := args; !current.IsEmpty(); current = current.Rest() {
		result, err = Eval(current.First(), env)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// UserFunction represents a user-defined function (closure)
type UserFunction struct {
	Params *List
	Body   Value
	Env    *Environment
}

func (f *UserFunction) Call(args []Value, callEnv *Environment) (Value, error) {
	// Create new environment with closure's environment as parent
	env := NewEnvironment(f.Env)

	// Bind parameters to arguments
	if len(args) != f.Params.Length() {
		return nil, fmt.Errorf("function expects %d arguments, got %d", f.Params.Length(), len(args))
	}

	i := 0
	for current := f.Params; !current.IsEmpty(); current = current.Rest() {
		param, ok := current.First().(Symbol)
		if !ok {
			return nil, fmt.Errorf("parameter must be a symbol, got %T", current.First())
		}
		env.Set(param, args[i])
		i++
	}

	return Eval(f.Body, env)
}

func (f *UserFunction) String() string {
	return "<user-function>"
}
