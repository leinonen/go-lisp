package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/leinonen/go-lisp/pkg/kernel"
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

	// Create a bootstrapped REPL environment
	repl := kernel.NewBootstrappedREPL()

	// Handle -e flag: evaluate code directly
	if *eval != "" {
		// Wrap the code in a 'do' block to allow multiple statements
		wrappedCode := "(do " + *eval + ")"

		// Parse the wrapped code
		expr, err := kernel.Parse(wrappedCode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing code: %v\n", err)
			os.Exit(1)
		}

		// Evaluate the parsed expression
		result, err := kernel.Eval(expr, repl.Env)
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
		_, err := kernel.LoadFile(*filename, repl.Env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file %s: %v\n", *filename, err)
			os.Exit(1)
		}
		return
	}

	// Check for legacy positional argument (backward compatibility)
	if len(flag.Args()) > 0 {
		legacyFilename := flag.Args()[0]
		_, err := kernel.LoadFile(legacyFilename, repl.Env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file %s: %v\n", legacyFilename, err)
			os.Exit(1)
		}
		return
	}

	// If no arguments provided, start REPL

	repl.Run()
}
