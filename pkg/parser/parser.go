package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/leinonen/go-lisp/pkg/types"
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
	case types.KEYWORD:
		return p.parseKeyword()
	case types.LPAREN:
		return p.parseList()
	case types.LBRACKET:
		return p.parseBracket()
	case types.QUOTE:
		return p.parseQuote()
	case types.RPAREN:
		return nil, fmt.Errorf("line %d, column %d: unexpected closing parenthesis", p.current.Position.Line, p.current.Position.Column)
	case types.RBRACKET:
		return nil, fmt.Errorf("line %d, column %d: unexpected closing bracket", p.current.Position.Line, p.current.Position.Column)
	default:
		return nil, fmt.Errorf("line %d, column %d: unexpected token: %v", p.current.Position.Line, p.current.Position.Column, p.current)
	}
}

func (p *Parser) parseNumber() (types.Expr, error) {
	// First try to parse as a regular float64
	value, err := strconv.ParseFloat(p.current.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("line %d, column %d: invalid number: %s", p.current.Position.Line, p.current.Position.Column, p.current.Value)
	}

	// Check if this number might lose precision when converted to float64
	// If the string representation is a large integer, use BigNumberExpr
	originalStr := p.current.Value
	pos := p.current.Position

	// Check if it's an integer without decimal point
	if !strings.Contains(originalStr, ".") && !strings.Contains(originalStr, "e") && !strings.Contains(originalStr, "E") {
		// Check if it's too large for safe float64 integer representation
		if len(originalStr) > 15 || (len(originalStr) == 16 && originalStr[0] > '1') {
			// This is a large integer that might lose precision in float64
			expr := &types.BigNumberExpr{Value: originalStr, Position: pos}
			p.readToken()
			return expr, nil
		}
	}

	// Use regular NumberExpr for smaller numbers or floating point numbers
	expr := &types.NumberExpr{Value: value, Position: pos}
	p.readToken()
	return expr, nil
}

func (p *Parser) parseString() (types.Expr, error) {
	expr := &types.StringExpr{Value: p.current.Value, Position: p.current.Position}
	p.readToken()
	return expr, nil
}

func (p *Parser) parseBoolean() (types.Expr, error) {
	var value bool
	switch p.current.Value {
	case "true":
		value = true
	case "false":
		value = false
	default:
		return nil, fmt.Errorf("line %d, column %d: invalid boolean value: %s", p.current.Position.Line, p.current.Position.Column, p.current.Value)
	}
	expr := &types.BooleanExpr{Value: value, Position: p.current.Position}
	p.readToken()
	return expr, nil
}

func (p *Parser) parseSymbol() (types.Expr, error) {
	expr := &types.SymbolExpr{Name: p.current.Value, Position: p.current.Position}
	p.readToken()
	return expr, nil
}

func (p *Parser) parseKeyword() (types.Expr, error) {
	expr := &types.KeywordExpr{Value: p.current.Value, Position: p.current.Position}
	p.readToken()
	return expr, nil
}

func (p *Parser) parseList() (types.Expr, error) {
	listPos := p.current.Position // Capture position of opening parenthesis
	p.readToken()                 // consume '('

	elements := make([]types.Expr, 0)

	for p.current.Type != types.RPAREN {
		if p.current.Type == types.TokenType(-1) { // EOF
			return nil, fmt.Errorf("line %d, column %d: unmatched opening parenthesis", listPos.Line, listPos.Column)
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		elements = append(elements, expr)
	}

	p.readToken() // consume ')'

	// Check for special module forms
	if len(elements) > 0 {
		if symbolExpr, ok := elements[0].(*types.SymbolExpr); ok {
			switch symbolExpr.Name {
			case "module":
				return p.parseModuleFromElements(elements, listPos)
			case "import":
				return p.parseImportFromElements(elements, listPos)
			case "load":
				return p.parseLoadFromElements(elements, listPos)
			case "require":
				return p.parseRequireFromElements(elements, listPos)
			}
		}
	}

	return &types.ListExpr{Elements: elements, Position: listPos}, nil
}

func (p *Parser) parseModuleFromElements(elements []types.Expr, pos types.Position) (types.Expr, error) {
	// (module name (export sym1 sym2...) body...)
	if len(elements) < 4 {
		return nil, fmt.Errorf("line %d, column %d: module requires at least name, export list, and body", pos.Line, pos.Column)
	}

	// Get module name
	nameExpr, ok := elements[1].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("line %d, column %d: module name must be a symbol", pos.Line, pos.Column)
	}

	// Get export list
	exportListExpr, ok := elements[2].(*types.ListExpr)
	if !ok {
		return nil, fmt.Errorf("line %d, column %d: module export list must be a list", pos.Line, pos.Column)
	}

	// Check export list format: (export symbol1 symbol2...)
	if len(exportListExpr.Elements) < 1 {
		return nil, fmt.Errorf("line %d, column %d: export list cannot be empty", pos.Line, pos.Column)
	}

	exportKeyword, ok := exportListExpr.Elements[0].(*types.SymbolExpr)
	if !ok || exportKeyword.Name != "export" {
		return nil, fmt.Errorf("line %d, column %d: export list must start with 'export'", pos.Line, pos.Column)
	}

	// Parse exported symbols
	exports := make([]string, len(exportListExpr.Elements)-1)
	for i, expr := range exportListExpr.Elements[1:] {
		symExpr, ok := expr.(*types.SymbolExpr)
		if !ok {
			return nil, fmt.Errorf("line %d, column %d: exported names must be symbols", pos.Line, pos.Column)
		}
		exports[i] = symExpr.Name
	}

	// Get body expressions
	body := elements[3:]

	return &types.ModuleExpr{
		Name:     nameExpr.Name,
		Exports:  exports,
		Body:     body,
		Position: pos,
	}, nil
}

func (p *Parser) parseImportFromElements(elements []types.Expr, pos types.Position) (types.Expr, error) {
	// (import module-name)
	if len(elements) != 2 {
		return nil, fmt.Errorf("line %d, column %d: import requires exactly one module name", pos.Line, pos.Column)
	}

	moduleNameExpr, ok := elements[1].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("line %d, column %d: import module name must be a symbol", pos.Line, pos.Column)
	}

	return &types.ImportExpr{ModuleName: moduleNameExpr.Name, Position: pos}, nil
}

func (p *Parser) parseLoadFromElements(elements []types.Expr, pos types.Position) (types.Expr, error) {
	// (load "filename.lisp")
	if len(elements) != 2 {
		return nil, fmt.Errorf("line %d, column %d: load requires exactly one filename", pos.Line, pos.Column)
	}

	filenameExpr, ok := elements[1].(*types.StringExpr)
	if !ok {
		return nil, fmt.Errorf("line %d, column %d: load filename must be a string", pos.Line, pos.Column)
	}

	return &types.LoadExpr{Filename: filenameExpr.Value, Position: pos}, nil
}

func (p *Parser) parseRequireFromElements(elements []types.Expr, pos types.Position) (types.Expr, error) {
	// Supported syntaxes:
	// (require "filename.lisp")
	// (require "filename.lisp" :as alias)
	// (require "filename.lisp" :only [symbol1 symbol2])

	if len(elements) < 2 {
		return nil, fmt.Errorf("line %d, column %d: require requires at least a filename", pos.Line, pos.Column)
	}

	filenameExpr, ok := elements[1].(*types.StringExpr)
	if !ok {
		return nil, fmt.Errorf("line %d, column %d: require filename must be a string", pos.Line, pos.Column)
	}

	requireExpr := &types.RequireExpr{
		Filename: filenameExpr.Value,
		Position: pos,
	}

	// Parse optional modifiers
	if len(elements) > 2 {
		if len(elements) < 4 {
			return nil, fmt.Errorf("line %d, column %d: require modifier requires an argument", pos.Line, pos.Column)
		}

		keywordExpr, ok := elements[2].(*types.KeywordExpr)
		if !ok {
			return nil, fmt.Errorf("line %d, column %d: require modifier must be a keyword (:as or :only)", pos.Line, pos.Column)
		}

		switch keywordExpr.Value {
		case "as":
			// (require "file.lisp" :as alias)
			if len(elements) != 4 {
				return nil, fmt.Errorf("line %d, column %d: require :as expects exactly one alias symbol", pos.Line, pos.Column)
			}
			aliasExpr, ok := elements[3].(*types.SymbolExpr)
			if !ok {
				return nil, fmt.Errorf("line %d, column %d: require :as alias must be a symbol", pos.Line, pos.Column)
			}
			requireExpr.AsAlias = aliasExpr.Name

		case "only":
			// (require "file.lisp" :only [symbol1 symbol2])
			if len(elements) != 4 {
				return nil, fmt.Errorf("line %d, column %d: require :only expects exactly one symbol list", pos.Line, pos.Column)
			}
			listExpr, ok := elements[3].(*types.ListExpr)
			if !ok {
				return nil, fmt.Errorf("line %d, column %d: require :only expects a list of symbols", pos.Line, pos.Column)
			}

			onlyList := make([]string, len(listExpr.Elements))
			for i, elem := range listExpr.Elements {
				symbolExpr, ok := elem.(*types.SymbolExpr)
				if !ok {
					return nil, fmt.Errorf("line %d, column %d: require :only list must contain only symbols", pos.Line, pos.Column)
				}
				onlyList[i] = symbolExpr.Name
			}
			requireExpr.OnlyList = onlyList

		default:
			return nil, fmt.Errorf("line %d, column %d: require modifier must be :as or :only, got %s", pos.Line, pos.Column, keywordExpr.Value)
		}
	}

	return requireExpr, nil
}

func (p *Parser) parseQuote() (types.Expr, error) {
	quotePos := p.current.Position // Capture position of quote
	// Consume the quote token
	p.readToken()

	// Parse the expression that follows the quote
	expr, err := p.parseExpr()
	if err != nil {
		return nil, fmt.Errorf("line %d, column %d: error parsing quoted expression: %v", quotePos.Line, quotePos.Column, err)
	}

	// Convert 'expr to (quote expr)
	quoteSymbol := &types.SymbolExpr{Name: "quote", Position: quotePos}
	return &types.ListExpr{Elements: []types.Expr{quoteSymbol, expr}, Position: quotePos}, nil
}

func (p *Parser) parseBracket() (types.Expr, error) {
	bracketPos := p.current.Position // Capture position of opening bracket
	p.readToken()                    // consume '['

	elements := make([]types.Expr, 0)

	for p.current.Type != types.RBRACKET {
		if p.current.Type == types.TokenType(-1) { // EOF
			return nil, fmt.Errorf("line %d, column %d: unmatched opening bracket", bracketPos.Line, bracketPos.Column)
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		elements = append(elements, expr)
	}

	p.readToken() // consume ']'

	return &types.BracketExpr{Elements: elements, Position: bracketPos}, nil
}
