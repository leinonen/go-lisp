package tokenizer

import (
	"fmt"
	"unicode"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

type Tokenizer struct {
	input    string
	position int
	current  rune
}

func NewTokenizer(input string) *Tokenizer {
	t := &Tokenizer{
		input: input,
	}
	t.readChar()
	return t
}

func (t *Tokenizer) readChar() {
	if t.position >= len(t.input) {
		t.current = 0 // ASCII NUL character represents "EOF"
	} else {
		t.current = rune(t.input[t.position])
	}
	t.position++
}

func (t *Tokenizer) peekChar() rune {
	if t.position >= len(t.input) {
		return 0
	}
	return rune(t.input[t.position])
}

func (t *Tokenizer) skipWhitespace() {
	for unicode.IsSpace(t.current) {
		t.readChar()
	}
}

func (t *Tokenizer) skipComment() {
	// Skip until end of line or end of input
	for t.current != '\n' && t.current != 0 {
		t.readChar()
	}
}

func (t *Tokenizer) readString() (string, error) {
	position := t.position
	for {
		t.readChar()
		if t.current == '"' || t.current == 0 {
			break
		}
	}

	if t.current == 0 {
		return "", fmt.Errorf("unterminated string")
	}

	return t.input[position : t.position-1], nil
}

func (t *Tokenizer) readNumber() string {
	position := t.position - 1

	// Handle negative numbers
	if t.current == '-' && unicode.IsDigit(t.peekChar()) {
		t.readChar()
	}

	for unicode.IsDigit(t.current) || t.current == '.' {
		t.readChar()
	}

	return t.input[position : t.position-1]
}

func (t *Tokenizer) readSymbol() string {
	position := t.position - 1

	for isSymbolChar(t.current) {
		t.readChar()
	}

	return t.input[position : t.position-1]
}

func isSymbolChar(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) ||
		ch == '+' || ch == '-' || ch == '*' || ch == '/' ||
		ch == '=' || ch == '<' || ch == '>' || ch == '!' ||
		ch == '?' || ch == '_' || ch == '.' || ch == '%'
}

func (t *Tokenizer) Tokenize() []types.Token {
	tokens, _ := t.TokenizeWithError()
	return tokens
}

func (t *Tokenizer) TokenizeWithError() ([]types.Token, error) {
	tokens := make([]types.Token, 0)

	for t.current != 0 {
		t.skipWhitespace()

		if t.current == 0 {
			break
		}

		// Handle comments
		if t.current == ';' {
			t.skipComment()
			continue
		}

		switch t.current {
		case '(':
			tokens = append(tokens, types.Token{Type: types.LPAREN, Value: "("})
			t.readChar()
		case ')':
			tokens = append(tokens, types.Token{Type: types.RPAREN, Value: ")"})
			t.readChar()
		case '[':
			tokens = append(tokens, types.Token{Type: types.LBRACKET, Value: "["})
			t.readChar()
		case ']':
			tokens = append(tokens, types.Token{Type: types.RBRACKET, Value: "]"})
			t.readChar()
		case '\'':
			tokens = append(tokens, types.Token{Type: types.QUOTE, Value: "'"})
			t.readChar()
		case '"':
			str, err := t.readString()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, types.Token{Type: types.STRING, Value: str})
			t.readChar() // consume closing quote

		case ':':
			// Keywords
			t.readChar() // consume ':'
			if !isSymbolChar(t.current) {
				return nil, fmt.Errorf("invalid keyword: colon must be followed by symbol characters")
			}
			keyword := t.readSymbol()
			if keyword == "" {
				return nil, fmt.Errorf("invalid keyword: empty keyword name")
			}
			tokens = append(tokens, types.Token{Type: types.KEYWORD, Value: keyword})
		default:
			if unicode.IsDigit(t.current) || (t.current == '-' && unicode.IsDigit(t.peekChar())) {
				number := t.readNumber()
				tokens = append(tokens, types.Token{Type: types.NUMBER, Value: number})
			} else if isSymbolChar(t.current) {
				symbol := t.readSymbol()
				// Check for boolean literals
				if symbol == "true" || symbol == "false" {
					tokens = append(tokens, types.Token{Type: types.BOOLEAN, Value: symbol})
				} else {
					tokens = append(tokens, types.Token{Type: types.SYMBOL, Value: symbol})
				}
			} else {
				return nil, fmt.Errorf("invalid character: %c", t.current)
			}
		}
	}

	return tokens, nil
}
