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
(defn square [x] (* x x))         ; define function
(square 5)                        ; 25
[1 2 3]                           ; vectors
{:name "GoLisp" :lang "Clojure"}  ; maps
(loop [n 5 acc 1] (if (= n 0) acc (recur (- n 1) (* acc n))))  ; factorial
```

## License

MIT
