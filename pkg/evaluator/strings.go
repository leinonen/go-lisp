// Package evaluator_strings contains string manipulation and processing functionality
package evaluator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// String manipulation functions

// String concatenation
func (e *Evaluator) evalStringConcat(args []types.Expr) (types.Value, error) {
	if len(args) == 0 {
		return types.StringValue(""), nil
	}

	var result strings.Builder
	for _, arg := range args {
		value, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}

		// Convert value to string representation
		result.WriteString(value.String())
	}

	return types.StringValue(result.String()), nil
}

// String length
func (e *Evaluator) evalStringLength(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-length requires exactly 1 argument")
	}

	value, err := e.Eval(args[0])
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
func (e *Evaluator) evalStringSubstring(args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("string-substring requires exactly 3 arguments: string, start, end")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	startValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}
	endValue, err := e.Eval(args[2])
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
func (e *Evaluator) evalStringCharAt(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-char-at requires exactly 2 arguments: string, index")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	indexValue, err := e.Eval(args[1])
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
func (e *Evaluator) evalStringUpper(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-upper requires exactly 1 argument")
	}

	value, err := e.Eval(args[0])
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
func (e *Evaluator) evalStringLower(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-lower requires exactly 1 argument")
	}

	value, err := e.Eval(args[0])
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
func (e *Evaluator) evalStringTrim(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-trim requires exactly 1 argument")
	}

	value, err := e.Eval(args[0])
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
func (e *Evaluator) evalStringSplit(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-split requires exactly 2 arguments: string, separator")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	sepValue, err := e.Eval(args[1])
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
func (e *Evaluator) evalStringJoin(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-join requires exactly 2 arguments: list, separator")
	}

	listValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	sepValue, err := e.Eval(args[1])
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
func (e *Evaluator) evalStringContains(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-contains? requires exactly 2 arguments: string, substring")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	substrValue, err := e.Eval(args[1])
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
func (e *Evaluator) evalStringStartsWith(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-starts-with? requires exactly 2 arguments: string, prefix")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	prefixValue, err := e.Eval(args[1])
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
func (e *Evaluator) evalStringEndsWith(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-ends-with? requires exactly 2 arguments: string, suffix")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	suffixValue, err := e.Eval(args[1])
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
func (e *Evaluator) evalStringReplace(args []types.Expr) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("string-replace requires exactly 3 arguments: string, old, new")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	oldValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}
	newValue, err := e.Eval(args[2])
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
func (e *Evaluator) evalStringIndexOf(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-index-of requires exactly 2 arguments: string, substring")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	substrValue, err := e.Eval(args[1])
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
func (e *Evaluator) evalStringToNumber(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string->number requires exactly 1 argument")
	}

	value, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	stringValue, ok := value.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string->number argument must be a string, got %T", value)
	}

	num, err := strconv.ParseFloat(string(stringValue), 64)
	if err != nil {
		return nil, fmt.Errorf("string->number: invalid number format '%s'", string(stringValue))
	}

	return types.NumberValue(num), nil
}

// Number to string conversion
func (e *Evaluator) evalNumberToString(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("number->string requires exactly 1 argument")
	}

	value, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	switch v := value.(type) {
	case types.NumberValue:
		return types.StringValue(strconv.FormatFloat(float64(v), 'g', -1, 64)), nil
	case *types.BigNumberValue:
		return types.StringValue(v.String()), nil
	default:
		return nil, fmt.Errorf("number->string argument must be a number, got %T", value)
	}
}

// Regular expression match
func (e *Evaluator) evalStringRegexMatch(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-regex-match? requires exactly 2 arguments: string, pattern")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	patternValue, err := e.Eval(args[1])
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

	matched, err := regexp.MatchString(string(pattern), string(str))
	if err != nil {
		return nil, fmt.Errorf("string-regex-match?: invalid regex pattern '%s': %v", string(pattern), err)
	}

	return types.BooleanValue(matched), nil
}

// Regular expression find all
func (e *Evaluator) evalStringRegexFindAll(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-regex-find-all requires exactly 2 arguments: string, pattern")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	patternValue, err := e.Eval(args[1])
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
func (e *Evaluator) evalStringRepeat(args []types.Expr) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-repeat requires exactly 2 arguments: string, count")
	}

	strValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	countValue, err := e.Eval(args[1])
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
func (e *Evaluator) evalStringPredicate(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string? requires exactly 1 argument")
	}

	value, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	_, isString := value.(types.StringValue)
	return types.BooleanValue(isString), nil
}

// String empty predicate
func (e *Evaluator) evalStringEmpty(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string-empty? requires exactly 1 argument")
	}

	value, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	stringValue, ok := value.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string-empty? argument must be a string, got %T", value)
	}

	return types.BooleanValue(len(string(stringValue)) == 0), nil
}
