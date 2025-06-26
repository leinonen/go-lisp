package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/leinonen/go-lisp/pkg/minimal"
)

func main() {
	// Define command line flags
	var filename string
	flag.StringVar(&filename, "f", "", "Load and execute a Lisp file")
	flag.Parse()

	// Handle legacy "examples" argument for backward compatibility
	if len(os.Args) > 1 && os.Args[1] == "examples" {
		// Run examples to demonstrate the architecture
		minimal.RunExamples()
		return
	}

	// Create a bootstrapped REPL environment
	repl := minimal.NewBootstrappedREPL()

	// If a filename is provided, load and execute it
	if filename != "" {
		fmt.Printf("Loading file: %s\n", filename)
		_, err := minimal.LoadFile(filename, repl.Env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading file '%s': %v\n", filename, err)
			os.Exit(1)
		}

		fmt.Println("File loaded successfully.")
		return
	}

	// Start interactive REPL if no file was specified
	fmt.Println("Starting Minimal Lisp REPL")
	fmt.Println("This demonstrates the micro kernel architecture")
	fmt.Println("Use -f <filename> to load and execute a file")
	fmt.Println()

	repl.Run()
}
