package minimal

import (
	"testing"
)

func TestREPLIsBalanced(t *testing.T) {
	repl := NewREPL()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty input",
			input:    "",
			expected: true,
		},
		{
			name:     "simple balanced parens",
			input:    "(+ 1 2)",
			expected: true,
		},
		{
			name:     "unbalanced opening paren",
			input:    "(+ 1 2",
			expected: false,
		},
		{
			name:     "unbalanced closing paren",
			input:    "+ 1 2)",
			expected: false,
		},
		{
			name:     "nested balanced parens",
			input:    "(+ (* 2 3) (- 5 1))",
			expected: true,
		},
		{
			name:     "nested unbalanced parens",
			input:    "(+ (* 2 3) (- 5 1)",
			expected: false,
		},
		{
			name:     "balanced brackets",
			input:    "[1 2 3]",
			expected: true,
		},
		{
			name:     "unbalanced brackets",
			input:    "[1 2 3",
			expected: false,
		},
		{
			name:     "mixed balanced parens and brackets",
			input:    "(vector [1 2] [3 4])",
			expected: true,
		},
		{
			name:     "mixed unbalanced parens and brackets",
			input:    "(vector [1 2] [3 4)",
			expected: false,
		},
		{
			name:     "string with parens inside",
			input:    `"This is a (test) string"`,
			expected: true,
		},
		{
			name:     "unbalanced parens with string containing parens",
			input:    `(println "This is a (test) string"`,
			expected: false,
		},
		{
			name:     "balanced with string containing parens",
			input:    `(println "This is a (test) string")`,
			expected: true,
		},
		{
			name:     "multiline balanced",
			input:    "(+ 1\n   2\n   3)",
			expected: true,
		},
		{
			name:     "multiline unbalanced",
			input:    "(+ 1\n   2\n   3",
			expected: false,
		},
		{
			name:     "quoted strings with quotes",
			input:    `"He said \"Hello\""`,
			expected: true,
		},
		{
			name:     "unclosed string",
			input:    `"unclosed string`,
			expected: false,
		},
		{
			name:     "wrong bracket type",
			input:    "(vector [1 2 3)",
			expected: false,
		},
		{
			name:     "wrong paren type",
			input:    "[vector (1 2 3]",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repl.isBalanced(tt.input)
			if result != tt.expected {
				t.Errorf("isBalanced(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestREPLParse(t *testing.T) {
	repl := NewREPL()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "simple expression",
			input:   "(+ 1 2)",
			wantErr: false,
		},
		{
			name:    "nested expression",
			input:   "(+ (* 2 3) (- 5 1))",
			wantErr: false,
		},
		{
			name:    "vector expression",
			input:   "[1 2 3]",
			wantErr: false,
		},
		{
			name:    "mixed expression",
			input:   "(vector [1 2] [3 4])",
			wantErr: false,
		},
		{
			name:    "multiline expression",
			input:   "(+ 1\n   2\n   3)",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repl.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
