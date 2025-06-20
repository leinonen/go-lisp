package evaluator

import (
	"math"

	"github.com/leinonen/go-lisp/pkg/types"
)

// valuesEqual is a helper function for comparing values in tests
func valuesEqual(a, b types.Value) bool {
	switch va := a.(type) {
	case types.NumberValue:
		if vb, ok := b.(types.NumberValue); ok {
			return math.Abs(float64(va-vb)) < 1e-9
		}
	case types.StringValue:
		if vb, ok := b.(types.StringValue); ok {
			return va == vb
		}
	case types.BooleanValue:
		if vb, ok := b.(types.BooleanValue); ok {
			return va == vb
		}
	case types.KeywordValue:
		if vb, ok := b.(types.KeywordValue); ok {
			return va == vb
		}
	case *types.ListValue:
		if vb, ok := b.(*types.ListValue); ok {
			if len(va.Elements) != len(vb.Elements) {
				return false
			}
			for i, elem := range va.Elements {
				if !valuesEqual(elem, vb.Elements[i]) {
					return false
				}
			}
			return true
		}
	case types.FunctionValue:
		if vb, ok := b.(types.FunctionValue); ok {
			// For functions, compare parameter lists and body string representation
			if len(va.Params) != len(vb.Params) {
				return false
			}
			for i, param := range va.Params {
				if param != vb.Params[i] {
					return false
				}
			}
			return va.Body.String() == vb.Body.String()
		}
	case *types.HashMapValue:
		if vb, ok := b.(*types.HashMapValue); ok {
			if len(va.Elements) != len(vb.Elements) {
				return false
			}
			for key, valueA := range va.Elements {
				valueB, exists := vb.Elements[key]
				if !exists || !valuesEqual(valueA, valueB) {
					return false
				}
			}
			return true
		}
	case *types.NilValue:
		_, ok := b.(*types.NilValue)
		return ok
	case *types.AtomValue:
		if vb, ok := b.(*types.AtomValue); ok {
			// For atoms, compare their current values
			return valuesEqual(va.Value(), vb.Value())
		}
	case nil:
		return b == nil
	}
	return false
}
