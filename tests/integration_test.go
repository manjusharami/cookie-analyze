 package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLIIntegration(t *testing.T) {
	tempDir := t.TempDir()
	binaryName := "cookie_analyzer"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "../cmd/main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
	}

	logFile := filepath.Join(tempDir, "test_log.csv")
	csvContent := []byte(
		"cookie,timestamp\n" +
			"AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00\n" +
			"SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00\n" +
			"AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00\n",
	)
	if err := os.WriteFile(logFile, csvContent, 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	t.Run("End-to-End Success", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "-f", logFile, "-d", "2018-12-09")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Command failed: %v, stderr: %s", err, stderr.String())
		}

		output := strings.TrimSpace(stdout.String())
		expected := "AtY0laUfhglK3lC7"
		if output != expected {
			t.Errorf("got %q, want %q", output, expected)
		}
	})

	t.Run("Missing Flag Error", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "-f", logFile)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err == nil {
			t.Fatal("Expected command failure, but got success")
		}

		if !strings.Contains(stderr.String(), "Error") {
			t.Errorf("Expected error message in stderr, got: %s", stderr.String())
		}
	})
}