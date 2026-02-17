//Danny Radosevich
//COSC3750

//Example on creating a service/daemon
/*
can start with a simple init.d script

#!/bin/bash
"path/to/daemon" $1
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

// Global variables for daemon management
var (
	bin     string // Path to the daemon binary
	cmd     string // Command passed as argument (install, start, stop, etc.)
	ErrSudo error  // Error message prompting for sudo if permission denied
)

// Path to the init.d script that will be created to manage the daemon
const initdFile = "/etc/init.d/mygodaemon"

// Constants for daemon state and logging
const (
	varDir  = "/var/mydaemon/" // Directory where pid file and logs are stored
	pidFile = "mydaemon.pid"   // File name for storing the daemon's process ID
)

// writePid writes the daemon's process ID to the pid file for tracking
func writePid(pid int) (err error) {
	// Open or create the pid file with write permissions
	f, err := os.OpenFile(filepath.Join(varDir, pidFile), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// Write the process ID as a string to the file
	if _, err = fmt.Fprintf(f, "%d", pid); err != nil {
		return err
	}
	return nil
}

// getPid retrieves the daemon's process ID from the pid file
func getPid() (pid int, err error) {
	// Read the entire pid file
	b, err := os.ReadFile(filepath.Join(varDir, pidFile))
	if err != nil {
		return 0, err
	}
	// Convert the string content to an integer PID
	if pid, err = strconv.Atoi(string(b)); err != nil {
		return 0, fmt.Errorf("Invalid PID value: %s", string(b))
	}
	return pid, nil
}

// init runs before main() to initialize global variables
func init() {
	// Get the absolute path to the daemon binary itself
	p, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	bin = p
	// Extract the command from command line arguments (if provided)
	if len(os.Args) != 1 {
		cmd = os.Args[1]
	}
	// Create a helpful error message that suggests using sudo if needed
	ErrSudo = fmt.Errorf("Give `sudo %s %s` a try", bin, cmd)
}

// install creates an init.d script to enable daemon management via system service commands
func install() error {
	// Check if the init.d file already exists
	_, err := os.Stat(initdFile)
	if err == nil {
		return errors.New(fmt.Sprintf("File %s already exists", initdFile))
	}
	// Create the init.d file with executable permissions (0755)
	f, err := os.OpenFile(initdFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		if !os.IsPermission(err) {
			return err
		}
		// If permission denied, suggest using sudo
		return ErrSudo
	}
	defer f.Close()
	// Write the bash script that will call the daemon with the command argument
	if _, err := f.WriteString(fmt.Sprintf("#!/bin/bash\n%s %s $1\n", bin, cmd)); err != nil {
		return err
	}
	fmt.Println("Daemon", bin, "installed")
	return nil
}

// uninstall removes the init.d script, disabling daemon system service integration
func uninstall() error {
	// Check if the init.d file exists
	_, err := os.Stat(initdFile)
	if err != nil && os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("File %s does not exist", initdFile))
	}
	// Remove the init.d file
	if err := os.Remove(initdFile); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		// If permission denied, suggest using sudo
		return ErrSudo
	}
	fmt.Println("Daemon", bin, "uninstalled")
	return nil
}

// status checks if the daemon is currently running and displays its PID
func status() (err error) {
	var pid int
	// Defer function to print status at the end (either running or not running)
	defer func() {
		if pid == 0 {
			fmt.Println("Daemon is not running")
		}
		fmt.Printf("Daemon is running with PID: %d\n", pid)
	}()
	// Get the PID from the pid file
	pid, err = getPid()
	if err != nil {
		return err
	}
	// Find the process with this PID
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	// Send signal 0 to check if the process exists (does not actually kill)
	if err := p.Signal(os.Signal(syscall.Signal(0))); err != nil {
		fmt.Println(pid, "not found - removing pid file")
		os.Remove(filepath.Join(varDir, pidFile))
	}
	return nil
}

// start launches the daemon as a background process and records its PID
func start() (err error) {
	// Permissions for creating log files (create if missing, append mode, write-only)
	const perm = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	// Create the /var/mydaemon directory if it doesn't exist
	if err = os.MkdirAll(varDir, 0755); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		// If permission denied, suggest using sudo
		return ErrSudo
	}
	// Create a command to run the daemon with "run" argument
	cmd := exec.Command(bin, "run")
	// Redirect stdout to the log file
	cmd.Stdout, err = os.OpenFile(filepath.Join(varDir, "mydaemon.log"), perm, 0644)
	if err != nil {
		return err
	}
	// Redirect stderr to the error log file
	cmd.Stderr, err = os.OpenFile(filepath.Join(varDir, "mydaemon.err"), perm, 0644)
	if err != nil {
		return err
	}
	// Set working directory to root to avoid keeping any directory open
	cmd.Dir = "/"
	// Start the daemon in the background
	if err = cmd.Start(); err != nil {
		return err
	}
	// Save the daemon's process ID for later management
	if err = writePid(cmd.Process.Pid); err != nil {
		return err
	}
	fmt.Printf("Daemon started with PID: %d\n", cmd.Process.Pid)
	return nil
}

// stop terminates the running daemon process and removes its PID file
func stop() (err error) {
	// Retrieve the daemon's PID from the pid file
	pid, err := getPid()
	if err != nil {
		return err
	}
	// Find the process with this PID
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	// Send a KILL signal to terminate the daemon
	if err = p.Signal(os.Kill); err != nil {
		return err
	}
	// Remove the pid file since the daemon is no longer running
	if err = os.Remove(filepath.Join(varDir, pidFile)); err != nil {
		return err
	}
	fmt.Printf("Daemon with PID %d stopped\n", pid)
	return nil
}

// runDaemon is the main loop that runs continuously when the daemon is active
// This is where actual daemon work would be implemented
func runDaemon() error {
	fmt.Println("Run....")
	// Infinite loop: the daemon runs until stopped externally
	for {
		// Sleep for 1 second between iterations (placeholder for actual work)
		time.Sleep(time.Second)
	}
	return nil
}

// main is the entry point that routes commands to appropriate daemon management functions
func main() {
	var err error
	// Switch on the command provided as an argument
	switch cmd {
	case "run":
		// Run the daemon's main loop (called by the start function)
		err = runDaemon()
	case "install":
		// Create the init.d script for system integration (requires sudo)
		err = install()
	case "uninstall":
		// Remove the init.d script (requires sudo)
		err = uninstall()
	case "status":
		// Check if the daemon is currently running
		err = status()
	case "start":
		// Start the daemon as a background process (may require sudo)
		err = start()
	case "stop":
		// Stop the running daemon process (may require sudo)
		err = stop()
	default:
		// Show usage if an invalid command is provided
		fmt.Printf("Usage: %s [install|uninstall|status|start|stop]\n", bin)
		return
	}

	// Display any errors that occurred
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}
