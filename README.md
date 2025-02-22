# Unix Tools in Go

A collection of reimagined Unix command-line tools implemented in Go. This project includes Go versions of popular commands like `echo`, `cat`, and `ls`, each written with clarity and modular design in mind.

---

## Overview

This repository is an exploration into reimplementing classic Unix tools using the Go programming language. The goal is to deepen understanding of system programming and the inner workings of command-line utilities. The project also demonstrates how Go's robust standard library and third-party packages (like `golang.org/x/term`) can be leveraged to build tools that are both powerful and efficient.

**Included Tools:**
- **echo**: Outputs the provided text to standard output.
- **cat**: Reads files (or standard input) and outputs their content, with optional line numbering and formatted headers.
- **ls**: Lists directory contents. Supports both a default multi-column format and a detailed long format (similar to `ls -l`), including file permissions, user/group ownership, file size, modification time, and a rich set of file icons for visual enhancement.

---

## Features

- **Modular Architecture:**  
  Each command is encapsulated in its own package, making the codebase easy to navigate and extend.
  
- **Rich Output Formatting:**  
  - **ls Command:** Displays files in a multi-column layout that adapts to terminal width. The long format provides detailed file metadata.
  - **File Icons:** A wide variety of file types are identified with appropriate icons (for images, videos, audio, archives, documents, source code, and more), offering an enhanced visual experience.
  
- **Terminal Awareness:**  
  Dynamic formatting based on terminal size ensures that the output remains neat and readable.

---

## Getting Started

### Prerequisites

- **Go:** Version 1.16 or newer is required.
- **Operating System:** Unix-like systems are recommended for the best experience.

### Installation

Clone the repository:

```bash
git clone https://github.com/drunkleen/unix-tools-go.git
cd unix-tools-go
```

**use `Makefile` to build:**

```bash
make all
```

**Build echo:**

```bash
go build -o bin/echo ./cmd/echo
```
or
```bash
make echo
```

**Build cat:**

```bash
go build -o bin/cat ./cmd/cat
```
or
```bash
make cat
```

**Build ls:**

```bash
go build -o bin/ls ./cmd/ls
```
or
```bash
make ls
```

---

## Usage

After building the project, you can use each tool as follows:

### echo

Prints the provided arguments to standard output.

```bash
./bin/echo Hello, world!
```

### cat

Reads file contents and prints them. Supports an optional `-n` flag for line numbering.

```bash
# Read from a file with line numbers
./bin/cat -n file.txt

# Read from standard input
./bin/cat file.txt | ./cat
```

### ls

Lists directory contents. Use `-l` for a detailed view.

```bash
# Default multi-column listing of the current directory:
./bin/ls

# Detailed listing (similar to ls -l):
./bin/ls -l /path/to/directory
```

---

## Video Demonstration

I created a detailed video on my YouTube channel that walks through the development of these tools in Persian Language.

[![Watch the Video](https://i9.ytimg.com/vi/7hc2LzaOt6o/maxresdefault.jpg?v=67b99eae&sqp=CMzC5r0G&rs=AOn4CLC9Orc8j7x0v-9xA6WGDbar6mdbng)](https://youtu.be/7hc2LzaOt6o)

---

## Project Structure

```plaintext
unix-tools-go/
├── cmd
│   ├── cat
│   │   └── main.go
│   ├── echo
│   │   └── main.go
│   └── ls
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── cat
│   │   └── cat.go
│   ├── echo
│   │   └── echo.go
│   └── ls
│       └── ls.go
└── README.md

```

---

## Contributing

Contributions are welcome! If you have ideas to improve these tools or want to add new features, feel free to fork the repository and submit a pull request. Please ensure that your changes follow standard Go practices and include proper tests and documentation.

---

## License

This project is licensed under the GNU License. See the [LICENSE](LICENSE) file for details.

---

Enjoy exploring these Unix tools reimagined in Go, and happy coding!