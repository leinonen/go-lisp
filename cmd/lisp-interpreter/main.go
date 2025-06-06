package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/leinonen/lisp-interpreter/pkg/interpreter"
)

func main() {
	fmt.Println("Welcome to the Lisp Interpreter!")
	fmt.Println("Type expressions to evaluate them, or 'quit' to exit.")
	fmt.Println("Examples:")
	fmt.Println("  42")
	fmt.Println("  (+ 1 2 3)")
	fmt.Println("  (* (+ 2 3) 4)")
	fmt.Println("  (if (< 3 5) \"yes\" \"no\")")
	fmt.Println()

	interpreter := interpreter.NewInterpreter()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("lisp> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "quit" || input == "exit" {
			break
		}

		result, err := interpreter.Interpret(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("=> %s\n", result.String())
		}
	}

	fmt.Println("Goodbye!")
}
