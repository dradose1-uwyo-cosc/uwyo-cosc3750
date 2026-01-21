//Danny Radosevich
// Values and Pointers for COSC3750

package main // package main defines the entry point for the program.
// Every executable Go program must have a main package.

import "fmt" // import the fmt package for formatted I/O functions

func increment(a int) {
	a = a + 1
}
func incrementPointer(a *int) {
	*a = *a + 1
}

func main() {
	// Demonstrate passing by value
	x := 10
	fmt.Println("Before incrementing by value:", x)
	increment(x)
	fmt.Println("After incrementing by value:", x)

	// Demonstrate passing by pointer
	y := 10
	fmt.Println("Before incrementing by pointer:", y)
	incrementPointer(&y)
	fmt.Println("After incrementing by pointer:", y)
}
