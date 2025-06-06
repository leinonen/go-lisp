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
