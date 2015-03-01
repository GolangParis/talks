package main

import "fmt"

func main() {
	c := make(chan int, 2) // HLxxx
	c <- 1
	c <- 2
	c3 := func() { c <- 3 } // HLxxx
	go c3()                 // HLxxx
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
}
