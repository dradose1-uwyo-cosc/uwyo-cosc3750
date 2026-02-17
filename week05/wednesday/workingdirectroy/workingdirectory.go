//Danny Radosevich
//COSC3750

//Example on working directory

package main

import (
	"fmt"
	"os"
)

func main() {
	//pull args, if there is one create a directory with that name and change to it
	args := os.Args
	fmt.Printf("Arguments: %v\n", args)
	if len(args) < 2 {
		fmt.Printf("No directory name provided, using current directory\n")
		return
	}
	//Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return
	}

	fmt.Printf("Current Working Directory: %s\n", cwd)
	dirName := args[1] //only one argument, the directory name
	err = os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
	fmt.Printf("Directory %s created successfully\n", dirName)
	err = os.Chdir(dirName)
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		return
	}
	fmt.Printf("Changed to directory: %s\n", dirName)
	//Get the current working directory again
	cwd, err = os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return
	}
	fmt.Printf("Current Working Directory: %s\n", cwd)
}
