package kernel

import (
	"strings"
	"testing"
)

// TestParser provides comprehensive testing of the parser
func TestParser(t *testing.T) {
	t.Run("BasicTokens", func(t *testing.T) {
		input := `42 3.14 "hello" symbol true false nil`
		tokens, err := NewLexer(input, "<test>").Tokenize()
		if err != nil {
			t.Fatalf("Tokenization error: %v", err)
		}

		expectedTypes := []TokenType{
			TokenNumber, TokenNumber, TokenString, TokenSymbol,
			TokenSymbol, TokenSymbol, TokenSymbol, TokenEOF,
		}

		if len(tokens) != len(expectedTypes) {
			t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
		}

		for i, expected := range expectedTypes {
			if tokens[i].Type != expected {
				t.Errorf("Token %d: expected type %v, got %v", i, expected, tokens[i].Type)
			}
		}

		// Test specific token values
		if tokens[0].Value != "42" {
			t.Errorf("Expected first token value '42', got '%s'", tokens[0].Value)
		}
		if tokens[1].Value != "3.14" {
			t.Errorf("Expected second token value '3.14', got '%s'", tokens[1].Value)
		}
		if tokens[2].Value != "hello" {
			t.Errorf("Expected third token value 'hello', got '%s'", tokens[2].Value)
		}
	})

	t.Run("ParenthesesAndBrackets", func(t *testing.T) {
		input := `() [] (a b) [x y]`
		tokens, err := NewLexer(input, "<test>").Tokenize()
		if err != nil {
			t.Fatalf("Tokenization error: %v", err)
		}

		expectedTypes := []TokenType{
			TokenLeftParen, TokenRightParen,
			TokenLeftBracket, TokenRightBracket,
			TokenLeftParen, TokenSymbol, TokenSymbol, TokenRightParen,
			TokenLeftBracket, TokenSymbol, TokenSymbol, TokenRightBracket,
			TokenEOF,
		}

		if len(tokens) != len(expectedTypes) {
			t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
		}

		for i, expected := range expectedTypes {
			if tokens[i].Type != expected {
				t.Errorf("Token %d: expected type %v, got %v", i, expected, tokens[i].Type)
			}
		}
	})

	t.Run("QuasiquoteAndUnquote", func(t *testing.T) {
		input := "`(a ~b)"
		tokens, err := NewLexer(input, "<test>").Tokenize()
		if err != nil {
			t.Fatalf("Tokenization error: %v", err)
		}

		expectedTypes := []TokenType{
			TokenQuasiquote, TokenLeftParen, TokenSymbol, TokenUnquote, TokenSymbol, TokenRightParen, TokenEOF,
		}

		if len(tokens) != len(expectedTypes) {
			t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
		}

		for i, expected := range expectedTypes {
			if tokens[i].Type != expected {
				t.Errorf("Token %d: expected type %v, got %v", i, expected, tokens[i].Type)
			}
		}
	})

	t.Run("Comments", func(t *testing.T) {
		input := `42 ; this is a comment
; another comment
"string"`
		tokens, err := NewLexer(input, "<test>").Tokenize()
		if err != nil {
			t.Fatalf("Tokenization error: %v", err)
		}

		// Comments should be skipped
		expectedTypes := []TokenType{TokenNumber, TokenString, TokenEOF}

		if len(tokens) != len(expectedTypes) {
			t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
		}

		if tokens[0].Value != "42" {
			t.Errorf("Expected first token '42', got '%s'", tokens[0].Value)
		}
		if tokens[1].Value != "string" {
			t.Errorf("Expected second token 'string', got '%s'", tokens[1].Value)
		}
	})

	t.Run("Whitespace", func(t *testing.T) {
		input := `  42   
		"hello"   
		symbol  `
		tokens, err := NewLexer(input, "<test>").Tokenize()
		if err != nil {
			t.Fatalf("Tokenization error: %v", err)
		}

		// Whitespace should be skipped
		expectedTypes := []TokenType{TokenNumber, TokenString, TokenSymbol, TokenEOF}

		if len(tokens) != len(expectedTypes) {
			t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
		}
	})

	t.Run("ParseNumbers", func(t *testing.T) {
		cases := []struct {
			input    string
			expected float64
		}{
			{"42", 42.0},
			{"0", 0.0},
			{"3.14159", 3.14159},
			{"-5", -5.0},
			{"-2.5", -2.5},
			{"0.1", 0.1},
		}

		for _, tc := range cases {
			parsed, _, err := ParseWithPositions(tc.input, "<test>")
			if err != nil {
				t.Fatalf("Parse error for '%s': %v", tc.input, err)
			}

			num, ok := parsed.(Number)
			if !ok {
				t.Fatalf("Expected Number for '%s', got %T", tc.input, parsed)
			}

			if float64(num) != tc.expected {
				t.Errorf("Expected %f for '%s', got %f", tc.expected, tc.input, float64(num))
			}
		}
	})

	t.Run("ParseStrings", func(t *testing.T) {
		cases := []struct {
			input    string
			expected string
		}{
			{`"hello"`, "hello"},
			{`""`, ""},
			{`"hello world"`, "hello world"},
			{`"with spaces  "`, "with spaces  "},
		}

		for _, tc := range cases {
			parsed, _, err := ParseWithPositions(tc.input, "<test>")
			if err != nil {
				t.Fatalf("Parse error for %s: %v", tc.input, err)
			}

			str, ok := parsed.(String)
			if !ok {
				t.Fatalf("Expected String for %s, got %T", tc.input, parsed)
			}

			if string(str) != tc.expected {
				t.Errorf("Expected '%s' for %s, got '%s'", tc.expected, tc.input, string(str))
			}
		}
	})

	t.Run("ParseSymbols", func(t *testing.T) {
		cases := []string{
			"symbol",
			"+",
			"-",
			"*",
			"/",
			"test-symbol",
			"name?",
			"set!",
			"->vector",
		}

		for _, tc := range cases {
			parsed, _, err := ParseWithPositions(tc, "<test>")
			if err != nil {
				t.Fatalf("Parse error for '%s': %v", tc, err)
			}

			sym, ok := parsed.(Symbol)
			if !ok {
				t.Fatalf("Expected Symbol for '%s', got %T", tc, parsed)
			}

			if string(sym) != tc {
				t.Errorf("Expected symbol '%s', got '%s'", tc, string(sym))
			}
		}
	})

	t.Run("ParseLists", func(t *testing.T) {
		// Empty list
		parsed, _, err := ParseWithPositions("()", "<test>")
		if err != nil {
			t.Fatalf("Parse error for empty list: %v", err)
		}

		list, ok := parsed.(*List)
		if !ok {
			t.Fatalf("Expected List, got %T", parsed)
		}

		if !list.IsEmpty() {
			t.Error("Expected empty list")
		}

		// Simple list
		parsed, _, err = ParseWithPositions("(1 2 3)", "<test>")
		if err != nil {
			t.Fatalf("Parse error for simple list: %v", err)
		}

		list, ok = parsed.(*List)
		if !ok {
			t.Fatalf("Expected List, got %T", parsed)
		}

		if list.Length() != 3 {
			t.Errorf("Expected list length 3, got %d", list.Length())
		}

		// Nested list
		parsed, _, err = ParseWithPositions("(+ (* 2 3) 4)", "<test>")
		if err != nil {
			t.Fatalf("Parse error for nested list: %v", err)
		}

		list, ok = parsed.(*List)
		if !ok {
			t.Fatalf("Expected List, got %T", parsed)
		}

		if list.Length() != 3 {
			t.Errorf("Expected outer list length 3, got %d", list.Length())
		}

		// Check nested structure
		second := list.Rest().First()
		nestedList, ok := second.(*List)
		if !ok {
			t.Fatalf("Expected nested List, got %T", second)
		}

		if nestedList.Length() != 3 {
			t.Errorf("Expected nested list length 3, got %d", nestedList.Length())
		}
	})

	t.Run("ParseVectors", func(t *testing.T) {
		// Empty vector
		parsed, _, err := ParseWithPositions("[]", "<test>")
		if err != nil {
			t.Fatalf("Parse error for empty vector: %v", err)
		}

		vector, ok := parsed.(*Vector)
		if !ok {
			t.Fatalf("Expected Vector, got %T", parsed)
		}

		if !vector.IsEmpty() {
			t.Error("Expected empty vector")
		}

		// Simple vector
		parsed, _, err = ParseWithPositions("[1 2 3]", "<test>")
		if err != nil {
			t.Fatalf("Parse error for simple vector: %v", err)
		}

		vector, ok = parsed.(*Vector)
		if !ok {
			t.Fatalf("Expected Vector, got %T", parsed)
		}

		if vector.Length() != 3 {
			t.Errorf("Expected vector length 3, got %d", vector.Length())
		}

		// Nested vector with list
		parsed, _, err = ParseWithPositions("[1 (2 3) 4]", "<test>")
		if err != nil {
			t.Fatalf("Parse error for nested vector: %v", err)
		}

		vector, ok = parsed.(*Vector)
		if !ok {
			t.Fatalf("Expected Vector, got %T", parsed)
		}

		if vector.Length() != 3 {
			t.Errorf("Expected vector length 3, got %d", vector.Length())
		}

		// Check nested list inside vector
		nested := vector.Get(1)
		nestedList, ok := nested.(*List)
		if !ok {
			t.Fatalf("Expected nested List, got %T", nested)
		}

		if nestedList.Length() != 2 {
			t.Errorf("Expected nested list length 2, got %d", nestedList.Length())
		}
	})

	t.Run("ParseSpecialValues", func(t *testing.T) {
		cases := []struct {
			input    string
			expected Value
		}{
			{"true", Boolean(true)},
			{"false", Boolean(false)},
			{"nil", Nil{}},
		}

		for _, tc := range cases {
			parsed, _, err := ParseWithPositions(tc.input, "<test>")
			if err != nil {
				t.Fatalf("Parse error for '%s': %v", tc.input, err)
			}

			if parsed != tc.expected {
				t.Errorf("Expected %v for '%s', got %v", tc.expected, tc.input, parsed)
			}
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		errorCases := []struct {
			input       string
			shouldError bool
			errorCheck  func(error) bool
		}{
			{
				input:       "(unclosed",
				shouldError: true,
				errorCheck: func(err error) bool {
					return strings.Contains(err.Error(), "unclosed") ||
						strings.Contains(err.Error(), "unexpected")
				},
			},
			{
				input:       "[unclosed",
				shouldError: true,
				errorCheck: func(err error) bool {
					return strings.Contains(err.Error(), "unclosed") ||
						strings.Contains(err.Error(), "unexpected")
				},
			},
			{
				input:       "extra)",
				shouldError: true,
				errorCheck: func(err error) bool {
					return strings.Contains(err.Error(), "unexpected")
				},
			},
			{
				input:       "extra]",
				shouldError: true,
				errorCheck: func(err error) bool {
					return strings.Contains(err.Error(), "unexpected")
				},
			},
			{
				input:       `"unclosed string`,
				shouldError: true,
				errorCheck: func(err error) bool {
					return strings.Contains(err.Error(), "string") ||
						strings.Contains(err.Error(), "unclosed")
				},
			},
		}

		for _, tc := range errorCases {
			_, _, err := ParseWithPositions(tc.input, "<test>")

			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected error for input '%s', but got none", tc.input)
				} else if !tc.errorCheck(err) {
					t.Errorf("Error check failed for input '%s': %v", tc.input, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
			}
		}
	})

	t.Run("PositionTracking", func(t *testing.T) {
		input := `(def test-fn
  (fn [x]
    (* x 2)))`

		_, pos, err := ParseWithPositions(input, "test.lisp")
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		if pos == nil {
			t.Fatal("Expected position information")
		}

		if pos.File != "test.lisp" {
			t.Errorf("Expected file 'test.lisp', got '%s'", pos.File)
		}

		if pos.Line != 1 {
			t.Errorf("Expected line 1, got %d", pos.Line)
		}

		if pos.Column != 1 {
			t.Errorf("Expected column 1, got %d", pos.Column)
		}
	})

	t.Run("ComplexNesting", func(t *testing.T) {
		input := `(defmacro when [condition body]
  (quasiquote (if (unquote condition) 
                 (unquote body) 
                 nil)))`

		parsed, _, err := ParseWithPositions(input, "<test>")
		if err != nil {
			t.Fatalf("Parse error for complex nesting: %v", err)
		}

		list, ok := parsed.(*List)
		if !ok {
			t.Fatalf("Expected List, got %T", parsed)
		}

		if list.Length() != 4 {
			t.Errorf("Expected top level length 4, got %d", list.Length())
		}

		// Check that it's defmacro
		first := list.First()
		if sym, ok := first.(Symbol); !ok || string(sym) != "defmacro" {
			t.Errorf("Expected first element to be 'defmacro', got %v", first)
		}
	})
}
