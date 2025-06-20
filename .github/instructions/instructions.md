---
applyTo: '**'
---
- New features should be implemented with TDD
- Plan new features and analyze existing codebase  before implementing new features
- Work systematically, to prevent features getting too large.
- New lisp functions should use polymorphism, so that we dont create duplicate functions for different data types.
- New language features should take inspiration from Clojure, because I want the final implementation to be very similar to Clojure, except for the Java-related parts of Clojure.
- Only create examples (lisp) when tests pass. Also verify each example carefully.
- Use `make` to build the binary, and the binary is located here: `./bin/golisp`
