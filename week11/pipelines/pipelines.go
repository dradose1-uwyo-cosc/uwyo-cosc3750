//Danny Radosevich
//COSC3750
//Pipeline Example
//Adapted from the book

// This program demonstrates the pipeline concurrency pattern in Go.
// A pipeline is a series of stages connected by channels, where each stage
// is a goroutine that reads from an input channel, does some work, and writes
// to an output channel. Stages are chained together so data flows through them
// in sequence. Using context.Context lets any stage cancel the whole pipeline
// cleanly (e.g. on timeout or user interrupt).

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// stage is a generic pipeline stage that reads integers from `in`, increments
// each one by 1, and sends the result to the returned output channel.
//
// The goroutine runs until `in` is closed (all upstream values consumed) or
// the context is cancelled, whichever comes first. Closing `out` via defer
// signals the next stage downstream that there is no more data.
func stage(ctx context.Context, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out) // always close out when this goroutine exits
		for v := range in {
			v = v.(int) + 1 // type-assert to int and increment
			select {
			case out <- v:   // forward the processed value downstream
			case <-ctx.Done(): // stop early if the pipeline is cancelled
				return
			}
		}
	}()
	return out
}

// sourceLine reads lines one at a time from an io.ReadCloser (e.g. a file or
// an in-memory string reader) and emits each line as a string on the returned
// channel. It is the first stage (the "source") of the text pipeline.
//
// Both the reader and the output channel are closed when the goroutine exits,
// so downstream stages naturally stop when there is no more input.
func sourceLine(ctx context.Context, r io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		// Deferred cleanup: close the reader (releases any file handles) and
		// close the channel (signals downstream that the source is exhausted).
		defer func() { r.Close(); close(ch) }()

		scnr := bufio.NewScanner(r) // Scanner splits input on newlines by default
		for scnr.Scan() {
			line := scnr.Text() // grab the current line without the trailing newline
			select {
			case ch <- line:     // send the line downstream
			case <-ctx.Done(): // abort if the context was cancelled
				return
			}
		}
	}()
	return ch
}

// textFilter is a pipeline stage that sits between sourceLine and printer.
// It reads every line from `in` and only forwards lines that contain the
// given `filter` substring. Lines that do not match are silently dropped.
//
// This is the classic "grep-like" stage in a text-processing pipeline.
func textFilter(ctx context.Context, in <-chan string, filter string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch) // signal printer that filtering is done
		for line := range in {
			if strings.Contains(line, filter) { // only pass lines that match
				select {
				case ch <- line:     // forward the matching line
				case <-ctx.Done(): // stop if cancelled
					return
				}
			}
		}
	}()
	return ch
}

// printer is the final stage (the "sink") of the pipeline. It consumes every
// line from `in` and writes it to the io.Writer `w` (here, os.Stdout).
//
// If a line contains `highlight`, that substring is wrapped in ANSI escape
// codes so it appears in the given `color` in the terminal. For example,
// color=31 produces red text.
//
// ANSI escape code format:
//   \x1b[<code>m  — turn on an attribute (e.g. \x1b[31m = red foreground)
//   \x1b[39m      — reset foreground color to default
func printer(ctx context.Context, in <-chan string, color int, highlight string, w io.Writer) {
	const close = "\x1b[39m"                   // reset-color escape sequence
	open := fmt.Sprintf("\x1b[%dm", color)     // color-on escape sequence built from `color`

	for {
		select {
		case <-ctx.Done(): // pipeline was cancelled — stop printing
			return
		case line, ok := <-in:
			if !ok {
				// Channel was closed by textFilter — no more lines to print.
				return
			}

			// Look for the highlight word in the line.
			i := strings.Index(line, highlight)
			if i >= 0 {
				// Print: text before match | colored match | text after match | newline
				fmt.Fprint(w, line[:i], open, highlight, close, line[i+len(highlight):], "\n")
			} else {
				// Line passed the filter but highlight wasn't found — print as-is.
				fmt.Fprintln(w, line)
			}
		}
	}
}

func main() {
	var search string

	// context.Background() is the root context with no cancellation or deadline.
	// In a real app you might use context.WithCancel or context.WithTimeout here
	// so the whole pipeline can be shut down (e.g. on Ctrl-C).
	ctx := context.Background()

	fmt.Print("Enter a search term: ")
	fmt.Scanln(&search)

	// For demonstration, the "file" is a hard-coded multi-line string wrapped in
	// io.NopCloser so it satisfies the io.ReadCloser interface that sourceLine expects.
	lines := sourceLine(ctx, io.NopCloser(strings.NewReader(`This is a test.
This line contains the word test.
This line does not contain the word.
Another test line.`)))

	// Chain the stages together:
	//   sourceLine → textFilter → printer
	// Each stage runs in its own goroutine; channels connect them.
	filtered := textFilter(ctx, lines, search)
	printer(ctx, filtered, 31, search, os.Stdout) // 31 = red; blocks until pipeline is done
}
