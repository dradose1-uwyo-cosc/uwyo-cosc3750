// Danny Radosevich
// COSC3750
// serve some files
package main

import (
	"net/http"
	"os"
)

func main() {
	dir := "."
	s, err := os.Stat(dir)
	if err != nil {
		panic(err)
	}
	if !s.IsDir() {
		panic("not a directory")
	}
	http.Handle("/", http.FileServer(http.Dir(dir)))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
