package minimal

import (
	"testing"
)

func TestFileLoadingAndStandardLibrary(t *testing.T) {
	// Create a bootstrapped REPL
	repl := NewBootstrappedREPL()

	// Test loading the standard library
	t.Log("Loading standard library...")
	result, err := LoadFile("stdlib.lisp", repl.Env)
	if err != nil {
		t.Fatalf("Error loading stdlib: %v", err)
	}
	t.Logf("Standard library loaded: %v", result)

	// Test some standard library functions
	testCases := []struct {
		expr     string
		expected string
	}{
		{"(def numbers [1 2 3 4 5])", "defined"},
		{"(nth numbers 2)", "3"},
		{"(conj numbers 6)", "[1 2 3 4 5 6]"},
		{"(sum numbers)", "15"},
	}

	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			parsed, err := repl.Parse(tc.expr)
			if err != nil {
				t.Fatalf("Parse error for '%s': %v", tc.expr, err)
			}

			result, err := Eval(parsed, repl.Env)
			if err != nil {
				t.Fatalf("Eval error for '%s': %v", tc.expr, err)
			}

			if result.String() != tc.expected {
				t.Errorf("Expected '%s', got '%s' for expression '%s'", tc.expected, result.String(), tc.expr)
			}
		})
	}
}
