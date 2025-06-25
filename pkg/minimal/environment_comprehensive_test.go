package minimal

import (
	"testing"
)

// TestEnvironmentComprehensive provides comprehensive testing of the environment system
func TestEnvironmentComprehensive(t *testing.T) {
	t.Run("BasicOperations", func(t *testing.T) {
		env := NewEnvironment(nil)

		// Test setting and getting a value
		sym := Intern("test-var")
		val := Number(42)
		env.Set(sym, val)

		retrieved, err := env.Get(sym)
		if err != nil {
			t.Fatalf("Error getting value: %v", err)
		}

		if retrieved != val {
			t.Errorf("Expected %v, got %v", val, retrieved)
		}
	})

	t.Run("UndefinedSymbol", func(t *testing.T) {
		env := NewEnvironment(nil)
		sym := Intern("undefined")

		_, err := env.Get(sym)
		if err == nil {
			t.Error("Expected error for undefined symbol")
		}

		if err.Error() != "undefined symbol: undefined" {
			t.Errorf("Expected undefined symbol error, got: %v", err)
		}
	})

	t.Run("ParentEnvironment", func(t *testing.T) {
		parent := NewEnvironment(nil)
		child := NewEnvironment(parent)

		// Set value in parent
		parentSym := Intern("parent-var")
		parentVal := Number(100)
		parent.Set(parentSym, parentVal)

		// Set value in child
		childSym := Intern("child-var")
		childVal := Number(200)
		child.Set(childSym, childVal)

		// Child should be able to access parent's variable
		retrieved, err := child.Get(parentSym)
		if err != nil {
			t.Fatalf("Error getting parent variable from child: %v", err)
		}
		if retrieved != parentVal {
			t.Errorf("Expected %v, got %v", parentVal, retrieved)
		}

		// Child should access its own variable
		retrieved, err = child.Get(childSym)
		if err != nil {
			t.Fatalf("Error getting child variable: %v", err)
		}
		if retrieved != childVal {
			t.Errorf("Expected %v, got %v", childVal, retrieved)
		}

		// Parent should not access child's variable
		_, err = parent.Get(childSym)
		if err == nil {
			t.Error("Parent should not access child's variable")
		}
	})

	t.Run("ShadowingVariables", func(t *testing.T) {
		parent := NewEnvironment(nil)
		child := NewEnvironment(parent)

		// Set same symbol in both environments
		sym := Intern("shadowed-var")
		parentVal := Number(100)
		childVal := Number(200)

		parent.Set(sym, parentVal)
		child.Set(sym, childVal)

		// Child should see its own value (shadowing parent)
		retrieved, err := child.Get(sym)
		if err != nil {
			t.Fatalf("Error getting shadowed variable: %v", err)
		}
		if retrieved != childVal {
			t.Errorf("Expected child value %v, got %v", childVal, retrieved)
		}

		// Parent should still see its own value
		retrieved, err = parent.Get(sym)
		if err != nil {
			t.Fatalf("Error getting parent variable: %v", err)
		}
		if retrieved != parentVal {
			t.Errorf("Expected parent value %v, got %v", parentVal, retrieved)
		}
	})

	t.Run("DefineGlobal", func(t *testing.T) {
		global := NewEnvironment(nil)
		middle := NewEnvironment(global)
		local := NewEnvironment(middle)

		// Define in local environment should go to global
		sym := Intern("global-var")
		val := Number(42)
		local.Define(sym, val)

		// All levels should see the global definition
		retrieved, err := global.Get(sym)
		if err != nil {
			t.Fatalf("Error getting global variable from global: %v", err)
		}
		if retrieved != val {
			t.Errorf("Expected %v, got %v", val, retrieved)
		}

		retrieved, err = middle.Get(sym)
		if err != nil {
			t.Fatalf("Error getting global variable from middle: %v", err)
		}
		if retrieved != val {
			t.Errorf("Expected %v, got %v", val, retrieved)
		}

		retrieved, err = local.Get(sym)
		if err != nil {
			t.Fatalf("Error getting global variable from local: %v", err)
		}
		if retrieved != val {
			t.Errorf("Expected %v, got %v", val, retrieved)
		}
	})

	t.Run("GetAllBindings", func(t *testing.T) {
		parent := NewEnvironment(nil)
		child := NewEnvironment(parent)

		// Set some variables
		parent.Set(Intern("parent1"), Number(1))
		parent.Set(Intern("parent2"), Number(2))
		parent.Set(Intern("shared"), Number(100))

		child.Set(Intern("child1"), Number(3))
		child.Set(Intern("shared"), Number(200)) // Shadows parent

		// Get all bindings from child
		all := child.GetAllBindings()

		// Should have 4 unique symbols
		if len(all) != 4 {
			t.Errorf("Expected 4 bindings, got %d", len(all))
		}

		// Check that child's shadowing takes precedence
		if val, ok := all[Intern("shared")]; !ok {
			t.Error("Expected 'shared' in all bindings")
		} else if num, ok := val.(Number); !ok || float64(num) != 200.0 {
			t.Errorf("Expected shared value 200, got %v", val)
		}

		// Check parent values are present
		if val, ok := all[Intern("parent1")]; !ok {
			t.Error("Expected 'parent1' in all bindings")
		} else if num, ok := val.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected parent1 value 1, got %v", val)
		}

		// Check child values are present
		if val, ok := all[Intern("child1")]; !ok {
			t.Error("Expected 'child1' in all bindings")
		} else if num, ok := val.(Number); !ok || float64(num) != 3.0 {
			t.Errorf("Expected child1 value 3, got %v", val)
		}
	})

	t.Run("GetLocalBindings", func(t *testing.T) {
		parent := NewEnvironment(nil)
		child := NewEnvironment(parent)

		// Set variables in both
		parent.Set(Intern("parent-var"), Number(1))
		child.Set(Intern("child-var"), Number(2))

		// Parent local bindings
		parentLocal := parent.GetLocalBindings()
		if len(parentLocal) != 1 {
			t.Errorf("Expected 1 parent local binding, got %d", len(parentLocal))
		}
		if val, ok := parentLocal[Intern("parent-var")]; !ok {
			t.Error("Expected 'parent-var' in parent local bindings")
		} else if num, ok := val.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected parent-var value 1, got %v", val)
		}

		// Child local bindings
		childLocal := child.GetLocalBindings()
		if len(childLocal) != 1 {
			t.Errorf("Expected 1 child local binding, got %d", len(childLocal))
		}
		if val, ok := childLocal[Intern("child-var")]; !ok {
			t.Error("Expected 'child-var' in child local bindings")
		} else if num, ok := val.(Number); !ok || float64(num) != 2.0 {
			t.Errorf("Expected child-var value 2, got %v", val)
		}

		// Child local should not contain parent variables
		if _, ok := childLocal[Intern("parent-var")]; ok {
			t.Error("Child local bindings should not contain parent variables")
		}
	})

	t.Run("DeepNesting", func(t *testing.T) {
		// Create a deep nesting chain
		env1 := NewEnvironment(nil)
		env2 := NewEnvironment(env1)
		env3 := NewEnvironment(env2)
		env4 := NewEnvironment(env3)
		env5 := NewEnvironment(env4)

		// Set variables at different levels
		env1.Set(Intern("level1"), Number(1))
		env3.Set(Intern("level3"), Number(3))
		env5.Set(Intern("level5"), Number(5))

		// Deepest level should access all
		val, err := env5.Get(Intern("level1"))
		if err != nil {
			t.Fatalf("Error accessing level1 from level5: %v", err)
		}
		if num, ok := val.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected level1 value 1, got %v", val)
		}

		val, err = env5.Get(Intern("level3"))
		if err != nil {
			t.Fatalf("Error accessing level3 from level5: %v", err)
		}
		if num, ok := val.(Number); !ok || float64(num) != 3.0 {
			t.Errorf("Expected level3 value 3, got %v", val)
		}

		val, err = env5.Get(Intern("level5"))
		if err != nil {
			t.Fatalf("Error accessing level5 from level5: %v", err)
		}
		if num, ok := val.(Number); !ok || float64(num) != 5.0 {
			t.Errorf("Expected level5 value 5, got %v", val)
		}

		// Middle level should not access deeper levels
		_, err = env3.Get(Intern("level5"))
		if err == nil {
			t.Error("env3 should not access env5 variables")
		}
	})

	t.Run("MultipleSymbolTypes", func(t *testing.T) {
		env := NewEnvironment(nil)

		// Test different value types
		env.Set(Intern("number"), Number(42))
		env.Set(Intern("string"), String("hello"))
		env.Set(Intern("boolean"), Boolean(true))
		env.Set(Intern("nil"), Nil{})
		env.Set(Intern("list"), NewList(Number(1), Number(2)))
		env.Set(Intern("vector"), NewVector(Number(3), Number(4)))

		// Verify all types can be retrieved correctly
		cases := []struct {
			sym      Symbol
			expected Value
		}{
			{Intern("number"), Number(42)},
			{Intern("string"), String("hello")},
			{Intern("boolean"), Boolean(true)},
			{Intern("nil"), Nil{}},
		}

		for _, tc := range cases {
			val, err := env.Get(tc.sym)
			if err != nil {
				t.Fatalf("Error getting %s: %v", tc.sym, err)
			}
			if val != tc.expected {
				t.Errorf("Expected %v for %s, got %v", tc.expected, tc.sym, val)
			}
		}

		// Test list
		val, err := env.Get(Intern("list"))
		if err != nil {
			t.Fatalf("Error getting list: %v", err)
		}
		if list, ok := val.(*List); !ok {
			t.Errorf("Expected List, got %T", val)
		} else if list.Length() != 2 {
			t.Errorf("Expected list length 2, got %d", list.Length())
		}

		// Test vector
		val, err = env.Get(Intern("vector"))
		if err != nil {
			t.Fatalf("Error getting vector: %v", err)
		}
		if vector, ok := val.(*Vector); !ok {
			t.Errorf("Expected Vector, got %T", val)
		} else if vector.Length() != 2 {
			t.Errorf("Expected vector length 2, got %d", vector.Length())
		}
	})
}
