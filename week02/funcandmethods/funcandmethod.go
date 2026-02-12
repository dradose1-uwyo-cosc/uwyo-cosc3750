//Danny Radosevich
//Functions and Methods for COSC3750

package main // package main defines the entry point for the program.
// Every executable Go program must have a main package.

import "fmt" // import the fmt package for formatted I/O functions

// foo is a function that takes an integer 'a' as input
// and returns an integer 'b' and a boolean 'err'.
// It demonstrates named return values.
func foo(a int) (b int, err bool) {
	if a < 0 {
		err = true
		return
	}
	b = a + 5
	return b, false
}

type S string

type Pair struct {
	first, second int
}

// upper is a method with a receiver of type S.
// It converts the string to uppercase and returns it.
func (msg S) upper() string {
	to_return := ""
	for _, chr := range msg {
		if chr >= 'a' && chr <= 'z' {
			to_return += string(chr - ('a' - 'A'))
		} else {
			to_return += string(chr)
		}
	}
	return to_return
}

func main() {
	foo(5)
	// Anonymous function assigned to a variable 'print'
	print := func(msg string) {
		fmt.Println(msg)
	}
	print("Functions and Methods in Go")
	a := S("GO Pokes")
	s := a.upper() // Call the upper method on type S
	print(s)

}
