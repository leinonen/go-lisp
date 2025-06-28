package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/leinonen/go-lisp/pkg/core"
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

	// Create a REPL with bootstrapped environment
	repl, err := core.NewREPL()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating REPL: %v\n", err)
		os.Exit(1)
	}

	// Handle -e flag: evaluate code directly
	if *eval != "" {
		// Evaluate the code directly
		result, err := repl.EvalString(*eval)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error evaluating code: %v\n", err)
			os.Exit(1)
		}

		// Don't print nil values (used by print functions to avoid duplicate output)
		if result != nil && result.String() != "nil" {
			fmt.Println(result)
		}
		return
	}

	// Handle -f flag: execute a file
	if *filename != "" {
		err := repl.LoadFile(*filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file %s: %v\n", *filename, err)
			os.Exit(1)
		}
		return
	}

	// Check for legacy positional argument (backward compatibility)
	if len(flag.Args()) > 0 {
		legacyFilename := flag.Args()[0]
		err := repl.LoadFile(legacyFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file %s: %v\n", legacyFilename, err)
			os.Exit(1)
		}
		return
	}

	// If no arguments provided, start REPL
	err = repl.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "REPL error: %v\n", err)
		os.Exit(1)
	}
}
