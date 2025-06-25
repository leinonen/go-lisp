package minimal

import (
	"testing"
)

// TestCoreTypesComprehensive provides comprehensive testing of core data types
func TestCoreTypesComprehensive(t *testing.T) {
	t.Run("Symbol", func(t *testing.T) {
		// Test symbol interning
		sym1 := Intern("test-symbol")
		sym2 := Intern("test-symbol")
		sym3 := Intern("different-symbol")

		if sym1 != sym2 {
			t.Error("Symbol interning failed - same string should return same symbol")
		}

		if sym1 == sym3 {
			t.Error("Different symbols should not be equal")
		}

		if sym1.String() != "test-symbol" {
			t.Errorf("Expected symbol string 'test-symbol', got '%s'", sym1.String())
		}
	})

	t.Run("Number", func(t *testing.T) {
		// Test integer representation
		intNum := Number(42)
		if intNum.String() != "42" {
			t.Errorf("Expected integer number string '42', got '%s'", intNum.String())
		}

		// Test float representation
		floatNum := Number(3.14159)
		if floatNum.String() != "3.14159" {
			t.Errorf("Expected float number string '3.14159', got '%s'", floatNum.String())
		}

		// Test zero
		zero := Number(0)
		if zero.String() != "0" {
			t.Errorf("Expected zero string '0', got '%s'", zero.String())
		}

		// Test negative numbers
		neg := Number(-5.5)
		if neg.String() != "-5.5" {
			t.Errorf("Expected negative number string '-5.5', got '%s'", neg.String())
		}
	})

	t.Run("Boolean", func(t *testing.T) {
		trueVal := Boolean(true)
		falseVal := Boolean(false)

		if trueVal.String() != "true" {
			t.Errorf("Expected true string 'true', got '%s'", trueVal.String())
		}

		if falseVal.String() != "false" {
			t.Errorf("Expected false string 'false', got '%s'", falseVal.String())
		}
	})

	t.Run("String", func(t *testing.T) {
		str := String("hello world")
		expected := `"hello world"`
		if str.String() != expected {
			t.Errorf("Expected string %s, got '%s'", expected, str.String())
		}

		// Test empty string
		empty := String("")
		if empty.String() != `""` {
			t.Errorf("Expected empty string '\"\"', got '%s'", empty.String())
		}
	})

	t.Run("Nil", func(t *testing.T) {
		nilVal := Nil{}
		if nilVal.String() != "nil" {
			t.Errorf("Expected nil string 'nil', got '%s'", nilVal.String())
		}
	})

	t.Run("List", func(t *testing.T) {
		// Test empty list
		empty := NewList()
		if !empty.IsEmpty() {
			t.Error("Empty list should report as empty")
		}
		if empty.Length() != 0 {
			t.Errorf("Empty list should have length 0, got %d", empty.Length())
		}
		if empty.String() != "()" {
			t.Errorf("Expected empty list string '()', got '%s'", empty.String())
		}

		// Test single element list
		single := NewList(Number(42))
		if single.IsEmpty() {
			t.Error("Single element list should not be empty")
		}
		if single.Length() != 1 {
			t.Errorf("Single element list should have length 1, got %d", single.Length())
		}
		if single.String() != "(42)" {
			t.Errorf("Expected single list string '(42)', got '%s'", single.String())
		}

		// Test multi-element list
		multi := NewList(Number(1), Number(2), Number(3))
		if multi.Length() != 3 {
			t.Errorf("Multi element list should have length 3, got %d", multi.Length())
		}
		if multi.String() != "(1 2 3)" {
			t.Errorf("Expected multi list string '(1 2 3)', got '%s'", multi.String())
		}

		// Test First() and Rest()
		first := multi.First()
		if num, ok := first.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected first element to be 1, got %v", first)
		}

		rest := multi.Rest()
		if rest.Length() != 2 {
			t.Errorf("Expected rest length 2, got %d", rest.Length())
		}
		if rest.String() != "(2 3)" {
			t.Errorf("Expected rest string '(2 3)', got '%s'", rest.String())
		}

		// Test Rest() on empty list
		emptyRest := empty.Rest()
		if !emptyRest.IsEmpty() {
			t.Error("Rest of empty list should be empty")
		}

		// Test First() on empty list
		emptyFirst := empty.First()
		if emptyFirst != nil {
			t.Errorf("First of empty list should be nil, got %v", emptyFirst)
		}
	})

	t.Run("Vector", func(t *testing.T) {
		// Test empty vector
		empty := NewVector()
		if !empty.IsEmpty() {
			t.Error("Empty vector should report as empty")
		}
		if empty.Length() != 0 {
			t.Errorf("Empty vector should have length 0, got %d", empty.Length())
		}
		if empty.String() != "[]" {
			t.Errorf("Expected empty vector string '[]', got '%s'", empty.String())
		}

		// Test single element vector
		single := NewVector(Number(42))
		if single.IsEmpty() {
			t.Error("Single element vector should not be empty")
		}
		if single.Length() != 1 {
			t.Errorf("Single element vector should have length 1, got %d", single.Length())
		}
		if single.String() != "[42]" {
			t.Errorf("Expected single vector string '[42]', got '%s'", single.String())
		}

		// Test multi-element vector
		multi := NewVector(Number(1), Number(2), Number(3))
		if multi.Length() != 3 {
			t.Errorf("Multi element vector should have length 3, got %d", multi.Length())
		}
		if multi.String() != "[1 2 3]" {
			t.Errorf("Expected multi vector string '[1 2 3]', got '%s'", multi.String())
		}

		// Test Get() operations
		val := multi.Get(0)
		if num, ok := val.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected Get(0) to return 1, got %v", val)
		}

		val = multi.Get(2)
		if num, ok := val.(Number); !ok || float64(num) != 3.0 {
			t.Errorf("Expected Get(2) to return 3, got %v", val)
		}

		// Test out of bounds Get()
		val = multi.Get(10)
		if _, ok := val.(Nil); !ok {
			t.Errorf("Expected Get(10) to return Nil, got %v", val)
		}

		val = multi.Get(-1)
		if _, ok := val.(Nil); !ok {
			t.Errorf("Expected Get(-1) to return Nil, got %v", val)
		}

		// Test Append()
		appended := multi.Append(Number(4))
		if appended.Length() != 4 {
			t.Errorf("Expected appended vector length 4, got %d", appended.Length())
		}
		if appended.String() != "[1 2 3 4]" {
			t.Errorf("Expected appended vector '[1 2 3 4]', got '%s'", appended.String())
		}

		// Original should be unchanged
		if multi.Length() != 3 {
			t.Error("Original vector should be unchanged after append")
		}

		// Test Update()
		updated := multi.Update(1, Number(99))
		if updated.Get(1).(Number) != 99 {
			t.Errorf("Expected updated vector[1] to be 99, got %v", updated.Get(1))
		}
		if updated.String() != "[1 99 3]" {
			t.Errorf("Expected updated vector '[1 99 3]', got '%s'", updated.String())
		}

		// Original should be unchanged
		if multi.Get(1).(Number) != 2 {
			t.Error("Original vector should be unchanged after update")
		}

		// Test out of bounds Update()
		unchanged := multi.Update(10, Number(999))
		if unchanged.String() != multi.String() {
			t.Error("Out of bounds update should return unchanged vector")
		}

		// Test ToList()
		list := multi.ToList()
		if list.Length() != 3 {
			t.Errorf("Expected ToList() length 3, got %d", list.Length())
		}
		if list.String() != "(1 2 3)" {
			t.Errorf("Expected ToList() string '(1 2 3)', got '%s'", list.String())
		}

		// Test First() and Rest()
		first := multi.First()
		if num, ok := first.(Number); !ok || float64(num) != 1.0 {
			t.Errorf("Expected first element to be 1, got %v", first)
		}

		rest := multi.Rest()
		if rest.Length() != 2 {
			t.Errorf("Expected rest length 2, got %d", rest.Length())
		}
		if rest.String() != "[2 3]" {
			t.Errorf("Expected rest string '[2 3]', got '%s'", rest.String())
		}
	})

	t.Run("HashMap", func(t *testing.T) {
		// Test empty hashmap
		empty := NewHashMap()
		if empty.Length() != 0 {
			t.Errorf("Empty hashmap should have length 0, got %d", empty.Length())
		}
		if empty.String() != "{}" {
			t.Errorf("Expected empty hashmap string '{}', got '%s'", empty.String())
		}

		// Test Get() on empty hashmap
		val := empty.Get("nonexistent")
		if _, ok := val.(Nil); !ok {
			t.Errorf("Expected Get() on empty hashmap to return Nil, got %v", val)
		}

		// Test Put() and Get()
		hm := empty.Put("name", String("Alice"))
		hm = hm.Put("age", Number(30))

		if hm.Length() != 2 {
			t.Errorf("Expected hashmap length 2, got %d", hm.Length())
		}

		name := hm.Get("name")
		if str, ok := name.(String); !ok || string(str) != "Alice" {
			t.Errorf("Expected name to be 'Alice', got %v", name)
		}

		age := hm.Get("age")
		if num, ok := age.(Number); !ok || float64(num) != 30.0 {
			t.Errorf("Expected age to be 30, got %v", age)
		}

		// Test overwriting key
		hm = hm.Put("age", Number(31))
		newAge := hm.Get("age")
		if num, ok := newAge.(Number); !ok || float64(num) != 31.0 {
			t.Errorf("Expected updated age to be 31, got %v", newAge)
		}

		// Test Keys()
		keys := hm.Keys()
		if keys.Length() != 2 {
			t.Errorf("Expected 2 keys, got %d", keys.Length())
		}

		// Test Values()
		values := hm.Values()
		if values.Length() != 2 {
			t.Errorf("Expected 2 values, got %d", values.Length())
		}

		// Original should be unchanged
		if empty.Length() != 0 {
			t.Error("Original hashmap should be unchanged")
		}
	})

	t.Run("Set", func(t *testing.T) {
		// Test empty set
		empty := NewSet()
		if empty.Length() != 0 {
			t.Errorf("Empty set should have length 0, got %d", empty.Length())
		}
		if empty.String() != "#{}" {
			t.Errorf("Expected empty set string '#{}', got '%s'", empty.String())
		}

		// Test Add() and Contains()
		set := empty.Add(Number(1))
		set = set.Add(Number(2))
		set = set.Add(Number(1)) // Duplicate

		if set.Length() != 2 {
			t.Errorf("Expected set length 2 after adding duplicates, got %d", set.Length())
		}

		if !set.Contains(Number(1)) {
			t.Error("Set should contain 1")
		}

		if !set.Contains(Number(2)) {
			t.Error("Set should contain 2")
		}

		if set.Contains(Number(3)) {
			t.Error("Set should not contain 3")
		}

		// Test Remove()
		smaller := set.Remove(Number(1))
		if smaller.Length() != 1 {
			t.Errorf("Expected set length 1 after removal, got %d", smaller.Length())
		}

		if smaller.Contains(Number(1)) {
			t.Error("Set should not contain 1 after removal")
		}

		if !smaller.Contains(Number(2)) {
			t.Error("Set should still contain 2 after removing 1")
		}

		// Original should be unchanged
		if !set.Contains(Number(1)) {
			t.Error("Original set should be unchanged after remove")
		}
	})

	t.Run("DefinedValue", func(t *testing.T) {
		defined := DefinedValue{}
		if defined.String() != "defined" {
			t.Errorf("Expected DefinedValue string 'defined', got '%s'", defined.String())
		}
	})
}
