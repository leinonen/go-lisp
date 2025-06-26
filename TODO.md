# Go Lisp - Clojure Feature Implementation TODO

This document outlines the functions to implement to bring Go Lisp closer to Clojure's feature set (excluding Java interop).

## High Priority - Core Sequence Operations

### Essential Collection Operations
- [ ] **`cons`** - Already implemented, ensure compatibility with all collection types
- [ ] **`conj`** - Add elements to collections (vector: append, list: prepend)
- [ ] **`concat`** - Concatenate multiple sequences
- [ ] **`count`** - Get length of collections (more idiomatic than `length`)
- [ ] **`empty?`** - Check if collection is empty
- [ ] **`seq`** - Convert collections to sequences
- [ ] **`next`** - Like `rest` but returns nil for empty sequences
- [ ] **`nth`** - Get element at index with optional default
- [ ] **`take`** - Take first n elements
- [ ] **`drop`** - Drop first n elements
- [ ] **`reverse`** - Reverse a sequence

### Collection Transformations
- [ ] **`map`** - Transform collections (improve existing stdlib version)
- [ ] **`filter`** - Filter elements by predicate
- [ ] **`reduce`** - Already implemented, ensure robustness
- [ ] **`apply`** - Apply function to collection as arguments

## High Priority - Map/Dictionary Operations

- [ ] **`zipmap`** - Create map from keys and values sequences
- [ ] **`keys`** - Get map keys (improve existing)
- [ ] **`vals`** - Get map values
- [ ] **`assoc`** - Associate key-value pairs in maps
- [ ] **`dissoc`** - Remove keys from maps
- [ ] **`get`** - Get value from map/vector with default
- [ ] **`contains?`** - Check if collection contains key/element

## High Priority - Type Predicates

### Basic Type Checks
- [ ] **`nil?`** - Check if value is nil
- [ ] **`string?`** - Check if value is string
- [ ] **`number?`** - Check if value is number
- [ ] **`vector?`** - Check if value is vector
- [ ] **`list?`** - Check if value is list
- [ ] **`map?`** - Check if value is map
- [ ] **`set?`** - Check if value is set
- [ ] **`fn?`** - Check if value is function

### Numeric Predicates
- [ ] **`zero?`** - Check if number is zero
- [ ] **`pos?`** - Check if number is positive
- [ ] **`neg?`** - Check if number is negative
- [ ] **`even?`** - Check if number is even
- [ ] **`odd?`** - Check if number is odd

## Medium Priority - Logic & Flow Control

### Logical Operations
- [ ] **`and`** - Logical and (short-circuiting)
- [ ] **`or`** - Logical or (short-circuiting)
- [ ] **`not`** - Logical not

### Control Flow
- [ ] **`when-not`** - Conditional execution when false
- [ ] **`cond`** - Multi-way conditional
- [ ] **`case`** - Pattern matching
- [ ] **`let`** - Local bindings

## Medium Priority - String Operations

- [ ] **`str`** - String concatenation
- [ ] **`subs`** - Substring
- [ ] **`split`** - Split strings
- [ ] **`join`** - Join strings
- [ ] **`trim`** - Trim whitespace

## Medium Priority - Math Operations

### Basic Math
- [ ] **`inc`** - Increment number
- [ ] **`dec`** - Decrement number
- [ ] **`max`** - Maximum of numbers
- [ ] **`min`** - Minimum of numbers
- [ ] **`mod`** - Modulo operation
- [ ] **`quot`** - Integer division
- [ ] **`rem`** - Remainder

## Medium Priority - Advanced Sequence Operations

### Sequence Utilities
- [ ] **`partition`** - Partition sequence into chunks
- [ ] **`take-while`** - Take while predicate is true
- [ ] **`drop-while`** - Drop while predicate is true
- [ ] **`some`** - Check if any element satisfies predicate
- [ ] **`every?`** - Check if all elements satisfy predicate

### Infinite Sequences
- [ ] **`repeatedly`** - Generate infinite sequence
- [ ] **`repeat`** - Repeat value n times
- [ ] **`cycle`** - Cycle through sequence infinitely

## Lower Priority - Set Operations

- [ ] **`union`** - Set union
- [ ] **`intersection`** - Set intersection
- [ ] **`difference`** - Set difference
- [ ] **`subset?`** - Check if set is subset
- [ ] **`superset?`** - Check if set is superset

## Lower Priority - I/O Operations

- [ ] **`slurp`** - Read file contents
- [ ] **`spit`** - Write to file
- [ ] **`println`** - Print with newline
- [ ] **`prn`** - Print for reading back
- [ ] **`read-string`** - Parse string as Lisp data

## Lower Priority - Meta Programming

- [ ] **`eval`** - Evaluate data as code
- [ ] **`macroexpand`** - Expand macro once
- [ ] **`macroexpand-1`** - Expand macro completely

## Lower Priority - Enhanced Features

### Destructuring Support
- [ ] **Enhanced `let`** with destructuring
- [ ] **Enhanced function parameters** with destructuring

### Additional Conveniences
- [ ] **`partial`** - Partial function application
- [ ] **`comp`** - Function composition
- [ ] **`constantly`** - Return constant function
- [ ] **`identity`** - Identity function

## Implementation Status

### ✅ Already Implemented
- Basic arithmetic (`+`, `-`, `*`)
- Comparison operators (`=`, `<`, `<=`, `>`, `>=`, `!=`)
- Core special forms (`if`, `fn`, `def`, `do`, `quote`, `loop`, `recur`)
- Basic collections (lists, vectors, hashmaps, sets)
- Sequence operations (`first`, `rest`, `list`, `vector`)
- HashMap operations (`hash-map`, `hash-map-get`, `hash-map-put`)
- Basic I/O (`print`, `load`)
- Macros (`defmacro`, `defn`)

### 🔄 Partially Implemented
- `map` (basic version in stdlib)
- `reduce` (implemented in bootstrap)
- Collection access functions (basic versions exist)
