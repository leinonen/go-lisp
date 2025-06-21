package evaluator

import (
	"github.com/leinonen/go-lisp/pkg/interfaces"
	"github.com/leinonen/go-lisp/pkg/types"
)

// Ensure Environment implements the required interfaces
var _ interfaces.EnvironmentReader = (*Environment)(nil)
var _ interfaces.EnvironmentWriter = (*Environment)(nil)

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
	// These wrapper functions allow arithmetic operations to be used in fn expressions
	env.bindings["+"] = &types.ArithmeticFunctionValue{Operation: "+"}
	env.bindings["-"] = &types.ArithmeticFunctionValue{Operation: "-"}
	env.bindings["*"] = &types.ArithmeticFunctionValue{Operation: "*"}
	env.bindings["/"] = &types.ArithmeticFunctionValue{Operation: "/"}
	env.bindings["%"] = &types.ArithmeticFunctionValue{Operation: "%"}

	// Register other built-in functions as callable functions
	env.bindings["list"] = &types.BuiltinFunctionValue{Name: "list"}
	env.bindings["first"] = &types.BuiltinFunctionValue{Name: "first"}
	env.bindings["rest"] = &types.BuiltinFunctionValue{Name: "rest"}
	env.bindings["cons"] = &types.BuiltinFunctionValue{Name: "cons"}
	env.bindings["length"] = &types.BuiltinFunctionValue{Name: "length"}
	env.bindings["empty?"] = &types.BuiltinFunctionValue{Name: "empty?"}
	env.bindings["map"] = &types.BuiltinFunctionValue{Name: "map"}
	env.bindings["filter"] = &types.BuiltinFunctionValue{Name: "filter"}
	env.bindings["reduce"] = &types.BuiltinFunctionValue{Name: "reduce"}
	env.bindings["append"] = &types.BuiltinFunctionValue{Name: "append"}
	env.bindings["reverse"] = &types.BuiltinFunctionValue{Name: "reverse"}
	env.bindings["nth"] = &types.BuiltinFunctionValue{Name: "nth"}

	// File operations
	env.bindings["read-file"] = &types.BuiltinFunctionValue{Name: "read-file"}
	env.bindings["write-file"] = &types.BuiltinFunctionValue{Name: "write-file"}
	env.bindings["file-exists?"] = &types.BuiltinFunctionValue{Name: "file-exists?"}

	return env
}

// NewEnvironmentWithParent creates a new environment with a parent environment
func NewEnvironmentWithParent(parent *Environment) *Environment {
	env := &Environment{
		bindings:    make(map[string]types.Value),
		parent:      parent,
		modules:     make(map[string]*types.ModuleValue),
		loadedFiles: make(map[string]bool),
	}

	// Copy modules and loaded files from parent to child
	if parent != nil {
		// Copy modules
		for name, module := range parent.modules {
			env.modules[name] = module
		}
		// Copy loaded files
		for filename, loaded := range parent.loadedFiles {
			env.loadedFiles[filename] = loaded
		}
	}

	return env
}

func (e *Environment) Set(name string, value types.Value) {
	e.bindings[name] = value
}

func (e *Environment) Get(name string) (types.Value, bool) {
	if value, exists := e.bindings[name]; exists {
		return value, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

// Has checks if a binding exists in the environment or any parent environment
func (e *Environment) Has(name string) bool {
	_, exists := e.Get(name)
	return exists
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

// Methods for completion support

// GetBindings returns a copy of the current environment's bindings
func (e *Environment) GetBindings() map[string]types.Value {
	bindings := make(map[string]types.Value)
	for name, value := range e.bindings {
		bindings[name] = value
	}
	return bindings
}

// GetParent returns the parent environment (can be nil)
func (e *Environment) GetParent() *Environment {
	return e.parent
}

// GetModules returns a copy of the modules map
func (e *Environment) GetModules() map[string]*types.ModuleValue {
	modules := make(map[string]*types.ModuleValue)
	for name, module := range e.modules {
		modules[name] = module
	}
	return modules
}

// Interface implementations for EnvironmentReader and EnvironmentWriter

// ListBindings returns a copy of all bindings including parent environments
func (e *Environment) ListBindings() map[string]types.Value {
	bindings := make(map[string]types.Value)

	// Start from root and work down
	environments := []*Environment{}
	current := e
	for current != nil {
		environments = append([]*Environment{current}, environments...)
		current = current.parent
	}

	// Apply bindings from root to current (children override parents)
	for _, env := range environments {
		for name, value := range env.bindings {
			bindings[name] = value
		}
	}

	return bindings
}

// Parent returns the parent environment as EnvironmentReader
func (e *Environment) Parent() interfaces.EnvironmentReader {
	if e.parent == nil {
		return nil
	}
	return e.parent
}

// Delete removes a binding from the current environment
func (e *Environment) Delete(name string) {
	delete(e.bindings, name)
}
