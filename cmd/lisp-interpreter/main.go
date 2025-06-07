package main

import (
	"bufio"
	"os"

	"github.com/leinonen/lisp-interpreter/pkg/interpreter"
	"github.com/leinonen/lisp-interpreter/pkg/repl"
)

func main() {
	interpreter := interpreter.NewInterpreter()
	scanner := bufio.NewScanner(os.Stdin)

	repl.REPL(interpreter, scanner)
}
