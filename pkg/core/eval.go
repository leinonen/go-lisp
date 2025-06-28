package core

import (
	"fmt"
	"os"
	"strings"
)

// Function interface for callable values
type Function interface {
	Call(args []Value, env *Environment) (Value, error)
}

// BuiltinFunction represents a built-in function
type BuiltinFunction struct {
	Name string
	Fn   func(args []Value, env *Environment) (Value, error)
}

func (bf *BuiltinFunction) Call(args []Value, env *Environment) (Value, error) {
	return bf.Fn(args, env)
}

func (bf *BuiltinFunction) String() string {
	return fmt.Sprintf("#<builtin:%s>", bf.Name)
}

// UserFunction represents a user-defined function
type UserFunction struct {
	Params *List
	Body   Value
	Env    *Environment
}

func (uf *UserFunction) Call(args []Value, env *Environment) (Value, error) {
	// Create new environment for function execution
	fnEnv := NewEnvironment(uf.Env)
	
	// Bind parameters to arguments
	err := bindParams(uf.Params, args, fnEnv)
	if err != nil {
		return nil, err
	}
	
	// Evaluate function body
	return Eval(uf.Body, fnEnv)
}

func (uf *UserFunction) String() string {
	return fmt.Sprintf("#<function>")
}

// bindParams binds function parameters to arguments
func bindParams(params *List, args []Value, env *Environment) error {
	paramList := listToSlice(params)
	
	if len(paramList) != len(args) {
		return fmt.Errorf("function expects %d arguments, got %d", len(paramList), len(args))
	}
	
	for i, param := range paramList {
		if sym, ok := param.(Symbol); ok {
			env.Set(sym, args[i])
		} else {
			return fmt.Errorf("parameter must be a symbol, got %T", param)
		}
	}
	
	return nil
}

// listToSlice converts a List to a slice of Values
func listToSlice(list *List) []Value {
	var result []Value
	current := list
	
	for current != nil {
		result = append(result, current.First())
		current = current.Rest()
	}
	
	return result
}

// Eval evaluates a Lisp expression
func Eval(expr Value, env *Environment) (Value, error) {
	switch v := expr.(type) {
	case Symbol:
		// Look up symbol in environment
		return env.Get(v)
		
	case *List:
		if v.IsEmpty() {
			return v, nil // Empty list evaluates to itself
		}
		
		// Check if first element is a special form
		if sym, ok := v.First().(Symbol); ok {
			result, err := evalSpecialForm(sym, v.Rest(), env)
			if err == nil {
				return result, nil
			}
			// If it's a recognized special form but had an error, return the error
			if isSpecialForm(sym) {
				return nil, err
			}
		}
		
		// Regular function call
		return evalFunctionCall(v, env)
		
	case Number, String, Keyword, *Vector:
		// These evaluate to themselves
		return expr, nil
		
	default:
		return expr, nil
	}
}

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

// evalFunctionCall evaluates a function call
func evalFunctionCall(list *List, env *Environment) (Value, error) {
	// Evaluate the function
	fn, err := Eval(list.First(), env)
	if err != nil {
		return nil, err
	}
	
	// Check if it's callable
	callable, ok := fn.(Function)
	if !ok {
		return nil, fmt.Errorf("cannot call non-function: %T", fn)
	}
	
	// Evaluate arguments
	var args []Value
	current := list.Rest()
	
	for current != nil {
		arg, err := Eval(current.First(), env)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		current = current.Rest()
	}
	
	// Call the function
	return callable.Call(args, env)
}

// isTruthy determines if a value is truthy
func isTruthy(v Value) bool {
	switch val := v.(type) {
	case Nil:
		return false
	case Number:
		if val.IsInteger() {
			return val.ToInt() != 0
		}
		return val.ToFloat() != 0.0
	case String:
		return string(val) != ""
	default:
		return true
	}
}

// NewCoreEnvironment creates an environment with core primitives
func NewCoreEnvironment() *Environment {
	env := NewEnvironment(nil)
	
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
					return nil, fmt.Errorf("+ expects numbers, got %T", arg)
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
				return nil, fmt.Errorf("- expects at least 1 argument")
			}
			
			first, ok := args[0].(Number)
			if !ok {
				return nil, fmt.Errorf("- expects numbers, got %T", args[0])
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
				// 1/x
				if first.ToFloat() == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return NewNumber(1.0 / first.ToFloat()), nil
			}
			
			// Division always returns float
			result := first.ToFloat()
			
			for _, arg := range args[1:] {
				if num, ok := arg.(Number); ok {
					divisor := num.ToFloat()
					if divisor == 0 {
						return nil, fmt.Errorf("division by zero")
					}
					result /= divisor
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
	
	// List operations
	env.Set(Intern("cons"), &BuiltinFunction{
		Name: "cons",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("cons expects 2 arguments")
			}
			
			if list, ok := args[1].(*List); ok {
				return &List{head: args[0], tail: list}, nil
			}
			
			// cons with non-list creates a new list
			return NewList(args[0], args[1]), nil
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
				return nil, fmt.Errorf("first expects list or vector, got %T", args[0])
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
					return (*List)(nil), nil  // Return nil list, which prints as "()"
				}
				rest := coll.Rest()
				if rest == nil {
					return (*List)(nil), nil  // Return nil list, which prints as "()"
				}
				return rest, nil
			case *Vector:
				if coll.Count() == 0 {
					return (*List)(nil), nil  // Return empty list
				}
				// Convert vector rest to list
				var elements []Value
				for i := 1; i < coll.Count(); i++ {
					elements = append(elements, coll.Get(i))
				}
				return NewList(elements...), nil
			case Nil:
				return (*List)(nil), nil  // Return nil list, which prints as "()"
			default:
				return nil, fmt.Errorf("rest expects list or vector, got %T", args[0])
			}
		},
	})
	
	// Meta-programming
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
	
	// Type predicates
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
	
	env.Set(Intern("string-length"), &BuiltinFunction{
		Name: "string-length",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("string-length expects 1 argument")
			}
			
			str, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("string-length expects string")
			}
			
			return NewNumber(int64(len(string(str)))), nil
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
			var result []Value
			for _, part := range parts {
				result = append(result, String(part))
			}
			
			return NewVector(result...), nil
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
	
	// I/O operations
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
	
	// List construction function
	env.Set(Intern("list"), &BuiltinFunction{
		Name: "list",
		Fn: func(args []Value, env *Environment) (Value, error) {
			return NewList(args...), nil
		},
	})
	
	// Add built-in symbols
	env.Set(Intern("nil"), Nil{})
	env.Set(Intern("true"), Symbol("true"))
	
	return env
}

// valuesEqual compares two values for equality
func valuesEqual(a, b Value) bool {
	switch va := a.(type) {
	case Symbol:
		if vb, ok := b.(Symbol); ok {
			return va == vb
		}
	case String:
		if vb, ok := b.(String); ok {
			return va == vb
		}
	case Number:
		if vb, ok := b.(Number); ok {
			return va.ToFloat() == vb.ToFloat()
		}
	case Keyword:
		if vb, ok := b.(Keyword); ok {
			return va == vb
		}
	case Nil:
		_, ok := b.(Nil)
		return ok
	}
	return false
}