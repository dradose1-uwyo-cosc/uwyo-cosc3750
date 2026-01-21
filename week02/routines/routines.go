// Danny radosevich
// Channels and Routines for COSC3750
// https://go.dev/tour/concurrency/1
package main // package main defines the entry point for the program.
// Every executable Go program must have a main package.

import (
	"fmt"  // import the fmt package for formatted I/O functions
	"time" // import the time package for time-related functions
)

func print(msg string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(msg)
	}
}

func main() {
	// Start a new goroutine to run the print function concurrently
	go print("go")

	// Call the print function in the main goroutine
	print("pokes")
}
