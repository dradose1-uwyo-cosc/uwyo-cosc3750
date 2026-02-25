//Danny Radosevich
//COSC3750

//example of using anonymous pipes for inter-process communication

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Create an anonymous pipe
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Printf("Error creating pipe: %v\n", err)
		return
	}
	//we've had one pipe yes
	//but what about second pipe?
	//we can create another pipe for the child to write back to the parent
	r2, w2, err := os.Pipe()
	if err != nil {
		fmt.Printf("Error creating second pipe: %v\n", err)
		return
	}
	var cmds = []*exec.Cmd{
		exec.Command("cat", "anonpipes.go"),
		exec.Command("grep", "pipe"),
		exec.Command("wc", "-l"),
	}
	// Set the stdout of the first command to the write end of the first pipe
	cmds[0].Stdout, cmds[1].Stdin = w, r
	//set the out of the second command to the write end of the second pipe
	cmds[1].Stdout, cmds[2].Stdin = w2, r2
	//and set the out of the third command to os.Stdout so we can see the output
	cmds[2].Stdout = os.Stdout
	//start the commands
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting command: %v\n", err)
			return
		}
	}

	for _, closer := range []interface{ Close() error }{w, r, w2, r2} {
		if err := closer.Close(); err != nil {
			fmt.Printf("Error closing pipe: %v\n", err)
			return
		}
	}
	for i, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			if i == 1 {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
					continue
				}
			}
			fmt.Printf("Error waiting for command: %v\n", err)
			return
		}
	}

}
