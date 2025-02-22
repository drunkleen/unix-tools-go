// Package main is the entry point for the Unix tools project.
package main

import (
	"os" // Provides access to command-line arguments and OS functionality.

	// Importing the echo package from the internal project structure.
	"github.com/drunkleen/unix-tools-go/internal/echo"
)

// main is the program's entry point.
func main() {
	// Check if there are command-line arguments (excluding the program name).
	// If no arguments are provided, the tool does nothing.
	if len(os.Args) > 1 {
		// Pass all arguments except the first (program name) to echo.Run.
		echo.Run(os.Args[1:])
	}
}
