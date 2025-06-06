package parser

import (
	"fmt"
	"strconv"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

type Parser struct {
	tokens   []types.Token
	position int
	current  types.Token
}

func NewParser(tokens []types.Token) *Parser {
	p := &Parser{
		tokens: tokens,
	}
	p.readToken()
	return p
}

func (p *Parser) readToken() {
	if p.position >= len(p.tokens) {
		p.current = types.Token{Type: types.TokenType(-1), Value: ""} // EOF token
	} else {
		p.current = p.tokens[p.position]
	}
	p.position++
}

func (p *Parser) Parse() (types.Expr, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Check for remaining tokens (like unmatched closing parentheses)
	if p.current.Type != types.TokenType(-1) { // Not EOF
		return nil, fmt.Errorf("unexpected token after expression: %v", p.current)
	}

	return expr, nil
}

func (p *Parser) parseExpr() (types.Expr, error) {
	switch p.current.Type {
	case types.NUMBER:
		return p.parseNumber()
	case types.STRING:
		return p.parseString()
	case types.BOOLEAN:
		return p.parseBoolean()
	case types.SYMBOL:
		return p.parseSymbol()
	case types.LPAREN:
		return p.parseList()
	case types.RPAREN:
		return nil, fmt.Errorf("unexpected closing parenthesis")
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.current)
	}
}

func (p *Parser) parseNumber() (types.Expr, error) {
	value, err := strconv.ParseFloat(p.current.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", p.current.Value)
	}

	expr := &types.NumberExpr{Value: value}
	p.readToken()
	return expr, nil
}

func (p *Parser) parseString() (types.Expr, error) {
	expr := &types.StringExpr{Value: p.current.Value}
	p.readToken()
	return expr, nil
}

func (p *Parser) parseBoolean() (types.Expr, error) {
	value := p.current.Value == "#t"
	expr := &types.BooleanExpr{Value: value}
	p.readToken()
	return expr, nil
}

func (p *Parser) parseSymbol() (types.Expr, error) {
	expr := &types.SymbolExpr{Name: p.current.Value}
	p.readToken()
	return expr, nil
}

func (p *Parser) parseList() (types.Expr, error) {
	p.readToken() // consume '('

	elements := make([]types.Expr, 0)

	for p.current.Type != types.RPAREN {
		if p.current.Type == types.TokenType(-1) { // EOF
			return nil, fmt.Errorf("unmatched opening parenthesis")
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		elements = append(elements, expr)
	}

	p.readToken() // consume ')'

	return &types.ListExpr{Elements: elements}, nil
}
