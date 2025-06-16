package parser

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestParserEdgeCases(t *testing.T) {
	t.Run("deeply nested expressions", func(t *testing.T) {
		// Create deeply nested parentheses: ((((1))))
		tokens := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.LPAREN, Value: "("},
			{Type: types.LPAREN, Value: "("},
			{Type: types.LPAREN, Value: "("},
			{Type: types.NUMBER, Value: "1"},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RPAREN, Value: ")"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should be nested list expressions
		outerList, ok := result.(*types.ListExpr)
		if !ok {
			t.Fatalf("expected ListExpr, got %T", result)
		}

		if len(outerList.Elements) != 1 {
			t.Errorf("expected 1 element in outer list, got %d", len(outerList.Elements))
		}
	})

	t.Run("mixed bracket and paren nesting", func(t *testing.T) {
		// (fn [x] (+ x [1 2]))
		tokens := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.SYMBOL, Value: "fn"},
			{Type: types.LBRACKET, Value: "["},
			{Type: types.SYMBOL, Value: "x"},
			{Type: types.RBRACKET, Value: "]"},
			{Type: types.LPAREN, Value: "("},
			{Type: types.SYMBOL, Value: "+"},
			{Type: types.SYMBOL, Value: "x"},
			{Type: types.LBRACKET, Value: "["},
			{Type: types.NUMBER, Value: "1"},
			{Type: types.NUMBER, Value: "2"},
			{Type: types.RBRACKET, Value: "]"},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.RPAREN, Value: ")"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		list, ok := result.(*types.ListExpr)
		if !ok {
			t.Fatalf("expected ListExpr, got %T", result)
		}

		if len(list.Elements) != 3 {
			t.Errorf("expected 3 elements, got %d", len(list.Elements))
		}

		// Check that second element is a bracket expression
		_, ok = list.Elements[1].(*types.BracketExpr)
		if !ok {
			t.Errorf("expected second element to be BracketExpr, got %T", list.Elements[1])
		}
	})

	t.Run("empty nested structures", func(t *testing.T) {
		// ([] () [[]])
		tokens := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.LBRACKET, Value: "["},
			{Type: types.RBRACKET, Value: "]"},
			{Type: types.LPAREN, Value: "("},
			{Type: types.RPAREN, Value: ")"},
			{Type: types.LBRACKET, Value: "["},
			{Type: types.LBRACKET, Value: "["},
			{Type: types.RBRACKET, Value: "]"},
			{Type: types.RBRACKET, Value: "]"},
			{Type: types.RPAREN, Value: ")"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		list, ok := result.(*types.ListExpr)
		if !ok {
			t.Fatalf("expected ListExpr, got %T", result)
		}

		if len(list.Elements) != 3 {
			t.Errorf("expected 3 elements, got %d", len(list.Elements))
		}

		// First element should be empty bracket
		bracket1, ok := list.Elements[0].(*types.BracketExpr)
		if !ok {
			t.Errorf("expected first element to be BracketExpr, got %T", list.Elements[0])
		} else if len(bracket1.Elements) != 0 {
			t.Errorf("expected empty bracket, got %d elements", len(bracket1.Elements))
		}

		// Second element should be empty list
		list2, ok := list.Elements[1].(*types.ListExpr)
		if !ok {
			t.Errorf("expected second element to be ListExpr, got %T", list.Elements[1])
		} else if len(list2.Elements) != 0 {
			t.Errorf("expected empty list, got %d elements", len(list2.Elements))
		}
	})

	t.Run("large numbers", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.NUMBER, Value: "123456789012345678901234567890"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should handle large numbers (might be parsed as BigNumberExpr or NumberExpr)
		switch result.(type) {
		case *types.NumberExpr, *types.BigNumberExpr:
			// Either is acceptable
		default:
			t.Errorf("expected NumberExpr or BigNumberExpr, got %T", result)
		}
	})

	t.Run("negative numbers in various contexts", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.NUMBER, Value: "-42"},
			{Type: types.NUMBER, Value: "-3.14"},
			{Type: types.NUMBER, Value: "-0"},
			{Type: types.RPAREN, Value: ")"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		list, ok := result.(*types.ListExpr)
		if !ok {
			t.Fatalf("expected ListExpr, got %T", result)
		}

		if len(list.Elements) != 3 {
			t.Errorf("expected 3 elements, got %d", len(list.Elements))
		}

		// All should be number expressions
		for i, elem := range list.Elements {
			if _, ok := elem.(*types.NumberExpr); !ok {
				t.Errorf("expected element %d to be NumberExpr, got %T", i, elem)
			}
		}
	})

	t.Run("special symbols", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.SYMBOL, Value: "+"},
			{Type: types.SYMBOL, Value: "-"},
			{Type: types.SYMBOL, Value: "*"},
			{Type: types.SYMBOL, Value: "/"},
			{Type: types.SYMBOL, Value: "="},
			{Type: types.SYMBOL, Value: "<"},
			{Type: types.SYMBOL, Value: ">"},
			{Type: types.SYMBOL, Value: "<="},
			{Type: types.SYMBOL, Value: ">="},
			{Type: types.SYMBOL, Value: "!="},
			{Type: types.RPAREN, Value: ")"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		list, ok := result.(*types.ListExpr)
		if !ok {
			t.Fatalf("expected ListExpr, got %T", result)
		}

		expectedSymbols := []string{"+", "-", "*", "/", "=", "<", ">", "<=", ">=", "!="}
		if len(list.Elements) != len(expectedSymbols) {
			t.Errorf("expected %d elements, got %d", len(expectedSymbols), len(list.Elements))
		}

		for i, elem := range list.Elements {
			symbol, ok := elem.(*types.SymbolExpr)
			if !ok {
				t.Errorf("expected element %d to be SymbolExpr, got %T", i, elem)
				continue
			}
			if i < len(expectedSymbols) && symbol.Name != expectedSymbols[i] {
				t.Errorf("expected symbol %s, got %s", expectedSymbols[i], symbol.Name)
			}
		}
	})

	t.Run("strings with escape sequences", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.STRING, Value: "hello\nworld"},
			{Type: types.STRING, Value: "tab\there"},
			{Type: types.STRING, Value: "quote\"inside"},
			{Type: types.STRING, Value: ""}, // empty string
			{Type: types.RPAREN, Value: ")"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		list, ok := result.(*types.ListExpr)
		if !ok {
			t.Fatalf("expected ListExpr, got %T", result)
		}

		if len(list.Elements) != 4 {
			t.Errorf("expected 4 elements, got %d", len(list.Elements))
		}

		for i, elem := range list.Elements {
			if _, ok := elem.(*types.StringExpr); !ok {
				t.Errorf("expected element %d to be StringExpr, got %T", i, elem)
			}
		}
	})

	t.Run("keywords", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.KEYWORD, Value: "name"},
			{Type: types.KEYWORD, Value: "age"},
			{Type: types.KEYWORD, Value: "123"}, // numeric keyword
			{Type: types.RPAREN, Value: ")"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		list, ok := result.(*types.ListExpr)
		if !ok {
			t.Fatalf("expected ListExpr, got %T", result)
		}

		if len(list.Elements) != 3 {
			t.Errorf("expected 3 elements, got %d", len(list.Elements))
		}

		for i, elem := range list.Elements {
			if _, ok := elem.(*types.KeywordExpr); !ok {
				t.Errorf("expected element %d to be KeywordExpr, got %T", i, elem)
			}
		}
	})
}

func TestParserErrorRecovery(t *testing.T) {
	t.Run("mismatched bracket types", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.LPAREN, Value: "("},
			{Type: types.NUMBER, Value: "1"},
			{Type: types.RBRACKET, Value: "]"}, // wrong closing bracket
		}

		parser := NewParser(tokens)
		_, err := parser.Parse()
		if err == nil {
			t.Error("expected error for mismatched bracket types")
		}
	})

	t.Run("unmatched opening bracket", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.LBRACKET, Value: "["},
			{Type: types.NUMBER, Value: "1"},
			{Type: types.NUMBER, Value: "2"},
		}

		parser := NewParser(tokens)
		_, err := parser.Parse()
		if err == nil {
			t.Error("expected error for unmatched opening bracket")
		}
	})

	t.Run("unexpected closing bracket", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.NUMBER, Value: "1"},
			{Type: types.RBRACKET, Value: "]"},
		}

		parser := NewParser(tokens)
		_, err := parser.Parse()
		if err == nil {
			t.Error("expected error for unexpected closing bracket")
		}
	})

	t.Run("multiple expressions without wrapping", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.NUMBER, Value: "1"},
			{Type: types.NUMBER, Value: "2"},
		}

		parser := NewParser(tokens)
		_, err := parser.Parse()
		if err == nil {
			t.Error("expected error for multiple top-level expressions")
		}
	})
}

func TestParserFloatingPointNumbers(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected float64
	}{
		{"simple float", "3.14", 3.14},
		{"zero float", "0.0", 0.0},
		{"negative float", "-2.5", -2.5},
		{"scientific notation", "1.23e10", 1.23e10},
		{"negative scientific", "-4.56e-3", -4.56e-3},
		{"integer as float", "42.0", 42.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := []types.Token{
				{Type: types.NUMBER, Value: tt.token},
			}

			parser := NewParser(tokens)
			result, err := parser.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			numExpr, ok := result.(*types.NumberExpr)
			if !ok {
				t.Fatalf("expected NumberExpr, got %T", result)
			}

			if numExpr.Value != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, numExpr.Value)
			}
		})
	}
}

func TestParserQuotedExpressions(t *testing.T) {
	t.Run("quoted symbol", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.QUOTE, Value: "'"},
			{Type: types.SYMBOL, Value: "symbol"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should create a quote expression
		listExpr, ok := result.(*types.ListExpr)
		if !ok {
			t.Fatalf("expected ListExpr for quote, got %T", result)
		}

		if len(listExpr.Elements) != 2 {
			t.Errorf("expected 2 elements in quote expression, got %d", len(listExpr.Elements))
		}

		// First element should be quote symbol
		quoteSymbol, ok := listExpr.Elements[0].(*types.SymbolExpr)
		if !ok {
			t.Errorf("expected first element to be SymbolExpr, got %T", listExpr.Elements[0])
		} else if quoteSymbol.Name != "quote" {
			t.Errorf("expected quote symbol, got %s", quoteSymbol.Name)
		}
	})

	t.Run("quoted list", func(t *testing.T) {
		tokens := []types.Token{
			{Type: types.QUOTE, Value: "'"},
			{Type: types.LPAREN, Value: "("},
			{Type: types.NUMBER, Value: "1"},
			{Type: types.NUMBER, Value: "2"},
			{Type: types.RPAREN, Value: ")"},
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should create a quote expression containing the list
		listExpr, ok := result.(*types.ListExpr)
		if !ok {
			t.Fatalf("expected ListExpr for quote, got %T", result)
		}

		if len(listExpr.Elements) != 2 {
			t.Errorf("expected 2 elements in quote expression, got %d", len(listExpr.Elements))
		}

		// Second element should be the quoted list
		quotedList, ok := listExpr.Elements[1].(*types.ListExpr)
		if !ok {
			t.Errorf("expected second element to be ListExpr, got %T", listExpr.Elements[1])
		} else if len(quotedList.Elements) != 2 {
			t.Errorf("expected quoted list to have 2 elements, got %d", len(quotedList.Elements))
		}
	})
}
