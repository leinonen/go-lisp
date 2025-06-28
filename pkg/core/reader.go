package core

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Token represents a token
type Token struct {
	Type  TokenType
	Value string
}

// TokenType represents the type of a token
type TokenType int

const (
	TokenSymbol TokenType = iota
	TokenNumber
	TokenString
	TokenKeyword
	TokenLeftParen
	TokenRightParen
	TokenLeftBracket
	TokenRightBracket
	TokenQuote
	TokenEOF
)

// Lexer tokenizes input
type Lexer struct {
	input    string
	position int
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:    input,
		position: 0,
	}
}

// Tokenize converts input into tokens
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for l.position < len(l.input) {
		// Skip whitespace
		if unicode.IsSpace(l.current()) {
			l.advance()
			continue
		}

		// Skip comments
		if l.current() == ';' {
			l.skipComment()
			continue
		}

		token, err := l.nextToken()
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	// Add EOF token
	tokens = append(tokens, Token{Type: TokenEOF})
	return tokens, nil
}

func (l *Lexer) current() rune {
	if l.position >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position])
}

func (l *Lexer) advance() {
	l.position++
}

func (l *Lexer) skipComment() {
	for l.position < len(l.input) && l.current() != '\n' {
		l.advance()
	}
}

func (l *Lexer) nextToken() (Token, error) {
	char := l.current()

	switch char {
	case '(':
		l.advance()
		return Token{Type: TokenLeftParen, Value: "("}, nil
	case ')':
		l.advance()
		return Token{Type: TokenRightParen, Value: ")"}, nil
	case '[':
		l.advance()
		return Token{Type: TokenLeftBracket, Value: "["}, nil
	case ']':
		l.advance()
		return Token{Type: TokenRightBracket, Value: "]"}, nil
	case '\'':
		l.advance()
		return Token{Type: TokenQuote, Value: "'"}, nil
	case '"':
		return l.readString()
	case ':':
		return l.readKeyword()
	default:
		if unicode.IsDigit(char) || (char == '-' && unicode.IsDigit(l.peek())) {
			return l.readNumber()
		}
		if isSymbolStart(char) {
			return l.readSymbol()
		}
		return Token{}, fmt.Errorf("unexpected character: %c", char)
	}
}

func (l *Lexer) peek() rune {
	if l.position+1 >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position+1])
}

func (l *Lexer) readString() (Token, error) {
	l.advance() // Skip opening quote
	start := l.position

	for l.position < len(l.input) && l.current() != '"' {
		if l.current() == '\\' {
			l.advance() // Skip escape character
		}
		l.advance()
	}

	if l.position >= len(l.input) {
		return Token{}, fmt.Errorf("unterminated string")
	}

	value := l.input[start:l.position]
	l.advance() // Skip closing quote

	// Basic escape handling
	value = strings.ReplaceAll(value, "\\\"", "\"")
	value = strings.ReplaceAll(value, "\\n", "\n")
	value = strings.ReplaceAll(value, "\\t", "\t")
	value = strings.ReplaceAll(value, "\\\\", "\\")

	return Token{Type: TokenString, Value: value}, nil
}

func (l *Lexer) readKeyword() (Token, error) {
	l.advance() // Skip ':'
	start := l.position

	for l.position < len(l.input) && isSymbolChar(l.current()) {
		l.advance()
	}

	value := l.input[start:l.position]
	return Token{Type: TokenKeyword, Value: value}, nil
}

func (l *Lexer) readNumber() (Token, error) {
	start := l.position

	if l.current() == '-' {
		l.advance()
	}

	for l.position < len(l.input) && unicode.IsDigit(l.current()) {
		l.advance()
	}

	// Handle decimal point
	if l.position < len(l.input) && l.current() == '.' {
		l.advance()
		for l.position < len(l.input) && unicode.IsDigit(l.current()) {
			l.advance()
		}
	}

	value := l.input[start:l.position]
	return Token{Type: TokenNumber, Value: value}, nil
}

func (l *Lexer) readSymbol() (Token, error) {
	start := l.position

	for l.position < len(l.input) && isSymbolChar(l.current()) {
		l.advance()
	}

	value := l.input[start:l.position]
	return Token{Type: TokenSymbol, Value: value}, nil
}

func isSymbolStart(char rune) bool {
	return unicode.IsLetter(char) || char == '_' || char == '+' || char == '-' ||
		char == '*' || char == '/' || char == '=' || char == '<' || char == '>' ||
		char == '!' || char == '?' || char == '%'
}

func isSymbolChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' ||
		char == '-' || char == '+' || char == '*' || char == '/' || char == '=' ||
		char == '<' || char == '>' || char == '!' || char == '?' || char == '%'
}

// Parser converts tokens to AST
type Parser struct {
	tokens   []Token
	position int
}

// NewParser creates a new parser
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

// Parse converts tokens to Lisp values
func (p *Parser) Parse() (Value, error) {
	if p.position >= len(p.tokens) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	return p.parseExpression()
}

// ParseAll parses all expressions from tokens
func (p *Parser) ParseAll() ([]Value, error) {
	var expressions []Value

	for p.position < len(p.tokens) && p.tokens[p.position].Type != TokenEOF {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}

func (p *Parser) parseExpression() (Value, error) {
	token := p.tokens[p.position]

	switch token.Type {
	case TokenLeftParen:
		return p.parseList()
	case TokenLeftBracket:
		return p.parseVector()
	case TokenQuote:
		p.position++
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return NewList(Intern("quote"), expr), nil
	case TokenSymbol:
		p.position++
		return Intern(token.Value), nil
	case TokenKeyword:
		p.position++
		return InternKeyword(token.Value), nil
	case TokenString:
		p.position++
		return String(token.Value), nil
	case TokenNumber:
		p.position++
		return p.parseNumber(token.Value)
	case TokenEOF:
		return nil, fmt.Errorf("unexpected end of input")
	default:
		return nil, fmt.Errorf("unexpected token: %s", token.Value)
	}
}

func (p *Parser) parseList() (Value, error) {
	p.position++ // Skip '('

	var elements []Value

	for p.position < len(p.tokens) && p.tokens[p.position].Type != TokenRightParen {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, expr)
	}

	if p.position >= len(p.tokens) {
		return nil, fmt.Errorf("unterminated list")
	}

	p.position++ // Skip ')'
	return NewList(elements...), nil
}

func (p *Parser) parseVector() (Value, error) {
	p.position++ // Skip '['

	var elements []Value

	for p.position < len(p.tokens) && p.tokens[p.position].Type != TokenRightBracket {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, expr)
	}

	if p.position >= len(p.tokens) {
		return nil, fmt.Errorf("unterminated vector")
	}

	p.position++ // Skip ']'
	return NewVector(elements...), nil
}

func (p *Parser) parseNumber(value string) (Value, error) {
	if strings.Contains(value, ".") {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float: %s", value)
		}
		return NewNumber(f), nil
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer: %s", value)
	}
	return NewNumber(i), nil
}

// ReadString parses a string into a Lisp value
func ReadString(input string) (Value, error) {
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	return parser.Parse()
}
