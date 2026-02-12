// Danny Radosevich
// COSC3750
// reading files
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	//open a file
	file, err := os.Open("ragtime.txt") //read only
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() //make sure to close the file when done
	//defer waits until the surrounding function returns
	//be careful with defer in loops, or when dealing with many files
	//read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	//print the content to the console
	fmt.Println(string(content))
}
