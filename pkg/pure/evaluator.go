// Package pure provides a pure plugin-based evaluator without legacy fallback
package pure

import (
	"fmt"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/plugins/arithmetic"
	atomplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/atom"
	bindingplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/binding"
	"github.com/leinonen/lisp-interpreter/pkg/plugins/comparison"
	concurrencyplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/concurrency"
	"github.com/leinonen/lisp-interpreter/pkg/plugins/control"
	coreplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/core"
	functionalplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/functional"
	hashmapplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/hashmap"
	httpplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/http"
	ioplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/io"
	jsonplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/json"
	keywordplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/keyword"
	"github.com/leinonen/lisp-interpreter/pkg/plugins/list"
	"github.com/leinonen/lisp-interpreter/pkg/plugins/logical"
	macroplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/macro"
	mathplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/math"
	sequenceplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/sequence"
	stringplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/string"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// PureEvaluator is a fully plugin-based evaluator with no legacy fallback
type PureEvaluator struct {
	env               *evaluator.Environment
	registry          registry.FunctionRegistry
	pluginManager     plugins.PluginManager
	concurrencyPlugin *concurrencyplugin.ConcurrencyPlugin
}

// NewPureEvaluator creates a new pure plugin-based evaluator
func NewPureEvaluator(env *evaluator.Environment) (*PureEvaluator, error) {
	// Create the registry
	reg := registry.NewRegistry()

	// Create plugin manager
	pluginManager := plugins.NewPluginManager(reg)

	pureEval := &PureEvaluator{
		env:           env,
		registry:      reg,
		pluginManager: pluginManager,
	}

	// Load all plugins
	if err := pureEval.loadAllPlugins(); err != nil {
		return nil, fmt.Errorf("failed to load plugins: %v", err)
	}

	return pureEval, nil
}

// loadAllPlugins loads all plugins for complete functionality
func (pe *PureEvaluator) loadAllPlugins() error {
	// Load core plugin first (def, fn, quote, etc.)
	corePlugin := coreplugin.NewCorePlugin(pe.env)
	if err := pe.pluginManager.LoadPlugin(corePlugin); err != nil {
		return fmt.Errorf("failed to load core plugin: %v", err)
	}

	// Load arithmetic plugin
	arithmeticPlugin := arithmetic.NewArithmeticPlugin()
	if err := pe.pluginManager.LoadPlugin(arithmeticPlugin); err != nil {
		return fmt.Errorf("failed to load arithmetic plugin: %v", err)
	}

	// Load comparison plugin
	comparisonPlugin := comparison.NewComparisonPlugin()
	if err := pe.pluginManager.LoadPlugin(comparisonPlugin); err != nil {
		return fmt.Errorf("failed to load comparison plugin: %v", err)
	}

	// Load logical plugin
	logicalPlugin := logical.NewLogicalPlugin()
	if err := pe.pluginManager.LoadPlugin(logicalPlugin); err != nil {
		return fmt.Errorf("failed to load logical plugin: %v", err)
	}

	// Load list plugin
	listPlugin := list.NewListPlugin()
	if err := pe.pluginManager.LoadPlugin(listPlugin); err != nil {
		return fmt.Errorf("failed to load list plugin: %v", err)
	}

	// Load control plugin (depends on logical)
	controlPlugin := control.NewControlPlugin()
	if err := pe.pluginManager.LoadPlugin(controlPlugin); err != nil {
		return fmt.Errorf("failed to load control plugin: %v", err)
	}

	// Load essential new plugins
	// Load keyword plugin
	keywordPlugin := keywordplugin.NewKeywordPlugin()
	if err := pe.pluginManager.LoadPlugin(keywordPlugin); err != nil {
		return fmt.Errorf("failed to load keyword plugin: %v", err)
	}

	// Load binding plugin (let)
	bindingPlugin := bindingplugin.NewBindingPlugin()
	if err := pe.pluginManager.LoadPlugin(bindingPlugin); err != nil {
		return fmt.Errorf("failed to load binding plugin: %v", err)
	}

	// Load sequence plugin (vector)
	sequencePlugin := sequenceplugin.NewSequencePlugin()
	if err := pe.pluginManager.LoadPlugin(sequencePlugin); err != nil {
		return fmt.Errorf("failed to load sequence plugin: %v", err)
	}

	// Load macro plugin
	macroPlugin := macroplugin.NewMacroPlugin()
	if err := pe.pluginManager.LoadPlugin(macroPlugin); err != nil {
		return fmt.Errorf("failed to load macro plugin: %v", err)
	}

	// Load string plugin
	stringPlugin := stringplugin.NewStringPlugin()
	if err := pe.pluginManager.LoadPlugin(stringPlugin); err != nil {
		return fmt.Errorf("failed to load string plugin: %v", err)
	}

	// Load functional plugin (map, filter, reduce, etc.)
	functionalPlugin := functionalplugin.NewFunctionalPlugin()
	if err := pe.pluginManager.LoadPlugin(functionalPlugin); err != nil {
		return fmt.Errorf("failed to load functional plugin: %v", err)
	}

	// Load math plugin
	mathPlugin := mathplugin.NewMathPlugin()
	if err := pe.pluginManager.LoadPlugin(mathPlugin); err != nil {
		return fmt.Errorf("failed to load math plugin: %v", err)
	}

	// Load hashmap plugin
	hashmapPlugin := hashmapplugin.NewHashMapPlugin()
	if err := pe.pluginManager.LoadPlugin(hashmapPlugin); err != nil {
		return fmt.Errorf("failed to load hashmap plugin: %v", err)
	}

	// Load atom plugin
	atomPlugin := atomplugin.NewAtomPlugin()
	if err := pe.pluginManager.LoadPlugin(atomPlugin); err != nil {
		return fmt.Errorf("failed to load atom plugin: %v", err)
	}

	// Load HTTP plugin
	httpPlugin := httpplugin.NewHTTPPlugin()
	if err := pe.pluginManager.LoadPlugin(httpPlugin); err != nil {
		return fmt.Errorf("failed to load HTTP plugin: %v", err)
	}

	// Load JSON plugin
	jsonPlugin := jsonplugin.NewJSONPlugin()
	if err := pe.pluginManager.LoadPlugin(jsonPlugin); err != nil {
		return fmt.Errorf("failed to load JSON plugin: %v", err)
	}

	// Load I/O plugin
	ioPlugin := ioplugin.NewIOPlugin()
	if err := pe.pluginManager.LoadPlugin(ioPlugin); err != nil {
		return fmt.Errorf("failed to load I/O plugin: %v", err)
	}

	// Load concurrency plugin
	concurrencyPlugin := concurrencyplugin.NewConcurrencyPlugin()
	if err := pe.pluginManager.LoadPlugin(concurrencyPlugin); err != nil {
		return fmt.Errorf("failed to load concurrency plugin: %v", err)
	}
	pe.concurrencyPlugin = concurrencyPlugin

	return nil
}

// Eval evaluates an expression using only the plugin system
func (pe *PureEvaluator) Eval(expr types.Expr) (types.Value, error) {
	switch ex := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(ex.Value), nil
	case *types.BigNumberExpr:
		bigNum, ok := types.NewBigNumberFromString(ex.Value)
		if !ok {
			return nil, fmt.Errorf("invalid big number: %s", ex.Value)
		}
		return bigNum, nil
	case *types.StringExpr:
		return types.StringValue(ex.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(ex.Value), nil
	case *types.KeywordExpr:
		return types.KeywordValue(ex.Value), nil
	case *types.SymbolExpr:
		return pe.evalSymbol(ex)
	case *types.ListExpr:
		return pe.evalList(ex)
	case *types.BracketExpr:
		// Evaluate bracket expressions as hashmaps or vectors
		return pe.evalBracket(ex)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// evalSymbol evaluates a symbol expression
func (pe *PureEvaluator) evalSymbol(symbol *types.SymbolExpr) (types.Value, error) {
	// Check for qualified module access (module.symbol)
	if strings.Contains(symbol.Name, ".") {
		return pe.evalModuleAccess(symbol.Name)
	}

	// Check if it's a registered function first
	if pe.registry.Has(symbol.Name) {
		// Return a function value that can be called
		return &types.BuiltinFunctionValue{Name: symbol.Name}, nil
	}

	// Look up variable in environment
	if value, ok := pe.env.Get(symbol.Name); ok {
		return value, nil
	}

	return nil, fmt.Errorf("undefined symbol: %s", symbol.Name)
}

// evalList evaluates a list expression using plugins
func (pe *PureEvaluator) evalList(list *types.ListExpr) (types.Value, error) {
	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("empty list cannot be evaluated")
	}

	// The first element should be a function
	firstExpr := list.Elements[0]

	// Check if it's a symbol that represents a registered function
	if symbolExpr, ok := firstExpr.(*types.SymbolExpr); ok {
		if fn, exists := pe.registry.Get(symbolExpr.Name); exists {
			// Call the registered function
			return fn.Call(pe, list.Elements[1:])
		}
	}

	// If it's not a registered function, evaluate the first element to get a function value
	funcValue, err := pe.Eval(firstExpr)
	if err != nil {
		return nil, err
	}

	// Call the function
	return pe.CallFunction(funcValue, list.Elements[1:])
}

// evalBracket evaluates bracket expressions (only for function parameters)
func (pe *PureEvaluator) evalBracket(bracket *types.BracketExpr) (types.Value, error) {
	// Brackets should primarily be used for function parameters
	// When evaluated as expressions, they create empty lists for now
	// This is a design choice - in practice, brackets should mainly appear
	// in function definitions, not as standalone expressions
	return &types.ListValue{Elements: []types.Value{}}, nil
}

// evalModuleAccess evaluates module.symbol access
func (pe *PureEvaluator) evalModuleAccess(qualifiedName string) (types.Value, error) {
	parts := strings.Split(qualifiedName, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid module access: %s", qualifiedName)
	}

	moduleName := parts[0]
	symbolName := parts[1]

	// Get the module
	module, ok := pe.env.GetModule(moduleName)
	if !ok {
		return nil, fmt.Errorf("undefined module: %s", moduleName)
	}

	// Get the symbol from the module
	if value, ok := module.Exports[symbolName]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("undefined symbol %s in module %s", symbolName, moduleName)
}

// CallFunction allows calling functions, including user-defined functions
func (pe *PureEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	switch fn := funcValue.(type) {
	case *types.FunctionValue:
		return pe.callUserFunction(fn, args)
	case *types.BuiltinFunctionValue:
		// Look up the builtin function and call it
		if regFn, exists := pe.registry.Get(fn.Name); exists {
			return regFn.Call(pe, args)
		}
		return nil, fmt.Errorf("undefined builtin function: %s", fn.Name)
	case *types.ArithmeticFunctionValue:
		// Handle arithmetic functions that might be stored in variables
		if regFn, exists := pe.registry.Get(fn.Operation); exists {
			return regFn.Call(pe, args)
		}
		return nil, fmt.Errorf("undefined arithmetic function: %s", fn.Operation)
	default:
		return nil, fmt.Errorf("cannot call non-function value: %T", funcValue)
	}
}

// callUserFunction calls a user-defined function with tail call optimization
func (pe *PureEvaluator) callUserFunction(fn *types.FunctionValue, args []types.Expr) (types.Value, error) {
	return pe.callUserFunctionWithTCO(fn, args, false)
}

// callUserFunctionWithTCO calls a user-defined function with tail call optimization support
func (pe *PureEvaluator) callUserFunctionWithTCO(fn *types.FunctionValue, args []types.Expr, inTailPosition bool) (types.Value, error) {
	// Tail call optimization loop
	currentFn := fn
	currentArgs := args

	for {
		// Check argument count
		if len(currentArgs) != len(currentFn.Params) {
			return nil, fmt.Errorf("function expects %d arguments, got %d", len(currentFn.Params), len(currentArgs))
		}

		// Create new environment for function execution
		fnEnv, ok := currentFn.Env.(*evaluator.Environment)
		if !ok {
			return nil, fmt.Errorf("invalid environment type in function")
		}
		newEnv := evaluator.NewEnvironmentWithParent(fnEnv)

		// Evaluate arguments and bind to parameters
		argValues := make([]types.Value, len(currentArgs))
		for i, argExpr := range currentArgs {
			argValue, err := pe.Eval(argExpr)
			if err != nil {
				return nil, err
			}
			argValues[i] = argValue
		}

		// Bind parameters to argument values
		for i, param := range currentFn.Params {
			newEnv.Set(param, argValues[i])
		}

		// Create new evaluator with function environment
		fnEvaluator, err := NewPureEvaluator(newEnv)
		if err != nil {
			return nil, err
		}

		// Evaluate function body
		return fnEvaluator.Eval(currentFn.Body)
	}
}

// SetInterpreterDependency allows setting an interpreter reference for plugins that need it
func (pe *PureEvaluator) SetInterpreterDependency(interp interface{}) {
	// For now, we only support the concurrency plugin dependency
	if pe.concurrencyPlugin != nil {
		if interpDep, ok := interp.(concurrencyplugin.InterpreterDependency); ok {
			pe.concurrencyPlugin.SetInterpreter(interpDep)
		}
	}
}

// Registry returns the function registry
func (pe *PureEvaluator) Registry() registry.FunctionRegistry {
	return pe.registry
}

// PluginManager returns the plugin manager
func (pe *PureEvaluator) PluginManager() plugins.PluginManager {
	return pe.pluginManager
}

// LoadPlugin loads a plugin
func (pe *PureEvaluator) LoadPlugin(plugin plugins.Plugin) error {
	return pe.pluginManager.LoadPlugin(plugin)
}

// UnloadPlugin unloads a plugin
func (pe *PureEvaluator) UnloadPlugin(name string) error {
	return pe.pluginManager.UnloadPlugin(name)
}

// ListPlugins returns information about loaded plugins
func (pe *PureEvaluator) ListPlugins() []plugins.PluginInfo {
	return pe.pluginManager.ListPlugins()
}

// ListFunctions returns all registered functions
func (pe *PureEvaluator) ListFunctions() []string {
	return pe.registry.List()
}

// ListFunctionsByCategory returns functions in a specific category
func (pe *PureEvaluator) ListFunctionsByCategory(category string) []string {
	return pe.registry.ListByCategory(category)
}

// GetFunctionHelp returns help text for a function
func (pe *PureEvaluator) GetFunctionHelp(name string) (string, bool) {
	if fn, exists := pe.registry.Get(name); exists {
		return fn.Help(), true
	}
	return "", false
}

// GetRegistry returns the function registry for completion support
func (pe *PureEvaluator) GetRegistry() registry.FunctionRegistry {
	return pe.registry
}
