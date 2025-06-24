package minimal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// REPL provides a simple read-eval-print loop for the minimal kernel
type REPL struct {
	Env *Environment
}

// NewREPL creates a new REPL instance
func NewREPL() *REPL {
	env := NewEnvironment(nil)
	repl := &REPL{Env: env}
	repl.setupBuiltins()
	return repl
}

// Run starts the REPL
func (r *REPL) Run() {
	fmt.Println("Minimal Lisp REPL - Type 'exit' to quit")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("minimal> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "exit" || input == "quit" {
			break
		}

		if input == "" {
			continue
		}

		// Parse and evaluate
		expr, err := r.parse(input)
		if err != nil {
			fmt.Printf("Parse error: %v\n", err)
			continue
		}

		result, err := Eval(expr, r.Env)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("=> %s\n", result.String())
	}
}

// Simple parser for the minimal kernel
func (r *REPL) parse(input string) (Value, error) {
	tokens := r.tokenize(input)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	return r.parseTokens(tokens)
}

func (r *REPL) tokenize(input string) []string {
	// Simple tokenizer - splits on whitespace and handles parentheses and square brackets
	var tokens []string
	var current strings.Builder

	for _, char := range input {
		switch char {
		case '(', ')', '[', ']':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(char))
		case ' ', '\t', '\n', '\r':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func (r *REPL) parseTokens(tokens []string) (Value, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("unexpected end of input")
	}

	token := tokens[0]

	if token == "(" {
		// Parse list
		return r.parseList(tokens[1:])
	}

	if token == "[" {
		// Parse vector
		return r.parseVector(tokens[1:])
	}

	// Parse atom
	return r.parseAtom(token), nil
}

func (r *REPL) parseList(tokens []string) (Value, error) {
	var elements []Value
	i := 0

	for i < len(tokens) {
		if tokens[i] == ")" {
			return NewList(elements...), nil
		}

		if tokens[i] == "(" {
			// Nested list
			subList, consumed, err := r.parseNestedList(tokens[i:])
			if err != nil {
				return nil, err
			}
			elements = append(elements, subList)
			i += consumed
		} else {
			// Atom
			elements = append(elements, r.parseAtom(tokens[i]))
			i++
		}
	}

	return nil, fmt.Errorf("unclosed list")
}

func (r *REPL) parseNestedList(tokens []string) (Value, int, error) {
	if tokens[0] != "(" {
		return nil, 0, fmt.Errorf("expected '('")
	}

	var elements []Value
	i := 1
	depth := 1

	for i < len(tokens) && depth > 0 {
		switch tokens[i] {
		case "(":
			depth++
			// Parse nested list recursively
			subList, consumed, err := r.parseNestedList(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			elements = append(elements, subList)
			i += consumed
		case "[":
			// Parse nested vector
			subVector, consumed, err := r.parseNestedVector(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			elements = append(elements, subVector)
			i += consumed
		case ")":
			depth--
			if depth == 0 {
				return NewList(elements...), i + 1, nil
			}
		default:
			elements = append(elements, r.parseAtom(tokens[i]))
			i++
		}
	}

	return nil, 0, fmt.Errorf("unclosed list")
}

func (r *REPL) parseVector(tokens []string) (Value, error) {
	var elements []Value
	i := 0

	for i < len(tokens) {
		if tokens[i] == "]" {
			return NewVector(elements...), nil
		}

		if tokens[i] == "(" {
			// Nested list
			subList, consumed, err := r.parseNestedList(tokens[i:])
			if err != nil {
				return nil, err
			}
			elements = append(elements, subList)
			i += consumed
		} else if tokens[i] == "[" {
			// Nested vector
			subVector, consumed, err := r.parseNestedVector(tokens[i:])
			if err != nil {
				return nil, err
			}
			elements = append(elements, subVector)
			i += consumed
		} else {
			// Atom
			elements = append(elements, r.parseAtom(tokens[i]))
			i++
		}
	}

	return nil, fmt.Errorf("unclosed vector")
}

func (r *REPL) parseNestedVector(tokens []string) (Value, int, error) {
	if tokens[0] != "[" {
		return nil, 0, fmt.Errorf("expected '['")
	}

	var elements []Value
	i := 1
	depth := 1

	for i < len(tokens) && depth > 0 {
		switch tokens[i] {
		case "[":
			depth++
			// Parse nested vector recursively
			subVector, consumed, err := r.parseNestedVector(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			elements = append(elements, subVector)
			i += consumed
		case "]":
			depth--
			if depth == 0 {
				return NewVector(elements...), i + 1, nil
			}
		case "(":
			// Parse nested list
			subList, consumed, err := r.parseNestedList(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			elements = append(elements, subList)
			i += consumed
		default:
			elements = append(elements, r.parseAtom(tokens[i]))
			i++
		}
	}

	return nil, 0, fmt.Errorf("unclosed vector")
}

func (r *REPL) parseAtom(token string) Value {
	// Try parsing as number
	if num, err := strconv.ParseFloat(token, 64); err == nil {
		return Number(num)
	}

	// Try parsing as boolean
	if token == "true" {
		return Boolean(true)
	}
	if token == "false" {
		return Boolean(false)
	}

	// Try parsing as nil
	if token == "nil" {
		return Nil{}
	}

	// Try parsing as string (simple quote handling)
	if len(token) >= 2 && token[0] == '"' && token[len(token)-1] == '"' {
		return String(token[1 : len(token)-1])
	}

	// Otherwise it's a symbol
	return Intern(token)
}

// setupBuiltins adds basic arithmetic functions to the environment
func (r *REPL) setupBuiltins() {
	// Addition
	r.Env.Set(Intern("+"), &BuiltinFunction{
		Name: "+",
		Fn: func(args []Value, env *Environment) (Value, error) {
			result := 0.0
			for _, arg := range args {
				if num, ok := arg.(Number); ok {
					result += float64(num)
				} else {
					return nil, fmt.Errorf("+ requires numbers, got %T", arg)
				}
			}
			return Number(result), nil
		},
	})

	// Subtraction
	r.Env.Set(Intern("-"), &BuiltinFunction{
		Name: "-",
		Fn: func(args []Value, env *Environment) (Value, error) {
			if len(args) == 0 {
				return nil, fmt.Errorf("- requires at least 1 argument")
			}

			first, ok := args[0].(Number)
			if !ok {
				return nil, fmt.Errorf("- requires numbers, got %T", args[0])
			}

			if len(args) == 1 {
				return Number(-float64(first)), nil
			}

			result := float64(first)
			for _, arg := range args[1:] {
				if num, ok := arg.(Number); ok {
					result -= float64(num)
				} else {
					return nil, fmt.Errorf("- requires numbers, got %T", arg)
				}
			}
			return Number(result), nil
		},
	})

	// Multiplication
	r.Env.Set(Intern("*"), &BuiltinFunction{
		Name: "*",
		Fn: func(args []Value, env *Environment) (Value, error) {
			result := 1.0
			for _, arg := range args {
				if num, ok := arg.(Number); ok {
					result *= float64(num)
				} else {
					return nil, fmt.Errorf("* requires numbers, got %T", arg)
				}
			}
			return Number(result), nil
		},
	})
}

// BuiltinFunction represents a built-in function
type BuiltinFunction struct {
	Name string
	Fn   func([]Value, *Environment) (Value, error)
}

func (bf *BuiltinFunction) Call(args []Value, env *Environment) (Value, error) {
	return bf.Fn(args, env)
}

func (bf *BuiltinFunction) String() string {
	return fmt.Sprintf("<builtin:%s>", bf.Name)
}
