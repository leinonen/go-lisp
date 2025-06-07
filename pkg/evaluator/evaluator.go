// Package evaluator provides the core evaluation functionality for the Lisp interpreter
package evaluator

import (
	"fmt"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// TailCallInfo represents information needed for a tail call
type TailCallInfo struct {
	Function types.FunctionValue
	Args     []types.Value // Already evaluated arguments
}

// Evaluator evaluates expressions
type Evaluator struct {
	env        *Environment
	tailCall   *TailCallInfo // Tail call information, if any
	tailCallOK bool          // Indicates if tail call is allowed
}

func NewEvaluator(env *Environment) *Evaluator {
	return &Evaluator{env: env}
}

func (e *Evaluator) Eval(expr types.Expr) (types.Value, error) {
	switch ex := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(ex.Value), nil
	case *types.StringExpr:
		return types.StringValue(ex.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(ex.Value), nil
	case *types.SymbolExpr:
		// Check for qualified module access (module.symbol)
		if strings.Contains(ex.Name, ".") {
			return e.evalQualifiedSymbol(ex.Name)
		}
		value, ok := e.env.Get(ex.Name)
		if !ok {
			return nil, fmt.Errorf("undefined symbol: %s", ex.Name)
		}
		return value, nil
	case *types.ListExpr:
		return e.evalList(ex)
	case *types.ModuleExpr:
		return e.evalModule(ex)
	case *types.ImportExpr:
		return e.evalImport(ex)
	case *types.LoadExpr:
		return e.evalLoad(ex)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (e *Evaluator) evalList(list *types.ListExpr) (types.Value, error) {
	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("empty list cannot be evaluated")
	}

	// The first element could be a special form, built-in function, or user-defined function
	firstExpr := list.Elements[0]

	// Check if it's a symbol (special form or function name)
	if symbolExpr, ok := firstExpr.(*types.SymbolExpr); ok {
		switch symbolExpr.Name {
		case "+":
			return e.evalArithmetic(list.Elements[1:], func(a, b float64) float64 { return a + b })
		case "-":
			return e.evalArithmetic(list.Elements[1:], func(a, b float64) float64 { return a - b })
		case "*":
			return e.evalArithmetic(list.Elements[1:], func(a, b float64) float64 { return a * b })
		case "/":
			return e.evalDivision(list.Elements[1:])
		case "=":
			return e.evalEquality(list.Elements[1:])
		case "<":
			return e.evalComparison(list.Elements[1:], func(a, b float64) bool { return a < b })
		case ">":
			return e.evalComparison(list.Elements[1:], func(a, b float64) bool { return a > b })
		case "<=":
			return e.evalComparison(list.Elements[1:], func(a, b float64) bool { return a <= b })
		case ">=":
			return e.evalComparison(list.Elements[1:], func(a, b float64) bool { return a >= b })
		case "and":
			return e.evalAnd(list.Elements[1:])
		case "or":
			return e.evalOr(list.Elements[1:])
		case "not":
			return e.evalNot(list.Elements[1:])
		case "if":
			return e.evalIf(list.Elements[1:])
		case "define":
			return e.evalDefine(list.Elements[1:])
		case "lambda":
			return e.evalLambda(list.Elements[1:])
		case "defun":
			return e.evalDefun(list.Elements[1:])
		case "list":
			return e.evalListConstruction(list.Elements[1:])
		case "first":
			return e.evalFirst(list.Elements[1:])
		case "rest":
			return e.evalRest(list.Elements[1:])
		case "cons":
			return e.evalCons(list.Elements[1:])
		case "length":
			return e.evalLength(list.Elements[1:])
		case "empty?":
			return e.evalEmpty(list.Elements[1:])
		case "map":
			return e.evalMap(list.Elements[1:])
		case "filter":
			return e.evalFilter(list.Elements[1:])
		case "reduce":
			return e.evalReduce(list.Elements[1:])
		case "append":
			return e.evalAppend(list.Elements[1:])
		case "reverse":
			return e.evalReverse(list.Elements[1:])
		case "nth":
			return e.evalNth(list.Elements[1:])
		case "env":
			return e.evalEnv(list.Elements[1:])
		case "modules":
			return e.evalModules(list.Elements[1:])
		case "builtins":
			return e.evalBuiltins(list.Elements[1:])
		default:
			// Try to call it as a user-defined function
			return e.evalFunctionCall(symbolExpr.Name, list.Elements[1:])
		}
	}

	// If first element is not a symbol, evaluate it (could be a lambda expression)
	funcValue, err := e.Eval(firstExpr)
	if err != nil {
		return nil, err
	}

	// Call the function
	return e.callFunctionWithTailCheck(funcValue, list.Elements[1:])
}
