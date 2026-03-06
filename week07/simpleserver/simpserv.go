//Danny Rdosevich
//COSC3750
//Simple http server

package main

import (
	"fmt"
	"net/http"
)

type customHandler int

func (c *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", *c)
	*c++
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})
	mux.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Goodbye!")
	})
	mux.Handle("/custom", new(customHandler))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Server error:", err)
	}
}
