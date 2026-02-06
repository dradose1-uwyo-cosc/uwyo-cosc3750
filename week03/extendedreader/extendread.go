//Danny Radosevich
//Examples of extending the reader
//Also including output colors

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// QueryWriter wraps an io.Writer and highlights query matches in color
type QueryWriter struct {
	Query     []byte // The byte sequence to search for
	io.Writer        // Underlying writer (embedded interface)
}

// Write implements io.Writer interface with custom query highlighting
// Searches each line for the query string and outlines all matches
func (qw QueryWriter) Write(b []byte) (n int, err error) {
	// Split input into lines
	lines := bytes.Split(b, []byte{'\n'})
	// Process each line
	for _, line := range lines {
		if len(qw.Query) == 0 {
			// Empty query, write line as-is
			_, err = qw.Writer.Write(append(line, '\n'))
			if err != nil {
				return 0, err
			}
			continue
		}
		// Find the query in this line and highlight all matches
		// Write the line outlining all matches:
		// 1. Text before match
		// 2. Opening bracket
		// 3. Match (gold)
		// 4. Closing bracket
		// Repeat for each match, then add newline
		start := 0
		for {
			idx := bytes.Index(line[start:], qw.Query)
			if idx == -1 {
				break
			}
			idx += start
			for _, s := range [][]byte{
				line[start:idx],                 // Text before query
				[]byte("\x1b[38;2;255;196;37m"), // ANSI escape: gold foreground (#FFC425)
				line[idx : idx+len(qw.Query)],   // Query (gold)
				[]byte("\x1b[0m"),               // ANSI escape: reset to default
			} {
				_, err = qw.Writer.Write(s)
				if err != nil {
					return 0, err
				}
			}
			start = idx + len(qw.Query)
		}
		// Write any remaining text and newline
		_, err = qw.Writer.Write(append(line[start:], '\n'))
		if err != nil {
			return 0, err
		}
	}
	return len(b), nil
}

func main() {
	// Check for required command-line arguments
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <file> <query> \n", os.Args[0])
		return
	}

	// Parse arguments
	query := []byte(strings.Join(os.Args[2:], " ")) // Convert query to bytes (supports multi-word)
	filename := os.Args[1]                          // File to search

	fmt.Printf("Searching for %q in file %q\n", query, filename)

	// Open the input file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create QueryWriter that wraps stdout
	qw := QueryWriter{
		Query:  query,
		Writer: os.Stdout,
	}

	// Copy file contents through QueryWriter
	// io.Copy reads from file and calls qw.Write() with each chunk
	_, err = io.Copy(qw, file)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}
}
