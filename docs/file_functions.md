# File Functions Documentation

## Overview
The Lisp interpreter provides built-in file system operations for reading, writing, and checking file existence. These functions enable scripts to interact with the file system directly from Lisp code.

## Available Functions

### `read-file`
Reads the contents of a file and returns it as a string.

**Syntax:**
```lisp
(read-file filename)
```

**Parameters:**
- `filename` - String: The path to the file to read

**Returns:**
- String: The complete contents of the file

**Examples:**
```lisp
; Read a text file
(read-file "data.txt")
; => "Hello, World!\nThis is file content."

; Read configuration
(def config (read-file "config.txt"))
(println! "Config:" config)

; Read and process file content
(def lines (string-split (read-file "numbers.txt") "\n"))
(map string->number lines)
```

**Error Conditions:**
- File does not exist: `"failed to read file filename: no such file or directory"`
- Permission denied: `"failed to read file filename: permission denied"`
- Non-string filename: `"read-file filename must be a string, got <type>"`
- Wrong number of arguments: `"read-file requires exactly 1 argument, got <count>"`

### `write-file`
Writes content to a file, creating the file if it doesn't exist or overwriting if it does.

**Syntax:**
```lisp
(write-file filename content)
```

**Parameters:**
- `filename` - String: The path to the file to write
- `content` - String: The content to write to the file

**Returns:**
- Boolean: `true` on successful write

**Examples:**
```lisp
; Write simple text
(write-file "output.txt" "Hello, World!")
; => true

; Write computed content
(def result (+ 10 20))
(write-file "result.txt" (number->string result))

; Write formatted data
(def numbers (list 1 2 3 4 5))
(def content (string-join (map number->string numbers) "\n"))
(write-file "numbers.txt" content)

; Create a log entry
(write-file "log.txt" 
  (string-concat 
    (current-time) ": " 
    "Process completed successfully\n"))
```

**Error Conditions:**
- Cannot create/write file: `"failed to write file filename: <error>"`
- Non-string filename: `"write-file filename must be a string, got <type>"`
- Non-string content: `"write-file content must be a string, got <type>"`
- Wrong number of arguments: `"write-file requires exactly 2 arguments, got <count>"`
- Directory doesn't exist: `"failed to write file filename: no such file or directory"`

**File Permissions:**
Files are created with permissions `0644` (readable by owner and group, writable by owner only).

### `file-exists?`
Checks whether a file exists at the specified path.

**Syntax:**
```lisp
(file-exists? filename)
```

**Parameters:**
- `filename` - String: The path to check

**Returns:**
- Boolean: `true` if file exists, `false` if it doesn't exist or cannot be accessed

**Examples:**
```lisp
; Simple existence check
(file-exists? "data.txt")
; => true (if file exists)
; => false (if file doesn't exist)

; Conditional file operations
(if (file-exists? "config.txt")
    (read-file "config.txt")
    "default configuration")

; Guard against missing files
(defn safe-read-file [filename]
  (if (file-exists? filename)
      (read-file filename)
      (error (string-concat "File not found: " filename))))

; Check multiple files
(def required-files (list "data.txt" "config.txt" "template.txt"))
(def missing-files 
  (filter (fn [f] (not (file-exists? f))) required-files))
(if (not (empty? missing-files))
    (error (string-concat "Missing files: " (string-join missing-files ", "))))
```

**Error Conditions:**
- Non-string filename: `"file-exists? filename must be a string, got <type>"`
- Wrong number of arguments: `"file-exists? requires exactly 1 argument, got <count>"`

**Note:** This function returns `false` for permission errors or other I/O errors, not just missing files.

## Usage Patterns

### Configuration Management
```lisp
; Load configuration with fallback
(def config-file "app.conf")
(def default-config "debug=false\nport=8080")

(def config 
  (if (file-exists? config-file)
      (read-file config-file)  
      default-config))

(println! "Using configuration:" config)
```

### Data Processing Pipeline
```lisp
; Read → Process → Write pipeline
(defn process-data-file [input-file output-file transform-fn]
  (if (not (file-exists? input-file))
      (error (string-concat "Input file not found: " input-file)))
  
  (def raw-data (read-file input-file))
  (def processed-data (transform-fn raw-data))
  (write-file output-file processed-data)
  
  (println! "Processed" input-file "→" output-file))

; Example usage
(process-data-file 
  "input.txt" 
  "output.txt"
  (fn [data] (string-upper data)))
```

### Backup and Versioning
```lisp
; Create backup before modifying
(defn safe-write-file [filename content]
  (if (file-exists? filename)
      (do 
        (def backup-name (string-concat filename ".backup"))
        (def original-content (read-file filename))
        (write-file backup-name original-content)
        (println! "Created backup:" backup-name)))
  
  (write-file filename content)
  (println! "File updated:" filename))
```

### Log File Management
```lisp
; Append to log file (simulation)
(defn append-to-log [log-file message]
  (def timestamp (current-time))  ; Hypothetical time function
  (def log-entry (string-concat timestamp ": " message "\n"))
  
  (def existing-content 
    (if (file-exists? log-file)
        (read-file log-file)
        ""))
  
  (def new-content (string-concat existing-content log-entry))
  (write-file log-file new-content))

; Usage
(append-to-log "app.log" "Application started")
(append-to-log "app.log" "Processing user request")
```

### File Validation
```lisp
; Validate file structure
(defn validate-csv-file [filename]
  (if (not (file-exists? filename))
      (error "CSV file not found"))
  
  (def content (read-file filename))
  (def lines (string-split content "\n"))
  
  (if (empty? lines)
      (error "CSV file is empty"))
  
  (def header (first lines))
  (def expected-columns 3)
  (def actual-columns (length (string-split header ",")))
  
  (if (not (= actual-columns expected-columns))
      (error (string-concat 
               "Expected " (number->string expected-columns)
               " columns, got " (number->string actual-columns))))
  
  (println! "CSV file validation passed"))
```

## Integration with Module System

File functions work seamlessly with the module system:

```lisp
; Save module to file
(def module-code 
  "(module math-utils
     (export square cube)
     (defn square [x] (* x x))
     (defn cube [x] (* x x x)))")

(write-file "math-utils.lisp" module-code)

; Later, load the saved module
(load "math-utils.lisp")
(import math-utils)
(square 5)  ; => 25
```

## Error Handling Best Practices

### Defensive Programming
```lisp
(defn robust-file-operation [filename]
  (if (not (string? filename))
      (error "Filename must be a string"))
  
  (if (= filename "")
      (error "Filename cannot be empty"))
  
  (if (not (file-exists? filename))
      (error (string-concat "File not found: " filename)))
  
  (read-file filename))
```

### Try-Catch Pattern (Using Error Function)
```lisp
(defn safe-read-config [filename default-content]
  (if (file-exists? filename)
      (do
        ; Attempt to read, but provide fallback for other errors
        (def content (read-file filename))
        (if (= content "")
            default-content
            content))
      default-content))
```

## Performance Considerations

- **File Size**: Large files are read entirely into memory. For very large files, consider implementing streaming or chunked processing at the application level.
- **File Permissions**: Files are created with standard permissions (0644). If different permissions are needed, use system tools or shell commands.
- **Atomic Operations**: Write operations are not atomic. For critical applications, consider writing to a temporary file and then renaming.
- **Encoding**: Files are read and written as UTF-8 strings. Binary file operations are not directly supported.

## Common Pitfalls

1. **Path Separators**: Use forward slashes (`/`) or double backslashes (`\\`) in file paths
2. **Relative Paths**: Paths are relative to the working directory where the interpreter was started
3. **File Locking**: No built-in file locking mechanism - avoid concurrent writes to the same file
4. **Error Recovery**: Always check `file-exists?` before attempting operations on files that might not exist

## See Also

- [String Functions](string_functions.md) - For processing file content
- [Module System](modules.md) - For loading and managing Lisp modules from files  
- [Error Handling](error_function.md) - For robust error management in file operations
- [I/O Operations](operations.md#io-operations) - Other input/output functions
- [Examples](examples.md) - More file operation examples
