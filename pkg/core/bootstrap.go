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
	
	// Load standard library files
	stdlibFiles := []string{
		"lisp/stdlib/core.lisp",      // Re-enabled after fixing function conflicts
		"lisp/stdlib/enhanced.lisp",  // Re-enabled for testing
	}
	
	for _, filename := range stdlibFiles {
		// Look for the stdlib file
		stdlibPath := filepath.Join(cwd, filename)
		
		// Check if file exists
		if _, err := os.Stat(stdlibPath); os.IsNotExist(err) {
			// Try alternative paths
			for _, path := range []string{
				"../../" + filename,
				"../../../" + filename,
				"./" + filename,
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
			// If we can't find the file, just continue to next file
			// This allows the minimal core to work without some stdlib files
			continue
		}
		
		// Parse and evaluate the standard library
		err = loadLibraryContent(string(content), env)
		if err != nil {
			return fmt.Errorf("failed to load %s: %v", filename, err)
		}
	}
	
	return nil
}

func loadLibraryContent(content string, env *Environment) error {
	// Parse and evaluate the standard library
	lexer := NewLexer(content)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return fmt.Errorf("failed to tokenize: %v", err)
	}
	
	parser := NewParser(tokens)
	expressions, err := parser.ParseAll()
	if err != nil {
		return fmt.Errorf("failed to parse: %v", err)
	}
	
	// Evaluate each expression in the standard library
	for _, expr := range expressions {
		_, err := Eval(expr, env)
		if err != nil {
			return fmt.Errorf("failed to evaluate expression: %v", err)
		}
	}
	
	return nil
}

// CreateBootstrappedEnvironment creates a core environment with standard library loaded
func CreateBootstrappedEnvironment() (*Environment, error) {
	env := NewCoreEnvironment()
	
	err := LoadStandardLibrary(env)
	if err != nil {
		return nil, err
	}
	
	return env, nil
}