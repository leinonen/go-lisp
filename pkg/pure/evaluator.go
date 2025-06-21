// Package pure provides a pure plugin-based evaluator without legacy fallback
package pure

import (
	"fmt"
	"strings"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/interfaces"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/plugins/arithmetic"
	"github.com/leinonen/go-lisp/pkg/plugins/atom"
	"github.com/leinonen/go-lisp/pkg/plugins/binding"
	"github.com/leinonen/go-lisp/pkg/plugins/comparison"
	"github.com/leinonen/go-lisp/pkg/plugins/concurrency"
	"github.com/leinonen/go-lisp/pkg/plugins/control"
	"github.com/leinonen/go-lisp/pkg/plugins/core"
	"github.com/leinonen/go-lisp/pkg/plugins/functional"
	"github.com/leinonen/go-lisp/pkg/plugins/hashmap"
	"github.com/leinonen/go-lisp/pkg/plugins/http"
	"github.com/leinonen/go-lisp/pkg/plugins/io"
	"github.com/leinonen/go-lisp/pkg/plugins/json"
	"github.com/leinonen/go-lisp/pkg/plugins/keyword"
	"github.com/leinonen/go-lisp/pkg/plugins/list"
	"github.com/leinonen/go-lisp/pkg/plugins/logical"
	"github.com/leinonen/go-lisp/pkg/plugins/macro"
	"github.com/leinonen/go-lisp/pkg/plugins/math"
	"github.com/leinonen/go-lisp/pkg/plugins/polymorphic"
	"github.com/leinonen/go-lisp/pkg/plugins/sequence"
	stringplugin "github.com/leinonen/go-lisp/pkg/plugins/string"
	"github.com/leinonen/go-lisp/pkg/plugins/utils"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// Ensure PureEvaluator implements the required interfaces
var _ interfaces.Evaluator = (*PureEvaluator)(nil)
var _ interfaces.EnvironmentProvider = (*PureEvaluator)(nil)

// PureEvaluator is a fully plugin-based evaluator with no legacy fallback
type PureEvaluator struct {
	env               *evaluator.Environment
	registry          registry.FunctionRegistry
	pluginManager     plugins.PluginManager
	concurrencyPlugin *concurrency.ConcurrencyPlugin
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
	// Load core plugin first (def, fn, quote, etc.) with dependency injection
	corePlugin := core.NewCorePlugin(pe.env)
	if err := pe.pluginManager.LoadPlugin(corePlugin); err != nil {
		return fmt.Errorf("failed to load core plugin: %v", err)
	}

	// Load arithmetic plugin with dependency injection
	arithmeticPlugin := arithmetic.NewArithmeticPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(arithmeticPlugin); err != nil {
		return fmt.Errorf("failed to load arithmetic plugin: %v", err)
	}

	// Load comparison plugin with dependency injection
	comparisonPlugin := comparison.NewComparisonPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(comparisonPlugin); err != nil {
		return fmt.Errorf("failed to load comparison plugin: %v", err)
	}

	// Load logical plugin with dependency injection
	logicalPlugin := logical.NewLogicalPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(logicalPlugin); err != nil {
		return fmt.Errorf("failed to load logical plugin: %v", err)
	}

	// Load polymorphic plugin for advanced features with dependency injection
	polymorphicPlugin := polymorphic.NewPolymorphicPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(polymorphicPlugin); err != nil {
		return fmt.Errorf("failed to load polymorphic plugin: %v", err)
	}

	// Load list plugin with dependency injection
	listPlugin := list.NewListPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(listPlugin); err != nil {
		return fmt.Errorf("failed to load list plugin: %v", err)
	}

	// Load control plugin (depends on logical) with dependency injection
	controlPlugin := control.NewControlPlugin(pe, pe)
	if err := pe.pluginManager.LoadPlugin(controlPlugin); err != nil {
		return fmt.Errorf("failed to load control plugin: %v", err)
	}

	// Load essential new plugins
	// Load keyword plugin with dependency injection
	keywordPlugin := keyword.NewKeywordPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(keywordPlugin); err != nil {
		return fmt.Errorf("failed to load keyword plugin: %v", err)
	}

	// Load binding plugin (let bindings) with dependency injection
	bindingPlugin := binding.NewBindingPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(bindingPlugin); err != nil {
		return fmt.Errorf("failed to load binding plugin: %v", err)
	}

	// Load sequence plugin with dependency injection
	sequencePlugin := sequence.NewSequencePlugin(pe)
	if err := pe.pluginManager.LoadPlugin(sequencePlugin); err != nil {
		return fmt.Errorf("failed to load sequence plugin: %v", err)
	}

	// Load macro plugin with dependency injection
	macroPlugin := macro.NewMacroPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(macroPlugin); err != nil {
		return fmt.Errorf("failed to load macro plugin: %v", err)
	}

	// Load string plugin with dependency injection
	stringPlugin := stringplugin.NewStringPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(stringPlugin); err != nil {
		return fmt.Errorf("failed to load string plugin: %v", err)
	}

	// Load utils plugin with dependency injection
	utilsPlugin := utils.NewUtilsPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(utilsPlugin); err != nil {
		return fmt.Errorf("failed to load utils plugin: %v", err)
	}

	// Load functional plugin (map, filter, reduce, etc.) with dependency injection
	functionalPlugin := functional.NewFunctionalPlugin(pe, pe)
	if err := pe.pluginManager.LoadPlugin(functionalPlugin); err != nil {
		return fmt.Errorf("failed to load functional plugin: %v", err)
	}

	// Load math plugin with dependency injection
	mathPlugin := math.NewMathPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(mathPlugin); err != nil {
		return fmt.Errorf("failed to load math plugin: %v", err)
	}

	// Load hashmap plugin with dependency injection
	hashmapPlugin := hashmap.NewHashMapPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(hashmapPlugin); err != nil {
		return fmt.Errorf("failed to load hashmap plugin: %v", err)
	}

	// Load atom plugin with dependency injection
	atomPlugin := atom.NewAtomPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(atomPlugin); err != nil {
		return fmt.Errorf("failed to load atom plugin: %v", err)
	}

	// Load HTTP plugin with dependency injection
	httpPlugin := http.NewHTTPPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(httpPlugin); err != nil {
		return fmt.Errorf("failed to load HTTP plugin: %v", err)
	}

	// Load JSON plugin with dependency injection
	jsonPlugin := json.NewJSONPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(jsonPlugin); err != nil {
		return fmt.Errorf("failed to load JSON plugin: %v", err)
	}

	// Load I/O plugin with dependency injection
	ioPlugin := io.NewIOPlugin(pe)
	if err := pe.pluginManager.LoadPlugin(ioPlugin); err != nil {
		return fmt.Errorf("failed to load I/O plugin: %v", err)
	}

	// Load concurrency plugin with dependency injection
	concurrencyPlugin := concurrency.NewConcurrencyPlugin(pe)
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
	case *types.HashMapExpr:
		return pe.evalHashMap(ex)
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

// evalBracket evaluates bracket expressions as vectors
func (pe *PureEvaluator) evalBracket(bracket *types.BracketExpr) (types.Value, error) {
	// Evaluate all elements in the bracket expression
	var elements []types.Value
	for _, elem := range bracket.Elements {
		value, err := pe.Eval(elem)
		if err != nil {
			return nil, err
		}
		elements = append(elements, value)
	}

	// Return a VectorValue
	return types.NewVectorValue(elements), nil
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
	case *types.PartialFunctionValue:
		// Handle partial functions by combining partial args with new args
		return pe.callPartialFunction(fn, args)
	case *types.ComplementFunctionValue:
		// Handle complement functions by negating the result
		return pe.callComplementFunction(fn, args)
	case *types.JuxtFunctionValue:
		// Handle juxt functions by calling all functions and collecting results
		return pe.callJuxtFunction(fn, args)
	case *types.CompFunctionValue:
		// Handle comp functions by composing functions (right to left)
		return pe.callCompFunction(fn, args)
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

// callPartialFunction handles calling a partially applied function
func (pe *PureEvaluator) callPartialFunction(partialFn *types.PartialFunctionValue, args []types.Expr) (types.Value, error) {
	// Convert the new arguments to values
	newArgs := make([]types.Value, len(args))
	for i, arg := range args {
		val, err := pe.Eval(arg)
		if err != nil {
			return nil, err
		}
		newArgs[i] = val
	}

	// Combine partial arguments with new arguments
	allArgs := make([]types.Value, 0, len(partialFn.PartialArgs)+len(newArgs))
	allArgs = append(allArgs, partialFn.PartialArgs...)
	allArgs = append(allArgs, newArgs...)

	// Convert values back to expressions for the function call
	// This is a bit of a hack - we need to create expressions from values
	argExprs := make([]types.Expr, len(allArgs))
	for i, val := range allArgs {
		argExprs[i] = pe.valueToExpr(val)
	}

	// Call the original function with combined arguments
	return pe.CallFunction(partialFn.OriginalFunction, argExprs)
}

// valueToExpr converts a value back to an expression (helper for partial function calls)
func (pe *PureEvaluator) valueToExpr(val types.Value) types.Expr {
	switch v := val.(type) {
	case types.NumberValue:
		return &types.NumberExpr{Value: float64(v)}
	case *types.BigNumberValue:
		return &types.BigNumberExpr{Value: v.Value.String()}
	case types.StringValue:
		return &types.StringExpr{Value: string(v)}
	case types.BooleanValue:
		return &types.BooleanExpr{Value: bool(v)}
	case types.KeywordValue:
		return &types.KeywordExpr{Value: string(v)}
	case *types.ListValue:
		exprs := make([]types.Expr, len(v.Elements))
		for i, elem := range v.Elements {
			exprs[i] = pe.valueToExpr(elem)
		}
		return &types.ListExpr{Elements: exprs}
	case *types.VectorValue:
		exprs := make([]types.Expr, len(v.Elements))
		for i, elem := range v.Elements {
			exprs[i] = pe.valueToExpr(elem)
		}
		return &types.BracketExpr{Elements: exprs}
	default:
		// For complex types, we'll create a symbol that can be resolved later
		// This is not ideal but works for most cases
		return &types.SymbolExpr{Name: fmt.Sprintf("#<value:%s>", val.String())}
	}
}

// evalHashMap evaluates hash map expressions
func (pe *PureEvaluator) evalHashMap(hashMap *types.HashMapExpr) (types.Value, error) {
	// Evaluate all elements in the hash map expression
	elements := make(map[string]types.Value)

	// Elements are stored as [key1, value1, key2, value2, ...]
	for i := 0; i < len(hashMap.Elements); i += 2 {
		keyExpr := hashMap.Elements[i]
		valueExpr := hashMap.Elements[i+1]

		// Evaluate key
		keyValue, err := pe.Eval(keyExpr)
		if err != nil {
			return nil, err
		}

		// Evaluate value
		value, err := pe.Eval(valueExpr)
		if err != nil {
			return nil, err
		}

		// Convert key to string for storage
		var keyStr string
		switch kv := keyValue.(type) {
		case types.StringValue:
			keyStr = string(kv)
		case types.KeywordValue:
			keyStr = string(kv)
		default:
			return nil, fmt.Errorf("hash map keys must be strings or keywords, got %T", keyValue)
		}

		elements[keyStr] = value
	}

	return &types.HashMapValue{Elements: elements}, nil
}

// SetInterpreterDependency allows setting an interpreter reference for plugins that need it
func (pe *PureEvaluator) SetInterpreterDependency(interp interface{}) {
	// For now, we only support the concurrency plugin dependency
	if pe.concurrencyPlugin != nil {
		if interpDep, ok := interp.(concurrency.InterpreterDependency); ok {
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

// EvalWithBindings evaluates an expression with additional temporary bindings
func (pe *PureEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	// Create a child environment with the additional bindings
	childEnv := pe.env.NewChildEnvironment()

	// Set the additional bindings in the child environment
	for name, value := range bindings {
		childEnv.Set(name, value)
	}

	// Create a temporary evaluator with the child environment
	tempEvaluator := &PureEvaluator{
		env:           childEnv.(*evaluator.Environment),
		registry:      pe.registry,
		pluginManager: pe.pluginManager,
	}

	// Evaluate the expression with the temporary evaluator
	return tempEvaluator.Eval(expr)
}

// callComplementFunction handles calling a complement function (negates the result)
func (pe *PureEvaluator) callComplementFunction(complementFn *types.ComplementFunctionValue, args []types.Expr) (types.Value, error) {
	// Call the original predicate function
	result, err := pe.CallFunction(complementFn.PredicateFunction, args)
	if err != nil {
		return nil, err
	}

	// Negate the result (convert to boolean first)
	switch r := result.(type) {
	case types.BooleanValue:
		return !r, nil
	case *types.NilValue:
		return types.BooleanValue(true), nil // nil is falsy, so complement is true
	default:
		// In Lisp, everything except nil and false is truthy
		// so the complement would be false
		return types.BooleanValue(false), nil
	}
}

// callJuxtFunction handles calling a juxt function (applies all functions to same args)
func (pe *PureEvaluator) callJuxtFunction(juxtFn *types.JuxtFunctionValue, args []types.Expr) (types.Value, error) {
	// Apply each function to the same arguments and collect results
	results := make([]types.Value, len(juxtFn.Functions))

	for i, fn := range juxtFn.Functions {
		result, err := pe.CallFunction(fn, args)
		if err != nil {
			return nil, fmt.Errorf("error calling function %d in juxt: %v", i, err)
		}
		results[i] = result
	}

	// Return the results as a vector
	return &types.VectorValue{Elements: results}, nil
}

// callCompFunction handles calling a comp function (function composition, right to left)
func (pe *PureEvaluator) callCompFunction(compFn *types.CompFunctionValue, args []types.Expr) (types.Value, error) {
	if len(compFn.Functions) == 0 {
		return &types.NilValue{}, nil
	}

	// Apply functions from right to left (composition)
	// Start with the rightmost function
	result, err := pe.CallFunction(compFn.Functions[len(compFn.Functions)-1], args)
	if err != nil {
		return nil, fmt.Errorf("error calling rightmost function in comp: %v", err)
	}

	// Apply remaining functions from right to left, each taking the result as a single argument
	for i := len(compFn.Functions) - 2; i >= 0; i-- {
		// Convert the result back to an expression for the next function call
		resultExpr := pe.valueToExpr(result)
		result, err = pe.CallFunction(compFn.Functions[i], []types.Expr{resultExpr})
		if err != nil {
			return nil, fmt.Errorf("error calling function %d in comp: %v", i, err)
		}
	}

	return result, nil
}

// EnvironmentProvider interface implementation

// GetEnvironment returns the current environment as EnvironmentReader
func (pe *PureEvaluator) GetEnvironment() interfaces.EnvironmentReader {
	return pe.env
}

// CreateChildEnvironment creates a new child environment
func (pe *PureEvaluator) CreateChildEnvironment() interfaces.EnvironmentWriter {
	childEnv := pe.env.NewChildEnvironment()
	return childEnv.(*evaluator.Environment)
}
