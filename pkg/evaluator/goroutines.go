package evaluator

import (
	"fmt"
	"runtime/debug"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// evalGo starts a new goroutine to evaluate an expression
func (e *Evaluator) evalGo(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go requires exactly 1 argument")
	}

	expr := args[0]
	future := types.NewFuture()

	// Start goroutine
	go func() {
		// Recover from panics and convert to errors
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("goroutine panic: %v\nStack trace:\n%s", r, debug.Stack())
				future.SetError(err)
			}
		}()

		// Create a new environment that inherits from the current one
		// This ensures goroutines have access to the same variables but can't interfere
		goroutineEnv := NewEnvironmentWithParent(e.env)
		goroutineEvaluator := NewEvaluator(goroutineEnv)

		// Evaluate the expression
		result, err := goroutineEvaluator.Eval(expr)
		if err != nil {
			future.SetError(err)
		} else {
			future.SetResult(result)
		}
	}()

	return future, nil
}

// evalGoWait waits for a goroutine (future) to complete
func (e *Evaluator) evalGoWait(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go-wait requires exactly 1 argument")
	}

	// Evaluate the future argument
	futureValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Check if it's a future
	future, ok := futureValue.(*types.FutureValue)
	if !ok {
		return nil, fmt.Errorf("go-wait requires a future, got %T", futureValue)
	}

	// Wait for the result
	result, waitErr := future.Wait()
	if waitErr != nil {
		return nil, waitErr
	}

	return result, nil
}

// evalGoWaitAll waits for multiple goroutines (futures) to complete
func (e *Evaluator) evalGoWaitAll(args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return &types.ListValue{Elements: []types.Value{}}, nil
	}

	// Collect all futures from arguments
	var futures []*types.FutureValue
	var results []types.Value

	for _, arg := range args {
		// Evaluate the argument
		value, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}

		// Handle different argument types
		switch v := value.(type) {
		case *types.FutureValue:
			// Single future
			futures = append(futures, v)
		case *types.ListValue:
			// List of futures
			for i, elem := range v.Elements {
				if future, ok := elem.(*types.FutureValue); ok {
					futures = append(futures, future)
				} else {
					return nil, fmt.Errorf("go-wait-all: element %d in list is not a future, got %T", i, elem)
				}
			}
		default:
			return nil, fmt.Errorf("go-wait-all requires futures or lists of futures, got %T", value)
		}
	}

	// Wait for all futures
	for _, future := range futures {
		result, waitErr := future.Wait()
		if waitErr != nil {
			return nil, waitErr
		}
		results = append(results, result)
	}

	return &types.ListValue{Elements: results}, nil
}

// evalWaitGroup creates a new wait group
func (e *Evaluator) evalWaitGroup(args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("wait-group requires no arguments")
	}

	return types.NewWaitGroup(), nil
}

// evalWaitGroupAdd adds to a wait group counter
func (e *Evaluator) evalWaitGroupAdd(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wait-group-add! requires exactly 2 arguments")
	}

	// Evaluate wait group
	wgValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	wg, ok := wgValue.(*types.WaitGroupValue)
	if !ok {
		return nil, fmt.Errorf("wait-group-add! requires a wait group, got %T", wgValue)
	}

	// Evaluate delta
	deltaValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	delta, ok := deltaValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("wait-group-add! requires a number, got %T", deltaValue)
	}

	wg.Add(int(delta))
	return wg, nil
}

// evalWaitGroupDone marks a wait group task as done
func (e *Evaluator) evalWaitGroupDone(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wait-group-done! requires exactly 1 argument")
	}

	// Evaluate wait group
	wgValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	wg, ok := wgValue.(*types.WaitGroupValue)
	if !ok {
		return nil, fmt.Errorf("wait-group-done! requires a wait group, got %T", wgValue)
	}

	wg.Done()
	return wg, nil
}

// evalWaitGroupWait waits for all wait group tasks to complete
func (e *Evaluator) evalWaitGroupWait(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wait-group-wait! requires exactly 1 argument")
	}

	// Evaluate wait group
	wgValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	wg, ok := wgValue.(*types.WaitGroupValue)
	if !ok {
		return nil, fmt.Errorf("wait-group-wait! requires a wait group, got %T", wgValue)
	}

	wg.Wait()
	return types.BooleanValue(true), nil
}
