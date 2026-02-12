// Danny Radosevich
// peeking.go
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run peeking.go <filename>")
		return
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	linecount := 1
	for err == nil {
		var b byte
		b, err = reader.ReadByte()
		if string(b) == "\n" {
			fmt.Print(string(b))
			linecount++
		} else {
			fmt.Print(string(b))
		}
	}
	fmt.Printf("\nTotal lines: %d\n", linecount)
}
