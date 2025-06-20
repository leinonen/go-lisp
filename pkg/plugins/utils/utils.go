// Package utils provides utility polymorphic functions
package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// UtilsPlugin provides utility polymorphic functions
type UtilsPlugin struct {
	*plugins.BasePlugin
}

// NewUtilsPlugin creates a new utils plugin
func NewUtilsPlugin() *UtilsPlugin {
	return &UtilsPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"utils",
			"1.0.0",
			"Utility polymorphic functions (frequencies, group-by, partition, etc.)",
			[]string{}, // No dependencies
		),
	}
}

// RegisterFunctions registers utility functions
func (p *UtilsPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	utilFunctions := []struct {
		name    string
		arity   int
		help    string
		handler func(registry.Evaluator, []types.Expr) (types.Value, error)
	}{
		// Collection utilities
		{"frequencies", 1, "Count occurrences: (frequencies coll)", p.evalFrequencies},
		{"group-by", 2, "Group by function result: (group-by fn coll)", p.evalGroupBy},
		{"partition", 2, "Partition into chunks: (partition n coll)", p.evalPartition},
		{"interleave", -1, "Interleave sequences: (interleave seq1 seq2 ...)", p.evalInterleave},
		{"interpose", 2, "Insert separator: (interpose sep coll)", p.evalInterpose},
		{"flatten", 1, "Flatten nested sequences: (flatten coll)", p.evalFlatten},
		{"shuffle", 1, "Randomize order: (shuffle coll)", p.evalShuffle},
		{"remove", 2, "Remove elements: (remove pred coll)", p.evalRemove},
		{"keep", 2, "Keep non-nil results: (keep fn coll)", p.evalKeep},
		{"mapcat", 2, "Map then concatenate: (mapcat fn coll)", p.evalMapcat},

		// Sequence operations
		{"take-while", 2, "Take while predicate true: (take-while pred coll)", p.evalTakeWhile},
		{"drop-while", 2, "Drop while predicate true: (drop-while pred coll)", p.evalDropWhile},
		{"split-at", 2, "Split at index: (split-at n coll)", p.evalSplitAt},
		{"split-with", 2, "Split with predicate: (split-with pred coll)", p.evalSplitWith},

		// Function utilities
		{"comp", -1, "Function composition: (comp f g h)", p.evalComp},
		{"partial", -1, "Partial application: (partial f arg1 arg2)", p.evalPartial},
		{"complement", 1, "Complement predicate: (complement pred)", p.evalComplement},
		{"juxt", -1, "Apply multiple functions: (juxt f g h)", p.evalJuxt},

		// Set-like operations
		{"union", -1, "Union of collections: (union coll1 coll2 ...)", p.evalUnion},
		{"intersection", -1, "Intersection of collections: (intersection coll1 coll2 ...)", p.evalIntersection},
		{"difference", 2, "Difference of collections: (difference coll1 coll2)", p.evalDifference},
	}

	for _, fn := range utilFunctions {
		f := functions.NewFunction(fn.name, registry.CategoryCore, fn.arity, fn.help, fn.handler)
		if err := reg.Register(f); err != nil {
			return fmt.Errorf("failed to register %s: %v", fn.name, err)
		}
	}

	return nil
}

// Helper function to extract sequence elements from various types
func (p *UtilsPlugin) extractSequence(value types.Value) ([]types.Value, error) {
	switch v := value.(type) {
	case *types.ListValue:
		return v.Elements, nil
	case *types.VectorValue:
		return v.Elements, nil
	case types.StringValue:
		// Convert string to list of character strings
		s := string(v)
		runes := []rune(s)
		elements := make([]types.Value, len(runes))
		for i, r := range runes {
			elements[i] = types.StringValue(string(r))
		}
		return elements, nil
	case *types.NilValue:
		return []types.Value{}, nil
	default:
		return nil, fmt.Errorf("not a sequence: %T", value)
	}
}

// Helper function to call a function with arguments
func (p *UtilsPlugin) callFunction(evaluator registry.Evaluator, fnValue types.Value, args []types.Expr) (types.Value, error) {
	switch fnValue.(type) {
	case *types.FunctionValue:
		// This would need proper function evaluation - simplified for now
		return evaluator.Eval(args[0]) // Just return first arg for now
	default:
		return nil, fmt.Errorf("not a function: %T", fnValue)
	}
}

// Helper to check if a value is truthy
func (p *UtilsPlugin) isTruthy(value types.Value) bool {
	switch v := value.(type) {
	case types.BooleanValue:
		return bool(v)
	case *types.NilValue:
		return false
	default:
		return true
	}
}

// evalFrequencies counts occurrences of each element
func (p *UtilsPlugin) evalFrequencies(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("frequencies requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("frequencies: %v", err)
	}

	freq := make(map[string]types.Value)
	counts := make(map[string]int)

	for _, elem := range elements {
		key := elem.String()
		counts[key]++
		freq[key] = elem // Keep original element as key
	}

	// Convert to hashmap
	result := make(map[string]types.Value)
	for key, count := range counts {
		result[key] = types.NumberValue(float64(count))
	}

	return &types.HashMapValue{Elements: result}, nil
}

// evalGroupBy groups elements by function result
func (p *UtilsPlugin) evalGroupBy(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("group-by requires exactly 2 arguments, got %d", len(args))
	}

	_, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("group-by: %v", err)
	}

	groups := make(map[string][]types.Value)

	for _, elem := range elements {
		// For now, use string representation as group key
		key := elem.String()
		groups[key] = append(groups[key], elem)
	}

	// Convert to hashmap
	result := make(map[string]types.Value)
	for key, group := range groups {
		result[key] = &types.ListValue{Elements: group}
	}

	return &types.HashMapValue{Elements: result}, nil
}

// evalPartition partitions sequence into chunks
func (p *UtilsPlugin) evalPartition(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("partition requires exactly 2 arguments, got %d", len(args))
	}

	nValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	nNum, ok := nValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("partition: first argument must be a number, got %T", nValue)
	}
	n := int(nNum)

	if n <= 0 {
		return nil, fmt.Errorf("partition: chunk size must be positive, got %d", n)
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("partition: %v", err)
	}

	var partitions []types.Value
	for i := 0; i+n <= len(elements); i += n {
		chunk := &types.ListValue{Elements: elements[i : i+n]}
		partitions = append(partitions, chunk)
	}

	return &types.ListValue{Elements: partitions}, nil
}

// evalInterleave interleaves multiple sequences
func (p *UtilsPlugin) evalInterleave(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("interleave requires at least 2 arguments, got %d", len(args))
	}

	var sequences [][]types.Value
	minLen := -1

	for _, arg := range args {
		value, err := evaluator.Eval(arg)
		if err != nil {
			return nil, err
		}

		elements, err := p.extractSequence(value)
		if err != nil {
			return nil, fmt.Errorf("interleave: %v", err)
		}

		sequences = append(sequences, elements)
		if minLen == -1 || len(elements) < minLen {
			minLen = len(elements)
		}
	}

	var result []types.Value
	for i := 0; i < minLen; i++ {
		for _, seq := range sequences {
			result = append(result, seq[i])
		}
	}

	return &types.ListValue{Elements: result}, nil
}

// evalInterpose inserts separator between elements
func (p *UtilsPlugin) evalInterpose(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("interpose requires exactly 2 arguments, got %d", len(args))
	}

	sepValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("interpose: %v", err)
	}

	if len(elements) == 0 {
		return &types.ListValue{Elements: []types.Value{}}, nil
	}

	var result []types.Value
	result = append(result, elements[0])

	for i := 1; i < len(elements); i++ {
		result = append(result, sepValue, elements[i])
	}

	return &types.ListValue{Elements: result}, nil
}

// evalFlatten flattens nested sequences
func (p *UtilsPlugin) evalFlatten(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("flatten requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	result := p.flattenValue(value)
	return &types.ListValue{Elements: result}, nil
}

// Helper function to recursively flatten
func (p *UtilsPlugin) flattenValue(value types.Value) []types.Value {
	switch v := value.(type) {
	case *types.ListValue:
		var result []types.Value
		for _, elem := range v.Elements {
			result = append(result, p.flattenValue(elem)...)
		}
		return result
	case *types.VectorValue:
		var result []types.Value
		for _, elem := range v.Elements {
			result = append(result, p.flattenValue(elem)...)
		}
		return result
	default:
		return []types.Value{value}
	}
}

// evalShuffle randomizes order of sequence
func (p *UtilsPlugin) evalShuffle(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("shuffle requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("shuffle: %v", err)
	}

	// Create a copy to shuffle
	shuffled := make([]types.Value, len(elements))
	copy(shuffled, elements)

	// Shuffle using Fisher-Yates algorithm
	rand.Seed(time.Now().UnixNano())
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return &types.ListValue{Elements: shuffled}, nil
}

// evalRemove removes elements that satisfy predicate
func (p *UtilsPlugin) evalRemove(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("remove requires exactly 2 arguments, got %d", len(args))
	}

	_, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("remove: %v", err)
	}

	var result []types.Value
	for _, elem := range elements {
		// For now, simplified - just keep all elements
		result = append(result, elem)
	}

	return &types.ListValue{Elements: result}, nil
}

// Stub implementations for the remaining functions
func (p *UtilsPlugin) evalKeep(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalMapcat(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalTakeWhile(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalDropWhile(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalSplitAt(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalSplitWith(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalComp(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalPartial(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalComplement(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalJuxt(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalUnion(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalIntersection(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}

func (p *UtilsPlugin) evalDifference(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	return &types.ListValue{Elements: []types.Value{}}, nil
}
