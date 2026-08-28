 package analyzer_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"most_active_cookie/internal/analyzer"
)

// mockCSV provides a standard time-sorted log sample for testing cookie parsing.
const mockCSV = `cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00
5UAVanZf6UtGyKVS,2018-12-09T07:25:00+00:00
AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-08T22:03:00+00:00
4sMM2LxV07bPJzwf,2018-12-08T21:30:00+00:00
fbcn5UAVanZf6UtG,2018-12-08T09:30:00+00:00
4sMM2LxV07bPJzwf,2018-12-07T23:30:00+00:00`

// errReader implements io.Reader to simulate a stream read failure during scanner processing.
type errReader struct{}

// Read always returns a simulated error to trigger scanner error handling paths.
func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read failure")
}

// TestGetMostActiveCookies verifies the functionality and error handling of the CookieAnalyzer.
func TestGetMostActiveCookies(t *testing.T) {
	proc := analyzer.New()

	// Verifies that a single clear winner is correctly identified for a specific date.
	t.Run("Single most active cookie", func(t *testing.T) {
		reader := strings.NewReader(mockCSV)
		got, err := proc.GetMostActiveCookies(reader, "2018-12-09")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 || got[0] != "AtY0laUfhglK3lC7" {
			t.Errorf("got %v, want ['AtY0laUfhglK3lC7']", got)
		}
	})

	// Verifies that all tied cookies are returned when multiple share the max frequency.
	t.Run("Tied most active cookies", func(t *testing.T) {
		reader := strings.NewReader(mockCSV)
		got, err := proc.GetMostActiveCookies(reader, "2018-12-08")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 3 {
			t.Fatalf("got %d cookies, expected 3", len(got))
		}
	})

	// Verifies that querying a date with no matching logs returns an empty slice without error.
	t.Run("Non-existent date", func(t *testing.T) {
		reader := strings.NewReader(mockCSV)
		got, err := proc.GetMostActiveCookies(reader, "2018-12-10")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("got %v, want empty slice", got)
		}
	})

	// Verifies that providing a non YYYY-MM-DD date string triggers a validation error.
	t.Run("Invalid date string format", func(t *testing.T) {
		reader := strings.NewReader(mockCSV)
		_, err := proc.GetMostActiveCookies(reader, "invalid-date")
		if err == nil {
			t.Error("expected error for malformed date, got nil")
		}
	})

	// Verifies that blank lines, malformed rows, and short timestamps are safely ignored.
	t.Run("Handles malformed lines and short timestamps gracefully", func(t *testing.T) {
		malformedCSV := "cookie,timestamp\n\nmalformed_row\ncookie_1,2018-12-09\ncookie_2,short\n"
		reader := strings.NewReader(malformedCSV)
		got, err := proc.GetMostActiveCookies(reader, "2018-12-09")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 || got[0] != "cookie_1" {
			t.Errorf("got %v, want ['cookie_1']", got)
		}
	})

	// Verifies that attempting to read a non-existent file path returns an OS error.
	t.Run("GetMostActiveCookiesFromFile missing file error", func(t *testing.T) {
		_, err := proc.GetMostActiveCookiesFromFile("non_existent_file.csv", "2018-12-09")
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})

	// Verifies file-based parsing using a temporary file on disk.
	t.Run("GetMostActiveCookiesFromFile valid file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "log.csv")
		if err := os.WriteFile(filePath, []byte(mockCSV), 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		got, err := proc.GetMostActiveCookiesFromFile(filePath, "2018-12-09")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 || got[0] != "AtY0laUfhglK3lC7" {
			t.Errorf("got %v, want ['AtY0laUfhglK3lC7']", got)
		}
	})

	// Verifies stream reading error handling when bufio.Scanner fails.
	t.Run("Reader stream error returns error", func(t *testing.T) {
		reader := &errReader{}
		_, err := proc.GetMostActiveCookies(reader, "2018-12-09")
		if err == nil {
			t.Error("expected error from failed reader stream, got nil")
		}
	})
}