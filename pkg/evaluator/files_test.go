package evaluator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/parser"
	"github.com/leinonen/lisp-interpreter/pkg/tokenizer"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestFileOperations(t *testing.T) {
	// Create temporary directory for test files
	tempDir := t.TempDir()

	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Helper function to evaluate expressions
	evalExpr := func(input string) (types.Value, error) {
		tok := tokenizer.NewTokenizer(input)
		tokens, err := tok.TokenizeWithError()
		if err != nil {
			return nil, err
		}

		p := parser.NewParser(tokens)
		expr, err := p.Parse()
		if err != nil {
			return nil, err
		}

		return evaluator.Eval(expr)
	}

	t.Run("write-file function", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "test-write.txt")
		testContent := "Hello, World!\nThis is a test file."

		// Test writing a file
		result, err := evalExpr(`(write-file "` + testFile + `" "` + testContent + `")`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// write-file should return #t on success
		if result.String() != "#t" {
			t.Errorf("expected #t, got %s", result.String())
		}

		// Verify file was actually written
		content, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("failed to read test file: %v", err)
		}

		if string(content) != testContent {
			t.Errorf("expected %q, got %q", testContent, string(content))
		}
	})

	t.Run("read-file function", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "test-read.txt")
		testContent := "Hello from read test!\nLine 2\nLine 3"

		// Create test file
		err := os.WriteFile(testFile, []byte(testContent), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Test reading the file
		result, err := evalExpr(`(read-file "` + testFile + `")`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should return string content
		if result.String() != testContent {
			t.Errorf("expected %q, got %q", testContent, result.String())
		}
	})

	t.Run("file-exists? function", func(t *testing.T) {
		existingFile := filepath.Join(tempDir, "existing.txt")
		nonExistingFile := filepath.Join(tempDir, "not-existing.txt")

		// Create existing file
		err := os.WriteFile(existingFile, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Test existing file
		result, err := evalExpr(`(file-exists? "` + existingFile + `")`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.String() != "#t" {
			t.Errorf("expected #t for existing file, got %s", result.String())
		}

		// Test non-existing file
		result, err = evalExpr(`(file-exists? "` + nonExistingFile + `")`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.String() != "#f" {
			t.Errorf("expected #f for non-existing file, got %s", result.String())
		}
	})

	t.Run("read-file error handling", func(t *testing.T) {
		nonExistingFile := filepath.Join(tempDir, "does-not-exist.txt")

		// Test reading non-existing file should return error
		_, err := evalExpr(`(read-file "` + nonExistingFile + `")`)
		if err == nil {
			t.Error("expected error when reading non-existing file")
		}
	})

	t.Run("write-file error handling", func(t *testing.T) {
		// Test writing to invalid path (directory that doesn't exist)
		invalidPath := filepath.Join(tempDir, "nonexistent-dir", "test.txt")

		_, err := evalExpr(`(write-file "` + invalidPath + `" "test content")`)
		if err == nil {
			t.Error("expected error when writing to invalid path")
		}
	})

	t.Run("argument validation", func(t *testing.T) {
		// Test missing arguments
		_, err := evalExpr(`(read-file)`)
		if err == nil {
			t.Error("expected error for read-file with no arguments")
		}

		_, err = evalExpr(`(write-file "test.txt")`)
		if err == nil {
			t.Error("expected error for write-file with only one argument")
		}

		_, err = evalExpr(`(file-exists?)`)
		if err == nil {
			t.Error("expected error for file-exists? with no arguments")
		}

		// Test too many arguments
		_, err = evalExpr(`(read-file "test1.txt" "test2.txt")`)
		if err == nil {
			t.Error("expected error for read-file with too many arguments")
		}

		_, err = evalExpr(`(write-file "test.txt" "content" "extra")`)
		if err == nil {
			t.Error("expected error for write-file with too many arguments")
		}

		_, err = evalExpr(`(file-exists? "test1.txt" "test2.txt")`)
		if err == nil {
			t.Error("expected error for file-exists? with too many arguments")
		}
	})

	t.Run("non-string arguments", func(t *testing.T) {
		// Test non-string filename
		_, err := evalExpr(`(read-file 123)`)
		if err == nil {
			t.Error("expected error for read-file with non-string filename")
		}

		_, err = evalExpr(`(write-file 123 "content")`)
		if err == nil {
			t.Error("expected error for write-file with non-string filename")
		}

		_, err = evalExpr(`(write-file "test.txt" 123)`)
		if err == nil {
			t.Error("expected error for write-file with non-string content")
		}

		_, err = evalExpr(`(file-exists? #t)`)
		if err == nil {
			t.Error("expected error for file-exists? with non-string filename")
		}
	})

	t.Run("integration test - read what we write", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "integration-test.txt")
		originalContent := "Line 1\nLine 2\nSpecial chars: !@#$%^&*()\nEmoji: 😀🎉"

		// Write the file
		_, err := evalExpr(`(write-file "` + testFile + `" "` + originalContent + `")`)
		if err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		// Read it back
		result, err := evalExpr(`(read-file "` + testFile + `")`)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		// Should be identical
		if result.String() != originalContent {
			t.Errorf("content mismatch after write/read cycle:\nexpected: %q\ngot: %q", originalContent, result.String())
		}

		// Verify it exists
		exists, err := evalExpr(`(file-exists? "` + testFile + `")`)
		if err != nil {
			t.Fatalf("failed to check file existence: %v", err)
		}

		if exists.String() != "#t" {
			t.Error("file should exist after writing")
		}
	})

	t.Run("empty file operations", func(t *testing.T) {
		emptyFile := filepath.Join(tempDir, "empty.txt")

		// Write empty content
		_, err := evalExpr(`(write-file "` + emptyFile + `" "")`)
		if err != nil {
			t.Fatalf("failed to write empty file: %v", err)
		}

		// Read empty file
		result, err := evalExpr(`(read-file "` + emptyFile + `")`)
		if err != nil {
			t.Fatalf("failed to read empty file: %v", err)
		}

		if result.String() != "" {
			t.Errorf("expected empty string, got %q", result.String())
		}

		// Should still exist
		exists, err := evalExpr(`(file-exists? "` + emptyFile + `")`)
		if err != nil {
			t.Fatalf("failed to check empty file existence: %v", err)
		}

		if exists.String() != "#t" {
			t.Error("empty file should exist after writing")
		}
	})
}
