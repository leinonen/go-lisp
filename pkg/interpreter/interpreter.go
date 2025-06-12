package interpreter

import (
	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Interpreter combines tokenizer, parser, and evaluator
type Interpreter struct {
	env *evaluator.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: evaluator.NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(input string) (types.Value, error) {
	// Tokenize
	tokenizer := tokenizer.NewTokenizer(input)
	tokens, err := tokenizer.TokenizeWithError()
	if err != nil {
		return nil, err
	}

	// Parse
	parser := parser.NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	// Evaluate
	evaluator := evaluator.NewEvaluator(i.env)
	return evaluator.Eval(ast)
}

// GetEnvironment returns the interpreter's environment for completion support
func (i *Interpreter) GetEnvironment() *evaluator.Environment {
	return i.env
}
