// Package ls implements the functionality for the "ls" Unix tool.
package ls

import (
	"flag"    // For parsing command-line flags.
	"fmt"     // For formatted I/O.
	"os"      // For file system and OS interaction.
	"os/user" // To lookup user and group information.
	"strings"

	// For manipulating file paths.
	"sort" // For sorting directory entries.
	// For string manipulation.
	"syscall" // To access low-level system calls and file metadata.
	"time"    // For handling time and date formatting.

	"golang.org/x/term" // To retrieve terminal size.
)

// Run executes the ls command, handling both default and long format listings.
func Run(args []string) {
	// Create a new FlagSet to handle command-line options for ls.
	fs := flag.NewFlagSet("ls", flag.ExitOnError)
	// Define the `-l` flag for long format listing.
	longFormat := fs.Bool("l", false, "Use a long listing format")
	// Parse the provided arguments.
	fs.Parse(args)

	// Set the target directory; default to the current directory.
	dir := "."
	if fs.NArg() > 0 {
		dir = fs.Arg(0) // Use the first non-flag argument as the directory.
	}

	// Read all entries in the target directory.
	entries, err := os.ReadDir(dir)
	if err != nil {
		// Report error if directory cannot be accessed.
		fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", dir, err)
		return
	}

	// Sort directory entries alphabetically by their name.
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	// Depending on the flag, choose the output format.
	if *longFormat {
		// In long format, first print the total disk blocks used.
		printTotalBlocks(entries)
		// Then print detailed information for each entry.
		for _, entry := range entries {
			printDetailedEntry(entry)
		}
	} else {
		// Otherwise, print entries in a multi-column layout.
		printMultiColumn(entries)
	}
}

// printTotalBlocks calculates and prints the total number of disk blocks used by the files.
func printTotalBlocks(entries []os.DirEntry) {
	var totalBlocks int64

	// Iterate over each entry to accumulate its disk block usage.
	for _, entry := range entries {
		info, err := entry.Info() // Retrieve file info.
		if err != nil {
			continue // Skip if file information cannot be obtained.
		}
		// Access the underlying system stat structure.
		stat := info.Sys().(*syscall.Stat_t)
		totalBlocks += stat.Blocks // Sum up the block count.
	}

	// Print the total blocks converted from 512-byte units to 1K blocks.
	fmt.Printf("total %d\n", totalBlocks/2)
}

// printDetailedEntry prints a detailed listing for a single file, similar to `ls -l`.
func printDetailedEntry(entry os.DirEntry) {
	// Get file info for the entry.
	info, err := entry.Info()
	if err != nil {
		fmt.Printf("ls: error reading file info for %s: %v\n", entry.Name(), err)
		return
	}

	// Convert file info to a syscall.Stat_t to access additional metadata.
	stat := info.Sys().(*syscall.Stat_t)

	// Construct the permissions string.
	perms := info.Mode().Perm().String()
	if info.IsDir() {
		perms = "d" + perms[1:] // Prefix with 'd' for directories.
	} else {
		perms = "-" + perms[1:] // Prefix with '-' for regular files.
	}

	// Retrieve UID and GID as strings.
	uid := fmt.Sprint(stat.Uid)
	gid := fmt.Sprint(stat.Gid)

	// Lookup the username associated with the UID.
	usr, err := user.LookupId(uid)
	if err != nil {
		usr = &user.User{Username: uid} // Fall back to raw UID if lookup fails.
	}

	// Lookup the group name associated with the GID.
	grp, err := user.LookupGroupId(gid)
	if err != nil {
		grp = &user.Group{Name: gid} // Fall back to raw GID if lookup fails.
	}

	// Format the modification time.
	// Use a different format if the file is older than approximately 6 months.
	timeFormat := "Jan _2 15:04"
	if time.Since(info.ModTime()).Hours() > 6*30*24 {
		timeFormat = "Jan _2 2006"
	}

	// Print file details in a format similar to `ls -l`.
	fmt.Printf("%s %d %s %s %4d %s %s\n",
		perms,                             // Permissions string.
		stat.Nlink,                        // Number of hard links.
		usr.Username,                      // Owner's username.
		grp.Name,                          // Group name.
		info.Size(),                       // File size in bytes.
		info.ModTime().Format(timeFormat), // Formatted modification time.
		getFileNameWithIcon(entry),        // File name with an associated icon.
	)
}

// printMultiColumn arranges file entries into a multi-column layout based on the terminal width.
func printMultiColumn(entries []os.DirEntry) {
	// Attempt to get the terminal width.
	width, _, err := term.GetSize(int(syscall.Stdin))
	if err != nil || width < 20 {
		width = 80 // Default to 80 columns if the terminal size is not available.
	}
	var names []string
	maxLen := 0 // Track the longest filename length.
	// Collect file names along with their icons.
	for _, entry := range entries {
		name := getFileNameWithIcon(entry)
		names = append(names, name)

		if len(name) > maxLen {
			maxLen = len(name) // Update max length for padding.
		}
	}

	colWidth := maxLen + 2   // Add padding to the maximum name length.
	cols := width / colWidth // Calculate the number of columns that fit in the terminal.
	if cols == 0 {
		cols = 1 // Ensure at least one column is used.
	}

	// Loop through the names and print them in columns.
	for i, name := range names {
		fmt.Printf("%-*s", colWidth, name) // Left-align within the column width.
		// Insert a newline when a row is complete or at the end of the list.
		if (i+1)%cols == 0 || i == len(names)-1 {
			fmt.Println()
		}
	}
}
