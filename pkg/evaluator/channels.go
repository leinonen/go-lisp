package evaluator

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// evalChan creates a new channel
func (e *Evaluator) evalChan(args []types.Expr) (types.Value, error) {
	size := 0 // Default to unbuffered channel

	if len(args) == 1 {
		// Evaluate the buffer size
		sizeValue, err := e.Eval(args[0])
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

// evalChanSend sends a value to a channel
func (e *Evaluator) evalChanSend(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("chan-send! requires exactly 2 arguments")
	}

	// Evaluate the channel
	channelValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	channel, ok := channelValue.(*types.ChannelValue)
	if !ok {
		return nil, fmt.Errorf("chan-send! requires a channel, got %T", channelValue)
	}

	// Evaluate the value to send
	value, err := e.Eval(args[1])
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

// evalChanRecv receives a value from a channel (blocking)
func (e *Evaluator) evalChanRecv(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("chan-recv! requires exactly 1 argument")
	}

	// Evaluate the channel
	channelValue, err := e.Eval(args[0])
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

// evalChanTryRecv tries to receive a value from a channel (non-blocking)
func (e *Evaluator) evalChanTryRecv(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("chan-try-recv! requires exactly 1 argument")
	}

	// Evaluate the channel
	channelValue, err := e.Eval(args[0])
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
		return nil, nil
	}

	return value, nil
}

// evalChanClose closes a channel
func (e *Evaluator) evalChanClose(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("chan-close! requires exactly 1 argument")
	}

	// Evaluate the channel
	channelValue, err := e.Eval(args[0])
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

// evalChanClosed checks if a channel is closed
func (e *Evaluator) evalChanClosed(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("chan-closed? requires exactly 1 argument")
	}

	// Evaluate the channel
	channelValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	channel, ok := channelValue.(*types.ChannelValue)
	if !ok {
		return nil, fmt.Errorf("chan-closed? requires a channel, got %T", channelValue)
	}

	return types.BooleanValue(channel.IsClosed()), nil
}
