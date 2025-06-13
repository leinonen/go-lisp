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
	case *types.BracketExpr:
		// Evaluate bracket expressions as lists
		return e.evalBracket(ex)
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
		case "def":
			return e.evalDef(list.Elements[1:])
		case "fn":
			return e.evalFn(list.Elements[1:])
		case "defn":
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
		case "last":
			return e.evalLast(list.Elements[1:])
		case "butlast":
			return e.evalButlast(list.Elements[1:])
		case "flatten":
			return e.evalFlatten(list.Elements[1:])
		case "zip":
			return e.evalZip(list.Elements[1:])
		case "sort":
			return e.evalSort(list.Elements[1:])
		case "distinct":
			return e.evalDistinct(list.Elements[1:])
		case "concat":
			return e.evalConcat(list.Elements[1:])
		case "partition":
			return e.evalPartition(list.Elements[1:])
		case "env":
			return e.evalEnv(list.Elements[1:])
		case "modules":
			return e.evalModules(list.Elements[1:])
		case "help":
			return e.evalBuiltins(list.Elements[1:])
		// Atom operations
		case "atom":
			return e.evalAtom(list.Elements[1:])
		case "deref":
			return e.evalDeref(list.Elements[1:])
		case "swap!":
			return e.evalSwap(list.Elements[1:])
		case "reset!":
			return e.evalReset(list.Elements[1:])
		// Hash map operations
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
		case "print!":
			return e.evalPrint(list.Elements[1:])
		case "println!":
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
		// Goroutine functions
		case "go":
			return e.evalGo(list.Elements[1:])
		case "go-wait":
			return e.evalGoWait(list.Elements[1:])
		case "go-wait-all":
			return e.evalGoWaitAll(list.Elements[1:])
		// Wait group functions
		case "wait-group":
			return e.evalWaitGroup(list.Elements[1:])
		case "wait-group-add!":
			return e.evalWaitGroupAdd(list.Elements[1:])
		case "wait-group-done!":
			return e.evalWaitGroupDone(list.Elements[1:])
		case "wait-group-wait!":
			return e.evalWaitGroupWait(list.Elements[1:])
		// Channel functions
		case "chan":
			return e.evalChan(list.Elements[1:])
		case "chan-send!":
			return e.evalChanSend(list.Elements[1:])
		case "chan-recv!":
			return e.evalChanRecv(list.Elements[1:])
		case "chan-try-recv!":
			return e.evalChanTryRecv(list.Elements[1:])
		case "chan-close!":
			return e.evalChanClose(list.Elements[1:])
		case "chan-closed?":
			return e.evalChanClosed(list.Elements[1:])
		// Control flow
		case "do":
			return e.evalDo(list.Elements[1:])
		// Mathematical functions
		case "sqrt":
			return e.evalSqrt(list.Elements[1:])
		case "pow":
			return e.evalPow(list.Elements[1:])
		case "sin":
			return e.evalSin(list.Elements[1:])
		case "cos":
			return e.evalCos(list.Elements[1:])
		case "tan":
			return e.evalTan(list.Elements[1:])
		case "log":
			return e.evalLog(list.Elements[1:])
		case "exp":
			return e.evalExp(list.Elements[1:])
		case "floor":
			return e.evalFloor(list.Elements[1:])
		case "ceil":
			return e.evalCeil(list.Elements[1:])
		case "round":
			return e.evalRound(list.Elements[1:])
		case "abs":
			return e.evalAbs(list.Elements[1:])
		case "min":
			return e.evalMin(list.Elements[1:])
		case "max":
			return e.evalMax(list.Elements[1:])
		case "random":
			return e.evalRandom(list.Elements[1:])
		case "pi":
			return e.evalPi(list.Elements[1:])
		case "e":
			return e.evalE(list.Elements[1:])
		// Additional trigonometric functions
		case "asin":
			return e.evalAsin(list.Elements[1:])
		case "acos":
			return e.evalAcos(list.Elements[1:])
		case "atan":
			return e.evalAtan(list.Elements[1:])
		case "atan2":
			return e.evalAtan2(list.Elements[1:])
		// Hyperbolic functions
		case "sinh":
			return e.evalSinh(list.Elements[1:])
		case "cosh":
			return e.evalCosh(list.Elements[1:])
		case "tanh":
			return e.evalTanh(list.Elements[1:])
		// Angle conversion functions
		case "degrees":
			return e.evalDegrees(list.Elements[1:])
		case "radians":
			return e.evalRadians(list.Elements[1:])
		// Additional logarithm functions
		case "log10":
			return e.evalLog10(list.Elements[1:])
		case "log2":
			return e.evalLog2(list.Elements[1:])
		// Additional utility functions
		case "trunc":
			return e.evalTrunc(list.Elements[1:])
		case "sign":
			return e.evalSign(list.Elements[1:])
		case "mod":
			return e.evalMod(list.Elements[1:])
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

	// If first element is not a symbol, evaluate it (could be a fn expression)
	funcValue, err := e.Eval(firstExpr)
	if err != nil {
		return nil, err
	}

	// Call the function
	return e.callFunctionWithTailCheck(funcValue, list.Elements[1:])
}

func (e *Evaluator) evalBracket(bracket *types.BracketExpr) (types.Value, error) {
	// Evaluate bracket expressions as lists - they create list values
	values := make([]types.Value, len(bracket.Elements))
	for i, elem := range bracket.Elements {
		val, err := e.Eval(elem)
		if err != nil {
			return nil, err
		}
		values[i] = val
	}
	return &types.ListValue{Elements: values}, nil
}
