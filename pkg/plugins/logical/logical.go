// Package logical provides logical operations as a plugin
package logical

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// LogicalPlugin provides logical operations
type LogicalPlugin struct {
	*plugins.BasePlugin
}

// NewLogicalPlugin creates a new logical plugin
func NewLogicalPlugin() *LogicalPlugin {
	return &LogicalPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"logical",
			"1.0.0",
			"Logical operations (and, or, not)",
			[]string{}, // No dependencies
		),
	}
}

// RegisterFunctions registers logical functions
func (lp *LogicalPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// And
	andFunc := functions.NewFunction(
		"and",
		registry.CategoryLogical,
		-1, // Variadic
		"Logical AND: (and true true false) => false, (and true true) => true",
		lp.evalAnd,
	)
	if err := reg.Register(andFunc); err != nil {
		return err
	}

	// Or
	orFunc := functions.NewFunction(
		"or",
		registry.CategoryLogical,
		-1, // Variadic
		"Logical OR: (or false false true) => true, (or false false) => false",
		lp.evalOr,
	)
	if err := reg.Register(orFunc); err != nil {
		return err
	}

	// Not
	notFunc := functions.NewFunction(
		"not",
		registry.CategoryLogical,
		1, // Exactly 1 argument
		"Logical NOT: (not true) => false, (not false) => true",
		lp.evalNot,
	)
	return reg.Register(notFunc)
}

// evalAnd implements logical AND with short-circuit evaluation
func (lp *LogicalPlugin) evalAnd(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return types.BooleanValue(true), nil // Empty AND is true
	}

	// Short-circuit evaluation: stop at first false value
	for _, arg := range args {
		value, err := evaluator.Eval(arg)
		if err != nil {
			return nil, err
		}

		if !functions.IsTruthy(value) {
			return types.BooleanValue(false), nil
		}
	}

	return types.BooleanValue(true), nil
}

// evalOr implements logical OR with short-circuit evaluation
func (lp *LogicalPlugin) evalOr(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return types.BooleanValue(false), nil // Empty OR is false
	}

	// Short-circuit evaluation: stop at first true value
	for _, arg := range args {
		value, err := evaluator.Eval(arg)
		if err != nil {
			return nil, err
		}

		if functions.IsTruthy(value) {
			return types.BooleanValue(true), nil
		}
	}

	return types.BooleanValue(false), nil
}

// evalNot implements logical NOT
func (lp *LogicalPlugin) evalNot(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("not requires exactly 1 argument, got %d", len(args))
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	return types.BooleanValue(!functions.IsTruthy(values[0])), nil
}
