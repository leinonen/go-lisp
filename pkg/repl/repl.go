package repl

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Interpreter interface for dependency injection in tests
type Interpreter interface {
	Interpret(input string) (types.Value, error)
}

// REPL starts a Read-Eval-Print Loop for the Lisp interpreter
func REPL(interp Interpreter, scanner *bufio.Scanner) {
	printWelcomeMessage()

	for {
		input := readCompleteExpression(scanner)
		if input == "" {
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if input == "quit" || input == "exit" {
			break
		}

		result, err := interp.Interpret(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("=> %s\n", result.String())
		}
	}

	printGoodbyeMessage()
}

// printWelcomeMessage prints a welcome message and instructions for the REPL
func printWelcomeMessage() {
	fmt.Println("Welcome to the Lisp Interpreter!")
	fmt.Println("Type expressions to evaluate them, or 'quit' to exit.")
	fmt.Println("Multi-line expressions are supported - the REPL will wait for balanced parentheses.")
	fmt.Println()
	fmt.Println("Helpful commands:")
	fmt.Println("  (builtins)        - List all available built-in functions")
	fmt.Println("  (builtins <func>) - Get help for a specific function")
	fmt.Println("  (env)             - Show current environment variables")
	fmt.Println("  (modules)         - Show loaded modules")
	fmt.Println()
}

// printGoodbyeMessage prints a goodbye message when the REPL ends
func printGoodbyeMessage() {
	fmt.Println("Goodbye!")
}

// readCompleteExpression reads input until we have a complete s-expression
// with balanced parentheses, or until the user enters a simple command
func readCompleteExpression(scanner *bufio.Scanner) string {
	var lines []string
	parenCount := 0
	inString := false
	escaped := false
	isFirstLine := true

	for {
		if isFirstLine {
			fmt.Print("lisp> ")
			isFirstLine = false
		} else {
			fmt.Print("...   ")
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Printf("Scanner error: %v\n", err)
			}
			return strings.Join(lines, "\n")
		}

		line := scanner.Text()
		lines = append(lines, line)

		// Check if this is a simple quit/exit command
		trimmed := strings.TrimSpace(line)
		if len(lines) == 1 && (trimmed == "quit" || trimmed == "exit") {
			return trimmed
		}

		// Count parentheses, respecting strings and escapes
		for _, ch := range line {
			if escaped {
				escaped = false
				continue
			}

			switch ch {
			case '\\':
				if inString {
					escaped = true
				}
			case '"':
				inString = !inString
			case '(':
				if !inString {
					parenCount++
				}
			case ')':
				if !inString {
					parenCount--
				}
			case ';':
				if !inString {
					// Comment - ignore rest of line
					break
				}
			}
		}

		// If we have balanced parentheses and at least one complete expression, we're done
		if parenCount == 0 && containsExpression(strings.Join(lines, "\n")) {
			break
		}

		// If parentheses count goes negative, we have unmatched closing parens
		if parenCount < 0 {
			break
		}
	}

	return strings.Join(lines, "\n")
}

// containsExpression checks if the input contains at least one meaningful expression
func containsExpression(input string) bool {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return false
	}

	// Remove comments and check if there's anything left
	lines := strings.Split(trimmed, "\n")
	for _, line := range lines {
		// Find comment position, respecting strings
		inString := false
		escaped := false
		for i, ch := range line {
			if escaped {
				escaped = false
				continue
			}

			switch ch {
			case '\\':
				if inString {
					escaped = true
				}
			case '"':
				inString = !inString
			case ';':
				if !inString {
					line = line[:i]
					break
				}
			}
		}

		// Check if line has non-whitespace content
		if strings.TrimSpace(line) != "" {
			return true
		}
	}

	return false
}
