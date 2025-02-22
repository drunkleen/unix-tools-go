// Package cat implements the functionality for the "cat" Unix tool.
package cat

import (
	"bufio"   // Provides buffered I/O for efficient reading.
	"flag"    // Used to parse command-line flags.
	"fmt"     // For formatted I/O operations.
	"os"      // For interacting with the file system and OS I/O.
	"strings" // Provides functions for string manipulation.

	"golang.org/x/term" // For obtaining terminal dimensions.
)

var (
	width int // Terminal width, used for formatting output.
)

// Run is the entry point for the cat functionality.
// It parses flags, determines the source(s) of input (files or stdin),
// and then prints the file contents (optionally with line numbers).
func Run(args []string) {
	var err error
	// Obtain the terminal dimensions to format the output header.
	width, _, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Print error to standard error if terminal size cannot be determined.
		fmt.Fprintf(os.Stderr, "cat: error getting terminal size: %v\n", err)
		return
	}

	// Create a new FlagSet for parsing command-line options specific to "cat".
	fs := flag.NewFlagSet("cat", flag.ExitOnError)
	// Define a boolean flag "-n" to indicate if line numbers should be printed.
	lineNumbers := fs.Bool("n", false, "print line numbers")
	// Parse the provided arguments according to the defined flags.
	fs.Parse(args)

	// Retrieve non-flag arguments, which are interpreted as file names.
	files := fs.Args()
	// If no files are provided, read from standard input.
	for len(files) == 0 {
		printFromReader(os.Stdin, lineNumbers)
		return
	}

	// Iterate over each provided file name.
	for _, file := range files {
		// Process each file and print its contents.
		err := printFile(file, lineNumbers)
		if err != nil {
			// If there's an error opening or reading a file, print it to stderr.
			fmt.Fprintf(os.Stderr, "cat: %v\n", err)
		}
	}
}

// printFile opens the specified file, prints its contents to stdout,
// and optionally adds line numbers. It returns an error if file access fails.
func printFile(fileName string, lineNumbers *bool) error {
	// Open the file in read-only mode.
	file, err := os.Open(fileName)
	if err != nil {
		return err // Return the error to the caller for handling.
	}
	// Ensure the file is closed after processing to free resources.
	defer file.Close()

	// Read from the file and print its contents.
	printFromReader(file, lineNumbers)
	return nil
}

// printFromReader reads from the provided file and prints its content to stdout.
// It prints a header with the file's name and optionally prefixes each line with its line number.
func printFromReader(reader *os.File, lineNumbers *bool) {
	// Create a new scanner to read the input line by line.
	scanner := bufio.NewScanner(reader)
	if *lineNumbers {
		// The header includes a border and centers the file name.
		fmt.Print(
			strings.Repeat("─", 7), "┬", strings.Repeat("─", width-8), "\n",
			strings.Repeat(" ", 7), "│ File: ",
			reader.Name(), "\n",
			strings.Repeat("─", 7), "┼", strings.Repeat("─", width-8), "\n",
		)
	}

	lineCounter := 1 // Initialize a counter for line numbering.
	if *lineNumbers {
		// Iterate over each line of the input.
		for scanner.Scan() {
			// If line numbering is enabled, format the output with a fixed width for numbers.
			// fmt.Printf("%6d │ %s\n", lineCounter, scanner.Text())
			fmt.Printf("%6d │ ", lineCounter)

			for i, t := range scanner.Text() {
				if i%(width-9) == 0 && i != 0 {
					println()
					fmt.Printf("       │ ")
				}
				print(string(t))
			}

			fmt.Println()

			lineCounter++
		}
		// The header includes a border and centers the file name.
		fmt.Print(
			strings.Repeat("─", 7), "┴", strings.Repeat("─", width-8), "\n",
		)
	} else {
		// Iterate over each line of the input.
		for scanner.Scan() {
			// Otherwise, simply print the line.
			fmt.Println(scanner.Text())
		}
	}
	// Check for errors that occurred during scanning.
	if err := scanner.Err(); err != nil {
		// Print any read error to stderr.
		fmt.Fprintf(os.Stderr, "cat: error reading input: %v\n", err)
	}
}
