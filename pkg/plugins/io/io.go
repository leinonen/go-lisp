// Package io provides I/O functionality for the Lisp interpreter
package io

import (
	"fmt"
	"os"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// IOPlugin provides I/O functions
type IOPlugin struct {
	*plugins.BasePlugin
}

// NewIOPlugin creates a new I/O plugin
func NewIOPlugin() *IOPlugin {
	return &IOPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"io",
			"1.0.0",
			"I/O operations (print, file operations)",
			[]string{}, // No dependencies
		),
	}
}

// RegisterFunctions registers I/O functions
func (p *IOPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// Print function
	printFunc := functions.NewFunction(
		"print!",
		registry.CategoryIO,
		-1, // Variadic
		"Print values to stdout without newline: (print! \"hello\" \"world\") => prints \"hello world\"",
		p.print,
	)
	if err := reg.Register(printFunc); err != nil {
		return err
	}

	// Println function
	printlnFunc := functions.NewFunction(
		"println!",
		registry.CategoryIO,
		-1, // Variadic
		"Print values to stdout with newline: (println! \"hello\" \"world\") => prints \"hello world\\n\"",
		p.println,
	)
	if err := reg.Register(printlnFunc); err != nil {
		return err
	}

	// Read file function
	readFileFunc := functions.NewFunction(
		"read-file",
		registry.CategoryIO,
		1,
		"Read file contents as string: (read-file \"test.txt\") => file contents",
		p.readFile,
	)
	if err := reg.Register(readFileFunc); err != nil {
		return err
	}

	// Write file function
	writeFileFunc := functions.NewFunction(
		"write-file",
		registry.CategoryIO,
		2,
		"Write string to file: (write-file \"test.txt\" \"content\") => true",
		p.writeFile,
	)
	if err := reg.Register(writeFileFunc); err != nil {
		return err
	}

	// File exists function
	fileExistsFunc := functions.NewFunction(
		"file-exists?",
		registry.CategoryIO,
		1,
		"Check if file exists: (file-exists? \"test.txt\") => true/false",
		p.fileExists,
	)
	if err := reg.Register(fileExistsFunc); err != nil {
		return err
	}

	return nil
}

// Print! function - outputs values to stdout without newline (side effect)
func (p *IOPlugin) print(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return &types.NilValue{}, nil
	}

	// Evaluate and print all arguments
	for i, arg := range args {
		value, err := evaluator.Eval(arg)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(p.valueToString(value))
	}

	return &types.NilValue{}, nil
}

// Println! function - outputs values to stdout with newline (side effect)
func (p *IOPlugin) println(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		fmt.Println()
		return &types.NilValue{}, nil
	}

	// Evaluate and print all arguments
	for i, arg := range args {
		value, err := evaluator.Eval(arg)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(p.valueToString(value))
	}
	fmt.Println()

	return &types.NilValue{}, nil
}

// Read-file reads the contents of a file and returns it as a string
func (p *IOPlugin) readFile(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("read-file requires exactly 1 argument, got %d", len(args))
	}

	// Evaluate the filename argument
	filenameValue, err := evaluator.Eval(args[0])
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

// Write-file writes content to a file
func (p *IOPlugin) writeFile(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("write-file requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the filename argument
	filenameValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Ensure filename is a string
	filename, ok := filenameValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("write-file filename must be a string, got %T", filenameValue)
	}

	// Evaluate the content argument
	contentValue, err := evaluator.Eval(args[1])
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

// File-exists? checks if a file exists
func (p *IOPlugin) fileExists(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("file-exists? requires exactly 1 argument, got %d", len(args))
	}

	// Evaluate the filename argument
	filenameValue, err := evaluator.Eval(args[0])
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

// Helper function to convert a value to its string representation for printing
func (p *IOPlugin) valueToString(value types.Value) string {
	switch v := value.(type) {
	case *types.NilValue:
		return "nil"
	case types.StringValue:
		return string(v)
	case types.NumberValue:
		return fmt.Sprintf("%.0f", float64(v))
	case types.BooleanValue:
		if v {
			return "true"
		}
		return "false"
	case *types.ListValue:
		result := "("
		for i, elem := range v.Elements {
			if i > 0 {
				result += " "
			}
			result += p.valueToString(elem)
		}
		result += ")"
		return result
	case *types.FunctionValue:
		paramNames := make([]string, len(v.Params))
		for i, param := range v.Params {
			paramNames[i] = param
		}
		return fmt.Sprintf("#<function([%s])>", fmt.Sprintf("%v", paramNames))
	case *types.HashMapValue:
		result := "{"
		first := true
		for key, val := range v.Elements {
			if !first {
				result += ", "
			}
			result += fmt.Sprintf("%s: %s", key, p.valueToString(val))
			first = false
		}
		result += "}"
		return result
	case *types.BigNumberValue:
		return v.String()
	case *types.FutureValue:
		return v.String()
	case *types.ChannelValue:
		return v.String()
	case *types.WaitGroupValue:
		return v.String()
	case *types.AtomValue:
		return v.String()
	case types.KeywordValue:
		return ":" + string(v)
	default:
		return fmt.Sprintf("%v", value)
	}
}
