package core

import (
	"fmt"
)

// evalSpecialForm handles special forms
func evalSpecialForm(sym Symbol, args *List, env *Environment) (Value, error) {
	switch sym {
	case "quote":
		argSlice := listToSlice(args)
		if len(argSlice) != 1 {
			return nil, NewArityError("quote expects 1 argument, got %d", len(argSlice))
		}
		return argSlice[0], nil

	case "quasiquote":
		argSlice := listToSlice(args)
		if len(argSlice) != 1 {
			return nil, fmt.Errorf("quasiquote expects 1 argument, got %d", len(argSlice))
		}
		return evalQuasiquote(argSlice[0], env)

	case "if":
		argSlice := listToSlice(args)
		if len(argSlice) < 2 || len(argSlice) > 3 {
			return nil, NewArityError("if expects 2-3 arguments, got %d", len(argSlice))
		}

		condition, err := Eval(argSlice[0], env)
		if err != nil {
			return nil, err
		}

		if isTruthy(condition) {
			return Eval(argSlice[1], env)
		} else if len(argSlice) == 3 {
			return Eval(argSlice[2], env)
		}
		return Nil{}, nil

	case "def":
		argSlice := listToSlice(args)
		if len(argSlice) != 2 {
			return nil, NewArityError("def expects 2 arguments, got %d", len(argSlice))
		}

		sym, ok := argSlice[0].(Symbol)
		if !ok {
			return nil, fmt.Errorf("def expects symbol as first argument, got %T", argSlice[0])
		}

		value, err := Eval(argSlice[1], env)
		if err != nil {
			return nil, err
		}

		env.Set(sym, value)
		return sym, nil

	case "fn":
		argSlice := listToSlice(args)
		if len(argSlice) < 2 {
			return nil, fmt.Errorf("fn expects at least 2 arguments, got %d", len(argSlice))
		}

		// Handle both lists and vectors for parameters
		var params *List
		switch p := argSlice[0].(type) {
		case *List:
			params = p
		case *Vector:
			// Convert vector to list
			var elements []Value
			for i := 0; i < p.Count(); i++ {
				elements = append(elements, p.Get(i))
			}
			params = NewList(elements...)
		default:
			return nil, fmt.Errorf("fn expects list or vector as first argument, got %T", argSlice[0])
		}

		// Handle multiple body expressions by wrapping in 'do'
		var body Value
		if len(argSlice) == 2 {
			body = argSlice[1]
		} else {
			// Multiple body expressions - wrap in do
			bodyExprs := argSlice[1:]
			doList := make([]Value, len(bodyExprs)+1)
			doList[0] = Symbol("do")
			copy(doList[1:], bodyExprs)
			body = NewList(doList...)
		}

		return &UserFunction{
			Params: params,
			Body:   body,
			Env:    env,
		}, nil

	case "do":
		argSlice := listToSlice(args)
		var result Value = Nil{}

		for _, expr := range argSlice {
			var err error
			result, err = Eval(expr, env)
			if err != nil {
				return nil, err
			}
		}

		return result, nil

	case "let":
		argSlice := listToSlice(args)
		if len(argSlice) < 2 {
			return nil, fmt.Errorf("let expects at least 2 arguments")
		}

		// Create new environment for let bindings
		letEnv := NewEnvironment(env)

		// Process bindings
		bindings := argSlice[0]
		var bindingList []Value

		switch b := bindings.(type) {
		case *List:
			bindingList = listToSlice(b)
		case *Vector:
			for i := 0; i < b.Count(); i++ {
				bindingList = append(bindingList, b.Get(i))
			}
		default:
			return nil, fmt.Errorf("let expects vector or list for bindings")
		}

		if len(bindingList)%2 != 0 {
			return nil, fmt.Errorf("let bindings must be even number of forms")
		}

		// Bind variables
		for i := 0; i < len(bindingList); i += 2 {
			sym, ok := bindingList[i].(Symbol)
			if !ok {
				return nil, fmt.Errorf("let binding names must be symbols")
			}

			value, err := Eval(bindingList[i+1], letEnv)
			if err != nil {
				return nil, err
			}

			letEnv.Set(sym, value)
		}

		// Evaluate body expressions
		var result Value = Nil{}
		for _, expr := range argSlice[1:] {
			var err error
			result, err = Eval(expr, letEnv)
			if err != nil {
				return nil, err
			}
		}

		return result, nil

	case "defmacro":
		argSlice := listToSlice(args)
		if len(argSlice) != 3 {
			return nil, fmt.Errorf("defmacro expects 3 arguments (name params body), got %d", len(argSlice))
		}

		sym, ok := argSlice[0].(Symbol)
		if !ok {
			return nil, fmt.Errorf("defmacro expects symbol as first argument, got %T", argSlice[0])
		}

		// Handle both lists and vectors for parameters
		var params *List
		switch p := argSlice[1].(type) {
		case *List:
			params = p
		case *Vector:
			// Convert vector to list
			var elements []Value
			for i := 0; i < p.Count(); i++ {
				elements = append(elements, p.Get(i))
			}
			params = NewList(elements...)
		default:
			return nil, fmt.Errorf("defmacro expects list or vector as second argument, got %T", argSlice[1])
		}

		macro := &Macro{
			Name:   sym,
			Params: params,
			Body:   argSlice[2],
			Env:    env,
		}

		env.Set(sym, macro)
		return sym, nil

	case "defn":
		argSlice := listToSlice(args)
		if len(argSlice) < 3 {
			return nil, fmt.Errorf("defn expects at least 3 arguments (name params body...), got %d", len(argSlice))
		}

		sym, ok := argSlice[0].(Symbol)
		if !ok {
			return nil, fmt.Errorf("defn expects symbol as first argument, got %T", argSlice[0])
		}

		// Handle both lists and vectors for parameters
		var params *List
		switch p := argSlice[1].(type) {
		case *List:
			params = p
		case *Vector:
			// Convert vector to list
			var elements []Value
			for i := 0; i < p.Count(); i++ {
				elements = append(elements, p.Get(i))
			}
			params = NewList(elements...)
		default:
			return nil, fmt.Errorf("defn expects list or vector as second argument, got %T", argSlice[1])
		}

		// Handle multiple body expressions by wrapping in 'do'
		var body Value
		if len(argSlice) == 3 {
			body = argSlice[2]
		} else {
			// Multiple body expressions - wrap in do
			bodyExprs := argSlice[2:]
			doList := make([]Value, len(bodyExprs)+1)
			doList[0] = Symbol("do")
			copy(doList[1:], bodyExprs)
			body = NewList(doList...)
		}

		function := &UserFunction{
			Params: params,
			Body:   body,
			Env:    env,
		}

		env.Set(sym, function)
		return sym, nil

	case "cond":
		argSlice := listToSlice(args)
		if len(argSlice)%2 != 0 && len(argSlice) > 0 {
			// Check if last argument is :else
			if len(argSlice)%2 == 1 {
				if sym, ok := argSlice[len(argSlice)-2].(Symbol); !ok || sym != ":else" {
					return nil, fmt.Errorf("cond expects even number of arguments or :else clause")
				}
			}
		}

		// Evaluate condition/expression pairs
		for i := 0; i < len(argSlice); i += 2 {
			if i+1 >= len(argSlice) {
				return nil, fmt.Errorf("cond: missing expression for condition")
			}

			condition := argSlice[i]

			// Special case for :else
			if sym, ok := condition.(Symbol); ok && sym == ":else" {
				return Eval(argSlice[i+1], env)
			}

			// Evaluate condition
			condResult, err := Eval(condition, env)
			if err != nil {
				return nil, err
			}

			// If condition is truthy, evaluate and return the expression
			if isTruthy(condResult) {
				return Eval(argSlice[i+1], env)
			}
		}

		// No condition matched
		return Nil{}, nil

	case "and":
		argSlice := listToSlice(args)
		if len(argSlice) == 0 {
			return Symbol("true"), nil
		}

		// Short-circuiting: evaluate expressions left-to-right
		// Return first falsy value or last value if all are truthy
		for i, expr := range argSlice {
			result, err := Eval(expr, env)
			if err != nil {
				return nil, err
			}

			// If falsy, return this value (short-circuit)
			if !isTruthy(result) {
				return result, nil
			}

			// If this is the last expression, return its value
			if i == len(argSlice)-1 {
				return result, nil
			}
		}

		// Should never reach here
		return Symbol("true"), nil

	case "or":
		argSlice := listToSlice(args)
		if len(argSlice) == 0 {
			return Nil{}, nil
		}

		// Short-circuiting: evaluate expressions left-to-right
		// Return first truthy value or last value if all are falsy
		for i, expr := range argSlice {
			result, err := Eval(expr, env)
			if err != nil {
				return nil, err
			}

			// If truthy, return this value (short-circuit)
			if isTruthy(result) {
				return result, nil
			}

			// If this is the last expression, return its value
			if i == len(argSlice)-1 {
				return result, nil
			}
		}

		// Should never reach here
		return Nil{}, nil

	case "loop":
		argSlice := listToSlice(args)
		if len(argSlice) < 2 {
			return nil, fmt.Errorf("loop expects at least 2 arguments (bindings body...)")
		}

		// Create new environment for loop bindings
		loopEnv := NewEnvironment(env)

		// Process bindings (similar to let)
		bindings := argSlice[0]
		var bindingList []Value

		switch b := bindings.(type) {
		case *List:
			bindingList = listToSlice(b)
		case *Vector:
			for i := 0; i < b.Count(); i++ {
				bindingList = append(bindingList, b.Get(i))
			}
		default:
			return nil, fmt.Errorf("loop expects vector or list for bindings")
		}

		if len(bindingList)%2 != 0 {
			return nil, fmt.Errorf("loop bindings must be even number of forms")
		}

		// Extract parameter names for recur
		var paramNames []Symbol
		for i := 0; i < len(bindingList); i += 2 {
			sym, ok := bindingList[i].(Symbol)
			if !ok {
				return nil, fmt.Errorf("loop binding names must be symbols")
			}
			paramNames = append(paramNames, sym)
		}

		// Initial binding values
		var initialValues []Value
		for i := 1; i < len(bindingList); i += 2 {
			value, err := Eval(bindingList[i], env)
			if err != nil {
				return nil, err
			}
			initialValues = append(initialValues, value)
		}

		// Loop execution with recur handling
		currentValues := initialValues
		for {
			// Bind current values
			for i, sym := range paramNames {
				loopEnv.Set(sym, currentValues[i])
			}

			// Evaluate body expressions
			var result Value = Nil{}
			for _, expr := range argSlice[1:] {
				var err error
				result, err = Eval(expr, loopEnv)
				if err != nil {
					return nil, err
				}

				// Check if result is a recur
				if recurVal, ok := result.(*RecurValue); ok {
					// Validate recur arity
					if len(recurVal.Values) != len(paramNames) {
						return nil, fmt.Errorf("recur expects %d arguments, got %d", len(paramNames), len(recurVal.Values))
					}
					// Update values for next iteration
					currentValues = recurVal.Values
					goto continueLoop
				}
			}

			// No recur found, return the final result
			return result, nil

		continueLoop:
			// Continue to next iteration
		}

	case "recur":
		argSlice := listToSlice(args)
		
		// Evaluate all arguments
		var values []Value
		for _, arg := range argSlice {
			value, err := Eval(arg, env)
			if err != nil {
				return nil, err
			}
			values = append(values, value)
		}

		// Return a RecurValue to be caught by loop
		return &RecurValue{Values: values}, nil
	}

	return nil, fmt.Errorf("unknown special form: %s", sym)
}

// isSpecialForm checks if a symbol is a special form
func isSpecialForm(sym Symbol) bool {
	switch sym {
	case "quote", "quasiquote", "if", "def", "fn", "do", "let", "defmacro", "defn", "cond", "and", "or", "loop", "recur":
		return true
	default:
		return false
	}
}

// evalQuasiquote handles quasiquote evaluation
func evalQuasiquote(expr Value, env *Environment) (Value, error) {
	return quasiQuoteExpand(expr, env)
}

// quasiQuoteExpand recursively expands quasiquoted expressions
func quasiQuoteExpand(expr Value, env *Environment) (Value, error) {
	switch v := expr.(type) {
	case *List:
		if v.IsEmpty() {
			return v, nil
		}

		// Check if this is an unquote form
		if first := v.First(); first != nil {
			if sym, ok := first.(Symbol); ok && sym == "unquote" {
				rest := v.Rest()
				args := listToSlice(rest)
				if len(args) != 1 {
					return nil, fmt.Errorf("unquote expects 1 argument, got %d", len(args))
				}
				return Eval(args[0], env)
			}
		}

		// Expand list elements, handling unquote-splicing
		var result []Value
		current := v
		for !current.IsEmpty() {
			elem := current.First()

			// Check for unquote-splicing
			if elemList, ok := elem.(*List); ok && !elemList.IsEmpty() {
				if sym, ok := elemList.First().(Symbol); ok && sym == "unquote-splicing" {
					rest := elemList.Rest()
					args := listToSlice(rest)
					if len(args) != 1 {
						return nil, fmt.Errorf("unquote-splicing expects 1 argument, got %d", len(args))
					}

					// Evaluate the spliced expression
					spliced, err := Eval(args[0], env)
					if err != nil {
						return nil, err
					}

					// Convert to slice and append all elements
					switch s := spliced.(type) {
					case *List:
						splicedSlice := listToSlice(s)
						result = append(result, splicedSlice...)
					case *Vector:
						for i := 0; i < s.Count(); i++ {
							result = append(result, s.Get(i))
						}
					default:
						return nil, fmt.Errorf("unquote-splicing can only splice sequences, got %T", spliced)
					}
				} else {
					// Regular element, expand recursively
					expanded, err := quasiQuoteExpand(elem, env)
					if err != nil {
						return nil, err
					}
					result = append(result, expanded)
				}
			} else {
				// Regular element, expand recursively
				expanded, err := quasiQuoteExpand(elem, env)
				if err != nil {
					return nil, err
				}
				result = append(result, expanded)
			}

			current = current.Rest()
		}

		return NewList(result...), nil

	case *Vector:
		// Expand vector elements
		var result []Value
		for i := 0; i < v.Count(); i++ {
			elem := v.Get(i)

			// Check for unquote-splicing in vectors
			if elemList, ok := elem.(*List); ok && !elemList.IsEmpty() {
				if sym, ok := elemList.First().(Symbol); ok && sym == "unquote-splicing" {
					rest := elemList.Rest()
					args := listToSlice(rest)
					if len(args) != 1 {
						return nil, fmt.Errorf("unquote-splicing expects 1 argument, got %d", len(args))
					}

					// Evaluate the spliced expression
					spliced, err := Eval(args[0], env)
					if err != nil {
						return nil, err
					}

					// Convert to slice and append all elements
					switch s := spliced.(type) {
					case *List:
						splicedSlice := listToSlice(s)
						result = append(result, splicedSlice...)
					case *Vector:
						for j := 0; j < s.Count(); j++ {
							result = append(result, s.Get(j))
						}
					default:
						return nil, fmt.Errorf("unquote-splicing can only splice sequences, got %T", spliced)
					}
				} else {
					// Regular element, expand recursively
					expanded, err := quasiQuoteExpand(elem, env)
					if err != nil {
						return nil, err
					}
					result = append(result, expanded)
				}
			} else {
				// Regular element, expand recursively
				expanded, err := quasiQuoteExpand(elem, env)
				if err != nil {
					return nil, err
				}
				result = append(result, expanded)
			}
		}

		return NewVector(result...), nil

	case *HashMap:
		// Expand hash map entries by iterating over internal keys
		var result []Value

		// Access the keys field through the structure
		for _, key := range v.keys {
			value := v.Get(key)

			// Expand key
			expandedKey, err := quasiQuoteExpand(key, env)
			if err != nil {
				return nil, err
			}

			// Expand value
			expandedValue, err := quasiQuoteExpand(value, env)
			if err != nil {
				return nil, err
			}

			result = append(result, expandedKey, expandedValue)
		}

		return NewHashMapWithPairs(result...), nil

	default:
		// Atoms (symbols, numbers, strings, etc.) are returned as-is
		return expr, nil
	}
}
