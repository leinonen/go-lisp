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

	// Vector operations
	env.Set(Intern("vector-get"), &BuiltinFunction{
		Name: "vector-get",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("vector-get expects 2 arguments, got %d", len(args))
			}

			vec, ok := args[0].(*Vector)
			if !ok {
				return nil, fmt.Errorf("first argument to vector-get must be a vector, got %T", args[0])
			}

			indexVal, ok := args[1].(Number)
			if !ok {
				return nil, fmt.Errorf("second argument to vector-get must be a number, got %T", args[1])
			}

			return vec.Get(int(indexVal)), nil
		},
	})

	env.Set(Intern("vector-append"), &BuiltinFunction{
		Name: "vector-append",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("vector-append expects 2 arguments, got %d", len(args))
			}

			vec, ok := args[0].(*Vector)
			if !ok {
				return nil, fmt.Errorf("first argument to vector-append must be a vector, got %T", args[0])
			}

			return vec.Append(args[1]), nil
		},
	})

	env.Set(Intern("vector-update"), &BuiltinFunction{
		Name: "vector-update",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 3 {
				return nil, fmt.Errorf("vector-update expects 3 arguments, got %d", len(args))
			}

			vec, ok := args[0].(*Vector)
			if !ok {
				return nil, fmt.Errorf("first argument to vector-update must be a vector, got %T", args[0])
			}

			indexVal, ok := args[1].(Number)
			if !ok {
				return nil, fmt.Errorf("second argument to vector-update must be a number, got %T", args[1])
			}

			return vec.Update(int(indexVal), args[2]), nil
		},
	})

	// HashMap operations
	env.Set(Intern("hash-map"), &BuiltinFunction{
		Name: "hash-map",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args)%2 != 0 {
				return nil, fmt.Errorf("hash-map expects even number of arguments (key-value pairs), got %d", len(args))
			}

			hm := NewHashMap()
			for i := 0; i < len(args); i += 2 {
				key := args[i].String()
				val := args[i+1]
				hm = hm.Put(key, val)
			}

			return hm, nil
		},
	})

	env.Set(Intern("hash-map-get"), &BuiltinFunction{
		Name: "hash-map-get",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("hash-map-get expects 2 arguments, got %d", len(args))
			}

			hm, ok := args[0].(*HashMap)
			if !ok {
				return nil, fmt.Errorf("first argument to hash-map-get must be a hash-map, got %T", args[0])
			}

			key := args[1].String()
			return hm.Get(key), nil
		},
	})

	env.Set(Intern("hash-map-put"), &BuiltinFunction{
		Name: "hash-map-put",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 3 {
				return nil, fmt.Errorf("hash-map-put expects 3 arguments, got %d", len(args))
			}

			hm, ok := args[0].(*HashMap)
			if !ok {
				return nil, fmt.Errorf("first argument to hash-map-put must be a hash-map, got %T", args[0])
			}

			key := args[1].String()
			val := args[2]
			return hm.Put(key, val), nil
		},
	})

	env.Set(Intern("hash-map-keys"), &BuiltinFunction{
		Name: "hash-map-keys",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("hash-map-keys expects 1 argument, got %d", len(args))
			}

			hm, ok := args[0].(*HashMap)
			if !ok {
				return nil, fmt.Errorf("argument to hash-map-keys must be a hash-map, got %T", args[0])
			}

			return hm.Keys(), nil
		},
	})

	// Set operations
	env.Set(Intern("set"), &BuiltinFunction{
		Name: "set",
		Fn: func(args []Value, env *Environment) (Value, error) {
			s := NewSet()
			for _, arg := range args {
				s = s.Add(arg)
			}
			return s, nil
		},
	})

	env.Set(Intern("set-add"), &BuiltinFunction{
		Name: "set-add",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("set-add expects 2 arguments, got %d", len(args))
			}

			s, ok := args[0].(*Set)
			if !ok {
				return nil, fmt.Errorf("first argument to set-add must be a set, got %T", args[0])
			}

			return s.Add(args[1]), nil
		},
	})

	env.Set(Intern("set-contains?"), &BuiltinFunction{
		Name: "set-contains?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("set-contains? expects 2 arguments, got %d", len(args))
			}

			s, ok := args[0].(*Set)
			if !ok {
				return nil, fmt.Errorf("first argument to set-contains? must be a set, got %T", args[0])
			}

			return Boolean(s.Contains(args[1])), nil
		},
	})

	// Define 'cons' function to construct lists
	env.Set(Intern("cons"), &BuiltinFunction{
		Name: "cons",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("cons requires exactly 2 arguments")
			}

			// cons creates a new list with the first argument as head
			// and the second argument as tail
			if list, ok := args[1].(*List); ok {
				return NewList(append([]Value{args[0]}, list.elements...)...), nil
			}

			// If second arg is not a list, create a dotted pair (simple list)
			return NewList(args[0], args[1]), nil
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
