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
			"Control flow operations (if, do, cond, when, when-not, loop, recur)",
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
	if err := reg.Register(whenNotFunc); err != nil {
		return err
	}

	// loop - establish recursion point
	loopFunc := functions.NewFunction(
		"loop",
		registry.CategoryControl,
		-1, // At least 1 argument: binding vector and body
		"Establish recursion point: (loop [bindings*] body*)",
		cp.evalLoop,
	)
	if err := reg.Register(loopFunc); err != nil {
		return err
	}

	// recur - jump back to loop with new values
	recurFunc := functions.NewFunction(
		"recur",
		registry.CategoryControl,
		-1, // Variadic - number of args must match loop bindings
		"Jump back to loop: (recur args*)",
		cp.evalRecur,
	)
	return reg.Register(recurFunc)
}

// evalLoop implements loop/recur functionality
// Loop establishes a recursion point with local bindings
// Syntax: (loop [var1 init1 var2 init2 ...] body...)
// The body is evaluated with the bindings in scope
// If recur is called, execution jumps back to the loop with new binding values
func (cp *ControlPlugin) evalLoop(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("loop requires at least 1 argument (binding vector)")
	}

	// First argument should be a binding vector [name1 val1 name2 val2 ...]
	bindings, ok := args[0].(*types.BracketExpr)
	if !ok {
		return nil, fmt.Errorf("loop first argument must be a binding vector")
	}

	if len(bindings.Elements)%2 != 0 {
		return nil, fmt.Errorf("loop binding vector must have even number of elements")
	}

	// Extract binding names and initial values
	var bindingNames []string
	var bindingValues []types.Value

	for i := 0; i < len(bindings.Elements); i += 2 {
		nameExpr, ok := bindings.Elements[i].(*types.SymbolExpr)
		if !ok {
			return nil, fmt.Errorf("loop binding names must be symbols")
		}
		bindingNames = append(bindingNames, nameExpr.Name)

		// Evaluate the initial value
		value, err := evaluator.Eval(bindings.Elements[i+1])
		if err != nil {
			return nil, fmt.Errorf("loop binding value evaluation failed: %v", err)
		}
		bindingValues = append(bindingValues, value)
	}

	// Body expressions (everything after the binding vector)
	body := args[1:]

	// Main loop - continue until no recur exception is thrown
	for {
		// Create a map of current bindings for this iteration
		currentBindings := make(map[string]types.Value)
		for i, name := range bindingNames {
			currentBindings[name] = bindingValues[i]
		}

		var result types.Value = types.BooleanValue(false)
		var err error
		var recurCalled bool

		// Evaluate each body expression with the current bindings
		for _, expr := range body {
			result, err = evaluator.EvalWithBindings(expr, currentBindings)
			if err != nil {
				// Check if it's a recur exception
				if recurErr, ok := err.(*types.RecurException); ok {
					// Validate argument count matches binding count
					if len(recurErr.Args) != len(bindingNames) {
						return nil, fmt.Errorf("recur argument count (%d) doesn't match loop binding count (%d)",
							len(recurErr.Args), len(bindingNames))
					}

					// Update binding values for next iteration
					bindingValues = recurErr.Args
					recurCalled = true
					break // Break out of body evaluation and continue loop
				}
				return nil, err
			}
		}

		// If we completed the body without a recur, return the result
		if !recurCalled {
			return result, nil
		}

		// If we got here due to recur, continue the loop with updated bindings
		// (bindingValues has been updated above)
	}
}

// evalRecur implements recur functionality
// Recur jumps back to the nearest enclosing loop with new binding values
// Syntax: (recur arg1 arg2 ...)
// The number of arguments must match the number of bindings in the loop
// This throws a RecurException that is caught by the loop construct
func (cp *ControlPlugin) evalRecur(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	// Evaluate all arguments
	var values []types.Value
	for _, arg := range args {
		value, err := evaluator.Eval(arg)
		if err != nil {
			return nil, fmt.Errorf("recur argument evaluation failed: %v", err)
		}
		values = append(values, value)
	}

	// Throw a RecurException to be caught by the nearest loop
	// This is how we implement the control flow jump
	return nil, types.NewRecurException(values)
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
