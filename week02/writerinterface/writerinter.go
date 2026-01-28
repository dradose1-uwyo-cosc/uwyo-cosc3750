//Danny Radosevich
// writerinter.go

package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run writerinter.go <source> <destination>")
		return
	}
	source := os.Args[1]
	destination := os.Args[2]
	//open source file
	src, err := os.Open(source)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}

	//go to end of source file
	curr, err := src.Seek(0, io.SeekEnd)
	if err != nil {
		fmt.Println("Error seeking to end of source file:", err)
		return
	}

	//create destination file
	destFile, err := os.Create(destination)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}
	defer destFile.Close()

	//copy from source to destination in reverse
	step := int64(1)
	for curr > 0 {
		if curr < step {
			step = curr
		}
		_, err = src.Seek(curr-step, io.SeekStart)
		if err != nil {
			fmt.Println("Error seeking in source file:", err)
			return
		}
		var b [1]byte
		_, err = src.Read(b[:])
		if err != nil {
			fmt.Println("Error reading from source file:", err)
			return
		}
		_, err = destFile.Write(b[:])
		if err != nil {
			fmt.Println("Error writing to destination file:", err)
			return
		}
		curr--
	}
	fmt.Println("Successfully copied in reverse from", source, "to", destination)

}
