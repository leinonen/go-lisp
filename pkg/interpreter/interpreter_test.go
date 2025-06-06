package interpreter

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:     "simple number",
			input:    "42",
			expected: types.NumberValue(42),
		},
		{
			name:     "simple addition",
			input:    "(+ 1 2)",
			expected: types.NumberValue(3),
		},
		{
			name:     "nested expression",
			input:    "(+ (* 2 3) 4)",
			expected: types.NumberValue(10),
		},
		{
			name:     "boolean expression",
			input:    "(< 3 5)",
			expected: types.BooleanValue(true),
		},
		{
			name:     "string literal",
			input:    `"hello world"`,
			expected: types.StringValue("hello world"),
		},
		{
			name:     "if expression",
			input:    "(if (< 3 5) 42 0)",
			expected: types.NumberValue(42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()
			result, err := interpreter.Interpret(tt.input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInterpreterDefine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Value
	}{
		{
			name:     "simple variable definition",
			input:    "(define x 42)",
			expected: types.NumberValue(42),
		},
		{
			name:     "define with expression",
			input:    "(define y (+ 10 20))",
			expected: types.NumberValue(30),
		},
		{
			name:     "define string variable",
			input:    `(define greeting "hello world")`,
			expected: types.StringValue("hello world"),
		},
		{
			name:     "define boolean variable",
			input:    "(define flag #t)",
			expected: types.BooleanValue(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()
			result, err := interpreter.Interpret(tt.input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valuesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInterpreterDefineAndUse(t *testing.T) {
	interpreter := NewInterpreter()

	// Define a variable
	_, err := interpreter.Interpret("(define x 10)")
	if err != nil {
		t.Fatalf("unexpected error defining x: %v", err)
	}

	// Use the variable
	result, err := interpreter.Interpret("x")
	if err != nil {
		t.Fatalf("unexpected error accessing x: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(10)) {
		t.Errorf("expected 10, got %v", result)
	}

	// Use the variable in an expression
	result, err = interpreter.Interpret("(+ x 5)")
	if err != nil {
		t.Fatalf("unexpected error in expression: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(15)) {
		t.Errorf("expected 15, got %v", result)
	}

	// Define another variable using the first
	result, err = interpreter.Interpret("(define y (* x 3))")
	if err != nil {
		t.Fatalf("unexpected error defining y: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(30)) {
		t.Errorf("expected 30, got %v", result)
	}

	// Use both variables
	result, err = interpreter.Interpret("(+ x y)")
	if err != nil {
		t.Fatalf("unexpected error in final expression: %v", err)
	}

	if !valuesEqual(result, types.NumberValue(40)) {
		t.Errorf("expected 40, got %v", result)
	}
}

// Helper function to compare values
func valuesEqual(a, b types.Value) bool {
	switch va := a.(type) {
	case types.NumberValue:
		if vb, ok := b.(types.NumberValue); ok {
			return va == vb
		}
	case types.StringValue:
		if vb, ok := b.(types.StringValue); ok {
			return va == vb
		}
	case types.BooleanValue:
		if vb, ok := b.(types.BooleanValue); ok {
			return va == vb
		}
	}
	return false
}
