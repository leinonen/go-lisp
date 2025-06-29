package core

import "fmt"

// setupArithmeticOperations adds arithmetic and comparison operations to the environment
func setupArithmeticOperations(env *Environment) {
	// Arithmetic operations
	env.Set(Intern("+"), &BuiltinFunction{
		Name: "+",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) == 0 {
				return NewNumber(int64(0)), nil
			}

			result := int64(0)
			isFloat := false
			floatResult := 0.0

			for _, arg := range args {
				if num, ok := arg.(Number); ok {
					if num.IsFloat() || isFloat {
						if !isFloat {
							floatResult = float64(result)
							isFloat = true
						}
						floatResult += num.ToFloat()
					} else {
						result += num.ToInt()
					}
				} else {
					return nil, NewTypeError("+ expects numbers, got %T", arg)
				}
			}

			if isFloat {
				return NewNumber(floatResult), nil
			}
			return NewNumber(result), nil
		},
	})

	env.Set(Intern("-"), &BuiltinFunction{
		Name: "-",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) == 0 {
				return nil, NewArityError("- expects at least 1 argument")
			}

			first, ok := args[0].(Number)
			if !ok {
				return nil, NewTypeError("- expects numbers, got %T", args[0])
			}

			if len(args) == 1 {
				// Unary minus
				if first.IsFloat() {
					return NewNumber(-first.ToFloat()), nil
				}
				return NewNumber(-first.ToInt()), nil
			}

			// Binary and n-ary minus
			result := first.ToInt()
			isFloat := first.IsFloat()
			floatResult := first.ToFloat()

			for _, arg := range args[1:] {
				if num, ok := arg.(Number); ok {
					if num.IsFloat() || isFloat {
						if !isFloat {
							floatResult = float64(result)
							isFloat = true
						}
						floatResult -= num.ToFloat()
					} else {
						result -= num.ToInt()
					}
				} else {
					return nil, fmt.Errorf("- expects numbers, got %T", arg)
				}
			}

			if isFloat {
				return NewNumber(floatResult), nil
			}
			return NewNumber(result), nil
		},
	})

	env.Set(Intern("*"), &BuiltinFunction{
		Name: "*",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) == 0 {
				return NewNumber(int64(1)), nil
			}

			result := int64(1)
			isFloat := false
			floatResult := 1.0

			for _, arg := range args {
				if num, ok := arg.(Number); ok {
					if num.IsFloat() || isFloat {
						if !isFloat {
							floatResult = float64(result)
							isFloat = true
						}
						floatResult *= num.ToFloat()
					} else {
						result *= num.ToInt()
					}
				} else {
					return nil, fmt.Errorf("* expects numbers, got %T", arg)
				}
			}

			if isFloat {
				return NewNumber(floatResult), nil
			}
			return NewNumber(result), nil
		},
	})

	env.Set(Intern("/"), &BuiltinFunction{
		Name: "/",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) == 0 {
				return nil, fmt.Errorf("/ expects at least 1 argument")
			}

			first, ok := args[0].(Number)
			if !ok {
				return nil, fmt.Errorf("/ expects numbers, got %T", args[0])
			}

			if len(args) == 1 {
				// Reciprocal
				if first.ToFloat() == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return NewNumber(1.0 / first.ToFloat()), nil
			}

			// Division
			result := first.ToFloat()

			for _, arg := range args[1:] {
				if num, ok := arg.(Number); ok {
					if num.ToFloat() == 0 {
						return nil, fmt.Errorf("division by zero")
					}
					result /= num.ToFloat()
				} else {
					return nil, fmt.Errorf("/ expects numbers, got %T", arg)
				}
			}

			return NewNumber(result), nil
		},
	})

	env.Set(Intern("%"), &BuiltinFunction{
		Name: "%",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("%% expects 2 arguments")
			}

			n1, ok1 := args[0].(Number)
			n2, ok2 := args[1].(Number)

			if !ok1 || !ok2 {
				return nil, fmt.Errorf("%% expects numbers")
			}

			if !n1.IsInteger() || !n2.IsInteger() {
				return nil, fmt.Errorf("%% expects integers")
			}

			divisor := n2.ToInt()
			if divisor == 0 {
				return nil, fmt.Errorf("modulo by zero")
			}

			result := n1.ToInt() % divisor
			return NewNumber(result), nil
		},
	})

	// Comparison operations
	env.Set(Intern("="), &BuiltinFunction{
		Name: "=",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("= expects at least 2 arguments")
			}

			first := args[0]
			for _, arg := range args[1:] {
				if !valuesEqual(first, arg) {
					return Nil{}, nil
				}
			}
			return Symbol("true"), nil
		},
	})

	env.Set(Intern("<"), &BuiltinFunction{
		Name: "<",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("< expects 2 arguments")
			}

			n1, ok1 := args[0].(Number)
			n2, ok2 := args[1].(Number)

			if !ok1 || !ok2 {
				return nil, fmt.Errorf("< expects numbers")
			}

			if n1.ToFloat() < n2.ToFloat() {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern(">"), &BuiltinFunction{
		Name: ">",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("> expects 2 arguments")
			}

			n1, ok1 := args[0].(Number)
			n2, ok2 := args[1].(Number)

			if !ok1 || !ok2 {
				return nil, fmt.Errorf("> expects numbers")
			}

			if n1.ToFloat() > n2.ToFloat() {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern(">="), &BuiltinFunction{
		Name: ">=",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf(">= expects 2 arguments")
			}

			n1, ok1 := args[0].(Number)
			n2, ok2 := args[1].(Number)

			if !ok1 || !ok2 {
				return nil, fmt.Errorf(">= expects numbers")
			}

			if n1.ToFloat() >= n2.ToFloat() {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})

	env.Set(Intern("<="), &BuiltinFunction{
		Name: "<=",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("<= expects 2 arguments")
			}

			n1, ok1 := args[0].(Number)
			n2, ok2 := args[1].(Number)

			if !ok1 || !ok2 {
				return nil, fmt.Errorf("<= expects numbers")
			}

			if n1.ToFloat() <= n2.ToFloat() {
				return Symbol("true"), nil
			}
			return Nil{}, nil
		},
	})
}
