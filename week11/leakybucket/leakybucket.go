//Danny Radosevich
//COSC 3750
//Leaky Bucket Example
//Adapted from hands on system programming with go
/*
Documentation generated with Claude
*/

// The leaky bucket algorithm is a classic rate-limiting technique.
// Imagine a bucket that holds a fixed number of tokens. Clients consume
// tokens to make requests. If there are not enough tokens the request is
// either denied or partially fulfilled. At a fixed interval the bucket is
// "refilled" back to capacity, allowing the next burst of requests.
//
// This implementation uses atomic operations instead of a mutex so that
// multiple goroutines can safely read and update the bucket's token count
// without blocking each other on a lock.

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

// bucket represents the leaky bucket itself.
// Both fields are uint64 so they can be manipulated with sync/atomic helpers,
// which require 64-bit aligned values.
type bucket struct {
	capacity uint64 // maximum number of tokens the bucket can hold
	status   uint64 // current number of tokens available (0 = empty / rate-limited)
}

// client represents one user or service that is competing for tokens.
type client struct {
	name  string        // human-readable label used in log output
	max   int           // upper bound on how many tokens to request per attempt
	b     *bucket       // shared bucket all clients draw from
	sleep time.Duration // how long to wait between attempts (controls request frequency)
}

// newBucket creates a bucket with the given capacity and starts a background
// goroutine that refills the bucket to full capacity every `rate` duration.
//
// Parameters:
//
//	ctx      — cancellation signal; when done the refill goroutine exits cleanly
//	cap      — maximum (and initial) token count
//	rate     — how often the bucket is refilled (e.g. time.Second/5 = every 200ms)
func newBucket(ctx context.Context, cap uint64, rate time.Duration) *bucket {
	// Start the bucket full so clients can make requests immediately.
	b := bucket{capacity: cap, status: cap}

	go func() {
		// time.NewTicker fires repeatedly at the given interval.
		ticker := time.NewTicker(rate)
		for {
			select {
			case <-ticker.C:
				// Refill: atomically reset status back to full capacity.
				// atomic.StoreUint64 is used instead of a plain assignment so
				// the write is safe even while other goroutines are reading status.
				atomic.StoreUint64(&b.status, b.capacity)

			case <-ctx.Done():
				// Context was cancelled (timeout or explicit cancel) — stop the ticker
				// to release its resources and let this goroutine exit.
				ticker.Stop()
				return
			}
		}
	}()
	return &b
}

// add attempts to take n tokens from the bucket and returns how many were
// actually granted. If the bucket is empty it returns 0. If fewer than n
// tokens remain the client gets only what is left (partial fulfillment).
//
// The function uses a compare-and-swap (CAS) loop to avoid a mutex:
//  1. Read the current token count.
//  2. Calculate the new count after taking tokens.
//  3. Atomically swap old → new, but only if status hasn't changed since step 1.
//  4. If another goroutine changed status between steps 1 and 3, retry.
func (b *bucket) add(n uint64) uint64 {
	for {
		// Atomically load the current available tokens.
		r := atomic.LoadUint64(&b.status)

		if r == 0 {
			// Bucket is empty — deny the request entirely.
			return 0
		}

		if n > r {
			// Requested more tokens than are available; cap to what's left.
			// The client is partially served rather than fully denied.
			n = r
		}

		// CompareAndSwapUint64(addr, old, new) atomically:
		//   if *addr == old { *addr = new; return true }
		//   else            { return false }
		// If it returns false, another goroutine changed status first — loop
		// and try again with the freshly loaded value.
		if !atomic.CompareAndSwapUint64(&b.status, r, r-n) {
			continue
		}

		// CAS succeeded — we successfully deducted n tokens.
		return n
	}
}

// run is the main loop for a client. It repeatedly:
//  1. Picks a random number of tokens to request (between 1 and c.max-1).
//  2. Sleeps for c.sleep to simulate think-time / network delay.
//  3. Calls add() to try to take tokens from the shared bucket.
//  4. Logs the attempt so we can observe rate-limiting in action.
//
// The loop exits when the context is cancelled (the program's 1-second timeout).
func (c client) run(ctx context.Context, start time.Time) {
	for {
		select {
		case <-ctx.Done():
			// Timeout or cancel — stop this client cleanly.
			return
		default:
			// Choose a random request size: [1, c.max)
			n := 1 + rand.Intn(c.max-1)

			// Pause to simulate the client doing work before its next request.
			time.Sleep(c.sleep)

			// Elapsed seconds since the program started — useful for reading logs.
			e := time.Since(start).Seconds()

			// Try to claim n tokens; a will be ≤ n (possibly 0 if the bucket is empty).
			a := c.b.add(uint64(n))

			log.Printf("%s tries to take %d after %.02fs, takes %d", c.name, n, e, a)
		}
	}
}

func main() {
	// Run the whole demo for exactly 1 second, then cancel everything.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // belt-and-suspenders: cancel even if timeout fires first

	start := time.Now()

	// t is the sleep interval for each client: 100ms (time.Second / 10).
	// With 5 clients each sleeping 100ms, requests arrive roughly every 20ms collectively.
	t := time.Second / 10

	// Create a shared bucket: capacity=10 tokens, refilled every 200ms (time.Second/5).
	// Refilling every 200ms means 5 refills per second → up to 50 tokens/sec available.
	b := newBucket(ctx, 10, time.Second/5)

	// Spawn 5 clients, all sharing the same bucket.
	// Each client requests between 1 and 4 tokens per attempt.
	for i := 0; i < 5; i++ {
		c := client{
			name:  fmt.Sprintf("Client%d", i),
			max:   5,
			b:     b,
			sleep: t,
		}
		go c.run(ctx, start)
	}

	// Block until the 1-second context timeout fires.
	// At that point all client goroutines and the refill goroutine will exit.
	<-ctx.Done()
}
