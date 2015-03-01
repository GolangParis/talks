package main

import "fmt"

func main() {
	c := make(chan int, 2) // HLxxx
	c <- 1
	c <- 2
	fmt.Println(<-c)
	fmt.Println(<-c)
}
