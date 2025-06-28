package core

import (
	"fmt"
	"os"
)

// setupIOOperations adds I/O and file operations to the environment
func setupIOOperations(env *Environment) {
	// Console I/O
	env.Set(Intern("println"), &BuiltinFunction{
		Name: "println",
		Fn: func(args []Value, env *Environment) (Value, error) {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				switch v := arg.(type) {
				case String:
					fmt.Print(string(v))
				case Symbol:
					fmt.Print(string(v))
				case Keyword:
					fmt.Print(v.String())
				case Number:
					fmt.Print(v.String())
				case Nil:
					fmt.Print("nil")
				default:
					fmt.Print(arg.String())
				}
			}
			fmt.Println()
			return Nil{}, nil
		},
	})

	env.Set(Intern("prn"), &BuiltinFunction{
		Name: "prn",
		Fn: func(args []Value, env *Environment) (Value, error) {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(arg.String())
			}
			fmt.Println()
			return Nil{}, nil
		},
	})

	// File I/O
	env.Set(Intern("slurp"), &BuiltinFunction{
		Name: "slurp",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("slurp expects 1 argument")
			}

			filename, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("slurp expects string, got %T", args[0])
			}

			content, err := os.ReadFile(string(filename))
			if err != nil {
				return nil, fmt.Errorf("slurp error: %v", err)
			}

			return String(content), nil
		},
	})

	env.Set(Intern("spit"), &BuiltinFunction{
		Name: "spit",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("spit expects 2 arguments")
			}

			filename, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("spit expects string as first argument, got %T", args[0])
			}

			content, ok := args[1].(String)
			if !ok {
				return nil, fmt.Errorf("spit expects string as second argument, got %T", args[1])
			}

			err := os.WriteFile(string(filename), []byte(content), 0644)
			if err != nil {
				return nil, fmt.Errorf("spit error: %v", err)
			}

			return String(filename), nil
		},
	})

	// File system operations
	env.Set(Intern("file-exists?"), &BuiltinFunction{
		Name: "file-exists?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("file-exists? expects 1 argument")
			}

			filename, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("file-exists? expects string, got %T", args[0])
			}

			if _, err := os.Stat(string(filename)); err == nil {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("list-dir"), &BuiltinFunction{
		Name: "list-dir",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("list-dir expects 1 argument")
			}

			dirname, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("list-dir expects string, got %T", args[0])
			}

			entries, err := os.ReadDir(string(dirname))
			if err != nil {
				return nil, fmt.Errorf("list-dir error: %v", err)
			}

			elements := make([]Value, len(entries))
			for i, entry := range entries {
				elements[i] = String(entry.Name())
			}

			return NewVector(elements...), nil
		},
	})
}
