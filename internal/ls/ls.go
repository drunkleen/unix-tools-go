// Package ls implements the functionality for the "ls" Unix tool.
package ls

import (
	"flag"          // For parsing command-line flags.
	"fmt"           // For formatted I/O.
	"os"            // For file system and OS interaction.
	"os/user"       // To lookup user and group information.
	"path/filepath" // For manipulating file paths.
	"sort"          // For sorting directory entries.
	"strings"       // For string manipulation.
	"syscall"       // To access low-level system calls and file metadata.
	"time"          // For handling time and date formatting.

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
		return entries[i].Name() < entries[j].Name()
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

// getFileNameWithIcon returns a filename prefixed with an icon based on its type.
func getFileNameWithIcon(entry os.DirEntry) string {
	if entry.IsDir() {
		// Use a folder icon for directories.
		return " " + entry.Name()
	}

	// Convert file extension to lowercase for case-insensitive matching.
	ext := strings.ToLower(filepath.Ext(entry.Name()))

	// Return file name with an appropriate icon based on its extension or specific file names.
	switch {

	case ext == ".go", entry.Name() == "go.mod", entry.Name() == "go.sum":
		return "󰟓 " + entry.Name() // Icon for Go files.

	case entry.Name() == "Dockerfile", entry.Name() == "docker-compose.yml", entry.Name() == ".dockerignore":
		return " " + entry.Name() // Icon for Dockerfile.

	case ext == ".rs", entry.Name() == "cargo.toml":
		return " " + entry.Name() // Icon for Rust files.

	case ext == ".md":
		return " " + entry.Name() // Icon for Markdown files.

	case ext == ".json":
		return " " + entry.Name() // Icon for JSON files.

	case ext == ".toml":
		return " " + entry.Name() // Icon for TOML files.

	case ext == ".css":
		return " " + entry.Name() // Icon for CSS files.

	case ext == ".html":
		return " " + entry.Name() // Icon for HTML files.

	case ext == ".js":
		return " " + entry.Name() // Icon for JavaScript files.

	case ext == ".pdf":
		return " " + entry.Name() // Icon for PDF files.

	case ext == ".txt":
		return "󰦨 " + entry.Name() // Icon for text files.

	case ext == ".git", ext == ".github", entry.Name() == ".gitignore":
		return " " + entry.Name() // Icon for Git-related files.

	// Additional universal file types
	case ext == ".png", ext == ".jpg", ext == ".jpeg", ext == ".gif", ext == ".bmp", ext == ".svg":
		return " " + entry.Name() // Icon for image files.

	case ext == ".mp4", ext == ".mkv", ext == ".avi", ext == ".mov", ext == ".wmv":
		return "󰃽 " + entry.Name() // Icon for video files.

	case ext == ".mp3", ext == ".wav", ext == ".flac", ext == ".ogg", ext == ".aac":
		return " " + entry.Name() // Icon for audio files.

	case ext == ".zip", ext == ".tar", ext == ".gz", ext == ".rar", ext == ".7z":
		return " " + entry.Name() // Icon for archive files.

	case ext == ".doc", ext == ".docx":
		return " " + entry.Name() // Icon for document files.

	case ext == ".xls", ext == ".xlsx":
		return " " + entry.Name() // Icon for spreadsheet files.

	case ext == ".ppt", ext == ".pptx":
		return "󱎐 " + entry.Name() // Icon for presentation files.

	case ext == ".sh":
		return " " + entry.Name() // Icon for shell scripts.

	case ext == ".c", ext == ".cpp", ext == ".h", ext == ".hpp":
		return " " + entry.Name() // Icon for C/C++ source files.

	case ext == ".py":
		return " " + entry.Name() // Icon for Python files.

	case ext == ".java":
		return " " + entry.Name() // Icon for Java files.

	case entry.Name() == "Makefile":
		return " " + entry.Name() // Icon for Makefile.

	default:
		return " " + entry.Name() // Default file icon.

	}
}
