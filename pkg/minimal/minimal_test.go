package minimal

import (
	"testing"
)

func TestMinimalKernel(t *testing.T) {
	repl := NewREPL() // This sets up builtins
	env := repl.Env
	Bootstrap(env)

	// Test basic arithmetic
	expr := NewList(Intern("+"), Number(1), Number(2))
	result, err := Eval(expr, env)
	if err != nil {
		t.Fatalf("Error evaluating (+1 2): %v", err)
	}

	if num, ok := result.(Number); !ok || float64(num) != 3.0 {
		t.Errorf("Expected 3, got %v", result)
	}

	// Test function definition and calling
	// (define square (fn [x] (* x x)))
	defineExpr := NewList(
		Intern("define"),
		Intern("square"),
		NewList(
			Intern("fn"),
			NewVector(Intern("x")),
			NewList(Intern("*"), Intern("x"), Intern("x")),
		),
	)

	_, err = Eval(defineExpr, env)
	if err != nil {
		t.Fatalf("Error defining square function: %v", err)
	}

	// (square 4)
	callExpr := NewList(Intern("square"), Number(4))
	result, err = Eval(callExpr, env)
	if err != nil {
		t.Fatalf("Error calling square function: %v", err)
	}

	if num, ok := result.(Number); !ok || float64(num) != 16.0 {
		t.Errorf("Expected 16, got %v", result)
	}

	// Test conditionals
	// (if true 42 0)
	ifExpr := NewList(Intern("if"), Boolean(true), Number(42), Number(0))
	result, err = Eval(ifExpr, env)
	if err != nil {
		t.Fatalf("Error evaluating if expression: %v", err)
	}

	if num, ok := result.(Number); !ok || float64(num) != 42.0 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestSymbolInterning(t *testing.T) {
	sym1 := Intern("test")
	sym2 := Intern("test")

	if sym1 != sym2 {
		t.Error("Symbol interning failed - same string should return same symbol")
	}
}

func TestListOperations(t *testing.T) {
	list := NewList(Number(1), Number(2), Number(3))

	if list.IsEmpty() {
		t.Error("List should not be empty")
	}

	if list.Length() != 3 {
		t.Errorf("Expected length 3, got %d", list.Length())
	}

	first := list.First()
	if num, ok := first.(Number); !ok || float64(num) != 1.0 {
		t.Errorf("Expected first element to be 1, got %v", first)
	}

	rest := list.Rest()
	if rest.Length() != 2 {
		t.Errorf("Expected rest length 2, got %d", rest.Length())
	}
}

func TestBootstrappedFunctions(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)

	// Test list function
	listExpr := NewList(Intern("list"), Number(1), Number(2), Number(3))
	result, err := Eval(listExpr, env)
	if err != nil {
		t.Fatalf("Error evaluating list function: %v", err)
	}

	if list, ok := result.(*List); !ok || list.Length() != 3 {
		t.Errorf("Expected list of length 3, got %v", result)
	}

	// Test first function
	firstExpr := NewList(Intern("first"), NewList(Intern("quote"), result))
	firstResult, err := Eval(firstExpr, env)
	if err != nil {
		t.Fatalf("Error evaluating first function: %v", err)
	}

	if num, ok := firstResult.(Number); !ok || float64(num) != 1.0 {
		t.Errorf("Expected first element to be 1, got %v", firstResult)
	}
}
