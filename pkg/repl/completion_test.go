package repl

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestCompletionProvider(t *testing.T) {
	// Create a test environment
	env := evaluator.NewEnvironment()

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

	// Create completion provider
	cp := NewCompletionProvider(env)
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

	t.Run("builtins function should be in completions", func(t *testing.T) {
		// Should include 'builtins' when typing 'buil'
		completions := cp.GetCompletions("(buil", 5)

		// Should include 'builtins'
		found := false
		for _, comp := range completions {
			if comp == "builtins" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'builtins' in completions for '(buil'")
		}
	})
}

func TestExtractCurrentWord(t *testing.T) {
	env := evaluator.NewEnvironment()
	cp := NewCompletionProvider(env)

	tests := []struct {
		line     string
		pos      int
		expected string
	}{
		{"(+ 1 ma", 7, "ma"},
		{"(defun test-fun", 15, "test-fun"},
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
	cp := NewCompletionProvider(env)

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
	cp := NewCompletionProvider(env)

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
	cp := NewCompletionProvider(env)

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

func TestBuiltinsArgumentCompletion(t *testing.T) {
	// Create a test environment
	env := evaluator.NewEnvironment()
	cp := NewCompletionProvider(env)

	t.Run("completion for builtins function arguments", func(t *testing.T) {
		// Should complete function names when typing builtins arguments
		completions := cp.GetCompletions("(builtins ma", 12)

		// Should include 'map' builtin function
		found := false
		for _, comp := range completions {
			if comp == "map" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'map' in completions for '(builtins ma'")
		}
	})

	t.Run("completion for builtins with plus sign", func(t *testing.T) {
		// Should complete '+' when typing it as argument to builtins
		completions := cp.GetCompletions("(builtins +", 11)

		// Should include '+'
		found := false
		for _, comp := range completions {
			if comp == "+" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected '+' in completions for '(builtins +'")
		}
	})

	t.Run("completion for builtins with filter prefix", func(t *testing.T) {
		// Should complete 'filter' when typing 'fil' as argument to builtins
		completions := cp.GetCompletions("(builtins fil", 13)

		// Should include 'filter'
		found := false
		for _, comp := range completions {
			if comp == "filter" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'filter' in completions for '(builtins fil'")
		}
	})

	t.Run("completion for builtins with spaces", func(t *testing.T) {
		// Should work with extra spaces
		completions := cp.GetCompletions("(builtins  ma", 13)

		// Should include 'map'
		found := false
		for _, comp := range completions {
			if comp == "map" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'map' in completions for '(builtins  ma' (with extra space)")
		}
	})

	t.Run("no completion for second argument to builtins", func(t *testing.T) {
		// Should not complete second argument to builtins (it only takes one)
		completions := cp.GetCompletions("(builtins map ma", 16)

		// Should be empty since builtins only takes one argument
		if len(completions) > 0 {
			t.Errorf("Expected no completions for second argument to builtins, got %v", completions)
		}
	})
}
