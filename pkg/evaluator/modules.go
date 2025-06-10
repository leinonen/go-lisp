// Package evaluator_modules contains module system functionality
package evaluator

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Module system methods

func (e *Evaluator) evalQualifiedSymbol(name string) (types.Value, error) {
	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid qualified symbol: %s", name)
	}

	moduleName := parts[0]
	symbolName := parts[1]

	// First try to get module from the module registry
	module, ok := e.env.GetModule(moduleName)
	if !ok {
		// If not found in module registry, check if it's stored as a regular variable (for aliases)
		moduleValue, ok := e.env.Get(moduleName)
		if !ok {
			return nil, fmt.Errorf("module not found: %s", moduleName)
		}

		// Check if the variable is actually a module
		module, ok = moduleValue.(*types.ModuleValue)
		if !ok {
			return nil, fmt.Errorf("symbol %s is not a module", moduleName)
		}
	}

	value, ok := module.Exports[symbolName]
	if !ok {
		return nil, fmt.Errorf("symbol %s not exported by module %s", symbolName, moduleName)
	}

	return value, nil
}

func (e *Evaluator) evalModule(moduleExpr *types.ModuleExpr) (types.Value, error) {
	// Create a new environment for the module
	moduleEnv := e.env.NewChildEnvironment()
	moduleEvaluator := NewEvaluator(moduleEnv.(*Environment))

	// Evaluate the module body
	for _, expr := range moduleExpr.Body {
		_, err := moduleEvaluator.Eval(expr)
		if err != nil {
			return nil, fmt.Errorf("error in module %s: %v", moduleExpr.Name, err)
		}
	}

	// Create the module value with exported bindings
	module := &types.ModuleValue{
		Name:    moduleExpr.Name,
		Exports: make(map[string]types.Value),
		Env:     moduleEnv,
	}

	// Collect exported symbols
	for _, exportName := range moduleExpr.Exports {
		value, ok := moduleEnv.Get(exportName)
		if !ok {
			return nil, fmt.Errorf("exported symbol %s not found in module %s", exportName, moduleExpr.Name)
		}
		module.Exports[exportName] = value
	}

	// Register the module in the global environment
	e.env.SetModule(moduleExpr.Name, module)

	return module, nil
}

func (e *Evaluator) evalImport(importExpr *types.ImportExpr) (types.Value, error) {
	module, ok := e.env.GetModule(importExpr.ModuleName)
	if !ok {
		return nil, fmt.Errorf("module not found: %s", importExpr.ModuleName)
	}

	// Import all exported symbols into current environment
	for name, value := range module.Exports {
		e.env.Set(name, value)
	}

	return module, nil
}

func (e *Evaluator) evalLoad(loadExpr *types.LoadExpr) (types.Value, error) {
	// Read the file
	content, err := ioutil.ReadFile(loadExpr.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", loadExpr.Filename, err)
	}

	// Tokenize
	tokenizer := tokenizer.NewTokenizer(string(content))
	tokens, err := tokenizer.TokenizeWithError()
	if err != nil {
		return nil, fmt.Errorf("tokenization error in %s: %v", loadExpr.Filename, err)
	}

	// Parse and evaluate each expression in the file
	var lastValue types.Value = types.BooleanValue(true) // default return value

	// Parse expressions one by one until we reach the end
	i := 0
	for i < len(tokens) {
		// Find the end of the current expression
		if tokens[i].Type == types.TokenType(-1) { // EOF
			break
		}

		// Extract tokens for this expression
		exprTokens, newIndex := e.extractExpression(tokens, i)
		if len(exprTokens) == 0 {
			break
		}

		// Parse the expression
		parser := parser.NewParser(exprTokens)
		expr, err := parser.Parse()
		if err != nil {
			return nil, fmt.Errorf("parse error in %s: %v", loadExpr.Filename, err)
		}

		// Evaluate the expression
		value, err := e.Eval(expr)
		if err != nil {
			return nil, fmt.Errorf("evaluation error in %s: %v", loadExpr.Filename, err)
		}

		lastValue = value
		i = newIndex
	}

	// Mark file as loaded
	e.env.MarkFileLoaded(loadExpr.Filename)

	return lastValue, nil
}

// Helper function to extract a complete expression from tokens
func (e *Evaluator) extractExpression(tokens []types.Token, start int) ([]types.Token, int) {
	if start >= len(tokens) {
		return nil, start
	}

	// Handle single token expressions (numbers, strings, booleans, symbols)
	if tokens[start].Type != types.LPAREN {
		return tokens[start : start+1], start + 1
	}

	// Handle list expressions - find matching closing paren
	parenCount := 0
	end := start
	for end < len(tokens) {
		if tokens[end].Type == types.LPAREN {
			parenCount++
		} else if tokens[end].Type == types.RPAREN {
			parenCount--
			if parenCount == 0 {
				end++
				break
			}
		}
		end++
	}

	if parenCount != 0 {
		// Unmatched parentheses - return what we have
		return tokens[start:end], end
	}

	return tokens[start:end], end
}

func (e *Evaluator) evalModules(args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("modules requires no arguments")
	}

	// Create a list of (module-name exports-list) pairs for all modules
	var elements []types.Value

	for name, module := range e.env.modules {
		// Create a list of exported symbol names
		var exports []types.Value
		for exportName := range module.Exports {
			exports = append(exports, types.StringValue(exportName))
		}

		// Create a pair (module-name exports-list)
		pair := &types.ListValue{
			Elements: []types.Value{
				types.StringValue(name),
				&types.ListValue{Elements: exports},
			},
		}
		elements = append(elements, pair)
	}

	return &types.ListValue{Elements: elements}, nil
}

func (e *Evaluator) evalRequire(requireExpr *types.RequireExpr) (types.Value, error) {
	// Check if file is already loaded to avoid re-evaluation
	if e.env.IsFileLoaded(requireExpr.Filename) {
		// File already loaded, we need to find the module(s) it defined
		// This is a simplified approach - in a real implementation, we'd track
		// which modules each file created
		return e.handleRequireForLoadedFile(requireExpr)
	}

	// Track modules before loading
	modulesBefore := make(map[string]*types.ModuleValue)
	for name, module := range e.env.modules {
		modulesBefore[name] = module
	}

	// Load the file using existing load functionality
	loadExpr := &types.LoadExpr{Filename: requireExpr.Filename}
	_, err := e.evalLoad(loadExpr)
	if err != nil {
		return nil, fmt.Errorf("failed to load file %s: %v", requireExpr.Filename, err)
	}

	// Find newly created modules
	var newModules []*types.ModuleValue
	for name, module := range e.env.modules {
		if _, existed := modulesBefore[name]; !existed {
			newModules = append(newModules, module)
		}
	}

	if len(newModules) == 0 {
		return nil, fmt.Errorf("no module found in file %s", requireExpr.Filename)
	}

	// For simplicity, use the first module found (most files define one module)
	module := newModules[0]

	// Handle different require modes
	if requireExpr.AsAlias != "" {
		// (require "file.lisp" :as alias) - create qualified access only
		e.env.Set(requireExpr.AsAlias, module)
	} else if len(requireExpr.OnlyList) > 0 {
		// (require "file.lisp" :only [fn1 fn2]) - import only specified functions
		for _, symbolName := range requireExpr.OnlyList {
			if exportedValue, exists := module.Exports[symbolName]; exists {
				e.env.Set(symbolName, exportedValue)
			} else {
				return nil, fmt.Errorf("symbol %s not exported by module %s", symbolName, module.Name)
			}
		}
	} else {
		// (require "file.lisp") - import all exports
		for name, value := range module.Exports {
			e.env.Set(name, value)
		}
	}

	return module, nil
}

func (e *Evaluator) handleRequireForLoadedFile(requireExpr *types.RequireExpr) (types.Value, error) {
	// For already loaded files, we need to guess which module they created
	// This is a heuristic approach - a better implementation would track file->module mappings

	// Try to find a module by filename-based heuristic
	moduleName := getModuleNameFromPath(requireExpr.Filename)
	if module, exists := e.env.GetModule(moduleName); exists {
		return e.applyRequireMode(requireExpr, module)
	}

	// If heuristic fails, look for any module that might match
	// This is imperfect but handles most common cases
	var possibleModule *types.ModuleValue
	for _, module := range e.env.modules {
		possibleModule = module
		break // Take the first one as fallback
	}

	if possibleModule == nil {
		return nil, fmt.Errorf("no suitable module found for already loaded file %s", requireExpr.Filename)
	}

	return e.applyRequireMode(requireExpr, possibleModule)
}

func (e *Evaluator) applyRequireMode(requireExpr *types.RequireExpr, module *types.ModuleValue) (types.Value, error) {
	// Handle different require modes
	if requireExpr.AsAlias != "" {
		// (require "file.lisp" :as alias) - create qualified access only
		e.env.Set(requireExpr.AsAlias, module)
	} else if len(requireExpr.OnlyList) > 0 {
		// (require "file.lisp" :only [fn1 fn2]) - import only specified functions
		for _, symbolName := range requireExpr.OnlyList {
			if exportedValue, exists := module.Exports[symbolName]; exists {
				e.env.Set(symbolName, exportedValue)
			} else {
				return nil, fmt.Errorf("symbol %s not exported by module %s", symbolName, module.Name)
			}
		}
	} else {
		// (require "file.lisp") - import all exports
		for name, value := range module.Exports {
			e.env.Set(name, value)
		}
	}

	return module, nil
}

// Helper function to extract module name from file path
func getModuleNameFromPath(filename string) string {
	// For now, use a simple heuristic - this could be improved
	// to actually parse the file and get the module name
	parts := strings.Split(filename, "/")
	baseName := parts[len(parts)-1]

	// Remove .lisp extension if present
	if strings.HasSuffix(baseName, ".lisp") {
		baseName = baseName[:len(baseName)-5]
	}

	return baseName
}
