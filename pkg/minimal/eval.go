package minimal

// Core evaluation logic for the minimal Lisp kernel

import (
	"fmt"
	"os"
	"strings"
)

// Eval evaluates a Lisp expression in the given environment
func Eval(expr Value, env *Environment) (Value, error) {
	return EvalWithContext(expr, env, NewEvaluationContext())
}

// EvalWithContext evaluates a Lisp expression with evaluation context for better error reporting
func EvalWithContext(expr Value, env *Environment, ctx *EvaluationContext) (Value, error) {
	// Set current expression for error reporting
	ctx.SetExpression(expr.String())

	switch v := expr.(type) {
	case Symbol:
		// Look up symbol in environment
		ctx.PushFrame(fmt.Sprintf("looking up symbol: %s", v))
		result, err := env.Get(v)
		ctx.PopFrame()

		if err != nil {
			return nil, ctx.WrapError(err)
		}
		return result, nil

	case *List:
		if v.IsEmpty() {
			return v, nil // Empty list evaluates to itself
		}

		// Check if first element is a special form
		if sym, ok := v.First().(Symbol); ok {
			ctx.PushFrame(fmt.Sprintf("evaluating special form: %s", sym))
			result, err := evalSpecialFormWithContext(sym, v.Rest(), env, ctx)
			ctx.PopFrame()

			if err == nil {
				return result, nil
			}
			// Not a special form - check if the error is our specific "not a special form" error
			if err.Error() != "not a special form" {
				return nil, ctx.WrapError(err) // Return actual errors from special forms
			}
			// Fall through to function application
		}

		// Regular function application
		ctx.PushFrame(fmt.Sprintf("applying function: %s", v.First()))
		result, err := ApplyWithContext(v, env, ctx)
		ctx.PopFrame()

		if err != nil {
			return nil, ctx.WrapError(err)
		}
		return result, nil

	default:
		// Self-evaluating expressions (numbers, strings, booleans, etc.)
		return expr, nil
	}
}

// Apply applies a function to arguments
func Apply(list *List, env *Environment) (Value, error) {
	return ApplyWithContext(list, env, NewEvaluationContext())
}

// ApplyWithContext applies a function to arguments with evaluation context
func ApplyWithContext(list *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if list.IsEmpty() {
		return nil, ctx.CreateError("cannot apply empty list")
	}

	// Evaluate the function
	ctx.PushFrame("evaluating function")
	fn, err := EvalWithContext(list.First(), env, ctx)
	ctx.PopFrame()

	if err != nil {
		return nil, err
	}

	// Check if it's a macro - macros get unevaluated arguments
	if macro, ok := fn.(*Macro); ok {
		ctx.PushFrame("expanding macro")
		// Collect unevaluated arguments for macro
		args := make([]Value, 0)
		for current := list.Rest(); !current.IsEmpty(); current = current.Rest() {
			args = append(args, current.First())
		}
		result, err := macro.CallWithContext(args, env, ctx)
		ctx.PopFrame()
		return result, err
	}

	function, ok := fn.(Function)
	if !ok {
		return nil, ctx.CreateError(fmt.Sprintf("%v is not a function", fn))
	}

	// Evaluate all arguments for regular functions
	ctx.PushFrame("evaluating arguments")
	args := make([]Value, 0)
	for current := list.Rest(); !current.IsEmpty(); current = current.Rest() {
		evaluated, err := EvalWithContext(current.First(), env, ctx)
		if err != nil {
			ctx.PopFrame()
			return nil, err
		}
		args = append(args, evaluated)
	}
	ctx.PopFrame()

	// Call the function
	ctx.PushFrame(fmt.Sprintf("calling function with %d arguments", len(args)))
	result, err := function.CallWithContext(args, env, ctx)
	ctx.PopFrame()

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Function interface for all callable functions
type Function interface {
	Value
	Call(args []Value, env *Environment) (Value, error)
	CallWithContext(args []Value, env *Environment, ctx *EvaluationContext) (Value, error)
}

// evalSpecialForm handles evaluation of special forms
// Returns (nil, nil) if the symbol is not a special form
func evalSpecialForm(name Symbol, args *List, env *Environment) (Value, error) {
	return evalSpecialFormWithContext(name, args, env, NewEvaluationContext())
}

// evalSpecialFormWithContext handles evaluation of special forms with context
func evalSpecialFormWithContext(name Symbol, args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	switch name {
	case Intern("quote"):
		return specialQuoteWithContext(args, env, ctx)
	case Intern("quasiquote"):
		return specialQuasiquoteWithContext(args, env, ctx)
	case Intern("unquote"):
		return nil, ctx.CreateError("unquote not inside quasiquote")
	case Intern("if"):
		return specialIfWithContext(args, env, ctx)
	case Intern("fn"):
		return specialFnWithContext(args, env, ctx)
	case Intern("def"):
		return specialDefineWithContext(args, env, ctx)
	case Intern("defmacro"):
		return specialDefmacroWithContext(args, env, ctx)
	case Intern("load"):
		return specialLoadWithContext(args, env, ctx)
	case Intern("do"):
		return specialDoWithContext(args, env, ctx)
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
		return nil, fmt.Errorf("def requires exactly 2 arguments")
	}

	// Get symbol name
	name, ok := args.First().(Symbol)
	if !ok {
		return nil, fmt.Errorf("def first argument must be a symbol")
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

// CallWithContext calls the function with evaluation context
func (f *UserFunction) CallWithContext(args []Value, callEnv *Environment, ctx *EvaluationContext) (Value, error) {
	// Create new environment with closure's environment as parent
	env := NewEnvironment(f.Env)

	// Bind parameters to arguments
	if len(args) != f.Params.Length() {
		return nil, ctx.CreateError(fmt.Sprintf("function expects %d arguments, got %d", f.Params.Length(), len(args)))
	}

	i := 0
	for current := f.Params; !current.IsEmpty(); current = current.Rest() {
		param, ok := current.First().(Symbol)
		if !ok {
			return nil, ctx.CreateError(fmt.Sprintf("parameter must be a symbol, got %T", current.First()))
		}
		env.Set(param, args[i])
		i++
	}

	ctx.PushFrame("executing user function body")
	result, err := EvalWithContext(f.Body, env, ctx)
	ctx.PopFrame()

	return result, err
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

// CallWithContext calls the macro with evaluation context
func (m *Macro) CallWithContext(args []Value, callEnv *Environment, ctx *EvaluationContext) (Value, error) {
	// Create new environment with macro's environment as parent
	env := NewEnvironment(m.Env)

	// Bind parameters to arguments (unevaluated)
	if len(args) != m.Params.Length() {
		return nil, ctx.CreateError(fmt.Sprintf("macro expects %d arguments, got %d", m.Params.Length(), len(args)))
	}

	i := 0
	for current := m.Params; !current.IsEmpty(); current = current.Rest() {
		param, ok := current.First().(Symbol)
		if !ok {
			return nil, ctx.CreateError(fmt.Sprintf("parameter must be a symbol, got %T", current.First()))
		}
		env.Set(param, args[i])
		i++
	}

	// Evaluate macro body to get the expansion
	ctx.PushFrame("expanding macro body")
	expansion, err := EvalWithContext(m.Body, env, ctx)
	ctx.PopFrame()

	if err != nil {
		return nil, err
	}

	// Evaluate the expansion in the calling environment
	ctx.PushFrame("evaluating macro expansion")
	result, err := EvalWithContext(expansion, callEnv, ctx)
	ctx.PopFrame()

	return result, err
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

// Context-aware versions of special forms for enhanced error handling

func specialQuoteWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if args.Length() != 1 {
		return nil, ctx.CreateError("quote requires exactly 1 argument")
	}
	return args.First(), nil
}

func specialQuasiquoteWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if args.Length() != 1 {
		return nil, ctx.CreateError("quasiquote requires exactly 1 argument")
	}
	return quasiExpandWithContext(args.First(), env, ctx)
}

// quasiExpandWithContext expands a quasiquoted expression with context
func quasiExpandWithContext(expr Value, env *Environment, ctx *EvaluationContext) (Value, error) {
	// Check if this is an unquote
	if list, ok := expr.(*List); ok && !list.IsEmpty() {
		if sym, ok := list.First().(Symbol); ok && sym == Intern("unquote") {
			if list.Length() != 2 {
				return nil, ctx.CreateError("unquote requires exactly 1 argument")
			}
			// Evaluate the unquoted expression
			ctx.PushFrame("evaluating unquoted expression")
			result, err := EvalWithContext(list.Rest().First(), env, ctx)
			ctx.PopFrame()
			return result, err
		}

		// Recursively expand list elements
		expanded := make([]Value, 0)
		for current := list; !current.IsEmpty(); current = current.Rest() {
			expandedElem, err := quasiExpandWithContext(current.First(), env, ctx)
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

func specialIfWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if args.Length() < 2 || args.Length() > 3 {
		return nil, ctx.CreateError("if requires 2 or 3 arguments")
	}

	// Evaluate condition
	ctx.PushFrame("evaluating if condition")
	condition, err := EvalWithContext(args.First(), env, ctx)
	ctx.PopFrame()

	if err != nil {
		return nil, err
	}

	// Check if condition is truthy
	if isTruthy(condition) {
		ctx.PushFrame("evaluating if then-branch")
		result, err := EvalWithContext(args.Rest().First(), env, ctx)
		ctx.PopFrame()
		return result, err
	} else if args.Length() == 3 {
		ctx.PushFrame("evaluating if else-branch")
		result, err := EvalWithContext(args.Rest().Rest().First(), env, ctx)
		ctx.PopFrame()
		return result, err
	}

	// Return nil if no else branch
	return Nil{}, nil
}

func specialFnWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if args.Length() < 2 {
		return nil, ctx.CreateError("fn requires at least 2 arguments: (fn [params] body...)")
	}

	params := args.First()
	var paramList *List

	// Accept both lists and vectors for parameters
	if list, ok := params.(*List); ok {
		paramList = list
	} else if vector, ok := params.(*Vector); ok {
		// Convert vector to list for consistency
		elements := make([]Value, 0)
		for i := 0; i < vector.Length(); i++ {
			elements = append(elements, vector.Get(i))
		}
		paramList = NewList(elements...)
	} else {
		return nil, ctx.CreateError("fn parameters must be a list or vector")
	}

	bodyArgs := args.Rest()
	var body Value

	if bodyArgs.Length() == 1 {
		body = bodyArgs.First()
	} else {
		// Multiple body expressions - wrap in 'do'
		bodyElements := []Value{Intern("do")}
		for current := bodyArgs; !current.IsEmpty(); current = current.Rest() {
			bodyElements = append(bodyElements, current.First())
		}
		body = NewList(bodyElements...)
	}

	// Return a closure
	return &UserFunction{
		Params: paramList,
		Body:   body,
		Env:    env,
	}, nil
}

func specialDefineWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if args.Length() != 2 {
		return nil, ctx.CreateError("def requires exactly 2 arguments")
	}

	symbol, ok := args.First().(Symbol)
	if !ok {
		return nil, ctx.CreateError("first argument to def must be a symbol")
	}

	ctx.PushFrame(fmt.Sprintf("evaluating definition of %s", symbol))
	value, err := EvalWithContext(args.Rest().First(), env, ctx)
	ctx.PopFrame()

	if err != nil {
		return nil, err
	}

	env.Set(symbol, value)
	return Intern("defined"), nil
}

func specialDefmacroWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if args.Length() != 3 {
		return nil, ctx.CreateError("defmacro requires exactly 3 arguments: (defmacro name [params] body)")
	}

	// Get macro name
	name, ok := args.First().(Symbol)
	if !ok {
		return nil, ctx.CreateError("macro name must be a symbol")
	}

	// Get parameters - accept both lists and vectors
	params := args.Rest().First()
	var paramList *List

	if list, ok := params.(*List); ok {
		paramList = list
	} else if vector, ok := params.(*Vector); ok {
		// Convert vector to list for consistency
		elements := make([]Value, 0)
		for i := 0; i < vector.Length(); i++ {
			elements = append(elements, vector.Get(i))
		}
		paramList = NewList(elements...)
	} else {
		return nil, ctx.CreateError("macro parameters must be a list or vector")
	}

	// Get body
	body := args.Rest().Rest().First()

	// Create macro
	macro := &Macro{
		Params: paramList,
		Body:   body,
		Env:    env,
	}

	// Define in environment
	env.Set(name, macro)
	return Intern("defined"), nil
}

func specialDoWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if args.Length() == 0 {
		return Nil{}, nil
	}

	var result Value = Nil{}
	for current := args; !current.IsEmpty(); current = current.Rest() {
		ctx.PushFrame("evaluating do expression")
		val, err := EvalWithContext(current.First(), env, ctx)
		ctx.PopFrame()

		if err != nil {
			return nil, err
		}
		result = val
	}

	return result, nil
}

func specialLoadWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if args.Length() != 1 {
		return nil, ctx.CreateError("load requires exactly 1 argument")
	}

	// Evaluate the filename argument
	ctx.PushFrame("evaluating load filename")
	filenameValue, err := EvalWithContext(args.First(), env, ctx)
	ctx.PopFrame()

	if err != nil {
		return nil, err
	}

	filename, ok := filenameValue.(String)
	if !ok {
		return nil, ctx.CreateError("load argument must be a string")
	}

	// Set filename in context for better error reporting
	oldFilename := ctx.Filename
	ctx.Filename = string(filename)
	defer func() { ctx.Filename = oldFilename }()

	ctx.PushFrame(fmt.Sprintf("loading file: %s", filename))
	result, err := loadFileWithContext(string(filename), env, ctx)
	ctx.PopFrame()

	return result, err
}

// loadFileWithContext loads and evaluates a Lisp file with context
func loadFileWithContext(filename string, env *Environment, ctx *EvaluationContext) (Value, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, ctx.CreateError(fmt.Sprintf("could not read file %s: %v", filename, err))
	}

	// Parse with position tracking
	expr, pos, err := ParseWithPositions(string(content), filename)
	if err != nil {
		return nil, err
	}

	// Set source location in context
	if pos != nil {
		ctx.SetLocation(pos.Line, pos.Column, pos.File)
	}

	ctx.PushFrame(fmt.Sprintf("evaluating file contents: %s", filename))
	result, err := EvalWithContext(expr, env, ctx)
	ctx.PopFrame()

	return result, err
}
