package minimal

import (
	"strings"
	"testing"
)

// TestMacroSystemComprehensive provides comprehensive testing of the macro system
func TestMacroSystemComprehensive(t *testing.T) {
	t.Run("BasicQuasiquote", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Basic quasiquote without unquote should work like quote
		quasiExpr := NewList(Intern("quasiquote"), NewList(Number(1), Number(2), Number(3)))
		result, err := Eval(quasiExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating quasiquote: %v", err)
		}

		list, ok := result.(*List)
		if !ok {
			t.Fatalf("Expected List, got %T", result)
		}

		if list.Length() != 3 {
			t.Errorf("Expected list length 3, got %d", list.Length())
		}

		// Verify elements are as expected
		first := list.First()
		if num, ok := first.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected first element 1, got %v", first)
		}
	})

	t.Run("QuasiquoteWithUnquote", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test quasiquote with unquote
		// `(1 ~(+ 2 3) 4) should become (1 5 4)
		innerExpr := NewList(Intern("+"), Number(2), Number(3))
		unquoteExpr := NewList(Intern("unquote"), innerExpr)
		quasiExpr := NewList(Intern("quasiquote"),
			NewList(Number(1), unquoteExpr, Number(4)))

		result, err := Eval(quasiExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating quasiquote with unquote: %v", err)
		}

		list, ok := result.(*List)
		if !ok {
			t.Fatalf("Expected List, got %T", result)
		}

		if list.Length() != 3 {
			t.Errorf("Expected list length 3, got %d", list.Length())
		}

		// Check elements
		elements := []float64{1.0, 5.0, 4.0}
		current := list
		for i, expected := range elements {
			if current.IsEmpty() {
				t.Fatalf("List too short at element %d", i)
			}

			if num, ok := current.First().(Number); !ok || float64(num) != expected {
				t.Errorf("Expected element %d to be %f, got %v", i, expected, current.First())
			}

			current = current.Rest()
		}
	})

	t.Run("NestedQuasiquote", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test nested structures
		// `(outer ~(+ 1 2) (inner ~(* 3 4)))
		innerAdd := NewList(Intern("+"), Number(1), Number(2))
		innerMul := NewList(Intern("*"), Number(3), Number(4))

		quasiExpr := NewList(Intern("quasiquote"),
			NewList(
				Intern("outer"),
				NewList(Intern("unquote"), innerAdd),
				NewList(
					Intern("inner"),
					NewList(Intern("unquote"), innerMul),
				),
			))

		result, err := Eval(quasiExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating nested quasiquote: %v", err)
		}

		list, ok := result.(*List)
		if !ok {
			t.Fatalf("Expected List, got %T", result)
		}

		if list.Length() != 3 {
			t.Errorf("Expected outer list length 3, got %d", list.Length())
		}

		// Check second element (should be 3)
		second := list.Rest().First()
		if num, ok := second.(Number); !ok || float64(num) != 3.0 {
			t.Errorf("Expected second element 3, got %v", second)
		}

		// Check third element (should be inner list with 12)
		third := list.Rest().Rest().First()
		innerList, ok := third.(*List)
		if !ok {
			t.Fatalf("Expected inner list, got %T", third)
		}

		if innerList.Length() != 2 {
			t.Errorf("Expected inner list length 2, got %d", innerList.Length())
		}

		innerSecond := innerList.Rest().First()
		if num, ok := innerSecond.(Number); !ok || float64(num) != 12.0 {
			t.Errorf("Expected inner second element 12, got %v", innerSecond)
		}
	})

	t.Run("UnquoteOutsideQuasiquote", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Unquote outside quasiquote should be an error
		unquoteExpr := NewList(Intern("unquote"), Number(42))
		_, err := Eval(unquoteExpr, repl.Env)
		if err == nil {
			t.Error("Expected error for unquote outside quasiquote")
		}

		if !strings.Contains(err.Error(), "unquote") {
			t.Errorf("Expected unquote error message, got: %v", err)
		}
	})

	t.Run("DefmacroBasic", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a simple macro
		// (defmacro when [condition body] `(if ~condition ~body nil))
		defmacroExpr := NewList(
			Intern("defmacro"),
			Intern("when"),
			NewVector(Intern("condition"), Intern("body")),
			NewList(
				Intern("quasiquote"),
				NewList(
					Intern("if"),
					NewList(Intern("unquote"), Intern("condition")),
					NewList(Intern("unquote"), Intern("body")),
					Nil{},
				),
			),
		)

		result, err := Eval(defmacroExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining macro: %v", err)
		}

		// Should return DefinedValue
		if _, ok := result.(DefinedValue); !ok {
			t.Errorf("Expected DefinedValue from defmacro, got %T", result)
		}

		// Verify the macro is defined
		whenMacro, err := repl.Env.Get(Intern("when"))
		if err != nil {
			t.Fatalf("Error getting defined macro: %v", err)
		}

		if _, ok := whenMacro.(*Macro); !ok {
			t.Errorf("Expected Macro, got %T", whenMacro)
		}
	})

	t.Run("MacroExpansion", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define unless macro
		// (defmacro unless [condition body] `(if ~condition nil ~body))
		defmacroExpr := NewList(
			Intern("defmacro"),
			Intern("unless"),
			NewVector(Intern("condition"), Intern("body")),
			NewList(
				Intern("quasiquote"),
				NewList(
					Intern("if"),
					NewList(Intern("unquote"), Intern("condition")),
					Nil{},
					NewList(Intern("unquote"), Intern("body")),
				),
			),
		)

		_, err := Eval(defmacroExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining unless macro: %v", err)
		}

		// Use the macro: (unless false 42)
		macroCallExpr := NewList(Intern("unless"), Boolean(false), Number(42))
		result, err := Eval(macroCallExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling unless macro: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected macro result 42, got %v", result)
		}

		// Test with true condition
		macroCallExpr2 := NewList(Intern("unless"), Boolean(true), Number(42))
		result, err = Eval(macroCallExpr2, repl.Env)
		if err != nil {
			t.Fatalf("Error calling unless macro with true condition: %v", err)
		}

		if _, ok := result.(Nil); !ok {
			t.Errorf("Expected nil for true condition, got %v", result)
		}
	})

	t.Run("MacroWithVariableCapture", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a macro that uses variables from its environment
		// (def x 100)
		// (defmacro use-x [] x)

		defineXExpr := NewList(Intern("def"), Intern("x"), Number(100))
		_, err := Eval(defineXExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining x: %v", err)
		}

		defmacroExpr := NewList(
			Intern("defmacro"),
			Intern("use-x"),
			NewVector(),
			Intern("x"),
		)

		_, err = Eval(defmacroExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining use-x macro: %v", err)
		}

		// Call the macro
		macroCallExpr := NewList(Intern("use-x"))
		result, err := Eval(macroCallExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling use-x macro: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 100.0 {
			t.Errorf("Expected macro to capture x=100, got %v", result)
		}
	})

	t.Run("MacroRecursion", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a recursive macro (should be possible but may not terminate)
		// For safety, we'll define a macro that expands to non-recursive code
		// (defmacro countdown [n] `(if (= ~n 0) "done" ~n))

		defmacroExpr := NewList(
			Intern("defmacro"),
			Intern("countdown"),
			NewVector(Intern("n")),
			NewList(
				Intern("quasiquote"),
				NewList(
					Intern("if"),
					NewList(Intern("="), NewList(Intern("unquote"), Intern("n")), Number(0)),
					String("done"),
					NewList(Intern("unquote"), Intern("n")),
				),
			),
		)

		_, err := Eval(defmacroExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining countdown macro: %v", err)
		}

		// Test with 0
		macroCall0 := NewList(Intern("countdown"), Number(0))
		result, err := Eval(macroCall0, repl.Env)
		if err != nil {
			t.Fatalf("Error calling countdown(0): %v", err)
		}

		if str, ok := result.(String); !ok || string(str) != "done" {
			t.Errorf("Expected 'done' for countdown(0), got %v", result)
		}

		// Test with non-zero
		macroCall5 := NewList(Intern("countdown"), Number(5))
		result, err = Eval(macroCall5, repl.Env)
		if err != nil {
			t.Fatalf("Error calling countdown(5): %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 5.0 {
			t.Errorf("Expected 5 for countdown(5), got %v", result)
		}
	})

	t.Run("MacroWithComplexExpansion", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define a macro that creates a more complex expansion
		// (defmacro let-one [var val body] `((fn [~var] ~body) ~val))

		defmacroExpr := NewList(
			Intern("defmacro"),
			Intern("let-one"),
			NewVector(Intern("var"), Intern("val"), Intern("body")),
			NewList(
				Intern("quasiquote"),
				NewList(
					NewList(
						Intern("fn"),
						NewVector(NewList(Intern("unquote"), Intern("var"))),
						NewList(Intern("unquote"), Intern("body")),
					),
					NewList(Intern("unquote"), Intern("val")),
				),
			),
		)

		_, err := Eval(defmacroExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining let-one macro: %v", err)
		}

		// Use the macro: (let-one x 42 (* x 2))
		macroCallExpr := NewList(
			Intern("let-one"),
			Intern("x"),
			Number(42),
			NewList(Intern("*"), Intern("x"), Number(2)),
		)

		result, err := Eval(macroCallExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling let-one macro: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 84.0 {
			t.Errorf("Expected let-one result 84, got %v", result)
		}
	})

	t.Run("MacroErrorHandling", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test defmacro with wrong number of arguments
		errorCases := []*List{
			NewList(Intern("defmacro")),                              // no args
			NewList(Intern("defmacro"), Intern("name")),              // only name
			NewList(Intern("defmacro"), Intern("name"), NewVector()), // no body
		}

		for _, tc := range errorCases {
			_, err := Eval(tc, repl.Env)
			if err == nil {
				t.Errorf("Expected error for defmacro with wrong args: %s", tc.String())
			}
		}

		// Test defmacro with non-symbol name
		nonSymbolName := NewList(
			Intern("defmacro"),
			Number(42), // not a symbol
			NewVector(),
			Number(1),
		)

		_, err := Eval(nonSymbolName, repl.Env)
		if err == nil {
			t.Error("Expected error for defmacro with non-symbol name")
		}

		// Test defmacro with invalid parameters (not a list or vector)
		nonVectorParams := NewList(
			Intern("defmacro"),
			Intern("test"),
			String("not-a-list-or-vector"), // neither list nor vector
			Number(1),
		)

		_, err = Eval(nonVectorParams, repl.Env)
		if err == nil {
			t.Error("Expected error for defmacro with non-vector parameters")
		}
	})

	t.Run("MacroHygiene", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test that macro expansions don't interfere with local bindings
		// This is a basic test - full hygiene would require gensym support

		// Define a potentially problematic macro
		// (defmacro bad-let [body] `((fn [x] ~body) 999))

		defmacroExpr := NewList(
			Intern("defmacro"),
			Intern("bad-let"),
			NewVector(Intern("body")),
			NewList(
				Intern("quasiquote"),
				NewList(
					NewList(
						Intern("fn"),
						NewVector(Intern("x")),
						NewList(Intern("unquote"), Intern("body")),
					),
					Number(999),
				),
			),
		)

		_, err := Eval(defmacroExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining bad-let macro: %v", err)
		}

		// Use the macro in a context where x is already bound
		// (def x 42)
		// (bad-let x) should return 999, not 42

		defineXExpr := NewList(Intern("def"), Intern("x"), Number(42))
		_, err = Eval(defineXExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error defining x: %v", err)
		}

		macroCallExpr := NewList(Intern("bad-let"), Intern("x"))
		result, err := Eval(macroCallExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling bad-let macro: %v", err)
		}

		// The macro should capture the x parameter, not the global x
		if num, ok := result.(Number); !ok || float64(num) != 999.0 {
			t.Errorf("Expected macro to use parameter x=999, got %v", result)
		}
	})
}
