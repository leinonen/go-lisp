package evaluator

import (
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// TestBasicChannel tests basic channel creation and operations
func TestBasicChannel(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test (chan) - unbuffered channel
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan"},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	channel, ok := result.(*types.ChannelValue)
	if !ok {
		t.Fatalf("Expected ChannelValue, got: %T", result)
	}

	if channel.IsClosed() {
		t.Error("New channel should not be closed")
	}
}

// TestBufferedChannel tests buffered channel creation
func TestBufferedChannel(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test (chan 3) - buffered channel with size 3
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan"},
			&types.NumberExpr{Value: 3},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	channel, ok := result.(*types.ChannelValue)
	if !ok {
		t.Fatalf("Expected ChannelValue, got: %T", result)
	}

	if channel.IsClosed() {
		t.Error("New channel should not be closed")
	}
}

// TestChannelSendReceive tests channel send and receive operations
func TestChannelSendReceive(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create a buffered channel
	channel := types.NewChannel(2)
	env.Set("test-chan", channel)

	// Test sending values
	sendExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan-send!"},
			&types.SymbolExpr{Name: "test-chan"},
			&types.NumberExpr{Value: 42},
		},
	}

	result, err := evaluator.Eval(sendExpr)
	if err != nil {
		t.Fatalf("Expected no error sending to channel, got: %v", err)
	}

	if result != types.BooleanValue(true) {
		t.Errorf("Expected true from chan-send!, got: %v", result)
	}

	// Test receiving values
	recvExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan-recv!"},
			&types.SymbolExpr{Name: "test-chan"},
		},
	}

	result, err = evaluator.Eval(recvExpr)
	if err != nil {
		t.Fatalf("Expected no error receiving from channel, got: %v", err)
	}

	if num, ok := result.(types.NumberValue); !ok || num != 42 {
		t.Errorf("Expected 42, got: %v", result)
	}
}

// TestChannelTryReceive tests non-blocking channel receive
func TestChannelTryReceive(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create an empty channel
	channel := types.NewChannel(1)
	env.Set("test-chan", channel)

	// Test try-receive on empty channel
	tryRecvExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan-try-recv!"},
			&types.SymbolExpr{Name: "test-chan"},
		},
	}

	result, err := evaluator.Eval(tryRecvExpr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil from empty channel, got: %v", result)
	}

	// Send a value and try again
	sendExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan-send!"},
			&types.SymbolExpr{Name: "test-chan"},
			&types.StringExpr{Value: "hello"},
		},
	}

	_, err = evaluator.Eval(sendExpr)
	if err != nil {
		t.Fatalf("Expected no error sending, got: %v", err)
	}

	result, err = evaluator.Eval(tryRecvExpr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if str, ok := result.(types.StringValue); !ok || str != "hello" {
		t.Errorf("Expected 'hello', got: %v", result)
	}
}

// TestChannelClose tests channel closing
func TestChannelClose(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create a channel
	channel := types.NewChannel(1)
	env.Set("test-chan", channel)

	// Test channel is not closed initially
	closedExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan-closed?"},
			&types.SymbolExpr{Name: "test-chan"},
		},
	}

	result, err := evaluator.Eval(closedExpr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != types.BooleanValue(false) {
		t.Errorf("Expected false, got: %v", result)
	}

	// Close the channel
	closeExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan-close!"},
			&types.SymbolExpr{Name: "test-chan"},
		},
	}

	result, err = evaluator.Eval(closeExpr)
	if err != nil {
		t.Fatalf("Expected no error closing channel, got: %v", err)
	}

	if result != types.BooleanValue(true) {
		t.Errorf("Expected true from chan-close!, got: %v", result)
	}

	// Test channel is closed
	result, err = evaluator.Eval(closedExpr)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != types.BooleanValue(true) {
		t.Errorf("Expected true, got: %v", result)
	}
}

// TestChannelSendToClosedChannel tests error handling when sending to closed channel
func TestChannelSendToClosedChannel(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create and close a channel
	channel := types.NewChannel(1)
	channel.Close()
	env.Set("test-chan", channel)

	// Try to send to closed channel
	sendExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan-send!"},
			&types.SymbolExpr{Name: "test-chan"},
			&types.NumberExpr{Value: 42},
		},
	}

	_, err := evaluator.Eval(sendExpr)
	if err == nil {
		t.Fatal("Expected error sending to closed channel")
	}
}

// TestChannelWithGoroutines tests channels with goroutines
func TestChannelWithGoroutines(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create a channel
	channel := types.NewChannel(0) // unbuffered
	env.Set("test-chan", channel)

	// Start a goroutine that sends a value
	goExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "go"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "chan-send!"},
					&types.SymbolExpr{Name: "test-chan"},
					&types.NumberExpr{Value: 100},
				},
			},
		},
	}

	future, err := evaluator.Eval(goExpr)
	if err != nil {
		t.Fatalf("Expected no error starting goroutine, got: %v", err)
	}

	// Receive the value (this will block until the goroutine sends)
	recvExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan-recv!"},
			&types.SymbolExpr{Name: "test-chan"},
		},
	}

	result, err := evaluator.Eval(recvExpr)
	if err != nil {
		t.Fatalf("Expected no error receiving, got: %v", err)
	}

	if num, ok := result.(types.NumberValue); !ok || num != 100 {
		t.Errorf("Expected 100, got: %v", result)
	}

	// Wait for the goroutine to complete
	futureVal, ok := future.(*types.FutureValue)
	if !ok {
		t.Fatalf("Expected FutureValue, got: %T", future)
	}

	_, waitErr := futureVal.Wait()
	if waitErr != nil {
		t.Fatalf("Expected no error waiting for future, got: %v", waitErr)
	}
}

// TestProducerConsumerPattern tests a producer-consumer pattern
func TestProducerConsumerPattern(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create channels
	dataChan := types.NewChannel(5)
	resultChan := types.NewChannel(1)
	env.Set("data-chan", dataChan)
	env.Set("result-chan", resultChan)

	// Start producer goroutine
	producerExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "go"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "do"},
					// Send numbers 1-5
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "chan-send!"},
							&types.SymbolExpr{Name: "data-chan"},
							&types.NumberExpr{Value: 1},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "chan-send!"},
							&types.SymbolExpr{Name: "data-chan"},
							&types.NumberExpr{Value: 2},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "chan-send!"},
							&types.SymbolExpr{Name: "data-chan"},
							&types.NumberExpr{Value: 3},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "chan-close!"},
							&types.SymbolExpr{Name: "data-chan"},
						},
					},
				},
			},
		},
	}

	producerFuture, err := evaluator.Eval(producerExpr)
	if err != nil {
		t.Fatalf("Expected no error starting producer, got: %v", err)
	}

	// Start consumer goroutine
	consumerExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "go"},
			&types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "do"},
					// Sum all received values
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "def"},
							&types.SymbolExpr{Name: "sum"},
							&types.NumberExpr{Value: 0},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "def"},
							&types.SymbolExpr{Name: "val1"},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "chan-recv!"},
									&types.SymbolExpr{Name: "data-chan"},
								},
							},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "def"},
							&types.SymbolExpr{Name: "sum"},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "+"},
									&types.SymbolExpr{Name: "sum"},
									&types.SymbolExpr{Name: "val1"},
								},
							},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "def"},
							&types.SymbolExpr{Name: "val2"},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "chan-recv!"},
									&types.SymbolExpr{Name: "data-chan"},
								},
							},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "def"},
							&types.SymbolExpr{Name: "sum"},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "+"},
									&types.SymbolExpr{Name: "sum"},
									&types.SymbolExpr{Name: "val2"},
								},
							},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "def"},
							&types.SymbolExpr{Name: "val3"},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "chan-recv!"},
									&types.SymbolExpr{Name: "data-chan"},
								},
							},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "def"},
							&types.SymbolExpr{Name: "sum"},
							&types.ListExpr{
								Elements: []types.Expr{
									&types.SymbolExpr{Name: "+"},
									&types.SymbolExpr{Name: "sum"},
									&types.SymbolExpr{Name: "val3"},
								},
							},
						},
					},
					&types.ListExpr{
						Elements: []types.Expr{
							&types.SymbolExpr{Name: "chan-send!"},
							&types.SymbolExpr{Name: "result-chan"},
							&types.SymbolExpr{Name: "sum"},
						},
					},
				},
			},
		},
	}

	consumerFuture, err := evaluator.Eval(consumerExpr)
	if err != nil {
		t.Fatalf("Expected no error starting consumer, got: %v", err)
	}

	// Wait for result
	resultExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "chan-recv!"},
			&types.SymbolExpr{Name: "result-chan"},
		},
	}

	result, err := evaluator.Eval(resultExpr)
	if err != nil {
		t.Fatalf("Expected no error receiving result, got: %v", err)
	}

	if num, ok := result.(types.NumberValue); !ok || num != 6 {
		t.Errorf("Expected 6 (1+2+3), got: %v", result)
	}

	// Wait for goroutines to complete
	if future, ok := producerFuture.(*types.FutureValue); ok {
		future.Wait()
	}
	if future, ok := consumerFuture.(*types.FutureValue); ok {
		future.Wait()
	}
}
