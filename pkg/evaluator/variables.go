// Package evaluator_variables contains variable definition and environment inspection functionality
package evaluator

import (
	"fmt"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Variable definition

func (e *Evaluator) evalDefine(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("define requires exactly 2 arguments: name and value")
	}

	// First argument must be a symbol (variable name)
	nameExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("define first argument must be a symbol")
	}

	// Evaluate the second argument (the value)
	value, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	// Set the variable in the environment
	e.env.Set(nameExpr.Name, value)

	// Return the value that was defined
	return value, nil
}

// Environment inspection methods

func (e *Evaluator) evalEnv(args []types.Expr) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("env requires no arguments")
	}

	// Create a list of (name value) pairs for all bindings in the environment
	var elements []types.Value

	for name, value := range e.env.bindings {
		// Create a pair (name value)
		pair := &types.ListValue{
			Elements: []types.Value{
				types.StringValue(name),
				value,
			},
		}
		elements = append(elements, pair)
	}

	return &types.ListValue{Elements: elements}, nil
}

func (e *Evaluator) evalBuiltins(args []types.Expr) (types.Value, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("builtins requires 0 or 1 arguments")
	}

	// If no arguments, list all built-in functions
	if len(args) == 0 {
		// List all built-in functions and special forms
		builtinNames := []string{
			// Arithmetic operations
			"+", "-", "*", "/",
			// Comparison operations
			"=", "<", ">", "<=", ">=",
			// Logical operations
			"and", "or", "not",
			// Control flow
			"if",
			// Variable and function definition
			"define", "lambda", "defun",
			// List operations
			"list", "first", "rest", "cons", "length", "empty?",
			// Higher-order functions
			"map", "filter", "reduce",
			// List manipulation
			"append", "reverse", "nth",
			// Environment inspection
			"env", "modules", "builtins",
		}

		// Convert to list of string values
		var elements []types.Value
		for _, name := range builtinNames {
			elements = append(elements, types.StringValue(name))
		}

		return &types.ListValue{Elements: elements}, nil
	}

	// If one argument, show help for that function
	funcNameExpr, ok := args[0].(*types.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("builtins argument must be a symbol")
	}

	funcName := funcNameExpr.Name
	helpText := e.getBuiltinHelp(funcName)
	if helpText == "" {
		return nil, fmt.Errorf("no help available for '%s' (not a built-in function)", funcName)
	}

	return types.StringValue(helpText), nil
}

func (e *Evaluator) getBuiltinHelp(funcName string) string {
	helpMap := map[string]string{
		// Arithmetic operations
		"+": "(+ num1 num2 ...)\nAddition with multiple operands.\nExample: (+ 1 2 3) => 6",
		"-": "(- num1 num2)\nSubtraction of two numbers.\nExample: (- 10 3) => 7",
		"*": "(* num1 num2 ...)\nMultiplication with multiple operands.\nExample: (* 2 3 4) => 24",
		"/": "(/ num1 num2)\nDivision of two numbers.\nExample: (/ 15 3) => 5",

		// Comparison operations
		"=":  "(= val1 val2)\nEquality comparison.\nExample: (= 5 5) => #t",
		"<":  "(< num1 num2)\nLess than comparison.\nExample: (< 3 5) => #t",
		">":  "(> num1 num2)\nGreater than comparison.\nExample: (> 7 3) => #t",
		"<=": "(<= num1 num2)\nLess than or equal comparison.\nExample: (<= 3 5) => #t, (<= 5 5) => #t",
		">=": "(>= num1 num2)\nGreater than or equal comparison.\nExample: (>= 7 3) => #t, (>= 5 5) => #t",

		// Logical operations
		"and": "(and expr1 expr2 ...)\nLogical AND - returns #t if all expressions are true.\nExample: (and #t #t) => #t, (and #t #f) => #f",
		"or":  "(or expr1 expr2 ...)\nLogical OR - returns #t if any expression is true.\nExample: (or #f #t) => #t, (or #f #f) => #f",
		"not": "(not expr)\nLogical NOT - returns the opposite of the expression.\nExample: (not #t) => #f, (not #f) => #t",

		// Control flow
		"if": "(if condition then-expr else-expr)\nConditional expression.\nExample: (if (< 3 5) \"yes\" \"no\") => \"yes\"",

		// Variable and function definition
		"define": "(define name value)\nDefine a variable with a name and value.\nExample: (define x 42)",
		"lambda": "(lambda (params) body)\nCreate an anonymous function.\nExample: (lambda (x) (+ x 1))",
		"defun":  "(defun name (params) body)\nDefine a named function.\nExample: (defun square (x) (* x x))",

		// List operations
		"list":   "(list elem1 elem2 ...)\nCreate a list with the given elements.\nExample: (list 1 2 3) => (1 2 3)",
		"first":  "(first lst)\nGet the first element of a list.\nExample: (first (list 1 2 3)) => 1",
		"rest":   "(rest lst)\nGet all elements except the first.\nExample: (rest (list 1 2 3)) => (2 3)",
		"cons":   "(cons elem lst)\nPrepend an element to a list.\nExample: (cons 0 (list 1 2)) => (0 1 2)",
		"length": "(length lst)\nGet the number of elements in a list.\nExample: (length (list 1 2 3)) => 3",
		"empty?": "(empty? lst)\nCheck if a list is empty.\nExample: (empty? (list)) => #t",

		// Higher-order functions
		"map":    "(map func lst)\nApply a function to each element of a list.\nExample: (map (lambda (x) (* x x)) (list 1 2 3)) => (1 4 9)",
		"filter": "(filter predicate lst)\nKeep only elements that satisfy a predicate.\nExample: (filter (lambda (x) (> x 0)) (list -1 2 -3 4)) => (2 4)",
		"reduce": "(reduce func init lst)\nReduce a list to a single value using a function.\nExample: (reduce (lambda (acc x) (+ acc x)) 0 (list 1 2 3)) => 6",

		// List manipulation
		"append":  "(append lst1 lst2)\nCombine two lists into one.\nExample: (append (list 1 2) (list 3 4)) => (1 2 3 4)",
		"reverse": "(reverse lst)\nReverse the order of elements in a list.\nExample: (reverse (list 1 2 3)) => (3 2 1)",
		"nth":     "(nth index lst)\nGet the element at a specific index (0-based).\nExample: (nth 1 (list \"a\" \"b\" \"c\")) => \"b\"",

		// Environment inspection
		"env":      "(env)\nShow all variables and functions in the current environment.\nExample: (env) => ((x 42) (square #<function([x])>))",
		"modules":  "(modules)\nShow all loaded modules and their exported symbols.\nExample: (modules) => ((math (square cube)))",
		"builtins": "(builtins) or (builtins func-name)\nShow all built-in functions or help for a specific function.\nExample: (builtins) => (+ - * / ...) or (builtins reduce) => help for reduce",
	}

	return helpMap[funcName]
}
