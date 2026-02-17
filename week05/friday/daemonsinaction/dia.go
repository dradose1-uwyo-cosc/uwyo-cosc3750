//Danny Radosevich

//COSC3750

//Example on daemons in action

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

/*
Kill command:
kill $(ps aux | grep "[d]ia" | awk '{print $2}')
*/

var pid = os.Getpid()

func runDaemon() {
	for {
		fmt.Printf("Daemon is running...\n")
		time.Sleep(5 * time.Second)
	}
}

func forkProcess() error {
	cmd := exec.Command(os.Args[0], "daemon")
	cmd.Stdout, cmd.Stderr, cmd.Dir = os.Stdout, os.Stderr, "/"
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Create new session (detach from terminal)
	}
	return cmd.Start()
}

func releaseProcess() error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Release()
}

func main() {
	//first get pid

	fmt.Printf("PID: %d\n", pid)
	//get parent id
	ppid := os.Getppid()
	fmt.Printf("PPID: %d\n", ppid)

	if len(os.Args) != 1 {
		runDaemon()
		return
	}
	if err := forkProcess(); err != nil {
		fmt.Printf("Error forking process: %v\n", err)
		return
	}
	if err := releaseProcess(); err != nil {
		fmt.Printf("Error releasing process: %v\n", err)
		return
	}
	fmt.Printf("Daemon process started successfully\n")

}
