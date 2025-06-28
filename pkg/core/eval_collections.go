package core

import "fmt"

// setupCollectionOperations adds collection operations and type predicates to the environment
func setupCollectionOperations(env *Environment) {
	// Collection operations
	env.Set(Intern("count"), &BuiltinFunction{
		Name: "count",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("count expects 1 argument")
			}

			switch coll := args[0].(type) {
			case *List:
				count := int64(0)
				current := coll
				for current != nil {
					count++
					current = current.Rest()
				}
				return NewNumber(count), nil
			case *Vector:
				return NewNumber(int64(coll.Count())), nil
			case String:
				return NewNumber(int64(len(string(coll)))), nil
			case Nil:
				return NewNumber(int64(0)), nil
			default:
				return nil, fmt.Errorf("count expects collection, got %T", args[0])
			}
		},
	})

	env.Set(Intern("empty?"), &BuiltinFunction{
		Name: "empty?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("empty? expects 1 argument")
			}

			switch coll := args[0].(type) {
			case *List:
				if coll.IsEmpty() {
					return Symbol("true"), nil
				}
				return Nil{}, nil
			case *Vector:
				if coll.Count() == 0 {
					return Symbol("true"), nil
				}
				return Nil{}, nil
			case String:
				if len(string(coll)) == 0 {
					return Symbol("true"), nil
				}
				return Nil{}, nil
			case Nil:
				return Symbol("true"), nil
			default:
				return nil, fmt.Errorf("empty? expects collection, got %T", args[0])
			}
		},
	})

	env.Set(Intern("nth"), &BuiltinFunction{
		Name: "nth",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 2 || len(args) > 3 {
				return nil, fmt.Errorf("nth expects 2-3 arguments")
			}

			n, ok := args[1].(Number)
			if !ok {
				return nil, fmt.Errorf("nth expects number as second argument")
			}

			index := int(n.ToInt())

			switch coll := args[0].(type) {
			case *List:
				current := coll
				for i := 0; i < index && current != nil; i++ {
					current = current.Rest()
				}
				if current == nil {
					if len(args) == 3 {
						return args[2], nil // Return default value
					}
					return nil, fmt.Errorf("index %d out of bounds", index)
				}
				return current.First(), nil
			case *Vector:
				if index < 0 || index >= coll.Count() {
					if len(args) == 3 {
						return args[2], nil // Return default value
					}
					return nil, fmt.Errorf("index %d out of bounds", index)
				}
				return coll.Get(index), nil
			case String:
				s := string(coll)
				if index < 0 || index >= len(s) {
					if len(args) == 3 {
						return args[2], nil // Return default value
					}
					return nil, fmt.Errorf("index %d out of bounds", index)
				}
				return String(string(s[index])), nil
			default:
				return nil, fmt.Errorf("nth expects collection, got %T", args[0])
			}
		},
	})

	env.Set(Intern("conj"), &BuiltinFunction{
		Name: "conj",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("conj expects at least 2 arguments")
			}

			coll := args[0]
			elements := args[1:]

			switch c := coll.(type) {
			case *List:
				result := c
				// For lists, conj adds to the front
				for i := len(elements) - 1; i >= 0; i-- {
					result = &List{head: elements[i], tail: result}
				}
				return result, nil
			case *Vector:
				// For vectors, conj adds to the end
				newElements := make([]Value, c.Count()+len(elements))
				for i := 0; i < c.Count(); i++ {
					newElements[i] = c.Get(i)
				}
				for i, elem := range elements {
					newElements[c.Count()+i] = elem
				}
				return NewVector(newElements...), nil
			case Nil:
				// Conj on nil creates a list
				result := (*List)(nil)
				for i := len(elements) - 1; i >= 0; i-- {
					result = &List{head: elements[i], tail: result}
				}
				return result, nil
			default:
				return nil, fmt.Errorf("conj expects collection, got %T", coll)
			}
		},
	})

	// List construction and access functions (these are already in core)
	env.Set(Intern("cons"), &BuiltinFunction{
		Name: "cons",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("cons expects 2 arguments")
			}

			// If the second argument is nil, create a list with nil as the second element
			if _, isNil := args[1].(Nil); isNil {
				return NewList(args[0], Nil{}), nil
			}

			return &List{head: args[0], tail: toList(args[1])}, nil
		},
	})

	env.Set(Intern("first"), &BuiltinFunction{
		Name: "first",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("first expects 1 argument")
			}

			switch coll := args[0].(type) {
			case *List:
				if coll.IsEmpty() {
					return Nil{}, nil
				}
				return coll.First(), nil
			case *Vector:
				if coll.Count() == 0 {
					return Nil{}, nil
				}
				return coll.Get(0), nil
			case Nil:
				return Nil{}, nil
			default:
				return nil, fmt.Errorf("first expects collection, got %T", args[0])
			}
		},
	})

	env.Set(Intern("rest"), &BuiltinFunction{
		Name: "rest",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("rest expects 1 argument")
			}

			switch coll := args[0].(type) {
			case *List:
				if coll.IsEmpty() {
					return (*List)(nil), nil
				}
				return coll.Rest(), nil
			case *Vector:
				if coll.Count() <= 1 {
					return (*List)(nil), nil
				}
				elements := make([]Value, coll.Count()-1)
				for i := 1; i < coll.Count(); i++ {
					elements[i-1] = coll.Get(i)
				}
				return NewList(elements...), nil
			case Nil:
				return (*List)(nil), nil
			default:
				return nil, fmt.Errorf("rest expects collection, got %T", args[0])
			}
		},
	})

	env.Set(Intern("list"), &BuiltinFunction{
		Name: "list",
		Fn: func(args []Value, env *Environment) (Value, error) {
			return NewList(args...), nil
		},
	})

	// Type predicates for collections
	env.Set(Intern("list?"), &BuiltinFunction{
		Name: "list?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("list? expects 1 argument")
			}

			if _, ok := args[0].(*List); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("vector?"), &BuiltinFunction{
		Name: "vector?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("vector? expects 1 argument")
			}

			if _, ok := args[0].(*Vector); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})
}

// Helper function to convert a value to a list
func toList(v Value) *List {
	if list, ok := v.(*List); ok {
		return list
	}
	if v == nil {
		return nil
	}
	// If it's a Nil value, return nil list (proper termination)
	if _, ok := v.(Nil); ok {
		return nil
	}
	// For other types, treat as empty list
	return nil
}
