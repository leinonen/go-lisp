package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/leinonen/go-lisp/pkg/core"
)

func main() {
	var filename = flag.String("f", "", "Execute a Lisp file")
	var evalStr = flag.String("e", "", "Evaluate a Lisp expression")
	var help = flag.Bool("h", false, "Show help")
	flag.Parse()

	if *help {
		fmt.Println("GoLisp Minimal Core Interpreter")
		fmt.Println("Usage:")
		fmt.Println("  golisp-core           # Start REPL")
		fmt.Println("  golisp-core -f file   # Execute file")
		fmt.Println("  golisp-core -e expr   # Evaluate expression")
		fmt.Println("  golisp-core -h        # Show help")
		return
	}

	// Create REPL instance
	repl, err := core.NewREPL()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating REPL: %v\n", err)
		os.Exit(1)
	}

	if *evalStr != "" {
		// Evaluate expression mode
		result, err := repl.EvalString(*evalStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result.String())
		return
	}

	if *filename != "" {
		// File execution mode
		err := repl.LoadFile(*filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading file: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// REPL mode
	err = repl.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "REPL error: %v\n", err)
		os.Exit(1)
	}
}