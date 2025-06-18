package io

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// Mock evaluator for testing
type mockEvaluator struct {
	env *evaluator.Environment
}

func newMockEvaluator() *mockEvaluator {
	return &mockEvaluator{
		env: evaluator.NewEnvironment(),
	}
}

func (me *mockEvaluator) Eval(expr types.Expr) (types.Value, error) {
	// For simple literals, convert them directly
	switch e := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(e.Value), nil
	case *types.StringExpr:
		return types.StringValue(e.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(e.Value), nil
	case *types.KeywordExpr:
		return types.KeywordValue(e.Value), nil
	case *types.SymbolExpr:
		// Try to resolve from environment
		if value, exists := me.env.Get(e.Name); exists {
			return value, nil
		}
		return nil, nil
	default:
		// For values wrapped in valueExpr, return them as-is
		if ve, ok := expr.(valueExpr); ok {
			return ve.value, nil
		}
		// For values, return them as-is
		if val, ok := expr.(types.Value); ok {
			return val, nil
		}
		return nil, nil
	}
}

func (me *mockEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	return nil, nil // Not used in these tests
}

func (me *mockEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	// For testing purposes, just call regular Eval
	// In a real implementation, this would use the bindings
	return me.Eval(expr)
}

// Helper function to wrap values as expressions
func wrapValue(value types.Value) types.Expr {
	return valueExpr{value}
}

type valueExpr struct {
	value types.Value
}

func (ve valueExpr) String() string {
	return ve.value.String()
}

func (ve valueExpr) GetPosition() types.Position {
	return types.Position{Line: 1, Column: 1}
}

// Helper function to capture stdout output
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestIOPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewIOPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"print", "println", "read-file", "write-file", "file-exists?"}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestIOPlugin_Print(t *testing.T) {
	plugin := NewIOPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name           string
		args           []types.Expr
		expectedOutput string
	}{
		{
			name:           "no arguments",
			args:           []types.Expr{},
			expectedOutput: "",
		},
		{
			name:           "single string",
			args:           []types.Expr{&types.StringExpr{Value: "hello"}},
			expectedOutput: "hello",
		},
		{
			name: "multiple arguments",
			args: []types.Expr{
				&types.StringExpr{Value: "hello"},
				&types.StringExpr{Value: "world"},
			},
			expectedOutput: "hello world",
		},
		{
			name: "mixed types",
			args: []types.Expr{
				&types.StringExpr{Value: "The answer is"},
				&types.NumberExpr{Value: 42},
			},
			expectedOutput: "The answer is 42",
		},
		{
			name: "boolean values",
			args: []types.Expr{
				&types.StringExpr{Value: "Flag:"},
				&types.BooleanExpr{Value: true},
			},
			expectedOutput: "Flag: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				result, err := plugin.print(evaluator, tt.args)
				if err != nil {
					t.Fatalf("print failed: %v", err)
				}

				// Should return nil value
				if _, ok := result.(*types.NilValue); !ok {
					t.Errorf("Expected NilValue, got %T", result)
				}
			})

			if output != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestIOPlugin_Println(t *testing.T) {
	plugin := NewIOPlugin()
	evaluator := newMockEvaluator()

	tests := []struct {
		name           string
		args           []types.Expr
		expectedOutput string
	}{
		{
			name:           "no arguments",
			args:           []types.Expr{},
			expectedOutput: "\n",
		},
		{
			name:           "single string",
			args:           []types.Expr{&types.StringExpr{Value: "hello"}},
			expectedOutput: "hello\n",
		},
		{
			name: "multiple arguments",
			args: []types.Expr{
				&types.StringExpr{Value: "hello"},
				&types.StringExpr{Value: "world"},
			},
			expectedOutput: "hello world\n",
		},
		{
			name: "mixed types",
			args: []types.Expr{
				&types.StringExpr{Value: "The answer is"},
				&types.NumberExpr{Value: 42},
			},
			expectedOutput: "The answer is 42\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				result, err := plugin.println(evaluator, tt.args)
				if err != nil {
					t.Fatalf("println failed: %v", err)
				}

				// Should return nil value
				if _, ok := result.(*types.NilValue); !ok {
					t.Errorf("Expected NilValue, got %T", result)
				}
			})

			if output != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestIOPlugin_ReadFile(t *testing.T) {
	plugin := NewIOPlugin()
	evaluator := newMockEvaluator()

	// Create a temporary test file
	tempFile := "test_read_file.txt"
	testContent := "Hello, World!\nThis is a test file."
	err := os.WriteFile(tempFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile)

	tests := []struct {
		name        string
		args        []types.Expr
		expected    string
		expectError bool
	}{
		{
			name:     "read existing file",
			args:     []types.Expr{&types.StringExpr{Value: tempFile}},
			expected: testContent,
		},
		{
			name:        "read non-existent file",
			args:        []types.Expr{&types.StringExpr{Value: "non_existent_file.txt"}},
			expectError: true,
		},
		{
			name:        "non-string filename",
			args:        []types.Expr{&types.NumberExpr{Value: 42}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.readFile(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("readFile failed: %v", err)
			}

			stringResult, ok := result.(types.StringValue)
			if !ok {
				t.Fatalf("Expected StringValue, got %T", result)
			}

			if string(stringResult) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(stringResult))
			}
		})
	}
}

func TestIOPlugin_WriteFile(t *testing.T) {
	plugin := NewIOPlugin()
	evaluator := newMockEvaluator()

	tempFile := "test_write_file.txt"
	defer os.Remove(tempFile)

	tests := []struct {
		name        string
		args        []types.Expr
		content     string
		expectError bool
	}{
		{
			name: "write string to file",
			args: []types.Expr{
				&types.StringExpr{Value: tempFile},
				&types.StringExpr{Value: "Hello, World!"},
			},
			content: "Hello, World!",
		},
		{
			name: "write empty string to file",
			args: []types.Expr{
				&types.StringExpr{Value: tempFile},
				&types.StringExpr{Value: ""},
			},
			content: "",
		},
		{
			name: "non-string filename",
			args: []types.Expr{
				&types.NumberExpr{Value: 42},
				&types.StringExpr{Value: "content"},
			},
			expectError: true,
		},
		{
			name: "non-string content",
			args: []types.Expr{
				&types.StringExpr{Value: tempFile},
				&types.NumberExpr{Value: 42},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.writeFile(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("writeFile failed: %v", err)
			}

			// Should return true on success
			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}

			if !bool(boolResult) {
				t.Errorf("Expected true, got false")
			}

			// Verify file contents
			if !tt.expectError {
				fileContent, err := os.ReadFile(tempFile)
				if err != nil {
					t.Fatalf("Failed to read written file: %v", err)
				}

				if string(fileContent) != tt.content {
					t.Errorf("Expected file content %q, got %q", tt.content, string(fileContent))
				}
			}
		})
	}
}

func TestIOPlugin_FileExists(t *testing.T) {
	plugin := NewIOPlugin()
	evaluator := newMockEvaluator()

	// Create a temporary test file
	tempFile := "test_file_exists.txt"
	err := os.WriteFile(tempFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile)

	tests := []struct {
		name        string
		args        []types.Expr
		expected    bool
		expectError bool
	}{
		{
			name:     "existing file",
			args:     []types.Expr{&types.StringExpr{Value: tempFile}},
			expected: true,
		},
		{
			name:     "non-existent file",
			args:     []types.Expr{&types.StringExpr{Value: "non_existent_file.txt"}},
			expected: false,
		},
		{
			name:        "non-string filename",
			args:        []types.Expr{&types.NumberExpr{Value: 42}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.fileExists(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("fileExists failed: %v", err)
			}

			boolResult, ok := result.(types.BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(boolResult))
			}
		})
	}
}

func TestIOPlugin_ValueToString(t *testing.T) {
	plugin := NewIOPlugin()

	tests := []struct {
		name     string
		value    types.Value
		expected string
	}{
		{
			name:     "nil value",
			value:    &types.NilValue{},
			expected: "nil",
		},
		{
			name:     "string value",
			value:    types.StringValue("hello"),
			expected: "hello",
		},
		{
			name:     "number value (integer)",
			value:    types.NumberValue(42),
			expected: "42",
		},
		{
			name:     "number value (decimal)",
			value:    types.NumberValue(3.14),
			expected: "3",
		},
		{
			name:     "boolean true",
			value:    types.BooleanValue(true),
			expected: "true",
		},
		{
			name:     "boolean false",
			value:    types.BooleanValue(false),
			expected: "false",
		},
		{
			name: "list value",
			value: &types.ListValue{Elements: []types.Value{
				types.NumberValue(1),
				types.NumberValue(2),
				types.NumberValue(3),
			}},
			expected: "(1 2 3)",
		},
		{
			name:     "empty list",
			value:    &types.ListValue{Elements: []types.Value{}},
			expected: "()",
		},
		{
			name:     "keyword value",
			value:    types.KeywordValue("test"),
			expected: ":test",
		},
		{
			name: "function value",
			value: &types.FunctionValue{
				Params: []string{"x", "y"},
				Body:   &types.NumberExpr{Value: 1},
				Env:    nil,
			},
			expected: "#<function([[x y]])>",
		},
		{
			name: "hashmap value",
			value: &types.HashMapValue{Elements: map[string]types.Value{
				"key1": types.StringValue("value1"),
				"key2": types.NumberValue(42),
			}},
			// Note: order may vary due to map iteration
			expected: "{", // We'll check if it contains the expected parts
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := plugin.valueToString(tt.value)

			// Special case for hashmap due to non-deterministic order
			if tt.name == "hashmap value" {
				if !strings.HasPrefix(result, "{") || !strings.HasSuffix(result, "}") {
					t.Errorf("Expected hashmap format, got %s", result)
				}
				if !strings.Contains(result, "key1: value1") {
					t.Errorf("Expected hashmap to contain 'key1: value1', got %s", result)
				}
				if !strings.Contains(result, "key2: 42") {
					t.Errorf("Expected hashmap to contain 'key2: 42', got %s", result)
				}
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIOPlugin_ErrorHandling(t *testing.T) {
	plugin := NewIOPlugin()
	evaluator := newMockEvaluator()

	// Test read-file with wrong number of arguments
	t.Run("read-file wrong arity", func(t *testing.T) {
		err := plugin.RegisterFunctions(registry.NewRegistry())
		if err != nil {
			t.Fatalf("Failed to register functions: %v", err)
		}
		// The arity checking is handled by the function framework,
		// so we test the actual function implementation
		_, err = plugin.readFile(evaluator, []types.Expr{})
		if err == nil {
			t.Error("Expected error for wrong number of arguments")
		}
	})

	// Test write-file with wrong number of arguments
	t.Run("write-file wrong arity", func(t *testing.T) {
		_, err := plugin.writeFile(evaluator, []types.Expr{&types.StringExpr{Value: "test"}})
		if err == nil {
			t.Error("Expected error for wrong number of arguments")
		}
	})

	// Test file-exists with wrong number of arguments
	t.Run("file-exists wrong arity", func(t *testing.T) {
		_, err := plugin.fileExists(evaluator, []types.Expr{})
		if err == nil {
			t.Error("Expected error for wrong number of arguments")
		}
	})
}

func TestIOPlugin_PluginInfo(t *testing.T) {
	plugin := NewIOPlugin()

	if plugin.Name() != "io" {
		t.Errorf("Expected plugin name 'io', got %s", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected plugin version '1.0.0', got %s", plugin.Version())
	}

	if plugin.Description() != "I/O operations (print, file operations)" {
		t.Errorf("Expected specific description, got %s", plugin.Description())
	}

	deps := plugin.Dependencies()
	if len(deps) != 0 {
		t.Errorf("Expected no dependencies, got %v", deps)
	}
}
