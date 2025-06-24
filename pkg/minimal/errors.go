package minimal

import (
	"fmt"
	"strings"
)

// Position represents a location in source code
type Position struct {
	Line   int
	Column int
	File   string
}

func (p Position) String() string {
	if p.File != "" {
		return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Column)
	}
	return fmt.Sprintf("line %d, column %d", p.Line, p.Column)
}

// EvaluationError represents an error during evaluation with enhanced context
type EvaluationError struct {
	Message        string
	StackTrace     []string
	SourceLocation Position
	Expression     string
}

func (e *EvaluationError) Error() string {
	var result strings.Builder

	// Main error message
	result.WriteString(e.Message)

	// Source location if available
	if e.SourceLocation.Line > 0 {
		result.WriteString(fmt.Sprintf(" at %s", e.SourceLocation))
	}

	// Expression if available
	if e.Expression != "" {
		result.WriteString(fmt.Sprintf("\n  in expression: %s", e.Expression))
	}

	// Stack trace if available
	if len(e.StackTrace) > 0 {
		result.WriteString("\nStack trace:")
		for i, frame := range e.StackTrace {
			result.WriteString(fmt.Sprintf("\n  %d: %s", i, frame))
		}
	}

	return result.String()
}

// ParseError represents an error during parsing
type ParseError struct {
	Message        string
	SourceLocation Position
	Input          string
}

func (e *ParseError) Error() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("Parse error: %s", e.Message))

	if e.SourceLocation.Line > 0 {
		result.WriteString(fmt.Sprintf(" at %s", e.SourceLocation))
	}

	if e.Input != "" {
		// Show the problematic line with a caret indicator
		lines := strings.Split(e.Input, "\n")
		if e.SourceLocation.Line > 0 && e.SourceLocation.Line <= len(lines) {
			problemLine := lines[e.SourceLocation.Line-1]
			result.WriteString(fmt.Sprintf("\n  %s", problemLine))

			// Add caret indicator
			if e.SourceLocation.Column > 0 && e.SourceLocation.Column <= len(problemLine)+1 {
				indicator := strings.Repeat(" ", e.SourceLocation.Column-1) + "^"
				result.WriteString(fmt.Sprintf("\n  %s", indicator))
			}
		}
	}

	return result.String()
}

// EvaluationContext tracks context during evaluation for better error reporting
type EvaluationContext struct {
	StackTrace     []string
	SourceLocation Position
	CurrentExpr    string
	Filename       string
}

// NewEvaluationContext creates a new evaluation context
func NewEvaluationContext() *EvaluationContext {
	return &EvaluationContext{
		StackTrace: make([]string, 0),
	}
}

// PushFrame adds a new frame to the evaluation stack
func (ctx *EvaluationContext) PushFrame(frame string) {
	ctx.StackTrace = append(ctx.StackTrace, frame)
}

// PopFrame removes the top frame from the evaluation stack
func (ctx *EvaluationContext) PopFrame() {
	if len(ctx.StackTrace) > 0 {
		ctx.StackTrace = ctx.StackTrace[:len(ctx.StackTrace)-1]
	}
}

// SetLocation sets the current source location
func (ctx *EvaluationContext) SetLocation(line, column int, file string) {
	ctx.SourceLocation = Position{Line: line, Column: column, File: file}
}

// SetExpression sets the current expression being evaluated
func (ctx *EvaluationContext) SetExpression(expr string) {
	ctx.CurrentExpr = expr
}

// CreateError creates an EvaluationError with context
func (ctx *EvaluationContext) CreateError(message string) *EvaluationError {
	// Create a copy of the stack trace
	stackCopy := make([]string, len(ctx.StackTrace))
	copy(stackCopy, ctx.StackTrace)

	return &EvaluationError{
		Message:        message,
		StackTrace:     stackCopy,
		SourceLocation: ctx.SourceLocation,
		Expression:     ctx.CurrentExpr,
	}
}

// WrapError wraps a regular error with evaluation context
func (ctx *EvaluationContext) WrapError(err error) *EvaluationError {
	if evalErr, ok := err.(*EvaluationError); ok {
		// Already an evaluation error, just return it
		return evalErr
	}

	return ctx.CreateError(err.Error())
}
