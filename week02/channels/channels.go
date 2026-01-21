//Danny Radosevich
//Channels for COSC3750

package main // package main defines the entry point for the program.
// Every executable Go program must have a main package.

import (
	"fmt" // import the fmt package for formatted I/O functions
)

func main() {
	// Create a new channel of type string
	messages := make(chan string)

	// Start a new goroutine that sends a message to the channel

	go func() {
		for i := 0; i < 5; i++ {
			messages <- fmt.Sprintf("ping %d", i)
			fmt.Println(fmt.Sprintf("ping %d", i))
		}
	}()

	// Start a new goroutine that receives messages from the channel

	//So why isn't this a goroutine?
	for i := 0; i < 5; i++ {
		//receive the message from the channel
		msg := <-messages
		fmt.Println(msg)
	}
	//effectively it could be, but then main would exit right away,
	// we leave it out of a routine to block

}
