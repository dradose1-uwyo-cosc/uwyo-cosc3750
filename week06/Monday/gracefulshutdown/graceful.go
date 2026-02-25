//Danny Radosevich
//COSC3750

//graceful shutdown example

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Application started with PID %d\n", os.Getpid())
	// Wait for a signal
	sig := <-sigChan
	log.Printf("Received signal: %s. Shutting down gracefully...\n", sig)

	// Perform any necessary cleanup here (e.g., close database connections, stop goroutines, etc.)

	log.Println("Cleanup complete. Exiting.")
}
