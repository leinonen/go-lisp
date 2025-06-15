// Package core provides core language functionality for the Lisp interpreter
package core

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// CorePlugin provides core language functionality
type CorePlugin struct {
	*plugins.BasePlugin
	env *evaluator.Environment
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
		env: env,
	}
}

// RegisterFunctions registers core language functions
func (p *CorePlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
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
