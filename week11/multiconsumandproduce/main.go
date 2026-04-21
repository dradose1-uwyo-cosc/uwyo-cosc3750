//Danny Radosevich
//COSC 3750
//Multi-Producer, Multi-Consumer Example
//Adapted from hands on system programming with go

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	const numProducers = 3
	const numConsumers = 5

	wgConsumers := sync.WaitGroup{}
	wgConsumers.Add(numConsumers)
	wgProducers := sync.WaitGroup{}
	wgProducers.Add(numProducers)

	ch := make(chan string) // Buffered channel to hold produced items

	// Start producer goroutines
	for i := range numProducers {
		go func(id int) {
			for j := range 10 {
				ch <- fmt.Sprintf("Producer %d: item %d", id, j) // Produce an item and send it to the channel
			}
			wgProducers.Done() // Signal that this producer is done
		}(i)

	}
	// Start consumer goroutines
	for i := range numConsumers {
		go func(id int) {

			for item := range ch { // Consume items from the channel until it's closed
				fmt.Printf("Consumer %d: consumed %s\n", id, item)
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Simulate variable processing time
			}
			wgConsumers.Done() // Signal that this consumer is done
		}(i)
	}

	wgProducers.Wait() // Wait for all producers to finish
	close(ch)          // Close the channel to signal consumers that no more items will be produced
	wgConsumers.Wait() // Wait for all consumers to finish processing

}
