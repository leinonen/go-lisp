package core

import (
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

func TestCorePlugin_RegisterFunctions(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"def", "fn", "defn", "quote", "help", "env", "plugins"}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestCorePlugin_DefFunc(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)
	evaluator := newMockEvaluator()
	evaluator.env = env

	tests := []struct {
		name        string
		args        []types.Expr
		expectedVal types.Value
		expectError bool
	}{
		{
			name:        "wrong number of arguments (0)",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (1)",
			args:        []types.Expr{&types.SymbolExpr{Name: "x"}},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (3)",
			args:        []types.Expr{&types.SymbolExpr{Name: "x"}, &types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 10}},
			expectError: true,
		},
		{
			name:        "first argument not a symbol",
			args:        []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 10}},
			expectError: true,
		},
		{
			name:        "define number variable",
			args:        []types.Expr{&types.SymbolExpr{Name: "x"}, &types.NumberExpr{Value: 42}},
			expectedVal: types.NumberValue(42),
		},
		{
			name:        "define string variable",
			args:        []types.Expr{&types.SymbolExpr{Name: "name"}, &types.StringExpr{Value: "John"}},
			expectedVal: types.StringValue("John"),
		},
		{
			name:        "define boolean variable",
			args:        []types.Expr{&types.SymbolExpr{Name: "flag"}, &types.BooleanExpr{Value: true}},
			expectedVal: types.BooleanValue(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.defFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("defFunc failed: %v", err)
			}

			if result.String() != tt.expectedVal.String() {
				t.Errorf("Expected %v, got %v", tt.expectedVal, result)
			}

			// Check that the variable was set in the environment
			if len(tt.args) >= 2 {
				if symbolExpr, ok := tt.args[0].(*types.SymbolExpr); ok {
					value, exists := env.Get(symbolExpr.Name)
					if !exists {
						t.Errorf("Variable %s was not set in environment", symbolExpr.Name)
					} else if value.String() != tt.expectedVal.String() {
						t.Errorf("Environment variable %s: expected %v, got %v", symbolExpr.Name, tt.expectedVal, value)
					}
				}
			}
		})
	}
}

func TestCorePlugin_FnFunc(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)
	evaluator := newMockEvaluator()

	tests := []struct {
		name           string
		args           []types.Expr
		expectedParams []string
		expectError    bool
	}{
		{
			name:        "wrong number of arguments (0)",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (1)",
			args:        []types.Expr{&types.BracketExpr{Elements: []types.Expr{}}},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (3)",
			args:        []types.Expr{&types.BracketExpr{Elements: []types.Expr{}}, &types.NumberExpr{Value: 5}, &types.StringExpr{Value: "extra"}},
			expectError: true,
		},
		{
			name:        "invalid parameter list type",
			args:        []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 10}},
			expectError: true,
		},
		{
			name:        "non-symbol parameter",
			args:        []types.Expr{&types.BracketExpr{Elements: []types.Expr{&types.NumberExpr{Value: 5}}}, &types.NumberExpr{Value: 10}},
			expectError: true,
		},
		{
			name:           "empty parameter list",
			args:           []types.Expr{&types.BracketExpr{Elements: []types.Expr{}}, &types.NumberExpr{Value: 42}},
			expectedParams: []string{},
		},
		{
			name: "single parameter",
			args: []types.Expr{
				&types.BracketExpr{Elements: []types.Expr{&types.SymbolExpr{Name: "x"}}},
				&types.NumberExpr{Value: 42},
			},
			expectedParams: []string{"x"},
		},
		{
			name: "multiple parameters",
			args: []types.Expr{
				&types.BracketExpr{Elements: []types.Expr{
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "y"},
					&types.SymbolExpr{Name: "z"},
				}},
				&types.NumberExpr{Value: 42},
			},
			expectedParams: []string{"x", "y", "z"},
		},
		{
			name: "function with list instead of bracket (backward compatibility)",
			args: []types.Expr{
				&types.ListExpr{Elements: []types.Expr{&types.SymbolExpr{Name: "a"}}},
				&types.StringExpr{Value: "body"},
			},
			expectedParams: []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.fnFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("fnFunc failed: %v", err)
			}

			fnValue, ok := result.(*types.FunctionValue)
			if !ok {
				t.Fatalf("Expected FunctionValue, got %T", result)
			}

			if len(fnValue.Params) != len(tt.expectedParams) {
				t.Errorf("Expected %d parameters, got %d", len(tt.expectedParams), len(fnValue.Params))
			}

			for i, expectedParam := range tt.expectedParams {
				if i >= len(fnValue.Params) {
					t.Errorf("Missing parameter %d: %s", i, expectedParam)
				} else if fnValue.Params[i] != expectedParam {
					t.Errorf("Parameter %d: expected %s, got %s", i, expectedParam, fnValue.Params[i])
				}
			}

			// Check that environment is captured
			if fnValue.Env != env {
				t.Errorf("Function environment not properly captured")
			}
		})
	}
}

func TestCorePlugin_DefnFunc(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)
	evaluator := newMockEvaluator()
	evaluator.env = env

	tests := []struct {
		name        string
		args        []types.Expr
		funcName    string
		expectError bool
	}{
		{
			name:        "wrong number of arguments (0)",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (2)",
			args:        []types.Expr{&types.SymbolExpr{Name: "add"}, &types.BracketExpr{Elements: []types.Expr{}}},
			expectError: true,
		},
		{
			name:        "first argument not a symbol",
			args:        []types.Expr{&types.NumberExpr{Value: 5}, &types.BracketExpr{Elements: []types.Expr{}}, &types.NumberExpr{Value: 42}},
			expectError: true,
		},
		{
			name: "define named function",
			args: []types.Expr{
				&types.SymbolExpr{Name: "add"},
				&types.BracketExpr{Elements: []types.Expr{
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "y"},
				}},
				&types.ListExpr{Elements: []types.Expr{
					&types.SymbolExpr{Name: "+"},
					&types.SymbolExpr{Name: "x"},
					&types.SymbolExpr{Name: "y"},
				}},
			},
			funcName: "add",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.defnFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("defnFunc failed: %v", err)
			}

			_, ok := result.(*types.FunctionValue)
			if !ok {
				t.Fatalf("Expected FunctionValue, got %T", result)
			}

			// Check that the function was set in the environment
			value, exists := env.Get(tt.funcName)
			if !exists {
				t.Errorf("Function %s was not set in environment", tt.funcName)
			} else if _, ok := value.(*types.FunctionValue); !ok {
				t.Errorf("Environment value for %s is not a function", tt.funcName)
			}
		})
	}
}

func TestCorePlugin_QuoteFunc(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)
	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expected    string
		expectError bool
	}{
		{
			name:        "wrong number of arguments (0)",
			args:        []types.Expr{},
			expectError: true,
		},
		{
			name:        "wrong number of arguments (2)",
			args:        []types.Expr{&types.NumberExpr{Value: 5}, &types.NumberExpr{Value: 10}},
			expectError: true,
		},
		{
			name:     "quote number",
			args:     []types.Expr{&types.NumberExpr{Value: 42}},
			expected: "42",
		},
		{
			name:     "quote string",
			args:     []types.Expr{&types.StringExpr{Value: "hello"}},
			expected: "hello",
		},
		{
			name:     "quote boolean",
			args:     []types.Expr{&types.BooleanExpr{Value: true}},
			expected: "true",
		},
		{
			name:     "quote keyword",
			args:     []types.Expr{&types.KeywordExpr{Value: "test"}},
			expected: ":test",
		},
		{
			name:     "quote symbol",
			args:     []types.Expr{&types.SymbolExpr{Name: "x"}},
			expected: "x",
		},
		{
			name: "quote list",
			args: []types.Expr{&types.ListExpr{Elements: []types.Expr{
				&types.NumberExpr{Value: 1},
				&types.NumberExpr{Value: 2},
				&types.NumberExpr{Value: 3},
			}}},
			expected: "(1 2 3)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.quoteFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("quoteFunc failed: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestCorePlugin_HelpFunc(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)
	reg := registry.NewRegistry()

	// Register the plugin to set up the registry
	plugin.RegisterFunctions(reg)

	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expectError bool
		checkResult func(string) bool
	}{
		{
			name: "help with no arguments",
			args: []types.Expr{},
			checkResult: func(result string) bool {
				return strings.Contains(result, "Available functions by category") &&
					strings.Contains(result, "=== core ===")
			},
		},
		{
			name: "help for specific function",
			args: []types.Expr{&types.SymbolExpr{Name: "def"}},
			checkResult: func(result string) bool {
				return strings.Contains(result, "Function: def") &&
					strings.Contains(result, "Category: core") &&
					strings.Contains(result, "Define a variable")
			},
		},
		{
			name: "help for non-existent function",
			args: []types.Expr{&types.SymbolExpr{Name: "nonexistent"}},
			checkResult: func(result string) bool {
				return strings.Contains(result, "Function not found: nonexistent")
			},
		},
		{
			name:        "help with too many arguments",
			args:        []types.Expr{&types.SymbolExpr{Name: "def"}, &types.SymbolExpr{Name: "extra"}},
			expectError: true,
		},
		{
			name:        "help with non-symbol argument",
			args:        []types.Expr{&types.NumberExpr{Value: 5}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.helpFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("helpFunc failed: %v", err)
			}

			stringResult, ok := result.(types.StringValue)
			if !ok {
				t.Fatalf("Expected StringValue, got %T", result)
			}

			if tt.checkResult != nil && !tt.checkResult(string(stringResult)) {
				t.Errorf("Result check failed. Result: %s", string(stringResult))
			}
		})
	}
}

func TestCorePlugin_EnvFunc(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)
	evaluator := newMockEvaluator()
	evaluator.env = env

	// Set up some test variables
	env.Set("x", types.NumberValue(42))
	env.Set("name", types.StringValue("test"))
	env.Set("myfunc", &types.FunctionValue{
		Params: []string{"a", "b"},
		Body:   &types.NumberExpr{Value: 1},
		Env:    env,
	})

	tests := []struct {
		name        string
		args        []types.Expr
		expectError bool
		checkResult func(string) bool
	}{
		{
			name: "env with no arguments",
			args: []types.Expr{},
			checkResult: func(result string) bool {
				return strings.Contains(result, "Environment Variables and Functions") &&
					strings.Contains(result, "x = 42") &&
					strings.Contains(result, "name = test") &&
					strings.Contains(result, "myfunc(a b)")
			},
		},
		{
			name:        "env with arguments",
			args:        []types.Expr{&types.SymbolExpr{Name: "extra"}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.envFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("envFunc failed: %v", err)
			}

			stringResult, ok := result.(types.StringValue)
			if !ok {
				t.Fatalf("Expected StringValue, got %T", result)
			}

			if tt.checkResult != nil && !tt.checkResult(string(stringResult)) {
				t.Errorf("Result check failed. Result: %s", string(stringResult))
			}
		})
	}
}

func TestCorePlugin_PluginsFunc(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)
	reg := registry.NewRegistry()

	// Register the plugin to set up the registry
	plugin.RegisterFunctions(reg)

	evaluator := newMockEvaluator()

	tests := []struct {
		name        string
		args        []types.Expr
		expectError bool
		checkResult func(string) bool
	}{
		{
			name: "plugins with no arguments",
			args: []types.Expr{},
			checkResult: func(result string) bool {
				return strings.Contains(result, "Loaded plugin categories") &&
					strings.Contains(result, "core")
			},
		},
		{
			name:        "plugins with arguments",
			args:        []types.Expr{&types.SymbolExpr{Name: "extra"}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.pluginsFunc(evaluator, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("pluginsFunc failed: %v", err)
			}

			stringResult, ok := result.(types.StringValue)
			if !ok {
				t.Fatalf("Expected StringValue, got %T", result)
			}

			if tt.checkResult != nil && !tt.checkResult(string(stringResult)) {
				t.Errorf("Result check failed. Result: %s", string(stringResult))
			}
		})
	}
}

func TestCorePlugin_ExprToValue(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)

	tests := []struct {
		name     string
		expr     types.Expr
		expected string
	}{
		{
			name:     "number expression",
			expr:     &types.NumberExpr{Value: 42},
			expected: "42",
		},
		{
			name:     "string expression",
			expr:     &types.StringExpr{Value: "hello"},
			expected: "hello",
		},
		{
			name:     "boolean expression",
			expr:     &types.BooleanExpr{Value: true},
			expected: "true",
		},
		{
			name:     "keyword expression",
			expr:     &types.KeywordExpr{Value: "test"},
			expected: ":test",
		},
		{
			name:     "symbol expression",
			expr:     &types.SymbolExpr{Name: "x"},
			expected: "x",
		},
		{
			name: "list expression",
			expr: &types.ListExpr{Elements: []types.Expr{
				&types.NumberExpr{Value: 1},
				&types.NumberExpr{Value: 2},
			}},
			expected: "(1 2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := plugin.exprToValue(tt.expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestCorePlugin_PluginInfo(t *testing.T) {
	env := evaluator.NewEnvironment()
	plugin := NewCorePlugin(env)

	if plugin.Name() != "core" {
		t.Errorf("Expected plugin name 'core', got %s", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected plugin version '1.0.0', got %s", plugin.Version())
	}

	if plugin.Description() != "Core language functionality (def, fn, quote, variables)" {
		t.Errorf("Expected specific description, got %s", plugin.Description())
	}

	deps := plugin.Dependencies()
	if len(deps) != 0 {
		t.Errorf("Expected no dependencies, got %v", deps)
	}
}
