package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// main START OMIT
func main() {
	response := make(chan *http.Response, 1)
	errors := make(chan *error)

	go func() {
		resp, err := http.Get("http://matt.aimonetti.net/")
		if err != nil {
			errors <- &err
		}
		response <- resp
	}()
	for {
		select {
		case r := <-response:                      // HLxxx
			fmt.Printf("%#v", r.Body)
			return
		case err := <-errors:                      // HLxxx
			log.Fatal(err)
		case <-time.After(200 * time.Millisecond): // HLxxx
			fmt.Printf("Timed out!")
			return
		}
	}
}
// main END OMIT
