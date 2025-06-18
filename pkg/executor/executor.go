package executor

import (
	"fmt"
	"io/ioutil"

	"github.com/leinonen/go-lisp/pkg/interpreter"
	"github.com/leinonen/go-lisp/pkg/tokenizer"
	"github.com/leinonen/go-lisp/pkg/types"
)

// ExecuteFile reads and executes a Lisp file using the provided interpreter
func ExecuteFile(interpreter *interpreter.Interpreter, filename string) error {
	// Read the file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Tokenize the content
	tokenizer := tokenizer.NewTokenizer(string(content))
	tokens, err := tokenizer.TokenizeWithError()
	if err != nil {
		return fmt.Errorf("tokenization error in %s: %v", filename, err)
	}

	// Parse and evaluate each expression in the file
	i := 0
	for i < len(tokens) {
		// Find the end of the current expression
		if tokens[i].Type == types.TokenType(-1) { // EOF
			break
		}

		// Extract tokens for this expression
		exprTokens, newIndex := extractExpression(tokens, i)
		if len(exprTokens) == 0 {
			break
		}

		// Prevent infinite loop - if we didn't advance, move forward by 1
		if newIndex <= i {
			i++
			continue
		}

		// Convert tokens back to string and interpret
		exprString := tokensToString(exprTokens)
		if exprString != "" {
			result, err := interpreter.Interpret(exprString)
			if err != nil {
				return fmt.Errorf("evaluation error in %s: %v", filename, err)
			}
			// Print the result of each expression (suppress nil values from print functions)
			if result.String() != "nil" {
				fmt.Println(result)
			}
		}

		i = newIndex
	}

	return nil
}

// extractExpression extracts a complete expression from tokens
func extractExpression(tokens []types.Token, start int) ([]types.Token, int) {
	if start >= len(tokens) {
		return nil, start
	}

	// Handle single token expressions (numbers, strings, booleans, symbols)
	if tokens[start].Type != types.LPAREN {
		return tokens[start : start+1], start + 1
	}

	// Handle list expressions - find matching closing paren
	parenCount := 0
	end := start
	for end < len(tokens) {
		if tokens[end].Type == types.LPAREN {
			parenCount++
		} else if tokens[end].Type == types.RPAREN {
			parenCount--
			if parenCount == 0 {
				end++
				break
			}
		}
		end++
	}

	if parenCount != 0 {
		// Unmatched parentheses - return what we have
		return tokens[start:end], end
	}

	return tokens[start:end], end
}

// tokensToString converts tokens back to string
func tokensToString(tokens []types.Token) string {
	var result string
	for i, token := range tokens {
		if i > 0 {
			// Add space between tokens (except for parentheses)
			if token.Type != types.RPAREN && tokens[i-1].Type != types.LPAREN {
				result += " "
			}
		}

		// Handle different token types appropriately
		switch token.Type {
		case types.STRING:
			// Add quotes back for string tokens
			result += `"` + token.Value + `"`
		case types.KEYWORD:
			// Add colon prefix back for keyword tokens
			result += ":" + token.Value
		default:
			result += token.Value
		}
	}
	return result
}
