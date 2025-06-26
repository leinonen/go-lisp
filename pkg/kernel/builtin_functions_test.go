package kernel

import (
	"testing"
)

// TestBuiltinFunctions provides comprehensive testing of all built-in functions
func TestBuiltinFunctions(t *testing.T) {
	t.Run("ArithmeticFunctions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test addition
		cases := []struct {
			expr     *List
			expected float64
		}{
			{NewList(Intern("+")), 0.0},                                  // + with no args
			{NewList(Intern("+"), Number(5)), 5.0},                       // + with one arg
			{NewList(Intern("+"), Number(1), Number(2)), 3.0},            // basic addition
			{NewList(Intern("+"), Number(1), Number(2), Number(3)), 6.0}, // multiple args
			{NewList(Intern("+"), Number(-5), Number(3)), -2.0},          // negative numbers
			{NewList(Intern("+"), Number(1.5), Number(2.5)), 4.0},        // floats
		}

		for _, tc := range cases {
			result, err := Eval(tc.expr, repl.Env)
			if err != nil {
				t.Fatalf("Error evaluating %s: %v", tc.expr.String(), err)
			}

			if num, ok := result.(Number); !ok || float64(num) != tc.expected {
				t.Errorf("Expected %f for %s, got %v", tc.expected, tc.expr.String(), result)
			}
		}

		// Test subtraction
		subCases := []struct {
			expr     *List
			expected float64
		}{
			{NewList(Intern("-"), Number(5)), -5.0},                       // unary minus
			{NewList(Intern("-"), Number(5), Number(3)), 2.0},             // binary subtraction
			{NewList(Intern("-"), Number(10), Number(3), Number(2)), 5.0}, // multiple args
			{NewList(Intern("-"), Number(-5)), 5.0},                       // unary minus on negative
		}

		for _, tc := range subCases {
			result, err := Eval(tc.expr, repl.Env)
			if err != nil {
				t.Fatalf("Error evaluating %s: %v", tc.expr.String(), err)
			}

			if num, ok := result.(Number); !ok || float64(num) != tc.expected {
				t.Errorf("Expected %f for %s, got %v", tc.expected, tc.expr.String(), result)
			}
		}

		// Test multiplication
		mulCases := []struct {
			expr     *List
			expected float64
		}{
			{NewList(Intern("*")), 1.0},                                   // * with no args
			{NewList(Intern("*"), Number(5)), 5.0},                        // * with one arg
			{NewList(Intern("*"), Number(2), Number(3)), 6.0},             // basic multiplication
			{NewList(Intern("*"), Number(2), Number(3), Number(4)), 24.0}, // multiple args
			{NewList(Intern("*"), Number(-2), Number(3)), -6.0},           // negative numbers
			{NewList(Intern("*"), Number(2.5), Number(4)), 10.0},          // floats
		}

		for _, tc := range mulCases {
			result, err := Eval(tc.expr, repl.Env)
			if err != nil {
				t.Fatalf("Error evaluating %s: %v", tc.expr.String(), err)
			}

			if num, ok := result.(Number); !ok || float64(num) != tc.expected {
				t.Errorf("Expected %f for %s, got %v", tc.expected, tc.expr.String(), result)
			}
		}

		// Test error cases
		errorCases := []*List{
			NewList(Intern("+"), String("not-a-number")),
			NewList(Intern("-"), Boolean(true)),
			NewList(Intern("*"), Nil{}),
			NewList(Intern("-")), // - with no args should error
		}

		for _, tc := range errorCases {
			_, err := Eval(tc, repl.Env)
			if err == nil {
				t.Errorf("Expected error for %s", tc.String())
			}
		}
	})

	t.Run("ComparisonFunctions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test equality
		eqCases := []struct {
			expr     *List
			expected bool
		}{
			{NewList(Intern("="), Number(5), Number(5)), true},
			{NewList(Intern("="), Number(5), Number(3)), false},
			{NewList(Intern("="), String("hello"), String("hello")), true},
			{NewList(Intern("="), String("hello"), String("world")), false},
			{NewList(Intern("="), Boolean(true), Boolean(true)), true},
			{NewList(Intern("="), Boolean(true), Boolean(false)), false},
			{NewList(Intern("="), Number(5), String("5")), false}, // different types
		}

		for _, tc := range eqCases {
			result, err := Eval(tc.expr, repl.Env)
			if err != nil {
				t.Fatalf("Error evaluating %s: %v", tc.expr.String(), err)
			}

			if boolean, ok := result.(Boolean); !ok || bool(boolean) != tc.expected {
				t.Errorf("Expected %t for %s, got %v", tc.expected, tc.expr.String(), result)
			}
		}

		// Test less than
		ltCases := []struct {
			expr     *List
			expected bool
		}{
			{NewList(Intern("<"), Number(3), Number(5)), true},
			{NewList(Intern("<"), Number(5), Number(3)), false},
			{NewList(Intern("<"), Number(5), Number(5)), false},
			{NewList(Intern("<"), Number(-2), Number(1)), true},
		}

		for _, tc := range ltCases {
			result, err := Eval(tc.expr, repl.Env)
			if err != nil {
				t.Fatalf("Error evaluating %s: %v", tc.expr.String(), err)
			}

			if boolean, ok := result.(Boolean); !ok || bool(boolean) != tc.expected {
				t.Errorf("Expected %t for %s, got %v", tc.expected, tc.expr.String(), result)
			}
		}

		// Test other comparison operators
		otherCases := []struct {
			expr     *List
			expected bool
		}{
			{NewList(Intern("<="), Number(3), Number(5)), true},
			{NewList(Intern("<="), Number(5), Number(5)), true},
			{NewList(Intern("<="), Number(5), Number(3)), false},
			{NewList(Intern(">"), Number(5), Number(3)), true},
			{NewList(Intern(">"), Number(3), Number(5)), false},
			{NewList(Intern(">="), Number(5), Number(3)), true},
			{NewList(Intern(">="), Number(5), Number(5)), true},
			{NewList(Intern(">="), Number(3), Number(5)), false},
			{NewList(Intern("!="), Number(5), Number(3)), true},
			{NewList(Intern("!="), Number(5), Number(5)), false},
		}

		for _, tc := range otherCases {
			result, err := Eval(tc.expr, repl.Env)
			if err != nil {
				t.Fatalf("Error evaluating %s: %v", tc.expr.String(), err)
			}

			if boolean, ok := result.(Boolean); !ok || bool(boolean) != tc.expected {
				t.Errorf("Expected %t for %s, got %v", tc.expected, tc.expr.String(), result)
			}
		}
	})

	t.Run("ListFunctions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test list creation
		listExpr := NewList(Intern("list"), Number(1), Number(2), Number(3))
		result, err := Eval(listExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error creating list: %v", err)
		}

		list, ok := result.(*List)
		if !ok {
			t.Fatalf("Expected List, got %T", result)
		}

		if list.Length() != 3 {
			t.Errorf("Expected list length 3, got %d", list.Length())
		}

		// Test first function
		firstExpr := NewList(Intern("first"), NewList(Intern("quote"), list))
		result, err = Eval(firstExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error getting first: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected first element 1, got %v", result)
		}

		// Test first on vector
		vec := NewVector(Number(10), Number(20))
		firstVecExpr := NewList(Intern("first"), NewList(Intern("quote"), vec))
		result, err = Eval(firstVecExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error getting first of vector: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 10.0 {
			t.Errorf("Expected first vector element 10, got %v", result)
		}

		// Test first on empty list
		emptyFirstExpr := NewList(Intern("first"), NewList(Intern("quote"), NewList()))
		result, err = Eval(emptyFirstExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error getting first of empty list: %v", err)
		}

		if _, ok := result.(Nil); !ok {
			t.Errorf("Expected nil for first of empty list, got %v", result)
		}

		// Test rest function
		restExpr := NewList(Intern("rest"), NewList(Intern("quote"), list))
		result, err = Eval(restExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error getting rest: %v", err)
		}

		restList, ok := result.(*List)
		if !ok {
			t.Fatalf("Expected List from rest, got %T", result)
		}

		if restList.Length() != 2 {
			t.Errorf("Expected rest length 2, got %d", restList.Length())
		}

		// Test cons function
		consExpr := NewList(Intern("cons"), Number(0), NewList(Intern("quote"), list))
		result, err = Eval(consExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling cons: %v", err)
		}

		consList, ok := result.(*List)
		if !ok {
			t.Fatalf("Expected List from cons, got %T", result)
		}

		if consList.Length() != 4 {
			t.Errorf("Expected cons result length 4, got %d", consList.Length())
		}

		if num, ok := consList.First().(Number); !ok || float64(num) != 0.0 {
			t.Errorf("Expected cons first element 0, got %v", consList.First())
		}
	})

	t.Run("VectorFunctions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		vec := NewVector(Number(1), Number(2), Number(3))

		// Test vector-get
		getExpr := NewList(Intern("vector-get"), NewList(Intern("quote"), vec), Number(1))
		result, err := Eval(getExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling vector-get: %v", err)
		}

		if num, ok := result.(Number); !ok || float64(num) != 2.0 {
			t.Errorf("Expected vector-get result 2, got %v", result)
		}

		// Test vector-get out of bounds
		getOOBExpr := NewList(Intern("vector-get"), NewList(Intern("quote"), vec), Number(10))
		result, err = Eval(getOOBExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling vector-get out of bounds: %v", err)
		}

		if _, ok := result.(Nil); !ok {
			t.Errorf("Expected nil for out of bounds vector-get, got %v", result)
		}

		// Test vector-append
		appendExpr := NewList(Intern("vector-append"), NewList(Intern("quote"), vec), Number(4))
		result, err = Eval(appendExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling vector-append: %v", err)
		}

		newVec, ok := result.(*Vector)
		if !ok {
			t.Fatalf("Expected Vector from vector-append, got %T", result)
		}

		if newVec.Length() != 4 {
			t.Errorf("Expected appended vector length 4, got %d", newVec.Length())
		}

		// Test vector-update
		updateExpr := NewList(Intern("vector-update"), NewList(Intern("quote"), vec), Number(1), Number(99))
		result, err = Eval(updateExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling vector-update: %v", err)
		}

		updatedVec, ok := result.(*Vector)
		if !ok {
			t.Fatalf("Expected Vector from vector-update, got %T", result)
		}

		if num, ok := updatedVec.Get(1).(Number); !ok || float64(num) != 99.0 {
			t.Errorf("Expected updated element 99, got %v", updatedVec.Get(1))
		}
	})

	t.Run("HashMapFunctions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test hash-map creation
		hmExpr := NewList(Intern("hash-map"), String("name"), String("Alice"), String("age"), Number(30))
		result, err := Eval(hmExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error creating hash-map: %v", err)
		}

		hm, ok := result.(*HashMap)
		if !ok {
			t.Fatalf("Expected HashMap, got %T", result)
		}

		if hm.Length() != 2 {
			t.Errorf("Expected hash-map length 2, got %d", hm.Length())
		}

		// Test hash-map-get
		getExpr := NewList(Intern("hash-map-get"), NewList(Intern("quote"), hm), String("name"))
		result, err = Eval(getExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling hash-map-get: %v", err)
		}

		if str, ok := result.(String); !ok || string(str) != "Alice" {
			t.Errorf("Expected hash-map-get result 'Alice', got %v", result)
		}

		// Test hash-map-get with missing key
		getMissingExpr := NewList(Intern("hash-map-get"), NewList(Intern("quote"), hm), String("missing"))
		result, err = Eval(getMissingExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling hash-map-get with missing key: %v", err)
		}

		if _, ok := result.(Nil); !ok {
			t.Errorf("Expected nil for missing key, got %v", result)
		}

		// Test hash-map-put
		putExpr := NewList(Intern("hash-map-put"), NewList(Intern("quote"), hm), String("city"), String("NYC"))
		result, err = Eval(putExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling hash-map-put: %v", err)
		}

		newHm, ok := result.(*HashMap)
		if !ok {
			t.Fatalf("Expected HashMap from hash-map-put, got %T", result)
		}

		if newHm.Length() != 3 {
			t.Errorf("Expected new hash-map length 3, got %d", newHm.Length())
		}

		// Test hash-map-keys
		keysExpr := NewList(Intern("hash-map-keys"), NewList(Intern("quote"), hm))
		result, err = Eval(keysExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling hash-map-keys: %v", err)
		}

		keys, ok := result.(*Vector)
		if !ok {
			t.Fatalf("Expected Vector from hash-map-keys, got %T", result)
		}

		if keys.Length() != 2 {
			t.Errorf("Expected 2 keys, got %d", keys.Length())
		}
	})

	t.Run("SetFunctions", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test set creation
		setExpr := NewList(Intern("set"), Number(1), Number(2), Number(1)) // duplicate 1
		result, err := Eval(setExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error creating set: %v", err)
		}

		set, ok := result.(*Set)
		if !ok {
			t.Fatalf("Expected Set, got %T", result)
		}

		if set.Length() != 2 { // duplicates should be removed
			t.Errorf("Expected set length 2, got %d", set.Length())
		}

		// Test set-add
		addExpr := NewList(Intern("set-add"), NewList(Intern("quote"), set), Number(3))
		result, err = Eval(addExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling set-add: %v", err)
		}

		newSet, ok := result.(*Set)
		if !ok {
			t.Fatalf("Expected Set from set-add, got %T", result)
		}

		if newSet.Length() != 3 {
			t.Errorf("Expected new set length 3, got %d", newSet.Length())
		}

		// Test set-contains?
		containsExpr := NewList(Intern("set-contains?"), NewList(Intern("quote"), set), Number(1))
		result, err = Eval(containsExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling set-contains?: %v", err)
		}

		if boolean, ok := result.(Boolean); !ok || !bool(boolean) {
			t.Errorf("Expected true for set-contains?, got %v", result)
		}

		// Test set-contains? with missing element
		containsMissingExpr := NewList(Intern("set-contains?"), NewList(Intern("quote"), set), Number(99))
		result, err = Eval(containsMissingExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling set-contains? with missing element: %v", err)
		}

		if boolean, ok := result.(Boolean); !ok || bool(boolean) {
			t.Errorf("Expected false for missing element, got %v", result)
		}
	})

	t.Run("PrintFunction", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test print function (mainly checking it doesn't error)
		printExpr := NewList(Intern("print"), String("Hello"), Number(42))
		result, err := Eval(printExpr, repl.Env)
		if err != nil {
			t.Fatalf("Error calling print: %v", err)
		}

		// Print should return nil
		if _, ok := result.(Nil); !ok {
			t.Errorf("Expected nil from print, got %v", result)
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		repl := NewREPL()
		Bootstrap(repl.Env)

		// Test wrong number of arguments
		errorCases := []*List{
			NewList(Intern("first")),                       // no args
			NewList(Intern("first"), Number(1), Number(2)), // too many args
			NewList(Intern("rest")),                        // no args
			NewList(Intern("vector-get"), NewVector()),     // not enough args
			NewList(Intern("hash-map-get"), NewHashMap()),  // not enough args
		}

		for _, tc := range errorCases {
			_, err := Eval(tc, repl.Env)
			if err == nil {
				t.Errorf("Expected error for %s", tc.String())
			}
		}

		// Test wrong argument types
		typeCases := []*List{
			NewList(Intern("first"), Number(42)),          // first requires collection
			NewList(Intern("rest"), String("not-a-list")), // rest requires collection
			NewList(Intern("vector-get"), String("not-a-vector"), Number(0)),
			NewList(Intern("hash-map-get"), Number(42), String("key")),
		}

		for _, tc := range typeCases {
			_, err := Eval(tc, repl.Env)
			if err == nil {
				t.Errorf("Expected type error for %s", tc.String())
			}
		}
	})
}
