// Package sequence provides vector and sequence operations for the Lisp interpreter
package sequence

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
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
	if err := reg.Register(seqFunc); err != nil {
		return err
	}

	// vector? predicate
	vectorPredFunc := functions.NewFunction(
		"vector?",
		registry.CategoryList,
		1,
		"Check if value is a vector: (vector? [1 2 3]) => true",
		p.evalVectorPredicate,
	)
	if err := reg.Register(vectorPredFunc); err != nil {
		return err
	}

	// conj for vectors (add to end)
	conjFunc := functions.NewFunction(
		"conj",
		registry.CategoryList,
		-1, // Variadic: (conj coll & items)
		"Add elements to collection: (conj [1 2] 3 4) => [1 2 3 4]",
		p.evalConj,
	)
	if err := reg.Register(conjFunc); err != nil {
		return err
	}

	// count for vectors
	countFunc := functions.NewFunction(
		"count",
		registry.CategoryList,
		1,
		"Get count of elements: (count [1 2 3]) => 3",
		p.evalCount,
	)
	return reg.Register(countFunc)
}

// evalVector creates a vector from arguments
func (p *SequencePlugin) evalVector(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}
	// Create a VectorValue
	return types.NewVectorValue(values), nil
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

	// Handle different input types
	switch v := value.(type) {
	case *types.ListValue:
		return types.NewVectorValue(v.Elements), nil
	case *types.VectorValue:
		// Already a vector, return as is
		return v, nil
	default:
		return nil, fmt.Errorf("vec: cannot convert %T to vector", value)
	}
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

	// Convert various types to sequences (lists)
	switch v := value.(type) {
	case *types.ListValue:
		return v, nil
	case *types.VectorValue:
		// Convert vector to list
		return &types.ListValue{Elements: v.Elements}, nil
	default:
		return nil, fmt.Errorf("seq: cannot convert %T to sequence", value)
	}
}

// evalVectorPredicate checks if value is a vector
func (p *SequencePlugin) evalVectorPredicate(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("vector? requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	if _, ok := value.(*types.VectorValue); ok {
		return types.BooleanValue(true), nil
	}
	return types.BooleanValue(false), nil
}

// evalConj adds elements to a collection
func (p *SequencePlugin) evalConj(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("conj requires at least 1 argument, got %d", len(args))
	}

	coll, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Evaluate the items to add
	var items []types.Value
	for _, arg := range args[1:] {
		item, err := evaluator.Eval(arg)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	switch c := coll.(type) {
	case *types.VectorValue:
		// For vectors, add to the end
		newElements := make([]types.Value, len(c.Elements)+len(items))
		copy(newElements, c.Elements)
		copy(newElements[len(c.Elements):], items)
		return types.NewVectorValue(newElements), nil
	case *types.ListValue:
		// For lists, add to the front (typical Clojure behavior)
		newElements := make([]types.Value, len(items)+len(c.Elements))
		copy(newElements, items)
		copy(newElements[len(items):], c.Elements)
		return &types.ListValue{Elements: newElements}, nil
	default:
		return nil, fmt.Errorf("conj: cannot add to %T", coll)
	}
}

// evalCount gets count of elements
func (p *SequencePlugin) evalCount(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("count requires exactly 1 argument, got %d", len(args))
	}

	coll, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	switch c := coll.(type) {
	case *types.VectorValue:
		return types.NumberValue(len(c.Elements)), nil
	case *types.ListValue:
		return types.NumberValue(len(c.Elements)), nil
	case *types.HashMapValue:
		return types.NumberValue(len(c.Elements)), nil
	case types.StringValue:
		return types.NumberValue(len(string(c))), nil
	default:
		return nil, fmt.Errorf("count: cannot count %T", coll)
	}
}
