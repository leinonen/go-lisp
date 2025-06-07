package evaluator

import "github.com/leinonen/lisp-interpreter/pkg/types"

// Environment represents a variable binding environment
type Environment struct {
	bindings map[string]types.Value
	parent   *Environment
	modules  map[string]*types.ModuleValue // module registry
}

func NewEnvironment() *Environment {
	return &Environment{
		bindings: make(map[string]types.Value),
		parent:   nil,
		modules:  make(map[string]*types.ModuleValue),
	}
}

func (e *Environment) Set(name string, value types.Value) {
	e.bindings[name] = value
}

func (e *Environment) Get(name string) (types.Value, bool) {
	if value, ok := e.bindings[name]; ok {
		return value, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

// NewChildEnvironment creates a new environment with this environment as parent
func (e *Environment) NewChildEnvironment() types.Environment {
	return &Environment{
		bindings: make(map[string]types.Value),
		parent:   e,
		modules:  e.modules, // share module registry with parent
	}
}

// Module-related methods
func (e *Environment) GetModule(name string) (*types.ModuleValue, bool) {
	if module, ok := e.modules[name]; ok {
		return module, true
	}
	if e.parent != nil {
		return e.parent.GetModule(name)
	}
	return nil, false
}

func (e *Environment) SetModule(name string, module *types.ModuleValue) {
	if e.modules == nil {
		e.modules = make(map[string]*types.ModuleValue)
	}
	e.modules[name] = module
}
