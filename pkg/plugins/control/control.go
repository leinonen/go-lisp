// Package control provides control flow operations as a plugin
package control

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
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
	if err := reg.Register(doFunc); err != nil {
		return err
	}

	// cond - multi-branch conditional
	condFunc := functions.NewFunction(
		"cond",
		registry.CategoryControl,
		-1, // Variadic pairs
		"Multi-branch conditional: (cond test1 expr1 test2 expr2 :else default)",
		cp.evalCond,
	)
	if err := reg.Register(condFunc); err != nil {
		return err
	}

	// when - conditional execution
	whenFunc := functions.NewFunction(
		"when",
		registry.CategoryControl,
		-1, // At least 1 argument
		"Conditional execution: (when test expr1 expr2 ...)",
		cp.evalWhen,
	)
	if err := reg.Register(whenFunc); err != nil {
		return err
	}

	// when-not - negated conditional execution
	whenNotFunc := functions.NewFunction(
		"when-not",
		registry.CategoryControl,
		-1, // At least 1 argument
		"Negated conditional execution: (when-not test expr1 expr2 ...)",
		cp.evalWhenNot,
	)
	return reg.Register(whenNotFunc)
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

// evalCond implements multi-branch conditional evaluation
func (cp *ControlPlugin) evalCond(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return types.BooleanValue(false), nil
	}

	// Process condition-expression pairs
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			return nil, fmt.Errorf("cond requires condition-expression pairs")
		}

		// Check for :else keyword
		if symbol, ok := args[i].(*types.SymbolExpr); ok && symbol.Name == ":else" {
			return evaluator.Eval(args[i+1])
		}

		// Evaluate condition
		condition, err := evaluator.Eval(args[i])
		if err != nil {
			return nil, fmt.Errorf("cond condition evaluation failed: %v", err)
		}

		// If condition is truthy, evaluate and return the expression
		if functions.IsTruthy(condition) {
			return evaluator.Eval(args[i+1])
		}
	}

	// No condition matched
	return types.BooleanValue(false), nil
}

// evalWhen implements conditional execution
func (cp *ControlPlugin) evalWhen(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("when requires at least 1 argument")
	}

	// Evaluate condition
	condition, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("when condition evaluation failed: %v", err)
	}

	// If condition is truthy, evaluate all expressions
	if functions.IsTruthy(condition) {
		var result types.Value = types.BooleanValue(true)
		for _, expr := range args[1:] {
			result, err = evaluator.Eval(expr)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}

	// Condition is false
	return types.BooleanValue(false), nil
}

// evalWhenNot implements negated conditional execution
func (cp *ControlPlugin) evalWhenNot(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("when-not requires at least 1 argument")
	}

	// Evaluate condition
	condition, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("when-not condition evaluation failed: %v", err)
	}

	// If condition is falsy, evaluate all expressions
	if !functions.IsTruthy(condition) {
		var result types.Value = types.BooleanValue(true)
		for _, expr := range args[1:] {
			result, err = evaluator.Eval(expr)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}

	// Condition is true
	return types.BooleanValue(false), nil
}
