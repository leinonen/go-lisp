package core_test

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/core"
)

func TestSymbol(t *testing.T) {
	sym := core.Symbol("test")
	if sym.String() != "test" {
		t.Errorf("Expected 'test', got '%s'", sym.String())
	}
}

func TestIntern(t *testing.T) {
	sym1 := core.Intern("test")
	sym2 := core.Intern("test")

	if sym1 != sym2 {
		t.Error("Intern should return the same symbol for the same string")
	}

	if sym1.String() != "test" {
		t.Errorf("Expected 'test', got '%s'", sym1.String())
	}
}

func TestKeyword(t *testing.T) {
	kw := core.Keyword("test")
	if kw.String() != ":test" {
		t.Errorf("Expected ':test', got '%s'", kw.String())
	}
}

func TestInternKeyword(t *testing.T) {
	kw1 := core.InternKeyword("test")
	kw2 := core.InternKeyword("test")

	if kw1 != kw2 {
		t.Error("InternKeyword should return the same keyword for the same string")
	}

	if kw1.String() != ":test" {
		t.Errorf("Expected ':test', got '%s'", kw1.String())
	}
}

func TestNumber(t *testing.T) {
	// Test integer
	intNum := core.NewNumber(int64(42))
	if !intNum.IsInteger() {
		t.Error("Expected integer number")
	}
	if intNum.IsFloat() {
		t.Error("Expected not float number")
	}
	if intNum.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", intNum.ToInt())
	}
	if intNum.ToFloat() != 42.0 {
		t.Errorf("Expected 42.0, got %f", intNum.ToFloat())
	}
	if intNum.String() != "42" {
		t.Errorf("Expected '42', got '%s'", intNum.String())
	}

	// Test float
	floatNum := core.NewNumber(3.14)
	if floatNum.IsInteger() {
		t.Error("Expected not integer number")
	}
	if !floatNum.IsFloat() {
		t.Error("Expected float number")
	}
	if floatNum.ToFloat() != 3.14 {
		t.Errorf("Expected 3.14, got %f", floatNum.ToFloat())
	}
	if floatNum.ToInt() != 3 {
		t.Errorf("Expected 3, got %d", floatNum.ToInt())
	}
}

func TestString(t *testing.T) {
	str := core.String("hello")
	expected := "\"hello\""
	if str.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str.String())
	}
}

func TestNil(t *testing.T) {
	nil1 := core.Nil{}
	nil2 := core.Nil{}

	if nil1.String() != "nil" {
		t.Errorf("Expected 'nil', got '%s'", nil1.String())
	}

	// Test equality
	if nil1 != nil2 {
		t.Error("Nil values should be equal")
	}
}

func TestList(t *testing.T) {
	// Test empty list
	emptyList := core.NewList()
	if emptyList != nil {
		t.Error("Empty list should be nil")
	}

	// Test single element list
	singleList := core.NewList(core.NewNumber(int64(1)))
	if singleList.IsEmpty() {
		t.Error("Single element list should not be empty")
	}
	if singleList.First().String() != "1" {
		t.Errorf("Expected '1', got '%s'", singleList.First().String())
	}
	if singleList.Rest() != nil {
		t.Error("Single element list rest should be nil")
	}

	// Test multi-element list
	multiList := core.NewList(core.NewNumber(int64(1)), core.NewNumber(int64(2)), core.NewNumber(int64(3)))
	if multiList.IsEmpty() {
		t.Error("Multi-element list should not be empty")
	}
	if multiList.First().String() != "1" {
		t.Errorf("Expected '1', got '%s'", multiList.First().String())
	}

	rest := multiList.Rest()
	if rest.First().String() != "2" {
		t.Errorf("Expected '2', got '%s'", rest.First().String())
	}

	// Test string representation
	expected := "(1 2 3)"
	if multiList.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, multiList.String())
	}
}

func TestVector(t *testing.T) {
	// Test empty vector
	emptyVec := core.NewVector()
	if emptyVec.Count() != 0 {
		t.Errorf("Expected count 0, got %d", emptyVec.Count())
	}
	if emptyVec.String() != "[]" {
		t.Errorf("Expected '[]', got '%s'", emptyVec.String())
	}

	// Test vector with elements
	vec := core.NewVector(core.NewNumber(int64(1)), core.NewNumber(int64(2)), core.NewNumber(int64(3)))
	if vec.Count() != 3 {
		t.Errorf("Expected count 3, got %d", vec.Count())
	}

	if vec.Get(0).String() != "1" {
		t.Errorf("Expected '1', got '%s'", vec.Get(0).String())
	}
	if vec.Get(1).String() != "2" {
		t.Errorf("Expected '2', got '%s'", vec.Get(1).String())
	}
	if vec.Get(2).String() != "3" {
		t.Errorf("Expected '3', got '%s'", vec.Get(2).String())
	}

	// Test out of bounds
	if vec.Get(-1).String() != "nil" {
		t.Error("Out of bounds access should return nil")
	}
	if vec.Get(10).String() != "nil" {
		t.Error("Out of bounds access should return nil")
	}

	// Test string representation
	expected := "[1 2 3]"
	if vec.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, vec.String())
	}
}

func TestHashMap(t *testing.T) {
	// Test empty hash-map
	emptyHM := core.NewHashMap()
	if emptyHM.Count() != 0 {
		t.Errorf("Expected count 0, got %d", emptyHM.Count())
	}
	if emptyHM.String() != "{}" {
		t.Errorf("Expected '{}', got '%s'", emptyHM.String())
	}

	// Test hash-map with pairs
	hm := core.NewHashMapWithPairs(
		core.InternKeyword("name"), core.String("Alice"),
		core.InternKeyword("age"), core.NewNumber(int64(30)),
	)
	if hm.Count() != 2 {
		t.Errorf("Expected count 2, got %d", hm.Count())
	}

	// Test get
	nameKey := core.InternKeyword("name")
	nameVal := hm.Get(nameKey)
	if nameVal.String() != "\"Alice\"" {
		t.Errorf("Expected '\"Alice\"', got '%s'", nameVal.String())
	}

	ageKey := core.InternKeyword("age")
	ageVal := hm.Get(ageKey)
	if ageVal.String() != "30" {
		t.Errorf("Expected '30', got '%s'", ageVal.String())
	}

	// Test get non-existent key
	nonExistentVal := hm.Get(core.InternKeyword("nonexistent"))
	if nonExistentVal.String() != "nil" {
		t.Errorf("Expected 'nil', got '%s'", nonExistentVal.String())
	}

	// Test contains key
	if !hm.ContainsKey(nameKey) {
		t.Error("Expected hash-map to contain :name key")
	}
	if hm.ContainsKey(core.InternKeyword("nonexistent")) {
		t.Error("Expected hash-map to not contain :nonexistent key")
	}

	// Test set new value
	hm.Set(core.InternKeyword("city"), core.String("New York"))
	if hm.Count() != 3 {
		t.Errorf("Expected count 3, got %d", hm.Count())
	}
	cityVal := hm.Get(core.InternKeyword("city"))
	if cityVal.String() != "\"New York\"" {
		t.Errorf("Expected '\"New York\"', got '%s'", cityVal.String())
	}
}

func TestSet(t *testing.T) {
	// Test empty set
	emptySet := core.NewSet()
	if emptySet.Count() != 0 {
		t.Errorf("Expected count 0, got %d", emptySet.Count())
	}
	if emptySet.String() != "#{}" {
		t.Errorf("Expected '#{}', got '%s'", emptySet.String())
	}

	// Test set with elements
	set := core.NewSetWithElements(
		core.NewNumber(int64(1)),
		core.NewNumber(int64(2)),
		core.NewNumber(int64(3)),
	)
	if set.Count() != 3 {
		t.Errorf("Expected count 3, got %d", set.Count())
	}

	// Test contains
	if !set.Contains(core.NewNumber(int64(1))) {
		t.Error("Expected set to contain 1")
	}
	if !set.Contains(core.NewNumber(int64(2))) {
		t.Error("Expected set to contain 2")
	}
	if !set.Contains(core.NewNumber(int64(3))) {
		t.Error("Expected set to contain 3")
	}
	if set.Contains(core.NewNumber(int64(4))) {
		t.Error("Expected set to not contain 4")
	}

	// Test add duplicate (should not increase count)
	set.Add(core.NewNumber(int64(1)))
	if set.Count() != 3 {
		t.Errorf("Expected count to remain 3, got %d", set.Count())
	}

	// Test add new element
	set.Add(core.NewNumber(int64(4)))
	if set.Count() != 4 {
		t.Errorf("Expected count 4, got %d", set.Count())
	}
	if !set.Contains(core.NewNumber(int64(4))) {
		t.Error("Expected set to contain 4 after adding")
	}

	// Test remove
	set.Remove(core.NewNumber(int64(2)))
	if set.Count() != 3 {
		t.Errorf("Expected count 3 after removal, got %d", set.Count())
	}
	if set.Contains(core.NewNumber(int64(2))) {
		t.Error("Expected set to not contain 2 after removal")
	}

	// Test remove non-existent element
	set.Remove(core.NewNumber(int64(99)))
	if set.Count() != 3 {
		t.Errorf("Expected count to remain 3, got %d", set.Count())
	}
}

func TestEnvironment(t *testing.T) {
	// Test basic environment
	env := core.NewEnvironment(nil)

	// Test setting and getting
	sym := core.Intern("test")
	val := core.NewNumber(int64(42))
	env.Set(sym, val)

	result, err := env.Get(sym)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.String() != "42" {
		t.Errorf("Expected '42', got '%s'", result.String())
	}

	// Test undefined symbol
	undefinedSym := core.Intern("undefined")
	_, err = env.Get(undefinedSym)
	if err == nil {
		t.Error("Expected error for undefined symbol")
	}

	// Test parent environment
	parentEnv := core.NewEnvironment(nil)
	parentSym := core.Intern("parent")
	parentVal := core.String("parent-value")
	parentEnv.Set(parentSym, parentVal)

	childEnv := core.NewEnvironment(parentEnv)
	result, err = childEnv.Get(parentSym)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.String() != "\"parent-value\"" {
		t.Errorf("Expected '\"parent-value\"', got '%s'", result.String())
	}

	// Test shadowing
	childSym := core.Intern("test")
	childVal := core.String("child-value")
	childEnv.Set(childSym, childVal)

	result, err = childEnv.Get(childSym)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.String() != "\"child-value\"" {
		t.Errorf("Expected '\"child-value\"', got '%s'", result.String())
	}
}

func TestValueInterface(t *testing.T) {
	// Test that all types implement Value interface
	values := []core.Value{
		core.Symbol("test"),
		core.Keyword("test"),
		core.NewNumber(int64(42)),
		core.String("test"),
		core.Nil{},
		core.NewList(core.NewNumber(int64(1))),
		core.NewVector(core.NewNumber(int64(1))),
		core.NewHashMap(),
		core.NewSet(),
	}

	for _, val := range values {
		// Just calling String() should not panic
		_ = val.String()
	}
}
