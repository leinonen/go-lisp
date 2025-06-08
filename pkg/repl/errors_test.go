package repl

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorFormatter_categorizeError(t *testing.T) {
	ef := NewErrorFormatter()

	tests := []struct {
		name         string
		errorMsg     string
		expectedType ErrorType
	}{
		// Syntax errors
		{
			name:         "unexpected token",
			errorMsg:     "unexpected token after expression",
			expectedType: ErrorTypeSyntax,
		},
		{
			name:         "parse error",
			errorMsg:     "parse error: unmatched closing parenthesis",
			expectedType: ErrorTypeSyntax,
		},
		{
			name:         "empty input",
			errorMsg:     "empty input",
			expectedType: ErrorTypeSyntax,
		},

		// Undefined symbol errors
		{
			name:         "undefined symbol",
			errorMsg:     "undefined symbol: foo",
			expectedType: ErrorTypeUndefined,
		},
		{
			name:         "unknown function",
			errorMsg:     "unknown function called",
			expectedType: ErrorTypeUndefined,
		},
		{
			name:         "no help available",
			errorMsg:     "no help available for 'unknown-function'",
			expectedType: ErrorTypeUndefined,
		},

		// Type errors
		{
			name:         "requires a boolean",
			errorMsg:     "if condition must be a boolean",
			expectedType: ErrorTypeTypeError,
		},
		{
			name:         "not a list",
			errorMsg:     "first requires a list argument",
			expectedType: ErrorTypeTypeError,
		},
		{
			name:         "not a function",
			errorMsg:     "value is not a function",
			expectedType: ErrorTypeTypeError,
		},

		// Runtime errors
		{
			name:         "division by zero",
			errorMsg:     "division by zero",
			expectedType: ErrorTypeRuntime,
		},
		{
			name:         "wrong number of arguments",
			errorMsg:     "if requires exactly 3 arguments",
			expectedType: ErrorTypeRuntime,
		},
		{
			name:         "empty list access",
			errorMsg:     "first called on empty list",
			expectedType: ErrorTypeRuntime,
		},
		{
			name:         "too many arguments",
			errorMsg:     "too many arguments provided",
			expectedType: ErrorTypeRuntime,
		},

		// File system errors
		{
			name:         "file not found",
			errorMsg:     "file not found: missing.lisp",
			expectedType: ErrorTypeFileSystem,
		},
		{
			name:         "permission denied",
			errorMsg:     "permission denied accessing file",
			expectedType: ErrorTypeFileSystem,
		},

		// Module errors
		{
			name:         "module not found",
			errorMsg:     "module 'math' not found",
			expectedType: ErrorTypeModule,
		},
		{
			name:         "import error",
			errorMsg:     "cannot import undefined symbol",
			expectedType: ErrorTypeModule,
		},

		// General errors
		{
			name:         "general error",
			errorMsg:     "something went wrong",
			expectedType: ErrorTypeGeneral,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ef.categorizeError(tt.errorMsg)
			if result != tt.expectedType {
				t.Errorf("categorizeError(%q) = %v, want %v", tt.errorMsg, result, tt.expectedType)
			}
		})
	}
}

func TestErrorFormatter_getErrorTypeLabel(t *testing.T) {
	ef := NewErrorFormatter()

	tests := []struct {
		errorType     ErrorType
		expectedLabel string
	}{
		{ErrorTypeSyntax, "Syntax Error"},
		{ErrorTypeRuntime, "Runtime Error"},
		{ErrorTypeUndefined, "Undefined Symbol"},
		{ErrorTypeTypeError, "Type Error"},
		{ErrorTypeFileSystem, "File System Error"},
		{ErrorTypeModule, "Module Error"},
		{ErrorTypeGeneral, "Error"},
	}

	for _, tt := range tests {
		t.Run(tt.expectedLabel, func(t *testing.T) {
			result := ef.getErrorTypeLabel(tt.errorType)
			if result != tt.expectedLabel {
				t.Errorf("getErrorTypeLabel(%v) = %q, want %q", tt.errorType, result, tt.expectedLabel)
			}
		})
	}
}

func TestErrorFormatter_FormatError(t *testing.T) {
	ef := NewErrorFormatter()

	tests := []struct {
		name     string
		err      error
		contains []string // Strings that should be present in the output
	}{
		{
			name:     "syntax error",
			err:      errors.New("unexpected token after expression"),
			contains: []string{"Syntax Error:", "unexpected token"},
		},
		{
			name:     "undefined symbol error",
			err:      errors.New("undefined symbol: foo"),
			contains: []string{"Undefined Symbol:", "undefined symbol: foo"},
		},
		{
			name:     "runtime error",
			err:      errors.New("division by zero"),
			contains: []string{"Runtime Error:", "division by zero"},
		},
		{
			name:     "nil error",
			err:      nil,
			contains: []string{}, // Should return empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ef.FormatError(tt.err)

			if tt.err == nil {
				if result != "" {
					t.Errorf("FormatError(nil) = %q, want empty string", result)
				}
				return
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("FormatError(%v) = %q, should contain %q", tt.err, result, substr)
				}
			}
		})
	}
}

func TestErrorFormatter_generateSuggestion(t *testing.T) {
	ef := NewErrorFormatter()

	tests := []struct {
		name           string
		errorMsg       string
		expectedSubstr string // Expected substring in suggestion
	}{
		{
			name:           "undefined symbol",
			errorMsg:       "undefined symbol: foo",
			expectedSubstr: "defined or imported",
		},
		{
			name:           "wrong number of arguments",
			errorMsg:       "if requires exactly 3 arguments",
			expectedSubstr: "function signature",
		},
		{
			name:           "unmatched parentheses",
			errorMsg:       "unmatched closing parenthesis",
			expectedSubstr: "balanced parentheses",
		},
		{
			name:           "division by zero",
			errorMsg:       "division by zero",
			expectedSubstr: "divisor is not zero",
		},
		{
			name:           "empty list",
			errorMsg:       "first called on empty list",
			expectedSubstr: "list has elements",
		},
		{
			name:           "not a function",
			errorMsg:       "value is not a function",
			expectedSubstr: "calling a function",
		},
		{
			name:           "file error",
			errorMsg:       "file not found",
			expectedSubstr: "file path exists",
		},
		{
			name:           "no suggestion",
			errorMsg:       "random error message",
			expectedSubstr: "", // Should return empty suggestion
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ef.generateSuggestion(tt.errorMsg)

			if tt.expectedSubstr == "" {
				if result != "" {
					t.Errorf("generateSuggestion(%q) = %q, want empty string", tt.errorMsg, result)
				}
				return
			}

			if !strings.Contains(result, tt.expectedSubstr) {
				t.Errorf("generateSuggestion(%q) = %q, should contain %q", tt.errorMsg, result, tt.expectedSubstr)
			}
		})
	}
}

func TestErrorFormatter_FormatErrorWithSuggestion(t *testing.T) {
	ef := NewErrorFormatter()

	tests := []struct {
		name       string
		err        error
		suggestion string
		contains   []string
	}{
		{
			name:       "error with suggestion",
			err:        errors.New("undefined symbol: foo"),
			suggestion: "Check if the symbol is defined",
			contains:   []string{"Undefined Symbol:", "undefined symbol: foo", "Suggestion:", "Check if the symbol is defined"},
		},
		{
			name:       "error without suggestion",
			err:        errors.New("some error"),
			suggestion: "",
			contains:   []string{"Error:", "some error"},
		},
		{
			name:       "nil error",
			err:        nil,
			suggestion: "Some suggestion",
			contains:   []string{}, // Should return empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ef.FormatErrorWithSuggestion(tt.err, tt.suggestion)

			if tt.err == nil {
				if result != "" {
					t.Errorf("FormatErrorWithSuggestion(nil, %q) = %q, want empty string", tt.suggestion, result)
				}
				return
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("FormatErrorWithSuggestion(%v, %q) = %q, should contain %q", tt.err, tt.suggestion, result, substr)
				}
			}

			// If no suggestion provided, should not contain "Suggestion:"
			if tt.suggestion == "" && strings.Contains(result, "Suggestion:") {
				t.Errorf("FormatErrorWithSuggestion(%v, %q) = %q, should not contain 'Suggestion:' when no suggestion provided", tt.err, tt.suggestion, result)
			}
		})
	}
}

func TestErrorFormatter_FormatErrorWithSmartSuggestion(t *testing.T) {
	ef := NewErrorFormatter()

	tests := []struct {
		name     string
		err      error
		contains []string
	}{
		{
			name:     "undefined symbol with auto suggestion",
			err:      errors.New("undefined symbol: foo"),
			contains: []string{"Undefined Symbol:", "undefined symbol: foo", "Suggestion:", "defined or imported"},
		},
		{
			name:     "syntax error with auto suggestion",
			err:      errors.New("unmatched closing parenthesis"),
			contains: []string{"Syntax Error:", "unmatched closing parenthesis", "Suggestion:", "balanced parentheses"},
		},
		{
			name:     "error without suggestion",
			err:      errors.New("random error"),
			contains: []string{"Error:", "random error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ef.FormatErrorWithSmartSuggestion(tt.err)

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("FormatErrorWithSmartSuggestion(%v) = %q, should contain %q", tt.err, result, substr)
				}
			}
		})
	}
}
