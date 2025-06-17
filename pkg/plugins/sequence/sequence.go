// Package sequence provides vector and sequence operations for the Lisp interpreter
package sequence

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// SequencePlugin provides vector and sequence operations
type SequencePlugin struct {
	*plugins.BasePlugin
}

// NewSequencePlugin creates a new sequence plugin
func NewSequencePlugin() *SequencePlugin {
	return &SequencePlugin{
		BasePlugin: plugins.NewBasePlugin(
			"sequence",
			"1.0.0",
			"Vectors and sequences (vector, vec, seq)",
			[]string{"list"}, // Depends on list
		),
	}
}

// RegisterFunctions registers sequence functions
func (p *SequencePlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// vector function
	vectorFunc := functions.NewFunction(
		"vector",
		registry.CategoryList,
		-1, // Variadic
		"Create a vector: (vector 1 2 3) => [1 2 3]",
		p.evalVector,
	)
	if err := reg.Register(vectorFunc); err != nil {
		return err
	}

	// vec function (convert list to vector)
	vecFunc := functions.NewFunction(
		"vec",
		registry.CategoryList,
		1,
		"Convert to vector: (vec '(1 2 3)) => [1 2 3]",
		p.evalVec,
	)
	if err := reg.Register(vecFunc); err != nil {
		return err
	}

	// seq function (convert to sequence)
	seqFunc := functions.NewFunction(
		"seq",
		registry.CategoryList,
		1,
		"Convert to sequence: (seq [1 2 3]) => (1 2 3)",
		p.evalSeq,
	)
	return reg.Register(seqFunc)
}

// evalVector creates a vector from arguments
func (p *SequencePlugin) evalVector(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}
	// For now, use ListValue but could create a VectorValue type
	return &types.ListValue{Elements: values}, nil
}

// evalVec converts a list to vector
func (p *SequencePlugin) evalVec(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("vec requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	list, err := functions.ExtractList(value)
	if err != nil {
		return nil, fmt.Errorf("vec: %v", err)
	}

	return &types.ListValue{Elements: list.Elements}, nil
}

// evalSeq converts to sequence
func (p *SequencePlugin) evalSeq(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("seq requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Convert various types to sequences
	switch v := value.(type) {
	case *types.ListValue:
		return v, nil
	default:
		return nil, fmt.Errorf("seq: cannot convert %T to sequence", value)
	}
}
