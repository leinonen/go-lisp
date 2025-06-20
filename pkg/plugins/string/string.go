// Package string implements string manipulation functions as a plugin
package string

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// StringPlugin implements string manipulation functions
type StringPlugin struct {
	*plugins.BasePlugin
}

// NewStringPlugin creates a new string plugin
func NewStringPlugin() *StringPlugin {
	return &StringPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"string",
			"1.0.0",
			"String manipulation and processing functions",
			[]string{}, // No dependencies
		),
	}
}

// Functions returns the list of functions provided by this plugin
func (p *StringPlugin) Functions() []string {
	return []string{
		"string-concat", "string-length", "string-substring", "string-char-at",
		"string-upper", "string-lower", "string-trim", "string-split", "string-join",
		"string-contains?", "string-starts-with?", "string-ends-with?", "string-replace",
		"string-index-of", "string->number", "number->string", "string-regex-match?",
		"string-regex-find-all", "string-repeat", "string?", "string-empty?",
		// Clojure-style aliases
		"str", "subs", "split", "join", "replace", "trim", "upper-case", "lower-case",
	}
}

// RegisterFunctions registers all string functions with the registry
func (p *StringPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// String concatenation
	concatFunc := functions.NewFunction(
		"string-concat",
		registry.CategoryString,
		-1, // Variadic
		"Concatenate strings: (string-concat \"Hello\" \" \" \"World\")",
		p.evalStringConcat,
	)
	if err := reg.Register(concatFunc); err != nil {
		return err
	}

	// String length
	lengthFunc := functions.NewFunction(
		"string-length",
		registry.CategoryString,
		1,
		"Get string length: (string-length \"Hello\") => 5",
		p.evalStringLength,
	)
	if err := reg.Register(lengthFunc); err != nil {
		return err
	}

	// String substring
	substringFunc := functions.NewFunction(
		"string-substring",
		registry.CategoryString,
		3,
		"Extract substring: (string-substring \"Hello\" 1 4) => \"ell\"",
		p.evalStringSubstring,
	)
	if err := reg.Register(substringFunc); err != nil {
		return err
	}

	// String character at index
	charAtFunc := functions.NewFunction(
		"string-char-at",
		registry.CategoryString,
		2,
		"Get character at index: (string-char-at \"Hello\" 1) => \"e\"",
		p.evalStringCharAt,
	)
	if err := reg.Register(charAtFunc); err != nil {
		return err
	}

	// String uppercase
	upperFunc := functions.NewFunction(
		"string-upper",
		registry.CategoryString,
		1,
		"Convert to uppercase: (string-upper \"hello\") => \"HELLO\"",
		p.evalStringUpper,
	)
	if err := reg.Register(upperFunc); err != nil {
		return err
	}

	// String lowercase
	lowerFunc := functions.NewFunction(
		"string-lower",
		registry.CategoryString,
		1,
		"Convert to lowercase: (string-lower \"HELLO\") => \"hello\"",
		p.evalStringLower,
	)
	if err := reg.Register(lowerFunc); err != nil {
		return err
	}

	// String trim
	trimFunc := functions.NewFunction(
		"string-trim",
		registry.CategoryString,
		1,
		"Trim whitespace: (string-trim \"  hello  \") => \"hello\"",
		p.evalStringTrim,
	)
	if err := reg.Register(trimFunc); err != nil {
		return err
	}

	// String split
	splitFunc := functions.NewFunction(
		"string-split",
		registry.CategoryString,
		2,
		"Split string: (string-split \"a,b,c\" \",\") => (\"a\" \"b\" \"c\")",
		p.evalStringSplit,
	)
	if err := reg.Register(splitFunc); err != nil {
		return err
	}

	// String join
	joinFunc := functions.NewFunction(
		"string-join",
		registry.CategoryString,
		2,
		"Join strings: (string-join (list \"a\" \"b\" \"c\") \",\") => \"a,b,c\"",
		p.evalStringJoin,
	)
	if err := reg.Register(joinFunc); err != nil {
		return err
	}

	// String contains
	containsFunc := functions.NewFunction(
		"string-contains?",
		registry.CategoryString,
		2,
		"Check if string contains substring: (string-contains? \"hello\" \"ell\") => true",
		p.evalStringContains,
	)
	if err := reg.Register(containsFunc); err != nil {
		return err
	}

	// String starts with
	startsWithFunc := functions.NewFunction(
		"string-starts-with?",
		registry.CategoryString,
		2,
		"Check if string starts with prefix: (string-starts-with? \"hello\" \"he\") => true",
		p.evalStringStartsWith,
	)
	if err := reg.Register(startsWithFunc); err != nil {
		return err
	}

	// String ends with
	endsWithFunc := functions.NewFunction(
		"string-ends-with?",
		registry.CategoryString,
		2,
		"Check if string ends with suffix: (string-ends-with? \"hello\" \"lo\") => true",
		p.evalStringEndsWith,
	)
	if err := reg.Register(endsWithFunc); err != nil {
		return err
	}

	// String replace
	replaceFunc := functions.NewFunction(
		"string-replace",
		registry.CategoryString,
		3,
		"Replace all occurrences: (string-replace \"hello\" \"l\" \"x\") => \"hexxo\"",
		p.evalStringReplace,
	)
	if err := reg.Register(replaceFunc); err != nil {
		return err
	}

	// String index of
	indexOfFunc := functions.NewFunction(
		"string-index-of",
		registry.CategoryString,
		2,
		"Find first index of substring: (string-index-of \"hello\" \"ell\") => 1",
		p.evalStringIndexOf,
	)
	if err := reg.Register(indexOfFunc); err != nil {
		return err
	}

	// String to number
	stringToNumberFunc := functions.NewFunction(
		"string->number",
		registry.CategoryString,
		1,
		"Convert string to number: (string->number \"42.5\") => 42.5",
		p.evalStringToNumber,
	)
	if err := reg.Register(stringToNumberFunc); err != nil {
		return err
	}

	// Number to string
	numberToStringFunc := functions.NewFunction(
		"number->string",
		registry.CategoryString,
		1,
		"Convert number to string: (number->string 42.5) => \"42.5\"",
		p.evalNumberToString,
	)
	if err := reg.Register(numberToStringFunc); err != nil {
		return err
	}

	// String regex match
	regexMatchFunc := functions.NewFunction(
		"string-regex-match?",
		registry.CategoryString,
		2,
		"Test if string matches regex: (string-regex-match? \"hello\" \"h.*o\") => true",
		p.evalStringRegexMatch,
	)
	if err := reg.Register(regexMatchFunc); err != nil {
		return err
	}

	// String regex find all
	regexFindAllFunc := functions.NewFunction(
		"string-regex-find-all",
		registry.CategoryString,
		2,
		"Find all regex matches: (string-regex-find-all \"abc123def456\" \"[0-9]+\") => (\"123\" \"456\")",
		p.evalStringRegexFindAll,
	)
	if err := reg.Register(regexFindAllFunc); err != nil {
		return err
	}

	// String repeat
	repeatFunc := functions.NewFunction(
		"string-repeat",
		registry.CategoryString,
		2,
		"Repeat string: (string-repeat \"Hi\" 3) => \"HiHiHi\"",
		p.evalStringRepeat,
	)
	if err := reg.Register(repeatFunc); err != nil {
		return err
	}

	// String predicate
	stringPredicateFunc := functions.NewFunction(
		"string?",
		registry.CategoryString,
		1,
		"Check if value is string: (string? \"hello\") => true",
		p.evalStringPredicate,
	)
	if err := reg.Register(stringPredicateFunc); err != nil {
		return err
	}

	// String empty predicate
	stringEmptyFunc := functions.NewFunction(
		"string-empty?",
		registry.CategoryString,
		1,
		"Check if string is empty: (string-empty? \"\") => true",
		p.evalStringEmpty,
	)
	if err := reg.Register(stringEmptyFunc); err != nil {
		return err
	}

	// Clojure-style aliases
	// str function (alias for string-concat)
	strFunc := functions.NewFunction(
		"str",
		registry.CategoryString,
		-1, // Variadic
		"Concatenate strings: (str \"Hello\" \" \" \"World\")",
		p.evalStringConcat,
	)
	if err := reg.Register(strFunc); err != nil {
		return err
	}

	// subs function (alias for string-substring)
	subsFunc := functions.NewFunction(
		"subs",
		registry.CategoryString,
		3,
		"Extract substring: (subs \"Hello\" 1 4) => \"ell\"",
		p.evalStringSubstring,
	)
	if err := reg.Register(subsFunc); err != nil {
		return err
	}

	// split function (alias for string-split)
	clojureSplitFunc := functions.NewFunction(
		"split",
		registry.CategoryString,
		2,
		"Split string: (split \"a,b,c\" \",\") => (\"a\" \"b\" \"c\")",
		p.evalStringSplit,
	)
	if err := reg.Register(clojureSplitFunc); err != nil {
		return err
	}

	// join function (alias for string-join)
	clojureJoinFunc := functions.NewFunction(
		"join",
		registry.CategoryString,
		2,
		"Join strings: (join \",\" (list \"a\" \"b\" \"c\")) => \"a,b,c\"",
		p.evalStringJoin,
	)
	if err := reg.Register(clojureJoinFunc); err != nil {
		return err
	}

	// replace function (alias for string-replace)
	clojureReplaceFunc := functions.NewFunction(
		"replace",
		registry.CategoryString,
		3,
		"Replace all occurrences: (replace \"hello\" \"l\" \"x\") => \"hexxo\"",
		p.evalStringReplace,
	)
	if err := reg.Register(clojureReplaceFunc); err != nil {
		return err
	}

	// trim function (alias for string-trim)
	clojureTrimFunc := functions.NewFunction(
		"trim",
		registry.CategoryString,
		1,
		"Trim whitespace: (trim \"  hello  \") => \"hello\"",
		p.evalStringTrim,
	)
	if err := reg.Register(clojureTrimFunc); err != nil {
		return err
	}

	// upper-case function (alias for string-upper)
	upperCaseFunc := functions.NewFunction(
		"upper-case",
		registry.CategoryString,
		1,
		"Convert to uppercase: (upper-case \"hello\") => \"HELLO\"",
		p.evalStringUpper,
	)
	if err := reg.Register(upperCaseFunc); err != nil {
		return err
	}

	// lower-case function (alias for string-lower)
	lowerCaseFunc := functions.NewFunction(
		"lower-case",
		registry.CategoryString,
		1,
		"Convert to lowercase: (lower-case \"HELLO\") => \"hello\"",
		p.evalStringLower,
	)
	if err := reg.Register(lowerCaseFunc); err != nil {
		return err
	}

	return nil
}

// String concatenation
func (p *StringPlugin) evalStringConcat(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return types.StringValue(""), nil
	}

	var result strings.Builder
	for _, arg := range args {
		value, err := evaluator.Eval(arg)
		if err != nil {
			return nil, err
		}

		result.WriteString(value.String())
	}

	return types.StringValue(result.String()), nil
}

// String length
func (p *StringPlugin) evalStringLength(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-length requires exactly 1 argument")
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	stringValue, ok := value.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-length argument must be a string, got %T", value)
	}

	return types.NumberValue(float64(len(string(stringValue)))), nil
}

// String substring
func (p *StringPlugin) evalStringSubstring(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("string-substring requires exactly 3 arguments: string, start, end")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	startValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}
	endValue, err := evaluator.Eval(args[2])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-substring first argument must be a string, got %T", strValue)
	}

	start, ok := startValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("string-substring second argument must be a number, got %T", startValue)
	}

	end, ok := endValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("string-substring third argument must be a number, got %T", endValue)
	}

	s := string(str)
	startIdx := int(start)
	endIdx := int(end)

	if startIdx < 0 || endIdx < 0 || startIdx > len(s) || endIdx > len(s) || startIdx > endIdx {
		return nil, fmt.Errorf("string-substring indices out of bounds")
	}

	return types.StringValue(s[startIdx:endIdx]), nil
}

// String character at index
func (p *StringPlugin) evalStringCharAt(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-char-at requires exactly 2 arguments: string, index")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	indexValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-char-at first argument must be a string, got %T", strValue)
	}

	index, ok := indexValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("string-char-at second argument must be a number, got %T", indexValue)
	}

	s := string(str)
	idx := int(index)

	if idx < 0 || idx >= len(s) {
		return nil, fmt.Errorf("string-char-at index out of bounds")
	}

	return types.StringValue(string(s[idx])), nil
}

// String uppercase
func (p *StringPlugin) evalStringUpper(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-upper requires exactly 1 argument")
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	stringValue, ok := value.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-upper argument must be a string, got %T", value)
	}

	return types.StringValue(strings.ToUpper(string(stringValue))), nil
}

// String lowercase
func (p *StringPlugin) evalStringLower(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-lower requires exactly 1 argument")
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	stringValue, ok := value.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-lower argument must be a string, got %T", value)
	}

	return types.StringValue(strings.ToLower(string(stringValue))), nil
}

// String trim whitespace
func (p *StringPlugin) evalStringTrim(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-trim requires exactly 1 argument")
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	stringValue, ok := value.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-trim argument must be a string, got %T", value)
	}

	return types.StringValue(strings.TrimSpace(string(stringValue))), nil
}

// String split
func (p *StringPlugin) evalStringSplit(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-split requires exactly 2 arguments: string, separator")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	sepValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-split first argument must be a string, got %T", strValue)
	}

	sep, ok := sepValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-split second argument must be a string, got %T", sepValue)
	}

	parts := strings.Split(string(str), string(sep))
	var elements []types.Value
	for _, part := range parts {
		elements = append(elements, types.StringValue(part))
	}

	return &types.ListValue{Elements: elements}, nil
}

// String join
func (p *StringPlugin) evalStringJoin(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-join requires exactly 2 arguments: list, separator")
	}

	listValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	sepValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	list, ok := listValue.(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("string-join first argument must be a list, got %T", listValue)
	}

	sep, ok := sepValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-join second argument must be a string, got %T", sepValue)
	}

	var parts []string
	for _, elem := range list.Elements {
		parts = append(parts, elem.String())
	}

	return types.StringValue(strings.Join(parts, string(sep))), nil
}

// String contains
func (p *StringPlugin) evalStringContains(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-contains? requires exactly 2 arguments: string, substring")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	substrValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-contains? first argument must be a string, got %T", strValue)
	}

	substr, ok := substrValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-contains? second argument must be a string, got %T", substrValue)
	}

	return types.BooleanValue(strings.Contains(string(str), string(substr))), nil
}

// String starts with
func (p *StringPlugin) evalStringStartsWith(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-starts-with? requires exactly 2 arguments: string, prefix")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	prefixValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-starts-with? first argument must be a string, got %T", strValue)
	}

	prefix, ok := prefixValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-starts-with? second argument must be a string, got %T", prefixValue)
	}

	return types.BooleanValue(strings.HasPrefix(string(str), string(prefix))), nil
}

// String ends with
func (p *StringPlugin) evalStringEndsWith(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-ends-with? requires exactly 2 arguments: string, suffix")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	suffixValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-ends-with? first argument must be a string, got %T", strValue)
	}

	suffix, ok := suffixValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-ends-with? second argument must be a string, got %T", suffixValue)
	}

	return types.BooleanValue(strings.HasSuffix(string(str), string(suffix))), nil
}

// String replace
func (p *StringPlugin) evalStringReplace(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("string-replace requires exactly 3 arguments: string, old, new")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	oldValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}
	newValue, err := evaluator.Eval(args[2])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-replace first argument must be a string, got %T", strValue)
	}

	old, ok := oldValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-replace second argument must be a string, got %T", oldValue)
	}

	new, ok := newValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-replace third argument must be a string, got %T", newValue)
	}

	return types.StringValue(strings.ReplaceAll(string(str), string(old), string(new))), nil
}

// String index of
func (p *StringPlugin) evalStringIndexOf(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-index-of requires exactly 2 arguments: string, substring")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	substrValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-index-of first argument must be a string, got %T", strValue)
	}

	substr, ok := substrValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-index-of second argument must be a string, got %T", substrValue)
	}

	index := strings.Index(string(str), string(substr))
	return types.NumberValue(float64(index)), nil
}

// String to number conversion
func (p *StringPlugin) evalStringToNumber(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string->number requires exactly 1 argument")
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	stringValue, ok := value.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string->number argument must be a string, got %T", value)
	}

	num, err := strconv.ParseFloat(string(stringValue), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number format: %s", string(stringValue))
	}

	return types.NumberValue(num), nil
}

// Number to string conversion
func (p *StringPlugin) evalNumberToString(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("number->string requires exactly 1 argument")
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	numberValue, ok := value.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("number->string argument must be a number, got %T", value)
	}

	return types.StringValue(strconv.FormatFloat(float64(numberValue), 'g', -1, 64)), nil
}

// Regular expression match
func (p *StringPlugin) evalStringRegexMatch(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-regex-match? requires exactly 2 arguments: string, pattern")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	patternValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-regex-match? first argument must be a string, got %T", strValue)
	}

	pattern, ok := patternValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-regex-match? second argument must be a string, got %T", patternValue)
	}

	re, err := regexp.Compile(string(pattern))
	if err != nil {
		return nil, fmt.Errorf("string-regex-match?: invalid regex pattern '%s': %v", string(pattern), err)
	}

	return types.BooleanValue(re.MatchString(string(str))), nil
}

// Regular expression find all
func (p *StringPlugin) evalStringRegexFindAll(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-regex-find-all requires exactly 2 arguments: string, pattern")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	patternValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-regex-find-all first argument must be a string, got %T", strValue)
	}

	pattern, ok := patternValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-regex-find-all second argument must be a string, got %T", patternValue)
	}

	re, err := regexp.Compile(string(pattern))
	if err != nil {
		return nil, fmt.Errorf("string-regex-find-all: invalid regex pattern '%s': %v", string(pattern), err)
	}

	matches := re.FindAllString(string(str), -1)
	var elements []types.Value
	for _, match := range matches {
		elements = append(elements, types.StringValue(match))
	}

	return &types.ListValue{Elements: elements}, nil
}

// String repeat
func (p *StringPlugin) evalStringRepeat(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-repeat requires exactly 2 arguments: string, count")
	}

	strValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}
	countValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	str, ok := strValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-repeat first argument must be a string, got %T", strValue)
	}

	count, ok := countValue.(types.NumberValue)
	if !ok {
		return nil, fmt.Errorf("string-repeat second argument must be a number, got %T", countValue)
	}

	if count < 0 {
		return nil, fmt.Errorf("string-repeat count must be non-negative")
	}

	return types.StringValue(strings.Repeat(string(str), int(count))), nil
}

// String predicate - check if value is string
func (p *StringPlugin) evalStringPredicate(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string? requires exactly 1 argument")
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	_, isString := value.(types.StringValue)
	return types.BooleanValue(isString), nil
}

// String empty predicate
func (p *StringPlugin) evalStringEmpty(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-empty? requires exactly 1 argument")
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	stringValue, ok := value.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-empty? argument must be a string, got %T", value)
	}

	return types.BooleanValue(len(string(stringValue)) == 0), nil
}
