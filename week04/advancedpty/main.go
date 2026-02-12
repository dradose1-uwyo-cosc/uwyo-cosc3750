//Danny Radoseivich
//Buffered input example

package main

import (
	"advancedpty/commands"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	lineScanner := bufio.NewScanner(os.Stdin) //scanner to read lines from standard input
	my_writer := os.Stdout                    //writer to write to standard output

	for {
		fmt.Print("$: ")
		// Read one line from input
		if !lineScanner.Scan() {
			return // EOF reached
		}

		// Parse the line into tokens using scanArgs
		tokenScanner := bufio.NewScanner(strings.NewReader(lineScanner.Text()))
		tokenScanner.Split(commands.ScanArgs)

		var cmd string    // The command to execute (the first token)
		var args []string // The arguments to the command (the remaining tokens)

		if tokenScanner.Scan() {
			cmd = string(tokenScanner.Bytes()) // The first token is the command
		}
		for tokenScanner.Scan() {
			// The remaining tokens are the arguments
			args = append(args, string(tokenScanner.Bytes()))
		}

		// Check if the command is a registered command
		found := commands.GetCommand(cmd)
		if found != nil {
			// If it is a registered command, execute it
			found.Execute(nil, my_writer, args...)
		} else {
			//if not registered use an exec
			cmd := exec.Command(cmd, args...) // Create a new command with the given name and arguments
			cmd.Stdout = my_writer            // Set the command's standard output to our writer
			cmd.Stderr = my_writer            // Set the command's standard error to our writer
			err := cmd.Run()                  // Run the command and wait for it to finish
			if err != nil {
				fmt.Fprintf(my_writer, "Error executing command: %v\n", err)
			}
		}
	}
}
