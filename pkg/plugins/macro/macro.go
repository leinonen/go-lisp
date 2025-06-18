// Package macro provides macro support for the Lisp interpreter
package macro

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// MacroPlugin provides macro functionality
type MacroPlugin struct {
	*plugins.BasePlugin
}

// NewMacroPlugin creates a new macro plugin
func NewMacroPlugin() *MacroPlugin {
	return &MacroPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"macro",
			"1.0.0",
			"Macro support (defmacro, macroexpand)",
			[]string{"core"}, // Depends on core
		),
	}
}

// RegisterFunctions registers macro functions
func (p *MacroPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// defmacro function
	defmacroFunc := functions.NewFunction(
		"defmacro",
		registry.CategoryCore,
		3,
		"Define a macro: (defmacro name [params] body)",
		p.evalDefmacro,
	)
	if err := reg.Register(defmacroFunc); err != nil {
		return err
	}

	// macroexpand function
	macroexpandFunc := functions.NewFunction(
		"macroexpand",
		registry.CategoryCore,
		1,
		"Expand a macro: (macroexpand macro-call)",
		p.evalMacroexpand,
	)
	if err := reg.Register(macroexpandFunc); err != nil {
		return err
	}

	// Note: quote is already provided by core plugin

	// unquote function (for use within quasiquote)
	unquoteFunc := functions.NewFunction(
		"unquote",
		registry.CategoryCore,
		1,
		"Unquote an expression (for use in quasiquote): (unquote expr)",
		p.evalUnquote,
	)
	return reg.Register(unquoteFunc)
}

// evalDefmacro defines a macro
func (p *MacroPlugin) evalDefmacro(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("defmacro requires exactly 3 arguments, got %d", len(args))
	}

	// First argument should be a symbol (macro name)
	nameExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("defmacro requires a symbol as first argument, got %T", args[0])
	}
	macroName := nameExpr.Name

	// Second argument should be parameter list
	var params []string
	switch paramExpr := args[1].(type) {
	case *types.BracketExpr:
		for i, param := range paramExpr.Elements {
			if symbol, ok := param.(*types.SymbolExpr); ok {
				params = append(params, symbol.Name)
			} else {
				return nil, fmt.Errorf("defmacro parameter %d must be a symbol, got %T", i, param)
			}
		}
	default:
		return nil, fmt.Errorf("defmacro requires a bracket expression for parameters, got %T", args[1])
	}

	// Third argument is the macro body
	body := args[2]

	// Create macro value
	macro := &types.MacroValue{
		Params: params,
		Body:   body,
	}

	// TODO: Store macro in environment with macroName
	// This would need access to the environment to store the macro
	// For now, return the macro value
	_ = macroName // Acknowledge we have the name but can't use it yet
	return macro, nil
}

// evalMacroexpand expands a macro call
func (p *MacroPlugin) evalMacroexpand(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("macroexpand requires exactly 1 argument, got %d", len(args))
	}

	// Check if the argument is a macro call
	if listExpr, ok := args[0].(*types.ListExpr); ok && len(listExpr.Elements) > 0 {
		if symbol, ok := listExpr.Elements[0].(*types.SymbolExpr); ok {
			// In a real implementation, we would look up the macro in the environment
			// For now, create a simple example expansion
			macroName := symbol.Name
			macroArgs := listExpr.Elements[1:]

			// Example: expand (when condition body) to (if condition body nil)
			if macroName == "when" && len(macroArgs) == 2 {
				expanded := &types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						macroArgs[0],                   // condition
						macroArgs[1],                   // body
						&types.SymbolExpr{Name: "nil"}, // else clause
					},
				}
				return &types.QuotedValue{Value: expanded}, nil
			}

			// Example: expand (unless condition body) to (if condition nil body)
			if macroName == "unless" && len(macroArgs) == 2 {
				expanded := &types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "if"},
						macroArgs[0],                   // condition
						&types.SymbolExpr{Name: "nil"}, // then clause
						macroArgs[1],                   // else clause (body)
					},
				}
				return &types.QuotedValue{Value: expanded}, nil
			}
		}
	}

	// If not a recognized macro, return it quoted
	return &types.QuotedValue{Value: args[0]}, nil
}

// evalUnquote evaluates an expression (opposite of quote)
func (p *MacroPlugin) evalUnquote(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("unquote requires exactly 1 argument, got %d", len(args))
	}

	// Unquote evaluates the expression
	return evaluator.Eval(args[0])
}

// expandMacro expands a macro call with given arguments
func (p *MacroPlugin) expandMacro(macro *types.MacroValue, args []types.Expr) (types.Expr, error) {
	if len(args) != len(macro.Params) {
		return nil, fmt.Errorf("macro expects %d arguments, got %d", len(macro.Params), len(args))
	}

	// Simple macro expansion: substitute parameters in the body
	// This is a basic implementation - a full macro system would need
	// proper hygiene and more sophisticated substitution
	return p.substituteInExpr(macro.Body, macro.Params, args), nil
}

// substituteInExpr performs simple parameter substitution in an expression
func (p *MacroPlugin) substituteInExpr(expr types.Expr, params []string, args []types.Expr) types.Expr {
	switch e := expr.(type) {
	case *types.SymbolExpr:
		// Check if this symbol is a parameter
		for i, param := range params {
			if e.Name == param {
				return args[i] // Replace with the corresponding argument
			}
		}
		return e // Not a parameter, return as-is

	case *types.ListExpr:
		// Recursively substitute in list elements
		var newElements []types.Expr
		for _, elem := range e.Elements {
			newElements = append(newElements, p.substituteInExpr(elem, params, args))
		}
		return &types.ListExpr{Elements: newElements}

	case *types.BracketExpr:
		// Recursively substitute in bracket elements
		var newElements []types.Expr
		for _, elem := range e.Elements {
			newElements = append(newElements, p.substituteInExpr(elem, params, args))
		}
		return &types.BracketExpr{Elements: newElements}

	default:
		// For literals (numbers, strings, etc.), return as-is
		return e
	}
}
