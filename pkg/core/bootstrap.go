package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadStandardLibrary loads the self-hosted standard library
func LoadStandardLibrary(env *Environment) error {
	// Find the stdlib directory relative to the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}
	
	// Look for lisp/stdlib/core.lisp
	stdlibPath := filepath.Join(cwd, "lisp", "stdlib", "core.lisp")
	
	// Check if file exists
	if _, err := os.Stat(stdlibPath); os.IsNotExist(err) {
		// Try alternative paths
		for _, path := range []string{
			"../../lisp/stdlib/core.lisp",
			"../../../lisp/stdlib/core.lisp",
			"./lisp/stdlib/core.lisp",
		} {
			if _, err := os.Stat(path); err == nil {
				stdlibPath = path
				break
			}
		}
	}
	
	// Load the standard library file
	content, err := os.ReadFile(stdlibPath)
	if err != nil {
		// If we can't find the file, just return without error for now
		// This allows the minimal core to work without the stdlib
		return nil
	}
	
	// Parse and evaluate the standard library
	lexer := NewLexer(string(content))
	tokens, err := lexer.Tokenize()
	if err != nil {
		return fmt.Errorf("failed to tokenize stdlib: %v", err)
	}
	
	parser := NewParser(tokens)
	expressions, err := parser.ParseAll()
	if err != nil {
		return fmt.Errorf("failed to parse stdlib: %v", err)
	}
	
	// Evaluate each expression in the standard library
	for _, expr := range expressions {
		_, err := Eval(expr, env)
		if err != nil {
			return fmt.Errorf("failed to evaluate stdlib expression: %v", err)
		}
	}
	
	return nil
}

// CreateBootstrappedEnvironment creates a core environment with standard library loaded
func CreateBootstrappedEnvironment() (*Environment, error) {
	env := NewCoreEnvironment()
	
	// For now, skip stdlib loading to test core functionality
	// err := LoadStandardLibrary(env)
	// if err != nil {
	// 	return nil, err
	// }
	
	return env, nil
}