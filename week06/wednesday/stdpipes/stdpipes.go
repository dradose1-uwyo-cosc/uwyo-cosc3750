//Danny Radosevich
//COSC3750
//an example on pipes
//This version demonstrates a fan-out pipeline pattern:
//  cat file -> multiple grep processes (one per target word) -> wc -l per grep
//The key lessons here are pipe wiring, process start order, and wait/close order.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

func main() {
	var (
		// Words we want to count in the source file. Each word gets its own grep|wc pipeline.
		words = []string{"hello", "wyoming", "hoenn"}
		// cmds[i][0] is grep for words[i], cmds[i][1] is wc -l for words[i].
		cmds = make([][2]*exec.Cmd, len(words))
		// writers holds each grep's stdin pipe writer (owned by parent process).
		// Parent writes data into these by connecting cat's stdout to a MultiWriter.
		writers = make([]io.WriteCloser, len(words))
		// buffers collects stdout from each wc command so we can print counts later.
		buffers = make([]bytes.Buffer, len(words))
		err     error
	)

	// Build N independent grep|wc pipelines (one per word).
	for i, word := range words {
		// Create grep <word>
		cmds[i][0] = exec.Command("grep", word)

		// Create a pipe to feed grep's stdin.
		// The returned WriteCloser belongs to the parent process.
		if writers[i], err = cmds[i][0].StdinPipe(); err != nil {
			panic(err)
		}

		// Create wc -l to count matching lines produced by grep.
		cmds[i][1] = exec.Command("wc", "-l")

		// Wire grep stdout -> wc stdin.
		// StdoutPipe must be configured before Start().
		if cmds[i][1].Stdin, err = cmds[i][0].StdoutPipe(); err != nil {
			panic(err)
		}

		// Capture wc output into a bytes.Buffer (instead of printing immediately).
		cmds[i][1].Stdout = &buffers[i]
	}

	// io.MultiWriter requires []io.Writer, but writers is []io.WriteCloser.
	// Convert element-by-element (Go slices are not covariant).
	mwWriters := make([]io.Writer, len(writers))
	for i := range writers {
		mwWriters[i] = writers[i]
	}

	// cat will read the file once and broadcast bytes to every grep stdin writer.
	cat := exec.Command("cat", "further_filtered_rockyou.txt")
	cat.Stdout = io.MultiWriter(mwWriters...)

	// Start all downstream commands first so readers are ready before producer writes.
	// If producer starts too early, you can get broken pipes or blocked writes.
	for i := range cmds {
		if err := cmds[i][0].Start(); err != nil {
			panic(err)
		}
		if err := cmds[i][1].Start(); err != nil {
			panic(err)
		}
	}

	// Start producer after consumers are running.
	if err := cat.Start(); err != nil {
		panic(err)
	}

	// Wait for cat to finish writing file contents.
	if err := cat.Wait(); err != nil {
		panic(err)
	}

	// Close each grep stdin writer so grep receives EOF and can terminate.
	// Without this, grep might continue waiting for more input and appear to hang.
	for i := range cmds {
		if err := writers[i].Close(); err != nil {
			panic(err)
		}
	}

	// Wait for each grep process.
	// grep exits with code 1 when there are no matches, which is expected and not fatal.
	for i := range cmds {
		if err := cmds[i][0].Wait(); err != nil {
			var exitErr *exec.ExitError
			if !(errors.As(err, &exitErr) && exitErr.ExitCode() == 1) {
				panic(err)
			}
		}
	}

	// Wait for wc, then print a clean count per search word.
	for i := range cmds {
		if err := cmds[i][1].Wait(); err != nil {
			panic(err)
		}
		// Trim trailing newline from wc output (wc prints like "123\n").
		count := bytes.TrimSpace(buffers[i].Bytes())
		fmt.Printf("%s: %s\n", words[i], count)
	}

}
