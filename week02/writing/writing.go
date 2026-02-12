//Danny Radosevich

// writing.go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run writing.go <destination> <contents>")
		return
	}
	destination := os.Args[1]
	contents := os.Args[2]
	err := os.WriteFile(destination, []byte(contents), 0644)
	//0644 is the file permission
	/*
		0644 means:
		- The owner of the file has read and write permissions (6 = 4 + 2).
		- The group members have read-only permissions (4).
		- Others (everyone else) also have read-only permissions (4).
	*/
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("Successfully wrote to", destination)

}
