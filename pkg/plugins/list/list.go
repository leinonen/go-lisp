// Package list provides list operations as a plugin
package list

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// ListPlugin provides list operations
type ListPlugin struct {
	*plugins.BasePlugin
}

// NewListPlugin creates a new list plugin
func NewListPlugin() *ListPlugin {
	return &ListPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"list",
			"1.0.0",
			"List creation and manipulation functions (list, cons, length, append)",
			[]string{}, // No dependencies
		),
	}
}

// RegisterFunctions registers list functions
func (lp *ListPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	listFunctions := []struct {
		name    string
		arity   int
		help    string
		handler func(registry.Evaluator, []types.Expr) (types.Value, error)
	}{
		{"list", -1, "Create a list: (list 1 2 3) => (1 2 3)", lp.evalList},
		{"cons", 2, "Prepend element: (cons 0 '(1 2)) => (0 1 2)", lp.evalCons},
		{"length", 1, "Get list length: (length '(1 2 3)) => 3", lp.evalLength},
		{"append", -1, "Append lists: (append '(1 2) '(3 4)) => (1 2 3 4)", lp.evalAppend},
		// Clojure-style aliases
		{"concat", -1, "Concatenate lists: (concat '(1 2) '(3 4)) => (1 2 3 4)", lp.evalAppend},
	}

	for _, fn := range listFunctions {
		f := functions.NewFunction(fn.name, registry.CategoryList, fn.arity, fn.help, fn.handler)
		if err := reg.Register(f); err != nil {
			return fmt.Errorf("failed to register %s: %v", fn.name, err)
		}
	}

	return nil
}

// evalList creates a new list from arguments
func (lp *ListPlugin) evalList(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}
	return &types.ListValue{Elements: values}, nil
}

// evalCons prepends an element to a list
func (lp *ListPlugin) evalCons(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	element := values[0]
	list, err := functions.ExtractList(values[1])
	if err != nil {
		return nil, fmt.Errorf("cons: second argument must be a list, got %T", values[1])
	}

	newElements := make([]types.Value, len(list.Elements)+1)
	newElements[0] = element
	copy(newElements[1:], list.Elements)

	return &types.ListValue{Elements: newElements}, nil
}

// evalLength returns the length of a list
func (lp *ListPlugin) evalLength(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	list, err := functions.ExtractList(values[0])
	if err != nil {
		return nil, fmt.Errorf("length: %v", err)
	}

	return types.NumberValue(len(list.Elements)), nil
}

// evalAppend concatenates multiple lists
func (lp *ListPlugin) evalAppend(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return &types.ListValue{Elements: []types.Value{}}, nil
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	var allElements []types.Value
	for i, value := range values {
		list, err := functions.ExtractList(value)
		if err != nil {
			return nil, fmt.Errorf("append: argument %d must be a list, got %T", i+1, value)
		}
		allElements = append(allElements, list.Elements...)
	}

	return &types.ListValue{Elements: allElements}, nil
}
