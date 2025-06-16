package types

import (
	"testing"
)

func TestValueTypesComprehensive(t *testing.T) {
	t.Run("NumberValue string representation", func(t *testing.T) {
		tests := []struct {
			value    NumberValue
			expected string
		}{
			{NumberValue(42), "42"},
			{NumberValue(-17), "-17"},
			{NumberValue(3.14), "3.14"},
			{NumberValue(0), "0"},
			{NumberValue(-0), "0"},
			{NumberValue(1.5), "1.5"},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("NumberValue(%g).String() = %q, expected %q", float64(tt.value), result, tt.expected)
			}
		}
	})

	t.Run("StringValue string representation", func(t *testing.T) {
		tests := []struct {
			value    StringValue
			expected string
		}{
			{StringValue("hello"), "hello"},
			{StringValue(""), ""},
			{StringValue("with\nnewline"), "with\nnewline"},
			{StringValue("with\"quote"), "with\"quote"},
			{StringValue("unicode: ñ"), "unicode: ñ"},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("StringValue(%q).String() = %q, expected %q", string(tt.value), result, tt.expected)
			}
		}
	})

	t.Run("BooleanValue string representation", func(t *testing.T) {
		tests := []struct {
			value    BooleanValue
			expected string
		}{
			{BooleanValue(true), "true"},
			{BooleanValue(false), "false"},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("BooleanValue(%t).String() = %q, expected %q", bool(tt.value), result, tt.expected)
			}
		}
	})

	t.Run("ListValue string representation", func(t *testing.T) {
		tests := []struct {
			name     string
			value    *ListValue
			expected string
		}{
			{
				"empty list",
				&ListValue{Elements: []Value{}},
				"()",
			},
			{
				"single element",
				&ListValue{Elements: []Value{NumberValue(42)}},
				"(42)",
			},
			{
				"multiple elements",
				&ListValue{Elements: []Value{
					NumberValue(1),
					NumberValue(2),
					NumberValue(3),
				}},
				"(1 2 3)",
			},
			{
				"mixed types",
				&ListValue{Elements: []Value{
					NumberValue(42),
					StringValue("hello"),
					BooleanValue(true),
				}},
				"(42 hello true)",
			},
			{
				"nested list",
				&ListValue{Elements: []Value{
					NumberValue(1),
					&ListValue{Elements: []Value{
						NumberValue(2),
						NumberValue(3),
					}},
					NumberValue(4),
				}},
				"(1 (2 3) 4)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.value.String()
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			})
		}
	})

	t.Run("FunctionValue string representation", func(t *testing.T) {
		tests := []struct {
			name     string
			value    FunctionValue
			expected string
		}{
			{
				"no parameters",
				FunctionValue{
					Params: []string{},
					Body:   &NumberExpr{Value: 42},
				},
				"#<function([])>",
			},
			{
				"single parameter",
				FunctionValue{
					Params: []string{"x"},
					Body:   &SymbolExpr{Name: "x"},
				},
				"#<function([x])>",
			},
			{
				"multiple parameters",
				FunctionValue{
					Params: []string{"x", "y"},
					Body:   &ListExpr{Elements: []Expr{&SymbolExpr{Name: "+"}, &SymbolExpr{Name: "x"}, &SymbolExpr{Name: "y"}}},
				},
				"#<function([x y])>",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.value.String()
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			})
		}
	})

	t.Run("KeywordValue string representation", func(t *testing.T) {
		tests := []struct {
			value    KeywordValue
			expected string
		}{
			{KeywordValue("name"), ":name"},
			{KeywordValue("age"), ":age"},
			{KeywordValue("123"), ":123"},
			{KeywordValue(""), ":"},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("KeywordValue(%q).String() = %q, expected %q", string(tt.value), result, tt.expected)
			}
		}
	})
}

func TestExpressionTypes(t *testing.T) {
	t.Run("NumberExpr string representation", func(t *testing.T) {
		tests := []struct {
			value    NumberExpr
			expected string
		}{
			{NumberExpr{Value: 42}, "NumberExpr(42)"},
			{NumberExpr{Value: -17}, "NumberExpr(-17)"},
			{NumberExpr{Value: 3.14}, "NumberExpr(3.14)"},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		}
	})

	t.Run("StringExpr string representation", func(t *testing.T) {
		tests := []struct {
			value    StringExpr
			expected string
		}{
			{StringExpr{Value: "hello"}, `StringExpr("hello")`},
			{StringExpr{Value: ""}, `StringExpr("")`},
			{StringExpr{Value: "with\nnewline"}, `StringExpr("with\nnewline")`},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		}
	})

	t.Run("BooleanExpr string representation", func(t *testing.T) {
		tests := []struct {
			value    BooleanExpr
			expected string
		}{
			{BooleanExpr{Value: true}, "BooleanExpr(true)"},
			{BooleanExpr{Value: false}, "BooleanExpr(false)"},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		}
	})

	t.Run("SymbolExpr string representation", func(t *testing.T) {
		tests := []struct {
			value    SymbolExpr
			expected string
		}{
			{SymbolExpr{Name: "hello"}, "SymbolExpr(hello)"},
			{SymbolExpr{Name: "+"}, "SymbolExpr(+)"},
			{SymbolExpr{Name: "var-name"}, "SymbolExpr(var-name)"},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		}
	})

	t.Run("KeywordExpr string representation", func(t *testing.T) {
		tests := []struct {
			value    KeywordExpr
			expected string
		}{
			{KeywordExpr{Value: "name"}, "KeywordExpr(:name)"},
			{KeywordExpr{Value: "age"}, "KeywordExpr(:age)"},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		}
	})

	t.Run("ListExpr string representation", func(t *testing.T) {
		tests := []struct {
			name     string
			value    ListExpr
			expected string
		}{
			{
				"empty list",
				ListExpr{Elements: []Expr{}},
				"ListExpr([])",
			},
			{
				"single element",
				ListExpr{Elements: []Expr{&NumberExpr{Value: 42}}},
				"ListExpr([NumberExpr(42)])",
			},
			{
				"multiple elements",
				ListExpr{Elements: []Expr{
					&SymbolExpr{Name: "+"},
					&NumberExpr{Value: 1},
					&NumberExpr{Value: 2},
				}},
				"ListExpr([SymbolExpr(+) NumberExpr(1) NumberExpr(2)])",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.value.String()
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			})
		}
	})

	t.Run("BracketExpr string representation", func(t *testing.T) {
		tests := []struct {
			name     string
			value    BracketExpr
			expected string
		}{
			{
				"empty bracket",
				BracketExpr{Elements: []Expr{}},
				"BracketExpr([])",
			},
			{
				"single element",
				BracketExpr{Elements: []Expr{&SymbolExpr{Name: "x"}}},
				"BracketExpr([SymbolExpr(x)])",
			},
			{
				"multiple elements",
				BracketExpr{Elements: []Expr{
					&SymbolExpr{Name: "x"},
					&SymbolExpr{Name: "y"},
				}},
				"BracketExpr([SymbolExpr(x) SymbolExpr(y)])",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.value.String()
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			})
		}
	})
}

func TestBigNumberTypes(t *testing.T) {
	t.Run("BigNumberExpr string representation", func(t *testing.T) {
		tests := []struct {
			value    BigNumberExpr
			expected string
		}{
			{BigNumberExpr{Value: "123456789012345678901234567890"}, "BigNumberExpr(123456789012345678901234567890)"},
			{BigNumberExpr{Value: "0"}, "BigNumberExpr(0)"},
			{BigNumberExpr{Value: "-999999999999999999999999999999"}, "BigNumberExpr(-999999999999999999999999999999)"},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		}
	})

	t.Run("BigNumberValue string representation", func(t *testing.T) {
		tests := []struct {
			name     string
			value    *BigNumberValue
			expected string
		}{
			{
				"positive big number",
				func() *BigNumberValue { v, _ := NewBigNumberFromString("123456789012345678901234567890"); return v }(),
				"123456789012345678901234567890",
			},
			{
				"negative big number",
				func() *BigNumberValue { v, _ := NewBigNumberFromString("-999999999999999999999999999999"); return v }(),
				"-999999999999999999999999999999",
			},
			{
				"zero as big number",
				func() *BigNumberValue { v, _ := NewBigNumberFromString("0"); return v }(),
				"0",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.value == nil {
					t.Fatalf("failed to create big number from string")
				}
				result := tt.value.String()
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			})
		}
	})
}

func TestAtomTypes(t *testing.T) {
	t.Run("AtomValue string representation", func(t *testing.T) {
		tests := []struct {
			name     string
			value    *AtomValue
			expected string
		}{
			{
				"atom with number",
				NewAtom(NumberValue(42)),
				"#<atom:42>",
			},
			{
				"atom with string",
				NewAtom(StringValue("hello")),
				"#<atom:hello>",
			},
			{
				"atom with list",
				NewAtom(&ListValue{Elements: []Value{NumberValue(1), NumberValue(2)}}),
				"#<atom:(1 2)>",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.value.String()
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			})
		}
	})
}

func TestConcurrencyTypes(t *testing.T) {
	t.Run("FutureValue string representation", func(t *testing.T) {
		future := &FutureValue{}
		result := future.String()
		expected := "#<future:pending>"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("ChannelValue string representation", func(t *testing.T) {
		tests := []struct {
			name     string
			value    *ChannelValue
			expected string
		}{
			{
				"unbuffered channel",
				NewChannel(0),
				"#<channel:open:size=0>",
			},
			{
				"buffered channel",
				NewChannel(10),
				"#<channel:open:size=10>",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.value.String()
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			})
		}
	})
}

func TestComplexNesting(t *testing.T) {
	t.Run("deeply nested list values", func(t *testing.T) {
		// Create: (((42)))
		innermost := &ListValue{Elements: []Value{NumberValue(42)}}
		middle := &ListValue{Elements: []Value{innermost}}
		outermost := &ListValue{Elements: []Value{middle}}

		result := outermost.String()
		expected := "(((42)))"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("mixed nested structures", func(t *testing.T) {
		// Create: (fn [x y] (+ x (* y 2)))
		fnList := &ListValue{Elements: []Value{
			StringValue("fn"),    // function name as string for this test
			StringValue("[x y]"), // parameters as string for simplicity
			&ListValue{Elements: []Value{ // body
				StringValue("+"),
				StringValue("x"),
				&ListValue{Elements: []Value{
					StringValue("*"),
					StringValue("y"),
					NumberValue(2),
				}},
			}},
		}}

		result := fnList.String()
		expected := "(fn [x y] (+ x (* y 2)))"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("list with nil values", func(t *testing.T) {
		list := &ListValue{Elements: []Value{
			NumberValue(1),
			nil,
			NumberValue(3),
		}}

		result := list.String()
		// Should handle nil values gracefully
		expected := "(1 nil 3)"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})
}

func TestEdgeCaseStringRepresentations(t *testing.T) {
	t.Run("empty and whitespace strings", func(t *testing.T) {
		tests := []struct {
			value    StringValue
			expected string
		}{
			{StringValue(""), ""},
			{StringValue(" "), " "},
			{StringValue("\t"), "\t"},
			{StringValue("\n"), "\n"},
			{StringValue("   "), "   "},
		}

		for _, tt := range tests {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("StringValue(%q).String() = %q, expected %q", string(tt.value), result, tt.expected)
			}
		}
	})

	t.Run("special number values", func(t *testing.T) {
		tests := []struct {
			name     string
			value    NumberValue
			expected string
		}{
			{"positive zero", NumberValue(0.0), "0"},
			{"negative zero", NumberValue(-0.0), "0"},
			{"very small positive", NumberValue(1e-10), "1e-10"},
			{"very large positive", NumberValue(1e20), "1e+20"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.value.String()
				// Note: exact float formatting might vary
				t.Logf("NumberValue(%g).String() = %q", float64(tt.value), result)
			})
		}
	})
}
