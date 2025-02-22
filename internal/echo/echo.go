// Package echo implements the functionality for the "echo" Unix tool.
package echo

import (
	"os"      // Used for interacting with standard I/O and process exit.
	"strings" // Provides string manipulation functions.
)

// Run concatenates the provided arguments and writes them to standard output.
func Run(args []string) {
	// Join all command-line arguments into a single string separated by spaces.
	output := strings.Join(args, " ")

	// Write the output to standard output (stdout) with a trailing newline.
	// os.Stdout.WriteString returns the number of bytes written and an error.
	// We check if there was an error during the write operation.
	if _, err := os.Stdout.WriteString(output + "\n"); err != nil {
		// If there is an error, exit the program with a non-zero status code.
		// Note: Exiting immediately may bypass deferred clean-up.
		os.Exit(1)
	}
}
