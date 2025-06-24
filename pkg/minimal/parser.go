package minimal

import (
	"fmt"
	"strconv"
	"strings"
)

// Token represents a token with position information
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
	TokenLeftParen
	TokenRightParen
	TokenLeftBracket
	TokenRightBracket
	TokenQuasiquote
	TokenUnquote
	TokenEOF
)

// Lexer tokenizes input with position tracking
type Lexer struct {
	input    string
	position int
	line     int
	column   int
	filename string
}

// NewLexer creates a new lexer
func NewLexer(input, filename string) *Lexer {
	return &Lexer{
		input:    input,
		position: 0,
		line:     1,
		column:   1,
		filename: filename,
	}
}

// Tokenize converts input into tokens with position information
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for l.position < len(l.input) {
		// Skip whitespace
		if l.isWhitespace(l.current()) {
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
			return nil, &ParseError{
				Message:        err.Error(),
				SourceLocation: Position{Line: l.line, Column: l.column, File: l.filename},
				Input:          l.input,
			}
		}

		tokens = append(tokens, token)
	}

	// Add EOF token
	tokens = append(tokens, Token{
		Type:     TokenEOF,
		Position: Position{Line: l.line, Column: l.column, File: l.filename},
	})

	return tokens, nil
}

func (l *Lexer) current() rune {
	if l.position >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position])
}

func (l *Lexer) peek() rune {
	if l.position+1 >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position+1])
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

func (l *Lexer) isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func (l *Lexer) skipComment() {
	for l.position < len(l.input) && l.current() != '\n' {
		l.advance()
	}
}

func (l *Lexer) nextToken() (Token, error) {
	startLine := l.line
	startColumn := l.column

	char := l.current()

	switch char {
	case '(':
		l.advance()
		return Token{
			Type:     TokenLeftParen,
			Value:    "(",
			Position: Position{Line: startLine, Column: startColumn, File: l.filename},
		}, nil

	case ')':
		l.advance()
		return Token{
			Type:     TokenRightParen,
			Value:    ")",
			Position: Position{Line: startLine, Column: startColumn, File: l.filename},
		}, nil

	case '[':
		l.advance()
		return Token{
			Type:     TokenLeftBracket,
			Value:    "[",
			Position: Position{Line: startLine, Column: startColumn, File: l.filename},
		}, nil

	case ']':
		l.advance()
		return Token{
			Type:     TokenRightBracket,
			Value:    "]",
			Position: Position{Line: startLine, Column: startColumn, File: l.filename},
		}, nil

	case '`':
		l.advance()
		return Token{
			Type:     TokenQuasiquote,
			Value:    "`",
			Position: Position{Line: startLine, Column: startColumn, File: l.filename},
		}, nil

	case '~':
		l.advance()
		return Token{
			Type:     TokenUnquote,
			Value:    "~",
			Position: Position{Line: startLine, Column: startColumn, File: l.filename},
		}, nil

	case '"':
		return l.readString(startLine, startColumn)

	default:
		if l.isDigit(char) || (char == '-' && l.isDigit(l.peek())) {
			return l.readNumber(startLine, startColumn)
		}
		return l.readSymbol(startLine, startColumn)
	}
}

func (l *Lexer) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l *Lexer) readString(startLine, startColumn int) (Token, error) {
	var value strings.Builder

	// Skip opening quote
	l.advance()

	for l.position < len(l.input) {
		char := l.current()

		if char == '"' {
			l.advance()
			return Token{
				Type:     TokenString,
				Value:    value.String(),
				Position: Position{Line: startLine, Column: startColumn, File: l.filename},
			}, nil
		}

		if char == '\\' {
			l.advance()
			if l.position >= len(l.input) {
				return Token{}, fmt.Errorf("unterminated string escape")
			}

			escaped := l.current()
			switch escaped {
			case 'n':
				value.WriteRune('\n')
			case 't':
				value.WriteRune('\t')
			case 'r':
				value.WriteRune('\r')
			case '\\':
				value.WriteRune('\\')
			case '"':
				value.WriteRune('"')
			default:
				value.WriteRune(escaped)
			}
			l.advance()
		} else {
			value.WriteRune(char)
			l.advance()
		}
	}

	return Token{}, fmt.Errorf("unterminated string")
}

func (l *Lexer) readNumber(startLine, startColumn int) (Token, error) {
	var value strings.Builder

	// Handle negative sign
	if l.current() == '-' {
		value.WriteRune('-')
		l.advance()
	}

	for l.position < len(l.input) {
		char := l.current()
		if l.isDigit(char) || char == '.' {
			value.WriteRune(char)
			l.advance()
		} else {
			break
		}
	}

	return Token{
		Type:     TokenNumber,
		Value:    value.String(),
		Position: Position{Line: startLine, Column: startColumn, File: l.filename},
	}, nil
}

func (l *Lexer) readSymbol(startLine, startColumn int) (Token, error) {
	var value strings.Builder

	for l.position < len(l.input) {
		char := l.current()
		if l.isWhitespace(char) || char == '(' || char == ')' ||
			char == '[' || char == ']' || char == '`' || char == '~' ||
			char == '"' || char == ';' {
			break
		}
		value.WriteRune(char)
		l.advance()
	}

	if value.Len() == 0 {
		return Token{}, fmt.Errorf("unexpected character: %c", l.current())
	}

	return Token{
		Type:     TokenSymbol,
		Value:    value.String(),
		Position: Position{Line: startLine, Column: startColumn, File: l.filename},
	}, nil
}

// Parser parses tokens into expressions with position tracking
type Parser struct {
	tokens   []Token
	position int
	filename string
}

// NewParser creates a new parser
func NewParser(tokens []Token, filename string) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
		filename: filename,
	}
}

// Parse parses tokens into a Value
func (p *Parser) Parse() (Value, *Position, error) {
	if p.position >= len(p.tokens) {
		return nil, nil, &ParseError{
			Message:        "unexpected end of input",
			SourceLocation: Position{File: p.filename},
		}
	}

	return p.parseExpression()
}

func (p *Parser) parseExpression() (Value, *Position, error) {
	if p.position >= len(p.tokens) {
		return nil, nil, &ParseError{
			Message:        "unexpected end of input",
			SourceLocation: Position{File: p.filename},
		}
	}

	token := p.tokens[p.position]
	pos := &token.Position

	switch token.Type {
	case TokenLeftParen:
		return p.parseList()

	case TokenLeftBracket:
		return p.parseVector()

	case TokenQuasiquote:
		p.position++
		expr, _, err := p.parseExpression()
		if err != nil {
			return nil, pos, err
		}
		return NewList(Intern("quasiquote"), expr), pos, nil

	case TokenUnquote:
		p.position++
		expr, _, err := p.parseExpression()
		if err != nil {
			return nil, pos, err
		}
		return NewList(Intern("unquote"), expr), pos, nil

	case TokenSymbol:
		p.position++
		return p.parseAtom(token.Value), pos, nil

	case TokenNumber:
		p.position++
		return p.parseAtom(token.Value), pos, nil

	case TokenString:
		p.position++
		return String(token.Value), pos, nil

	default:
		return nil, pos, &ParseError{
			Message:        fmt.Sprintf("unexpected token: %s", token.Value),
			SourceLocation: token.Position,
		}
	}
}

func (p *Parser) parseList() (Value, *Position, error) {
	startPos := &p.tokens[p.position].Position
	p.position++ // Skip '('

	var elements []Value

	for p.position < len(p.tokens) && p.tokens[p.position].Type != TokenRightParen {
		expr, _, err := p.parseExpression()
		if err != nil {
			return nil, startPos, err
		}
		elements = append(elements, expr)
	}

	if p.position >= len(p.tokens) {
		return nil, startPos, &ParseError{
			Message:        "unclosed list",
			SourceLocation: *startPos,
		}
	}

	p.position++ // Skip ')'
	return NewList(elements...), startPos, nil
}

func (p *Parser) parseVector() (Value, *Position, error) {
	startPos := &p.tokens[p.position].Position
	p.position++ // Skip '['

	var elements []Value

	for p.position < len(p.tokens) && p.tokens[p.position].Type != TokenRightBracket {
		expr, _, err := p.parseExpression()
		if err != nil {
			return nil, startPos, err
		}
		elements = append(elements, expr)
	}

	if p.position >= len(p.tokens) {
		return nil, startPos, &ParseError{
			Message:        "unclosed vector",
			SourceLocation: *startPos,
		}
	}

	p.position++ // Skip ']'
	return NewVector(elements...), startPos, nil
}

func (p *Parser) parseAtom(value string) Value {
	// Try to parse as number
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return Number(num)
	}

	// Parse as boolean
	switch value {
	case "true":
		return Boolean(true)
	case "false":
		return Boolean(false)
	case "nil":
		return Nil{}
	}

	// Parse as symbol
	return Intern(value)
}

// ParseWithPositions parses input and returns the expression with position information
func ParseWithPositions(input, filename string) (Value, *Position, error) {
	lexer := NewLexer(input, filename)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, nil, err
	}

	parser := NewParser(tokens, filename)
	return parser.Parse()
}
