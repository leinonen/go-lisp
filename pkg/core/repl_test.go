package core

import (
	"testing"
)

// Test the isBalanced function
func TestIsBalanced(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple balanced", "(+ 1 2)", true},
		{"simple unbalanced", "(+ 1 2", false},
		{"extra closing", "(+ 1 2))", false},
		{"nested balanced", "(map + (list 1 2 3) (list 4 5 6))", true},
		{"nested unbalanced", "(map + (list 1 2 3) (list 4 5 6)", false},
		{"empty string", "", true},
		{"only whitespace", "   \n\t  ", true},
		{"mixed brackets", "(vector [1 2 3] {key value})", true},
		{"mixed brackets unbalanced", "(vector [1 2 3] {key value)", false},
		{"string with parens", "(println \"hello (world)\")", true},
		{"string with escaped quote", "(println \"say \\\"hello\\\"\")", true},
		{"comment only", "; this is a comment", true},
		{"expression with comment", "(+ 1 2) ; comment", true},
		{"multiline with comment", "(defn test []\n  ; comment\n  (+ 1 2))", true},
		{"multiline unbalanced", "(defn test []\n  (+ 1 2", false},
		{"comment interrupts balance", "(+ 1 ; comment\n   2)", true},
		{"semicolon in string", "(println \"hello; world\")", true},
		{"backslash in comment", "; this is a \\ comment", true},
		{"complex nested", "(((()))))", false}, // 4 opens, 5 closes - unbalanced
		{"complex nested unbalanced", "(((())))", true}, // 4 opens, 4 closes - balanced
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBalanced(tt.input)
			if result != tt.expected {
				t.Errorf("isBalanced(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test the hasNonWhitespaceContent function
func TestHasNonWhitespaceContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple expression", "(+ 1 2)", true},
		{"empty string", "", false},
		{"only whitespace", "   \n\t  ", false},
		{"only comment", "; this is a comment", false},
		{"comment with whitespace", "  ; comment  \n  ", false},
		{"multiple comments", "; first\n; second\n  ; third  ", false},
		{"expression with comment", "(+ 1 2) ; comment", true},
		{"comment then expression", "; comment\n(+ 1 2)", true},
		{"string literal", "\"hello\"", true},
		{"string with parens", "\"hello (world)\"", true},
		{"escaped string", "\"say \\\"hello\\\"\"", true},
		{"multiline string", "\"line1\nline2\"", true},
		{"symbol only", "symbol", true},
		{"keyword only", ":keyword", true},
		{"number only", "42", true},
		{"nil literal", "nil", true},
		{"boolean literal", "true", true},
		{"whitespace around content", "  (+ 1 2)  ", true},
		{"comment around content", "; comment\n(+ 1 2)\n; another", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNonWhitespaceContent(tt.input)
			if result != tt.expected {
				t.Errorf("hasNonWhitespaceContent(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Helper function to count brackets for force evaluation testing
func countBrackets(input string) (openCount, closeCount int) {
	inString := false
	inComment := false
	escapeNext := false
	
	for _, char := range input {
		if escapeNext {
			escapeNext = false
			continue
		}
		
		if char == '\\' && inString {
			escapeNext = true
			continue
		}
		
		if char == ';' && !inString {
			inComment = true
			continue
		}
		
		if char == '\n' {
			inComment = false
			continue
		}
		
		if inComment {
			continue
		}
		
		if char == '"' {
			inString = !inString
			continue
		}
		
		if inString {
			continue
		}
		
		switch char {
		case '(', '[', '{':
			openCount++
		case ')', ']', '}':
			closeCount++
		}
	}
	
	return openCount, closeCount
}

// Test the bracket counting logic for force evaluation
func TestCountBrackets(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectedOpen int
		expectedClose int
	}{
		{"empty", "", 0, 0},
		{"balanced simple", "(+ 1 2)", 1, 1},
		{"unbalanced open", "(+ 1 2", 1, 0},
		{"unbalanced close", "+ 1 2)", 0, 1},
		{"nested balanced", "(map + (list 1))", 2, 2},
		{"nested unbalanced", "(map + (list 1)", 2, 1},
		{"mixed brackets", "(vector [1] {k v})", 3, 3},
		{"string with brackets", "(println \"(hello)\")", 1, 1},
		{"comment with brackets", "(+ 1 ; (comment)\n   2)", 1, 1},
		{"complex nesting", "((()))", 3, 3},
		{"escaped in string", "(println \"\\\"(hello)\\\"\")", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			open, close := countBrackets(tt.input)
			if open != tt.expectedOpen {
				t.Errorf("countBrackets(%q) open = %d, want %d", tt.input, open, tt.expectedOpen)
			}
			if close != tt.expectedClose {
				t.Errorf("countBrackets(%q) close = %d, want %d", tt.input, close, tt.expectedClose)
			}
		})
	}
}

// Test evaluation condition logic
func TestShouldEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"complete expression", "(+ 1 2)", true},
		{"incomplete expression", "(+ 1 2", false},
		{"empty input", "", false},
		{"whitespace only", "   \n  ", false},
		{"comment only", "; just comment", false},
		{"expression with comment", "(+ 1 2) ; comment", true},
		{"multiline complete", "(defn test []\n  (+ 1 2))", true},
		{"multiline incomplete", "(defn test []\n  (+ 1 2", false},
		{"string literal", "\"hello\"", true},
		{"keyword", ":test", true},
		{"symbol", "test-symbol", true},
		{"number", "42", true},
		{"nil", "nil", true},
		{"boolean", "true", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNonWhitespaceContent(tt.input) && isBalanced(tt.input)
			if result != tt.expected {
				t.Errorf("shouldEvaluate(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test REPL creation and environment access
func TestREPLCreationAndEnvironment(t *testing.T) {
	repl, err := NewREPL()
	if err != nil {
		t.Fatalf("Failed to create REPL: %v", err)
	}
	defer repl.rl.Close() // Clean up readline instance

	// Test that environment is accessible
	env := repl.GetEnv()
	if env == nil {
		t.Fatal("REPL environment is nil")
	}

	// Test that environment has expected built-in symbols
	symbols := env.GetAllSymbols()
	if len(symbols) == 0 {
		t.Fatal("Environment has no symbols")
	}

	// Check for some expected built-in functions
	expectedSymbols := []string{"+", "-", "*", "/", "cons", "first", "rest", "map", "filter"}
	symbolsMap := make(map[string]bool)
	for _, sym := range symbols {
		symbolsMap[sym] = true
	}

	for _, expected := range expectedSymbols {
		if !symbolsMap[expected] {
			t.Errorf("Expected symbol %q not found in environment", expected)
		}
	}
}

// Test REPL evaluation functionality
func TestREPLEvaluation(t *testing.T) {
	repl, err := NewREPL()
	if err != nil {
		t.Fatalf("Failed to create REPL: %v", err)
	}
	defer repl.rl.Close()

	tests := []struct {
		name     string
		input    string
		wantErr  bool
	}{
		{"simple addition", "(+ 1 2)", false},
		{"function call", "(cons 1 (list 2 3))", false},
		{"syntax error", "(+ 1 2", true}, // Incomplete expression should error
		{"undefined symbol", "undefined-symbol", true},
		{"string literal", "\"hello world\"", false},
		{"nil literal", "nil", false},
		{"boolean literal", "true", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repl.Eval(tt.input)
			hasErr := err != nil
			if hasErr != tt.wantErr {
				t.Errorf("REPL.Eval(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}