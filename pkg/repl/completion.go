// Package repl provides completion functionality for the REPL
package repl

import (
	"sort"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// CompletionProvider provides tab completion functionality for the REPL
type CompletionProvider struct {
	env *evaluator.Environment
}

// NewCompletionProvider creates a new completion provider
func NewCompletionProvider(env *evaluator.Environment) *CompletionProvider {
	return &CompletionProvider{env: env}
}

// CompletionContext represents the context where completion is happening
type CompletionContext struct {
	inFunctionPosition bool   // true if we're in a position where a function name is expected
	afterOpenParen     bool   // true if we're right after an opening parenthesis
	parenDepth         int    // current parenthesis nesting depth
	inSpecialFunction  string // name of special function if we're completing its arguments (e.g., "builtins")
}

// GetCompletions returns a list of possible completions for the given prefix
func (cp *CompletionProvider) GetCompletions(line string, pos int) []string {
	// Extract the current word being typed and analyze context
	prefix := cp.extractCurrentWord(line, pos)
	context := cp.analyzeContext(line, pos)

	// Check if we're in a special function that takes function names as arguments
	if context.inSpecialFunction != "" {
		return cp.getSpecialFunctionCompletions(context.inSpecialFunction, prefix)
	}

	// Only provide completions if we're in a function position (after '(')
	if !context.inFunctionPosition {
		return nil
	}

	var completions []string

	// Get built-in functions first
	builtins := cp.getBuiltinFunctions()
	for _, builtin := range builtins {
		if strings.HasPrefix(builtin, prefix) {
			completions = append(completions, builtin)
		}
	}

	// Get user-defined functions
	userFunctions := cp.getUserDefinedFunctions()
	for _, funcName := range userFunctions {
		if strings.HasPrefix(funcName, prefix) {
			completions = append(completions, funcName)
		}
	}

	// Get functions from loaded modules
	moduleFunctions := cp.getModuleFunctions()
	for _, modFunc := range moduleFunctions {
		if strings.HasPrefix(modFunc, prefix) {
			completions = append(completions, modFunc)
		}
	}

	// Remove duplicates and sort
	completions = cp.removeDuplicates(completions)
	sort.Strings(completions)

	return completions
}

// extractCurrentWord extracts the word being completed from the input line
func (cp *CompletionProvider) extractCurrentWord(line string, pos int) string {
	if pos > len(line) {
		pos = len(line)
	}

	// Find the start of the current word
	start := pos
	for start > 0 && cp.isSymbolChar(rune(line[start-1])) {
		start--
	}

	// Find the end of the current word
	end := pos
	for end < len(line) && cp.isSymbolChar(rune(line[end])) {
		end++
	}

	word := line[start:end]

	// Always return the word (even if empty) since we'll check context separately
	return word
}

// isSymbolChar checks if a character can be part of a Lisp symbol
func (cp *CompletionProvider) isSymbolChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '-' || ch == '_' || ch == '?' || ch == '!' ||
		ch == '+' || ch == '*' || ch == '/' || ch == '=' ||
		ch == '<' || ch == '>' || ch == '.' || ch == '%'
}

// getBuiltinFunctions returns a list of all built-in function names
func (cp *CompletionProvider) getBuiltinFunctions() []string {
	return []string{
		// Arithmetic operations
		"+", "-", "*", "/", "%",
		// Comparison operations
		"=", "<", ">", "<=", ">=",
		// Logical operations
		"and", "or", "not",
		// Control flow
		"if",
		// Variable and function definition
		"define", "lambda", "defun",
		// Macro system
		"defmacro", "quote",
		// List operations
		"list", "first", "rest", "cons", "length", "empty?",
		// Higher-order functions
		"map", "filter", "reduce",
		// List manipulation
		"append", "reverse", "nth",
		// Hash map operations
		"hash-map", "hash-map-get", "hash-map-put", "hash-map-remove",
		"hash-map-contains?", "hash-map-keys", "hash-map-values",
		"hash-map-size", "hash-map-empty?",
		// String operations
		"string-concat", "string-length", "string-substring", "string-char-at",
		"string-upper", "string-lower", "string-trim", "string-split", "string-join",
		"string-contains?", "string-starts-with?", "string-ends-with?", "string-replace",
		"string-index-of", "string->number", "number->string", "string-regex-match?",
		"string-regex-find-all", "string-repeat", "string?", "string-empty?",
		// Module system
		"load", "import",
		// Environment inspection
		"env", "modules", "builtins",
		// Print functions
		"print", "println",
		// Constants
		"nil",
		// Error handling
		"error",
	}
}

// getUserDefinedSymbols returns all user-defined symbols from the environment
func (cp *CompletionProvider) getUserDefinedSymbols() []string {
	var symbols []string

	// Walk through the environment chain to collect all bindings
	env := cp.env
	for env != nil {
		for name, value := range env.GetBindings() {
			// Include all symbols, but prefer functions
			switch value.(type) {
			case types.FunctionValue, *types.BuiltinFunctionValue, *types.ArithmeticFunctionValue:
				// Functions get priority
				symbols = append(symbols, name)
			default:
				// Variables and other values
				symbols = append(symbols, name)
			}
		}
		env = env.GetParent()
	}

	return symbols
}

// getUserDefinedFunctions returns only user-defined functions from the environment
func (cp *CompletionProvider) getUserDefinedFunctions() []string {
	var functions []string

	// Walk through the environment chain to collect function bindings
	env := cp.env
	for env != nil {
		for name, value := range env.GetBindings() {
			// Include only function values
			switch value.(type) {
			case types.FunctionValue:
				functions = append(functions, name)
			}
		}
		env = env.GetParent()
	}

	return functions
}

// getModuleFunctions returns functions from loaded modules
func (cp *CompletionProvider) getModuleFunctions() []string {
	var functions []string

	// Get all modules
	modules := cp.env.GetModules()
	for moduleName, module := range modules {
		// Add qualified function names (module.function)
		for exportName := range module.Exports {
			functions = append(functions, moduleName+"."+exportName)
		}
	}

	return functions
}

// getSpecialFunctionCompletions returns completions for arguments to special functions
func (cp *CompletionProvider) getSpecialFunctionCompletions(functionName, prefix string) []string {
	var completions []string

	switch functionName {
	case "builtins":
		// For builtins function, complete with all built-in function names
		builtins := cp.getBuiltinFunctions()
		for _, builtin := range builtins {
			if strings.HasPrefix(builtin, prefix) {
				completions = append(completions, builtin)
			}
		}
	default:
		// For other functions, no special completion
		return nil
	}

	// Remove duplicates and sort
	completions = cp.removeDuplicates(completions)
	sort.Strings(completions)

	return completions
}

// findSpecialFunction checks if we're completing arguments for a special function like "builtins"
func (cp *CompletionProvider) findSpecialFunction(line string, wordStart int) string {
	// Look backwards to find the function name at the start of the current expression

	// Find the most recent '(' before our position
	lastParen := -1
	inString := false
	escaped := false

	for i := 0; i < wordStart; i++ {
		ch := rune(line[i])

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
				lastParen = i
			}
		case ')':
			if !inString {
				lastParen = -1 // Reset when we see a closing paren
			}
		}
	}

	if lastParen == -1 {
		return ""
	}

	// Extract the function name after the '('
	funcStart := lastParen + 1
	// Skip whitespace
	for funcStart < len(line) && (line[funcStart] == ' ' || line[funcStart] == '\t') {
		funcStart++
	}

	if funcStart >= len(line) {
		return ""
	}

	// Find the end of the function name
	funcEnd := funcStart
	for funcEnd < len(line) && cp.isSymbolChar(rune(line[funcEnd])) {
		funcEnd++
	}

	if funcEnd <= funcStart {
		return ""
	}

	functionName := line[funcStart:funcEnd]

	// Check if we're past the function name (in argument position)
	// The word we're completing should start after the function name
	if wordStart > funcEnd {
		// Check if this is a special function that takes function names as arguments
		switch functionName {
		case "builtins":
			// For builtins, only complete the first argument
			// Count how many completed arguments we've already seen (excluding the current word being completed)
			argCount := 0
			pos := funcEnd

			// Skip to the start of arguments
			for pos < wordStart && (line[pos] == ' ' || line[pos] == '\t') {
				pos++
			}

			// Count completed arguments by counting symbol boundaries
			// We need to be careful not to count the word we're currently completing
			for pos < wordStart {
				// Skip whitespace
				for pos < wordStart && (line[pos] == ' ' || line[pos] == '\t') {
					pos++
				}

				if pos < wordStart {
					// Found start of a symbol
					symbolEnd := pos

					// Skip to end of this symbol
					for symbolEnd < len(line) && cp.isSymbolChar(rune(line[symbolEnd])) {
						symbolEnd++
					}

					// Only count this as a completed argument if the symbol ends before our word starts
					// This excludes the word we're currently completing
					if symbolEnd <= wordStart {
						argCount++
					}

					// Move position to after this symbol
					pos = symbolEnd
				}
			}

			// Only provide completions for the first argument
			// argCount == 0 means we're completing the first argument
			// argCount == 1 means we're completing the second argument (which should not be allowed)
			if argCount == 0 {
				return "builtins"
			}
			return ""
		default:
			return ""
		}
	}

	return ""
}

// analyzeContext analyzes the completion context based on the input line and position
func (cp *CompletionProvider) analyzeContext(line string, pos int) CompletionContext {
	context := CompletionContext{}

	if pos > len(line) {
		pos = len(line)
	}

	// Count parentheses and analyze position
	parenCount := 0
	inString := false
	escaped := false

	for i := 0; i < pos; i++ {
		ch := rune(line[i])

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
		}
	}

	context.parenDepth = parenCount

	// Extract current word and find its start position
	wordStart := pos
	for wordStart > 0 && cp.isSymbolChar(rune(line[wordStart-1])) {
		wordStart--
	}

	// Check if we're in an argument position for a special function
	specialFunc := cp.findSpecialFunction(line, wordStart)
	if specialFunc != "" {
		context.inSpecialFunction = specialFunc
		return context
	}

	// Look backwards from word start to find the most recent '('
	searchPos := wordStart - 1
	for searchPos >= 0 && (line[searchPos] == ' ' || line[searchPos] == '\t') {
		searchPos--
	}

	// Check if we found a '(' immediately before (possibly with whitespace)
	if searchPos >= 0 && line[searchPos] == '(' {
		context.afterOpenParen = (wordStart == searchPos+1) // true if no space between '(' and word
		context.inFunctionPosition = true
	} else {
		// Look for the pattern "( symbol" to see if we're still in function position
		// Find the most recent '(' before our position
		lastParen := -1
		inStr := false
		esc := false

		for i := 0; i < wordStart; i++ {
			ch := rune(line[i])

			if esc {
				esc = false
				continue
			}

			switch ch {
			case '\\':
				if inStr {
					esc = true
				}
			case '"':
				inStr = !inStr
			case '(':
				if !inStr {
					lastParen = i
				}
			case ')':
				if !inStr {
					lastParen = -1 // Reset when we see a closing paren
				}
			}
		}

		if lastParen >= 0 {
			// Check if there are any symbols between the '(' and our position
			symbolCount := 0
			tempPos := lastParen + 1

			for tempPos < wordStart {
				// Skip whitespace
				for tempPos < wordStart && (line[tempPos] == ' ' || line[tempPos] == '\t') {
					tempPos++
				}

				if tempPos < wordStart {
					// Found a symbol, count it
					symbolCount++
					// Skip to end of this symbol
					for tempPos < wordStart && cp.isSymbolChar(rune(line[tempPos])) {
						tempPos++
					}
				}
			}

			// We're in function position if this is the first symbol after '('
			context.inFunctionPosition = (symbolCount == 0)
		}
	}

	return context
}

// removeDuplicates removes duplicate strings from a slice
func (cp *CompletionProvider) removeDuplicates(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// lispCompleter implements readline.AutoCompleter for Lisp-aware completion
type lispCompleter struct {
	provider *CompletionProvider
}

// NewLispCompleter creates a new Lisp-aware completer
func NewLispCompleter(provider *CompletionProvider) *lispCompleter {
	return &lispCompleter{provider: provider}
}

// Do implements the readline.AutoCompleter interface
func (lc *lispCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	// Convert rune slice to string for processing
	lineStr := string(line)

	// Get completions
	completions := lc.provider.GetCompletions(lineStr, pos)

	if len(completions) == 0 {
		return nil, 0
	}

	// Extract the current word to determine how much to replace
	currentWord := lc.provider.extractCurrentWord(lineStr, pos)
	replaceLength := len(currentWord)

	// Convert completions to rune slices
	var suggestions [][]rune
	for _, completion := range completions {
		// Only include the part that extends beyond the current prefix
		if len(completion) > len(currentWord) {
			extension := completion[len(currentWord):]
			suggestions = append(suggestions, []rune(extension))
		} else if completion == currentWord {
			// Exact match - add as-is for display
			suggestions = append(suggestions, []rune(completion))
		}
	}

	return suggestions, replaceLength
}
