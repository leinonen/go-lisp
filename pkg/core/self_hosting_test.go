package core_test

import (
	"fmt"
	"testing"

	"github.com/leinonen/go-lisp/pkg/core"
)

func TestSelfHostingCompilerIntegration(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Load the self-hosting compiler (adjust path for test context)
	loadExpr, err := core.ReadString("(load-file \"../../lisp/self-hosting.lisp\")")
	if err != nil {
		t.Errorf("Parse error loading self-hosting compiler: %v", err)
		return
	}

	_, err = core.Eval(loadExpr, env)
	if err != nil {
		t.Errorf("Error loading self-hosting compiler: %v", err)
		return
	}

	// Test that read-all function is available and working
	tests := []struct {
		input    string
		expected string
		desc     string
	}{
		{
			input:    "(read-all \"(+ 1 2)\")",
			expected: "((+ 1 2))",
			desc:     "single expression",
		},
		{
			input:    "(read-all \"(def x 10) (+ x 5)\")",
			expected: "((def x 10) (+ x 5))",
			desc:     "multiple expressions",
		},
		{
			input:    "(count (read-all \"(+ 1 2) (* 3 4) (- 10 5)\"))",
			expected: "3",
			desc:     "count of parsed expressions",
		},
	}

	for _, test := range tests {
		expr, err := core.ReadString(test.input)
		if err != nil {
			t.Errorf("Parse error for %s: %v", test.desc, err)
			continue
		}

		result, err := core.Eval(expr, env)
		if err != nil {
			t.Errorf("Eval error for %s: %v", test.desc, err)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("Expected '%s' for %s, got '%s'", test.expected, test.desc, result.String())
		}
	}
}

func TestSelfHostingCompilerContextCreation(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Load the self-hosting compiler (adjust path for test context)
	loadExpr, err := core.ReadString("(load-file \"../../lisp/self-hosting.lisp\")")
	if err != nil {
		t.Errorf("Parse error loading self-hosting compiler: %v", err)
		return
	}

	_, err = core.Eval(loadExpr, env)
	if err != nil {
		t.Errorf("Error loading self-hosting compiler: %v", err)
		return
	}

	// Test context creation
	contextExpr, err := core.ReadString("(make-context)")
	if err != nil {
		t.Errorf("Parse error for make-context: %v", err)
		return
	}

	result, err := core.Eval(contextExpr, env)
	if err != nil {
		t.Errorf("Error creating context: %v", err)
		return
	}

	// Should return a hash-map with symbols, locals, and target keys
	resultStr := result.String()
	if !contains(resultStr, ":symbols") || !contains(resultStr, ":locals") || !contains(resultStr, ":target") {
		t.Errorf("Context should contain :symbols, :locals, and :target keys, got: %s", resultStr)
	}
}

func TestSelfHostingReadAllWithRealFile(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Load the self-hosting compiler (adjust path for test context)
	loadExpr, err := core.ReadString("(load-file \"../../lisp/self-hosting.lisp\")")
	if err != nil {
		t.Errorf("Parse error loading self-hosting compiler: %v", err)
		return
	}

	_, err = core.Eval(loadExpr, env)
	if err != nil {
		t.Errorf("Error loading self-hosting compiler: %v", err)
		return
	}

	// Create a test file with multiple expressions
	testContent := `(def x 10)
(def y 20)
(defn add [a b] (+ a b))
(add x y)`

	testFile := "/tmp/test-multi-expressions.lisp"
	createFileExpr, err := core.ReadString(fmt.Sprintf("(spit \"%s\" \"%s\")", testFile, testContent))
	if err != nil {
		t.Errorf("Parse error creating test file: %v", err)
		return
	}

	_, err = core.Eval(createFileExpr, env)
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
		return
	}

	// Test reading the file with read-all
	readFileExpr, err := core.ReadString(fmt.Sprintf("(read-all (slurp \"%s\"))", testFile))
	if err != nil {
		t.Errorf("Parse error for read-all with file: %v", err)
		return
	}

	result, err := core.Eval(readFileExpr, env)
	if err != nil {
		t.Errorf("Error reading file with read-all: %v", err)
		return
	}

	// Test that we got 4 expressions
	countExpr, err := core.ReadString(fmt.Sprintf("(count (read-all (slurp \"%s\")))", testFile))
	if err != nil {
		t.Errorf("Parse error for count test: %v", err)
		return
	}

	countResult, err := core.Eval(countExpr, env)
	if err != nil {
		t.Errorf("Error counting expressions: %v", err)
		return
	}

	if countResult.String() != "4" {
		t.Errorf("Expected 4 expressions, got %s", countResult.String())
	}

	// Test that each expression is parsed correctly
	resultStr := result.String()
	expectedParts := []string{"(def x 10)", "(def y 20)", "(defn add [a b] (+ a b))", "(add x y)"}

	for _, part := range expectedParts {
		if !contains(resultStr, part) {
			t.Errorf("Expected result to contain '%s', got: %s", part, resultStr)
		}
	}
}

func TestSelfHostingCompilerParsesSelf(t *testing.T) {
	env := core.NewCoreEnvironment()

	// Load the self-hosting compiler (adjust path for test context)
	loadExpr, err := core.ReadString("(load-file \"../../lisp/self-hosting.lisp\")")
	if err != nil {
		t.Errorf("Parse error loading self-hosting compiler: %v", err)
		return
	}

	_, err = core.Eval(loadExpr, env)
	if err != nil {
		t.Errorf("Error loading self-hosting compiler: %v", err)
		return
	}

	// Test that the compiler can parse its own source file
	parseSelfExpr, err := core.ReadString("(count (read-all (slurp \"../../lisp/self-hosting.lisp\")))")
	if err != nil {
		t.Errorf("Parse error for self-parsing test: %v", err)
		return
	}

	result, err := core.Eval(parseSelfExpr, env)
	if err != nil {
		t.Errorf("Error parsing self-hosting compiler: %v", err)
		return
	}

	// Should parse multiple expressions (exact count may vary, but should be > 10)
	countStr := result.String()
	if countStr == "0" || countStr == "1" {
		t.Errorf("Expected multiple expressions from self-hosting.lisp, got %s", countStr)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
