package main

import (
	"fmt"
)

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x // HLxxx
		x, y = y, x+y
	}
	close(c) // HLxxx
}

func main() {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c { // HLxxx
		fmt.Println(i)
	}
}

