package repl

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// createTestRegistry creates a registry with some basic functions for testing
func createTestRegistry() registry.FunctionRegistry {
	reg := registry.NewRegistry()

	// Add some common functions for testing
	testFunctions := []struct {
		name     string
		category string
		arity    int
		help     string
	}{
		{"map", "functional", 2, "Apply function to list"},
		{"filter", "functional", 2, "Filter list with predicate"},
		{"reduce", "functional", 3, "Reduce list with function"},
		{"+", "arithmetic", -1, "Add numbers"},
		{"-", "arithmetic", -1, "Subtract numbers"},
		{"*", "arithmetic", -1, "Multiply numbers"},
		{"help", "core", -1, "Show help"},
		{"def", "core", 2, "Define variable"},
		{"if", "control", 3, "Conditional expression"},
	}

	for _, tf := range testFunctions {
		fn := functions.NewFunction(tf.name, tf.category, tf.arity, tf.help, nil)
		reg.Register(fn)
	}

	return reg
}

func TestCompletionProvider(t *testing.T) {
	// Create a test environment
	env := evaluator.NewEnvironment()

	// Create a test registry
	reg := createTestRegistry()

	// Add some user-defined functions
	env.Set("my-function", types.FunctionValue{
		Params: []string{"x"},
		Body:   &types.NumberExpr{Value: 42},
		Env:    env,
	})
	env.Set("another-func", types.FunctionValue{
		Params: []string{"x", "y"},
		Body:   &types.NumberExpr{Value: 24},
		Env:    env,
	})
	env.Set("my-variable", types.NumberValue(123))

	// Create completion provider with registry
	cp := NewCompletionProviderWithRegistry(env, reg)
	t.Run("complete builtin functions", func(t *testing.T) {
		// Should only complete when in function position (after '(')
		completions := cp.GetCompletions("(ma", 3)

		// Should include 'map' builtin
		found := false
		for _, comp := range completions {
			if comp == "map" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'map' in completions for '(ma'")
		}
	})

	t.Run("complete user-defined functions", func(t *testing.T) {
		// Should only complete when in function position
		completions := cp.GetCompletions("(my-", 4)

		// Should include 'my-function'
		foundFunc := false
		for _, comp := range completions {
			if comp == "my-function" {
				foundFunc = true
			}
		}
		if !foundFunc {
			t.Error("Expected 'my-function' in completions for '(my-'")
		}
	})

	t.Run("complete arithmetic functions", func(t *testing.T) {
		// Should only complete when in function position
		completions := cp.GetCompletions("(+", 2)

		// Should include '+'
		found := false
		for _, comp := range completions {
			if comp == "+" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected '+' in completions for '(+'")
		}
	})

	t.Run("no completions outside function position", func(t *testing.T) {
		// Should not complete when not after '('
		completions := cp.GetCompletions("ma", 2)

		if len(completions) > 0 {
			t.Errorf("Expected no completions for 'ma' (not after paren), got %v", completions)
		}
	})

	t.Run("no completions in argument position", func(t *testing.T) {
		// Should not complete function arguments
		completions := cp.GetCompletions("(+ ma", 5)

		if len(completions) > 0 {
			t.Errorf("Expected no completions in argument position, got %v", completions)
		}
	})

	t.Run("completion right after open paren with no prefix", func(t *testing.T) {
		// Should complete all functions when right after '(' with no prefix
		completions := cp.GetCompletions("(", 1)

		// Should include some built-in functions
		foundMap := false
		foundPlus := false
		for _, comp := range completions {
			if comp == "map" {
				foundMap = true
			}
			if comp == "+" {
				foundPlus = true
			}
		}
		if !foundMap {
			t.Error("Expected 'map' in completions after '('")
		}
		if !foundPlus {
			t.Error("Expected '+' in completions after '('")
		}
	})

	t.Run("help function should be in completions", func(t *testing.T) {
		// Should include 'help' when typing 'hel'
		completions := cp.GetCompletions("(hel", 4)

		// Should include 'help'
		found := false
		for _, comp := range completions {
			if comp == "help" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'help' in completions for '(hel'")
		}
	})
}

func TestExtractCurrentWord(t *testing.T) {
	env := evaluator.NewEnvironment()
	reg := createTestRegistry()
	cp := NewCompletionProviderWithRegistry(env, reg)

	tests := []struct {
		line     string
		pos      int
		expected string
	}{
		{"(+ 1 ma", 7, "ma"},
		{"(defn test-fun", 15, "test-fun"},
		{"(map filt", 9, "filt"},
		{"hello world", 5, "hello"},
		{"hello world", 11, "world"},
		{"(+ (* 2 3) red", 14, "red"},
		{"", 0, ""},
		{"()", 1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := cp.extractCurrentWord(tt.line, tt.pos)
			if result != tt.expected {
				t.Errorf("extractCurrentWord(%q, %d) = %q, expected %q",
					tt.line, tt.pos, result, tt.expected)
			}
		})
	}
}

func TestIsSymbolChar(t *testing.T) {
	env := evaluator.NewEnvironment()
	reg := createTestRegistry()
	cp := NewCompletionProviderWithRegistry(env, reg)

	validChars := []rune{'a', 'Z', '0', '9', '-', '_', '?', '!', '+', '*', '/', '=', '<', '>', '.', '%'}
	invalidChars := []rune{' ', '\t', '\n', '(', ')', '[', ']', '{', '}', '"', '\'', ';', ','}

	for _, ch := range validChars {
		if !cp.isSymbolChar(ch) {
			t.Errorf("Expected '%c' to be a valid symbol character", ch)
		}
	}

	for _, ch := range invalidChars {
		if cp.isSymbolChar(ch) {
			t.Errorf("Expected '%c' to be an invalid symbol character", ch)
		}
	}
}

func TestLispAwareCompletion(t *testing.T) {
	// Create a test environment
	env := evaluator.NewEnvironment()

	// Add some user-defined functions
	env.Set("my-add", types.FunctionValue{
		Params: []string{"x", "y"},
		Body:   &types.NumberExpr{Value: 42},
		Env:    env,
	})
	env.Set("my-variable", types.NumberValue(123))

	// Create completion provider
	reg := createTestRegistry()
	cp := NewCompletionProviderWithRegistry(env, reg)

	t.Run("completion after open paren", func(t *testing.T) {
		// Should complete functions when right after '('
		completions := cp.GetCompletions("(ma", 3)

		// Should include 'map' builtin function
		found := false
		for _, comp := range completions {
			if comp == "map" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'map' in completions after '(ma'")
		}
	})

	t.Run("completion in function position with spaces", func(t *testing.T) {
		// Should complete functions when after '( ' (with space)
		completions := cp.GetCompletions("( my", 4)

		// Should include user-defined function 'my-add'
		found := false
		for _, comp := range completions {
			if comp == "my-add" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'my-add' in completions after '( my'")
		}
	})

	t.Run("completion in argument position", func(t *testing.T) {
		// Should NOT complete in argument position
		completions := cp.GetCompletions("(+ my", 5)

		// Should be empty since we're not in function position
		if len(completions) > 0 {
			t.Errorf("Expected no completions in argument position, got %v", completions)
		}
	})

	t.Run("nested expression completion", func(t *testing.T) {
		// Should complete functions at the start of nested expressions
		completions := cp.GetCompletions("(+ 1 (ma", 8)

		// Should include 'map' in nested function position
		found := false
		for _, comp := range completions {
			if comp == "map" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'map' in completions in nested expression")
		}
	})
}

func TestCompletionContext(t *testing.T) {
	env := evaluator.NewEnvironment()
	reg := createTestRegistry()
	cp := NewCompletionProviderWithRegistry(env, reg)

	tests := []struct {
		line               string
		pos                int
		expectedFuncPos    bool
		expectedAfterParen bool
		description        string
	}{
		{"(ma", 3, true, true, "right after open paren"},
		{"( ma", 4, true, false, "after paren with space"},
		{"(+ 1 2", 6, false, false, "in argument position"},
		{"(+ (ma", 6, true, true, "nested function position"},
		{"(define x (ma", 13, true, true, "nested in define"},
		{"", 0, false, false, "empty line"},
		{"ma", 2, false, false, "top level symbol"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			context := cp.analyzeContext(tt.line, tt.pos)

			if context.inFunctionPosition != tt.expectedFuncPos {
				t.Errorf("For %q at pos %d: expected inFunctionPosition=%v, got %v",
					tt.line, tt.pos, tt.expectedFuncPos, context.inFunctionPosition)
			}

			if context.afterOpenParen != tt.expectedAfterParen {
				t.Errorf("For %q at pos %d: expected afterOpenParen=%v, got %v",
					tt.line, tt.pos, tt.expectedAfterParen, context.afterOpenParen)
			}
		})
	}
}

func TestHelpArgumentCompletion(t *testing.T) {
	// Create a test environment
	env := evaluator.NewEnvironment()
	reg := createTestRegistry()
	cp := NewCompletionProviderWithRegistry(env, reg)

	t.Run("completion for help function arguments", func(t *testing.T) {
		// Should complete function names when typing help arguments
		completions := cp.GetCompletions("(help ma", 8)

		// Should include 'map' builtin function
		found := false
		for _, comp := range completions {
			if comp == "map" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'map' in completions for '(help ma'")
		}
	})

	t.Run("completion for help with plus sign", func(t *testing.T) {
		// Should complete '+' when typing it as argument to help
		completions := cp.GetCompletions("(help +", 7)

		// Should include '+'
		found := false
		for _, comp := range completions {
			if comp == "+" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected '+' in completions for '(help +'")
		}
	})

	t.Run("completion for help with filter prefix", func(t *testing.T) {
		// Should complete 'filter' when typing 'fil' as argument to help
		completions := cp.GetCompletions("(help fil", 9)

		// Should include 'filter'
		found := false
		for _, comp := range completions {
			if comp == "filter" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'filter' in completions for '(help fil'")
		}
	})

	t.Run("completion for help with spaces", func(t *testing.T) {
		// Should work with extra spaces
		completions := cp.GetCompletions("(help  ma", 9)

		// Should include 'map'
		found := false
		for _, comp := range completions {
			if comp == "map" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'map' in completions for '(help  ma' (with extra space)")
		}
	})

	t.Run("no completion for second argument to help", func(t *testing.T) {
		// Should not complete second argument to help (it only takes one)
		completions := cp.GetCompletions("(help map ma", 12)

		// Should be empty since help only takes one argument
		if len(completions) > 0 {
			t.Errorf("Expected no completions for second argument to help, got %v", completions)
		}
	})
}
