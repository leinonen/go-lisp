// Package advanced provides advanced binding features like destructuring
package advanced

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// AdvancedPlugin provides advanced binding features
type AdvancedPlugin struct {
	*plugins.BasePlugin
}

// NewAdvancedPlugin creates a new advanced binding plugin
func NewAdvancedPlugin() *AdvancedPlugin {
	return &AdvancedPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"advanced",
			"1.0.0",
			"Advanced binding features (destructuring, pattern matching)",
			[]string{"core", "binding"}, // Depends on core and basic binding
		),
	}
}

// RegisterFunctions registers advanced binding functions
func (p *AdvancedPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// let-destructure function for destructuring bindings
	letDestructureFunc := functions.NewFunction(
		"let-destructure",
		registry.CategoryCore,
		2,
		"Destructuring let: (let-destructure [[a b] [1 2]] (+ a b))",
		p.evalLetDestructure,
	)
	if err := reg.Register(letDestructureFunc); err != nil {
		return err
	}

	// fn-destructure function for destructuring in function parameters
	fnDestructureFunc := functions.NewFunction(
		"fn-destructure",
		registry.CategoryCore,
		2,
		"Destructuring function: (fn-destructure [[a b]] (+ a b))",
		p.evalFnDestructure,
	)
	return reg.Register(fnDestructureFunc)
}

// evalLetDestructure implements destructuring let bindings
func (p *AdvancedPlugin) evalLetDestructure(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("let-destructure requires exactly 2 arguments, got %d", len(args))
	}

	// First argument should be a binding vector with destructuring patterns
	bindingsExpr, ok := args[0].(*types.BracketExpr)
	if !ok {
		return nil, fmt.Errorf("let-destructure requires a bracket expression for bindings, got %T", args[0])
	}

	// Bindings should come in pairs: [pattern value pattern value ...]
	if len(bindingsExpr.Elements)%2 != 0 {
		return nil, fmt.Errorf("let-destructure bindings must come in pairs, got %d elements", len(bindingsExpr.Elements))
	}

	// Expand destructuring patterns into regular bindings
	var expandedBindings []types.Expr
	for i := 0; i < len(bindingsExpr.Elements); i += 2 {
		pattern := bindingsExpr.Elements[i]
		value := bindingsExpr.Elements[i+1]

		bindings, err := p.expandPattern(pattern, value)
		if err != nil {
			return nil, fmt.Errorf("failed to expand destructuring pattern: %v", err)
		}
		expandedBindings = append(expandedBindings, bindings...)
	}

	// Create a regular let expression with expanded bindings
	letExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "let"},
			&types.BracketExpr{Elements: expandedBindings},
			args[1], // body
		},
	}

	return evaluator.Eval(letExpr)
}

// evalFnDestructure implements destructuring in function parameters
func (p *AdvancedPlugin) evalFnDestructure(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("fn-destructure requires exactly 2 arguments, got %d", len(args))
	}

	// First argument should be parameter patterns
	paramsExpr, ok := args[0].(*types.BracketExpr)
	if !ok {
		return nil, fmt.Errorf("fn-destructure requires a bracket expression for parameters, got %T", args[0])
	}

	// For now, create a simple function that takes the first parameter as a list
	// and destructures it. This is a simplified implementation.
	if len(paramsExpr.Elements) != 1 {
		return nil, fmt.Errorf("fn-destructure currently supports only one destructuring parameter")
	}

	pattern := paramsExpr.Elements[0]

	// Create a function that accepts one argument and destructures it
	fnBody := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "let-destructure"},
			&types.BracketExpr{Elements: []types.Expr{
				pattern,
				&types.SymbolExpr{Name: "arg0"}, // The function argument
			}},
			args[1], // original body
		},
	}

	// Create a regular function
	fnExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "fn"},
			&types.BracketExpr{Elements: []types.Expr{
				&types.SymbolExpr{Name: "arg0"},
			}},
			fnBody,
		},
	}

	return evaluator.Eval(fnExpr)
}

// expandPattern expands a destructuring pattern into individual bindings
func (p *AdvancedPlugin) expandPattern(pattern types.Expr, value types.Expr) ([]types.Expr, error) {
	switch pat := pattern.(type) {
	case *types.SymbolExpr:
		// Simple binding - no destructuring needed
		return []types.Expr{pat, value}, nil

	case *types.BracketExpr:
		// Vector destructuring: [a b c] destructures a sequence
		var bindings []types.Expr
		for i, elem := range pat.Elements {
			if symbol, ok := elem.(*types.SymbolExpr); ok {
				// Create (nth value i) expression
				nthExpr := &types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "nth"},
						value,
						&types.NumberExpr{Value: float64(i)},
					},
				}
				bindings = append(bindings, symbol, nthExpr)
			} else {
				return nil, fmt.Errorf("destructuring pattern elements must be symbols, got %T", elem)
			}
		}
		return bindings, nil

	case *types.ListExpr:
		// List destructuring - similar to vector but uses list functions
		var bindings []types.Expr
		for i, elem := range pat.Elements {
			if symbol, ok := elem.(*types.SymbolExpr); ok {
				// Create (nth value i) expression
				nthExpr := &types.ListExpr{
					Elements: []types.Expr{
						&types.SymbolExpr{Name: "nth"},
						value,
						&types.NumberExpr{Value: float64(i)},
					},
				}
				bindings = append(bindings, symbol, nthExpr)
			} else {
				return nil, fmt.Errorf("destructuring pattern elements must be symbols, got %T", elem)
			}
		}
		return bindings, nil

	default:
		return nil, fmt.Errorf("unsupported destructuring pattern type: %T", pattern)
	}
}
