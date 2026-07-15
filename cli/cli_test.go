package cli

import (
	"path/filepath"
	"testing"
)

// TestCLI_Run_Help checks that passing an invalid command doesn't panic
// and correctly exits or returns an error code.
func TestCLI_Run_Help(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "chain.json")

	// We pass a nonexistent command. Run() should just return an error code or 0.
	// We mainly just want to cover the parser logic to satisfy the coverage check.
	code := Run([]string{"-data=" + testPath, "invalid-cmd"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for invalid command, got %d", code)
	}
}

// TestCLI_Run_Print checks that a valid read-only command executes successfully.
func TestCLI_Run_Print(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "chain.json")

	code := Run([]string{"-data=" + testPath, "print"})
	if code != 0 {
		t.Fatalf("expected exit code 0 for print command, got %d", code)
	}
}
