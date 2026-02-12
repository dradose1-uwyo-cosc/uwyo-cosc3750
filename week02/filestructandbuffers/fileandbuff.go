//Danny Radosevich
//COSC3750
//File and Buffer example

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("ragtime.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//base byte array buffer
	buffer := make([]byte, 16) //16 bytes
	for n := 0; err == nil; {
		n, err = file.Read(buffer)
		if err == nil {
			fmt.Printf("Read %d bytes: %s\n", n, string(buffer[:n]))
		}

	}
	if err != nil && err != io.EOF {
		fmt.Println("Error reading file:", err)
	}

	//using bufio for buffered reading
	fmt.Println("\nUsing bufio for buffered reading:")
	file.Seek(0, 0) //reset file pointer to beginning
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Print(line)
	}
}
