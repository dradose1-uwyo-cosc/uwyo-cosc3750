//Danny Radoseivich
//Buffered input example

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode"
	"unicode/utf8"
)

// List of colors
type color string

const (
	Reset   color = "\x1b[0m"
	Red     color = "\x1b[31m"
	Green   color = "\x1b[32m"
	Yellow  color = "\x1b[33m"
	Blue    color = "\x1b[34m"
	Magenta color = "\x1b[35m"
	Cyan    color = "\x1b[36m"
	White   color = "\x1b[37m"
	Gold    color = "\x1b[38;2;255;196;47m" // #FFC42F
)

func (c color) colorStart(w io.Writer) {
	fmt.Fprint(w, c)
}

func (c color) colorEnd(w io.Writer) {
	fmt.Fprint(w, Reset)
}

func echo(args []string, w io.Writer) {
	for _, arg := range args {
		w.Write([]byte(arg + " "))
	}
	w.Write([]byte("\n"))
}

func isQuote(r rune) bool {
	return r == '"' || r == '\''
}

func scanArgs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start, first := 0, rune(0) //start is the index of the start of the token,
	// first is the first character of the token
	for width := 0; start < len(data); start += width {
		// Iterate over the input data, looking for the start of a token
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		// Decode the next rune (character) from the input data
		if !unicode.IsSpace(r) {
			// If the character is not a space, it is the start of a token
			first = r
			break
		}
	}
	if isQuote(first) {
		start++ //skip the opening quote
	}
	for width, i := 0, start; i < len(data); i += width {
		// Iterate over the input data, looking for the end of the token
		var r rune
		r, width = utf8.DecodeRune(data[i:]) // Decode the next rune (character) from the input data
		if ok := isQuote(first); !ok && unicode.IsSpace(r) || ok && r == first {
			// If we are not in a quoted token and we encounter a space,
			// or if we are in a quoted token and we encounter the closing quote,
			// then we have reached the end of the token
			return i + width, data[start:i], nil
			// Return the index of the next token, the current token, and no error
		}
	}
	if atEOF && len(data) > start {
		// If we have reached the end of the input data and there is still a token to return,
		// return it
		if isQuote(first) {
			err = fmt.Errorf("unterminated quote: %q", first)
		}
		return len(data), data[start:], err
	}
	if isQuote(first) {
		start-- // if we are at the end of the data and we have an unterminated quote,
		// we want to include the opening quote in the token
	}
	return start, nil, nil // Return the index of the next token, no token, and no error
}

func main() {
	lineScanner := bufio.NewScanner(os.Stdin) //scanner to read lines from standard input
	my_writer := os.Stdout                    //writer to write to standard output
	fmt.Print("$: ")
	for {
		// Read one line from input
		if !lineScanner.Scan() {
			return // EOF reached
		}

		// Parse the line into tokens using scanArgs
		tokenScanner := bufio.NewScanner(strings.NewReader(lineScanner.Text()))
		tokenScanner.Split(scanArgs)

		var cmd string    // The command to execute (the first token)
		var args []string // The arguments to the command (the remaining tokens)

		if tokenScanner.Scan() {
			cmd = string(tokenScanner.Bytes()) // The first token is the command
		}
		for tokenScanner.Scan() {
			// The remaining tokens are the arguments
			args = append(args, string(tokenScanner.Bytes()))
		}
		switch cmd {
		case "exit":
			return
		case "color":
			// three options
			// color start <color> - set the color for subsequent output
			// color end - reset the color to default
			// color list - list all available colors
			if len(args) == 0 {
				fmt.Fprintln(my_writer, "Usage: color <start|end|list> [color]")
				continue
			}
			switch args[0] {
			case "start":
				if len(args) != 2 {
					fmt.Fprintln(my_writer, "Usage: color start <color>")
					continue
				}
				switch args[1] {
				case "red":
					Red.colorStart(my_writer)
				case "green":
					Green.colorStart(my_writer)
				case "yellow":
					Yellow.colorStart(my_writer)
				case "blue":
					Blue.colorStart(my_writer)
				case "magenta":
					Magenta.colorStart(my_writer)
				case "cyan":
					Cyan.colorStart(my_writer)
				case "white":
					White.colorStart(my_writer)
				case "gold":
					Gold.colorStart(my_writer)
				default:
					fmt.Fprintf(my_writer, "Unknown color: %s\n", args[1])
				}
			case "end":
				Reset.colorStart(my_writer)
			case "list":
				fmt.Fprintln(my_writer, "Available colors: red, green, yellow, blue, magenta, cyan, white, gold")
			default:
				fmt.Fprintln(my_writer, "Usage: color <start|end|list> [color]")
			}
		case "echo":
			// Echo the arguments back to the user
			echo(args, my_writer)

		default:
			//for now we are going to punt the command to exec
			cmd := exec.Command(cmd, args...) // Create a new command with the given name and arguments
			cmd.Stdout = my_writer            // Set the command's standard output to our writer
			cmd.Stderr = my_writer            // Set the command's standard error to our writer
			err := cmd.Run()                  // Run the command and wait for it to finish
			if err != nil {
				fmt.Fprintf(my_writer, "Error executing command: %v\n", err)
			}
		}
		fmt.Print("$: ")
	}
}
