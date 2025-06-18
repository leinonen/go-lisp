// Package concurrency provides concurrency functionality for the Lisp interpreter
package concurrency

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// InterpreterDependency is an optional interface that the interpreter can implement
// to provide additional context to plugins
type InterpreterDependency interface {
	Interpret(input string) (types.Value, error)
	GetEnvironment() interface{} // Returns environment for context sharing
}

// ConcurrencyPlugin provides concurrency functions
type ConcurrencyPlugin struct {
	*plugins.BasePlugin
	interpreter InterpreterDependency // Optional interpreter reference
}

// NewConcurrencyPlugin creates a new concurrency plugin
func NewConcurrencyPlugin() *ConcurrencyPlugin {
	return &ConcurrencyPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"concurrency",
			"1.0.0",
			"Concurrency operations (goroutines, channels)",
			[]string{}, // No dependencies
		),
		interpreter: nil, // Will be set later if needed
	}
}

// SetInterpreter sets an optional interpreter dependency
// This allows the plugin to access full interpreter functionality if needed
func (p *ConcurrencyPlugin) SetInterpreter(interp InterpreterDependency) {
	p.interpreter = interp
}

// HasInterpreter returns true if an interpreter dependency is available
func (p *ConcurrencyPlugin) HasInterpreter() bool {
	return p.interpreter != nil
}

// RegisterFunctions registers concurrency functions
func (p *ConcurrencyPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// Go function
	goFunc := functions.NewFunction(
		"go",
		registry.CategoryConcurrency,
		1,
		"Start goroutine: (go expression) => future",
		p.goFunc,
	)
	if err := reg.Register(goFunc); err != nil {
		return err
	}

	// Go wait function
	goWaitFunc := functions.NewFunction(
		"go-wait",
		registry.CategoryConcurrency,
		1,
		"Wait for goroutine result: (go-wait future) => result",
		p.goWait,
	)
	if err := reg.Register(goWaitFunc); err != nil {
		return err
	}

	// Go wait all function
	goWaitAllFunc := functions.NewFunction(
		"go-wait-all",
		registry.CategoryConcurrency,
		-1, // Variadic
		"Wait for multiple goroutines: (go-wait-all future1 future2) => (result1 result2)",
		p.goWaitAll,
	)
	if err := reg.Register(goWaitAllFunc); err != nil {
		return err
	}

	// Channel creation function
	chanFunc := functions.NewFunction(
		"chan",
		registry.CategoryConcurrency,
		-1, // 0 or 1 args
		"Create channel: (chan) => unbuffered, (chan 5) => buffered",
		p.chan_,
	)
	if err := reg.Register(chanFunc); err != nil {
		return err
	}

	// Channel send function
	chanSendFunc := functions.NewFunction(
		"chan-send!",
		registry.CategoryConcurrency,
		2,
		"Send to channel: (chan-send! channel value) => true",
		p.chanSend,
	)
	if err := reg.Register(chanSendFunc); err != nil {
		return err
	}

	// Channel receive function
	chanRecvFunc := functions.NewFunction(
		"chan-recv!",
		registry.CategoryConcurrency,
		1,
		"Receive from channel (blocking): (chan-recv! channel) => value",
		p.chanRecv,
	)
	if err := reg.Register(chanRecvFunc); err != nil {
		return err
	}

	// Channel try receive function
	chanTryRecvFunc := functions.NewFunction(
		"chan-try-recv!",
		registry.CategoryConcurrency,
		1,
		"Try receive from channel (non-blocking): (chan-try-recv! channel) => value or nil",
		p.chanTryRecv,
	)
	if err := reg.Register(chanTryRecvFunc); err != nil {
		return err
	}

	// Channel close function
	chanCloseFunc := functions.NewFunction(
		"chan-close!",
		registry.CategoryConcurrency,
		1,
		"Close channel: (chan-close! channel) => true",
		p.chanClose,
	)
	if err := reg.Register(chanCloseFunc); err != nil {
		return err
	}

	// Channel closed check function
	chanClosedFunc := functions.NewFunction(
		"chan-closed?",
		registry.CategoryConcurrency,
		1,
		"Check if channel is closed: (chan-closed? channel) => true/false",
		p.chanClosed,
	)
	if err := reg.Register(chanClosedFunc); err != nil {
		return err
	}

	return nil
}

// Go function - starts a new goroutine to evaluate an expression
func (p *ConcurrencyPlugin) goFunc(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go requires exactly 1 argument, got %d", len(args))
	}

	// Create a future to hold the result
	future := types.NewFuture()

	// Start goroutine to evaluate the expression
	go func() {
		// Recover from panics and convert to errors
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("goroutine panic: %v", r)
				future.SetError(err)
			}
		}()

		// Evaluate the expression in the goroutine
		result, err := evaluator.Eval(args[0])
		if err != nil {
			future.SetError(err)
		} else {
			future.SetResult(result)
		}
	}()

	return future, nil
}

// Go-wait waits for a goroutine (future) to complete
func (p *ConcurrencyPlugin) goWait(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go-wait requires exactly 1 argument, got %d", len(args))
	}

	// Evaluate the future argument
	futureValue, err := evaluator.Eval(args[0])
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

// Go-wait-all waits for multiple goroutines (futures) to complete
func (p *ConcurrencyPlugin) goWaitAll(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return &types.ListValue{Elements: []types.Value{}}, nil
	}

	// Collect all futures from arguments
	var futures []*types.FutureValue
	var results []types.Value

	for _, arg := range args {
		// Evaluate the argument
		value, err := evaluator.Eval(arg)
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
					return nil, fmt.Errorf("go-wait-all list element %d is not a future, got %T", i, elem)
				}
			}
		default:
			return nil, fmt.Errorf("go-wait-all argument must be a future or list of futures, got %T", v)
		}
	}

	// Wait for all futures and collect results
	for _, future := range futures {
		result, err := future.Wait()
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return &types.ListValue{Elements: results}, nil
}

// Chan creates a new channel
func (p *ConcurrencyPlugin) chan_(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	size := 0 // Default to unbuffered channel

	if len(args) == 1 {
		// Evaluate the buffer size
		sizeValue, err := evaluator.Eval(args[0])
		if err != nil {
			return nil, err
		}

		sizeNum, ok := sizeValue.(types.NumberValue)
		if !ok {
			return nil, fmt.Errorf("chan buffer size must be a number, got %T", sizeValue)
		}

		if sizeNum < 0 {
			return nil, fmt.Errorf("chan buffer size cannot be negative: %v", sizeNum)
		}

		size = int(sizeNum)
	} else if len(args) > 1 {
		return nil, fmt.Errorf("chan requires 0 or 1 arguments, got %d", len(args))
	}

	return types.NewChannel(size), nil
}

// Chan-send! sends a value to a channel
func (p *ConcurrencyPlugin) chanSend(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("chan-send! requires exactly 2 arguments, got %d", len(args))
	}

	// Evaluate channel argument
	channelValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	channel, ok := channelValue.(*types.ChannelValue)
	if !ok {
		return nil, fmt.Errorf("chan-send! requires a channel, got %T", channelValue)
	}

	// Evaluate value argument
	value, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Send the value
	err = channel.Send(value)
	if err != nil {
		return nil, err
	}

	return types.BooleanValue(true), nil
}

// Chan-recv! receives a value from a channel (blocking)
func (p *ConcurrencyPlugin) chanRecv(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("chan-recv! requires exactly 1 argument, got %d", len(args))
	}

	// Evaluate channel argument
	channelValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	channel, ok := channelValue.(*types.ChannelValue)
	if !ok {
		return nil, fmt.Errorf("chan-recv! requires a channel, got %T", channelValue)
	}

	// Receive the value
	value, ok := channel.Receive()
	if !ok {
		// Channel was closed and no more values
		return &types.NilValue{}, nil
	}

	return value, nil
}

// Chan-try-recv! tries to receive a value from a channel (non-blocking)
func (p *ConcurrencyPlugin) chanTryRecv(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	// Evaluate channel argument
	channelValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	channel, ok := channelValue.(*types.ChannelValue)
	if !ok {
		return nil, fmt.Errorf("chan-try-recv! requires a channel, got %T", channelValue)
	}

	// Try to receive the value
	value, ok := channel.TryReceive()
	if !ok {
		// No value available or channel closed
		return &types.NilValue{}, nil
	}

	return value, nil
}

// Chan-close! closes a channel
func (p *ConcurrencyPlugin) chanClose(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	// Evaluate channel argument
	channelValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	channel, ok := channelValue.(*types.ChannelValue)
	if !ok {
		return nil, fmt.Errorf("chan-close! requires a channel, got %T", channelValue)
	}

	// Close the channel
	channel.Close()

	return types.BooleanValue(true), nil
}

// Chan-closed? checks if a channel is closed
func (p *ConcurrencyPlugin) chanClosed(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	// Evaluate channel argument
	channelValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	channel, ok := channelValue.(*types.ChannelValue)
	if !ok {
		return nil, fmt.Errorf("chan-closed? requires a channel, got %T", channelValue)
	}

	// Check if channel is closed
	return types.BooleanValue(channel.IsClosed()), nil
}
