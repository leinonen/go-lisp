package evaluator

import (
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// evalDo evaluates multiple expressions in sequence and returns the last result
func (e *Evaluator) evalDo(args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return &types.NilValue{}, nil
	}

	var result types.Value = &types.NilValue{}
	var err error

	// Evaluate each expression in sequence
	for _, expr := range args {
		result, err = e.Eval(expr)
		if err != nil {
			return nil, err
		}
	}

	// Return the result of the last expression
	return result, nil
}
