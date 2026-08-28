package main

import (
	"os"

	"most_active_cookie/internal/cli"
)

// Hooks for testing process execution and exit codes without calling os.Exit directly.
var (
	runCLI   = cli.Run
	exitFunc = os.Exit
)

func main() {
	exitCode := runCLI(os.Args[1:], os.Stdout, os.Stderr)
	exitFunc(exitCode)
}