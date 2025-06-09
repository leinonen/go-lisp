// Package evaluator_macros contains macro system functionality
package evaluator

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// evalDefmacro handles macro definitions
func (e *Evaluator) evalDefmacro(args []types.Expr) (types.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("defmacro requires at least 3 arguments: name, parameters, and body")
	}

	// First argument must be a symbol (macro name)
	nameExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("defmacro first argument must be a symbol")
	}

	// Second argument must be a list of parameter names
	paramsExpr, ok := args[1].(*types.ListExpr)
	if !ok {
		return nil, fmt.Errorf("defmacro second argument must be a parameter list")
	}

	// Extract parameter names
	params := make([]string, len(paramsExpr.Elements))
	for i, paramExpr := range paramsExpr.Elements {
		symbolExpr, ok := paramExpr.(*types.SymbolExpr)
		if !ok {
			return nil, fmt.Errorf("defmacro parameter must be a symbol, got %T", paramExpr)
		}
		params[i] = symbolExpr.Name
	}

	// For now, we'll just use the last expression as the body
	// In a more complete implementation, we'd want to evaluate all expressions
	// and return the last one (similar to a 'do' or 'progn' form)
	var body types.Expr
	if len(args) == 3 {
		body = args[2]
	} else {
		body = args[len(args)-1]
	}

	// Create the macro value with captured environment
	macro := types.MacroValue{
		Params: params,
		Body:   body,
		Env:    e.env, // capture current environment for closures
	}

	// Set the macro in the environment
	e.env.Set(nameExpr.Name, macro)

	return macro, nil
}

// evalQuote handles quote special form
func (e *Evaluator) evalQuote(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("quote requires exactly 1 argument")
	}

	// Convert the expression to a value without evaluating it
	return e.exprToQuotedValue(args[0])
}

// exprToQuotedValue converts an expression to a quoted value
func (e *Evaluator) exprToQuotedValue(expr types.Expr) (types.Value, error) {
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
		// Symbols in quotes become quoted symbols (preserving that they're symbols)
		return &types.QuotedValue{Value: ex}, nil
	case *types.ListExpr:
		// Convert list elements to quoted values
		elements := make([]types.Value, len(ex.Elements))
		for i, elem := range ex.Elements {
			value, err := e.exprToQuotedValue(elem)
			if err != nil {
				return nil, err
			}
			elements[i] = value
		}
		return &types.ListValue{Elements: elements}, nil
	default:
		return nil, fmt.Errorf("cannot quote expression of type %T", expr)
	}
}

// expandMacro expands a macro call
func (e *Evaluator) expandMacro(macro types.MacroValue, args []types.Expr) (types.Expr, error) {
	// Create a new environment for macro expansion
	macroEnv := macro.Env.NewChildEnvironment()

	// Handle special case: single parameter macros that should collect all arguments
	// This is specifically for macros like 'progn' that need variadic behavior
	if len(macro.Params) == 1 && len(args) > 1 {
		// Check if this is a known variadic macro by parameter name
		paramName := macro.Params[0]
		if paramName == "exprs" || paramName == "expressions" || paramName == "clauses" {
			// Convert all arguments into a list for the single parameter
			argValues := make([]types.Value, len(args))
			for i, arg := range args {
				quotedArg, err := e.exprToQuotedValue(arg)
				if err != nil {
					return nil, fmt.Errorf("error preparing macro argument: %v", err)
				}
				argValues[i] = quotedArg
			}

			// Create a list value containing all arguments
			argsList := &types.ListValue{Elements: argValues}
			macroEnv.Set(macro.Params[0], argsList)
		} else {
			// Standard argument count check for non-variadic macros
			return nil, fmt.Errorf("macro %v expects %d arguments, got %d",
				macro.Params, len(macro.Params), len(args))
		}
	} else {
		// Standard macro with exact parameter matching
		if len(args) != len(macro.Params) {
			return nil, fmt.Errorf("macro %v expects %d arguments, got %d",
				macro.Params, len(macro.Params), len(args))
		}

		// Bind macro parameters to arguments (without evaluating the arguments)
		for i, param := range macro.Params {
			// Convert the argument expression to a quoted value for the macro
			quotedArg, err := e.exprToQuotedValue(args[i])
			if err != nil {
				return nil, fmt.Errorf("error preparing macro argument: %v", err)
			}
			macroEnv.Set(param, quotedArg)
		}
	}

	// Evaluate the macro body in the macro environment to get the expansion
	macroEvaluator := NewEvaluator(macroEnv.(*Environment))
	expansion, err := macroEvaluator.Eval(macro.Body)
	if err != nil {
		return nil, fmt.Errorf("error expanding macro: %v", err)
	}

	// Convert the expansion result back to an expression
	return e.valueToExpr(expansion)
}

// valueToExpr converts a value back to an expression
func (e *Evaluator) valueToExpr(val types.Value) (types.Expr, error) {
	switch v := val.(type) {
	case types.NumberValue:
		return &types.NumberExpr{Value: float64(v)}, nil
	case types.StringValue:
		// String values should remain as string expressions
		return &types.StringExpr{Value: string(v)}, nil
	case types.BooleanValue:
		return &types.BooleanExpr{Value: bool(v)}, nil
	case types.KeywordValue:
		return &types.KeywordExpr{Value: string(v)}, nil
	case *types.BigNumberValue:
		return &types.BigNumberExpr{Value: v.String()}, nil
	case *types.ListValue:
		// Convert list elements to expressions
		elements := make([]types.Expr, len(v.Elements))
		for i, elem := range v.Elements {
			expr, err := e.valueToExpr(elem)
			if err != nil {
				return nil, err
			}
			elements[i] = expr
		}
		return &types.ListExpr{Elements: elements}, nil
	case *types.QuotedValue:
		// Quoted values should return their contained expression
		return v.Value, nil
	case *types.NilValue:
		return &types.SymbolExpr{Name: "nil"}, nil
	default:
		return nil, fmt.Errorf("cannot convert value of type %T to expression", val)
	}
}

// isMacroCall checks if a function call is actually a macro call
func (e *Evaluator) isMacroCall(funcName string) (types.MacroValue, bool) {
	value, ok := e.env.Get(funcName)
	if !ok {
		return types.MacroValue{}, false
	}

	macro, ok := value.(types.MacroValue)
	return macro, ok
}
