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
			case *HashMap:
				return NewNumber(int64(coll.Count())), nil
			case *Set:
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

	// Add length as alias for count
	env.Set(Intern("length"), &BuiltinFunction{
		Name: "length",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("length expects 1 argument")
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
			case *HashMap:
				return NewNumber(int64(coll.Count())), nil
			case *Set:
				return NewNumber(int64(coll.Count())), nil
			case String:
				return NewNumber(int64(len(string(coll)))), nil
			case Nil:
				return NewNumber(int64(0)), nil
			default:
				return nil, fmt.Errorf("length expects collection, got %T", args[0])
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
			case *HashMap:
				if coll.Count() == 0 {
					return Symbol("true"), nil
				}
				return Nil{}, nil
			case *Set:
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

	env.Set(Intern("vector"), &BuiltinFunction{
		Name: "vector",
		Fn: func(args []Value, env *Environment) (Value, error) {
			return NewVector(args...), nil
		},
	})

	env.Set(Intern("hash-map"), &BuiltinFunction{
		Name: "hash-map",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args)%2 != 0 {
				return nil, fmt.Errorf("hash-map expects even number of arguments")
			}
			return NewHashMapWithPairs(args...), nil
		},
	})

	env.Set(Intern("hash-map?"), &BuiltinFunction{
		Name: "hash-map?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("hash-map? expects 1 argument")
			}

			if _, ok := args[0].(*HashMap); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("set"), &BuiltinFunction{
		Name: "set",
		Fn: func(args []Value, env *Environment) (Value, error) {
			return NewSetWithElements(args...), nil
		},
	})

	env.Set(Intern("set?"), &BuiltinFunction{
		Name: "set?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("set? expects 1 argument")
			}

			if _, ok := args[0].(*Set); ok {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("get"), &BuiltinFunction{
		Name: "get",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 2 || len(args) > 3 {
				return nil, fmt.Errorf("get expects 2-3 arguments")
			}

			switch coll := args[0].(type) {
			case *HashMap:
				value := coll.Get(args[1])
				if _, isNil := value.(Nil); isNil && len(args) == 3 {
					return args[2], nil // Return default value
				}
				return value, nil
			case *Vector:
				if n, ok := args[1].(Number); ok {
					index := int(n.ToInt())
					if index < 0 || index >= coll.Count() {
						if len(args) == 3 {
							return args[2], nil // Return default value
						}
						return Nil{}, nil
					}
					return coll.Get(index), nil
				}
				return nil, fmt.Errorf("get expects number index for vector")
			default:
				return nil, fmt.Errorf("get expects hash-map or vector, got %T", args[0])
			}
		},
	})

	env.Set(Intern("assoc"), &BuiltinFunction{
		Name: "assoc",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 3 || len(args)%2 == 0 {
				return nil, fmt.Errorf("assoc expects odd number of arguments (at least 3)")
			}

			if hm, ok := args[0].(*HashMap); ok {
				// Create a new hash-map with the same pairs
				newHM := NewHashMap()
				for _, key := range hm.keys {
					newHM.Set(key, hm.Get(key))
				}
				// Add new pairs
				for i := 1; i < len(args)-1; i += 2 {
					newHM.Set(args[i], args[i+1])
				}
				return newHM, nil
			}
			return nil, fmt.Errorf("assoc expects hash-map as first argument")
		},
	})

	env.Set(Intern("dissoc"), &BuiltinFunction{
		Name: "dissoc",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("dissoc expects at least 2 arguments")
			}

			if hm, ok := args[0].(*HashMap); ok {
				newHM := NewHashMap()
				keysToRemove := make(map[string]bool)
				for i := 1; i < len(args); i++ {
					keysToRemove[hm.keyToString(args[i])] = true
				}
				for _, key := range hm.keys {
					if !keysToRemove[hm.keyToString(key)] {
						newHM.Set(key, hm.Get(key))
					}
				}
				return newHM, nil
			}
			return nil, fmt.Errorf("dissoc expects hash-map as first argument")
		},
	})

	env.Set(Intern("contains?"), &BuiltinFunction{
		Name: "contains?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("contains? expects 2 arguments")
			}

			switch coll := args[0].(type) {
			case *HashMap:
				if coll.ContainsKey(args[1]) {
					return Symbol("true"), nil
				}
				return Nil{}, nil
			case *Set:
				if coll.Contains(args[1]) {
					return Symbol("true"), nil
				}
				return Nil{}, nil
			default:
				return nil, fmt.Errorf("contains? expects hash-map or set, got %T", args[0])
			}
		},
	})

	// Add hash-map-put as alias for assoc
	env.Set(Intern("hash-map-put"), &BuiltinFunction{
		Name: "hash-map-put",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 3 {
				return nil, fmt.Errorf("hash-map-put expects at least 3 arguments (map, key, value, ...)")
			}

			hm, ok := args[0].(*HashMap)
			if !ok {
				return nil, fmt.Errorf("hash-map-put expects hash-map as first argument, got %T", args[0])
			}

			// Create new hash-map with all existing mappings
			newHM := NewHashMap()
			for _, key := range hm.keys {
				newHM.Set(key, hm.Get(key))
			}

			// Add new key-value pairs (must be even number of key-value arguments)
			if (len(args)-1)%2 != 0 {
				return nil, fmt.Errorf("hash-map-put expects even number of key-value arguments")
			}

			for i := 1; i < len(args); i += 2 {
				key := args[i]
				value := args[i+1]
				newHM.Set(key, value)
			}

			return newHM, nil
		},
	})

	// Map keys and values functions
	env.Set(Intern("keys"), &BuiltinFunction{
		Name: "keys",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("keys expects 1 argument")
			}

			switch coll := args[0].(type) {
			case *HashMap:
				if coll.Count() == 0 {
					return (*List)(nil), nil
				}
				return NewList(coll.keys...), nil
			default:
				return nil, fmt.Errorf("keys expects hash-map, got %T", args[0])
			}
		},
	})

	env.Set(Intern("vals"), &BuiltinFunction{
		Name: "vals",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("vals expects 1 argument")
			}

			switch coll := args[0].(type) {
			case *HashMap:
				if coll.Count() == 0 {
					return (*List)(nil), nil
				}
				values := make([]Value, 0, len(coll.keys))
				for _, key := range coll.keys {
					values = append(values, coll.Get(key))
				}
				return NewList(values...), nil
			default:
				return nil, fmt.Errorf("vals expects hash-map, got %T", args[0])
			}
		},
	})

	env.Set(Intern("zipmap"), &BuiltinFunction{
		Name: "zipmap",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("zipmap expects 2 arguments")
			}

			// Convert both arguments to slices of values
			keys, err := collectionToSlice(args[0])
			if err != nil {
				return nil, fmt.Errorf("zipmap expects collection as first argument: %v", err)
			}

			values, err := collectionToSlice(args[1])
			if err != nil {
				return nil, fmt.Errorf("zipmap expects collection as second argument: %v", err)
			}

			// Create hash-map from key-value pairs
			result := NewHashMap()
			minLen := len(keys)
			if len(values) < minLen {
				minLen = len(values)
			}

			for i := 0; i < minLen; i++ {
				result.Set(keys[i], values[i])
			}

			return result, nil
		},
	})

	// Set operations
	env.Set(Intern("union"), &BuiltinFunction{
		Name: "union",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("union expects at least 2 arguments")
			}

			// All arguments must be sets
			for i, arg := range args {
				if _, ok := arg.(*Set); !ok {
					return nil, fmt.Errorf("union expects set as argument %d, got %T", i+1, arg)
				}
			}

			result := NewSet()

			// Add all elements from all sets
			for _, arg := range args {
				set := arg.(*Set)
				for _, elem := range set.order {
					result.Add(elem)
				}
			}

			return result, nil
		},
	})

	env.Set(Intern("intersection"), &BuiltinFunction{
		Name: "intersection",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("intersection expects at least 2 arguments")
			}

			// All arguments must be sets
			for i, arg := range args {
				if _, ok := arg.(*Set); !ok {
					return nil, fmt.Errorf("intersection expects set as argument %d, got %T", i+1, arg)
				}
			}

			// Start with elements from the first set
			firstSet := args[0].(*Set)
			result := NewSet()

			// Check each element from the first set
			for _, elem := range firstSet.order {
				inAllSets := true
				// Check if element exists in all other sets
				for i := 1; i < len(args); i++ {
					otherSet := args[i].(*Set)
					if !otherSet.Contains(elem) {
						inAllSets = false
						break
					}
				}
				if inAllSets {
					result.Add(elem)
				}
			}

			return result, nil
		},
	})

	env.Set(Intern("difference"), &BuiltinFunction{
		Name: "difference",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("difference expects at least 2 arguments")
			}

			// All arguments must be sets
			for i, arg := range args {
				if _, ok := arg.(*Set); !ok {
					return nil, fmt.Errorf("difference expects set as argument %d, got %T", i+1, arg)
				}
			}

			// Start with elements from the first set
			firstSet := args[0].(*Set)
			result := NewSet()

			// Add elements from first set that don't exist in any other set
			for _, elem := range firstSet.order {
				inOtherSets := false
				// Check if element exists in any other set
				for i := 1; i < len(args); i++ {
					otherSet := args[i].(*Set)
					if otherSet.Contains(elem) {
						inOtherSets = true
						break
					}
				}
				if !inOtherSets {
					result.Add(elem)
				}
			}

			return result, nil
		},
	})

	env.Set(Intern("subset?"), &BuiltinFunction{
		Name: "subset?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("subset? expects 2 arguments")
			}

			set1, ok1 := args[0].(*Set)
			set2, ok2 := args[1].(*Set)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("subset? expects two sets")
			}

			// Check if all elements in set1 are in set2
			for _, elem := range set1.order {
				if !set2.Contains(elem) {
					return Nil{}, nil
				}
			}
			return Symbol("true"), nil
		},
	})

	env.Set(Intern("superset?"), &BuiltinFunction{
		Name: "superset?",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("superset? expects 2 arguments")
			}

			set1, ok1 := args[0].(*Set)
			set2, ok2 := args[1].(*Set)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("superset? expects two sets")
			}

			// Check if all elements in set2 are in set1
			for _, elem := range set2.order {
				if !set1.Contains(elem) {
					return Nil{}, nil
				}
			}
			return Symbol("true"), nil
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

// Helper function to convert collection to slice of values
func collectionToSlice(coll Value) ([]Value, error) {
	switch c := coll.(type) {
	case *List:
		var result []Value
		current := c
		for current != nil && !current.IsEmpty() {
			result = append(result, current.First())
			current = current.Rest()
		}
		return result, nil
	case *Vector:
		result := make([]Value, c.Count())
		for i := 0; i < c.Count(); i++ {
			result[i] = c.Get(i)
		}
		return result, nil
	case *Set:
		result := make([]Value, 0, c.Count())
		for _, elem := range c.order {
			result = append(result, elem)
		}
		return result, nil
	case Nil:
		return []Value{}, nil
	default:
		return nil, fmt.Errorf("expected collection, got %T", coll)
	}
}
