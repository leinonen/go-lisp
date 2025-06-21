// Package interfaces provides core interfaces for dependency injection
package interfaces

import "github.com/leinonen/go-lisp/pkg/types"

// CoreEvaluator provides basic expression evaluation
type CoreEvaluator interface {
	Eval(expr types.Expr) (types.Value, error)
}

// FunctionCaller provides function call capabilities
type FunctionCaller interface {
	CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error)
}

// BindingEvaluator provides evaluation with custom bindings
type BindingEvaluator interface {
	EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error)
}

// EnvironmentProvider provides access to environment management
type EnvironmentProvider interface {
	GetEnvironment() EnvironmentReader
	CreateChildEnvironment() EnvironmentWriter
}

// EnvironmentReader provides read-only access to variable bindings
type EnvironmentReader interface {
	Get(name string) (types.Value, bool)
	Has(name string) bool
	ListBindings() map[string]types.Value
}

// EnvironmentWriter provides read-write access to variable bindings
type EnvironmentWriter interface {
	EnvironmentReader
	Set(name string, value types.Value)
	Delete(name string)
	Parent() EnvironmentReader

	// Additional methods needed by plugins
	GetBindings() map[string]types.Value // For introspection
}

// Evaluator combines all evaluation capabilities for backward compatibility
type Evaluator interface {
	CoreEvaluator
	FunctionCaller
	BindingEvaluator
}
