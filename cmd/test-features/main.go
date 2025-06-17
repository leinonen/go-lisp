package main

import (
	"fmt"
	"log"

	"github.com/leinonen/lisp-interpreter/pkg/interpreter"
)

func main() {
	// Create interpreter with new features
	interp, err := interpreter.NewInterpreter()
	if err != nil {
		log.Fatalf("Failed to create interpreter: %v", err)
	}

	// Test cases for new features
	testCases := []struct {
		name  string
		input string
	}{
		// Keywords
		{"keyword literal", ":name"},
		{"keyword function", "(keyword \"test\")"},
		{"keyword predicate", "(keyword? :name)"},

		// Enhanced control flow
		{"cond simple", "(cond (> 3 2) \"yes\" :else \"no\")"},
		{"when true", "(when true 1 2 3)"},
		{"when false", "(when false 1 2 3)"},
		{"when-not false", "(when-not false 1 2 3)"},

		// Let bindings (simplified version)
		{"let simple", "(let [x 10 y 20] (+ x y))"},

		// Vectors
		{"vector creation", "(vector 1 2 3)"},
		{"vec conversion", "(vec (list 1 2 3))"},
	}

	fmt.Println("Testing new Go-Lisp features:")
	fmt.Println("=============================")

	for _, test := range testCases {
		fmt.Printf("\n%s: %s\n", test.name, test.input)
		result, err := interp.Interpret(test.input)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Result: %v\n", result)
		}
	}
}
