; Module System Demo
; This file demonstrates the module system with imports, exports, and qualified access

; Create a math utilities module
(module math-utils
  (export square cube add-squares power)
  
  (defn square [x] [* x x])
  (defn cube [x] [* x x x])
  (defn add-squares [x y] (+ (square x) [square y]))
  (defn power [base exp]
    (if (= exp 0) 1 (* base (power base (- exp 1)))))
    
  ; Private helper function (not exported)
  (defn helper [x] [+ x 1]))

; Create a list utilities module  
(module list-utils
  (export double-all sum-list reverse-and-double)
  
  (defn double-all [lst] (map (fn [x] [* x 2]) lst))
  (defn sum-list [lst] (reduce (fn [acc x] [+ acc x]) 0 lst))
  (defn reverse-and-double [lst] (double-all [reverse lst])))

; Demonstrate qualified access (without importing)
(math-utils.square 5)                   ; => 25
(math-utils.cube 3)                     ; => 27
(list-utils.double-all (list 1 2 3))    ; => (2 4 6)

; Import a module to use functions directly
(import math-utils)
(square 7)                              ; => 49
(add-squares 3 4)                       ; => 25

; Import another module
(import list-utils)
(def test-list (list 1 2 3 4 5))
(sum-list test-list)                    ; => 15
(reverse-and-double test-list)          ; => (10 8 6 4 2)

; Modules can still be accessed with qualified names after import
(math-utils.power 2 8)                  ; => 256

; Check what modules are loaded
(modules)

; Private functions are not accessible
; (helper 5)                            ; Would cause error
; (math-utils.helper 5)                 ; Would also cause error
