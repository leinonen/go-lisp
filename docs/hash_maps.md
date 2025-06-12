# Hash Maps

Hash maps are key-value associative data structures that provide efficient storage and retrieval of data by keys. This implementation uses immutable operations, meaning that all modification operations return new hash maps rather than modifying existing ones.

## Overview

Hash maps in this Lisp interpreter:
- Use strings as keys and any Lisp value as values
- Support immutable operations (functional style)
- Provide O(1) average case lookup, insertion, and deletion
- Display in a readable format: `{key1: value1, key2: value2, ...}`
- Return `nil` for missing keys

```lisp
;; Clear function parameter syntax
(def get-with-default (lambda [map key default]
  (if (hash-map-contains? map key)
      (hash-map-get map key)
      default)))

;; Nested functions are easier to distinguish
(def transform-values (lambda [map transformer]
  (reduce (lambda [acc key]
    (hash-map-put acc key (transformer (hash-map-get map key))))
    (hash-map) (hash-map-keys map))))
```

## Creating Hash Maps

### `hash-map`
Creates a new hash map from alternating key-value pairs.

**Syntax:**
```lisp
(hash-map key1 value1 key2 value2 ...)
```

**Examples:**
```lisp
lisp> (hash-map)
=> {}

lisp> (hash-map "name" "Alice")
=> {name: Alice}

lisp> (hash-map "name" "Alice" "age" 30 "city" "Boston")
=> {name: Alice, age: 30, city: Boston}

lisp> (hash-map "numbers" (list 1 2 3) "nested" (hash-map "inner" "value"))
=> {numbers: (1 2 3), nested: {inner: value}}
```

**Error Cases:**
```lisp
lisp> (hash-map "key")
error: hash-map requires an even number of arguments

lisp> (hash-map 123 "value")
error: hash-map keys must be strings, got number
```

## Accessing Values

### `hash-map-get`
Retrieves a value from a hash map by key.

**Syntax:**
```lisp
(hash-map-get hash-map key)
```

**Examples:**
```lisp
lisp> (def my-map (hash-map "name" "Alice" "age" 30))
=> {name: Alice, age: 30}

lisp> (hash-map-get my-map "name")
=> Alice

lisp> (hash-map-get my-map "age")
=> 30

lisp> (hash-map-get my-map "missing")
=> nil
```

**Error Cases:**
```lisp
lisp> (hash-map-get "not-a-map" "key")
error: expected hash-map as first argument to hash-map-get

lisp> (hash-map-get my-map 123)
error: expected string key for hash-map-get
```

## Modifying Hash Maps

### `hash-map-put`
Creates a new hash map with an additional or updated key-value pair.

**Syntax:**
```lisp
(hash-map-put hash-map key value)
```

**Examples:**
```lisp
lisp> (def original (hash-map "name" "Alice"))
=> {name: Alice}

lisp> (hash-map-put original "age" 30)
=> {name: Alice, age: 30}

lisp> (hash-map-put original "name" "Bob")  ; Updates existing key
=> {name: Bob}

; Original hash map is unchanged (immutable)
lisp> original
=> {name: Alice}
```

### `hash-map-remove`
Creates a new hash map with a key-value pair removed.

**Syntax:**
```lisp
(hash-map-remove hash-map key)
```

**Examples:**
```lisp
lisp> (def my-map (hash-map "name" "Alice" "age" 30 "city" "Boston"))
=> {name: Alice, age: 30, city: Boston}

lisp> (hash-map-remove my-map "age")
=> {name: Alice, city: Boston}

lisp> (hash-map-remove my-map "missing")  ; No error for missing keys
=> {name: Alice, age: 30, city: Boston}
```

## Querying Hash Maps

### `hash-map-contains?`
Tests whether a hash map contains a specific key.

**Syntax:**
```lisp
(hash-map-contains? hash-map key)
```

**Examples:**
```lisp
lisp> (def my-map (hash-map "name" "Alice" "age" nil))
=> {name: Alice, age: nil}

lisp> (hash-map-contains? my-map "name")
=> #t

lisp> (hash-map-contains? my-map "age")    ; Returns true even for nil values
=> #t

lisp> (hash-map-contains? my-map "missing")
=> #f
```

### `hash-map-keys`
Returns a list of all keys in the hash map.

**Syntax:**
```lisp
(hash-map-keys hash-map)
```

**Examples:**
```lisp
lisp> (hash-map-keys (hash-map))
=> ()

lisp> (hash-map-keys (hash-map "name" "Alice" "age" 30))
=> (name age)  ; Order may vary
```

### `hash-map-values`
Returns a list of all values in the hash map.

**Syntax:**
```lisp
(hash-map-values hash-map)
```

**Examples:**
```lisp
lisp> (hash-map-values (hash-map))
=> ()

lisp> (hash-map-values (hash-map "name" "Alice" "age" 30))
=> (Alice 30)  ; Order corresponds to keys
```

### `hash-map-size`
Returns the number of key-value pairs in the hash map.

**Syntax:**
```lisp
(hash-map-size hash-map)
```

**Examples:**
```lisp
lisp> (hash-map-size (hash-map))
=> 0

lisp> (hash-map-size (hash-map "name" "Alice" "age" 30))
=> 2
```

### `hash-map-empty?`
Tests whether a hash map is empty.

**Syntax:**
```lisp
(hash-map-empty? hash-map)
```

**Examples:**
```lisp
lisp> (hash-map-empty? (hash-map))
=> #t

lisp> (hash-map-empty? (hash-map "key" "value"))
=> #f
```

## Practical Examples

### Building a Person Record
```lisp
lisp> (def person (hash-map "name" "Alice" "age" 30 "email" "alice@example.com"))
=> {name: Alice, age: 30, email: alice@example.com}

lisp> (def updated-person (hash-map-put person "phone" "555-1234"))
=> {name: Alice, age: 30, email: alice@example.com, phone: 555-1234}

lisp> (hash-map-get updated-person "phone")
=> 555-1234
```

### Configuration Management
```lisp
lisp> (def config (hash-map "debug" #t "port" 8080 "host" "localhost"))
=> {debug: #t, port: 8080, host: localhost}

lisp> (if (hash-map-get config "debug")
        (hash-map-put config "log-level" "verbose")
        config)
=> {debug: #t, port: 8080, host: localhost, log-level: verbose}
```

### Data Processing
```lisp
lisp> (def inventory (hash-map "apples" 10 "oranges" 5 "bananas" 8))
=> {apples: 10, oranges: 5, bananas: 8}

lisp> (def total-items 
        (reduce + 0 (hash-map-values inventory)))
=> 23

lisp> (hash-map-keys inventory)
=> (apples oranges bananas)
```

### Nested Hash Maps
```lisp
lisp> (def nested (hash-map 
        "user" (hash-map "name" "Alice" "id" 123)
        "settings" (hash-map "theme" "dark" "notifications" #t)))
=> {user: {name: Alice, id: 123}, settings: {theme: dark, notifications: #t}}

lisp> (hash-map-get (hash-map-get nested "user") "name")
=> Alice
```

## Common Patterns

### Safe Access with Default Values
```lisp
lisp> (def get-with-default (lambda [map key default]
        (if (hash-map-contains? map key)
            (hash-map-get map key)
            default)))

lisp> (get-with-default my-map "missing" "not found")
=> not found
```

### Bulk Updates
```lisp
lisp> (def bulk-put (lambda [map pairs]
        (if (empty? pairs)
            map
            (bulk-put 
              (hash-map-put map (first pairs) (first (rest pairs)))
              (rest (rest pairs))))))

lisp> (bulk-put (hash-map) (list "a" 1 "b" 2 "c" 3))
=> {a: 1, b: 2, c: 3}
```

### Filtering Hash Maps
```lisp
lisp> (def filter-map (lambda [predicate map]
        (def filtered-pairs 
          (filter (lambda [key] (predicate (hash-map-get map key)))
                  (hash-map-keys map)))
        (bulk-put (hash-map) 
          (reduce append () 
            (map (lambda [key] (list key (hash-map-get map key))) 
                 filtered-pairs)))))

lisp> (filter-map (lambda [val] (> val 5)) 
                  (hash-map "a" 3 "b" 7 "c" 10))
=> {b: 7, c: 10}
```

## Performance Characteristics

- **Creation**: O(n) where n is the number of key-value pairs
- **Access**: O(1) average case for `hash-map-get`
- **Insertion**: O(n) due to immutable copying (creates new hash map)
- **Removal**: O(n) due to immutable copying (creates new hash map)
- **Size queries**: O(1) for `hash-map-size` and `hash-map-empty?`
- **Key/value enumeration**: O(n) for `hash-map-keys` and `hash-map-values`

## Design Principles

1. **Immutability**: Operations return new hash maps, preserving functional programming principles
2. **Type Safety**: Keys must be strings; values can be any Lisp type
3. **Nil Handling**: Missing keys return `nil`; `nil` values are valid and distinct from missing keys
4. **Consistency**: Function naming follows established patterns (`hash-map-*`)
5. **Error Handling**: Clear error messages for type mismatches and invalid arguments

## Integration with Other Features

Hash maps work seamlessly with other Lisp features:

- **Lists**: Keys and values can be extracted as lists for processing with `map`, `filter`, `reduce`
- **Conditionals**: Use `hash-map-contains?` and `hash-map-empty?` in conditional expressions
- **Functions**: Store functions as values in hash maps
- **Modules**: Hash maps can be used within modules and exported/imported
- **Recursion**: Safe to use in recursive functions with tail call optimization

This implementation provides a solid foundation for associative data structures in the Lisp interpreter, enabling more sophisticated data manipulation and storage patterns.
