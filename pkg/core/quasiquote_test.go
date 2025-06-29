package core_test

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/core"
)

func TestEvalQuasiquote(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Set up test variables
	expr, _ := core.ReadString("(def x 42)")
	core.Eval(expr, env)
	expr, _ = core.ReadString("(def y (list 1 2 3))")
	core.Eval(expr, env)

	tests := []struct {
		input    string
		expected string
	}{
		// Basic quasiquote - should act like quote
		{"`42", "42"},
		{"`hello", "hello"},
		{"`(a b c)", "(a b c)"},

		// Simple unquote
		{"`~x", "42"},
		{"`(a ~x c)", "(a 42 c)"},
		{"`[a ~x c]", "[a 42 c]"},

		// Unquote-splicing with lists
		{"`(a ~@y d)", "(a 1 2 3 d)"},
		{"`[a ~@y d]", "[a 1 2 3 d]"},

		// Nested quasiquotes and complex examples
		{"`(+ 1 ~(+ 2 3))", "(+ 1 5)"},
		{"`(list ~@(list 1 2) 3)", "(list 1 2 3)"},

		// HashMap quasiquote
		{"`{:a ~x :b 2}", "{:a 42 :b 2}"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for '%s': %v", test.input, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, result.String())
		}
	}
}

func TestQuasiquoteReader(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Test that the reader correctly parses quasiquote syntax
		{"`x", "(quasiquote x)"},
		{"~x", "(unquote x)"},
		{"~@x", "(unquote-splicing x)"},
		{"`(a ~b)", "(quasiquote (a (unquote b)))"},
		{"`[a ~@b]", "(quasiquote [a (unquote-splicing b)])"},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for '%s': %v", test.input, err)
			continue
		}

		if expr.String() != test.expected {
			t.Errorf("Expected '%s' for input '%s', got '%s'", test.expected, test.input, expr.String())
		}
	}
}
