package executor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/interpreter"
)

func TestExecuteFileSimple(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Create test file
	filename := filepath.Join(tempDir, "test.lisp")
	err := os.WriteFile(filename, []byte("(+ 1 2)"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Execute the file
	interpreter := interpreter.NewInterpreter()
	err = ExecuteFile(interpreter, filename)
	if err != nil {
		t.Fatalf("ExecuteFile failed: %v", err)
	}
}

func TestExecuteFileError(t *testing.T) {
	interpreter := interpreter.NewInterpreter()

	// Test with non-existent file
	err := ExecuteFile(interpreter, "non-existent-file.lisp")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read file") {
		t.Errorf("Expected 'failed to read file' error, got: %v", err)
	}
}
