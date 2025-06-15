// Package control provides control flow operations as a plugin
package control

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// ControlPlugin provides control flow operations
type ControlPlugin struct {
	*plugins.BasePlugin
}

// NewControlPlugin creates a new control plugin
func NewControlPlugin() *ControlPlugin {
	return &ControlPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"control",
			"1.0.0",
			"Control flow operations (if, do)",
			[]string{"logical"}, // Depends on logical for truthiness
		),
	}
}

// RegisterFunctions registers control flow functions
func (cp *ControlPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// If
	ifFunc := functions.NewFunction(
		"if",
		registry.CategoryControl,
		-1, // 2 or 3 arguments
		"Conditional: (if condition then-expr else-expr?) => evaluated result",
		cp.evalIf,
	)
	if err := reg.Register(ifFunc); err != nil {
		return err
	}

	// Do - sequential evaluation
	doFunc := functions.NewFunction(
		"do",
		registry.CategoryControl,
		-1, // Variadic
		"Sequential evaluation: (do expr1 expr2 expr3) => result of last expr",
		cp.evalDo,
	)
	return reg.Register(doFunc)
}

// evalIf implements conditional evaluation
func (cp *ControlPlugin) evalIf(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("if requires 2 or 3 arguments, got %d", len(args))
	}

	// Evaluate condition
	condition, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("if condition evaluation failed: %v", err)
	}

	// Check if condition is truthy
	if functions.IsTruthy(condition) {
		// Evaluate then branch
		return evaluator.Eval(args[1])
	} else if len(args) == 3 {
		// Evaluate else branch if provided
		return evaluator.Eval(args[2])
	}

	// No else branch and condition is false
	return types.BooleanValue(false), nil
}

// evalDo implements sequential evaluation
func (cp *ControlPlugin) evalDo(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("do requires at least 1 argument")
	}

	var result types.Value = types.BooleanValue(false) // Default result

	// Evaluate each expression in sequence
	for _, expr := range args {
		var err error
		result, err = evaluator.Eval(expr)
		if err != nil {
			return nil, err
		}
	}

	// Return the result of the last expression
	return result, nil
}
