package core

import (
	"fmt"
	"sync/atomic"
)

// Global counter for gensym
var gensymCounter int64

// setupMetaProgramming adds meta-programming functions and type predicates to the environment
func setupMetaProgramming(env *Environment) {
	// Basic language literals
	env.Set(Intern("nil"), Nil{})
	env.Set(Intern("true"), Symbol("true"))
	env.Set(Intern("false"), Nil{})

	// Meta-programming functions
	env.Set(Intern("eval"), &BuiltinFunction{
		Name: "eval",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("eval expects 1 argument")
			}

			return Eval(args[0], env)
		},
	})

	env.Set(Intern("read-string"), &BuiltinFunction{
		Name: "read-string",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("read-string expects 1 argument")
			}

			str, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("read-string expects string, got %T", args[0])
			}

			return ReadString(string(str))
		},
	})

	env.Set(Intern("read-all-string"), &BuiltinFunction{
		Name: "read-all-string",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("read-all-string expects 1 argument")
			}

			str, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("read-all-string expects string, got %T", args[0])
			}

			lexer := NewLexer(string(str))
			tokens, err := lexer.Tokenize()
			if err != nil {
				return nil, fmt.Errorf("failed to tokenize: %v", err)
			}

			parser := NewParser(tokens)
			expressions, err := parser.ParseAll()
			if err != nil {
				return nil, fmt.Errorf("failed to parse: %v", err)
			}

			// Convert slice to List
			return NewList(expressions...), nil
		},
	})

	env.Set(Intern("throw"), &BuiltinFunction{
		Name: "throw",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("throw expects 1 argument")
			}

			// Convert the argument to a string for the error message
			var msg string
			if str, ok := args[0].(String); ok {
				msg = string(str)
			} else {
				msg = args[0].String()
			}

			return nil, fmt.Errorf("%s", msg)
		},
	})

	env.Set(Intern("gensym"), &BuiltinFunction{
		Name: "gensym",
		Fn: func(args []Value, env *Environment) (Value, error) {
			var prefix string
			if len(args) == 0 {
				prefix = "G__"
			} else if len(args) == 1 {
				if str, ok := args[0].(String); ok {
					prefix = string(str)
				} else if sym, ok := args[0].(Symbol); ok {
					prefix = string(sym)
				} else {
					return nil, fmt.Errorf("gensym expects string or symbol as prefix, got %T", args[0])
				}
			} else {
				return nil, fmt.Errorf("gensym expects 0 or 1 arguments, got %d", len(args))
			}

			// Atomically increment the counter to ensure uniqueness
			id := atomic.AddInt64(&gensymCounter, 1)
			return Symbol(fmt.Sprintf("%s%d", prefix, id)), nil
		},
	})

	env.Set(Intern("macroexpand"), &BuiltinFunction{
		Name: "macroexpand",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("macroexpand expects 1 argument")
			}

			return macroExpand(args[0], env)
		},
	})

	// Basic type predicates
	env.Set(Intern("symbol?"), &BuiltinFunction{
		Name: "symbol?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("symbol? expects 1 argument")
			}

			if _, ok := args[0].(Symbol); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("number?"), &BuiltinFunction{
		Name: "number?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("number? expects 1 argument")
			}

			if _, ok := args[0].(Number); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("keyword?"), &BuiltinFunction{
		Name: "keyword?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("keyword? expects 1 argument")
			}

			if _, ok := args[0].(Keyword); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("nil?"), &BuiltinFunction{
		Name: "nil?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("nil? expects 1 argument")
			}

			if _, ok := args[0].(Nil); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("fn?"), &BuiltinFunction{
		Name: "fn?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("fn? expects 1 argument")
			}

			if _, ok := args[0].(Function); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})
}

// macroExpand performs macro expansion on an expression
func macroExpand(expr Value, env *Environment) (Value, error) {
	// Only expand if expr is a list starting with a macro
	list, ok := expr.(*List)
	if !ok || list.IsEmpty() {
		return expr, nil
	}

	// Get the first element
	first := list.First()
	sym, ok := first.(Symbol)
	if !ok {
		return expr, nil
	}

	// Look up the symbol in the environment
	value, err := env.Get(sym)
	if err != nil {
		return expr, nil
	}

	// Check if it's a macro
	macro, ok := value.(*Macro)
	if !ok {
		return expr, nil
	}

	// Collect arguments for macro expansion
	args := listToSlice(list.Rest())

	// Create environment for macro expansion
	macroEnv := NewEnvironment(macro.Env)

	// Bind macro parameters to arguments
	err = bindParams(macro.Params, args, macroEnv)
	if err != nil {
		return nil, fmt.Errorf("macro expansion error: %v", err)
	}

	// Evaluate the macro body to get the expanded form
	expanded, err := Eval(macro.Body, macroEnv)
	if err != nil {
		return nil, fmt.Errorf("macro expansion error: %v", err)
	}

	return expanded, nil
}
