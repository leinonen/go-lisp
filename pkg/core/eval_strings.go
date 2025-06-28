package core

import (
	"fmt"
	"strings"
)

// setupStringOperations adds string operations and predicates to the environment
func setupStringOperations(env *Environment) {
	// String operations
	env.Set(Intern("str"), &BuiltinFunction{
		Name: "str",
		Fn: func(args []Value, env *Environment) (Value, error) {
			result := ""
			for _, arg := range args {
				switch v := arg.(type) {
				case String:
					result += string(v)
				case Symbol:
					result += string(v)
				case Keyword:
					result += v.String()
				case Number:
					result += v.String()
				case Nil:
					result += ""
				default:
					result += arg.String()
				}
			}
			return String(result), nil
		},
	})

	env.Set(Intern("substring"), &BuiltinFunction{
		Name: "substring",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 2 || len(args) > 3 {
				return nil, fmt.Errorf("substring expects 2-3 arguments")
			}

			str, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("substring expects string as first argument")
			}

			start, ok := args[1].(Number)
			if !ok {
				return nil, fmt.Errorf("substring expects number as second argument")
			}

			s := string(str)
			startIdx := int(start.ToInt())

			if startIdx < 0 || startIdx > len(s) {
				return String(""), nil
			}

			if len(args) == 2 {
				return String(s[startIdx:]), nil
			}

			end, ok := args[2].(Number)
			if !ok {
				return nil, fmt.Errorf("substring expects number as third argument")
			}

			endIdx := int(end.ToInt())
			if endIdx < startIdx || endIdx > len(s) {
				endIdx = len(s)
			}

			return String(s[startIdx:endIdx]), nil
		},
	})

	env.Set(Intern("string-split"), &BuiltinFunction{
		Name: "string-split",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("string-split expects 2 arguments")
			}

			str, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("string-split expects string as first argument")
			}

			sep, ok := args[1].(String)
			if !ok {
				return nil, fmt.Errorf("string-split expects string as second argument")
			}

			parts := strings.Split(string(str), string(sep))
			elements := make([]Value, len(parts))
			for i, part := range parts {
				elements[i] = String(part)
			}

			return NewVector(elements...), nil
		},
	})

	env.Set(Intern("string-replace"), &BuiltinFunction{
		Name: "string-replace",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 3 {
				return nil, fmt.Errorf("string-replace expects 3 arguments")
			}

			str, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("string-replace expects string as first argument")
			}

			old, ok := args[1].(String)
			if !ok {
				return nil, fmt.Errorf("string-replace expects string as second argument")
			}

			new, ok := args[2].(String)
			if !ok {
				return nil, fmt.Errorf("string-replace expects string as third argument")
			}

			result := strings.ReplaceAll(string(str), string(old), string(new))
			return String(result), nil
		},
	})

	env.Set(Intern("string-contains?"), &BuiltinFunction{
		Name: "string-contains?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("string-contains? expects 2 arguments")
			}

			str, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("string-contains? expects string as first argument")
			}

			substr, ok := args[1].(String)
			if !ok {
				return nil, fmt.Errorf("string-contains? expects string as second argument")
			}

			if strings.Contains(string(str), string(substr)) {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("string-trim"), &BuiltinFunction{
		Name: "string-trim",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("string-trim expects 1 argument")
			}

			str, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("string-trim expects string")
			}

			result := strings.TrimSpace(string(str))
			return String(result), nil
		},
	})

	// String predicate
	env.Set(Intern("string?"), &BuiltinFunction{
		Name: "string?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("string? expects 1 argument")
			}

			if _, ok := args[0].(String); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})
}
