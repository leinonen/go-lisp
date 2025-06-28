package core

import (
	"testing"
)

func TestLexerTokenization(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{"(", []TokenType{TokenLeftParen, TokenEOF}},
		{")", []TokenType{TokenRightParen, TokenEOF}},
		{"[", []TokenType{TokenLeftBracket, TokenEOF}},
		{"]", []TokenType{TokenRightBracket, TokenEOF}},
		{"'", []TokenType{TokenQuote, TokenEOF}},
		{"42", []TokenType{TokenNumber, TokenEOF}},
		{"-42", []TokenType{TokenNumber, TokenEOF}},
		{"3.14", []TokenType{TokenNumber, TokenEOF}},
		{"hello", []TokenType{TokenSymbol, TokenEOF}},
		{":keyword", []TokenType{TokenKeyword, TokenEOF}},
		{"\"string\"", []TokenType{TokenString, TokenEOF}},
		{"(+ 1 2)", []TokenType{TokenLeftParen, TokenSymbol, TokenNumber, TokenNumber, TokenRightParen, TokenEOF}},
		{"[1 2 3]", []TokenType{TokenLeftBracket, TokenNumber, TokenNumber, TokenNumber, TokenRightBracket, TokenEOF}},
	}

	for _, test := range tests {
		lexer := NewLexer(test.input)
		tokens, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}

		if len(tokens) != len(test.expected) {
			t.Errorf("Expected %d tokens for '%s', got %d", len(test.expected), test.input, len(tokens))
			continue
		}

		for i, token := range tokens {
			if token.Type != test.expected[i] {
				t.Errorf("Expected token type %v at position %d for '%s', got %v", test.expected[i], i, test.input, token.Type)
			}
		}
	}
}

func TestLexerTokenValues(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue string
		expectedType  TokenType
	}{
		{"42", "42", TokenNumber},
		{"-42", "-42", TokenNumber},
		{"3.14", "3.14", TokenNumber},
		{"hello", "hello", TokenSymbol},
		{":keyword", "keyword", TokenKeyword},
		{"\"hello world\"", "hello world", TokenString},
		{"+", "+", TokenSymbol},
		{"test-symbol", "test-symbol", TokenSymbol},
		{"*special*", "*special*", TokenSymbol},
	}

	for _, test := range tests {
		lexer := NewLexer(test.input)
		tokens, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}

		if len(tokens) < 1 {
			t.Errorf("Expected at least 1 token for '%s'", test.input)
			continue
		}

		token := tokens[0]
		if token.Type != test.expectedType {
			t.Errorf("Expected token type %v for '%s', got %v", test.expectedType, test.input, token.Type)
		}
		if token.Value != test.expectedValue {
			t.Errorf("Expected token value '%s' for '%s', got '%s'", test.expectedValue, test.input, token.Value)
		}
	}
}

func TestLexerStringEscapes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"hello\\nworld\"", "hello\nworld"},
		{"\"hello\\tworld\"", "hello\tworld"},
		{"\"hello\\\"world\"", "hello\"world"},
		{"\"hello\\\\world\"", "hello\\world"},
	}

	for _, test := range tests {
		lexer := NewLexer(test.input)
		tokens, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}

		if len(tokens) < 1 || tokens[0].Type != TokenString {
			t.Errorf("Expected string token for '%s'", test.input)
			continue
		}

		if tokens[0].Value != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, tokens[0].Value)
		}
	}
}

func TestLexerComments(t *testing.T) {
	input := "; This is a comment\n(+ 1 2) ; Another comment\n"
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := []TokenType{TokenLeftParen, TokenSymbol, TokenNumber, TokenNumber, TokenRightParen, TokenEOF}
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, token := range tokens {
		if token.Type != expected[i] {
			t.Errorf("Expected token type %v at position %d, got %v", expected[i], i, token.Type)
		}
	}
}

func TestParserBasicExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"42", "42"},
		{"-42", "-42"},
		{"3.14", "3.14"},
		{"hello", "hello"},
		{":keyword", ":keyword"},
		{"\"string\"", "\"string\""},
		{"()", "()"},
		{"[]", "[]"},
		{"(+ 1 2)", "(+ 1 2)"},
		{"[1 2 3]", "[1 2 3]"},
		{"'x", "(quote x)"},
		{"'(+ 1 2)", "(quote (+ 1 2))"},
	}

	for _, test := range tests {
		result, err := ReadString(test.input)
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestParserNestedExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"((+ 1 2))", "((+ 1 2))"},
		{"(+ (* 2 3) 4)", "(+ (* 2 3) 4)"},
		{"[1 [2 3] 4]", "[1 [2 3] 4]"},
		{"(fn [x] (+ x 1))", "(fn [x] (+ x 1))"},
		{"(if true 1 2)", "(if true 1 2)"},
	}

	for _, test := range tests {
		result, err := ReadString(test.input)
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestParserParseAll(t *testing.T) {
	input := "(def x 42) (+ x 1) (* x 2)"

	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	parser := NewParser(tokens)
	expressions, err := parser.ParseAll()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(expressions) != 3 {
		t.Errorf("Expected 3 expressions, got %d", len(expressions))
	}

	expected := []string{"(def x 42)", "(+ x 1)", "(* x 2)"}
	for i, expr := range expressions {
		if expr.String() != expected[i] {
			t.Errorf("Expected '%s' at position %d, got '%s'", expected[i], i, expr.String())
		}
	}
}

func TestParserErrors(t *testing.T) {
	tests := []string{
		"(",              // Unterminated list
		"[",              // Unterminated vector
		")",              // Unexpected closing paren
		"]",              // Unexpected closing bracket
		"\"unterminated", // Unterminated string
		"'",              // Quote without expression
	}

	for _, test := range tests {
		_, err := ReadString(test)
		if err == nil {
			t.Errorf("Expected error for input '%s', but got none", test)
		}
	}
}

func TestNumberParsing(t *testing.T) {
	tests := []struct {
		input      string
		isInteger  bool
		intValue   int64
		floatValue float64
	}{
		{"42", true, 42, 42.0},
		{"-42", true, -42, -42.0},
		{"0", true, 0, 0.0},
		{"3.14", false, 3, 3.14},
		{"-3.14", false, -3, -3.14},
		{"0.0", false, 0, 0.0},
	}

	for _, test := range tests {
		result, err := ReadString(test.input)
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}

		num, ok := result.(Number)
		if !ok {
			t.Errorf("Expected Number for input '%s', got %T", test.input, result)
			continue
		}

		if num.IsInteger() != test.isInteger {
			t.Errorf("Expected isInteger=%v for input '%s', got %v", test.isInteger, test.input, num.IsInteger())
		}

		if num.ToInt() != test.intValue {
			t.Errorf("Expected int value %d for input '%s', got %d", test.intValue, test.input, num.ToInt())
		}

		if num.ToFloat() != test.floatValue {
			t.Errorf("Expected float value %f for input '%s', got %f", test.floatValue, test.input, num.ToFloat())
		}
	}
}
