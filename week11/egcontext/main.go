//Danny Radosevich
//COSC 3750
//Error Group Example
//Adapted from hands on system programming with go

package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"golang.org/x/sync/errgroup"
)

func sender(ctx context.Context, ch chan<- string, n int) func() error {
	return func() (err error) {
		for i := 0; ; i++ {
			if rand.Intn(100) == 42 {
				//I'll be honest, VSCode autocompleted the joke and made it better than I would have
				return errors.New("The Answer to the Ultimate Question of Life, The Universe, and Everything is 42. But this sender is tired of sending it.")
			}
			select {
			case ch <- fmt.Sprintf("Sender %d: %d", n, i):
			case <-ctx.Done():
				return nil // Stop sending if the context is cancelled
			}
		}
	}
}

func main() {
	eg, egctx := errgroup.WithContext(context.Background())
	ch := make(chan string)

	for i := 0; i < 10; i++ {
		eg.Go(sender(egctx, ch, i))
	}

	tries := 0
	go func() {
		for s := range ch {
			tries++
			fmt.Println(s)
		}
	}()
	if err := eg.Wait(); err != nil {
		fmt.Println("Error:", err)
		fmt.Printf("Senders *failed* after %d tries\n", tries)
	} else {
		fmt.Printf("All senders succeeded, took %d tries\n", tries)
	}
}
