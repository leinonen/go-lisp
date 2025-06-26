package kernel

import (
	"os"
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
	// (def square (fn [x] (* x x)))
	defineExpr := NewList(
		Intern("def"),
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

func TestMacroSystemAdvanced(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)

	// Test 1: Basic quasiquote
	quoteExpr := NewList(Intern("quasiquote"), NewList(Number(1), Number(2), Number(3)))
	result, err := Eval(quoteExpr, env)
	if err != nil {
		t.Fatalf("Error evaluating quasiquote: %v", err)
	}
	if list, ok := result.(*List); !ok || list.Length() != 3 {
		t.Errorf("Expected list of length 3, got %v", result)
	}

	// Test 2: Quasiquote with unquote
	innerExpr := NewList(Intern("+"), Number(2), Number(3))
	unquoteExpr := NewList(Intern("unquote"), innerExpr)
	quasiUnquoteExpr := NewList(Intern("quasiquote"),
		NewList(Number(1), unquoteExpr, Number(4)))

	result, err = Eval(quasiUnquoteExpr, env)
	if err != nil {
		t.Fatalf("Error evaluating quasiquote with unquote: %v", err)
	}

	// Should result in (1 5 4)
	if list, ok := result.(*List); ok {
		first := list.First()
		second := list.Rest().First()
		third := list.Rest().Rest().First()

		if num, ok := first.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected first element to be 1, got %v", first)
		}
		if num, ok := second.(Number); !ok || float64(num) != 5.0 {
			t.Errorf("Expected second element to be 5, got %v", second)
		}
		if num, ok := third.(Number); !ok || float64(num) != 4.0 {
			t.Errorf("Expected third element to be 4, got %v", third)
		}
	} else {
		t.Errorf("Expected list result, got %v", result)
	}

	// Test 3: Define and use a macro
	// (defmacro unless [condition body] `(if ,condition nil ,body))
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

	_, err = Eval(defmacroExpr, env)
	if err != nil {
		t.Fatalf("Error defining macro: %v", err)
	}

	// Use the macro: (unless false 42)
	macroCallExpr := NewList(Intern("unless"), Boolean(false), Number(42))
	result, err = Eval(macroCallExpr, env)
	if err != nil {
		t.Fatalf("Error calling macro: %v", err)
	}

	if num, ok := result.(Number); !ok || float64(num) != 42.0 {
		t.Errorf("Expected macro result to be 42, got %v", result)
	}

	// Test 4: Macro with true condition should return nil
	macroCallExpr2 := NewList(Intern("unless"), Boolean(true), Number(42))
	result, err = Eval(macroCallExpr2, env)
	if err != nil {
		t.Fatalf("Error calling macro with true condition: %v", err)
	}

	if _, ok := result.(Nil); !ok {
		t.Errorf("Expected macro result to be nil, got %v", result)
	}
}

func TestAdvancedDataStructures(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)

	// Test Vector operations
	vec := NewVector(Number(1), Number(2), Number(3))

	// Test vector-get
	result := vec.Get(1)
	if num, ok := result.(Number); !ok || float64(num) != 2.0 {
		t.Errorf("Expected vector-get to return 2, got %v", result)
	}

	// Test vector-append
	newVec := vec.Append(Number(4))
	if newVec.Length() != 4 {
		t.Errorf("Expected appended vector length 4, got %d", newVec.Length())
	}

	// Test vector-update
	updatedVec := vec.Update(1, Number(99))
	result = updatedVec.Get(1)
	if num, ok := result.(Number); !ok || float64(num) != 99.0 {
		t.Errorf("Expected updated vector element to be 99, got %v", result)
	}

	// Test HashMap operations
	hm := NewHashMap()
	hm = hm.Put("name", String("Alice"))
	hm = hm.Put("age", Number(30))

	// Test hash-map-get
	result = hm.Get("name")
	if str, ok := result.(String); !ok || string(str) != "Alice" {
		t.Errorf("Expected hash-map get to return 'Alice', got %v", result)
	}

	// Test hash-map-keys
	keys := hm.Keys()
	if keys.Length() != 2 {
		t.Errorf("Expected 2 keys, got %d", keys.Length())
	}

	// Test Set operations
	set := NewSet()
	set = set.Add(Number(1))
	set = set.Add(Number(2))
	set = set.Add(Number(1)) // Duplicate should be ignored

	if set.Length() != 2 {
		t.Errorf("Expected set length 2, got %d", set.Length())
	}

	if !set.Contains(Number(1)) {
		t.Error("Expected set to contain 1")
	}

	if set.Contains(Number(3)) {
		t.Error("Expected set to not contain 3")
	}
}

func TestFileLoading(t *testing.T) {
	repl := NewREPL()
	env := repl.Env
	Bootstrap(env)

	// Create a temporary test file
	testContent := `(def test-var 42)
(def test-fn (fn [x] (* x 2)))`

	err := os.WriteFile("test-temp.lisp", []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove("test-temp.lisp")

	// Test loading the file
	_, err = LoadFile("test-temp.lisp", env)
	if err != nil {
		t.Fatalf("Failed to load file: %v", err)
	}

	// Verify the definitions were loaded
	testVar, err := env.Get(Intern("test-var"))
	if err != nil {
		t.Fatalf("Failed to get test-var: %v", err)
	}

	if num, ok := testVar.(Number); !ok || float64(num) != 42.0 {
		t.Errorf("Expected test-var to be 42, got %v", testVar)
	}

	testFn, err := env.Get(Intern("test-fn"))
	if err != nil {
		t.Fatalf("Failed to get test-fn: %v", err)
	}

	if _, ok := testFn.(*UserFunction); !ok {
		t.Errorf("Expected test-fn to be a function, got %T", testFn)
	}
}
