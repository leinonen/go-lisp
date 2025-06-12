// Package evaluator_variables contains variable definition and environment inspection functionality
package evaluator

import (
	"fmt"
	"sort"

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
		return nil, fmt.Errorf("help requires 0 or 1 arguments")
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
			// Macro system
			"defmacro", "quote",
			// List operations
			"list", "first", "rest", "cons", "length", "empty?",
			// Higher-order functions
			"map", "filter", "reduce",
			// List manipulation
			"append", "reverse", "nth",
			// Hash map operations
			"hash-map", "hash-map-get", "hash-map-put", "hash-map-remove",
			"hash-map-contains?", "hash-map-keys", "hash-map-values",
			"hash-map-size", "hash-map-empty?",
			// String operations
			"string-concat", "string-length", "string-substring", "string-char-at",
			"string-upper", "string-lower", "string-trim", "string-split", "string-join",
			"string-contains?", "string-starts-with?", "string-ends-with?", "string-replace",
			"string-index-of", "string->number", "number->string", "string-regex-match?",
			"string-regex-find-all", "string-repeat", "string?", "string-empty?",
			// Module system
			"load", "import",
			// Environment inspection
			"env", "modules", "help",
			// Print functions
			"print", "println",
			// Constants
			"nil",
			// Error handling
			"error",
		}

		// Sort the builtin names alphabetically
		sort.Strings(builtinNames)

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
		return nil, fmt.Errorf("help argument must be a symbol")
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
		"lambda": "(lambda [params] body)\nCreate an anonymous function.\nExample: (lambda [x] (+ x 1))",
		"defun":  "(defun name [params] body)\nDefine a named function.\nExample: (defun square [x] (* x x))",

		// Macro system
		"defmacro": "(defmacro name [params] body)\nDefine a macro that transforms code at evaluation time.\nExample: (defmacro when [condition body] (list 'if condition body 'nil))",
		"quote":    "(quote expr) or 'expr\nReturn expression without evaluating it.\nExample: (quote (+ 1 2)) => (+ 1 2), '(+ 1 2) => (+ 1 2)",

		// List operations
		"list":   "(list elem1 elem2 ...)\nCreate a list with the given elements.\nExample: (list 1 2 3) => (1 2 3)",
		"first":  "(first lst)\nGet the first element of a list.\nExample: (first (list 1 2 3)) => 1",
		"rest":   "(rest lst)\nGet all elements except the first.\nExample: (rest (list 1 2 3)) => (2 3)",
		"cons":   "(cons elem lst)\nPrepend an element to a list.\nExample: (cons 0 (list 1 2)) => (0 1 2)",
		"length": "(length lst)\nGet the number of elements in a list.\nExample: (length (list 1 2 3)) => 3",
		"empty?": "(empty? lst)\nCheck if a list is empty.\nExample: (empty? (list)) => #t",

		// Higher-order functions
		"map":    "(map func lst)\nApply a function to each element of a list.\nExample: (map (lambda [x] (* x x)) (list 1 2 3)) => (1 4 9)",
		"filter": "(filter predicate lst)\nKeep only elements that satisfy a predicate.\nExample: (filter (lambda [x] (> x 0)) (list -1 2 -3 4)) => (2 4)",
		"reduce": "(reduce func init lst)\nReduce a list to a single value using a function.\nExample: (reduce (lambda [acc x] (+ acc x)) 0 (list 1 2 3)) => 6",

		// List manipulation
		"append":  "(append lst1 lst2)\nCombine two lists into one.\nExample: (append (list 1 2) (list 3 4)) => (1 2 3 4)",
		"reverse": "(reverse lst)\nReverse the order of elements in a list.\nExample: (reverse (list 1 2 3)) => (3 2 1)",
		"nth":     "(nth index lst)\nGet the element at a specific index (0-based).\nExample: (nth 1 (list \"a\" \"b\" \"c\")) => \"b\"",

		// Hash map operations
		"hash-map":           "(hash-map key1 value1 key2 value2 ...)\nCreate a hash map with key-value pairs.\nExample: (hash-map \"name\" \"Alice\" \"age\" 30) => {\"name\" Alice \"age\" 30}",
		"hash-map-get":       "(hash-map-get hashmap key)\nGet a value from a hash map by key.\nExample: (hash-map-get {\"x\" 42} \"x\") => 42",
		"hash-map-put":       "(hash-map-put hashmap key value)\nAdd or update a key-value pair (returns new hash map).\nExample: (hash-map-put {} \"x\" 42) => {\"x\" 42}",
		"hash-map-remove":    "(hash-map-remove hashmap key)\nRemove a key-value pair (returns new hash map).\nExample: (hash-map-remove {\"x\" 42 \"y\" 99} \"x\") => {\"y\" 99}",
		"hash-map-contains?": "(hash-map-contains? hashmap key)\nCheck if a hash map contains a key.\nExample: (hash-map-contains? {\"x\" 42} \"x\") => #t",
		"hash-map-keys":      "(hash-map-keys hashmap)\nGet all keys from a hash map as a list.\nExample: (hash-map-keys {\"a\" 1 \"b\" 2}) => (\"a\" \"b\")",
		"hash-map-values":    "(hash-map-values hashmap)\nGet all values from a hash map as a list.\nExample: (hash-map-values {\"a\" 1 \"b\" 2}) => (1 2)",
		"hash-map-size":      "(hash-map-size hashmap)\nGet the number of key-value pairs in a hash map.\nExample: (hash-map-size {\"a\" 1 \"b\" 2}) => 2",
		"hash-map-empty?":    "(hash-map-empty? hashmap)\nCheck if a hash map is empty.\nExample: (hash-map-empty? {}) => #t",

		// String operations
		"string-concat":         "(string-concat str1 str2 ...)\nConcatenate multiple strings or values.\nExample: (string-concat \"Hello\" \" \" \"World\") => \"Hello World\"",
		"string-length":         "(string-length str)\nGet the length of a string.\nExample: (string-length \"hello\") => 5",
		"string-substring":      "(string-substring str start end)\nExtract a substring from start to end (exclusive).\nExample: (string-substring \"hello\" 1 4) => \"ell\"",
		"string-char-at":        "(string-char-at str index)\nGet the character at a specific index.\nExample: (string-char-at \"hello\" 1) => \"e\"",
		"string-upper":          "(string-upper str)\nConvert string to uppercase.\nExample: (string-upper \"hello\") => \"HELLO\"",
		"string-lower":          "(string-lower str)\nConvert string to lowercase.\nExample: (string-lower \"HELLO\") => \"hello\"",
		"string-trim":           "(string-trim str)\nRemove whitespace from both ends of string.\nExample: (string-trim \"  hello  \") => \"hello\"",
		"string-split":          "(string-split str separator)\nSplit string by separator into a list.\nExample: (string-split \"a,b,c\" \",\") => (\"a\" \"b\" \"c\")",
		"string-join":           "(string-join lst separator)\nJoin list elements into a string with separator.\nExample: (string-join (list \"a\" \"b\" \"c\") \",\") => \"a,b,c\"",
		"string-contains?":      "(string-contains? str substring)\nCheck if string contains substring.\nExample: (string-contains? \"hello\" \"ell\") => #t",
		"string-starts-with?":   "(string-starts-with? str prefix)\nCheck if string starts with prefix.\nExample: (string-starts-with? \"hello\" \"he\") => #t",
		"string-ends-with?":     "(string-ends-with? str suffix)\nCheck if string ends with suffix.\nExample: (string-ends-with? \"hello\" \"lo\") => #t",
		"string-replace":        "(string-replace str old new)\nReplace all occurrences of old with new.\nExample: (string-replace \"hello\" \"l\" \"x\") => \"hexxo\"",
		"string-index-of":       "(string-index-of str substring)\nFind first index of substring (-1 if not found).\nExample: (string-index-of \"hello\" \"ell\") => 1",
		"string->number":        "(string->number str)\nConvert string to number.\nExample: (string->number \"42.5\") => 42.5",
		"number->string":        "(number->string num)\nConvert number to string.\nExample: (number->string 42) => \"42\"",
		"string-regex-match?":   "(string-regex-match? str pattern)\nCheck if string matches regex pattern.\nExample: (string-regex-match? \"hello123\" \"[a-z]+[0-9]+\") => #t",
		"string-regex-find-all": "(string-regex-find-all str pattern)\nFind all regex matches in string.\nExample: (string-regex-find-all \"abc123def456\" \"[0-9]+\") => (\"123\" \"456\")",
		"string-repeat":         "(string-repeat str count)\nRepeat string count times.\nExample: (string-repeat \"ha\" 3) => \"hahaha\"",
		"string?":               "(string? value)\nCheck if value is a string.\nExample: (string? \"hello\") => #t",
		"string-empty?":         "(string-empty? str)\nCheck if string is empty.\nExample: (string-empty? \"\") => #t",

		// Module system
		"load":   "(load \"filename\")\nLoad and execute a Lisp file.\nExample: (load \"library/core.lisp\") => loads the core library",
		"import": "(import module-name)\nImport all exported symbols from a module.\nExample: (import math) => imports all functions from math module",

		// Environment inspection
		"env":     "(env)\nShow all variables and functions in the current environment.\nExample: (env) => ((x 42) (square #<function([x])>))",
		"modules": "(modules)\nShow all loaded modules and their exported symbols.\nExample: (modules) => ((math (square cube)))",
		"help":    "(help) or (help func-name)\nShow all built-in functions or help for a specific function.\nExample: (help) => (+ - * / ...) or (help reduce) => help for reduce",

		// Print functions
		"print":   "(print value1 value2 ...)\nOutput values to stdout without newline.\nExample: (print \"Hello\" \" \" \"World\") outputs: Hello World",
		"println": "(println value1 value2 ...)\nOutput values to stdout with newline.\nExample: (println \"Hello World\") outputs: Hello World\\n",

		// Constants
		"nil": "nil\nBuilt-in constant representing empty/null values.\nExample: (hash-map-get {} \"missing\") => nil",

		// Error handling
		"error": "(error message)\nRaise an error with the given message.\nExample: (error \"Something went wrong!\") => Error: Something went wrong!",
	}

	return helpMap[funcName]
}

// Error function

func (e *Evaluator) evalError(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("error requires exactly 1 argument")
	}

	// Evaluate the message argument
	messageValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	// Convert the message to a string
	var message string
	switch msgVal := messageValue.(type) {
	case types.StringValue:
		message = string(msgVal)
	default:
		message = msgVal.String()
	}

	// Return an error with the message
	return nil, fmt.Errorf("%s", message)
}
