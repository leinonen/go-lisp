// Package evaluator_files contains file system operation functionality
package evaluator

import (
	"fmt"
	"os"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// File operations

// evalReadFile reads the contents of a file and returns it as a string
func (e *Evaluator) evalReadFile(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("read-file requires exactly 1 argument, got %d", len(args))
	}

	// Evaluate the filename argument
	filenameValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Ensure it's a string
	filename, ok := filenameValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("read-file filename must be a string, got %T", filenameValue)
	}

	// Read the file
	content, err := os.ReadFile(string(filename))
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	return types.StringValue(content), nil
}

// evalWriteFile writes content to a file
func (e *Evaluator) evalWriteFile(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("write-file requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the filename argument
	filenameValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Ensure filename is a string
	filename, ok := filenameValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("write-file filename must be a string, got %T", filenameValue)
	}

	// Evaluate the content argument
	contentValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Ensure content is a string
	content, ok := contentValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("write-file content must be a string, got %T", contentValue)
	}

	// Write the file with 0644 permissions
	err = os.WriteFile(string(filename), []byte(content), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write file %s: %v", filename, err)
	}

	// Return true on success
	return types.BooleanValue(true), nil
}

// evalFileExists checks if a file exists
func (e *Evaluator) evalFileExists(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("file-exists? requires exactly 1 argument, got %d", len(args))
	}

	// Evaluate the filename argument
	filenameValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Ensure it's a string
	filename, ok := filenameValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("file-exists? filename must be a string, got %T", filenameValue)
	}

	// Check if file exists using os.Stat
	_, err = os.Stat(string(filename))
	if err != nil {
		if os.IsNotExist(err) {
			return types.BooleanValue(false), nil
		}
		// For other errors (permissions, etc.), return false
		return types.BooleanValue(false), nil
	}

	return types.BooleanValue(true), nil
}
