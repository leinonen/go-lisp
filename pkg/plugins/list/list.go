// Package list provides list operations as a plugin
package list

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
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
			"List manipulation functions (list, first, rest, cons, length, empty?)",
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
		{"first", 1, "Get first element: (first '(1 2 3)) => 1", lp.evalFirst},
		{"rest", 1, "Get rest of list: (rest '(1 2 3)) => (2 3)", lp.evalRest},
		{"cons", 2, "Prepend element: (cons 0 '(1 2)) => (0 1 2)", lp.evalCons},
		{"length", 1, "Get list length: (length '(1 2 3)) => 3", lp.evalLength},
		{"empty?", 1, "Check if list is empty: (empty? '()) => true", lp.evalEmpty},
		{"append", -1, "Append lists: (append '(1 2) '(3 4)) => (1 2 3 4)", lp.evalAppend},
		{"reverse", 1, "Reverse list: (reverse '(1 2 3)) => (3 2 1)", lp.evalReverse},
		{"nth", 2, "Get nth element: (nth 1 '(a b c)) => b", lp.evalNth},
		{"last", 1, "Get last element: (last '(1 2 3)) => 3", lp.evalLast},
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

// evalFirst returns the first element of a list
func (lp *ListPlugin) evalFirst(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	list, err := functions.ExtractList(values[0])
	if err != nil {
		return nil, fmt.Errorf("first: %v", err)
	}

	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("first: cannot get first element of empty list")
	}

	return list.Elements[0], nil
}

// evalRest returns all elements except the first
func (lp *ListPlugin) evalRest(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	list, err := functions.ExtractList(values[0])
	if err != nil {
		return nil, fmt.Errorf("rest: %v", err)
	}

	if len(list.Elements) == 0 {
		return &types.ListValue{Elements: []types.Value{}}, nil
	}

	return &types.ListValue{Elements: list.Elements[1:]}, nil
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

// evalEmpty checks if a list is empty
func (lp *ListPlugin) evalEmpty(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	list, err := functions.ExtractList(values[0])
	if err != nil {
		return nil, fmt.Errorf("empty?: %v", err)
	}

	return types.BooleanValue(len(list.Elements) == 0), nil
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

// evalReverse reverses a list
func (lp *ListPlugin) evalReverse(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	list, err := functions.ExtractList(values[0])
	if err != nil {
		return nil, fmt.Errorf("reverse: %v", err)
	}

	reversed := make([]types.Value, len(list.Elements))
	for i, elem := range list.Elements {
		reversed[len(list.Elements)-1-i] = elem
	}

	return &types.ListValue{Elements: reversed}, nil
}

// evalNth returns the nth element of a list (0-indexed)
func (lp *ListPlugin) evalNth(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	index, err := functions.ExtractFloat64(values[0])
	if err != nil {
		return nil, fmt.Errorf("nth: first argument must be a number, got %T", values[0])
	}

	list, err := functions.ExtractList(values[1])
	if err != nil {
		return nil, fmt.Errorf("nth: second argument must be a list, got %T", values[1])
	}

	idx := int(index)
	if idx < 0 || idx >= len(list.Elements) {
		return nil, fmt.Errorf("nth: index %d out of bounds for list of length %d", idx, len(list.Elements))
	}

	return list.Elements[idx], nil
}

// evalLast returns the last element of a list
func (lp *ListPlugin) evalLast(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	list, err := functions.ExtractList(values[0])
	if err != nil {
		return nil, fmt.Errorf("last: %v", err)
	}

	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("last: cannot get last element of empty list")
	}

	return list.Elements[len(list.Elements)-1], nil
}
