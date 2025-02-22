// Package main serves as the entry point for the Unix tools project.
package main

import (
	"os" // Provides access to command-line arguments.

	// Importing the cat package from the internal project structure.
	"github.com/drunkleen/unix-tools-go/internal/cat"
)

// main is the starting point of the application.
func main() {
	// Check if command-line arguments are provided.
	if len(os.Args) > 1 {
		// Pass all arguments except the program name to cat.Run.
		cat.Run(os.Args[1:])
	}
}
