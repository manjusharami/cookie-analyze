 package cli

import (
	"flag"
	"fmt"
	"io"
	"most_active_cookie/internal/analyzer"
)

// Run parses command-line flags and executes the cookie analyzer tool.
// It writes standard output to stdout, error messages to stderr, and returns an exit code (0 for success, 1 for failure).
func Run(args []string, stdout, stderr io.Writer) int {
	// Configure custom flag set to manage argument parsing and direct usage messages to stderr.
	flags := flag.NewFlagSet("most_active_cookie", flag.ContinueOnError)
	flags.SetOutput(stderr)

	// Define command-line flags for file path (-f) and target date (-d).
	filePath := flags.String("f", "", "Path to log CSV file")
	targetDate := flags.String("d", "", "Date in YYYY-MM-DD format (UTC)")

	// Parse incoming arguments and return exit code 1 if flag parsing fails.
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Enforce required flags (-f and -d).
	if *filePath == "" || *targetDate == "" {
		fmt.Fprintln(stderr, "Error: both -f and -d arguments are required.")
		return 1
	}

	// Instantiate analyzer and process the log file for the specified date.
	proc := analyzer.New()
	results, err := proc.GetMostActiveCookiesFromFile(*filePath, *targetDate)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}

	// Output each active cookie result on a new line to stdout.
	for _, cookie := range results {
		fmt.Fprintln(stdout, cookie)
	}

	return 0
}