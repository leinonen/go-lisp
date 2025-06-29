package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// REPL represents a Read-Eval-Print-Loop
type REPL struct {
	env *Environment
	ctx *EvaluationContext
}

// NewREPL creates a new REPL with bootstrapped environment
func NewREPL() (*REPL, error) {
	env, err := CreateBootstrappedEnvironment()
	if err != nil {
		return nil, err
	}

	return &REPL{
		env: env,
		ctx: NewEvaluationContext(),
	}, nil
}

// Run starts the REPL
func (r *REPL) Run() error {
	fmt.Println("GoLisp Minimal Core REPL")
	fmt.Println("Type 'exit' to quit")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("core> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			break
		}

		result, err := r.Eval(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("%s\n", result.String())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("REPL error: %v", err)
	}

	return nil
}

// Eval evaluates a string expression
func (r *REPL) Eval(input string) (Value, error) {
	// Parse the input
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	// Evaluate the expression with context
	return EvalWithContext(expr, r.env, r.ctx)
}

// LoadFile loads and evaluates a Lisp file
func (r *REPL) LoadFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	// Parse the file content
	lexer := NewLexer(string(content))
	tokens, err := lexer.Tokenize()
	if err != nil {
		return fmt.Errorf("failed to tokenize file %s: %v", filename, err)
	}

	parser := NewParser(tokens)
	expressions, err := parser.ParseAll()
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %v", filename, err)
	}

	// Set the file context for better error reporting
	r.ctx.Position.File = filename
	
	// Evaluate each expression
	for _, expr := range expressions {
		_, err := EvalWithContext(expr, r.env, r.ctx)
		if err != nil {
			return fmt.Errorf("failed to evaluate expression in file %s: %v", filename, err)
		}
	}

	return nil
}

// EvalString evaluates a string and returns the result
func (r *REPL) EvalString(input string) (Value, error) {
	return r.Eval(input)
}
