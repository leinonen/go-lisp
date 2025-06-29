package core

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

// createDynamicCompleter creates a completer based on the current environment
func (r *REPL) createDynamicCompleter() readline.AutoCompleter {
	// Get all symbols from the environment
	symbols := r.env.GetAllSymbols()
	
	// Static special forms that always need parentheses
	specialForms := []string{
		"def", "defn", "if", "fn", "let", "do", "loop", "recur",
		"when", "unless", "cond", "quote", "quasiquote", "unquote",
		"unquote-splicing", "defmacro", "macroexpand",
	}
	
	// Static literals that don't need parentheses
	literals := []string{
		"nil", "true", "false", "exit", "quit",
	}
	
	var items []readline.PrefixCompleterInterface
	
	// Add special forms with parentheses
	for _, form := range specialForms {
		items = append(items, readline.PcItem("("+form))
	}
	
	// Add all environment symbols
	for _, symbol := range symbols {
		// Skip if it's a special form (already added)
		isSpecialForm := false
		for _, form := range specialForms {
			if symbol == form {
				isSpecialForm = true
				break
			}
		}
		if isSpecialForm {
			continue
		}
		
		// Check if it's a literal (no parentheses needed)
		isLiteral := false
		for _, literal := range literals {
			if symbol == literal {
				isLiteral = true
				break
			}
		}
		
		if isLiteral {
			items = append(items, readline.PcItem(symbol))
		} else {
			// Most functions/variables get parentheses
			items = append(items, readline.PcItem("("+symbol))
		}
	}
	
	// Add literals
	for _, literal := range literals {
		items = append(items, readline.PcItem(literal))
	}
	
	return readline.NewPrefixCompleter(items...)
}

// updateCompleter refreshes the autocomplete with current environment symbols
// Note: Due to limitations in the readline library, we currently don't update 
// the completer during the session. The completer is set once at startup.
func (r *REPL) updateCompleter() {
	// TODO: Implement real-time completion updates when readline library supports it
	// For now, the dynamic completer is created once at REPL startup
}

// GetEnv returns the REPL's environment (for testing purposes)
func (r *REPL) GetEnv() *Environment {
	return r.env
}

// isBalanced checks if parentheses, brackets, and braces are balanced
func isBalanced(input string) bool {
	stack := 0
	inString := false
	inComment := false
	escapeNext := false
	
	for _, char := range input {
		if escapeNext {
			escapeNext = false
			continue
		}
		
		if char == '\\' && inString {
			escapeNext = true
			continue
		}
		
		// Handle comments (; to end of line)
		if char == ';' && !inString {
			inComment = true
			continue
		}
		
		if char == '\n' {
			inComment = false
			continue
		}
		
		if inComment {
			continue
		}
		
		if char == '"' {
			inString = !inString
			continue
		}
		
		if inString {
			continue
		}
		
		switch char {
		case '(', '[', '{':
			stack++
		case ')', ']', '}':
			stack--
			if stack < 0 {
				return false // More closing than opening
			}
		}
	}
	
	return stack == 0
}

// hasNonWhitespaceContent checks if the input has actual content (not just whitespace/comments)
func hasNonWhitespaceContent(input string) bool {
	inString := false
	inComment := false
	escapeNext := false
	
	for _, char := range input {
		if escapeNext {
			escapeNext = false
			if inString {
				return true // Escaped character in string counts as content
			}
			continue
		}
		
		if char == '\\' && inString {
			escapeNext = true
			continue
		}
		
		if char == ';' && !inString {
			inComment = true
			continue
		}
		
		if char == '\n' {
			inComment = false
			continue
		}
		
		if inComment {
			continue
		}
		
		if char == '"' {
			inString = !inString
			return true // String content counts
		}
		
		if inString {
			return true // Any character in string counts
		}
		
		// Check for non-whitespace characters outside of comments and strings
		if char != ' ' && char != '\t' && char != '\n' && char != '\r' {
			return true
		}
	}
	
	return false
}

// REPL represents a Read-Eval-Print-Loop
type REPL struct {
	env *Environment
	ctx *EvaluationContext
	rl  *readline.Instance
}

// NewREPL creates a new REPL with bootstrapped environment
func NewREPL() (*REPL, error) {
	env, err := CreateBootstrappedEnvironment()
	if err != nil {
		return nil, err
	}

	// Create REPL instance first (we need it to create the dynamic completer)
	repl := &REPL{
		env: env,
		ctx: NewEvaluationContext(),
	}

	// Configure readline with history and completion
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "GoLisp> ",
		AutoComplete:    repl.createDynamicCompleter(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create readline: %v", err)
	}

	// Set the readline instance on the REPL
	repl.rl = rl
	
	return repl, nil
}

// Run starts the REPL
func (r *REPL) Run() error {
	defer r.rl.Close()

	fmt.Println("GoLisp Enhanced REPL")
	fmt.Println("Type 'exit' or 'quit' to quit")
	fmt.Println("Multi-line expressions supported - press Enter on incomplete expressions")
	fmt.Println("Type ')' on empty line during multi-line input to force evaluation")

	var inputBuffer strings.Builder
	isMultiLine := false

	for {
		// Set prompt based on whether we're in multi-line mode
		prompt := "GoLisp> "
		if isMultiLine {
			prompt = "      > "
		}
		r.rl.SetPrompt(prompt)

		line, err := r.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					if isMultiLine {
						// Cancel multi-line input
						inputBuffer.Reset()
						isMultiLine = false
						fmt.Println("Cancelled")
						continue
					}
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			return fmt.Errorf("REPL error: %v", err)
		}

		trimmedLine := strings.TrimSpace(line)

		// Handle exit commands even in multi-line mode
		if !isMultiLine && (trimmedLine == "exit" || trimmedLine == "quit") {
			break
		}

		// Handle force evaluation with ')' on empty line
		if isMultiLine && trimmedLine == ")" {
			currentInput := inputBuffer.String()
			
			// Only try to balance if there's actually content and unclosed parens
			if hasNonWhitespaceContent(currentInput) {
				// Count how many opening brackets we have vs closing ones
				openCount := 0
				closeCount := 0
				inString := false
				inComment := false
				escapeNext := false
				
				for _, char := range currentInput {
					if escapeNext {
						escapeNext = false
						continue
					}
					
					if char == '\\' && inString {
						escapeNext = true
						continue
					}
					
					if char == ';' && !inString {
						inComment = true
						continue
					}
					
					if char == '\n' {
						inComment = false
						continue
					}
					
					if inComment {
						continue
					}
					
					if char == '"' {
						inString = !inString
						continue
					}
					
					if inString {
						continue
					}
					
					switch char {
					case '(', '[', '{':
						openCount++
					case ')', ']', '}':
						closeCount++
					}
				}
				
				// Add closing parens only if we have unclosed opening ones
				if openCount > closeCount {
					for i := 0; i < (openCount - closeCount); i++ {
						currentInput += ")"
					}
				}
				
				// Now evaluate if we have content
				if hasNonWhitespaceContent(currentInput) {
					result, err := r.Eval(currentInput)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					} else {
						fmt.Printf("%s\n", result.String())
						r.updateCompleter()
					}
				}
			}
			
			// Reset for next input
			inputBuffer.Reset()
			isMultiLine = false
			r.rl.SetPrompt("GoLisp> ")
			continue
		}

		// Handle invalid input when not in multi-line mode
		if !isMultiLine && trimmedLine == ")" {
			fmt.Println("Error: Unexpected closing parenthesis")
			continue
		}

		// Skip empty lines when not in multi-line mode
		if !isMultiLine && trimmedLine == "" {
			continue
		}

		// Add line to buffer
		if inputBuffer.Len() > 0 {
			inputBuffer.WriteString("\n")
		}
		inputBuffer.WriteString(line)

		currentInput := inputBuffer.String()

		// Check if expression has content and is balanced
		if hasNonWhitespaceContent(currentInput) && isBalanced(currentInput) {
			// Expression is complete, evaluate it
			result, err := r.Eval(currentInput)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("%s\n", result.String())
				// Update completer after successful evaluation
				r.updateCompleter()
			}
			
			// Reset for next input
			inputBuffer.Reset()
			isMultiLine = false
			r.rl.SetPrompt("GoLisp> ")
		} else if !hasNonWhitespaceContent(currentInput) {
			// Only whitespace/comments, just continue
			continue
		} else {
			// Expression is incomplete, continue multi-line input
			isMultiLine = true
		}
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
