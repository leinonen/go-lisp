package core

import "fmt"

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
