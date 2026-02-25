//Danny Radosevich
//COSC3750

//Child exit codes

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func exitStatus(state *os.ProcessState) int {
	if status, ok := state.Sys().(syscall.WaitStatus); ok {
		return status.ExitStatus()
	}
	return -1 // Return -1 if we can't get the exit status
}
func ProcessState(e error) *os.ProcessState {
	if exitErr, ok := e.(*exec.ExitError); ok {
		return exitErr.ProcessState
	}
	return nil
}

func main() {
	args := os.Args
	cmd := exec.Command("ls", args[1])
	if err := cmd.Run(); err != nil {
		if status := exitStatus(cmd.ProcessState); status != -1 {
			fmt.Printf("Child process exited with status: %d\n", status)
		} else {
			fmt.Printf("Error running command: %v\n", err)
		}
	} else {
		fmt.Printf("Child process exited successfully\n")
		status := exitStatus(cmd.ProcessState)
		fmt.Printf("Child process exit status: %d\n", status)
	}
}
