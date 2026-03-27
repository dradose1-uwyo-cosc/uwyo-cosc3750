//Danny Radosevich
//COSC3750
//Basic Mutex Example

package main

import (
	"fmt"
	"sync"
)

// counter bundles a value with a mutex so concurrent access is always protected.
type counter struct {
	mew   sync.Mutex
	value int
}

// increment adds 1 to value. The mutex ensures only one goroutine modifies
// value at a time, preventing a data race.
func (c *counter) increment() {
	c.mew.Lock()
	c.value++
	c.mew.Unlock()
}

// decrement subtracts 1 from value, also protected by the mutex.
func (c *counter) decrement() {
	c.mew.Lock()
	c.value--
	c.mew.Unlock()
}

// getValue returns the current value. defer guarantees the mutex is unlocked
// even if a panic occurs before the return.
func (c *counter) getValue() int {
	c.mew.Lock()
	defer c.mew.Unlock()
	return c.value
}

func main() {
	c := counter{}

	// Buffered channel used as a signal: each goroutine sends one empty struct
	// when it finishes, allowing main to know when all work is done.
	done := make(chan struct{}, 10000)

	// Spawn 10,000 goroutines. Even indices increment, odd indices decrement.
	// Without the mutex, these concurrent writes to c.value would be a data race.
	for i := 0; i < cap(done); i++ {
		go func(i int) {
			if i%2 == 0 {
				c.increment()
			} else {
				c.decrement()
			}
			done <- struct{}{} // signal that this goroutine is done
		}(i)
	}

	//comment out the following for to see the
	// random results
	for i := 0; i < cap(done); i++ {
		<-done
	}

	fmt.Println(c.getValue())
}
