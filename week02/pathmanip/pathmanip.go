//Danny Radosevich
//COSC3750
//Path manipulation example

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	//we're going to take a file path as an arg
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run pathmanip.go <file-path>")
		return
	}
	inputPath := os.Args[1]

	//get the absolute path
	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}
	fmt.Println("Absolute path:", absPath)

}
