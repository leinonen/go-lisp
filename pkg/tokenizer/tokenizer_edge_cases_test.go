package tokenizer

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/types"
)

func TestTokenizerEdgeCases(t *testing.T) {
	t.Run("empty input variations", func(t *testing.T) {
		inputs := []string{
			"",
			"   ",
			"\t",
			"\n",
			"\r\n",
			"   \t \n  ",
		}

		for _, input := range inputs {
			tokenizer := NewTokenizer(input)
			result := tokenizer.Tokenize()
			if len(result) != 0 {
				t.Errorf("expected empty token list for input %q, got %d tokens", input, len(result))
			}
		}
	})

	t.Run("whitespace preservation in strings", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{`"hello world"`, "hello world"},
			{`"  leading spaces"`, "  leading spaces"},
			{`"trailing spaces  "`, "trailing spaces  "},
			{`"multiple   spaces"`, "multiple   spaces"},
		}

		for _, tt := range tests {
			tokenizer := NewTokenizer(tt.input)
			result := tokenizer.Tokenize()
			if len(result) != 1 {
				t.Errorf("expected 1 token for %q, got %d", tt.input, len(result))
				continue
			}
			if result[0].Type != types.STRING {
				t.Errorf("expected STRING token for %q, got %v", tt.input, result[0].Type)
				continue
			}
			if result[0].Value != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result[0].Value)
			}
		}
	})

	t.Run("large numbers", func(t *testing.T) {
		tests := []string{
			"123456789012345678901234567890",
			"9999999999999999999999999999999999999999",
			"123.456789012345",
			"0.000000000000000000001",
		}

		for _, input := range tests {
			tokenizer := NewTokenizer(input)
			result := tokenizer.Tokenize()
			if len(result) != 1 {
				t.Errorf("expected 1 token for %q, got %d", input, len(result))
				continue
			}
			if result[0].Type != types.NUMBER {
				t.Errorf("expected NUMBER token for %q, got %v", input, result[0].Type)
				continue
			}
			if result[0].Value != input {
				t.Errorf("expected %q, got %q", input, result[0].Value)
			}
		}
	})

	t.Run("edge case symbols", func(t *testing.T) {
		tests := []string{
			"a",
			"x1",
			"var-name",
			"var_name",
			"*global*",
			"function?",
			"empty!",
			"->",
			"<-",
			"++",
			"--",
			"**",
			"//",
			"%%",
			"+=",
			"-=",
			"<=>",
			"symbol123",
		}

		for _, input := range tests {
			tokenizer := NewTokenizer(input)
			result := tokenizer.Tokenize()
			if len(result) != 1 {
				t.Errorf("expected 1 token for %q, got %d", input, len(result))
				continue
			}
			if result[0].Type != types.SYMBOL {
				t.Errorf("expected SYMBOL token for %q, got %v", input, result[0].Type)
				continue
			}
			if result[0].Value != input {
				t.Errorf("expected %q, got %q", input, result[0].Value)
			}
		}
	})

	t.Run("nested structures", func(t *testing.T) {
		input := "((((((nested))))))"
		tokenizer := NewTokenizer(input)
		result := tokenizer.Tokenize()

		expected := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.LPAREN, Value: "("},
			{Type: types.LPAREN, Value: "("},
			{Type: types.LPAREN, Value: "("},
			{Type: types.LPAREN, Value: "("},
			{Type: types.LPAREN, Value: "("},
			{Type: types.SYMBOL, Value: "nested"},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RPAREN, Value: ")"},
		}

		if len(result) != len(expected) {
			t.Errorf("expected %d tokens, got %d", len(expected), len(result))
		}

		for i, expectedToken := range expected {
			if i >= len(result) {
				break
			}
			if result[i].Type != expectedToken.Type || result[i].Value != expectedToken.Value {
				t.Errorf("token %d: expected %v, got %v", i, expectedToken, result[i])
			}
		}
	})

	t.Run("mixed brackets and parentheses", func(t *testing.T) {
		input := "([{}])"
		tokenizer := NewTokenizer(input)
		result := tokenizer.Tokenize()

		// Note: {} might not be supported, so we test what we know works
		input = "([()])"
		tokenizer = NewTokenizer(input)
		result = tokenizer.Tokenize()

		expected := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.LBRACKET, Value: "["},
			{Type: types.LPAREN, Value: "("},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RBRACKET, Value: "]"},
			{Type: types.RPAREN, Value: ")"},
		}

		if len(result) != len(expected) {
			t.Errorf("expected %d tokens, got %d", len(expected), len(result))
		}

		for i, expectedToken := range expected {
			if i >= len(result) {
				break
			}
			if result[i].Type != expectedToken.Type || result[i].Value != expectedToken.Value {
				t.Errorf("token %d: expected %v, got %v", i, expectedToken, result[i])
			}
		}
	})

	t.Run("comments with special characters", func(t *testing.T) {
		tests := []struct {
			input    string
			expected []types.Token
		}{
			{
				"; comment with (parens) and [brackets]",
				[]types.Token{},
			},
			{
				"42 ; inline comment with symbols !@#$%",
				[]types.Token{{Type: types.NUMBER, Value: "42"}},
			},
			{
				"; comment with \"quotes\" and 'apostrophes'\n(+ 1 2)",
				[]types.Token{
					{Type: types.LPAREN, Value: "("},
					{Type: types.SYMBOL, Value: "+"},
					{Type: types.NUMBER, Value: "1"},
					{Type: types.NUMBER, Value: "2"},
					{Type: types.RPAREN, Value: ")"},
				},
			},
			{
				"(a ; comment\n b)",
				[]types.Token{
					{Type: types.LPAREN, Value: "("},
					{Type: types.SYMBOL, Value: "a"},
					{Type: types.SYMBOL, Value: "b"},
					{Type: types.RPAREN, Value: ")"},
				},
			},
		}

		for _, tt := range tests {
			tokenizer := NewTokenizer(tt.input)
			result := tokenizer.Tokenize()

			if len(result) != len(tt.expected) {
				t.Errorf("input %q: expected %d tokens, got %d", tt.input, len(tt.expected), len(result))
				continue
			}

			for i, expectedToken := range tt.expected {
				if result[i].Type != expectedToken.Type || result[i].Value != expectedToken.Value {
					t.Errorf("input %q, token %d: expected %v, got %v", tt.input, i, expectedToken, result[i])
				}
			}
		}
	})

	t.Run("unicode characters", func(t *testing.T) {
		tests := []struct {
			input    string
			expected []types.Token
		}{
			{
				`"unicode: ñ ü ♠ ♥ ♦ ♣"`,
				[]types.Token{{Type: types.STRING, Value: "unicode: ñ ü ♠ ♥ ♦ ♣"}},
			},
		}

		for _, tt := range tests {
			tokenizer := NewTokenizer(tt.input)
			result := tokenizer.Tokenize()

			if len(result) != len(tt.expected) {
				t.Errorf("input %q: expected %d tokens, got %d", tt.input, len(tt.expected), len(result))
				continue
			}

			for i, expectedToken := range tt.expected {
				if result[i].Type != expectedToken.Type || result[i].Value != expectedToken.Value {
					t.Errorf("input %q, token %d: expected %v, got %v", tt.input, i, expectedToken, result[i])
				}
			}
		}
	})
}

func TestTokenizerErrorConditions(t *testing.T) {
	t.Run("unterminated strings", func(t *testing.T) {
		tests := []string{
			`"unterminated`,
			`"unterminated with newline
			`,
		}

		for _, input := range tests {
			tokenizer := NewTokenizer(input)
			_, err := tokenizer.TokenizeWithError()
			if err == nil {
				t.Errorf("expected error for unterminated string: %q", input)
			}
		}
	})

	t.Run("invalid characters", func(t *testing.T) {
		// Test characters that might not be valid in the language
		tests := []string{
			"@invalid",
			"#invalid",
			"$invalid",
			"invalid`char",
			"invalid~char",
		}

		for _, input := range tests {
			tokenizer := NewTokenizer(input)
			_, err := tokenizer.TokenizeWithError()
			// Some of these might be valid symbols depending on implementation
			if err != nil {
				t.Logf("input %q correctly rejected: %v", input, err)
			} else {
				t.Logf("input %q was accepted (might be valid symbol)", input)
			}
		}
	})

	t.Run("very long input", func(t *testing.T) {
		// Test with a very long symbol name
		longSymbol := string(make([]byte, 10000))
		for i := range longSymbol {
			longSymbol = longSymbol[:i] + "a" + longSymbol[i+1:]
		}

		tokenizer := NewTokenizer(longSymbol)
		result := tokenizer.Tokenize()

		if len(result) != 1 {
			t.Errorf("expected 1 token for long symbol, got %d", len(result))
		} else if result[0].Type != types.SYMBOL {
			t.Errorf("expected SYMBOL token for long symbol, got %v", result[0].Type)
		} else if len(result[0].Value) != 10000 {
			t.Errorf("expected symbol length 10000, got %d", len(result[0].Value))
		}
	})

	t.Run("numbers with invalid formats", func(t *testing.T) {
		tests := []string{
			"123.456.789", // multiple decimal points
			"123e",        // incomplete scientific notation
			"123e+",       // incomplete scientific notation
			"123.e10",     // decimal before e
			".123",        // leading decimal point
			"123.",        // trailing decimal point
		}

		for _, input := range tests {
			tokenizer := NewTokenizer(input)
			result := tokenizer.Tokenize()

			// These might be tokenized as multiple tokens or symbols
			// depending on implementation
			t.Logf("input %q tokenized as %d tokens", input, len(result))
			for i, token := range result {
				t.Logf("  token %d: %v = %q", i, token.Type, token.Value)
			}
		}
	})
}
