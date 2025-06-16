package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestEnvironment(t *testing.T) {
	t.Run("new environment creation", func(t *testing.T) {
		env := NewEnvironment()
		if env == nil {
			t.Fatal("NewEnvironment() returned nil")
		}
	})

	t.Run("set and get values", func(t *testing.T) {
		env := NewEnvironment()

		// Test setting and getting a number
		env.Set("x", types.NumberValue(42))
		val, exists := env.Get("x")
		if !exists {
			t.Error("expected variable 'x' to exist")
		}
		if !valuesEqual(val, types.NumberValue(42)) {
			t.Errorf("expected 42, got %v", val)
		}

		// Test setting and getting a string
		env.Set("message", types.StringValue("hello"))
		val, exists = env.Get("message")
		if !exists {
			t.Error("expected variable 'message' to exist")
		}
		if !valuesEqual(val, types.StringValue("hello")) {
			t.Errorf("expected 'hello', got %v", val)
		}
	})

	t.Run("overwrite existing values", func(t *testing.T) {
		env := NewEnvironment()

		// Set initial value
		env.Set("x", types.NumberValue(42))

		// Overwrite with different value
		env.Set("x", types.StringValue("hello"))

		val, exists := env.Get("x")
		if !exists {
			t.Error("expected variable 'x' to exist")
		}
		if !valuesEqual(val, types.StringValue("hello")) {
			t.Errorf("expected 'hello', got %v", val)
		}
	})

	t.Run("get non-existent variable", func(t *testing.T) {
		env := NewEnvironment()

		val, exists := env.Get("nonexistent")
		if exists {
			t.Error("expected variable 'nonexistent' to not exist")
		}
		if val != nil {
			t.Errorf("expected nil value for non-existent variable, got %v", val)
		}
	})

	t.Run("environment with parent", func(t *testing.T) {
		parent := NewEnvironment()
		parent.Set("parent_var", types.NumberValue(100))

		child := NewEnvironmentWithParent(parent)

		// Child should find parent variables
		val, exists := child.Get("parent_var")
		if !exists {
			t.Error("expected child to find parent variable")
		}
		if !valuesEqual(val, types.NumberValue(100)) {
			t.Errorf("expected 100, got %v", val)
		}

		// Child can override parent variables
		child.Set("parent_var", types.NumberValue(200))
		val, exists = child.Get("parent_var")
		if !exists {
			t.Error("expected child variable to exist")
		}
		if !valuesEqual(val, types.NumberValue(200)) {
			t.Errorf("expected 200, got %v", val)
		}

		// Parent should still have original value
		val, exists = parent.Get("parent_var")
		if !exists {
			t.Error("expected parent variable to still exist")
		}
		if !valuesEqual(val, types.NumberValue(100)) {
			t.Errorf("expected parent to still have 100, got %v", val)
		}
	})

	t.Run("nested environments", func(t *testing.T) {
		grandparent := NewEnvironment()
		grandparent.Set("level", types.NumberValue(1))

		parent := NewEnvironmentWithParent(grandparent)
		parent.Set("level", types.NumberValue(2))

		child := NewEnvironmentWithParent(parent)
		child.Set("level", types.NumberValue(3))

		// Child should see its own value
		val, exists := child.Get("level")
		if !exists {
			t.Error("expected child variable to exist")
		}
		if !valuesEqual(val, types.NumberValue(3)) {
			t.Errorf("expected child to have level 3, got %v", val)
		}

		// Test variable only in grandparent
		grandparent.Set("grandparent_only", types.StringValue("gp"))
		val, exists = child.Get("grandparent_only")
		if !exists {
			t.Error("expected child to find grandparent variable")
		}
		if !valuesEqual(val, types.StringValue("gp")) {
			t.Errorf("expected 'gp', got %v", val)
		}
	})

	t.Run("functions in environment", func(t *testing.T) {
		env := NewEnvironment()

		fn := types.FunctionValue{
			Params: []string{"x"},
			Body:   &types.NumberExpr{Value: 42},
			Env:    env,
		}

		env.Set("test_fn", fn)

		val, exists := env.Get("test_fn")
		if !exists {
			t.Error("expected function to exist in environment")
		}

		retrievedFn, ok := val.(types.FunctionValue)
		if !ok {
			t.Errorf("expected FunctionValue, got %T", val)
		}

		if len(retrievedFn.Params) != 1 || retrievedFn.Params[0] != "x" {
			t.Errorf("function parameters not preserved correctly")
		}
	})

	t.Run("complex values in environment", func(t *testing.T) {
		env := NewEnvironment()

		// Test list values
		list := &types.ListValue{
			Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			},
		}
		env.Set("my_list", list)

		val, exists := env.Get("my_list")
		if !exists {
			t.Error("expected list to exist in environment")
		}

		retrievedList, ok := val.(*types.ListValue)
		if !ok {
			t.Errorf("expected *ListValue, got %T", val)
		}

		if len(retrievedList.Elements) != 3 {
			t.Errorf("expected list with 3 elements, got %d", len(retrievedList.Elements))
		}
	})
}

func TestEnvironmentEdgeCases(t *testing.T) {
	t.Run("empty string key", func(t *testing.T) {
		env := NewEnvironment()
		env.Set("", types.NumberValue(42))

		val, exists := env.Get("")
		if !exists {
			t.Error("expected empty string key to work")
		}
		if !valuesEqual(val, types.NumberValue(42)) {
			t.Errorf("expected 42, got %v", val)
		}
	})

	t.Run("special character keys", func(t *testing.T) {
		env := NewEnvironment()

		specialKeys := []string{
			"var-with-dashes",
			"var_with_underscores",
			"var123",
			"!@#$%",
			"unicode-ñ",
		}

		for i, key := range specialKeys {
			val := types.NumberValue(float64(i))
			env.Set(key, val)

			retrieved, exists := env.Get(key)
			if !exists {
				t.Errorf("expected key '%s' to exist", key)
			}
			if !valuesEqual(retrieved, val) {
				t.Errorf("value mismatch for key '%s': expected %v, got %v", key, val, retrieved)
			}
		}
	})

	t.Run("nil values", func(t *testing.T) {
		env := NewEnvironment()

		// Setting a nil value should work
		env.Set("nil_var", nil)

		val, exists := env.Get("nil_var")
		if !exists {
			t.Error("expected nil variable to exist")
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("case sensitivity", func(t *testing.T) {
		env := NewEnvironment()

		env.Set("Variable", types.NumberValue(1))
		env.Set("variable", types.NumberValue(2))
		env.Set("VARIABLE", types.NumberValue(3))

		// All should be different
		val1, exists1 := env.Get("Variable")
		val2, exists2 := env.Get("variable")
		val3, exists3 := env.Get("VARIABLE")

		if !exists1 || !exists2 || !exists3 {
			t.Error("expected all case variations to exist")
		}

		if valuesEqual(val1, val2) || valuesEqual(val2, val3) || valuesEqual(val1, val3) {
			t.Error("expected case-sensitive keys to have different values")
		}
	})
}
