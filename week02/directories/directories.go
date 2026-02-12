//Danny Radosevich
//COSC3750
//Directories example

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	//get the current working directory
	dir, err := os.Getwd() //built into os
	//two returns, always best practice to check for errors
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Current working directory:", dir)
	}
	//change directory
	err = os.Chdir("..") //go up one directory
	if err != nil {
		log.Fatal(err)
	}

	if dir, err = os.Getwd(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("New working directory:", dir)
	}

}
