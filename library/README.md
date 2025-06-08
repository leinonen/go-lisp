# String Library

This directory contains higher-level string manipulation utilities built on top of the core string functions that are implemented as built-in primitives.

## Architecture

### Built-in String Functions (Go Primitives)
These functions are implemented directly in Go for performance and are available immediately:

- **Basic Operations**: `string-concat`, `string-length`, `string-substring`, `string-char-at`
- **Case Conversion**: `string-upper`, `string-lower`
- **Whitespace**: `string-trim`
- **Splitting/Joining**: `string-split`, `string-join`
- **Search Operations**: `string-contains?`, `string-starts-with?`, `string-ends-with?`
- **Manipulation**: `string-replace`, `string-index-of`, `string-repeat`
- **Type Conversion**: `string->number`, `number->string`
- **Regular Expressions**: `string-regex-match?`, `string-regex-find-all`
- **Validation**: `string?`, `string-empty?`

### Library Functions (Lisp Compositions)
The functions in `strings.lisp` are higher-level utilities that combine the primitives:

- **Convenience Aliases**: `str-concat`, `str-empty`
- **Enhanced Validation**: `str-blank?`, `str-non-empty?`, `str-numeric?`, `str-alpha?`, `str-alnum?`
- **Text Processing**: `str-words`, `str-lines`, `str-char-count`
- **Transformation**: `str-reverse`, `str-capitalize`, `str-title-case`
- **Formatting**: `str-pad-left`, `str-pad-right`, `str-center`

## Usage

```lisp
; Built-in functions are always available
(string-concat "Hello" " " "World")  ; => "Hello World"
(string-length "Hello")              ; => 5
(string-upper "hello")               ; => "HELLO"

; Load the library for higher-level functions
(load "library/strings.lisp")
(import strings)

; Use composed functions
(str-capitalize "hello world")       ; => "Hello world"
(str-title-case "hello world")       ; => "Hello World"
(str-reverse "hello")                ; => "olleh"
(str-numeric? "123")                 ; => #t
```

## Benefits of This Architecture

1. **Performance**: Critical operations are implemented in Go for speed
2. **Extensibility**: Higher-level operations can be easily added in Lisp
3. **Readability**: Library functions provide more descriptive names and behaviors
4. **Modularity**: Users can choose to load only the functions they need
5. **Educational**: Shows how to build complex operations from simple primitives

The built-in functions provide the essential building blocks, while the library functions demonstrate idiomatic patterns and provide convenient abstractions for common tasks.
