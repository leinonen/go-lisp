package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/interpreter"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Mock interpreter for testing REPL functionality
type mockInterpreter struct {
	responses []interpretResponse
	callIndex int
}

type interpretResponse struct {
	result types.Value
	err    error
}

func (m *mockInterpreter) Interpret(input string) (types.Value, error) {
	if m.callIndex >= len(m.responses) {
		return types.StringValue(fmt.Sprintf("echo: %s", input)), nil
	}
	response := m.responses[m.callIndex]
	m.callIndex++
	return response.result, response.err
}

func newMockInterpreter(responses ...interpretResponse) *mockInterpreter {
	return &mockInterpreter{
		responses: responses,
		callIndex: 0,
	}
}

// Test helper to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestContainsExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "whitespace only",
			input:    "   \n\t  ",
			expected: false,
		},
		{
			name:     "simple expression",
			input:    "(+ 1 2)",
			expected: true,
		},
		{
			name:     "symbol only",
			input:    "foo",
			expected: true,
		},
		{
			name:     "number only",
			input:    "42",
			expected: true,
		},
		{
			name:     "string only",
			input:    `"hello"`,
			expected: true,
		},
		{
			name:     "comment only",
			input:    "; this is a comment",
			expected: false,
		},
		{
			name:     "multiple comment lines",
			input:    "; comment 1\n; comment 2\n; comment 3",
			expected: false,
		},
		{
			name:     "expression with comment",
			input:    "(+ 1 2) ; add numbers",
			expected: true,
		},
		{
			name:     "comment before expression",
			input:    "; comment\n(+ 1 2)",
			expected: true,
		},
		{
			name:     "expression in string with semicolon",
			input:    `"hello; world"`,
			expected: true,
		},
		{
			name:     "comment with semicolon in string",
			input:    `; "hello; world"`,
			expected: false,
		},
		{
			name:     "multiline with mixed comments and expressions",
			input:    "; comment\n(+ 1 2)\n; another comment",
			expected: true,
		},
		{
			name:     "escaped quote in string",
			input:    `"hello \"world\""`,
			expected: true,
		},
		{
			name:     "semicolon in escaped string should not be comment",
			input:    `"hello; \"escaped; quote\""`,
			expected: true,
		},
		{
			name:     "whitespace with comment",
			input:    "   ; just a comment  \n  ",
			expected: false,
		},
		{
			name:     "complex multiline expression",
			input:    "(defn factorial [n]\n  (if (= n 0)\n      1\n      (* n (factorial (- n 1)))))",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsExpression(tt.input)
			if result != tt.expected {
				t.Errorf("containsExpression(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReadCompleteExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple expression",
			input:    "(+ 1 2)\n",
			expected: "(+ 1 2)",
		},
		{
			name:     "quit command",
			input:    "quit\n",
			expected: "quit",
		},
		{
			name:     "exit command",
			input:    "exit\n",
			expected: "exit",
		},
		{
			name:     "quit with whitespace",
			input:    "  quit  \n",
			expected: "quit",
		},
		{
			name:     "multiline expression",
			input:    "(+\n  1\n  2)\n",
			expected: "(+\n  1\n  2)",
		},
		{
			name:     "nested parentheses",
			input:    "(+ (* 2 3) (/ 8 4))\n",
			expected: "(+ (* 2 3) (/ 8 4))",
		},
		{
			name:     "string with parentheses",
			input:    `"hello (world)"` + "\n",
			expected: `"hello (world)"`,
		},
		{
			name:     "comment at end of line",
			input:    "(+ 1 2) ; add numbers\n",
			expected: "(+ 1 2) ; add numbers",
		},
		{
			name:     "multiline with comments",
			input:    "; calculate sum\n(+ 1 2)\n",
			expected: "; calculate sum\n(+ 1 2)",
		},
		{
			name:     "unbalanced opening parentheses",
			input:    "((+ 1 2\n  3)\n",
			expected: "((+ 1 2\n  3)",
		},
		{
			name:     "expression with escaped quote",
			input:    `"hello \"world\""` + "\n",
			expected: `"hello \"world\""`,
		},
		{
			name:     "multiple complete expressions",
			input:    "(+ 1 2)\n(* 3 4)\n",
			expected: "(+ 1 2)",
		},
		{
			name:     "empty lines before expression",
			input:    "\n\n(+ 1 2)\n",
			expected: "\n\n(+ 1 2)",
		},
		{
			name:     "comment-only lines",
			input:    "; comment 1\n; comment 2\n",
			expected: "; comment 1\n; comment 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			result := readCompleteExpression(scanner)
			if result != tt.expected {
				t.Errorf("readCompleteExpression() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestReadCompleteExpressionUnbalancedParens(t *testing.T) {
	// Test case for unbalanced closing parentheses
	input := "(+ 1 2))\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	result := readCompleteExpression(scanner)
	expected := "(+ 1 2))"
	if result != expected {
		t.Errorf("readCompleteExpression with unbalanced closing parens = %q, expected %q", result, expected)
	}
}

func TestREPLFlow(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		responses      []interpretResponse
		expectedOutput []string // Substrings that should appear in output
	}{
		{
			name:  "simple expression evaluation",
			input: "(+ 1 2)\nquit\n",
			responses: []interpretResponse{
				{result: types.NumberValue(3), err: nil},
			},
			expectedOutput: []string{
				"Welcome to Go Lisp!",
				"lisp> ",
				"=> 3",
				"lisp> ",
				"Goodbye!",
			},
		},
		{
			name:  "error handling",
			input: "(undefined-function)\nquit\n",
			responses: []interpretResponse{
				{result: nil, err: fmt.Errorf("undefined function: undefined-function")},
			},
			expectedOutput: []string{
				"Welcome to Go Lisp!",
				"lisp> ",
				"Error: undefined function: undefined-function",
				"lisp> ",
				"Goodbye!",
			},
		},
		{
			name:  "multiline expression",
			input: "(+\n  1\n  2)\nquit\n",
			responses: []interpretResponse{
				{result: types.NumberValue(3), err: nil},
			},
			expectedOutput: []string{
				"Welcome to Go Lisp!",
				"lisp> ",
				"...   ",
				"...   ",
				"=> 3",
				"lisp> ",
				"Goodbye!",
			},
		},
		{
			name:  "empty input handling",
			input: "\n\n  \n(+ 1 1)\nquit\n",
			responses: []interpretResponse{
				{result: types.NumberValue(2), err: nil},
			},
			expectedOutput: []string{
				"Welcome to Go Lisp!",
				"lisp> ",
				"lisp> ",
				"lisp> ",
				"lisp> ",
				"=> 2",
				"lisp> ",
				"Goodbye!",
			},
		},
		{
			name:      "exit command",
			input:     "exit\n",
			responses: []interpretResponse{},
			expectedOutput: []string{
				"Welcome to Go Lisp!",
				"lisp> ",
				"Goodbye!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInterp := newMockInterpreter(tt.responses...)
			scanner := bufio.NewScanner(strings.NewReader(tt.input))

			output := captureOutput(func() {
				REPLWithOptions(mockInterp, scanner, false) // Disable colors for tests
			})

			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, but got:\n%s", expected, output)
				}
			}
		})
	}
}

func TestPrintWelcomeMessage(t *testing.T) {
	output := captureOutput(func() {
		printWelcomeMessageNoColor()
	})

	expectedParts := []string{
		"Welcome to Go Lisp!",
		"Type expressions to evaluate them",
		"Multi-line expressions are supported",
		"(help)",
		"(env)",
		"(modules)",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("Expected welcome message to contain %q, but got:\n%s", part, output)
		}
	}
}

func TestPrintGoodbyeMessage(t *testing.T) {
	output := captureOutput(func() {
		printGoodbyeMessageNoColor()
	})

	expected := "Goodbye!"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected goodbye message to contain %q, but got:\n%s", expected, output)
	}
}

// Integration test with real interpreter
func TestREPLIntegration(t *testing.T) {
	// Test with actual interpreter for basic functionality
	interp, err := interpreter.NewInterpreter()
	if err != nil {
		t.Fatalf("Failed to create interpreter: %v", err)
	}
	input := "(+ 1 2)\nquit\n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	output := captureOutput(func() {
		REPLWithOptions(interp, scanner, false) // Disable colors for tests
	})

	expectedParts := []string{
		"Welcome to Go Lisp!",
		"lisp> ",
		"=> 3",
		"Goodbye!",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("Expected integration test output to contain %q, but got:\n%s", part, output)
		}
	}
}

// Test edge cases for string and comment parsing
func TestStringAndCommentParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "string with embedded semicolon and quotes",
			input:    `"test; \"quoted; text\""`,
			expected: true,
		},
		{
			name:     "nested escaped quotes",
			input:    `"outer \"inner \\\"nested\\\" inner\" outer"`,
			expected: true,
		},
		{
			name:     "backslash at end of string",
			input:    `"test\\"`,
			expected: true,
		},
		{
			name:     "comment after string with semicolon",
			input:    `"test;" ; this is a comment`,
			expected: true,
		},
		{
			name:     "multiple strings and comments",
			input:    `"first;" "second;" ; comment`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsExpression(tt.input)
			if result != tt.expected {
				t.Errorf("containsExpression(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkContainsExpression(b *testing.B) {
	input := `; This is a comment
(defn factorial [n]
  (if (= n 0)
      1
      (* n (factorial (- n 1)))))
; Another comment`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		containsExpression(input)
	}
}

func BenchmarkReadCompleteExpression(b *testing.B) {
	input := "(defn factorial [n]\n  (if (= n 0)\n      1\n      (* n (factorial (- n 1)))))\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner := bufio.NewScanner(strings.NewReader(input))
		readCompleteExpression(scanner)
	}
}
