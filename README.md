# GoLisp

![GoLisp logo](./docs/img/golisp-logo.png)

A minimalist Lisp interpreter written in Go, inspired by Clojure.

## Usage

```bash
make build
./bin/golisp
```

## Examples

```lisp
(+ 1 2 3)                         ; 6
(def square (fn [x] (* x x)))     ; define function
(square 5)                        ; 25
[1 2 3]                           ; vectors
```

## License

MIT
