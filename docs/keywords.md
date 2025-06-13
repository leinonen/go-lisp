# Keywords

Keywords are symbolic constants that start with a colon (`:`) and evaluate to themselves. They are commonly used as keys in hash maps due to their self-evaluating nature and improved readability compared to strings.

## Overview

Keywords in this Lisp interpreter:
- Start with a colon `:` followed by valid symbol characters
- Are self-evaluating (they evaluate to themselves)
- Can contain letters, numbers, and common symbols like `-`, `_`, `?`, `!`
- Are distinct from strings (`:name` is different from `"name"`)
- Display with the colon prefix: `:keyword-name`

```lisp
;; Clear distinction between function parameters and function calls
(def create-user (fn [name email] ...))  ; Function definition
(create-user "Alice" "alice@example.com")       ; Function call
```

## Syntax

```lisp
:keyword-name
:first-name
:user-id
:item-123
:valid?
:has-value!
```

**Valid characters in keywords:**
- Letters (a-z, A-Z)
- Numbers (0-9)
- Hyphens (-)
- Underscores (_)
- Question marks (?)
- Exclamation marks (!)
- Other symbol characters: `+`, `*`, `/`, `=`, `<`, `>`, `.`, `%`

## Basic Usage

### Self-Evaluation
Keywords evaluate to themselves:

```lisp
lisp> :name
=> :name

lisp> :user-status
=> :user-status
```

### In Variables
Keywords can be stored in variables:

```lisp
lisp> (def status-key :status)
=> :status

lisp> status-key
=> :status
```

### In Data Structures
Keywords work well in lists and other data structures:

```lisp
lisp> (list :name :age :email)
=> (:name :age :email)
```

## Hash Maps with Keywords

Keywords are particularly useful as hash map keys because they're more readable and idiomatic than strings:

### Creating Hash Maps
```lisp
lisp> (hash-map :name "Alice" :age 30 :city "Boston")
=> {:name Alice, :age 30, :city Boston}
```

### Accessing Values
```lisp
lisp> (def person (hash-map :name "Alice" :age 30))
=> {:name Alice, :age 30}

lisp> (hash-map-get person :name)
=> Alice

lisp> (hash-map-get person :age)
=> 30

lisp> (hash-map-get person :missing)
=> nil
```

### Modifying Hash Maps
```lisp
lisp> (hash-map-put person :email "alice@example.com")
=> {:name Alice, :age 30, :email alice@example.com}

lisp> (hash-map-remove person :age)
=> {:name Alice}
```

### Querying Hash Maps
```lisp
lisp> (hash-map-contains? person :name)
=> true

lisp> (hash-map-contains? person :missing)
=> false
```

## Keywords vs Strings

Keywords and strings are different types:

```lisp
lisp> (hash-map "name" "string-key" :name "keyword-key")
=> {name string-key, :name keyword-key}

lisp> (hash-map-get mixed-map "name")    ; String key
=> string-key

lisp> (hash-map-get mixed-map :name)     ; Keyword key
=> keyword-key
```

## Common Patterns

### Configuration Objects
Keywords make configuration more readable:

```lisp
lisp> (def config (hash-map
        :debug true
        :port 8080
        :host "localhost"
        :max-connections 100))
=> {:debug true, :port 8080, :host localhost, :max-connections 100}

lisp> (hash-map-get config :debug)
=> true
```

### Record Creation
Use keywords for creating structured data:

```lisp
lisp> (def create-user (fn [name email]
        (hash-map
          :name name
          :email email
          :created-at "2024-01-15"
          :active true)))

lisp> (create-user "Alice" "alice@example.com")
=> {:name Alice, :email alice@example.com, :created-at 2024-01-15, :active true}
```

### Nested Structures
Keywords work well for nested hash maps:

```lisp
lisp> (def app-state (hash-map
        :user (hash-map :name "Alice" :id 123)
        :settings (hash-map :theme :dark :lang :en)))
=> {:user {:name Alice, :id 123}, :settings {:theme :dark, :lang :en}}

lisp> (hash-map-get (hash-map-get app-state :user) :name)
=> Alice
```

## Error Cases

```lisp
lisp> :
error: invalid keyword: colon must be followed by symbol characters

lisp> : 
error: invalid keyword: colon must be followed by symbol characters
```

## Performance

- **Keyword evaluation**: O(1) - keywords are self-evaluating constants
- **Hash map operations**: Same performance as string keys
- **Memory**: Keywords are stored as strings internally with a `:` prefix

## Integration with Other Features

Keywords work seamlessly with:
- **Variables**: Can be stored and passed around like any value
- **Functions**: Can be function parameters and return values
- **Lists**: Can be elements in lists for processing
- **Conditionals**: Can be compared and used in conditional expressions
- **Modules**: Can be exported and imported between modules

## Best Practices

1. **Use keywords for hash map keys** instead of strings when the key represents a symbolic constant
2. **Use consistent naming** with kebab-case (`:first-name` not `:firstName`)
3. **Use descriptive names** that clearly indicate the purpose (`:user-id` not `:id`)
4. **Prefer keywords over strings** for configuration keys and record fields
5. **Group related keywords** with common prefixes (`:user-name`, `:user-email`, `:user-status`)

This implementation provides a clean and idiomatic way to work with symbolic constants and structured data in the Lisp interpreter.
