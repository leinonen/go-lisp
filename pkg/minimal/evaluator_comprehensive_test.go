package minimal

import (
	"testing"
)

// TestEvaluatorComprehensive provides comprehensive testing of the evaluator
func TestEvaluatorComprehensive(t *testing.T) {
	t.Run("SelfEvaluatingExpressions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		cases := []struct {
			input    Value
			expected Value
		}{
			{Number(42), Number(42)},
			{Number(-3.14), Number(-3.14)},
			{String("hello"), String("hello")},
			{Boolean(true), Boolean(true)},
			{Boolean(false), Boolean(false)},
			{Nil{}, Nil{}},
		}

		for _, tc := range cases {
			result, err := Eval(tc.input, repl.Env)
			if err != nil {
				t.Fatalf("Error evaluating %v: %v", tc.input, err)
			}

			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		}
	})

	t.Run("SymbolLookup", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test built-in symbols
		sym := Intern("+")
		result, err := Eval(sym, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating +: %v", err)
		}

		if _, ok := result.(*BuiltinFunction); !ok {
			t.Errorf("Expected BuiltinFunction for +, got %T", result)
		}

		// Test undefined symbol
		undefined := Intern("undefined-symbol")
		_, err = Eval(undefined, repl.Env)
		if err == nil {
			t.Error("Expected error for undefined symbol")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		empty := NewList()
		result, err := Eval(empty, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating empty list: %v", err)
		}

		if resultList, ok := result.(*List); !ok || !resultList.IsEmpty() {
			t.Errorf("Expected empty list, got %v", result)
		}
	})

	t.Run("QuoteSpecialForm", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// (quote x) should return x unevaluated
		quoted := NewList(Intern("quote"), Intern("undefined-symbol"))
		result, err := Eval(quoted, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating quote: %v", err)
		}

		if sym, ok := result.(Symbol); !ok || string(sym) != "undefined-symbol" {
			t.Errorf("Expected symbol 'undefined-symbol', got %v", result)
		}

		// (quote (1 2 3)) should return list unevaluated
		quotedList := NewList(Intern("quote"), NewList(Number(1), Number(2), Number(3)))
		result, err = Eval(quotedList, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating quoted list: %v", err)
		}

		if list, ok := result.(*List); !ok || list.Length() != 3 {
			t.Errorf("Expected list of length 3, got %v", result)
		}
	})

	t.Run("IfSpecialForm", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// (if true 42 0)
		ifTrue := NewList(Intern("if"), Boolean(true), Number(42), Number(0))
		result, err := Eval(ifTrue, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating if true: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected 42, got %v", result)
		}

		// (if false 42 0)
		ifFalse := NewList(Intern("if"), Boolean(false), Number(42), Number(0))
		result, err = Eval(ifFalse, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating if false: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 0.0 {
			t.Errorf("Expected 0, got %v", result)
		}

		// (if nil 42 0) - nil is falsy
		ifNil := NewList(Intern("if"), Nil{}, Number(42), Number(0))
		result, err = Eval(ifNil, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating if nil: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 0.0 {
			t.Errorf("Expected 0 for nil condition, got %v", result)
		}

		// (if 1 42 0) - numbers are truthy
		ifNumber := NewList(Intern("if"), Number(1), Number(42), Number(0))
		result, err = Eval(ifNumber, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating if number: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected 42 for number condition, got %v", result)
		}

		// Test if with only two arguments (no else clause)
		ifTwoArgs := NewList(Intern("if"), Boolean(true), Number(42))
		result, err = Eval(ifTwoArgs, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating if with two args: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected 42, got %v", result)
		}

		// Test if false with no else clause should return nil
		ifFalseNoElse := NewList(Intern("if"), Boolean(false), Number(42))
		result, err = Eval(ifFalseNoElse, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating if false no else: %v", err)
		}

		if _, ok := result.(Nil); !ok {
			t.Errorf("Expected nil for false condition with no else, got %v", result)
		}
	})

	t.Run("DefSpecialForm", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// (def x 42)
		defExpr := NewList(Intern("def"), Intern("x"), Number(42))
		result, err := Eval(defExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating def: %v", err)
		}

		// Should return DefinedValue
		if _, ok := result.(DefinedValue); !ok {
			t.Errorf("Expected DefinedValue, got %T", result)
		}

		// Verify the variable is defined
		x, err := repl.Env.Get(Intern("x"))
		if err != nil {
			t.Fatalf("Error getting defined variable: %v", err)
		}

		if num, ok := x.(Number); !ok || float64(num) != 42.0 {
			t.Errorf("Expected x to be 42, got %v", x)
		}

		// Test def with expression
		defExprComp := NewList(Intern("def"), Intern("y"),
			NewList(Intern("+"), Number(10), Number(20)))
		_, err = Eval(defExprComp, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating def with expression: %v", err)
		}

		y, err := repl.Env.Get(Intern("y"))
		if err != nil {
			t.Fatalf("Error getting computed variable: %v", err)
		}

		if num, ok := y.(Number); !ok || float64(num) != 30.0 {
			t.Errorf("Expected y to be 30, got %v", y)
		}
	})

	t.Run("FnSpecialForm", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// (fn [x] (* x 2))
		fnExpr := NewList(
			Intern("fn"),
			NewVector(Intern("x")),
			NewList(Intern("*"), Intern("x"), Number(2)),
		)

		result, err := Eval(fnExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating fn: %v", err)
		}

		fn, ok := result.(*UserFunction)
		if !ok {
			t.Fatalf("Expected UserFunction, got %T", result)
		}

		// Test function parameters
		if fn.Params.Length() != 1 {
			t.Errorf("Expected 1 parameter, got %d", fn.Params.Length())
		}

		param := fn.Params.First()
		if sym, ok := param.(Symbol); !ok || string(sym) != "x" {
			t.Errorf("Expected parameter 'x', got %v", param)
		}

		// Test calling the function
		callResult, err := fn.Call([]Value{Number(5)}, repl.Env)
		if err != nil {
			t.Fatalf("Error calling function: %v", err)
		}

		if num, ok := callResult.(Number); !ok || float64(num) != 10.0 {
			t.Errorf("Expected function result 10, got %v", callResult)
		}

		// Test function with multiple parameters
		fnMulti := NewList(
			Intern("fn"),
			NewVector(Intern("x"), Intern("y")),
			NewList(Intern("+"), Intern("x"), Intern("y")),
		)

		result, err = Eval(fnMulti, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating multi-param fn: %v", err)
		}

		fnObj, ok := result.(*UserFunction)
		if !ok {
			t.Fatalf("Expected UserFunction, got %T", result)
		}

		callResult, err = fnObj.Call([]Value{Number(3), Number(7)}, repl.Env)
		if err != nil {
			t.Fatalf("Error calling multi-param function: %v", err)
		}

		if num, ok := callResult.(Number); !ok || float64(num) != 10.0 {
			t.Errorf("Expected multi-param function result 10, got %v", callResult)
		}
	})

	t.Run("DoSpecialForm", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// (do (def x 1) (def y 2) (+ x y))
		doExpr := NewList(
			Intern("do"),
			NewList(Intern("def"), Intern("x"), Number(1)),
			NewList(Intern("def"), Intern("y"), Number(2)),
			NewList(Intern("+"), Intern("x"), Intern("y")),
		)

		result, err := Eval(doExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating do: %v", err)
		}

		// Should return result of last expression
		if num, ok := result.(Number); !ok || float64(num) != 3.0 {
			t.Errorf("Expected do result 3, got %v", result)
		}

		// Verify side effects happened
		x, err := repl.Env.Get(Intern("x"))
		if err != nil {
			t.Fatalf("Error getting x after do: %v", err)
		}
		if num, ok := x.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected x to be 1, got %v", x)
		}

		y, err := repl.Env.Get(Intern("y"))
		if err != nil {
			t.Fatalf("Error getting y after do: %v", err)
		}
		if num, ok := y.(Number); !ok || float64(num) != 2.0 {
			t.Errorf("Expected y to be 2, got %v", y)
		}
	})

	t.Run("FunctionApplication", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test built-in function application
		addExpr := NewList(Intern("+"), Number(1), Number(2), Number(3))
		result, err := Eval(addExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating addition: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 6.0 {
			t.Errorf("Expected 6, got %v", result)
		}

		// Test nested function application
		nestedExpr := NewList(Intern("+"),
			NewList(Intern("*"), Number(2), Number(3)),
			NewList(Intern("-"), Number(10), Number(5)))

		result, err = Eval(nestedExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error evaluating nested expression: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 11.0 {
			t.Errorf("Expected 11, got %v", result)
		}

		// Test user function application
		// First define a function
		defFn := NewList(
			Intern("def"),
			Intern("square"),
			NewList(
				Intern("fn"),
				NewVector(Intern("x")),
				NewList(Intern("*"), Intern("x"), Intern("x")),
			),
		)

		_, err = Eval(defFn, repl.Env)
		if err != nil {
			t.Fatalf("Error defining square function: %v", err)
		}

		// Call the function
		callExpr := NewList(Intern("square"), Number(4))
		result, err = Eval(callExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling square function: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 16.0 {
			t.Errorf("Expected 16, got %v", result)
		}
	})

	t.Run("Closures", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test closure capturing outer variables
		// (def make-adder (fn [x] (fn [y] (+ x y))))
		makeAdderDef := NewList(
			Intern("def"),
			Intern("make-adder"),
			NewList(
				Intern("fn"),
				NewVector(Intern("x")),
				NewList(
					Intern("fn"),
					NewVector(Intern("y")),
					NewList(Intern("+"), Intern("x"), Intern("y")),
				),
			),
		)

		_, err := Eval(makeAdderDef, repl.Env)
		if err != nil {
			t.Fatalf("Error defining make-adder: %v", err)
		}

		// (def add5 (make-adder 5))
		add5Def := NewList(
			Intern("def"),
			Intern("add5"),
			NewList(Intern("make-adder"), Number(5)),
		)

		_, err = Eval(add5Def, repl.Env)
		if err != nil {
			t.Fatalf("Error creating add5: %v", err)
		}

		// (add5 3) should return 8
		result, err := Eval(NewList(Intern("add5"), Number(3)), repl.Env)
		if err != nil {
			t.Fatalf("Error calling add5: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 8.0 {
			t.Errorf("Expected 8 from closure, got %v", result)
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test calling non-function
		invalidCall := NewList(Number(42), Number(1))
		_, err := Eval(invalidCall, repl.Env)
		if err == nil {
			t.Error("Expected error calling number as function")
		}

		// Test wrong number of arguments
		wrongArgs := NewList(Intern("+")) // + with no arguments
		_, err = Eval(wrongArgs, repl.Env)
		// This should succeed since + with no args returns 0
		if err != nil {
			t.Fatalf("Unexpected error for + with no args: %v", err)
		}

		// Test wrong argument types
		wrongTypes := NewList(Intern("+"), Number(1), String("not-a-number"))
		_, err = Eval(wrongTypes, repl.Env)
		if err == nil {
			t.Error("Expected error for wrong argument types")
		}
	})

	t.Run("RecursiveFunction", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Define factorial function
		// (def factorial (fn [n] (if (= n 0) 1 (* n (factorial (- n 1))))))
		factorialDef := NewList(
			Intern("def"),
			Intern("factorial"),
			NewList(
				Intern("fn"),
				NewVector(Intern("n")),
				NewList(
					Intern("if"),
					NewList(Intern("="), Intern("n"), Number(0)),
					Number(1),
					NewList(
						Intern("*"),
						Intern("n"),
						NewList(
							Intern("factorial"),
							NewList(Intern("-"), Intern("n"), Number(1)),
						),
					),
				),
			),
		)

		_, err := Eval(factorialDef, repl.Env)
		if err != nil {
			t.Fatalf("Error defining factorial: %v", err)
		}

		// Test factorial(5) = 120
		result, err := Eval(NewList(Intern("factorial"), Number(5)), repl.Env)
		if err != nil {
			t.Fatalf("Error calling factorial: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 120.0 {
			t.Errorf("Expected factorial(5) = 120, got %v", result)
		}
	})
}
