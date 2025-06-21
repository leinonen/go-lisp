// Package binding provides let bindings for the Lisp interpreter
package binding

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/interfaces"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// BindingPlugin provides let binding functionality
type BindingPlugin struct {
	*plugins.BasePlugin
	evaluator interfaces.CoreEvaluator
}

// NewBindingPlugin creates a new binding plugin with dependencies
func NewBindingPlugin(evaluator interfaces.CoreEvaluator) *BindingPlugin {
	return &BindingPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"binding",
			"1.0.0",
			"Local variable bindings (let)",
			[]string{"core"}, // Depends on core
		),
		evaluator: evaluator,
	}
}

// RegisterFunctions registers binding functions
func (p *BindingPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// let function
	letFunc := functions.NewFunction(
		"let",
		registry.CategoryCore,
		2,
		"Local bindings: (let [var1 val1 var2 val2] body)",
		p.evalLet,
	)
	return reg.Register(letFunc)
}

// evalLet implements local variable bindings
func (p *BindingPlugin) evalLet(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("let requires exactly 2 arguments, got %d", len(args))
	}

	// First argument should be a binding vector [var1 val1 var2 val2]
	bindingsExpr, ok := args[0].(*types.BracketExpr)
	if !ok {
		return nil, fmt.Errorf("let requires a bracket expression for bindings, got %T", args[0])
	}

	// Bindings should come in pairs
	if len(bindingsExpr.Elements)%2 != 0 {
		return nil, fmt.Errorf("let bindings must come in pairs, got %d elements", len(bindingsExpr.Elements))
	}

	// For now, create a simple let implementation using do and def
	// This is a temporary solution until we can properly integrate with environments

	// Extract bindings and body
	var doArgs []types.Expr

	// Add def statements for each binding
	for i := 0; i < len(bindingsExpr.Elements); i += 2 {
		varExpr := bindingsExpr.Elements[i]
		valueExpr := bindingsExpr.Elements[i+1]

		// Create (def var value) expression
		defExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "def"},
				varExpr,
				valueExpr,
			},
		}
		doArgs = append(doArgs, defExpr)
	}

	// Add the body expression
	doArgs = append(doArgs, args[1])

	// Create and evaluate (do def1 def2 ... body)
	doExpr := &types.ListExpr{
		Elements: append([]types.Expr{&types.SymbolExpr{Name: "do"}}, doArgs...),
	}

	// Use helper method instead
	return p.eval(evaluator, doExpr)
}

// Use injected evaluator when available
func (p *BindingPlugin) eval(evaluator registry.Evaluator, expr types.Expr) (types.Value, error) {
	if p.evaluator != nil {
		return p.evaluator.Eval(expr)
	}
	return evaluator.Eval(expr)
}
