package minimal

// Bootstrap demonstrates extending the minimal kernel with higher-level constructs
// This shows how to implement features in Lisp itself, following the architecture in future.md

import "fmt"

// Bootstrap loads higher-level constructs into the environment
func Bootstrap(env *Environment) error {
	// Define 'when' macro-like function using core 'if'
	// (define when (fn [condition & body] (if condition (do body) nil)))
	whenParams := NewVector(Intern("condition"), Intern("body")).ToList()
	whenBody := NewList(
		Intern("if"),
		Intern("condition"),
		NewList(Intern("do"), Intern("body")),
		Nil{},
	)
	whenFn := &UserFunction{
		Params: whenParams,
		Body:   whenBody,
		Env:    env,
	}
	env.Set(Intern("when"), whenFn)

	// Define 'unless' using 'if'
	// (define unless (fn [condition body] (if condition nil body)))
	unlessParams := NewVector(Intern("condition"), Intern("body")).ToList()
	unlessBody := NewList(
		Intern("if"),
		Intern("condition"),
		Nil{},
		Intern("body"),
	)
	unlessFn := &UserFunction{
		Params: unlessParams,
		Body:   unlessBody,
		Env:    env,
	}
	env.Set(Intern("unless"), unlessFn)

	// Define 'list' function to create lists dynamically
	env.Set(Intern("list"), &BuiltinFunction{
		Name: "list",
		Fn: func(args []Value, env *Environment) (Value, error) {
			return NewList(args...), nil
		},
	})

	// Define 'first' function
	env.Set(Intern("first"), &BuiltinFunction{
		Name: "first",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("first requires exactly 1 argument")
			}

			if list, ok := args[0].(*List); ok {
				if list.IsEmpty() {
					return Nil{}, nil
				}
				return list.First(), nil
			}

			return nil, fmt.Errorf("first requires a list, got %T", args[0])
		},
	})

	// Define 'rest' function
	env.Set(Intern("rest"), &BuiltinFunction{
		Name: "rest",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("rest requires exactly 1 argument")
			}

			if list, ok := args[0].(*List); ok {
				return list.Rest(), nil
			}

			return nil, fmt.Errorf("rest requires a list, got %T", args[0])
		},
	})

	// Define comparison operators
	env.Set(Intern("="), &BuiltinFunction{
		Name: "=",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("= requires exactly 2 arguments")
			}

			// Simple equality check
			return Boolean(args[0].String() == args[1].String()), nil
		},
	})

	env.Set(Intern("<"), &BuiltinFunction{
		Name: "<",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("< requires exactly 2 arguments")
			}

			num1, ok1 := args[0].(Number)
			num2, ok2 := args[1].(Number)

			if !ok1 || !ok2 {
				return nil, fmt.Errorf("< requires numbers")
			}

			return Boolean(float64(num1) < float64(num2)), nil
		},
	})

	// Print function for output
	env.Set(Intern("print"), &BuiltinFunction{
		Name: "print",
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

	return nil
}

// NewBootstrappedREPL creates a REPL with bootstrapped functions
func NewBootstrappedREPL() *REPL {
	repl := NewREPL()
	Bootstrap(repl.Env)
	return repl
}
