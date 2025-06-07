# Supported Operations

This document provides a comprehensive reference for all operations supported by the Lisp interpreter.

## Arithmetic

- `(+ 1 2 3)` - Addition with multiple operands
- `(- 10 3)` - Subtraction
- `(* 2 3 4)` - Multiplication with multiple operands
- `(/ 15 3)` - Division

## Comparison

- `(= 5 5)` - Equality
- `(< 3 5)` - Less than
- `(> 7 3)` - Greater than

## Conditional

- `(if condition then-expr else-expr)` - If expression

## Variables

- `(define name value)` - Define a variable with a name and value

## Functions

- `(lambda (params) body)` - Create an anonymous function
- `(defun name (params) body)` - Define a named function (combines define and lambda)
- `(funcname args...)` - Call a user-defined function

## Lists

- `(list)` - Create an empty list
- `(list 1 2 3)` - Create a list with elements
- `(first lst)` - Get the first element of a list
- `(rest lst)` - Get all elements except the first
- `(cons elem lst)` - Prepend an element to a list
- `(length lst)` - Get the number of elements in a list
- `(empty? lst)` - Check if a list is empty
- `(append lst1 lst2)` - Combine two lists into one
- `(reverse lst)` - Reverse the order of elements in a list
- `(nth index lst)` - Get the element at a specific index (0-based)

## Higher-Order Functions

- `(map func lst)` - Apply a function to each element of a list
- `(filter predicate lst)` - Keep only elements that satisfy a predicate
- `(reduce func init lst)` - Reduce a list to a single value using a function

## Comments

- `;` - Line comments (from semicolon to end of line are ignored)
- Comments can appear anywhere in the code
- Useful for documenting code and adding explanations

## Module System

- `(module name (export sym1 sym2...) body...)` - Define a module with exported symbols
- `(import module-name)` - Import all exported symbols from a module into current scope
- `(load "filename.lisp")` - Load and execute a Lisp file
- `module.symbol` - Qualified access to module symbols without importing

## Environment Inspection

- `(env)` - Show all variables and functions in the current environment
- `(modules)` - Show all loaded modules and their exported symbols
- `(builtins)` - Show all available built-in functions and special forms
- `(builtins func-name)` - Get detailed help for a specific built-in function
