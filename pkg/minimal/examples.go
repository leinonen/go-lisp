package minimal

// Example demonstrates the minimal kernel architecture from future.md

import "fmt"

// RunExamples demonstrates the minimal kernel and bootstrap concepts
func RunExamples() {
	fmt.Println("=== Minimal Lisp Kernel Examples ===")
	fmt.Println("Based on the architecture described in future.md")

	// Create environment with minimal core
	repl := NewREPL()
	env := repl.Env

	fmt.Println("1. Minimal Core Primitives:")
	fmt.Println("   - Symbols and interning")
	fmt.Println("   - Lists")
	fmt.Println("   - Environment (lexical scope)")
	fmt.Println("   - Eval/apply logic")
	fmt.Println("   - Basic special forms: if, def, fn, quote, do")

	// Test core special forms
	fmt.Println("\n2. Testing Core Special Forms:")

	// Quote
	fmt.Println("   (quote hello) =>",
		evalAndPrint(NewList(Intern("quote"), Intern("hello")), env))

	// Define
	evalAndPrint(NewList(Intern("def"), Intern("x"), Number(42)), env)
	fmt.Println("   (def x 42) => defined")

	// Symbol lookup
	fmt.Println("   x =>", evalAndPrint(Intern("x"), env))

	// If
	fmt.Println("   (if true \"yes\" \"no\") =>",
		evalAndPrint(NewList(Intern("if"), Boolean(true), String("yes"), String("no")), env))

	// Function definition
	fmt.Println("\n3. User-Defined Functions (Closures):")

	// (def add (fn [a b] (+ a b)))
	addFn := NewList(
		Intern("fn"),
		NewVector(Intern("a"), Intern("b")),
		NewList(Intern("+"), Intern("a"), Intern("b")),
	)
	evalAndPrint(NewList(Intern("def"), Intern("add"), addFn), env)
	fmt.Println("   (def add (fn [a b] (+ a b))) => defined")

	// Bootstrap higher-level functions
	fmt.Println("\n4. Bootstrap Language in Itself:")
	Bootstrap(env)

	// Use arithmetic
	fmt.Println("   (add 3 4) =>",
		evalAndPrint(NewList(Intern("add"), Number(3), Number(4)), env))

	// Test bootstrapped functions
	fmt.Println("   (list 1 2 3) =>",
		evalAndPrint(NewList(Intern("list"), Number(1), Number(2), Number(3)), env))

	fmt.Println("   (first (list 1 2 3)) =>",
		evalAndPrint(NewList(Intern("first"),
			NewList(Intern("list"), Number(1), Number(2), Number(3))), env))

	// Show higher-level construct built from primitives
	fmt.Println("\n5. Higher-Level Constructs:")

	// Define factorial using recursion
	factorialFn := NewList(
		Intern("fn"),
		NewVector(Intern("n")),
		NewList(
			Intern("if"),
			NewList(Intern("<"), Intern("n"), Number(2)),
			Number(1),
			NewList(Intern("*"), Intern("n"),
				NewList(Intern("factorial"),
					NewList(Intern("-"), Intern("n"), Number(1)))),
		),
	)

	evalAndPrint(NewList(Intern("def"), Intern("factorial"), factorialFn), env)
	fmt.Println("   (def factorial (fn [n] ...)) => defined")

	fmt.Println("   (factorial 5) =>",
		evalAndPrint(NewList(Intern("factorial"), Number(5)), env))

	fmt.Println("\n6. Key Architectural Benefits:")
	fmt.Println("   ✓ Minimal core (< 200 lines core eval logic)")
	fmt.Println("   ✓ Self-hosting potential")
	fmt.Println("   ✓ Clean separation: core vs. library")
	fmt.Println("   ✓ Easy to extend and test")
	fmt.Println("   ✓ Ready for macro system")

	fmt.Println("\n=== Ready for Next Steps ===")
	fmt.Println("• Macro system (quasiquote, unquote)")
	fmt.Println("• Module system")
	fmt.Println("• Advanced data structures")
	fmt.Println("• Integration with existing codebase")
}

func evalAndPrint(expr Value, env *Environment) string {
	result, err := Eval(expr, env)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return result.String()
}
