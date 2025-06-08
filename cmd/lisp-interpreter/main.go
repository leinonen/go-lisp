package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/leinonen/lisp-interpreter/pkg/executor"
	"github.com/leinonen/lisp-interpreter/pkg/interpreter"
	"github.com/leinonen/lisp-interpreter/pkg/repl"
)

func main() {
	var (
		help     = flag.Bool("help", false, "Show help message")
		eval     = flag.String("e", "", "Evaluate code directly instead of reading from a file")
		filename = flag.String("f", "", "File to execute")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                     # Start interactive REPL\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f script.lisp      # Execute a file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -e '(+ 1 2 3)'      # Evaluate code directly\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -help               # Show this help message\n", os.Args[0])
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	interpreter := interpreter.NewInterpreter()

	// Handle -e flag: evaluate code directly
	if *eval != "" {
		result, err := interpreter.Interpret(*eval)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error evaluating code: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
		return
	}

	// Handle -f flag: execute a file
	if *filename != "" {
		err := executor.ExecuteFile(interpreter, *filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file %s: %v\n", *filename, err)
			os.Exit(1)
		}
		return
	}

	// Check for legacy positional argument (backward compatibility)
	if len(flag.Args()) > 0 {
		legacyFilename := flag.Args()[0]
		err := executor.ExecuteFile(interpreter, legacyFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file %s: %v\n", legacyFilename, err)
			os.Exit(1)
		}
		return
	}

	// If no arguments provided, start REPL
	scanner := bufio.NewScanner(os.Stdin)
	repl.REPL(interpreter, scanner)
}
