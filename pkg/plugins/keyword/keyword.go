// Package keyword provides keyword utility functions for the Lisp interpreter
package keyword

import (
	"fmt"

	"github.com/leinonen/go-lisp/pkg/functions"
	"github.com/leinonen/go-lisp/pkg/plugins"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// KeywordPlugin provides keyword utility functions
type KeywordPlugin struct {
	*plugins.BasePlugin
}

// NewKeywordPlugin creates a new keyword plugin
func NewKeywordPlugin() *KeywordPlugin {
	return &KeywordPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"keyword",
			"1.0.0",
			"Keyword utility functions (keyword, keyword?)",
			[]string{}, // No dependencies
		),
	}
}

// RegisterFunctions registers keyword functions
func (p *KeywordPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// keyword function - convert string to keyword
	keywordFunc := functions.NewFunction(
		"keyword",
		registry.CategoryCore,
		1,
		"Convert string to keyword: (keyword \"name\") => :name",
		p.evalKeyword,
	)
	if err := reg.Register(keywordFunc); err != nil {
		return err
	}

	// keyword? predicate
	keywordPredicateFunc := functions.NewFunction(
		"keyword?",
		registry.CategoryCore,
		1,
		"Check if value is a keyword: (keyword? :name) => true",
		p.evalKeywordPredicate,
	)
	return reg.Register(keywordPredicateFunc)
}

// evalKeyword converts a string to a keyword
func (p *KeywordPlugin) evalKeyword(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("keyword requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	switch v := value.(type) {
	case types.StringValue:
		str := string(v)
		// Remove leading : if present
		if len(str) > 0 && str[0] == ':' {
			str = str[1:]
		}
		return types.KeywordValue(str), nil
	case types.KeywordValue:
		// Already a keyword, return as-is
		return v, nil
	default:
		return nil, fmt.Errorf("keyword requires a string argument, got %T", value)
	}
}

// evalKeywordPredicate checks if a value is a keyword
func (p *KeywordPlugin) evalKeywordPredicate(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("keyword? requires exactly 1 argument, got %d", len(args))
	}

	value, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	_, isKeyword := value.(types.KeywordValue)
	return types.BooleanValue(isKeyword), nil
}
