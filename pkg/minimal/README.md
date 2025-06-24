# Core Architecture Overview

1. Minimal Core Interpreter (Kernel)
Build a tiny, trusted set of primitive operations. This includes:

- Symbols and interning
- Lists (linked or array-based)
- Environments (lexical scope + symbol table)
- Eval/apply logic
- Basic special forms: if, define, fn, quote, do
- A minimal REPL

💡 Design principle: The core should be so small that you could hold it in your head. Think of this like a "Lisp microkernel".

2. Bootstrap Language in Itself
Implement higher-level constructs using the language itself:

- Build macros using the core
- Extend with user-defined functions
- Implement conditionals, loops, data structures, and module systems in the language

This makes the language self-hosting and keeps the core clean.

3. Macro System (Hygienic Optional)
Macros are Lisp's superpower. Design a macro system that:

- Can manipulate code-as-data (quasiquote, unquote, etc.)
- Optionally supports hygienic macros (via symbol renaming or gensym)
- This makes your language infinitely extensible.