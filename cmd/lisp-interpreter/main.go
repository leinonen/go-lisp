package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/leinonen/lisp-interpreter/pkg/executor"
	"github.com/leinonen/lisp-interpreter/pkg/interpreter"
	"github.com/leinonen/lisp-interpreter/pkg/repl"
)

func main() {
	interpreter := interpreter.NewInterpreter()

	// Check if a file argument was provided
	if len(os.Args) > 1 {
		filename := os.Args[1]
		err := executor.ExecuteFile(interpreter, filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file %s: %v\n", filename, err)
			os.Exit(1)
		}
		return
	}

	// If no file argument, start REPL
	scanner := bufio.NewScanner(os.Stdin)
	repl.REPL(interpreter, scanner)
}
