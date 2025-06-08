package repl

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Interpreter interface for dependency injection in tests
type Interpreter interface {
	Interpret(input string) (types.Value, error)
}

// REPL starts a Read-Eval-Print Loop for the Lisp interpreter
func REPL(interp Interpreter, scanner *bufio.Scanner) {
	REPLWithOptions(interp, scanner, true)
}

// REPLWithOptions starts a REPL with configurable options
func REPLWithOptions(interp Interpreter, scanner *bufio.Scanner, enableColors bool) {
	// Configure colors based on the enableColors parameter
	if !enableColors {
		color.NoColor = true
		printWelcomeMessageNoColor()
	} else {
		printWelcomeMessage()
	}

	// Create error formatter for colored output
	errorFormatter := NewErrorFormatter()

	for {
		input := readCompleteExpressionWithColors(scanner, enableColors)
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
			// Use colored error formatting with smart suggestions
			fmt.Println(errorFormatter.FormatErrorWithSmartSuggestion(err))
		} else {
			// Color the result output
			resultColor := color.New(color.FgGreen)
			fmt.Printf("=> %s\n", resultColor.Sprint(result.String()))
		}
	}

	if enableColors {
		printGoodbyeMessage()
	} else {
		printGoodbyeMessageNoColor()
	}
}

// printWelcomeMessage prints a welcome message and instructions for the REPL
func printWelcomeMessage() {
	titleColor := color.New(color.FgCyan, color.Bold)
	instructionColor := color.New(color.FgYellow)
	commandColor := color.New(color.FgGreen)

	titleColor.Println("Welcome to the Lisp Interpreter!")
	instructionColor.Println("Type expressions to evaluate them, or 'quit' to exit.")
	instructionColor.Println("Multi-line expressions are supported - the REPL will wait for balanced parentheses.")
	fmt.Println()
	commandColor.Println("Helpful commands:")
	fmt.Println("  (builtins)        - List all available built-in functions")
	fmt.Println("  (builtins <func>) - Get help for a specific function")
	fmt.Println("  (env)             - Show current environment variables")
	fmt.Println("  (modules)         - Show loaded modules")
	fmt.Println()
	instructionColor.Println("✨ Errors are now color-coded by type with helpful suggestions!")
	fmt.Println()
}

// printGoodbyeMessage prints a goodbye message when the REPL ends
func printGoodbyeMessage() {
	goodbyeColor := color.New(color.FgMagenta, color.Bold)
	goodbyeColor.Println("Goodbye!")
}

// printWelcomeMessageNoColor prints welcome message without colors (for testing)
func printWelcomeMessageNoColor() {
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
	fmt.Println("✨ Errors are now color-coded by type with helpful suggestions!")
	fmt.Println()
}

// printGoodbyeMessageNoColor prints goodbye message without colors (for testing)
func printGoodbyeMessageNoColor() {
	fmt.Println("Goodbye!")
}

// readCompleteExpression reads input until we have a complete s-expression
// with balanced parentheses, or until the user enters a simple command
func readCompleteExpression(scanner *bufio.Scanner) string {
	return readCompleteExpressionWithColors(scanner, true)
}

// readCompleteExpressionWithColors reads input with optional colored prompts
func readCompleteExpressionWithColors(scanner *bufio.Scanner, enableColors bool) string {
	var lines []string
	parenCount := 0
	inString := false
	escaped := false
	isFirstLine := true

	// Colors for prompts
	primaryPromptColor := color.New(color.FgBlue, color.Bold)
	continuationPromptColor := color.New(color.FgHiBlack)

	for {
		if isFirstLine {
			if enableColors {
				primaryPromptColor.Print("lisp> ")
			} else {
				fmt.Print("lisp> ")
			}
			isFirstLine = false
		} else {
			if enableColors {
				continuationPromptColor.Print("...   ")
			} else {
				fmt.Print("...   ")
			}
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
