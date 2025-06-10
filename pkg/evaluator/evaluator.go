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
	case *types.BigNumberExpr:
		// Convert string to BigNumberValue
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
	case *types.RequireExpr:
		return e.evalRequire(ex)
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
			return e.evalSubtraction(list.Elements[1:])
		case "*":
			return e.evalMultiplication(list.Elements[1:])
		case "/":
			return e.evalDivision(list.Elements[1:])
		case "%":
			return e.evalModulo(list.Elements[1:])
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
		case "defmacro":
			return e.evalDefmacro(list.Elements[1:])
		case "quote":
			return e.evalQuote(list.Elements[1:])
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
		case "hash-map":
			return e.evalHashMap(list.Elements[1:])
		case "hash-map-get":
			return e.evalHashMapGet(list.Elements[1:])
		case "hash-map-put":
			return e.evalHashMapPut(list.Elements[1:])
		case "hash-map-remove":
			return e.evalHashMapRemove(list.Elements[1:])
		case "hash-map-contains?":
			return e.evalHashMapContains(list.Elements[1:])
		case "hash-map-keys":
			return e.evalHashMapKeys(list.Elements[1:])
		case "hash-map-values":
			return e.evalHashMapValues(list.Elements[1:])
		case "hash-map-size":
			return e.evalHashMapSize(list.Elements[1:])
		case "hash-map-empty?":
			return e.evalHashMapEmpty(list.Elements[1:])
		case "error":
			return e.evalError(list.Elements[1:])
		// Print functions
		case "print":
			return e.evalPrint(list.Elements[1:])
		case "println":
			return e.evalPrintln(list.Elements[1:])
		// String functions
		case "string-concat":
			return e.evalStringConcat(list.Elements[1:])
		case "string-length":
			return e.evalStringLength(list.Elements[1:])
		case "string-substring":
			return e.evalStringSubstring(list.Elements[1:])
		case "string-char-at":
			return e.evalStringCharAt(list.Elements[1:])
		case "string-upper":
			return e.evalStringUpper(list.Elements[1:])
		case "string-lower":
			return e.evalStringLower(list.Elements[1:])
		case "string-trim":
			return e.evalStringTrim(list.Elements[1:])
		case "string-split":
			return e.evalStringSplit(list.Elements[1:])
		case "string-join":
			return e.evalStringJoin(list.Elements[1:])
		case "string-contains?":
			return e.evalStringContains(list.Elements[1:])
		case "string-starts-with?":
			return e.evalStringStartsWith(list.Elements[1:])
		case "string-ends-with?":
			return e.evalStringEndsWith(list.Elements[1:])
		case "string-replace":
			return e.evalStringReplace(list.Elements[1:])
		case "string-index-of":
			return e.evalStringIndexOf(list.Elements[1:])
		case "string->number":
			return e.evalStringToNumber(list.Elements[1:])
		case "number->string":
			return e.evalNumberToString(list.Elements[1:])
		case "string-regex-match?":
			return e.evalStringRegexMatch(list.Elements[1:])
		case "string-regex-find-all":
			return e.evalStringRegexFindAll(list.Elements[1:])
		case "string-repeat":
			return e.evalStringRepeat(list.Elements[1:])
		case "string?":
			return e.evalStringPredicate(list.Elements[1:])
		case "string-empty?":
			return e.evalStringEmpty(list.Elements[1:])
		default:
			// Check if it's a macro call first
			if macro, isMacro := e.isMacroCall(symbolExpr.Name); isMacro {
				// Expand the macro and evaluate the result
				expanded, err := e.expandMacro(macro, list.Elements[1:])
				if err != nil {
					return nil, err
				}
				return e.Eval(expanded)
			}
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
