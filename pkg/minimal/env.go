package minimal

// Environment for lexical scoping

import "fmt"

// Environment represents a lexical environment with symbol bindings
type Environment struct {
	bindings map[Symbol]Value
	parent   *Environment
}

// NewEnvironment creates a new environment with optional parent
func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		bindings: make(map[Symbol]Value),
		parent:   parent,
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

// Define sets a value in the global environment
func (env *Environment) Define(sym Symbol, val Value) {
	// Find global environment
	global := env
	for global.parent != nil {
		global = global.parent
	}
	global.Set(sym, val)
}
