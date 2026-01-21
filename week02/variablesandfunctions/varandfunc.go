// Danny Radosevich
// Variables and Functions for COSC3750
// Package main defines the entry point for the program.
// Every executable Go program must have a main package.
package main

// Import the fmt package, which provides formatted I/O
// functions such as printing to standard output.
import "fmt"

// print outputs the given message to standard output.
// It is a helper function that wraps fmt.Println.
func print(msg string) {
	// Print the message followed by a newline
	fmt.Println(msg)
}

// add takes two integers as input and returns their sum.
// This demonstrates a function with parameters and a return value.
func add(a int, b int) int {
	return a + b
}

// main is the entry point of the program.
// Execution begins here when the program is run.
func main() {
	// Variable declaration using the var keyword with an explicit type.
	// The variable 'message' is initialized to a string value.
	var message string = "Hello, COSC3750!"
	print(message)

	// Reassign the variable 'message' with a new string value.
	message = "Variables and Functions in Go"
	print(message)

	// Short variable declaration using :=.
	// The compiler infers the type of 'new_message' from the assigned
	// value.
	new_message := "Using short variable declaration"
	print(new_message)
}
