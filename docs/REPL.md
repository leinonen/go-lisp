# GoLisp Enhanced REPL

The GoLisp REPL (Read-Eval-Print-Loop) provides a sophisticated interactive development environment with modern features for productive Lisp programming.

## Quick Start

```bash
# Start the REPL
./bin/golisp

# You'll see:
GoLisp Enhanced REPL
Type 'exit' or 'quit' to quit
Multi-line expressions supported - press Enter on incomplete expressions
Type ')' on empty line during multi-line input to force evaluation
GoLisp>
```

## Core Features

### Multi-line Expression Support

The REPL automatically detects incomplete expressions and seamlessly continues on multiple lines:

```lisp
GoLisp> (defn factorial [n]
      >   (if (= n 0)
      >       1
      >       (* n (factorial (- n 1)))))
#<function>

GoLisp> (factorial 5)
120
```

**How it works:**
- **Smart Detection**: Parses parentheses, brackets `[]`, and braces `{}`
- **String Handling**: Ignores parentheses inside string literals
- **Comment Support**: Handles Lisp comments (`;` to end of line)
- **Escape Sequences**: Properly handles escaped quotes in strings
- **Visual Feedback**: Prompt changes to `      > ` for continuation lines

### Dynamic Auto-completion

Tab completion is environment-aware and includes all loaded functions:

```lisp
GoLisp> (ma<TAB>         ; Completes to "(map"
GoLisp> (fi<TAB>         ; Completes to "(filter"  
GoLisp> (red<TAB>        ; Completes to "(reduce"
GoLisp> nil<TAB>         ; Completes to "nil" (no parentheses)
```

**Available Completions:**
- **117+ Symbols**: All built-in and standard library functions
- **Smart Categories**:
  - Functions: `(map`, `(filter`, `(+`, `(cons`
  - Special forms: `(def`, `(defn`, `(if`, `(fn`
  - Literals: `nil`, `true`, `false` (no parentheses)
  - REPL commands: `exit`, `quit`

**Features:**
- **Environment-Aware**: Only shows actually loaded functions
- **Real-time Updates**: Framework in place for updating after new definitions
- **Proper Syntax**: Automatically adds opening parenthesis for functions

### History and Navigation

Professional terminal handling with full navigation support:

- **↑/↓ arrows**: Navigate through command history
- **←/→ arrows**: Move cursor within the current line
- **Home/End**: Jump to beginning/end of line
- **Ctrl+A/E**: Alternative home/end navigation
- **Ctrl+U**: Clear current line
- **Ctrl+C**: Cancel multi-line input or exit REPL
- **Ctrl+D/EOF**: Exit REPL

### Force Evaluation

For incomplete expressions, you can force evaluation:

```lisp
GoLisp> (map + (list 1 2 3
      >        (list 4 5 6          ; Missing closing parens
      > )                           ; Force evaluation
[5 7 9]
```

**How it works:**
- Type `)` on an empty line during multi-line input
- Automatically calculates and adds missing closing parentheses
- Only works when there's actual content to evaluate
- Safe: won't create invalid expressions

## Advanced Features

### Error Handling

The REPL provides comprehensive error feedback:

```lisp
GoLisp> )
Error: Unexpected closing parenthesis

GoLisp> (+ 1 "hello")
TypeError: + expects numbers, got core.String

GoLisp> (def)
ArityError: def expects 2 arguments, got 0

GoLisp> undefined-function
NameError: undefined symbol: undefined-function
```

### Comment Support

Full support for Lisp comments in multi-line expressions:

```lisp
GoLisp> (defn factorial [n]
      >   ; Base case
      >   (if (= n 0)
      >       1
      >       ; Recursive case  
      >       (* n (factorial (- n 1)))))
#<function>
```

### Complex Data Structures

Multi-line support works with all GoLisp data types:

```lisp
GoLisp> {:users [{:name "Alice" :age 30}
      >          {:name "Bob" :age 25}]
      >  :count 2}
{:users [{:name "Alice" :age 30} {:name "Bob" :age 25}] :count 2}

GoLisp> (vector (range 5)
      >         (map (fn [x] (* x x)) (range 5)))
[[0 1 2 3 4] [0 1 4 9 16]]
```

## Technical Implementation

### Parentheses Balancing Algorithm

The REPL uses a sophisticated balancing algorithm:

1. **State Tracking**: Maintains state for strings, comments, and escape sequences
2. **Bracket Counting**: Tracks `()`, `[]`, `{}` with proper nesting
3. **Comment Handling**: Ignores content after `;` until newline
4. **String Literals**: Ignores brackets inside `"..."` strings
5. **Escape Sequences**: Handles `\"` and `\\` properly

### Content Detection

Distinguishes between meaningful content and whitespace/comments:

```lisp
# These don't trigger evaluation:
GoLisp>           ; Just whitespace
GoLisp> ; Just a comment

# These do:
GoLisp> 42        ; Number literal
GoLisp> "hello"   ; String literal
GoLisp> (+ 1 2)   ; Expression
```

### Auto-completion System

1. **Environment Introspection**: Queries the Lisp environment for all available symbols
2. **Dynamic Generation**: Creates completion list from actually loaded functions
3. **Smart Categorization**: Applies appropriate prefixes based on symbol type
4. **Performance**: Efficient lookup with sorted symbol lists

## Tips and Best Practices

### Productive REPL Usage

1. **Experiment Freely**: Use the REPL to test expressions before putting them in files
2. **Multi-line Functions**: Write complex functions directly in the REPL
3. **Tab Completion**: Use tab completion to discover available functions
4. **History Navigation**: Use arrow keys to recall and modify previous expressions
5. **Force Evaluation**: Use `)` trick when you have unbalanced expressions

### Common Patterns

```lisp
# Define and test functions interactively
GoLisp> (defn square [x] (* x x))
GoLisp> (square 5)
25

# Test with various inputs
GoLisp> (map square [1 2 3 4 5])
[1 4 9 16 25]

# Iterate and refine
GoLisp> (defn square [x] 
      >   (println "Squaring" x)
      >   (* x x))
```

### Error Recovery

```lisp
# Fix syntax errors easily
GoLisp> (+ 1 2
      > 3)                 ; Fix missing parenthesis
6

# Cancel problematic input
GoLisp> (some-complex-expression
      > ^C                 ; Ctrl+C to cancel
Cancelled
GoLisp>
```

## Comparison with Other REPLs

| Feature | GoLisp REPL | Basic REPL | Notes |
|---------|-------------|------------|--------|
| Multi-line | ✅ Automatic | ❌ Manual | Smart parentheses detection |
| History | ✅ Arrow keys | ❌ None | Full terminal support |
| Autocomplete | ✅ Dynamic | ❌ None | Environment-aware |
| Error handling | ✅ Professional | ❌ Basic | Categorized errors |
| Force evaluation | ✅ Smart | ❌ None | Safe completion |
| Comments | ✅ Full support | ⚠️ Limited | Multi-line comments |

## Future Enhancements

The REPL is designed for extensibility. Planned features include:

- **Real-time Completion Updates**: Immediate updates after `def`/`defn`
- **Syntax Highlighting**: Colorized input for better readability  
- **Bracket Matching**: Visual matching of corresponding brackets
- **Documentation Integration**: Inline help for functions
- **Session Save/Load**: Persist REPL sessions to files
- **Debugging Integration**: Step-through debugging in REPL

## Troubleshooting

### Common Issues

**Q: Tab completion doesn't work**
A: Ensure you're using a terminal that supports readline. Most modern terminals (bash, zsh, PowerShell Core) work fine.

**Q: Arrow keys show strange characters**
A: Your terminal may not support ANSI escape sequences. Try a different terminal or SSH to a Linux system.

**Q: Multi-line input gets stuck**
A: Use Ctrl+C to cancel, or type `)` on an empty line to force evaluation.

**Q: History doesn't persist between sessions**
A: This is by design. History is kept in memory during the session but not saved to disk.

### Platform Notes

- **Linux/macOS**: Full feature support
- **Windows**: Basic support; may have limitations with advanced terminal features
- **SSH**: Full support when connecting to Linux systems

## Contributing

The REPL implementation is thoroughly tested with 79 unit tests covering:

- Parentheses balancing (33 test cases)
- Content detection (19 test cases)  
- Bracket counting (11 test cases)
- Evaluation logic (16 test cases)
- Integration tests (REPL creation, evaluation)

When contributing REPL enhancements:

1. Add comprehensive tests for new functionality
2. Ensure backward compatibility
3. Test on multiple platforms
4. Update this documentation

## Related Documentation

- **[CLAUDE.md](../CLAUDE.md)**: Complete project architecture
- **[README.md](../README.md)**: Project overview and examples
- **[examples.md](examples.md)**: Comprehensive function examples

The GoLisp Enhanced REPL provides a modern, productive environment for interactive Lisp development with features that rival professional development tools.