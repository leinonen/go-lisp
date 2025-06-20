// Package polymorphic provides polymorphic functions that work across multiple data types
package polymorphic

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// PolymorphicPlugin provides polymorphic functions that work across data types
type PolymorphicPlugin struct {
	*plugins.BasePlugin
}

// NewPolymorphicPlugin creates a new polymorphic plugin
func NewPolymorphicPlugin() *PolymorphicPlugin {
	return &PolymorphicPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"polymorphic",
			"1.0.0",
			"Polymorphic functions that work across multiple data types",
			[]string{}, // No dependencies
		),
	}
}

// RegisterFunctions registers polymorphic functions
func (p *PolymorphicPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	funcs := []struct {
		name    string
		arity   int
		help    string
		handler func(registry.Evaluator, []types.Expr) (types.Value, error)
	}{
		// Sequence functions
		{"first", 1, "Get first element of sequence: (first coll) - works on lists, vectors, strings", p.evalFirst},
		{"rest", 1, "Get rest of sequence: (rest coll) - works on lists, vectors, strings", p.evalRest},
		{"last", 1, "Get last element of sequence: (last coll) - works on lists, vectors, strings", p.evalLast},
		{"nth", 2, "Get nth element of sequence: (nth coll n) - works on lists, vectors, strings", p.evalNth},
		{"second", 1, "Get second element of sequence: (second coll) - works on lists, vectors, strings", p.evalSecond},
		{"empty?", 1, "Check if collection is empty: (empty? coll) - works on all collections", p.evalEmpty},
		{"seq", 1, "Convert to sequence: (seq coll) - works on lists, vectors, strings", p.evalSeq},

		// Collection functions
		{"take", 2, "Take first n elements: (take n coll) - works on all sequences", p.evalTake},
		{"drop", 2, "Drop first n elements: (drop n coll) - works on all sequences", p.evalDrop},
		{"reverse", 1, "Reverse sequence: (reverse coll) - works on lists, vectors, strings", p.evalReverse},
		{"distinct", 1, "Remove duplicates: (distinct coll) - works on all sequences", p.evalDistinct},
		{"sort", 1, "Sort sequence: (sort coll) - works on all sequences", p.evalSort},
		{"into", 2, "Merge collections: (into to from) - works on all collections", p.evalInto},

		// Predicate functions
		{"seq?", 1, "Check if sequential: (seq? x)", p.evalSeqPredicate},
		{"coll?", 1, "Check if collection: (coll? x)", p.evalCollPredicate},
		{"sequential?", 1, "Check if sequential: (sequential? x)", p.evalSequentialPredicate},
		{"indexed?", 1, "Check if indexed: (indexed? x)", p.evalIndexedPredicate},

		// Utility functions
		{"identity", 1, "Return argument unchanged: (identity x)", p.evalIdentity},
		{"constantly", 1, "Return constant function: (constantly x)", p.evalConstantly},
	}

	for _, fn := range funcs {
		f := functions.NewFunction(fn.name, registry.CategoryCore, fn.arity, fn.help, fn.handler)
		if err := reg.Register(f); err != nil {
			return fmt.Errorf("failed to register %s: %v", fn.name, err)
		}
	}

	return nil
}

// Helper function to extract sequence elements from various types
func (p *PolymorphicPlugin) extractSequence(value types.Value) ([]types.Value, error) {
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

// Helper function to convert elements back to the same type as input
func (p *PolymorphicPlugin) createSameType(input types.Value, elements []types.Value) types.Value {
	switch input.(type) {
	case *types.ListValue:
		return &types.ListValue{Elements: elements}
	case *types.VectorValue:
		return types.NewVectorValue(elements)
	case types.StringValue:
		// Convert back to string
		var result string
		for _, elem := range elements {
			if str, ok := elem.(types.StringValue); ok {
				result += string(str)
			}
		}
		return types.StringValue(result)
	default:
		// Default to list
		return &types.ListValue{Elements: elements}
	}
}

// evalFirst returns the first element of any sequence
func (p *PolymorphicPlugin) evalFirst(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("first requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("first: %v", err)
	}

	if len(elements) == 0 {
		return &types.NilValue{}, nil
	}

	return elements[0], nil
}

// evalRest returns all elements except the first
func (p *PolymorphicPlugin) evalRest(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("rest requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("rest: %v", err)
	}

	if len(elements) <= 1 {
		return &types.ListValue{Elements: []types.Value{}}, nil
	}

	// Always return a list for rest (Clojure behavior)
	return &types.ListValue{Elements: elements[1:]}, nil
}

// evalLast returns the last element of any sequence
func (p *PolymorphicPlugin) evalLast(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("last requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("last: %v", err)
	}

	if len(elements) == 0 {
		return &types.NilValue{}, nil
	}

	return elements[len(elements)-1], nil
}

// evalNth returns the nth element of any indexed sequence
func (p *PolymorphicPlugin) evalNth(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("nth requires exactly 2 arguments, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	indexValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Extract index
	indexNum, ok := indexValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("nth: index must be a number, got %T", indexValue)
	}
	index := int(indexNum)

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("nth: %v", err)
	}

	if index < 0 || index >= len(elements) {
		return nil, fmt.Errorf("nth: index %d out of bounds for sequence of length %d", index, len(elements))
	}

	return elements[index], nil
}

// evalSecond returns the second element of any sequence
func (p *PolymorphicPlugin) evalSecond(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("second requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("second: %v", err)
	}

	if len(elements) < 2 {
		return &types.NilValue{}, nil
	}

	return elements[1], nil
}

// evalEmpty checks if any collection is empty
func (p *PolymorphicPlugin) evalEmpty(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("empty? requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	switch v := value.(type) {
	case *types.ListValue:
		return types.BooleanValue(len(v.Elements) == 0), nil
	case *types.VectorValue:
		return types.BooleanValue(len(v.Elements) == 0), nil
	case *types.HashMapValue:
		return types.BooleanValue(len(v.Elements) == 0), nil
	case types.StringValue:
		return types.BooleanValue(len(string(v)) == 0), nil
	case *types.NilValue:
		return types.BooleanValue(true), nil
	default:
		return nil, fmt.Errorf("empty?: not a collection: %T", value)
	}
}

// evalSeq converts any collection to a sequence (list)
func (p *PolymorphicPlugin) evalSeq(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("seq requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("seq: %v", err)
	}

	if len(elements) == 0 {
		return &types.NilValue{}, nil
	}

	return &types.ListValue{Elements: elements}, nil
}

// evalTake takes the first n elements from any sequence
func (p *PolymorphicPlugin) evalTake(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("take requires exactly 2 arguments, got %d", len(args))
	}

	nValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Extract n
	nNum, ok := nValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("take: first argument must be a number, got %T", nValue)
	}
	n := int(nNum)

	if n <= 0 {
		return &types.ListValue{Elements: []types.Value{}}, nil
	}

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("take: %v", err)
	}

	if n > len(elements) {
		n = len(elements)
	}

	return &types.ListValue{Elements: elements[:n]}, nil
}

// evalDrop drops the first n elements from any sequence
func (p *PolymorphicPlugin) evalDrop(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("drop requires exactly 2 arguments, got %d", len(args))
	}

	nValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	collValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Extract n
	nNum, ok := nValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("drop: first argument must be a number, got %T", nValue)
	}
	n := int(nNum)

	elements, err := p.extractSequence(collValue)
	if err != nil {
		return nil, fmt.Errorf("drop: %v", err)
	}

	if n <= 0 {
		return &types.ListValue{Elements: elements}, nil
	}

	if n >= len(elements) {
		return &types.ListValue{Elements: []types.Value{}}, nil
	}

	return &types.ListValue{Elements: elements[n:]}, nil
}

// evalReverse reverses any sequence
func (p *PolymorphicPlugin) evalReverse(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("reverse requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("reverse: %v", err)
	}

	reversed := make([]types.Value, len(elements))
	for i, elem := range elements {
		reversed[len(elements)-1-i] = elem
	}

	return p.createSameType(value, reversed), nil
}

// evalDistinct removes duplicates from any sequence
func (p *PolymorphicPlugin) evalDistinct(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("distinct requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("distinct: %v", err)
	}

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

// evalSort sorts any sequence
func (p *PolymorphicPlugin) evalSort(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sort requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	elements, err := p.extractSequence(value)
	if err != nil {
		return nil, fmt.Errorf("sort: %v", err)
	}

	// Create a copy to sort
	sorted := make([]types.Value, len(elements))
	copy(sorted, elements)

	// Simple string-based sorting
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].String() > sorted[j].String() {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return &types.ListValue{Elements: sorted}, nil
}

// evalInto merges one collection into another
func (p *PolymorphicPlugin) evalInto(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("into requires exactly 2 arguments, got %d", len(args))
	}

	toValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	fromValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	toElements, err := p.extractSequence(toValue)
	if err != nil {
		return nil, fmt.Errorf("into: first argument must be a collection, got %T", toValue)
	}

	fromElements, err := p.extractSequence(fromValue)
	if err != nil {
		return nil, fmt.Errorf("into: second argument must be a collection, got %T", fromValue)
	}

	// Merge elements
	result := make([]types.Value, len(toElements)+len(fromElements))
	copy(result, toElements)
	copy(result[len(toElements):], fromElements)

	return p.createSameType(toValue, result), nil
}

// Predicate functions
func (p *PolymorphicPlugin) evalSeqPredicate(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("seq? requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	switch value.(type) {
	case *types.ListValue, *types.VectorValue, types.StringValue:
		return types.BooleanValue(true), nil
	default:
		return types.BooleanValue(false), nil
	}
}

func (p *PolymorphicPlugin) evalCollPredicate(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("coll? requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	switch value.(type) {
	case *types.ListValue, *types.VectorValue, *types.HashMapValue:
		return types.BooleanValue(true), nil
	default:
		return types.BooleanValue(false), nil
	}
}

func (p *PolymorphicPlugin) evalSequentialPredicate(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sequential? requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	switch value.(type) {
	case *types.ListValue, *types.VectorValue:
		return types.BooleanValue(true), nil
	default:
		return types.BooleanValue(false), nil
	}
}

func (p *PolymorphicPlugin) evalIndexedPredicate(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("indexed? requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	switch value.(type) {
	case *types.VectorValue, types.StringValue:
		return types.BooleanValue(true), nil
	default:
		return types.BooleanValue(false), nil
	}
}

// Utility functions
func (p *PolymorphicPlugin) evalIdentity(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("identity requires exactly 1 argument, got %d", len(args))
	}

	return evaluator.Eval(args[0])
}

func (p *PolymorphicPlugin) evalConstantly(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("constantly requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Return a function that always returns the given value
	// For now, return a placeholder - this needs proper function implementation
	return value, nil
}
