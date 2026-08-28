 package analyzer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// CookieAnalyzer provides methods to analyze cookie activity from log streams.
type CookieAnalyzer struct{}

// New constructs and returns a new CookieAnalyzer instance.
func New() *CookieAnalyzer {
	return &CookieAnalyzer{}
}

// GetMostActiveCookiesFromFile opens the log file at the specified filePath
// and parses its contents to find the most active cookies for targetDate.
func (a *CookieAnalyzer) GetMostActiveCookiesFromFile(filePath string, targetDate string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return a.GetMostActiveCookies(file, targetDate)
}

// GetMostActiveCookies processes an input stream (io.Reader) line-by-line using a streaming parser.
// It aggregates frequency counts for cookies matching targetDate (YYYY-MM-DD in UTC).
// Because the input log is assumed to be sorted by timestamp in descending order,
// parsing terminates early once records drop below targetDate.
func (a *CookieAnalyzer) GetMostActiveCookies(r io.Reader, targetDate string) ([]string, error) {
	// Validate targetDate format before starting line processing.
	if _, err := time.Parse("2006-01-02", targetDate); err != nil {
		return nil, fmt.Errorf("invalid date format %q: expected YYYY-MM-DD", targetDate)
	}

	scanner := bufio.NewScanner(r)
	counts := make(map[string]int)

	// Skip the header row (expected: "cookie,timestamp").
	if scanner.Scan() {
		_ = scanner.Text()
	}

	// Stream file contents line by line to maintain O(1) space during parsing.
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}

		cookie := strings.TrimSpace(parts[0])
		timestampStr := strings.TrimSpace(parts[1])

		// Extract YYYY-MM-DD prefix from ISO 8601 timestamp string.
		if len(timestampStr) < 10 {
			continue
		}
		logDate := timestampStr[:10]

		if logDate == targetDate {
			counts[cookie]++
		} else if logDate < targetDate {
			// Early termination: since logs are sorted descending, subsequent entries will be older.
			break
		}
	}

	// Check for any I/O errors encountered by the scanner during reading.
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading stream: %w", err)
	}

	if len(counts) == 0 {
		return []string{}, nil
	}

	// Determine the maximum occurrence count among cookies for targetDate.
	maxCount := 0
	for _, count := range counts {
		if count > maxCount {
			maxCount = count
		}
	}

	// Collect all cookies matching the maximum occurrence count (handles ties).
	var result []string
	for cookie, count := range counts {
		if count == maxCount {
			result = append(result, cookie)
		}
	}

	return result, nil
}