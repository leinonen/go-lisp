package kernel

// Core evaluation logic for the minimal Lisp kernel

import (
	"fmt"
	"os"
	"path/filepath"
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
			// Enhance undefined symbol errors with suggestions
			enhancedErr := enhanceUndefinedSymbolError(err, v.String(), env)
			return nil, ctx.WrapError(enhancedErr)
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
	case Intern("loop"):
		return specialLoopWithContext(args, env, ctx)
	case Intern("recur"):
		return specialRecurWithContext(args, env, ctx)
	default:
		return nil, fmt.Errorf("not a special form") // Use a specific error
	}
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

// LoadFile loads and evaluates a Lisp file
func LoadFile(filename string, env *Environment) (Value, error) {
	return loadFileWithContext(filename, env, NewEvaluationContext())
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

func specialQuoteWithContext(args *List, _ *Environment, ctx *EvaluationContext) (Value, error) {
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

	// Handle vectors
	if vector, ok := expr.(*Vector); ok {
		expanded := make([]Value, 0)
		for i := 0; i < vector.Length(); i++ {
			expandedElem, err := quasiExpandWithContext(vector.Get(i), env, ctx)
			if err != nil {
				return nil, err
			}
			expanded = append(expanded, expandedElem)
		}
		return NewVector(expanded...), nil
	}

	// For non-lists and non-vectors, return as-is (quoted)
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

	// Parameters must be a vector (Clojure-style)
	if vector, ok := params.(*Vector); ok {
		// Convert vector to list for internal consistency
		elements := make([]Value, 0)
		for i := 0; i < vector.Length(); i++ {
			elements = append(elements, vector.Get(i))
		}
		paramList = NewList(elements...)
	} else {
		return nil, ctx.CreateError("fn parameters must be a vector")
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

	// Return a closure with captured environment
	return &UserFunction{
		Params: paramList,
		Body:   body,
		Env:    env.CreateClosure(),
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

	// In a closure context, def creates/updates local bindings
	// This enables proper lexical closure behavior
	if env.parent != nil {
		// We're in a local environment (function, let, etc.)
		// Check if variable exists in current scope chain
		if env.HasBinding(symbol) {
			// Update existing binding
			err := env.Update(symbol, value)
			if err != nil {
				// If update fails, create new local binding
				env.SetLocal(symbol, value)
			}
		} else {
			// Create new local binding
			env.SetLocal(symbol, value)
		}
	} else {
		// We're in global environment, use global definition
		env.Set(symbol, value)
	}

	return DefinedValue{}, nil
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

	// Get parameters - must be a vector (Clojure-style)
	params := args.Rest().First()
	var paramList *List

	if vector, ok := params.(*Vector); ok {
		// Convert vector to list for internal consistency
		elements := make([]Value, 0)
		for i := 0; i < vector.Length(); i++ {
			elements = append(elements, vector.Get(i))
		}
		paramList = NewList(elements...)
	} else {
		return nil, ctx.CreateError("macro parameters must be a vector")
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
	return DefinedValue{}, nil
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
	// Try to find the file in different locations
	var content []byte
	var err error
	var actualPath string

	// First try the filename as-is (for absolute paths or relative to current directory)
	content, err = os.ReadFile(filename)
	if err != nil {
		// If that fails, try looking in the lisp directory relative to the working directory
		lispPath := filepath.Join("lisp", filename)
		content, err = os.ReadFile(lispPath)
		if err != nil {
			// Try relative to the executable location
			if execPath, execErr := os.Executable(); execErr == nil {
				execDir := filepath.Dir(execPath)
				lispPath = filepath.Join(execDir, "..", "lisp", filename)
				content, err = os.ReadFile(lispPath)
				if err != nil {
					// Finally, try relative to the Go module root
					if wd, wdErr := os.Getwd(); wdErr == nil {
						// Walk up directories to find go.mod
						for dir := wd; dir != "/" && dir != "."; dir = filepath.Dir(dir) {
							if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
								lispPath = filepath.Join(dir, "lisp", filename)
								content, err = os.ReadFile(lispPath)
								if err == nil {
									actualPath = lispPath
									break
								}
							}
						}
					}
					if err != nil {
						return nil, ctx.CreateError(fmt.Sprintf("could not read file %s: %v", filename, err))
					}
				} else {
					actualPath = lispPath
				}
			} else {
				return nil, ctx.CreateError(fmt.Sprintf("could not read file %s: %v", filename, err))
			}
		} else {
			actualPath = lispPath
		}
	} else {
		actualPath = filename
	}

	// Remove comments for cleaner parsing
	cleanedContent := removeComments(string(content))

	// Create a lexer for the file content - use actual path for better error reporting
	lexer := NewLexer(cleanedContent, actualPath)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	// Create parser
	parser := NewParser(tokens, actualPath)

	ctx.PushFrame(fmt.Sprintf("evaluating file: %s", actualPath))
	defer ctx.PopFrame()

	// Parse and evaluate multiple expressions from the file
	for parser.position < len(parser.tokens) {
		// Skip EOF tokens
		if parser.position < len(parser.tokens) && parser.tokens[parser.position].Type == TokenEOF {
			break
		}

		// Parse one expression
		expr, pos, err := parser.Parse()
		if err != nil {
			return nil, err
		}

		// Set source location in context
		if pos != nil {
			ctx.SetLocation(pos.Line, pos.Column, pos.File)
		}

		// Evaluate the expression - we don't care about the result
		_, err = EvalWithContext(expr, env, ctx)
		if err != nil {
			return nil, err
		}
	}

	// Always return nil after successfully loading a file
	return Nil{}, nil
}

// Enhanced error handling functions

// getSimilarSymbols finds symbols similar to the given symbol name
func getSimilarSymbols(symbolName string, env *Environment) []string {
	var similar []string

	// Get all defined symbols in the environment
	bindings := env.GetAllBindings()

	for name := range bindings {
		nameStr := name.String()
		// Simple similarity check - same length or starts with same characters
		if len(nameStr) > 0 && (len(nameStr) == len(symbolName) ||
			strings.HasPrefix(nameStr, symbolName[:min(2, len(symbolName))]) ||
			strings.HasPrefix(symbolName, nameStr[:min(2, len(nameStr))])) {
			similar = append(similar, nameStr)
		}
	}

	return similar
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// enhanceUndefinedSymbolError adds suggestions to undefined symbol errors
func enhanceUndefinedSymbolError(err error, symbolName string, env *Environment) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "undefined symbol") {
		return err
	}

	similar := getSimilarSymbols(symbolName, env)
	if len(similar) > 0 {
		suggestions := strings.Join(similar[:min(5, len(similar))], ", ")
		return fmt.Errorf("%s\nDid you mean one of: %s", errMsg, suggestions)
	}

	// Show available functions from environment
	bindings := env.GetAllBindings()
	var functions []string
	for name, value := range bindings {
		if _, ok := value.(*UserFunction); ok {
			functions = append(functions, name.String())
		}
	}

	if len(functions) > 0 {
		available := strings.Join(functions[:min(10, len(functions))], ", ")
		return fmt.Errorf("%s\nAvailable functions: %s", errMsg, available)
	}

	return err
}

// specialLoopWithContext implements the loop special form
// Syntax: (loop [bindings...] body...)
func specialLoopWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	if args.Length() < 2 {
		return nil, ctx.CreateError("loop requires at least 2 arguments: (loop [bindings...] body...)")
	}

	// First argument should be a vector of bindings
	bindingsArg := args.First()
	bindingsVector, ok := bindingsArg.(*Vector)
	if !ok {
		return nil, ctx.CreateError("loop bindings must be a vector")
	}

	// Parse bindings - should be [name value name value ...]
	if bindingsVector.Length()%2 != 0 {
		return nil, ctx.CreateError("loop bindings must be pairs of [name value ...]")
	}

	// Extract binding names and initial values
	var bindingNames []Symbol
	var initialValues []Value

	for i := 0; i < bindingsVector.Length(); i += 2 {
		name, ok := bindingsVector.Get(i).(Symbol)
		if !ok {
			return nil, ctx.CreateError("loop binding names must be symbols")
		}
		bindingNames = append(bindingNames, name)

		// Evaluate initial value in current environment
		initValue, err := EvalWithContext(bindingsVector.Get(i+1), env, ctx)
		if err != nil {
			return nil, err
		}
		initialValues = append(initialValues, initValue)
	}

	// Get body expressions
	bodyArgs := args.Rest()
	if bodyArgs.IsEmpty() {
		return nil, ctx.CreateError("loop requires a body")
	}

	// Create new environment for loop
	loopEnv := NewEnvironment(env)

	// Set up loop context
	loopEnv.SetLoopContext(bindingNames)

	// Execute loop with tail call optimization
	return executeLoop(bindingNames, initialValues, bodyArgs, loopEnv, ctx)
}

// executeLoop implements the actual loop execution with tail call optimization
func executeLoop(bindingNames []Symbol, values []Value, body *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	for {
		// Bind current values to loop variables
		for i, name := range bindingNames {
			env.Set(name, values[i])
		}

		// Execute body expressions
		var result Value
		var err error

		for current := body; !current.IsEmpty(); current = current.Rest() {
			result, err = EvalWithContext(current.First(), env, ctx)
			if err != nil {
				return nil, err
			}

			// Check if this is a recur call
			if recurVal, ok := result.(*RecurValue); ok {
				// Validate argument count
				if len(recurVal.Values) != len(bindingNames) {
					return nil, ctx.CreateError(fmt.Sprintf("recur expects %d arguments, got %d",
						len(bindingNames), len(recurVal.Values)))
				}

				// Update values for next iteration
				values = recurVal.Values
				break // Break out of body evaluation and restart loop
			}
		}

		// If we didn't encounter a recur, return the result
		if _, ok := result.(*RecurValue); !ok {
			return result, nil
		}
	}
}

// specialRecurWithContext implements the recur special form
// Syntax: (recur args...)
func specialRecurWithContext(args *List, env *Environment, ctx *EvaluationContext) (Value, error) {
	// Check if we're in a loop context
	loopContext := env.GetLoopContext()
	if loopContext == nil {
		return nil, ctx.CreateError("recur can only be used inside a loop")
	}

	// Evaluate all arguments
	var values []Value
	for current := args; !current.IsEmpty(); current = current.Rest() {
		value, err := EvalWithContext(current.First(), env, ctx)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	// Return a RecurValue which will be handled by the loop
	return &RecurValue{Values: values}, nil
}
