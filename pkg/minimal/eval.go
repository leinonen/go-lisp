package minimal

// Core evaluation logic for the minimal Lisp kernel

import (
	"fmt"
	"os"
	"strings"
)

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

	// Check if it's a macro - macros get unevaluated arguments
	if macro, ok := fn.(*Macro); ok {
		// Collect unevaluated arguments for macro
		args := make([]Value, 0)
		for current := list.Rest(); !current.IsEmpty(); current = current.Rest() {
			args = append(args, current.First())
		}
		return macro.Call(args, env)
	}

	function, ok := fn.(Function)
	if !ok {
		return nil, fmt.Errorf("%v is not a function", fn)
	}

	// Evaluate all arguments for regular functions
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
	case Intern("quasiquote"):
		return specialQuasiquote(args, env)
	case Intern("unquote"):
		return nil, fmt.Errorf("unquote not inside quasiquote")
	case Intern("if"):
		return specialIf(args, env)
	case Intern("fn"):
		return specialFn(args, env)
	case Intern("define"):
		return specialDefine(args, env)
	case Intern("defmacro"):
		return specialDefmacro(args, env)
	case Intern("load"):
		return specialLoad(args, env)
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

func specialQuasiquote(args *List, env *Environment) (Value, error) {
	if args.Length() != 1 {
		return nil, fmt.Errorf("quasiquote requires exactly 1 argument")
	}
	return quasiExpand(args.First(), env)
}

// quasiExpand expands a quasiquoted expression
func quasiExpand(expr Value, env *Environment) (Value, error) {
	// Check if this is an unquote
	if list, ok := expr.(*List); ok && !list.IsEmpty() {
		if sym, ok := list.First().(Symbol); ok && sym == Intern("unquote") {
			if list.Length() != 2 {
				return nil, fmt.Errorf("unquote requires exactly 1 argument")
			}
			// Evaluate the unquoted expression
			return Eval(list.Rest().First(), env)
		}

		// Recursively expand list elements
		expanded := make([]Value, 0)
		for current := list; !current.IsEmpty(); current = current.Rest() {
			expandedElem, err := quasiExpand(current.First(), env)
			if err != nil {
				return nil, err
			}
			expanded = append(expanded, expandedElem)
		}
		return NewList(expanded...), nil
	}

	// For non-lists, return as-is (quoted)
	return expr, nil
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
	if args.Length() < 2 {
		return nil, fmt.Errorf("fn requires at least 2 arguments: (fn [params] body...)")
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

	// Get function body - if multiple expressions, wrap in 'do'
	var body Value
	bodyArgs := args.Rest()
	if bodyArgs.Length() == 1 {
		body = bodyArgs.First()
	} else {
		// Multiple body expressions - wrap in 'do'
		bodyElements := make([]Value, 0, bodyArgs.Length()+1)
		bodyElements = append(bodyElements, Intern("do"))
		for current := bodyArgs; !current.IsEmpty(); current = current.Rest() {
			bodyElements = append(bodyElements, current.First())
		}
		body = NewList(bodyElements...)
	}

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

	return DefinedValue{}, nil
}

// specialDefmacro implements defmacro functionality
func specialDefmacro(args *List, env *Environment) (Value, error) {
	if args.Length() != 3 {
		return nil, fmt.Errorf("defmacro requires exactly 3 arguments: name, parameters, body")
	}

	// Get macro name
	name, ok := args.First().(Symbol)
	if !ok {
		return nil, fmt.Errorf("macro name must be a symbol, got %T", args.First())
	}

	// Get parameters (must be a vector or list)
	var params *List
	paramsArg := args.Rest().First()
	if vec, ok := paramsArg.(*Vector); ok {
		params = vec.ToList()
	} else if list, ok := paramsArg.(*List); ok {
		params = list
	} else {
		return nil, fmt.Errorf("macro parameters must be a vector or list, got %T", paramsArg)
	}

	// Get body
	body := args.Rest().Rest().First()

	// Create macro
	macro := &Macro{
		Params: params,
		Body:   body,
		Env:    env,
	}

	// Store in environment
	env.Set(name, macro)

	return Intern("defined"), nil
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

// Macro represents a user-defined macro
type Macro struct {
	Params *List
	Body   Value
	Env    *Environment
}

func (m *Macro) Call(args []Value, callEnv *Environment) (Value, error) {
	// Create new environment with macro's environment as parent
	env := NewEnvironment(m.Env)

	// Bind parameters to arguments (unevaluated)
	if len(args) != m.Params.Length() {
		return nil, fmt.Errorf("macro expects %d arguments, got %d", m.Params.Length(), len(args))
	}

	i := 0
	for current := m.Params; !current.IsEmpty(); current = current.Rest() {
		param, ok := current.First().(Symbol)
		if !ok {
			return nil, fmt.Errorf("parameter must be a symbol, got %T", current.First())
		}
		env.Set(param, args[i])
		i++
	}

	// Evaluate macro body to get the expansion
	expansion, err := Eval(m.Body, env)
	if err != nil {
		return nil, err
	}

	// Evaluate the expansion in the calling environment
	return Eval(expansion, callEnv)
}

func (m *Macro) String() string {
	return "<macro>"
}

// specialLoad implements file loading
func specialLoad(args *List, env *Environment) (Value, error) {
	if args.Length() != 1 {
		return nil, fmt.Errorf("load requires exactly 1 argument (filename)")
	}

	filename, ok := args.First().(String)
	if !ok {
		return nil, fmt.Errorf("load argument must be a string, got %T", args.First())
	}

	return LoadFile(string(filename), env)
}

// LoadFile loads and evaluates a Lisp file
func LoadFile(filename string, env *Environment) (Value, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	// Create a temporary REPL for parsing
	repl := &REPL{Env: env}

	// Remove comments and clean up the content
	cleanContent := removeComments(string(content))

	// Parse the entire content as a sequence of expressions
	tokens := repl.tokenize(cleanContent)
	var lastResult Value = Nil{}
	pos := 0

	for pos < len(tokens) {
		// Skip empty tokens
		if tokens[pos] == "" {
			pos++
			continue
		}

		expr, newPos, err := repl.parseExpression(tokens, pos)
		if err != nil {
			return nil, fmt.Errorf("parse error in %s: %v", filename, err)
		}

		result, err := Eval(expr, env)
		if err != nil {
			return nil, fmt.Errorf("eval error in %s: %v", filename, err)
		}

		lastResult = result
		pos = newPos
	}

	return lastResult, nil
}

// removeComments removes Lisp comments from source code
func removeComments(content string) string {
	lines := strings.Split(content, "\n")
	var cleanLines []string

	for _, line := range lines {
		// Find comment start, but be careful about strings
		cleanLine := ""
		inString := false
		i := 0

		for i < len(line) {
			char := line[i]
			if char == '"' && (i == 0 || line[i-1] != '\\') {
				inString = !inString
				cleanLine += string(char)
			} else if char == ';' && !inString {
				// Found comment start outside of string
				break
			} else {
				cleanLine += string(char)
			}
			i++
		}

		cleanLine = strings.TrimSpace(cleanLine)
		if cleanLine != "" {
			cleanLines = append(cleanLines, cleanLine)
		}
	}

	return strings.Join(cleanLines, " ")
}
