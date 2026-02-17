//Danny Radosevich
//COSC3750
//Buffering to a child process using standard input

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	//create a buffer to hold the data we want to send to the child process
	b := bytes.NewBuffer(nil) //initally empty buffer
	//create the command
	cmd := exec.Command("cat") //cat will read from stdin and write to stdout
	//set the standard input of the command to the buffer
	cmd.Stdin = b
	cmd.Stdout = os.Stdout //set the standard output of the command to the standard output of the parent process
	//write some data to the buffer
	b.WriteString("Hello, World!\n")
	b.WriteString("Using the memory space\n")
	b.WriteString("of the parent process\n")
	// Print the bytes in the buffer (memory space)
	b.WriteString(fmt.Sprintf("Buffer memory contents: %v\n", b.Bytes()))
	//start the command
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}
	cmd.Wait() //wait for the command to finish
}
