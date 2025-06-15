// Package registry provides a dynamic function registry for built-in functions
package registry

import (
	"fmt"
	"sort"
	"sync"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Evaluator interface for plugin functions to call back to the evaluator
type Evaluator interface {
	Eval(expr types.Expr) (types.Value, error)
	CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error)
}

// BuiltinFunction represents a built-in function that can be registered
type BuiltinFunction interface {
	// Metadata
	Name() string
	Category() string
	Arity() int // -1 for variadic, 0+ for fixed arity
	Help() string

	// Execution
	Call(evaluator Evaluator, args []types.Expr) (types.Value, error)
}

// FunctionRegistry manages registered built-in functions
type FunctionRegistry interface {
	Register(fn BuiltinFunction) error
	Unregister(name string) error
	Get(name string) (BuiltinFunction, bool)
	List() []string
	ListByCategory(category string) []string
	Categories() []string
	Has(name string) bool
}

// registry implements FunctionRegistry
type registry struct {
	functions  map[string]BuiltinFunction
	categories map[string][]string
	mutex      sync.RWMutex
}

// NewRegistry creates a new function registry
func NewRegistry() FunctionRegistry {
	return &registry{
		functions:  make(map[string]BuiltinFunction),
		categories: make(map[string][]string),
	}
}

// Register adds a function to the registry
func (r *registry) Register(fn BuiltinFunction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	name := fn.Name()
	if name == "" {
		return fmt.Errorf("function name cannot be empty")
	}

	if _, exists := r.functions[name]; exists {
		return fmt.Errorf("function %s already registered", name)
	}

	r.functions[name] = fn

	// Update category index
	category := fn.Category()
	if category != "" {
		r.categories[category] = append(r.categories[category], name)
		sort.Strings(r.categories[category])
	}

	return nil
}

// Unregister removes a function from the registry
func (r *registry) Unregister(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	fn, exists := r.functions[name]
	if !exists {
		return fmt.Errorf("function %s not found", name)
	}

	delete(r.functions, name)

	// Update category index
	category := fn.Category()
	if category != "" {
		if funcs, exists := r.categories[category]; exists {
			for i, funcName := range funcs {
				if funcName == name {
					r.categories[category] = append(funcs[:i], funcs[i+1:]...)
					break
				}
			}
			if len(r.categories[category]) == 0 {
				delete(r.categories, category)
			}
		}
	}

	return nil
}

// Get retrieves a function from the registry
func (r *registry) Get(name string) (BuiltinFunction, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	fn, exists := r.functions[name]
	return fn, exists
}

// List returns all registered function names
func (r *registry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ListByCategory returns function names in a specific category
func (r *registry) ListByCategory(category string) []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if funcs, exists := r.categories[category]; exists {
		// Return a copy to avoid external modification
		result := make([]string, len(funcs))
		copy(result, funcs)
		return result
	}
	return []string{}
}

// Categories returns all available categories
func (r *registry) Categories() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	categories := make([]string, 0, len(r.categories))
	for category := range r.categories {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// Has checks if a function is registered
func (r *registry) Has(name string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.functions[name]
	return exists
}

// Function categories
const (
	CategoryArithmetic  = "arithmetic"
	CategoryComparison  = "comparison"
	CategoryLogical     = "logical"
	CategoryList        = "list"
	CategoryString      = "string"
	CategoryIO          = "io"
	CategoryHTTP        = "http"
	CategoryJSON        = "json"
	CategoryMath        = "math"
	CategoryControl     = "control"
	CategoryFunction    = "function"
	CategoryFunctional  = "functional"
	CategoryAtom        = "atom"
	CategoryHashMap     = "hashmap"
	CategoryConcurrency = "concurrency"
	CategoryEnvironment = "environment"
	CategoryModule      = "module"
	CategoryError       = "error"
	CategoryPrint       = "print"
	CategoryCore        = "core"
	CategoryFile        = "file"
)
