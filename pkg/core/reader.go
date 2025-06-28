package core

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Position represents a position in source code
type Position struct {
	Line   int
	Column int
	Offset int
}

// Token represents a token with source location
type Token struct {
	Type     TokenType
	Value    string
	Position Position
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
	TokenLeftBrace
	TokenRightBrace
	TokenHash
	TokenQuote
	TokenQuasiquote
	TokenUnquote
	TokenUnquoteSplicing
	TokenEOF
)

// Lexer tokenizes input
type Lexer struct {
	input    string
	position int
	line     int
	column   int
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:    input,
		position: 0,
		line:     1,
		column:   1,
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
	tokens = append(tokens, Token{Type: TokenEOF, Position: l.currentPosition()})
	return tokens, nil
}

func (l *Lexer) current() rune {
	if l.position >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position])
}

func (l *Lexer) advance() {
	if l.position < len(l.input) && l.input[l.position] == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	l.position++
}

func (l *Lexer) currentPosition() Position {
	return Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.position,
	}
}

func (l *Lexer) skipComment() {
	for l.position < len(l.input) && l.current() != '\n' {
		l.advance()
	}
}

func (l *Lexer) nextToken() (Token, error) {
	char := l.current()
	pos := l.currentPosition()

	switch char {
	case '(':
		l.advance()
		return Token{Type: TokenLeftParen, Value: "(", Position: pos}, nil
	case ')':
		l.advance()
		return Token{Type: TokenRightParen, Value: ")", Position: pos}, nil
	case '[':
		l.advance()
		return Token{Type: TokenLeftBracket, Value: "[", Position: pos}, nil
	case ']':
		l.advance()
		return Token{Type: TokenRightBracket, Value: "]", Position: pos}, nil
	case '{':
		l.advance()
		return Token{Type: TokenLeftBrace, Value: "{", Position: pos}, nil
	case '}':
		l.advance()
		return Token{Type: TokenRightBrace, Value: "}", Position: pos}, nil
	case '#':
		l.advance()
		return Token{Type: TokenHash, Value: "#", Position: pos}, nil
	case '\'':
		l.advance()
		return Token{Type: TokenQuote, Value: "'", Position: pos}, nil
	case '`':
		l.advance()
		return Token{Type: TokenQuasiquote, Value: "`", Position: pos}, nil
	case '~':
		l.advance()
		// Check for unquote-splicing (~@)
		if l.position < len(l.input) && l.current() == '@' {
			l.advance()
			return Token{Type: TokenUnquoteSplicing, Value: "~@", Position: pos}, nil
		}
		return Token{Type: TokenUnquote, Value: "~", Position: pos}, nil
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
	pos := l.currentPosition()
	l.advance() // Skip opening quote
	start := l.position

	for l.position < len(l.input) && l.current() != '"' {
		if l.current() == '\\' {
			l.advance() // Skip escape character
		}
		l.advance()
	}

	if l.position >= len(l.input) {
		return Token{}, fmt.Errorf("unterminated string at line %d, column %d", pos.Line, pos.Column)
	}

	value := l.input[start:l.position]
	l.advance() // Skip closing quote

	// Basic escape handling
	value = strings.ReplaceAll(value, "\\\"", "\"")
	value = strings.ReplaceAll(value, "\\n", "\n")
	value = strings.ReplaceAll(value, "\\t", "\t")
	value = strings.ReplaceAll(value, "\\\\", "\\")

	return Token{Type: TokenString, Value: value, Position: pos}, nil
}

func (l *Lexer) readKeyword() (Token, error) {
	pos := l.currentPosition()
	l.advance() // Skip ':'
	start := l.position

	for l.position < len(l.input) && isSymbolChar(l.current()) {
		l.advance()
	}

	value := l.input[start:l.position]
	return Token{Type: TokenKeyword, Value: value, Position: pos}, nil
}

func (l *Lexer) readNumber() (Token, error) {
	pos := l.currentPosition()
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
	return Token{Type: TokenNumber, Value: value, Position: pos}, nil
}

func (l *Lexer) readSymbol() (Token, error) {
	pos := l.currentPosition()
	start := l.position

	for l.position < len(l.input) && isSymbolChar(l.current()) {
		l.advance()
	}

	value := l.input[start:l.position]
	return Token{Type: TokenSymbol, Value: value, Position: pos}, nil
}

func isSymbolStart(char rune) bool {
	return unicode.IsLetter(char) || char == '_' || char == '+' || char == '-' ||
		char == '*' || char == '/' || char == '=' || char == '<' || char == '>' ||
		char == '!' || char == '?' || char == '%' || char == '&'
}

func isSymbolChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' ||
		char == '-' || char == '+' || char == '*' || char == '/' || char == '=' ||
		char == '<' || char == '>' || char == '!' || char == '?' || char == '%' || char == '&'
}

// Parser converts tokens to AST
type Parser struct {
	tokens   []Token
	position int
	source   string // Original source code for error reporting
}

// ParseError represents a parsing error with location information
type ParseError struct {
	Message  string
	Position Position
	Source   string
}

func (e *ParseError) Error() string {
	if e.Source != "" {
		lines := strings.Split(e.Source, "\n")
		if e.Position.Line > 0 && e.Position.Line <= len(lines) {
			line := lines[e.Position.Line-1]
			// Show the line and point to the error
			pointer := strings.Repeat(" ", e.Position.Column-1) + "^"
			return fmt.Sprintf("Parse error at line %d, column %d: %s\n%s\n%s", 
				e.Position.Line, e.Position.Column, e.Message, line, pointer)
		}
	}
	return fmt.Sprintf("Parse error at line %d, column %d: %s", 
		e.Position.Line, e.Position.Column, e.Message)
}

// NewParser creates a new parser
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

// NewParserWithSource creates a new parser with source code for error reporting
func NewParserWithSource(tokens []Token, source string) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
		source:   source,
	}
}

// Parse converts tokens to Lisp values
func (p *Parser) Parse() (Value, error) {
	if p.position >= len(p.tokens) {
		pos := Position{Line: 1, Column: 1}
		if len(p.tokens) > 0 {
			pos = p.tokens[len(p.tokens)-1].Position
		}
		return nil, &ParseError{
			Message:  "unexpected end of input",
			Position: pos,
			Source:   p.source,
		}
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
	case TokenLeftBrace:
		return p.parseHashMap()
	case TokenHash:
		return p.parseSet()
	case TokenQuote:
		p.position++
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return NewList(Intern("quote"), expr), nil
	case TokenQuasiquote:
		p.position++
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return NewList(Intern("quasiquote"), expr), nil
	case TokenUnquote:
		p.position++
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return NewList(Intern("unquote"), expr), nil
	case TokenUnquoteSplicing:
		p.position++
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return NewList(Intern("unquote-splicing"), expr), nil
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
		return nil, &ParseError{
			Message:  "unexpected end of input",
			Position: token.Position,
			Source:   p.source,
		}
	default:
		return nil, &ParseError{
			Message:  fmt.Sprintf("unexpected token: %s", token.Value),
			Position: token.Position,
			Source:   p.source,
		}
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
		// Get position of the opening paren
		openParenPos := p.tokens[0].Position // This is approximate - we could track better
		if len(p.tokens) > 1 {
			openParenPos = p.tokens[len(p.tokens)-1].Position
		}
		return nil, &ParseError{
			Message:  "unterminated list",
			Position: openParenPos,
			Source:   p.source,
		}
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

func (p *Parser) parseHashMap() (Value, error) {
	p.position++ // Skip '{'

	var elements []Value

	for p.position < len(p.tokens) && p.tokens[p.position].Type != TokenRightBrace {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, expr)
	}

	if p.position >= len(p.tokens) {
		return nil, fmt.Errorf("unterminated hash-map")
	}

	if len(elements)%2 != 0 {
		return nil, fmt.Errorf("hash-map literal requires even number of elements")
	}

	p.position++ // Skip '}'
	return NewHashMapWithPairs(elements...), nil
}

func (p *Parser) parseSet() (Value, error) {
	p.position++ // Skip '#'

	if p.position >= len(p.tokens) || p.tokens[p.position].Type != TokenLeftBrace {
		return nil, fmt.Errorf("expected '{' after '#'")
	}

	p.position++ // Skip '{'

	var elements []Value

	for p.position < len(p.tokens) && p.tokens[p.position].Type != TokenRightBrace {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, expr)
	}

	if p.position >= len(p.tokens) {
		return nil, fmt.Errorf("unterminated set")
	}

	p.position++ // Skip '}'
	return NewSetWithElements(elements...), nil
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

	parser := NewParserWithSource(tokens, input)
	return parser.Parse()
}
