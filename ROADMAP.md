# GoLisp Self-Hosting Roadmap

## Current Status ✅

Your GoLisp implementation already has a strong foundation for self-hosting:

### Core Infrastructure
- ✅ Evaluator with 900+ lines of robust evaluation logic
- ✅ Macro system with `defmacro`, quasiquote/unquote  
- ✅ File loading system (`LoadFile`)
- ✅ REPL with parsing and evaluation
- ✅ Lexical environments with proper scoping
- ✅ Tail call optimization via loop/recur
- ✅ Bootstrap system extending kernel in Lisp

### Data Types & Special Forms
- ✅ Core data types: symbols, lists, vectors, hash-maps, sets
- ✅ Special forms: `if`, `def`, `fn`, `quote`, `do`, `defmacro`, `loop`, `recur`
- ✅ Arithmetic and comparison operators
- ✅ Collection operations (`first`, `rest`, `cons`, `conj`, etc.)

## Phase 1: Meta-Programming Core ✅

### Completed
- ✅ `eval` - Evaluate data as code
- ✅ `read-string` - Parse string into Lisp data  
- ✅ `slurp` - Read entire file as string
- ✅ `spit` - Write string to file
- ✅ `str` - String concatenation
- ✅ Type predicates: `symbol?`, `string?`, `keyword?`, `list?`, `vector?`
- ✅ Symbol manipulation: `symbol`, `keyword`, `name`
- ✅ `macroexpand` - Expand macros for inspection
- ✅ `gensym` - Generate unique symbols  
- ✅ `list*` - List construction with spread

### Still Needed
- [ ] Error handling improvements

## Phase 2: Enhanced Standard Library 📚

### String Operations
- [ ] `split` - Split strings by delimiter
- [ ] `join` - Join strings with separator  
- [ ] `substring` - Extract substrings
- [ ] `trim` - Remove whitespace
- [ ] `replace` - String replacement

### Collection Operations  
- [ ] `map` - Enhanced version with multiple collections
- [ ] `filter` - Filter by predicate
- [ ] `reduce` - Enhanced reduce (already have basic version)
- [ ] `apply` - Apply function to collection as arguments
- [ ] `sort` - Sort collections
- [ ] `group-by` - Group by key function

### I/O Operations
- [ ] `println` - Print with newline
- [ ] `prn` - Print for reading back
- [ ] File system operations (directory listing, etc.)

## Phase 3: Self-Hosting Compiler 🚀

### Core Compiler (Created in `lisp/self-hosting.lisp`)
- ✅ Basic compilation framework
- ✅ Special form compilation (`def`, `fn`, `if`, `quote`, `do`, `let`)
- ✅ Symbol table management
- ✅ Local variable tracking
- ✅ Function application compilation
- ✅ Vector compilation
- [ ] Macro expansion during compilation
- [ ] Improved multi-expression parsing (read-all)
- [ ] Optimization passes
- [ ] Error reporting with source locations

### Code Generation
- [ ] Generate optimized Lisp code
- [ ] Generate Go code (advanced)
- [ ] Generate bytecode (advanced)
- [ ] Dead code elimination
- [ ] Constant folding

### Compilation Pipeline
- [ ] Multi-file compilation
- [ ] Dependency resolution
- [ ] Module system
- [ ] Package management

## Phase 4: Advanced Self-Hosting 🎯

### Performance Optimizations
- [ ] Inline function calls
- [ ] Tail call optimization (beyond loop/recur)
- [ ] Constant propagation
- [ ] Type inference
- [ ] Memory optimization

### Development Tools
- [ ] Debugger written in GoLisp
- [ ] Profiler written in GoLisp
- [ ] Documentation generator
- [ ] Test framework
- [ ] Package manager

### Language Extensions
- [ ] Pattern matching
- [ ] Advanced destructuring
- [ ] Async/await constructs
- [ ] Exception handling
- [ ] Namespaces/modules

## Implementation Strategy 💡

### Step 1: Verify Meta-Programming Core
```bash
cd /home/leinonen/code/go-lisp
make build
./bin/golisp -e "(eval '(+ 1 2 3))"  # Should output 6
./bin/golisp -e "(symbol? 'hello)"   # Should output true
```

### Step 2: Test Self-Hosting Compiler
```bash
./bin/golisp -f lisp/self-hosting.lisp
# Then in REPL:
# (bootstrap-self-hosting)
```

### Step 3: Bootstrap Process
1. Load self-hosting compiler
2. Compile standard library with itself
3. Compile compiler with itself
4. Verify compiled versions work identically

### Step 4: Iterative Improvement
- Add missing functions as needed
- Improve compiler optimizations
- Expand standard library
- Add development tools

## Key Benefits of Self-Hosting 🌟

1. **Language Evolution**: Easy to add new features
2. **Bootstrapping**: Compiler improvements apply to itself
3. **Dogfooding**: Using the language to build itself
4. **Educational**: Demonstrates language capabilities
5. **Portability**: Easier to port to new platforms
6. **Optimization**: Self-optimizing compiler

## Next Immediate Steps 🎯

1. **✅ Phase 1 Complete - All meta-programming functions implemented**
2. **Add Phase 2 standard library functions (map, filter, apply, etc.)**
3. **Improve multi-expression parsing in self-hosting compiler**
4. **Test self-hosting compiler with realistic examples**
5. **Add macro expansion during compilation**

## Architecture Refactoring Plan 🏗️

### Current State Analysis
The kernel currently has **8,822 lines** across core files with significant bloat:
- `bootstrap.go`: 1,062 lines (mostly built-in functions)
- `eval.go`: 901 lines (evaluation logic)
- `types.go`: 467 lines (data types)
- `parser.go`: 574 lines (lexer/parser)

### Phase 0: Minimal Core Refactoring 🎯

**Goal**: Reduce kernel from 8,822 lines to ~300 lines for a truly minimal, self-hosting core.

#### 1. Ultra-Minimal Kernel Structure
Move to a 3-file core architecture:

**`pkg/core/core.go`** (~100 lines) - Essential types:
```go
type Value interface { String() string }
type Symbol, Keyword, List, Vector, Number, String, Nil
type Environment (basic get/set)
```

**`pkg/core/reader.go`** (~100 lines) - Minimal parser:
```go
// Basic tokenization and s-expression parsing
// Just enough to read: symbols, lists, numbers, strings
```

**`pkg/core/eval.go`** (~100 lines) - Minimal evaluator:
```go
// Only handle: symbol lookup, function calls, special forms
// Special forms: if, def, fn, quote, do
```

#### 2. Self-Hosting Layer (Lisp Implementation)
Move from Go to Lisp:

**From `bootstrap.go` (1,062 lines → ~50 lines Go):**
- All built-in functions except core primitives → `lisp/stdlib/`
- Macros like `defn`, `when`, `unless` → `lisp/macros/`
- Collection operations → `lisp/collections/`
- String operations → `lisp/strings/`
- Math operations beyond basic arithmetic → `lisp/math/`

**Core primitives to keep in Go (~20 functions):**
- Arithmetic: `+`, `-`, `*`, `/`, `=`, `<`, `>`
- Lists: `cons`, `first`, `rest`
- Meta: `eval`, `read-string`
- I/O: `slurp`, `spit`
- Types: `symbol?`, `list?`, `vector?`

#### 3. Modular Extension System

```
pkg/
├── core/           # Minimal kernel (~300 lines)
│   ├── core.go
│   ├── reader.go
│   └── eval.go
├── extensions/     # Optional Go extensions
│   ├── io/
│   ├── math/
│   └── system/
lisp/
├── kernel/         # Self-hosted standard library
│   ├── stdlib.lisp
│   ├── macros.lisp
│   └── collections.lisp
└── compiler/       # Self-hosted compiler
    ├── compiler.lisp
    ├── optimizer.lisp
    └── codegen.lisp
```

#### 4. Refactoring Implementation Plan

**Phase 0.1: Extract Minimal Core**
1. Audit `bootstrap.go` functions: categorize as "core primitive" vs "standard library"
2. Create `pkg/core/` with only essential primitives (~20 functions)
3. Move complex functions to `lisp/stdlib/`

**Phase 0.2: Self-Host Standard Library**
1. Rewrite built-in functions in Lisp using core primitives
2. Implement in `lisp/stdlib/` modules
3. Bootstrap process: `minimal-core → loads stdlib.lisp → full functionality`

**Phase 0.3: Self-Hosting Compiler**
1. Complete existing `lisp/self-hosting.lisp`
2. Add optimization passes in Lisp
3. Bootstrap: `minimal-core → stdlib → compiler → self-hosting`

#### 5. Benefits of Minimal Core

- **Minimal attack surface**: Core reduced from 8,822 to ~300 lines
- **Language evolution**: New features added in Lisp, not Go
- **Self-improvement**: Compiler can optimize itself
- **Portability**: Easy to port minimal core to new platforms
- **Educational**: Demonstrates language capabilities
- **Bootstrapping**: True self-hosting from minimal foundation

#### 6. Migration Strategy

1. **Backward Compatibility**: Keep current API working during transition
2. **Gradual Migration**: Move functions one module at a time
3. **Testing**: Ensure self-hosted versions match Go implementations
4. **Performance**: Profile and optimize self-hosted implementations

### Current Architecture Strengths

Your current architecture is excellent for self-hosting:
- Clean separation between kernel (Go) and library (Lisp)
- Robust macro system for code transformation
- File loading system for modular development
- REPL for interactive development
- Strong error handling and reporting

**The foundation is solid - you're 80% there! The main refactoring needed is extracting the minimal core and moving built-ins to Lisp.**
