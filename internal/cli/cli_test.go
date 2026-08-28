 package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"most_active_cookie/internal/cli"
)

// TestRun tests all command-line flag handling, error scenarios, and output formatting for cli.Run.
func TestRun(t *testing.T) {
	// Setup a temporary log file for testing execution.
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test_log.csv")
	csvContent := []byte(
		"cookie,timestamp\n" +
			"AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00\n" +
			"SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00\n" +
			"AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00\n",
	)
	if err := os.WriteFile(logFile, csvContent, 0644); err != nil {
		t.Fatalf("failed to create temp log file: %v", err)
	}

	// Verifies successful parsing, logic execution, and stdout output.
	t.Run("Success path returns exit code 0 and prints results", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		args := []string{"-f", logFile, "-d", "2018-12-09"}

		code := cli.Run(args, &stdout, &stderr)

		if code != 0 {
			t.Errorf("got exit code %d, want 0", code)
		}
		expectedOutput := "AtY0laUfhglK3lC7\n"
		if stdout.String() != expectedOutput {
			t.Errorf("stdout got %q, want %q", stdout.String(), expectedOutput)
		}
		if stderr.Len() > 0 {
			t.Errorf("expected empty stderr, got %q", stderr.String())
		}
	})

	// Verifies that missing the required -f flag results in exit code 1 and error messaging.
	t.Run("Missing -f flag returns exit code 1", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		args := []string{"-d", "2018-12-09"}

		code := cli.Run(args, &stdout, &stderr)

		if code != 1 {
			t.Errorf("got exit code %d, want 1", code)
		}
		if !strings.Contains(stderr.String(), "Error: both -f and -d arguments are required.") {
			t.Errorf("stderr missing expected error message, got %q", stderr.String())
		}
	})

	// Verifies that missing the required -d flag results in exit code 1 and error messaging.
	t.Run("Missing -d flag returns exit code 1", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		args := []string{"-f", logFile}

		code := cli.Run(args, &stdout, &stderr)

		if code != 1 {
			t.Errorf("got exit code %d, want 1", code)
		}
		if !strings.Contains(stderr.String(), "Error: both -f and -d arguments are required.") {
			t.Errorf("stderr missing expected error message, got %q", stderr.String())
		}
	})

	// Verifies that passing unknown command-line flags triggers parsing failure and returns exit code 1.
	t.Run("Invalid flag argument returns exit code 1", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		args := []string{"-unknown-flag"}

		code := cli.Run(args, &stdout, &stderr)

		if code != 1 {
			t.Errorf("got exit code %d, want 1", code)
		}
	})

	// Verifies file open errors return exit code 1 with stderr output.
	t.Run("Non-existent log file returns exit code 1", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		args := []string{"-f", "missing_file.csv", "-d", "2018-12-09"}

		code := cli.Run(args, &stdout, &stderr)

		if code != 1 {
			t.Errorf("got exit code %d, want 1", code)
		}
		if !strings.Contains(stderr.String(), "Error:") {
			t.Errorf("stderr expected error message, got %q", stderr.String())
		}
	})

	// Verifies date formatting errors return exit code 1 with stderr output.
	t.Run("Invalid date string returns exit code 1", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		args := []string{"-f", logFile, "-d", "invalid-date"}

		code := cli.Run(args, &stdout, &stderr)

		if code != 1 {
			t.Errorf("got exit code %d, want 1", code)
		}
		if !strings.Contains(stderr.String(), "Error:") {
			t.Errorf("stderr expected error message, got %q", stderr.String())
		}
	})
}