# New List Functions Documentation

This document describes the newly added list functions to the Lisp interpreter.

## Functions

### `last`
**Syntax:** `(last list)`

Returns the last element of a list.

**Examples:**
```lisp
(last (list 1 2 3 4))      ; => 4
(last (list "a"))          ; => "a"
```

**Errors:**
- Throws an error if the list is empty
- Throws an error if the argument is not a list

---

### `butlast`
**Syntax:** `(butlast list)`

Returns all elements of a list except the last one.

**Examples:**
```lisp
(butlast (list 1 2 3 4))   ; => (1 2 3)
(butlast (list "a"))       ; => ()
(butlast (list))           ; => ()
```

**Errors:**
- Throws an error if the argument is not a list

---

### `flatten`
**Syntax:** `(flatten list)`

Recursively flattens nested lists into a single flat list.

**Examples:**
```lisp
(flatten (list 1 (list 2 3) 4))              ; => (1 2 3 4)
(flatten (list 1 (list 2 (list 3 4)) 5))     ; => (1 2 3 4 5)
(flatten (list 1 2 3))                       ; => (1 2 3)
```

**Errors:**
- Throws an error if the argument is not a list

---

### `zip`
**Syntax:** `(zip list1 list2 ...)`

Combines multiple lists by pairing corresponding elements. The result length is determined by the shortest input list.

**Examples:**
```lisp
(zip (list 1 2 3) (list "a" "b" "c"))        ; => ((1 "a") (2 "b") (3 "c"))
(zip (list 1 2) (list "a" "b" "c"))          ; => ((1 "a") (2 "b"))
(zip (list 1 2 3) (list 10 20 30) (list "x" "y" "z"))  ; => ((1 10 "x") (2 20 "y") (3 30 "z"))
```

**Errors:**
- Requires at least 2 arguments
- All arguments must be lists

---

### `sort`
**Syntax:** `(sort list [comparator])`

Sorts a list in ascending order. Optionally accepts a custom comparator function.

**Examples:**
```lisp
(sort (list 5 2 8 1 9 3))                    ; => (1 2 3 5 8 9)
(sort (list "banana" "apple" "cherry"))      ; => ("apple" "banana" "cherry")
(sort (list 5 2 8 1) (fn [a b] (> a b)))     ; => (8 5 2 1) (descending)
```

**Default comparison:**
- Numbers: numerical order
- Strings: lexicographical order
- Mixed types: converted to strings for comparison

**Errors:**
- Throws an error if the first argument is not a list
- If provided, the comparator must be a function that takes exactly 2 parameters

---

### `distinct`
**Syntax:** `(distinct list)`

Returns a new list with duplicate elements removed, preserving the order of first occurrence.

**Examples:**
```lisp
(distinct (list 1 2 2 3 1 4 3))              ; => (1 2 3 4)
(distinct (list "a" "b" "a" "c"))            ; => ("a" "b" "c")
(distinct (list 1 2 3))                      ; => (1 2 3)
```

**Errors:**
- Throws an error if the argument is not a list

---

### `concat`
**Syntax:** `(concat list1 list2 ...)`

Concatenates multiple lists into a single list.

**Examples:**
```lisp
(concat (list 1 2) (list 3 4))               ; => (1 2 3 4)
(concat (list 1 2) (list 3 4) (list 5 6))    ; => (1 2 3 4 5 6)
(concat (list) (list 1 2) (list))            ; => (1 2)
(concat)                                      ; => ()
```

**Errors:**
- All arguments must be lists

---

### `partition`
**Syntax:** `(partition n list)`

Splits a list into chunks of size `n`. The last chunk may be smaller if the list length is not evenly divisible by `n`.

**Examples:**
```lisp
(partition 2 (list 1 2 3 4 5))               ; => ((1 2) (3 4) (5))
(partition 3 (list 1 2 3 4 5 6 7 8 9))       ; => ((1 2 3) (4 5 6) (7 8 9))
(partition 1 (list "a" "b" "c"))             ; => (("a") ("b") ("c"))
```

**Errors:**
- The first argument must be a positive number
- The second argument must be a list

---

## Usage in Combination

These functions can be combined for powerful list processing:

```lisp
; Get the last chunk of a partitioned list
(last (partition 3 (list 1 2 3 4 5 6 7 8 9)))  ; => (7 8 9)

; Sort and remove duplicates
(sort (distinct (list 3 1 4 1 5 9 2 6 5)))      ; => (1 2 3 4 5 6 9)

; Flatten then zip with indices
(zip (list 0 1 2 3) (flatten (list (list 1 2) (list 3 4))))  ; => ((0 1) (1 2) (2 3) (3 4))
```

## Integration with Existing Functions

These new functions work seamlessly with existing list functions like `map`, `filter`, and `reduce`:

```lisp
; Sort a filtered list
(sort (filter (fn [x] (> x 5)) (list 8 3 12 1 9 4 15)))  ; => (8 9 12 15)

; Map over flattened nested structure
(map (fn [x] (* x 2)) (flatten (list (list 1 2) (list 3 4))))  ; => (2 4 6 8)
```
