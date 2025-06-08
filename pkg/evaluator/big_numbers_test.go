package evaluator

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Helper function to parse and evaluate a string expression
func parseAndEval(t *testing.T, input string) (types.Value, error) {
	tok := tokenizer.NewTokenizer(input)
	tokens, err := tok.TokenizeWithError()
	if err != nil {
		t.Fatalf("Tokenization error: %v", err)
	}

	p := parser.NewParser(tokens)
	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	env := NewEnvironment()
	evaluator := NewEvaluator(env)
	return evaluator.Eval(expr)
}

func TestBigNumberArithmetic(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "large number addition",
			expr:     "(+ 1000000000000000 1000000000000000)",
			expected: "2000000000000000",
		},
		{
			name:     "large number multiplication",
			expr:     "(* 123456789012345 987654321098765)",
			expected: "121932631137021071359549253925",
		},
		{
			name:     "very large number subtraction",
			expr:     "(- 9999999999999999999 1)",
			expected: "9999999999999999998",
		},
		{
			name:     "multiple operand multiplication",
			expr:     "(* 1000000 1000000 1000000)",
			expected: "1000000000000000000",
		},
		{
			name:     "chain multiplication",
			expr:     "(* 2 3 4 5 6 7 8 9 10 11 12 13 14 15)",
			expected: "1307674368000",
		},
		{
			name:     "mixed small and large numbers",
			expr:     "(+ (* 1000000000000000 1000) 123)",
			expected: "1000000000000000123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAndEval(t, tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestBigNumberComparisons(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{
			name:     "big number equality - true",
			expr:     "(= 1000000000000000000 1000000000000000000)",
			expected: true,
		},
		{
			name:     "big number equality - false",
			expr:     "(= 1000000000000000000 1000000000000000001)",
			expected: false,
		},
		{
			name:     "big number less than - true",
			expr:     "(< 1000000000000000000 1000000000000000001)",
			expected: true,
		},
		{
			name:     "big number less than - false",
			expr:     "(< 1000000000000000001 1000000000000000000)",
			expected: false,
		},
		{
			name:     "big number greater than - true",
			expr:     "(> 1000000000000000001 1000000000000000000)",
			expected: true,
		},
		{
			name:     "big number greater than - false",
			expr:     "(> 1000000000000000000 1000000000000000001)",
			expected: false,
		},
		{
			name:     "big number less than or equal - true (equal)",
			expr:     "(<= 1000000000000000000 1000000000000000000)",
			expected: true,
		},
		{
			name:     "big number less than or equal - true (less)",
			expr:     "(<= 1000000000000000000 1000000000000000001)",
			expected: true,
		},
		{
			name:     "big number greater than or equal - true (equal)",
			expr:     "(>= 1000000000000000000 1000000000000000000)",
			expected: true,
		},
		{
			name:     "big number greater than or equal - true (greater)",
			expr:     "(>= 1000000000000000001 1000000000000000000)",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAndEval(t, tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("expected BooleanValue, got %T", result)
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestBigNumberFactorial(t *testing.T) {
	// Create a shared environment and evaluator for this test
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Helper function that uses the shared environment
	evalWithEnv := func(input string) (types.Value, error) {
		tok := tokenizer.NewTokenizer(input)
		tokens, err := tok.TokenizeWithError()
		if err != nil {
			return nil, fmt.Errorf("tokenization error: %v", err)
		}

		p := parser.NewParser(tokens)
		expr, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("parse error: %v", err)
		}

		return evaluator.Eval(expr)
	}

	// First define the factorial function
	factorialDef := `(defun factorial (n acc)
		(if (= n 0)
			acc
			(factorial (- n 1) (* n acc))))`

	_, err := evalWithEnv(factorialDef)
	if err != nil {
		t.Fatalf("Failed to define factorial: %v", err)
	}

	tests := []struct {
		name     string
		n        int
		expected string
	}{
		{
			name:     "factorial of 20",
			n:        20,
			expected: "2432902008176640000",
		},
		{
			name:     "factorial of 30",
			n:        30,
			expected: "265252859812191058636308480000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := fmt.Sprintf("(factorial %d 1)", tt.n)
			result, err := evalWithEnv(expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestBigNumberTypeConversion(t *testing.T) {
	// Test creating BigNumberValue from different sources
	t.Run("from string", func(t *testing.T) {
		bigNum, ok := types.NewBigNumberFromString("123456789012345678901234567890")
		if !ok {
			t.Fatal("failed to create big number from string")
		}
		expected := "123456789012345678901234567890"
		if bigNum.String() != expected {
			t.Errorf("expected %s, got %s", expected, bigNum.String())
		}
	})

	t.Run("from int64", func(t *testing.T) {
		bigNum := types.NewBigNumberFromInt64(1000000000000000)
		expected := "1000000000000000"
		if bigNum.String() != expected {
			t.Errorf("expected %s, got %s", expected, bigNum.String())
		}
	})

	t.Run("from big.Int", func(t *testing.T) {
		bigInt := big.NewInt(1234567890)
		bigInt.Mul(bigInt, big.NewInt(1000000000)) // 1234567890000000000
		bigNum := types.NewBigNumberValue(bigInt)
		expected := "1234567890000000000"
		if bigNum.String() != expected {
			t.Errorf("expected %s, got %s", expected, bigNum.String())
		}
	})
}

func TestBigNumberMixedOperations(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "mix big and regular numbers in addition",
			expr:     "(+ 1000000000000000000 123)",
			expected: "1000000000000000123",
		},
		{
			name:     "mix big and regular numbers in multiplication",
			expr:     "(* 1000000000000000 456)",
			expected: "456000000000000000",
		},
		{
			name:     "nested expressions with big numbers",
			expr:     "(+ (* 1000000000000000 2) (* 500000000000000 3))",
			expected: "3500000000000000",
		},
		{
			name:     "subtraction resulting in big number",
			expr:     "(- 2000000000000000000 1)",
			expected: "1999999999999999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAndEval(t, tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestBigNumberOverflowDetection(t *testing.T) {
	// Test that multiplication automatically uses big integers when overflow is detected
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "multiplication that would overflow float64",
			expr:     "(* 1000000000000000 1000000000000000)",
			expected: "1000000000000000000000000000000",
		},
		{
			name:     "chain multiplication causing overflow",
			expr:     "(* 10000000000 10000000000 10000000000)",
			expected: "1000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAndEval(t, tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}
