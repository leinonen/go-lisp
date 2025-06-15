package interpreter

import (
	"github.com/leinonen/lisp-interpreter/pkg/evaluator"
	"github.com/leinonen/lisp-interpreter/pkg/parser"
	concurrencyplugin "github.com/leinonen/lisp-interpreter/pkg/plugins/concurrency"
	"github.com/leinonen/lisp-interpreter/pkg/pure"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// Interpreter combines tokenizer, parser, and evaluator
type Interpreter struct {
	env       *evaluator.Environment
	evaluator registry.Evaluator
}

// Ensure Interpreter implements the InterpreterDependency interface
var _ concurrencyplugin.InterpreterDependency = (*Interpreter)(nil)

func NewInterpreter() (*Interpreter, error) {
	// Use modular evaluator by default now
	return NewModularInterpreter()
}

// NewModularInterpreter creates an interpreter that uses the pure plugin-based evaluator
func NewModularInterpreter() (*Interpreter, error) {
	env := evaluator.NewEnvironment()
	pureEval, err := pure.NewPureEvaluator(env)
	if err != nil {
		return nil, err
	}

	interp := &Interpreter{
		env:       env,
		evaluator: pureEval,
	}

	// Set the interpreter as a dependency for plugins that need it
	pureEval.SetInterpreterDependency(interp)

	return interp, nil
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

	// Evaluate using modular evaluator (always available now)
	return i.evaluator.Eval(ast)
}

// GetEnvironment returns the interpreter's environment for completion support
func (i *Interpreter) GetEnvironment() interface{} {
	return i.env
}
