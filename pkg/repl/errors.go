package repl

import (
	"errors"
	"strings"

	"github.com/fatih/color"
	"github.com/leinonen/go-lisp/pkg/types"
)

// ErrorType represents different categories of errors for color coding
type ErrorType int

const (
	ErrorTypeSyntax ErrorType = iota
	ErrorTypeRuntime
	ErrorTypeUndefined
	ErrorTypeTypeError
	ErrorTypeFileSystem
	ErrorTypeModule
	ErrorTypeGeneral
)

// ErrorFormatter handles colored error output for the REPL
type ErrorFormatter struct {
	syntaxColor    *color.Color
	runtimeColor   *color.Color
	undefinedColor *color.Color
	typeColor      *color.Color
	fileColor      *color.Color
	moduleColor    *color.Color
	generalColor   *color.Color
	prefixColor    *color.Color
}

// NewErrorFormatter creates a new error formatter with predefined colors
func NewErrorFormatter() *ErrorFormatter {
	return &ErrorFormatter{
		syntaxColor:    color.New(color.FgRed, color.Bold),     // Bright red for syntax errors
		runtimeColor:   color.New(color.FgMagenta, color.Bold), // Magenta for runtime errors
		undefinedColor: color.New(color.FgYellow, color.Bold),  // Yellow for undefined symbols
		typeColor:      color.New(color.FgCyan, color.Bold),    // Cyan for type errors
		fileColor:      color.New(color.FgBlue, color.Bold),    // Blue for file system errors
		moduleColor:    color.New(color.FgGreen, color.Bold),   // Green for module errors
		generalColor:   color.New(color.FgWhite, color.Bold),   // White for general errors
		prefixColor:    color.New(color.FgRed, color.Bold),     // Red for "Error:" prefix
	}
}

// categorizeError determines the error type based on the error message
func (ef *ErrorFormatter) categorizeError(errMsg string) ErrorType {
	errLower := strings.ToLower(errMsg)

	// Syntax errors
	if strings.Contains(errLower, "unexpected token") ||
		strings.Contains(errLower, "syntax error") ||
		strings.Contains(errLower, "parse error") ||
		strings.Contains(errLower, "unmatched") ||
		strings.Contains(errLower, "expected") ||
		strings.Contains(errLower, "empty input") {
		return ErrorTypeSyntax
	}

	// Module errors (check before undefined symbols since import errors can contain "undefined")
	if strings.Contains(errLower, "module") ||
		strings.Contains(errLower, "import") ||
		strings.Contains(errLower, "export") ||
		strings.Contains(errLower, "load") {
		return ErrorTypeModule
	}

	// Undefined symbol errors
	if strings.Contains(errLower, "undefined symbol") ||
		strings.Contains(errLower, "undefined variable") ||
		strings.Contains(errLower, "unknown function") ||
		strings.Contains(errLower, "no help available") {
		return ErrorTypeUndefined
	}

	// Type errors
	if strings.Contains(errLower, "requires a") ||
		strings.Contains(errLower, "must be a") ||
		strings.Contains(errLower, "wrong type") ||
		strings.Contains(errLower, "type error") ||
		strings.Contains(errLower, "not a list") ||
		strings.Contains(errLower, "not a number") ||
		strings.Contains(errLower, "not a string") ||
		strings.Contains(errLower, "not a boolean") ||
		strings.Contains(errLower, "not a function") {
		return ErrorTypeTypeError
	}

	// Runtime errors
	if strings.Contains(errLower, "division by zero") ||
		strings.Contains(errLower, "out of bounds") ||
		strings.Contains(errLower, "empty list") ||
		strings.Contains(errLower, "wrong number of arguments") ||
		strings.Contains(errLower, "requires exactly") ||
		strings.Contains(errLower, "too few") ||
		strings.Contains(errLower, "too many") ||
		strings.Contains(errLower, "stack overflow") ||
		strings.Contains(errLower, "infinite") {
		return ErrorTypeRuntime
	}

	// File system errors
	if strings.Contains(errLower, "file") ||
		strings.Contains(errLower, "directory") ||
		strings.Contains(errLower, "no such") ||
		strings.Contains(errLower, "permission denied") ||
		strings.Contains(errLower, "cannot read") ||
		strings.Contains(errLower, "cannot write") {
		return ErrorTypeFileSystem
	}

	return ErrorTypeGeneral
}

// getColorForErrorType returns the appropriate color for an error type
func (ef *ErrorFormatter) getColorForErrorType(errorType ErrorType) *color.Color {
	switch errorType {
	case ErrorTypeSyntax:
		return ef.syntaxColor
	case ErrorTypeRuntime:
		return ef.runtimeColor
	case ErrorTypeUndefined:
		return ef.undefinedColor
	case ErrorTypeTypeError:
		return ef.typeColor
	case ErrorTypeFileSystem:
		return ef.fileColor
	case ErrorTypeModule:
		return ef.moduleColor
	default:
		return ef.generalColor
	}
}

// getErrorTypeLabel returns a human-readable label for the error type
func (ef *ErrorFormatter) getErrorTypeLabel(errorType ErrorType) string {
	switch errorType {
	case ErrorTypeSyntax:
		return "Syntax Error"
	case ErrorTypeRuntime:
		return "Runtime Error"
	case ErrorTypeUndefined:
		return "Undefined Symbol"
	case ErrorTypeTypeError:
		return "Type Error"
	case ErrorTypeFileSystem:
		return "File System Error"
	case ErrorTypeModule:
		return "Module Error"
	default:
		return "Error"
	}
}

// FormatError formats an error with appropriate colors and categorization
func (ef *ErrorFormatter) FormatError(err error) string {
	if err == nil {
		return ""
	}

	errMsg := err.Error()
	errorType := ef.categorizeError(errMsg)
	errorColor := ef.getColorForErrorType(errorType)
	errorLabel := ef.getErrorTypeLabel(errorType)

	// Check if this is a positional error
	var posErr *types.PositionalError
	if errors.As(err, &posErr) {
		// For positional errors, format with enhanced line information
		prefix := ef.prefixColor.Sprintf("%s:", errorLabel)
		locationColor := color.New(color.FgHiBlue, color.Bold)
		location := locationColor.Sprintf(" (line %d, column %d)", posErr.Position.Line, posErr.Position.Column)
		message := errorColor.Sprintf(" %s", posErr.Message)
		return prefix + location + message
	}

	// Check if the error message already contains line information
	if strings.Contains(errMsg, "line ") && strings.Contains(errMsg, "column ") {
		// Extract and format existing line/column information
		prefix := ef.prefixColor.Sprintf("%s:", errorLabel)
		message := errorColor.Sprintf(" %s", errMsg)
		return prefix + message
	}

	// Standard error formatting
	prefix := ef.prefixColor.Sprintf("%s:", errorLabel)
	message := errorColor.Sprintf(" %s", errMsg)

	return prefix + message
}

// FormatErrorWithSuggestion formats an error with a suggestion
func (ef *ErrorFormatter) FormatErrorWithSuggestion(err error, suggestion string) string {
	if err == nil {
		return ""
	}

	baseError := ef.FormatError(err)
	if suggestion == "" {
		return baseError
	}

	suggestionColor := color.New(color.FgHiBlack, color.Italic)
	suggestionText := suggestionColor.Sprintf("\n  Suggestion: %s", suggestion)

	return baseError + suggestionText
}

// generateSuggestion provides helpful suggestions based on the error message
func (ef *ErrorFormatter) generateSuggestion(errMsg string) string {
	errLower := strings.ToLower(errMsg)

	if strings.Contains(errLower, "undefined symbol") {
		return "Check if the symbol is defined or imported from a module"
	}

	if strings.Contains(errLower, "wrong number of arguments") ||
		strings.Contains(errLower, "requires exactly") ||
		strings.Contains(errLower, "too few") ||
		strings.Contains(errLower, "too many") {
		return "Check the function signature with (help function-name)"
	}

	if strings.Contains(errLower, "unmatched") || strings.Contains(errLower, "unexpected token") {
		return "Check for balanced parentheses and proper syntax"
	}

	if strings.Contains(errLower, "division by zero") {
		return "Ensure the divisor is not zero"
	}

	if strings.Contains(errLower, "empty list") {
		return "Check if the list has elements before accessing them"
	}

	if strings.Contains(errLower, "not a") && strings.Contains(errLower, "function") {
		return "Make sure you're calling a function, not a variable"
	}

	if strings.Contains(errLower, "file") || strings.Contains(errLower, "directory") {
		return "Check if the file path exists and is accessible"
	}

	return ""
}

// FormatErrorWithSmartSuggestion formats an error with an automatically generated suggestion
func (ef *ErrorFormatter) FormatErrorWithSmartSuggestion(err error) string {
	if err == nil {
		return ""
	}

	errMsg := err.Error()
	suggestion := ef.generateSuggestion(errMsg)
	return ef.FormatErrorWithSuggestion(err, suggestion)
}
