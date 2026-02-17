# Go Daemon Service

A simple daemon/service management program written in Go that demonstrates process forking, PID file management, and background process control.

## Building

First, compile the program:

```bash
go build -o service
```

This creates an executable named `service` in the current directory.

## Usage

The program accepts commands to manage the daemon. Most commands require `sudo` for permission to create files in `/var/mydaemon/` and `/etc/init.d/`.

### Commands

#### `start` - Start the daemon
Launches the daemon as a background process and records its PID.

```bash
sudo ./service start
```

**What it does:**
- Creates `/var/mydaemon/` directory (if it doesn't exist)
- Spawns the daemon process
- Writes the daemon's PID to `/var/mydaemon/mydaemon.pid`
- Redirects daemon output to log files (`mydaemon.log` and `mydaemon.err`)
- Parent process exits, daemon continues running independently

#### `stop` - Stop the daemon
Terminates the running daemon process.

```bash
sudo ./service stop
```

**What it does:**
- Reads the PID from `/var/mydaemon/mydaemon.pid`
- Sends a KILL signal to the process
- Removes the PID file

#### `status` - Check daemon status
Displays whether the daemon is running and shows its PID.

```bash
sudo ./service status
```

**What it does:**
- Reads the PID from the file
- Attempts to signal the process (without killing it)
- Reports if the daemon is running or not
- Cleans up stale PID files if the process is dead

#### `install` - Install system service integration
Creates an init.d script to enable system-level daemon management.

```bash
sudo ./service install
```

**What it does:**
- Creates `/etc/init.d/mygodaemon` script
- Allows the daemon to be managed via system commands like `service` or `systemctl`

#### `uninstall` - Remove system service integration
Removes the init.d script created by `install`.

```bash
sudo ./service uninstall
```

## Example Workflow

```bash
# 1. Build the program
go build -o service

# 2. Start the daemon (requires sudo)
sudo ./service start
# Output: Daemon started with PID: 12345

# 3. Check status (can do this anytime, even after exiting)
sudo ./service status
# Output: Daemon is running with PID: 12345

# 4. Close the terminal or even reboot—the daemon continues running

# 5. Stop the daemon
sudo ./service stop
# Output: Daemon with PID 12345 stopped

# 6. Verify it's stopped
sudo ./service status
# Output: Daemon is not running
```

## Key Concepts

### PID File
The program stores the daemon's process ID in `/var/mydaemon/mydaemon.pid`. This allows:
- Tracking the daemon across program invocations
- Starting the program, exiting, then restarting and still being able to stop the original daemon
- Checking daemon status without keeping the original process alive

### Background Process
When `start` is called, the daemon is launched with `cmd.Start()` (not `cmd.Run()`), which:
- Returns immediately without waiting for the process
- Allows the parent process to exit
- Leaves the daemon running independently

### Logging
When the daemon starts, its output is redirected to:
- `/var/mydaemon/mydaemon.log` - Standard output
- `/var/mydaemon/mydaemon.err` - Error output

## Permissions

Most operations require root privileges because they involve:
- Writing to `/var/mydaemon/` directory
- Creating/modifying `/etc/init.d/mygodaemon`

Use `sudo` when running commands that create or modify these system directories.

## Notes

- The daemon is a simple infinite loop that sleeps for 1 second between iterations
- Replace the `runDaemon()` function with actual work as needed
- The program demonstrates fundamental daemon concepts; production systems would use more robust service managers like systemd
