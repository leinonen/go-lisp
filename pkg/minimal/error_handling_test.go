package minimal

import (
	"strings"
	"testing"
)

func TestEnhancedErrorHandling(t *testing.T) {
	repl := NewREPL()

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "undefined symbol",
			input:   "undefined-symbol",
			wantErr: true,
			errCheck: func(err error) bool {
				evalErr, ok := err.(*EvaluationError)
				if !ok {
					return false
				}
				return strings.Contains(evalErr.Message, "undefined-symbol") &&
					evalErr.SourceLocation.Line > 0
			},
		},
		{
			name:    "parse error with position",
			input:   "(unclosed list",
			wantErr: true,
			errCheck: func(err error) bool {
				parseErr, ok := err.(*ParseError)
				if !ok {
					return false
				}
				return (strings.Contains(parseErr.Message, "unclosed") ||
					strings.Contains(parseErr.Message, "unexpected")) &&
					parseErr.SourceLocation.Line > 0
			},
		},
		{
			name:    "function call error with stack trace",
			input:   "(+ 1 \"not-a-number\")",
			wantErr: true,
			errCheck: func(err error) bool {
				evalErr, ok := err.(*EvaluationError)
				if !ok {
					return false
				}
				return len(evalErr.StackTrace) > 0 &&
					strings.Contains(evalErr.Message, "numbers")
			},
		},
		{
			name:    "nested function error",
			input:   "(+ 1 (+ 2 \"bad\"))",
			wantErr: true,
			errCheck: func(err error) bool {
				evalErr, ok := err.(*EvaluationError)
				if !ok {
					return false
				}
				return len(evalErr.StackTrace) > 1 // Should have nested calls
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, pos, err := ParseWithPositions(tt.input, "<test>")

			if tt.wantErr && err != nil {
				// Parse error - check it
				if !tt.errCheck(err) {
					t.Errorf("Parse error check failed: %v", err)
				}
				return
			} else if err != nil && !tt.wantErr {
				t.Errorf("Unexpected parse error: %v", err)
				return
			}

			// Evaluate with context
			ctx := NewEvaluationContext()
			if pos != nil {
				ctx.SetLocation(pos.Line, pos.Column, pos.File)
			}
			ctx.SetExpression(tt.input)

			_, err = EvalWithContext(expr, repl.Env, ctx)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !tt.errCheck(err) {
					t.Errorf("Error check failed: %v", err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestErrorContextPreservation(t *testing.T) {
	repl := NewREPL()

	// Define a function that will cause an error
	input := `(define error-fn (fn [x] (+ x "bad")))`
	expr, _, err := ParseWithPositions(input, "<test>")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	ctx := NewEvaluationContext()
	_, err = EvalWithContext(expr, repl.Env, ctx)
	if err != nil {
		t.Fatalf("Error defining function: %v", err)
	}

	// Now call the function and check error context
	input = `(error-fn 42)`
	expr, pos, err := ParseWithPositions(input, "<test>")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	ctx = NewEvaluationContext()
	if pos != nil {
		ctx.SetLocation(pos.Line, pos.Column, pos.File)
	}
	ctx.SetExpression(input)

	_, err = EvalWithContext(expr, repl.Env, ctx)
	if err == nil {
		t.Fatal("Expected error but got none")
	}

	evalErr, ok := err.(*EvaluationError)
	if !ok {
		t.Fatalf("Expected EvaluationError, got %T", err)
	}

	// Check that we have stack trace information
	if len(evalErr.StackTrace) == 0 {
		t.Error("Expected stack trace but got none")
	}

	// Check that we have source location
	if evalErr.SourceLocation.Line == 0 {
		t.Error("Expected source location but got none")
	}

	// Print the error to see it in action
	t.Logf("Error output: %v", evalErr)
}

func TestParseErrorFormatting(t *testing.T) {
	input := `(+ 1 2
(unclosed`

	_, _, err := ParseWithPositions(input, "test.lisp")
	if err == nil {
		t.Fatal("Expected parse error but got none")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("Expected ParseError, got %T", err)
	}

	errorStr := parseErr.Error()

	// Check that the error includes position information
	if !strings.Contains(errorStr, "test.lisp") {
		t.Error("Error should include filename")
	}

	if !strings.Contains(errorStr, ":") {
		t.Error("Error should include line/column information")
	}

	// We should get some indication of what went wrong
	if !strings.Contains(errorStr, "unclosed") && !strings.Contains(errorStr, "unexpected") {
		t.Error("Error should show information about the problem")
	}

	t.Logf("Parse error output: %v", parseErr)
}
