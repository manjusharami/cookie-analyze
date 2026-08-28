# 1. Run via binary directly
./cookie_analyzer -f cookie_log.csv -d 2018-12-07

# 2. Run via Makefile with default date (2018-12-09)
make run

# 3. Run via Makefile with a custom date
make run DATE=2018-12-08

# 4. Run via Makefile with a custom file and date
make run FILE=my_log.csv DATE=2018-12-07

# 5. Run test coverage suite (100% statement coverage)
make test-coverage

# 6. Open interactive coverage in your browser
go tool cover -html=coverage.out




# Most Active Cookie Analyzer

A high-performance, maintainable Go CLI tool built to parse cookie logs and determine the most active cookie(s) for a target UTC date.

---

## Problem Description

Given a cookie log file in the following format:

csv
cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00
5UAVanZf6UtGyKVS,2018-12-09T07:25:00+00:00
AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-08T22:03:00+00:00
4sMM2LxV07bPJzwf,2018-12-08T21:30:00+00:00
fbcn5UAVanZf6UtG,2018-12-08T09:30:00+00:00
4sMM2LxV07bPJzwf,2018-12-07T23:30:00+00:00

Write a command line program in your preferred language to process the log file and return the most active cookie for a specific day. Please include a -f parameter for the filename to process and a -d parameter to specify the date.

Execution
$ ./most_active_cookie -f cookie_log.csv -d 2018-12-09

Output to stdout:

AtY0laUfhglK3lC7

Problem Assumptions

Tied Frequency: If multiple cookies meet the criteria, all tied cookies are returned on separate lines.

Dependencies: Uses standard Go libraries only for logic and CLI parsing.

Time Zone: The -d parameter expects a date in UTC (YYYY-MM-DD).

Memory: Sufficient memory is available to process the file.

Sorting: The log file is pre-sorted by timestamp in descending order (most recent occurrence is on the first line).

 ## Project Structure

```text
most_active_cookie/
├── cmd/
│   ├── main.go              # CLI binary entry point
│   └── main_test.go         # Unit tests for entry point
├── internal/
│   ├── analyzer/
│   │   ├── analyzer.go      # Log processing and frequency analysis engine
│   │   └── analyzer_test.go # Comprehensive unit tests for analyzer
│   └── cli/
│       ├── cli.go           # Command-line flag parsing and execution runner
│       └── cli_test.go      # Unit tests for CLI runner and flags
├── tests/
│   └── integration_test.go  # End-to-end integration tests
├── go.mod                   # Go module definition
├── Makefile                 # Automation tasks for building and testing
└── README.md                # Project documentation

Prerequisites
Go 1.18+ installed.

How to Run the Project

1. Build the Binary
Compile the executable into the root directory:
go build -o most_active_cookie ./cmd
Alternatively, using the Makefile:
make build

2.Prepare Log Data (cookie_log.csv)
Create a sample log file in your current working directory:
cat << 'EOF' > cookie_log.csv
cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00
5UAVanZf6UtGyKVS,2018-12-09T07:25:00+00:00
AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-08T22:03:00+00:00
4sMM2LxV07bPJzwf,2018-12-08T21:30:00+00:00
fbcn5UAVanZf6UtG,2018-12-08T09:30:00+00:00
4sMM2LxV07bPJzwf,2018-12-07T23:30:00+00:00
EOF

3.Execute the CLI

./most_active_cookie -f cookie_log.csv -d 2018-12-09

Output:
AtY0laUfhglK3lC7

Tied cookies (returns multiple on separate lines):

./most_active_cookie -f cookie_log.csv -d 2018-12-08

Output:

SAZuXPGUrfbcn5UA
4sMM2LxV07bPJzwf
fbcn5UAVanZf6UtG

4.Testing & Coverage

Run All Unit & Integration Tests

go test ./...

Check 100% Statement Coverage
Generate coverage metrics across application packages:

go test -coverpkg=./cmd/...,./internal/... -coverprofile=coverage.out ./cmd/... ./internal/...
go tool cover -func=coverage.out

Interactive Visual Coverage Report
Open line-by-line statement coverage in your browser:
