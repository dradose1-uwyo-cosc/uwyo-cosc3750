//Danny Radosevich
//Hello world for COSC3750

package main // package main marks this file as part of the main package,
//  which builds an executable program, not a reusable library.

import "fmt" // import brings in other packages;
// "fmt" is the standard formatting package used for printing text
// to the screen.

// func main is the entry point of the program.
// When you run the compiled program, execution starts in this main
// function.
func main() {
	// fmt.Println calls the Println function from the fmt package.
	// It prints the string "Hello, World!" followed by a
	// newline to the console.
	fmt.Println("Hello, World!")
}
