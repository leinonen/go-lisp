// Package environment provides enhanced let bindings with proper environment scoping
package environment

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// EnvironmentPlugin provides proper environment scoping for let bindings
type EnvironmentPlugin struct {
	*plugins.BasePlugin
	env *evaluator.Environment
}

// NewEnvironmentPlugin creates a new environment plugin
func NewEnvironmentPlugin(env *evaluator.Environment) *EnvironmentPlugin {
	return &EnvironmentPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"environment",
			"1.0.0",
			"Enhanced environment scoping (let*, letfn)",
			[]string{"core", "binding"}, // Depends on core and basic binding
		),
		env: env,
	}
}

// RegisterFunctions registers environment functions
func (p *EnvironmentPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// let* function - sequential let with proper scoping
	letStarFunc := functions.NewFunction(
		"let*",
		registry.CategoryCore,
		2,
		"Sequential let with proper scoping: (let* [x 10 y (+ x 5)] y)",
		p.evalLetStar,
	)
	if err := reg.Register(letStarFunc); err != nil {
		return err
	}

	// letfn function - function bindings
	letfnFunc := functions.NewFunction(
		"letfn",
		registry.CategoryCore,
		2,
		"Local function bindings: (letfn [(f [x] (* x x))] (f 5))",
		p.evalLetfn,
	)
	return reg.Register(letfnFunc)
}

// evalLetStar implements sequential let bindings with proper scoping
func (p *EnvironmentPlugin) evalLetStar(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("let* requires exactly 2 arguments, got %d", len(args))
	}

	// First argument should be a binding vector [var1 val1 var2 val2]
	bindingsExpr, ok := args[0].(*types.BracketExpr)
	if !ok {
		return nil, fmt.Errorf("let* requires a bracket expression for bindings, got %T", args[0])
	}

	// Bindings should come in pairs
	if len(bindingsExpr.Elements)%2 != 0 {
		return nil, fmt.Errorf("let* bindings must come in pairs, got %d elements", len(bindingsExpr.Elements))
	}

	// For proper let* implementation, we would need:
	// 1. Create a new child environment
	// 2. Sequentially evaluate and bind each variable
	// 3. Evaluate the body in the new environment
	// 4. Return the result

	// Since we don't have direct access to environment creation in the evaluator interface,
	// we'll implement a sequential version using nested let expressions
	if len(bindingsExpr.Elements) == 0 {
		// No bindings, just evaluate body
		return evaluator.Eval(args[1])
	}

	// Build nested let expressions for sequential binding
	return p.buildSequentialLet(evaluator, bindingsExpr.Elements, args[1])
}

// buildSequentialLet builds nested let expressions for sequential binding
func (p *EnvironmentPlugin) buildSequentialLet(evaluator registry.Evaluator, bindings []types.Expr, body types.Expr) (types.Value, error) {
	if len(bindings) == 0 {
		return evaluator.Eval(body)
	}

	if len(bindings) < 2 {
		return nil, fmt.Errorf("invalid bindings for sequential let")
	}

	// Take the first binding pair
	varExpr := bindings[0]
	valueExpr := bindings[1]
	remainingBindings := bindings[2:]

	// Create the inner expression
	var innerExpr types.Expr
	if len(remainingBindings) > 0 {
		// More bindings to process - create nested let*
		innerLet := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "let*"},
				&types.BracketExpr{Elements: remainingBindings},
				body,
			},
		}
		innerExpr = innerLet
	} else {
		// Last binding - use the body directly
		innerExpr = body
	}

	// Create the outer let expression
	letExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "let"},
			&types.BracketExpr{Elements: []types.Expr{varExpr, valueExpr}},
			innerExpr,
		},
	}

	return evaluator.Eval(letExpr)
}

// evalLetfn implements local function bindings
func (p *EnvironmentPlugin) evalLetfn(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("letfn requires exactly 2 arguments, got %d", len(args))
	}

	// First argument should be a vector of function definitions
	bindingsExpr, ok := args[0].(*types.BracketExpr)
	if !ok {
		return nil, fmt.Errorf("letfn requires a bracket expression for function bindings, got %T", args[0])
	}

	// For now, implement as a series of defn calls within a do block
	var doArgs []types.Expr

	// Process each function definition
	for _, binding := range bindingsExpr.Elements {
		funcDef, ok := binding.(*types.BracketExpr)
		if !ok || len(funcDef.Elements) != 3 {
			return nil, fmt.Errorf("letfn function binding must be [name params body]")
		}

		// Create (defn name params body) expression
		defnExpr := &types.ListExpr{
			Elements: []types.Expr{
				&types.SymbolExpr{Name: "defn"},
				funcDef.Elements[0], // name
				funcDef.Elements[1], // params
				funcDef.Elements[2], // body
			},
		}
		doArgs = append(doArgs, defnExpr)
	}

	// Add the body expression
	doArgs = append(doArgs, args[1])

	// Create and evaluate (do defn1 defn2 ... body)
	doExpr := &types.ListExpr{
		Elements: append([]types.Expr{&types.SymbolExpr{Name: "do"}}, doArgs...),
	}

	return evaluator.Eval(doExpr)
}

// delegateToLet delegates to the basic let implementation
func (p *EnvironmentPlugin) delegateToLet(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	// Create (let bindings body) expression
	letExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "let"},
			args[0], // bindings
			args[1], // body
		},
	}
	return evaluator.Eval(letExpr)
}
