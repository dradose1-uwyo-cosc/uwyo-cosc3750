//Danny Radosevich
//COSC3750
//Basic WaitGroup Example

package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	for i := 0; rand.Intn(10) != 0; i++ {
		wg.Add(1) // Increment the WaitGroup counter for each goroutine
		go func(i int) {
			for j := 0; j < 10; j++ {
				fmt.Printf("Goroutine %d: iteration %d\n", i, j)
			}
			wg.Done() // Decrement the WaitGroup counter when the goroutine is done
		}(i)
	}
	wg.Wait() // Wait for all goroutines to finish
}
