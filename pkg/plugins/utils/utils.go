// Package utils provides utility polymorphic functions
package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/interfaces"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// UtilsPlugin provides utility polymorphic functions
type UtilsPlugin struct {
	*plugins.BasePlugin
	evaluator interfaces.CoreEvaluator
}

// NewUtilsPlugin creates a new utils plugin with dependency injection
func NewUtilsPlugin(evaluator interfaces.CoreEvaluator) *UtilsPlugin {
	return &UtilsPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"utils",
			"1.0.0",
			"Utility polymorphic functions (frequencies, group-by, partition, etc.)",
			[]string{}, // No dependencies
		),
		evaluator: evaluator,
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
	// Use the evaluator's CallFunction method if available
	if pureEval, ok := evaluator.(interface {
		CallFunction(types.Value, []types.Expr) (types.Value, error)
	}); ok {
		return pureEval.CallFunction(fnValue, args)
	}

	// Fallback for simple cases
	switch fn := fnValue.(type) {
	case *types.FunctionValue:
		// This is a simplified implementation - in reality we'd need full function evaluation
		if len(args) != len(fn.Params) {
			return nil, fmt.Errorf("function expects %d arguments, got %d", len(fn.Params), len(args))
		}
		// For now, just return the first argument (placeholder)
		if len(args) > 0 {
			return evaluator.Eval(args[0])
		}
		return &types.NilValue{}, nil
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

	// Evaluate the predicate function argument
	predVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating predicate in remove: %v", err)
	}

	// Evaluate the collection argument
	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, fmt.Errorf("error evaluating collection in remove: %v", err)
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("remove: %v", err)
	}

	var result []types.Value
	for _, elem := range elements {
		// Apply predicate to element
		elemExpr := p.valueToExpr(elem)
		predResult, err := p.callFunction(evaluator, predVal, []types.Expr{elemExpr})
		if err != nil {
			return nil, fmt.Errorf("error calling predicate in remove: %v", err)
		}

		// Keep elements where predicate is false (remove those where it's true)
		if !p.isTruthy(predResult) {
			result = append(result, elem)
		}
	}

	return &types.ListValue{Elements: result}, nil
}

// Stub implementations for the remaining functions
func (p *UtilsPlugin) evalKeep(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("keep requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the function argument
	fnVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating function in keep: %v", err)
	}

	// Evaluate the collection argument
	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, fmt.Errorf("error evaluating collection in keep: %v", err)
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("keep: %v", err)
	}

	var result []types.Value
	for _, elem := range elements {
		// Apply function to element
		elemExpr := p.valueToExpr(elem)
		fnResult, err := p.callFunction(evaluator, fnVal, []types.Expr{elemExpr})
		if err != nil {
			return nil, fmt.Errorf("error calling function in keep: %v", err)
		}

		// Keep non-nil results
		if _, isNil := fnResult.(*types.NilValue); !isNil {
			result = append(result, fnResult)
		}
	}

	return &types.ListValue{Elements: result}, nil
}

func (p *UtilsPlugin) evalMapcat(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("mapcat requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the function argument
	fnVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating function in mapcat: %v", err)
	}

	// Evaluate the collection argument
	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, fmt.Errorf("error evaluating collection in mapcat: %v", err)
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("mapcat: %v", err)
	}

	var result []types.Value
	for _, elem := range elements {
		// Apply function to element
		elemExpr := p.valueToExpr(elem)
		fnResult, err := p.callFunction(evaluator, fnVal, []types.Expr{elemExpr})
		if err != nil {
			return nil, fmt.Errorf("error calling function in mapcat: %v", err)
		}

		// Concatenate the result (flatten one level)
		resultElements, err := p.extractSequence(fnResult)
		if err != nil {
			// If it's not a sequence, treat it as a single element
			result = append(result, fnResult)
		} else {
			result = append(result, resultElements...)
		}
	}

	return &types.ListValue{Elements: result}, nil
}

func (p *UtilsPlugin) evalTakeWhile(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("take-while requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the predicate function argument
	predVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating predicate in take-while: %v", err)
	}

	// Evaluate the collection argument
	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, fmt.Errorf("error evaluating collection in take-while: %v", err)
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("take-while: %v", err)
	}

	var result []types.Value
	for _, elem := range elements {
		// Apply predicate to element
		elemExpr := p.valueToExpr(elem)
		predResult, err := p.callFunction(evaluator, predVal, []types.Expr{elemExpr})
		if err != nil {
			return nil, fmt.Errorf("error calling predicate in take-while: %v", err)
		}

		// Check if predicate is truthy
		if !p.isTruthy(predResult) {
			break // Stop taking when predicate becomes false
		}
		result = append(result, elem)
	}

	return &types.ListValue{Elements: result}, nil
}

func (p *UtilsPlugin) evalDropWhile(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("drop-while requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the predicate function argument
	predVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating predicate in drop-while: %v", err)
	}

	// Evaluate the collection argument
	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, fmt.Errorf("error evaluating collection in drop-while: %v", err)
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("drop-while: %v", err)
	}

	// Find the first element where predicate is false
	dropIndex := len(elements) // Default to dropping all if predicate is always true
	for i, elem := range elements {
		// Apply predicate to element
		elemExpr := p.valueToExpr(elem)
		predResult, err := p.callFunction(evaluator, predVal, []types.Expr{elemExpr})
		if err != nil {
			return nil, fmt.Errorf("error calling predicate in drop-while: %v", err)
		}

		// Check if predicate is truthy
		if !p.isTruthy(predResult) {
			dropIndex = i
			break
		}
	}

	// Return the remaining elements
	var result []types.Value
	if dropIndex < len(elements) {
		result = elements[dropIndex:]
	}

	return &types.ListValue{Elements: result}, nil
}

func (p *UtilsPlugin) evalSplitAt(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("split-at requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the index argument
	nValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating index in split-at: %v", err)
	}

	// Evaluate the collection argument
	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, fmt.Errorf("error evaluating collection in split-at: %v", err)
	}

	nNum, ok := nValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("split-at: first argument must be a number, got %T", nValue)
	}
	n := int(nNum)

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("split-at: %v", err)
	}

	// Ensure n is within bounds
	if n < 0 {
		n = 0
	}
	if n > len(elements) {
		n = len(elements)
	}

	// Split the collection
	firstPart := &types.ListValue{Elements: elements[:n]}
	secondPart := &types.ListValue{Elements: elements[n:]}

	// Return as a vector of two parts
	return &types.VectorValue{Elements: []types.Value{firstPart, secondPart}}, nil
}

func (p *UtilsPlugin) evalSplitWith(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("split-with requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the predicate function argument
	predVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating predicate in split-with: %v", err)
	}

	// Evaluate the collection argument
	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, fmt.Errorf("error evaluating collection in split-with: %v", err)
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("split-with: %v", err)
	}

	// Find the first element where predicate is false
	splitIndex := len(elements) // Default to splitting at end if predicate is always true
	for i, elem := range elements {
		// Apply predicate to element
		elemExpr := p.valueToExpr(elem)
		predResult, err := p.callFunction(evaluator, predVal, []types.Expr{elemExpr})
		if err != nil {
			return nil, fmt.Errorf("error calling predicate in split-with: %v", err)
		}

		// Check if predicate is truthy
		if !p.isTruthy(predResult) {
			splitIndex = i
			break
		}
	}

	// Split the collection
	firstPart := &types.ListValue{Elements: elements[:splitIndex]}
	secondPart := &types.ListValue{Elements: elements[splitIndex:]}

	// Return as a vector of two parts
	return &types.VectorValue{Elements: []types.Value{firstPart, secondPart}}, nil
}

func (p *UtilsPlugin) evalComp(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("comp requires at least 1 argument (function), got %d", len(args))
	}

	// Evaluate all function arguments
	functions := make([]types.Value, len(args))
	for i, arg := range args {
		fnVal, err := evaluator.Eval(arg)
		if err != nil {
			return nil, fmt.Errorf("error evaluating function %d in comp: %v", i, err)
		}
		functions[i] = fnVal
	}

	// Create a comp function value
	compFn := &types.CompFunctionValue{
		Functions: functions,
	}

	return compFn, nil
}

func (p *UtilsPlugin) evalPartial(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("partial requires at least 1 argument (function), got %d", len(args))
	}

	// Evaluate the function argument
	fnVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating function in partial: %v", err)
	}

	// Evaluate the partial arguments
	partialArgs := make([]types.Value, 0, len(args)-1)
	for i := 1; i < len(args); i++ {
		arg, err := evaluator.Eval(args[i])
		if err != nil {
			return nil, fmt.Errorf("error evaluating partial argument %d: %v", i, err)
		}
		partialArgs = append(partialArgs, arg)
	}

	// Create a partially applied function value
	partialBuiltin := &types.PartialFunctionValue{
		OriginalFunction: fnVal,
		PartialArgs:      partialArgs,
	}

	return partialBuiltin, nil
}

func (p *UtilsPlugin) evalComplement(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("complement requires exactly 1 argument (predicate function), got %d", len(args))
	}

	// Evaluate the predicate function argument
	predVal, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating predicate in complement: %v", err)
	}

	// Create a complement function value
	complementFn := &types.ComplementFunctionValue{
		PredicateFunction: predVal,
	}

	return complementFn, nil
}

func (p *UtilsPlugin) evalJuxt(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("juxt requires at least 1 argument (function), got %d", len(args))
	}

	// Evaluate all function arguments
	functions := make([]types.Value, len(args))
	for i, arg := range args {
		fnVal, err := evaluator.Eval(arg)
		if err != nil {
			return nil, fmt.Errorf("error evaluating function %d in juxt: %v", i, err)
		}
		functions[i] = fnVal
	}

	// Create a juxt function value
	juxtFn := &types.JuxtFunctionValue{
		Functions: functions,
	}

	return juxtFn, nil
}

func (p *UtilsPlugin) evalUnion(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("union requires at least 1 argument, got %d", len(args))
	}

	// Use a map to track unique elements
	seen := make(map[string]bool)
	var result []types.Value

	for i, arg := range args {
		// Evaluate the collection argument
		collValue, err := evaluator.Eval(arg)
		if err != nil {
			return nil, fmt.Errorf("error evaluating collection %d in union: %v", i, err)
		}

		elements, err := p.extractSequence(collValue)
		if err != nil {
			return nil, fmt.Errorf("union: collection %d: %v", i, err)
		}

		// Add unique elements to result
		for _, elem := range elements {
			key := elem.String()
			if !seen[key] {
				seen[key] = true
				result = append(result, elem)
			}
		}
	}

	return &types.ListValue{Elements: result}, nil
}

func (p *UtilsPlugin) evalIntersection(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("intersection requires at least 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		// Single collection - return it as-is (unique elements)
		collValue, err := evaluator.Eval(args[0])
		if err != nil {
			return nil, fmt.Errorf("error evaluating collection in intersection: %v", err)
		}

		elements, err := p.extractSequence(collValue)
		if err != nil {
			return nil, fmt.Errorf("intersection: %v", err)
		}

		// Return unique elements
		seen := make(map[string]bool)
		var result []types.Value
		for _, elem := range elements {
			key := elem.String()
			if !seen[key] {
				seen[key] = true
				result = append(result, elem)
			}
		}

		return &types.ListValue{Elements: result}, nil
	}

	// Multiple collections - find intersection
	// Start with the first collection
	collValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating first collection in intersection: %v", err)
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("intersection: first collection: %v", err)
	}

	// Create a map of candidates from the first collection
	candidates := make(map[string]types.Value)
	for _, elem := range elements {
		key := elem.String()
		candidates[key] = elem
	}

	// Check each subsequent collection
	for i := 1; i < len(args); i++ {
		collValue, err := evaluator.Eval(args[i])
		if err != nil {
			return nil, fmt.Errorf("error evaluating collection %d in intersection: %v", i, err)
		}

		elements, err := p.extractSequence(collValue)
		if err != nil {
			return nil, fmt.Errorf("intersection: collection %d: %v", i, err)
		}

		// Create a set of elements in this collection
		currentSet := make(map[string]bool)
		for _, elem := range elements {
			key := elem.String()
			currentSet[key] = true
		}

		// Keep only candidates that are in this collection
		newCandidates := make(map[string]types.Value)
		for key, val := range candidates {
			if currentSet[key] {
				newCandidates[key] = val
			}
		}
		candidates = newCandidates
	}

	// Convert map back to list
	var result []types.Value
	for _, val := range candidates {
		result = append(result, val)
	}

	return &types.ListValue{Elements: result}, nil
}

func (p *UtilsPlugin) evalDifference(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("difference requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate the first collection
	coll1Value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, fmt.Errorf("error evaluating first collection in difference: %v", err)
	}

	// Evaluate the second collection
	coll2Value, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, fmt.Errorf("error evaluating second collection in difference: %v", err)
	}

	elements1, err := p.extractSequence(coll1Value)
	if err != nil {
		return nil, fmt.Errorf("difference: first collection: %v", err)
	}

	elements2, err := p.extractSequence(coll2Value)
	if err != nil {
		return nil, fmt.Errorf("difference: second collection: %v", err)
	}

	// Create a set of elements in the second collection
	excludeSet := make(map[string]bool)
	for _, elem := range elements2 {
		key := elem.String()
		excludeSet[key] = true
	}

	// Keep only elements from first collection that are not in second collection
	var result []types.Value
	seen := make(map[string]bool) // To avoid duplicates
	for _, elem := range elements1 {
		key := elem.String()
		if !excludeSet[key] && !seen[key] {
			seen[key] = true
			result = append(result, elem)
		}
	}

	return &types.ListValue{Elements: result}, nil
}

// valueToExpr converts a value back to an expression (helper for function calls)
func (p *UtilsPlugin) valueToExpr(val types.Value) types.Expr {
	switch v := val.(type) {
	case types.NumberValue:
		return &types.NumberExpr{Value: float64(v)}
	case *types.BigNumberValue:
		return &types.BigNumberExpr{Value: v.Value.String()}
	case types.StringValue:
		return &types.StringExpr{Value: string(v)}
	case types.BooleanValue:
		return &types.BooleanExpr{Value: bool(v)}
	case types.KeywordValue:
		return &types.KeywordExpr{Value: string(v)}
	case *types.ListValue:
		exprs := make([]types.Expr, len(v.Elements))
		for i, elem := range v.Elements {
			exprs[i] = p.valueToExpr(elem)
		}
		return &types.ListExpr{Elements: exprs}
	case *types.VectorValue:
		exprs := make([]types.Expr, len(v.Elements))
		for i, elem := range v.Elements {
			exprs[i] = p.valueToExpr(elem)
		}
		return &types.BracketExpr{Elements: exprs}
	default:
		// For complex types, we'll create a symbol that can be resolved later
		return &types.SymbolExpr{Name: fmt.Sprintf("#<value:%s>", val.String())}
	}
}
