package evaluator

import "github.com/leinonen/lisp-interpreter/pkg/types"

// Environment represents a variable binding environment
type Environment struct {
	bindings    map[string]types.Value
	parent      *Environment
	modules     map[string]*types.ModuleValue // module registry
	loadedFiles map[string]bool               // track loaded files to avoid re-loading
}

func NewEnvironment() *Environment {
	env := &Environment{
		bindings:    make(map[string]types.Value),
		parent:      nil,
		modules:     make(map[string]*types.ModuleValue),
		loadedFiles: make(map[string]bool),
	}

	// Register built-in constants
	env.bindings["nil"] = &types.NilValue{}

	// Register arithmetic operations as callable functions
	// These wrapper functions allow arithmetic operations to be used in lambda expressions
	env.bindings["+"] = &types.ArithmeticFunctionValue{Operation: "+"}
	env.bindings["-"] = &types.ArithmeticFunctionValue{Operation: "-"}
	env.bindings["*"] = &types.ArithmeticFunctionValue{Operation: "*"}
	env.bindings["/"] = &types.ArithmeticFunctionValue{Operation: "/"}
	env.bindings["%"] = &types.ArithmeticFunctionValue{Operation: "%"}

	return env
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
		bindings:    make(map[string]types.Value),
		parent:      e,
		modules:     e.modules, // share module registry with parent
		loadedFiles: e.loadedFiles,
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

// File loading tracking methods
func (e *Environment) IsFileLoaded(filename string) bool {
	if e.loadedFiles == nil {
		return false
	}
	return e.loadedFiles[filename]
}

func (e *Environment) MarkFileLoaded(filename string) {
	if e.loadedFiles == nil {
		e.loadedFiles = make(map[string]bool)
	}
	e.loadedFiles[filename] = true
}
