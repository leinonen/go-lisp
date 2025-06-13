; Test file for new list functions
; This demonstrates the new list functions added to the Lisp interpreter

(println! "=== Testing New List Functions ===")

; Test data
(def numbers (list 1 2 3 4 5))
(def nested (list 1 (list 2 3) 4 (list 5 (list 6 7))))
(def duplicates (list 1 2 2 3 1 4 3))

; Test last function
(println! "\n--- Testing last ---")
(println! "last of" numbers "=>" (last numbers))
(println! "last of (list 42) =>" (last (list 42)))

; Test butlast function
(println! "\n--- Testing butlast ---")
(println! "butlast of" numbers "=>" (butlast numbers))
(println! "butlast of (list 42) =>" (butlast (list 42)))
(println! "butlast of empty list =>" (butlast (list)))

; Test flatten function
(println! "\n--- Testing flatten ---")
(println! "flatten of" nested "=>" (flatten nested))
(println! "flatten of simple list" numbers "=>" (flatten numbers))

; Test zip function
(println! "\n--- Testing zip ---")
(def letters (list "a" "b" "c"))
(println! "zip" numbers "and" letters "=>" (zip numbers letters))
(println! "zip with different lengths =>" (zip (list 1 2) (list "x" "y" "z")))

; Test sort function
(println! "\n--- Testing sort ---")
(def unsorted (list 5 2 8 1 9 3))
(println! "sort" unsorted "=>" (sort unsorted))
(println! "sort strings =>" (sort (list "banana" "apple" "cherry")))

; Test distinct function
(println! "\n--- Testing distinct ---")
(println! "distinct" duplicates "=>" (distinct duplicates))
(println! "distinct of already unique list =>" (distinct numbers))

; Test concat function
(println! "\n--- Testing concat ---")
(def list1 (list 1 2))
(def list2 (list 3 4))
(def list3 (list 5 6))
(println! "concat" list1 list2 list3 "=>" (concat list1 list2 list3))
(println! "concat with empty list =>" (concat (list) numbers (list)))

; Test partition function
(println! "\n--- Testing partition ---")
(def long-list (list 1 2 3 4 5 6 7 8 9))
(println! "partition" long-list "by 3 =>" (partition 3 long-list))
(println! "partition" numbers "by 2 =>" (partition 2 numbers))

; Demonstrate usage in combination
(println! "\n--- Combination Examples ---")
(println! "Last 3 elements:" (last (partition 3 numbers)))
(println! "Distinct sorted:" (sort (distinct duplicates)))
(println! "Flatten then take every 2nd:" (butlast (flatten nested)))

(println! "\n=== All tests completed! ===")
