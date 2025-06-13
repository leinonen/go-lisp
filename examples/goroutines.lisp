;; Goroutine and Channel Examples for Lisp Interpreter
;; This file demonstrates the new concurrency features

;; Basic Goroutine Example
;; Create a simple goroutine that calculates a square
(def square-future (go (* 5 5)))
(println! "Goroutine started...")
(println! "Result:" (go-wait square-future))

;; Multiple Goroutines
;; Start several goroutines and wait for all results
(def futures (list
  (go (+ 1 2 3))
  (go (* 4 5))
  (go (- 10 3))))

(def results (go-wait-all futures))
(println! "Multiple goroutine results:" results)

;; Goroutines with Shared State (Atoms)
;; Multiple goroutines incrementing a shared counter
(def counter (atom 0))

(def increment-futures (list
  (go (reset! counter (+ (deref counter) 1)))
  (go (reset! counter (+ (deref counter) 1)))
  (go (reset! counter (+ (deref counter) 1)))
  (go (reset! counter (+ (deref counter) 1)))
  (go (reset! counter (+ (deref counter) 1)))))

;; Wait for all increments to complete
(go-wait-all increment-futures)
(println! "Final counter value:" (deref counter))

;; Basic Channel Example
;; Create a channel and send/receive values
(def my-chan (chan 2)) ; buffered channel with size 2

;; Send some values
(chan-send! my-chan "Hello")
(chan-send! my-chan "World")

;; Receive values
(println! "Received:" (chan-recv! my-chan))
(println! "Received:" (chan-recv! my-chan))

;; Goroutines with Channels
;; Producer-Consumer pattern
(def data-chan (chan))
(def result-chan (chan))

;; Producer goroutine
(def producer (go 
  (chan-send! data-chan 10)))

;; Send more values
(def producer2 (go (chan-send! data-chan 20)))
(def producer3 (go (chan-send! data-chan 30)))

;; Wait for all sends to complete
(go-wait producer)
(go-wait producer2) 
(go-wait producer3)
(chan-close! data-chan)

;; Consumer goroutine that sums values
(def consumer (go
  (+ (+ (chan-recv! data-chan) (chan-recv! data-chan)) (chan-recv! data-chan))))

;; Get the result
(def sum-result (go-wait consumer))
(println! "Sum from producer-consumer:" sum-result)

;; Wait for goroutines to complete
(go-wait consumer)

;; Channel Operations
;; Demonstrate try-receive and channel closing
(def test-chan (chan 1))

;; Try to receive from empty channel
(def empty-result (chan-try-recv! test-chan))
(if (= empty-result nil)
  (println! "Channel is empty (as expected)")
  (println! "Unexpected value:" empty-result))

;; Send a value and try again
(chan-send! test-chan "test-value")
(def received-value (chan-try-recv! test-chan))
(println! "Try-received:" received-value)

;; Check if channel is closed
(println! "Channel closed?" (chan-closed? test-chan))
(chan-close! test-chan)
(println! "Channel closed?" (chan-closed? test-chan))

(println! "Goroutine and channel examples completed!")
