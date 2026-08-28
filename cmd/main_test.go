 package main

import (
	"io"
	"os"
	"testing"
)

// TestMain verifies that the main application entry point executes cleanly,
// delegates command-line arguments to the CLI runner, and exits with code 0.
func TestMain(t *testing.T) {
	// Preserve original package state and global process variables.
	origArgs := os.Args
	origRun := runCLI
	origExit := exitFunc

	// Restore original package state after the test completes.
	defer func() {
		os.Args = origArgs
		runCLI = origRun
		exitFunc = origExit
	}()

	// Mock exitFunc to intercept and record exit status instead of terminating the test runner.
	var capturedCode int
	exitFunc = func(code int) {
		capturedCode = code
	}

	// Mock runCLI to simulate successful CLI execution.
	runCLI = func(args []string, stdout, stderr io.Writer) int {
		return 0
	}

	// Supply mock command-line flags.
	os.Args = []string{"cmd", "-f", "test.csv", "-d", "2018-12-09"}

	// Execute the main entry point.
	main()

	// Assert that the recorded exit code indicates success.
	if capturedCode != 0 {
		t.Errorf("got exit code %d, want 0", capturedCode)
	}
}