;; File Functions Demo
;; This file demonstrates the file I/O capabilities of the Lisp interpreter

(println! "=== File Functions Demo ===\n")

;; 1. Basic File Writing
(println! "1. Writing to files")
(println! "-------------------")

; Write simple text
(def message "Hello, World!\nThis is a test file.\nWritten from Lisp!")
(write-file "test-output.txt" message)
(println! "✓ Written to test-output.txt")

; Write computed content
(def numbers (list 1 2 3 4 5))
(def number-sum (reduce (fn [acc x] (+ acc x)) 0 numbers))
(def result-text (string-concat "Sum of " (string-join (map number->string numbers) ", ") " = " (number->string number-sum)))
(write-file "calculation-result.txt" result-text)
(println! "✓ Written calculation result to calculation-result.txt")

;; 2. File Existence Checking
(println! "\n2. Checking file existence")
(println! "---------------------------")

(def files-to-check (list "test-output.txt" "calculation-result.txt" "nonexistent.txt"))
(map (fn [file] 
       (if (file-exists? file)
           (println! "✓" file "exists")
           (println! "✗" file "does not exist")))
     files-to-check)

;; 3. Reading Files
(println! "\n3. Reading file contents")
(println! "------------------------")

; Read back what we wrote
(if (file-exists? "test-output.txt")
    (do
      (def content (read-file "test-output.txt"))
      (println! "Content of test-output.txt:")
      (println! content))
    (println! "test-output.txt not found"))

(if (file-exists? "calculation-result.txt")
    (do
      (def result-content (read-file "calculation-result.txt"))
      (println! "\nContent of calculation-result.txt:")
      (println! result-content))
    (println! "calculation-result.txt not found"))

;; 4. Configuration File Pattern
(println! "\n4. Configuration file pattern")
(println! "------------------------------")

(def config-file "demo-config.txt")
(def default-config "# Demo Configuration\ndebug=true\nport=8080\nmax_connections=100")

; Create config if it doesn't exist
(if (not (file-exists? config-file))
    (do
      (write-file config-file default-config)
      (println! "✓ Created default configuration file"))
    (println! "✓ Configuration file already exists"))

; Read and display config
(def config (read-file config-file))
(println! "Current configuration:")
(println! config)

;; 5. Data Processing Pipeline
(println! "\n5. Data processing pipeline")
(println! "----------------------------")

; Create sample data file
(def sample-data "apple\nbanana\ncherry\ndate\nelderberry")
(write-file "fruits.txt" sample-data)
(println! "✓ Created sample data file: fruits.txt")

; Read and process the data
(def raw-data (read-file "fruits.txt"))
(def lines (string-split raw-data "\n"))
(def processed-lines (map string-upper lines))
(def processed-data (string-join processed-lines "\n"))

; Write processed data
(write-file "fruits-upper.txt" processed-data)
(println! "✓ Processed data written to fruits-upper.txt")

; Show both versions
(println! "\nOriginal fruits:")
(println! raw-data)
(println! "\nProcessed fruits (uppercase):")
(println! processed-data)

;; 6. Safe File Operations
(println! "\n6. Safe file operations")
(println! "-----------------------")

; Function that safely reads a file with fallback
(defn safe-read-file [filename fallback]
  (if (file-exists? filename)
      (read-file filename)
      fallback))

; Demonstrate safe reading
(def safe-content (safe-read-file "maybe-exists.txt" "Default content when file is missing"))
(println! "Safe read result:" safe-content)

; Function that creates backup before modifying
(defn safe-write-with-backup [filename content]
  (if (file-exists? filename)
      (do
        (def backup-name (string-concat filename ".backup"))
        (def original-content (read-file filename))
        (write-file backup-name original-content)
        (println! "✓ Created backup:" backup-name)))
  
  (write-file filename content)
  (println! "✓ File updated:" filename))

; Demonstrate backup functionality
(safe-write-with-backup "demo-config.txt" 
  "# Updated Demo Configuration\ndebug=false\nport=3000\nmax_connections=50\nupdated=true")

;; 7. Error Handling Demo
(println! "\n7. Error handling")
(println! "-----------------")

; This function demonstrates error handling for file operations
(defn demonstrate-error-handling [filename]
  (println! "Attempting to read:" filename)
  (if (file-exists? filename)
      (do
        (def content (read-file filename))
        (println! "✓ Successfully read" (number->string (length content)) "characters"))
      (println! "✗ File does not exist:" filename)))

(demonstrate-error-handling "test-output.txt")   ; Should succeed
(demonstrate-error-handling "missing-file.txt")  ; Should report missing

(println! "\n=== File Functions Demo Complete ===")
(println! "Check the created files in your current directory:")
(println! "- test-output.txt")
(println! "- calculation-result.txt") 
(println! "- demo-config.txt (and .backup)")
(println! "- fruits.txt")
(println! "- fruits-upper.txt")
