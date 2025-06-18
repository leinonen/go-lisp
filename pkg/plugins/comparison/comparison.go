// Package comparison provides comparison operations as a plugin
package comparison

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// ComparisonPlugin provides comparison operations
type ComparisonPlugin struct {
	*plugins.BasePlugin
}

// NewComparisonPlugin creates a new comparison plugin
func NewComparisonPlugin() *ComparisonPlugin {
	return &ComparisonPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"comparison",
			"1.0.0",
			"Comparison operations (=, <, >, <=, >=)",
			[]string{}, // No dependencies
		),
	}
}

// RegisterFunctions registers comparison functions
func (cp *ComparisonPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// Equality
	eqFunc := functions.NewFunction(
		"=",
		registry.CategoryComparison,
		-1, // Variadic
		"Test equality: (= 1 1 1) => true, (= 1 2) => false",
		cp.evalEquality,
	)
	if err := reg.Register(eqFunc); err != nil {
		return err
	}

	// Less than
	ltFunc := functions.NewFunction(
		"<",
		registry.CategoryComparison,
		-1, // Variadic
		"Test less than: (< 1 2 3) => true, (< 1 3 2) => false",
		cp.evalLessThan,
	)
	if err := reg.Register(ltFunc); err != nil {
		return err
	}

	// Greater than
	gtFunc := functions.NewFunction(
		">",
		registry.CategoryComparison,
		-1, // Variadic
		"Test greater than: (> 3 2 1) => true, (> 3 1 2) => false",
		cp.evalGreaterThan,
	)
	if err := reg.Register(gtFunc); err != nil {
		return err
	}

	// Less than or equal
	leFunc := functions.NewFunction(
		"<=",
		registry.CategoryComparison,
		-1, // Variadic
		"Test less than or equal: (<= 1 2 2) => true, (<= 2 1) => false",
		cp.evalLessThanOrEqual,
	)
	if err := reg.Register(leFunc); err != nil {
		return err
	}

	// Greater than or equal
	geFunc := functions.NewFunction(
		">=",
		registry.CategoryComparison,
		-1, // Variadic
		"Test greater than or equal: (>= 3 2 2) => true, (>= 2 3) => false",
		cp.evalGreaterThanOrEqual,
	)
	return reg.Register(geFunc)
}

// evalEquality implements equality comparison
func (cp *ComparisonPlugin) evalEquality(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("= requires at least 2 arguments")
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	// Compare all values to the first one
	first := values[0]
	for i := 1; i < len(values); i++ {
		if !cp.valuesEqual(first, values[i]) {
			return types.BooleanValue(false), nil
		}
	}
	return types.BooleanValue(true), nil
}

// evalLessThan implements less than comparison
func (cp *ComparisonPlugin) evalLessThan(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("< requires at least 2 arguments")
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	// Check that each value is less than the next
	for i := 0; i < len(values)-1; i++ {
		if !cp.valueLessThan(values[i], values[i+1]) {
			return types.BooleanValue(false), nil
		}
	}
	return types.BooleanValue(true), nil
}

// evalGreaterThan implements greater than comparison
func (cp *ComparisonPlugin) evalGreaterThan(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("> requires at least 2 arguments")
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	// Check that each value is greater than the next
	for i := 0; i < len(values)-1; i++ {
		if !cp.valueGreaterThan(values[i], values[i+1]) {
			return types.BooleanValue(false), nil
		}
	}
	return types.BooleanValue(true), nil
}

// evalLessThanOrEqual implements less than or equal comparison
func (cp *ComparisonPlugin) evalLessThanOrEqual(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("<= requires at least 2 arguments")
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	// Check that each value is less than or equal to the next
	for i := 0; i < len(values)-1; i++ {
		if cp.valueGreaterThan(values[i], values[i+1]) {
			return types.BooleanValue(false), nil
		}
	}
	return types.BooleanValue(true), nil
}

// evalGreaterThanOrEqual implements greater than or equal comparison
func (cp *ComparisonPlugin) evalGreaterThanOrEqual(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf(">= requires at least 2 arguments")
	}

	values, err := functions.EvalArgs(evaluator, args)
	if err != nil {
		return nil, err
	}

	// Check that each value is greater than or equal to the next
	for i := 0; i < len(values)-1; i++ {
		if cp.valueLessThan(values[i], values[i+1]) {
			return types.BooleanValue(false), nil
		}
	}
	return types.BooleanValue(true), nil
}

// Helper functions for value comparison

// valuesEqual compares two values for equality
func (cp *ComparisonPlugin) valuesEqual(a, b types.Value) bool {
	// Handle same type comparisons
	switch va := a.(type) {
	case types.NumberValue:
		if vb, ok := b.(types.NumberValue); ok {
			return float64(va) == float64(vb)
		}
		if vb, ok := b.(*types.BigNumberValue); ok {
			// Convert regular number to big number for comparison
			bigA := types.NewBigNumberFromInt64(int64(va))
			return bigA.Value.Cmp(vb.Value) == 0
		}
	case *types.BigNumberValue:
		if vb, ok := b.(*types.BigNumberValue); ok {
			return va.Value.Cmp(vb.Value) == 0
		}
		if vb, ok := b.(types.NumberValue); ok {
			// Convert regular number to big number for comparison
			bigB := types.NewBigNumberFromInt64(int64(vb))
			return va.Value.Cmp(bigB.Value) == 0
		}
	case types.StringValue:
		if vb, ok := b.(types.StringValue); ok {
			return string(va) == string(vb)
		}
	case types.BooleanValue:
		if vb, ok := b.(types.BooleanValue); ok {
			return bool(va) == bool(vb)
		}
	case types.KeywordValue:
		if vb, ok := b.(types.KeywordValue); ok {
			return string(va) == string(vb)
		}
	}
	return false
}

// valueLessThan compares if a < b for numeric values
func (cp *ComparisonPlugin) valueLessThan(a, b types.Value) bool {
	aNum, aErr := functions.ExtractFloat64(a)
	bNum, bErr := functions.ExtractFloat64(b)

	if aErr != nil || bErr != nil {
		return false // Can't compare non-numeric values
	}

	return aNum < bNum
}

// valueGreaterThan compares if a > b for numeric values
func (cp *ComparisonPlugin) valueGreaterThan(a, b types.Value) bool {
	aNum, aErr := functions.ExtractFloat64(a)
	bNum, bErr := functions.ExtractFloat64(b)

	if aErr != nil || bErr != nil {
		return false // Can't compare non-numeric values
	}

	return aNum > bNum
}
