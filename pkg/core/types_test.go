package core

import (
	"testing"
)

func TestSymbol(t *testing.T) {
	sym := Symbol("test")
	if sym.String() != "test" {
		t.Errorf("Expected 'test', got '%s'", sym.String())
	}
}

func TestIntern(t *testing.T) {
	sym1 := Intern("test")
	sym2 := Intern("test")
	
	if sym1 != sym2 {
		t.Error("Intern should return the same symbol for the same string")
	}
	
	if sym1.String() != "test" {
		t.Errorf("Expected 'test', got '%s'", sym1.String())
	}
}

func TestKeyword(t *testing.T) {
	kw := Keyword("test")
	if kw.String() != ":test" {
		t.Errorf("Expected ':test', got '%s'", kw.String())
	}
}

func TestInternKeyword(t *testing.T) {
	kw1 := InternKeyword("test")
	kw2 := InternKeyword("test")
	
	if kw1 != kw2 {
		t.Error("InternKeyword should return the same keyword for the same string")
	}
	
	if kw1.String() != ":test" {
		t.Errorf("Expected ':test', got '%s'", kw1.String())
	}
}

func TestNumber(t *testing.T) {
	// Test integer
	intNum := NewNumber(int64(42))
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
	floatNum := NewNumber(3.14)
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
	str := String("hello")
	expected := "\"hello\""
	if str.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str.String())
	}
}

func TestNil(t *testing.T) {
	nil1 := Nil{}
	nil2 := Nil{}
	
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
	emptyList := NewList()
	if emptyList != nil {
		t.Error("Empty list should be nil")
	}
	
	// Test single element list
	singleList := NewList(NewNumber(int64(1)))
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
	multiList := NewList(NewNumber(int64(1)), NewNumber(int64(2)), NewNumber(int64(3)))
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
	emptyVec := NewVector()
	if emptyVec.Count() != 0 {
		t.Errorf("Expected count 0, got %d", emptyVec.Count())
	}
	if emptyVec.String() != "[]" {
		t.Errorf("Expected '[]', got '%s'", emptyVec.String())
	}
	
	// Test vector with elements
	vec := NewVector(NewNumber(int64(1)), NewNumber(int64(2)), NewNumber(int64(3)))
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

func TestEnvironment(t *testing.T) {
	// Test basic environment
	env := NewEnvironment(nil)
	
	// Test setting and getting
	sym := Intern("test")
	val := NewNumber(int64(42))
	env.Set(sym, val)
	
	result, err := env.Get(sym)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.String() != "42" {
		t.Errorf("Expected '42', got '%s'", result.String())
	}
	
	// Test undefined symbol
	undefinedSym := Intern("undefined")
	_, err = env.Get(undefinedSym)
	if err == nil {
		t.Error("Expected error for undefined symbol")
	}
	
	// Test parent environment
	parentEnv := NewEnvironment(nil)
	parentSym := Intern("parent")
	parentVal := String("parent-value")
	parentEnv.Set(parentSym, parentVal)
	
	childEnv := NewEnvironment(parentEnv)
	result, err = childEnv.Get(parentSym)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.String() != "\"parent-value\"" {
		t.Errorf("Expected '\"parent-value\"', got '%s'", result.String())
	}
	
	// Test shadowing
	childSym := Intern("test")
	childVal := String("child-value")
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
	values := []Value{
		Symbol("test"),
		Keyword("test"),
		NewNumber(int64(42)),
		String("test"),
		Nil{},
		NewList(NewNumber(int64(1))),
		NewVector(NewNumber(int64(1))),
	}
	
	for _, val := range values {
		// Just calling String() should not panic
		_ = val.String()
	}
}