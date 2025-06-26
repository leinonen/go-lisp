package kernel

// Environment for lexical scoping

import "fmt"

// Environment represents a lexical environment with symbol bindings
type Environment struct {
	bindings    map[Symbol]Value
	parent      *Environment
	loopContext *LoopContext
}

// NewEnvironment creates a new environment with optional parent
func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		bindings:    make(map[Symbol]Value),
		parent:      parent,
		loopContext: nil,
	}
}

// Get looks up a symbol in the environment
func (env *Environment) Get(sym Symbol) (Value, error) {
	if val, ok := env.bindings[sym]; ok {
		return val, nil
	}

	if env.parent != nil {
		return env.parent.Get(sym)
	}

	return nil, fmt.Errorf("undefined symbol: %s", sym)
}

// Set binds a value to a symbol in the current environment
func (env *Environment) Set(sym Symbol, val Value) {
	env.bindings[sym] = val
}

// SetLocal sets a value in the current environment only (for lexical closures)
func (env *Environment) SetLocal(sym Symbol, val Value) {
	env.bindings[sym] = val
}

// Define sets a value in the global environment
func (env *Environment) Define(sym Symbol, val Value) {
	// Find global environment
	global := env
	for global.parent != nil {
		global = global.parent
	}
	global.Set(sym, val)
}

// Update modifies an existing binding in the current or parent environment
// This is crucial for lexical closures with mutable captured variables
func (env *Environment) Update(sym Symbol, val Value) error {
	if _, ok := env.bindings[sym]; ok {
		env.bindings[sym] = val
		return nil
	}

	if env.parent != nil {
		return env.parent.Update(sym, val)
	}

	return fmt.Errorf("undefined symbol: %s", sym)
}

// HasBinding checks if a symbol is bound in this environment or its parents
func (env *Environment) HasBinding(sym Symbol) bool {
	if _, ok := env.bindings[sym]; ok {
		return true
	}

	if env.parent != nil {
		return env.parent.HasBinding(sym)
	}

	return false
}

// CreateClosure creates a new environment that captures the current lexical scope
// This is used for proper lexical closure implementation
func (env *Environment) CreateClosure() *Environment {
	// Create a new environment with the current one as parent
	// This ensures proper lexical scoping
	return NewEnvironment(env)
}

// GetAllBindings returns all bindings in this environment and its parents
func (env *Environment) GetAllBindings() map[Symbol]Value {
	all := make(map[Symbol]Value)

	// Start from root and work down, so child bindings override parent bindings
	envs := []*Environment{}
	current := env
	for current != nil {
		envs = append([]*Environment{current}, envs...)
		current = current.parent
	}

	// Collect all bindings from root to current
	for _, e := range envs {
		for sym, val := range e.bindings {
			all[sym] = val
		}
	}

	return all
}

// GetLocalBindings returns only the bindings in this environment level
func (env *Environment) GetLocalBindings() map[Symbol]Value {
	result := make(map[Symbol]Value)
	for sym, val := range env.bindings {
		result[sym] = val
	}
	return result
}

// SetLoopContext sets the loop context for the current environment
func (env *Environment) SetLoopContext(bindings []Symbol) {
	env.loopContext = &LoopContext{
		Bindings: bindings,
		Parent:   env.loopContext,
	}
}

// ClearLoopContext removes the current loop context
func (env *Environment) ClearLoopContext() {
	if env.loopContext != nil {
		env.loopContext = env.loopContext.Parent
	}
}

// GetLoopContext returns the current loop context
func (env *Environment) GetLoopContext() *LoopContext {
	// Look for loop context in current environment first
	if env.loopContext != nil {
		return env.loopContext
	}

	// Then check parent environments
	if env.parent != nil {
		return env.parent.GetLoopContext()
	}

	return nil
}
