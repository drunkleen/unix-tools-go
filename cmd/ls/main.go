// Package main is the entry point for the Unix tools project.
package main

import (
	"os" // Provides access to command-line arguments.

	// Importing the ls package from the internal project structure.
	"github.com/drunkleen/unix-tools-go/internal/ls"
)

// main is the starting point of the application.
func main() {
	// Pass command-line arguments (excluding the program name) to ls.Run.
	ls.Run(os.Args[1:])
}
