//Danny Radosevich
//COSC 3750
//Error Group Example
//Adapted from hands on system programming with go

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
main.go:14:2: no required module provides package golang.org/x/sync/errgroup; to add it:
        go get golang.org/x/sync/errgroup
*/

func visitor(url string) func() error {
	return func() (err error) {
		s := time.Now()
		defer func() {
			log.Println(url, time.Since(s), err)
		}()
		var resp *http.Response
		if resp, err = http.Get(url); err != nil {
			return
		}
		return resp.Body.Close()
	}
}

func main() {
	eg := errgroup.Group{}
	urls := []string{
		"https://www.google.com",
		"https://www.fake.invalid",
		"https://go.dev",
	}
	for _, url := range urls {
		eg.Go(visitor(url))
	}
	if err := eg.Wait(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("All visitors succeeded")
	}

}
