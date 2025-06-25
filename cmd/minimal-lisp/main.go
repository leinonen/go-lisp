package main

import (
	"fmt"
	"os"

	"github.com/leinonen/go-lisp/pkg/minimal"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "examples" {
		// Run examples to demonstrate the architecture
		minimal.RunExamples()
		return
	}

	fmt.Println("Starting Minimal Lisp REPL")
	fmt.Println("This demonstrates the micro kernel architecture")
	fmt.Println()

	// Create and run the minimal REPL
	repl := minimal.NewBootstrappedREPL()
	repl.Run()
}
