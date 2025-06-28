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
			return nil, fmt.Errorf("quote expects 1 argument, got %d", len(argSlice))
		}
		return argSlice[0], nil

	case "if":
		argSlice := listToSlice(args)
		if len(argSlice) < 2 || len(argSlice) > 3 {
			return nil, fmt.Errorf("if expects 2-3 arguments, got %d", len(argSlice))
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
			return nil, fmt.Errorf("def expects 2 arguments, got %d", len(argSlice))
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
		if len(argSlice) != 2 {
			return nil, fmt.Errorf("fn expects 2 arguments, got %d", len(argSlice))
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

		return &UserFunction{
			Params: params,
			Body:   argSlice[1],
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
	}

	return nil, fmt.Errorf("unknown special form: %s", sym)
}

// isSpecialForm checks if a symbol is a special form
func isSpecialForm(sym Symbol) bool {
	switch sym {
	case "quote", "if", "def", "fn", "do", "let":
		return true
	default:
		return false
	}
}
