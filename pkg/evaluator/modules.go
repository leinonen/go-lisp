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

	module, ok := e.env.GetModule(moduleName)
	if !ok {
		return nil, fmt.Errorf("module not found: %s", moduleName)
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
