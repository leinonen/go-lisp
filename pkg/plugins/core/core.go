// Package core provides core language functionality for the Lisp interpreter
package core

import (
	"fmt"
	"sort"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// CorePlugin provides core language functionality
type CorePlugin struct {
	*plugins.BasePlugin
	env      *evaluator.Environment
	registry registry.FunctionRegistry
}

// NewCorePlugin creates a new core plugin
func NewCorePlugin(env *evaluator.Environment) *CorePlugin {
	return &CorePlugin{
		BasePlugin: plugins.NewBasePlugin(
			"core",
			"1.0.0",
			"Core language functionality (def, fn, quote, variables)",
			[]string{}, // No dependencies
		),
		env:      env,
		registry: nil, // Will be set during registration
	}
}

// RegisterFunctions registers core language functions
func (p *CorePlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// Store the registry for help functions
	p.registry = reg

	// def function
	defFunc := functions.NewFunction(
		"def",
		registry.CategoryCore,
		2,
		"Define a variable: (def name value)",
		p.defFunc,
	)
	if err := reg.Register(defFunc); err != nil {
		return err
	}

	// fn function
	fnFunc := functions.NewFunction(
		"fn",
		registry.CategoryCore,
		2,
		"Create a function: (fn [params] body)",
		p.fnFunc,
	)
	if err := reg.Register(fnFunc); err != nil {
		return err
	}

	// defn function
	defnFunc := functions.NewFunction(
		"defn",
		registry.CategoryCore,
		3,
		"Define a named function: (defn name [params] body)",
		p.defnFunc,
	)
	if err := reg.Register(defnFunc); err != nil {
		return err
	}

	// quote function
	quoteFunc := functions.NewFunction(
		"quote",
		registry.CategoryCore,
		1,
		"Quote an expression: (quote expr)",
		p.quoteFunc,
	)
	if err := reg.Register(quoteFunc); err != nil {
		return err
	}

	// help function
	helpFunc := functions.NewFunction(
		"help",
		registry.CategoryCore,
		-1, // Variable arity: 0 or 1 argument
		"Show help: (help) for all functions, (help function-name) for specific function",
		p.helpFunc,
	)
	if err := reg.Register(helpFunc); err != nil {
		return err
	}

	// env function
	envFunc := functions.NewFunction(
		"env",
		registry.CategoryCore,
		0,
		"Show environment variables: (env)",
		p.envFunc,
	)
	if err := reg.Register(envFunc); err != nil {
		return err
	}

	// modules function
	modulesFunc := functions.NewFunction(
		"modules",
		registry.CategoryCore,
		0,
		"Show loaded modules/plugins: (modules)",
		p.modulesFunc,
	)
	if err := reg.Register(modulesFunc); err != nil {
		return err
	}

	return nil
}

// defFunc implements variable definition
func (p *CorePlugin) defFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("def requires exactly 2 arguments, got %d", len(args))
	}

	// First argument should be a symbol
	symbolExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("def requires a symbol as first argument, got %T", args[0])
	}

	// Evaluate the second argument
	value, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Define the variable in the environment
	p.env.Set(symbolExpr.Name, value)
	return value, nil
}

// fnFunc creates a function
func (p *CorePlugin) fnFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("fn requires exactly 2 arguments, got %d", len(args))
	}

	// First argument should be a bracket expression for parameters [x y z]
	var params []string
	switch paramExpr := args[0].(type) {
	case *types.BracketExpr:
		// Extract parameter names from bracket expression
		for i, param := range paramExpr.Elements {
			if symbol, ok := param.(*types.SymbolExpr); ok {
				params = append(params, symbol.Name)
			} else {
				return nil, fmt.Errorf("fn parameter %d must be a symbol, got %T", i, param)
			}
		}
	case *types.ListExpr:
		// Also support list expression for backward compatibility
		for i, param := range paramExpr.Elements {
			if symbol, ok := param.(*types.SymbolExpr); ok {
				params = append(params, symbol.Name)
			} else {
				return nil, fmt.Errorf("fn parameter %d must be a symbol, got %T", i, param)
			}
		}
	default:
		return nil, fmt.Errorf("fn requires a parameter list as first argument, got %T", args[0])
	}

	// Body is the second argument
	body := args[1]

	// Create function value
	fn := &types.FunctionValue{
		Params: params,
		Body:   body,
		Env:    p.env, // Capture current environment
	}

	return fn, nil
}

// defnFunc defines a named function
func (p *CorePlugin) defnFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("defn requires exactly 3 arguments, got %d", len(args))
	}

	// First argument should be a symbol (function name)
	nameExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("defn requires a symbol as first argument, got %T", args[0])
	}

	// Create the function using fn logic
	fnValue, err := p.fnFunc(evaluator, args[1:])
	if err != nil {
		return nil, err
	}

	// Define the function in the environment
	p.env.Set(nameExpr.Name, fnValue)
	return fnValue, nil
}

// quoteFunc quotes an expression
func (p *CorePlugin) quoteFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("quote requires exactly 1 argument, got %d", len(args))
	}

	// Return the expression as a quoted value
	return p.exprToValue(args[0]), nil
}

// exprToValue converts an expression to a value for quoting
func (p *CorePlugin) exprToValue(expr types.Expr) types.Value {
	switch ex := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(ex.Value)
	case *types.StringExpr:
		return types.StringValue(ex.Value)
	case *types.BooleanExpr:
		return types.BooleanValue(ex.Value)
	case *types.KeywordExpr:
		return types.KeywordValue(ex.Value)
	case *types.SymbolExpr:
		return types.StringValue(ex.Name)
	case *types.ListExpr:
		var elements []types.Value
		for _, elem := range ex.Elements {
			elements = append(elements, p.exprToValue(elem))
		}
		return &types.ListValue{Elements: elements}
	default:
		// For other types, return as string representation
		return types.StringValue(fmt.Sprintf("%v", expr))
	}
}

// helpFunc shows help for functions
func (p *CorePlugin) helpFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	// If no arguments, show help for all functions
	if len(args) == 0 {
		return p.allFunctionsHelp(), nil
	}

	// If one argument, show help for specific function
	if len(args) == 1 {
		funcName, ok := args[0].(*types.SymbolExpr)
		if !ok {
			return nil, fmt.Errorf("help function name must be a symbol, got %T", args[0])
		}
		return p.specificFunctionHelp(funcName.Name), nil
	}

	return nil, fmt.Errorf("help requires 0 or 1 argument, got %d", len(args))
}

// allFunctionsHelp returns a list of all function names organized by category
func (p *CorePlugin) allFunctionsHelp() types.Value {
	if p.registry == nil {
		return types.StringValue("Registry not available")
	}

	categories := p.registry.Categories()
	var helpLines []string

	helpLines = append(helpLines, "Available functions by category:")
	helpLines = append(helpLines, "")

	for _, category := range categories {
		helpLines = append(helpLines, fmt.Sprintf("=== %s ===", category))
		functions := p.registry.ListByCategory(category)

		// List functions in rows for better readability
		var currentLine strings.Builder
		functionsPerLine := 3 // Reduce to 3 per line for better spacing
		for i, funcName := range functions {
			if i > 0 && i%functionsPerLine == 0 {
				helpLines = append(helpLines, currentLine.String())
				currentLine.Reset()
			}
			if currentLine.Len() > 0 {
				currentLine.WriteString("  ")
			}
			currentLine.WriteString(fmt.Sprintf("%-20s", funcName)) // Increase spacing
		}
		if currentLine.Len() > 0 {
			helpLines = append(helpLines, currentLine.String())
		}
		helpLines = append(helpLines, "")
	}

	helpLines = append(helpLines, "Use (help function-name) for detailed help on a specific function.")

	// Return as a single string, not a list
	return types.StringValue(strings.Join(helpLines, "\n"))
}

// specificFunctionHelp returns help information for a specific function
func (p *CorePlugin) specificFunctionHelp(name string) types.Value {
	if p.registry == nil {
		return types.StringValue("Registry not available")
	}

	if fn, exists := p.registry.Get(name); exists {
		arityStr := fmt.Sprintf("%d", fn.Arity())
		if fn.Arity() == -1 {
			arityStr = "variable"
		}

		helpText := fmt.Sprintf("Function: %s\nCategory: %s\nArity: %s\nHelp: %s",
			fn.Name(), fn.Category(), arityStr, fn.Help())
		return types.StringValue(helpText)
	}
	return types.StringValue(fmt.Sprintf("Function not found: %s", name))
}

// envFunc shows environment variables
func (p *CorePlugin) envFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("env requires no arguments, got %d", len(args))
	}

	var result strings.Builder
	result.WriteString("Environment Variables and Functions:\n\n")

	// Get bindings from the current environment
	bindings := p.env.GetBindings()

	// Separate variables and functions
	var variables []string
	var functions []string
	var builtins []string

	for name, value := range bindings {
		switch value.(type) {
		case *types.FunctionValue:
			functions = append(functions, name)
		case *types.ArithmeticFunctionValue, *types.BuiltinFunctionValue:
			builtins = append(builtins, name)
		default:
			variables = append(variables, fmt.Sprintf("%s = %s", name, value.String()))
		}
	}

	// Sort the slices for consistent output
	sort.Strings(variables)
	sort.Strings(functions)
	sort.Strings(builtins)

	// Display user-defined variables
	if len(variables) > 0 {
		result.WriteString("=== User-defined Variables ===\n")
		for _, variable := range variables {
			result.WriteString("  " + variable + "\n")
		}
		result.WriteString("\n")
	} else {
		result.WriteString("=== User-defined Variables ===\n")
		result.WriteString("  (none)\n\n")
	}

	// Display user-defined functions
	if len(functions) > 0 {
		result.WriteString("=== User-defined Functions ===\n")
		for _, function := range functions {
			if value, exists := bindings[function]; exists {
				if fn, ok := value.(*types.FunctionValue); ok {
					result.WriteString(fmt.Sprintf("  %s(%s)\n", function, strings.Join(fn.Params, " ")))
				} else {
					result.WriteString(fmt.Sprintf("  %s\n", function))
				}
			}
		}
		result.WriteString("\n")
	} else {
		result.WriteString("=== User-defined Functions ===\n")
		result.WriteString("  (none)\n\n")
	}

	// Display some built-in functions (limit to avoid clutter)
	if len(builtins) > 0 {
		result.WriteString("=== Built-in Functions (sample) ===\n")
		count := 0
		for _, builtin := range builtins {
			if count < 10 { // Show only first 10 to avoid clutter
				result.WriteString("  " + builtin + "\n")
				count++
			}
		}
		if len(builtins) > 10 {
			result.WriteString(fmt.Sprintf("  ... and %d more (use 'help' to see all)\n", len(builtins)-10))
		}
		result.WriteString("\n")
	}

	result.WriteString("Use (help) to see all available functions.")

	return types.StringValue(result.String()), nil
}

// modulesFunc shows loaded modules/plugins
func (p *CorePlugin) modulesFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("modules requires no arguments, got %d", len(args))
	}

	if p.registry == nil {
		return types.StringValue("Registry not available"), nil
	}

	categories := p.registry.Categories()
	var moduleLines []string

	moduleLines = append(moduleLines, "Loaded plugin categories:")
	for _, category := range categories {
		functions := p.registry.ListByCategory(category)
		moduleLines = append(moduleLines, fmt.Sprintf("  %s (%d functions)", category, len(functions)))
	}

	return types.StringValue(strings.Join(moduleLines, "\n")), nil
}
