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
	fmt.Println("Multi-line input supported - expressions are evaluated when parentheses are balanced")
	fmt.Println("Use ':reset' or ':clear' to discard incomplete input")
	scanner := bufio.NewScanner(os.Stdin)

	var inputBuffer strings.Builder

	for {
		// Determine the prompt based on whether we're continuing input
		prompt := "minimal> "
		if inputBuffer.Len() > 0 {
			prompt = "      ... "
		}

		fmt.Print(prompt)
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()

		// Check for exit commands
		if strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit" {
			if inputBuffer.Len() > 0 {
				fmt.Println("Discarding incomplete input")
				inputBuffer.Reset()
				continue
			}
			break
		}

		// Check for reset command to clear incomplete input
		if strings.TrimSpace(line) == ":reset" || strings.TrimSpace(line) == ":clear" {
			if inputBuffer.Len() > 0 {
				fmt.Println("Input buffer cleared")
				inputBuffer.Reset()
			}
			continue
		}

		// Add line to buffer
		if inputBuffer.Len() > 0 {
			inputBuffer.WriteString("\n")
		}
		inputBuffer.WriteString(line)

		input := inputBuffer.String()

		// Skip empty input
		if strings.TrimSpace(input) == "" {
			inputBuffer.Reset()
			continue
		}

		// Check if parentheses are balanced
		if !r.isBalanced(input) {
			// Continue reading more input
			continue
		}

		// Parse and evaluate with enhanced error handling
		expr, pos, err := ParseWithPositions(input, "<repl>")
		if err != nil {
			fmt.Printf("%v\n", err)
			inputBuffer.Reset()
			continue
		}

		// Create evaluation context
		ctx := NewEvaluationContext()
		if pos != nil {
			ctx.SetLocation(pos.Line, pos.Column, pos.File)
		}
		ctx.SetExpression(input)

		result, err := EvalWithContext(expr, r.Env, ctx)
		if err != nil {
			fmt.Printf("%v\n", err)
			inputBuffer.Reset()
			continue
		}

		fmt.Printf("=> %s\n", result.String())
		inputBuffer.Reset()
	}
}

// Simple parser for the minimal kernel
func (r *REPL) parse(input string) (Value, error) {
	tokens := r.tokenize(input)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	result, _, err := r.parseExpression(tokens, 0)
	return result, err
}

// Parse exposes the parser for external use
func (r *REPL) Parse(input string) (Value, error) {
	return r.parse(input)
}

func (r *REPL) parseExpression(tokens []string, pos int) (Value, int, error) {
	if pos >= len(tokens) {
		return nil, pos, fmt.Errorf("unexpected end of input")
	}

	token := tokens[pos]

	if token == "(" {
		// Parse list
		elements, newPos, err := r.parseListElements(tokens, pos+1)
		if err != nil {
			return nil, newPos, err
		}
		return NewList(elements...), newPos, nil
	}

	if token == "[" {
		// Parse vector
		elements, newPos, err := r.parseVectorElements(tokens, pos+1)
		if err != nil {
			return nil, newPos, err
		}
		return NewVector(elements...), newPos, nil
	}

	if token == "`" {
		// Parse quasiquote - `expr becomes (quasiquote expr)
		expr, newPos, err := r.parseExpression(tokens, pos+1)
		if err != nil {
			return nil, newPos, err
		}
		return NewList(Intern("quasiquote"), expr), newPos, nil
	}

	if token == "~" {
		// Parse unquote - ~expr becomes (unquote expr)
		expr, newPos, err := r.parseExpression(tokens, pos+1)
		if err != nil {
			return nil, newPos, err
		}
		return NewList(Intern("unquote"), expr), newPos, nil
	}

	// Parse atom
	return r.parseAtom(token), pos + 1, nil
}

func (r *REPL) parseListElements(tokens []string, pos int) ([]Value, int, error) {
	var elements []Value

	for pos < len(tokens) {
		if tokens[pos] == ")" {
			return elements, pos + 1, nil
		}

		expr, newPos, err := r.parseExpression(tokens, pos)
		if err != nil {
			return nil, newPos, err
		}
		elements = append(elements, expr)
		pos = newPos
	}

	return nil, pos, fmt.Errorf("unclosed list")
}

func (r *REPL) parseVectorElements(tokens []string, pos int) ([]Value, int, error) {
	var elements []Value

	for pos < len(tokens) {
		if tokens[pos] == "]" {
			return elements, pos + 1, nil
		}

		expr, newPos, err := r.parseExpression(tokens, pos)
		if err != nil {
			return nil, newPos, err
		}
		elements = append(elements, expr)
		pos = newPos
	}

	return nil, pos, fmt.Errorf("unclosed vector")
}

func (r *REPL) tokenize(input string) []string {
	// Enhanced tokenizer that properly handles string literals
	var tokens []string
	var current strings.Builder
	inString := false

	for i, char := range input {
		if char == '"' {
			// Handle string start/end
			current.WriteRune(char)
			if !inString {
				// Starting a string
				inString = true
			} else {
				// Ending a string (check for escape)
				if i > 0 && rune(input[i-1]) != '\\' {
					inString = false
					// Complete string token
					tokens = append(tokens, current.String())
					current.Reset()
				}
			}
		} else if inString {
			// Inside string - add everything including spaces
			current.WriteRune(char)
		} else {
			// Outside string - normal tokenization
			switch char {
			case '(', ')', '[', ']', '`', '~':
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
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func (r *REPL) parseTokens(tokens []string) (Value, error) {
	result, _, err := r.parseExpression(tokens, 0)
	return result, err
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

// CallWithContext calls the builtin function with evaluation context
func (bf *BuiltinFunction) CallWithContext(args []Value, env *Environment, ctx *EvaluationContext) (Value, error) {
	ctx.PushFrame(fmt.Sprintf("calling builtin: %s", bf.Name))
	result, err := bf.Fn(args, env)
	ctx.PopFrame()

	if err != nil {
		return nil, ctx.WrapError(err)
	}
	return result, nil
}

func (bf *BuiltinFunction) String() string {
	return fmt.Sprintf("<builtin:%s>", bf.Name)
}

// isBalanced checks if parentheses and brackets are balanced in the input
func (r *REPL) isBalanced(input string) bool {
	stack := make([]rune, 0)
	inString := false

	for i, char := range input {
		if char == '"' {
			// Handle string start/end (check for escape)
			if !inString {
				inString = true
			} else if i == 0 || rune(input[i-1]) != '\\' {
				inString = false
			}
			continue
		}

		if inString {
			// Skip characters inside strings
			continue
		}

		switch char {
		case '(', '[':
			stack = append(stack, char)
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false // Unmatched closing paren
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return false // Unmatched closing bracket
			}
			stack = stack[:len(stack)-1]
		}
	}

	// All parentheses and brackets should be matched
	return len(stack) == 0 && !inString
}
