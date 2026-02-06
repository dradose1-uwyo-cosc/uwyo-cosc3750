//COSC3750
// Modified from: https://go.dev/play/p/0YpRK25wFw_c

package main

import (
	"fmt"
	"io"
)

func main() {
	pr, pw := io.Pipe()
	//go routine so spinning off a writer
	go func(w io.WriteCloser) {
		for _, s := range []string{"a string", "another string", "last one"} {
			fmt.Printf("-> writing %q\n", s)
			fmt.Fprint(w, s)
		}
		w.Close()
	}(pw)
	var err error
	//the reader is not in a routine, so it will block main until done
	for n, b := 0, make([]byte, 100); err == nil; {
		fmt.Println("<- waiting...")
		n, err = pr.Read(b)
		if err == nil {
			fmt.Printf("<- received %q\n", string(b[:n]))
		}
	}
	if err != nil && err != io.EOF {
		fmt.Println("error:", err)
	}
}
